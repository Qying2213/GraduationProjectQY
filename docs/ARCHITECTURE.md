# 智能人才招聘平台 - 系统架构文档

## 1. 系统架构图

### 1.1 整体架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              客户端层 (Client Layer)                          │
├─────────────────────────────────┬───────────────────────────────────────────┤
│       求职者端 (Portal)          │           后台管理端 (Admin)                │
│   Vue3 + TypeScript + Vite      │      Vue3 + TypeScript + Element Plus     │
└─────────────────────────────────┴───────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              网关层 (Gateway Layer)                          │
│                                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   路由转发   │  │   认证鉴权   │  │   限流熔断   │  │   日志记录   │        │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
│                           API Gateway (Go + Gin)                            │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              服务层 (Service Layer)                          │
│                                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐    │
│  │ User Service │  │ Job Service  │  │Talent Service│  │Resume Service│    │
│  │   用户服务    │  │   职位服务    │  │   人才服务    │  │   简历服务    │    │
│  │  :8081       │  │  :8082       │  │  :8083       │  │  :8084       │    │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘    │
│                                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐                      │
│  │Interview Svc │  │ Message Svc  │  │Recommend Svc │                      │
│  │   面试服务    │  │   消息服务    │  │   推荐服务    │                      │
│  │  :8085       │  │  :8086       │  │  :8087       │                      │
│  └──────────────┘  └──────────────┘  └──────────────┘                      │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              数据层 (Data Layer)                             │
│                                                                             │
│  ┌──────────────────────┐  ┌──────────────────────┐  ┌─────────────────┐   │
│  │     PostgreSQL       │  │        Redis         │  │   File Storage  │   │
│  │      主数据库         │  │      缓存/会话        │  │    文件存储      │   │
│  └──────────────────────┘  └──────────────────────┘  └─────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 微服务架构详图

```mermaid
graph TB
    subgraph Client["客户端"]
        A1[求职者Web端]
        A2[HR管理端]
        A3[移动端]
    end

    subgraph Gateway["API网关"]
        B1[路由转发]
        B2[JWT认证]
        B3[限流控制]
        B4[日志记录]
    end

    subgraph Services["微服务集群"]
        C1[用户服务<br/>User Service]
        C2[职位服务<br/>Job Service]
        C3[人才服务<br/>Talent Service]
        C4[简历服务<br/>Resume Service]
        C5[面试服务<br/>Interview Service]
        C6[消息服务<br/>Message Service]
        C7[推荐服务<br/>Recommendation Service]
    end

    subgraph Data["数据存储"]
        D1[(PostgreSQL)]
        D2[(Redis)]
        D3[文件存储]
    end

    A1 --> B1
    A2 --> B1
    A3 --> B1
    
    B1 --> B2
    B2 --> B3
    B3 --> B4
    
    B4 --> C1
    B4 --> C2
    B4 --> C3
    B4 --> C4
    B4 --> C5
    B4 --> C6
    B4 --> C7
    
    C1 --> D1
    C2 --> D1
    C3 --> D1
    C4 --> D1
    C4 --> D3
    C5 --> D1
    C6 --> D1
    C6 --> D2
    C7 --> D1
    C7 --> D2
```

---

## 2. 数据库 ER 图

### 2.1 核心实体关系图

```mermaid
erDiagram
    USERS ||--o{ TALENTS : manages
    USERS ||--o{ JOBS : creates
    USERS ||--o{ MESSAGES : sends
    USERS ||--o{ INTERVIEWS : conducts
    
    TALENTS ||--o{ RESUMES : has
    TALENTS ||--o{ APPLICATIONS : submits
    
    JOBS ||--o{ APPLICATIONS : receives
    JOBS ||--o{ INTERVIEWS : schedules
    
    APPLICATIONS ||--o{ INTERVIEWS : leads_to
    
    INTERVIEWS ||--o{ INTERVIEW_FEEDBACKS : has

    USERS {
        int id PK
        string username UK
        string email UK
        string password
        string role
        string avatar
        string phone
        string department
        string position
        string status
        timestamp created_at
        timestamp updated_at
    }

    TALENTS {
        int id PK
        string name
        string email
        string phone
        array skills
        int experience
        string education
        string status
        array tags
        int user_id FK
        string location
        string salary
        text summary
        string gender
        int age
        timestamp created_at
        timestamp updated_at
    }

    JOBS {
        int id PK
        string title
        text description
        array requirements
        string salary
        string location
        string type
        string status
        int created_by FK
        string department
        string level
        array skills
        array benefits
        timestamp created_at
        timestamp updated_at
    }

    RESUMES {
        int id PK
        int talent_id FK
        int job_id FK
        string file_path
        string file_name
        string status
        int match_score
        text parsed_data
        timestamp created_at
        timestamp updated_at
    }

    APPLICATIONS {
        int id PK
        int job_id FK
        int talent_id FK
        int resume_id FK
        string status
        text cover_letter
        text notes
        timestamp created_at
        timestamp updated_at
    }

    INTERVIEWS {
        int id PK
        int candidate_id FK
        string candidate_name
        int position_id FK
        string position
        string type
        string date
        string time
        int duration
        int interviewer_id FK
        string interviewer
        string method
        string location
        string status
        text notes
        text feedback
        int rating
        int created_by FK
        timestamp created_at
        timestamp updated_at
    }

    INTERVIEW_FEEDBACKS {
        int id PK
        int interview_id FK
        int interviewer_id FK
        int rating
        text strengths
        text weaknesses
        text comments
        string recommendation
        timestamp created_at
    }

    MESSAGES {
        int id PK
        int sender_id FK
        int receiver_id FK
        string type
        string title
        text content
        boolean is_read
        timestamp created_at
    }
```

### 2.2 数据表说明

| 表名 | 说明 | 主要字段 |
|------|------|----------|
| users | 用户表 | 存储系统用户信息，包括管理员、HR、面试官等 |
| talents | 人才表 | 存储候选人信息，包括技能、经验、教育背景等 |
| jobs | 职位表 | 存储招聘职位信息，包括要求、薪资、福利等 |
| resumes | 简历表 | 存储简历文件信息和解析结果 |
| applications | 申请表 | 记录候选人的职位申请 |
| interviews | 面试表 | 存储面试安排信息 |
| interview_feedbacks | 面试反馈表 | 存储面试官的评价反馈 |
| messages | 消息表 | 存储系统消息和用户通知 |

---

## 3. 业务流程图

### 3.1 招聘流程

```mermaid
flowchart TD
    A[HR发布职位] --> B[职位上线]
    B --> C{求职者投递}
    C -->|投递简历| D[简历筛选]
    D -->|通过| E[安排面试]
    D -->|不通过| F[发送拒绝通知]
    E --> G[初试]
    G -->|通过| H[复试]
    G -->|不通过| F
    H -->|通过| I[终面/HR面]
    H -->|不通过| F
    I -->|通过| J[发送Offer]
    I -->|不通过| F
    J --> K{候选人接受?}
    K -->|接受| L[入职]
    K -->|拒绝| M[流程结束]
    F --> M
    L --> M
```

### 3.2 简历解析流程

```mermaid
flowchart LR
    A[上传简历] --> B[文件存储]
    B --> C[格式检测]
    C --> D{文件类型}
    D -->|PDF| E[PDF解析]
    D -->|Word| F[Word解析]
    E --> G[文本提取]
    F --> G
    G --> H[AI分析]
    H --> I[信息提取]
    I --> J[技能识别]
    J --> K[经验匹配]
    K --> L[生成结构化数据]
    L --> M[存储解析结果]
```

### 3.3 智能推荐流程

```mermaid
flowchart TD
    A[获取职位需求] --> B[提取关键技能]
    B --> C[查询人才库]
    C --> D[技能匹配计算]
    D --> E[经验匹配计算]
    E --> F[学历匹配计算]
    F --> G[综合评分]
    G --> H[排序筛选]
    H --> I[返回推荐列表]
```

---

## 4. 技术栈说明

### 4.1 前端技术栈

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

### 4.2 后端技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.21+ | 后端语言 |
| Gin | 1.9+ | Web框架 |
| GORM | 1.25+ | ORM框架 |
| JWT | - | 身份认证 |
| PostgreSQL | 14+ | 主数据库 |
| Redis | 6+ | 缓存/会话 |

### 4.3 部署技术

| 技术 | 用途 |
|------|------|
| Docker | 容器化 |
| Docker Compose | 容器编排 |
| Nginx | 反向代理 |
| GitHub Actions | CI/CD |

---

## 5. 安全设计

### 5.1 认证授权

- JWT Token 认证
- 基于角色的访问控制 (RBAC)
- 细粒度权限管理

### 5.2 数据安全

- 密码 bcrypt 加密存储
- 敏感数据脱敏显示
- SQL 注入防护
- XSS 攻击防护

### 5.3 传输安全

- HTTPS 加密传输
- API 请求签名
- 请求频率限制
