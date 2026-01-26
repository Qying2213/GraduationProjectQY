# Docker 部署指南

## 1. Dockerfile 配置

### 1.1 后端服务通用 Dockerfile

每个后端服务目录下已有 Dockerfile，以下是优化版本：

```dockerfile
# backend/user-service/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 复制 common 模块
COPY ../common ../common

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
```

### 1.2 前端 Dockerfile

```dockerfile
# frontend/Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# 生产阶段
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

## 2. Docker Compose 配置

### 2.1 开发环境 (docker-compose.dev.yml)

```yaml
version: '3.8'

services:
  # 数据库
  postgres:
    image: pgvector/pgvector:pg16
    container_name: talent-postgres
    environment:
      POSTGRES_USER: qinyang
      POSTGRES_DB: talent_platform
      POSTGRES_HOST_AUTH_METHOD: trust
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql
      - ./backend/database/mock_data.sql:/docker-entrypoint-initdb.d/02-mock_data.sql
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U qinyang -d talent_platform"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis
  redis:
    image: redis:7-alpine
    container_name: talent-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # Gateway
  gateway:
    build:
      context: ./backend/gateway
      dockerfile: Dockerfile
    container_name: talent-gateway
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # User Service
  user-service:
    build:
      context: ./backend/user-service
      dockerfile: Dockerfile
    container_name: talent-user-service
    ports:
      - "8081:8081"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # Job Service
  job-service:
    build:
      context: ./backend/job-service
      dockerfile: Dockerfile
    container_name: talent-job-service
    ports:
      - "8082:8082"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # Interview Service
  interview-service:
    build:
      context: ./backend/interview-service
      dockerfile: Dockerfile
    container_name: talent-interview-service
    ports:
      - "8083:8083"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # Resume Service
  resume-service:
    build:
      context: ./backend/resume-service
      dockerfile: Dockerfile
    container_name: talent-resume-service
    ports:
      - "8084:8084"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform
      - COZE_BASE_URL=https://api.coze.cn
      - COZE_TOKEN=${COZE_TOKEN}
      - COZE_WORKFLOW_ID=${COZE_WORKFLOW_ID}
      - ARK_API_KEY=${ARK_API_KEY}
    volumes:
      - resume_uploads:/app/uploads

  # Message Service
  message-service:
    build:
      context: ./backend/message-service
      dockerfile: Dockerfile
    container_name: talent-message-service
    ports:
      - "8085:8085"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # Talent Service
  talent-service:
    build:
      context: ./backend/talent-service
      dockerfile: Dockerfile
    container_name: talent-talent-service
    ports:
      - "8086:8086"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform

  # Recommendation Service
  recommendation-service:
    build:
      context: ./backend/recommendation-service
      dockerfile: Dockerfile
    container_name: talent-recommendation-service
    ports:
      - "8087:8087"
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=qinyang
      - DB_NAME=talent_platform
      - ARK_API_KEY=${ARK_API_KEY}
      - VOLC_ENDPOINT=https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal
      - VOLC_MODEL_ID=doubao-embedding-vision-251215

  # Frontend
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: talent-frontend
    ports:
      - "80:80"
    depends_on:
      - gateway

volumes:
  postgres_data:
  redis_data:
  resume_uploads:
```

### 2.2 生产环境 (docker-compose.prod.yml)

```yaml
version: '3.8'

services:
  postgres:
    image: pgvector/pgvector:pg16
    container_name: talent-postgres
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER:-qinyang}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: talent_platform
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - talent-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-qinyang} -d talent_platform"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: talent-redis
    restart: always
    networks:
      - talent-network

  gateway:
    image: ${DOCKER_USERNAME}/talent-platform-gateway:latest
    container_name: talent-gateway
    restart: always
    ports:
      - "8080:8080"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  user-service:
    image: ${DOCKER_USERNAME}/talent-platform-user-service:latest
    container_name: talent-user-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  job-service:
    image: ${DOCKER_USERNAME}/talent-platform-job-service:latest
    container_name: talent-job-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  interview-service:
    image: ${DOCKER_USERNAME}/talent-platform-interview-service:latest
    container_name: talent-interview-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  resume-service:
    image: ${DOCKER_USERNAME}/talent-platform-resume-service:latest
    container_name: talent-resume-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
      - COZE_BASE_URL=https://api.coze.cn
      - COZE_TOKEN=${COZE_TOKEN}
      - COZE_WORKFLOW_ID=${COZE_WORKFLOW_ID}
      - ARK_API_KEY=${ARK_API_KEY}
    volumes:
      - resume_uploads:/app/uploads
    networks:
      - talent-network

  message-service:
    image: ${DOCKER_USERNAME}/talent-platform-message-service:latest
    container_name: talent-message-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  talent-service:
    image: ${DOCKER_USERNAME}/talent-platform-talent-service:latest
    container_name: talent-talent-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
    networks:
      - talent-network

  recommendation-service:
    image: ${DOCKER_USERNAME}/talent-platform-recommendation-service:latest
    container_name: talent-recommendation-service
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      - DB_HOST=postgres
      - DB_USER=${DB_USER:-qinyang}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=talent_platform
      - ARK_API_KEY=${ARK_API_KEY}
      - VOLC_ENDPOINT=https://ark.cn-beijing.volces.com/api/v3/embeddings/multimodal
      - VOLC_MODEL_ID=doubao-embedding-vision-251215
    networks:
      - talent-network

  frontend:
    image: ${DOCKER_USERNAME}/talent-platform-frontend:latest
    container_name: talent-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - gateway
    networks:
      - talent-network

  # Nginx 反向代理 (可选)
  nginx:
    image: nginx:alpine
    container_name: talent-nginx
    restart: always
    ports:
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - frontend
      - gateway
    networks:
      - talent-network

networks:
  talent-network:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
  resume_uploads:
```

## 3. 常用 Docker 命令

```bash
# 构建所有服务
docker-compose -f docker-compose.dev.yml build

# 启动所有服务
docker-compose -f docker-compose.dev.yml up -d

# 查看日志
docker-compose -f docker-compose.dev.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.dev.yml logs -f resume-service

# 停止所有服务
docker-compose -f docker-compose.dev.yml down

# 重启特定服务
docker-compose -f docker-compose.dev.yml restart resume-service

# 清理未使用的镜像
docker system prune -a

# 查看容器状态
docker-compose -f docker-compose.dev.yml ps
```

## 4. 环境变量配置

创建 `.env` 文件：

```bash
# .env
DOCKER_USERNAME=your-dockerhub-username
DB_USER=qinyang
DB_PASSWORD=your-secure-password
COZE_TOKEN=pat_xxx
COZE_WORKFLOW_ID=7583886563420373019
ARK_API_KEY=your-ark-api-key
```

## 5. 健康检查

```bash
# 检查所有服务健康状态
for port in 8080 8081 8082 8083 8084 8085 8086 8087; do
    echo "Port $port: $(curl -s http://localhost:$port/health | jq -r '.status')"
done
```
