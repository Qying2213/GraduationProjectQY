#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import argparse
import json
import time
from pathlib import Path
from typing import Dict, List

import requests


def login(auth_base_url: str, username: str, password: str) -> str:
    resp = requests.post(
        f"{auth_base_url}/login",
        json={"username": username, "password": password},
        timeout=30,
    )
    resp.raise_for_status()
    data = resp.json()
    token = data.get("data", {}).get("token")
    if not token:
        raise RuntimeError(f"登录失败: {data}")
    return token


def build_jd_text(job_base_url: str, job_id: int) -> str:
    resp = requests.get(f"{job_base_url}/jobs/{job_id}", timeout=30)
    resp.raise_for_status()
    payload = resp.json().get("data", {})

    parts: List[str] = []
    for label, key in [
        ("职位名称", "title"),
        ("所属部门", "department"),
        ("工作地点", "location"),
        ("薪资范围", "salary"),
        ("岗位描述", "description"),
    ]:
        value = payload.get(key)
        if value:
            parts.append(f"{label}：{value}")

    requirements = payload.get("requirements") or []
    if requirements:
        if isinstance(requirements, list):
            parts.append("任职要求：\n- " + "\n- ".join(str(item) for item in requirements if item))
        else:
            parts.append(f"任职要求：{requirements}")

    skills = payload.get("skills") or []
    if skills:
        if isinstance(skills, list):
            parts.append("技能要求：" + "、".join(str(item) for item in skills if item))
        else:
            parts.append(f"技能要求：{skills}")

    return "\n\n".join(parts)


def upload_resume(resume_base_url: str, token: str, pdf_path: Path) -> Dict:
    with pdf_path.open("rb") as f:
        resp = requests.post(
            f"{resume_base_url}/resumes/upload",
            headers={"Authorization": f"Bearer {token}"},
            files={"file": (pdf_path.name, f, "application/pdf")},
            timeout=120,
        )
    resp.raise_for_status()
    data = resp.json()
    if data.get("code") != 0:
        raise RuntimeError(f"上传失败: {pdf_path.name} -> {data}")
    return data.get("data", {})


def batch_evaluate(resume_base_url: str, resume_ids: List[int], jd_text: str) -> Dict:
    resp = requests.post(
        f"{resume_base_url}/ai/evaluate/batch",
        json={"resume_ids": resume_ids, "jd_text": jd_text},
        timeout=1800,
    )
    resp.raise_for_status()
    data = resp.json()
    if data.get("code") != 0:
        raise RuntimeError(f"批量评估失败: {data}")
    return data.get("data", {})


def chunked(values: List[int], size: int) -> List[List[int]]:
    return [values[i : i + size] for i in range(0, len(values), size)]


def parse_resume_ids(raw: str) -> List[int]:
    values = []
    for part in raw.split(","):
        part = part.strip()
        if not part:
            continue
        values.append(int(part))
    return values


def main() -> None:
    parser = argparse.ArgumentParser(description="批量上传 PDF 简历并按 5 份一组进行 AI 评估")
    parser.add_argument("--auth-base-url", default="http://localhost:8080/api/v1", help="登录接口根地址")
    parser.add_argument("--job-base-url", default="http://localhost:8080/api/v1", help="岗位接口根地址")
    parser.add_argument("--resume-base-url", default="http://localhost:8080/api/v1", help="简历/AI 接口根地址")
    parser.add_argument("--input-dir", default="output_backend_100/后端开发工程师", help="PDF 目录")
    parser.add_argument("--username", default="candidate01", help="用于上传的账号")
    parser.add_argument("--password", default="123456", help="用于上传的密码")
    parser.add_argument("--job-id", type=int, default=1, help="目标岗位 ID")
    parser.add_argument("--batch-size", type=int, default=5, help="每批评估数量")
    parser.add_argument("--limit", type=int, default=100, help="最多处理多少份 PDF")
    parser.add_argument("--output", default="batch_eval_summary.json", help="结果汇总文件")
    parser.add_argument("--existing-resume-ids", default="", help="已上传的 resume_id 列表，逗号分隔；提供后将跳过上传")
    args = parser.parse_args()

    auth_base_url = args.auth_base_url.rstrip("/")
    job_base_url = args.job_base_url.rstrip("/")
    resume_base_url = args.resume_base_url.rstrip("/")
    input_dir = Path(args.input_dir).resolve()
    output_path = Path(args.output).resolve()

    jd_text = build_jd_text(job_base_url, args.job_id)
    print(f"[info] 已获取 job_id={args.job_id} 的 JD 文本，长度: {len(jd_text)}", flush=True)

    uploaded: List[Dict] = []
    if args.existing_resume_ids.strip():
        resume_ids = parse_resume_ids(args.existing_resume_ids)
        print(f"[info] 复用已上传简历，共 {len(resume_ids)} 份", flush=True)
    else:
        pdf_files = sorted(input_dir.glob("*.pdf"))[: args.limit]
        if not pdf_files:
            raise FileNotFoundError(f"未找到 PDF: {input_dir}")

        print(f"[info] 待上传 PDF 数量: {len(pdf_files)}", flush=True)
        print(f"[info] 输入目录: {input_dir}", flush=True)

        token = login(auth_base_url, args.username, args.password)
        print(f"[info] 登录成功: {args.username}", flush=True)

        for idx, pdf in enumerate(pdf_files, start=1):
            uploaded_item = upload_resume(resume_base_url, token, pdf)
            uploaded.append(
                {
                    "index": idx,
                    "filename": pdf.name,
                    "resume_id": uploaded_item.get("id"),
                    "talent_id": uploaded_item.get("talent_id"),
                    "job_id": uploaded_item.get("job_id"),
                }
            )
            if idx % 10 == 0 or idx == len(pdf_files):
                print(f"[upload] {idx}/{len(pdf_files)} 已上传", flush=True)

        resume_ids = [item["resume_id"] for item in uploaded if item.get("resume_id")]

    batches = chunked(resume_ids, args.batch_size)

    batch_results: List[Dict] = []
    for batch_index, batch_resume_ids in enumerate(batches, start=1):
        print(f"[eval] 开始第 {batch_index}/{len(batches)} 批，resume_ids={batch_resume_ids}", flush=True)
        started_at = time.time()
        result = batch_evaluate(resume_base_url, batch_resume_ids, jd_text)
        duration = round(time.time() - started_at, 2)
        batch_results.append(
            {
                "batch_index": batch_index,
                "resume_ids": batch_resume_ids,
                "duration_sec": duration,
                "result": result,
            }
        )
        print(
            f"[eval] 第 {batch_index} 批完成，用时 {duration}s，成功 {result.get('success')}，失败 {result.get('failed')}",
            flush=True,
        )

    summary = {
        "auth_base_url": auth_base_url,
        "job_base_url": job_base_url,
        "resume_base_url": resume_base_url,
        "username": args.username,
        "job_id": args.job_id,
        "input_dir": str(input_dir),
        "uploaded_count": len(uploaded),
        "resume_id_count": len(resume_ids),
        "batch_size": args.batch_size,
        "uploaded": uploaded,
        "batches": batch_results,
    }

    output_path.write_text(json.dumps(summary, ensure_ascii=False, indent=2), encoding="utf-8")
    print(f"[done] 已上传 {len(uploaded)} 份并完成 {len(batch_results)} 轮批量评估", flush=True)
    print(f"[done] 汇总文件: {output_path}", flush=True)


if __name__ == "__main__":
    main()
