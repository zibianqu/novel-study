package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/repository"
)

// GetProjectStatusTool 获取项目状态工具
type GetProjectStatusTool struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
}

// NewGetProjectStatusTool 创建获取项目状态工具
func NewGetProjectStatusTool(projectRepo *repository.ProjectRepository, chapterRepo *repository.ChapterRepository) *GetProjectStatusTool {
	return &GetProjectStatusTool{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
	}
}

func (t *GetProjectStatusTool) GetName() string {
	return "get_project_status"
}

func (t *GetProjectStatusTool) GetDescription() string {
	return "获取项目当前状态，包括项目信息、章节数、字数等。参数：project_id(项目ID)"
}

func (t *GetProjectStatusTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'project_id' parameter")
	}

	// 获取项目信息
	project, err := t.projectRepo.GetByID(ctx, int(projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// 获取章节列表
	chapters, err := t.chapterRepo.GetByProjectID(ctx, int(projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	// 统计信息
	totalWords := 0
	completedChapters := 0
	for _, ch := range chapters {
		totalWords += ch.WordCount
		if ch.Status == "completed" {
			completedChapters++
		}
	}

	// 返回结果
	response := map[string]interface{}{
		"project_id":         project.ID,
		"title":              project.Title,
		"type":               project.Type,
		"genre":              project.Genre,
		"status":             project.Status,
		"total_chapters":     len(chapters),
		"completed_chapters": completedChapters,
		"total_words":        totalWords,
		"duration_ms":        time.Since(start).Milliseconds(),
	}

	return response, nil
}
