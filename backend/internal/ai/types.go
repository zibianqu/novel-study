package ai

import (
	"context"
)

// AgentConfig Agent配置
type AgentConfig struct {
	AgentKey     string
	Name         string
	SystemPrompt string
	Model        string
	Temperature  float64
	MaxTokens    int
	Tools        []string
}

// AgentRequest Agent请求
type AgentRequest struct {
	UserID      int
	ProjectID   int
	ChapterID   *int
	Prompt      string
	Context     map[string]interface{}
	Stream      bool
}

// AgentResponse Agent响应
type AgentResponse struct {
	Content      string
	TokensUsed   int
	DurationMs   int64
	Metadata     map[string]interface{}
	Error        string
}

// Agent 接口
type Agent interface {
	GetConfig() *AgentConfig
	Execute(ctx context.Context, req *AgentRequest) (*AgentResponse, error)
	ExecuteStream(ctx context.Context, req *AgentRequest, callback func(string)) error
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// QualityScore 质量评分
type QualityScore struct {
	OverallScore int                    `json:"overall_score"`
	Passed       bool                   `json:"passed"`
	Dimensions   map[string]interface{} `json:"dimensions"`
	Feedback     map[string]string      `json:"feedback"`
}
