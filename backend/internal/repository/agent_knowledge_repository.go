package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// AgentKnowledgeCategory Agent知识库分类
type AgentKnowledgeCategory struct {
	ID          int       `json:"id"`
	AgentID     int       `json:"agent_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AgentKnowledgeItem Agent知识条目
type AgentKnowledgeItem struct {
	ID           int       `json:"id"`
	AgentID      int       `json:"agent_id"`
	CategoryID   int       `json:"category_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	QualityScore int       `json:"quality_score"`
	UsageCount   int       `json:"usage_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AgentKnowledgeRepository Agent知识库Repository
type AgentKnowledgeRepository struct {
	db *sql.DB
}

// NewAgentKnowledgeRepository 创建Repository
func NewAgentKnowledgeRepository(db *sql.DB) *AgentKnowledgeRepository {
	return &AgentKnowledgeRepository{db: db}
}

// CreateCategory 创建分类
func (r *AgentKnowledgeRepository) CreateCategory(ctx context.Context, category *AgentKnowledgeCategory) error {
	query := `
		INSERT INTO agent_knowledge_categories (agent_id, name, description, priority)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		category.AgentID,
		category.Name,
		category.Description,
		category.Priority,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

// GetCategoriesByAgentID 获取Agent的所有分类
func (r *AgentKnowledgeRepository) GetCategoriesByAgentID(ctx context.Context, agentID int) ([]*AgentKnowledgeCategory, error) {
	query := `
		SELECT id, agent_id, name, description, priority, created_at, updated_at
		FROM agent_knowledge_categories
		WHERE agent_id = $1
		ORDER BY priority DESC, id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*AgentKnowledgeCategory
	for rows.Next() {
		var cat AgentKnowledgeCategory
		err := rows.Scan(
			&cat.ID,
			&cat.AgentID,
			&cat.Name,
			&cat.Description,
			&cat.Priority,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}

	return categories, rows.Err()
}

// GetCategoryByID 根据ID获取分类
func (r *AgentKnowledgeRepository) GetCategoryByID(ctx context.Context, id int) (*AgentKnowledgeCategory, error) {
	query := `
		SELECT id, agent_id, name, description, priority, created_at, updated_at
		FROM agent_knowledge_categories
		WHERE id = $1
	`

	var cat AgentKnowledgeCategory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cat.ID,
		&cat.AgentID,
		&cat.Name,
		&cat.Description,
		&cat.Priority,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}

	return &cat, err
}

// CreateKnowledgeItem 创建知识条目
func (r *AgentKnowledgeRepository) CreateKnowledgeItem(ctx context.Context, item *AgentKnowledgeItem) error {
	query := `
		INSERT INTO agent_knowledge_items 
		(agent_id, category_id, title, content, tags, quality_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		item.AgentID,
		item.CategoryID,
		item.Title,
		item.Content,
		pq.Array(item.Tags),
		item.QualityScore,
	).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

// GetKnowledgeItemsByCategory 获取分类下的所有知识
func (r *AgentKnowledgeRepository) GetKnowledgeItemsByCategory(ctx context.Context, categoryID int, limit int) ([]*AgentKnowledgeItem, error) {
	query := `
		SELECT id, agent_id, category_id, title, content, tags, quality_score, usage_count, created_at, updated_at
		FROM agent_knowledge_items
		WHERE category_id = $1
		ORDER BY quality_score DESC, usage_count DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanKnowledgeItems(rows)
}

// GetKnowledgeItemsByAgentID 获取Agent的所有知识
func (r *AgentKnowledgeRepository) GetKnowledgeItemsByAgentID(ctx context.Context, agentID int, categoryID *int, limit int) ([]*AgentKnowledgeItem, error) {
	var query string
	var args []interface{}

	if categoryID != nil {
		query = `
			SELECT id, agent_id, category_id, title, content, tags, quality_score, usage_count, created_at, updated_at
			FROM agent_knowledge_items
			WHERE agent_id = $1 AND category_id = $2
			ORDER BY quality_score DESC, usage_count DESC
			LIMIT $3
		`
		args = []interface{}{agentID, *categoryID, limit}
	} else {
		query = `
			SELECT id, agent_id, category_id, title, content, tags, quality_score, usage_count, created_at, updated_at
			FROM agent_knowledge_items
			WHERE agent_id = $1
			ORDER BY quality_score DESC, usage_count DESC
			LIMIT $2
		`
		args = []interface{}{agentID, limit}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanKnowledgeItems(rows)
}

// SearchKnowledgeItems 搜索知识条目
func (r *AgentKnowledgeRepository) SearchKnowledgeItems(ctx context.Context, agentID int, keyword string, limit int) ([]*AgentKnowledgeItem, error) {
	query := `
		SELECT id, agent_id, category_id, title, content, tags, quality_score, usage_count, created_at, updated_at
		FROM agent_knowledge_items
		WHERE agent_id = $1 AND (title ILIKE $2 OR content ILIKE $2 OR $2 = ANY(tags))
		ORDER BY quality_score DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, agentID, "%"+keyword+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanKnowledgeItems(rows)
}

// IncrementUsageCount 增加使用次数
func (r *AgentKnowledgeRepository) IncrementUsageCount(ctx context.Context, itemID int) error {
	query := `
		UPDATE agent_knowledge_items
		SET usage_count = usage_count + 1, updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}

// UpdateKnowledgeItem 更新知识条目
func (r *AgentKnowledgeRepository) UpdateKnowledgeItem(ctx context.Context, item *AgentKnowledgeItem) error {
	query := `
		UPDATE agent_knowledge_items
		SET title = $1, content = $2, tags = $3, quality_score = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, query,
		item.Title,
		item.Content,
		pq.Array(item.Tags),
		item.QualityScore,
		item.ID,
	)
	return err
}

// DeleteKnowledgeItem 删除知识条目
func (r *AgentKnowledgeRepository) DeleteKnowledgeItem(ctx context.Context, itemID int) error {
	query := `DELETE FROM agent_knowledge_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}

// GetKnowledgeStats 获取知识库统计
func (r *AgentKnowledgeRepository) GetKnowledgeStats(ctx context.Context, agentID int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT category_id) as category_count,
			COUNT(*) as total_items,
			AVG(quality_score) as avg_quality_score,
			SUM(usage_count) as total_usage
		FROM agent_knowledge_items
		WHERE agent_id = $1
	`

	var categoryCount, totalItems, totalUsage int
	var avgQualityScore float64

	err := r.db.QueryRowContext(ctx, query, agentID).Scan(
		&categoryCount,
		&totalItems,
		&avgQualityScore,
		&totalUsage,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"category_count":     categoryCount,
		"total_items":        totalItems,
		"avg_quality_score":  avgQualityScore,
		"total_usage":        totalUsage,
	}, nil
}

// scanKnowledgeItems 扫描知识条目结果
func (r *AgentKnowledgeRepository) scanKnowledgeItems(rows *sql.Rows) ([]*AgentKnowledgeItem, error) {
	var items []*AgentKnowledgeItem

	for rows.Next() {
		var item AgentKnowledgeItem
		err := rows.Scan(
			&item.ID,
			&item.AgentID,
			&item.CategoryID,
			&item.Title,
			&item.Content,
			pq.Array(&item.Tags),
			&item.QualityScore,
			&item.UsageCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, rows.Err()
}
