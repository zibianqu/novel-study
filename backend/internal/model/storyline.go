package model

import (
	"time"
)

// Storyline 三线（天线/地线/剧情线）
type Storyline struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	LineType     string    `json:"line_type"` // skyline, groundline, plotline
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	ChapterRange string    `json:"chapter_range"` // 例如: "[1,10]"
	Status       string    `json:"status"`        // planned, ongoing, completed
	SortOrder    int       `json:"sort_order"`
	ParentID     *int      `json:"parent_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// StorylineConvergence 三线交汇点
type StorylineConvergence struct {
	ID                 int       `json:"id"`
	ProjectID          int       `json:"project_id"`
	Name               string    `json:"name"`
	SkylineMeaning     string    `json:"skyline_meaning"`
	GroundlineMeaning  string    `json:"groundline_meaning"`
	PlotlineMeaning    string    `json:"plotline_meaning"`
	ChapterID          *int      `json:"chapter_id"`
	CreatedAt          time.Time `json:"created_at"`
}
