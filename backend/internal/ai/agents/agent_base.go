package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

// GetName 获取Agent名称
func (a *BaseAgent) GetName() string {
	return a.config.Name
}

// GetDescription 获取Agent描述
func (a *BaseAgent) GetDescription() string {
	return a.config.Description
}

// Execute 执行Agent
func (a *BaseAgent) Execute(ctx context.Context, req *ai.AgentRequest) (*ai.AgentResponse, error) {
	start := time.Now()

	// 参数验证
	if req.Prompt == "" {
		return nil, fmt.Errorf("prompt cannot be empty")
	}

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

	log.Printf("[%s] Executing request: %s", a.config.Name, req.Prompt)

	// 调用 OpenAI API
	content, tokensUsed, err := a.callOpenAI(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	duration := time.Since(start).Milliseconds()

	return &ai.AgentResponse{
		Content:    content,
		TokensUsed: tokensUsed,
		DurationMs: duration,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// ExecuteStream 流式执行
func (a *BaseAgent) ExecuteStream(ctx context.Context, req *ai.AgentRequest, callback func(string)) error {
	// 参数验证
	if req.Prompt == "" {
		return fmt.Errorf("prompt cannot be empty")
	}

	// 构建消息
	messages := []ai.ChatMessage{
		{Role: "system", Content: a.config.SystemPrompt},
		{Role: "user", Content: req.Prompt},
	}

	if req.Context != nil && len(req.Context) > 0 {
		contextJSON, _ := json.Marshal(req.Context)
		contextMsg := fmt.Sprintf("\n\n上下文信息: %s", string(contextJSON))
		messages[1].Content += contextMsg
	}

	log.Printf("[%s] Executing stream request: %s", a.config.Name, req.Prompt)

	// 调用流式 API
	return a.callOpenAIStream(ctx, messages, callback)
}

// callOpenAI 调用OpenAI API
func (a *BaseAgent) callOpenAI(ctx context.Context, messages []ai.ChatMessage) (string, int, error) {
	// TODO: 实际集成 OpenAI API
	// 这里先返回模拟响应
	if len(messages) > 0 {
		lastMsg := messages[len(messages)-1]
		content := fmt.Sprintf("%s 处理结果: %s", a.config.Name, lastMsg.Content)
		return content, 100, nil // 模拟 100 tokens
	}
	return "模拟响应", 50, nil
}

// callOpenAIStream 流式调用OpenAI API
func (a *BaseAgent) callOpenAIStream(ctx context.Context, messages []ai.ChatMessage, callback func(string)) error {
	// TODO: 实现流式输出
	// 临时实现：模拟流式输出
	content, _, err := a.callOpenAI(ctx, messages)
	if err != nil {
		return err
	}

	// 模拟分段输出
	words := []rune(content)
	chunkSize := 10
	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		callback(string(words[i:end]))
		time.Sleep(50 * time.Millisecond) // 模拟延迟
	}

	return nil
}
