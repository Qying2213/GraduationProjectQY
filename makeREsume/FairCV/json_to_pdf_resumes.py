#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
将 FairCV 的 JSON 简历批量导出为中文 PDF。

特点：
1. 不依赖第三方 PDF 库，直接生成 PDF
2. 支持大文件流式读取（如 5GB+ 的 resumes.json）
3. 支持按数量限制导出，避免一次性生成过多文件

示例：
python3 json_to_pdf_resumes.py \
  --input data/resumes.json \
  --output-dir output_pdfs \
  --limit 100

如需导出全部：
python3 json_to_pdf_resumes.py \
  --input data/resumes.json \
  --output-dir output_pdfs \
  --all
"""

from __future__ import annotations

import argparse
import csv
import json
import math
import re
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, Generator, Iterable, List, Optional, Tuple


PAGE_WIDTH = 595
PAGE_HEIGHT = 842
MARGIN_X = 50
TOP_Y = 790
BOTTOM_Y = 60
BODY_FONT_SIZE = 11
SECTION_FONT_SIZE = 13
TITLE_FONT_SIZE = 18
LINE_GAP = 8


@dataclass
class StyledLine:
    text: str
    font_size: int
    gap_after: int = 4


class PDFWriter:
    def __init__(self) -> None:
        self.objects: List[Optional[bytes]] = [None]

    def reserve(self) -> int:
        self.objects.append(None)
        return len(self.objects) - 1

    def set_object(self, obj_id: int, content: str | bytes) -> None:
        if isinstance(content, str):
            content = content.encode("latin-1")
        self.objects[obj_id] = content

    def write(self, path: Path, root_id: int, info_id: Optional[int] = None) -> None:
        out = bytearray(b"%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")
        offsets = [0]

        for obj_id in range(1, len(self.objects)):
            obj = self.objects[obj_id]
            if obj is None:
                raise ValueError(f"PDF object {obj_id} is not set")
            offsets.append(len(out))
            out.extend(f"{obj_id} 0 obj\n".encode("ascii"))
            out.extend(obj)
            if not obj.endswith(b"\n"):
                out.extend(b"\n")
            out.extend(b"endobj\n")

        xref_start = len(out)
        out.extend(f"xref\n0 {len(self.objects)}\n".encode("ascii"))
        out.extend(b"0000000000 65535 f \n")
        for offset in offsets[1:]:
            out.extend(f"{offset:010d} 00000 n \n".encode("ascii"))

        trailer = f"<< /Size {len(self.objects)} /Root {root_id} 0 R"
        if info_id is not None:
            trailer += f" /Info {info_id} 0 R"
        trailer += " >>"

        out.extend(b"trailer\n")
        out.extend(trailer.encode("ascii"))
        out.extend(b"\nstartxref\n")
        out.extend(str(xref_start).encode("ascii"))
        out.extend(b"\n%%EOF")

        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(out)


def sanitize_filename(name: str, max_length: int = 80) -> str:
    name = re.sub(r'[\\/:*?"<>|]+', "_", name).strip()
    name = re.sub(r"\s+", "_", name)
    if not name:
        name = "resume"
    return name[:max_length]


def clean_markdown_line(line: str) -> Tuple[str, str]:
    raw = line.strip()
    if not raw:
        return "blank", ""
    if raw == "---":
        return "divider", ""
    if raw.startswith("### "):
        return "section", raw[4:].replace("**", "").strip()
    if raw.startswith("## "):
        return "section", raw[3:].replace("**", "").strip()
    if raw.startswith("# "):
        return "section", raw[2:].replace("**", "").strip()

    cleaned = raw.replace("**", "")
    cleaned = re.sub(r"^\-\s*", "• ", cleaned)
    return "body", cleaned


def is_cjk(char: str) -> bool:
    code = ord(char)
    return (
        0x4E00 <= code <= 0x9FFF
        or 0x3400 <= code <= 0x4DBF
        or 0x3000 <= code <= 0x303F
        or 0xFF00 <= code <= 0xFFEF
    )


def estimate_char_width(char: str, font_size: int) -> float:
    if char == "\t":
        return font_size * 2
    if is_cjk(char):
        return float(font_size)
    if char.isspace():
        return font_size * 0.35
    if char.isascii():
        return font_size * 0.56
    return font_size * 0.8


def wrap_text(text: str, font_size: int, max_width: float) -> List[str]:
    if not text:
        return [""]

    lines: List[str] = []
    current = []
    current_width = 0.0

    for char in text:
        char_width = estimate_char_width(char, font_size)
        if current and current_width + char_width > max_width:
            lines.append("".join(current).rstrip())
            current = [char]
            current_width = char_width
        else:
            current.append(char)
            current_width += char_width

    if current:
        lines.append("".join(current).rstrip())

    return lines or [""]


def build_resume_lines(resume: Dict[str, object]) -> List[StyledLine]:
    metadata = resume.get("metadata", {}) if isinstance(resume, dict) else {}
    content = str(resume.get("content", "") if isinstance(resume, dict) else "")
    position = str(metadata.get("position", "")).strip()

    lines: List[StyledLine] = []
    if position:
        lines.append(StyledLine(text=f"{position}简历", font_size=TITLE_FONT_SIZE, gap_after=12))

    for raw_line in content.splitlines():
        style, cleaned = clean_markdown_line(raw_line)
        if style == "blank":
            lines.append(StyledLine(text="", font_size=BODY_FONT_SIZE, gap_after=8))
        elif style == "divider":
            lines.append(StyledLine(text="—" * 28, font_size=BODY_FONT_SIZE, gap_after=8))
        elif style == "section":
            lines.append(StyledLine(text=cleaned, font_size=SECTION_FONT_SIZE, gap_after=8))
        else:
            lines.append(StyledLine(text=cleaned, font_size=BODY_FONT_SIZE, gap_after=4))

    return lines


def split_lines_into_pages(lines: Iterable[StyledLine]) -> List[List[StyledLine]]:
    max_width = PAGE_WIDTH - 2 * MARGIN_X
    pages: List[List[StyledLine]] = []
    current_page: List[StyledLine] = []
    y = TOP_Y

    for line in lines:
        wrapped = wrap_text(line.text, line.font_size, max_width)
        required_height = len(wrapped) * (line.font_size + 2) + line.gap_after + LINE_GAP
        if current_page and y - required_height < BOTTOM_Y:
            pages.append(current_page)
            current_page = []
            y = TOP_Y

        for idx, wrapped_line in enumerate(wrapped):
            gap_after = line.gap_after if idx == len(wrapped) - 1 else 2
            current_page.append(
                StyledLine(text=wrapped_line, font_size=line.font_size, gap_after=gap_after)
            )
            y -= line.font_size + 2 + gap_after

    if current_page:
        pages.append(current_page)

    return pages or [[StyledLine(text="空白简历", font_size=TITLE_FONT_SIZE)]]


def pdf_text_hex(text: str) -> str:
    return text.encode("utf-16-be").hex().upper()


def build_page_stream(lines: List[StyledLine]) -> bytes:
    parts = ["BT"]
    y = TOP_Y

    for line in lines:
        parts.append(f"/F1 {line.font_size} Tf")
        parts.append(f"1 0 0 1 {MARGIN_X} {y} Tm")
        if line.text:
            parts.append(f"<{pdf_text_hex(line.text)}> Tj")
        y -= line.font_size + 2 + line.gap_after

    parts.append("ET")
    content = "\n".join(parts).encode("ascii")
    return b"<< /Length " + str(len(content)).encode("ascii") + b" >>\nstream\n" + content + b"\nendstream"


def build_pdf(resume: Dict[str, object], output_path: Path) -> int:
    writer = PDFWriter()

    catalog_id = writer.reserve()
    pages_id = writer.reserve()
    font_id = writer.reserve()
    cidfont_id = writer.reserve()
    fontdesc_id = writer.reserve()
    info_id = writer.reserve()

    writer.set_object(
        fontdesc_id,
        "<< /Type /FontDescriptor /FontName /STSong-Light /Flags 4 /ItalicAngle 0 "
        "/Ascent 752 /Descent -271 /CapHeight 737 /StemV 80 /MissingWidth 500 >>",
    )
    writer.set_object(
        cidfont_id,
        f"<< /Type /Font /Subtype /CIDFontType0 /BaseFont /STSong-Light "
        f"/CIDSystemInfo << /Registry (Adobe) /Ordering (GB1) /Supplement 4 >> "
        f"/FontDescriptor {fontdesc_id} 0 R /DW 1000 >>",
    )
    writer.set_object(
        font_id,
        f"<< /Type /Font /Subtype /Type0 /BaseFont /STSong-Light "
        f"/Encoding /UniGB-UCS2-H /DescendantFonts [{cidfont_id} 0 R] >>",
    )

    lines = build_resume_lines(resume)
    page_lines = split_lines_into_pages(lines)

    page_ids: List[int] = []
    for page in page_lines:
        page_id = writer.reserve()
        content_id = writer.reserve()
        writer.set_object(content_id, build_page_stream(page))
        writer.set_object(
            page_id,
            f"<< /Type /Page /Parent {pages_id} 0 R "
            f"/MediaBox [0 0 {PAGE_WIDTH} {PAGE_HEIGHT}] "
            f"/Resources << /Font << /F1 {font_id} 0 R >> >> "
            f"/Contents {content_id} 0 R >>",
        )
        page_ids.append(page_id)

    kids = " ".join(f"{page_id} 0 R" for page_id in page_ids)
    writer.set_object(pages_id, f"<< /Type /Pages /Count {len(page_ids)} /Kids [{kids}] >>")
    writer.set_object(catalog_id, f"<< /Type /Catalog /Pages {pages_id} 0 R >>")
    writer.set_object(info_id, "<< /Producer (FairCV JSON to PDF) /Title (Chinese Resume Export) >>")

    writer.write(output_path, root_id=catalog_id, info_id=info_id)
    return len(page_ids)


def extract_resume_name(content: str) -> str:
    cleaned = content.replace("**", "")
    match = re.search(r"姓名[：:]\s*([^\n\r]+)", cleaned)
    if match:
        return match.group(1).strip().split()[0]
    return "未命名"


def iter_json_array(path: Path, chunk_size: int = 1024 * 1024) -> Generator[Dict[str, object], None, None]:
    decoder = json.JSONDecoder()
    with path.open("r", encoding="utf-8") as f:
        buffer = ""
        in_array = False
        eof = False

        while True:
            if not eof and len(buffer) < chunk_size:
                chunk = f.read(chunk_size)
                if chunk:
                    buffer += chunk
                else:
                    eof = True

            buffer = buffer.lstrip()
            if not in_array:
                if not buffer:
                    if eof:
                        return
                    continue
                if buffer[0] != "[":
                    raise ValueError("resumes.json 必须是 JSON 数组格式")
                buffer = buffer[1:]
                in_array = True
                continue

            buffer = buffer.lstrip()
            if not buffer:
                if eof:
                    return
                continue

            if buffer[0] == "]":
                return
            if buffer[0] == ",":
                buffer = buffer[1:]
                continue

            try:
                obj, index = decoder.raw_decode(buffer)
            except json.JSONDecodeError:
                if eof:
                    raise
                continue

            buffer = buffer[index:]
            if isinstance(obj, dict):
                yield obj


def iter_resumes(path: Path) -> Generator[Dict[str, object], None, None]:
    with path.open("r", encoding="utf-8") as f:
        first_char = ""
        while True:
            ch = f.read(1)
            if not ch:
                break
            if not ch.isspace():
                first_char = ch
                break

    if first_char == "[":
        yield from iter_json_array(path)
        return

    with path.open("r", encoding="utf-8") as f:
        data = json.load(f)
    if isinstance(data, dict) and isinstance(data.get("resumes"), list):
        for item in data["resumes"]:
            if isinstance(item, dict):
                yield item
    elif isinstance(data, list):
        for item in data:
            if isinstance(item, dict):
                yield item
    else:
        raise ValueError("不支持的 JSON 结构")


def resume_matches_position(resume: Dict[str, object], position_contains: Optional[str]) -> bool:
    if not position_contains:
        return True
    metadata = resume.get("metadata", {}) if isinstance(resume, dict) else {}
    position = str(metadata.get("position", "")).strip()
    return position_contains.lower() in position.lower()


def export_resumes_to_pdf(
    input_path: Path,
    output_dir: Path,
    limit: Optional[int],
    start_index: int,
    position_contains: Optional[str] = None,
) -> Tuple[int, Path]:
    output_dir.mkdir(parents=True, exist_ok=True)
    manifest_path = output_dir / "manifest.csv"

    exported = 0
    with manifest_path.open("w", encoding="utf-8", newline="") as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(["index", "filename", "name", "position", "pages"])

        for idx, resume in enumerate(iter_resumes(input_path), start=1):
            if idx < start_index:
                continue
            if not resume_matches_position(resume, position_contains):
                continue
            if limit is not None and exported >= limit:
                break

            metadata = resume.get("metadata", {}) if isinstance(resume, dict) else {}
            content = str(resume.get("content", "") if isinstance(resume, dict) else "")
            position = str(metadata.get("position", "未知岗位")).strip() or "未知岗位"
            name = extract_resume_name(content)

            position_dir = output_dir / sanitize_filename(position, max_length=40)
            filename = f"{idx:06d}_{sanitize_filename(position, 30)}_{sanitize_filename(name, 20)}.pdf"
            pdf_path = position_dir / filename

            pages = build_pdf(resume, pdf_path)
            writer.writerow([idx, str(pdf_path.relative_to(output_dir)), name, position, pages])
            exported += 1

            if exported % 10 == 0:
                print(f"[progress] 已导出 {exported} 份 PDF，最近文件：{pdf_path.name}")

    return exported, manifest_path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="批量将 FairCV JSON 简历导出为中文 PDF")
    parser.add_argument("--input", type=Path, default=Path("data/resumes.json"), help="输入 JSON 文件")
    parser.add_argument("--output-dir", type=Path, default=Path("output_pdfs"), help="PDF 输出目录")
    parser.add_argument("--limit", type=int, default=100, help="最多导出多少份，默认 100")
    parser.add_argument("--start-index", type=int, default=1, help="从第几份简历开始导出，默认 1")
    parser.add_argument("--all", action="store_true", help="导出全部简历（请谨慎使用）")
    parser.add_argument("--position-contains", type=str, default=None, help="仅导出岗位名称包含该关键词的简历，例如 后端开发工程师")
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    input_path = args.input
    if not input_path.is_absolute():
        input_path = (Path.cwd() / input_path).resolve()

    output_dir = args.output_dir
    if not output_dir.is_absolute():
        output_dir = (Path.cwd() / output_dir).resolve()

    if not input_path.exists():
        raise FileNotFoundError(f"输入文件不存在: {input_path}")

    limit = None if args.all else args.limit

    print(f"[info] 输入文件: {input_path}")
    print(f"[info] 输出目录: {output_dir}")
    print(f"[info] 起始索引: {args.start_index}")
    print(f"[info] 导出数量: {'全部' if limit is None else limit}")
    print(f"[info] 岗位筛选: {args.position_contains or '不过滤'}")

    exported, manifest_path = export_resumes_to_pdf(
        input_path=input_path,
        output_dir=output_dir,
        limit=limit,
        start_index=max(1, args.start_index),
        position_contains=args.position_contains,
    )

    print(f"[done] 已导出 {exported} 份 PDF")
    print(f"[done] 清单文件: {manifest_path}")


if __name__ == "__main__":
    main()
