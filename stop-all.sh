#!/bin/bash

# 智能人才运营平台 - 停止脚本
# 使用方法: ./stop-all.sh

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo "=========================================="
echo "   停止所有服务（后端+前端）"
echo "=========================================="
echo ""

# 定义后端服务端口
ports=(8080 8081 8082 8083 8084 8085 8086 8087 8088 8089)

echo "停止后端服务..."
for port in "${ports[@]}"; do
    pid=$(lsof -ti :$port 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}停止端口 $port 的服务 (PID: $pid)${NC}"
        kill -9 $pid 2>/dev/null
    fi
done

# 停止前端服务
echo ""
echo "停止前端服务..."

# 停止 5173 端口
frontend_pid=$(lsof -ti :5173 2>/dev/null)
if [ -n "$frontend_pid" ]; then
    echo -e "${YELLOW}停止前端 (端口 5173, PID: $frontend_pid)${NC}"
    kill -9 $frontend_pid 2>/dev/null
fi

# 杀死 vite 和 npm 相关进程
pkill -9 -f "vite" 2>/dev/null && echo -e "${YELLOW}停止 Vite 进程${NC}"
pkill -9 -f "npm run dev" 2>/dev/null && echo -e "${YELLOW}停止 npm 进程${NC}"

# 清理 PID 文件
rm -f .pids/*.pid 2>/dev/null

echo ""
echo -e "${GREEN}所有服务已停止${NC}"
echo ""
