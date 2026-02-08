package rag

import (
	"context"
	"fmt"
	"strings"
)

// Retriever RAG 检索器
type Retriever struct {
	embedding   *EmbeddingService
	vectorStore *VectorStore
}

// NewRetriever 创建检索器
func NewRetriever(embedding *EmbeddingService, vectorStore *VectorStore) *Retriever {
	return &Retriever{
		embedding:   embedding,
		vectorStore: vectorStore,
	}
}

// Retrieve 检索相关文档
func (r *Retriever) Retrieve(ctx context.Context, projectID int, query string, topK int) ([]*Document, error) {
	// 1. 生成查询向量
	queryEmbedding, err := r.embedding.EmbedSingle(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	// 2. 相似度搜索
	docs, err := r.vectorStore.SimilaritySearch(ctx, projectID, queryEmbedding, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return docs, nil
}

// BuildContext 构建 RAG 上下文
func (r *Retriever) BuildContext(docs []*Document) string {
	if len(docs) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("相关上下文信息：\n\n")

	for i, doc := range docs {
		builder.WriteString(fmt.Sprintf("[%d] 相关度: %.2f\n", i+1, doc.Score))
		builder.WriteString(doc.Content)
		builder.WriteString("\n\n---\n\n")
	}

	return builder.String()
}
