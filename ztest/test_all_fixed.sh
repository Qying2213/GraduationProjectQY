#!/bin/bash

# ============================================================
# 智能招聘系统 - 完整功能测试脚本 (修复版)
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 结果文件
RESULT_FILE="ztest/test_results_$(date +%Y%m%d_%H%M%S).log"

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1" | tee -a "$RESULT_FILE"
    ((PASSED_TESTS++))
    ((TOTAL_TESTS++))
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1" | tee -a "$RESULT_FILE"
    ((FAILED_TESTS++))
    ((TOTAL_TESTS++))
}

log_section() {
    echo "" | tee -a "$RESULT_FILE"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" | tee -a "$RESULT_FILE"
    echo -e "${CYAN} $1${NC}" | tee -a "$RESULT_FILE"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}" | tee -a "$RESULT_FILE"
}

echo "智能招聘系统 - 完整功能测试" > "$RESULT_FILE"
echo "测试时间: $(date)" >> "$RESULT_FILE"

# ============================================================
log_section "1. 服务健康检查"

for port in 8080 8081 8082 8083 8084 8085 8086 8087; do
    if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
        SERVICE=$(curl -s "http://localhost:$port/health" | grep -o '"service":"[^"]*"' | cut -d'"' -f4)
        log_pass "端口 $port ($SERVICE) - 运行中"
    else
        log_fail "端口 $port - 未运行"
    fi
done

# ============================================================
log_section "2. 用户服务测试"

# 登录 (正确路由: /api/v1/login)
LOGIN_RESP=$(curl -s -X POST "http://localhost:8081/api/v1/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')
if echo "$LOGIN_RESP" | grep -q '"token"'; then
    TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    log_pass "用户登录成功"
else
    log_fail "用户登录失败"
    TOKEN=""
fi

# 获取用户列表
USERS_RESP=$(curl -s "http://localhost:8081/api/v1/users")
if echo "$USERS_RESP" | grep -q '"code":0'; then
    TOTAL=$(echo "$USERS_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_pass "获取用户列表成功 (共 $TOTAL 用户)"
else
    log_fail "获取用户列表失败"
fi

# ============================================================
log_section "3. 职位服务测试"

# 获取职位列表
JOBS_RESP=$(curl -s "http://localhost:8082/api/v1/jobs")
if echo "$JOBS_RESP" | grep -q '"code":0'; then
    TOTAL=$(echo "$JOBS_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_pass "获取职位列表成功 (共 $TOTAL 职位)"
else
    log_fail "获取职位列表失败"
fi

# 创建职位
CREATE_JOB=$(curl -s -X POST "http://localhost:8082/api/v1/jobs" \
    -H "Content-Type: application/json" \
    -d '{"title":"测试职位_'$(date +%s)'","department":"技术部","location":"深圳","salary":"20-30K","description":"测试","requirements":"测试","skills":"Go,Docker"}')
if echo "$CREATE_JOB" | grep -q '"code":0'; then
    log_pass "创建职位成功"
else
    log_fail "创建职位失败"
fi

# 搜索职位
SEARCH_JOBS=$(curl -s "http://localhost:8082/api/v1/jobs?skills=Go")
if echo "$SEARCH_JOBS" | grep -q '"code":0'; then
    log_pass "搜索职位成功"
else
    log_fail "搜索职位失败"
fi

# ============================================================
log_section "4. 人才服务测试"

# 获取人才列表
TALENTS_RESP=$(curl -s "http://localhost:8086/api/v1/talents")
if echo "$TALENTS_RESP" | grep -q '"code":0'; then
    TOTAL=$(echo "$TALENTS_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_pass "获取人才列表成功 (共 $TOTAL 人才)"
else
    log_fail "获取人才列表失败"
fi

# 创建人才
CREATE_TALENT=$(curl -s -X POST "http://localhost:8086/api/v1/talents" \
    -H "Content-Type: application/json" \
    -d '{"name":"测试候选人_'$(date +%s)'","email":"test'$(date +%s)'@test.com","phone":"13800138000","skills":"Go,Python","experience":5,"education":"本科","location":"深圳"}')
if echo "$CREATE_TALENT" | grep -q '"code":0'; then
    log_pass "创建人才成功"
else
    log_fail "创建人才失败"
fi

# ============================================================
log_section "5. 简历服务测试"

# 获取简历列表
RESUMES_RESP=$(curl -s "http://localhost:8084/api/v1/resumes")
if echo "$RESUMES_RESP" | grep -q '"code":0'; then
    TOTAL=$(echo "$RESUMES_RESP" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_pass "获取简历列表成功 (共 $TOTAL 简历)"
else
    log_fail "获取简历列表失败"
fi

# AI配置检查
AI_CONFIG=$(curl -s "http://localhost:8084/api/v1/ai/config")
if echo "$AI_CONFIG" | grep -q '"configured":true'; then
    log_pass "AI服务已配置 (Coze)"
else
    log_fail "AI服务未配置"
fi

# 获取评估结果
EVALS_RESP=$(curl -s "http://localhost:8084/api/v1/evaluations")
if echo "$EVALS_RESP" | grep -q '"code":0'; then
    log_pass "获取评估结果列表成功"
else
    log_fail "获取评估结果列表失败"
fi

# ============================================================
log_section "6. 推荐服务测试"

# 语义匹配
SEMANTIC=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/semantic-match" \
    -H "Content-Type: application/json" \
    -d '{"text1":"Go语言后端开发","text2":"Golang服务端工程师"}')
if echo "$SEMANTIC" | grep -q '"similarity"'; then
    SIM=$(echo "$SEMANTIC" | grep -o '"similarity":[0-9.]*' | cut -d':' -f2)
    log_pass "语义匹配成功 (相似度: $SIM)"
else
    log_fail "语义匹配失败"
fi

# RAG查询
RAG=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/rag/query" \
    -H "Content-Type: application/json" \
    -d '{"query":"Go后端开发","top_k":3,"type":"talent"}')
if echo "$RAG" | grep -q '"results"'; then
    log_pass "RAG查询成功"
else
    log_fail "RAG查询失败"
fi

# 职位推荐
JOB_REC=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/jobs-for-talent" \
    -H "Content-Type: application/json" \
    -d '{"id":1,"name":"张伟","skills":["Go","Docker"],"experience":5,"education":"本科","location":"北京"}')
if echo "$JOB_REC" | grep -q '"code":0'; then
    log_pass "职位推荐成功"
else
    log_fail "职位推荐失败"
fi

# 人才推荐
TALENT_REC=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/talents-for-job" \
    -H "Content-Type: application/json" \
    -d '{"id":1,"title":"Go开发","skills":["Go","Docker"],"location":"北京","level":"senior"}')
if echo "$TALENT_REC" | grep -q '"code":0'; then
    log_pass "人才推荐成功"
else
    log_fail "人才推荐失败"
fi

# 推荐统计
REC_STATS=$(curl -s "http://localhost:8087/api/v1/recommendations/stats")
if echo "$REC_STATS" | grep -q '"code":0'; then
    log_pass "推荐统计成功"
else
    log_fail "推荐统计失败"
fi

# ============================================================
log_section "7. 面试服务测试"

# 获取面试列表
INTERVIEWS=$(curl -s "http://localhost:8083/api/v1/interviews")
if echo "$INTERVIEWS" | grep -q '"code":0'; then
    TOTAL=$(echo "$INTERVIEWS" | grep -o '"total":[0-9]*' | cut -d':' -f2)
    log_pass "获取面试列表成功 (共 $TOTAL 面试)"
else
    log_fail "获取面试列表失败"
fi

# 创建面试
CREATE_INTERVIEW=$(curl -s -X POST "http://localhost:8083/api/v1/interviews" \
    -H "Content-Type: application/json" \
    -d '{"job_id":1,"talent_id":1,"interviewer_id":1,"scheduled_time":"2026-02-01T10:00:00Z","location":"线上","type":"technical","status":"scheduled"}')
if echo "$CREATE_INTERVIEW" | grep -q '"code":0'; then
    log_pass "创建面试成功"
else
    log_fail "创建面试失败"
fi

# ============================================================
log_section "8. 消息服务测试"

# 获取消息列表
MESSAGES=$(curl -s "http://localhost:8085/api/v1/messages?user_id=1")
if echo "$MESSAGES" | grep -q '"code":0'; then
    log_pass "获取消息列表成功"
else
    log_fail "获取消息列表失败"
fi

# 发送消息
SEND_MSG=$(curl -s -X POST "http://localhost:8085/api/v1/messages" \
    -H "Content-Type: application/json" \
    -d '{"sender_id":1,"receiver_id":2,"content":"测试消息","type":"system"}')
if echo "$SEND_MSG" | grep -q '"code":0'; then
    log_pass "发送消息成功"
else
    log_fail "发送消息失败"
fi

# ============================================================
log_section "9. 数据库测试"

# PostgreSQL连接
if psql -U qinyang -d talent_platform -c "SELECT 1" > /dev/null 2>&1; then
    log_pass "PostgreSQL连接正常"
else
    log_fail "PostgreSQL连接失败"
fi

# pgvector扩展
if psql -U qinyang -d talent_platform -c "SELECT COUNT(*) FROM talent_embeddings" > /dev/null 2>&1; then
    COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM talent_embeddings" | tr -d ' ')
    log_pass "pgvector扩展正常 (talent_embeddings: $COUNT 条)"
else
    log_fail "pgvector扩展异常"
fi

# ============================================================
log_section "10. AI评估流程测试 (耗时约2分钟)"

echo "正在执行AI评估..."
AI_EVAL=$(curl -s --max-time 300 -X POST "http://localhost:8084/api/v1/ai/evaluate" \
    -H "Content-Type: application/json" \
    -d '{"resume_id":2,"jd_text":"招聘Go后端开发工程师，要求3年以上经验"}')

if echo "$AI_EVAL" | grep -q '"code":0'; then
    SCORE=$(echo "$AI_EVAL" | grep -o '"total_score":[0-9.]*' | cut -d':' -f2)
    GRADE=$(echo "$AI_EVAL" | grep -o '"grade":"[^"]*"' | cut -d'"' -f4)
    OCR=$(echo "$AI_EVAL" | grep -o '"ocr_extracted":[a-z]*' | cut -d':' -f2)
    log_pass "AI评估成功 (总分: $SCORE, 等级: $GRADE, OCR: $OCR)"
else
    log_fail "AI评估失败"
fi

# ============================================================
log_section "测试结果汇总"

echo "" | tee -a "$RESULT_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULT_FILE"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}" | tee -a "$RESULT_FILE"
echo -e "${RED}失败: $FAILED_TESTS${NC}" | tee -a "$RESULT_FILE"
echo "总计: $TOTAL_TESTS" | tee -a "$RESULT_FILE"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" | tee -a "$RESULT_FILE"

PASS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
echo "通过率: $PASS_RATE%" | tee -a "$RESULT_FILE"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过!${NC}" | tee -a "$RESULT_FILE"
    exit 0
else
    echo -e "${YELLOW}有 $FAILED_TESTS 个测试失败${NC}" | tee -a "$RESULT_FILE"
    exit 1
fi
