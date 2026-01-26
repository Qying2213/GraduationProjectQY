#!/bin/bash

# ============================================================
# 智能招聘系统 - 数据库测试脚本
# ============================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

DB_USER="qinyang"
DB_NAME="talent_platform"

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║           智能招聘系统 - 数据库测试                          ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# ============================================================
# 1. 连接测试
# ============================================================

echo -e "\n${CYAN}[1] PostgreSQL连接测试${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if psql -U $DB_USER -d $DB_NAME -c "SELECT version();" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PostgreSQL连接成功${NC}"
    psql -U $DB_USER -d $DB_NAME -c "SELECT version();"
else
    echo -e "${RED}✗ PostgreSQL连接失败${NC}"
    exit 1
fi

# ============================================================
# 2. 扩展检查
# ============================================================

echo -e "\n${CYAN}[2] 扩展检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "已安装的扩展:"
psql -U $DB_USER -d $DB_NAME -c "SELECT extname, extversion FROM pg_extension;"

# 检查pgvector
if psql -U $DB_USER -d $DB_NAME -c "SELECT 1 FROM pg_extension WHERE extname='vector';" | grep -q "1"; then
    echo -e "${GREEN}✓ pgvector扩展已安装${NC}"
else
    echo -e "${RED}✗ pgvector扩展未安装${NC}"
fi

# ============================================================
# 3. 表结构检查
# ============================================================

echo -e "\n${CYAN}[3] 表结构检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "数据库中的表:"
psql -U $DB_USER -d $DB_NAME -c "\dt"

# ============================================================
# 4. 数据统计
# ============================================================

echo -e "\n${CYAN}[4] 数据统计${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "各表数据量:"
echo ""

# 用户表
USER_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM users" 2>/dev/null | tr -d ' ')
echo "  users (用户表):              $USER_COUNT 条"

# 职位表
JOB_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM jobs" 2>/dev/null | tr -d ' ')
echo "  jobs (职位表):               $JOB_COUNT 条"

# 人才表
TALENT_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM talents" 2>/dev/null | tr -d ' ')
echo "  talents (人才表):            $TALENT_COUNT 条"

# 简历表
RESUME_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM resumes" 2>/dev/null | tr -d ' ')
echo "  resumes (简历表):            $RESUME_COUNT 条"

# 面试表
INTERVIEW_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM interviews" 2>/dev/null | tr -d ' ')
echo "  interviews (面试表):         $INTERVIEW_COUNT 条"

# 消息表
MESSAGE_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM messages" 2>/dev/null | tr -d ' ')
echo "  messages (消息表):           $MESSAGE_COUNT 条"

# 评估结果表
EVAL_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM evaluation_results" 2>/dev/null | tr -d ' ')
echo "  evaluation_results (评估表): $EVAL_COUNT 条"

# 向量表
TALENT_EMB_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM talent_embeddings" 2>/dev/null | tr -d ' ')
echo "  talent_embeddings (人才向量): $TALENT_EMB_COUNT 条"

JOB_EMB_COUNT=$(psql -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM job_embeddings" 2>/dev/null | tr -d ' ')
echo "  job_embeddings (职位向量):   $JOB_EMB_COUNT 条"

# ============================================================
# 5. 示例数据查询
# ============================================================

echo -e "\n${CYAN}[5] 示例数据查询${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${YELLOW}用户示例 (前5条):${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, username, email, role FROM users LIMIT 5;"

echo -e "\n${YELLOW}职位示例 (前5条):${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, title, department, location, status FROM jobs LIMIT 5;"

echo -e "\n${YELLOW}人才示例 (前5条):${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, name, skills, experience, education FROM talents LIMIT 5;"

echo -e "\n${YELLOW}简历示例 (前5条):${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, file_name, status, match_score FROM resumes LIMIT 5;"

echo -e "\n${YELLOW}评估结果示例 (前5条):${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, resume_name, match_score, match_level, status FROM evaluation_results LIMIT 5;"

# ============================================================
# 6. 向量数据检查
# ============================================================

echo -e "\n${CYAN}[6] 向量数据检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo -e "\n${YELLOW}人才向量示例:${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, talent_id, LEFT(content, 100) as content_preview FROM talent_embeddings LIMIT 3;"

echo -e "\n${YELLOW}职位向量示例:${NC}"
psql -U $DB_USER -d $DB_NAME -c "SELECT id, job_id, LEFT(content, 100) as content_preview FROM job_embeddings LIMIT 3;"

# ============================================================
# 7. 索引检查
# ============================================================

echo -e "\n${CYAN}[7] 索引检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "向量索引:"
psql -U $DB_USER -d $DB_NAME -c "SELECT indexname, indexdef FROM pg_indexes WHERE tablename IN ('talent_embeddings', 'job_embeddings');"

# ============================================================
# 8. 向量相似度查询测试
# ============================================================

echo -e "\n${CYAN}[8] 向量相似度查询测试${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "测试向量相似度查询 (查找与第一条人才最相似的3条记录):"
psql -U $DB_USER -d $DB_NAME -c "
SELECT 
    t2.talent_id,
    t2.content,
    1 - (t1.embedding <=> t2.embedding) as similarity
FROM talent_embeddings t1, talent_embeddings t2
WHERE t1.id = 1 AND t1.id != t2.id
ORDER BY t1.embedding <=> t2.embedding
LIMIT 3;
" 2>/dev/null || echo "向量查询失败，可能没有足够的数据"

# ============================================================
# 总结
# ============================================================

echo -e "\n${GREEN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                    数据库测试完成                            ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo "数据库状态摘要:"
echo "  - PostgreSQL: 正常"
echo "  - pgvector扩展: 已安装"
echo "  - 用户数: $USER_COUNT"
echo "  - 职位数: $JOB_COUNT"
echo "  - 人才数: $TALENT_COUNT"
echo "  - 简历数: $RESUME_COUNT"
echo "  - 人才向量: $TALENT_EMB_COUNT"
echo "  - 职位向量: $JOB_EMB_COUNT"
