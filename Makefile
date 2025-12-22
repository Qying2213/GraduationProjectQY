.PHONY: all build test clean docker-build docker-up docker-down help

# 变量
SERVICES := gateway user-service talent-service job-service resume-service recommendation-service message-service interview-service

# 默认目标
all: help

# 帮助信息
help:
	@echo "智能人才运营平台 - 开发命令"
	@echo ""
	@echo "使用方法: make [命令]"
	@echo ""
	@echo "开发命令:"
	@echo "  dev-backend     启动所有后端服务（开发模式）"
	@echo "  dev-frontend    启动前端开发服务器"
	@echo "  dev             同时启动前后端"
	@echo ""
	@echo "构建命令:"
	@echo "  build           构建所有服务"
	@echo "  build-frontend  构建前端"
	@echo "  build-backend   构建所有后端服务"
	@echo ""
	@echo "测试命令:"
	@echo "  test            运行所有测试"
	@echo "  test-backend    运行后端测试"
	@echo "  test-frontend   运行前端测试"
	@echo "  test-coverage   运行测试并生成覆盖率报告"
	@echo ""
	@echo "Docker 命令:"
	@echo "  docker-build    构建 Docker 镜像"
	@echo "  docker-up       启动 Docker 容器"
	@echo "  docker-down     停止 Docker 容器"
	@echo "  docker-logs     查看 Docker 日志"
	@echo ""
	@echo "数据库命令:"
	@echo "  db-init         初始化数据库"
	@echo "  db-reset        重置数据库"
	@echo "  db-migrate      运行数据库迁移"
	@echo ""
	@echo "其他命令:"
	@echo "  clean           清理构建产物"
	@echo "  lint            运行代码检查"
	@echo "  fmt             格式化代码"

# 开发模式
dev-backend:
	@echo "启动后端服务..."
	./start-backend.sh

dev-frontend:
	@echo "启动前端开发服务器..."
	cd frontend && npm run dev

dev:
	@echo "启动开发环境..."
	@make -j2 dev-backend dev-frontend

# 构建
build: build-backend build-frontend

build-frontend:
	@echo "构建前端..."
	cd frontend && npm ci && npm run build

build-backend:
	@echo "构建后端服务..."
	@for service in $(SERVICES); do \
		echo "构建 $$service..."; \
		cd backend/$$service && go build -o $$service . && cd ../..; \
	done

# 测试
test: test-backend test-frontend

test-backend:
	@echo "运行后端测试..."
	@for service in $(SERVICES); do \
		echo "测试 $$service..."; \
		cd backend/$$service && go test -v ./... && cd ../..; \
	done

test-frontend:
	@echo "运行前端测试..."
	cd frontend && npm run test

test-coverage:
	@echo "运行测试并生成覆盖率报告..."
	@for service in $(SERVICES); do \
		echo "测试 $$service..."; \
		cd backend/$$service && go test -coverprofile=coverage.out ./... && cd ../..; \
	done
	cd frontend && npm run test:coverage

# Docker
docker-build:
	@echo "构建 Docker 镜像..."
	docker-compose build

docker-up:
	@echo "启动 Docker 容器..."
	docker-compose up -d

docker-down:
	@echo "停止 Docker 容器..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --rmi local

# 数据库
db-init:
	@echo "初始化数据库..."
	psql -U postgres -c "DROP DATABASE IF EXISTS talent_platform;"
	psql -U postgres -c "CREATE DATABASE talent_platform;"
	psql -U postgres -d talent_platform -f backend/database/schema.sql
	psql -U postgres -d talent_platform -f backend/database/mock_data.sql

db-reset:
	@echo "重置数据库数据..."
	psql -U postgres -d talent_platform -f backend/database/mock_data.sql

db-migrate:
	@echo "运行数据库迁移..."
	psql -U postgres -d talent_platform -f backend/database/schema.sql

# 清理
clean:
	@echo "清理构建产物..."
	@for service in $(SERVICES); do \
		rm -f backend/$$service/$$service; \
	done
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.cache

# 代码质量
lint:
	@echo "运行代码检查..."
	@for service in $(SERVICES); do \
		echo "检查 $$service..."; \
		cd backend/$$service && golangci-lint run && cd ../..; \
	done
	cd frontend && npm run lint

fmt:
	@echo "格式化代码..."
	@for service in $(SERVICES); do \
		cd backend/$$service && go fmt ./... && cd ../..; \
	done

# 安装依赖
deps:
	@echo "安装依赖..."
	@for service in $(SERVICES); do \
		cd backend/$$service && go mod tidy && cd ../..; \
	done
	cd frontend && npm install

# 生成 API 文档
docs:
	@echo "生成 API 文档..."
	@echo "API 文档位于 backend/docs/swagger.yaml"
