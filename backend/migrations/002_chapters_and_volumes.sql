-- 章节和卷的表结构

-- 卷表
CREATE TABLE IF NOT EXISTS volumes (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    title           VARCHAR(200) NOT NULL,
    summary         TEXT,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 章节表
CREATE TABLE IF NOT EXISTS chapters (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    volume_id       INT REFERENCES volumes(id) ON DELETE SET NULL,
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    word_count      INT DEFAULT 0,
    sort_order      INT DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'draft',
    locked_by       INT REFERENCES users(id) ON DELETE SET NULL,
    locked_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_volumes_project_id ON volumes(project_id);
CREATE INDEX IF NOT EXISTS idx_chapters_project_id ON chapters(project_id);
CREATE INDEX IF NOT EXISTS idx_chapters_volume_id ON chapters(volume_id);
CREATE INDEX IF NOT EXISTS idx_chapters_status ON chapters(status);

-- 成功消息
DO $$
BEGIN
    RAISE NOTICE '章节和卷表创建完成！';
END $$;
