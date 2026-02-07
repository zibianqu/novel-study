package agents

import (
	"context"
	"fmt"
	"strings"

	"github.com/novelforge/backend/internal/ai"
	"github.com/novelforge/backend/internal/ai/prompts"
	"github.com/novelforge/backend/internal/model"
	"github.com/novelforge/backend/internal/repository"
)

// DirectorAgent 总导演Agent - 用户唯一对话入口
type DirectorAgent struct {
	executor      *AgentExecutor
	projectRepo   *repository.ProjectRepository
	chapterRepo   *repository.ChapterRepository
	storylineRepo *repository.StorylineRepository
	agentRepo     *repository.AgentRepository
}

// NewDirectorAgent 创建总导演Agent
func NewDirectorAgent(
	executor *AgentExecutor,
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	storylineRepo *repository.StorylineRepository,
	agentRepo *repository.AgentRepository,
) *DirectorAgent {
	return &DirectorAgent{
		executor:      executor,
		projectRepo:   projectRepo,
		chapterRepo:   chapterRepo,
		storylineRepo: storylineRepo,
		agentRepo:     agentRepo,
	}
}

// ChatRequest 总导演对话请求
type DirectorChatRequest struct {
	ProjectID   int    `json:"project_id"`
	Message     string `json:"message"`
	SessionID   string `json:"session_id"`
}

// Chat 与总导演对话（同步）
func (d *DirectorAgent) Chat(ctx context.Context, req DirectorChatRequest) (*ExecuteResponse, error) {
	// 1. 收集项目上下文
	projectCtx, err := d.gatherProjectContext(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("收集项目上下文失败: %w", err)
	}

	// 2. 执行总导演Agent
	return d.executor.Execute(ctx, ExecuteRequest{
		AgentKey:    "director",
		ProjectID:   req.ProjectID,
		Instruction: req.Message,
		Variables:   projectCtx,
	})
}

// ChatStream 与总导演对话（流式）
func (d *DirectorAgent) ChatStream(ctx context.Context, req DirectorChatRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	projectCtx, err := d.gatherProjectContext(ctx, req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("收集项目上下文失败: %w", err)
	}

	return d.executor.ExecuteStream(ctx, ExecuteRequest{
		AgentKey:    "director",
		ProjectID:   req.ProjectID,
		Instruction: req.Message,
		Variables:   projectCtx,
	}, callback)
}

// Forecast 多章推演
func (d *DirectorAgent) Forecast(ctx context.Context, projectID int, chapters int) (*ForecastResult, error) {
	projectCtx, err := d.gatherProjectContext(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 获取三线状态
	storylines, err := d.storylineRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	skyStatus, groundStatus, plotStatus := categorizeStorylines(storylines)

	vars := map[string]string{
		"forecast_chapters": fmt.Sprintf("%d", chapters),
		"skyline_status":    skyStatus,
		"groundline_status": groundStatus,
		"plotline_status":   plotStatus,
		"story_summary":     projectCtx["current_status"],
	}

	// 使用推演模板
	messages := prompts.AssembleAgentPrompt(
		prompts.ForecastPrompt,
		vars,
		"", "", "",
	)

	resp, err := d.executor.engine.Chat(ctx, ai.ChatRequest{
		Model:       "",
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
	})
	if err != nil {
		return nil, err
	}

	return &ForecastResult{
		Content:   resp.Content,
		Chapters:  chapters,
		ProjectID: projectID,
	}, nil
}

// ForecastResult 推演结果
type ForecastResult struct {
	Content   string `json:"content"`
	Chapters  int    `json:"chapters"`
	ProjectID int    `json:"project_id"`
}

// gatherProjectContext 收集项目上下文信息
func (d *DirectorAgent) gatherProjectContext(ctx context.Context, projectID int) (map[string]string, error) {
	vars := make(map[string]string)

	if projectID == 0 {
		vars["project_name"] = "未选择项目"
		vars["project_type"] = ""
		vars["project_genre"] = ""
		vars["chapter_count"] = "0"
		vars["word_count"] = "0"
		vars["current_status"] = "尚未开始"
		return vars, nil
	}

	project, err := d.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	vars["project_name"] = project.Title
	vars["project_type"] = project.Type
	vars["project_genre"] = project.Genre
	vars["chapter_count"] = fmt.Sprintf("%d", project.ChapterCount)
	vars["word_count"] = fmt.Sprintf("%d", project.WordCount)
	vars["current_status"] = project.Status

	// 获取三线状态
	storylines, err := d.storylineRepo.ListByProject(ctx, projectID)
	if err == nil {
		sky, ground, plot := categorizeStorylines(storylines)
		vars["current_status"] += fmt.Sprintf("\n\n天线状态：%s\n地线状态：%s\n剧情线状态：%s", sky, ground, plot)
	}

	return vars, nil
}

// categorizeStorylines 分类三线状态
func categorizeStorylines(storylines []model.Storyline) (sky, ground, plot string) {
	var skyParts, groundParts, plotParts []string
	for _, s := range storylines {
		summary := fmt.Sprintf("[%s] %s (Ch.%d-%d)", s.Status, s.Title, s.ChapterStart, s.ChapterEnd)
		switch s.LineType {
		case "skyline":
			skyParts = append(skyParts, summary)
		case "groundline":
			groundParts = append(groundParts, summary)
		case "plotline":
			plotParts = append(plotParts, summary)
		}
	}
	return strings.Join(skyParts, "\n"), strings.Join(groundParts, "\n"), strings.Join(plotParts, "\n")
}
