# 智能人才招聘管理平台 - 部署文档

> 📖 返回 [项目首页](../README.md) | 相关文档：[系统架构](ARCHITECTURE.md) | [快速启动](QUICKSTART.md)

---

## 1. 部署架构

```
                    ┌─────────────┐
                    │   Nginx     │
                    │  (反向代理)  │
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
          ▼                ▼                ▼
    ┌──────────┐    ┌──────────┐    ┌──────────┐
    │ Frontend │    │ Gateway  │    │ Services │
    │  :5173   │    │  :8080   │    │ :8081-90 │
    └──────────┘    └──────────┘    └──────────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
          ▼                ▼                ▼
    ┌──────────┐    ┌──────────┐    ┌──────────┐
    │PostgreSQL│    │   ES     │    │  Kibana  │
    │  :5432   │    │  :9200   │    │  :5601   │
    └──────────┘    └──────────┘    └──────────┘
```

---

## 2. Docker 部署

### 2.1 使用 Docker Compose

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 2.2 docker-compose.yml 配置

```yaml
version: '3.8'

services:
  # PostgreSQL
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: talent_platform
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/database/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql

  # Elasticsearch
  elasticsearch:
    image: elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  # Kibana
  kibana:
    image: kibana:8.11.0
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch

  # Frontend
  frontend:
    build: ./frontend
    ports:
      - "5173:80"

volumes:
  postgres_data:
  es_data:
```

---

## 3. Nginx 配置

### 3.1 生产环境配置

```nginx
upstream gateway {
    server localhost:8080;
}

server {
    listen 80;
    server_name your-domain.com;

    # 前端静态资源
    location / {
        root /var/www/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://gateway;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket 支持
    location /ws {
        proxy_pass http://gateway;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 4. 环境变量配置

### 4.1 后端服务

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=talent_platform

# Elasticsearch 配置
ES_HOST=localhost
ES_PORT=9200

# JWT 配置
JWT_SECRET=your-secret-key

# Coze AI 配置 (evaluator-service)
COZE_API_KEY=your-coze-api-key
COZE_WORKFLOW_ID=your-workflow-id
```

### 4.2 前端配置

```bash
# .env.production
VITE_API_BASE_URL=/api
```

---

## 5. 生产环境检查清单

- [ ] 数据库备份策略
- [ ] 日志轮转配置
- [ ] SSL 证书配置
- [ ] 防火墙规则
- [ ] 监控告警配置
- [ ] 限流配置
- [ ] 敏感信息加密

---

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [📖 项目首页](../README.md) | 项目概述 |
| [📐 系统架构](ARCHITECTURE.md) | 架构设计 |
| [🚀 快速启动](QUICKSTART.md) | 本地开发 |
| [🧪 测试指南](TEST_GUIDE.md) | 测试方法 |
