-- Agent和工作流表结构

-- Agent表
CREATE TABLE IF NOT EXISTS agents (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id) ON DELETE CASCADE,
    agent_key       VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(100) NOT NULL,
    icon            VARCHAR(50),
    description     TEXT,
    type            VARCHAR(20) NOT NULL,       -- 'core' 或 'extension'
    layer           VARCHAR(20) NOT NULL,       -- 'decision', 'strategy', 'execution', 'quality', 'auxiliary'
    system_prompt   TEXT NOT NULL,
    model           VARCHAR(50) DEFAULT 'gpt-4o',
    temperature     FLOAT DEFAULT 0.7,
    max_tokens      INT DEFAULT 4096,
    tools           JSONB DEFAULT '[]',
    input_schema    JSONB DEFAULT '{}',
    output_schema   JSONB DEFAULT '{}',
    permissions     JSONB DEFAULT '{}',
    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- AI交互日志
CREATE TABLE IF NOT EXISTS ai_interaction_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id) ON DELETE CASCADE,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    agent_id        INT REFERENCES agents(id) ON DELETE SET NULL,
    action_type     VARCHAR(50),
    input_prompt    TEXT,
    output_response TEXT,
    tokens_input    INT DEFAULT 0,
    tokens_output   INT DEFAULT 0,
    model           VARCHAR(50),
    duration_ms     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_agents_type ON agents(type);
CREATE INDEX IF NOT EXISTS idx_agents_is_active ON agents(is_active);
CREATE INDEX IF NOT EXISTS idx_ai_logs_user_id ON ai_interaction_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_logs_project_id ON ai_interaction_logs(project_id);

-- 插入核心Agent
INSERT INTO agents (agent_key, name, icon, description, type, layer, system_prompt, model, temperature, max_tokens, sort_order) VALUES
('agent_0_director', '总导演', '🎬', '全局调度、任务分配、质量把控', 'core', 'decision', '你是 NovelForge AI 的总导演...', 'gpt-4o', 0.5, 4096, 0),
('agent_1_narrator', '旁白叙述者', '🎙️', '环境/动作/心理描写、叙事', 'core', 'execution', '你是 NovelForge AI 的旁白叙述者...', 'gpt-4o', 0.8, 4096, 1),
('agent_2_character', '角色扉演者', '🎭', '角色对话、角色行为、多角色互动', 'core', 'execution', '你是 NovelForge AI 的角色扉演者...', 'gpt-4o', 0.9, 4096, 2),
('agent_3_quality', '审核导演', '👁️', '质量审核、一致性检查、修改指导', 'core', 'quality', '你是 NovelForge AI 的审核导漗...', 'gpt-4o', 0.3, 2048, 3)
ON CONFLICT (agent_key) DO NOTHING;

-- 成功消息
DO $$
BEGIN
    RAISE NOTICE 'Agent表创建完成！';
END $$;
