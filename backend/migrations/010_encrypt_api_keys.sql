-- API密钥加密存储
-- 执行时间: 2026-02-08

ALTER TABLE users ADD COLUMN IF NOT EXISTS encryption_version INT DEFAULT 1;
ALTER TABLE users ADD COLUMN IF NOT EXISTS key_updated_at TIMESTAMP;

CREATE TABLE IF NOT EXISTS api_key_audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(20) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_api_key_audit_user ON api_key_audit_logs(user_id, created_at DESC);

CREATE TABLE IF NOT EXISTS api_key_rotations (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    old_key_hash VARCHAR(64),
    rotation_reason VARCHAR(100),
    rotated_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION audit_api_key_change()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND NEW.api_key_encrypted IS DISTINCT FROM OLD.api_key_encrypted THEN
        INSERT INTO api_key_audit_logs (user_id, action, created_at)
        VALUES (NEW.id, 'updated', NOW());
        
        INSERT INTO api_key_rotations (user_id, rotation_reason)
        VALUES (NEW.id, 'manual_update');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_audit_api_key_change
AFTER UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION audit_api_key_change();

COMMENT ON TABLE api_key_audit_logs IS 'API密钥访问审计日志';
COMMENT ON TABLE api_key_rotations IS 'API密钥轮换记录';
