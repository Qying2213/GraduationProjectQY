# 智能人才招聘管理平台 - 数据库设计文档

> 📖 返回 [项目首页](../README.md) | 相关文档：[系统架构](ARCHITECTURE.md) | [系统设计](SYSTEM_DESIGN.md) | [快速启动](QUICKSTART.md)

---

## 1. 数据库概述

- **数据库类型**：PostgreSQL 14+
- **字符集**：UTF-8
- **时区**：Asia/Shanghai
- **核心表数量**：10张

---

## 2. ER图

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

---

## 3. 数据表详细设计

### 3.1 用户表 (users)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| username | VARCHAR(50) | NOT NULL UNIQUE | 用户名 |
| email | VARCHAR(100) | NOT NULL UNIQUE | 邮箱 |
| password | VARCHAR(255) | NOT NULL | 密码（bcrypt加密） |
| role | VARCHAR(20) | NOT NULL DEFAULT 'viewer' | 角色 |
| avatar | VARCHAR(500) | | 头像URL |
| phone | VARCHAR(20) | | 手机号 |
| department | VARCHAR(50) | | 部门 |
| position | VARCHAR(50) | | 职位 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_users_email (email)
- idx_users_role (role)
- idx_users_status (status)

**角色枚举**：admin, hr_manager, recruiter, interviewer, viewer

**状态枚举**：active, inactive, suspended


### 3.2 角色表 (roles)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | VARCHAR(50) | NOT NULL UNIQUE | 角色名称 |
| code | VARCHAR(50) | NOT NULL UNIQUE | 角色代码 |
| description | TEXT | | 描述 |
| permissions | TEXT[] | | 权限列表 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**预设角色**：
| code | name | permissions |
|------|------|-------------|
| admin | 超级管理员 | ['*'] |
| hr_manager | HR主管 | ['talent:*', 'job:*', 'resume:*', 'interview:*', 'message:*'] |
| recruiter | 招聘专员 | ['talent:view', 'talent:create', 'talent:edit', 'job:view', 'resume:*', 'interview:*'] |
| interviewer | 面试官 | ['talent:view', 'job:view', 'interview:view', 'interview:feedback'] |
| viewer | 只读用户 | ['talent:view', 'job:view', 'resume:view', 'interview:view'] |

### 3.3 职位表 (jobs)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| title | VARCHAR(200) | NOT NULL | 职位名称 |
| description | TEXT | | 职位描述 |
| requirements | TEXT[] | | 职位要求 |
| salary | VARCHAR(50) | | 薪资范围 |
| location | VARCHAR(50) | | 工作地点 |
| type | VARCHAR(20) | NOT NULL DEFAULT 'full-time' | 职位类型 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'open' | 状态 |
| created_by | INTEGER | REFERENCES users(id) | 创建人 |
| department | VARCHAR(50) | | 所属部门 |
| level | VARCHAR(20) | | 职级 |
| skills | TEXT[] | | 技能要求 |
| benefits | TEXT[] | | 福利待遇 |
| headcount | INTEGER | DEFAULT 1 | 招聘人数 |
| urgent | BOOLEAN | DEFAULT FALSE | 是否紧急 |
| deadline | DATE | | 截止日期 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_jobs_status (status)
- idx_jobs_type (type)
- idx_jobs_location (location)
- idx_jobs_created_by (created_by)

**类型枚举**：full-time, part-time, contract, internship

**状态枚举**：open, closed, filled, paused

**职级枚举**：junior, mid, senior, expert, management

### 3.4 人才表 (talents)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| name | VARCHAR(100) | NOT NULL | 姓名 |
| email | VARCHAR(100) | NOT NULL | 邮箱 |
| phone | VARCHAR(20) | | 手机号 |
| skills | TEXT[] | | 技能列表 |
| experience | INTEGER | DEFAULT 0 | 工作年限 |
| education | VARCHAR(20) | | 学历 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态 |
| tags | TEXT[] | | 标签 |
| user_id | INTEGER | REFERENCES users(id) | 关联用户 |
| location | VARCHAR(50) | | 所在地 |
| salary | VARCHAR(50) | | 期望薪资 |
| summary | TEXT | | 个人简介 |
| gender | VARCHAR(10) | | 性别 |
| age | INTEGER | | 年龄 |
| current_company | VARCHAR(100) | | 当前公司 |
| current_position | VARCHAR(100) | | 当前职位 |
| source | VARCHAR(50) | | 来源渠道 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_talents_email (email)
- idx_talents_status (status)
- idx_talents_skills (skills) - GIN索引
- idx_talents_location (location)

**状态枚举**：active, hired, pending, rejected

**学历枚举**：高中, 大专, 本科, 硕士, 博士


### 3.5 简历表 (resumes)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| talent_id | INTEGER | REFERENCES talents(id) ON DELETE CASCADE | 人才ID |
| job_id | INTEGER | REFERENCES jobs(id) ON DELETE SET NULL | 职位ID |
| file_path | VARCHAR(500) | | 文件路径 |
| file_name | VARCHAR(200) | | 文件名 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'pending' | 状态 |
| match_score | INTEGER | DEFAULT 0 | AI匹配分数(0-100) |
| parse_result | JSONB | | 解析结果 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_resumes_talent_id (talent_id)
- idx_resumes_job_id (job_id)
- idx_resumes_status (status)

**状态枚举**：pending, reviewing, interviewed, offered, hired, rejected

### 3.6 面试表 (interviews)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| candidate_id | INTEGER | NOT NULL | 候选人ID |
| candidate_name | VARCHAR(100) | NOT NULL | 候选人姓名 |
| position_id | INTEGER | NOT NULL | 职位ID |
| position | VARCHAR(200) | NOT NULL | 职位名称 |
| type | VARCHAR(20) | NOT NULL DEFAULT 'initial' | 面试类型 |
| date | VARCHAR(20) | NOT NULL | 面试日期 |
| time | VARCHAR(10) | NOT NULL | 面试时间 |
| duration | INTEGER | DEFAULT 60 | 时长(分钟) |
| interviewer_id | INTEGER | REFERENCES users(id) | 面试官ID |
| interviewer | VARCHAR(100) | NOT NULL | 面试官姓名 |
| method | VARCHAR(20) | NOT NULL DEFAULT 'onsite' | 面试方式 |
| location | VARCHAR(500) | | 面试地点/链接 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'scheduled' | 状态 |
| notes | TEXT | | 备注 |
| feedback | TEXT | | 反馈 |
| rating | INTEGER | DEFAULT 0 | 评分(1-5) |
| created_by | INTEGER | REFERENCES users(id) | 创建人 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_interviews_candidate_id (candidate_id)
- idx_interviews_interviewer_id (interviewer_id)
- idx_interviews_date (date)
- idx_interviews_status (status)

**类型枚举**：initial(初试), second(复试), final(终面), hr(HR面)

**方式枚举**：onsite(现场), video(视频), phone(电话)

**状态枚举**：scheduled, completed, cancelled, no_show

### 3.7 面试反馈表 (interview_feedbacks)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| interview_id | INTEGER | REFERENCES interviews(id) ON DELETE CASCADE | 面试ID |
| interviewer_id | INTEGER | REFERENCES users(id) | 面试官ID |
| rating | INTEGER | NOT NULL CHECK (1-5) | 评分 |
| strengths | TEXT | | 优势 |
| weaknesses | TEXT | | 不足 |
| comments | TEXT | | 评语 |
| recommendation | VARCHAR(50) | | 建议 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**索引**：
- idx_interview_feedbacks_interview_id (interview_id)

**建议枚举**：pass, fail, pending


### 3.8 消息表 (messages)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| sender_id | INTEGER | REFERENCES users(id) | 发送者ID |
| receiver_id | INTEGER | REFERENCES users(id) NOT NULL | 接收者ID |
| type | VARCHAR(20) | NOT NULL DEFAULT 'system' | 消息类型 |
| title | VARCHAR(200) | NOT NULL | 标题 |
| content | TEXT | | 内容 |
| is_read | BOOLEAN | DEFAULT FALSE | 是否已读 |
| read_at | TIMESTAMP | | 阅读时间 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**：
- idx_messages_receiver_id (receiver_id)
- idx_messages_sender_id (sender_id)
- idx_messages_is_read (is_read)
- idx_messages_type (type)

**类型枚举**：system, interview, feedback, offer, reminder, chat

### 3.9 应聘记录表 (applications)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| talent_id | INTEGER | REFERENCES talents(id) ON DELETE CASCADE | 人才ID |
| job_id | INTEGER | REFERENCES jobs(id) ON DELETE CASCADE | 职位ID |
| resume_id | INTEGER | REFERENCES resumes(id) | 简历ID |
| stage | VARCHAR(50) | NOT NULL DEFAULT 'applied' | 阶段 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态 |
| source | VARCHAR(50) | | 来源 |
| notes | TEXT | | 备注 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新时间 |

**约束**：UNIQUE(talent_id, job_id)

**索引**：
- idx_applications_talent_id (talent_id)
- idx_applications_job_id (job_id)
- idx_applications_stage (stage)

**阶段枚举**：applied, screening, interview, offer, hired, rejected

### 3.10 操作日志表 (operation_logs)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | SERIAL | PRIMARY KEY | 主键 |
| user_id | INTEGER | REFERENCES users(id) | 用户ID |
| action | VARCHAR(50) | NOT NULL | 操作类型 |
| resource_type | VARCHAR(50) | | 资源类型 |
| resource_id | INTEGER | | 资源ID |
| details | JSONB | | 详细信息 |
| ip_address | VARCHAR(50) | | IP地址 |
| user_agent | TEXT | | 用户代理 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |

**索引**：
- idx_operation_logs_user_id (user_id)
- idx_operation_logs_action (action)
- idx_operation_logs_created_at (created_at)

---

## 4. 触发器

### 4.1 自动更新 updated_at

所有包含 `updated_at` 字段的表都配置了触发器，在更新记录时自动更新该字段：

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
```

---

## 5. 数据库初始化

### 5.1 创建数据库

```bash
psql -U postgres -c "CREATE DATABASE talent_platform;"
```

### 5.2 导入表结构

```bash
psql -U postgres -d talent_platform -f backend/database/schema.sql
```

### 5.3 导入模拟数据

```bash
psql -U postgres -d talent_platform -f backend/database/mock_data.sql
```

---

## 📚 相关文档

| 文档 | 说明 |
|------|------|
| [📖 项目首页](../README.md) | 项目概述和快速入门 |
| [📐 系统架构](ARCHITECTURE.md) | 微服务架构设计 |
| [📋 系统设计](SYSTEM_DESIGN.md) | 功能模块设计 |
| [🚀 快速启动](QUICKSTART.md) | 环境配置、安装步骤 |
| [🐳 部署文档](DEPLOYMENT.md) | Docker部署、生产环境 |
| [🧪 测试指南](TEST_GUIDE.md) | API测试、功能测试 |
| [📝 代码规范](CODE_GUIDE.md) | 目录结构、开发指南 |
