# 智能人才招聘管理平台 - 快速启动指南

> 📖 返回 [项目首页](../README.md) | 相关文档：[系统架构](ARCHITECTURE.md) | [部署文档](DEPLOYMENT.md) | [测试指南](TEST_GUIDE.md)

---

## 1. 环境要求

| 软件 | 版本 | 必需 | 说明 |
|------|------|------|------|
| Node.js | 18+ | ✅ | 前端运行环境 |
| Go | 1.21+ | ✅ | 后端运行环境 |
| PostgreSQL | 14+ | ✅ | 主数据库 |
| Elasticsearch | 8.x | ⚪ | 日志功能（可选） |
| Docker | 20+ | ⚪ | 容器部署（可选） |

---

## 2. 快速启动步骤

### 2.1 克隆项目

```bash
git clone <repository-url>
cd talent-platform
```

### 2.2 初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE talent_platform;"

# 导入表结构
psql -U postgres -d talent_platform -f backend/database/schema.sql

# 导入模拟数据（可选，包含测试账号和示例数据）
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

### 2.3 安装后端依赖

```bash
cd backend
chmod +x setup-deps.sh
./setup-deps.sh
```

> 💡 `setup-deps.sh` 会自动配置国内镜像加速

### 2.4 启动后端服务

```bash
# 一键启动所有微服务
chmod +x start-backend.sh
./start-backend.sh
```

这会启动以下服务：
| 服务 | 端口 |
|------|------|
| API Gateway | 8080 |
| user-service | 8081 |
| job-service | 8082 |
| interview-service | 8083 |
| resume-service | 8084 |
| message-service | 8085 |
| talent-service | 8086 |
| recommendation-service | 8087 |
| log-service | 8088 |
| evaluator-service | 8090 |

### 2.5 启动前端

```bash
cd frontend
npm install
npm run dev
```

### 2.6 访问系统

| 入口 | 地址 | 说明 |
|------|------|------|
| 管理后台 | http://localhost:5173/login | HR/管理员登录 |
| 求职者门户 | http://localhost:5173/portal | 求职者浏览职位 |
| 数据大屏 | http://localhost:5173/data-screen | 数据可视化 |
| AI评估系统 | http://localhost:8090 | 独立AI评估入口 |

---

## 3. 测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 所有权限 |
| hr_zhang | password123 | HR主管 | 招聘全流程 |
| hr_li | password123 | 招聘专员 | 日常招聘 |
| tech_chen | password123 | 面试官 | 面试评估 |
| viewer_test | password123 | 只读用户 | 仅查看 |


---

## 4. 可选：启动 Elasticsearch

如需使用日志功能：

```bash
cd backend
chmod +x start-es.sh
./start-es.sh
```

或使用 Docker：

```bash
docker-compose up -d elasticsearch kibana
```

访问 Kibana：http://localhost:5601

---

## 5. 验证服务状态

### 5.1 检查后端服务

```bash
# 检查各服务健康状态
curl http://localhost:8080/health  # Gateway
curl http://localhost:8081/health  # User Service
curl http://localhost:8082/health  # Job Service
curl http://localhost:8083/health  # Interview Service
curl http://localhost:8084/health  # Resume Service
curl http://localhost:8085/health  # Message Service
curl http://localhost:8086/health  # Talent Service
curl http://localhost:8087/health  # Recommendation Service (如果启动)
curl http://localhost:8088/health  # Log Service
```

### 5.2 运行 API 测试

```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

---

## 6. 开发模式

### 6.1 单独启动某个服务

```bash
cd backend/user-service
go run main.go
```

### 6.2 前端开发

```bash
cd frontend
npm run dev      # 开发模式
npm run build    # 构建生产版本
npm run test     # 运行测试
```

### 6.3 数据库配置

各服务支持环境变量配置数据库连接：

```bash
export DB_HOST=localhost
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=talent_platform
export DB_PORT=5432
```

---

## 7. 常见问题

### Q: 端口被占用？

```bash
# 查看占用端口的进程
lsof -i :8081

# 杀掉进程
kill -9 <PID>

# 或一键杀掉所有后端服务
pkill -f "go run main.go"
```

### Q: Go 依赖下载慢？

```bash
# 使用国内镜像
export GOPROXY=https://goproxy.cn,direct
```

### Q: 数据库连接失败？

1. 确认 PostgreSQL 服务已启动
2. 确认数据库 `talent_platform` 已创建
3. 检查用户名和密码配置

### Q: 前端启动报错？

```bash
# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

---

## 8. 下一步

- 📐 了解 [系统架构](ARCHITECTURE.md)
- 📋 查看 [系统设计](SYSTEM_DESIGN.md)
- 🗄️ 了解 [数据库设计](DATABASE_DESIGN.md)
- 🐳 查看 [部署文档](DEPLOYMENT.md)
- 🧪 运行 [测试指南](TEST_GUIDE.md)
- 📝 阅读 [代码规范](CODE_GUIDE.md)
