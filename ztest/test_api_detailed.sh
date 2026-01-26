#!/bin/bash

# ============================================================
# 智能招聘系统 - API详细测试脚本
# ============================================================
# 更详细的API测试，包含请求和响应的完整输出
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# 输出目录
OUTPUT_DIR="ztest/output"
mkdir -p "$OUTPUT_DIR"

# 日志函数
log_test() {
    echo ""
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}测试: $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

log_request() {
    echo -e "${BLUE}请求:${NC} $1 $2"
    if [ -n "$3" ]; then
        echo -e "${BLUE}数据:${NC}"
        echo "$3" | python3 -m json.tool 2>/dev/null || echo "$3"
    fi
}

log_response() {
    echo -e "${GREEN}响应:${NC}"
    echo "$1" | python3 -m json.tool 2>/dev/null || echo "$1"
}

# ============================================================
# 1. 用户服务详细测试
# ============================================================

log_test "1.1 用户登录 - 获取JWT Token"
log_request "POST" "http://localhost:8081/api/v1/users/login" '{"username":"admin","password":"admin123"}'

LOGIN_RESP=$(curl -s -X POST "http://localhost:8081/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')
log_response "$LOGIN_RESP"

TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo -e "${YELLOW}提取的Token: ${TOKEN:0:50}...${NC}"

# ============================================================
log_test "1.2 获取用户列表"
log_request "GET" "http://localhost:8081/api/v1/users"

USERS_RESP=$(curl -s "http://localhost:8081/api/v1/users")
log_response "$USERS_RESP"

# ============================================================
log_test "1.3 获取单个用户"
log_request "GET" "http://localhost:8081/api/v1/users/1"

USER_RESP=$(curl -s "http://localhost:8081/api/v1/users/1")
log_response "$USER_RESP"

# ============================================================
# 2. 职位服务详细测试
# ============================================================

log_test "2.1 获取职位列表"
log_request "GET" "http://localhost:8082/api/v1/jobs"

JOBS_RESP=$(curl -s "http://localhost:8082/api/v1/jobs")
log_response "$JOBS_RESP"

# ============================================================
log_test "2.2 创建新职位"
JOB_DATA='{
    "title": "高级Go开发工程师",
    "department": "技术部",
    "location": "深圳",
    "salary": "30-50K",
    "description": "负责核心系统后端开发，参与架构设计",
    "requirements": "5年以上Go开发经验，熟悉微服务架构",
    "skills": "Go,gRPC,Kubernetes,Docker,Redis,PostgreSQL"
}'
log_request "POST" "http://localhost:8082/api/v1/jobs" "$JOB_DATA"

CREATE_JOB_RESP=$(curl -s -X POST "http://localhost:8082/api/v1/jobs" \
    -H "Content-Type: application/json" \
    -d "$JOB_DATA")
log_response "$CREATE_JOB_RESP"

# ============================================================
log_test "2.3 搜索职位 (按技能)"
log_request "GET" "http://localhost:8082/api/v1/jobs?skills=Go"

SEARCH_JOBS_RESP=$(curl -s "http://localhost:8082/api/v1/jobs?skills=Go")
log_response "$SEARCH_JOBS_RESP"

# ============================================================
# 3. 人才服务详细测试
# ============================================================

log_test "3.1 获取人才列表"
log_request "GET" "http://localhost:8086/api/v1/talents"

TALENTS_RESP=$(curl -s "http://localhost:8086/api/v1/talents")
log_response "$TALENTS_RESP"

# ============================================================
log_test "3.2 创建新人才"
TALENT_DATA='{
    "name": "王小明",
    "email": "wangxiaoming@test.com",
    "phone": "13900139000",
    "skills": "Go,Python,Docker,Kubernetes",
    "experience": 6,
    "education": "硕士",
    "location": "北京",
    "current_company": "某大厂",
    "current_position": "高级工程师",
    "salary_expectation": "40-60K"
}'
log_request "POST" "http://localhost:8086/api/v1/talents" "$TALENT_DATA"

CREATE_TALENT_RESP=$(curl -s -X POST "http://localhost:8086/api/v1/talents" \
    -H "Content-Type: application/json" \
    -d "$TALENT_DATA")
log_response "$CREATE_TALENT_RESP"

# ============================================================
# 4. 简历服务详细测试
# ============================================================

log_test "4.1 获取简历列表"
log_request "GET" "http://localhost:8084/api/v1/resumes"

RESUMES_RESP=$(curl -s "http://localhost:8084/api/v1/resumes")
log_response "$RESUMES_RESP"

# ============================================================
log_test "4.2 检查AI配置状态"
log_request "GET" "http://localhost:8084/api/v1/ai/config"

AI_CONFIG_RESP=$(curl -s "http://localhost:8084/api/v1/ai/config")
log_response "$AI_CONFIG_RESP"

# ============================================================
log_test "4.3 获取评估结果列表"
log_request "GET" "http://localhost:8084/api/v1/evaluations"

EVALS_RESP=$(curl -s "http://localhost:8084/api/v1/evaluations")
log_response "$EVALS_RESP"

# ============================================================
log_test "4.4 获取评估统计"
log_request "GET" "http://localhost:8084/api/v1/evaluations/stats"

EVAL_STATS_RESP=$(curl -s "http://localhost:8084/api/v1/evaluations/stats")
log_response "$EVAL_STATS_RESP"

# ============================================================
# 5. 推荐服务详细测试
# ============================================================

log_test "5.1 语义匹配测试"
SEMANTIC_DATA='{
    "text1": "Go语言后端开发工程师",
    "text2": "Golang服务端开发"
}'
log_request "POST" "http://localhost:8087/api/v1/recommendations/semantic-match" "$SEMANTIC_DATA"

SEMANTIC_RESP=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/semantic-match" \
    -H "Content-Type: application/json" \
    -d "$SEMANTIC_DATA")
log_response "$SEMANTIC_RESP"

# ============================================================
log_test "5.2 RAG查询 - 搜索相似人才"
RAG_DATA='{
    "query": "熟悉Go语言和微服务架构的后端开发工程师",
    "top_k": 5,
    "type": "talent"
}'
log_request "POST" "http://localhost:8087/api/v1/recommendations/rag/query" "$RAG_DATA"

RAG_RESP=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/rag/query" \
    -H "Content-Type: application/json" \
    -d "$RAG_DATA")
log_response "$RAG_RESP"

# ============================================================
log_test "5.3 RAG查询 - 搜索相似职位"
RAG_JOB_DATA='{
    "query": "高并发系统开发，分布式架构设计",
    "top_k": 5,
    "type": "job"
}'
log_request "POST" "http://localhost:8087/api/v1/recommendations/rag/query" "$RAG_JOB_DATA"

RAG_JOB_RESP=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/rag/query" \
    -H "Content-Type: application/json" \
    -d "$RAG_JOB_DATA")
log_response "$RAG_JOB_RESP"

# ============================================================
log_test "5.4 为人才推荐职位"
TALENT_PROFILE='{
    "id": 1,
    "name": "张伟",
    "skills": ["Go", "Docker", "Kubernetes", "Redis", "微服务"],
    "experience": 6,
    "education": "本科",
    "location": "北京",
    "salary": "30-40K"
}'
log_request "POST" "http://localhost:8087/api/v1/recommendations/jobs-for-talent" "$TALENT_PROFILE"

JOB_REC_RESP=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/jobs-for-talent" \
    -H "Content-Type: application/json" \
    -d "$TALENT_PROFILE")
log_response "$JOB_REC_RESP"

# ============================================================
log_test "5.5 为职位推荐人才"
JOB_PROFILE='{
    "id": 1,
    "title": "高级Go开发工程师",
    "skills": ["Go", "Docker", "Kubernetes"],
    "location": "北京",
    "level": "senior",
    "salary": "30-50K"
}'
log_request "POST" "http://localhost:8087/api/v1/recommendations/talents-for-job" "$JOB_PROFILE"

TALENT_REC_RESP=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/talents-for-job" \
    -H "Content-Type: application/json" \
    -d "$JOB_PROFILE")
log_response "$TALENT_REC_RESP"

# ============================================================
log_test "5.6 推荐统计"
log_request "GET" "http://localhost:8087/api/v1/recommendations/stats"

REC_STATS_RESP=$(curl -s "http://localhost:8087/api/v1/recommendations/stats")
log_response "$REC_STATS_RESP"

# ============================================================
# 6. 面试服务详细测试
# ============================================================

log_test "6.1 获取面试列表"
log_request "GET" "http://localhost:8083/api/v1/interviews"

INTERVIEWS_RESP=$(curl -s "http://localhost:8083/api/v1/interviews")
log_response "$INTERVIEWS_RESP"

# ============================================================
log_test "6.2 创建面试"
INTERVIEW_DATA='{
    "job_id": 1,
    "talent_id": 1,
    "interviewer_id": 1,
    "scheduled_time": "2026-02-15T14:00:00Z",
    "location": "线上-腾讯会议",
    "type": "technical",
    "status": "scheduled",
    "notes": "一面技术面试"
}'
log_request "POST" "http://localhost:8083/api/v1/interviews" "$INTERVIEW_DATA"

CREATE_INTERVIEW_RESP=$(curl -s -X POST "http://localhost:8083/api/v1/interviews" \
    -H "Content-Type: application/json" \
    -d "$INTERVIEW_DATA")
log_response "$CREATE_INTERVIEW_RESP"

# ============================================================
# 7. 消息服务详细测试
# ============================================================

log_test "7.1 获取消息列表"
log_request "GET" "http://localhost:8085/api/v1/messages?user_id=1"

MESSAGES_RESP=$(curl -s "http://localhost:8085/api/v1/messages?user_id=1")
log_response "$MESSAGES_RESP"

# ============================================================
log_test "7.2 发送消息"
MESSAGE_DATA='{
    "sender_id": 1,
    "receiver_id": 2,
    "content": "您好，您的简历已通过初筛，请准备技术面试。",
    "type": "notification"
}'
log_request "POST" "http://localhost:8085/api/v1/messages" "$MESSAGE_DATA"

SEND_MSG_RESP=$(curl -s -X POST "http://localhost:8085/api/v1/messages" \
    -H "Content-Type: application/json" \
    -d "$MESSAGE_DATA")
log_response "$SEND_MSG_RESP"

# ============================================================
echo ""
echo -e "${GREEN}============================================================${NC}"
echo -e "${GREEN}API详细测试完成!${NC}"
echo -e "${GREEN}============================================================${NC}"
