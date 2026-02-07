package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/novel-study/internal/service"
)

type AIHandler struct {
	service *service.AIService
}

func NewAIHandler(service *service.AIService) *AIHandler {
	return &AIHandler{service: service}
}

// Chat 与总导演对话
func (h *AIHandler) Chat(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int    `json:"project_id" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Chat(c.Request.Context(), userID, req.ProjectID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ChatStream 流式对话 (SSE)
func (h *AIHandler) ChatStream(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int    `json:"project_id" binding:"required"`
		Message   string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置SSE头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	// 流式输出回调
	callback := func(chunk string) {
		c.SSEvent("message", chunk)
		flusher.Flush()
	}

	err := h.service.ChatStream(c.Request.Context(), userID, req.ProjectID, req.Message, callback)
	if err != nil {
		c.SSEvent("error", err.Error())
		return
	}

	c.SSEvent("done", "")
}

// GenerateChapter 生成章节
func (h *AIHandler) GenerateChapter(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID    int    `json:"project_id" binding:"required"`
		ChapterTitle string `json:"chapter_title" binding:"required"`
		Outline      string `json:"outline"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.GenerateChapter(c.Request.Context(), userID, req.ProjectID, req.ChapterTitle, req.Outline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CheckQuality 质量检查
func (h *AIHandler) CheckQuality(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int    `json:"project_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CheckQuality(c.Request.Context(), userID, req.ProjectID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAgents 获取Agent列表
func (h *AIHandler) GetAgents(c *gin.Context) {
	agents, err := h.service.GetAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agents": agents})
}
