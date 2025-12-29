# 🎯 智能人才招聘管理平台

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.4-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Coze_AI-集成-FF6B6B?style=for-the-badge" alt="Coze AI">
</p>

<p align="center">
  基于 <strong>Go 微服务架构</strong> + <strong>Vue3</strong> + <strong>Coze AI</strong> 的企业级智能人才招聘管理系统
</p>

---

## 📖 项目简介

智能人才招聘管理平台是一个功能完善的招聘管理系统，采用**微服务架构**设计，包含**企业管理后台**和**求职者门户**两大模块。后端采用 Go 语言开发6个独立微服务，前端使用 Vue3 + TypeScript + Element Plus 构建，集成 Coze AI 实现智能简历评估。

### 系统入口

| 模块 | 地址 | 说明 |
|------|------|------|
| 企业管理后台 | http://localhost:5173/login | HR/管理员使用 |
| 求职者门户 | http://localhost:5173/portal | 求职者浏览职位、投递简历 |

---

## ✨ 功能特性

### 企业管理后台
- 📊 **仪表板** - 数据概览、趋势图表、快捷操作
- 🧑‍💼 **人才管理** - 人才库、搜索筛选、标签分类
- � **职位管理**理 - 职位发布、状态管理、技能要求配置
- � ***简历管理** - 简历上传、状态流转、AI智能评估
- 🤖 **智能推荐** - 人岗匹配算法
- 📅 **面试日历** - 面试安排、日程管理、反馈记录
- 🎯 **招聘看板** - 可视化招聘流程
- 💬 **消息中心** - 站内消息、系统通知
- � ***数据报表** - 招聘数据可视化
- � **权据限管理** - RBAC 角色权限控制
- ⚙️ **系统设置** - 系统配置、操作日志

### 求职者门户
- 🏠 **首页** - 平台介绍、热门职位推荐
- 🔍 **职位列表** - 多条件筛选、职位搜索
- 🏢 **企业招聘** - 企业信息展示、在招职位
- � ***职位详情** - 职位信息、一键投递
- � ***个人中心** - 我的投递、简历管理

### AI 简历评估
- 🤖 **智能评估** - 基于 Coze AI 的多维度简历评分
- 🎯 **JD匹配** - 自动计算简历与职位的匹配度
- 📝 **评估报告** - 生成详细的评估报告和录用建议

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                      前端展示层                              │
│  ┌─────────────────────┐  ┌─────────────────────────────┐  │
│  │   企业管理后台        │  │      求职者门户              │  │
│  │  (Vue3 + Element)    │  │   (Vue3 + Element)          │  │
│  └─────────────────────┘  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                        Vite Proxy / Nginx
                              │
┌─────────────────────────────────────────────────────────────┐
│                      后端服务层 (Go + Gin)                   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │用户服务   │ │职位服务   │ │面试服务   │ │简历服务   │      │
│  │ :8081    │ │ :8082    │ │ :8083    │ │ :8084    │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
│  ┌──────────┐ ┌──────────┐                                 │
│  │消息服务   │ │人才服务   │                                 │
│  │ :8085    │ │ :8086    │                                 │
│  └──────────┘ └──────────┘                                 │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              PostgreSQL 数据库                        │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      AI服务层                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Coze 工作流平台                          │   │
│  │         简历评估 | 人岗匹配 | 智能推荐                 │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ 技术栈

| 层级 | 技术 |
|------|------|
| **前端** | Vue 3 + TypeScript + Element Plus + Pinia + Vue Router + ECharts + Axios |
| **构建** | Vite |
| **后端** | Go 1.21 + Gin + GORM + JWT |
| **数据库** | PostgreSQL 14 |
| **AI** | Coze AI 工作流 |
| **测试** | Vitest (前端) + Shell脚本 (后端API) |
| **部署** | Docker + Nginx |

---

## 🚀 快速启动

### 环境要求
- Node.js 18+
- Go 1.21+
- PostgreSQL 14+

### 1. 克隆项目
```bash
git clone <repository-url>
cd graduate
```

### 2. 初始化数据库
```bash
# 创建数据库
psql -U qinyang -c "CREATE DATABASE talent_platform;"

# 导入表结构和模拟数据
cd backend/database
./import_mock_data.sh
```

### 3. 启动后端服务
```bash
# 一键启动6个微服务（会打开6个终端窗口）
./start-backend.sh
```

### 4. 启动前端
```bash
cd frontend
npm install
npm run dev
```

### 5. 访问系统
- 管理后台: http://localhost:5173/login
- 求职者门户: http://localhost:5173/portal

---

## 🔑 测试账号

| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 超级管理员 |
| hr_zhang | password123 | HR主管 |
| interviewer_chen | password123 | 面试官 |

---

## 📊 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 5173 | Vue3 开发服务器 |
| 用户服务 | 8081 | 认证、用户管理 |
| 职位服务 | 8082 | 职位CRUD |
| 面试服务 | 8083 | 面试安排 |
| 简历服务 | 8084 | 简历管理、AI评估 |
| 消息服务 | 8085 | 消息通知 |
| 人才服务 | 8086 | 人才库管理 |
| PostgreSQL | 5432 | 数据库 |

---

## 📁 项目结构

```
graduate/
├── frontend/                    # 前端项目
│   ├── src/
│   │   ├── api/                # API接口封装
│   │   ├── components/         # 公共组件
│   │   ├── views/              # 页面视图
│   │   │   ├── portal/        # 求职者门户页面
│   │   │   ├── dashboard/     # 仪表板
│   │   │   ├── talents/       # 人才管理
│   │   │   ├── jobs/          # 职位管理
│   │   │   └── ...
│   │   ├── store/              # Pinia状态管理
│   │   ├── router/             # 路由配置
│   │   └── utils/              # 工具函数
│   └── package.json
│
├── backend/                     # 后端微服务
│   ├── user-service/           # 用户服务 :8081
│   ├── job-service/            # 职位服务 :8082
│   ├── interview-service/      # 面试服务 :8083
│   ├── resume-service/         # 简历服务 :8084
│   ├── message-service/        # 消息服务 :8085
│   ├── talent-service/         # 人才服务 :8086
│   ├── common/                 # 公共模块
│   ├── database/               # 数据库脚本
│   └── test_api.sh             # API测试脚本
│
├── docs/                        # 项目文档
│   ├── ARCHITECTURE.md         # 架构文档
│   ├── QUICKSTART.md           # 快速启动
│   └── DEPLOYMENT.md           # 部署文档
│
├── issue/                       # 开题报告等文档
├── start-backend.sh            # 后端启动脚本
├── docker-compose.yml          # Docker编排
└── README.md                   # 项目说明
```

---

## 🧪 运行测试

```bash
# 后端API测试（74个测试用例）
cd backend
./test_api.sh

# 前端单元测试（81个测试用例）
cd frontend
npm run test
```

---

## 📚 文档

| 文档 | 说明 |
|------|------|
| [快速启动](docs/QUICKSTART.md) | 详细启动步骤 |
| [系统架构](docs/ARCHITECTURE.md) | 架构设计、ER图、流程图 |
| [系统设计](docs/SYSTEM_DESIGN.md) | 系统设计文档（毕业设计） |
| [数据库设计](docs/DATABASE_DESIGN.md) | 数据库设计文档 |
| [部署文档](docs/DEPLOYMENT.md) | Docker/Nginx部署 |
| [测试指南](docs/TEST_GUIDE.md) | 测试用例与方法 |
| [代码规范](docs/CODE_GUIDE.md) | 代码规范与目录结构 |

---

## 🎓 关于

本项目为大学毕业设计作品 —— **基于微服务架构的智能人才招聘管理平台设计与实现**

**技术亮点**：
- 微服务架构设计（6个独立服务）
- 前后端分离，RESTful API
- Coze AI 智能简历评估
- RBAC 权限控制
- 完整的单元测试覆盖

---

## 📄 License

MIT License
