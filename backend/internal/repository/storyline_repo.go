package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

// Storyline 三线模型
type Storyline struct {
	ID           int       `db:"id" json:"id"`
	ProjectID    int       `db:"project_id" json:"project_id"`
	LineType     string    `db:"line_type" json:"line_type"` // skyline/groundline/plotline
	Title        string    `db:"title" json:"title"`
	Content      string    `db:"content" json:"content"`
	ChapterRange string    `db:"chapter_range" json:"chapter_range"`
	Status       string    `db:"status" json:"status"`
	SortOrder    int       `db:"sort_order" json:"sort_order"`
	ParentID     *int      `db:"parent_id" json:"parent_id,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// StorylineRepository 三线仓库
type StorylineRepository struct {
	db *sqlx.DB
}

// NewStorylineRepository 创建三线仓库
func NewStorylineRepository(db *sqlx.DB) *StorylineRepository {
	return &StorylineRepository{
		db: db,
	}
}

// GetByID 按 ID 获取
func (r *StorylineRepository) GetByID(ctx context.Context, id int) (*Storyline, error) {
	var storyline Storyline
	query := `SELECT * FROM storylines WHERE id = $1`
	err := r.db.GetContext(ctx, &storyline, query, id)
	return &storyline, err
}

// GetByProjectID 按项目 ID 获取所有三线
func (r *StorylineRepository) GetByProjectID(ctx context.Context, projectID int) ([]*Storyline, error) {
	var storylines []*Storyline
	query := `SELECT * FROM storylines WHERE project_id = $1 ORDER BY line_type, sort_order`
	err := r.db.SelectContext(ctx, &storylines, query, projectID)
	return storylines, err
}

// GetByProjectIDAndType 按项目 ID 和类型获取
func (r *StorylineRepository) GetByProjectIDAndType(ctx context.Context, projectID int, lineType string) ([]*Storyline, error) {
	var storylines []*Storyline
	query := `SELECT * FROM storylines WHERE project_id = $1 AND line_type = $2 ORDER BY sort_order`
	err := r.db.SelectContext(ctx, &storylines, query, projectID, lineType)
	return storylines, err
}

// Update 更新三线
func (r *StorylineRepository) Update(ctx context.Context, storyline *Storyline) error {
	query := `
		UPDATE storylines 
		SET title = $1, content = $2, status = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, storyline.Title, storyline.Content, storyline.Status, storyline.ID)
	return err
}

// Create 创建三线
func (r *StorylineRepository) Create(ctx context.Context, storyline *Storyline) error {
	query := `
		INSERT INTO storylines (project_id, line_type, title, content, chapter_range, status, sort_order, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		storyline.ProjectID,
		storyline.LineType,
		storyline.Title,
		storyline.Content,
		storyline.ChapterRange,
		storyline.Status,
		storyline.SortOrder,
		storyline.ParentID,
	).Scan(&storyline.ID, &storyline.CreatedAt, &storyline.UpdatedAt)
}
