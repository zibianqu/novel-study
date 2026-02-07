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

// WriterAgents 写作执行层Agent集合（旁白 + 角色 + 审核）
type WriterAgents struct {
	executor    *AgentExecutor
	chapterRepo *repository.ChapterRepository
}

// NewWriterAgents 创建写作Agent集合
func NewWriterAgents(executor *AgentExecutor, chapterRepo *repository.ChapterRepository) *WriterAgents {
	return &WriterAgents{
		executor:    executor,
		chapterRepo: chapterRepo,
	}
}

// ==================== 续写 ====================

// ContinueRequest 续写请求
type ContinueRequest struct {
	ProjectID   int    `json:"project_id"`
	ChapterID   int    `json:"chapter_id"`
	Instruction string `json:"instruction"` // 用户指令
}

// Continue 续写（流式）
func (w *WriterAgents) Continue(ctx context.Context, req ContinueRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	// 获取当前章节的最近内容
	chapter, err := w.chapterRepo.GetByID(ctx, req.ChapterID)
	if err != nil {
		return nil, fmt.Errorf("获取章节失败: %w", err)
	}

	// 取最近2000字作为上下文
	recentContent := getRecentContent(chapter.Content, 2000)

	vars := map[string]string{
		"recent_content":   recentContent,
		"user_instruction": req.Instruction,
		"chapter_outline":  "", // TODO: 从大纲获取
	}

	// 使用旁白Agent续写
	return w.executor.ExecuteStream(ctx, ExecuteRequest{
		AgentKey:    "narrator",
		ProjectID:   req.ProjectID,
		Instruction: prompts.ContinueWritingPrompt.Render(vars),
		Variables: map[string]string{
			"writing_style":     "根据前文风格",
			"scene_description": "续接上文",
		},
	}, callback)
}

// ==================== 润色 ====================

// PolishRequest 润色请求
type PolishRequest struct {
	ProjectID   int    `json:"project_id"`
	ChapterID   int    `json:"chapter_id"`
	Text        string `json:"text"`
	Instruction string `json:"instruction"`
}

// Polish 润色（流式）
func (w *WriterAgents) Polish(ctx context.Context, req PolishRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	vars := map[string]string{
		"original_text":     req.Text,
		"polish_instruction": req.Instruction,
	}

	return w.executor.ExecuteStream(ctx, ExecuteRequest{
		AgentKey:    "narrator",
		ProjectID:   req.ProjectID,
		Instruction: prompts.PolishPrompt.Render(vars),
		Variables: map[string]string{
			"writing_style":     "保持原风格",
			"scene_description": "润色段落",
		},
	}, callback)
}

// ==================== 改写 ====================

// RewriteRequest 改写请求
type RewriteRequest struct {
	ProjectID   int    `json:"project_id"`
	ChapterID   int    `json:"chapter_id"`
	Text        string `json:"text"`
	Instruction string `json:"instruction"`
}

// Rewrite 改写（流式）
func (w *WriterAgents) Rewrite(ctx context.Context, req RewriteRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	vars := map[string]string{
		"original_text":      req.Text,
		"rewrite_instruction": req.Instruction,
	}

	return w.executor.ExecuteStream(ctx, ExecuteRequest{
		AgentKey:    "narrator",
		ProjectID:   req.ProjectID,
		Instruction: prompts.RewritePrompt.Render(vars),
		Variables: map[string]string{
			"writing_style":     "根据改写要求",
			"scene_description": "改写段落",
		},
	}, callback)
}

// ==================== 角色对话生成 ====================

// DialogueRequest 对话生成请求
type DialogueRequest struct {
	ProjectID    int      `json:"project_id"`
	ChapterID    int      `json:"chapter_id"`
	CharacterIDs []int    `json:"character_ids"`
	Scene        string   `json:"scene"`
	Instruction  string   `json:"instruction"`
}

// GenerateDialogue 生成角色对话（流式）
func (w *WriterAgents) GenerateDialogue(ctx context.Context, req DialogueRequest, callback ai.StreamCallback) (*ExecuteResponse, error) {
	vars := map[string]string{
		"scene_description":    req.Scene,
		"characters_info":      "", // TODO: 从数据库获取角色信息
		"dialogue_instruction": req.Instruction,
	}

	return w.executor.ExecuteStream(ctx, ExecuteRequest{
		AgentKey:    "character",
		ProjectID:   req.ProjectID,
		Instruction: prompts.DialoguePrompt.Render(vars),
		Variables: map[string]string{
			"characters_info":    "参与角色信息",
			"character_relations": "角色关系",
			"scene_description":  req.Scene,
		},
	}, callback)
}

// ==================== 审核 ====================

// ReviewRequest 审核请求
type ReviewRequest struct {
	ProjectID int    `json:"project_id"`
	ChapterID int    `json:"chapter_id"`
	Content   string `json:"content"` // 待审核内容
}

// ReviewResult 审核结果
type ReviewResult struct {
	OverallScore int                    `json:"overall_score"`
	Passed       bool                   `json:"passed"`
	Content      string                 `json:"content"` // 原始审核报告
	AgentOutput  *ExecuteResponse       `json:"agent_output"`
}

// Review 审核内容
func (w *WriterAgents) Review(ctx context.Context, req ReviewRequest) (*ReviewResult, error) {
	vars := map[string]string{
		"characters_info":     "", // TODO: 从数据库获取
		"world_settings":      "", // TODO: 从数据库获取
		"outline_requirement": "", // TODO: 从大纲获取
	}

	resp, err := w.executor.Execute(ctx, ExecuteRequest{
		AgentKey:    "reviewer",
		ProjectID:   req.ProjectID,
		Instruction: "请审核以下内容：\n\n" + req.Content,
		Variables:   vars,
	})
	if err != nil {
		return nil, err
	}

	// 简单解析审核分数（实际应该解析JSON）
	passed := true
	score := 80 // 默认分数
	if strings.Contains(resp.Content, "\"passed\": false") || strings.Contains(resp.Content, "\"passed\":false") {
		passed = false
		score = 60
	}

	return &ReviewResult{
		OverallScore: score,
		Passed:       passed,
		Content:      resp.Content,
		AgentOutput:  resp,
	}, nil
}

// ==================== 一致性检查 ====================

// ConsistencyCheck 一致性检查
func (w *WriterAgents) ConsistencyCheck(ctx context.Context, projectID, chapterID int) (*ExecuteResponse, error) {
	chapter, err := w.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	vars := map[string]string{
		"chapter_content": chapter.Content,
		"characters_info": "", // TODO
		"world_settings":  "", // TODO
	}

	return w.executor.Execute(ctx, ExecuteRequest{
		AgentKey:    "reviewer",
		ProjectID:   projectID,
		Instruction: prompts.ConsistencyCheckPrompt.Render(vars),
		Variables: map[string]string{
			"characters_info": "",
			"world_settings":  "",
		},
	})
}

// ==================== 章节创作完整流程 ====================

// ChapterWriteRequest 完整章节创作请求
type ChapterWriteRequest struct {
	ProjectID   int    `json:"project_id"`
	ChapterID   int    `json:"chapter_id"`
	Instruction string `json:"instruction"`
	MaxRetries  int    `json:"max_retries"` // 最大审核返工次数
}

// WriteChapter 完整的章节创作流程（旁白→对话→审核→返工循环）
func (w *WriterAgents) WriteChapter(ctx context.Context, req ChapterWriteRequest, callback ai.StreamCallback) (*ChapterWriteResult, error) {
	if req.MaxRetries == 0 {
		req.MaxRetries = 3
	}

	var lastContent string
	var lastReview *ReviewResult

	for attempt := 0; attempt <= req.MaxRetries; attempt++ {
		instruction := req.Instruction

		// 如果是返工，加入审核反馈
		if attempt > 0 && lastReview != nil {
			instruction = fmt.Sprintf("请根据以下审核意见修改内容：\n\n审核意见：%s\n\n原始指令：%s",
				lastReview.Content, req.Instruction)
		}

		// 1. 旁白Agent生成内容
		narratorResp, err := w.Continue(ctx, ContinueRequest{
			ProjectID:   req.ProjectID,
			ChapterID:   req.ChapterID,
			Instruction: instruction,
		}, callback)
		if err != nil {
			return nil, fmt.Errorf("第%d轮旁白生成失败: %w", attempt+1, err)
		}

		lastContent = narratorResp.Content

		// 2. 审核Agent检查
		review, err := w.Review(ctx, ReviewRequest{
			ProjectID: req.ProjectID,
			ChapterID: req.ChapterID,
			Content:   lastContent,
		})
		if err != nil {
			return nil, fmt.Errorf("第%d轮审核失败: %w", attempt+1, err)
		}

		lastReview = review

		// 3. 审核通过则结束
		if review.Passed {
			return &ChapterWriteResult{
				Content:  lastContent,
				Review:   review,
				Attempts: attempt + 1,
			}, nil
		}

		// 审核不通过，继续返工（最后一次不再返工）
		if attempt == req.MaxRetries {
			break
		}
	}

	// 超过最大返工次数，返回最后一版
	return &ChapterWriteResult{
		Content:  lastContent,
		Review:   lastReview,
		Attempts: req.MaxRetries + 1,
	}, nil
}

// ChapterWriteResult 章节创作结果
type ChapterWriteResult struct {
	Content  string        `json:"content"`
	Review   *ReviewResult `json:"review"`
	Attempts int           `json:"attempts"`
}

// ==================== 辅助函数 ====================

// getRecentContent 获取最近N个字符的内容
func getRecentContent(content string, maxChars int) string {
	runes := []rune(content)
	if len(runes) <= maxChars {
		return content
	}
	return string(runes[len(runes)-maxChars:])
}

// categorizeStorylines 已在 director.go 中定义，此处引用
// func categorizeStorylines 移除重复定义
