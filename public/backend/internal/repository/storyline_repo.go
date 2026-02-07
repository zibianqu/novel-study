package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type StorylineRepository struct {
	db *pgxpool.Pool
}

func NewStorylineRepository(db *pgxpool.Pool) *StorylineRepository {
	return &StorylineRepository{db: db}
}

func (r *StorylineRepository) ListByProject(ctx context.Context, projectID int) ([]model.Storyline, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, project_id, line_type, title, content, chapter_start, chapter_end, 
				status, sort_order, parent_id, created_at, updated_at 
		 FROM storylines WHERE project_id = $1 ORDER BY line_type, sort_order`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storylines []model.Storyline
	for rows.Next() {
		var s model.Storyline
		if err := rows.Scan(&s.ID, &s.ProjectID, &s.LineType, &s.Title, &s.Content,
			&s.ChapterStart, &s.ChapterEnd, &s.Status, &s.SortOrder, &s.ParentID,
			&s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		storylines = append(storylines, s)
	}
	return storylines, nil
}

func (r *StorylineRepository) Create(ctx context.Context, s *model.Storyline) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO storylines (project_id, line_type, title, content, chapter_start, chapter_end, 
			status, sort_order, parent_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`,
		s.ProjectID, s.LineType, s.Title, s.Content, s.ChapterStart, s.ChapterEnd,
		s.Status, s.SortOrder, s.ParentID,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StorylineRepository) Update(ctx context.Context, s *model.Storyline) error {
	_, err := r.db.Exec(ctx,
		`UPDATE storylines SET title=$1, content=$2, chapter_start=$3, chapter_end=$4, 
			status=$5, updated_at=NOW() WHERE id=$6`,
		s.Title, s.Content, s.ChapterStart, s.ChapterEnd, s.Status, s.ID)
	return err
}

func (r *StorylineRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM storylines WHERE id = $1", id)
	return err
}
