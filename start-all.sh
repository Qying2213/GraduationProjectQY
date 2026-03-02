#!/bin/bash

# 智能人才运营平台 - 一键启动脚本（多终端版）
# 使用方法: ./start-all.sh
# macOS 系统会为每个服务打开独立的终端窗口

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# 获取项目根目录的绝对路径
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$PROJECT_ROOT/backend/.env"

echo ""
echo "=========================================="
echo "   智能人才运营平台 - 启动脚本"
echo "=========================================="
echo ""
echo "项目路径: $PROJECT_ROOT"
echo ""

# 检查 .env 文件是否存在
if [ -f "$ENV_FILE" ]; then
    print_info "加载环境变量: $ENV_FILE"
    # 读取环境变量（用于显示）
    export $(grep -v '^#' "$ENV_FILE" | xargs)
else
    print_info "未找到 .env 文件，使用默认配置"
    print_info "可复制 backend/.env.example 为 backend/.env 并配置"
fi

# 检查是否是macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    print_info "检测到 macOS 系统，将为每个服务打开独立终端"
    echo ""
    
    # 定义服务列表: 目录名:端口:服务名:启动命令
    services=(
        "gateway:8080:Gateway网关:go run main.go"
        "user-service:8081:用户服务:go run main.go"
        "job-service:8082:职位服务:go run main.go"
        "interview-service:8083:面试服务:go run main.go"
        "resume-service:8084:简历服务:go run main.go"
        "message-service:8085:消息服务:go run main.go"
        "talent-service:8086:人才服务:go run main.go"
        "recommendation-service:8087:推荐服务:go run main.go"
        "log-service:8088:日志服务:go run main.go"
        "evaluator-service:8090:AI评估服务:go run cmd/server/main.go"
    )
    
    for service in "${services[@]}"; do
        IFS=':' read -r dir port name cmd <<< "$service"
        
        print_info "启动 $name (端口: $port)..."
        
        # 使用 osascript 打开新的 Terminal 窗口，并加载环境变量
        osascript <<EOF
tell application "Terminal"
    activate
    do script "cd '$PROJECT_ROOT/backend/$dir' && if [ -f '$ENV_FILE' ]; then export \$(grep -v '^#' '$ENV_FILE' | xargs); fi && echo '========================================' && echo '  $name - 端口 $port' && echo '========================================' && echo '' && $cmd"
end tell
EOF
        
        sleep 0.5
    done
    
    echo ""
    print_success "所有服务终端已打开！"
    echo ""
    echo "=========================================="
    echo "   服务端口一览"
    echo "=========================================="
    echo ""
    echo "  Gateway网关:     http://localhost:8080"
    echo "  用户服务:        http://localhost:8081"
    echo "  职位服务:        http://localhost:8082"
    echo "  面试服务:        http://localhost:8083"
    echo "  简历服务:        http://localhost:8084"
    echo "  消息服务:        http://localhost:8085"
    echo "  人才服务:        http://localhost:8086"
    echo "  推荐服务:        http://localhost:8087"
    echo "  日志服务:        http://localhost:8088"
    echo "  AI评估服务:      http://localhost:8090"
    echo ""
    echo "=========================================="
    echo ""
    echo "启动前端:"
    echo "  cd frontend && npm run dev"
    echo ""
    echo "停止所有服务:"
    echo "  ./stop-all.sh"
    echo ""

else
    # Linux 或其他系统，使用后台进程方式
    print_info "非 macOS 系统，使用后台进程方式启动"
    
    # 加载环境变量
    if [ -f "$ENV_FILE" ]; then
        export $(grep -v '^#' "$ENV_FILE" | xargs)
    fi
    
    services=(
        "gateway:8080:Gateway:go run main.go"
        "user-service:8081:User Service:go run main.go"
        "job-service:8082:Job Service:go run main.go"
        "interview-service:8083:Interview Service:go run main.go"
        "resume-service:8084:Resume Service:go run main.go"
        "message-service:8085:Message Service:go run main.go"
        "talent-service:8086:Talent Service:go run main.go"
        "recommendation-service:8087:Recommendation Service:go run main.go"
        "log-service:8088:Log Service:go run main.go"
        "evaluator-service:8090:Evaluator Service:go run cmd/server/main.go"
    )
    
    for service in "${services[@]}"; do
        IFS=':' read -r dir port name cmd <<< "$service"
        
        print_info "启动 $name (端口: $port)..."
        (cd "$PROJECT_ROOT/backend/$dir" && eval "$cmd" > "/tmp/${dir}.log" 2>&1 &)
        sleep 1
    done
    
    echo ""
    print_success "所有服务已在后台启动"
    echo "查看日志: tail -f /tmp/gateway.log"
fi
