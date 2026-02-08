package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/repository"
)

// GetChapterContentTool 获取章节内容工具
type GetChapterContentTool struct {
	chapterRepo *repository.ChapterRepository
}

// NewGetChapterContentTool 创建获取章节内容工具
func NewGetChapterContentTool(chapterRepo *repository.ChapterRepository) *GetChapterContentTool {
	return &GetChapterContentTool{
		chapterRepo: chapterRepo,
	}
}

func (t *GetChapterContentTool) GetName() string {
	return "get_chapter_content"
}

func (t *GetChapterContentTool) GetDescription() string {
	return "获取指定章节的内容。参数：chapter_id(章节ID)或 project_id+chapter_number(项目ID+章节号)"
}

func (t *GetChapterContentTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	var chapter *repository.Chapter
	var err error

	// 方式1：通过 chapter_id 获取
	if chapterID, ok := params["chapter_id"].(float64); ok {
		chapter, err = t.chapterRepo.GetByID(ctx, int(chapterID))
		if err != nil {
			return nil, fmt.Errorf("failed to get chapter by ID: %w", err)
		}
	} else {
		// 方式2：通过 project_id + chapter_number 获取
		projectID, ok1 := params["project_id"].(float64)
		chapterNumber, ok2 := params["chapter_number"].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("missing required parameters: chapter_id or (project_id + chapter_number)")
		}

		chapter, err = t.chapterRepo.GetByProjectIDAndNumber(ctx, int(projectID), int(chapterNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to get chapter: %w", err)
		}
	}

	// 返回结果
	response := map[string]interface{}{
		"chapter_id":   chapter.ID,
		"project_id":   chapter.ProjectID,
		"title":        chapter.Title,
		"content":      chapter.Content,
		"word_count":   chapter.WordCount,
		"status":       chapter.Status,
		"sort_order":   chapter.SortOrder,
		"duration_ms":  time.Since(start).Milliseconds(),
	}

	return response, nil
}
