package rag

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pgvector/pgvector-go"
)

// VectorStore 向量存储
type VectorStore struct {
	db *sql.DB
}

// NewVectorStore 创建向量存储
func NewVectorStore(db *sql.DB) *VectorStore {
	return &VectorStore{db: db}
}

// Document 文档结构
type Document struct {
	ID        int
	Content   string
	Metadata  map[string]interface{}
	Embedding []float32
	Score     float64
}

// AddDocument 添加文档
func (vs *VectorStore) AddDocument(ctx context.Context, projectID int, content string, embedding []float32, metadata map[string]interface{}) (int, error) {
	// 将 metadata 转为 JSON
	metadataJSON := "{}"
	if metadata != nil {
		bytes, err := json.Marshal(metadata)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJSON = string(bytes)
	}

	query := `
		INSERT INTO knowledge_vectors (project_id, content, embedding, metadata, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`

	var id int
	err := vs.db.QueryRowContext(ctx, query, projectID, content, pgvector.NewVector(embedding), metadataJSON).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert document: %w", err)
	}

	return id, nil
}

// SimilaritySearch 相似度搜索
func (vs *VectorStore) SimilaritySearch(ctx context.Context, projectID int, queryEmbedding []float32, topK int) ([]*Document, error) {
	if topK <= 0 {
		topK = 10
	}
	if topK > 100 {
		topK = 100
	}

	query := `
		SELECT id, content, metadata, 
		       1 - (embedding <=> $1) as similarity
		FROM knowledge_vectors
		WHERE project_id = $2
		ORDER BY embedding <=> $1
		LIMIT $3
	`

	rows, err := vs.db.QueryContext(ctx, query, pgvector.NewVector(queryEmbedding), projectID, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer rows.Close()

	var docs []*Document
	for rows.Next() {
		doc := &Document{
			Metadata: make(map[string]interface{}),
		}
		var metadataJSON string
		err := rows.Scan(&doc.ID, &doc.Content, &metadataJSON, &doc.Score)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// 解析 metadata JSON
		if metadataJSON != "" && metadataJSON != "{}" {
			if err := json.Unmarshal([]byte(metadataJSON), &doc.Metadata); err != nil {
				// 日志错误但不阻断
				fmt.Printf("Warning: failed to unmarshal metadata: %v\n", err)
			}
		}

		docs = append(docs, doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return docs, nil
}

// DeleteDocument 删除文档
func (vs *VectorStore) DeleteDocument(ctx context.Context, id int) error {
	_, err := vs.db.ExecContext(ctx, "DELETE FROM knowledge_vectors WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// DeleteProjectDocuments 删除项目所有文档
func (vs *VectorStore) DeleteProjectDocuments(ctx context.Context, projectID int) error {
	_, err := vs.db.ExecContext(ctx, "DELETE FROM knowledge_vectors WHERE project_id = $1", projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project documents: %w", err)
	}
	return nil
}
