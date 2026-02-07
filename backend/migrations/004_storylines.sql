-- 三线系统表结构

-- 三线表
CREATE TABLE IF NOT EXISTS storylines (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    line_type       VARCHAR(20) NOT NULL,    -- 'skyline', 'groundline', 'plotline'
    title           VARCHAR(200) NOT NULL,
    content         TEXT,
    chapter_range   VARCHAR(50),             -- JSON格式: '[1,10]'
    status          VARCHAR(20) DEFAULT 'planned',
    sort_order      INT DEFAULT 0,
    parent_id       INT REFERENCES storylines(id) ON DELETE SET NULL,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 三线交汇点
CREATE TABLE IF NOT EXISTS storyline_convergences (
    id                  SERIAL PRIMARY KEY,
    project_id          INT REFERENCES projects(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    skyline_meaning     TEXT,
    groundline_meaning  TEXT,
    plotline_meaning    TEXT,
    chapter_id          INT REFERENCES chapters(id) ON DELETE SET NULL,
    created_at          TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_storylines_project_id ON storylines(project_id);
CREATE INDEX IF NOT EXISTS idx_storylines_line_type ON storylines(line_type);
CREATE INDEX IF NOT EXISTS idx_convergences_project_id ON storyline_convergences(project_id);

-- 成功消息
DO $$
BEGIN
    RAISE NOTICE '三线表创建完成！';
END $$;
