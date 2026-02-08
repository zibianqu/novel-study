package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// ToolCall 工具调用记录
type ToolCall struct {
	ID           int                    `db:"id" json:"id"`
	AgentID      int                    `db:"agent_id" json:"agent_id"`
	ToolName     string                 `db:"tool_name" json:"tool_name"`
	InputParams  map[string]interface{} `db:"input_params" json:"input_params"`
	OutputResult interface{}            `db:"output_result" json:"output_result"`
	Success      bool                   `db:"success" json:"success"`
	ErrorMessage string                 `db:"error_message" json:"error_message"`
	DurationMs   int                    `db:"duration_ms" json:"duration_ms"`
	CreatedAt    time.Time              `db:"created_at" json:"created_at"`
}

// ToolCallRepository 工具调用仓库
type ToolCallRepository struct {
	db *sqlx.DB
}

// NewToolCallRepository 创建工具调用仓库
func NewToolCallRepository(db *sqlx.DB) *ToolCallRepository {
	return &ToolCallRepository{
		db: db,
	}
}

// Create 创建工具调用记录
func (r *ToolCallRepository) Create(ctx context.Context, toolCall *ToolCall) error {
	// 序列化 JSON 字段
	inputJSON, err := json.Marshal(toolCall.InputParams)
	if err != nil {
		return err
	}

	outputJSON, err := json.Marshal(toolCall.OutputResult)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO agent_tool_calls (agent_id, tool_name, input_params, output_result, success, error_message, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err = r.db.QueryRowContext(ctx, query,
		toolCall.AgentID,
		toolCall.ToolName,
		inputJSON,
		outputJSON,
		toolCall.Success,
		toolCall.ErrorMessage,
		toolCall.DurationMs,
	).Scan(&toolCall.ID, &toolCall.CreatedAt)

	return err
}

// GetByAgentID 按 Agent ID 获取工具调用记录
func (r *ToolCallRepository) GetByAgentID(ctx context.Context, agentID int, limit int) ([]*ToolCall, error) {
	query := `
		SELECT id, agent_id, tool_name, input_params, output_result, success, error_message, duration_ms, created_at
		FROM agent_tool_calls
		WHERE agent_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, agentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var toolCalls []*ToolCall
	for rows.Next() {
		var tc ToolCall
		var inputJSON, outputJSON []byte

		err := rows.Scan(
			&tc.ID,
			&tc.AgentID,
			&tc.ToolName,
			&inputJSON,
			&outputJSON,
			&tc.Success,
			&tc.ErrorMessage,
			&tc.DurationMs,
			&tc.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 反序列化 JSON
		if err := json.Unmarshal(inputJSON, &tc.InputParams); err != nil {
			tc.InputParams = make(map[string]interface{})
		}
		if err := json.Unmarshal(outputJSON, &tc.OutputResult); err != nil {
			tc.OutputResult = nil
		}

		toolCalls = append(toolCalls, &tc)
	}

	return toolCalls, rows.Err()
}

// GetStatistics 获取工具调用统计
func (r *ToolCallRepository) GetStatistics(ctx context.Context, agentID int, startTime, endTime time.Time) (map[string]interface{}, error) {
	query := `
		SELECT 
			tool_name,
			COUNT(*) as call_count,
			COUNT(CASE WHEN success THEN 1 END) as success_count,
			AVG(duration_ms) as avg_duration_ms,
			MAX(duration_ms) as max_duration_ms,
			MIN(duration_ms) as min_duration_ms
		FROM agent_tool_calls
		WHERE agent_id = $1 AND created_at BETWEEN $2 AND $3
		GROUP BY tool_name
	`

	rows, err := r.db.QueryContext(ctx, query, agentID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statistics := make(map[string]interface{})
	tools := []map[string]interface{}{}

	for rows.Next() {
		var toolName string
		var callCount, successCount int
		var avgDuration, maxDuration, minDuration sql.NullFloat64

		err := rows.Scan(&toolName, &callCount, &successCount, &avgDuration, &maxDuration, &minDuration)
		if err != nil {
			return nil, err
		}

		toolStat := map[string]interface{}{
			"tool_name":       toolName,
			"call_count":      callCount,
			"success_count":   successCount,
			"success_rate":    float64(successCount) / float64(callCount),
			"avg_duration_ms": avgDuration.Float64,
			"max_duration_ms": maxDuration.Float64,
			"min_duration_ms": minDuration.Float64,
		}
		tools = append(tools, toolStat)
	}

	statistics["tools"] = tools
	statistics["agent_id"] = agentID
	statistics["period"] = map[string]interface{}{
		"start": startTime,
		"end":   endTime,
	}

	return statistics, rows.Err()
}
