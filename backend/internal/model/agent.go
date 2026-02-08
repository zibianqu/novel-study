package model

import (
	"time"
)

type Agent struct {
	ID           int       `json:"id"`
	UserID       *int      `json:"user_id"`
	AgentKey     string    `json:"agent_key"`
	Name         string    `json:"name"`
	Icon         string    `json:"icon"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`  // core, extension
	Layer        string    `json:"layer"` // decision, strategy, execution, quality, auxiliary
	SystemPrompt string    `json:"system_prompt"`
	Model        string    `json:"model"`
	Temperature  float64   `json:"temperature"`
	MaxTokens    int       `json:"max_tokens"`
	Tools        string    `json:"tools"` // JSON array
	IsActive     bool      `json:"is_active"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AIInteractionLog struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ProjectID      *int      `json:"project_id"`
	AgentID        *int      `json:"agent_id"`
	ActionType     string    `json:"action_type"`
	InputPrompt    string    `json:"input_prompt"`
	OutputResponse string    `json:"output_response"`
	TokensInput    int       `json:"tokens_input"`
	TokensOutput   int       `json:"tokens_output"`
	Model          string    `json:"model"`
	DurationMs     int       `json:"duration_ms"`
	CreatedAt      time.Time `json:"created_at"`
}
