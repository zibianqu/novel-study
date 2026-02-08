package ai

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"novel-study/backend/internal/ai/agents"
	"novel-study/backend/internal/ai/rag"
	"novel-study/backend/internal/ai/tools"
	"novel-study/backend/internal/config"
	"novel-study/backend/internal/repository"
)

// Engine AI引擎
type Engine struct {
	config         *config.Config
	agents         map[string]Agent
	apiKey         string
	mu             sync.RWMutex // 保护并发访问
	toolRegistry   *tools.ToolRegistry
	retriever      *rag.Retriever
	projectRepo    *repository.ProjectRepository
	chapterRepo    *repository.ChapterRepository
	storylineRepo  *repository.StorylineRepository
	neo4jRepo      *repository.Neo4jRepository
}

// NewEngine 创建新的AI引擎
func NewEngine(
	cfg *config.Config,
	db *sql.DB,
	retriever *rag.Retriever,
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	storylineRepo *repository.StorylineRepository,
	neo4jRepo *repository.Neo4jRepository,
) *Engine {
	// 初始化工具注册表
	toolLogger := tools.NewDBToolCallLogger(db)
	toolRegistry := tools.NewToolRegistry(toolLogger)

	engine := &Engine{
		config:        cfg,
		agents:        make(map[string]Agent),
		apiKey:        cfg.OpenAIAPIKey,
		toolRegistry:  toolRegistry,
		retriever:     retriever,
		projectRepo:   projectRepo,
		chapterRepo:   chapterRepo,
		storylineRepo: storylineRepo,
		neo4jRepo:     neo4jRepo,
	}

	// 注册工具
	engine.RegisterTools()

	// 注册Agents
	engine.RegisterCoreAgents()

	fmt.Println("✅ AI 引擎初始化完成，已注册 7 个 Agent")
	fmt.Printf("✅ 工具系统初始化完成，已注册 %d 个工具\n", len(engine.toolRegistry.ListTools()))

	return engine
}

// RegisterTools 注册所有工具
func (e *Engine) RegisterTools() {
	// RAG 检索工具
	e.toolRegistry.Register(tools.NewRAGSearchTool(e.retriever))

	// Neo4j 查询工具
	e.toolRegistry.Register(tools.NewNeo4jQueryTool(e.neo4jRepo))

	// 项目状态工具
	e.toolRegistry.Register(tools.NewGetProjectStatusTool(e.projectRepo, e.chapterRepo))
	e.toolRegistry.Register(tools.NewGetChapterContentTool(e.chapterRepo))

	// 三线管理工具
	e.toolRegistry.Register(tools.NewGetStorylineStatusTool(e.storylineRepo))
	e.toolRegistry.Register(tools.NewUpdateStorylineTool(e.storylineRepo))
	e.toolRegistry.Register(tools.NewCreateStorylineTool(e.storylineRepo))
}

// RegisterCoreAgents 注册核心Agent
func (e *Engine) RegisterCoreAgents() {
	// Agent 0: 总导演
	e.RegisterAgent("agent_0_director", agents.NewDirectorAgent(e.apiKey, e.toolRegistry))

	// Agent 1: 旁白叙述者
	e.RegisterAgent("agent_1_narrator", agents.NewNarratorAgent(e.apiKey, e.toolRegistry))

	// Agent 2: 角色扮演者
	e.RegisterAgent("agent_2_character", agents.NewCharacterAgent(e.apiKey, e.toolRegistry))

	// Agent 3: 审核导演
	e.RegisterAgent("agent_3_quality", agents.NewQualityAgent(e.apiKey, e.toolRegistry))

	// Agent 4: 天线掌控者
	e.RegisterAgent("agent_4_skyline", agents.NewSkylineAgent(e.apiKey, e.toolRegistry))

	// Agent 5: 地线掌控者
	e.RegisterAgent("agent_5_groundline", agents.NewGroundlineAgent(e.apiKey, e.toolRegistry))

	// Agent 6: 剧情线掌控者
	e.RegisterAgent("agent_6_plotline", agents.NewPlotlineAgent(e.apiKey, e.toolRegistry))
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

// GetToolRegistry 获取工具注册表
func (e *Engine) GetToolRegistry() *tools.ToolRegistry {
	return e.toolRegistry
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

// ExecuteTool 执行工具
func (e *Engine) ExecuteTool(ctx context.Context, agentID int, toolName string, params map[string]interface{}) (interface{}, error) {
	return e.toolRegistry.Execute(ctx, agentID, toolName, params)
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

// ListTools 获取所有工具列表
func (e *Engine) ListTools() []tools.ToolInfo {
	return e.toolRegistry.ListTools()
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
