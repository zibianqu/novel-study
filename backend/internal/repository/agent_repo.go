package repository

import (
	"database/sql"
	"github.com/zibianqu/novel-study/internal/model"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(db *sql.DB) *AgentRepository {
	return &AgentRepository{db: db}
}

func (r *AgentRepository) GetCoreAgents() ([]*model.Agent, error) {
	query := `
		SELECT id, agent_key, name, icon, description, type, layer,
		       system_prompt, model, temperature, max_tokens, tools,
		       is_active, sort_order, created_at, updated_at
		FROM agents WHERE type = 'core' AND is_active = true
		ORDER BY sort_order ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*model.Agent
	for rows.Next() {
		agent := &model.Agent{}
		err := rows.Scan(
			&agent.ID,
			&agent.AgentKey,
			&agent.Name,
			&agent.Icon,
			&agent.Description,
			&agent.Type,
			&agent.Layer,
			&agent.SystemPrompt,
			&agent.Model,
			&agent.Temperature,
			&agent.MaxTokens,
			&agent.Tools,
			&agent.IsActive,
			&agent.SortOrder,
			&agent.CreatedAt,
			&agent.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

func (r *AgentRepository) LogInteraction(log *model.AIInteractionLog) error {
	query := `
		INSERT INTO ai_interaction_logs 
		(user_id, project_id, agent_id, action_type, input_prompt, output_response,
		 tokens_input, tokens_output, model, duration_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		log.UserID,
		log.ProjectID,
		log.AgentID,
		log.ActionType,
		log.InputPrompt,
		log.OutputResponse,
		log.TokensInput,
		log.TokensOutput,
		log.Model,
		log.DurationMs,
	).Scan(&log.ID, &log.CreatedAt)
}
