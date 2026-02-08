#!/bin/bash

# NovelForge AI - 数据库初始化脚本

set -e

echo "================================="
echo "NovelForge AI - Database Setup"
echo "================================="

# 配置
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-novelforge}"
DB_PASSWORD="${DB_PASSWORD:-novelforge_password}"
DB_NAME="${DB_NAME:-novelforge_db}"

export PGPASSWORD="$DB_PASSWORD"

# 检查 PostgreSQL 是否可用
echo "Checking PostgreSQL connection..."
until psql -h "$DB_HOST" -U "$DB_USER" -d postgres -c '\q' 2>/dev/null; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done
echo "✅ PostgreSQL is ready!"

# 创建数据库（如果不存在）
echo "Creating database if not exists..."
psql -h "$DB_HOST" -U "$DB_USER" -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    psql -h "$DB_HOST" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME"

echo "✅ Database '$DB_NAME' is ready!"

# 执行迁移
echo "Running migrations..."

MIGRATIONS_DIR="backend/migrations"

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "❌ Error: Migrations directory not found: $MIGRATIONS_DIR"
    exit 1
fi

for migration in "$MIGRATIONS_DIR"/*.sql; do
    if [ -f "$migration" ]; then
        echo "Applying migration: $(basename "$migration")..."
        psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$migration"
        echo "✅ Applied: $(basename "$migration")"
    fi
done

echo "================================="
echo "✅ Database setup completed!"
echo "================================="

# 显示表列表
echo "
Tables in database:"
psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "\dt"

unset PGPASSWORD
