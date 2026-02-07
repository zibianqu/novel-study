package openai

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
	"github.com/zibianqu/novel-study/internal/ai"
)

// Client OpenAI 客户端封装
type Client struct {
	client *openai.Client
}

// NewClient 创建 OpenAI 客户端
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		return nil
	}
	return &Client{
		client: openai.NewClient(apiKey),
	}
}

// ChatCompletion 聊天完成
func (c *Client) ChatCompletion(ctx context.Context, messages []ai.ChatMessage, model string, temperature float64, maxTokens int) (*ai.AgentResponse, error) {
	if c.client == nil {
		return nil, errors.New("OpenAI client not initialized")
	}

	// 转换消息格式
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 调用 API
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    openaiMessages,
		Temperature: float32(temperature),
		MaxTokens:   maxTokens,
	})

	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, errors.New("no response from OpenAI")
	}

	return &ai.AgentResponse{
		Content:    resp.Choices[0].Message.Content,
		TokensUsed: resp.Usage.TotalTokens,
		Metadata: map[string]interface{}{
			"model":             resp.Model,
			"finish_reason":     resp.Choices[0].FinishReason,
			"prompt_tokens":     resp.Usage.PromptTokens,
			"completion_tokens": resp.Usage.CompletionTokens,
		},
	}, nil
}

// ChatCompletionStream 流式聊天完成
func (c *Client) ChatCompletionStream(ctx context.Context, messages []ai.ChatMessage, model string, temperature float64, maxTokens int, callback func(string)) error {
	if c.client == nil {
		return errors.New("OpenAI client not initialized")
	}

	// 转换消息格式
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 创建流式请求
	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:       model,
		Messages:    openaiMessages,
		Temperature: float32(temperature),
		MaxTokens:   maxTokens,
		Stream:      true,
	})

	if err != nil {
		return err
	}
	defer stream.Close()

	// 接收流式响应
	for {
		response, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if len(response.Choices) > 0 {
			callback(response.Choices[0].Delta.Content)
		}
	}

	return nil
}

// CreateEmbedding 创建向量嵌入
func (c *Client) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	if c.client == nil {
		return nil, errors.New("OpenAI client not initialized")
	}

	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	// 调用 OpenAI Embedding API
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2, // text-embedding-ada-002
		Input: texts,
	})

	if err != nil {
		return nil, err
	}

	// 转换为 [][]float32
	embeddings := make([][]float32, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = data.Embedding
	}

	return embeddings, nil
}
