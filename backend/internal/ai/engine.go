package ai

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/zibianqu/novel-study/internal/ai/agents"
	"github.com/zibianqu/novel-study/internal/config"
)

// Engine AI引擎
type Engine struct {
	config  *config.Config
	agents  map[string]Agent
	apiKey  string
	mu      sync.RWMutex // 保护并发访问
}

// NewEngine 创建新的AI引擎
func NewEngine(cfg *config.Config) *Engine {
	engine := &Engine{
		config: cfg,
		agents: make(map[string]Agent),
		apiKey: cfg.OpenAIAPIKey,
	}

	// 注册Agents
	engine.RegisterCoreAgents()

	return engine
}

// RegisterCoreAgents 注册核心Agent
func (e *Engine) RegisterCoreAgents() {
	// Agent 0: 总导演
	e.RegisterAgent("agent_0_director", agents.NewDirectorAgent(e.apiKey))

	// Agent 1: 旁白叙述者
	e.RegisterAgent("agent_1_narrator", agents.NewNarratorAgent(e.apiKey))

	// Agent 2: 角色扮演者
	e.RegisterAgent("agent_2_character", agents.NewCharacterAgent(e.apiKey))

	// Agent 3: 审核导演
	e.RegisterAgent("agent_3_quality", agents.NewQualityAgent(e.apiKey))

	// Agent 4: 天线掌控者
	e.RegisterAgent("agent_4_skyline", agents.NewSkylineAgent(e.apiKey))

	// Agent 5: 地线掌控者
	e.RegisterAgent("agent_5_groundline", agents.NewGroundlineAgent(e.apiKey))

	// Agent 6: 剧情线掌控者
	e.RegisterAgent("agent_6_plotline", agents.NewPlotlineAgent(e.apiKey))
}

// RegisterAgent 注册Agent（线程安全）
func (e *Engine) RegisterAgent(key string, agent Agent) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.agents[key] = agent
}

// GetAgent 获取Agent（线程安全）
func (e *Engine) GetAgent(key string) (Agent, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	agent, ok := e.agents[key]
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", key)
	}
	return agent, nil
}

// ExecuteAgent 执行Agent
func (e *Engine) ExecuteAgent(ctx context.Context, agentKey string, req *AgentRequest) (*AgentResponse, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	agent, err := e.GetAgent(agentKey)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	resp, err := agent.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("agent execution failed: %w", err)
	}

	resp.DurationMs = time.Since(startTime).Milliseconds()
	return resp, nil
}

// ExecuteAgentStream 流式执行Agent
func (e *Engine) ExecuteAgentStream(ctx context.Context, agentKey string, req *AgentRequest, callback func(string)) error {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	agent, err := e.GetAgent(agentKey)
	if err != nil {
		return err
	}

	return agent.ExecuteStream(ctx, req, callback)
}

// ListAgents 获取所有Agent列表（线程安全）
func (e *Engine) ListAgents() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	keys := make([]string, 0, len(e.agents))
	for key := range e.agents {
		keys = append(keys, key)
	}
	return keys
}

// ChatCompletion 通用聊天完成
func (e *Engine) ChatCompletion(ctx context.Context, messages []ChatMessage, model string, temperature float64, maxTokens int) (string, error) {
	if e.apiKey == "" {
		return "", errors.New("OpenAI API key not configured")
	}

	// 检查上下文
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
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
	return fmt.Sprintf("模拟AI响应: 收到您的消息 '%s'", lastMsg.Content)
}
