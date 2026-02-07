package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zibianqu/novel-study/internal/config"
)

// Engine AI引擎
type Engine struct {
	config  *config.Config
	agents  map[string]Agent
	apiKey  string
}

// NewEngine 创建新的AI引擎
func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		config: cfg,
		agents: make(map[string]Agent),
		apiKey: cfg.OpenAIAPIKey,
	}

	// 注册核心Agent
	engine.RegisterCoreAgents()

	return engine
}

// RegisterCoreAgents 注册核心Agent
func (e *Engine) RegisterCoreAgents() {
	// Agent 0: 总导演
	e.RegisterAgent("agent_0_director", NewDirectorAgent(e.apiKey))

	// Agent 1: 旁白叙述者
	e.RegisterAgent("agent_1_narrator", NewNarratorAgent(e.apiKey))

	// Agent 2: 角色扉演者
	e.RegisterAgent("agent_2_character", NewCharacterAgent(e.apiKey))

	// Agent 3: 审核导演
	e.RegisterAgent("agent_3_quality", NewQualityAgent(e.apiKey))

	// TODO: Agent 4-6 在后续实现
}

// RegisterAgent 注册Agent
func (e *Engine) RegisterAgent(key string, agent Agent) {
	e.agents[key] = agent
}

// GetAgent 获取Agent
func (e *Engine) GetAgent(key string) (Agent, error) {
	agent, ok := e.agents[key]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", key)
	}
	return agent, nil
}

// ExecuteAgent 执行Agent
func (e *Engine) ExecuteAgent(ctx context.Context, agentKey string, req *AgentRequest) (*AgentResponse, error) {
	agent, err := e.GetAgent(agentKey)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	resp, err := agent.Execute(ctx, req)
	if err != nil {
		return nil, err
	}

	resp.DurationMs = time.Since(startTime).Milliseconds()
	return resp, nil
}

// ExecuteAgentStream 流式执行Agent
func (e *Engine) ExecuteAgentStream(ctx context.Context, agentKey string, req *AgentRequest, callback func(string)) error {
	agent, err := e.GetAgent(agentKey)
	if err != nil {
		return err
	}

	return agent.ExecuteStream(ctx, req, callback)
}

// ChatCompletion 通用聊天完成
func (e *Engine) ChatCompletion(ctx context.Context, messages []ChatMessage, model string, temperature float64, maxTokens int) (string, error) {
	if e.apiKey == "" {
		return "", errors.New("OpenAI API key not configured")
	}

	// TODO: 实际集成 OpenAI API
	// 这里先返回模拟响应
	return e.mockChatCompletion(messages), nil
}

// mockChatCompletion 模拟聊天完成（用于测试）
func (e *Engine) mockChatCompletion(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "这是一个模拟响应。"
	}

	lastMsg := messages[len(messages)-1]
	return fmt.Sprintf("模拟 AI 响应: 收到您的消息 '%s'", lastMsg.Content)
}
