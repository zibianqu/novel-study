package model

import (
	"time"
)

type Chapter struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	VolumeID  *int      `json:"volume_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	WordCount int       `json:"word_count"`
	SortOrder int       `json:"sort_order"`
	Status    string    `json:"status"` // draft, published
	LockedBy  *int      `json:"locked_by"`
	LockedAt  *time.Time `json:"locked_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Volume struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateChapterRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	VolumeID  *int   `json:"volume_id"`
	Title     string `json:"title" binding:"required,min=1,max=200"`
	Content   string `json:"content"`
	SortOrder int    `json:"sort_order"`
}

type UpdateChapterRequest struct {
	Title     string `json:"title" binding:"omitempty,min=1,max=200"`
	Content   string `json:"content"`
	Status    string `json:"status" binding:"omitempty,oneof=draft published"`
	SortOrder int    `json:"sort_order"`
}
