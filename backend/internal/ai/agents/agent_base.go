package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"novel-study/backend/internal/ai"
	"novel-study/backend/internal/ai/tools"
)

// BaseAgent 基础Agent实现
type BaseAgent struct {
	config       *ai.AgentConfig
	apiKey       string
	toolRegistry *tools.ToolRegistry
	agentID      int // 用于工具调用日志
}

// NewBaseAgent 创建基础Agent
func NewBaseAgent(config *ai.AgentConfig, apiKey string, toolRegistry *tools.ToolRegistry, agentID int) *BaseAgent {
	return &BaseAgent{
		config:       config,
		apiKey:       apiKey,
		toolRegistry: toolRegistry,
		agentID:      agentID,
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

// CanUseTool 检查是否可以使用指定工具
func (a *BaseAgent) CanUseTool(toolName string) bool {
	for _, t := range a.config.Tools {
		if t == toolName {
			return true
		}
	}
	return false
}

// CallTool 调用工具
func (a *BaseAgent) CallTool(ctx context.Context, toolName string, params map[string]interface{}) (interface{}, error) {
	if !a.CanUseTool(toolName) {
		return nil, fmt.Errorf("agent %s is not authorized to use tool: %s", a.config.Name, toolName)
	}

	if a.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not initialized")
	}

	log.Printf("[%s] 调用工具: %s", a.config.Name, toolName)
	return a.toolRegistry.Execute(ctx, a.agentID, toolName, params)
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

	// 添加工具信息
	if a.toolRegistry != nil && len(a.config.Tools) > 0 {
		toolsInfo := a.buildToolsInfo()
		messages[0].Content += toolsInfo
	}

	log.Printf("[%s] Executing request: %s", a.config.Name, req.Prompt)

	// ✨ 调用 OpenAI API (带重试)
	content, tokensUsed, err := a.callOpenAIWithRetry(ctx, messages, 3)
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

// buildToolsInfo 构建工具信息
func (a *BaseAgent) buildToolsInfo() string {
	if len(a.config.Tools) == 0 {
		return ""
	}

	var toolsDesc string
	toolsDesc = "\n\n你可以使用以下工具：\n"

	for _, toolName := range a.config.Tools {
		tool, err := a.toolRegistry.Get(toolName)
		if err == nil {
			toolsDesc += fmt.Sprintf("- %s: %s\n", toolName, tool.GetDescription())
		}
	}

	toolsDesc += "\n如果需要使用工具，请在响应中明确说明需要调用的工具和参数。"

	return toolsDesc
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

	// 添加工具信息
	if a.toolRegistry != nil && len(a.config.Tools) > 0 {
		toolsInfo := a.buildToolsInfo()
		messages[0].Content += toolsInfo
	}

	log.Printf("[%s] Executing stream request: %s", a.config.Name, req.Prompt)

	// 调用流式 API
	return a.callOpenAIStream(ctx, messages, callback)
}

// callOpenAIWithRetry 带重试的 OpenAI API 调用
func (a *BaseAgent) callOpenAIWithRetry(ctx context.Context, messages []ai.ChatMessage, maxRetries int) (string, int, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		content, tokensUsed, err := a.callOpenAI(ctx, messages)
		if err == nil {
			return content, tokensUsed, nil
		}

		lastErr = err
		log.Printf("[%s] API call failed (attempt %d/%d): %v", a.config.Name, i+1, maxRetries, err)

		// 最后一次失败不等待
		if i < maxRetries-1 {
			// 指数退避: 1s, 2s, 4s
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			log.Printf("[%s] Retrying in %v...", a.config.Name, waitTime)

			select {
			case <-time.After(waitTime):
				// 继续重试
			case <-ctx.Done():
				// Context 取消
				return "", 0, ctx.Err()
			}
		}
	}

	return "", 0, fmt.Errorf("all retries failed: %w", lastErr)
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
		// 检查 context 是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		callback(string(words[i:end]))
		time.Sleep(50 * time.Millisecond) // 模拟延迟
	}

	return nil
}
