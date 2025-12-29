#!/bin/bash
# =====================================================
# 智能人才运营平台 - 后端服务启动脚本
# 启动7个核心微服务，每个服务一个终端窗口
# =====================================================

echo "🚀 启动所有后端服务..."
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# macOS: 用 osascript 打开新终端窗口
open_terminal() {
    local name=$1
    local port=$2
    local service=$3
    osascript -e "tell application \"Terminal\" to do script \"cd '$SCRIPT_DIR/backend/$service' && echo '=== $name (端口 $port) ===' && go run main.go\""
}

# 启动7个核心后端服务
open_terminal "user-service" "8081" "user-service"
open_terminal "job-service" "8082" "job-service"
open_terminal "interview-service" "8083" "interview-service"
open_terminal "resume-service" "8084" "resume-service"
open_terminal "message-service" "8085" "message-service"
open_terminal "talent-service" "8086" "talent-service"
open_terminal "log-service" "8088" "log-service"

# 启动AI评估服务（需要配置Coze）
open_terminal_cmd() {
    local name=$1
    local port=$2
    local dir=$3
    local cmd=$4
    osascript -e "tell application \"Terminal\" to do script \"cd '$SCRIPT_DIR/backend/$dir' && echo '=== $name (端口 $port) ===' && $cmd\""
}
open_terminal_cmd "evaluator-service" "8090" "evaluator-service" "go run cmd/server/main.go"

echo "✅ 已启动8个后端服务终端"
echo ""
echo "服务列表："
echo "  ├── user-service        http://localhost:8081"
echo "  ├── job-service         http://localhost:8082"
echo "  ├── interview-service   http://localhost:8083"
echo "  ├── resume-service      http://localhost:8084"
echo "  ├── message-service     http://localhost:8085"
echo "  ├── talent-service      http://localhost:8086"
echo "  ├── log-service (ES)    http://localhost:8088"
echo "  └── evaluator-service   http://localhost:8090 (AI评估)"
echo ""
echo "前端启动: cd frontend && npm run dev"
echo "API测试:  cd backend && ./test_api.sh"
