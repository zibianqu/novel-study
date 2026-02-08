-- NovelForge AI - API密钥加密存储
-- 创建时间: 2026-02-08
-- 目的: 使用AES-256-GCM加密存储API密钥

-- ====================================
-- 修改API密钥表结构
-- ====================================

-- 添加加密字段
ALTER TABLE api_keys 
    ADD COLUMN IF NOT EXISTS key_encrypted TEXT,
    ADD COLUMN IF NOT EXISTS encryption_version INTEGER DEFAULT 1;

-- 添加轮换记录字段
ALTER TABLE api_keys 
    ADD COLUMN IF NOT EXISTS last_rotated_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS rotation_count INTEGER DEFAULT 0;

-- ====================================
-- 创建API密钥审计日志表
-- ====================================

CREATE TABLE IF NOT EXISTS api_key_audit_logs (
    id SERIAL PRIMARY KEY,
    api_key_id INTEGER NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'created', 'updated', 'rotated', 'deleted', 'accessed'
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_api_key_id ON api_key_audit_logs(api_key_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON api_key_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON api_key_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON api_key_audit_logs(created_at DESC);

-- ====================================
-- 创建API密钥使用统计表
-- ====================================

CREATE TABLE IF NOT EXISTS api_key_usage_stats (
    id SERIAL PRIMARY KEY,
    api_key_id INTEGER NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    request_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(api_key_id, date)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_usage_stats_api_key_id ON api_key_usage_stats(api_key_id);
CREATE INDEX IF NOT EXISTS idx_usage_stats_date ON api_key_usage_stats(date DESC);

-- ====================================
-- 更新函数
-- ====================================

-- 更新 updated_at 字段的触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为api_key_usage_stats添加触发器
DROP TRIGGER IF EXISTS update_api_key_usage_stats_updated_at ON api_key_usage_stats;
CREATE TRIGGER update_api_key_usage_stats_updated_at
    BEFORE UPDATE ON api_key_usage_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
