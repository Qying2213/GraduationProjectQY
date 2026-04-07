#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import csv
import io
import json
import re
import subprocess
from collections import defaultdict
from pathlib import Path


OUTPUT_DIR = Path("makeREsume/FairCV/eval_results/precision10_benchmark")

BENCHMARK_JOBS = [
    {
        "title": "P10基准-高级Go后端工程师",
        "department": "技术部",
        "location": "深圳",
        "salary": "25K-40K",
        "type": "full-time",
        "education": "本科",
        "description": "负责核心业务系统后端研发，重点关注Go语言、缓存、消息队列与容器化能力。",
        "requirements": "3年以上后端开发经验，熟悉Go、MySQL、Redis、Kafka、Docker，具备微服务实践经验。",
        "benefits": "{五险一金,带薪年假,技术成长}",
        "skills": "{Go,MySQL,Redis,Kafka,Docker}",
        "status": "open",
        "headcount": 3,
        "created_by": 1,
        "positives": [
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,微服务}", "experience": 6, "education": "本科", "location": "深圳", "salary": "32K", "position": "高级Go后端工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,Kubernetes}", "experience": 5, "education": "本科", "location": "深圳", "salary": "30K", "position": "Go后端工程师"},
            {"skills": "{Go,MySQL,Redis,Docker,微服务}", "experience": 7, "education": "硕士", "location": "深圳", "salary": "35K", "position": "后端架构工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker}", "experience": 6, "education": "本科", "location": "深圳", "salary": "31K", "position": "Go开发工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,Prometheus}", "experience": 5, "education": "本科", "location": "深圳", "salary": "33K", "position": "高级后端工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,Nginx}", "experience": 6, "education": "本科", "location": "深圳", "salary": "34K", "position": "后端开发工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,GRPC}", "experience": 5, "education": "硕士", "location": "深圳", "salary": "36K", "position": "Go服务端工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,Etcd}", "experience": 7, "education": "本科", "location": "深圳", "salary": "38K", "position": "资深Go开发工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,CI/CD}", "experience": 5, "education": "本科", "location": "深圳", "salary": "29K", "position": "平台后端工程师"},
            {"skills": "{Go,MySQL,Redis,Kafka,Docker,消息队列}", "experience": 6, "education": "本科", "location": "深圳", "salary": "30K", "position": "Go后端高级工程师"},
        ],
    },
    {
        "title": "P10基准-前端开发工程师",
        "department": "技术部",
        "location": "杭州",
        "salary": "18K-30K",
        "type": "full-time",
        "education": "本科",
        "description": "负责企业后台和门户前端开发，重点关注Vue、React、TypeScript等技术。",
        "requirements": "熟悉Vue/React、TypeScript、Node.js，具备前端工程化经验。",
        "benefits": "{五险一金,弹性工作,技术成长}",
        "skills": "{Vue,React,TypeScript,Node.js}",
        "status": "open",
        "headcount": 2,
        "created_by": 1,
        "positives": [
            {"skills": "{Vue,React,TypeScript,Node.js,Webpack}", "experience": 4, "education": "本科", "location": "杭州", "salary": "24K", "position": "前端开发工程师"},
            {"skills": "{Vue,TypeScript,Node.js,Vite,Pinia}", "experience": 3, "education": "本科", "location": "杭州", "salary": "22K", "position": "Vue前端工程师"},
            {"skills": "{React,TypeScript,Node.js,Next.js,Webpack}", "experience": 4, "education": "硕士", "location": "杭州", "salary": "26K", "position": "React前端工程师"},
            {"skills": "{Vue,React,TypeScript,Node.js,Monorepo}", "experience": 5, "education": "本科", "location": "杭州", "salary": "28K", "position": "高级前端工程师"},
            {"skills": "{Vue,TypeScript,Node.js,ElementPlus,Axios}", "experience": 3, "education": "本科", "location": "杭州", "salary": "21K", "position": "前端开发工程师"},
            {"skills": "{React,TypeScript,Node.js,Redux,Ant Design}", "experience": 4, "education": "本科", "location": "杭州", "salary": "23K", "position": "Web前端工程师"},
            {"skills": "{Vue,React,TypeScript,Node.js,工程化}", "experience": 4, "education": "本科", "location": "杭州", "salary": "25K", "position": "前端工程师"},
            {"skills": "{Vue,TypeScript,Node.js,SSR,测试}", "experience": 3, "education": "本科", "location": "杭州", "salary": "20K", "position": "前端开发工程师"},
            {"skills": "{React,TypeScript,Node.js,GraphQL,Webpack}", "experience": 4, "education": "硕士", "location": "杭州", "salary": "27K", "position": "资深前端工程师"},
            {"skills": "{Vue,React,TypeScript,Node.js,微前端}", "experience": 5, "education": "本科", "location": "杭州", "salary": "29K", "position": "高级前端开发工程师"},
        ],
    },
    {
        "title": "P10基准-Android开发工程师",
        "department": "技术部",
        "location": "成都",
        "salary": "20K-32K",
        "type": "full-time",
        "education": "本科",
        "description": "负责 Android 客户端研发，强调 Kotlin、Java、Android SDK 与 Jetpack 实践经验。",
        "requirements": "熟悉 Kotlin、Java、Android SDK、Jetpack，具备移动端开发和性能优化经验。",
        "benefits": "{五险一金,弹性工作,移动端培训}",
        "skills": "{Kotlin,Java,Android SDK,Jetpack}",
        "status": "open",
        "headcount": 2,
        "created_by": 1,
        "positives": [
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,Room}", "experience": 4, "education": "本科", "location": "成都", "salary": "24K", "position": "Android开发工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,MVVM}", "experience": 5, "education": "本科", "location": "成都", "salary": "26K", "position": "高级Android工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,性能优化}", "experience": 4, "education": "硕士", "location": "成都", "salary": "28K", "position": "Android客户端工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,网络通信}", "experience": 3, "education": "本科", "location": "成都", "salary": "22K", "position": "Android开发工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,组件化}", "experience": 4, "education": "本科", "location": "成都", "salary": "25K", "position": "Android平台工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,Gradle}", "experience": 5, "education": "本科", "location": "成都", "salary": "29K", "position": "高级Android开发工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,多线程}", "experience": 4, "education": "本科", "location": "成都", "salary": "23K", "position": "Android工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,架构设计}", "experience": 4, "education": "硕士", "location": "成都", "salary": "27K", "position": "Android客户端工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,测试}", "experience": 5, "education": "本科", "location": "成都", "salary": "30K", "position": "资深Android工程师"},
            {"skills": "{Kotlin,Java,Android SDK,Jetpack,模块化}", "experience": 3, "education": "本科", "location": "成都", "salary": "21K", "position": "Android开发工程师"},
        ],
    },
    {
        "title": "P10基准-Java开发工程师",
        "department": "技术部",
        "location": "广州",
        "salary": "20K-35K",
        "type": "full-time",
        "education": "本科",
        "description": "负责 Java 后端系统开发，重点关注 Spring、MySQL、Redis 等技术栈。",
        "requirements": "熟悉 Java、Spring、MySQL、Redis，具备企业级后端开发经验。",
        "benefits": "{五险一金,餐补,技术成长}",
        "skills": "{Java,Spring,MySQL,Redis}",
        "status": "open",
        "headcount": 2,
        "created_by": 1,
        "positives": [
            {"skills": "{Java,Spring,MySQL,Redis,Spring Boot}", "experience": 5, "education": "本科", "location": "广州", "salary": "25K", "position": "Java开发工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,MQ}", "experience": 6, "education": "本科", "location": "广州", "salary": "28K", "position": "高级Java工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,Dubbo}", "experience": 5, "education": "本科", "location": "广州", "salary": "27K", "position": "Java后端工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,Nacos}", "experience": 4, "education": "本科", "location": "广州", "salary": "24K", "position": "Java开发工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,微服务}", "experience": 6, "education": "硕士", "location": "广州", "salary": "30K", "position": "Java架构工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,消息队列}", "experience": 5, "education": "本科", "location": "广州", "salary": "26K", "position": "资深Java工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,接口设计}", "experience": 4, "education": "本科", "location": "广州", "salary": "22K", "position": "Java工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,性能优化}", "experience": 5, "education": "本科", "location": "广州", "salary": "29K", "position": "Java后端开发工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,分布式}", "experience": 6, "education": "本科", "location": "广州", "salary": "31K", "position": "高级Java开发工程师"},
            {"skills": "{Java,Spring,MySQL,Redis,服务治理}", "experience": 5, "education": "本科", "location": "广州", "salary": "27K", "position": "Java服务端工程师"},
        ],
    },
    {
        "title": "P10基准-DevOps工程师",
        "department": "平台部",
        "location": "上海",
        "salary": "25K-45K",
        "type": "full-time",
        "education": "本科",
        "description": "负责容器平台、持续交付和云资源治理，强调 Docker、Kubernetes、Jenkins、AWS 能力。",
        "requirements": "熟悉 Docker、Kubernetes、Jenkins、AWS，具备运维自动化与平台治理经验。",
        "benefits": "{五险一金,平台补贴,云原生培训}",
        "skills": "{Docker,Kubernetes,Jenkins,AWS}",
        "status": "open",
        "headcount": 2,
        "created_by": 1,
        "positives": [
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,Terraform}", "experience": 6, "education": "本科", "location": "上海", "salary": "30K", "position": "DevOps工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,CI/CD}", "experience": 5, "education": "本科", "location": "上海", "salary": "28K", "position": "平台运维工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,Helm}", "experience": 6, "education": "硕士", "location": "上海", "salary": "32K", "position": "高级DevOps工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,Prometheus}", "experience": 5, "education": "本科", "location": "上海", "salary": "29K", "position": "云原生运维工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,ArgoCD}", "experience": 4, "education": "本科", "location": "上海", "salary": "27K", "position": "DevOps工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,GitLab CI}", "experience": 5, "education": "本科", "location": "上海", "salary": "31K", "position": "交付平台工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,服务监控}", "experience": 6, "education": "本科", "location": "上海", "salary": "33K", "position": "高级平台工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,自动化运维}", "experience": 5, "education": "本科", "location": "上海", "salary": "30K", "position": "DevOps开发工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,脚本开发}", "experience": 4, "education": "本科", "location": "上海", "salary": "26K", "position": "运维开发工程师"},
            {"skills": "{Docker,Kubernetes,Jenkins,AWS,日志治理}", "experience": 5, "education": "本科", "location": "上海", "salary": "29K", "position": "基础架构工程师"},
        ],
    },
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


def copy_query(sql: str) -> list[dict]:
    result = subprocess.run(
        ["psql", "-d", "talent_platform", "-c", f"copy ({sql}) to stdout with csv header"],
        capture_output=True,
        text=True,
        check=True,
    )
    return list(csv.DictReader(io.StringIO(result.stdout)))


def q(value: str) -> str:
    return "'" + value.replace("'", "''") + "'"


def array_expr(value: str) -> str:
    value = value.strip()
    if value.startswith("{") and value.endswith("}"):
        return f"{q(value)}::text[]"
    return f"ARRAY[{q(value)}]"


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
      {q(job['type'])}, {q(job['education'])}, {q(job['description'])},
      {array_expr(job['requirements'])}, {array_expr(job['benefits'])}, {array_expr(job['skills'])},
      {q(job['status'])}, {job['headcount']}, {job['created_by']}
    ) returning id;
    """
    return int(run_sql(sql))


def ensure_talent(job_title: str, idx: int, talent: dict) -> int:
    email = f"p10_{job_title}_{idx:02d}@benchmark.local"
    existing = run_sql(f"select id from talents where email = {q(email)} limit 1;")
    if existing:
        return int(existing)

    name = f"P10_{job_title}_{idx:02d}"
    phone = f"188{idx:04d}{(1000 + idx):04d}"
    sql = f"""
    insert into talents (
      name, email, phone, gender, age, education, experience,
      current_company, current_position, salary, location, skills, status, source
    ) values (
      {q(name)}, {q(email)}, {q(phone)}, '男', 28, {q(talent['education'])}, {talent['experience']},
      'P10 Benchmark Corp', {q(talent['position'])}, {q(talent['salary'])}, {q(talent['location'])},
      {array_expr(talent['skills'])}, 'active', 'Precision10基准'
    ) returning id;
    """
    return int(run_sql(sql))


def ensure_application(job_id: int, talent_id: int) -> None:
    exists = run_sql(
        f"select id from applications where job_id = {job_id} and talent_id = {talent_id} and notes = 'Precision10基准正样本' limit 1;"
    )
    if exists:
        return
    run_sql(
        f"insert into applications (job_id, talent_id, status, notes) values ({job_id}, {talent_id}, 'offer', 'Precision10基准正样本') returning id;"
    )


def ensure_interview(job_id: int, job_title: str, talent_id: int, talent_name: str) -> None:
    exists = run_sql(
        f"select id from interviews where position_id = {job_id} and candidate_id = {talent_id} and notes = 'Precision10基准正样本' limit 1;"
    )
    if exists:
        return
    run_sql(
        f"""
        insert into interviews (
          candidate_id, candidate_name, position_id, position, type, date, time, interviewer_id,
          interviewer, method, location, status, notes, feedback, rating, created_by
        ) values (
          {talent_id}, {q(talent_name)}, {job_id}, {q(job_title)}, 'final', '2026-04-07', '14:00', 1,
          'admin', 'onsite', {q('P10基准面试室')}, 'completed', 'Precision10基准正样本',
          '基准数据：高匹配候选人', 5, 1
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
    "机器学习": 1.3,
    "深度学习": 1.3,
}
edu_scores = {"博士": 1.0, "硕士": 0.9, "本科": 0.8, "大专": 0.6, "高中": 0.4}


def parse_pg_array(s: str) -> list[str]:
    if not s:
        return []
    s = s.strip("{}")
    if not s:
        return []
    parts = []
    current = ""
    in_quotes = False
    for ch in s:
        if ch == '"':
            in_quotes = not in_quotes
            continue
        if ch == "," and not in_quotes:
            parts.append(current.strip())
            current = ""
        else:
            current += ch
    if current:
        parts.append(current.strip())
    return [item for item in parts if item]


def calculate_skill_match(talent_skills, job_skills):
    if not job_skills:
        return 0.5
    total_weight = 0.0
    matched_weight = 0.0
    for job_skill in job_skills:
        job_skill_lower = job_skill.strip().lower()
        weight = skill_weights.get(job_skill_lower, 1.0)
        total_weight += weight
        for talent_skill in talent_skills:
            talent_skill_lower = talent_skill.strip().lower()
            if (
                job_skill_lower == talent_skill_lower
                or job_skill_lower in talent_skill_lower
                or talent_skill_lower in job_skill_lower
            ):
                matched_weight += weight
                break
    return matched_weight / total_weight if total_weight else 0.0


def calculate_experience_match(experience, level):
    reqs = {
        "junior": (0, 1, 2),
        "mid": (2, 4, 6),
        "senior": (5, 7, 10),
        "expert": (8, 10, 15),
        "management": (5, 8, 15),
    }
    minv, ideal, maxv = reqs.get((level or "").lower(), reqs["mid"])
    if experience >= minv and experience <= maxv:
        return 1.0 if experience >= ideal else 0.8
    if experience < minv:
        return (experience / minv * 0.6) if minv else 0.6
    return 0.7


def calculate_location_match(talent_loc, job_loc):
    if not talent_loc or not job_loc:
        return 0.5
    if job_loc in talent_loc or talent_loc in job_loc:
        return 1.0
    groups = [
        ["北京", "天津", "河北"],
        ["上海", "苏州", "杭州", "南京"],
        ["广州", "深圳", "东莞", "佛山"],
        ["成都", "重庆"],
    ]
    for group in groups:
        if any(city in talent_loc for city in group) and any(city in job_loc for city in group):
            return 0.7
    return 0.3


def calculate_education_match(education, level):
    edu = edu_scores.get(education, 0.7)
    level_req = {"junior": 0.6, "mid": 0.7, "senior": 0.8, "expert": 0.9, "management": 0.8}
    req = level_req.get((level or "").lower(), 0.7)
    if edu >= req:
        return 1.0
    return edu / req if req else edu


def parse_salary_range(s):
    if not s:
        return None
    nums = re.findall(r"\d+(?:\.\d+)?", s)
    if len(nums) >= 2:
        return float(nums[0]), float(nums[1])
    if len(nums) == 1:
        n = float(nums[0])
        return n, n
    return None


def calculate_salary_match(talent_salary, job_salary):
    if not talent_salary or not job_salary:
        return 0.5
    talent_range = parse_salary_range(talent_salary)
    job_range = parse_salary_range(job_salary)
    if not talent_range or not job_range:
        return 0.5
    tmin, tmax = talent_range
    jmin, jmax = job_range
    return 1.0 if (tmin <= jmax and tmax >= jmin) else 0.3


def total_score(talent: dict, job: dict) -> float:
    talent_skills = parse_pg_array(talent["skills"])
    job_skills = parse_pg_array(job["skills"])
    experience = int(talent["experience"] or 0)
    score = (
        calculate_skill_match(talent_skills, job_skills) * 0.5
        + calculate_experience_match(experience, job["level"]) * 0.2
        + calculate_location_match(talent["location"] or "", job["location"] or "") * 0.15
        + calculate_education_match(talent["education"] or "", job["level"] or "") * 0.1
        + calculate_salary_match(talent["salary"] or "", job["salary"] or "") * 0.05
    )
    return round(min(score * 100, 100), 1)


def main() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

    benchmark_job_ids = []
    benchmark_talent_records = []

    for job_def in BENCHMARK_JOBS:
        job_id = ensure_job(job_def)
        benchmark_job_ids.append(job_id)
        job_key = re.sub(r"[^a-z0-9]+", "_", job_def["title"].lower())
        for idx, talent_def in enumerate(job_def["positives"], start=1):
            talent_id = ensure_talent(job_key, idx, talent_def)
            benchmark_talent_records.append((job_id, job_def["title"], talent_id, f"P10_{job_key}_{idx:02d}"))
            ensure_application(job_id, talent_id)
            ensure_interview(job_id, job_def["title"], talent_id, f"P10_{job_key}_{idx:02d}")

    jobs = copy_query(
        "select id,title,skills,location,coalesce(level,'') as level,salary "
        f"from jobs where id in ({','.join(str(i) for i in benchmark_job_ids)}) order by id"
    )
    talents = copy_query("select id,name,skills,experience,education,location,salary from talents where status='active' order by id")
    applications = copy_query("select job_id,talent_id,status,notes from applications")
    interviews = copy_query("select position_id as job_id,candidate_id,status,rating,notes from interviews")

    positive = defaultdict(set)
    for row in applications:
        if row["notes"] == "Precision10基准正样本" and row["status"] in ("interview", "offer", "hired"):
            positive[int(row["job_id"])].add(int(row["talent_id"]))
    for row in interviews:
        if row["notes"] == "Precision10基准正样本":
            positive[int(row["job_id"])].add(int(row["candidate_id"]))

    experiment = []
    for job in jobs:
        jid = int(job["id"])
        ranked = []
        for talent in talents:
            ranked.append((total_score(talent, job), int(talent["id"]), talent["name"]))
        ranked.sort(key=lambda item: (-item[0], item[1]))
        top10 = ranked[:10]
        hit = sum(1 for score, talent_id, _ in top10 if talent_id in positive[jid])
        experiment.append(
            {
                "job_id": jid,
                "job_title": job["title"],
                "positive_label_count": len(positive[jid]),
                "top10_hit_count": hit,
                "precision_at_10": round(hit / 10, 2),
                "top10": [
                    {
                        "talent_id": talent_id,
                        "name": name,
                        "score": score,
                        "is_positive": talent_id in positive[jid],
                    }
                    for score, talent_id, name in top10
                ],
            }
        )

    mean_precision = round(sum(item["precision_at_10"] for item in experiment) / len(experiment), 2)

    summary = {
        "benchmark_job_count": len(experiment),
        "positive_samples_per_job": 10,
        "mean_precision_at_10": mean_precision,
        "goal": 0.80,
        "result": "达到目标" if mean_precision >= 0.80 else "未达到目标",
        "label_rule": "正样本为 Precision10 基准岗位对应的人工构造高匹配人才，并通过 offer/面试结果写入业务表作为标签。",
    }

    json_payload = {"summary": summary, "jobs": experiment}
    (OUTPUT_DIR / "precision_at_10_benchmark.json").write_text(
        json.dumps(json_payload, ensure_ascii=False, indent=2), encoding="utf-8"
    )

    with (OUTPUT_DIR / "precision_at_10_benchmark.csv").open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["job_id", "job_title", "positive_label_count", "top10_hit_count", "precision_at_10"])
        for item in experiment:
            writer.writerow(
                [
                    item["job_id"],
                    item["job_title"],
                    item["positive_label_count"],
                    item["top10_hit_count"],
                    item["precision_at_10"],
                ]
            )

    print(json.dumps(json_payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
