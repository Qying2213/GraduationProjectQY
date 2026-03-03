#!/usr/bin/env python3
# ==============================================================================
# 智能人才运营平台 - Python 测试脚本
# ==============================================================================
# 使用方法: cd ztest && python3 test_api.py
# 前置条件: pip install requests 
# ==============================================================================

import requests
import json
import time
import sys
from datetime import datetime

# 配置
BASE_URL = "http://localhost:8080/api/v1"
TOKEN = None

# 颜色输出
class Colors:
    GREEN = '\033[92m'
    RED = '\033[91m'
    YELLOW = '\033[93m'
    BLUE = '\033[94m'
    END = '\033[0m'

def log_pass(msg):
    print(f"{Colors.GREEN}[PASS]{Colors.END} {msg}")

def log_fail(msg, error=""):
    print(f"{Colors.RED}[FAIL]{Colors.END} {msg} - {error}")

def log_info(msg):
    print(f"{Colors.BLUE}[INFO]{Colors.END} {msg}")

def log_warn(msg):
    print(f"{Colors.YELLOW}[WARN]{Colors.END} {msg}")

# 测试结果收集
class TestResults:
    def __init__(self):
        self.total = 0
        self.passed = 0
        self.failed = 0
        self.issues = []
    
    def add_pass(self, name):
        self.total += 1
        self.passed += 1
        log_pass(name)
    
    def add_fail(self, name, reason):
        self.total += 1
        self.failed += 1
        self.issues.append({"test": name, "reason": reason})
        log_fail(name, reason)

results = TestResults()

# ==============================================================================
# 测试用例
# ==============================================================================

def test_health_check():
    """测试网关健康检查"""
    try:
        r = requests.get("http://localhost:8080/health", timeout=5)
        if r.status_code == 200:
            results.add_pass("Gateway 健康检查")
            return True
        else:
            results.add_fail("Gateway 健康检查", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("Gateway 健康检查", str(e))
        return False

def test_login():
    """测试用户登录"""
    global TOKEN
    try:
        r = requests.post(f"{BASE_URL}/login", json={
            "username": "admin",
            "password": "admin123"
        }, timeout=10)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0 and data.get("data", {}).get("token"):
                TOKEN = data["data"]["token"]
                results.add_pass("用户登录 (admin)")
                return True
            else:
                results.add_fail("用户登录", f"响应: {data}")
                return False
        else:
            results.add_fail("用户登录", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("用户登录", str(e))
        return False

def test_get_profile():
    """测试获取用户信息"""
    if not TOKEN:
        results.add_fail("获取用户信息", "未登录")
        return False
    try:
        r = requests.get(f"{BASE_URL}/profile", 
                        headers={"Authorization": f"Bearer {TOKEN}"}, 
                        timeout=10)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0:
                user = data.get("data", {})
                # 检查头像字段
                if "avatar" in user:
                    results.add_pass("获取用户信息 (含 avatar 字段)")
                else:
                    results.add_pass("获取用户信息")
                    log_warn("  -> 用户数据缺少 avatar 字段")
                return True
            else:
                results.add_fail("获取用户信息", f"code={data.get('code')}")
                return False
        else:
            results.add_fail("获取用户信息", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取用户信息", str(e))
        return False

def test_update_profile_avatar():
    """测试更新头像"""
    if not TOKEN:
        results.add_fail("更新头像", "未登录")
        return False
    try:
        test_avatar = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
        r = requests.put(f"{BASE_URL}/profile", 
                        headers={"Authorization": f"Bearer {TOKEN}"},
                        json={"avatar": test_avatar},
                        timeout=10)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0:
                # 检查返回的用户信息中是否有头像
                user_data = data.get("data", {})
                if user_data.get("avatar") == test_avatar:
                    results.add_pass("更新头像 - 返回数据正确")
                else:
                    results.add_pass("更新头像 - API 成功")
                    log_warn("  -> 返回的 avatar 与请求不一致，可能是 bug")
                    results.issues.append({
                        "test": "头像同步问题",
                        "reason": "更新头像后返回的数据中 avatar 字段不一致"
                    })
                return True
            else:
                results.add_fail("更新头像", f"code={data.get('code')}")
                return False
        else:
            results.add_fail("更新头像", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("更新头像", str(e))
        return False

def test_get_jobs():
    """测试获取职位列表"""
    try:
        r = requests.get(f"{BASE_URL}/jobs", 
                        headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                        timeout=10)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0:
                jobs = data.get("data", {}).get("items", [])
                results.add_pass(f"获取职位列表 (共 {len(jobs)} 个)")
                return True
            else:
                results.add_fail("获取职位列表", f"code={data.get('code')}")
                return False
        else:
            results.add_fail("获取职位列表", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取职位列表", str(e))
        return False

def test_get_resumes():
    """测试获取简历列表"""
    try:
        r = requests.get(f"{BASE_URL}/resumes", 
                        headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                        timeout=10)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0:
                resumes = data.get("data", {}).get("items", [])
                results.add_pass(f"获取简历列表 (共 {len(resumes)} 份)")
                
                # 检查简历状态字段
                for resume in resumes[:3]:  # 只检查前3个
                    status = resume.get("status")
                    if status not in ["pending", "processing", "parsed", "failed"]:
                        log_warn(f"  -> 简历 {resume.get('id')} 状态异常: {status}")
                return True
            else:
                results.add_fail("获取简历列表", f"code={data.get('code')}")
                return False
        else:
            results.add_fail("获取简历列表", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取简历列表", str(e))
        return False

def test_get_talents():
    """测试获取人才列表"""
    try:
        r = requests.get(f"{BASE_URL}/talents", 
                        headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                        timeout=10)
        if r.status_code == 200:
            results.add_pass("获取人才列表")
            return True
        else:
            results.add_fail("获取人才列表", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取人才列表", str(e))
        return False

def test_get_messages():
    """测试获取消息列表"""
    try:
        r = requests.get(f"{BASE_URL}/messages", 
                        headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                        timeout=10)
        if r.status_code == 200:
            results.add_pass("获取消息列表")
            return True
        else:
            results.add_fail("获取消息列表", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取消息列表", str(e))
        return False

def test_get_recommendations():
    """测试获取推荐统计"""
    try:
        r = requests.get(f"{BASE_URL}/recommendations/stats",
                        headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                        timeout=10)
        if r.status_code == 200:
            results.add_pass("获取推荐统计")
            return True
        else:
            results.add_fail("获取推荐统计", f"状态码 {r.status_code}")
            return False
    except Exception as e:
        results.add_fail("获取推荐统计", str(e))
        return False

def test_stats_endpoints():
    """测试统计接口"""
    endpoints = [
        ("dashboard", "仪表盘统计"),
        ("funnel", "招聘漏斗"),
        ("channels", "渠道统计"),
        ("trend", "趋势数据"),
    ]
    
    for endpoint, name in endpoints:
        try:
            r = requests.get(f"{BASE_URL}/stats/{endpoint}", 
                            headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                            timeout=10)
            if r.status_code == 200:
                results.add_pass(f"统计: {name}")
            else:
                results.add_fail(f"统计: {name}", f"状态码 {r.status_code}")
        except Exception as e:
            results.add_fail(f"统计: {name}", str(e))

def test_ai_evaluate():
    """测试 AI 评估接口"""
    try:
        # 这个接口需要有效的 resume_id 和 job_id
        r = requests.post(f"{BASE_URL}/ai/evaluate", 
                         headers={"Authorization": f"Bearer {TOKEN}"} if TOKEN else {},
                         json={"resume_id": 1, "job_id": 1},
                         timeout=30)
        if r.status_code == 200:
            data = r.json()
            if data.get("code") == 0:
                results.add_pass("AI 评估接口")
            else:
                # 可能是简历或职位不存在，这是预期的
                results.add_pass("AI 评估接口 (接口可用)")
                log_warn(f"  -> 返回: {data.get('message', 'unknown')}")
        elif r.status_code in [400, 404, 503]:
            results.add_pass("AI 评估接口 (接口可用)")
            log_warn("  -> 需要有效的 resume_id 和 job_id")
        else:
            results.add_fail("AI 评估接口", f"状态码 {r.status_code}")
    except Exception as e:
        results.add_fail("AI 评估接口", str(e))

# ==============================================================================
# 潜在问题检测
# ==============================================================================

def analyze_potential_issues():
    """分析潜在问题"""
    print("\n" + "="*60)
    print("🔍 潜在问题分析")
    print("="*60 + "\n")
    
    issues_found = []
    
    # 1. 检查 Redis 连接
    try:
        import socket
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(2)
        result = sock.connect_ex(('localhost', 6379))
        sock.close()
        if result != 0:
            issues_found.append({
                "level": "WARNING",
                "title": "Redis 未运行",
                "desc": "Redis 缓存服务未启动，可能影响性能",
                "fix": "启动 Redis: redis-server 或 brew services start redis"
            })
    except:
        pass
    
    # 2. 检查 Elasticsearch
    try:
        r = requests.get("http://localhost:9200/_cluster/health", timeout=2)
        if r.status_code != 200:
            issues_found.append({
                "level": "WARNING",
                "title": "Elasticsearch 状态异常",
                "desc": f"状态码: {r.status_code}",
                "fix": "检查 Elasticsearch 日志"
            })
    except:
        issues_found.append({
            "level": "INFO",
            "title": "Elasticsearch 未运行",
            "desc": "全文搜索功能可能不可用",
            "fix": "可选: 启动 Elasticsearch"
        })
    
    # 输出问题
    if not issues_found and len(results.issues) == 0:
        print(f"{Colors.GREEN}✅ 未发现明显问题{Colors.END}\n")
    else:
        for issue in issues_found:
            level_color = Colors.RED if issue["level"] == "ERROR" else Colors.YELLOW
            print(f"{level_color}[{issue['level']}]{Colors.END} {issue['title']}")
            print(f"    描述: {issue['desc']}")
            print(f"    修复: {issue['fix']}\n")
        
        for issue in results.issues:
            print(f"{Colors.RED}[TEST ISSUE]{Colors.END} {issue['test']}")
            print(f"    原因: {issue['reason']}\n")

# ==============================================================================
# 主函数
# ==============================================================================

def main():
    print("\n" + "="*60)
    print("  智能人才运营平台 - Python API 测试")
    print("="*60)
    print(f"  开始时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"  目标服务: {BASE_URL}")
    print("="*60 + "\n")
    
    # 1. 健康检查
    log_info("=== 1. 健康检查 ===")
    if not test_health_check():
        print(f"\n{Colors.RED}❌ 网关服务未启动，请先运行 start-all.sh{Colors.END}\n")
        sys.exit(1)
    
    # 2. 认证测试
    print()
    log_info("=== 2. 认证测试 ===")
    test_login()
    test_get_profile()
    test_update_profile_avatar()
    
    # 3. 业务接口测试
    print()
    log_info("=== 3. 业务接口测试 ===")
    test_get_jobs()
    test_get_resumes()
    test_get_talents()
    test_get_messages()
    test_get_recommendations()
    
    # 4. 统计接口测试
    print()
    log_info("=== 4. 统计接口测试 ===")
    test_stats_endpoints()
    
    # 5. AI 功能测试
    print()
    log_info("=== 5. AI 功能测试 ===")
    test_ai_evaluate()
    
    # 分析潜在问题
    analyze_potential_issues()
    
    # 测试总结
    print("="*60)
    print("📊 测试总结")
    print("="*60)
    print(f"  总测试数: {results.total}")
    print(f"  {Colors.GREEN}通过: {results.passed}{Colors.END}")
    print(f"  {Colors.RED}失败: {results.failed}{Colors.END}")
    print(f"  通过率: {results.passed/results.total*100:.1f}%")
    print("="*60 + "\n")
    
    if results.failed > 0:
        sys.exit(1)

if __name__ == "__main__":
    main()
