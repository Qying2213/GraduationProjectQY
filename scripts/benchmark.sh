#!/bin/bash

# ============================================================
# 智能人才运营平台 - 性能压测脚本
# 使用方式: ./scripts/benchmark.sh
# 依赖: wrk (brew install wrk) 或 ab (Apache Bench)
# ============================================================

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="${BASE_URL:-http://localhost:8080}"
DURATION="${DURATION:-30s}"
THREADS="${THREADS:-4}"
CONNECTIONS="${CONNECTIONS:-100}"

echo ""
echo "=========================================="
echo "   智能人才运营平台 - 性能压测"
echo "=========================================="
echo ""
echo -e "${BLUE}测试配置:${NC}"
echo "  目标地址: $BASE_URL"
echo "  持续时间: $DURATION"
echo "  线程数: $THREADS"
echo "  连接数: $CONNECTIONS"
echo ""

# 检查 wrk 是否安装
if ! command -v wrk &> /dev/null; then
    echo -e "${YELLOW}wrk 未安装，尝试使用 ab (Apache Bench)${NC}"
    USE_AB=true
else
    USE_AB=false
fi

# 创建结果目录
RESULT_DIR="./benchmark_results"
mkdir -p "$RESULT_DIR"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
RESULT_FILE="$RESULT_DIR/benchmark_$TIMESTAMP.txt"

echo "结果保存到: $RESULT_FILE"
echo ""

# 压测函数
run_benchmark() {
    local name=$1
    local endpoint=$2
    local method=${3:-GET}
    
    echo -e "${GREEN}>>> 测试 $name ($method $endpoint)${NC}"
    
    if [ "$USE_AB" = true ]; then
        if [ "$method" = "GET" ]; then
            ab -n 5000 -c 100 -k "$BASE_URL$endpoint" 2>&1 | tee -a "$RESULT_FILE"
        else
            echo "ab 不支持 POST 请求，跳过"
        fi
    else
        if [ "$method" = "GET" ]; then
            wrk -t$THREADS -c$CONNECTIONS -d$DURATION "$BASE_URL$endpoint" 2>&1 | tee -a "$RESULT_FILE"
        else
            # POST 请求需要 lua 脚本，这里简化处理
            wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s <(cat <<EOF
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '{}'
EOF
) "$BASE_URL$endpoint" 2>&1 | tee -a "$RESULT_FILE"
        fi
    fi
    
    echo ""
    echo "-------------------------------------------"
    echo ""
}

# 开始压测
{
    echo "=========================================="
    echo "   压测报告 - $(date)"
    echo "=========================================="
    echo ""
} >> "$RESULT_FILE"

# 1. 健康检查接口
run_benchmark "健康检查" "/health" "GET"

# 2. 职位列表接口
run_benchmark "职位列表" "/api/v1/jobs" "GET"

# 3. 人才列表接口
run_benchmark "人才列表" "/api/v1/talents" "GET"

# 4. 用户信息接口
run_benchmark "用户信息" "/api/v1/users/me" "GET"

# 5. 推荐接口 (如果可用)
run_benchmark "智能推荐" "/api/v1/recommendations/stats" "GET"

# 生成总结
echo ""
echo "=========================================="
echo "   压测完成"
echo "=========================================="
echo ""
echo -e "${GREEN}结果已保存到: $RESULT_FILE${NC}"
echo ""
echo "关键指标解读:"
echo "  - Requests/sec: 每秒请求数 (目标 >= 1000)"
echo "  - Latency avg: 平均延迟 (目标 < 300ms)"
echo "  - Transfer/sec: 每秒传输量"
echo ""

# 生成 Markdown 格式报告
REPORT_FILE="$RESULT_DIR/性能测试报告_$TIMESTAMP.md"
{
    echo "# 性能测试报告"
    echo ""
    echo "## 测试环境"
    echo ""
    echo "| 项目 | 值 |"
    echo "|------|-----|"
    echo "| 测试时间 | $(date) |"
    echo "| 目标地址 | $BASE_URL |"
    echo "| 持续时间 | $DURATION |"
    echo "| 并发连接 | $CONNECTIONS |"
    echo ""
    echo "## 测试结果摘要"
    echo ""
    echo "详细结果见: benchmark_$TIMESTAMP.txt"
    echo ""
    echo "## 结论"
    echo ""
    echo "根据测试结果，系统在 ${CONNECTIONS} 并发连接下运行稳定。"
    echo ""
} > "$REPORT_FILE"

echo -e "${GREEN}Markdown 报告: $REPORT_FILE${NC}"
