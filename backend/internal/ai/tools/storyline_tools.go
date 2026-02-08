package tools

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/repository"
)

// GetStorylineStatusTool 获取三线状态工具
type GetStorylineStatusTool struct {
	storylineRepo *repository.StorylineRepository
}

// NewGetStorylineStatusTool 创建三线状态工具
func NewGetStorylineStatusTool(storylineRepo *repository.StorylineRepository) *GetStorylineStatusTool {
	return &GetStorylineStatusTool{
		storylineRepo: storylineRepo,
	}
}

func (t *GetStorylineStatusTool) GetName() string {
	return "get_storyline_status"
}

func (t *GetStorylineStatusTool) GetDescription() string {
	return "获取三线(天线/地线/剧情线)当前状态。参数: project_id(项目ID), line_type(线类型: skyline|天线, groundline|地线, plotline|剧情线, all|全部)"
}

func (t *GetStorylineStatusTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: project_id")
	}

	lineType := "all"
	if lt, ok := params["line_type"].(string); ok {
		lineType = lt
	}

	if lineType == "all" {
		// 获取所有三线
		skylines, err := t.storylineRepo.GetByType(ctx, int(projectID), "skyline")
		if err != nil {
			return nil, fmt.Errorf("failed to get skylines: %w", err)
		}

		groundlines, err := t.storylineRepo.GetByType(ctx, int(projectID), "groundline")
		if err != nil {
			return nil, fmt.Errorf("failed to get groundlines: %w", err)
		}

		plotlines, err := t.storylineRepo.GetByType(ctx, int(projectID), "plotline")
		if err != nil {
			return nil, fmt.Errorf("failed to get plotlines: %w", err)
		}

		return map[string]interface{}{
			"skylines":    skylines,
			"groundlines": groundlines,
			"plotlines":   plotlines,
		}, nil
	} else {
		// 获取指定类型的线
		lines, err := t.storylineRepo.GetByType(ctx, int(projectID), lineType)
		if err != nil {
			return nil, fmt.Errorf("failed to get %s: %w", lineType, err)
		}

		return map[string]interface{}{
			"line_type": lineType,
			"lines":     lines,
		}, nil
	}
}

// UpdateStorylineTool 更新三线工具
type UpdateStorylineTool struct {
	storylineRepo *repository.StorylineRepository
}

// NewUpdateStorylineTool 创建三线更新工具
func NewUpdateStorylineTool(storylineRepo *repository.StorylineRepository) *UpdateStorylineTool {
	return &UpdateStorylineTool{
		storylineRepo: storylineRepo,
	}
}

func (t *UpdateStorylineTool) GetName() string {
	return "update_storyline"
}

func (t *UpdateStorylineTool) GetDescription() string {
	return "更新三线规划内容。参数: storyline_id(线ID), title(标题), content(内容), status(状态)"
}

func (t *UpdateStorylineTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	storylineID, ok := params["storyline_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: storyline_id")
	}

	// 获取现有数据
	storyline, err := t.storylineRepo.GetByID(ctx, int(storylineID))
	if err != nil {
		return nil, fmt.Errorf("failed to get storyline: %w", err)
	}

	// 更新字段
	if title, ok := params["title"].(string); ok {
		storyline.Title = title
	}

	if content, ok := params["content"].(string); ok {
		storyline.Content = content
	}

	if status, ok := params["status"].(string); ok {
		storyline.Status = status
	}

	// 保存更新
	if err := t.storylineRepo.Update(ctx, storyline); err != nil {
		return nil, fmt.Errorf("failed to update storyline: %w", err)
	}

	return map[string]interface{}{
		"success":      true,
		"storyline_id": storylineID,
		"message":      "三线更新成功",
	}, nil
}

// CreateStorylineTool 创建三线工具
type CreateStorylineTool struct {
	storylineRepo *repository.StorylineRepository
}

// NewCreateStorylineTool 创建三线创建工具
func NewCreateStorylineTool(storylineRepo *repository.StorylineRepository) *CreateStorylineTool {
	return &CreateStorylineTool{
		storylineRepo: storylineRepo,
	}
}

func (t *CreateStorylineTool) GetName() string {
	return "create_storyline"
}

func (t *CreateStorylineTool) GetDescription() string {
	return "创建新的三线规划。参数: project_id(项目ID), line_type(线类型), title(标题), content(内容), chapter_range(章节范围, 可选)"
}

func (t *CreateStorylineTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: project_id")
	}

	lineType, ok := params["line_type"].(string)
	if !ok || lineType == "" {
		return nil, fmt.Errorf("missing required parameter: line_type")
	}

	title, ok := params["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("missing required parameter: title")
	}

	content, ok := params["content"].(string)
	if !ok {
		content = ""
	}

	storyline := &repository.Storyline{
		ProjectID: int(projectID),
		LineType:  lineType,
		Title:     title,
		Content:   content,
		Status:    "planned",
	}

	if err := t.storylineRepo.Create(ctx, storyline); err != nil {
		return nil, fmt.Errorf("failed to create storyline: %w", err)
	}

	return map[string]interface{}{
		"success":      true,
		"storyline_id": storyline.ID,
		"message":      "三线创建成功",
	}, nil
}
