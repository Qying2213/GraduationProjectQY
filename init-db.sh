#!/bin/bash
# 数据库初始化脚本

echo "=========================================="
echo "  智能人才运营平台 - 数据库初始化"
echo "=========================================="

DB_USER="${DB_USER:-qinyang}"
DB_NAME="${DB_NAME:-talent_platform}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

echo ""
echo "数据库配置："
echo "  主机: $DB_HOST:$DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"
echo ""

# 检查 PostgreSQL 是否运行
if ! pg_isready -h $DB_HOST -p $DB_PORT > /dev/null 2>&1; then
    echo "❌ PostgreSQL 未运行，请先启动 PostgreSQL"
    echo "   macOS: brew services start postgresql"
    exit 1
fi

echo "✅ PostgreSQL 已运行"

# 检查数据库是否存在
if psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME; then
    echo "✅ 数据库 $DB_NAME 已存在"
    
    read -p "是否重置数据库？(y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "正在重置数据库..."
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/database/schema.sql
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/database/mock_data.sql
        echo "✅ 数据库已重置"
    fi
else
    echo "正在创建数据库 $DB_NAME..."
    createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    
    echo "正在初始化表结构..."
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/database/schema.sql
    
    echo "正在导入测试数据..."
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f backend/database/mock_data.sql
    
    echo "✅ 数据库初始化完成"
fi

echo ""
echo "=========================================="
echo "  初始化完成！"
echo "=========================================="
echo ""
echo "测试账号："
echo "  用户名: admin / hr_li / hr_wang"
echo "  密码: password123"
echo ""
