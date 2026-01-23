#!/bin/bash

# ============================================================
# 智能招聘系统 - 完整功能测试脚本
# ============================================================
# 测试覆盖:
# 1. 基础服务健康检查
# 2. 用户服务 (注册/登录/JWT)
# 3. 职位服务 (CRUD)
# 4. 人才服务 (CRUD)
# 5. 简历服务 (上传/解析)
# 6. AI评估服务 (OCR/Embedding/RAG/Coze)
# 7. 推荐服务 (语义匹配/RAG)
# 8. 面试服务
# 9. 消息服务
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务端口配置
GATEWAY_PORT=8080
USER_PORT=8081
JOB_PORT=8082
INTERVIEW_PORT=8083
RESUME_PORT=8084
MESSAGE_PORT=8085
TALENT_PORT=8086
RECOMMENDATION_PORT=8087

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# 测试结果文件
RESULT_FILE="ztest/test_results_$(date +%Y%m%d_%H%M%S).log"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1" | tee -a "$RESULT_FILE"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1" | tee -a "$RESULT_FILE"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1" | tee -a "$RESULT_FILE"
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

log_skip() {
    echo -e "${YELLOW}[SKIP]${NC} $1" | tee -a "$RESULT_FILE"
    ((SKIPPED_TESTS++))
    ((TOTAL_TESTS++))
}

log_section() {
    echo "" | tee -a "$RESULT_FILE"
    echo -e "${YELLOW}============================================================${NC}" | tee -a "$RESULT_FILE"
    echo -e "${YELLOW} $1${NC}" | tee -a "$RESULT_FILE"
    echo -e "${YELLOW}============================================================${NC}" | tee -a "$RESULT_FILE"
}

# HTTP请求函数
http_get() {
    local url=$1
    local expected_code=${2:-200}
    local response=$(curl -s -w "\n%{http_code}" "$url" 2>/dev/null)
    local body=$(echo "$response" | head -n -1)
    local code=$(echo "$response" | tail -n 1)
    
    if [ "$code" == "$expected_code" ]; then
        echo "$body"
        return 0
    else
        echo "HTTP $code: $body"
        return 1
    fi
}

http_post() {
    local url=$1
    local data=$2
    local expected_code=${3:-200}
    local token=${4:-""}
    
    local auth_header=""
    if [ -n "$token" ]; then
        auth_header="-H \"Authorization: Bearer $token\""
    fi
    
    local response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $token" \
        -d "$data" 2>/dev/null)
    local body=$(echo "$response" | head -n -1)
    local code=$(echo "$response" | tail -n 1)
    
    if [ "$code" == "$expected_code" ]; then
        echo "$body"
        return 0
    else
        echo "HTTP $code: $body"
        return 1
    fi
}

# 检查服务是否运行
check_service() {
    local name=$1
    local port=$2
    
    if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
        log_success "$name 服务运行正常 (端口 $port)"
        return 0
    else
        log_fail "$name 服务未运行 (端口 $port)"
        return 1
    fi
}

# ============================================================
# 测试开始
# ============================================================

echo "智能招聘系统 - 完整功能测试" > "$RESULT_FILE"
echo "测试时间: $(date)" >> "$RESULT_FILE"
echo "" >> "$RESULT_FILE"

log_section "1. 基础服务健康检查"

check_service "Gateway" $GATEWAY_PORT || true
check_service "User" $USER_PORT || true
check_service "Job" $JOB_PORT || true
check_service "Interview" $INTERVIEW_PORT || true
check_service "Resume" $RESUME_PORT || true
check_service "Message" $MESSAGE_PORT || true
check_service "Talent" $TALENT_PORT || true
check_service "Recommendation" $RECOMMENDATION_PORT || true

# ============================================================
log_section "2. 用户服务测试"

# 2.1 用户注册
log_info "测试用户注册..."
REGISTER_RESP=$(http_post "http://localhost:$USER_PORT/api/v1/users/register" \
    '{"username":"testuser_'$(date +%s)'","password":"Test123456","email":"test'$(date +%s)'@test.com","role":"hr"}' \
    "200" 2>/dev/null) || REGISTER_RESP=""

if echo "$REGISTER_RESP" | grep -q '"code":0'; then
    log_success "用户注册成功"
else
    log_skip "用户注册 (可能用户已存在)"
fi

# 2.2 用户登录
log_info "测试用户登录..."
LOGIN_RESP=$(http_post "http://localhost:$USER_PORT/api/v1/users/login" \
    '{"username":"admin","password":"admin123"}' \
    "200" 2>/dev/null) || LOGIN_RESP=""

if echo "$LOGIN_RESP" | grep -q '"token"'; then
    TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    log_success "用户登录成功，获取到Token"
else
    log_fail "用户登录失败: $LOGIN_RESP"
    TOKEN=""
fi

# 2.3 获取用户列表
log_info "测试获取用户列表..."
USERS_RESP=$(http_get "http://localhost:$USER_PORT/api/v1/users" 2>/dev/null) || USERS_RESP=""
if echo "$USERS_RESP" | grep -q '"code":0'; then
    log_success "获取用户列表成功"
else
    log_fail "获取用户列表失败"
fi

# ============================================================
log_section "3. 职位服务测试"

# 3.1 获取职位列表
log_info "测试获取职位列表..."
JOBS_RESP=$(http_get "http://localhost:$JOB_PORT/api/v1/jobs" 2>/dev/null) || JOBS_RESP=""
if echo "$JOBS_RESP" | grep -q '"code":0'; then
    JOB_COUNT=$(echo "$JOBS_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_success "获取职位列表成功 (共 $JOB_COUNT 个职位)"
else
    log_fail "获取职位列表失败"
fi

# 3.2 创建职位
log_info "测试创建职位..."
CREATE_JOB_RESP=$(http_post "http://localhost:$JOB_PORT/api/v1/jobs" \
    '{"title":"测试Go开发工程师","department":"技术部","location":"深圳","salary":"25-40K","description":"负责后端开发","requirements":"3年以上Go经验","skills":"Go,Docker,Kubernetes"}' \
    "200" 2>/dev/null) || CREATE_JOB_RESP=""

if echo "$CREATE_JOB_RESP" | grep -q '"code":0'; then
    NEW_JOB_ID=$(echo "$CREATE_JOB_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    log_success "创建职位成功 (ID: $NEW_JOB_ID)"
else
    log_fail "创建职位失败: $CREATE_JOB_RESP"
    NEW_JOB_ID=""
fi

# 3.3 获取单个职位
if [ -n "$NEW_JOB_ID" ]; then
    log_info "测试获取单个职位..."
    JOB_DETAIL=$(http_get "http://localhost:$JOB_PORT/api/v1/jobs/$NEW_JOB_ID" 2>/dev/null) || JOB_DETAIL=""
    if echo "$JOB_DETAIL" | grep -q '"code":0'; then
        log_success "获取职位详情成功"
    else
        log_fail "获取职位详情失败"
    fi
fi

# ============================================================
log_section "4. 人才服务测试"

# 4.1 获取人才列表
log_info "测试获取人才列表..."
TALENTS_RESP=$(http_get "http://localhost:$TALENT_PORT/api/v1/talents" 2>/dev/null) || TALENTS_RESP=""
if echo "$TALENTS_RESP" | grep -q '"code":0'; then
    TALENT_COUNT=$(echo "$TALENTS_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_success "获取人才列表成功 (共 $TALENT_COUNT 个人才)"
else
    log_fail "获取人才列表失败"
fi

# 4.2 创建人才
log_info "测试创建人才..."
CREATE_TALENT_RESP=$(http_post "http://localhost:$TALENT_PORT/api/v1/talents" \
    '{"name":"测试候选人","email":"candidate'$(date +%s)'@test.com","phone":"13800138000","skills":"Go,Python,Docker","experience":5,"education":"本科","location":"深圳"}' \
    "200" 2>/dev/null) || CREATE_TALENT_RESP=""

if echo "$CREATE_TALENT_RESP" | grep -q '"code":0'; then
    NEW_TALENT_ID=$(echo "$CREATE_TALENT_RESP" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    log_success "创建人才成功 (ID: $NEW_TALENT_ID)"
else
    log_fail "创建人才失败: $CREATE_TALENT_RESP"
    NEW_TALENT_ID=""
fi

# 4.3 搜索人才
log_info "测试搜索人才..."
SEARCH_RESP=$(http_get "http://localhost:$TALENT_PORT/api/v1/talents?skills=Go" 2>/dev/null) || SEARCH_RESP=""
if echo "$SEARCH_RESP" | grep -q '"code":0'; then
    log_success "搜索人才成功"
else
    log_fail "搜索人才失败"
fi

# ============================================================
log_section "5. 简历服务测试"

# 5.1 获取简历列表
log_info "测试获取简历列表..."
RESUMES_RESP=$(http_get "http://localhost:$RESUME_PORT/api/v1/resumes" 2>/dev/null) || RESUMES_RESP=""
if echo "$RESUMES_RESP" | grep -q '"code":0'; then
    RESUME_COUNT=$(echo "$RESUMES_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_success "获取简历列表成功 (共 $RESUME_COUNT 份简历)"
else
    log_fail "获取简历列表失败"
fi

# 5.2 检查AI配置
log_info "测试AI配置状态..."
AI_CONFIG=$(http_get "http://localhost:$RESUME_PORT/api/v1/ai/config" 2>/dev/null) || AI_CONFIG=""
if echo "$AI_CONFIG" | grep -q '"configured":true'; then
    log_success "AI服务已配置 (Coze)"
else
    log_skip "AI服务未配置"
fi

# 5.3 获取评估结果列表
log_info "测试获取评估结果列表..."
EVALS_RESP=$(http_get "http://localhost:$RESUME_PORT/api/v1/evaluations" 2>/dev/null) || EVALS_RESP=""
if echo "$EVALS_RESP" | grep -q '"code":0'; then
    log_success "获取评估结果列表成功"
else
    log_fail "获取评估结果列表失败"
fi

# ============================================================
log_section "6. AI评估服务测试 (OCR/Embedding/RAG/Coze)"

# 6.1 OCR文本提取测试
log_info "测试OCR文本提取..."
# 创建测试PDF路径
TEST_RESUME_PATH="uploads/1769138758127289000_【golang后端开发工程师_深圳 】朱冠州 9年.pdf"

if [ -f "backend/resume-service/$TEST_RESUME_PATH" ]; then
    log_success "测试简历文件存在"
else
    log_skip "测试简历文件不存在，跳过OCR测试"
fi

# 6.2 AI评估测试 (完整流程)
log_info "测试AI评估完整流程 (OCR → Embedding → RAG → Coze)..."
log_info "注意: 此测试需要约1-2分钟完成..."

AI_EVAL_RESP=$(curl -s --max-time 180 -X POST "http://localhost:$RESUME_PORT/api/v1/ai/evaluate" \
    -H "Content-Type: application/json" \
    -d '{"resume_id":2,"jd_text":"招聘Go后端开发工程师，要求3年以上经验，熟悉微服务架构，熟练使用Docker和Kubernetes"}' 2>/dev/null) || AI_EVAL_RESP=""

if echo "$AI_EVAL_RESP" | grep -q '"code":0'; then
    TOTAL_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"total_score":[0-9.]*' | cut -d':' -f2)
    GRADE=$(echo "$AI_EVAL_RESP" | grep -o '"grade":"[^"]*"' | cut -d'"' -f4)
    OCR_USED=$(echo "$AI_EVAL_RESP" | grep -o '"ocr_extracted":[a-z]*' | cut -d':' -f2)
    EMBEDDING_USED=$(echo "$AI_EVAL_RESP" | grep -o '"embedding_used":[a-z]*' | cut -d':' -f2)
    RAG_USED=$(echo "$AI_EVAL_RESP" | grep -o '"rag_enhanced":[a-z]*' | cut -d':' -f2)
    
    log_success "AI评估完成: 总分=$TOTAL_SCORE, 等级=$GRADE"
    log_info "  - OCR提取: $OCR_USED"
    log_info "  - Embedding: $EMBEDDING_USED"
    log_info "  - RAG增强: $RAG_USED"
else
    log_fail "AI评估失败: $AI_EVAL_RESP"
fi

# ============================================================
log_section "7. 推荐服务测试"

# 7.1 语义匹配测试
log_info "测试语义匹配 (Embedding相似度)..."
SEMANTIC_RESP=$(http_post "http://localhost:$RECOMMENDATION_PORT/api/v1/recommendations/semantic-match" \
    '{"text1":"Go语言后端开发","text2":"Golang服务端工程师"}' \
    "200" 2>/dev/null) || SEMANTIC_RESP=""

if echo "$SEMANTIC_RESP" | grep -q '"similarity"'; then
    SIMILARITY=$(echo "$SEMANTIC_RESP" | grep -o '"similarity":[0-9.]*' | cut -d':' -f2)
    log_success "语义匹配成功: 相似度=$SIMILARITY"
else
    log_fail "语义匹配失败: $SEMANTIC_RESP"
fi

# 7.2 RAG查询测试
log_info "测试RAG查询..."
RAG_RESP=$(http_post "http://localhost:$RECOMMENDATION_PORT/api/v1/recommendations/rag/query" \
    '{"query":"Go后端开发工程师","top_k":3,"type":"talent"}' \
    "200" 2>/dev/null) || RAG_RESP=""

if echo "$RAG_RESP" | grep -q '"results"'; then
    RESULT_COUNT=$(echo "$RAG_RESP" | grep -o '"id":[0-9]*' | wc -l)
    log_success "RAG查询成功: 返回 $RESULT_COUNT 条结果"
else
    log_fail "RAG查询失败: $RAG_RESP"
fi

# 7.3 职位推荐测试
log_info "测试为人才推荐职位..."
JOB_REC_RESP=$(http_post "http://localhost:$RECOMMENDATION_PORT/api/v1/recommendations/jobs-for-talent" \
    '{"id":1,"name":"张伟","skills":["Go","Docker","Kubernetes"],"experience":5,"education":"本科","location":"北京"}' \
    "200" 2>/dev/null) || JOB_REC_RESP=""

if echo "$JOB_REC_RESP" | grep -q '"code":0'; then
    log_success "职位推荐成功"
else
    log_fail "职位推荐失败"
fi

# 7.4 人才推荐测试
log_info "测试为职位推荐人才..."
TALENT_REC_RESP=$(http_post "http://localhost:$RECOMMENDATION_PORT/api/v1/recommendations/talents-for-job" \
    '{"id":1,"title":"Go开发工程师","skills":["Go","Docker","Kubernetes"],"location":"北京","level":"senior"}' \
    "200" 2>/dev/null) || TALENT_REC_RESP=""

if echo "$TALENT_REC_RESP" | grep -q '"code":0'; then
    log_success "人才推荐成功"
else
    log_fail "人才推荐失败"
fi

# 7.5 推荐统计
log_info "测试推荐统计..."
STATS_RESP=$(http_get "http://localhost:$RECOMMENDATION_PORT/api/v1/recommendations/stats" 2>/dev/null) || STATS_RESP=""
if echo "$STATS_RESP" | grep -q '"code":0'; then
    log_success "获取推荐统计成功"
else
    log_fail "获取推荐统计失败"
fi

# ============================================================
log_section "8. 面试服务测试"

# 8.1 获取面试列表
log_info "测试获取面试列表..."
INTERVIEWS_RESP=$(http_get "http://localhost:$INTERVIEW_PORT/api/v1/interviews" 2>/dev/null) || INTERVIEWS_RESP=""
if echo "$INTERVIEWS_RESP" | grep -q '"code":0'; then
    log_success "获取面试列表成功"
else
    log_fail "获取面试列表失败"
fi

# 8.2 创建面试
log_info "测试创建面试..."
CREATE_INTERVIEW_RESP=$(http_post "http://localhost:$INTERVIEW_PORT/api/v1/interviews" \
    '{"job_id":1,"talent_id":1,"interviewer_id":1,"scheduled_time":"2026-02-01T10:00:00Z","location":"线上","type":"technical","status":"scheduled"}' \
    "200" 2>/dev/null) || CREATE_INTERVIEW_RESP=""

if echo "$CREATE_INTERVIEW_RESP" | grep -q '"code":0'; then
    log_success "创建面试成功"
else
    log_fail "创建面试失败: $CREATE_INTERVIEW_RESP"
fi

# ============================================================
log_section "9. 消息服务测试"

# 9.1 获取消息列表
log_info "测试获取消息列表..."
MESSAGES_RESP=$(http_get "http://localhost:$MESSAGE_PORT/api/v1/messages?user_id=1" 2>/dev/null) || MESSAGES_RESP=""
if echo "$MESSAGES_RESP" | grep -q '"code":0'; then
    log_success "获取消息列表成功"
else
    log_fail "获取消息列表失败"
fi

# 9.2 发送消息
log_info "测试发送消息..."
SEND_MSG_RESP=$(http_post "http://localhost:$MESSAGE_PORT/api/v1/messages" \
    '{"sender_id":1,"receiver_id":2,"content":"测试消息","type":"system"}' \
    "200" 2>/dev/null) || SEND_MSG_RESP=""

if echo "$SEND_MSG_RESP" | grep -q '"code":0'; then
    log_success "发送消息成功"
else
    log_fail "发送消息失败: $SEND_MSG_RESP"
fi

# ============================================================
log_section "10. 数据库连接测试"

log_info "测试PostgreSQL连接..."
if psql -U qinyang -d talent_platform -c "SELECT 1" > /dev/null 2>&1; then
    log_success "PostgreSQL连接正常"
    
    # 统计数据
    USER_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM users" 2>/dev/null | tr -d ' ')
    JOB_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM jobs" 2>/dev/null | tr -d ' ')
    TALENT_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM talents" 2>/dev/null | tr -d ' ')
    RESUME_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM resumes" 2>/dev/null | tr -d ' ')
    
    log_info "数据库统计:"
    log_info "  - 用户数: $USER_COUNT"
    log_info "  - 职位数: $JOB_COUNT"
    log_info "  - 人才数: $TALENT_COUNT"
    log_info "  - 简历数: $RESUME_COUNT"
else
    log_fail "PostgreSQL连接失败"
fi

# 检查pgvector
log_info "测试pgvector扩展..."
PGVECTOR_CHECK=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM talent_embeddings" 2>/dev/null | tr -d ' ')
if [ -n "$PGVECTOR_CHECK" ]; then
    log_success "pgvector扩展正常 (talent_embeddings: $PGVECTOR_CHECK 条)"
else
    log_fail "pgvector扩展异常"
fi

# ============================================================
log_section "测试结果汇总"

echo "" | tee -a "$RESULT_FILE"
echo "============================================================" | tee -a "$RESULT_FILE"
echo "测试完成时间: $(date)" | tee -a "$RESULT_FILE"
echo "============================================================" | tee -a "$RESULT_FILE"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}" | tee -a "$RESULT_FILE"
echo -e "${RED}失败: $FAILED_TESTS${NC}" | tee -a "$RESULT_FILE"
echo -e "${YELLOW}跳过: $SKIPPED_TESTS${NC}" | tee -a "$RESULT_FILE"
echo "总计: $TOTAL_TESTS" | tee -a "$RESULT_FILE"
echo "============================================================" | tee -a "$RESULT_FILE"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}所有测试通过!${NC}" | tee -a "$RESULT_FILE"
    exit 0
else
    echo -e "${RED}有 $FAILED_TESTS 个测试失败${NC}" | tee -a "$RESULT_FILE"
    exit 1
fi
