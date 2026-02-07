package repository

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type ChapterRepository struct {
	db *pgxpool.Pool
}

func NewChapterRepository(db *pgxpool.Pool) *ChapterRepository {
	return &ChapterRepository{db: db}
}

func (r *ChapterRepository) Create(ctx context.Context, ch *model.Chapter) error {
	ch.WordCount = utf8.RuneCountInString(ch.Content)
	return r.db.QueryRow(ctx,
		`INSERT INTO chapters (project_id, volume_id, title, content, word_count, sort_order, status) 
		 VALUES ($1, $2, $3, $4, $5, 
		 	COALESCE((SELECT MAX(sort_order)+1 FROM chapters WHERE project_id=$1), 0), 
		 	'draft') 
		 RETURNING id, sort_order, created_at, updated_at`,
		ch.ProjectID, ch.VolumeID, ch.Title, ch.Content, ch.WordCount,
	).Scan(&ch.ID, &ch.SortOrder, &ch.CreatedAt, &ch.UpdatedAt)
}

func (r *ChapterRepository) GetByID(ctx context.Context, id int) (*model.Chapter, error) {
	ch := &model.Chapter{}
	err := r.db.QueryRow(ctx,
		`SELECT id, project_id, volume_id, title, content, word_count, sort_order, 
				status, locked_by, locked_at, created_at, updated_at 
		 FROM chapters WHERE id = $1`, id,
	).Scan(&ch.ID, &ch.ProjectID, &ch.VolumeID, &ch.Title, &ch.Content, &ch.WordCount,
		&ch.SortOrder, &ch.Status, &ch.LockedBy, &ch.LockedAt, &ch.CreatedAt, &ch.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("章节不存在: %w", err)
	}
	return ch, nil
}

func (r *ChapterRepository) ListByProject(ctx context.Context, projectID int) ([]model.Volume, []model.Chapter, error) {
	// 获取卷
	volRows, err := r.db.Query(ctx,
		`SELECT id, project_id, title, summary, sort_order, created_at 
		 FROM volumes WHERE project_id = $1 ORDER BY sort_order`, projectID)
	if err != nil {
		return nil, nil, err
	}
	defer volRows.Close()

	var volumes []model.Volume
	for volRows.Next() {
		var v model.Volume
		if err := volRows.Scan(&v.ID, &v.ProjectID, &v.Title, &v.Summary, &v.SortOrder, &v.CreatedAt); err != nil {
			return nil, nil, err
		}
		volumes = append(volumes, v)
	}

	// 获取章节（不含content以减小数据量）
	chRows, err := r.db.Query(ctx,
		`SELECT id, project_id, volume_id, title, word_count, sort_order, status, 
				locked_by, locked_at, created_at, updated_at 
		 FROM chapters WHERE project_id = $1 ORDER BY sort_order`, projectID)
	if err != nil {
		return nil, nil, err
	}
	defer chRows.Close()

	var chapters []model.Chapter
	for chRows.Next() {
		var ch model.Chapter
		if err := chRows.Scan(&ch.ID, &ch.ProjectID, &ch.VolumeID, &ch.Title, &ch.WordCount,
			&ch.SortOrder, &ch.Status, &ch.LockedBy, &ch.LockedAt, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, nil, err
		}
		chapters = append(chapters, ch)
	}

	return volumes, chapters, nil
}

func (r *ChapterRepository) Update(ctx context.Context, ch *model.Chapter) error {
	ch.WordCount = utf8.RuneCountInString(ch.Content)
	_, err := r.db.Exec(ctx,
		`UPDATE chapters SET title=$1, content=$2, word_count=$3, status=$4, updated_at=NOW() 
		 WHERE id=$5`,
		ch.Title, ch.Content, ch.WordCount, ch.Status, ch.ID)
	return err
}

func (r *ChapterRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM chapters WHERE id = $1", id)
	return err
}

func (r *ChapterRepository) Lock(ctx context.Context, chapterID, userID int) error {
	_, err := r.db.Exec(ctx,
		"UPDATE chapters SET locked_by=$1, locked_at=NOW() WHERE id=$2 AND locked_by IS NULL",
		userID, chapterID)
	return err
}

func (r *ChapterRepository) Unlock(ctx context.Context, chapterID, userID int) error {
	_, err := r.db.Exec(ctx,
		"UPDATE chapters SET locked_by=NULL, locked_at=NULL WHERE id=$1 AND locked_by=$2",
		chapterID, userID)
	return err
}

// SaveVersion 保存章节版本
func (r *ChapterRepository) SaveVersion(ctx context.Context, v *model.ChapterVersion) error {
	// 获取下一个版本号
	var nextVersion int
	err := r.db.QueryRow(ctx,
		"SELECT COALESCE(MAX(version_num), 0) + 1 FROM chapter_versions WHERE chapter_id = $1",
		v.ChapterID,
	).Scan(&nextVersion)
	if err != nil {
		return err
	}

	v.VersionNum = nextVersion
	return r.db.QueryRow(ctx,
		`INSERT INTO chapter_versions (chapter_id, version_num, content, delta_content, 
			delta_position, agent_outputs, embedding_ids, graph_changes, created_by) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		 RETURNING id, created_at`,
		v.ChapterID, v.VersionNum, v.Content, v.DeltaContent,
		v.DeltaPosition, v.AgentOutputs, v.EmbeddingIDs, v.GraphChanges, v.CreatedBy,
	).Scan(&v.ID, &v.CreatedAt)
}

func (r *ChapterRepository) ListVersions(ctx context.Context, chapterID int) ([]model.ChapterVersion, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, chapter_id, version_num, delta_content, agent_outputs, created_by, created_at 
		 FROM chapter_versions WHERE chapter_id = $1 ORDER BY version_num DESC`, chapterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []model.ChapterVersion
	for rows.Next() {
		var v model.ChapterVersion
		if err := rows.Scan(&v.ID, &v.ChapterID, &v.VersionNum, &v.DeltaContent,
			&v.AgentOutputs, &v.CreatedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (r *ChapterRepository) GetVersion(ctx context.Context, versionID int) (*model.ChapterVersion, error) {
	v := &model.ChapterVersion{}
	err := r.db.QueryRow(ctx,
		`SELECT id, chapter_id, version_num, content, delta_content, 
				delta_position, agent_outputs, embedding_ids, graph_changes, created_by, created_at 
		 FROM chapter_versions WHERE id = $1`, versionID,
	).Scan(&v.ID, &v.ChapterID, &v.VersionNum, &v.Content, &v.DeltaContent,
		&v.DeltaPosition, &v.AgentOutputs, &v.EmbeddingIDs, &v.GraphChanges, &v.CreatedBy, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}
