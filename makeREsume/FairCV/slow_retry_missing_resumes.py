#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import argparse
import csv
import json
import subprocess
import time
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List

import requests


def now_str() -> str:
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def load_missing_ids(path: Path) -> List[int]:
    rows = list(csv.DictReader(path.open(encoding="utf-8-sig")))
    return [int(row["resume_id"]) for row in rows if row.get("accuracy_label") == "missing"]


def parse_resume_ids(raw: str) -> List[int]:
    values: List[int] = []
    for part in raw.split(","):
        piece = part.strip()
        if not piece:
            continue
        if "-" in piece:
            start_str, end_str = piece.split("-", 1)
            start = int(start_str.strip())
            end = int(end_str.strip())
            step = 1 if end >= start else -1
            values.extend(range(start, end + step, step))
        else:
            values.append(int(piece))
    return values


def safe_write_json(path: Path, payload: Dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    tmp = path.with_suffix(path.suffix + ".tmp")
    tmp.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    tmp.replace(path)


def call_evaluate(base_url: str, resume_id: int, job_id: int, timeout_sec: int) -> requests.Response:
    return requests.post(
        f"{base_url.rstrip('/')}/ai/evaluate/batch",
        json={"resume_ids": [resume_id], "job_id": job_id},
        timeout=timeout_sec,
    )


def refresh_reports(project_root: Path, cohort_dir: Path) -> None:
    status_csv = cohort_dir / "cohort_latest_eval_status.csv"

    export_sql = (
        "\\copy (with latest_eval as (select distinct on (resume_id) "
        "resume_id, id as evaluation_id, created_at as latest_eval_at, parsed_name, match_score, "
        "match_level, report_recommendation from evaluation_results "
        "where resume_id between 31 and 130 order by resume_id, created_at desc, id desc) "
        "select r.id as resume_id, r.file_name, r.status as resume_status, le.evaluation_id, "
        "le.latest_eval_at, coalesce(nullif(le.parsed_name, ''), split_part(r.file_name, '_', 3)) as parsed_name, "
        "le.match_score, le.match_level, le.report_recommendation, "
        "case when le.latest_eval_at >= timestamp with time zone '2026-04-03 16:27:04+08' then 'current_run' "
        "when le.evaluation_id is not null then 'historical' else 'missing' end as evaluation_source "
        "from resumes r left join latest_eval le on le.resume_id = r.id "
        "where r.id between 31 and 130 order by r.id) "
        f"to '{status_csv.as_posix()}' with csv header"
    )
    subprocess.run(
        ["psql", "-d", "talent_platform", "-c", export_sql],
        cwd=project_root,
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )

    subprocess.run(
        ["python3", "makeREsume/FairCV/analyze_cohort_accuracy.py"],
        cwd=project_root,
        check=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="极低频慢速补跑缺失简历样本，尽量避开 Coze 限流")
    parser.add_argument("--base-url", default="http://localhost:8084/api/v1", help="resume-service API 根地址")
    parser.add_argument("--job-id", type=int, default=1, help="目标岗位 ID")
    parser.add_argument(
        "--input",
        default="makeREsume/FairCV/eval_results/backend_100/cohort_accuracy_cleaned.csv",
        help="包含 missing 样本的清洗结果 CSV",
    )
    parser.add_argument(
        "--output",
        default="makeREsume/FairCV/eval_results/backend_100/slow_retry_progress.json",
        help="慢速补跑进度文件",
    )
    parser.add_argument("--resume-ids", default="", help="只补跑指定 resume_id，例如 86 或 86,87,90-95")
    parser.add_argument("--timeout-sec", type=int, default=420, help="单次请求超时秒数")
    parser.add_argument("--delay-sec", type=int, default=90, help="成功后下一次请求前冷却秒数")
    parser.add_argument("--rate-limit-wait-sec", type=int, default=300, help="命中限流后的等待秒数")
    parser.add_argument("--max-attempts-per-id", type=int, default=6, help="单个 resume_id 最多尝试次数")
    parser.add_argument("--project-root", default=".", help="项目根目录")
    args = parser.parse_args()

    project_root = Path(args.project_root).resolve()
    input_path = (project_root / args.input).resolve()
    output_path = (project_root / args.output).resolve()
    cohort_dir = input_path.parent

    missing_ids = parse_resume_ids(args.resume_ids) if args.resume_ids.strip() else load_missing_ids(input_path)
    progress: Dict[str, Any] = {
        "started_at": now_str(),
        "updated_at": now_str(),
        "config": {
            "base_url": args.base_url,
            "job_id": args.job_id,
            "timeout_sec": args.timeout_sec,
            "delay_sec": args.delay_sec,
            "rate_limit_wait_sec": args.rate_limit_wait_sec,
            "max_attempts_per_id": args.max_attempts_per_id,
        },
        "target_missing_ids": missing_ids,
        "completed_resume_ids": [],
        "still_missing_resume_ids": missing_ids[:],
        "attempt_log": [],
    }
    safe_write_json(output_path, progress)

    pending = missing_ids[:]
    print(f"[{now_str()}] 开始顺序慢速补跑，待处理 {len(pending)} 份", flush=True)

    for sequence_index, resume_id in enumerate(missing_ids, start=1):
        if resume_id not in pending:
            continue

        success = False
        for attempt in range(1, args.max_attempts_per_id + 1):
            started_at = time.time()
            error_text = ""
            http_status = 0

            try:
                resp = call_evaluate(args.base_url, resume_id, args.job_id, args.timeout_sec)
                http_status = resp.status_code
                payload = resp.json()
                results = payload.get("data", {}).get("results") or []
                errors = payload.get("data", {}).get("errors") or []
                if resp.status_code == 200 and payload.get("code") == 0 and results:
                    success = True
                else:
                    error_text = "; ".join(str(item) for item in errors) or str(payload)
            except Exception as exc:
                error_text = str(exc)

            elapsed_sec = round(time.time() - started_at, 2)
            print(
                f"[{now_str()}] resume_id={resume_id} "
                f"attempt={attempt}/{args.max_attempts_per_id} "
                f"{'成功' if success else '失败'} http={http_status} time={elapsed_sec}s "
                f"{'' if success else 'error=' + error_text}",
                flush=True,
            )

            progress["attempt_log"].append(
                {
                    "sequence_index": sequence_index,
                    "resume_id": resume_id,
                    "attempt": attempt,
                    "success": success,
                    "http_status": http_status,
                    "elapsed_sec": elapsed_sec,
                    "error": error_text,
                    "timestamp": now_str(),
                }
            )

            if success:
                progress["completed_resume_ids"].append(resume_id)
                progress["completed_resume_ids"] = sorted(set(progress["completed_resume_ids"]))
                pending = [item for item in pending if item != resume_id]
                refresh_reports(project_root, cohort_dir)
                progress["still_missing_resume_ids"] = pending[:]
                progress["updated_at"] = now_str()
                safe_write_json(output_path, progress)
                print(f"[{now_str()}] 成功冷却 {args.delay_sec}s", flush=True)
                time.sleep(args.delay_sec)
                break

            progress["still_missing_resume_ids"] = pending[:]
            progress["updated_at"] = now_str()
            safe_write_json(output_path, progress)
            lower_error = error_text.lower()
            if "concurrency limit" in lower_error or "rate" in lower_error:
                print(f"[{now_str()}] 命中限流，等待 {args.rate_limit_wait_sec}s 后继续", flush=True)
                time.sleep(args.rate_limit_wait_sec)
            else:
                time.sleep(args.delay_sec)

        if not success:
            print(
                f"[{now_str()}] resume_id={resume_id} 在 {args.max_attempts_per_id} 次尝试后仍未成功，暂时保留为缺失样本",
                flush=True,
            )

    refresh_reports(project_root, cohort_dir)
    print(f"[done] 仍缺失 {len(pending)} 份: {pending}", flush=True)
    print(f"[done] 进度文件: {output_path}", flush=True)


if __name__ == "__main__":
    main()
