#!/bin/bash

# 智能人才运营平台 - 停止脚本
# 使用方法: ./stop-all.sh

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo ""
echo "=========================================="
echo "   停止所有后端服务"
echo "=========================================="
echo ""

# 定义服务端口
ports=(8080 8081 8082 8083 8084 8085 8086 8087 8088 8089)

for port in "${ports[@]}"; do
    pid=$(lsof -ti :$port 2>/dev/null)
    if [ -n "$pid" ]; then
        echo -e "${YELLOW}停止端口 $port 的服务 (PID: $pid)${NC}"
        kill -9 $pid 2>/dev/null
    fi
done

echo ""
echo -e "${GREEN}所有服务已停止${NC}"
echo ""
