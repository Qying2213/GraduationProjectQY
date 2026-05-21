#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations

import argparse
import csv
import json
from collections import Counter
from datetime import datetime
from pathlib import Path
from statistics import mean, median, quantiles
from typing import Any, Dict, List


def now_str() -> str:
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S")


def load_rows(path: Path) -> List[Dict[str, Any]]:
    with path.open("r", encoding="utf-8-sig", newline="") as f:
        return list(csv.DictReader(f))


def to_float(value: Any) -> float:
    try:
        return float(value)
    except (TypeError, ValueError):
        return 0.0


def compute_threshold(scores: List[float]) -> Dict[str, float]:
    positive_scores = sorted(score for score in scores if score > 0)
    if len(positive_scores) >= 4:
        q1, _, q3 = quantiles(positive_scores, n=4, method="inclusive")
    elif positive_scores:
        q1 = positive_scores[0]
        q3 = positive_scores[-1]
    else:
        q1 = 0.0
        q3 = 0.0

    iqr = q3 - q1
    threshold = max(0.0, round(q1 - 1.5 * iqr, 1))
    return {
        "q1": round(q1, 2),
        "median": round(median(positive_scores), 2) if positive_scores else 0.0,
        "q3": round(q3, 2),
        "iqr": round(iqr, 2),
        "threshold": threshold,
    }


def classify_rows(rows: List[Dict[str, Any]], threshold: float) -> List[Dict[str, Any]]:
    cleaned: List[Dict[str, Any]] = []
    for row in rows:
        score = to_float(row.get("match_score"))
        source = row.get("evaluation_source", "")
        evaluated = source != "missing"
        obvious_error = evaluated and score <= threshold

        if source == "missing":
            label = "missing"
            note = "该简历目前没有最新评估结果"
        elif obvious_error:
            label = "error_data"
            note = f"匹配分 {score} 低于异常阈值 {threshold}"
        else:
            label = "valid"
            note = "结果落在正常分布区间"

        cleaned.append(
            {
                **row,
                "match_score": score,
                "is_evaluated": evaluated,
                "is_obvious_error": obvious_error,
                "accuracy_label": label,
                "analysis_note": note,
            }
        )
    return cleaned


def build_summary(cleaned_rows: List[Dict[str, Any]], threshold_info: Dict[str, float]) -> Dict[str, Any]:
    total = len(cleaned_rows)
    evaluated_rows = [row for row in cleaned_rows if row["is_evaluated"]]
    valid_rows = [row for row in cleaned_rows if row["accuracy_label"] == "valid"]
    error_rows = [row for row in cleaned_rows if row["accuracy_label"] == "error_data"]
    missing_rows = [row for row in cleaned_rows if row["accuracy_label"] == "missing"]

    scores = [row["match_score"] for row in evaluated_rows]
    source_counter = Counter(row.get("evaluation_source", "") for row in cleaned_rows)

    summary = {
        "generated_at": now_str(),
        "total_samples": total,
        "evaluated_samples": len(evaluated_rows),
        "valid_samples": len(valid_rows),
        "error_samples": len(error_rows),
        "missing_samples": len(missing_rows),
        "coverage_rate": round(len(evaluated_rows) / total * 100, 2) if total else 0.0,
        "accuracy_among_evaluated": round(len(valid_rows) / len(evaluated_rows) * 100, 2) if evaluated_rows else 0.0,
        "strict_accuracy_on_total": round(len(valid_rows) / total * 100, 2) if total else 0.0,
        "threshold_rule": "Tukey lower fence on positive latest match_score values",
        "threshold_info": threshold_info,
        "score_avg": round(mean(scores), 2) if scores else 0.0,
        "score_median": round(median(scores), 2) if scores else 0.0,
        "score_max": round(max(scores), 2) if scores else 0.0,
        "score_min": round(min(scores), 2) if scores else 0.0,
        "source_counter": dict(source_counter),
        "top_valid": sorted(
            (
                {
                    "resume_id": row["resume_id"],
                    "file_name": row["file_name"],
                    "match_score": row["match_score"],
                    "evaluation_source": row["evaluation_source"],
                }
                for row in valid_rows
            ),
            key=lambda item: (-item["match_score"], int(item["resume_id"])),
        )[:10],
        "error_examples": [
            {
                "resume_id": row["resume_id"],
                "file_name": row["file_name"],
                "match_score": row["match_score"],
                "evaluation_source": row["evaluation_source"],
            }
            for row in sorted(error_rows, key=lambda item: (item["match_score"], int(item["resume_id"])))
        ],
    }
    return summary


def write_cleaned_csv(path: Path, rows: List[Dict[str, Any]]) -> None:
    fieldnames = [
        "resume_id",
        "file_name",
        "resume_status",
        "evaluation_id",
        "latest_eval_at",
        "parsed_name",
        "match_score",
        "match_level",
        "report_recommendation",
        "evaluation_source",
        "is_evaluated",
        "is_obvious_error",
        "accuracy_label",
        "analysis_note",
    ]
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8-sig", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        for row in rows:
            writer.writerow({key: row.get(key, "") for key in fieldnames})


def write_markdown(path: Path, summary: Dict[str, Any]) -> None:
    threshold = summary["threshold_info"]["threshold"]
    lines = [
        "# 实验数据准确率分析",
        "",
        f"- 生成时间：{summary['generated_at']}",
        f"- 样本总数：{summary['total_samples']}",
        f"- 已有最新评估结果：{summary['evaluated_samples']}",
        f"- 明显异常低分错误数据：{summary['error_samples']}",
        f"- 有效结果数：{summary['valid_samples']}",
        f"- 缺失结果数：{summary['missing_samples']}",
        f"- 覆盖率：{summary['coverage_rate']}%",
        f"- 已评估样本准确率：{summary['accuracy_among_evaluated']}%",
        f"- 以100份总样本计的严格准确率：{summary['strict_accuracy_on_total']}%",
        "",
        "## 阈值规则",
        "",
        f"- 规则：{summary['threshold_rule']}",
        f"- Q1：{summary['threshold_info']['q1']}",
        f"- 中位数：{summary['threshold_info']['median']}",
        f"- Q3：{summary['threshold_info']['q3']}",
        f"- IQR：{summary['threshold_info']['iqr']}",
        f"- 明显异常低分阈值：`match_score <= {threshold}`",
        "",
        "## 结论",
        "",
        f"- 本轮已返回结果主要集中在 40-50 分区间，低于 {threshold} 分的记录可以视为明显异常数据。",
        f"- 以“异常低分算错、其余算对”为口径，当前已评估样本准确率为 {summary['accuracy_among_evaluated']}%。",
        f"- 如果按完整 100 份样本严格计算，目前有效率为 {summary['strict_accuracy_on_total']}%，主要受未补齐评估结果影响。",
        "",
        "## 异常样本",
        "",
        "| resume_id | 文件名 | 分数 | 来源 |",
        "| --- | --- | --- | --- |",
    ]

    for item in summary["error_examples"]:
        lines.append(
            f"| {item['resume_id']} | {item['file_name']} | {item['match_score']} | {item['evaluation_source']} |"
        )

    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser(description="分析实验样本最新评估结果的准确率，并标记明显异常低分")
    parser.add_argument(
        "--input",
        default="makeREsume/FairCV/eval_results/backend_100/cohort_latest_eval_status.csv",
        help="最新评估状态 CSV",
    )
    parser.add_argument(
        "--output-dir",
        default="makeREsume/FairCV/eval_results/backend_100",
        help="分析结果输出目录",
    )
    parser.add_argument(
        "--threshold",
        type=float,
        default=-1,
        help="手动指定异常低分阈值；默认按 Tukey lower fence 自动计算",
    )
    args = parser.parse_args()

    input_path = Path(args.input).resolve()
    output_dir = Path(args.output_dir).resolve()

    rows = load_rows(input_path)
    evaluated_scores = [to_float(row.get("match_score")) for row in rows if row.get("evaluation_source") != "missing"]
    threshold_info = compute_threshold(evaluated_scores)
    if args.threshold >= 0:
        threshold_info["threshold"] = round(args.threshold, 2)

    cleaned_rows = classify_rows(rows, threshold_info["threshold"])
    summary = build_summary(cleaned_rows, threshold_info)

    cleaned_csv_path = output_dir / "cohort_accuracy_cleaned.csv"
    summary_json_path = output_dir / "cohort_accuracy_summary.json"
    report_md_path = output_dir / "cohort_accuracy_report.md"

    write_cleaned_csv(cleaned_csv_path, cleaned_rows)
    summary_json_path.write_text(json.dumps(summary, ensure_ascii=False, indent=2), encoding="utf-8")
    write_markdown(report_md_path, summary)

    print(f"[done] cleaned csv: {cleaned_csv_path}")
    print(f"[done] summary json: {summary_json_path}")
    print(f"[done] report md: {report_md_path}")


if __name__ == "__main__":
    main()
