# 智能人才运营平台 - 快速启动指南

## 环境要求

| 软件 | 版本 | 说明 |
|-----|------|------|
| Go | 1.21+ | 后端开发语言 |
| Node.js | 18+ | 前端运行环境 |
| PostgreSQL | 14+ | 数据库 |
| npm | 9+ | 包管理器 |

## 一、数据库配置

### 1.1 创建数据库

```bash
# 登录PostgreSQL
psql -U postgres

# 创建数据库和用户
CREATE DATABASE talent_platform;
CREATE USER qinyang WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE talent_platform TO qinyang;
\q
```

### 1.2 初始化表结构

```bash
cd backend/database
psql -U qinyang -d talent_platform -f schema.sql
```

### 1.3 导入模拟数据（可选）

```bash
cd backend/database
./import_mock_data.sh
```

---

## 二、后端启动

### 2.1 一键启动（推荐）

```bash
cd backend
./start-all.sh
```

### 2.2 启动脚本命令

```bash
./start-all.sh start    # 启动所有服务
./start-all.sh stop     # 停止所有服务
./start-all.sh restart  # 重启所有服务
./start-all.sh status   # 查看服务状态
```

### 2.3 手动启动（单个服务）

```bash
# 在不同终端分别启动
cd backend/user-service && go run main.go      # 端口 8081
cd backend/job-service && go run main.go       # 端口 8082
cd backend/interview-service && go run main.go # 端口 8083
cd backend/resume-service && go run main.go    # 端口 8084
cd backend/message-service && go run main.go   # 端口 8085
cd backend/talent-service && go run main.go    # 端口 8086
```

### 2.4 服务端口

| 服务 | 端口 | 说明 |
|-----|------|------|
| user-service | 8081 | 用户认证、权限管理 |
| job-service | 8082 | 职位CRUD、搜索筛选 |
| interview-service | 8083 | 面试安排、反馈管理 |
| resume-service | 8084 | 简历上传、AI评估 |
| message-service | 8085 | 消息通知、未读统计 |
| talent-service | 8086 | 人才库管理、搜索 |

---

## 三、前端启动

### 3.1 安装依赖

```bash
cd frontend
npm install
```

### 3.2 启动开发服务器

```bash
npm run dev
```

访问地址: http://localhost:5173

### 3.3 构建生产版本

```bash
npm run build
```

---

## 四、测试账号

| 用户名 | 密码 | 角色 |
|-------|------|------|
| admin | password123 | 管理员 |
| hr_zhang | password123 | HR |
| hr_wang | password123 | HR |
| interviewer1 | password123 | 面试官 |

---

## 五、运行测试

### 5.1 后端API测试

```bash
cd backend
./test_api.sh
```

测试覆盖: 74个测试用例

### 5.2 前端测试

```bash
cd frontend
npm test -- --run
```

测试覆盖: 81个测试用例

---

## 六、AI评估配置（可选）

### 6.1 获取Coze API Key

1. 访问 [Coze平台](https://www.coze.cn)
2. 创建工作流
3. 获取API Token和Workflow ID

### 6.2 配置环境变量

```bash
cd backend/resume-service
cp .env.example .env

# 编辑 .env 文件
COZE_API_TOKEN=pat_xxxxxxxxxxxxxxxx
COZE_WORKFLOW_ID=7xxxxxxxxxxxxxxxxxx
```

---

## 七、常见问题

### Q1: 数据库连接失败

检查PostgreSQL服务是否启动:
```bash
brew services start postgresql  # macOS
sudo systemctl start postgresql # Linux
```

### Q2: 端口被占用

查看并关闭占用端口的进程:
```bash
lsof -i :8081
kill -9 <PID>
```

### Q3: Go模块下载失败

设置Go代理:
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### Q4: 前端依赖安装失败

清除缓存重试:
```bash
rm -rf node_modules package-lock.json
npm install
```

---

## 八、项目结构

```
graduate/
├── frontend/           # 前端项目 (Vue3 + TypeScript)
├── backend/            # 后端项目 (Go + Gin)
│   ├── user-service/       # 用户服务 :8081
│   ├── job-service/        # 职位服务 :8082
│   ├── interview-service/  # 面试服务 :8083
│   ├── resume-service/     # 简历服务 :8084
│   ├── message-service/    # 消息服务 :8085
│   ├── talent-service/     # 人才服务 :8086
│   ├── database/           # 数据库脚本
│   ├── start-all.sh        # 一键启动脚本
│   └── test_api.sh         # API测试脚本
└── docs/               # 文档
```
