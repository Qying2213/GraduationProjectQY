#!/bin/bash
# ==============================================================================
# 智能人才运营平台 - 全面测试脚本
# ==============================================================================
# 使用方法: cd ztest && ./run_all_tests.sh
# 前置条件: 确保所有后端服务已启动 (./start-all.sh)
# ==============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
BASE_URL="http://localhost:8080/api/v1"
REPORT_FILE="test_report_$(date +%Y%m%d_%H%M%S).md"

# 计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# 测试结果数组
declare -a TEST_RESULTS

# 日志函数
log_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
log_pass() { echo -e "${GREEN}[PASS]${NC} $1"; ((PASSED_TESTS++)); TEST_RESULTS+=("✅ $1"); }
log_fail() { echo -e "${RED}[FAIL]${NC} $1"; ((FAILED_TESTS++)); TEST_RESULTS+=("❌ $1 - $2"); }
log_skip() { echo -e "${YELLOW}[SKIP]${NC} $1"; ((SKIPPED_TESTS++)); TEST_RESULTS+=("⏭️ $1 - 跳过"); }

# HTTP 请求函数
http_get() {
    local url="$1"
    local token="$2"
    if [ -n "$token" ]; then
        curl -sS -w "\n%{http_code}" -H "Authorization: Bearer $token" "$url" 2>/dev/null
    else
        curl -sS -w "\n%{http_code}" "$url" 2>/dev/null
    fi
}

http_post() {
    local url="$1"
    local data="$2"
    local token="$3"
    if [ -n "$token" ]; then
        curl -sS -w "\n%{http_code}" -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d "$data" "$url" 2>/dev/null
    else
        curl -sS -w "\n%{http_code}" -X POST -H "Content-Type: application/json" -d "$data" "$url" 2>/dev/null
    fi
}

http_put() {
    local url="$1"
    local data="$2"
    local token="$3"
    curl -sS -w "\n%{http_code}" -X PUT -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d "$data" "$url" 2>/dev/null
}

http_delete() {
    local url="$1"
    local token="$2"
    curl -sS -w "\n%{http_code}" -X DELETE -H "Authorization: Bearer $token" "$url" 2>/dev/null
}

# 解析响应
parse_response() {
    local response="$1"
    local body=$(echo "$response" | head -n -1)
    local status=$(echo "$response" | tail -n 1)
    echo "$body|$status"
}

# 检查服务是否运行
check_service() {
    local name="$1"
    local port="$2"
    if nc -z localhost "$port" 2>/dev/null; then
        log_pass "服务 $name (:$port) 运行正常"
        return 0
    else
        log_fail "服务 $name (:$port) 未运行" "请启动服务"
        return 1
    fi
}

# ==============================================================================
# 测试用例
# ==============================================================================

echo ""
echo "╔══════════════════════════════════════════════════════════════════════╗"
echo "║        智能人才运营平台 - 全面测试脚本                                  ║"
echo "╚══════════════════════════════════════════════════════════════════════╝"
echo ""
echo "开始时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# ------------------------------------------------------------------------------
# 1. 服务健康检查
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📡 测试 1: 服务健康检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

((TOTAL_TESTS++))
check_service "Gateway" 8080
((TOTAL_TESTS++))
check_service "User Service" 8081
((TOTAL_TESTS++))
check_service "Job Service" 8082
((TOTAL_TESTS++))
check_service "Interview Service" 8083
((TOTAL_TESTS++))
check_service "Resume Service" 8084
((TOTAL_TESTS++))
check_service "Message Service" 8085
((TOTAL_TESTS++))
check_service "Talent Service" 8086
((TOTAL_TESTS++))
check_service "Recommendation Service" 8087

# Gateway 健康检查端点
((TOTAL_TESTS++))
response=$(http_get "http://localhost:8080/health")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "Gateway 健康检查端点 (/health)"
else
    log_fail "Gateway 健康检查端点" "状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 2. 用户服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "👤 测试 2: 用户服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 2.1 用户注册测试
((TOTAL_TESTS++))
RANDOM_EMAIL="test_$(date +%s)@example.com"
response=$(http_post "$BASE_URL/register" "{\"username\":\"testuser_$$\",\"email\":\"$RANDOM_EMAIL\",\"password\":\"Test123456\",\"role\":\"hr\"}")
parsed=$(parse_response "$response")
body=$(echo "$parsed" | cut -d'|' -f1)
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    code=$(echo "$body" | grep -o '"code":[0-9]*' | cut -d: -f2)
    if [ "$code" == "0" ]; then
        log_pass "用户注册 - 新用户注册成功"
    else
        log_fail "用户注册" "业务码 code=$code"
    fi
else
    log_fail "用户注册" "HTTP状态码 $status"
fi

# 2.2 用户登录测试 (使用 admin 账号)
((TOTAL_TESTS++))
response=$(http_post "$BASE_URL/login" '{"username":"admin","password":"admin123"}')
parsed=$(parse_response "$response")
body=$(echo "$parsed" | cut -d'|' -f1)
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    code=$(echo "$body" | grep -o '"code":[0-9]*' | cut -d: -f2)
    if [ "$code" == "0" ]; then
        TOKEN=$(echo "$body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
        log_pass "用户登录 - admin 登录成功"
    else
        log_fail "用户登录" "业务码 code=$code"
    fi
else
    log_fail "用户登录" "HTTP状态码 $status"
fi

# 2.3 获取用户信息测试
((TOTAL_TESTS++))
if [ -n "$TOKEN" ]; then
    response=$(http_get "$BASE_URL/profile" "$TOKEN")
    parsed=$(parse_response "$response")
    body=$(echo "$parsed" | cut -d'|' -f1)
    status=$(echo "$parsed" | cut -d'|' -f2)
    if [ "$status" == "200" ]; then
        code=$(echo "$body" | grep -o '"code":[0-9]*' | cut -d: -f2)
        if [ "$code" == "0" ]; then
            log_pass "获取用户信息 - GET /profile"
        else
            log_fail "获取用户信息" "业务码 code=$code"
        fi
    else
        log_fail "获取用户信息" "HTTP状态码 $status"
    fi
else
    log_skip "获取用户信息 - 未获取到Token"
fi

# 2.4 错误密码登录测试
((TOTAL_TESTS++))
response=$(http_post "$BASE_URL/login" '{"username":"admin","password":"wrongpassword"}')
parsed=$(parse_response "$response")
body=$(echo "$parsed" | cut -d'|' -f1)
status=$(echo "$parsed" | cut -d'|' -f2)
code=$(echo "$body" | grep -o '"code":[0-9]*' | cut -d: -f2)
if [ "$code" != "0" ]; then
    log_pass "错误密码登录 - 正确拒绝"
else
    log_fail "错误密码登录" "应该被拒绝但成功了"
fi

echo ""

# ------------------------------------------------------------------------------
# 3. 职位服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💼 测试 3: 职位服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 3.1 获取职位列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/jobs" "$TOKEN")
parsed=$(parse_response "$response")
body=$(echo "$parsed" | cut -d'|' -f1)
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取职位列表 - GET /jobs"
else
    log_fail "获取职位列表" "HTTP状态码 $status"
fi

# 3.2 创建职位
((TOTAL_TESTS++))
if [ -n "$TOKEN" ]; then
    JOB_DATA='{
        "title": "测试职位_'$$'",
        "department": "技术部",
        "location": "北京",
        "job_type": "full_time",
        "salary_min": 20000,
        "salary_max": 40000,
        "description": "这是测试职位描述",
        "requirements": "本科以上学历",
        "status": "open"
    }'
    response=$(http_post "$BASE_URL/jobs" "$JOB_DATA" "$TOKEN")
    parsed=$(parse_response "$response")
    body=$(echo "$parsed" | cut -d'|' -f1)
    status=$(echo "$parsed" | cut -d'|' -f2)
    if [ "$status" == "200" ] || [ "$status" == "201" ]; then
        code=$(echo "$body" | grep -o '"code":[0-9]*' | cut -d: -f2)
        if [ "$code" == "0" ]; then
            JOB_ID=$(echo "$body" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
            log_pass "创建职位 - POST /jobs"
        else
            log_fail "创建职位" "业务码 code=$code"
        fi
    else
        log_fail "创建职位" "HTTP状态码 $status"
    fi
else
    log_skip "创建职位 - 未获取到Token"
fi

echo ""

# ------------------------------------------------------------------------------
# 4. 简历服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📄 测试 4: 简历服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 4.1 获取简历列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/resumes" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取简历列表 - GET /resumes"
else
    log_fail "获取简历列表" "HTTP状态码 $status"
fi

# 4.2 AI 评估接口检查
((TOTAL_TESTS++))
response=$(http_post "$BASE_URL/ai/evaluate" '{"resume_id":1,"job_id":1}' "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ] || [ "$status" == "400" ] || [ "$status" == "404" ]; then
    log_pass "AI 评估接口 - 可访问 (POST /ai/evaluate)"
else
    log_fail "AI 评估接口" "HTTP状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 5. 人才服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "👥 测试 5: 人才服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 5.1 获取人才列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/talents" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取人才列表 - GET /talents"
else
    log_fail "获取人才列表" "HTTP状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 6. 消息服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💬 测试 6: 消息服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 6.1 获取消息列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/messages" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取消息列表 - GET /messages"
else
    log_fail "获取消息列表" "HTTP状态码 $status"
fi

# 6.2 获取未读消息数
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/messages/unread-count" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取未读消息数 - GET /messages/unread-count"
else
    log_fail "获取未读消息数" "HTTP状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 7. 推荐服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🎯 测试 7: 推荐服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 7.1 获取推荐列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/recommendations" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取推荐列表 - GET /recommendations"
else
    log_fail "获取推荐列表" "HTTP状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 8. 面试服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📅 测试 8: 面试服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 8.1 获取面试列表
((TOTAL_TESTS++))
response=$(http_get "$BASE_URL/interviews" "$TOKEN")
parsed=$(parse_response "$response")
status=$(echo "$parsed" | cut -d'|' -f2)
if [ "$status" == "200" ]; then
    log_pass "获取面试列表 - GET /interviews"
else
    log_fail "获取面试列表" "HTTP状态码 $status"
fi

echo ""

# ------------------------------------------------------------------------------
# 9. 统计服务测试
# ------------------------------------------------------------------------------
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 测试 9: 统计服务"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

STATS_ENDPOINTS=(
    "dashboard|仪表盘统计"
    "funnel|招聘漏斗"
    "channels|渠道统计"
    "department-progress|部门进度"
    "interviewer-rank|面试官排行"
    "trend|趋势数据"
    "job-rank|职位排行"
)

for endpoint_info in "${STATS_ENDPOINTS[@]}"; do
    endpoint=$(echo "$endpoint_info" | cut -d'|' -f1)
    name=$(echo "$endpoint_info" | cut -d'|' -f2)
    ((TOTAL_TESTS++))
    response=$(http_get "$BASE_URL/stats/$endpoint" "$TOKEN")
    parsed=$(parse_response "$response")
    status=$(echo "$parsed" | cut -d'|' -f2)
    if [ "$status" == "200" ]; then
        log_pass "$name - GET /stats/$endpoint"
    else
        log_fail "$name" "HTTP状态码 $status"
    fi
done

echo ""

# ==============================================================================
# 生成测试报告
# ==============================================================================
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 测试总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo -e "总测试数:   ${BLUE}$TOTAL_TESTS${NC}"
echo -e "通过:       ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败:       ${RED}$FAILED_TESTS${NC}"
echo -e "跳过:       ${YELLOW}$SKIPPED_TESTS${NC}"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✅ 全部测试通过!${NC}"
else
    echo -e "${RED}⚠️ 存在失败的测试，请检查日志${NC}"
fi

# 生成 Markdown 报告
cat > "$REPORT_FILE" << EOF
# 测试报告

**执行时间**: $(date '+%Y-%m-%d %H:%M:%S')

## 统计

| 指标 | 数量 |
|------|------|
| 总测试数 | $TOTAL_TESTS |
| ✅ 通过 | $PASSED_TESTS |
| ❌ 失败 | $FAILED_TESTS |
| ⏭️ 跳过 | $SKIPPED_TESTS |

## 测试结果详情

EOF

for result in "${TEST_RESULTS[@]}"; do
    echo "- $result" >> "$REPORT_FILE"
done

echo ""
echo "报告已保存到: $REPORT_FILE"
echo ""
echo "结束时间: $(date '+%Y-%m-%d %H:%M:%S')"
