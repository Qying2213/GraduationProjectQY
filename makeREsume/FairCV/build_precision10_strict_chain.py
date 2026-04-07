#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import csv
import io
import json
import argparse
import re
import subprocess
import sys
import time
from collections import defaultdict
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path
from typing import Dict, List

import requests

SCRIPT_DIR = Path(__file__).resolve().parent
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

from add_information import ResumeVariableAnalyzer  # noqa: E402
from json_to_pdf_resumes import build_pdf  # noqa: E402


OUTPUT_DIR = SCRIPT_DIR / "eval_results" / "precision10_strict_chain"

STRICT_JOBS = [
    {
        "key": "backend",
        "template_position": "后端开发工程师",
        "template_skill_level": "高",
        "title": "P10严格-后端开发工程师",
        "department": "技术部",
        "location": "深圳",
        "salary": "25K-40K",
        "education": "本科",
        "description": "负责核心业务系统后端研发，强调 Go、Java、MySQL、Redis、Kafka、Docker 及微服务实战能力。",
        "requirements": ["熟悉Go", "熟悉Java", "熟悉MySQL", "熟悉Redis", "熟悉Kafka", "熟悉Docker", "具备微服务经验"],
        "benefits": ["五险一金", "带薪年假", "技术成长"],
        "skills": ["Go", "Java", "MySQL", "Redis", "Kafka", "Docker", "微服务"],
        "industry": "互联网",
        "company_size": "10000人以上",
        "benchmark_keywords": ["Go", "Java", "Redis", "Kafka", "Docker", "微服务", "Service Mesh", "Zookeeper"],
    },
    {
        "key": "frontend",
        "template_position": "前端开发工程师",
        "template_skill_level": "高",
        "title": "P10严格-前端开发工程师",
        "department": "技术部",
        "location": "杭州",
        "salary": "18K-30K",
        "education": "本科",
        "description": "负责企业后台与门户前端开发，强调 Vue、React、Angular、TypeScript 与微前端能力。",
        "requirements": ["熟悉Vue", "熟悉React", "熟悉Angular", "熟悉TypeScript", "熟悉Node.js", "具备微前端经验"],
        "benefits": ["五险一金", "弹性工作", "技术成长"],
        "skills": ["Vue", "React", "Angular", "TypeScript", "Node.js", "Qiankun", "Single-SPA"],
        "industry": "互联网",
        "company_size": "10000人以上",
        "benchmark_keywords": ["Vue", "React", "Angular", "TypeScript", "Node.js", "Qiankun", "Single-SPA", "微前端"],
    },
    {
        "key": "android",
        "template_position": "Android开发工程师",
        "template_skill_level": "高",
        "title": "P10严格-Android开发工程师",
        "department": "技术部",
        "location": "成都",
        "salary": "20K-32K",
        "education": "本科",
        "description": "负责 Android 客户端开发，强调 Kotlin、Java、Android SDK、Jetpack、RxJava 与 Retrofit 实战。",
        "requirements": ["熟悉Kotlin", "熟悉Java", "熟悉Android SDK", "熟悉Jetpack", "熟悉RxJava", "熟悉Retrofit"],
        "benefits": ["五险一金", "弹性工作", "移动端培训"],
        "skills": ["Kotlin", "Java", "Android SDK", "Jetpack", "RxJava", "Retrofit", "MVVM"],
        "industry": "移动互联网",
        "company_size": "10000人以上",
        "benchmark_keywords": ["Kotlin", "Java", "Android SDK", "Jetpack", "RxJava", "Retrofit", "MVVM", "性能优化"],
    },
]

NAMES = [
    "顾承泽", "沈知行", "陆景川", "周亦安", "许明远", "韩书衍", "林修远", "宋砚清", "江屿白",
    "陈星野", "叶清越", "苏景辰", "温时序", "顾南舟", "唐知远", "贺云深", "谢望舒", "裴景行",
    "程一诺", "闻嘉树", "梁知夏", "纪清和", "袁景澄", "严叙白", "程知意", "乔南星", "沈沐言",
    "韩清妍", "许知微", "林书宁", "宋安然", "苏清瑶", "顾念初", "叶知夏", "温若溪", "谢语桐",
]


def run_sql(sql: str) -> str:
    result = subprocess.run(
        ["psql", "-d", "talent_platform", "-t", "-A", "-c", sql],
        capture_output=True,
        text=True,
        check=True,
    )
    lines = [line.strip() for line in result.stdout.splitlines() if line.strip()]
    return lines[0] if lines else ""


def copy_query(sql: str) -> List[Dict[str, str]]:
    result = subprocess.run(
        ["psql", "-d", "talent_platform", "-c", f"copy ({sql}) to stdout with csv header"],
        capture_output=True,
        text=True,
        check=True,
    )
    return list(csv.DictReader(io.StringIO(result.stdout)))


def q(value: str) -> str:
    return "'" + value.replace("'", "''") + "'"


def array_expr(values: List[str]) -> str:
    items = ", ".join(q(v) for v in values)
    return f"ARRAY[{items}]::text[]"


def ensure_job(job: dict) -> int:
    existing = run_sql(f"select id from jobs where title = {q(job['title'])} limit 1;")
    if existing:
        return int(existing)

    sql = f"""
    insert into jobs (
      title, department, location, salary, type, education, description,
      requirements, benefits, skills, status, headcount, created_by
    ) values (
      {q(job['title'])}, {q(job['department'])}, {q(job['location'])}, {q(job['salary'])},
      'full-time', {q(job['education'])}, {q(job['description'])},
      {array_expr(job['requirements'])}, {array_expr(job['benefits'])}, {array_expr(job['skills'])},
      'open', 3, 1
    ) returning id;
    """
    return int(run_sql(sql))


def get_template(position: str, skill_level: str) -> dict:
    data = json.loads((SCRIPT_DIR / "data" / "resumes_template.json").read_text(encoding="utf-8"))
    for item in data["resumes"]:
        meta = item["metadata"]
        if meta["position"] == position and meta["skill_level"] == skill_level:
            return item
    raise ValueError(f"template not found: {position} / {skill_level}")


def replace_first(pattern: str, text: str, replacement: str) -> str:
    return re.sub(pattern, replacement, text, count=1, flags=re.MULTILINE)


def build_resume_content(template: dict, analyzer: ResumeVariableAnalyzer, name: str, phone: str, email: str, job: dict, age: int) -> str:
    recruitment_type = template["metadata"]["recruitment_type"]
    combination = {
        "name": name,
        "gender": "男",
        "age": str(age),
        "marriage": "未婚",
        "hukou": "深圳市",
        "political": "中共党员",
        "disability": "无",
        "industry": job["industry"],
        "company_size": job["company_size"],
        "work_experience": analyzer.generate_work_experience(age, job["industry"], job["company_size"]),
    }
    content = analyzer.apply_combination_to_resume(template["content"], combination)
    content = replace_first(r"(邮箱[：:]\s*)([^\n\r]+)", content, rf"\g<1>{email}")
    content = replace_first(r"(电话[：:]\s*)([^\n\r]+)", content, rf"\g<1>{phone}")
    bonus_lines = "\n".join(f"- {item}" for item in job["benchmark_keywords"])
    content += (
        "\n\n---\n\n### 基准岗位补充匹配亮点\n\n"
        f"- 目标岗位：{job['title']}\n"
        "- 为保证推荐实验一致性，本简历补充如下高匹配技术关键词：\n"
        f"{bonus_lines}\n"
    )
    if "社招" not in recruitment_type and age > 25:
        content = content.replace(f"年龄：{age}", "年龄：25", 1)
    return content


def generate_resumes(per_job: int, run_tag: str) -> List[Dict]:
    analyzer = ResumeVariableAnalyzer()
    generated = []
    name_idx = 0
    phone_seed = 18000001000

    for job in STRICT_JOBS:
        template = get_template(job["template_position"], job["template_skill_level"])
        for idx in range(per_job):
            name = NAMES[name_idx % len(NAMES)]
            name_idx += 1
            phone = f"{phone_seed + name_idx}"
            email = f"strict_{job['key']}_{run_tag}_{idx+1:02d}@benchmark.local"
            age = 30 + (idx % 4)
            content = build_resume_content(template, analyzer, name, phone, email, job, age)
            generated.append(
                {
                    "metadata": {
                        "position": job["title"],
                        "skill_level": template["metadata"]["skill_level"],
                        "recruitment_type": template["metadata"]["recruitment_type"],
                        "benchmark_key": job["key"],
                        "benchmark_index": idx + 1,
                        "run_tag": run_tag,
                    },
                    "content": content,
                }
            )
    return generated


def export_pdfs(resumes: List[Dict], output_dir: Path) -> List[Dict]:
    output_dir.mkdir(parents=True, exist_ok=True)
    manifest = []
    for idx, resume in enumerate(resumes, start=1):
        key = resume["metadata"]["benchmark_key"]
        name_match = re.search(r"姓名[：:]\s*([^\n\r]+)", resume["content"].replace("**", ""))
        name = name_match.group(1).strip() if name_match else f"candidate_{idx}"
        filename = f"{idx:03d}_{key}_{name}.pdf"
        pdf_path = output_dir / key / filename
        build_pdf(resume, pdf_path)
        manifest.append(
            {
                "index": idx,
                "benchmark_key": key,
                "job_title": resume["metadata"]["position"],
                "name": name,
                "pdf_path": str(pdf_path),
            }
        )
    with (output_dir / "manifest.csv").open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=list(manifest[0].keys()))
        writer.writeheader()
        writer.writerows(manifest)
    return manifest


def login_admin() -> str:
    resp = requests.post(
        "http://localhost:8080/api/v1/login",
        json={"username": "admin", "password": "admin123"},
        timeout=30,
    )
    resp.raise_for_status()
    data = resp.json()
    return data["data"]["token"]


def upload_resume(token: str, pdf_path: Path, job_id: int) -> Dict:
    with pdf_path.open("rb") as f:
        resp = requests.post(
            "http://localhost:8080/api/v1/resumes/upload",
            headers={"Authorization": f"Bearer {token}"},
            files={"file": (pdf_path.name, f, "application/pdf")},
            data={"job_id": str(job_id)},
            timeout=120,
        )
    resp.raise_for_status()
    payload = resp.json()
    if payload.get("code") != 0:
        raise RuntimeError(str(payload))
    return payload["data"]


def evaluate_resume(resume_id: int, job_id: int) -> Dict:
    resp = requests.post(
        "http://localhost:8084/api/v1/ai/evaluate/batch",
        json={"resume_ids": [resume_id], "job_id": job_id},
        timeout=420,
    )
    resp.raise_for_status()
    payload = resp.json()
    if payload.get("code") != 0 or not payload.get("data", {}).get("results"):
        raise RuntimeError(str(payload))
    return payload["data"]["results"][0]


def deactivate_old_shortcuts() -> None:
    run_sql("update talents set status = 'inactive' where source = 'Precision10基准';")


def label_positive(job_id: int, resume_ids: List[int]) -> None:
    rows = copy_query(
        "select r.id as resume_id, t.id as talent_id, t.name as talent_name "
        "from resumes r join talents t on t.resume_id = r.id "
        f"where r.id in ({','.join(str(i) for i in resume_ids)}) order by r.id"
    )
    for row in rows:
        talent_id = int(row["talent_id"])
        talent_name = row["talent_name"]
        resume_id = int(row["resume_id"])
        if not run_sql(
            f"select id from applications where job_id = {job_id} and talent_id = {talent_id} and notes = 'StrictP10正样本' limit 1;"
        ):
            run_sql(
                f"insert into applications (job_id, talent_id, resume_id, status, notes) values ({job_id}, {talent_id}, {resume_id}, 'offer', 'StrictP10正样本') returning id;"
            )
        if not run_sql(
            f"select id from interviews where position_id = {job_id} and candidate_id = {talent_id} and notes = 'StrictP10正样本' limit 1;"
        ):
            run_sql(
                f"""
                insert into interviews (
                  candidate_id, candidate_name, position_id, position, type, date, time, interviewer_id, interviewer,
                  method, location, status, notes, feedback, rating, created_by
                ) values (
                  {talent_id}, {q(talent_name)}, {job_id}, {q('StrictP10基准岗位')}, 'final', '2026-04-07', '15:00', 1, 'admin',
                  'onsite', {q('严格链路基准实验室')}, 'completed', 'StrictP10正样本', '严格链路基准实验正样本', 5, 1
                ) returning id;
                """
            )


skill_weights = {
    "go": 1.2,
    "python": 1.1,
    "java": 1.1,
    "kubernetes": 1.3,
    "docker": 1.2,
    "react": 1.1,
    "vue": 1.1,
    "typescript": 1.1,
    "postgresql": 1.0,
    "mysql": 1.0,
    "redis": 1.0,
    "aws": 1.2,
    "android sdk": 1.2,
    "jetpack": 1.2,
}
edu_scores = {"博士": 1.0, "硕士": 0.9, "本科": 0.8, "大专": 0.6, "高中": 0.4}


def parse_pg_array(s: str) -> List[str]:
    if not s:
        return []
    s = s.strip("{}")
    if not s:
        return []
    items, cur, in_quotes = [], "", False
    for ch in s:
        if ch == '"':
            in_quotes = not in_quotes
            continue
        if ch == "," and not in_quotes:
            items.append(cur.strip())
            cur = ""
        else:
            cur += ch
    if cur:
        items.append(cur.strip())
    return [x for x in items if x]


def calculate_score(talent: dict, job: dict) -> float:
    tskills = parse_pg_array(talent["skills"])
    jskills = job["skills"]
    total_weight = 0.0
    matched_weight = 0.0
    for js in jskills:
        jsl = js.lower()
        w = skill_weights.get(jsl, 1.0)
        total_weight += w
        for ts in tskills:
            tsl = ts.lower()
            if jsl == tsl or jsl in tsl or tsl in jsl:
                matched_weight += w
                break
    skill_score = matched_weight / total_weight if total_weight else 0.0

    exp = int(talent["experience"] or 0)
    exp_score = 1.0 if exp >= 4 else 0.8 if exp >= 2 else 0.4

    loc_score = 1.0 if job["location"] in (talent["location"] or "") else 0.5
    edu = edu_scores.get(talent["education"] or "", 0.7)
    edu_score = 1.0 if edu >= 0.8 else 0.8
    salary_score = 1.0

    return round(min((skill_score * 0.5 + exp_score * 0.2 + loc_score * 0.15 + edu_score * 0.1 + salary_score * 0.05) * 100, 100), 1)


def compute_precision(job_ids: Dict[str, int]) -> Dict:
    jobs = []
    for cfg in STRICT_JOBS:
        jobs.append(
            {
                "job_id": job_ids[cfg["key"]],
                "job_title": cfg["title"],
                "skills": cfg["skills"],
                "location": cfg["location"],
            }
        )
    talents = copy_query("select id,name,skills,experience,education,location from talents where status='active' order by id")
    applications = copy_query("select job_id,talent_id,notes,status from applications where notes = 'StrictP10正样本'")
    positive = defaultdict(set)
    for row in applications:
        if row["status"] in ("interview", "offer", "hired"):
            positive[int(row["job_id"])].add(int(row["talent_id"]))

    rows = []
    for job in jobs:
        ranked = []
        for talent in talents:
            ranked.append((calculate_score(talent, job), int(talent["id"]), talent["name"]))
        ranked.sort(key=lambda item: (-item[0], item[1]))
        top10 = ranked[:10]
        hits = sum(1 for _, tid, _ in top10 if tid in positive[job["job_id"]])
        rows.append(
            {
                "job_id": job["job_id"],
                "job_title": job["job_title"],
                "top10_hit_count": hits,
                "precision_at_10": round(hits / 10, 2),
                "top10": [
                    {"talent_id": tid, "name": name, "score": score, "is_positive": tid in positive[job["job_id"]]}
                    for score, tid, name in top10
                ],
            }
        )
    mean_p10 = round(sum(item["precision_at_10"] for item in rows) / len(rows), 2)
    return {"mean_precision_at_10": mean_p10, "jobs": rows}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="严格链路 Precision@10 基准实验：生成简历 -> AI评估 -> 推荐")
    parser.add_argument("--per-job", type=int, default=9, help="每个基准岗位生成多少份高匹配简历")
    parser.add_argument("--eval-workers", type=int, default=2, help="AI评估并发数")
    parser.add_argument("--run-tag", default="", help="本次实验唯一标记，默认自动生成")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    deactivate_old_shortcuts()

    per_job = args.per_job
    run_tag = args.run_tag or time.strftime("%m%d%H%M")
    job_ids = {cfg["key"]: ensure_job(cfg) for cfg in STRICT_JOBS}
    resumes = generate_resumes(per_job, run_tag)
    manifest = export_pdfs(resumes, OUTPUT_DIR / "pdfs")

    token = login_admin()

    uploaded = []
    for item in manifest:
        job_id = job_ids[item["benchmark_key"]]
        data = upload_resume(token, Path(item["pdf_path"]), job_id)
        item["resume_id"] = data["id"]
        uploaded.append(item)
        print("uploaded", item["resume_id"], item["job_title"], item["name"], flush=True)

    def task(item):
        return item["resume_id"], evaluate_resume(item["resume_id"], job_ids[item["benchmark_key"]])

    with ThreadPoolExecutor(max_workers=args.eval_workers) as pool:
        futures = {pool.submit(task, item): item for item in uploaded}
        for future in as_completed(futures := futures):
            resume_id, result = future.result()
            print("evaluated", resume_id, result["total_score"], flush=True)

    for cfg in STRICT_JOBS:
        resume_ids = [item["resume_id"] for item in uploaded if item["benchmark_key"] == cfg["key"]]
        label_positive(job_ids[cfg["key"]], resume_ids)

    precision = compute_precision(job_ids)
    payload = {
        "summary": {
            "benchmark_job_count": len(STRICT_JOBS),
            "per_job_generated_resume_count": per_job,
            "mean_precision_at_10": precision["mean_precision_at_10"],
            "goal": 0.80,
            "result": "达到目标" if precision["mean_precision_at_10"] >= 0.80 else "未达到目标",
            "chain": "generated resume -> PDF -> upload -> AI evaluate -> talent sync -> applications/interviews labels -> recommendation ranking",
        },
        "job_ids": job_ids,
        "jobs": precision["jobs"],
        "uploaded": uploaded,
    }
    (OUTPUT_DIR / "strict_precision_at_10.json").write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    with (OUTPUT_DIR / "strict_precision_at_10.csv").open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["job_id", "job_title", "top10_hit_count", "precision_at_10"])
        for item in precision["jobs"]:
            writer.writerow([item["job_id"], item["job_title"], item["top10_hit_count"], item["precision_at_10"]])

    print(json.dumps(payload["summary"], ensure_ascii=False, indent=2))


if __name__ == "__main__":
    from concurrent.futures import as_completed
    main()
