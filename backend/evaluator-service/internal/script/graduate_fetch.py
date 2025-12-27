#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
从毕业设计后台获取简历数据
替代 wintalent_fetch.py，用于从本地毕业设计项目的 resume-service 获取简历
"""

from __future__ import annotations
import os
import json
import base64
import argparse
from typing import Dict, List, Optional, Any

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry


class GraduateFetcher:
    """从毕业设计后台获取简历的类"""
    
    def __init__(
        self,
        base_url: str = None,
        page_no: Optional[int] = None,
        page_size: Optional[int] = None,
    ):
        # 从环境变量或默认值获取后台地址
        self.base_url = base_url or os.getenv("GRADUATE_API_URL", "http://localhost:8084")
        self.page_no = page_no or 1
        self.page_size = page_size or 100
        self.session = self._build_session()
    
    def _build_session(self) -> requests.Session:
        """创建带重试机制的 session"""
        s = requests.Session()
        retries = Retry(
            total=3,
            connect=3,
            read=3,
            backoff_factor=0.5,
            status_forcelist=(429, 500, 502, 503, 504),
            allowed_methods=("GET", "POST"),
            raise_on_status=False,
        )
        adapter = HTTPAdapter(max_retries=retries)
        s.mount("https://", adapter)
        s.mount("http://", adapter)
        return s
    
    def _fetch_resumes(self) -> Dict[str, Any]:
        """从后台获取简历列表（用于评估）"""
        url = f"{self.base_url}/api/v1/resumes/evaluation"
        params = {
            "page": self.page_no,
            "page_size": self.page_size,
            "status": "pending",  # 只获取待评估的简历
        }
        
        resp = self.session.get(url, params=params, timeout=30)
        resp.raise_for_status()
        return resp.json()
    
    def _fetch_resume_file(self, resume_id: int) -> Optional[bytes]:
        """下载简历文件"""
        url = f"{self.base_url}/api/v1/resumes/{resume_id}/download"
        try:
            resp = self.session.get(url, timeout=60)
            if resp.status_code == 200:
                return resp.content
        except Exception:
            pass
        return None
    
    def _fetch_job_info(self, job_id: int) -> Optional[Dict[str, Any]]:
        """获取职位信息作为 JD"""
        if not job_id:
            return None
        url = f"{self.base_url}/api/v1/jobs/{job_id}"
        try:
            resp = self.session.get(url, timeout=30)
            if resp.status_code == 200:
                data = resp.json()
                return data.get("data", data)
        except Exception:
            pass
        return None
    
    def fetch_page_payload(self) -> Dict[str, Any]:
        """
        获取单页简历数据，返回格式与 wintalent_fetch.py 一致:
        { total: number, items: [{name, apply_id, resume_id, jd, resume_pdf_b64}] }
        """
        result = self._fetch_resumes()
        
        # 解析返回数据
        data = result.get("data", result)
        resumes = data.get("resumes", [])
        total = data.get("total", len(resumes))
        
        items: List[Dict[str, Any]] = []
        
        for resume in resumes:
            resume_id = resume.get("id")
            talent_id = resume.get("talent_id")
            job_id = resume.get("job_id")
            file_name = resume.get("file_name", "")
            
            # 从文件名提取姓名（假设格式为 "姓名_简历.pdf" 或直接使用 talent_id）
            name = file_name.split("_")[0] if "_" in file_name else f"候选人_{talent_id}"
            
            # 直接使用 API 返回的 base64 编码文件内容
            pdf_b64 = resume.get("file_base64")
            
            # 获取职位信息作为 JD
            jd = self._fetch_job_info(job_id) if job_id else None
            
            items.append({
                "name": name,
                "apply_id": str(talent_id) if talent_id else "",
                "resume_id": str(resume_id) if resume_id else "",
                "post_id": str(job_id) if job_id else "",
                "post_name": jd.get("title", "") if jd else "",
                "recruit_type": "",
                "jd": jd,
                "resume_pdf_b64": pdf_b64,
            })
        
        return {"total": total, "items": items}
    
    def fetch_all(self) -> List[Dict[str, Any]]:
        """获取所有简历（自动分页）"""
        all_items: List[Dict[str, Any]] = []
        page = 1
        max_pages = 100
        
        while page <= max_pages:
            self.page_no = page
            payload = self.fetch_page_payload()
            items = payload.get("items", [])
            
            if not items:
                break
            
            all_items.extend(items)
            
            total = payload.get("total", 0)
            if len(all_items) >= total:
                break
            
            page += 1
        
        return all_items


def main():
    parser = argparse.ArgumentParser(description="从毕业设计后台获取简历")
    parser.add_argument("--json-out", action="store_true", help="输出 JSON 格式")
    parser.add_argument("--page-no", type=int, default=1, help="页码")
    parser.add_argument("--page-size", type=int, default=100, help="每页数量")
    parser.add_argument("--base-url", type=str, help="后台 API 地址")
    
    args = parser.parse_args()
    
    fetcher = GraduateFetcher(
        base_url=args.base_url,
        page_no=args.page_no,
        page_size=args.page_size,
    )
    
    if args.json_out:
        # 输出 JSON 格式（与 wintalent_fetch.py 兼容）
        payload = fetcher.fetch_page_payload()
        print(json.dumps(payload, ensure_ascii=False))
    else:
        # 默认输出简要信息
        items = fetcher.fetch_all()
        print(f"获取到 {len(items)} 份简历")
        for item in items:
            print(f"  - {item['name']} (ID: {item['resume_id']})")


if __name__ == "__main__":
    main()
