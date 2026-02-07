package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *model.Project) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO projects (user_id, title, type, genre, description, status, settings) 
		 VALUES ($1, $2, $3, $4, $5, 'draft', '{}') 
		 RETURNING id, created_at, updated_at`,
		p.UserID, p.Title, p.Type, p.Genre, p.Description,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepository) GetByID(ctx context.Context, id int) (*model.Project, error) {
	p := &model.Project{}
	err := r.db.QueryRow(ctx,
		`SELECT p.id, p.user_id, p.title, p.type, p.genre, p.description, 
				p.cover_image, p.status, p.word_count, p.settings, p.created_at, p.updated_at,
				COALESCE((SELECT COUNT(*) FROM chapters WHERE project_id = p.id), 0) as chapter_count
		 FROM projects p WHERE p.id = $1`, id,
	).Scan(&p.ID, &p.UserID, &p.Title, &p.Type, &p.Genre, &p.Description,
		&p.CoverImage, &p.Status, &p.WordCount, &p.Settings, &p.CreatedAt, &p.UpdatedAt,
		&p.ChapterCount)
	if err != nil {
		return nil, fmt.Errorf("项目不存在: %w", err)
	}
	return p, nil
}

func (r *ProjectRepository) ListByUser(ctx context.Context, userID int, offset, limit int) ([]model.Project, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM projects WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT p.id, p.user_id, p.title, p.type, p.genre, p.description, 
				p.cover_image, p.status, p.word_count, p.settings, p.created_at, p.updated_at,
				COALESCE((SELECT COUNT(*) FROM chapters WHERE project_id = p.id), 0) as chapter_count
		 FROM projects p 
		 WHERE p.user_id = $1 
		 ORDER BY p.updated_at DESC 
		 LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Type, &p.Genre, &p.Description,
			&p.CoverImage, &p.Status, &p.WordCount, &p.Settings, &p.CreatedAt, &p.UpdatedAt,
			&p.ChapterCount); err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}
	return projects, total, nil
}

func (r *ProjectRepository) Update(ctx context.Context, p *model.Project) error {
	_, err := r.db.Exec(ctx,
		`UPDATE projects SET title=$1, genre=$2, description=$3, status=$4, 
		 word_count=$5, updated_at=NOW() WHERE id=$6`,
		p.Title, p.Genre, p.Description, p.Status, p.WordCount, p.ID)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}

func (r *ProjectRepository) UpdateWordCount(ctx context.Context, projectID int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE projects SET word_count = (
			SELECT COALESCE(SUM(word_count), 0) FROM chapters WHERE project_id = $1
		), updated_at = NOW() WHERE id = $1`, projectID)
	return err
}
