package repository

import (
	"database/sql"
	"github.com/zibianqu/novel-study/internal/model"
)

type ChapterRepository struct {
	db *sql.DB
}

func NewChapterRepository(db *sql.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) Create(chapter *model.Chapter) error {
	query := `
		INSERT INTO chapters (project_id, volume_id, title, content, word_count, sort_order, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		chapter.ProjectID,
		chapter.VolumeID,
		chapter.Title,
		chapter.Content,
		chapter.WordCount,
		chapter.SortOrder,
		chapter.Status,
	).Scan(&chapter.ID, &chapter.CreatedAt, &chapter.UpdatedAt)
}

func (r *ChapterRepository) GetByID(id int) (*model.Chapter, error) {
	chapter := &model.Chapter{}
	query := `
		SELECT id, project_id, volume_id, title, content, word_count, 
		       sort_order, status, locked_by, locked_at, created_at, updated_at
		FROM chapters WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&chapter.ID,
		&chapter.ProjectID,
		&chapter.VolumeID,
		&chapter.Title,
		&chapter.Content,
		&chapter.WordCount,
		&chapter.SortOrder,
		&chapter.Status,
		&chapter.LockedBy,
		&chapter.LockedAt,
		&chapter.CreatedAt,
		&chapter.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return chapter, nil
}

func (r *ChapterRepository) GetByProjectID(projectID int) ([]*model.Chapter, error) {
	query := `
		SELECT id, project_id, volume_id, title, content, word_count,
		       sort_order, status, locked_by, locked_at, created_at, updated_at
		FROM chapters WHERE project_id = $1
		ORDER BY sort_order ASC, created_at ASC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []*model.Chapter
	for rows.Next() {
		chapter := &model.Chapter{}
		err := rows.Scan(
			&chapter.ID,
			&chapter.ProjectID,
			&chapter.VolumeID,
			&chapter.Title,
			&chapter.Content,
			&chapter.WordCount,
			&chapter.SortOrder,
			&chapter.Status,
			&chapter.LockedBy,
			&chapter.LockedAt,
			&chapter.CreatedAt,
			&chapter.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func (r *ChapterRepository) Update(chapter *model.Chapter) error {
	query := `
		UPDATE chapters 
		SET title = $1, content = $2, word_count = $3, sort_order = $4,
		    status = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	return r.db.QueryRow(
		query,
		chapter.Title,
		chapter.Content,
		chapter.WordCount,
		chapter.SortOrder,
		chapter.Status,
		chapter.ID,
	).Scan(&chapter.UpdatedAt)
}

func (r *ChapterRepository) Delete(id int) error {
	query := `DELETE FROM chapters WHERE id = $1`
	result, err := r.db.Exec(query, id)
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

func (r *ChapterRepository) Lock(chapterID, userID int) error {
	query := `UPDATE chapters SET locked_by = $1, locked_at = NOW() WHERE id = $2 AND locked_by IS NULL`
	result, err := r.db.Exec(query, userID, chapterID)
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

func (r *ChapterRepository) Unlock(chapterID, userID int) error {
	query := `UPDATE chapters SET locked_by = NULL, locked_at = NULL WHERE id = $1 AND locked_by = $2`
	result, err := r.db.Exec(query, chapterID, userID)
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
