-- Agent 工具调用日志表
CREATE TABLE IF NOT EXISTS agent_tool_calls (
    id SERIAL PRIMARY KEY,
    agent_id INT,
    tool_name VARCHAR(100) NOT NULL,
    input_params JSONB,
    output_result JSONB,
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    duration_ms INT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_agent_tool_calls_agent_id ON agent_tool_calls(agent_id);
CREATE INDEX idx_agent_tool_calls_tool_name ON agent_tool_calls(tool_name);
CREATE INDEX idx_agent_tool_calls_created_at ON agent_tool_calls(created_at);

-- 为 agents 表添加默认工具配置
COMMENT ON COLUMN agents.tools IS '工具列表, 格式: ["rag_search", "query_neo4j", ...]';

-- 更新 Agent 0 (总导演) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j", 
    "get_project_status",
    "get_storyline_status",
    "get_chapter_content",
    "update_storyline",
    "create_storyline"
]'::jsonb
WHERE agent_key = 'agent_0';

-- 更新 Agent 1 (旁白叙述者) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j",
    "get_chapter_content"
]'::jsonb
WHERE agent_key = 'agent_1';

-- 更新 Agent 2 (角色扮演者) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j",
    "get_chapter_content"
]'::jsonb
WHERE agent_key = 'agent_2';

-- 更新 Agent 3 (审核导演) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "get_chapter_content"
]'::jsonb
WHERE agent_key = 'agent_3';

-- 更新 Agent 4 (天线掌控者) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline"
]'::jsonb
WHERE agent_key = 'agent_4';

-- 更新 Agent 5 (地线掌控者) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline"
]'::jsonb
WHERE agent_key = 'agent_5';

-- 更新 Agent 6 (剧情线掌控者) 的工具配置
UPDATE agents 
SET tools = '[
    "rag_search",
    "query_neo4j",
    "get_storyline_status",
    "update_storyline",
    "create_storyline"
]'::jsonb
WHERE agent_key = 'agent_6';
