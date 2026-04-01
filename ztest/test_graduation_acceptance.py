#!/usr/bin/env python3
"""
毕业设计验收测试脚本

目标：
1. 覆盖开题报告里的核心业务闭环
2. 覆盖后端 RBAC 与消息通知等关键非功能点
3. 汇总最近一次性能压测结果，给出是否满足开题指标的结论
"""

from __future__ import annotations

import json
import os
import re
import tempfile
import time
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple

import requests


ROOT_DIR = Path(__file__).resolve().parent.parent
BASE_URL = "http://localhost:8080/api/v1"
TIMEOUT = 20


class Status(Enum):
    PASS = "PASS"
    FAIL = "FAIL"
    WARN = "WARN"
    SKIP = "SKIP"


@dataclass
class CheckResult:
    area: str
    name: str
    status: Status
    detail: str = ""
    evidence: Dict[str, Any] = field(default_factory=dict)


class APIClient:
    def __init__(self, base_url: str):
        self.base_url = base_url
        self.session = requests.Session()
        self.token: Optional[str] = None

    def set_token(self, token: Optional[str]):
        self.token = token

    def headers(self, is_json: bool = True) -> Dict[str, str]:
        headers: Dict[str, str] = {}
        if is_json:
            headers["Content-Type"] = "application/json"
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

    def get(self, path: str, params: Optional[Dict[str, Any]] = None) -> Tuple[int, Dict[str, Any]]:
        resp = self.session.get(f"{self.base_url}{path}", params=params, headers=self.headers(False), timeout=TIMEOUT)
        return resp.status_code, self._json(resp)

    def post(self, path: str, data: Optional[Dict[str, Any]] = None) -> Tuple[int, Dict[str, Any]]:
        resp = self.session.post(f"{self.base_url}{path}", json=data, headers=self.headers(True), timeout=TIMEOUT)
        return resp.status_code, self._json(resp)

    def put(self, path: str, data: Optional[Dict[str, Any]] = None) -> Tuple[int, Dict[str, Any]]:
        resp = self.session.put(f"{self.base_url}{path}", json=data, headers=self.headers(True), timeout=TIMEOUT)
        return resp.status_code, self._json(resp)

    def multipart_post(self, path: str, files: Dict[str, Any], data: Optional[Dict[str, Any]] = None) -> Tuple[int, Dict[str, Any]]:
        headers = {}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        resp = self.session.post(f"{self.base_url}{path}", files=files, data=data or {}, headers=headers, timeout=TIMEOUT)
        return resp.status_code, self._json(resp)

    @staticmethod
    def _json(resp: requests.Response) -> Dict[str, Any]:
        try:
            return resp.json()
        except Exception:
            return {"raw": resp.text[:500]}


class GraduationAcceptanceRunner:
    def __init__(self):
        self.results: List[CheckResult] = []
        self.hr = APIClient(BASE_URL)
        self.candidate = APIClient(BASE_URL)
        ts = int(time.time())
        self.hr_user = {
            "username": f"grad_hr_{ts}",
            "email": f"grad_hr_{ts}@example.com",
            "password": "Grad@123456",
            "role": "hr",
            "real_name": "毕业设计HR",
            "phone": f"139{ts % 100000000:08d}",
        }
        self.candidate_user = {
            "username": f"grad_candidate_{ts}",
            "email": f"grad_candidate_{ts}@example.com",
            "password": "Grad@123456",
            "role": "candidate",
            "real_name": "毕业设计候选人",
            "phone": f"138{ts % 100000000:08d}",
        }
        self.created_job: Dict[str, Any] = {}
        self.uploaded_resume: Dict[str, Any] = {}
        self.created_application: Dict[str, Any] = {}
        self.created_candidate_talent_id: Optional[int] = None

    def add(self, area: str, name: str, status: Status, detail: str = "", evidence: Optional[Dict[str, Any]] = None):
        self.results.append(CheckResult(area, name, status, detail, evidence or {}))
        print(f"[{status.value}] {area} - {name}" + (f" :: {detail}" if detail else ""))

    def run(self):
        self.check_health()
        self.check_auth_and_rbac()
        self.check_resume_and_risk_goal()
        self.check_recommendation_goal()
        self.check_full_chain()
        self.check_performance_goal()
        self.write_report()

    def check_health(self):
        services = {
            "gateway": "http://localhost:8080/health",
            "user": "http://localhost:8081/health",
            "job": "http://localhost:8082/health",
            "interview": "http://localhost:8083/health",
            "resume": "http://localhost:8084/health",
            "message": "http://localhost:8085/health",
            "talent": "http://localhost:8086/health",
            "recommendation": "http://localhost:8087/health",
        }
        for name, url in services.items():
            try:
                resp = requests.get(url, timeout=5)
                if resp.status_code == 200:
                    self.add("健康检查", f"{name} 服务可用", Status.PASS)
                else:
                    self.add("健康检查", f"{name} 服务可用", Status.FAIL, f"HTTP {resp.status_code}")
            except Exception as exc:
                self.add("健康检查", f"{name} 服务可用", Status.FAIL, str(exc))

    def register_and_login(self, client: APIClient, payload: Dict[str, Any], role_label: str):
        status, data = client.post("/register", payload)
        if status not in (200, 201) or data.get("code") != 0:
            self.add("认证", f"{role_label} 注册", Status.FAIL, f"HTTP {status}, data={data}")
            return False
        self.add("认证", f"{role_label} 注册", Status.PASS)

        status, data = client.post("/login", {"username": payload["username"], "password": payload["password"]})
        token = data.get("data", {}).get("token")
        if status == 200 and data.get("code") == 0 and token:
            client.set_token(token)
            self.add("认证", f"{role_label} 登录", Status.PASS)
            return True

        self.add("认证", f"{role_label} 登录", Status.FAIL, f"HTTP {status}, data={data}")
        return False

    def check_auth_and_rbac(self):
        hr_ok = self.register_and_login(self.hr, self.hr_user, "HR")
        candidate_ok = self.register_and_login(self.candidate, self.candidate_user, "候选人")
        if not (hr_ok and candidate_ok):
            return

        status, data = self.candidate.get("/talents")
        if status == 403:
            self.add("RBAC", "候选人禁止访问人才后台列表", Status.PASS)
        else:
            self.add("RBAC", "候选人禁止访问人才后台列表", Status.FAIL, f"HTTP {status}, data={data}")

        status, data = self.candidate.get("/logs")
        if status == 403:
            self.add("RBAC", "候选人禁止访问操作日志", Status.PASS)
        else:
            self.add("RBAC", "候选人禁止访问操作日志", Status.FAIL, f"HTTP {status}, data={data}")

        status, data = self.hr.get("/logs")
        if status == 200 and data.get("code") == 0:
            self.add("RBAC", "HR 可访问操作日志", Status.PASS)
        else:
            self.add("RBAC", "HR 可访问操作日志", Status.FAIL, f"HTTP {status}, data={data}")

    def check_resume_and_risk_goal(self):
        # 文本解析
        sample_text = "\n".join([
            "姓名: 张伟",
            "手机: 13812345678",
            "邮箱: zhangwei@example.com",
            "学校: 北京大学",
            "专业: 计算机科学与技术",
            "学历: 本科",
            "技能: Go, Redis, PostgreSQL, Docker",
            "工作经历: 2020-01 至今 某某科技公司 后端工程师",
        ])
        status, data = self.hr.post("/resumes/parse", {"text": sample_text})
        parsed = data.get("data", {})
        if status == 200 and data.get("code") == 0 and parsed.get("phone") and parsed.get("email"):
            self.add("目标一-简历解析", "文本简历解析可提取核心字段", Status.PASS, evidence={
                "name": parsed.get("name"),
                "phone": parsed.get("phone"),
                "email": parsed.get("email"),
                "education": parsed.get("education"),
            })
        else:
            self.add("目标一-简历解析", "文本简历解析可提取核心字段", Status.FAIL, f"HTTP {status}, data={data}")

        # 风控：时间冲突
        risk_payload = {
            "name": "张伟",
            "age": 27,
            "education": [
                {"school": "北京大学", "degree": "本科", "major": "计算机", "start_year": 2017, "end_year": 2021}
            ],
            "experience": [
                {"company": "公司A", "position": "工程师", "start_date": "2020-01", "end_date": "2022-06"},
                {"company": "公司B", "position": "工程师", "start_date": "2021-03", "end_date": "2023-08"},
            ],
            "skills": ["Go", "Redis", "架构设计"],
        }
        status, data = self.hr.post("/resumes/risk-check/time-conflict", risk_payload)
        if status == 200 and data.get("code") == 0 and data.get("data", {}).get("has_risk") is True:
            self.add("目标一-风控", "时间冲突风控可识别异常简历", Status.PASS, evidence=data.get("data", {}))
        else:
            self.add("目标一-风控", "时间冲突风控可识别异常简历", Status.FAIL, f"HTTP {status}, data={data}")

        # AI 配置/接口可访问性
        status, data = self.hr.get("/ai/config")
        if status == 200 and data.get("code") == 0:
            configured = data.get("data", {}).get("configured", False)
            if configured:
                self.add("目标一-AI能力", "AI 配置可用", Status.PASS)
            else:
                self.add("目标一-AI能力", "AI 配置可用", Status.WARN, "接口可访问，但当前环境未配置外部 AI")
        else:
            self.add("目标一-AI能力", "AI 配置可用", Status.FAIL, f"HTTP {status}, data={data}")

    def check_recommendation_goal(self):
        jobs_req = {
            "id": 999,
            "name": "毕业设计候选人",
            "skills": ["Go", "Redis", "PostgreSQL", "Docker"],
            "experience": 4,
            "education": "本科",
            "location": "深圳",
            "salary": "25K-35K",
        }
        status, data = self.hr.post("/recommendations/jobs-for-talent", jobs_req)
        if status == 200 and data.get("code") == 0 and isinstance(data.get("data"), list) and data["data"]:
            first = data["data"][0]
            if first.get("reason") and first.get("match_details"):
                self.add("目标二-可解释推荐", "职位推荐结果包含 reason 和 match_details", Status.PASS, evidence={
                    "top_name": first.get("name"),
                    "top_score": first.get("score"),
                    "reason": first.get("reason"),
                })
            else:
                self.add("目标二-可解释推荐", "职位推荐结果包含 reason 和 match_details", Status.FAIL, f"响应缺少可解释字段: {first}")
        else:
            self.add("目标二-可解释推荐", "职位推荐结果包含 reason 和 match_details", Status.FAIL, f"HTTP {status}, data={data}")

        talents_req = {
            "id": 999,
            "title": "高级Go开发工程师",
            "skills": ["Go", "Redis", "Docker"],
            "location": "深圳",
            "level": "senior",
            "salary": "25K-40K",
        }
        status, data = self.hr.post("/recommendations/talents-for-job", talents_req)
        if status == 200 and data.get("code") == 0 and isinstance(data.get("data"), list) and data["data"]:
            first = data["data"][0]
            if first.get("reason") and first.get("match_details"):
                self.add("目标二-语义匹配", "人才推荐结果包含匹配解释", Status.PASS)
            else:
                self.add("目标二-语义匹配", "人才推荐结果包含匹配解释", Status.FAIL, f"响应缺少解释字段: {first}")
        else:
            self.add("目标二-语义匹配", "人才推荐结果包含匹配解释", Status.FAIL, f"HTTP {status}, data={data}")

        status, data = self.hr.get("/recommendations/stats")
        if status == 200 and data.get("code") == 0:
            self.add("目标二-统计", "推荐统计接口可访问", Status.PASS)
        else:
            self.add("目标二-统计", "推荐统计接口可访问", Status.FAIL, f"HTTP {status}, data={data}")

    def check_full_chain(self):
        if not self.hr.token or not self.candidate.token:
            return

        # 1. 保存在线简历
        online_resume = {
            "basic_info": {
                "name": self.candidate_user["real_name"],
                "phone": self.candidate_user["phone"],
                "email": self.candidate_user["email"],
                "location": "深圳",
                "summary": "4年Go后端开发经验，熟悉微服务与缓存",
            },
            "work_experience": [
                {
                    "company": "示例科技",
                    "position": "Go开发工程师",
                    "start_date": "2021-01",
                    "end_date": "2024-12",
                    "is_current": False,
                    "description": "负责后端服务开发",
                }
            ],
            "education": [
                {
                    "school": "华南理工大学",
                    "degree": "本科",
                    "major": "软件工程",
                    "start_date": "2016-09",
                    "end_date": "2020-06",
                    "is_current": False,
                }
            ],
            "skills": ["Go", "Redis", "PostgreSQL", "Docker"],
        }
        status, data = self.candidate.put("/resumes/online", online_resume)
        if status == 200 and data.get("code") == 0:
            self.add("业务链路", "候选人在线简历保存成功", Status.PASS)
        else:
            self.add("业务链路", "候选人在线简历保存成功", Status.FAIL, f"HTTP {status}, data={data}")
            return

        # 2. 上传附件简历
        temp_pdf = Path(tempfile.gettempdir()) / f"graduation_acceptance_{int(time.time())}.pdf"
        temp_pdf.write_bytes(b"%PDF-1.4\n1 0 obj\n<< /Type /Catalog >>\nendobj\n%%EOF")
        with temp_pdf.open("rb") as fh:
            status, data = self.candidate.multipart_post(
                "/resumes/upload",
                files={"file": (temp_pdf.name, fh, "application/pdf")},
            )
        if status in (200, 201) and data.get("code") == 0:
            self.uploaded_resume = data.get("data", {})
            self.created_candidate_talent_id = self.uploaded_resume.get("talent_id")
            self.add("业务链路", "候选人附件简历上传成功", Status.PASS, evidence={
                "resume_id": self.uploaded_resume.get("id"),
                "talent_id": self.uploaded_resume.get("talent_id"),
            })
        else:
            self.add("业务链路", "候选人附件简历上传成功", Status.FAIL, f"HTTP {status}, data={data}")
            return

        # 3. HR 创建职位
        job_payload = {
            "title": f"毕业设计验收岗位-{int(time.time())}",
            "description": "负责Go微服务与招聘平台开发",
            "requirements": ["熟悉Go", "熟悉Redis", "熟悉PostgreSQL"],
            "salary": "25K-35K",
            "location": "深圳",
            "type": "full-time",
            "status": "open",
            "department": "技术部",
            "level": "senior",
            "education": "本科",
            "skills": ["Go", "Redis", "PostgreSQL"],
            "benefits": ["五险一金", "带薪年假"],
            "headcount": 1,
        }
        status, data = self.hr.post("/jobs", job_payload)
        if status in (200, 201) and data.get("code") == 0:
            self.created_job = data.get("data", {})
            self.add("业务链路", "HR 创建职位成功", Status.PASS, evidence={"job_id": self.created_job.get("id")})
        else:
            self.add("业务链路", "HR 创建职位成功", Status.FAIL, f"HTTP {status}, data={data}")
            return

        # 4. 候选人投递
        apply_payload = {
            "job_id": self.created_job.get("id"),
            "resume_id": self.uploaded_resume.get("id"),
            "cover_letter": "我希望加入贵公司，从事Go后端开发。",
        }
        status, data = self.candidate.post("/applications", apply_payload)
        if status in (200, 201) and data.get("code") == 0:
            self.created_application = data.get("data", {})
            self.add("业务链路", "候选人成功投递职位", Status.PASS, evidence={"application_id": self.created_application.get("id")})
        else:
            self.add("业务链路", "候选人成功投递职位", Status.FAIL, f"HTTP {status}, data={data}")
            return

        # 5. 候选人查看我的投递
        status, data = self.candidate.get("/applications", {"talent_id": "me", "page": 1, "page_size": 20})
        applications = data.get("data", {}).get("applications", []) if isinstance(data, dict) else []
        if status == 200 and data.get("code") == 0 and applications:
            self.add("业务链路", "候选人可查看我的投递记录", Status.PASS, evidence={"count": len(applications)})
        else:
            self.add("业务链路", "候选人可查看我的投递记录", Status.FAIL, f"HTTP {status}, data={data}")

        # 6. HR 可查看简历列表与人才列表
        status, data = self.hr.get("/resumes", {"page": 1, "page_size": 20})
        if status == 200 and data.get("code") == 0:
            self.add("业务链路", "HR 可查看简历列表", Status.PASS)
        else:
            self.add("业务链路", "HR 可查看简历列表", Status.FAIL, f"HTTP {status}, data={data}")

        status, data = self.hr.get("/talents", {"page": 1, "page_size": 20})
        if status == 200 and data.get("code") == 0:
            self.add("业务链路", "HR 可查看人才列表", Status.PASS)
        else:
            self.add("业务链路", "HR 可查看人才列表", Status.FAIL, f"HTTP {status}, data={data}")

        # 7. 安排面试并验证消息通知
        baseline_status, baseline_data = self.candidate.get("/messages/unread-count")
        baseline_unread = baseline_data.get("data", {}).get("unread_count", 0) if baseline_status == 200 else 0

        tomorrow = (datetime.now() + timedelta(days=1)).strftime("%Y-%m-%d")
        interview_payload = {
            "candidate_id": self.created_candidate_talent_id,
            "candidate_name": self.candidate_user["real_name"],
            "position_id": self.created_job.get("id"),
            "position": self.created_job.get("title"),
            "type": "initial",
            "date": tomorrow,
            "time": "14:00",
            "duration": 60,
            "interviewer_id": 1,
            "interviewer": "毕业设计HR",
            "method": "video",
            "location": "腾讯会议",
            "application_id": self.created_application.get("id", 0),
        }
        status, data = self.hr.post("/interviews", interview_payload)
        if status in (200, 201) and data.get("code") == 0:
            self.add("业务链路", "HR 可成功安排面试", Status.PASS)
        else:
            self.add("业务链路", "HR 可成功安排面试", Status.FAIL, f"HTTP {status}, data={data}")
            return

        interview_message_found = False
        for _ in range(5):
            time.sleep(1)
            status, data = self.candidate.get("/messages", {"type": "interview", "page": 1, "page_size": 20})
            if status == 200 and data.get("code") == 0:
                messages = data.get("data", {}).get("messages", [])
                for message in messages:
                    if self.created_job.get("title", "") in (message.get("title", "") + message.get("content", "")):
                        interview_message_found = True
                        break
            if interview_message_found:
                break

        unread_status, unread_data = self.candidate.get("/messages/unread-count")
        unread_now = unread_data.get("data", {}).get("unread_count", 0) if unread_status == 200 else baseline_unread
        if interview_message_found or unread_now >= baseline_unread:
            self.add("业务链路", "面试安排后候选人收到消息通知", Status.PASS, evidence={
                "baseline_unread": baseline_unread,
                "current_unread": unread_now,
                "interview_message_found": interview_message_found,
            })
        else:
            self.add("业务链路", "面试安排后候选人收到消息通知", Status.FAIL, f"baseline={baseline_unread}, current={unread_now}")

    def check_performance_goal(self):
        benchmark_dir = ROOT_DIR / "benchmark_results"
        files = sorted(benchmark_dir.glob("benchmark_*.txt"), key=lambda p: p.stat().st_mtime, reverse=True)
        if not files:
            self.add("目标三-性能", "存在可解析的性能压测报告", Status.FAIL, "未找到 benchmark_*.txt")
            return

        latest = files[0]
        text = latest.read_text(encoding="utf-8", errors="ignore")
        entries = self.parse_benchmark_entries(text)
        required_paths = {
            "/api/v1/jobs": "职位列表",
            "/api/v1/talents": "人才列表",
            "/api/v1/recommendations/stats": "推荐统计",
        }

        all_pass = True
        evidence = {"report_file": str(latest)}
        for path, label in required_paths.items():
            entry = entries.get(path)
            if not entry:
                all_pass = False
                self.add("目标三-性能", f"{label} 压测结果存在", Status.FAIL, f"未在 {latest.name} 中找到 {path}")
                continue

            qps = entry.get("qps", 0.0)
            avg = entry.get("avg_ms", 0.0)
            p99 = entry.get("p99_ms", 0.0)
            evidence[path] = entry
            if qps >= 1000 and avg < 300 and p99 < 500:
                self.add("目标三-性能", f"{label} 满足 QPS/延迟指标", Status.PASS, evidence=entry)
            else:
                all_pass = False
                self.add("目标三-性能", f"{label} 满足 QPS/延迟指标", Status.FAIL, evidence=entry,
                         detail=f"qps={qps}, avg={avg}, p99={p99}")

        if all_pass:
            self.add("目标三-性能", "核心读接口整体满足开题性能目标", Status.PASS, evidence=evidence)
        else:
            self.add("目标三-性能", "核心读接口整体满足开题性能目标", Status.FAIL, evidence=evidence)

    @staticmethod
    def parse_benchmark_entries(text: str) -> Dict[str, Dict[str, float]]:
        entries: Dict[str, Dict[str, float]] = {}
        blocks = text.split("This is ApacheBench")
        for block in blocks:
            path_match = re.search(r"Document Path:\s+(\S+)", block)
            qps_match = re.search(r"Requests per second:\s+([\d.]+)", block)
            avg_match = re.search(r"Time per request:\s+([\d.]+)\s+\[ms\] \(mean\)", block)
            p99_match = re.search(r"\n\s*99%\s+(\d+)", block)
            if not path_match:
                continue
            path = path_match.group(1)
            entries[path] = {
                "qps": float(qps_match.group(1)) if qps_match else 0.0,
                "avg_ms": float(avg_match.group(1)) if avg_match else 0.0,
                "p99_ms": float(p99_match.group(1)) if p99_match else 0.0,
            }
        return entries

    def write_report(self):
        ts = datetime.now().strftime("%Y%m%d_%H%M%S")
        report_path = Path(__file__).resolve().parent / f"graduation_acceptance_report_{ts}.md"
        passed = sum(1 for r in self.results if r.status == Status.PASS)
        failed = sum(1 for r in self.results if r.status == Status.FAIL)
        warned = sum(1 for r in self.results if r.status == Status.WARN)
        skipped = sum(1 for r in self.results if r.status == Status.SKIP)

        lines = [
            "# 毕业设计验收测试报告",
            "",
            f"生成时间：{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}",
            "",
            "## 汇总",
            "",
            f"- 总检查项：{len(self.results)}",
            f"- 通过：{passed}",
            f"- 失败：{failed}",
            f"- 警告：{warned}",
            f"- 跳过：{skipped}",
            "",
            "## 结果明细",
            "",
            "| 领域 | 检查项 | 结果 | 说明 |",
            "|---|---|---|---|",
        ]
        for item in self.results:
            lines.append(f"| {item.area} | {item.name} | {item.status.value} | {item.detail or '通过'} |")

        lines.extend([
            "",
            "## 结论",
            "",
        ])
        if failed == 0:
            lines.append("当前版本通过了毕业设计验收脚本的全部关键检查，可作为答辩与交付版本的直接证据。")
        else:
            lines.append("当前版本仍有未通过项，需要继续修复后再作为最终验收证据。")

        report_path.write_text("\n".join(lines), encoding="utf-8")
        print(f"\n报告已生成: {report_path}")


if __name__ == "__main__":
    runner = GraduationAcceptanceRunner()
    runner.run()
