package tools

import (
	"context"
	"fmt"
	"time"

	"novel-study/backend/internal/repository"
)

// UpdateStorylineTool 更新三线规划工具
type UpdateStorylineTool struct {
	storylineRepo *repository.StorylineRepository
}

// NewUpdateStorylineTool 创建更新三线工具
func NewUpdateStorylineTool(storylineRepo *repository.StorylineRepository) *UpdateStorylineTool {
	return &UpdateStorylineTool{
		storylineRepo: storylineRepo,
	}
}

func (t *UpdateStorylineTool) GetName() string {
	return "update_storyline"
}

func (t *UpdateStorylineTool) GetDescription() string {
	return "更新三线规划内容。参数：storyline_id(三线ID), title(标题), content(内容), status(状态)"
}

func (t *UpdateStorylineTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	start := time.Now()

	// 解析参数
	storylineID, ok := params["storyline_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing or invalid 'storyline_id' parameter")
	}

	// 获取原有数据
	storyline, err := t.storylineRepo.GetByID(ctx, int(storylineID))
	if err != nil {
		return nil, fmt.Errorf("failed to get storyline: %w", err)
	}

	// 更新字段
	if title, ok := params["title"].(string); ok && title != "" {
		storyline.Title = title
	}
	if content, ok := params["content"].(string); ok && content != "" {
		storyline.Content = content
	}
	if status, ok := params["status"].(string); ok && status != "" {
		storyline.Status = status
	}

	// 保存更新
	err = t.storylineRepo.Update(ctx, storyline)
	if err != nil {
		return nil, fmt.Errorf("failed to update storyline: %w", err)
	}

	// 返回结果
	response := map[string]interface{}{
		"storyline_id": storyline.ID,
		"title":        storyline.Title,
		"status":       storyline.Status,
		"updated":      true,
		"duration_ms":  time.Since(start).Milliseconds(),
	}

	return response, nil
}
