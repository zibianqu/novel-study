package tools

import (
	"context"
	"fmt"
	"time"
)

// DispatchAgentTool 调度其他Agent的工具
type DispatchAgentTool struct {
	// 这里需要依赖 AgentExecutor 或者 Scheduler
	// 为了避免循环依赖，我们使用接口
	executor AgentExecutor
}

// AgentExecutor Agent执行器接口
type AgentExecutor interface {
	ExecuteAgent(ctx context.Context, agentKey string, input map[string]interface{}) (map[string]interface{}, error)
}

// NewDispatchAgentTool 创建调度Agent工具
func NewDispatchAgentTool(executor AgentExecutor) *DispatchAgentTool {
	return &DispatchAgentTool{
		executor: executor,
	}
}

func (t *DispatchAgentTool) GetName() string {
	return "dispatch_agent"
}

func (t *DispatchAgentTool) GetDescription() string {
	return "调度指定Agent执行任务。参数：agent_key(Agent标识), input(输入数据)。仅总导演Agent可用。"
}

func (t *DispatchAgentTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	agentKey, ok := params["agent_key"].(string)
	if !ok || agentKey == "" {
		return nil, fmt.Errorf("missing or invalid 'agent_key' parameter")
	}

	// 获取输入数据
	input := make(map[string]interface{})
	if inp, ok := params["input"].(map[string]interface{}); ok {
		input = inp
	}

	// 执行Agent
	result, err := t.executor.ExecuteAgent(ctx, agentKey, input)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	// 返回结果
	response := map[string]interface{}{
		"agent_key":   agentKey,
		"result":      result,
		"duration_ms": time.Since(start).Milliseconds(),
	}

	return response, nil
}
