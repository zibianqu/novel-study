package rag

const (
	// EmbeddingDimension OpenAI text-embedding-ada-002 的维度
	EmbeddingDimension = 1536

	// DefaultTopK 默认返回结果数
	DefaultTopK = 10

	// MaxTopK 最大返回结果数
	MaxTopK = 100

	// MinSimilarityScore 最小相似度阈值
	MinSimilarityScore = 0.5
)

// EmbeddingModel Embedding 模型配置
type EmbeddingModel struct {
	Name      string
	Dimension int
	MaxTokens int
}

var (
	// OpenAIAdaV2 OpenAI Ada v2 模型
	OpenAIAdaV2 = EmbeddingModel{
		Name:      "text-embedding-ada-002",
		Dimension: 1536,
		MaxTokens: 8191,
	}
)
