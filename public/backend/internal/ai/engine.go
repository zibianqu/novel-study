package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"github.com/novelforge/backend/internal/config"
)

// Engine AI引擎 - 封装OpenAI调用
type Engine struct {
	client   *openai.Client
	cfg      config.OpenAIConfig
	aiCfg    config.AIConfig
}

// NewEngine 创建AI引擎
func NewEngine(cfg config.OpenAIConfig, aiCfg config.AIConfig) *Engine {
	clientCfg := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		clientCfg.BaseURL = cfg.BaseURL
	}

	return &Engine{
		client: openai.NewClientWithConfig(clientCfg),
		cfg:    cfg,
		aiCfg:  aiCfg,
	}
}

// GetClient 获取原始OpenAI客户端
func (e *Engine) GetClient() *openai.Client {
	return e.client
}

// GetConfig 获取AI配置
func (e *Engine) GetConfig() config.AIConfig {
	return e.aiCfg
}

// ==================== 同步调用 ====================

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"` // system / user / assistant
	Content string `json:"content"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Content      string `json:"content"`
	TokensInput  int    `json:"tokens_input"`
	TokensOutput int    `json:"tokens_output"`
	Model        string `json:"model"`
}

// Chat 同步聊天（非流式）
func (e *Engine) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = e.cfg.DefaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 4096
	}

	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := e.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return nil, fmt.Errorf("OpenAI调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI返回为空")
	}

	return &ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		TokensInput:  resp.Usage.PromptTokens,
		TokensOutput: resp.Usage.CompletionTokens,
		Model:        resp.Model,
	}, nil
}

// ==================== 流式调用 ====================

// StreamCallback 流式回调函数
type StreamCallback func(chunk string, done bool, err error)

// ChatStream 流式聊天（SSE用）
func (e *Engine) ChatStream(ctx context.Context, req ChatRequest, callback StreamCallback) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = e.cfg.DefaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 4096
	}

	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, m := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	stream, err := e.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: float32(req.Temperature),
		MaxTokens:   req.MaxTokens,
		Stream:      true,
	})
	if err != nil {
		return nil, fmt.Errorf("创建流式请求失败: %w", err)
	}
	defer stream.Close()

	var fullContent strings.Builder
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			callback("", true, nil)
			break
		}
		if err != nil {
			callback("", true, err)
			return nil, fmt.Errorf("流式接收失败: %w", err)
		}

		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta.Content
			fullContent.WriteString(delta)
			callback(delta, false, nil)
		}
	}

	return &ChatResponse{
		Content: fullContent.String(),
		Model:   req.Model,
	}, nil
}
