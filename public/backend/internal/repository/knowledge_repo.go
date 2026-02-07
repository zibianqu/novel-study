package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type KnowledgeRepository struct {
	db *pgxpool.Pool
}

func NewKnowledgeRepository(db *pgxpool.Pool) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

// ========== 分类 ==========

func (r *KnowledgeRepository) ListCategories(ctx context.Context, agentID int) ([]model.KnowledgeCategory, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.agent_id, c.parent_id, c.name, c.description, c.sort_order, c.created_at,
				COALESCE((SELECT COUNT(*) FROM agent_knowledge_items WHERE category_id = c.id), 0)
		 FROM agent_knowledge_categories c 
		 WHERE c.agent_id = $1 ORDER BY c.sort_order`, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.KnowledgeCategory
	for rows.Next() {
		var c model.KnowledgeCategory
		if err := rows.Scan(&c.ID, &c.AgentID, &c.ParentID, &c.Name, &c.Description,
			&c.SortOrder, &c.CreatedAt, &c.ItemCount); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *KnowledgeRepository) CreateCategory(ctx context.Context, c *model.KnowledgeCategory) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO agent_knowledge_categories (agent_id, parent_id, name, description, sort_order)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		c.AgentID, c.ParentID, c.Name, c.Description, c.SortOrder,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *KnowledgeRepository) UpdateCategory(ctx context.Context, c *model.KnowledgeCategory) error {
	_, err := r.db.Exec(ctx,
		`UPDATE agent_knowledge_categories SET name=$1, description=$2, sort_order=$3 WHERE id=$4`,
		c.Name, c.Description, c.SortOrder, c.ID)
	return err
}

func (r *KnowledgeRepository) DeleteCategory(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM agent_knowledge_categories WHERE id = $1", id)
	return err
}

// ========== 知识条目 ==========

func (r *KnowledgeRepository) ListItems(ctx context.Context, agentID int, categoryID *int, offset, limit int) ([]model.KnowledgeItem, int, error) {
	var total int
	query := "SELECT COUNT(*) FROM agent_knowledge_items WHERE agent_id = $1"
	args := []any{agentID}
	if categoryID != nil {
		query += " AND category_id = $2"
		args = append(args, *categoryID)
	}
	if err := r.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	selectQuery := `SELECT id, agent_id, category_id, title, content, tags, source, 
					quality_score, use_count, is_active, created_by, created_at, updated_at 
					FROM agent_knowledge_items WHERE agent_id = $1`
	selectArgs := []any{agentID}
	argIdx := 2
	if categoryID != nil {
		selectQuery += fmt.Sprintf(" AND category_id = $%d", argIdx)
		selectArgs = append(selectArgs, *categoryID)
		argIdx++
	}
	selectQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	selectArgs = append(selectArgs, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.KnowledgeItem
	for rows.Next() {
		var item model.KnowledgeItem
		if err := rows.Scan(&item.ID, &item.AgentID, &item.CategoryID, &item.Title, &item.Content,
			&item.Tags, &item.Source, &item.QualityScore, &item.UseCount, &item.IsActive,
			&item.CreatedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, nil
}

func (r *KnowledgeRepository) CreateItem(ctx context.Context, item *model.KnowledgeItem) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, source, created_by)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		item.AgentID, item.CategoryID, item.Title, item.Content, item.Tags, item.Source, item.CreatedBy,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *KnowledgeRepository) GetItem(ctx context.Context, id int) (*model.KnowledgeItem, error) {
	item := &model.KnowledgeItem{}
	err := r.db.QueryRow(ctx,
		`SELECT id, agent_id, category_id, title, content, tags, source, 
				quality_score, use_count, is_active, created_by, created_at, updated_at 
		 FROM agent_knowledge_items WHERE id = $1`, id,
	).Scan(&item.ID, &item.AgentID, &item.CategoryID, &item.Title, &item.Content,
		&item.Tags, &item.Source, &item.QualityScore, &item.UseCount, &item.IsActive,
		&item.CreatedBy, &item.CreatedAt, &item.UpdatedAt)
	return item, err
}

func (r *KnowledgeRepository) UpdateItem(ctx context.Context, item *model.KnowledgeItem) error {
	_, err := r.db.Exec(ctx,
		`UPDATE agent_knowledge_items SET title=$1, content=$2, tags=$3, is_active=$4, updated_at=NOW() 
		 WHERE id=$5`,
		item.Title, item.Content, item.Tags, item.IsActive, item.ID)
	return err
}

func (r *KnowledgeRepository) DeleteItem(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM agent_knowledge_items WHERE id = $1", id)
	return err
}

func (r *KnowledgeRepository) IncrementUseCount(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx,
		"UPDATE agent_knowledge_items SET use_count = use_count + 1 WHERE id = $1", id)
	return err
}

// fmt is needed in ListItems
func fmt_Sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
