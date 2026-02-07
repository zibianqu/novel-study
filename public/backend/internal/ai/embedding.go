package ai

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
	openai "github.com/sashabaranov/go-openai"

	"github.com/novelforge/backend/internal/config"
)

// EmbeddingService 向量化服务
type EmbeddingService struct {
	client *openai.Client
	db     *pgxpool.Pool
	cfg    config.EmbeddingConfig
	ragCfg config.RAGConfig
}

func NewEmbeddingService(client *openai.Client, db *pgxpool.Pool, cfg config.AIConfig) *EmbeddingService {
	return &EmbeddingService{
		client: client,
		db:     db,
		cfg:    cfg.Embedding,
		ragCfg: cfg.RAG,
	}
}

// ==================== 文本分块 ====================

// ChunkText 将长文本分块
func (s *EmbeddingService) ChunkText(text string, chunkSize, overlap int) []string {
	if chunkSize == 0 {
		chunkSize = s.ragCfg.ChunkSize
	}
	if overlap == 0 {
		overlap = s.ragCfg.ChunkOverlap
	}

	runes := []rune(text)
	length := len(runes)
	if length <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0
	for start < length {
		end := start + chunkSize
		if end > length {
			end = length
		}

		// 尝试在句子边界处分割
		chunk := string(runes[start:end])
		if end < length {
			// 找最后一个句号/问号/叹号作为分割点
			lastSentEnd := strings.LastIndexAny(chunk, "。！？\n")
			if lastSentEnd > chunkSize/2 {
				end = start + utf8.RuneCountInString(chunk[:lastSentEnd+3]) // +3 for 中文标点的UTF8长度
				chunk = string(runes[start:end])
			}
		}

		chunks = append(chunks, strings.TrimSpace(chunk))
		start = end - overlap
		if start < 0 {
			start = 0
		}
	}

	return chunks
}

// ==================== Embedding 生成 ====================

// CreateEmbedding 生成文本的向量表示
func (s *EmbeddingService) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(s.cfg.Model),
	})
	if err != nil {
		return nil, fmt.Errorf("创建Embedding失败: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("Embedding响应为空")
	}
	return resp.Data[0].Embedding, nil
}

// CreateEmbeddings 批量生成向量
func (s *EmbeddingService) CreateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	resp, err := s.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: texts,
		Model: openai.EmbeddingModel(s.cfg.Model),
	})
	if err != nil {
		return nil, fmt.Errorf("批量创建Embedding失败: %w", err)
	}

	embeddings := make([][]float32, len(resp.Data))
	for i, d := range resp.Data {
		embeddings[i] = d.Embedding
	}
	return embeddings, nil
}

// ==================== 内容向量入库 ====================

// IndexChapterContent 将章节内容向量化并存入 pgvector
func (s *EmbeddingService) IndexChapterContent(ctx context.Context, projectID, chapterID int, content string) ([]int, error) {
	// 1. 删除旧的向量记录
	_, err := s.db.Exec(ctx,
		"DELETE FROM content_embeddings WHERE chapter_id = $1", chapterID)
	if err != nil {
		return nil, fmt.Errorf("删除旧向量失败: %w", err)
	}

	// 2. 分块
	chunks := s.ChunkText(content, 0, 0)
	if len(chunks) == 0 {
		return nil, nil
	}

	// 3. 批量生成向量
	embeddings, err := s.CreateEmbeddings(ctx, chunks)
	if err != nil {
		return nil, err
	}

	// 4. 批量写入数据库
	var ids []int
	for i, chunk := range chunks {
		var id int
		err := s.db.QueryRow(ctx,
			`INSERT INTO content_embeddings (project_id, chapter_id, chunk_text, chunk_index, embedding)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			projectID, chapterID, chunk, i, pgvectorString(embeddings[i]),
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("写入向量失败: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// IndexAgentKnowledge 将Agent知识条目向量化并存入
func (s *EmbeddingService) IndexAgentKnowledge(ctx context.Context, agentID, itemID int, content string) error {
	// 删除旧记录
	_, err := s.db.Exec(ctx,
		"DELETE FROM agent_knowledge_embeddings WHERE item_id = $1", itemID)
	if err != nil {
		return fmt.Errorf("删除旧知识向量失败: %w", err)
	}

	// 分块
	chunks := s.ChunkText(content, 0, 0)
	if len(chunks) == 0 {
		return nil
	}

	// 生成向量
	embeddings, err := s.CreateEmbeddings(ctx, chunks)
	if err != nil {
		return err
	}

	// 写入
	for i, chunk := range chunks {
		_, err := s.db.Exec(ctx,
			`INSERT INTO agent_knowledge_embeddings (item_id, agent_id, chunk_text, chunk_index, embedding)
			 VALUES ($1, $2, $3, $4, $5)`,
			itemID, agentID, chunk, i, pgvectorString(embeddings[i]),
		)
		if err != nil {
			return fmt.Errorf("写入知识向量失败: %w", err)
		}
	}

	return nil
}

// ==================== RAG 检索 ====================

// SearchContent 在项目内容中进行向量检索
func (s *EmbeddingService) SearchContent(ctx context.Context, projectID int, query string, topK int) ([]RAGResult, error) {
	if topK == 0 {
		topK = s.ragCfg.TopK
	}

	// 生成查询向量
	queryEmb, err := s.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	// 向量检索
	rows, err := s.db.Query(ctx,
		`SELECT id, chapter_id, chunk_text, chunk_index, 
				1 - (embedding <=> $1::vector) as similarity
		 FROM content_embeddings 
		 WHERE project_id = $2 
		 ORDER BY embedding <=> $1::vector 
		 LIMIT $3`,
		pgvectorString(queryEmb), projectID, topK,
	)
	if err != nil {
		return nil, fmt.Errorf("向量检索失败: %w", err)
	}
	defer rows.Close()

	var results []RAGResult
	for rows.Next() {
		var r RAGResult
		if err := rows.Scan(&r.ID, &r.ChapterID, &r.Text, &r.ChunkIndex, &r.Similarity); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// SearchAgentKnowledge 在Agent知识库中进行向量检索
func (s *EmbeddingService) SearchAgentKnowledge(ctx context.Context, agentID int, query string, topK int) ([]RAGResult, error) {
	if topK == 0 {
		topK = s.ragCfg.TopK
	}

	queryEmb, err := s.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT e.id, e.item_id, e.chunk_text, e.chunk_index,
				1 - (e.embedding <=> $1::vector) as similarity
		 FROM agent_knowledge_embeddings e
		 WHERE e.agent_id = $2
		 ORDER BY e.embedding <=> $1::vector
		 LIMIT $3`,
		pgvectorString(queryEmb), agentID, topK,
	)
	if err != nil {
		return nil, fmt.Errorf("知识库检索失败: %w", err)
	}
	defer rows.Close()

	var results []RAGResult
	for rows.Next() {
		var r RAGResult
		if err := rows.Scan(&r.ID, &r.ItemID, &r.Text, &r.ChunkIndex, &r.Similarity); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

// RAGResult 检索结果
type RAGResult struct {
	ID         int     `json:"id"`
	ChapterID  *int    `json:"chapter_id,omitempty"`
	ItemID     *int    `json:"item_id,omitempty"`
	Text       string  `json:"text"`
	ChunkIndex int     `json:"chunk_index"`
	Similarity float64 `json:"similarity"`
}

// pgvectorString 将 float32 切片转为 pgvector 格式字符串
func pgvectorString(v []float32) string {
	parts := make([]string, len(v))
	for i, f := range v {
		parts[i] = fmt.Sprintf("%f", f)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
