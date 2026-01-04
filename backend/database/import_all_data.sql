-- 导入所有种子数据（追加模式）
-- 使用方法: psql -U postgres -d talent_platform -f import_all_data.sql
-- 或者在 psql 中执行: \i import_all_data.sql
-- 注意：这些SQL会追加数据，不会覆盖现有数据

-- 设置客户端编码
SET client_encoding = 'UTF8';

-- 开始事务
BEGIN;

-- 导入顺序很重要，需要先导入被依赖的表

-- 1. 导入用户数据
\echo '正在追加用户数据...'
\i seed_data_users.sql

-- 2. 导入职位数据
\echo '正在追加职位数据...'
\i seed_data_jobs.sql

-- 3. 导入人才数据
\echo '正在追加人才数据...'
\i seed_data_talents.sql

-- 4. 导入简历数据
\echo '正在追加简历数据...'
\i seed_data_resumes.sql

-- 5. 导入申请数据
\echo '正在追加申请数据...'
\i seed_data_applications.sql

-- 6. 导入面试数据
\echo '正在追加面试数据...'
\i seed_data_interviews.sql

-- 7. 导入消息数据
\echo '正在追加消息数据...'
\i seed_data_messages.sql

-- 提交事务
COMMIT;

\echo '所有数据追加完成！'

-- 显示数据统计
\echo ''
\echo '=== 数据统计 ==='
SELECT 'users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'jobs', COUNT(*) FROM jobs
UNION ALL
SELECT 'talents', COUNT(*) FROM talents
UNION ALL
SELECT 'resumes', COUNT(*) FROM resumes
UNION ALL
SELECT 'applications', COUNT(*) FROM applications
UNION ALL
SELECT 'interviews', COUNT(*) FROM interviews
UNION ALL
SELECT 'messages', COUNT(*) FROM messages
ORDER BY table_name;
