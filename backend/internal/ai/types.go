package ai

import "context"

// Agent 接口定义
type Agent interface {
	Execute(ctx context.Context, req *AgentRequest) (*AgentResponse, error)
	ExecuteStream(ctx context.Context, req *AgentRequest, callback func(string)) error
	GetName() string
	GetDescription() string
}

// AgentRequest Agent请求
type AgentRequest struct {
	Prompt      string                 `json:"prompt"`
	Context     string                 `json:"context"`
	ProjectID   int                    `json:"project_id"`
	Metadata    map[string]interface{} `json:"metadata"`
	MaxTokens   int                    `json:"max_tokens"`
	Temperature float64                `json:"temperature"`
}

// AgentResponse Agent响应
type AgentResponse struct {
	Content    string                 `json:"content"`
	TokensUsed int                    `json:"tokens_used"`
	DurationMs int64                  `json:"duration_ms"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"`
}
