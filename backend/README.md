# 智能人才运营平台 - 后端说明

## 概述

本目录包含毕业设计项目的后端部分，采用 Go + Gin 微服务架构实现。当前后端由 8 个业务微服务、1 个 API 网关和 1 个 AI 评估服务组成，覆盖用户、职位、简历、人才、推荐、面试、消息、日志等核心业务。

## 服务清单

| 服务 | 端口 | 主要职责 |
|---|---:|---|
| gateway | 8080 | API 网关、统一路由、限流、操作日志接入 |
| user-service | 8081 | 用户认证、JWT、RBAC、用户资料 |
| job-service | 8082 | 职位发布、职位搜索、职位统计 |
| interview-service | 8083 | 面试安排、反馈记录、面试状态流转 |
| resume-service | 8084 | 简历上传、OCR、AI 评估、风控、评估链路追踪 |
| message-service | 8085 | 站内消息、未读统计、通知能力 |
| talent-service | 8086 | 人才库管理、人才搜索与筛选 |
| recommendation-service | 8087 | 推荐统计、人岗匹配、推荐解释 |
| log-service | 8088 | 操作日志查询、日志统计、审计追踪 |
| evaluator-service | 8090 | 独立 AI 评估入口、钉钉推送、候选人评估管理 |

## 目录结构

```text
backend/
├── gateway/
├── user-service/
├── job-service/
├── interview-service/
├── resume-service/
├── message-service/
├── talent-service/
├── recommendation-service/
├── log-service/
├── evaluator-service/
├── common/
├── databaseSQL/
├── cozeWorkflow/
└── README.md
```

## 核心 AI 链路

当前简历智能评估链路为：

```text
PDF / DOC / DOCX
  -> OCR / 文本提取
  -> Embedding 向量化
  -> RAG 检索增强
  -> Coze 工作流评估
  -> 结构化结果入库
```

说明：

- `resume-service` 负责在线评估主链路
- `cozeWorkflow/workflow/` 存放导出的 Coze 工作流定义
- 当前工作流已支持同时接收 `resume_file` 与 `resume_text`

## 环境变量

推荐直接复制：

```bash
cp backend/.env.example backend/.env
```

关键变量说明：

| 变量 | 说明 |
|---|---|
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` | PostgreSQL 连接配置 |
| `REDIS_HOST` / `REDIS_PORT` | Redis 配置 |
| `ES_URL` | Elasticsearch 地址 |
| `JWT_SECRET` | JWT 密钥 |
| `COZE_BASE_URL` | Coze API 基础地址 |
| `COZE_TOKEN` | `resume-service` 评估工作流使用的 Coze Token |
| `COZE_WORKFLOW_ID` | `resume-service` 评估工作流 ID |
| `COZE_API_TOKEN` | 其他 AI 归因/推荐链路兼容使用的 Coze Token |
| `COZE_ATTRIBUTION_WORKFLOW_ID` | 推荐归因工作流 ID |
| `ARK_API_KEY` | 火山引擎 Embedding API Key |
| `VOLC_ENDPOINT` | Embedding 接口地址 |
| `VOLC_MODEL_ID` | Embedding 模型 ID |

## 本地启动

### 方式一：一键启动

项目根目录执行：

```bash
./start-all.sh
```

停止：

```bash
./stop-all.sh
```

### 方式二：逐个服务启动

```bash
cd backend/gateway && go run main.go
cd backend/user-service && go run main.go
cd backend/job-service && go run main.go
cd backend/interview-service && go run main.go
cd backend/resume-service && go run main.go
cd backend/message-service && go run main.go
cd backend/talent-service && go run main.go
cd backend/recommendation-service && go run main.go
cd backend/log-service && go run main.go
cd backend/evaluator-service && go run ./cmd/server
```

## 数据库初始化

推荐在项目根目录执行：

```bash
./init-db.sh
```

手动执行时：

```bash
createdb talent_platform
psql -d talent_platform -f backend/databaseSQL/schema.sql
psql -d talent_platform -f backend/databaseSQL/init_data.sql
```

## 常用测试

### 后端单元/服务测试

```bash
cd backend/resume-service && go test ./...
cd backend/recommendation-service && go test ./...
cd backend/interview-service && go test ./...
```

### 全链路接口测试

```bash
cd ztest
./run_all_tests.sh
```

### 毕业设计专项验收

```bash
python3 ztest/test_graduation_acceptance.py
```

## 当前交付状态

根据 2026-04-01 的验收材料，当前后端相关结论如下：

- 全链路接口测试 33/33 通过
- 毕业设计专项验收 34/34 通过
- 核心读接口压测达到开题性能目标
- AI 评估链路已完成 OCR、Embedding、RAG 与 Coze 工作流整合

可直接参考：

- [专项验收报告](/Users/qinyang/Desktop/GraduationProjectQY/ztest/graduation_acceptance_report_20260401_115415.md)
- [交付验收报告](/Users/qinyang/Desktop/GraduationProjectQY/docs/毕业设计交付验收报告_20260401.md)
- [详细测试报告](/Users/qinyang/Desktop/GraduationProjectQY/docs/毕业设计详细测试报告_20260401.md)
