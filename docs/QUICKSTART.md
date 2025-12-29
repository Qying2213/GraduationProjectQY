# 快速启动指南

## 环境要求

| 软件 | 版本 | 说明 |
|------|------|------|
| Node.js | 18+ | 前端运行环境 |
| Go | 1.21+ | 后端运行环境 |
| PostgreSQL | 14+ | 数据库 |
| Elasticsearch | 8.x | 日志存储（可选） |
| Docker | 20+ | 容器化部署（可选） |

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

## 二、安装后端依赖

```bash
cd backend
chmod +x setup-deps.sh
./setup-deps.sh
```

该脚本会使用国内镜像（goproxy.cn）下载所有Go依赖。

---

## 三、启动 Elasticsearch（可选，用于日志功能）

```bash
cd backend
chmod +x start-es.sh
./start-es.sh
```

或使用 Docker Compose：

```bash
docker-compose up -d elasticsearch kibana
```

访问地址：
- Elasticsearch: http://localhost:9200
- Kibana: http://localhost:5601

> 注意：如果不启动ES，日志功能会自动降级，不影响其他功能正常使用。

---

## 四、启动后端服务

### 方式1：一键启动（推荐）

```bash
# 在项目根目录执行
chmod +x start-backend.sh
./start-backend.sh
```

这会自动打开7个终端窗口，分别运行7个微服务。

### 方式2：手动启动

打开7个终端，分别执行：

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

# 终端7 - 日志服务（需要ES）
cd backend/log-service && go run main.go
```

---

## 五、验证服务状态

```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

预期结果：**90个测试全部通过，通过率100%**

测试覆盖：
- 用户服务（10个）：登录、注册、用户列表
- 职位服务（14个）：职位CRUD、筛选、搜索
- 人才服务（14个）：人才CRUD、搜索
- 简历服务（14个）：简历管理、状态更新
- 面试服务（16个）：面试安排、反馈
- 消息服务（10个）：消息发送、已读标记
- 日志服务（8个）：ES日志查询、统计
- 综合测试（4个）：ES日志验证

---

## 六、启动前端

```bash
cd frontend
npm install
npm run dev
```

---

## 七、访问系统

| 入口 | 地址 | 账号 |
|------|------|------|
| 管理后台 | http://localhost:5173/login | admin / password123 |
| 求职者门户 | http://localhost:5173/portal | 无需登录 |

### 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 超级管理员 |
| hr_zhang | password123 | HR主管 |
| hr_li | password123 | 招聘专员 |
| tech_chen | password123 | 面试官 |
| viewer_test | password123 | 只读用户 |

---

## 八、服务端口一览

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 5173 | Vite开发服务器 |
| user-service | 8081 | 用户认证服务 |
| job-service | 8082 | 职位管理服务 |
| interview-service | 8083 | 面试管理服务 |
| resume-service | 8084 | 简历管理服务 |
| message-service | 8085 | 消息通知服务 |
| talent-service | 8086 | 人才库服务 |
| log-service | 8088 | 日志查询服务 |
| PostgreSQL | 5432 | 数据库 |
| Elasticsearch | 9200 | 日志存储 |
| Kibana | 5601 | 日志可视化 |

---

## 九、日志功能说明

系统集成了 Elasticsearch 日志功能，所有 API 请求都会自动记录到 ES 中。

### 日志查询 API

```bash
# 查询所有日志
curl http://localhost:8088/api/v1/logs

# 按服务筛选
curl http://localhost:8088/api/v1/logs?service=user-service

# 按方法筛选
curl http://localhost:8088/api/v1/logs?method=POST

# 分页查询
curl "http://localhost:8088/api/v1/logs?page=1&page_size=20"

# 获取统计信息
curl http://localhost:8088/api/v1/logs/stats

# 获取服务列表
curl http://localhost:8088/api/v1/logs/services

# 获取操作类型列表
curl http://localhost:8088/api/v1/logs/actions
```

### 使用 Kibana 查看日志

1. 访问 http://localhost:5601
2. 创建索引模式：`operation_logs*`
3. 在 Discover 中查看和搜索日志

---

## 十、常见问题

### Q: 端口被占用怎么办？

```bash
# 查看占用端口的进程
lsof -i :8081

# 杀掉进程
kill -9 <PID>

# 或一键杀掉所有后端服务
pkill -f "go run main.go"
```

### Q: Go依赖下载慢？

使用国内镜像：

```bash
export GOPROXY=https://goproxy.cn,direct
```

或运行 `backend/setup-deps.sh` 脚本。

### Q: 数据库连接失败？

检查 `backend/*/main.go` 中的数据库配置：

```go
dsn := "host=localhost user=qinyang dbname=talent_platform port=5432 sslmode=disable"
```

确保用户名和数据库名正确。

### Q: 前端API请求404？

确保所有7个后端服务都已启动，可以用测试脚本验证：

```bash
cd backend && ./test_api.sh
```

### Q: ES连接失败？

如果没有启动 Elasticsearch，日志功能会自动降级，不影响其他功能正常使用。

启动 ES：
```bash
cd backend && ./start-es.sh
```

### Q: Coze AI评估功能不可用？

需要配置 Coze API Key，在 `backend/resume-service/main.go` 中设置：

```go
cozeAPIKey := os.Getenv("COZE_API_KEY")
```

或设置环境变量：
```bash
export COZE_API_KEY=your_api_key
```
