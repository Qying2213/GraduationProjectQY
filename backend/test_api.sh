#!/bin/bash
# =====================================================
# 智能人才运营平台 - 完整API功能测试脚本
# 测试所有后端服务的主要功能（除Coze工作流外）
# 共90个测试用例
# =====================================================

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 服务地址
USER_SERVICE="http://localhost:8081"
JOB_SERVICE="http://localhost:8082"
INTERVIEW_SERVICE="http://localhost:8083"
RESUME_SERVICE="http://localhost:8084"
MESSAGE_SERVICE="http://localhost:8085"
TALENT_SERVICE="http://localhost:8086"
LOG_SERVICE="http://localhost:8088"
ES_SERVICE="http://localhost:9200"

# 计数器
TOTAL=0
PASSED=0
FAILED=0

# 静默测试函数
test_api_silent() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_code="$5"
    
    TOTAL=$((TOTAL + 1))
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$url" \
            -H "Content-Type: application/json" -d "$data" 2>/dev/null)
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
        echo -e "  响应: $(echo "$body" | head -c 100)"
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# 检查服务
check_service() {
    local name="$1"
    local url="$2"
    response=$(curl -s --connect-timeout 2 -o /dev/null -w "%{http_code}" "$url/health" 2>/dev/null)
    if [ "$response" != "000" ] && [ "$response" != "" ]; then
        echo -e "${GREEN}✓${NC} $name 运行中"
        return 0
    fi
    echo -e "${RED}✗${NC} $name 未运行"
    return 1
}

# 检查ES
check_es() {
    response=$(curl -s --connect-timeout 2 -o /dev/null -w "%{http_code}" "$ES_SERVICE" 2>/dev/null)
    if [ "$response" = "200" ]; then
        echo -e "${GREEN}✓${NC} Elasticsearch (9200) 运行中"
        return 0
    fi
    echo -e "${YELLOW}!${NC} Elasticsearch (9200) 未运行 (日志功能降级)"
    return 1
}

echo ""
echo -e "${BLUE}================================================================${NC}"
echo -e "${BLUE}         智能人才运营平台 - 完整API功能测试${NC}"
echo -e "${BLUE}              共90个测试用例（含ES日志功能）${NC}"
echo -e "${BLUE}================================================================${NC}"
echo ""

# 1. 检查服务状态
echo -e "${CYAN}[1/9] 检查服务状态${NC}"
echo "----------------------------------------"
check_service "user-service (8081)" "$USER_SERVICE"
check_service "job-service (8082)" "$JOB_SERVICE"
check_service "interview-service (8083)" "$INTERVIEW_SERVICE"
check_service "resume-service (8084)" "$RESUME_SERVICE"
check_service "message-service (8085)" "$MESSAGE_SERVICE"
check_service "talent-service (8086)" "$TALENT_SERVICE"
check_service "log-service (8088)" "$LOG_SERVICE"
check_es
echo ""

TIMESTAMP=$(date +%s)

# 2. 用户服务测试
echo -e "${CYAN}[2/9] 用户服务测试 - 10个${NC}"
echo "----------------------------------------"
test_api_silent "登录-admin" "POST" "$USER_SERVICE/api/v1/login" '{"username":"admin","password":"password123"}' "200"
test_api_silent "登录-hr_zhang" "POST" "$USER_SERVICE/api/v1/login" '{"username":"hr_zhang","password":"password123"}' "200"
test_api_silent "登录-错误密码" "POST" "$USER_SERVICE/api/v1/login" '{"username":"admin","password":"wrong"}' "401"
test_api_silent "登录-不存在用户" "POST" "$USER_SERVICE/api/v1/login" '{"username":"notexist","password":"123"}' "401"
test_api_silent "注册-新用户" "POST" "$USER_SERVICE/api/v1/register" '{"username":"test_'$TIMESTAMP'","email":"t'$TIMESTAMP'@t.com","password":"test123","role":"hr"}' "201"
test_api_silent "注册-重复用户名" "POST" "$USER_SERVICE/api/v1/register" '{"username":"admin","email":"new@t.com","password":"test123"}' "400"
test_api_silent "注册-重复邮箱" "POST" "$USER_SERVICE/api/v1/register" '{"username":"newu123","email":"admin@company.com","password":"test123"}' "400"
test_api_silent "用户列表" "GET" "$USER_SERVICE/api/v1/users" "" "200"
test_api_silent "用户列表-分页" "GET" "$USER_SERVICE/api/v1/users?page=1&page_size=5" "" "200"
test_api_silent "用户列表-第二页" "GET" "$USER_SERVICE/api/v1/users?page=2&page_size=3" "" "200"
echo ""

# 3. 职位服务测试
echo -e "${CYAN}[3/9] 职位服务测试 - 14个${NC}"
echo "----------------------------------------"
test_api_silent "职位列表" "GET" "$JOB_SERVICE/api/v1/jobs" "" "200"
test_api_silent "职位列表-分页" "GET" "$JOB_SERVICE/api/v1/jobs?page=1&page_size=5" "" "200"
test_api_silent "职位列表-第二页" "GET" "$JOB_SERVICE/api/v1/jobs?page=2&page_size=5" "" "200"
test_api_silent "职位-状态筛选open" "GET" "$JOB_SERVICE/api/v1/jobs?status=open" "" "200"
test_api_silent "职位-状态筛选closed" "GET" "$JOB_SERVICE/api/v1/jobs?status=closed" "" "200"
test_api_silent "职位-地点筛选" "GET" "$JOB_SERVICE/api/v1/jobs?location=北京" "" "200"
test_api_silent "职位-类型筛选" "GET" "$JOB_SERVICE/api/v1/jobs?type=full-time" "" "200"
test_api_silent "职位-关键词搜索" "GET" "$JOB_SERVICE/api/v1/jobs?keyword=前端" "" "200"
test_api_silent "职位详情-ID1" "GET" "$JOB_SERVICE/api/v1/jobs/1" "" "200"
test_api_silent "职位详情-ID2" "GET" "$JOB_SERVICE/api/v1/jobs/2" "" "200"
test_api_silent "职位-不存在" "GET" "$JOB_SERVICE/api/v1/jobs/9999" "" "404"
test_api_silent "职位统计" "GET" "$JOB_SERVICE/api/v1/jobs/stats" "" "200"
test_api_silent "创建职位" "POST" "$JOB_SERVICE/api/v1/jobs" '{"title":"测试职位_'$TIMESTAMP'","description":"测试","location":"上海","salary":"15K","type":"full-time","status":"open","department":"技术部","created_by":1}' "201"
test_api_silent "更新职位" "PUT" "$JOB_SERVICE/api/v1/jobs/1" '{"title":"高级前端(更新)","status":"open"}' "200"
echo ""

# 4. 人才服务测试
echo -e "${CYAN}[4/9] 人才服务测试 - 14个${NC}"
echo "----------------------------------------"
test_api_silent "人才列表" "GET" "$TALENT_SERVICE/api/v1/talents" "" "200"
test_api_silent "人才列表-分页" "GET" "$TALENT_SERVICE/api/v1/talents?page=1&page_size=5" "" "200"
test_api_silent "人才列表-第二页" "GET" "$TALENT_SERVICE/api/v1/talents?page=2&page_size=5" "" "200"
test_api_silent "人才-状态筛选" "GET" "$TALENT_SERVICE/api/v1/talents?status=active" "" "200"
test_api_silent "人才-搜索" "GET" "$TALENT_SERVICE/api/v1/talents?search=张" "" "200"
test_api_silent "人才详情-ID1" "GET" "$TALENT_SERVICE/api/v1/talents/1" "" "200"
test_api_silent "人才详情-ID5" "GET" "$TALENT_SERVICE/api/v1/talents/5" "" "200"
test_api_silent "人才-不存在" "GET" "$TALENT_SERVICE/api/v1/talents/9999" "" "404"
test_api_silent "搜索人才-前端" "GET" "$TALENT_SERVICE/api/v1/talents/search?keyword=前端" "" "200"
test_api_silent "搜索人才-Go" "GET" "$TALENT_SERVICE/api/v1/talents/search?keyword=Go" "" "200"
test_api_silent "搜索人才-经验" "GET" "$TALENT_SERVICE/api/v1/talents/search?min_experience=3&max_experience=8" "" "200"
test_api_silent "搜索人才-地点" "GET" "$TALENT_SERVICE/api/v1/talents/search?location=北京" "" "200"
test_api_silent "创建人才" "POST" "$TALENT_SERVICE/api/v1/talents" '{"name":"测试_'$TIMESTAMP'","email":"t'$TIMESTAMP'@t.com","phone":"13900000099","skills":["Java"],"experience":3,"education":"本科","status":"active","location":"深圳"}' "201"
test_api_silent "更新人才" "PUT" "$TALENT_SERVICE/api/v1/talents/1" '{"status":"active","experience":7}' "200"
echo ""

# 5. 简历服务测试
echo -e "${CYAN}[5/9] 简历服务测试 - 14个${NC}"
echo "----------------------------------------"
test_api_silent "简历列表" "GET" "$RESUME_SERVICE/api/v1/resumes" "" "200"
test_api_silent "简历列表-分页" "GET" "$RESUME_SERVICE/api/v1/resumes?page=1&page_size=5" "" "200"
test_api_silent "简历列表-第二页" "GET" "$RESUME_SERVICE/api/v1/resumes?page=2&page_size=5" "" "200"
test_api_silent "简历-状态pending" "GET" "$RESUME_SERVICE/api/v1/resumes?status=pending" "" "200"
test_api_silent "简历-状态reviewing" "GET" "$RESUME_SERVICE/api/v1/resumes?status=reviewing" "" "200"
test_api_silent "简历-状态interviewed" "GET" "$RESUME_SERVICE/api/v1/resumes?status=interviewed" "" "200"
test_api_silent "简历-排序降序" "GET" "$RESUME_SERVICE/api/v1/resumes?sort_by=created_at&sort_order=desc" "" "200"
test_api_silent "简历-排序升序" "GET" "$RESUME_SERVICE/api/v1/resumes?sort_by=created_at&sort_order=asc" "" "200"
test_api_silent "简历详情-ID1" "GET" "$RESUME_SERVICE/api/v1/resumes/1" "" "200"
test_api_silent "简历详情-ID5" "GET" "$RESUME_SERVICE/api/v1/resumes/5" "" "200"
test_api_silent "简历-不存在" "GET" "$RESUME_SERVICE/api/v1/resumes/9999" "" "404"
test_api_silent "更新简历状态" "PUT" "$RESUME_SERVICE/api/v1/resumes/1/status" '{"status":"reviewing"}' "200"
test_api_silent "AI配置状态" "GET" "$RESUME_SERVICE/api/v1/ai/config" "" "200"
test_api_silent "评估用简历列表" "GET" "$RESUME_SERVICE/api/v1/resumes/evaluation" "" "200"
echo ""

# 6. 面试服务测试
echo -e "${CYAN}[6/9] 面试服务测试 - 16个${NC}"
echo "----------------------------------------"
test_api_silent "面试列表" "GET" "$INTERVIEW_SERVICE/api/v1/interviews" "" "200"
test_api_silent "面试列表-分页" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?page=1&page_size=5" "" "200"
test_api_silent "面试列表-第二页" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?page=2&page_size=5" "" "200"
test_api_silent "面试-状态scheduled" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?status=scheduled" "" "200"
test_api_silent "面试-状态completed" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?status=completed" "" "200"
test_api_silent "面试-状态cancelled" "GET" "$INTERVIEW_SERVICE/api/v1/interviews?status=cancelled" "" "200"
test_api_silent "面试详情-ID1" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/1" "" "200"
test_api_silent "面试详情-ID5" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/5" "" "200"
test_api_silent "面试-不存在" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/9999" "" "404"
test_api_silent "面试统计" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/stats" "" "200"
test_api_silent "今日面试" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/today" "" "200"
test_api_silent "创建面试" "POST" "$INTERVIEW_SERVICE/api/v1/interviews" '{"candidate_id":1,"candidate_name":"测试","position_id":1,"position":"测试职位","type":"initial","date":"2025-12-30","time":"14:00","duration":60,"interviewer_id":5,"interviewer":"陈强","method":"video","location":"https://meet.test/123","created_by":1}' "201"
test_api_silent "面试官日程-ID5" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/interviewer/5" "" "200"
test_api_silent "面试官日程-ID6" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/interviewer/6" "" "200"
test_api_silent "候选人面试-ID1" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/candidate/1" "" "200"
test_api_silent "面试反馈-ID1" "GET" "$INTERVIEW_SERVICE/api/v1/interviews/1/feedback" "" "200"
echo ""

# 7. 消息服务测试
echo -e "${CYAN}[7/9] 消息服务测试 - 10个${NC}"
echo "----------------------------------------"
test_api_silent "消息列表-user1" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1" "" "200"
test_api_silent "消息列表-user2" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=2" "" "200"
test_api_silent "消息列表-分页" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1&page=1&page_size=5" "" "200"
test_api_silent "消息-类型system" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1&type=system" "" "200"
test_api_silent "消息-类型interview" "GET" "$MESSAGE_SERVICE/api/v1/messages?user_id=1&type=interview" "" "200"
test_api_silent "未读数-user1" "GET" "$MESSAGE_SERVICE/api/v1/messages/unread-count?user_id=1" "" "200"
test_api_silent "未读数-user3" "GET" "$MESSAGE_SERVICE/api/v1/messages/unread-count?user_id=3" "" "200"
test_api_silent "发送消息" "POST" "$MESSAGE_SERVICE/api/v1/messages" '{"sender_id":1,"receiver_id":2,"type":"chat","title":"测试","content":"测试消息"}' "201"
test_api_silent "发送系统消息" "POST" "$MESSAGE_SERVICE/api/v1/messages" '{"receiver_id":3,"type":"system","title":"通知","content":"系统通知"}' "201"
test_api_silent "标记已读" "PUT" "$MESSAGE_SERVICE/api/v1/messages/1/read" "" "200"
echo ""

# 8. 日志服务测试 (ES)
echo -e "${CYAN}[8/9] 日志服务测试 (ES) - 8个${NC}"
echo "----------------------------------------"
test_api_silent "日志服务健康检查" "GET" "$LOG_SERVICE/health" "" "200"
test_api_silent "查询所有日志" "GET" "$LOG_SERVICE/api/v1/logs" "" "200"
test_api_silent "日志-分页" "GET" "$LOG_SERVICE/api/v1/logs?page=1&page_size=10" "" "200"
test_api_silent "日志-按服务筛选" "GET" "$LOG_SERVICE/api/v1/logs?service=user-service" "" "200"
test_api_silent "日志-按方法筛选" "GET" "$LOG_SERVICE/api/v1/logs?method=POST" "" "200"
test_api_silent "获取日志统计" "GET" "$LOG_SERVICE/api/v1/logs/stats" "" "200"
test_api_silent "获取服务列表" "GET" "$LOG_SERVICE/api/v1/logs/services" "" "200"
test_api_silent "获取操作类型" "GET" "$LOG_SERVICE/api/v1/logs/actions" "" "200"
echo ""

# 9. 综合功能测试
echo -e "${CYAN}[9/9] 综合功能测试 - 4个${NC}"
echo "----------------------------------------"
test_api_silent "应用记录列表" "GET" "$RESUME_SERVICE/api/v1/applications" "" "200"
test_api_silent "应用记录-分页" "GET" "$RESUME_SERVICE/api/v1/applications?page=1&page_size=5" "" "200"

# 等待日志写入ES
sleep 1

# 验证ES日志记录
ES_COUNT=$(curl -s "$ES_SERVICE/operation_logs/_count" 2>/dev/null | grep -o '"count":[0-9]*' | grep -o '[0-9]*')
if [ -n "$ES_COUNT" ] && [ "$ES_COUNT" -gt 0 ]; then
    TOTAL=$((TOTAL + 1))
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ PASS${NC} - ES日志记录验证 (已记录 $ES_COUNT 条)"
else
    TOTAL=$((TOTAL + 1))
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ FAIL${NC} - ES日志记录验证"
fi

# 验证日志包含多个服务
SERVICES=$(curl -s "$LOG_SERVICE/api/v1/logs" 2>/dev/null | grep -o '"service":"[^"]*"' | sort -u | wc -l)
if [ "$SERVICES" -gt 1 ]; then
    TOTAL=$((TOTAL + 1))
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}✓ PASS${NC} - 多服务日志验证 (记录了 $SERVICES 个服务)"
else
    TOTAL=$((TOTAL + 1))
    FAILED=$((FAILED + 1))
    echo -e "${RED}✗ FAIL${NC} - 多服务日志验证"
fi
echo ""

# 测试结果汇总
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
fi
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    echo ""
    echo "测试覆盖:"
    echo "  用户服务(10) 职位服务(14) 人才服务(14)"
    echo "  简历服务(14) 面试服务(16) 消息服务(10)"
    echo "  日志服务(8)  综合测试(4)"
    echo ""
    echo -e "${YELLOW}注意: Coze工作流评估功能未测试（需要API Key）${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ 部分测试失败，请检查相关服务${NC}"
    echo "排查: 1.启动服务 2.检查数据库 3.导入模拟数据 4.启动ES"
    exit 1
fi
