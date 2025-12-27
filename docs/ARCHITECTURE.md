# 智能人才招聘管理平台 - 系统架构文档

## 1. 系统架构概述

本系统采用**微服务架构**设计，将业务功能拆分为6个独立的微服务，每个服务独立部署、独立数据库，通过RESTful API进行通信。

### 1.1 架构特点

- **服务独立部署**：每个微服务可独立启动、停止、升级
- **故障隔离**：单个服务故障不影响其他服务
- **技术栈灵活**：各服务可选择最适合的技术方案
- **弹性扩展**：可根据负载对单个服务进行水平扩展

---

## 2. 系统架构图

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                      前端展示层                              │
│  ┌─────────────────────┐  ┌─────────────────────────────┐  │
│  │   企业管理后台        │  │      求职者门户              │  │
│  │  (Vue3 + Element)    │  │   (Vue3 + Element)          │  │
│  │      :5173           │  │        :5173                │  │
│  └─────────────────────┘  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                    Vite Proxy (开发) / Nginx (生产)
                              │
┌─────────────────────────────────────────────────────────────┐
│                   后端服务层 (Go + Gin)                      │
│                                                             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐      │
│  │用户服务   │ │职位服务   │ │面试服务   │ │简历服务   │      │
│  │ :8081    │ │ :8082    │ │ :8083    │ │ :8084    │      │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘      │
│                                                             │
│  ┌──────────┐ ┌──────────┐                                 │
│  │消息服务   │ │人才服务   │                                 │
│  │ :8085    │ │ :8086    │                                 │
│  └──────────┘ └──────────┘                                 │
└─────────────────────────────────────────────────────────────┘
                              │
                            GORM
                              │
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                              │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              PostgreSQL 数据库 :5432                  │   │
│  │   users | jobs | talents | resumes | interviews     │   │
│  │   messages | applications | interview_feedbacks     │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                              │
                          API 调用
                              │
┌─────────────────────────────────────────────────────────────┐
│                      AI服务层                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Coze 工作流平台                          │   │
│  │         简历评估 | 人岗匹配 | 智能推荐                 │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 服务端口配置

| 服务名称 | 端口 | 说明 |
|---------|------|------|
| Frontend | 5173 | 前端开发服务器 |
| User Service | 8081 | 用户认证、权限管理 |
| Job Service | 8082 | 职位CRUD、搜索筛选 |
| Interview Service | 8083 | 面试安排、反馈管理 |
| Resume Service | 8084 | 简历上传、AI评估 |
| Message Service | 8085 | 消息通知、未读统计 |
| Talent Service | 8086 | 人才库管理、搜索 |
| PostgreSQL | 5432 | 主数据库 |

---

## 3. 微服务详细设计

### 3.1 用户服务 (User Service) :8081

**职责**：用户认证、权限管理

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/v1/login | 用户登录 |
| POST | /api/v1/register | 用户注册 |
| GET | /api/v1/users | 获取用户列表 |
| GET | /api/v1/users/:id | 获取用户详情 |
| PUT | /api/v1/users/:id | 更新用户信息 |

### 3.2 职位服务 (Job Service) :8082

**职责**：职位管理、职位统计

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/jobs | 获取职位列表 |
| GET | /api/v1/jobs/:id | 获取职位详情 |
| POST | /api/v1/jobs | 创建职位 |
| PUT | /api/v1/jobs/:id | 更新职位 |
| DELETE | /api/v1/jobs/:id | 删除职位 |
| GET | /api/v1/jobs/stats | 获取职位统计 |

### 3.3 面试服务 (Interview Service) :8083

**职责**：面试安排、反馈管理

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/interviews | 获取面试列表 |
| GET | /api/v1/interviews/:id | 获取面试详情 |
| POST | /api/v1/interviews | 创建面试 |
| PUT | /api/v1/interviews/:id | 更新面试 |
| GET | /api/v1/interviews/stats | 面试统计 |
| GET | /api/v1/interviews/today | 今日面试 |
| GET | /api/v1/interviews/:id/feedback | 获取面试反馈 |

### 3.4 简历服务 (Resume Service) :8084

**职责**：简历管理、AI评估集成

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/resumes | 获取简历列表 |
| GET | /api/v1/resumes/:id | 获取简历详情 |
| POST | /api/v1/resumes | 上传简历 |
| PUT | /api/v1/resumes/:id/status | 更新简历状态 |
| GET | /api/v1/applications | 获取申请列表 |
| POST | /api/v1/applications | 创建申请 |
| GET | /api/v1/ai/config | 检查AI配置 |
| POST | /api/v1/ai/evaluate | AI评估简历 |

### 3.5 消息服务 (Message Service) :8085

**职责**：站内消息、通知管理

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/messages | 获取消息列表 |
| POST | /api/v1/messages | 发送消息 |
| PUT | /api/v1/messages/:id/read | 标记已读 |
| GET | /api/v1/messages/unread-count | 未读数量 |
| DELETE | /api/v1/messages/:id | 删除消息 |

### 3.6 人才服务 (Talent Service) :8086

**职责**：人才库管理、搜索

**API接口**：
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/v1/talents | 获取人才列表 |
| GET | /api/v1/talents/:id | 获取人才详情 |
| POST | /api/v1/talents | 创建人才 |
| PUT | /api/v1/talents/:id | 更新人才 |
| DELETE | /api/v1/talents/:id | 删除人才 |

---

## 4. 数据库设计

### 4.1 ER图

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│    users     │       │    jobs      │       │   talents    │
├──────────────┤       ├──────────────┤       ├──────────────┤
│ id           │       │ id           │       │ id           │
│ username     │──┐    │ title        │    ┌──│ name         │
│ email        │  │    │ description  │    │  │ email        │
│ password     │  │    │ salary       │    │  │ skills       │
│ role         │  │    │ location     │    │  │ experience   │
│ status       │  │    │ status       │    │  │ education    │
└──────────────┘  │    │ created_by   │──┐ │  │ status       │
                  │    └──────────────┘  │ │  └──────────────┘
                  │           │          │ │         │
                  │           │          │ │         │
                  │    ┌──────▼──────┐   │ │  ┌──────▼──────┐
                  │    │ applications │   │ │  │   resumes   │
                  │    ├─────────────┤   │ │  ├─────────────┤
                  │    │ id          │   │ │  │ id          │
                  │    │ job_id      │───┘ │  │ talent_id   │
                  │    │ talent_id   │─────┘  │ file_path   │
                  │    │ resume_id   │───────►│ status      │
                  │    │ status      │        │ match_score │
                  │    └─────────────┘        └─────────────┘
                  │
                  │    ┌──────────────┐       ┌──────────────┐
                  │    │  interviews  │       │  messages    │
                  │    ├──────────────┤       ├──────────────┤
                  └───►│ interviewer  │       │ sender_id    │
                       │ candidate_id │       │ receiver_id  │
                       │ position_id  │       │ title        │
                       │ status       │       │ content      │
                       │ feedback     │       │ is_read      │
                       └──────────────┘       └──────────────┘
```

### 4.2 核心数据表

| 表名 | 说明 |
|------|------|
| users | 用户表（管理员、HR、面试官） |
| talents | 人才表（候选人信息） |
| jobs | 职位表（招聘职位） |
| resumes | 简历表（简历文件和解析结果） |
| applications | 申请表（职位申请记录） |
| interviews | 面试表（面试安排） |
| interview_feedbacks | 面试反馈表 |
| messages | 消息表（站内消息） |

---

## 5. 业务流程

### 5.1 招聘流程

```
发布职位 → 候选人投递 → 简历筛选 → AI评估 → 安排面试 → 面试反馈 → 录用决策
    │           │           │          │          │           │          │
    ▼           ▼           ▼          ▼          ▼           ▼          ▼
 Job Service  Resume    Resume     Coze AI   Interview  Interview   Message
              Service   Service              Service    Service     Service
```

### 5.2 AI评估流程

```
上传简历 → 文件存储 → 调用Coze API → AI分析 → 返回评估结果 → 更新简历状态
```

---

## 6. 技术栈

### 6.1 前端

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue.js | 3.x | 前端框架 |
| TypeScript | 5.x | 类型安全 |
| Vite | 5.x | 构建工具 |
| Element Plus | 2.x | UI组件库 |
| Pinia | 2.x | 状态管理 |
| Vue Router | 4.x | 路由管理 |
| ECharts | 5.x | 数据可视化 |
| Axios | 1.x | HTTP客户端 |
| Vitest | 1.x | 单元测试 |

### 6.2 后端

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端语言 |
| Gin | 1.9+ | Web框架 |
| GORM | 1.25+ | ORM框架 |
| JWT | - | 身份认证 |
| bcrypt | - | 密码加密 |
| PostgreSQL | 14+ | 主数据库 |

### 6.3 AI服务

| 技术 | 用途 |
|------|------|
| Coze API | AI工作流平台 |
| Coze Workflow | 简历智能评估 |

### 6.4 部署

| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 容器编排 |
| Nginx | 反向代理/API网关 |

---

## 7. Nginx 网关配置（生产环境）

```nginx
upstream user-service {
    server localhost:8081;
}
upstream job-service {
    server localhost:8082;
}
upstream interview-service {
    server localhost:8083;
}
upstream resume-service {
    server localhost:8084;
}
upstream message-service {
    server localhost:8085;
}
upstream talent-service {
    server localhost:8086;
}

server {
    listen 80;
    server_name localhost;

    # 前端静态资源
    location / {
        root /var/www/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 路由转发
    location /api/v1/login { proxy_pass http://user-service; }
    location /api/v1/register { proxy_pass http://user-service; }
    location /api/v1/users { proxy_pass http://user-service; }
    location /api/v1/jobs { proxy_pass http://job-service; }
    location /api/v1/interviews { proxy_pass http://interview-service; }
    location /api/v1/resumes { proxy_pass http://resume-service; }
    location /api/v1/applications { proxy_pass http://resume-service; }
    location /api/v1/ai { proxy_pass http://resume-service; }
    location /api/v1/messages { proxy_pass http://message-service; }
    location /api/v1/talents { proxy_pass http://talent-service; }
}
```

---

## 8. 安全设计

### 8.1 认证授权
- JWT Token 认证
- 基于角色的访问控制 (RBAC)
- Token 过期机制

### 8.2 数据安全
- 密码 bcrypt 加密存储
- SQL 注入防护（GORM参数化查询）
- XSS 攻击防护

### 8.3 文件安全
- 文件类型白名单
- 文件大小限制
- 存储路径隔离

---

## 9. 测试覆盖

### 9.1 后端测试
- API测试脚本: `backend/test_api.sh`
- 测试用例: 74个
- 覆盖: 全部6个微服务

### 9.2 前端测试
- 测试框架: Vitest
- 测试用例: 81个
- 覆盖: API、组件、工具函数、状态管理

---

## 10. 目录结构

```
graduate/
├── frontend/                    # 前端项目
│   └── src/
│       ├── api/                # API接口
│       ├── components/         # 组件
│       ├── views/              # 页面
│       ├── store/              # 状态管理
│       └── router/             # 路由
│
├── backend/                     # 后端微服务
│   ├── user-service/           # 用户服务 :8081
│   ├── job-service/            # 职位服务 :8082
│   ├── interview-service/      # 面试服务 :8083
│   ├── resume-service/         # 简历服务 :8084
│   ├── message-service/        # 消息服务 :8085
│   ├── talent-service/         # 人才服务 :8086
│   ├── common/                 # 公共模块
│   └── database/               # 数据库脚本
│
├── docs/                        # 文档
└── issue/                       # 开题报告
```
