# 数据库设计文档

## 1. 数据库概述

### 1.1 数据库选型
- **主数据库**：PostgreSQL 14+
- **日志存储**：Elasticsearch 8.x
- **字符集**：UTF-8

### 1.2 数据库信息
- **数据库名**：talent_platform
- **用户名**：qinyang
- **端口**：5432

---

## 2. ER 图

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    users    │     │    roles    │     │    jobs     │
├─────────────┤     ├─────────────┤     ├─────────────┤
│ id (PK)     │     │ id (PK)     │     │ id (PK)     │
│ username    │────>│ name        │     │ title       │
│ email       │     │ code        │     │ description │
│ password    │     │ permissions │     │ salary      │
│ role (FK)   │     └─────────────┘     │ location    │
│ status      │                         │ created_by  │──┐
└─────────────┘                         └─────────────┘  │
       │                                       │         │
       │                                       │         │
       ▼                                       ▼         │
┌─────────────┐     ┌─────────────┐     ┌─────────────┐  │
│   talents   │     │   resumes   │     │ applications│  │
├─────────────┤     ├─────────────┤     ├─────────────┤  │
│ id (PK)     │<────│ talent_id   │     │ id (PK)     │  │
│ name        │     │ id (PK)     │     │ resume_id   │──┘
│ email       │     │ file_path   │     │ job_id      │
│ skills      │     │ status      │     │ status      │
│ experience  │     │ ai_score    │     │ created_at  │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │
       │                   │
       ▼                   ▼
┌─────────────┐     ┌─────────────┐
│  interviews │     │  feedbacks  │
├─────────────┤     ├─────────────┤
│ id (PK)     │<────│interview_id │
│ candidate_id│     │ id (PK)     │
│ position_id │     │ score       │
│ interviewer │     │ comment     │
│ date/time   │     │ result      │
│ status      │     └─────────────┘
└─────────────┘
       │
       ▼
┌─────────────┐
│  messages   │
├─────────────┤
│ id (PK)     │
│ sender_id   │
│ receiver_id │
│ type        │
│ content     │
│ is_read     │
└─────────────┘
```

---

## 3. 数据表设计

### 3.1 用户表 (users)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 用户ID |
| username | VARCHAR(50) | UNIQUE, NOT NULL | 用户名 |
| email | VARCHAR(100) | UNIQUE, NOT NULL | 邮箱 |
| password | VARCHAR(255) | NOT NULL | 密码(bcrypt) |
| role | VARCHAR(20) | NOT NULL | 角色代码 |
| avatar | VARCHAR(255) | | 头像URL |
| phone | VARCHAR(20) | | 手机号 |
| department | VARCHAR(50) | | 部门 |
| position | VARCHAR(50) | | 职位 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |
| real_name | VARCHAR(50) | | 真实姓名 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer',
    avatar VARCHAR(255),
    phone VARCHAR(20),
    department VARCHAR(50),
    position VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    real_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### 3.2 角色表 (roles)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 角色ID |
| name | VARCHAR(50) | NOT NULL | 角色名称 |
| code | VARCHAR(20) | UNIQUE, NOT NULL | 角色代码 |
| description | TEXT | | 角色描述 |
| permissions | TEXT[] | | 权限列表 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |

```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    description TEXT,
    permissions TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3.3 职位表 (jobs)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 职位ID |
| title | VARCHAR(100) | NOT NULL | 职位名称 |
| description | TEXT | | 职位描述 |
| requirements | TEXT[] | | 任职要求 |
| salary | VARCHAR(50) | | 薪资范围 |
| location | VARCHAR(50) | | 工作地点 |
| type | VARCHAR(20) | | 工作类型 |
| status | VARCHAR(20) | DEFAULT 'open' | 状态 |
| created_by | INTEGER | REFERENCES users(id) | 创建人 |
| department | VARCHAR(50) | | 所属部门 |
| level | VARCHAR(20) | | 职级 |
| skills | TEXT[] | | 技能要求 |
| benefits | TEXT[] | | 福利待遇 |
| headcount | INTEGER | DEFAULT 1 | 招聘人数 |
| urgent | BOOLEAN | DEFAULT FALSE | 是否紧急 |
| deadline | DATE | | 截止日期 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    requirements TEXT[],
    salary VARCHAR(50),
    location VARCHAR(50),
    type VARCHAR(20),
    status VARCHAR(20) DEFAULT 'open',
    created_by INTEGER REFERENCES users(id),
    department VARCHAR(50),
    level VARCHAR(20),
    skills TEXT[],
    benefits TEXT[],
    headcount INTEGER DEFAULT 1,
    urgent BOOLEAN DEFAULT FALSE,
    deadline DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_jobs_department ON jobs(department);
```

### 3.4 人才表 (talents)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 人才ID |
| name | VARCHAR(50) | NOT NULL | 姓名 |
| email | VARCHAR(100) | UNIQUE | 邮箱 |
| phone | VARCHAR(20) | | 手机号 |
| skills | TEXT[] | | 技能列表 |
| experience | INTEGER | | 工作年限 |
| education | VARCHAR(50) | | 学历 |
| status | VARCHAR(20) | DEFAULT 'active' | 状态 |
| location | VARCHAR(50) | | 所在地 |
| salary | VARCHAR(50) | | 期望薪资 |
| source | VARCHAR(50) | | 来源渠道 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE talents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    skills TEXT[],
    experience INTEGER,
    education VARCHAR(50),
    status VARCHAR(20) DEFAULT 'active',
    location VARCHAR(50),
    salary VARCHAR(50),
    source VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_talents_status ON talents(status);
CREATE INDEX idx_talents_location ON talents(location);
CREATE INDEX idx_talents_skills ON talents USING GIN(skills);
```

### 3.5 简历表 (resumes)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 简历ID |
| talent_id | INTEGER | REFERENCES talents(id) | 人才ID |
| file_path | VARCHAR(255) | | 文件路径 |
| file_name | VARCHAR(100) | | 文件名 |
| status | VARCHAR(20) | DEFAULT 'pending' | 状态 |
| ai_score | DECIMAL(5,2) | | AI评分 |
| ai_analysis | JSONB | | AI分析结果 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE resumes (
    id SERIAL PRIMARY KEY,
    talent_id INTEGER REFERENCES talents(id),
    file_path VARCHAR(255),
    file_name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending',
    ai_score DECIMAL(5,2),
    ai_analysis JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_resumes_status ON resumes(status);
CREATE INDEX idx_resumes_talent_id ON resumes(talent_id);
```

### 3.6 申请表 (applications)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 申请ID |
| resume_id | INTEGER | REFERENCES resumes(id) | 简历ID |
| job_id | INTEGER | REFERENCES jobs(id) | 职位ID |
| status | VARCHAR(20) | DEFAULT 'pending' | 状态 |
| match_score | DECIMAL(5,2) | | 匹配度 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    resume_id INTEGER REFERENCES resumes(id),
    job_id INTEGER REFERENCES jobs(id),
    status VARCHAR(20) DEFAULT 'pending',
    match_score DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_applications_status ON applications(status);
CREATE INDEX idx_applications_job_id ON applications(job_id);
```

### 3.7 面试表 (interviews)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 面试ID |
| candidate_id | INTEGER | | 候选人ID |
| candidate_name | VARCHAR(50) | | 候选人姓名 |
| position_id | INTEGER | | 职位ID |
| position | VARCHAR(100) | | 职位名称 |
| type | VARCHAR(20) | | 面试类型 |
| date | DATE | NOT NULL | 面试日期 |
| time | VARCHAR(10) | | 面试时间 |
| duration | INTEGER | DEFAULT 60 | 时长(分钟) |
| interviewer_id | INTEGER | | 面试官ID |
| interviewer | VARCHAR(50) | | 面试官姓名 |
| method | VARCHAR(20) | | 面试方式 |
| location | VARCHAR(255) | | 面试地点/链接 |
| status | VARCHAR(20) | DEFAULT 'scheduled' | 状态 |
| created_by | INTEGER | | 创建人ID |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT NOW() | 更新时间 |

```sql
CREATE TABLE interviews (
    id SERIAL PRIMARY KEY,
    candidate_id INTEGER,
    candidate_name VARCHAR(50),
    position_id INTEGER,
    position VARCHAR(100),
    type VARCHAR(20),
    date DATE NOT NULL,
    time VARCHAR(10),
    duration INTEGER DEFAULT 60,
    interviewer_id INTEGER,
    interviewer VARCHAR(50),
    method VARCHAR(20),
    location VARCHAR(255),
    status VARCHAR(20) DEFAULT 'scheduled',
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interviews_date ON interviews(date);
CREATE INDEX idx_interviews_status ON interviews(status);
CREATE INDEX idx_interviews_interviewer_id ON interviews(interviewer_id);
```

### 3.8 面试反馈表 (interview_feedbacks)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 反馈ID |
| interview_id | INTEGER | REFERENCES interviews(id) | 面试ID |
| interviewer_id | INTEGER | | 面试官ID |
| score | INTEGER | | 评分(1-5) |
| technical_score | INTEGER | | 技术评分 |
| communication_score | INTEGER | | 沟通评分 |
| culture_fit_score | INTEGER | | 文化匹配 |
| strengths | TEXT | | 优势 |
| weaknesses | TEXT | | 不足 |
| comment | TEXT | | 评语 |
| result | VARCHAR(20) | | 结果 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |

```sql
CREATE TABLE interview_feedbacks (
    id SERIAL PRIMARY KEY,
    interview_id INTEGER REFERENCES interviews(id),
    interviewer_id INTEGER,
    score INTEGER CHECK (score >= 1 AND score <= 5),
    technical_score INTEGER,
    communication_score INTEGER,
    culture_fit_score INTEGER,
    strengths TEXT,
    weaknesses TEXT,
    comment TEXT,
    result VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feedbacks_interview_id ON interview_feedbacks(interview_id);
```

### 3.9 消息表 (messages)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | SERIAL | PRIMARY KEY | 消息ID |
| sender_id | INTEGER | | 发送者ID |
| receiver_id | INTEGER | NOT NULL | 接收者ID |
| type | VARCHAR(20) | | 消息类型 |
| title | VARCHAR(100) | | 标题 |
| content | TEXT | | 内容 |
| is_read | BOOLEAN | DEFAULT FALSE | 是否已读 |
| created_at | TIMESTAMP | DEFAULT NOW() | 创建时间 |

```sql
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER,
    receiver_id INTEGER NOT NULL,
    type VARCHAR(20),
    title VARCHAR(100),
    content TEXT,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
CREATE INDEX idx_messages_is_read ON messages(is_read);
```

### 3.10 操作日志表 (operation_logs) - Elasticsearch

| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | keyword | 日志ID |
| timestamp | date | 时间戳 |
| service | keyword | 服务名称 |
| user_id | integer | 用户ID |
| username | keyword | 用户名 |
| ip | ip | IP地址 |
| method | keyword | 请求方法 |
| path | keyword | 请求路径 |
| query | text | 查询参数 |
| status_code | integer | 状态码 |
| duration | integer | 耗时(ms) |
| user_agent | text | 用户代理 |
| action | keyword | 操作类型 |
| module | keyword | 模块 |
| level | keyword | 日志级别 |

---

## 4. 索引设计

### 4.1 主键索引
所有表都使用 SERIAL 类型的自增主键。

### 4.2 唯一索引
- users.username
- users.email
- talents.email
- roles.code

### 4.3 普通索引
- 状态字段 (status)
- 外键字段 (xxx_id)
- 时间字段 (created_at)
- 常用查询字段 (location, department)

### 4.4 GIN 索引
- talents.skills (数组字段全文搜索)

---

## 5. 数据字典

### 5.1 用户状态 (user.status)
| 值 | 说明 |
|-----|------|
| active | 活跃 |
| inactive | 停用 |
| locked | 锁定 |

### 5.2 职位状态 (job.status)
| 值 | 说明 |
|-----|------|
| open | 招聘中 |
| closed | 已关闭 |
| filled | 已招满 |
| paused | 已暂停 |

### 5.3 人才状态 (talent.status)
| 值 | 说明 |
|-----|------|
| active | 活跃 |
| hired | 已雇佣 |
| pending | 待处理 |
| rejected | 已拒绝 |

### 5.4 简历状态 (resume.status)
| 值 | 说明 |
|-----|------|
| pending | 待处理 |
| reviewing | 审核中 |
| interviewed | 已面试 |
| hired | 已录用 |
| rejected | 已拒绝 |

### 5.5 面试状态 (interview.status)
| 值 | 说明 |
|-----|------|
| scheduled | 已安排 |
| completed | 已完成 |
| cancelled | 已取消 |
| no_show | 未出席 |

### 5.6 消息类型 (message.type)
| 值 | 说明 |
|-----|------|
| system | 系统通知 |
| interview | 面试邀约 |
| application | 简历投递 |
| notification | 提醒通知 |
| chat | 聊天消息 |

### 5.7 角色代码 (role.code)
| 值 | 说明 | 权限 |
|-----|------|------|
| admin | 超级管理员 | 所有权限 |
| hr_manager | HR主管 | 招聘全流程 |
| recruiter | 招聘专员 | 日常招聘 |
| interviewer | 面试官 | 面试评估 |
| viewer | 只读用户 | 仅查看 |

---

## 6. 数据初始化

### 6.1 初始化脚本
```bash
# 创建数据库
psql -U qinyang -c "CREATE DATABASE talent_platform;"

# 导入表结构
psql -U qinyang -d talent_platform -f backend/database/schema.sql

# 导入模拟数据
psql -U qinyang -d talent_platform -f backend/database/mock_data.sql
psql -U qinyang -d talent_platform -f backend/database/mock_data_2_talents.sql
psql -U qinyang -d talent_platform -f backend/database/mock_data_3_resumes.sql
psql -U qinyang -d talent_platform -f backend/database/mock_data_4_interviews.sql
psql -U qinyang -d talent_platform -f backend/database/mock_data_5_messages.sql
```

### 6.2 测试账号
| 用户名 | 密码 | 角色 |
|--------|------|------|
| admin | password123 | 超级管理员 |
| hr_zhang | password123 | HR主管 |
| hr_li | password123 | 招聘专员 |
| tech_chen | password123 | 面试官 |
| viewer_test | password123 | 只读用户 |

---

## 7. 备份与恢复

### 7.1 备份命令
```bash
# 完整备份
pg_dump -U qinyang -d talent_platform > backup_$(date +%Y%m%d).sql

# 仅数据备份
pg_dump -U qinyang -d talent_platform --data-only > data_backup.sql
```

### 7.2 恢复命令
```bash
# 恢复数据库
psql -U qinyang -d talent_platform < backup_20250101.sql
```

---

## 8. 性能优化建议

1. **定期 VACUUM**：清理死元组
2. **定期 ANALYZE**：更新统计信息
3. **监控慢查询**：开启 pg_stat_statements
4. **连接池**：使用 PgBouncer
5. **读写分离**：主从复制（生产环境）
