#!/bin/bash

# ============================================================
# 智能招聘系统 - 完整功能测试
# ============================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

PASS=0
FAIL=0

test_api() {
    local name=$1
    local method=$2
    local url=$3
    local data=$4
    local expect=$5
    
    if [ "$method" == "GET" ]; then
        resp=$(curl -s --max-time 10 "$url")
    else
        resp=$(curl -s --max-time 10 -X POST "$url" -H "Content-Type: application/json" -d "$data")
    fi
    
    if echo "$resp" | grep -q "$expect"; then
        echo -e "${GREEN}[PASS]${NC} $name"
        ((PASS++))
    else
        echo -e "${RED}[FAIL]${NC} $name"
        ((FAIL++))
    fi
}

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           智能招聘系统 - 完整功能测试                        ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# ============================================================
echo -e "\n${CYAN}[1] 服务健康检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "Gateway (8080)" "GET" "http://localhost:8080/health" "" "healthy"
test_api "User Service (8081)" "GET" "http://localhost:8081/health" "" "healthy"
test_api "Job Service (8082)" "GET" "http://localhost:8082/health" "" "healthy"
test_api "Interview Service (8083)" "GET" "http://localhost:8083/health" "" "healthy"
test_api "Resume Service (8084)" "GET" "http://localhost:8084/health" "" "healthy"
test_api "Message Service (8085)" "GET" "http://localhost:8085/health" "" "healthy"
test_api "Talent Service (8086)" "GET" "http://localhost:8086/health" "" "healthy"
test_api "Recommendation Service (8087)" "GET" "http://localhost:8087/health" "" "healthy"

# ============================================================
echo -e "\n${CYAN}[2] 用户服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "用户登录" "POST" "http://localhost:8081/api/v1/login" '{"username":"admin","password":"admin123"}' "token"
test_api "获取用户列表" "GET" "http://localhost:8081/api/v1/users" "" '"code":0'

# ============================================================
echo -e "\n${CYAN}[3] 职位服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "获取职位列表" "GET" "http://localhost:8082/api/v1/jobs" "" '"code":0'
test_api "创建职位" "POST" "http://localhost:8082/api/v1/jobs" '{"title":"测试职位","department":"技术部","location":"深圳","salary":"20K","description":"测试","requirements":"测试","skills":"Go"}' '"code":0'
test_api "搜索职位" "GET" "http://localhost:8082/api/v1/jobs?skills=Go" "" '"code":0'

# ============================================================
echo -e "\n${CYAN}[4] 人才服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "获取人才列表" "GET" "http://localhost:8086/api/v1/talents" "" '"code":0'
test_api "创建人才" "POST" "http://localhost:8086/api/v1/talents" '{"name":"测试人才","email":"test@test.com","phone":"13800138000","skills":"Go,Python","experience":5,"education":"本科","location":"深圳"}' '"code":0'

# ============================================================
echo -e "\n${CYAN}[5] 简历服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "获取简历列表" "GET" "http://localhost:8084/api/v1/resumes" "" '"code":0'
test_api "AI配置检查" "GET" "http://localhost:8084/api/v1/ai/config" "" '"configured":true'
test_api "获取评估结果" "GET" "http://localhost:8084/api/v1/evaluations" "" '"code":0'

# ============================================================
echo -e "\n${CYAN}[6] 推荐服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "语义匹配(Embedding)" "POST" "http://localhost:8087/api/v1/recommendations/semantic-match" '{"text1":"Go开发","text2":"Golang工程师"}' '"similarity"'
test_api "RAG查询(pgvector)" "POST" "http://localhost:8087/api/v1/recommendations/rag/query" '{"query":"Go后端开发","top_k":3,"type":"talent"}' '"results"'
test_api "职位推荐" "POST" "http://localhost:8087/api/v1/recommendations/jobs-for-talent" '{"id":1,"name":"张伟","skills":["Go"],"experience":5,"education":"本科","location":"北京"}' '"code":0'
test_api "人才推荐" "POST" "http://localhost:8087/api/v1/recommendations/talents-for-job" '{"id":1,"title":"Go开发","skills":["Go"],"location":"北京","level":"senior"}' '"code":0'
test_api "推荐统计" "GET" "http://localhost:8087/api/v1/recommendations/stats" "" '"code":0'

# ============================================================
echo -e "\n${CYAN}[7] 面试服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "获取面试列表" "GET" "http://localhost:8083/api/v1/interviews" "" '"code":0'
test_api "创建面试" "POST" "http://localhost:8083/api/v1/interviews" '{"job_id":1,"talent_id":1,"interviewer_id":1,"scheduled_time":"2026-02-01T10:00:00Z","location":"线上","type":"technical","status":"scheduled"}' '"code":0'

# ============================================================
echo -e "\n${CYAN}[8] 消息服务${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

test_api "获取消息列表" "GET" "http://localhost:8085/api/v1/messages?user_id=1" "" '"code":0'
test_api "发送消息" "POST" "http://localhost:8085/api/v1/messages" '{"sender_id":1,"receiver_id":2,"content":"测试消息","type":"system"}' '"code":0'

# ============================================================
echo -e "\n${CYAN}[9] 数据库检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if psql -U qinyang -d talent_platform -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}[PASS]${NC} PostgreSQL连接"
    ((PASS++))
else
    echo -e "${RED}[FAIL]${NC} PostgreSQL连接"
    ((FAIL++))
fi

if psql -U qinyang -d talent_platform -c "SELECT COUNT(*) FROM talent_embeddings" > /dev/null 2>&1; then
    COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM talent_embeddings" | tr -d ' ')
    echo -e "${GREEN}[PASS]${NC} pgvector扩展 (向量数: $COUNT)"
    ((PASS++))
else
    echo -e "${RED}[FAIL]${NC} pgvector扩展"
    ((FAIL++))
fi

# ============================================================
echo -e "\n${CYAN}[10] AI评估流程测试 (OCR→Embedding→RAG→Coze)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${YELLOW}注意: 此测试需要1-2分钟...${NC}"

AI_RESP=$(curl -s --max-time 180 -X POST "http://localhost:8084/api/v1/ai/evaluate" \
    -H "Content-Type: application/json" \
    -d '{"resume_id":2,"jd_text":"招聘Go后端开发工程师，要求3年以上经验，熟悉微服务架构"}')

if echo "$AI_RESP" | grep -q '"code":0'; then
    SCORE=$(echo "$AI_RESP" | grep -o '"total_score":[0-9.]*' | cut -d':' -f2)
    GRADE=$(echo "$AI_RESP" | grep -o '"grade":"[^"]*"' | cut -d'"' -f4)
    OCR=$(echo "$AI_RESP" | grep -o '"ocr_extracted":[a-z]*' | cut -d':' -f2)
    echo -e "${GREEN}[PASS]${NC} AI评估完成 (总分: $SCORE, 等级: $GRADE, OCR: $OCR)"
    ((PASS++))
else
    echo -e "${RED}[FAIL]${NC} AI评估失败"
    ((FAIL++))
fi

# ============================================================
echo -e "\n${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${CYAN}                      测试结果汇总${NC}"
echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

TOTAL=$((PASS + FAIL))
RATE=$((PASS * 100 / TOTAL))

echo ""
echo -e "  ${GREEN}通过: $PASS${NC}"
echo -e "  ${RED}失败: $FAIL${NC}"
echo -e "  总计: $TOTAL"
echo -e "  通过率: ${RATE}%"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过!${NC}"
else
    echo -e "${YELLOW}有 $FAIL 个测试失败${NC}"
fi
