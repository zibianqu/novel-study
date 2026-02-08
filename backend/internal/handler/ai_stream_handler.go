package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/zibianqu/novel-study/internal/ai"

	"github.com/gin-gonic/gin"
)

// AIStreamHandler AI 流式生成Handler
type AIStreamHandler struct {
	aiEngine *ai.Engine
}

// NewAIStreamHandler 创建Handler
func NewAIStreamHandler(aiEngine *ai.Engine) *AIStreamHandler {
	return &AIStreamHandler{
		aiEngine: aiEngine,
	}
}

// ContinueWriteRequest 续写请求
type ContinueWriteRequest struct {
	ProjectID     int                    `json:"project_id" binding:"required"`
	ChapterID     int                    `json:"chapter_id"`
	Context       string                 `json:"context"`       // 上下文
	Length        int                    `json:"length"`        // 生成长度
	Style         string                 `json:"style"`         // 风格要求
	CustomPrompt  string                 `json:"custom_prompt"` // 自定义提示
	AgentID       int                    `json:"agent_id"`      // 指定 Agent (0=自动选择)
	ExtraContext  map[string]interface{} `json:"extra_context"` // 额外上下文
}

// PolishRequest 润色请求
type PolishRequest struct {
	ProjectID    int                    `json:"project_id" binding:"required"`
	Content      string                 `json:"content" binding:"required"`
	PolishType   string                 `json:"polish_type"` // "grammar", "style", "clarity", "all"
	CustomPrompt string                 `json:"custom_prompt"`
	ExtraContext map[string]interface{} `json:"extra_context"`
}

// RewriteRequest 改写请求
type RewriteRequest struct {
	ProjectID    int                    `json:"project_id" binding:"required"`
	Content      string                 `json:"content" binding:"required"`
	Instruction  string                 `json:"instruction" binding:"required"` // 改写指令
	Style        string                 `json:"style"`
	ExtraContext map[string]interface{} `json:"extra_context"`
}

// ChatRequest 对话请求
type ChatRequest struct {
	ProjectID    int                    `json:"project_id"`
	Message      string                 `json:"message" binding:"required"`
	AgentID      int                    `json:"agent_id"` // 指定 Agent (0=总导演)
	History      []ChatMessage          `json:"history"`
	ExtraContext map[string]interface{} `json:"extra_context"`
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`    // user, assistant
	Content string `json:"content"`
}

// ContinueWrite 续写接口 (SSE)
// @Summary 流式续写内容
// @Tags AI Stream
// @Accept json
// @Produce text/event-stream
// @Param request body ContinueWriteRequest true "续写请求"
// @Router /api/v1/ai/stream/continue [post]
func (h *AIStreamHandler) ContinueWrite(c *gin.Context) {
	var req ContinueWriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建 SSE 写入器
	stream := NewSSEStreamHandler(c)

	// 构建提示词
	prompt := h.buildContinuePrompt(&req)

	// 准备上下文
	context := req.ExtraContext
	if context == nil {
		context = make(map[string]interface{})
	}
	context["project_id"] = req.ProjectID
	context["chapter_id"] = req.ChapterID
	context["length"] = req.Length
	context["style"] = req.Style

	// 执行流式生成
	h.executeStream(c.Request.Context(), stream, req.AgentID, prompt, context)
}

// Polish 润色接口 (SSE)
// @Summary 流式润色内容
// @Tags AI Stream
// @Accept json
// @Produce text/event-stream
// @Param request body PolishRequest true "润色请求"
// @Router /api/v1/ai/stream/polish [post]
func (h *AIStreamHandler) Polish(c *gin.Context) {
	var req PolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream := NewSSEStreamHandler(c)
	prompt := h.buildPolishPrompt(&req)

	context := req.ExtraContext
	if context == nil {
		context = make(map[string]interface{})
	}
	context["project_id"] = req.ProjectID
	context["polish_type"] = req.PolishType
	context["original_content"] = req.Content

	// 使用审核导演 (Agent 3) 进行润色
	h.executeStream(c.Request.Context(), stream, 3, prompt, context)
}

// Rewrite 改写接口 (SSE)
// @Summary 流式改写内容
// @Tags AI Stream
// @Accept json
// @Produce text/event-stream
// @Param request body RewriteRequest true "改写请求"
// @Router /api/v1/ai/stream/rewrite [post]
func (h *AIStreamHandler) Rewrite(c *gin.Context) {
	var req RewriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream := NewSSEStreamHandler(c)
	prompt := h.buildRewritePrompt(&req)

	context := req.ExtraContext
	if context == nil {
		context = make(map[string]interface{})
	}
	context["project_id"] = req.ProjectID
	context["instruction"] = req.Instruction
	context["original_content"] = req.Content

	// 默认使用旁白叙述者 (Agent 1)
	h.executeStream(c.Request.Context(), stream, 1, prompt, context)
}

// Chat 对话接口 (SSE)
// @Summary 流式对话
// @Tags AI Stream
// @Accept json
// @Produce text/event-stream
// @Param request body ChatRequest true "对话请求"
// @Router /api/v1/ai/stream/chat [post]
func (h *AIStreamHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream := NewSSEStreamHandler(c)

	context := req.ExtraContext
	if context == nil {
		context = make(map[string]interface{})
	}
	if req.ProjectID > 0 {
		context["project_id"] = req.ProjectID
	}
	context["history"] = req.History

	// 默认使用总导演 (Agent 0)
	agentID := req.AgentID
	if agentID == 0 {
		agentID = 0
	}

	h.executeStream(c.Request.Context(), stream, agentID, req.Message, context)
}

// executeStream 执行流式生成
func (h *AIStreamHandler) executeStream(
	ctx context.Context,
	stream *SSEStreamHandler,
	agentID int,
	prompt string,
	context map[string]interface{},
) {
	start := time.Now()

	// 创建请求
	req := &ai.AgentRequest{
		Prompt:  prompt,
		Context: context,
	}

	// 执行流式生成
	err := h.aiEngine.ExecuteAgentStream(ctx, agentID, req, func(chunk string) {
		if err := stream.OnChunk(chunk); err != nil {
			log.Printf("Failed to send chunk: %v", err)
		}
	})

	if err != nil {
		stream.OnError(err)
		return
	}

	// 发送完成信号
	stream.OnComplete(map[string]interface{}{
		"duration_ms": time.Since(start).Milliseconds(),
		"agent_id":    agentID,
	})
}

// buildContinuePrompt 构建续写提示词
func (h *AIStreamHandler) buildContinuePrompt(req *ContinueWriteRequest) string {
	prompt := "请续写以下内容\n\n"

	if req.Context != "" {
		prompt += "上下文：\n" + req.Context + "\n\n"
	}

	if req.Style != "" {
		prompt += "风格要求：" + req.Style + "\n"
	}

	if req.Length > 0 {
		prompt += fmt.Sprintf("字数要求：%d 字左右\n", req.Length)
	}

	if req.CustomPrompt != "" {
		prompt += "\n" + req.CustomPrompt + "\n"
	}

	prompt += "\n请开始续写："

	return prompt
}

// buildPolishPrompt 构建润色提示词
func (h *AIStreamHandler) buildPolishPrompt(req *PolishRequest) string {
	prompt := "请对以下内容进行润色\n\n"

	prompt += "原文：\n" + req.Content + "\n\n"

	switch req.PolishType {
	case "grammar":
		prompt += "润色要求：修正语法错误和语句不通的地方\n"
	case "style":
		prompt += "润色要求：优化文笔风格，增加文学性\n"
	case "clarity":
		prompt += "润色要求：提高表达的清晰度和准确性\n"
	default:
		prompt += "润色要求：全面优化文字质量\n"
	}

	if req.CustomPrompt != "" {
		prompt += "\n" + req.CustomPrompt + "\n"
	}

	prompt += "\n请输出润色后的内容："

	return prompt
}

// buildRewritePrompt 构建改写提示词
func (h *AIStreamHandler) buildRewritePrompt(req *RewriteRequest) string {
	prompt := "请根据以下指令改写内容\n\n"

	prompt += "原文：\n" + req.Content + "\n\n"
	prompt += "改写指令：" + req.Instruction + "\n"

	if req.Style != "" {
		prompt += "风格要求：" + req.Style + "\n"
	}

	prompt += "\n请输出改写后的内容："

	return prompt
}

import "fmt"
