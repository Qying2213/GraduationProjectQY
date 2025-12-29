# 🎯 智能人才招聘管理平台

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Vue-3.4-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/TypeScript-5.x-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/PostgreSQL-14+-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Elasticsearch-8.x-005571?style=for-the-badge&logo=elasticsearch&logoColor=white" alt="Elasticsearch">
  <img src="https://img.shields.io/badge/Coze_AI-集成-FF6B6B?style=for-the-badge" alt="Coze AI">
</p>

<p align="center">
  <strong>🚀 基于 Go 微服务架构 + Vue3 + Coze AI 的企业级智能人才招聘管理系统</strong>
</p>

<p align="center">
  <a href="#-功能特性">功能特性</a> •
  <a href="#-系统架构">系统架构</a> •
  <a href="#-快速启动">快速启动</a> •
  <a href="#-技术栈">技术栈</a> •
  <a href="#-项目文档">项目文档</a>
</p>

---

## 📖 项目简介

智能人才招聘管理平台是一个功能完善的企业级招聘管理系统，采用**微服务架构**设计，包含**企业管理后台**、**求职者门户**和**数据大屏**三大模块。

- 🏢 **后端**：Go 语言开发 **9个独立微服务 + API网关**，支持独立部署和弹性扩展
- � **前端***：Vue3 + TypeScript + Element Plus，响应式设计，支持暗色主题
- 🤖 **AI能力**：集成 Coze AI 实现智能简历评估、人岗匹配推荐
- 📊 **数据分析**：ECharts 数据可视化，实时招聘数据大屏
- � **日志系析统**：Elasticsearch 全文搜索，操作日志可追溯
- �  **权限管理**：RBAC 角色权限控制，5种预设角色


### 🌐 系统入口

| 模块 | 地址 | 说明 |
|------|------|------|
| 🏢 企业管理后台 | http://localhost:5173/login | HR/管理员使用，招聘全流程管理 |
| 👤 求职者门户 | http://localhost:5173/portal | 求职者浏览职位、投递简历、管理个人简历 |
| 📊 数据大屏 | http://localhost:5173/data-screen | 招聘数据可视化展示（全屏独立页面） |
| 🤖 AI评估系统 | http://localhost:8090 | 智能简历评估独立入口（支持钉钉推送） |

---

## ✨ 功能特性

### 🏢 企业管理后台

| 模块 | 功能描述 |
|------|----------|
| 📊 **仪表板** | 核心数据概览、招聘漏斗、渠道统计、趋势图表、快捷操作入口 |
| 🧑‍💼 **人才管理** | 人才库维护、多维度搜索筛选、标签分类、技能匹配 |
| 💼 **职位管理** | 职位发布/编辑/关闭、技能要求配置、招聘进度跟踪、职位统计 |
| 📄 **简历管理** | 简历上传解析、状态流转、AI智能评估、匹配度分析、批量评估 |
| 🤖 **智能推荐** | 基于AI的人岗匹配算法、为职位推荐人才、为人才推荐职位、多维度匹配分析 |
| 📅 **面试日历** | 面试安排、日程管理、面试反馈记录、面试官评分、面试改期/取消 |
| 🎯 **招聘看板** | 可视化招聘流程、拖拽式状态管理、招聘阶段统计 |
| 💬 **消息中心** | 站内消息、系统通知、面试提醒、简历投递提醒 |
| 📈 **数据报表** | 招聘数据统计、漏斗分析、趋势报表、部门招聘进度 |
| 👥 **权限管理** | RBAC角色权限控制、5种预设角色（超级管理员/HR主管/招聘专员/面试官/只读用户） |
| ⚙️ **系统设置** | 系统配置、操作日志查询（ES全文搜索）、审计追踪 |

### 👤 求职者门户

| 模块 | 功能描述 |
|------|----------|
| 🏠 **首页** | 平台介绍、热门职位推荐、企业展示 |
| 🔍 **职位搜索** | 多条件筛选（地点/薪资/经验/学历/职位类型）、关键词搜索 |
| 🏢 **企业招聘** | 企业信息展示、在招职位列表 |
| 📋 **职位详情** | 职位信息、要求说明、一键投递 |
| 👤 **个人中心** | 我的投递记录、简历管理、投递状态跟踪 |


### 🤖 AI 智能能力

| 功能 | 描述 |
|------|------|
| 📊 **智能评估** | 基于 Coze AI 的多维度简历评分（技术能力/项目经验/教育背景/JD匹配度） |
| 🎯 **智能推荐** | 多维度匹配算法（技能匹配50%+经验匹配20%+地理位置15%+学历10%+薪资5%） |
| 📝 **评估报告** | 生成详细的评估报告、优劣势分析、录用建议、候选人对比 |
| 🔄 **批量处理** | 支持批量简历评估，提升筛选效率 |
| 📱 **钉钉推送** | 评估结果自动推送到钉钉群，支持多配置管理 |

---

## 🏗️ 系统架构

### 整体架构图

```
┌────────────────────────────────────────────────────────────────────────────┐
│                              客户端层                                       │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐         │
│  │   企业管理后台    │  │    求职者门户     │  │    数据大屏      │         │
│  │  (Vue3+Element)  │  │  (Vue3+Element)  │  │   (ECharts)     │         │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬────────┘         │
└───────────┼─────────────────────┼─────────────────────┼──────────────────┘
            │                     │                     │
            └─────────────────────┼─────────────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │      API Gateway :8080    │
                    │   (路由分发/限流/日志)     │
                    └─────────────┬─────────────┘
                                  │
┌─────────────────────────────────▼──────────────────────────────────────────┐
│                           微服务层 (Go + Gin)                               │
│                                                                            │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │ 用户服务   │ │ 职位服务   │ │ 面试服务   │ │ 简历服务   │              │
│  │  :8081    │ │  :8082    │ │  :8083    │ │  :8084    │              │
│  │ 认证/权限  │ │ 职位CRUD  │ │ 面试安排   │ │ 简历/AI评估 │              │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘              │
│                                                                            │
│  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐              │
│  │ 消息服务   │ │ 人才服务   │ │ 推荐服务   │ │ 日志服务   │              │
│  │  :8085    │ │  :8086    │ │  :8087    │ │  :8088    │              │
│  │ 站内消息   │ │ 人才库    │ │ 智能匹配   │ │ ES日志    │              │
│  └────────────┘ └────────────┘ └────────────┘ └────────────┘              │
│                                                                            │
│                        ┌────────────────────┐                              │
│                        │   AI评估服务 :8090  │                              │
│                        │ Coze AI/钉钉推送   │                              │
│                        └────────────────────┘                              │
└────────────────────────────────┬───────────────────────────────────────────┘
                                 │
            ┌────────────────────┼────────────────────┐
            │                    │                    │
            ▼                    ▼                    ▼
┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐
│    PostgreSQL     │ │   Elasticsearch   │ │     Coze AI       │
│    (业务数据)      │ │    (日志数据)      │ │   (智能评估)       │
│     :5432         │ │      :9200        │ │    API调用        │
└───────────────────┘ └───────────────────┘ └───────────────────┘
```


### 数据库 ER 图

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│    users    │       │    roles    │       │    jobs     │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │       │ id (PK)     │       │ id (PK)     │
│ username    │──────>│ name        │       │ title       │
│ email       │       │ code        │       │ description │
│ password    │       │ permissions │       │ requirements│
│ role        │       └─────────────┘       │ salary      │
│ department  │                             │ location    │
│ status      │                             │ type        │
└──────┬──────┘                             │ status      │
       │                                    │ skills[]    │
       │                                    │ created_by  │──┐
       │                                    └──────┬──────┘  │
       │                                           │         │
       ▼                                           ▼         │
┌─────────────┐       ┌─────────────┐       ┌─────────────┐  │
│   talents   │       │   resumes   │       │applications │  │
├─────────────┤       ├─────────────┤       ├─────────────┤  │
│ id (PK)     │<──────│ talent_id   │       │ id (PK)     │  │
│ name        │       │ id (PK)     │       │ talent_id   │  │
│ email       │       │ job_id      │       │ job_id      │──┘
│ phone       │       │ file_path   │       │ resume_id   │
│ skills[]    │       │ status      │       │ stage       │
│ experience  │       │ match_score │       │ status      │
│ education   │       │ parse_result│       │ source      │
│ location    │       └─────────────┘       └─────────────┘
│ salary      │
└──────┬──────┘
       │
       ▼
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│  interviews │       │  feedbacks  │       │  messages   │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ id (PK)     │<──────│interview_id │       │ id (PK)     │
│ candidate_id│       │ id (PK)     │       │ sender_id   │
│ position_id │       │interviewer_id│      │ receiver_id │
│ interviewer │       │ rating      │       │ type        │
│ type        │       │ strengths   │       │ title       │
│ date/time   │       │ weaknesses  │       │ content     │
│ method      │       │ comments    │       │ is_read     │
│ status      │       │recommendation│      │ created_at  │
│ feedback    │       └─────────────┘       └─────────────┘
└─────────────┘

┌─────────────────┐
│ operation_logs  │
├─────────────────┤
│ id (PK)         │
│ user_id         │
│ action          │
│ resource_type   │
│ resource_id     │
│ details (JSONB) │
│ ip_address      │
│ created_at      │
└─────────────────┘
```


### 招聘业务流程

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ 发布职位  │───>│ 候选人   │───>│ 简历筛选  │───>│ AI评估   │───>│ 安排面试  │───>│ 录用决策  │
│          │    │ 投递简历  │    │          │    │          │    │          │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
     │               │               │               │               │               │
     ▼               ▼               ▼               ▼               ▼               ▼
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│Job       │    │Resume    │    │Resume    │    │Evaluator │    │Interview │    │Message   │
│Service   │    │Service   │    │Service   │    │Service   │    │Service   │    │Service   │
│:8082     │    │:8084     │    │:8084     │    │:8090     │    │:8083     │    │:8085     │
└──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘    └──────────┘
```

### 智能推荐算法

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          智能推荐匹配算法                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  总匹配度 = 技能匹配(50%) + 经验匹配(20%) + 地理位置(15%)                    │
│           + 学历匹配(10%) + 薪资匹配(5%)                                     │
│                                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ 技能匹配    │  │ 经验匹配    │  │ 位置匹配    │  │ 学历匹配    │        │
│  │ 权重: 50%   │  │ 权重: 20%   │  │ 权重: 15%   │  │ 权重: 10%   │        │
│  ├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤        │
│  │ Go: 1.2x    │  │ junior: 0-2年│ │ 完全匹配:100%│ │ 博士: 100%  │        │
│  │ K8s: 1.3x   │  │ mid: 2-6年  │  │ 同城市群:70%│  │ 硕士: 90%   │        │
│  │ Python: 1.1x│  │ senior: 5-10年│ │ 不匹配: 30% │  │ 本科: 80%   │        │
│  │ Vue: 1.1x   │  │ expert: 8-15年│ │             │  │ 大专: 60%   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│                                                                             │
│  匹配等级: ≥80% 高度匹配 | ≥60% 中等匹配 | ≥40% 基本匹配 | <40% 匹配度较低  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 🛠️ 技术栈

### 前端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue.js | 3.4+ | 渐进式 JavaScript 框架 |
| TypeScript | 5.x | 类型安全的 JavaScript 超集 |
| Vite | 5.x | 下一代前端构建工具 |
| Element Plus | 2.x | Vue3 UI 组件库 |
| Pinia | 2.x | Vue3 状态管理 |
| Vue Router | 4.x | Vue3 官方路由 |
| ECharts | 5.x | 数据可视化图表库 |
| Axios | 1.x | HTTP 请求库 |
| Vitest | 1.x | Vue3 单元测试框架 |
| SCSS | - | CSS 预处理器 |


### 后端技术

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 高性能后端语言 |
| Gin | 1.9+ | 轻量级 Web 框架 |
| GORM | 1.25+ | Go ORM 框架 |
| JWT | - | JSON Web Token 认证 |
| bcrypt | - | 密码加密 |
| Elasticsearch Client | 8.x | ES 日志存储客户端 |

### 数据存储

| 技术 | 版本 | 说明 |
|------|------|------|
| PostgreSQL | 14+ | 主数据库，存储业务数据 |
| Elasticsearch | 8.x | 日志存储与全文搜索 |

### AI 服务

| 技术 | 说明 |
|------|------|
| Coze API | 字节跳动 AI 工作流平台 |
| Coze Workflow | 简历智能评估工作流 |
| 钉钉机器人 | 评估结果推送通知 |

### DevOps

| 技术 | 说明 |
|------|------|
| Docker | 容器化部署 |
| Docker Compose | 容器编排 |
| Nginx | 反向代理/负载均衡 |
| Kibana | 日志可视化 |

---

## 🚀 快速启动

### 环境要求

| 软件 | 版本 | 必需 |
|------|------|------|
| Node.js | 18+ | ✅ |
| Go | 1.21+ | ✅ |
| PostgreSQL | 14+ | ✅ |
| Elasticsearch | 8.x | ⚪ 可选（日志功能） |
| Docker | 20+ | ⚪ 可选（容器部署） |

### 1️⃣ 克隆项目

```bash
git clone <repository-url>
cd talent-platform
```

### 2️⃣ 初始化数据库

```bash
# 创建数据库
psql -U postgres -c "CREATE DATABASE talent_platform;"

# 导入表结构
psql -U postgres -d talent_platform -f backend/database/schema.sql

# 导入模拟数据（可选）
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

### 3️⃣ 安装后端依赖

```bash
cd backend
chmod +x setup-deps.sh
./setup-deps.sh  # 使用国内镜像加速
```

### 4️⃣ 启动后端服务

```bash
# 一键启动所有微服务（会打开多个终端窗口）
chmod +x start-backend.sh
./start-backend.sh
```

### 5️⃣ 启动前端

```bash
cd frontend
npm install
npm run dev
```

### 6️⃣ 访问系统

| 入口 | 地址 |
|------|------|
| 管理后台 | http://localhost:5173/login |
| 求职者门户 | http://localhost:5173/portal |
| 数据大屏 | http://localhost:5173/data-screen |


---

## 🔑 测试账号

| 用户名 | 密码 | 角色 | 权限说明 |
|--------|------|------|----------|
| admin | password123 | 超级管理员 | 所有权限 |
| hr_zhang | password123 | HR主管 | 招聘全流程管理 |
| hr_li | password123 | 招聘专员 | 日常招聘操作 |
| tech_chen | password123 | 面试官 | 面试评估 |
| viewer_test | password123 | 只读用户 | 仅查看 |

---

## 📊 服务端口一览

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端 | 5173 | Vite 开发服务器 |
| API Gateway | 8080 | 统一网关（路由分发/限流/日志） |
| user-service | 8081 | 用户认证、权限管理、用户信息 |
| job-service | 8082 | 职位 CRUD、搜索筛选、职位统计 |
| interview-service | 8083 | 面试安排、反馈管理、面试统计 |
| resume-service | 8084 | 简历管理、状态流转、AI评估接口 |
| message-service | 8085 | 消息通知、站内信、未读统计 |
| talent-service | 8086 | 人才库管理、人才搜索 |
| recommendation-service | 8087 | 智能推荐、人岗匹配算法 |
| log-service | 8088 | 操作日志查询（ES） |
| evaluator-service | 8090 | AI 简历评估、钉钉推送、候选人管理 |
| PostgreSQL | 5432 | 主数据库 |
| Elasticsearch | 9200 | 日志存储 |
| Kibana | 5601 | 日志可视化 |

---

## 📁 项目结构

```
talent-platform/
├── 📂 frontend/                     # 前端项目 (Vue3 + TypeScript)
│   ├── src/
│   │   ├── api/                    # API 接口封装
│   │   ├── components/             # 公共组件
│   │   │   ├── layout/            # 布局组件 (MainLayout/PortalLayout)
│   │   │   ├── common/            # 通用组件
│   │   │   └── charts/            # 图表组件
│   │   ├── views/                  # 页面视图
│   │   │   ├── auth/              # 登录注册
│   │   │   ├── dashboard/         # 仪表板 + 数据大屏
│   │   │   ├── talents/           # 人才管理
│   │   │   ├── jobs/              # 职位管理
│   │   │   ├── resumes/           # 简历管理
│   │   │   ├── recommend/         # 智能推荐
│   │   │   ├── interviews/        # 面试详情
│   │   │   ├── calendar/          # 面试日历
│   │   │   ├── kanban/            # 招聘看板
│   │   │   ├── messages/          # 消息中心
│   │   │   ├── reports/           # 数据报表
│   │   │   ├── portal/            # 求职者门户
│   │   │   ├── profile/           # 个人中心
│   │   │   └── system/            # 系统设置 (权限/日志/设置)
│   │   ├── store/                  # Pinia 状态管理
│   │   ├── router/                 # 路由配置
│   │   ├── types/                  # TypeScript 类型定义
│   │   ├── utils/                  # 工具函数
│   │   └── styles/                 # 全局样式
│   └── package.json
│
├── 📂 backend/                      # 后端微服务 (Go + Gin)
│   ├── gateway/                    # API 网关 :8080 (路由分发/限流/统计)
│   ├── user-service/               # 用户服务 :8081
│   ├── job-service/                # 职位服务 :8082
│   ├── interview-service/          # 面试服务 :8083
│   ├── resume-service/             # 简历服务 :8084 (含AI评估接口)
│   ├── message-service/            # 消息服务 :8085
│   ├── talent-service/             # 人才服务 :8086
│   ├── recommendation-service/     # 推荐服务 :8087 (智能匹配算法)
│   ├── log-service/                # 日志服务 :8088 (ES查询)
│   ├── evaluator-service/          # AI评估服务 :8090 (Coze AI/钉钉)
│   ├── common/                     # 公共模块
│   │   ├── config/                # 配置管理
│   │   ├── elasticsearch/         # ES 客户端
│   │   ├── middleware/            # 中间件 (CORS/日志/认证)
│   │   └── response/              # 统一响应
│   ├── database/                   # 数据库脚本
│   │   ├── schema.sql             # 表结构 (10张核心表)
│   │   └── mock_data*.sql         # 模拟数据
│   └── test_api.sh                 # API 测试脚本
│
├── 📂 docs/                         # 项目文档
│   ├── ARCHITECTURE.md             # 架构设计
│   ├── SYSTEM_DESIGN.md            # 系统设计
│   ├── DATABASE_DESIGN.md          # 数据库设计
│   ├── QUICKSTART.md               # 快速启动
│   ├── DEPLOYMENT.md               # 部署文档
│   ├── TEST_GUIDE.md               # 测试指南
│   └── CODE_GUIDE.md               # 代码规范
│
├── docker-compose.yml              # Docker 编排
├── start-backend.sh                # 后端启动脚本
├── init-db.sh                      # 数据库初始化
├── Makefile                        # 构建命令
└── README.md                       # 项目说明
```


---

## 🧪 测试

### 后端 API 测试

```bash
cd backend
chmod +x test_api.sh
./test_api.sh
```

### 前端单元测试

```bash
cd frontend
npm run test
```

---

## 📚 项目文档

> 点击下方链接查看详细文档

| 文档 | 说明 | 链接 |
|------|------|------|
| 📐 **系统架构** | 微服务架构设计、服务通信、技术选型 | [查看文档](docs/ARCHITECTURE.md) |
| 📋 **系统设计** | 功能模块设计、接口设计、安全设计 | [查看文档](docs/SYSTEM_DESIGN.md) |
| 🗄️ **数据库设计** | 表结构设计、ER图、数据字典、索引设计 | [查看文档](docs/DATABASE_DESIGN.md) |
| 🚀 **快速启动** | 环境配置、安装步骤、服务启动 | [查看文档](docs/QUICKSTART.md) |
| 🐳 **部署文档** | Docker部署、Nginx配置、生产环境 | [查看文档](docs/DEPLOYMENT.md) |
| 🧪 **测试指南** | API测试、功能测试、测试用例 | [查看文档](docs/TEST_GUIDE.md) |
| 📝 **代码规范** | 目录结构、代码风格、开发指南 | [查看文档](docs/CODE_GUIDE.md) |

---

## 🎯 项目亮点

### 架构设计
- ✅ **微服务架构**：9个独立服务 + API网关，职责清晰，支持独立部署和扩展
- ✅ **前后端分离**：RESTful API 设计，前后端完全解耦
- ✅ **统一网关**：API Gateway 实现路由分发、限流、日志记录

### 功能特性
- ✅ **AI 智能评估**：集成 Coze AI，实现简历智能评分和人岗匹配
- ✅ **智能推荐算法**：多维度匹配（技能/经验/位置/学历/薪资）
- ✅ **实时日志**：Elasticsearch 存储，支持全文搜索和可视化
- ✅ **数据可视化**：ECharts 数据大屏，直观展示招聘数据
- ✅ **RBAC 权限**：5种预设角色，细粒度权限控制
- ✅ **钉钉集成**：评估结果自动推送，支持多配置管理

### 工程质量
- ✅ **TypeScript**：前端全面使用 TypeScript，类型安全
- ✅ **代码规范**：统一的代码风格和目录结构
- ✅ **完整测试**：前后端测试覆盖

### 部署运维
- ✅ **Docker 支持**：一键容器化部署
- ✅ **日志追踪**：操作日志可追溯，支持审计

---

## 🔧 常见问题

<details>
<summary><b>Q: 端口被占用怎么办？</b></summary>

```bash
# 查看占用端口的进程
lsof -i :8081

# 杀掉进程
kill -9 <PID>

# 或一键杀掉所有后端服务
pkill -f "go run main.go"
```
</details>

<details>
<summary><b>Q: Go 依赖下载慢？</b></summary>

```bash
# 使用国内镜像
export GOPROXY=https://goproxy.cn,direct

# 或运行脚本
cd backend && ./setup-deps.sh
```
</details>

<details>
<summary><b>Q: 数据库连接失败？</b></summary>

检查数据库配置，确保用户名和数据库名正确。各服务支持环境变量配置：
```bash
export DB_HOST=localhost
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=talent_platform
export DB_PORT=5432
```
</details>

<details>
<summary><b>Q: ES 日志功能不可用？</b></summary>

Elasticsearch 是可选组件，不启动不影响其他功能。如需启用：
```bash
cd backend && ./start-es.sh
```
</details>

---

## 🎓 关于项目

本项目为大学毕业设计作品 —— **基于微服务架构的智能人才招聘管理平台设计与实现**

### 技术创新点
1. **微服务架构实践**：将单体应用拆分为9个独立微服务 + API网关
2. **AI 能力集成**：利用 Coze AI 实现智能简历评估
3. **智能推荐算法**：多维度加权匹配算法实现人岗智能匹配
4. **全链路日志**：Elasticsearch 实现操作日志全文搜索

---

## 📄 License

MIT License

---

<p align="center">
  <sub>Made with ❤️ for graduation project</sub>
</p>
