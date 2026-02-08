-- 添加性能优化索引
-- 执行时间: 2026-02-08

-- 用户表索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- 项目表索引
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_projects_user_status ON projects(user_id, status);

-- 章节表索引
CREATE INDEX IF NOT EXISTS idx_chapters_project_id ON chapters(project_id);
CREATE INDEX IF NOT EXISTS idx_chapters_volume_id ON chapters(volume_id);
CREATE INDEX IF NOT EXISTS idx_chapters_status ON chapters(status);
CREATE INDEX IF NOT EXISTS idx_chapters_updated_at ON chapters(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_chapters_project_sort ON chapters(project_id, sort_order);

-- 知识库索引
CREATE INDEX IF NOT EXISTS idx_knowledge_project_id ON knowledge_base(project_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_category ON knowledge_base(category);
CREATE INDEX IF NOT EXISTS idx_knowledge_project_category ON knowledge_base(project_id, category);

-- Agent相关索引
CREATE INDEX IF NOT EXISTS idx_agents_user_id ON agents(user_id) WHERE user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_agents_type ON agents(type);
CREATE INDEX IF NOT EXISTS idx_agents_active ON agents(is_active);

-- Agent知识库索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_agent_id ON agent_knowledge_items(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_category ON agent_knowledge_items(category_id);
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_active ON agent_knowledge_items(is_active);

-- 三线管理索引
CREATE INDEX IF NOT EXISTS idx_storylines_project_id ON storylines(project_id);
CREATE INDEX IF NOT EXISTS idx_storylines_type ON storylines(line_type);
CREATE INDEX IF NOT EXISTS idx_storylines_project_type ON storylines(project_id, line_type);

-- 向量搜索优化
DROP INDEX IF EXISTS content_embeddings_embedding_idx;
CREATE INDEX content_embeddings_embedding_idx ON content_embeddings 
USING hnsw (embedding vector_cosine_ops) 
WITH (m = 16, ef_construction = 64);

DROP INDEX IF EXISTS agent_knowledge_embeddings_embedding_idx;
CREATE INDEX agent_knowledge_embeddings_embedding_idx ON agent_knowledge_embeddings 
USING hnsw (embedding vector_cosine_ops) 
WITH (m = 16, ef_construction = 64);

-- 复合索引优化查询
CREATE INDEX IF NOT EXISTS idx_content_embeddings_project ON content_embeddings(project_id, chapter_id);
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_embeddings_agent ON agent_knowledge_embeddings(agent_id, item_id);

COMMENT ON INDEX idx_projects_user_status IS '用户项目状态复合查询优化';
COMMENT ON INDEX content_embeddings_embedding_idx IS 'HNSW向量相似度搜索优化';
