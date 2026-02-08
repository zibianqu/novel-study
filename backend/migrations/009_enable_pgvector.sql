-- NovelForge AI - 启用 pgvector 扫展
-- 创建时间: 2026-02-08
-- 目的: 优化向量搜索性能

-- ====================================
-- 启用 pgvector 扫展
-- ====================================

CREATE EXTENSION IF NOT EXISTS vector;

-- ====================================
-- 创建向量 Embedding 表
-- ====================================

-- 内容 Embedding 表
CREATE TABLE IF NOT EXISTS content_embeddings (
    id SERIAL PRIMARY KEY,
    content_id INTEGER NOT NULL,
    content_type VARCHAR(50) NOT NULL, -- 'chapter', 'knowledge', etc
    embedding vector(1536) NOT NULL,   -- OpenAI ada-002
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Agent 知识 Embedding 表
CREATE TABLE IF NOT EXISTS agent_knowledge_embeddings (
    id SERIAL PRIMARY KEY,
    agent_type VARCHAR(50) NOT NULL,
    knowledge_text TEXT NOT NULL,
    embedding vector(1536) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ====================================
-- 创建 HNSW 索引 (高性能向量搜索)
-- ====================================

-- content_embeddings 表的 HNSW 索引
CREATE INDEX IF NOT EXISTS content_embeddings_embedding_idx 
    ON content_embeddings 
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

-- agent_knowledge_embeddings 表的 HNSW 索引
CREATE INDEX IF NOT EXISTS agent_knowledge_embeddings_embedding_idx 
    ON agent_knowledge_embeddings 
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

-- ====================================
-- 标准索引
-- ====================================

-- content_id 索引
CREATE INDEX IF NOT EXISTS idx_content_embeddings_content_id 
    ON content_embeddings(content_id);

-- content_type 索引
CREATE INDEX IF NOT EXISTS idx_content_embeddings_type 
    ON content_embeddings(content_type);

-- agent_type 索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_type 
    ON agent_knowledge_embeddings(agent_type);

-- ====================================
-- 分析表
-- ====================================

ANALYZE content_embeddings;
ANALYZE agent_knowledge_embeddings;

COMMIT;
