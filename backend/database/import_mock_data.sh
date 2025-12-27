#!/bin/bash
# 导入模拟数据脚本
# 使用方法: ./import_mock_data.sh

DB_USER="qinyang"
DB_NAME="talent_platform"

echo "开始导入模拟数据..."

# 按顺序执行SQL文件
psql -U $DB_USER -d $DB_NAME -f mock_data.sql
psql -U $DB_USER -d $DB_NAME -f mock_data_2_talents.sql
psql -U $DB_USER -d $DB_NAME -f mock_data_3_resumes.sql
psql -U $DB_USER -d $DB_NAME -f mock_data_4_interviews.sql
psql -U $DB_USER -d $DB_NAME -f mock_data_5_messages.sql

echo "模拟数据导入完成！"
