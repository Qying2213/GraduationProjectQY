#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from __future__ import annotations
import os
import json
import re
import binascii
import base64
from pathlib import Path
from typing import Dict, List, Optional, Any

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry
from requests_toolbelt.multipart.encoder import MultipartEncoder
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad


class WintalentFetcher:
    BASE = "https://www.wintalent.cn"
    CREATE_TOKEN_URL = f"{BASE}/interviewer/common/createToken"
    LOGIN_URL = f"{BASE}/interviewer/login"
    INFO_URL = f"{BASE}/interviewer/loginInfo"
    
    RECOMMEND_LIST_PATH = "/interviewer/interviewPlatform/screen/recommendToMe"
    RESUME_ORIGINAL_PATH = "/interviewPlatform/getResumeOriginalInfo"
    SHOW_POST_JD_PATH = "/interviewer/common/data/showPostJD"
    REQUEST_TIMEOUT = 30
    REQUEST_TIMEOUT_PDF = 60
    
    def __init__(
        self,
        corp_code: str,
        username: Optional[str],
        password: Optional[str],
        current_code: str = "0/430400/430447/430448/430449",
        page_no: Optional[int] = None,
        page_size: Optional[int] = None,
        output_dir: str = "out",
        session_cookie: Optional[str] = None,
        recommend_begin_time_str: Optional[str] = None,
        recommend_end_time_str: Optional[str] = None,
    ):
        self.corp_code = corp_code
        self.current_code = current_code
        self.page_no = page_no
        self.page_size = page_size
        self.output_dir = Path(output_dir)
        self.session_cookie = session_cookie
        self.recommend_begin_time_str = recommend_begin_time_str
        self.recommend_end_time_str = recommend_end_time_str
        self.username = username or os.getenv("WT_USERNAME")
        self.password = password or os.getenv("WT_PASSWORD")
        self.session = self._build_session()
        self._ensure_logged_in()
    
    def _build_session(self) -> requests.Session:
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
    
    def _get_login_info(self, corp_code: str) -> Optional[Dict]:
        params = {"corpCode": corp_code}
        url = f"{self.INFO_URL}?r=0.123456789"
        resp = self.session.post(url, data=params, headers={
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "X-Requested-With": "XMLHttpRequest",
            "Origin": self.BASE,
            "Referer": f"{self.BASE}/interviewer/login"
        })
        if resp.status_code == 200:
            try:
                return resp.json()
            except:
                pass
        return None
    
    def _aes_encrypt(self, text: str, key_string: str) -> str:
        if not key_string:
            raise ValueError("Encryption key is empty")
        key = key_string.encode("utf-8")
        data = text.encode("utf-8")
        cipher = AES.new(key, AES.MODE_ECB)
        encrypted_bytes = cipher.encrypt(pad(data, AES.block_size))
        return binascii.hexlify(encrypted_bytes).decode("utf-8")
    
    def _login(self, username: str, password: str) -> Dict:
        info = self._get_login_info(self.corp_code)
        if not info:
            raise ValueError("Failed to get login info")
        
        time_sign = info.get("time")
        pass_key = info.get("interviewerPassKey")
        pass_type = info.get("interviewerPassType")
        
        if str(pass_type) == "3":
            raise ValueError("Server requires SM4 encryption. Script currently only supports AES.")
        
        encrypted_password = self._aes_encrypt(password, pass_key)
        
        payload = {
            "remremberMeInput": "true",
            "timeSign": time_sign,
            "actycoStatus": "",
            "loginMode": "",
            "localeType": "1",
            "login-select": "1",
            "corpCode": self.corp_code,
            "userName": username,
            "password": encrypted_password,
            "verifyCode": "",
            "phoneNum": "",
            "verification": "",
            "remremberMe": "true"
        }
        
        login_resp = self.session.post(self.LOGIN_URL, data=payload, headers={
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "X-Requested-With": "XMLHttpRequest",
            "Origin": self.BASE,
            "Referer": f"{self.BASE}/interviewer/login"
        })
        
        login_resp.raise_for_status()
        return self.session.cookies.get_dict()
    
    def _ensure_logged_in(self) -> None:
        if self.session_cookie:
            self.session.cookies.set("SESSION", self.session_cookie, domain="www.wintalent.cn", path="/")
        elif self.username and self.password:
            cookies = self._login(self.username, self.password)
            if not cookies or "SESSION" not in cookies:
                raise ValueError("登录失败，未返回 SESSION cookie")
            self.session.cookies.update(cookies)
        else:
            env_sess = os.getenv("WT_SESSION") or os.getenv("SESSION")
            if env_sess:
                self.session.cookies.set("SESSION", env_sess, domain="www.wintalent.cn", path="/")
            else:
                raise ValueError("未提供登录凭据（需要提供 username 和 password，或 session_cookie，或设置 WT_SESSION 环境变量）")
        
        session_cookies = [c for c in self.session.cookies if c.name == "SESSION"]
        if not session_cookies:
            raise ValueError("SESSION cookie 缺失")
    
    def _fetch_token_for_endpoint(self, endpoint: str) -> Optional[str]:
        headers = {
            "accept": "application/json, text/javascript, */*; q=0.01",
            "accept-language": "zh-CN,zh;q=0.9",
            "content-type": "application/x-www-form-urlencoded",
            "origin": self.BASE,
            "priority": "u=1, i",
            "referer": f"{self.BASE}/interviewer/interviewPlatform/newpc/antPage/screen.html?currentCode={self.current_code}",
            "sec-ch-ua": '"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"',
            "sec-ch-ua-mobile": "?0",
            "sec-ch-ua-platform": '"macOS"',
            "sec-fetch-dest": "empty",
            "sec-fetch-mode": "cors",
            "sec-fetch-site": "same-origin",
            "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
            "x-requested-with": "XMLHttpRequest",
        }
        
        r = self.session.post(
            self.CREATE_TOKEN_URL,
            headers=headers,
            data={"url": endpoint},
            timeout=self.REQUEST_TIMEOUT
        )
        r.raise_for_status()
        
        data = r.json()
        token_url = data.get("tokenUrl")
        if token_url:
            if not token_url.startswith("http"):
                return f"{self.BASE}{token_url}" if token_url.startswith("/") else f"{self.BASE}/{token_url}"
            return token_url
        return None
    
    def _fetch_full_urls(self) -> Dict[str, str]:
        endpoints = [self.RECOMMEND_LIST_PATH, self.RESUME_ORIGINAL_PATH, self.SHOW_POST_JD_PATH]
        resolved = {}
        for endpoint in endpoints:
            full_url = self._fetch_token_for_endpoint(endpoint)
            if full_url:
                resolved[endpoint] = full_url
        return resolved
    
    def _recommend_to_me(self, recommend_full_url: str) -> Dict:
        page = str(self.page_no) if self.page_no is not None else "1"
        size = str(self.page_size) if self.page_size is not None else "100"

        fields = {
            "isSearchButton": "true",
            "currentCode": self.current_code,
            "currentPage": page,
            "rowSize": size,
        }
        if self.recommend_begin_time_str:
            fields["recommendBeginTimeStr"] = self.recommend_begin_time_str
        if self.recommend_end_time_str:
            fields["recommendEndTimeStr"] = self.recommend_end_time_str

        mp = MultipartEncoder(fields=fields)
        headers = {
            "Accept": "application/json, text/javascript, */*; q=0.01",
            "Origin": self.BASE,
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
            "Accept-Language": "zh-CN,zh;q=0.9",
            "Sec-Fetch-Dest": "empty",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Site": "same-origin",
            "X-Requested-With": "XMLHttpRequest",
            "Referer": f"{self.BASE}/interviewer/interviewPlatform/newpc/antPage/screen.html",
            "Content-Type": mp.content_type,
        }
        r = self.session.post(recommend_full_url, headers=headers, data=mp, timeout=self.REQUEST_TIMEOUT)
        if r.status_code in (401, 403):
            raise PermissionError(f"列表接口未授权: HTTP {r.status_code}")
        r.raise_for_status()
        return r.json()
    
    def _fetch_resume_pdf(self, resume_original_full_url: str, apply_id: str | int, resume_id: str | int) -> bytes:
        connector = "&" if "?" in resume_original_full_url else "?"
        url = f"{resume_original_full_url}{connector}lanType=1&applyId={apply_id}&resumeId={resume_id}&showPdf=true&fileType=pdf"
        headers = {
            "Accept": "application/json, text/javascript, */*; q=0.01",
            "Origin": self.BASE,
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
            "Accept-Language": "zh-CN,zh;q=0.9",
            "Sec-Fetch-Dest": "empty",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Site": "same-origin",
            "X-Requested-With": "XMLHttpRequest",
            "Referer": f"{self.BASE}/interviewer/interviewPlatform/newpc/jsp/resume/resumeInfo.html",
        }
        r = self.session.get(url, headers=headers, timeout=self.REQUEST_TIMEOUT_PDF)
        if r.status_code in (401, 403):
            raise PermissionError(f"原始简历接口未授权: HTTP {r.status_code}")
        r.raise_for_status()
        return r.content
    
    def _fetch_show_post_jd(self, show_post_jd_full_url: str, post_id: str = "100703", recruit_type: str = "2") -> Dict:
        """调用 SHOW_POST_JD 接口，返回原样 JSON。"""
        fields = {
            "postId": str(post_id) if post_id else "100703",
            "recruitType": str(recruit_type) if recruit_type else "2",
        }
        mp = MultipartEncoder(fields=fields)
        headers = {
            "Accept": "application/json, text/javascript, */*; q=0.01",
            "Origin": self.BASE,
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36",
            "Accept-Language": "zh-CN,zh;q=0.9",
            "Sec-Fetch-Dest": "empty",
            "Sec-Fetch-Mode": "cors",
            "Sec-Fetch-Site": "same-origin",
            "X-Requested-With": "XMLHttpRequest",
            "Referer": f"{self.BASE}/interviewer/interviewPlatform/newpc/antPage/screen.html",
            "Content-Type": mp.content_type,
        }
        r = self.session.post(show_post_jd_full_url, headers=headers, data=mp, timeout=self.REQUEST_TIMEOUT)
        if r.status_code in (401, 403):
            raise PermissionError(f"岗位JD接口未授权: HTTP {r.status_code}")
        r.raise_for_status()
        return r.json()
    
    @staticmethod
    def _sanitize_filename(name: Optional[str]) -> str:
        if not name:
            return "unknown"
        return re.sub(r"[^\w\-.]+", "_", name)
    
    def _extract_candidates(self, list_json: Dict) -> List[Dict]:
        items = list_json.get("rowList") or []
        results = []
        for it in items:
            results.append({
                "postId": it.get("postId"),  # 岗位ID
                "postName": it.get("enPostName") or it.get("cnPostName") or it.get("postName"),
                "recruitType": it.get("recruitType"),  # 招聘类型
                "name": it.get("name"),
                "applyId": it.get("applyId"),
                "resumeId": it.get("resumeId"),
            })
        return results
    
    def _save_pdf(self, candidate: Dict, content: bytes) -> Path:
        self.output_dir.mkdir(parents=True, exist_ok=True)
        name = candidate.get("name") or "unknown"
        fname = f"{self._sanitize_filename(name)}.pdf"
        path = self.output_dir / fname
        with path.open("wb") as f:
            f.write(content)
        return path
    
    def fetch_all(self) -> List[Dict[str, Any]]:
        """抓取候选人列表，并为每条记录返回 name、岗位JD原样JSON、简历PDF二进制。
        - 不落盘保存PDF；
        - 对单条失败的接口返回 None；
        - 保持现有分页逻辑。
        """
        full_urls = self._fetch_full_urls()
        recommend_full = full_urls.get(self.RECOMMEND_LIST_PATH)
        resume_original_full = full_urls.get(self.RESUME_ORIGINAL_PATH)
        show_post_jd_full = full_urls.get(self.SHOW_POST_JD_PATH)

        if not recommend_full:
            raise ValueError("未获取到推荐列表接口完整 URL")

        results: List[Dict[str, Any]] = []
        # 自动分页：从第 1 页开始，基于 rowCount 判定终止
        page = 1
        page_size = int(self.page_size or 100)
        total = None  # type: Optional[int]
        fetched = 0
        max_pages = 1000  # 安全上限，避免异常情况下无限循环

        while page <= max_pages:
            self.page_no = page
            self.page_size = page_size
            try:
                list_json = self._recommend_to_me(recommend_full)
            except PermissionError:
                recommend_full = self._fetch_token_for_endpoint(self.RECOMMEND_LIST_PATH)
                if not recommend_full:
                    raise ValueError("无法获取推荐列表接口 URL")
                list_json = self._recommend_to_me(recommend_full)

            # 读取总数
            if total is None:
                try:
                    total = int(list_json.get("rowCount") or 0)
                except Exception:
                    total = 0

            candidates = self._extract_candidates(list_json)
            if not candidates:
                break

            for c in candidates:
                apply_id, resume_id = c.get("applyId"), c.get("resumeId")
                post_id, recruit_type = c.get("postId"), c.get("recruitType")
                post_name = c.get("postName")
                name = c.get("name")
                jd: Optional[Dict] = None
                pdf_bytes: Optional[bytes] = None

                # 简历 PDF
                if apply_id and resume_id:
                    try:
                        if not resume_original_full:
                            resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH)
                        if resume_original_full:
                            pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                    except PermissionError:
                        resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH) or resume_original_full
                        try:
                            if resume_original_full:
                                pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                        except requests.RequestException:
                            pdf_bytes = None
                    except requests.RequestException:
                        pdf_bytes = None

                # 岗位 JD（使用候选人的 postId 和 recruitType）
                try:
                    if not show_post_jd_full:
                        show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH)
                    if show_post_jd_full:
                        jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
                except PermissionError:
                    show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH) or show_post_jd_full
                    try:
                        if show_post_jd_full:
                            jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
                    except requests.RequestException:
                        jd = None
                except requests.RequestException:
                    jd = None

                results.append({
                    "name": name,
                    "apply_id": str(apply_id) if apply_id else "",
                    "resume_id": str(resume_id) if resume_id else "",
                    "post_id": str(post_id) if post_id else "",
                    "post_name": post_name or "",
                    "recruit_type": str(recruit_type) if recruit_type else "",
                    "jd": jd,
                    "resume_pdf": pdf_bytes,
                })

            fetched += len(candidates)
            # 如果服务端给了 rowCount，达到总数则停止
            if total and fetched >= total:
                break

            page += 1

        return results

    def fetch_page_payload(self) -> Dict[str, Any]:
        """抓取单页，返回 { total: rowCount, items: [{name, jd, resume_pdf_b64}] }"""
        full_urls = self._fetch_full_urls()
        recommend_full = full_urls.get(self.RECOMMEND_LIST_PATH)
        resume_original_full = full_urls.get(self.RESUME_ORIGINAL_PATH)
        show_post_jd_full = full_urls.get(self.SHOW_POST_JD_PATH)
        if not recommend_full:
            raise ValueError("未获取到推荐列表接口完整 URL")

        # 只抓取当前 page_no 指定的这一页
        try:
            list_json = self._recommend_to_me(recommend_full)
        except PermissionError:
            recommend_full = self._fetch_token_for_endpoint(self.RECOMMEND_LIST_PATH)
            if not recommend_full:
                raise ValueError("无法获取推荐列表接口 URL")
            list_json = self._recommend_to_me(recommend_full)

        total = 0
        try:
            # recommendToMe 返回字段 rowCount
            total = int(list_json.get("rowCount") or 0)
        except Exception:
            total = 0

        candidates = self._extract_candidates(list_json)
        items: List[Dict[str, Any]] = []

        for c in candidates:
            apply_id, resume_id = c.get("applyId"), c.get("resumeId")
            post_id, recruit_type = c.get("postId"), c.get("recruitType")
            post_name = c.get("postName")
            name = c.get("name")
            jd: Optional[Dict] = None
            pdf_bytes: Optional[bytes] = None

            # 简历 PDF
            if apply_id and resume_id:
                try:
                    if not resume_original_full:
                        resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH)
                    if resume_original_full:
                        pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                except PermissionError:
                    resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH) or resume_original_full
                    try:
                        if resume_original_full:
                            pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                    except requests.RequestException:
                        pdf_bytes = None
                except requests.RequestException:
                    pdf_bytes = None

            # 岗位 JD（使用候选人的 postId 和 recruitType）
            try:
                if not show_post_jd_full:
                    show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH)
                if show_post_jd_full:
                    jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
            except PermissionError:
                show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH) or show_post_jd_full
                try:
                    if show_post_jd_full:
                        jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
                except requests.RequestException:
                    jd = None
            except requests.RequestException:
                jd = None

            b64 = base64.b64encode(pdf_bytes).decode("utf-8") if isinstance(pdf_bytes, (bytes, bytearray)) else None
            items.append({
                "name": name,
                "apply_id": str(apply_id) if apply_id else "",
                "resume_id": str(resume_id) if resume_id else "",
                "post_id": str(post_id) if post_id else "",
                "post_name": post_name or "",
                "recruit_type": str(recruit_type) if recruit_type else "",
                "jd": jd,
                "resume_pdf_b64": b64,
            })

        # 如果服务端没给 rowCount，则以当前页条数作为兜底，不影响调用方按 < pageSize 规则停止
        if total <= 0:
            total = len(items)
        return {"total": total, "items": items}

    def fetch_positions(self) -> List[Dict[str, Any]]:
        """抓取所有候选人，提取并去重岗位信息（含 JD）。
        返回格式：[{post_id, post_name, recruit_type, service_condition, work_content}]
        """
        full_urls = self._fetch_full_urls()
        recommend_full = full_urls.get(self.RECOMMEND_LIST_PATH)
        show_post_jd_full = full_urls.get(self.SHOW_POST_JD_PATH)

        if not recommend_full:
            raise ValueError("未获取到推荐列表接口完整 URL")

        # 用于去重的字典，key 为 postId
        positions_map: Dict[str, Dict[str, Any]] = {}
        
        # 自动分页抓取所有候选人
        page = 1
        page_size = int(self.page_size or 100)
        total = None
        fetched = 0
        max_pages = 1000

        while page <= max_pages:
            self.page_no = page
            self.page_size = page_size
            try:
                list_json = self._recommend_to_me(recommend_full)
            except PermissionError:
                recommend_full = self._fetch_token_for_endpoint(self.RECOMMEND_LIST_PATH)
                if not recommend_full:
                    raise ValueError("无法获取推荐列表接口 URL")
                list_json = self._recommend_to_me(recommend_full)

            if total is None:
                try:
                    total = int(list_json.get("rowCount") or 0)
                except Exception:
                    total = 0

            candidates = self._extract_candidates(list_json)
            if not candidates:
                break

            for c in candidates:
                post_id = c.get("postId")
                if not post_id:
                    continue
                post_id_str = str(post_id)
                
                # 如果已经处理过这个岗位，跳过
                if post_id_str in positions_map:
                    continue
                
                post_name = c.get("postName") or ""
                recruit_type = c.get("recruitType")
                
                # 获取岗位 JD
                jd: Optional[Dict] = None
                try:
                    if not show_post_jd_full:
                        show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH)
                    if show_post_jd_full:
                        jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
                except PermissionError:
                    show_post_jd_full = self._fetch_token_for_endpoint(self.SHOW_POST_JD_PATH) or show_post_jd_full
                    try:
                        if show_post_jd_full:
                            jd = self._fetch_show_post_jd(show_post_jd_full, post_id, recruit_type)
                    except requests.RequestException:
                        jd = None
                except requests.RequestException:
                    jd = None
                
                # 提取 JD 内容
                service_condition = ""
                work_content = ""
                if jd:
                    service_condition = jd.get("serviceCondition") or ""
                    work_content = jd.get("workContent") or ""
                
                positions_map[post_id_str] = {
                    "post_id": post_id_str,
                    "post_name": post_name,
                    "recruit_type": str(recruit_type) if recruit_type else "",
                    "service_condition": service_condition,
                    "work_content": work_content,
                }

            fetched += len(candidates)
            if total and fetched >= total:
                break
            page += 1

        return list(positions_map.values())

    def run(self) -> None:
        full_urls = self._fetch_full_urls()
        recommend_full = full_urls.get(self.RECOMMEND_LIST_PATH)
        resume_original_full = full_urls.get(self.RESUME_ORIGINAL_PATH)

        if not recommend_full:
            raise ValueError("未获取到推荐列表接口完整 URL")

        ok = 0
        # 如果显式指定了 page_no，则只拉取这一页；否则自动分页从第 1 页开始直到无数据
        start_page = self.page_no if self.page_no is not None else 1

        page = start_page
        while True:
            self.page_no = page
            try:
                list_json = self._recommend_to_me(recommend_full)
            except PermissionError:
                recommend_full = self._fetch_token_for_endpoint(self.RECOMMEND_LIST_PATH)
                if not recommend_full:
                    raise ValueError("无法获取推荐列表接口 URL")
                list_json = self._recommend_to_me(recommend_full)

            candidates = self._extract_candidates(list_json)
            if not candidates:
                # 当前页已经没有数据，结束循环
                break

            for c in candidates:
                apply_id, resume_id = c.get("applyId"), c.get("resumeId")
                if not apply_id or not resume_id:
                    continue

                try:
                    if not resume_original_full:
                        resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH)
                    if resume_original_full:
                        pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                        self._save_pdf(c, pdf_bytes)
                        ok += 1
                except PermissionError:
                    resume_original_full = self._fetch_token_for_endpoint(self.RESUME_ORIGINAL_PATH) or resume_original_full
                    if resume_original_full:
                        pdf_bytes = self._fetch_resume_pdf(resume_original_full, apply_id, resume_id)
                        self._save_pdf(c, pdf_bytes)
                        ok += 1
                except requests.RequestException:
                    continue

            # 如果用户只指定了某一页，处理完这一页即停止
            if self.page_no is not None:
                break

            page += 1


def login_only(corp_code: str, username: str, password: str) -> None:
    """仅执行登录验证，不抓取简历。成功输出 JSON 并退出码 0，失败输出错误 JSON 并退出码非 0。"""
    import sys
    try:
        # 创建 session 并尝试登录
        session = requests.Session()
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
        session.mount("https://", adapter)
        session.mount("http://", adapter)
        
        # 获取登录信息
        info_url = f"{WintalentFetcher.INFO_URL}?r=0.123456789"
        resp = session.post(info_url, data={"corpCode": corp_code}, headers={
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "X-Requested-With": "XMLHttpRequest",
            "Origin": WintalentFetcher.BASE,
            "Referer": f"{WintalentFetcher.BASE}/interviewer/login"
        })
        
        if resp.status_code != 200:
            print(json.dumps({"status": "error", "message": f"获取登录信息失败: HTTP {resp.status_code}"}, ensure_ascii=False))
            sys.exit(1)
        
        try:
            info = resp.json()
        except Exception:
            print(json.dumps({"status": "error", "message": "解析登录信息失败"}, ensure_ascii=False))
            sys.exit(1)
        
        time_sign = info.get("time")
        pass_key = info.get("interviewerPassKey")
        pass_type = info.get("interviewerPassType")
        
        if str(pass_type) == "3":
            print(json.dumps({"status": "error", "message": "服务器要求 SM4 加密，当前脚本仅支持 AES"}, ensure_ascii=False))
            sys.exit(3)
        
        if not pass_key:
            print(json.dumps({"status": "error", "message": "未获取到加密密钥"}, ensure_ascii=False))
            sys.exit(3)
        
        # AES 加密密码
        key = pass_key.encode("utf-8")
        data = password.encode("utf-8")
        cipher = AES.new(key, AES.MODE_ECB)
        encrypted_bytes = cipher.encrypt(pad(data, AES.block_size))
        encrypted_password = binascii.hexlify(encrypted_bytes).decode("utf-8")
        
        # 执行登录
        payload = {
            "remremberMeInput": "true",
            "timeSign": time_sign,
            "actycoStatus": "",
            "loginMode": "",
            "localeType": "1",
            "login-select": "1",
            "corpCode": corp_code,
            "userName": username,
            "password": encrypted_password,
            "verifyCode": "",
            "phoneNum": "",
            "verification": "",
            "remremberMe": "true"
        }
        
        login_resp = session.post(WintalentFetcher.LOGIN_URL, data=payload, headers={
            "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
            "X-Requested-With": "XMLHttpRequest",
            "Origin": WintalentFetcher.BASE,
            "Referer": f"{WintalentFetcher.BASE}/interviewer/login"
        })
        
        if login_resp.status_code != 200:
            print(json.dumps({"status": "error", "message": f"登录请求失败: HTTP {login_resp.status_code}"}, ensure_ascii=False))
            sys.exit(1)
        
        # 检查是否获取到 SESSION cookie
        cookies = session.cookies.get_dict()
        if "SESSION" not in cookies:
            # 尝试解析响应体获取错误信息
            try:
                resp_data = login_resp.json()
                error_msg = resp_data.get("message") or resp_data.get("msg") or "用户名或密码错误"
            except Exception:
                error_msg = "登录失败，未返回 SESSION cookie"
            print(json.dumps({"status": "error", "message": error_msg}, ensure_ascii=False))
            sys.exit(1)
        
        # 登录成功
        print(json.dumps({
            "status": "success",
            "user": {
                "corp_code": corp_code,
                "username": username
            }
        }, ensure_ascii=False))
        sys.exit(0)
        
    except requests.RequestException as e:
        print(json.dumps({"status": "error", "message": f"网络连接失败: {str(e)}"}, ensure_ascii=False))
        sys.exit(2)
    except Exception as e:
        print(json.dumps({"status": "error", "message": str(e)}, ensure_ascii=False))
        sys.exit(1)


if __name__ == "__main__":
    import argparse
    import sys
    
    parser = argparse.ArgumentParser(description="Wintalent简历抓取工具")
    parser.add_argument("--corp-code", help="组织代码（如：motern）")
    parser.add_argument("--username", help="登录账号")
    parser.add_argument("--password", help="登录密码")
    parser.add_argument("--current-code", default="0/430400/430447/430448/430449", help="当前代码路径")
    parser.add_argument("--page-no", type=int, help="页码")
    parser.add_argument("--page-size", type=int, help="每页大小")
    parser.add_argument("--output-dir", default="out", help="输出目录")
    parser.add_argument("--session-cookie", help="会话Cookie（可选，如果提供则跳过登录）")
    parser.add_argument("--json-out", action="store_true", help="以JSON打印抓取结果（resume_pdf以base64编码）")
    parser.add_argument("--login-only", action="store_true", help="仅执行登录验证，不抓取简历")
    parser.add_argument("--positions-only", action="store_true", help="仅抓取岗位列表（去重），不抓取简历")
    
    args = parser.parse_args()

    # 默认写死凭据与分页参数（可被传入参数覆盖）
    corp_code = args.corp_code or "motern"
    username = args.username or os.getenv("WT_USERNAME")
    password = args.password or os.getenv("WT_PASSWORD")
    page_no = args.page_no or 1
    page_size = args.page_size or 100
    
    # 登录验证模式
    if args.login_only:
        if not username or not password:
            print(json.dumps({"status": "error", "message": "登录验证需要提供 username 和 password"}, ensure_ascii=False))
            sys.exit(1)
        login_only(corp_code, username, password)
    else:
        fetcher = WintalentFetcher(
            corp_code=corp_code,
            username=username,
            password=password,
            current_code=args.current_code,
            page_no=page_no,
            page_size=page_size,
            output_dir=args.output_dir,
            session_cookie=args.session_cookie,
        )

        if args.positions_only:
            # 仅抓取岗位列表（去重）
            positions = fetcher.fetch_positions()
            print(json.dumps({"total": len(positions), "positions": positions}, ensure_ascii=False))
        elif args.json_out:
            # 直接在脚本内自动分页抓取全部数据并返回 { total, items }
            data = fetcher.fetch_all()
            items: List[Dict[str, Any]] = []
            for item in data:
                pdf_bytes = item.get("resume_pdf")
                b64 = base64.b64encode(pdf_bytes).decode("utf-8") if isinstance(pdf_bytes, (bytes, bytearray)) else None
                items.append({
                    "name": item.get("name"),
                    "apply_id": item.get("apply_id", ""),
                    "resume_id": item.get("resume_id", ""),
                    "post_id": item.get("post_id", ""),
                    "post_name": item.get("post_name", ""),
                    "recruit_type": item.get("recruit_type", ""),
                    "jd": item.get("jd"),
                    "resume_pdf_b64": b64,
                })
            print(json.dumps({"total": len(items), "items": items}, ensure_ascii=False))
        else:
            fetcher.run()
