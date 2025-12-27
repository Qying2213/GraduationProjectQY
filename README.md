# 🎯 智能人才招聘平台

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.4-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Coze_AI-集成-FF6B6B?style=for-the-badge" alt="Coze AI">
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
</p>

<p align="center">
  基于 <strong>Go 微服务架构</strong> + <strong>Vue3</strong> + <strong>Coze AI</strong> 的企业级人才招聘管理系统
</p>

---

## 📖 项目简介

智能人才招聘平台是一个功能完善的招聘管理系统，包含**后台管理系统**、**前台求职端**和**AI简历评估系统**三大模块。后端采用 Go 语言微服务架构，前端使用 Vue3 + TypeScript + Element Plus 构建，集成 Coze AI 实现智能简历评估。

### 系统入口

| 模块 | 地址 | 说明 |
|------|------|------|
| 后台管理 | http://localhost:5173/login | HR/管理员使用 |
| 前台求职端 | http://localhost:5173/portal | 求职者浏览职位、投递简历 |
| AI评估系统 | http://localhost:8090 | 智能简历评估 |

---

## ✨ 功能特性

### 后台管理系统
- 🧑‍💼 **人才管理** - 人才库、搜索筛选、标签分类、批量导入导出
- 📋 **职位管理** - 职位发布、状态管理、技能要求配置
- 📄 **简历管理** - 简历上传、AI智能解析、匹配度计算
- 🤖 **智能推荐** - 多维度人岗匹配算法（技能/经验/学历/地点）
- 📅 **面试管理** - 面试安排、日历视图、反馈提交
- 🎯 **招聘看板** - 可视化招聘流程、拖拽式状态变更
- 💬 **消息中心** - 实时通知、面试提醒
- 📊 **数据报表** - 招聘数据可视化、ECharts 图表
- 🔐 **权限管理** - RBAC 角色权限、5种预设角色

### 前台求职端
- 🏠 **首页** - 职位搜索、热门分类、推荐职位
- 🔍 **职位列表** - 多条件筛选、职位搜索
- 📝 **职位详情** - 职位信息、一键投递
- 👤 **个人中心** - 我的投递、简历管理
- 🔑 **用户认证** - 求职者注册、登录

### AI 简历评估系统
- 🤖 **智能评估** - 基于 Coze AI 的多维度简历评分
- 📊 **六维评分** - 年龄、经验、学历、公司、技术、项目
- 🎯 **JD匹配** - 自动计算简历与职位的匹配度
- 📝 **评估报告** - 生成详细的评估报告和录用建议
- 🔔 **钉钉推送** - 评估结果自动推送到钉钉群
- 📈 **批量处理** - 支持批量简历评估，实时进度显示

---

## 🚀 快速启动

### 环境要求
- Node.js 18+
- Go 1.21+
- PostgreSQL 14+
- Python 3.8+ (AI评估系统需要)

### 一键启动（推荐）

```bash
# 1. 初始化数据库
psql -U postgres -c "CREATE DATABASE talent_platform;"
psql -U postgres -d talent_platform -f backend/database/schema.sql
psql -U postgres -d talent_platform -f backend/database/mock_data.sql

# 2. 一键启动所有服务
chmod +x start-all.sh
./start-all.sh
```

### 或使用 Make 命令

```bash
make dev-all    # 启动所有服务（含 AI 评估系统）
make dev        # 只启动前后端（不含评估系统）
```

### 或使用 Docker

```bash
docker-compose up -d
```

详细启动说明请查看 [快速启动文档](docs/QUICKSTART.md)

---

## 🔑 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 超级管理员 |
| hr_zhang | password123 | HR主管 |
| viewer_test | password123 | 只读用户 |

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                     客户端 (Client)                          │
│         后台管理 + 前台求职端 (Vue3 + Element Plus)           │
└─────────────────────────────┬───────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────┐
│                    API 网关 (Gateway)                        │
│                Go + Gin + JWT (Port: 8080)                   │
└─────────────────────────────┬───────────────────────────────┘
                              │
    ┌──────────┬──────────┬───┴───┬──────────┬──────────┬──────────┐
    ▼          ▼          ▼       ▼          ▼          ▼          ▼
┌────────┐┌────────┐┌────────┐┌────────┐┌────────┐┌────────┐┌────────┐
│ 用户   ││ 人才   ││ 职位   ││ 简历   ││ 推荐   ││ 面试   ││ AI评估 │
│ 服务   ││ 服务   ││ 服务   ││ 服务   ││ 服务   ││ 服务   ││ 服务   │
│ :8081  ││ :8082  ││ :8083  ││ :8084  ││ :8085  ││ :8087  ││ :8090  │
└────┬───┘└────┬───┘└────┬───┘└────┬───┘└────┬───┘└────┬───┘└────┬───┘
     └─────────┴─────────┴────┬────┴─────────┴─────────┘         │
                              ▼                                   │
                    ┌─────────────────┐              ┌────────────▼────┐
                    │   PostgreSQL    │              │    Coze AI      │
                    └─────────────────┘              └─────────────────┘
```

---

## 🛠️ 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3.4 + TypeScript + Element Plus + Pinia + ECharts |
| 后端 | Go 1.21 + Gin + GORM |
| AI | Coze AI 工作流 + 钉钉机器人 |
| 数据库 | PostgreSQL 14 + SQLite (评估系统) |
| 测试 | Vitest (前端) + Go Test (后端) |
| 部署 | Docker + Docker Compose |

---

## 📁 项目结构

```
├── backend/                    # 后端微服务
│   ├── gateway/               # API 网关 (:8080)
│   ├── user-service/          # 用户服务 (:8081)
│   ├── talent-service/        # 人才服务 (:8082)
│   ├── job-service/           # 职位服务 (:8083)
│   ├── resume-service/        # 简历服务 (:8084)
│   ├── recommendation-service/ # 推荐服务 (:8085)
│   ├── message-service/       # 消息服务 (:8086)
│   ├── interview-service/     # 面试服务 (:8087)
│   ├── evaluator-service/     # AI评估服务 (:8090) ⭐
│   │   ├── internal/
│   │   │   ├── ai/           # AI 客户端
│   │   │   ├── api/          # API 路由和处理器
│   │   │   ├── config/       # 配置管理
│   │   │   ├── dingtalk/     # 钉钉集成
│   │   │   ├── models/       # 数据模型
│   │   │   ├── script/       # Python 脚本
│   │   │   ├── service/      # 业务逻辑
│   │   │   └── thirdparty/   # 第三方集成 (Coze)
│   │   └── static/           # 静态文件
│   ├── common/                # 公共模块
│   └── database/              # 数据库脚本
├── frontend/                   # 前端应用
│   └── src/
│       ├── views/portal/      # 前台求职端页面
│       └── views/             # 后台管理页面
├── docs/                       # 项目文档
└── docker-compose.yml          # Docker 编排
```

---

## 🧪 运行测试

```bash
# 前端测试
cd frontend
npm run test              # 运行测试
npm run test:coverage     # 覆盖率报告

# 后端测试
cd backend/resume-service && go test ./... -v
cd backend/user-service && go test ./... -v
```

---

## 📚 文档

| 文档 | 说明 |
|------|------|
| [快速启动](docs/QUICKSTART.md) | 3步启动项目 |
| [系统架构](docs/ARCHITECTURE.md) | 架构图、ER图、流程图 |
| [部署文档](docs/DEPLOYMENT.md) | 本地/Docker/生产部署 |

---

## 📄 License

MIT License

---

## 🎓 关于

本项目为大学毕业设计作品，展示了完整的企业级招聘系统开发能力。

**技术亮点**：
- 微服务架构设计（9个独立服务）
- Coze AI 智能简历评估
- 多维度人岗匹配算法
- 钉钉机器人集成
- RBAC 权限控制
- 完整的单元测试

---

## 📊 服务端口一览

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 5173 | Vue3 开发服务器 |
| API网关 | 8080 | 统一入口 |
| 用户服务 | 8081 | 认证/用户管理 |
| 人才服务 | 8082 | 人才库 |
| 职位服务 | 8083 | 职位管理 |
| 简历服务 | 8084 | 简历解析 |
| 推荐服务 | 8085 | AI推荐 |
| 消息服务 | 8086 | 通知 |
| 面试服务 | 8087 | 面试管理 |
| AI评估服务 | 8090 | Coze AI 简历评估 |
| PostgreSQL | 5432 | 数据库 |
