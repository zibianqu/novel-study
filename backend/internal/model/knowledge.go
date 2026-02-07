package model

import (
	"time"
)

// KnowledgeBase 知识库
type KnowledgeBase struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Type        string    `json:"type"` // character, worldview, plot, custom
	Tags        string    `json:"tags"` // JSON array
	IsVectorized bool     `json:"is_vectorized"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// KnowledgeVector 知识向量
type KnowledgeVector struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Content   string    `json:"content"`
	Embedding []float32 `json:"-"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateKnowledgeRequest 创建知识请求
type CreateKnowledgeRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Tags      string `json:"tags"`
}
