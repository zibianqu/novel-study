package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/repository"
)

// ToolRegistry 工具注册表（完善版）
type ToolRegistry struct {
	tools map[string]Tool
	logger ToolCallLogger // 日志记录器
}

// ToolCallLogger 工具调用日志记录器接口
type ToolCallLogger interface {
	LogToolCall(ctx context.Context, agentID int, toolName string, input map[string]interface{}, output interface{}, success bool, err error, duration time.Duration) error
}

// NewToolRegistry 创建工具注册表
func NewToolRegistry(logger ToolCallLogger) *ToolRegistry {
	return &ToolRegistry{
		tools:  make(map[string]Tool),
		logger: logger,
	}
}

// Register 注册工具
func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.GetName()] = tool
}

// Get 获取工具
func (r *ToolRegistry) Get(name string) (Tool, error) {
	tool, ok := r.tools[name]
	if !ok {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	return tool, nil
}

// Execute 执行工具并记录日志
func (r *ToolRegistry) Execute(ctx context.Context, agentID int, toolName string, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 获取工具
	tool, err := r.Get(toolName)
	if err != nil {
		return nil, err
	}

	// 执行工具
	result, err := tool.Execute(ctx, params)
	duration := time.Since(start)

	// 记录日志
	if r.logger != nil {
		logErr := r.logger.LogToolCall(ctx, agentID, toolName, params, result, err == nil, err, duration)
		if logErr != nil {
			// 记录失败不影响工具执行，但要记录一下
			fmt.Printf("Failed to log tool call: %v\n", logErr)
		}
	}

	return result, err
}

// ListTools 列出所有工具
func (r *ToolRegistry) ListTools() []string {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// GetToolDescriptions 获取所有工具的描述
func (r *ToolRegistry) GetToolDescriptions() map[string]string {
	descriptions := make(map[string]string)
	for name, tool := range r.tools {
		descriptions[name] = tool.GetDescription()
	}
	return descriptions
}

// ToolCallLoggerImpl 工具调用日志记录器实现
type ToolCallLoggerImpl struct {
	repo *repository.ToolCallRepository
}

// NewToolCallLogger 创建日志记录器
func NewToolCallLogger(repo *repository.ToolCallRepository) *ToolCallLoggerImpl {
	return &ToolCallLoggerImpl{
		repo: repo,
	}
}

func (l *ToolCallLoggerImpl) LogToolCall(ctx context.Context, agentID int, toolName string, input map[string]interface{}, output interface{}, success bool, err error, duration time.Duration) error {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	return l.repo.Create(ctx, &repository.ToolCall{
		AgentID:      agentID,
		ToolName:     toolName,
		InputParams:  input,
		OutputResult: output,
		Success:      success,
		ErrorMessage: errorMsg,
		DurationMs:   int(duration.Milliseconds()),
	})
}
