#!/bin/bash

echo "=========================================="
echo "智能人才运营平台 - 服务启动脚本"
echo "=========================================="
echo ""

# 检查PostgreSQL
echo "检查PostgreSQL连接..."
if command -v psql &> /dev/null; then
    echo "✓ PostgreSQL 已安装"
else
    echo "✗ 请先安装PostgreSQL"
    exit 1
fi

# 检查Go
echo "检查Go环境..."
if command -v go &> /dev/null; then
    echo "✓ Go $(go version) 已安装"
else
    echo "✗ 请先安装Go 1.21+"
    exit 1
fi

echo ""
echo "准备启动7个后端微服务..."
echo "提示: 每个服务将在新的终端窗口中启动"
echo ""

# macOS使用osascript启动新终端
if [[ "$OSTYPE" == "darwin"* ]]; then
    # 用户服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/user-service && echo '启动用户服务 (8081)...' && go mod tidy && go run main.go\""
    
    # 人才服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/talent-service && echo '启动人才服务 (8082)...' && go mod tidy && go run main.go\""
    
    # 职位服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/job-service && echo '启动职位服务 (8083)...' && go mod tidy && go run main.go\""
    
    # 简历服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/resume-service && echo '启动简历服务 (8084)...' && go mod tidy && go run main.go\""
    
    # 推荐服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/recommendation-service && echo '启动推荐服务 (8085)...' && go mod tidy && go run main.go\""
    
    # 消息服务
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/message-service && echo '启动消息服务 (8086)...' && go mod tidy && go run main.go\""
    
    # API网关
    osascript -e "tell app \"Terminal\" to do script \"cd $(pwd)/backend/gateway && echo '启动API网关 (8080)...' && go mod tidy && go run main.go\""
else
    echo "请手动在7个终端窗口中启动各个服务:"
    echo ""
    echo "终端1: cd backend/user-service && go mod tidy && go run main.go"
    echo "终端2: cd backend/talent-service && go mod tidy && go run main.go"
    echo "终端3: cd backend/job-service && go mod tidy && go run main.go"
    echo "终端4: cd backend/resume-service && go mod tidy && go run main.go"
    echo "终端5: cd backend/recommendation-service && go mod tidy && go run main.go"
    echo "终端6: cd backend/message-service && go mod tidy && go run main.go"
    echo "终端7: cd backend/gateway && go mod tidy && go run main.go"
    echo ""
fi

echo ""
echo "✓ 后端服务启动命令已发送"
echo ""
echo "等待服务启动完成后，请运行以下命令启动前端:"
echo "  cd frontend && npm install && npm run dev"
echo ""
echo "访问地址:"
echo "  前端: http://localhost:3000"
echo "  API网关: http://localhost:8080"
echo ""
echo "=========================================="
