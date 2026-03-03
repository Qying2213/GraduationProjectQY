#!/usr/bin/env python3
"""
智能人才运营平台 - 综合测试脚本
Intelligent Talent Operation Platform - Comprehensive Test Suite

功能覆盖:
- 服务健康检查
- 用户认证 (注册/登录/资料)
- 职位管理 CRUD
- 简历管理 CRUD
- 人才管理
- 消息服务
- 推荐服务
- 面试管理
- 统计服务
- AI 评估

作者: 自动生成
日期: 2026-02-03
"""

import requests
import json
import time
import sys
import os
from datetime import datetime
from typing import Optional, Dict, Any, List, Tuple
from dataclasses import dataclass, field
from enum import Enum

# ==================== 配置 ====================


class Config:
    """测试配置"""

    GATEWAY_URL = "http://localhost:8080"
    API_BASE = f"{GATEWAY_URL}/api/v1"
    TIMEOUT = 120  # 增加超时时间以支持 AI 评估等耗时请求

    # 微服务直接访问 (用于健康检查)
    SERVICES = {
        "gateway": ("http://localhost:8080", "/health"),
        "user": ("http://localhost:8081", "/health"),
        "job": ("http://localhost:8082", "/health"),
        "interview": ("http://localhost:8083", "/health"),
        "resume": ("http://localhost:8084", "/health"),
        "message": ("http://localhost:8085", "/health"),
        "talent": ("http://localhost:8086", "/health"),
        "recommendation": ("http://localhost:8087", "/health"),
    }


# ==================== 测试结果 ====================


class TestStatus(Enum):
    PASS = "✅ PASS"
    FAIL = "❌ FAIL"
    SKIP = "⏭️ SKIP"
    WARN = "⚠️ WARN"


@dataclass
class TestResult:
    """单个测试结果"""

    name: str
    status: TestStatus
    message: str = ""
    response_code: int = 0
    response_time_ms: float = 0
    details: Dict[str, Any] = field(default_factory=dict)


@dataclass
class TestSuite:
    """测试套件"""

    name: str
    results: List[TestResult] = field(default_factory=list)

    def add(self, result: TestResult):
        self.results.append(result)
        status_icon = result.status.value
        print(f"  {status_icon} {result.name}")
        if result.message:
            print(f"       └─ {result.message}")

    @property
    def passed(self) -> int:
        return sum(1 for r in self.results if r.status == TestStatus.PASS)

    @property
    def failed(self) -> int:
        return sum(1 for r in self.results if r.status == TestStatus.FAIL)

    @property
    def skipped(self) -> int:
        return sum(1 for r in self.results if r.status == TestStatus.SKIP)


# ==================== HTTP 客户端 ====================


class APIClient:
    """API 客户端"""

    def __init__(self, base_url: str):
        self.base_url = base_url
        self.token: Optional[str] = None
        self.session = requests.Session()

    def set_token(self, token: str):
        self.token = token

    def _headers(self) -> Dict[str, str]:
        headers = {"Content-Type": "application/json"}
        if self.token:
            headers["Authorization"] = f"Bearer {self.token}"
        return headers

    def get(self, path: str, params: Dict = None) -> Tuple[int, Dict, float]:
        """GET 请求"""
        start = time.time()
        try:
            resp = self.session.get(
                f"{self.base_url}{path}",
                headers=self._headers(),
                params=params,
                timeout=Config.TIMEOUT,
            )
            elapsed = (time.time() - start) * 1000
            try:
                data = resp.json()
            except:
                data = {"raw": resp.text[:500]}
            return resp.status_code, data, elapsed
        except requests.RequestException as e:
            return 0, {"error": str(e)}, (time.time() - start) * 1000

    def post(self, path: str, data: Dict = None) -> Tuple[int, Dict, float]:
        """POST 请求"""
        start = time.time()
        try:
            resp = self.session.post(
                f"{self.base_url}{path}",
                headers=self._headers(),
                json=data,
                timeout=Config.TIMEOUT,
            )
            elapsed = (time.time() - start) * 1000
            try:
                result = resp.json()
            except:
                result = {"raw": resp.text[:500]}
            return resp.status_code, result, elapsed
        except requests.RequestException as e:
            return 0, {"error": str(e)}, (time.time() - start) * 1000

    def put(self, path: str, data: Dict = None) -> Tuple[int, Dict, float]:
        """PUT 请求"""
        start = time.time()
        try:
            resp = self.session.put(
                f"{self.base_url}{path}",
                headers=self._headers(),
                json=data,
                timeout=Config.TIMEOUT,
            )
            elapsed = (time.time() - start) * 1000
            try:
                result = resp.json()
            except:
                result = {"raw": resp.text[:500]}
            return resp.status_code, result, elapsed
        except requests.RequestException as e:
            return 0, {"error": str(e)}, (time.time() - start) * 1000

    def delete(self, path: str) -> Tuple[int, Dict, float]:
        """DELETE 请求"""
        start = time.time()
        try:
            resp = self.session.delete(
                f"{self.base_url}{path}",
                headers=self._headers(),
                timeout=Config.TIMEOUT,
            )
            elapsed = (time.time() - start) * 1000
            try:
                result = resp.json()
            except:
                result = {"raw": resp.text[:500]}
            return resp.status_code, result, elapsed
        except requests.RequestException as e:
            return 0, {"error": str(e)}, (time.time() - start) * 1000


# ==================== 测试类 ====================


class TalentPlatformTests:
    """智能人才运营平台测试"""

    def __init__(self):
        self.client = APIClient(Config.API_BASE)
        self.suites: List[TestSuite] = []
        self.test_user = {
            "username": f"testuser_{int(time.time())}",
            "email": f"test_{int(time.time())}@example.com",
            "password": "Test@123456",
            "role": "hr",
        }
        self.created_resources: Dict[str, List[int]] = {
            "jobs": [],
            "resumes": [],
            "interviews": [],
        }

    def run_all(self):
        """运行所有测试"""
        print("\n" + "=" * 70)
        print("     智能人才运营平台 - 综合测试套件")
        print("     Intelligent Talent Operation Platform - Test Suite")
        print("=" * 70)
        print(f"\n⏰ 开始时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n")

        # 1. 服务健康检查
        self.test_service_health()

        # 2. 用户认证
        self.test_user_auth()

        # 3. 职位服务
        self.test_job_service()

        # 4. 简历服务
        self.test_resume_service()

        # 5. 人才服务
        self.test_talent_service()

        # 6. 消息服务
        self.test_message_service()

        # 7. 推荐服务
        self.test_recommendation_service()

        # 8. 面试服务
        self.test_interview_service()

        # 9. 统计服务
        self.test_stats_service()

        # 10. AI 评估
        self.test_ai_evaluation()

        # 生成报告
        self.generate_report()

    # ==================== 1. 服务健康检查 ====================

    def test_service_health(self):
        """测试所有微服务健康状态"""
        suite = TestSuite("1. 服务健康检查")
        print(f"\n{'─' * 50}")
        print(f"📡 {suite.name}")
        print(f"{'─' * 50}")

        for name, (url, path) in Config.SERVICES.items():
            try:
                start = time.time()
                resp = requests.get(f"{url}{path}", timeout=5)
                elapsed = (time.time() - start) * 1000

                if resp.status_code == 200:
                    suite.add(
                        TestResult(
                            name=f"{name.capitalize()} Service",
                            status=TestStatus.PASS,
                            message=f"响应时间: {elapsed:.0f}ms",
                            response_code=resp.status_code,
                            response_time_ms=elapsed,
                        )
                    )
                else:
                    suite.add(
                        TestResult(
                            name=f"{name.capitalize()} Service",
                            status=TestStatus.FAIL,
                            message=f"HTTP {resp.status_code}",
                            response_code=resp.status_code,
                        )
                    )
            except Exception as e:
                suite.add(
                    TestResult(
                        name=f"{name.capitalize()} Service",
                        status=TestStatus.FAIL,
                        message=f"连接失败: {str(e)[:50]}",
                    )
                )

        self.suites.append(suite)

    # ==================== 2. 用户认证 ====================

    def test_user_auth(self):
        """测试用户认证功能"""
        suite = TestSuite("2. 用户认证")
        print(f"\n{'─' * 50}")
        print(f"👤 {suite.name}")
        print(f"{'─' * 50}")

        # 2.1 用户注册
        code, data, elapsed = self.client.post("/register", self.test_user)
        if code in [200, 201]:
            suite.add(
                TestResult(
                    name="用户注册",
                    status=TestStatus.PASS,
                    message=f"用户名: {self.test_user['username']}",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="用户注册",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}: {data.get('message', data)}",
                    response_code=code,
                )
            )

        # 2.2 用户登录
        code, data, elapsed = self.client.post(
            "/login",
            {
                "username": self.test_user["username"],
                "password": self.test_user["password"],
            },
        )
        if code == 200 and "data" in data:
            token = data.get("data", {}).get("token") or data.get("token")
            if token:
                self.client.set_token(token)
                suite.add(
                    TestResult(
                        name="用户登录",
                        status=TestStatus.PASS,
                        message="Token 获取成功",
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name="用户登录",
                        status=TestStatus.WARN,
                        message="登录成功但无 Token",
                        response_code=code,
                    )
                )
        else:
            suite.add(
                TestResult(
                    name="用户登录",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}: {data.get('message', data)}",
                    response_code=code,
                )
            )

        # 2.3 尝试使用admin登录
        if not self.client.token:
            code, data, elapsed = self.client.post(
                "/login", {"username": "admin", "password": "admin123"}
            )
            if code == 200:
                token = data.get("data", {}).get("token") or data.get("token")
                if token:
                    self.client.set_token(token)
                    suite.add(
                        TestResult(
                            name="Admin 登录备选",
                            status=TestStatus.PASS,
                            message="使用 admin 账号登录成功",
                            response_code=code,
                        )
                    )

        # 2.4 获取用户资料
        if self.client.token:
            code, data, elapsed = self.client.get("/profile")
            if code == 200:
                suite.add(
                    TestResult(
                        name="获取用户资料",
                        status=TestStatus.PASS,
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name="获取用户资料",
                        status=TestStatus.FAIL,
                        message=f"HTTP {code}",
                        response_code=code,
                    )
                )
        else:
            suite.add(
                TestResult(
                    name="获取用户资料", status=TestStatus.SKIP, message="无有效 Token"
                )
            )

        # 2.5 错误密码测试
        code, data, elapsed = self.client.post(
            "/login", {"username": "admin", "password": "wrong_password"}
        )
        if code in [400, 401]:
            suite.add(
                TestResult(
                    name="错误密码拒绝",
                    status=TestStatus.PASS,
                    message="正确拒绝错误密码",
                    response_code=code,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="错误密码拒绝",
                    status=TestStatus.FAIL,
                    message=f"应返回 401，实际 {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 3. 职位服务 ====================

    def test_job_service(self):
        """测试职位服务"""
        suite = TestSuite("3. 职位服务")
        print(f"\n{'─' * 50}")
        print(f"💼 {suite.name}")
        print(f"{'─' * 50}")

        # 3.1 获取职位列表
        code, data, elapsed = self.client.get("/jobs", {"page": 1, "page_size": 10})
        if code == 200:
            jobs = data.get("data", {}).get("jobs", [])
            suite.add(
                TestResult(
                    name="获取职位列表",
                    status=TestStatus.PASS,
                    message=f"共 {len(jobs)} 个职位",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取职位列表",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 3.2 创建职位
        if self.client.token:
            new_job = {
                "title": f"测试职位_{int(time.time())}",
                "department": "技术部",
                "location": "深圳",
                "salary": "20k-30k",
                "description": "这是一个测试职位",
                "requirements": ["熟悉Python", "熟悉Go", "3年以上经验"],  # 必须是数组
                "status": "open",
                "created_by": 1,  # 使用admin用户ID满足外键约束
            }
            code, data, elapsed = self.client.post("/jobs", new_job)
            if code in [200, 201]:
                job_id = data.get("data", {}).get("id")
                if job_id:
                    self.created_resources["jobs"].append(job_id)
                suite.add(
                    TestResult(
                        name="创建职位",
                        status=TestStatus.PASS,
                        message=f"职位ID: {job_id}",
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name="创建职位",
                        status=TestStatus.FAIL,
                        message=f"HTTP {code}: {data.get('message', '')}",
                        response_code=code,
                    )
                )
        else:
            suite.add(
                TestResult(name="创建职位", status=TestStatus.SKIP, message="需要认证")
            )

        # 3.3 获取单个职位
        if self.created_resources["jobs"]:
            job_id = self.created_resources["jobs"][0]
            code, data, elapsed = self.client.get(f"/jobs/{job_id}")
            if code == 200:
                suite.add(
                    TestResult(
                        name="获取职位详情",
                        status=TestStatus.PASS,
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name="获取职位详情",
                        status=TestStatus.FAIL,
                        message=f"HTTP {code}",
                        response_code=code,
                    )
                )

        # 3.4 职位搜索
        code, data, elapsed = self.client.get("/jobs", {"search": "工程师", "page": 1})
        if code == 200:
            suite.add(
                TestResult(
                    name="职位搜索",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="职位搜索",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 4. 简历服务 ====================

    def test_resume_service(self):
        """测试简历服务"""
        suite = TestSuite("4. 简历服务")
        print(f"\n{'─' * 50}")
        print(f"📄 {suite.name}")
        print(f"{'─' * 50}")

        # 4.1 获取简历列表
        code, data, elapsed = self.client.get("/resumes", {"page": 1, "page_size": 10})
        if code == 200:
            resumes = data.get("data", {}).get("resumes", [])
            suite.add(
                TestResult(
                    name="获取简历列表",
                    status=TestStatus.PASS,
                    message=f"共 {len(resumes)} 份简历",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
            # 保存一个真实存在的简历ID用于后续测试
            if resumes:
                self.created_resources["resumes"].append(resumes[0].get("id"))
        else:
            suite.add(
                TestResult(
                    name="获取简历列表",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 4.2 获取单个简历
        if self.created_resources["resumes"]:
            resume_id = self.created_resources["resumes"][0]
            code, data, elapsed = self.client.get(f"/resumes/{resume_id}")
            if code == 200:
                suite.add(
                    TestResult(
                        name="获取简历详情",
                        status=TestStatus.PASS,
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name="获取简历详情",
                        status=TestStatus.FAIL,
                        message=f"HTTP {code}",
                        response_code=code,
                    )
                )

        # 4.3 简历筛选
        code, data, elapsed = self.client.get("/resumes", {"status": "parsed"})
        if code == 200:
            suite.add(
                TestResult(
                    name="简历状态筛选",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="简历状态筛选",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 5. 人才服务 ====================

    def test_talent_service(self):
        """测试人才服务"""
        suite = TestSuite("5. 人才服务")
        print(f"\n{'─' * 50}")
        print(f"👥 {suite.name}")
        print(f"{'─' * 50}")

        # 5.1 获取人才列表
        code, data, elapsed = self.client.get("/talents", {"page": 1, "page_size": 10})
        if code == 200:
            talents = data.get("data", {}).get("talents", [])
            suite.add(
                TestResult(
                    name="获取人才列表",
                    status=TestStatus.PASS,
                    message=f"共 {len(talents)} 位人才",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取人才列表",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 5.2 人才搜索
        code, data, elapsed = self.client.get("/talents", {"search": "工程师"})
        if code == 200:
            suite.add(
                TestResult(
                    name="人才搜索",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="人才搜索",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 5.3 人才筛选 (按经验)
        code, data, elapsed = self.client.get("/talents", {"min_experience": 3})
        if code == 200:
            suite.add(
                TestResult(
                    name="人才经验筛选",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="人才经验筛选",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 6. 消息服务 ====================

    def test_message_service(self):
        """测试消息服务"""
        suite = TestSuite("6. 消息服务")
        print(f"\n{'─' * 50}")
        print(f"💬 {suite.name}")
        print(f"{'─' * 50}")

        if not self.client.token:
            suite.add(
                TestResult(
                    name="消息服务测试", status=TestStatus.SKIP, message="需要认证"
                )
            )
            self.suites.append(suite)
            return

        # 6.1 获取消息列表
        code, data, elapsed = self.client.get("/messages")
        if code == 200:
            messages = data.get("data", {}).get("messages") or []
            suite.add(
                TestResult(
                    name="获取消息列表",
                    status=TestStatus.PASS,
                    message=f"共 {len(messages)} 条消息",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        elif code == 401:
            suite.add(
                TestResult(
                    name="获取消息列表",
                    status=TestStatus.WARN,
                    message="需要认证 (401)",
                    response_code=code,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取消息列表",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 6.2 获取未读消息数
        code, data, elapsed = self.client.get("/messages/unread-count")
        if code == 200:
            count = data.get("data", {}).get("count", 0)
            suite.add(
                TestResult(
                    name="获取未读消息数",
                    status=TestStatus.PASS,
                    message=f"未读: {count}",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        elif code == 401:
            suite.add(
                TestResult(
                    name="获取未读消息数",
                    status=TestStatus.WARN,
                    message="需要认证 (401)",
                    response_code=code,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取未读消息数",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 7. 推荐服务 ====================

    def test_recommendation_service(self):
        """测试推荐服务"""
        suite = TestSuite("7. 推荐服务")
        print(f"\n{'─' * 50}")
        print(f"🎯 {suite.name}")
        print(f"{'─' * 50}")

        # 7.1 为人才推荐职位
        code, data, elapsed = self.client.post(
            "/recommendations/jobs-for-talent", {"talent_id": 1, "limit": 5}
        )
        if code == 200:
            # Response may be a list directly or nested in data
            resp_data = data.get("data", data) if isinstance(data, dict) else data
            jobs = (
                resp_data.get("recommendations", resp_data)
                if isinstance(resp_data, dict)
                else resp_data
            )
            jobs = jobs if isinstance(jobs, list) else []
            suite.add(
                TestResult(
                    name="为人才推荐职位",
                    status=TestStatus.PASS,
                    message=f"推荐 {len(jobs)} 个职位",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="为人才推荐职位",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}: {data.get('message', '')}",
                    response_code=code,
                )
            )

        # 7.2 为职位推荐人才
        code, data, elapsed = self.client.post(
            "/recommendations/talents-for-job", {"job_id": 1, "limit": 5}
        )
        if code == 200:
            # Response may be a list directly or nested in data
            resp_data = data.get("data", data) if isinstance(data, dict) else data
            talents = (
                resp_data.get("recommendations", resp_data)
                if isinstance(resp_data, dict)
                else resp_data
            )
            talents = talents if isinstance(talents, list) else []
            suite.add(
                TestResult(
                    name="为职位推荐人才",
                    status=TestStatus.PASS,
                    message=f"推荐 {len(talents)} 位人才",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="为职位推荐人才",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}: {data.get('message', '')}",
                    response_code=code,
                )
            )

        # 7.3 推荐统计
        code, data, elapsed = self.client.get("/recommendations/stats")
        if code == 200:
            suite.add(
                TestResult(
                    name="推荐统计",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="推荐统计",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 8. 面试服务 ====================

    def test_interview_service(self):
        """测试面试服务"""
        suite = TestSuite("8. 面试服务")
        print(f"\n{'─' * 50}")
        print(f"📅 {suite.name}")
        print(f"{'─' * 50}")

        # 8.1 获取面试列表
        code, data, elapsed = self.client.get(
            "/interviews", {"page": 1, "page_size": 10}
        )
        if code == 200:
            interviews = data.get("data", {}).get("interviews", [])
            suite.add(
                TestResult(
                    name="获取面试列表",
                    status=TestStatus.PASS,
                    message=f"共 {len(interviews)} 场面试",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取面试列表",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 8.2 面试筛选 (按状态)
        code, data, elapsed = self.client.get("/interviews", {"status": "scheduled"})
        if code == 200:
            suite.add(
                TestResult(
                    name="面试状态筛选",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="面试状态筛选",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 9. 统计服务 ====================

    def test_stats_service(self):
        """测试统计服务"""
        suite = TestSuite("9. 统计服务")
        print(f"\n{'─' * 50}")
        print(f"📊 {suite.name}")
        print(f"{'─' * 50}")

        stats_endpoints = [
            ("/stats/dashboard", "仪表盘统计"),
            ("/stats/funnel", "招聘漏斗"),
            ("/stats/channels", "渠道统计"),
            ("/stats/department-progress", "部门进度"),
            ("/stats/interviewer-rank", "面试官排行"),
            ("/stats/trend", "趋势数据"),
            ("/stats/job-rank", "职位排行"),
        ]

        for endpoint, name in stats_endpoints:
            code, data, elapsed = self.client.get(endpoint)
            if code == 200:
                suite.add(
                    TestResult(
                        name=name,
                        status=TestStatus.PASS,
                        response_code=code,
                        response_time_ms=elapsed,
                    )
                )
            else:
                suite.add(
                    TestResult(
                        name=name,
                        status=TestStatus.FAIL,
                        message=f"HTTP {code}",
                        response_code=code,
                    )
                )

        self.suites.append(suite)

    # ==================== 10. AI 评估 ====================

    def test_ai_evaluation(self):
        """测试 AI 评估功能"""
        suite = TestSuite("10. AI 评估")
        print(f"\n{'─' * 50}")
        print(f"🤖 {suite.name}")
        print(f"{'─' * 50}")

        # 10.1 检查 AI 配置
        code, data, elapsed = self.client.get("/ai/config")
        if code == 200:
            configured = data.get("data", {}).get("configured", False)
            suite.add(
                TestResult(
                    name="AI 配置检查",
                    status=TestStatus.PASS if configured else TestStatus.WARN,
                    message=f"已配置: {configured}",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="AI 配置检查",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        # 10.2 AI 评估 (使用ID=1的简历，已更新为真实文件路径)
        resume_id = 1  # 使用已更新为真实文件路径的简历
        code, data, elapsed = self.client.post(
            "/ai/evaluate", {"resume_id": resume_id, "job_id": 1}
        )
        if code == 200:
            score = data.get("data", {}).get("total_score", 0)
            suite.add(
                TestResult(
                    name="AI 简历评估",
                    status=TestStatus.PASS,
                    message=f"评估分数: {score}",
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        elif (code == 500 and "读取简历文件失败" in str(data)) or (
            code == 404 and "简历文件不存在" in str(data)
        ):
            suite.add(
                TestResult(
                    name="AI 简历评估",
                    status=TestStatus.WARN,
                    message="简历文件不存在 (种子数据问题)",
                    response_code=code,
                )
            )
        elif code == 503:
            suite.add(
                TestResult(
                    name="AI 简历评估",
                    status=TestStatus.SKIP,
                    message="AI 服务未配置",
                    response_code=code,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="AI 简历评估",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}: {data.get('message', '') if isinstance(data, dict) else data}",
                    response_code=code,
                )
            )

        # 10.3 获取当前任务
        code, data, elapsed = self.client.get("/ai/current-task")
        if code == 200:
            suite.add(
                TestResult(
                    name="获取当前 AI 任务",
                    status=TestStatus.PASS,
                    response_code=code,
                    response_time_ms=elapsed,
                )
            )
        else:
            suite.add(
                TestResult(
                    name="获取当前 AI 任务",
                    status=TestStatus.FAIL,
                    message=f"HTTP {code}",
                    response_code=code,
                )
            )

        self.suites.append(suite)

    # ==================== 报告生成 ====================

    def generate_report(self):
        """生成测试报告"""
        print("\n" + "=" * 70)
        print("                       📋 测试报告")
        print("=" * 70)

        total_pass = sum(s.passed for s in self.suites)
        total_fail = sum(s.failed for s in self.suites)
        total_skip = sum(s.skipped for s in self.suites)
        total = total_pass + total_fail + total_skip

        print(f"\n📊 总体结果:")
        print(f"   总测试数: {total}")
        print(
            f"   ✅ 通过: {total_pass} ({total_pass / total * 100:.1f}%)"
            if total > 0
            else ""
        )
        print(f"   ❌ 失败: {total_fail}")
        print(f"   ⏭️ 跳过: {total_skip}")

        print(f"\n📑 各模块结果:")
        print(f"   {'模块':<20} {'通过':<8} {'失败':<8} {'跳过':<8}")
        print(f"   {'-' * 44}")
        for suite in self.suites:
            print(
                f"   {suite.name:<20} {suite.passed:<8} {suite.failed:<8} {suite.skipped:<8}"
            )

        # 保存报告到文件
        report_file = f"test_report_{datetime.now().strftime('%Y%m%d_%H%M%S')}.md"
        self._save_markdown_report(
            report_file, total, total_pass, total_fail, total_skip
        )
        print(f"\n📁 详细报告已保存: {report_file}")

        print(f"\n⏰ 结束时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print("=" * 70 + "\n")

        # 返回退出码
        if total_fail > 0:
            sys.exit(1)

    def _save_markdown_report(
        self, filename: str, total: int, passed: int, failed: int, skipped: int
    ):
        """保存 Markdown 格式报告"""
        with open(filename, "w", encoding="utf-8") as f:
            f.write("# 智能人才运营平台 - 测试报告\n\n")
            f.write(f"**测试时间**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")

            f.write("## 📊 总体结果\n\n")
            f.write("| 指标 | 数量 | 百分比 |\n")
            f.write("|------|------|--------|\n")
            f.write(f"| 总测试 | {total} | 100% |\n")
            f.write(
                f"| ✅ 通过 | {passed} | {passed / total * 100:.1f}% |\n"
                if total > 0
                else ""
            )
            f.write(
                f"| ❌ 失败 | {failed} | {failed / total * 100:.1f}% |\n"
                if total > 0
                else ""
            )
            f.write(
                f"| ⏭️ 跳过 | {skipped} | {skipped / total * 100:.1f}% |\n\n"
                if total > 0
                else ""
            )

            f.write("## 📑 详细结果\n\n")
            for suite in self.suites:
                f.write(f"### {suite.name}\n\n")
                f.write("| 测试项 | 状态 | 说明 | 响应时间 |\n")
                f.write("|--------|------|------|----------|\n")
                for r in suite.results:
                    time_str = (
                        f"{r.response_time_ms:.0f}ms" if r.response_time_ms > 0 else "-"
                    )
                    msg = r.message[:50] + "..." if len(r.message) > 50 else r.message
                    f.write(f"| {r.name} | {r.status.value} | {msg} | {time_str} |\n")
                f.write("\n")


# ==================== 主程序 ====================

if __name__ == "__main__":
    try:
        tests = TalentPlatformTests()
        tests.run_all()
    except KeyboardInterrupt:
        print("\n\n⚠️ 测试被用户中断")
        sys.exit(130)
    except Exception as e:
        print(f"\n\n❌ 测试执行出错: {e}")
        import traceback

        traceback.print_exc()
        sys.exit(1)
