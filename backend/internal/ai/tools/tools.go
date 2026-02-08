package tools

import (
	"context"
	"fmt"
	"time"
)

// Tool Agent工具接口
type Tool interface {
	GetName() string
	GetDescription() string
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)
}

// ToolRegistry 工具注册表
type ToolRegistry struct {
	tools map[string]Tool
	logger ToolCallLogger
}

// ToolCallLogger 工具调用日志接口
type ToolCallLogger interface {
	Log(agentID int, toolName string, params map[string]interface{}, result interface{}, err error, duration time.Duration)
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
	fmt.Printf("✅ 工具注册成功: %s - %s\n", tool.GetName(), tool.GetDescription())
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
	tool, err := r.Get(toolName)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	result, execErr := tool.Execute(ctx, params)
	duration := time.Since(startTime)

	// 记录日志
	if r.logger != nil {
		r.logger.Log(agentID, toolName, params, result, execErr, duration)
	}

	return result, execErr
}

// ListTools 列出所有可用工具
func (r *ToolRegistry) ListTools() []ToolInfo {
	tools := make([]ToolInfo, 0, len(r.tools))
	for name, tool := range r.tools {
		tools = append(tools, ToolInfo{
			Name:        name,
			Description: tool.GetDescription(),
		})
	}
	return tools
}

// ToolInfo 工具信息
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToolCallResult 工具调用结果
type ToolCallResult struct {
	ToolName   string                 `json:"tool_name"`
	Success    bool                   `json:"success"`
	Result     interface{}            `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	DurationMs int64                  `json:"duration_ms"`
	Params     map[string]interface{} `json:"params,omitempty"`
}
