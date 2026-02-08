package service

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/repository"
)

// AIService AI服务
type AIService struct {
	engine      *ai.Engine
	agentRepo   *repository.AgentRepository
	projectRepo *repository.ProjectRepository
}

// NewAIService 创建AI服务
func NewAIService(engine *ai.Engine, agentRepo *repository.AgentRepository, projectRepo *repository.ProjectRepository) *AIService {
	return &AIService{
		engine:      engine,
		agentRepo:   agentRepo,
		projectRepo: projectRepo,
	}
}

// Chat 与总导演对话
func (s *AIService) Chat(ctx context.Context, userID, projectID int, message string) (*ai.AgentResponse, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问此项目")
	}

	// 构建Agent请求
	req := &ai.AgentRequest{
		UserID:    userID,
		ProjectID: projectID,
		Prompt:    message,
		Context: map[string]interface{}{
			"project_title": project.Title,
			"project_type":  project.Type,
			"project_genre": project.Genre,
		},
	}

	// 调用Agent 0 (总导演)
	resp, err := s.engine.ExecuteAgent(ctx, "agent_0_director", req)
	if err != nil {
		return nil, err
	}

	// 记录日志
	s.logInteraction(userID, projectID, 0, "chat", message, resp)

	return resp, nil
}

// ChatStream 流式对话
func (s *AIService) ChatStream(ctx context.Context, userID, projectID int, message string, callback func(string)) error {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return fmt.Errorf("无权访问此项目")
	}

	// 构建Agent请求
	req := &ai.AgentRequest{
		UserID:    userID,
		ProjectID: projectID,
		Prompt:    message,
		Stream:    true,
		Context: map[string]interface{}{
			"project_title": project.Title,
		},
	}

	// 调用Agent 0 流式输出
	return s.engine.ExecuteAgentStream(ctx, "agent_0_director", req, callback)
}

// GenerateChapter 生成章节
func (s *AIService) GenerateChapter(ctx context.Context, userID, projectID int, chapterTitle, outline string) (*ai.AgentResponse, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问此项目")
	}

	// 构建 Prompt
	prompt := fmt.Sprintf("请创作章节：%s\n\n大纲：%s", chapterTitle, outline)

	req := &ai.AgentRequest{
		UserID:    userID,
		ProjectID: projectID,
		Prompt:    prompt,
		Context: map[string]interface{}{
			"chapter_title": chapterTitle,
			"outline":       outline,
		},
	}

	// 调用Agent 1 (旁白叙述者)
	resp, err := s.engine.ExecuteAgent(ctx, "agent_1_narrator", req)
	if err != nil {
		return nil, err
	}

	// 记录日志
	s.logInteraction(userID, projectID, 1, "generate_chapter", prompt, resp)

	return resp, nil
}

// CheckQuality 质量检查
func (s *AIService) CheckQuality(ctx context.Context, userID, projectID int, content string) (*ai.AgentResponse, error) {
	req := &ai.AgentRequest{
		UserID:    userID,
		ProjectID: projectID,
		Prompt:    fmt.Sprintf("请审核以下内容：\n\n%s", content),
	}

	// 调用Agent 3 (审核导演)
	resp, err := s.engine.ExecuteAgent(ctx, "agent_3_quality", req)
	if err != nil {
		return nil, err
	}

	// 记录日志
	s.logInteraction(userID, projectID, 3, "quality_check", content, resp)

	return resp, nil
}

// GetAgents 获取Agent列表
func (s *AIService) GetAgents() ([]string, error) {
	return s.engine.ListAgents(), nil
}

// logInteraction 记录AI交互日志
func (s *AIService) logInteraction(userID, projectID, agentID int, actionType, input string, resp *ai.AgentResponse) {
	log := &model.AIInteractionLog{
		UserID:         userID,
		ProjectID:      &projectID,
		AgentID:        &agentID,
		ActionType:     actionType,
		InputPrompt:    input,
		OutputResponse: resp.Content,
		TokensInput:    0, // TODO: 计算
		TokensOutput:   resp.TokensUsed,
		Model:          "gpt-4o",
		DurationMs:     int(resp.DurationMs),
	}

	s.agentRepo.LogInteraction(log)
}
