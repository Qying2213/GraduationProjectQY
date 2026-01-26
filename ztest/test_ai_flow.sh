#!/bin/bash

# ============================================================
# 智能招聘系统 - AI流程专项测试
# ============================================================
# 测试完整的AI评估流程:
# PDF → OCR文本提取 → Embedding向量化 → RAG检索 → Coze AI评估
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m'

echo -e "${MAGENTA}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           智能招聘系统 - AI流程专项测试                      ║"
echo "║                                                              ║"
echo "║  测试流程: PDF → OCR → Embedding → RAG → Coze评估           ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# ============================================================
# 1. 检查服务状态
# ============================================================

echo -e "\n${CYAN}[1/7] 检查服务状态${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

check_service() {
    local name=$1
    local port=$2
    if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
        echo -e "  ${GREEN}✓${NC} $name (端口 $port) - 运行中"
        return 0
    else
        echo -e "  ${RED}✗${NC} $name (端口 $port) - 未运行"
        return 1
    fi
}

check_service "Resume Service" 8084
check_service "Recommendation Service" 8087

# ============================================================
# 2. 检查AI配置
# ============================================================

echo -e "\n${CYAN}[2/7] 检查AI配置${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

AI_CONFIG=$(curl -s "http://localhost:8084/api/v1/ai/config")
echo "AI配置响应:"
echo "$AI_CONFIG" | python3 -m json.tool 2>/dev/null || echo "$AI_CONFIG"

if echo "$AI_CONFIG" | grep -q '"configured":true'; then
    echo -e "  ${GREEN}✓${NC} Coze AI 已配置"
else
    echo -e "  ${RED}✗${NC} Coze AI 未配置"
    echo "请检查环境变量: COZE_TOKEN, COZE_WORKFLOW_ID"
fi

# ============================================================
# 3. 检查Embedding服务
# ============================================================

echo -e "\n${CYAN}[3/7] 测试Embedding服务 (Volcengine Doubao)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

EMBEDDING_TEST=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/semantic-match" \
    -H "Content-Type: application/json" \
    -d '{"text1":"Go语言后端开发","text2":"Golang服务端工程师"}')

echo "Embedding测试响应:"
echo "$EMBEDDING_TEST" | python3 -m json.tool 2>/dev/null || echo "$EMBEDDING_TEST"

if echo "$EMBEDDING_TEST" | grep -q '"similarity"'; then
    SIMILARITY=$(echo "$EMBEDDING_TEST" | grep -o '"similarity":[0-9.]*' | cut -d':' -f2)
    echo -e "  ${GREEN}✓${NC} Embedding服务正常，相似度: $SIMILARITY"
else
    echo -e "  ${RED}✗${NC} Embedding服务异常"
fi

# ============================================================
# 4. 检查RAG服务
# ============================================================

echo -e "\n${CYAN}[4/7] 测试RAG服务 (pgvector)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

RAG_TEST=$(curl -s -X POST "http://localhost:8087/api/v1/recommendations/rag/query" \
    -H "Content-Type: application/json" \
    -d '{"query":"Go后端开发工程师，熟悉微服务架构","top_k":3,"type":"talent"}')

echo "RAG查询响应:"
echo "$RAG_TEST" | python3 -m json.tool 2>/dev/null || echo "$RAG_TEST"

if echo "$RAG_TEST" | grep -q '"results"'; then
    RESULT_COUNT=$(echo "$RAG_TEST" | grep -o '"id":[0-9]*' | wc -l | tr -d ' ')
    echo -e "  ${GREEN}✓${NC} RAG服务正常，返回 $RESULT_COUNT 条相似人才"
else
    echo -e "  ${RED}✗${NC} RAG服务异常"
fi

# ============================================================
# 5. 检查数据库中的简历
# ============================================================

echo -e "\n${CYAN}[5/7] 检查测试数据${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "数据库中的简历:"
psql -U qinyang -d talent_platform -c "SELECT id, file_name, file_path, status FROM resumes LIMIT 5;" 2>/dev/null || echo "无法连接数据库"

echo ""
echo "向量数据库统计:"
TALENT_EMB_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM talent_embeddings" 2>/dev/null | tr -d ' ')
JOB_EMB_COUNT=$(psql -U qinyang -d talent_platform -t -c "SELECT COUNT(*) FROM job_embeddings" 2>/dev/null | tr -d ' ')
echo "  - talent_embeddings: $TALENT_EMB_COUNT 条"
echo "  - job_embeddings: $JOB_EMB_COUNT 条"

# ============================================================
# 6. 执行完整AI评估流程
# ============================================================

echo -e "\n${CYAN}[6/7] 执行完整AI评估流程${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "${YELLOW}注意: 此测试需要1-2分钟完成，请耐心等待...${NC}"
echo ""
echo "测试参数:"
echo "  - 简历ID: 2"
echo "  - JD: 招聘Go后端开发工程师，要求3年以上经验，熟悉微服务架构"
echo ""

echo "发送评估请求..."
START_TIME=$(date +%s)

AI_EVAL_RESP=$(curl -s --max-time 300 -X POST "http://localhost:8084/api/v1/ai/evaluate" \
    -H "Content-Type: application/json" \
    -d '{
        "resume_id": 2,
        "jd_text": "招聘Go后端开发工程师，要求3年以上经验，熟悉微服务架构，熟练使用Docker和Kubernetes，有高并发系统开发经验优先"
    }')

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "评估完成! 耗时: ${DURATION}秒"
echo ""
echo "完整响应:"
echo "$AI_EVAL_RESP" | python3 -m json.tool 2>/dev/null || echo "$AI_EVAL_RESP"

# ============================================================
# 7. 解析评估结果
# ============================================================

echo -e "\n${CYAN}[7/7] 评估结果分析${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if echo "$AI_EVAL_RESP" | grep -q '"code":0'; then
    echo -e "${GREEN}评估成功!${NC}"
    echo ""
    
    # 提取关键信息
    TOTAL_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"total_score":[0-9.]*' | cut -d':' -f2)
    GRADE=$(echo "$AI_EVAL_RESP" | grep -o '"grade":"[^"]*"' | cut -d'"' -f4)
    JD_MATCH=$(echo "$AI_EVAL_RESP" | grep -o '"jd_match_score":[0-9]*' | cut -d':' -f2)
    OCR_USED=$(echo "$AI_EVAL_RESP" | grep -o '"ocr_extracted":[a-z]*' | cut -d':' -f2)
    EMBEDDING_USED=$(echo "$AI_EVAL_RESP" | grep -o '"embedding_used":[a-z]*' | cut -d':' -f2)
    RAG_USED=$(echo "$AI_EVAL_RESP" | grep -o '"rag_enhanced":[a-z]*' | cut -d':' -f2)
    RECOMMENDATION=$(echo "$AI_EVAL_RESP" | grep -o '"recommendation":"[^"]*"' | cut -d'"' -f4)
    
    echo "┌─────────────────────────────────────────────────────────────┐"
    echo "│                      评估结果摘要                          │"
    echo "├─────────────────────────────────────────────────────────────┤"
    printf "│  %-20s │  %-35s │\n" "总分" "$TOTAL_SCORE"
    printf "│  %-20s │  %-35s │\n" "等级" "$GRADE"
    printf "│  %-20s │  %-35s │\n" "JD匹配度" "$JD_MATCH"
    echo "├─────────────────────────────────────────────────────────────┤"
    echo "│                      流程执行情况                          │"
    echo "├─────────────────────────────────────────────────────────────┤"
    printf "│  %-20s │  %-35s │\n" "OCR文本提取" "$OCR_USED"
    printf "│  %-20s │  %-35s │\n" "Embedding向量化" "$EMBEDDING_USED"
    printf "│  %-20s │  %-35s │\n" "RAG检索增强" "$RAG_USED"
    echo "├─────────────────────────────────────────────────────────────┤"
    printf "│  %-20s │  %-35s │\n" "录用建议" "$RECOMMENDATION"
    echo "└─────────────────────────────────────────────────────────────┘"
    
    # 提取各维度得分
    AGE_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"age_score":[0-9]*' | cut -d':' -f2)
    EXP_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"experience_score":[0-9]*' | cut -d':' -f2)
    EDU_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"education_score":[0-9]*' | cut -d':' -f2)
    COMPANY_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"company_score":[0-9]*' | cut -d':' -f2)
    TECH_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"tech_score":[0-9]*' | cut -d':' -f2)
    PROJECT_SCORE=$(echo "$AI_EVAL_RESP" | grep -o '"project_score":[0-9]*' | cut -d':' -f2)
    
    echo ""
    echo "各维度得分 (满分10分):"
    echo "  年龄适配:   $AGE_SCORE/10"
    echo "  工作经验:   $EXP_SCORE/10"
    echo "  学历背景:   $EDU_SCORE/10"
    echo "  公司背景:   $COMPANY_SCORE/10"
    echo "  技术能力:   $TECH_SCORE/10"
    echo "  项目经验:   $PROJECT_SCORE/10"
    
else
    echo -e "${RED}评估失败!${NC}"
    echo "错误信息: $AI_EVAL_RESP"
fi

# ============================================================
# 总结
# ============================================================

echo ""
echo -e "${MAGENTA}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                      测试完成                                ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo "AI评估流程说明:"
echo ""
echo "  ┌─────────┐    ┌───────────┐    ┌─────────┐    ┌──────────┐"
echo "  │   PDF   │ -> │    OCR    │ -> │Embedding│ -> │   RAG    │"
echo "  │  简历   │    │ 文本提取  │    │ 向量化  │    │ 检索增强 │"
echo "  └─────────┘    └───────────┘    └─────────┘    └──────────┘"
echo "                                                       │"
echo "                                                       v"
echo "  ┌─────────────────────────────────────────────────────────┐"
echo "  │                    Coze AI 评估                         │"
echo "  │  - 基本信息解析                                         │"
echo "  │  - 6维度评分 (年龄/经验/学历/公司/技术/项目)            │"
echo "  │  - JD匹配度分析                                         │"
echo "  │  - 综合评价与录用建议                                   │"
echo "  └─────────────────────────────────────────────────────────┘"
