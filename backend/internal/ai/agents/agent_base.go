package agents

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zibianqu/novel-study/internal/ai"
)

// BaseAgent 基础Agent实现
type BaseAgent struct {
	config *ai.AgentConfig
	apiKey string
}

// NewBaseAgent 创建基础Agent
func NewBaseAgent(config *ai.AgentConfig, apiKey string) *BaseAgent {
	return &BaseAgent{
		config: config,
		apiKey: apiKey,
	}
}

// GetConfig 获取配置
func (a *BaseAgent) GetConfig() *ai.AgentConfig {
	return a.config
}

// Execute 执行Agent
func (a *BaseAgent) Execute(ctx context.Context, req *ai.AgentRequest) (*ai.AgentResponse, error) {
	// 构建消息
	messages := []ai.ChatMessage{
		{Role: "system", Content: a.config.SystemPrompt},
		{Role: "user", Content: req.Prompt},
	}

	// 添加上下文信息
	if req.Context != nil && len(req.Context) > 0 {
		contextJSON, _ := json.Marshal(req.Context)
		contextMsg := fmt.Sprintf("\n\n上下文信息: %s", string(contextJSON))
		messages[1].Content += contextMsg
	}

	// 调用 OpenAI API
	content, err := a.callOpenAI(ctx, messages)
	if err != nil {
		return nil, err
	}

	return &ai.AgentResponse{
		Content:    content,
		TokensUsed: 0, // TODO: 计算token消耗
		Metadata:   make(map[string]interface{}),
	}, nil
}

// ExecuteStream 流式执行
func (a *BaseAgent) ExecuteStream(ctx context.Context, req *ai.AgentRequest, callback func(string)) error {
	// TODO: 实现流式输出
	resp, err := a.Execute(ctx, req)
	if err != nil {
		return err
	}
	callback(resp.Content)
	return nil
}

// callOpenAI 调用OpenAI API
func (a *BaseAgent) callOpenAI(ctx context.Context, messages []ai.ChatMessage) (string, error) {
	// TODO: 实际集成 OpenAI API
	// 这里先返回模拟响应
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		return fmt.Sprintf("%s 处理结果: %s", a.config.Name, lastMsg.Content), nil
	}
	return "模拟响应", nil
}
