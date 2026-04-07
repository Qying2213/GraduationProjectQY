#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import argparse
import csv
import json
import statistics
import threading
import time
from collections import Counter
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, Iterable, List, Tuple

import requests


GRADE_ORDER = {"S": 0, "A": 1, "B": 2, "C": 3, "D": 4, "E": 5, "F": 6}
DIMENSION_FIELDS = [
    "total_score",
    "jd_match_score",
    "age_score",
    "experience_score",
    "education_score",
    "company_score",
    "tech_score",
    "project_score",
]
thread_local = threading.local()


def now_str() -> str:
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def parse_resume_ids(raw: str, start_id: int, count: int) -> List[int]:
    if not raw.strip():
        return list(range(start_id, start_id + count))

    values: List[int] = []
    for part in raw.split(","):
        piece = part.strip()
        if not piece:
            continue
        if "-" in piece:
            left, right = piece.split("-", 1)
            start = int(left.strip())
            end = int(right.strip())
            step = 1 if end >= start else -1
            values.extend(range(start, end + step, step))
        else:
            values.append(int(piece))
    return values


def get_session() -> requests.Session:
    session = getattr(thread_local, "session", None)
    if session is None:
        session = requests.Session()
        thread_local.session = session
    return session


def safe_write_json(path: Path, payload: Dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    temp_path = path.with_suffix(path.suffix + ".tmp")
    temp_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    temp_path.replace(path)


def load_manifest(manifest_path: Path, resume_id_start: int) -> Dict[int, Dict[str, Any]]:
    mapping: Dict[int, Dict[str, Any]] = {}
    with manifest_path.open("r", encoding="utf-8-sig", newline="") as f:
        reader = csv.DictReader(f)
        for row in reader:
            index = int(row["index"])
            resume_id = resume_id_start + index - 1
            full_name = row["filename"]
            mapping[resume_id] = {
                "index": index,
                "resume_id": resume_id,
                "filename": Path(full_name).name,
                "relative_path": full_name,
                "candidate_name": row.get("name", ""),
                "position": row.get("position", ""),
                "pages": int(row.get("pages") or 0),
            }
    return mapping


def normalize_result(resume_id: int, result: Dict[str, Any], manifest_row: Dict[str, Any]) -> Dict[str, Any]:
    matched_skills = result.get("matched_skills") or []
    missing_skills = result.get("missing_skills") or []
    candidate_name = manifest_row.get("candidate_name") or result.get("candidate_name") or ""

    return {
        "resume_id": resume_id,
        "index": manifest_row.get("index"),
        "filename": manifest_row.get("filename"),
        "relative_path": manifest_row.get("relative_path"),
        "candidate_name": candidate_name,
        "position": manifest_row.get("position"),
        "pages": manifest_row.get("pages"),
        "total_score": round(float(result.get("total_score") or 0), 2),
        "grade": result.get("grade") or "",
        "jd_match_score": int(result.get("jd_match_score") or 0),
        "age_score": int(result.get("age_score") or 0),
        "experience_score": int(result.get("experience_score") or 0),
        "education_score": int(result.get("education_score") or 0),
        "company_score": int(result.get("company_score") or 0),
        "tech_score": int(result.get("tech_score") or 0),
        "project_score": int(result.get("project_score") or 0),
        "recommendation": result.get("recommendation") or "",
        "matched_skills": matched_skills,
        "missing_skills": missing_skills,
        "summary": result.get("summary") or "",
    }


def evaluate_resume(base_url: str, resume_id: int, job_id: int, timeout_sec: int) -> Tuple[Dict[str, Any], float]:
    started_at = time.time()
    session = get_session()
    resp = session.post(
        f"{base_url.rstrip('/')}/ai/evaluate/batch",
        json={"resume_ids": [resume_id], "job_id": job_id},
        timeout=timeout_sec,
    )
    resp.raise_for_status()

    payload = resp.json()
    if payload.get("code") != 0:
        raise RuntimeError(f"接口返回失败: {payload}")

    data = payload.get("data", {})
    results = data.get("results") or []
    errors = data.get("errors") or []
    if not results:
        raise RuntimeError("未返回评估结果: " + "; ".join(str(item) for item in errors))

    return results[0], round(time.time() - started_at, 2)


def build_checkpoint(
    config: Dict[str, Any],
    resume_ids: List[int],
    manifest_map: Dict[int, Dict[str, Any]],
    succeeded: Dict[int, Dict[str, Any]],
    failed: Dict[int, Dict[str, Any]],
    started_at: str,
) -> Dict[str, Any]:
    ordered_records: List[Dict[str, Any]] = []
    for resume_id in resume_ids:
        if resume_id in succeeded:
            ordered_records.append(succeeded[resume_id])
            continue
        if resume_id in failed:
            ordered_records.append(failed[resume_id])
            continue

        manifest_row = manifest_map.get(resume_id, {})
        ordered_records.append(
            {
                "resume_id": resume_id,
                "index": manifest_row.get("index"),
                "filename": manifest_row.get("filename"),
                "candidate_name": manifest_row.get("candidate_name"),
                "position": manifest_row.get("position"),
                "status": "pending",
                "attempt": 0,
                "elapsed_sec": 0,
                "error": "",
                "evaluated_at": "",
            }
        )

    return {
        "config": config,
        "started_at": started_at,
        "updated_at": now_str(),
        "target_count": len(resume_ids),
        "success_count": len(succeeded),
        "failure_count": len(failed),
        "pending_count": len(resume_ids) - len(succeeded) - len(failed),
        "records": ordered_records,
    }


def load_existing_records(checkpoint_path: Path) -> Tuple[Dict[int, Dict[str, Any]], Dict[int, Dict[str, Any]]]:
    if not checkpoint_path.exists():
        return {}, {}

    payload = json.loads(checkpoint_path.read_text(encoding="utf-8"))
    succeeded: Dict[int, Dict[str, Any]] = {}
    failed: Dict[int, Dict[str, Any]] = {}

    for item in payload.get("records", []):
        resume_id = int(item["resume_id"])
        status = item.get("status")
        if status == "success":
            succeeded[resume_id] = item
        elif status == "failed":
            failed[resume_id] = item

    return succeeded, failed


def sort_success_records(records: Iterable[Dict[str, Any]]) -> List[Dict[str, Any]]:
    ranked = sorted(
        records,
        key=lambda row: (
            -float(row.get("total_score") or 0),
            GRADE_ORDER.get(str(row.get("grade") or ""), 999),
            -int(row.get("jd_match_score") or 0),
            -int(row.get("tech_score") or 0),
            int(row.get("resume_id") or 0),
        ),
    )
    for idx, row in enumerate(ranked, start=1):
        row["rank"] = idx
    return ranked


def build_summary_metrics(success_rows: List[Dict[str, Any]], failure_rows: List[Dict[str, Any]], target_count: int) -> List[Tuple[str, Any]]:
    metrics: List[Tuple[str, Any]] = [
        ("target_count", target_count),
        ("success_count", len(success_rows)),
        ("failure_count", len(failure_rows)),
        ("success_rate", round(len(success_rows) / target_count * 100, 2) if target_count else 0),
    ]

    if not success_rows:
        return metrics

    grade_counter = Counter(str(row.get("grade") or "") for row in success_rows)
    recommendation_counter = Counter(str(row.get("recommendation") or "") for row in success_rows)

    for field in DIMENSION_FIELDS:
        values = [float(row.get(field) or 0) for row in success_rows]
        metrics.extend(
            [
                (f"{field}_avg", round(statistics.mean(values), 2)),
                (f"{field}_median", round(statistics.median(values), 2)),
                (f"{field}_max", round(max(values), 2)),
                (f"{field}_min", round(min(values), 2)),
            ]
        )

    for grade in sorted(grade_counter, key=lambda item: GRADE_ORDER.get(item, 999)):
        metrics.append((f"grade_{grade}", grade_counter[grade]))

    for recommendation, count in recommendation_counter.items():
        label = recommendation or "empty"
        metrics.append((f"recommendation_{label}", count))

    top_row = max(success_rows, key=lambda row: float(row.get("total_score") or 0))
    bottom_row = min(success_rows, key=lambda row: float(row.get("total_score") or 0))
    metrics.extend(
        [
            ("top_candidate", top_row.get("candidate_name") or top_row.get("filename") or ""),
            ("top_score", top_row.get("total_score") or 0),
            ("bottom_candidate", bottom_row.get("candidate_name") or bottom_row.get("filename") or ""),
            ("bottom_score", bottom_row.get("total_score") or 0),
        ]
    )

    return metrics


def write_detail_csv(path: Path, success_rows: List[Dict[str, Any]], failure_rows: List[Dict[str, Any]]) -> None:
    headers = [
        "rank",
        "resume_id",
        "index",
        "candidate_name",
        "filename",
        "relative_path",
        "position",
        "pages",
        "total_score",
        "grade",
        "jd_match_score",
        "age_score",
        "experience_score",
        "education_score",
        "company_score",
        "tech_score",
        "project_score",
        "recommendation",
        "matched_skills",
        "missing_skills",
        "summary",
        "status",
        "attempt",
        "elapsed_sec",
        "evaluated_at",
        "error",
    ]
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=headers)
        writer.writeheader()

        for row in success_rows:
            writer.writerow(
                {
                    **row,
                    "matched_skills": "、".join(row.get("matched_skills") or []),
                    "missing_skills": "、".join(row.get("missing_skills") or []),
                }
            )

        for row in failure_rows:
            writer.writerow(
                {
                    "rank": "",
                    "resume_id": row.get("resume_id"),
                    "index": row.get("index"),
                    "candidate_name": row.get("candidate_name"),
                    "filename": row.get("filename"),
                    "relative_path": row.get("relative_path"),
                    "position": row.get("position"),
                    "pages": row.get("pages"),
                    "status": row.get("status"),
                    "attempt": row.get("attempt"),
                    "elapsed_sec": row.get("elapsed_sec"),
                    "evaluated_at": row.get("evaluated_at"),
                    "error": row.get("error"),
                }
            )


def write_summary_csv(path: Path, metrics: List[Tuple[str, Any]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["metric", "value"])
        for key, value in metrics:
            writer.writerow([key, value])


def write_markdown_report(path: Path, metrics: Dict[str, Any], success_rows: List[Dict[str, Any]], failure_rows: List[Dict[str, Any]]) -> None:
    target_count = metrics.get("target_count", 0)
    lines = [
        f"# {target_count}份简历AI评估实验报告",
        "",
        f"- 生成时间：{now_str()}",
        f"- 目标样本数：{metrics.get('target_count', 0)}",
        f"- 成功评估：{metrics.get('success_count', 0)}",
        f"- 失败数量：{metrics.get('failure_count', 0)}",
        f"- 成功率：{metrics.get('success_rate', 0)}%",
        "",
        "## 核心统计",
        "",
        f"- 平均总分：{metrics.get('total_score_avg', 0)}",
        f"- 平均JD匹配：{metrics.get('jd_match_score_avg', 0)}",
        f"- 平均技术能力：{metrics.get('tech_score_avg', 0)}",
        f"- 平均项目经验：{metrics.get('project_score_avg', 0)}",
        f"- 最高分候选人：{metrics.get('top_candidate', '')} / {metrics.get('top_score', 0)}",
        f"- 最低分候选人：{metrics.get('bottom_candidate', '')} / {metrics.get('bottom_score', 0)}",
        "",
        "## Top 10",
        "",
        "| 排名 | resume_id | 姓名 | 文件名 | 总分 | 等级 | JD匹配 | 技术 | 项目 | 结论 |",
        "| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |",
    ]

    for row in success_rows[:10]:
        lines.append(
            f"| {row.get('rank', '')} | {row.get('resume_id', '')} | "
            f"{row.get('candidate_name', '')} | {row.get('filename', '')} | "
            f"{row.get('total_score', 0)} | {row.get('grade', '')} | "
            f"{row.get('jd_match_score', 0)} | {row.get('tech_score', 0)} | "
            f"{row.get('project_score', 0)} | {row.get('recommendation', '')} |"
        )

    if failure_rows:
        lines.extend(["", "## 失败记录", ""])
        for row in failure_rows:
            lines.append(
                f"- resume_id={row.get('resume_id')} / {row.get('filename', '')} / "
                f"attempt={row.get('attempt')} / error={row.get('error', '')}"
            )

    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser(description="对已上传到 resume-service 的简历做并发 AI 评估，并导出表格")
    parser.add_argument("--base-url", default="http://localhost:8084/api/v1", help="resume-service API 根地址")
    parser.add_argument("--job-id", type=int, default=1, help="目标岗位 ID")
    parser.add_argument("--resume-id-start", type=int, default=31, help="实验简历起始 resume_id")
    parser.add_argument("--count", type=int, default=100, help="实验简历数量")
    parser.add_argument("--resume-ids", default="", help="自定义 resume_id 范围，例如 31-130 或 31,32,40")
    parser.add_argument("--manifest", default="makeREsume/FairCV/output_backend_100/manifest.csv", help="实验简历 manifest 文件")
    parser.add_argument("--output-dir", default="makeREsume/FairCV/eval_results/backend_100", help="结果输出目录")
    parser.add_argument("--concurrency", type=int, default=10, help="并发评估数")
    parser.add_argument("--timeout-sec", type=int, default=400, help="单次请求超时时间（秒）")
    parser.add_argument("--retries", type=int, default=2, help="失败重试次数")
    args = parser.parse_args()

    output_dir = Path(args.output_dir).resolve()
    manifest_path = Path(args.manifest).resolve()
    checkpoint_path = output_dir / "evaluation_checkpoint.json"
    final_json_path = output_dir / "evaluation_results.json"
    detail_csv_path = output_dir / "evaluation_details.csv"
    summary_csv_path = output_dir / "evaluation_summary.csv"
    report_md_path = output_dir / "evaluation_report.md"

    resume_ids = parse_resume_ids(args.resume_ids, args.resume_id_start, args.count)
    manifest_map = load_manifest(manifest_path, args.resume_id_start)
    started_at = now_str()

    config = {
        "base_url": args.base_url,
        "job_id": args.job_id,
        "resume_id_start": args.resume_id_start,
        "count": args.count,
        "resume_ids": resume_ids,
        "manifest": str(manifest_path),
        "output_dir": str(output_dir),
        "concurrency": args.concurrency,
        "timeout_sec": args.timeout_sec,
        "retries": args.retries,
    }

    succeeded, failed = load_existing_records(checkpoint_path)
    pending_ids = [resume_id for resume_id in resume_ids if resume_id not in succeeded]

    print(f"[{now_str()}] 目标 resume_id 数量: {len(resume_ids)}", flush=True)
    print(f"[{now_str()}] 已有成功缓存: {len(succeeded)}，待处理: {len(pending_ids)}", flush=True)

    for attempt in range(1, args.retries + 2):
        if not pending_ids:
            break

        print(
            f"[{now_str()}] 开始第 {attempt}/{args.retries + 1} 轮评估，"
            f"本轮待处理 {len(pending_ids)} 份，并发数={args.concurrency}",
            flush=True,
        )

        round_failures: Dict[int, Dict[str, Any]] = {}
        with ThreadPoolExecutor(max_workers=args.concurrency) as executor:
            futures = {
                executor.submit(evaluate_resume, args.base_url, resume_id, args.job_id, args.timeout_sec): resume_id
                for resume_id in pending_ids
            }

            completed_count = 0
            total_count = len(pending_ids)
            for future in as_completed(futures):
                resume_id = futures[future]
                manifest_row = manifest_map.get(resume_id, {})
                completed_count += 1

                try:
                    raw_result, elapsed_sec = future.result()
                    normalized = normalize_result(resume_id, raw_result, manifest_row)
                    normalized.update(
                        {
                            "status": "success",
                            "attempt": attempt,
                            "elapsed_sec": elapsed_sec,
                            "error": "",
                            "evaluated_at": now_str(),
                        }
                    )
                    succeeded[resume_id] = normalized
                    failed.pop(resume_id, None)
                    print(
                        f"[{now_str()}] 完成 {completed_count}/{total_count}: "
                        f"resume_id={resume_id} score={normalized['total_score']} grade={normalized['grade']} "
                        f"time={elapsed_sec}s",
                        flush=True,
                    )
                except Exception as exc:
                    failure_row = {
                        "resume_id": resume_id,
                        "index": manifest_row.get("index"),
                        "filename": manifest_row.get("filename"),
                        "relative_path": manifest_row.get("relative_path"),
                        "candidate_name": manifest_row.get("candidate_name"),
                        "position": manifest_row.get("position"),
                        "pages": manifest_row.get("pages"),
                        "status": "failed",
                        "attempt": attempt,
                        "elapsed_sec": 0,
                        "error": str(exc),
                        "evaluated_at": now_str(),
                    }
                    round_failures[resume_id] = failure_row
                    failed[resume_id] = failure_row
                    print(
                        f"[{now_str()}] 失败 {completed_count}/{total_count}: "
                        f"resume_id={resume_id} error={exc}",
                        flush=True,
                    )

                checkpoint = build_checkpoint(config, resume_ids, manifest_map, succeeded, failed, started_at)
                safe_write_json(checkpoint_path, checkpoint)

        pending_ids = [resume_id for resume_id in round_failures if attempt <= args.retries]
        if pending_ids:
            print(f"[{now_str()}] 本轮结束，待重试 {len(pending_ids)} 份", flush=True)

    success_rows = sort_success_records([row for row in succeeded.values() if row.get("status") == "success"])
    failure_rows = sorted(
        [row for row in failed.values() if row.get("status") == "failed" and row["resume_id"] not in succeeded],
        key=lambda row: int(row.get("resume_id") or 0),
    )

    summary_metrics = build_summary_metrics(success_rows, failure_rows, len(resume_ids))
    summary_dict = {key: value for key, value in summary_metrics}

    final_payload = {
        "config": config,
        "started_at": started_at,
        "completed_at": now_str(),
        "target_count": len(resume_ids),
        "success_count": len(success_rows),
        "failure_count": len(failure_rows),
        "summary": summary_dict,
        "success_records": success_rows,
        "failure_records": failure_rows,
    }

    safe_write_json(final_json_path, final_payload)
    write_detail_csv(detail_csv_path, success_rows, failure_rows)
    write_summary_csv(summary_csv_path, summary_metrics)
    write_markdown_report(report_md_path, summary_dict, success_rows, failure_rows)

    print(f"[done] JSON: {final_json_path}", flush=True)
    print(f"[done] 明细表: {detail_csv_path}", flush=True)
    print(f"[done] 汇总表: {summary_csv_path}", flush=True)
    print(f"[done] 报告: {report_md_path}", flush=True)


if __name__ == "__main__":
    main()
