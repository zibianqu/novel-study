-- NovelForge AI 数据库初始化脚本

-- 启用 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    avatar          VARCHAR(500),
    settings        JSONB DEFAULT '{}',
    api_key_encrypted VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(200) NOT NULL,
    type            VARCHAR(20) NOT NULL,
    genre           VARCHAR(50),
    description     TEXT,
    cover_image     VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'draft',
    word_count      INT DEFAULT 0,
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);

-- 成功消息
DO $$
BEGIN
    RAISE NOTICE '数据库初始化完成！';
END $$;
