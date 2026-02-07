package agents

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/novelforge/backend/internal/ai"
	"github.com/novelforge/backend/internal/ai/prompts"
	"github.com/novelforge/backend/internal/model"
	"github.com/novelforge/backend/internal/repository"
)

// AgentExecutor Agent执行器 - 调用指定Agent完成任务
type AgentExecutor struct {
	engine       *ai.Engine
	embedding    *ai.EmbeddingService
	agentRepo    *repository.AgentRepository
	db           *pgxpool.Pool
}

// NewAgentExecutor 创建Agent执行器
func NewAgentExecutor(
	engine *ai.Engine,
	embedding *ai.EmbeddingService,
	agentRepo *repository.AgentRepository,
	db *pgxpool.Pool,
) *AgentExecutor {
	return &AgentExecutor{
		engine:    engine,
		embedding: embedding,
		agentRepo: agentRepo,
		db:        db,
	}
}

// ExecuteRequest Agent执行请求
type ExecuteRequest struct {
	AgentID     int               `json:"agent_id"`
	AgentKey    string            `json:"agent_key"`     // 或通过key查找
	ProjectID   int              `json:"project_id"`
	Instruction string           `json:"instruction"`    // 用户/总导演指令
	Variables   map[string]string `json:"variables"`      // 模板变量
	Context     string           `json:"context"`         // 额外上下文
	Stream      bool             `json:"stream"`          // 是否流式
}

// ExecuteResponse Agent执行结果
type ExecuteResponse struct {
	AgentID      int       `json:"agent_id"`
	AgentName    string    `json:"agent_name"`
	Content      string    `json:"content"`
	TokensInput  int       `json:"tokens_input"`
	TokensOutput int       `json:"tokens_output"`
	DurationMs   int64     `json:"duration_ms"`
	Model        string    `json:"model"`
}

// Execute 执行Agent任务（同步）
func (e *AgentExecutor) Execute(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	start := time.Now()

	// 1. 获取Agent配置
	agent, err := e.getAgent(ctx, req)
	if err != nil {
		return nil, err
	}

	// 2. 检索Agent专属知识
	knowledgeCtx := ""
	if req.Instruction != "" {
		results, err := e.embedding.SearchAgentKnowledge(ctx, agent.ID, req.Instruction, 5)
		if err != nil {
			log.Printf("Agent知识检索失败(非致命): %v", err)
		} else if len(results) > 0 {
			var parts []string
			for _, r := range results {
				parts = append(parts, r.Text)
			}
			knowledgeCtx = strings.Join(parts, "\n---\n")
		}
	}

	// 3. 检索项目内容（RAG）
	ragCtx := ""
	if req.ProjectID > 0 && req.Instruction != "" {
		results, err := e.embedding.SearchContent(ctx, req.ProjectID, req.Instruction, 5)
		if err != nil {
			log.Printf("内容RAG检索失败(非致命): %v", err)
		} else if len(results) > 0 {
			var parts []string
			for _, r := range results {
				parts = append(parts, r.Text)
			}
			ragCtx = strings.Join(parts, "\n---\n")
		}
	}

	// 4. 组装Prompt
	systemTemplate := &prompts.Template{
		Name:     agent.AgentKey,
		Template: agent.SystemPrompt,
	}

	vars := req.Variables
	if vars == nil {
		vars = make(map[string]string)
	}

	messages := prompts.AssembleAgentPrompt(
		systemTemplate,
		vars,
		knowledgeCtx,
		ragCtx,
		req.Instruction,
	)

	// 加入额外上下文
	if req.Context != "" {
		messages = append(messages[:1], append([]ai.Message{
			{Role: "user", Content: "额外参考信息：\n" + req.Context},
			{Role: "assistant", Content: "好的，我已了解这些信息，将在后续创作中参考使用。"},
		}, messages[1:]...)...)
	}

	// 5. 调用AI
	chatReq := ai.ChatRequest{
		Model:       agent.Model,
		Messages:    messages,
		Temperature: agent.Temperature,
		MaxTokens:   agent.MaxTokens,
	}

	resp, err := e.engine.Chat(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("Agent[%s]执行失败: %w", agent.Name, err)
	}

	duration := time.Since(start).Milliseconds()

	// 6. 记录日志
	e.logInteraction(ctx, agent, req, resp, duration)

	return &ExecuteResponse{
		AgentID:      agent.ID,
		AgentName:    agent.Name,
		Content:      resp.Content,
		TokensInput:  resp.TokensInput,
		TokensOutput: resp.TokensOutput,
		DurationMs:   duration,
		Model:        resp.Model,
	}, nil
}

// ExecuteStream 执行Agent任务（流式）
func (e *AgentExecutor) ExecuteStream(ctx context.Context, req ExecuteRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	start := time.Now()

	agent, err := e.getAgent(ctx, req)
	if err != nil {
		return nil, err
	}

	// 知识检索（同 Execute）
	knowledgeCtx := ""
	if req.Instruction != "" {
		results, _ := e.embedding.SearchAgentKnowledge(ctx, agent.ID, req.Instruction, 5)
		if len(results) > 0 {
			var parts []string
			for _, r := range results {
				parts = append(parts, r.Text)
			}
			knowledgeCtx = strings.Join(parts, "\n---\n")
		}
	}

	ragCtx := ""
	if req.ProjectID > 0 && req.Instruction != "" {
		results, _ := e.embedding.SearchContent(ctx, req.ProjectID, req.Instruction, 5)
		if len(results) > 0 {
			var parts []string
			for _, r := range results {
				parts = append(parts, r.Text)
			}
			ragCtx = strings.Join(parts, "\n---\n")
		}
	}

	systemTemplate := &prompts.Template{
		Name:     agent.AgentKey,
		Template: agent.SystemPrompt,
	}

	vars := req.Variables
	if vars == nil {
		vars = make(map[string]string)
	}

	messages := prompts.AssembleAgentPrompt(systemTemplate, vars, knowledgeCtx, ragCtx, req.Instruction)

	chatReq := ai.ChatRequest{
		Model:       agent.Model,
		Messages:    messages,
		Temperature: agent.Temperature,
		MaxTokens:   agent.MaxTokens,
	}

	resp, err := e.engine.ChatStream(ctx, chatReq, callback)
	if err != nil {
		return nil, fmt.Errorf("Agent[%s]流式执行失败: %w", agent.Name, err)
	}

	duration := time.Since(start).Milliseconds()

	return &ExecuteResponse{
		AgentID:    agent.ID,
		AgentName:  agent.Name,
		Content:    resp.Content,
		DurationMs: duration,
		Model:      resp.Model,
	}, nil
}

// getAgent 获取Agent配置
func (e *AgentExecutor) getAgent(ctx context.Context, req ExecuteRequest) (*model.Agent, error) {
	if req.AgentID > 0 {
		return e.agentRepo.GetByID(ctx, req.AgentID)
	}
	if req.AgentKey != "" {
		return e.agentRepo.GetByKey(ctx, req.AgentKey)
	}
	return nil, fmt.Errorf("必须提供agent_id或agent_key")
}

// logInteraction 记录AI交互日志
func (e *AgentExecutor) logInteraction(ctx context.Context, agent *model.Agent, req ExecuteRequest, resp *ai.ChatResponse, durationMs int64) {
	_, err := e.db.Exec(ctx,
		`INSERT INTO ai_interaction_logs 
		 (agent_id, project_id, action_type, input_prompt, output_response, 
		  tokens_input, tokens_output, model, duration_ms)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		agent.ID, nilIfZero(req.ProjectID), "agent_execute",
		req.Instruction, resp.Content,
		resp.TokensInput, resp.TokensOutput, resp.Model, durationMs,
	)
	if err != nil {
		log.Printf("记录AI日志失败(非致命): %v", err)
	}
}

func nilIfZero(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}
