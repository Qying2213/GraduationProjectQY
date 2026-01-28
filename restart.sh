#!/bin/bash

# 智能人才运营平台 - 快速重启脚本
# 先关闭旧终端，再启动新终端

echo "🔄 正在重启智能人才运营平台..."
echo ""

# 获取项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$PROJECT_ROOT/backend/.env"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. 停止所有服务
echo -e "${YELLOW}[1/3] 停止所有服务...${NC}"

# 停止后端服务（通过端口查找进程）
PORTS=(8080 8081 8082 8083 8084 8085 8086 8087)
for port in "${PORTS[@]}"; do
    pid=$(lsof -ti:$port 2>/dev/null)
    if [ -n "$pid" ]; then
        kill -9 $pid 2>/dev/null
        echo "  ✓ 停止端口 $port 的服务"
    fi
done

# 停止前端服务
pid=$(lsof -ti:5173 2>/dev/null)
if [ -n "$pid" ]; then
    kill -9 $pid 2>/dev/null
    echo "  ✓ 停止前端服务"
fi

sleep 1
echo -e "${GREEN}  所有服务已停止${NC}"
echo ""

# 2. 关闭旧的服务终端窗口
echo -e "${YELLOW}[2/3] 关闭旧终端窗口...${NC}"
osascript <<EOF
tell application "Terminal"
    set windowsToClose to {}
    repeat with w in windows
        repeat with t in tabs of w
            set tabName to name of t
            if tabName contains "Gateway" or tabName contains "Service" or tabName contains "Frontend" or tabName contains "8080" or tabName contains "8081" or tabName contains "8082" or tabName contains "8083" or tabName contains "8084" or tabName contains "8085" or tabName contains "8086" or tabName contains "8087" or tabName contains "5173" then
                set end of windowsToClose to w
                exit repeat
            end if
        end repeat
    end repeat
    repeat with w in windowsToClose
        close w
    end repeat
end tell
EOF
echo -e "${GREEN}  旧终端已关闭${NC}"
echo ""

# 3. 启动所有服务（每个服务一个终端窗口）
echo -e "${YELLOW}[3/3] 启动所有服务...${NC}"

# 构建环境变量加载命令
ENV_CMD=""
if [ -f "$ENV_FILE" ]; then
    ENV_CMD="export \$(grep -v '^#' $ENV_FILE | xargs) && "
fi

# 启动 Gateway
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/gateway' && $ENV_CMD echo '🚀 Gateway (8080)' && go run .\""
echo "  ✓ Gateway (8080)"

# 启动 User Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/user-service' && $ENV_CMD echo '🚀 User Service (8081)' && go run .\""
echo "  ✓ User Service (8081)"

# 启动 Job Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/job-service' && $ENV_CMD echo '🚀 Job Service (8082)' && go run .\""
echo "  ✓ Job Service (8082)"

# 启动 Interview Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/interview-service' && $ENV_CMD echo '🚀 Interview Service (8083)' && go run .\""
echo "  ✓ Interview Service (8083)"

# 启动 Resume Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/resume-service' && $ENV_CMD echo '🚀 Resume Service (8084)' && go run .\""
echo "  ✓ Resume Service (8084)"

# 启动 Message Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/message-service' && $ENV_CMD echo '🚀 Message Service (8085)' && go run .\""
echo "  ✓ Message Service (8085)"

# 启动 Talent Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/talent-service' && $ENV_CMD echo '🚀 Talent Service (8086)' && go run .\""
echo "  ✓ Talent Service (8086)"

# 启动 Recommendation Service
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/backend/recommendation-service' && $ENV_CMD echo '🚀 Recommendation Service (8087)' && go run .\""
echo "  ✓ Recommendation Service (8087)"

# 启动前端
osascript -e "tell application \"Terminal\" to do script \"cd '$PROJECT_ROOT/frontend' && echo '🚀 Frontend (5173)' && npm run dev\""
echo "  ✓ Frontend (5173)"

echo ""
echo -e "${GREEN}✅ 所有服务已在独立终端窗口中重启！${NC}"
echo ""
echo "服务地址:"
echo "  前端:     http://localhost:5173"
echo "  Gateway:  http://localhost:8080"
echo ""
echo "提示: 关闭终端窗口即可停止对应服务"
