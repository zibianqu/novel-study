package repository

import (
	"database/sql"
	"github.com/zibianqu/novel-study/internal/model"
)

type KnowledgeRepository struct {
	db *sql.DB
}

func NewKnowledgeRepository(db *sql.DB) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

func (r *KnowledgeRepository) Create(kb *model.KnowledgeBase) error {
	query := `
		INSERT INTO knowledge_base (project_id, title, content, type, tags, is_vectorized, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, false, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		kb.ProjectID,
		kb.Title,
		kb.Content,
		kb.Type,
		kb.Tags,
	).Scan(&kb.ID, &kb.CreatedAt, &kb.UpdatedAt)
}

func (r *KnowledgeRepository) GetByID(id int) (*model.KnowledgeBase, error) {
	kb := &model.KnowledgeBase{}
	query := `
		SELECT id, project_id, title, content, type, tags, is_vectorized, created_at, updated_at
		FROM knowledge_base WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&kb.ID,
		&kb.ProjectID,
		&kb.Title,
		&kb.Content,
		&kb.Type,
		&kb.Tags,
		&kb.IsVectorized,
		&kb.CreatedAt,
		&kb.UpdatedAt,
	)
	return kb, err
}

func (r *KnowledgeRepository) GetByProjectID(projectID int) ([]*model.KnowledgeBase, error) {
	query := `
		SELECT id, project_id, title, content, type, tags, is_vectorized, created_at, updated_at
		FROM knowledge_base WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.KnowledgeBase
	for rows.Next() {
		kb := &model.KnowledgeBase{}
		err := rows.Scan(
			&kb.ID,
			&kb.ProjectID,
			&kb.Title,
			&kb.Content,
			&kb.Type,
			&kb.Tags,
			&kb.IsVectorized,
			&kb.CreatedAt,
			&kb.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, kb)
	}
	return items, nil
}

func (r *KnowledgeRepository) Update(kb *model.KnowledgeBase) error {
	query := `
		UPDATE knowledge_base 
		SET title = $1, content = $2, type = $3, tags = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		kb.Title,
		kb.Content,
		kb.Type,
		kb.Tags,
		kb.ID,
	).Scan(&kb.UpdatedAt)
}

func (r *KnowledgeRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM knowledge_base WHERE id = $1", id)
	return err
}

func (r *KnowledgeRepository) MarkVectorized(id int) error {
	_, err := r.db.Exec("UPDATE knowledge_base SET is_vectorized = true WHERE id = $1", id)
	return err
}
