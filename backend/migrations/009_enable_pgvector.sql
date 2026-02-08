-- NovelForge AI - 启用 pgvector 扩展
-- 创建日期: 2026-02-08

-- =====================================================
-- 1. 启用 pgvector 扩展
-- =====================================================

-- 检查并启用 pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- 验证扩展已安装
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_extension WHERE extname = 'vector'
    ) THEN
        RAISE EXCEPTION 'pgvector extension not installed. Please install it first.';
    END IF;
END $$;

-- =====================================================
-- 2. 创建/更新 content_embeddings 表
-- =====================================================

-- 如果表不存在，创建它
CREATE TABLE IF NOT EXISTS content_embeddings (
    id              SERIAL PRIMARY KEY,
    project_id      INTEGER NOT NULL,
    chapter_id      INTEGER REFERENCES chapters(id) ON DELETE CASCADE,
    chunk_text      TEXT NOT NULL,
    chunk_index     INTEGER DEFAULT 0,
    embedding       VECTOR(1536),  -- OpenAI text-embedding-3-small
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 3. 创建 HNSW 索引（高性能向量检索）
-- =====================================================

-- 删除旧的 IVFFlat 索引（如果存在）
DROP INDEX IF EXISTS content_embeddings_embedding_idx;

-- 创建 HNSW 索引（更快的查询速度）
CREATE INDEX IF NOT EXISTS idx_content_embeddings_hnsw ON content_embeddings 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- m = 16: 每个节点的最大连接数（默认值，平衡性能和精度）
-- ef_construction = 64: 构建时的搜索宽度（较高的值提高精度）

-- =====================================================
-- 4. 添加其他必要索引
-- =====================================================

-- project_id 索引（用于过滤）
CREATE INDEX IF NOT EXISTS idx_content_embeddings_project ON content_embeddings(project_id);

-- chapter_id 索引
CREATE INDEX IF NOT EXISTS idx_content_embeddings_chapter ON content_embeddings(chapter_id);

-- 组合索引（project + chapter）
CREATE INDEX IF NOT EXISTS idx_content_embeddings_project_chapter 
ON content_embeddings(project_id, chapter_id);

-- =====================================================
-- 5. 创建 Agent Knowledge Embeddings 表
-- =====================================================

CREATE TABLE IF NOT EXISTS agent_knowledge_embeddings (
    id              SERIAL PRIMARY KEY,
    item_id         INTEGER REFERENCES agent_knowledge_items(id) ON DELETE CASCADE,
    agent_id        INTEGER NOT NULL,
    chunk_text      TEXT NOT NULL,
    chunk_index     INTEGER DEFAULT 0,
    embedding       VECTOR(1536),
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW()
);

-- Agent Knowledge HNSW 索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_embeddings_hnsw 
ON agent_knowledge_embeddings 
USING hnsw (embedding vector_cosine_ops)
WITH (m = 16, ef_construction = 64);

-- agent_id 索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_embeddings_agent 
ON agent_knowledge_embeddings(agent_id);

-- item_id 索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_embeddings_item 
ON agent_knowledge_embeddings(item_id);

-- =====================================================
-- 6. 创建向量搜索函数
-- =====================================================

-- 搜索项目内容
CREATE OR REPLACE FUNCTION search_project_content(
    p_project_id INTEGER,
    p_query_embedding VECTOR(1536),
    p_limit INTEGER DEFAULT 5,
    p_threshold FLOAT DEFAULT 0.7
)
RETURNS TABLE (
    id INTEGER,
    chapter_id INTEGER,
    chunk_text TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ce.id,
        ce.chapter_id,
        ce.chunk_text,
        1 - (ce.embedding <=> p_query_embedding) AS similarity
    FROM content_embeddings ce
    WHERE ce.project_id = p_project_id
        AND 1 - (ce.embedding <=> p_query_embedding) >= p_threshold
    ORDER BY ce.embedding <=> p_query_embedding
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- 搜索 Agent 知识库
CREATE OR REPLACE FUNCTION search_agent_knowledge(
    p_agent_id INTEGER,
    p_query_embedding VECTOR(1536),
    p_limit INTEGER DEFAULT 3,
    p_threshold FLOAT DEFAULT 0.7
)
RETURNS TABLE (
    id INTEGER,
    item_id INTEGER,
    chunk_text TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ake.id,
        ake.item_id,
        ake.chunk_text,
        1 - (ake.embedding <=> p_query_embedding) AS similarity
    FROM agent_knowledge_embeddings ake
    WHERE ake.agent_id = p_agent_id
        AND 1 - (ake.embedding <=> p_query_embedding) >= p_threshold
    ORDER BY ake.embedding <=> p_query_embedding
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql;

-- =====================================================
-- 7. 更新统计信息
-- =====================================================

ANALYZE content_embeddings;
ANALYZE agent_knowledge_embeddings;

-- =====================================================
-- 8. 添加注释
-- =====================================================

COMMENT ON TABLE content_embeddings IS '项目内容向量嵌入表';
COMMENT ON TABLE agent_knowledge_embeddings IS 'Agent知识库向量嵌入表';
COMMENT ON INDEX idx_content_embeddings_hnsw IS 'HNSW向量索引 - 高性能相似度搜索';
COMMENT ON FUNCTION search_project_content IS '项目内容向量搜索函数';
COMMENT ON FUNCTION search_agent_knowledge IS 'Agent知识库向量搜索函数';

-- 完成
SELECT 'pgvector 启用并优化完成' AS status;
