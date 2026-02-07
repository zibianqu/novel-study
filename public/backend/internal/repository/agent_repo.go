package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/novelforge/backend/internal/model"
)

type AgentRepository struct {
	db *pgxpool.Pool
}

func NewAgentRepository(db *pgxpool.Pool) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) Create(ctx context.Context, a *model.Agent) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO agents (user_id, agent_key, name, icon, description, type, layer, 
			system_prompt, model, temperature, max_tokens, tools, input_schema, output_schema, permissions) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) 
		 RETURNING id, created_at, updated_at`,
		a.UserID, a.AgentKey, a.Name, a.Icon, a.Description, a.Type, a.Layer,
		a.SystemPrompt, a.Model, a.Temperature, a.MaxTokens,
		a.Tools, a.InputSchema, a.OutputSchema, a.Permissions,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AgentRepository) GetByID(ctx context.Context, id int) (*model.Agent, error) {
	a := &model.Agent{}
	err := r.db.QueryRow(ctx,
		`SELECT a.id, a.user_id, a.agent_key, a.name, a.icon, a.description, a.type, a.layer,
				a.system_prompt, a.model, a.temperature, a.max_tokens, a.tools, 
				a.input_schema, a.output_schema, a.permissions, a.is_active, a.sort_order,
				a.created_at, a.updated_at,
				COALESCE((SELECT COUNT(*) FROM agent_knowledge_items WHERE agent_id = a.id), 0)
		 FROM agents a WHERE a.id = $1`, id,
	).Scan(&a.ID, &a.UserID, &a.AgentKey, &a.Name, &a.Icon, &a.Description, &a.Type, &a.Layer,
		&a.SystemPrompt, &a.Model, &a.Temperature, &a.MaxTokens, &a.Tools,
		&a.InputSchema, &a.OutputSchema, &a.Permissions, &a.IsActive, &a.SortOrder,
		&a.CreatedAt, &a.UpdatedAt, &a.KnowledgeCount)
	if err != nil {
		return nil, fmt.Errorf("Agent不存在: %w", err)
	}
	return a, nil
}

func (r *AgentRepository) GetByKey(ctx context.Context, key string) (*model.Agent, error) {
	a := &model.Agent{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, agent_key, name, icon, description, type, layer,
				system_prompt, model, temperature, max_tokens, tools, 
				input_schema, output_schema, permissions, is_active, sort_order,
				created_at, updated_at
		 FROM agents WHERE agent_key = $1`, key,
	).Scan(&a.ID, &a.UserID, &a.AgentKey, &a.Name, &a.Icon, &a.Description, &a.Type, &a.Layer,
		&a.SystemPrompt, &a.Model, &a.Temperature, &a.MaxTokens, &a.Tools,
		&a.InputSchema, &a.OutputSchema, &a.Permissions, &a.IsActive, &a.SortOrder,
		&a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Agent不存在: %w", err)
	}
	return a, nil
}

func (r *AgentRepository) List(ctx context.Context) ([]model.Agent, error) {
	rows, err := r.db.Query(ctx,
		`SELECT a.id, a.user_id, a.agent_key, a.name, a.icon, a.description, a.type, a.layer,
				a.system_prompt, a.model, a.temperature, a.max_tokens, a.is_active, a.sort_order,
				a.created_at, a.updated_at,
				COALESCE((SELECT COUNT(*) FROM agent_knowledge_items WHERE agent_id = a.id), 0)
		 FROM agents a ORDER BY a.sort_order, a.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []model.Agent
	for rows.Next() {
		var a model.Agent
		if err := rows.Scan(&a.ID, &a.UserID, &a.AgentKey, &a.Name, &a.Icon, &a.Description,
			&a.Type, &a.Layer, &a.SystemPrompt, &a.Model, &a.Temperature, &a.MaxTokens,
			&a.IsActive, &a.SortOrder, &a.CreatedAt, &a.UpdatedAt, &a.KnowledgeCount); err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, nil
}

func (r *AgentRepository) Update(ctx context.Context, a *model.Agent) error {
	_, err := r.db.Exec(ctx,
		`UPDATE agents SET name=$1, icon=$2, description=$3, system_prompt=$4, 
			model=$5, temperature=$6, max_tokens=$7, is_active=$8, updated_at=NOW() 
		 WHERE id=$9`,
		a.Name, a.Icon, a.Description, a.SystemPrompt,
		a.Model, a.Temperature, a.MaxTokens, a.IsActive, a.ID)
	return err
}

func (r *AgentRepository) Delete(ctx context.Context, id int) error {
	// 只允许删除扩展Agent
	result, err := r.db.Exec(ctx,
		"DELETE FROM agents WHERE id = $1 AND type = 'extension'", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("无法删除核心Agent或Agent不存在")
	}
	return nil
}
