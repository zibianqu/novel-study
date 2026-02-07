package repository

import (
	"database/sql"
	"github.com/zibianqu/novel-study/internal/model"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *model.Project) error {
	query := `
		INSERT INTO projects (user_id, title, type, genre, description, cover_image, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		project.UserID,
		project.Title,
		project.Type,
		project.Genre,
		project.Description,
		project.CoverImage,
		project.Status,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
}

func (r *ProjectRepository) GetByID(id int) (*model.Project, error) {
	project := &model.Project{}
	query := `
		SELECT id, user_id, title, type, genre, description, cover_image, 
		       status, word_count, settings, created_at, updated_at
		FROM projects WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.UserID,
		&project.Title,
		&project.Type,
		&project.Genre,
		&project.Description,
		&project.CoverImage,
		&project.Status,
		&project.WordCount,
		&project.Settings,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) GetByUserID(userID int) ([]*model.Project, error) {
	query := `
		SELECT id, user_id, title, type, genre, description, cover_image,
		       status, word_count, settings, created_at, updated_at
		FROM projects WHERE user_id = $1
		ORDER BY updated_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		project := &model.Project{}
		err := rows.Scan(
			&project.ID,
			&project.UserID,
			&project.Title,
			&project.Type,
			&project.Genre,
			&project.Description,
			&project.CoverImage,
			&project.Status,
			&project.WordCount,
			&project.Settings,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (r *ProjectRepository) Update(project *model.Project) error {
	query := `
		UPDATE projects 
		SET title = $1, genre = $2, description = $3, cover_image = $4, 
		    status = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		project.Title,
		project.Genre,
		project.Description,
		project.CoverImage,
		project.Status,
		project.ID,
		project.UserID,
	).Scan(&project.UpdatedAt)
}

func (r *ProjectRepository) Delete(id, userID int) error {
	query := `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProjectRepository) UpdateWordCount(projectID, wordCount int) error {
	query := `UPDATE projects SET word_count = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, wordCount, projectID)
	return err
}
