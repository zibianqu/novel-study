-- 启用pgvector扩展
-- 执行时间: 2026-02-08

CREATE EXTENSION IF NOT EXISTS vector;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector') THEN
        RAISE EXCEPTION 'pgvector extension not installed';
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'content_embeddings' 
        AND column_name = 'embedding' 
        AND data_type = 'USER-DEFINED'
    ) THEN
        ALTER TABLE content_embeddings ALTER COLUMN embedding TYPE vector(1536) USING embedding::vector(1536);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'agent_knowledge_embeddings' 
        AND column_name = 'embedding' 
        AND data_type = 'USER-DEFINED'
    ) THEN
        ALTER TABLE agent_knowledge_embeddings ALTER COLUMN embedding TYPE vector(1536) USING embedding::vector(1536);
    END IF;
END $$;

ALTER TABLE content_embeddings ADD CONSTRAINT check_embedding_dim 
    CHECK (vector_dims(embedding) = 1536);
    
ALTER TABLE agent_knowledge_embeddings ADD CONSTRAINT check_agent_embedding_dim 
    CHECK (vector_dims(embedding) = 1536);

CREATE OR REPLACE FUNCTION search_similar_content(
    query_embedding vector(1536),
    p_project_id INT,
    limit_count INT DEFAULT 5
)
RETURNS TABLE (
    id INT,
    chapter_id INT,
    chunk_text TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ce.id,
        ce.chapter_id,
        ce.chunk_text,
        1 - (ce.embedding <=> query_embedding) as similarity
    FROM content_embeddings ce
    WHERE ce.project_id = p_project_id
    ORDER BY ce.embedding <=> query_embedding
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION search_agent_knowledge(
    query_embedding vector(1536),
    p_agent_id INT,
    limit_count INT DEFAULT 3
)
RETURNS TABLE (
    id INT,
    item_id INT,
    chunk_text TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ake.id,
        ake.item_id,
        ake.chunk_text,
        1 - (ake.embedding <=> query_embedding) as similarity
    FROM agent_knowledge_embeddings ake
    WHERE ake.agent_id = p_agent_id
    ORDER BY ake.embedding <=> query_embedding
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION search_similar_content IS 'RAG内容相似度搜索';
COMMENT ON FUNCTION search_agent_knowledge IS 'Agent知识库相似度搜索';
