#!/bin/bash
# =====================================================
# 智能人才运营平台 - 完整API功能测试脚本
# 测试所有后端服务的主要功能（除Coze工作流外）
# 共74个测试用例
# =====================================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 服务地址 - 根据实际配置调整
USER_SERVICE="http://localhost:8081"
JOB_SERVICE="http://localhost:8082"
INTERVIEW_SERVICE="http://localhost:8083"
RESUME_SERVICE="http://localhost:8084"
MESSAGE_SERVICE="http://localhost:8085"
TALENT_SERVICE="http://localhost:8086"

# 计数器
TOTAL=0
PASSED=0
FAILED=0

# 保存创建的资源ID
CREATED_USER_ID=""
CREATED_JOB_ID=""
CREATED_TALENT_ID=""
CREATED_INTERVIEW_ID=""
CREATED_MESSAGE_ID=""

# 测试函数
test_api() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_code="$5"
    
    TOTAL=$((TOTAL + 1))
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -d "$data" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} - $name (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        echo "$body"  # 返回body供后续使用
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} - $name (期望: $expected_code, 实际: $http_code)"
        echo -e "  响应: $(echo "$body" | head -c 150)"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 静默测试（不输出body）
test_api_silent() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_code="$5"
    
    TOTAL=$((TOTAL + 1))
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" \
            -d "$data" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" 2>/dev/null)
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" = "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} - $name (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} - $name (期望: $expected_code, 实际: $http_code)"
        echo -e "  响应: $(echo "$body" | head -c 150)"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 检查服务是否运行（尝试health端点或根路径）
check_service() {
    local name="$1"
    local url="$2"
    
    # 先尝试health端点
    if curl -s --connect-timeout 2 "$url/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $name 运行中"
        return 0
    fi
    
    # 再尝试api端点
    if curl -s --connect-timeout 2 "$url/api/v1" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} $name 运行中"
        return 0
    fi
    
    # 最后尝试根路径
    response=$(curl -s --connect-timeout 2 -o /dev/null -w "%{http_code}" "$url/" 2>/dev/null)
    if [ "$response" != "000" ]; then
        echo -e "${GREEN}✓${NC} $name 运行中"
        return 0
    fi
    
    echo -e "${RED}✗${NC} $name 未运行"
    return 1
}

echo ""
echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}         智能人才运营平台 - 完整API功能测试${NC}"
echo -e "${BLUE}         共74个测试用例（除Coze工作流外）${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

# =====================================================
# 1. 检查服务状态
# =====================================================
echo -e "${CYAN}[1/8] 检查服务状态${NC}"
echo "----------------------------------------"

services_ok=true
check_service "user-service (8081)" "$USER_SERVICE" || services_ok=false
check_service "job-service (8082)" "$JOB_SERVICE" || services_ok=false
check_service "interview-service (8083)" "$INTERVIEW_SERVICE" || services_ok=false
check_service "resume-service (8084)" "$RESUME_SERVICE" || services_ok=false
check_service "message-service (8085)" "$MESSAGE_SERVICE" || services_ok=false
check_service "talent-service (8086)" "$TALENT_SERVICE" || services_ok=false

echo ""

if [ "$services_ok" = false ]; then
    echo -e "${YELLOW}警告: 部分服务未运行，相关测试可能失败${NC}"
    echo -e "${YELLOW}请先启动所有服务: ./start-backend.sh${NC}"
    echo ""
fi

# =====================================================
# 2. 用户服务测试 (user-service)
# =====================================================
echo -e "${CYAN}[2/8] 用户服务测试 (user-service) - 8个测试${NC}"
echo "----------------------------------------"

# 2.1 用户登录
test_api_silent "用户登录 - 正确密码" "POST" "$USER_SERVICE/api/v1/login" \
    '{"username":"admin","password":"password123"}' "200"

test_api_silent "用户登录 - 错误密码" "POST" "$USER_SERVICE/api/v1/login" \
    '{"username":"admin","password":"wrongpassword"}' "401"

test_api_silent "用户登录 - 不存在的用户" "POST" "$USER_SERVICE/api/v1/login" \
    '{"username":"notexist","password":"password123"}' "401"

# 2.2 用户注册
TIMESTAMP=$(date +%s)
test_api_silent "用户注册 - 新用户" "POST" "$USER_SERVICE/api/v1/register" \
    '{"username":"testuser_'$TIMESTAMP'","email":"test'$TIMESTAMP'@test.com","password":"test123456","role":"hr"}' "201"

test_api_silent "用户注册 - 重复用户名" "POST" "$USER_SERVICE/api/v1/register" \
    '{"username":"admin","email":"new@test.com","password":"test123456"}' "400"

test_api_silent "用户注册 - 重复邮箱" "POST" "$USER_SERVICE/api/v1/register" \
    '{"username":"newuser123","email":"admin@company.com","password":"test123456"}' "400"

# 2.3 获取用户列表
test_api_silent "获取用户列表" "GET" "$USER_SERVICE/api/v1/users" "" "200"

test_api_silent "获取用户列表 - 分页" "GET" "$USER_SERVICE/api/v1/users?page=1&page_size=5" "" "200"

echo ""

# =====================================================
# 3. 职位服务测试 (job-service)
# =====================================================
echo -e "${CYAN}[3/8] 职位服务测试 (job-service) - 14个测试${NC}"
echo "----------------------------------------"

# 3.1 获取职位列表
test_api_silent "获取职位列表" "GET" "$JOB_SERVICE/api/v1/jobs" "" "200"

test_api_silent "获取职位列表 - 分页" "GET" "$JOB_SERVICE/api/v1/jobs?page=1&page_size=5" "" "200"

test_api_silent "获取职位列表 - 第二页" "GET" "$JOB_SERVICE/api/v1/jobs?page=2&page_size=5" "" "200"

test_api_silent "获取职位列表 - 按状态筛选(open)" "GET" "$JOB_SERVICE/api/v1/jobs?status=open" "" "200"

test_api_silent "获取职位列表 - 按状态筛选(closed)" "GET" "$JOB_SERVICE/api/v1/jobs?status=closed" "" "200"

test_api_silent "获取职位列表 - 按地点筛选" "GET" "$JOB_SERVICE/api/v1/jobs?location=北京" "" "200"

test_api_silent "获取职位列表 - 按类型筛选" "GET" "$JOB_SERVICE/api/v1/jobs?type=full-time" "" "200"

test_api_silent "获取职位列表 - 关键词搜索" "GET" "$JOB_SERVICE/api/v1/jobs?keyword=前端" "" "200"

# 3.2 获取职位详情
test_api_silent "获取职位详情 (ID=1)" "GET" "$JOB_SERVICE/api/v1/jobs/1" "" "200"

test_api_silent "获取职位详情 (ID=2)" "GET" "$JOB_SERVICE/api/v1/jobs/2" "" "200"

test_api_silent "获取不存在的职位" "GET" "$JOB_SERVICE/api/v1/jobs/9999" "" "404"

# 3.3 获取职位统计
test_api_silent "获取职位统计" "GET" "$JOB_SERVICE/api/v1/jobs/stats" "" "200"

# 3.4 创建职位
test_api_silent "创建新职位" "POST" "$JOB_SERVICE/api/v1/jobs" \
    '{"title":"测试工程师_'$TIMESTAMP'","description":"这是一个测试职位","location":"上海","salary":"15-25K","type":"full-time","status":"open","department":"技术部","created_by":1}' "201"

# 3.5 更新职位
test_api_silent "更新职位" "PUT" "$JOB_SERVICE/api/v1/jobs/1" \
    '{"title":"高级前端工程师(更新)","status":"open"}' "200"

echo ""

# =====================================================
# 4. 人才服务测试 (talent-service)
# =====================================================
echo -e "${CYAN}[4/8] 人才服务测试 (talent-service) - 12个测试${NC}"
echo "----------------------------------------"

# 4.1 获取人才列表
test_api_silent "获取人才列表" "GET" "$TALENT_SERVICE/api/v1/talents" "" "200"

test_api_silent "获取人才列表 - 分页" "GET" "$TALENT_SERVICE/api/v1/talents?page=1&page_size=5" "" "200"

test_api_silent "获取人才列表 - 第二页" "GET" "$TALENT_SERVICE/api/v1/talents?page=2&page_size=5" "" "200"

test_api_silent "获取人才列表 - 按状态筛选" "GET" "$TALENT_SERVICE/api/v1/talents?status=active" "" "200"

# 4.2 获取人才详情
test_api_silent "获取人才详情 (ID=1)" "GET" "$TALENT_SERVICE/api/v1/talents/1" "" "200"

test_api_silent "获取人才详情 (ID=5)" "GET" "$TALENT_SERVICE/api/v1/talents/5" "" "200"

test_api_silent "获取不存在的人才" "GET" "$TALENT_SERVICE/api/v1/talents/9999" "" "404"

# 4.3 搜索人才
test_api_silent "搜索人才 - 关键词(前端)" "GET" "$TALENT_SERVICE/api/v1/talents/search?keyword=前端" "" "200"

test_api_silent "搜索人才 - 关键词(Go)" "GET" "$TALENT_SERVICE/api/v1/talents/search?keyword=Go" "" "200"

test_api_silent "搜索人才 - 按经验筛选" "GET" "$TALENT_SERVICE/api/v1/talents/search?min_experience=3&max_experience=8" "" "200"

# 4.4 创建人才
test_api_silent "创建新人才" "POST" "$TALENT_SERVICE/api/v1/talents" \
    '{"name":"测试人才_'$TIMESTAMP'","email":"talent'$TIMESTAMP'@test.com","phone":"13900000099","skills":["Java","Python"],"experience":3,"education":"本科","status":"active","location":"深圳"}' "201"

# 4.5 更新人才
test_api_silent "更新人才信息" "PUT" "$TALENT_SERVICE/api/v1/talents/1" \
    '{"status":"active","experience":7}' "200"

echo ""

# =====================================================
# 5. 简历服务测试 (resume-service)
# =====================================================
echo -e "${CYAN}[5/8] 简历服务测试 (resume-service) - 14个测试${NC}"
echo "----------------------------------------"

# 5.1 获取简历列表
test_api_silent "获取简历列表" "GET" "$RESUME_SERVICE/api/v1/resumes" "" "200"

test_api_silent "获取简历列表 - 分页" "GET" "$RESUME_SERVICE/api/v1/resumes?page=1&page_size=5" "" "200"

test_api_silent "获取简历列表 - 第二页" "GET" "$RESUME_SERVICE/api/v1/resumes?page=2&page_size=5" "" "200"

test_api_silent "获取简历列表 - 按状态筛选(pending)" "GET" "$RESUME_SERVICE/api/v1/resumes?status=pending" "" "200"

test_api_silent "获取简历列表 - 按状态筛选(reviewing)" "GET" "$RESUME_SERVICE/api/v1/resumes?status=reviewing" "" "200"

test_api_silent "获取简历列表 - 按状态筛选(interviewed)" "GET" "$RESUME_SERVICE/api/v1/resumes?status=interviewed" "" "200"

test_api_silent "获取简历列表 - 按时间排序(降序)" "GET" "$RESUME_SERVICE/api/v1/resumes?sort_by=created_at&sort_order=desc" "" "200"

test_api_silent "获取简历列表 - 按时间排序(升序)" "GET" "$RESUME_SERVICE/api/v1/resumes?sort_by=created_at&sort_order=asc" "" "200"

# 5.2 获取简历详情
test_api_silent "获取简历详情 (ID=1)" "GET" "$RESUME_SERVICE/api/v1/resumes/1" "" "200"

test_api_silent "获取简历详情 (ID=5)" "GET" "$RESUME_SERVICE/api/v1/resumes/5" "" "200"

test_api_silent "获取不存在的简历" "GET" "$RESUME_SERVICE/api/v1/resumes/9999" "" "404"

# 5.3 更新简历状态
test_api_silent "更新简历状态" "PUT" "$RESUME_SERVICE/api/v1/resumes/1/status" \
    '{"status":"reviewing"}' "200"

# 5.4 AI配置检查（不测试实际评估）
test_api_silent "获取AI配置状态" "GET" "$RESUME_SERVICE/api/v1/ai/config" "" "200"

# 5.5 获取评估用简历列表
test_api_silent "获取评估用简历列表" "GET" "$RESUME_SERVICE/api/v1/resumes/evaluation" "" "200"

echo ""

# =====================================================
# 6. 面试服务测试 (interview-service)
# =====================================================
echo -e "${CYAN}[6/8] 面试服务测试 (interview-service) - 14个测试${NC}"
echo "----------------------------------------"

# 6.1 获取面试列表
test_api_silent "获取面试列表" "GET" "$INTERVIEW_SERVICE/api/v1/interviews" "" "200"

test_api_silent "获取面试列表 - 分页" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?page=1&page_size=5" "" "200"

test_api_silent "获取面试列表 - 第二页" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?page=2&page_size=5" "" "200"

test_api_silent "获取面试列表 - 按状态筛选(scheduled)" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?status=scheduled" "" "200"

test_api_silent "获取面试列表 - 按状态筛选(completed)" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?status=completed" "" "200"

# 6.2 获取面试详情
test_api_silent "获取面试详情 (ID=1)" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/1" "" "200"

test_api_silent "获取不存在的面试" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/9999" "" "404"

# 6.3 获取面试统计
test_api_silent "获取面试统计" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/stats" "" "200"

# 6.4 获取今日面试
test_api_silent "获取今日面试" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/today" "" "200"

# 6.5 创建面试
test_api_silent "创建新面试" "POST" "$INTERVIEW_SERVICE/api/v1/interviews" \
    '{"candidate_id":1,"candidate_name":"测试候选人","position_id":1,"position":"测试职位","type":"initial","date":"2025-12-30","time":"14:00","duration":60,"interviewer_id":5,"interviewer":"陈强","method":"video","location":"https://meeting.test.com/123","created_by":1}' "201"

# 6.6 获取面试官日程
test_api_silent "获取面试官日程" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/interviewer/5" "" "200"

# 6.7 获取候选人面试记录
test_api_silent "获取候选人面试记录" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/candidate/1" "" "200"

# 6.8 更新面试
test_api_silent "更新面试信息" "PUT" "$INTERVIEW_SERVICE/api/v1/interviews/1" \
    '{"notes":"更新的备注信息"}' "200"

# 6.9 获取面试反馈
test_api_silent "获取面试反馈" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/1/feedback" "" "200"

echo ""

# =====================================================
# 7. 消息服务测试 (message-service)
# =====================================================
echo -e "${CYAN}[7/8] 消息服务测试 (message-service) - 9个测试${NC}"
echo "----------------------------------------"

# 7.1 获取消息列表
test_api_silent "获取消息列表 (user_id=1)" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1" "" "200"

test_api_silent "获取消息列表 (user_id=2)" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=2" "" "200"

test_api_silent "获取消息列表 - 分页" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1&page=1&page_size=5" "" "200"

test_api_silent "获取消息列表 - 按类型筛选" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1&type=system" "" "200"

# 7.2 获取未读消息数
test_api_silent "获取未读消息数 (user_id=1)" "GET" "$MESSAGE_SERVICE/api/v1/messages/unread-count?user_id=1" "" "200"

test_api_silent "获取未读消息数 (user_id=3)" "GET" "$MESSAGE_SERVICE/api/v1/messages/unread-count?user_id=3" "" "200"

# 7.3 发送消息
test_api_silent "发送新消息" "POST" "$MESSAGE_SERVICE/api/v1/messages" \
    '{"sender_id":1,"receiver_id":2,"type":"chat","title":"测试消息","content":"这是一条测试消息"}' "201"

test_api_silent "发送系统消息" "POST" "$MESSAGE_SERVICE/api/v1/messages" \
    '{"receiver_id":3,"type":"system","title":"系统通知","content":"这是一条系统通知"}' "201"

# 7.4 标记消息已读
test_api_silent "标记消息已读" "PUT" "$MESSAGE_SERVICE/api/v1/messages/1/read" "" "200"

echo ""

# =====================================================
# 8. 综合功能测试
# =====================================================
echo -e "${CYAN}[8/8] 综合功能测试 - 3个测试${NC}"
echo "----------------------------------------"

# 8.1 健康检查
test_api_silent "interview-service 健康检查" "GET" "$INTERVIEW_SERVICE/health" "" "200"

# 8.2 应用记录测试
test_api_silent "获取应用记录列表" "GET" "$RESUME_SERVICE/api/v1/applications" "" "200"

test_api_silent "获取应用记录列表 - 分页" "GET" "$RESUME_SERVICE/api/v1/applications?page=1&page_size=5" "" "200"

echo ""

# =====================================================
# 测试结果汇总
# =====================================================
echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}                    测试结果汇总${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""
echo -e "总测试数: ${TOTAL}"
echo -e "${GREEN}通过: ${PASSED}${NC}"
echo -e "${RED}失败: ${FAILED}${NC}"
echo ""

if [ $TOTAL -gt 0 ]; then
    PASS_RATE=$((PASSED * 100 / TOTAL))
    echo -e "通过率: ${PASS_RATE}%"
else
    echo -e "通过率: 0%"
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo ""
    echo "测试覆盖范围（共74个测试用例）:"
    echo "  - 用户服务 (8个): 登录、注册、用户列表"
    echo "  - 职位服务 (14个): CRUD、筛选、搜索、统计"
    echo "  - 人才服务 (12个): CRUD、搜索、筛选"
    echo "  - 简历服务 (14个): 列表、详情、状态更新、排序"
    echo "  - 面试服务 (14个): CRUD、反馈、日程、统计"
    echo "  - 消息服务 (9个): 发送、列表、已读、未读数"
    echo "  - 综合测试 (3个): 健康检查、应用记录"
    echo ""
    echo -e "${YELLOW}注意: Coze工作流评估功能未测试（需要API Key）${NC}"
    echo ""
    exit 0
else
    echo -e "${YELLOW}⚠ 部分测试失败，请检查相关服务${NC}"
    echo ""
    echo "常见问题排查:"
    echo "  1. 确保所有服务已启动: ./start-backend.sh"
    echo "  2. 检查数据库连接是否正常"
    echo "  3. 确认模拟数据已导入: cd database && ./import_mock_data.sh"
    echo ""
    exit 1
fi
