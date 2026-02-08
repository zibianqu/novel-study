package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/repository"
)

// GetStorylineStatusTool 获取三线状态工具
type GetStorylineStatusTool struct {
	storylineRepo *repository.StorylineRepository
}

// NewGetStorylineStatusTool 创建获取三线状态工具
func NewGetStorylineStatusTool(storylineRepo *repository.StorylineRepository) *GetStorylineStatusTool {
	return &GetStorylineStatusTool{
		storylineRepo: storylineRepo,
	}
}

func (t *GetStorylineStatusTool) GetName() string {
	return "get_storyline_status"
}

func (t *GetStorylineStatusTool) GetDescription() string {
	return "获取三线当前状态（天线/地线/剧情线）。参数：project_id(项目ID), line_type(可选，skyline/groundline/plotline)"
}

func (t *GetStorylineStatusTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'project_id' parameter")
	}

	// 可选线类型
	lineType := ""
	if lt, ok := params["line_type"].(string); ok {
		lineType = lt
	}

	// 获取三线数据
	var storylines []*repository.Storyline
	var err error

	if lineType == "" {
		storylines, err = t.storylineRepo.GetByProjectID(ctx, int(projectID))
	} else {
		storylines, err = t.storylineRepo.GetByProjectIDAndType(ctx, int(projectID), lineType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get storylines: %w", err)
	}

	// 按类型分组
	grouped := map[string][]map[string]interface{}{
		"skyline":    {},
		"groundline": {},
		"plotline":   {},
	}

	for _, sl := range storylines {
		item := map[string]interface{}{
			"id":            sl.ID,
			"title":         sl.Title,
			"content":       sl.Content,
			"status":        sl.Status,
			"chapter_range": sl.ChapterRange,
			"sort_order":    sl.SortOrder,
		}
		grouped[sl.LineType] = append(grouped[sl.LineType], item)
	}

	// 返回结果
	response := map[string]interface{}{
		"project_id":  int(projectID),
		"storylines":  grouped,
		"total_count": len(storylines),
		"duration_ms": time.Since(start).Milliseconds(),
	}

	return response, nil
}
