# 数据库设置指南

## 环境要求

- PostgreSQL 14+
- 数据库名称: `talent_platform`

## 快速开始

### 1. 创建数据库

```bash
# 连接到 PostgreSQL
psql -U postgres

# 创建数据库
CREATE DATABASE talent_platform;

# 退出
\q
```

### 2. 执行表结构脚本

```bash
psql -U postgres -d talent_platform -f schema.sql
```

### 3. 导入 Mock 数据

```bash
psql -U postgres -d talent_platform -f mock_data.sql
```

### 一键执行

```bash
# 创建数据库并导入所有数据
psql -U postgres -c "DROP DATABASE IF EXISTS talent_platform;"
psql -U postgres -c "CREATE DATABASE talent_platform;"
psql -U postgres -d talent_platform -f schema.sql
psql -U postgres -d talent_platform -f mock_data.sql
```

## 数据库配置

在各服务的环境变量中配置:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=talent_platform
```

## 表结构说明

### 核心表

| 表名 | 说明 | 记录数 |
|------|------|--------|
| users | 用户表 | 9 |
| jobs | 职位表 | 12 |
| talents | 人才表 | 20 |
| resumes | 简历表 | 17 |
| interviews | 面试表 | 20 |
| messages | 消息表 | 19 |

### 辅助表

| 表名 | 说明 |
|------|------|
| interview_feedbacks | 面试反馈 |
| applications | 应聘记录 |
| roles | 角色权限 |
| operation_logs | 操作日志 |

## Mock 数据说明

### 用户账号

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | password123 | 超级管理员 | 系统管理员 |
| hr_zhang | password123 | HR主管 | HR总监 |
| hr_li | password123 | 招聘专员 | 招聘专员 |
| hr_wang | password123 | 招聘专员 | 招聘专员 |
| tech_chen | password123 | 面试官 | 技术总监 |
| tech_liu | password123 | 面试官 | 前端负责人 |
| tech_zhao | password123 | 面试官 | 后端负责人 |
| product_sun | password123 | 面试官 | 产品总监 |
| viewer_test | password123 | 只读用户 | 运营专员 |

### 职位状态分布

- 招聘中 (open): 10 个
- 已关闭 (closed): 1 个
- 已满员 (filled): 1 个

### 人才状态分布

- 活跃 (active): 16 人
- 已雇佣 (hired): 1 人
- 待处理 (pending): 1 人
- 已拒绝 (rejected): 1 人

### 面试数据

- 今天的面试: 3 场
- 明天的面试: 2 场
- 后天的面试: 2 场
- 已完成的面试: 12 场
- 已取消的面试: 1 场

### 消息数据

- 系统通知
- 面试邀请
- 面试反馈
- Offer 审批
- 未读消息: 5 条

## 数据关系

```
users
  ├── jobs (created_by)
  ├── talents (user_id)
  ├── interviews (interviewer_id, created_by)
  └── messages (sender_id, receiver_id)

talents
  ├── resumes (talent_id)
  └── applications (talent_id)

jobs
  ├── resumes (job_id)
  └── applications (job_id)

interviews
  └── interview_feedbacks (interview_id)
```

## 常用查询

### 查看统计数据

```sql
SELECT '用户数' as metric, COUNT(*) as count FROM users
UNION ALL SELECT '职位数', COUNT(*) FROM jobs
UNION ALL SELECT '人才数', COUNT(*) FROM talents
UNION ALL SELECT '简历数', COUNT(*) FROM resumes
UNION ALL SELECT '面试数', COUNT(*) FROM interviews
UNION ALL SELECT '消息数', COUNT(*) FROM messages;
```

### 查看今日面试

```sql
SELECT candidate_name, position, time, interviewer, location
FROM interviews
WHERE date = CURRENT_DATE::text AND status = 'scheduled'
ORDER BY time;
```

### 查看未读消息

```sql
SELECT m.*, u.username as sender_name
FROM messages m
LEFT JOIN users u ON m.sender_id = u.id
WHERE m.receiver_id = 3 AND m.is_read = false
ORDER BY m.created_at DESC;
```

### 查看招聘漏斗

```sql
SELECT
    status,
    COUNT(*) as count
FROM resumes
GROUP BY status
ORDER BY
    CASE status
        WHEN 'pending' THEN 1
        WHEN 'reviewing' THEN 2
        WHEN 'interviewed' THEN 3
        WHEN 'offered' THEN 4
        WHEN 'hired' THEN 5
        WHEN 'rejected' THEN 6
    END;
```

## 重置数据

如果需要重置数据，执行:

```bash
psql -U postgres -d talent_platform -f mock_data.sql
```

这会清空所有表并重新插入 Mock 数据。
