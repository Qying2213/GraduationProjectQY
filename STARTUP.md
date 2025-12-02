# 智能人才运营平台 - 快速启动指南

## 环境要求

- **Node.js** 18+
- **Go** 1.21+
- **PostgreSQL** 14+

## 启动步骤

### 第一步：初始化数据库

```bash
# 1. 确保 PostgreSQL 服务已启动

# 2. 创建数据库并导入数据
psql -U postgres -c "DROP DATABASE IF EXISTS talent_platform;"
psql -U postgres -c "CREATE DATABASE talent_platform;"
psql -U postgres -d talent_platform -f backend/database/schema.sql
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

验证数据导入成功：
```bash
psql -U postgres -d talent_platform -c "SELECT COUNT(*) FROM users;"
# 应该返回 9
```

### 第二步：启动后端服务

**方式一：使用启动脚本（推荐）**

```bash
cd /Users/qinyang/Desktop/abc
./start-backend.sh
```

**方式二：手动启动各服务**

在 7 个不同的终端窗口中分别执行：

```bash
# 终端1 - API网关 (8080)
cd backend/gateway && go mod tidy && go run main.go

# 终端2 - 用户服务 (8081)
cd backend/user-service && go mod tidy && go run main.go

# 终端3 - 人才服务 (8082)
cd backend/talent-service && go mod tidy && go run main.go

# 终端4 - 职位服务 (8083)
cd backend/job-service && go mod tidy && go run main.go

# 终端5 - 简历服务 (8084)
cd backend/resume-service && go mod tidy && go run main.go

# 终端6 - 推荐服务 (8085)
cd backend/recommendation-service && go mod tidy && go run main.go

# 终端7 - 消息服务 (8086)
cd backend/message-service && go mod tidy && go run main.go
```

### 第三步：启动前端

```bash
cd frontend

# 首次运行需要安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 访问地址

| 服务 | 地址 |
|------|------|
| 前端应用 | http://localhost:3000 |
| API网关 | http://localhost:8080 |

## 测试账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 拥有所有权限 |
| hr_zhang | password123 | HR主管 | HR总监 |
| hr_li | password123 | 招聘专员 | 招聘专员 |
| tech_chen | password123 | 面试官 | 技术总监 |
| tech_liu | password123 | 面试官 | 前端负责人 |
| viewer_test | password123 | 只读用户 | 仅查看权限 |

## 服务端口一览

| 服务名称 | 端口 | 说明 |
|----------|------|------|
| gateway | 8080 | API网关，统一入口 |
| user-service | 8081 | 用户认证与管理 |
| talent-service | 8082 | 人才库管理 |
| job-service | 8083 | 职位管理 |
| resume-service | 8084 | 简历解析与管理 |
| recommendation-service | 8085 | AI推荐服务 |
| message-service | 8086 | 消息通知服务 |
| interview-service | 8087 | 面试管理服务 |

## 数据库配置

默认配置（可在各服务的环境变量中修改）：

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=talent_platform
```

## 常见问题

### Q: 数据库连接失败？
A: 检查 PostgreSQL 是否启动，以及用户名密码是否正确。

### Q: 端口被占用？
A: 使用 `lsof -i :端口号` 查看占用进程，使用 `kill -9 PID` 终止。

### Q: 前端启动报错？
A: 删除 `node_modules` 文件夹后重新执行 `npm install`。

### Q: 后端服务启动失败？
A: 确保 Go 环境已安装，执行 `go mod tidy` 更新依赖。

## 重置数据

如需重置所有数据：

```bash
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

## 停止服务

- 前端：在终端按 `Ctrl + C`
- 后端：关闭各个终端窗口，或使用 `pkill -f "go run main.go"`
