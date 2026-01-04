# 数据库种子数据

本目录包含用于追加数据库数据的 SQL 文件。所有 SQL 都使用 `ON CONFLICT DO NOTHING` 或条件判断，不会覆盖现有数据。

## 文件说明

| 文件名 | 说明 | 新增数据量 |
|--------|------|-----------|
| `seed_data_users.sql` | 追加用户（HR、面试官） | ~6条 |
| `seed_data_jobs.sql` | 追加职位（各部门各级别） | ~13条 |
| `seed_data_talents.sql` | 追加人才（各种技能） | ~30条 |
| `seed_data_resumes.sql` | 追加简历 | ~10条 |
| `seed_data_applications.sql` | 追加申请记录 | ~15条 |
| `seed_data_interviews.sql` | 追加面试记录 | ~15条 |
| `seed_data_messages.sql` | 追加消息通知 | ~20条 |

## 导入方法

### 方法一：使用主导入脚本（推荐）

```bash
cd backend/database
psql -U postgres -d talent_platform -f import_all_data.sql
```

### 方法二：单独导入各文件

按以下顺序导入（有依赖关系）：

```bash
cd backend/database
psql -U postgres -d talent_platform -f seed_data_users.sql
psql -U postgres -d talent_platform -f seed_data_jobs.sql
psql -U postgres -d talent_platform -f seed_data_talents.sql
psql -U postgres -d talent_platform -f seed_data_resumes.sql
psql -U postgres -d talent_platform -f seed_data_applications.sql
psql -U postgres -d talent_platform -f seed_data_interviews.sql
psql -U postgres -d talent_platform -f seed_data_messages.sql
```

### 方法三：在 psql 交互模式中导入

```sql
\c talent_platform
\i import_all_data.sql
```

## 新增测试账号

密码都是 `123456`：

| 用户名 | 角色 | 说明 |
|--------|------|------|
| `hr_wang` | HR | 人力资源 |
| `hr_chen` | HR | 人力资源 |
| `interviewer_zhao` | 面试官 | 技术面试官 |
| `interviewer_sun` | 面试官 | 技术面试官 |

## 数据特点

1. **追加模式**：所有 SQL 使用 `ON CONFLICT DO NOTHING`，不会覆盖现有数据
2. **动态关联**：申请、面试、简历等数据使用子查询获取ID，自动关联已有数据
3. **真实场景**：包含各种状态的数据（待处理、面试中、已发offer、已入职等）

## 注意事项

- 可以多次执行，不会产生重复数据
- 如果某些关联数据不存在，相关记录会被跳过
- 建议先备份数据库再导入
