package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/service"
)

type KnowledgeHandler struct {
	service *service.KnowledgeService
}

func NewKnowledgeHandler(service *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{service: service}
}

// CreateKnowledge 创建知识
func (h *KnowledgeHandler) CreateKnowledge(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CreateKnowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kb, err := h.service.CreateKnowledge(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, kb)
}

// GetKnowledge 获取知识详情
func (h *KnowledgeHandler) GetKnowledge(c *gin.Context) {
	userID := c.GetInt("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	kb, err := h.service.GetKnowledge(id, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kb)
}

// GetProjectKnowledge 获取项目知识列表
func (h *KnowledgeHandler) GetProjectKnowledge(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	items, err := h.service.GetProjectKnowledge(projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"knowledge": items})
}

// SearchKnowledge 搜索知识
func (h *KnowledgeHandler) SearchKnowledge(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int    `json:"project_id" binding:"required"`
		Query     string `json:"query" binding:"required"`
		TopK      int    `json:"top_k"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TopK == 0 {
		req.TopK = 5
	}

	docs, err := h.service.SearchKnowledge(c.Request.Context(), req.ProjectID, userID, req.Query, req.TopK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": docs})
}

// DeleteKnowledge 删除知识
func (h *KnowledgeHandler) DeleteKnowledge(c *gin.Context) {
	userID := c.GetInt("user_id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.service.DeleteKnowledge(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
