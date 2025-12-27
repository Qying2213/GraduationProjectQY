# 快速启动指南

## 环境要求

| 软件 | 版本 | 说明 |
|------|------|------|
| Node.js | 18+ | 前端运行环境 |
| Go | 1.21+ | 后端运行环境 |
| PostgreSQL | 14+ | 数据库 |

---

## 一、数据库初始化

### 1. 创建数据库

```bash
# 使用你的数据库用户
psql -U qinyang -c "CREATE DATABASE talent_platform;"
```

### 2. 导入数据

```bash
cd backend/database
chmod +x import_mock_data.sh
./import_mock_data.sh
```

或手动导入：

```bash
psql -U qinyang -d talent_platform -f schema.sql
psql -U qinyang -d talent_platform -f mock_data.sql
psql -U qinyang -d talent_platform -f mock_data_2_talents.sql
psql -U qinyang -d talent_platform -f mock_data_3_resumes.sql
psql -U qinyang -d talent_platform -f mock_data_4_interviews.sql
psql -U qinyang -d talent_platform -f mock_data_5_messages.sql
```

---

## 二、启动后端服务

### 方式1：一键启动（推荐）

```bash
# 在项目根目录执行
chmod +x start-backend.sh
./start-backend.sh
```

这会自动打开6个终端窗口，分别运行6个微服务。

### 方式2：手动启动

打开6个终端，分别执行：

```bash
# 终端1 - 用户服务
cd backend/user-service && go run main.go

# 终端2 - 职位服务
cd backend/job-service && go run main.go

# 终端3 - 面试服务
cd backend/interview-service && go run main.go

# 终端4 - 简历服务
cd backend/resume-service && go run main.go

# 终端5 - 消息服务
cd backend/message-service && go run main.go

# 终端6 - 人才服务
cd backend/talent-service && go run main.go
```

### 验证服务状态

```bash
# 运行API测试
cd backend
./test_api.sh
```

预期结果：74个测试全部通过

---

## 三、启动前端

```bash
cd frontend
npm install
npm run dev
```

---

## 四、访问系统

| 入口 | 地址 | 账号 |
|------|------|------|
| 管理后台 | http://localhost:5173/login | admin / password123 |
| 求职者门户 | http://localhost:5173/portal | 无需登录 |

---

## 五、服务端口一览

| 服务 | 端口 |
|------|------|
| 前端 | 5173 |
| 用户服务 | 8081 |
| 职位服务 | 8082 |
| 面试服务 | 8083 |
| 简历服务 | 8084 |
| 消息服务 | 8085 |
| 人才服务 | 8086 |
| PostgreSQL | 5432 |

---

## 六、常见问题

### Q: 端口被占用怎么办？

```bash
# 查看占用端口的进程
lsof -i :8081

# 杀掉进程
kill -9 <PID>

# 或一键杀掉所有后端服务
pkill -f "go run main.go"
```

### Q: 数据库连接失败？

检查 `backend/*/main.go` 中的数据库配置：

```go
dsn := "host=localhost user=qinyang dbname=talent_platform port=5432 sslmode=disable"
```

确保用户名和数据库名正确。

### Q: 前端API请求404？

确保所有6个后端服务都已启动，可以用测试脚本验证：

```bash
cd backend && ./test_api.sh
```
