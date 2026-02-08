package tools

import (
	"context"
	"fmt"
	"novel-study/backend/internal/repository"
)

// GetProjectStatusTool 获取项目状态工具
type GetProjectStatusTool struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
}

// NewGetProjectStatusTool 创建项目状态工具
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
	return "获取项目当前状态，包括总章节数、总字数、最后更新时间等。参数: project_id(项目ID)"
}

func (t *GetProjectStatusTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: project_id")
	}

	// 获取项目基本信息
	project, err := t.projectRepo.GetByID(ctx, int(projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// 获取章节统计
	chapters, err := t.chapterRepo.GetByProjectID(ctx, int(projectID))
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	totalWordCount := 0
	publishedCount := 0
	draftCount := 0

	for _, ch := range chapters {
		totalWordCount += ch.WordCount
		if ch.Status == "published" {
			publishedCount++
		} else if ch.Status == "draft" {
			draftCount++
		}
	}

	return map[string]interface{}{
		"project_id":       project.ID,
		"title":            project.Title,
		"type":             project.Type,
		"genre":            project.Genre,
		"total_chapters":   len(chapters),
		"published_count":  publishedCount,
		"draft_count":      draftCount,
		"total_word_count": totalWordCount,
		"status":           project.Status,
		"created_at":       project.CreatedAt,
		"updated_at":       project.UpdatedAt,
	}, nil
}

// GetChapterContentTool 获取章节内容工具
type GetChapterContentTool struct {
	chapterRepo *repository.ChapterRepository
}

// NewGetChapterContentTool 创建章节内容工具
func NewGetChapterContentTool(chapterRepo *repository.ChapterRepository) *GetChapterContentTool {
	return &GetChapterContentTool{
		chapterRepo: chapterRepo,
	}
}

func (t *GetChapterContentTool) GetName() string {
	return "get_chapter_content"
}

func (t *GetChapterContentTool) GetDescription() string {
	return "获取指定章节的内容。参数: chapter_id(章节ID) 或 chapter_number(章节号) + project_id(项目ID)"
}

func (t *GetChapterContentTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	var chapter interface{}
	var err error

	// 支持两种查询方式
	if chapterID, ok := params["chapter_id"].(float64); ok {
		chapter, err = t.chapterRepo.GetByID(ctx, int(chapterID))
	} else if chapterNum, ok := params["chapter_number"].(float64); ok {
		projectID, ok := params["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id required when using chapter_number")
		}
		chapter, err = t.chapterRepo.GetByNumber(ctx, int(projectID), int(chapterNum))
	} else {
		return nil, fmt.Errorf("either chapter_id or (chapter_number + project_id) required")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get chapter: %w", err)
	}

	return chapter, nil
}
