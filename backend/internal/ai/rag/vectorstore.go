package rag

import (
	"context"
	"database/sql"
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
		// TODO: 实际应该序列化为 JSON
	}

	query := `
		INSERT INTO knowledge_vectors (project_id, content, embedding, metadata, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`

	var id int
	err := vs.db.QueryRowContext(ctx, query, projectID, content, pgvector.NewVector(embedding), metadataJSON).Scan(&id)
	return id, err
}

// SimilaritySearch 相似度搜索
func (vs *VectorStore) SimilaritySearch(ctx context.Context, projectID int, queryEmbedding []float32, topK int) ([]*Document, error) {
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
		return nil, err
	}
	defer rows.Close()

	var docs []*Document
	for rows.Next() {
		doc := &Document{}
		var metadataJSON string
		err := rows.Scan(&doc.ID, &doc.Content, &metadataJSON, &doc.Score)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

// DeleteDocument 删除文档
func (vs *VectorStore) DeleteDocument(ctx context.Context, id int) error {
	_, err := vs.db.ExecContext(ctx, "DELETE FROM knowledge_vectors WHERE id = $1", id)
	return err
}

// DeleteProjectDocuments 删除项目所有文档
func (vs *VectorStore) DeleteProjectDocuments(ctx context.Context, projectID int) error {
	_, err := vs.db.ExecContext(ctx, "DELETE FROM knowledge_vectors WHERE project_id = $1", projectID)
	return err
}
