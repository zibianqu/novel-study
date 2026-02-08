package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// DBToolCallLogger 数据库工具调用日志记录器
type DBToolCallLogger struct {
	db *sql.DB
}

// NewDBToolCallLogger 创建数据库日志记录器
func NewDBToolCallLogger(db *sql.DB) *DBToolCallLogger {
	return &DBToolCallLogger{db: db}
}

// Log 记录工具调用日志
func (l *DBToolCallLogger) Log(agentID int, toolName string, params map[string]interface{}, result interface{}, err error, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	success := err == nil
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error()
	}

	// 序列化参数和结果
	paramsJSON, _ := json.Marshal(params)
	resultJSON, _ := json.Marshal(result)

	durationMs := duration.Milliseconds()

	query := `
		INSERT INTO agent_tool_calls 
		(agent_id, tool_name, input_params, output_result, success, error_message, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, insertErr := l.db.ExecContext(ctx, query, agentID, toolName, paramsJSON, resultJSON, success, errorMessage, durationMs)
	if insertErr != nil {
		fmt.Printf("⚠️  工具调用日志保存失败: %v\n", insertErr)
	}
}

// GetToolCallStats 获取工具调用统计
func (l *DBToolCallLogger) GetToolCallStats(ctx context.Context, agentID int, since time.Time) (*ToolCallStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_calls,
			COUNT(*) FILTER (WHERE success = true) as success_calls,
			COUNT(*) FILTER (WHERE success = false) as failed_calls,
			AVG(duration_ms) as avg_duration_ms,
			MAX(duration_ms) as max_duration_ms
		FROM agent_tool_calls
		WHERE agent_id = $1 AND created_at >= $2
	`

	var stats ToolCallStats
	err := l.db.QueryRowContext(ctx, query, agentID, since).Scan(
		&stats.TotalCalls,
		&stats.SuccessCalls,
		&stats.FailedCalls,
		&stats.AvgDurationMs,
		&stats.MaxDurationMs,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get tool call stats: %w", err)
	}

	return &stats, nil
}

// GetRecentToolCalls 获取最近的工具调用记录
func (l *DBToolCallLogger) GetRecentToolCalls(ctx context.Context, agentID int, limit int) ([]ToolCallRecord, error) {
	query := `
		SELECT id, tool_name, input_params, output_result, success, error_message, duration_ms, created_at
		FROM agent_tool_calls
		WHERE agent_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := l.db.QueryContext(ctx, query, agentID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query tool calls: %w", err)
	}
	defer rows.Close()

	var records []ToolCallRecord
	for rows.Next() {
		var r ToolCallRecord
		var paramsJSON, resultJSON []byte

		err := rows.Scan(
			&r.ID,
			&r.ToolName,
			&paramsJSON,
			&resultJSON,
			&r.Success,
			&r.ErrorMessage,
			&r.DurationMs,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tool call: %w", err)
		}

		// 反序列化 JSON
		if len(paramsJSON) > 0 {
			json.Unmarshal(paramsJSON, &r.InputParams)
		}
		if len(resultJSON) > 0 {
			json.Unmarshal(resultJSON, &r.OutputResult)
		}

		records = append(records, r)
	}

	return records, nil
}

// ToolCallStats 工具调用统计
type ToolCallStats struct {
	TotalCalls    int     `json:"total_calls"`
	SuccessCalls  int     `json:"success_calls"`
	FailedCalls   int     `json:"failed_calls"`
	AvgDurationMs float64 `json:"avg_duration_ms"`
	MaxDurationMs int     `json:"max_duration_ms"`
}

// ToolCallRecord 工具调用记录
type ToolCallRecord struct {
	ID           int                    `json:"id"`
	ToolName     string                 `json:"tool_name"`
	InputParams  map[string]interface{} `json:"input_params"`
	OutputResult interface{}            `json:"output_result"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	DurationMs   int                    `json:"duration_ms"`
	CreatedAt    time.Time              `json:"created_at"`
}
