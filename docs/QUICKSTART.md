# 智能人才招聘平台 - 快速启动指南

## 一、项目简介

智能人才招聘平台是一个完整的招聘管理系统，包含：
- **后台管理系统**：HR 管理人才、职位、面试等
- **前台求职端**：求职者浏览职位、投递简历

技术栈：Go + Gin + PostgreSQL（后端） | Vue3 + TypeScript + Element Plus（前端）

---

## 二、环境要求

| 软件 | 版本 | 必需 |
|------|------|------|
| Node.js | 18.0+ | ✅ |
| Go | 1.21+ | ✅ |
| PostgreSQL | 14.0+ | ✅ |
| Redis | 6.0+ | ❌ 可选 |
| Docker | 20.0+ | ❌ 可选 |

---

## 三、快速启动（3步完成）

### 步骤 1：初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE talent_platform;"

# 导入表结构和测试数据
psql -U postgres -d talent_platform -f backend/database/schema.sql
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

### 步骤 2：启动后端

```bash
# 使用启动脚本（推荐）
chmod +x start-backend.sh
./start-backend.sh
```

或手动启动各服务（需要8个终端）：
```bash
cd backend/gateway && go run main.go              # 端口 8080
cd backend/user-service && go run main.go         # 端口 8081
cd backend/talent-service && go run main.go       # 端口 8082
cd backend/job-service && go run main.go          # 端口 8083
cd backend/resume-service && go run main.go       # 端口 8084
cd backend/recommendation-service && go run main.go # 端口 8085
cd backend/message-service && go run main.go      # 端口 8086
cd backend/interview-service && go run main.go    # 端口 8087
```

### 步骤 3：启动前端

```bash
cd frontend
npm install
npm run dev
```

---

## 四、访问系统

| 入口 | 地址 | 说明 |
|------|------|------|
| 后台管理 | http://localhost:3000/login | HR/管理员使用 |
| 前台求职端 | http://localhost:3000/portal | 求职者使用 |
| API 网关 | http://localhost:8080 | 后端接口 |

---

## 五、测试账号

| 用户名 | 密码 | 角色 | 权限 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 全部权限 |
| hr_zhang | password123 | HR主管 | 人才/职位/面试管理 |
| hr_li | password123 | 招聘专员 | 人才/简历管理 |
| tech_chen | password123 | 面试官 | 面试/反馈 |
| viewer_test | password123 | 只读用户 | 仅查看 |

---

## 六、运行测试

### 前端测试
```bash
cd frontend
npm run test              # 运行测试
npm run test:coverage     # 生成覆盖率报告
```

### 后端测试
```bash
# 简历解析模块
cd backend/resume-service && go test ./parser/... -v

# 用户服务
cd backend/user-service && go test ./handlers/... -v

# 面试服务
cd backend/interview-service && go test ./handlers/... -v

# 推荐服务
cd backend/recommendation-service && go test ./handlers/... -v
```

---

## 七、Docker 部署（可选）

```bash
# 一键启动
docker-compose up -d

# 查看状态
docker-compose ps

# 停止服务
docker-compose down
```

---

## 八、服务端口一览

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 3000 | Vue3 应用 |
| API网关 | 8080 | 统一入口 |
| 用户服务 | 8081 | 认证/用户管理 |
| 人才服务 | 8082 | 人才库 |
| 职位服务 | 8083 | 职位管理 |
| 简历服务 | 8084 | 简历解析 |
| 推荐服务 | 8085 | AI推荐 |
| 消息服务 | 8086 | 通知 |
| 面试服务 | 8087 | 面试管理 |
| PostgreSQL | 5432 | 数据库 |

---

## 九、常见问题

### Q1: 数据库连接失败
```bash
# 检查 PostgreSQL 是否运行
pg_isready

# 检查数据库是否存在
psql -U postgres -c "\l" | grep talent_platform
```

### Q2: 端口被占用
```bash
# 查看端口占用
lsof -i :8080

# 终止进程
kill -9 <PID>
```

### Q3: 前端依赖安装失败
```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### Q4: 登录失败
```bash
# 重新导入测试数据
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

---

## 十、项目文档

| 文档 | 路径 | 说明 |
|------|------|------|
| 系统架构 | docs/ARCHITECTURE.md | 架构图、ER图、流程图 |
| 部署文档 | docs/DEPLOYMENT.md | 详细部署指南 |
| API文档 | backend/docs/swagger.yaml | 接口说明 |

---

**版本**: v1.0  
**更新日期**: 2024年12月
