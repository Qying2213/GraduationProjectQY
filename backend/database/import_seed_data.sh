#!/bin/bash
# 导入种子数据脚本

DB_NAME="talent_platform"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "开始导入种子数据到 $DB_NAME..."

psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_users.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_jobs.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_talents.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_resumes.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_applications.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_interviews.sql"
psql -d $DB_NAME -f "$SCRIPT_DIR/seed_data_messages.sql"

echo ""
echo "=== 数据统计 ==="
psql -d $DB_NAME -c "SELECT 'users' as table_name, COUNT(*) as count FROM users UNION ALL SELECT 'jobs', COUNT(*) FROM jobs UNION ALL SELECT 'talents', COUNT(*) FROM talents UNION ALL SELECT 'resumes', COUNT(*) FROM resumes UNION ALL SELECT 'applications', COUNT(*) FROM applications UNION ALL SELECT 'interviews', COUNT(*) FROM interviews UNION ALL SELECT 'messages', COUNT(*) FROM messages ORDER BY table_name;"

echo ""
echo "导入完成！"
