package rag

import (
	"context"

	"github.com/zibianqu/novel-study/internal/ai/openai"
)

// EmbeddingService 向量嵌入服务
type EmbeddingService struct {
	client *openai.Client
}

// NewEmbeddingService 创建嵌入服务
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		client: openai.NewClient(apiKey),
	}
}

// Embed 生成向量嵌入
func (s *EmbeddingService) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return s.client.CreateEmbedding(ctx, texts)
}

// EmbedSingle 生成单个文本的向量
func (s *EmbeddingService) EmbedSingle(ctx context.Context, text string) ([]float32, error) {
	embeddings, err := s.Embed(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) == 0 {
		return nil, nil
	}
	return embeddings[0], nil
}
