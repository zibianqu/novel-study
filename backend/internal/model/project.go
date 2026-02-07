package model

import (
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"` // novel_long, novel_short, copywriting
	Genre       string    `json:"genre"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	Status      string    `json:"status"` // draft, writing, completed
	WordCount   int       `json:"word_count"`
	Settings    string    `json:"settings"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Type        string `json:"type" binding:"required,oneof=novel_long novel_short copywriting"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

type UpdateProjectRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=200"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Status      string `json:"status" binding:"omitempty,oneof=draft writing completed"`
}
