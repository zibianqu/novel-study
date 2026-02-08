-- 知识库表结构

-- 启用 pgvector 扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 知识库表
CREATE TABLE IF NOT EXISTS knowledge_base (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    title           VARCHAR(200) NOT NULL,
    content         TEXT NOT NULL,
    type            VARCHAR(50) NOT NULL,    -- 'character', 'worldview', 'plot', 'custom'
    tags            JSONB DEFAULT '[]',
    is_vectorized   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 知识向量表（pgvector）
CREATE TABLE IF NOT EXISTS knowledge_vectors (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    embedding       vector(1536),            -- OpenAI text-embedding-3-small 维度
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_knowledge_project_id ON knowledge_base(project_id);
CREATE INDEX IF NOT EXISTS idx_knowledge_type ON knowledge_base(type);
CREATE INDEX IF NOT EXISTS idx_vectors_project_id ON knowledge_vectors(project_id);

-- 向量相似度索引（HNSW）
CREATE INDEX IF NOT EXISTS idx_vectors_embedding ON knowledge_vectors 
USING hnsw (embedding vector_cosine_ops);

-- 成功消息
DO $$
BEGIN
    RAISE NOTICE '知识库表创建完成！';
END $$;
