package handler

import (
	"net/http"
	"strconv"

	"github.com/zibianqu/novel-study/internal/repository"

	"github.com/gin-gonic/gin"
)

// AgentKnowledgeHandler Agent知识库Handler
type AgentKnowledgeHandler struct {
	knowledgeRepo *repository.AgentKnowledgeRepository
}

// NewAgentKnowledgeHandler 创建Handler
func NewAgentKnowledgeHandler(knowledgeRepo *repository.AgentKnowledgeRepository) *AgentKnowledgeHandler {
	return &AgentKnowledgeHandler{
		knowledgeRepo: knowledgeRepo,
	}
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	AgentID     int    `json:"agent_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
}

// CreateKnowledgeItemRequest 创建知识条目请求
type CreateKnowledgeItemRequest struct {
	AgentID      int      `json:"agent_id" binding:"required"`
	CategoryID   int      `json:"category_id" binding:"required"`
	Title        string   `json:"title" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Tags         []string `json:"tags"`
	QualityScore int      `json:"quality_score"`
}

// CreateCategory 创建分类
// @Summary 创建 Agent 知识库分类
// @Tags Agent Knowledge
// @Accept json
// @Produce json
// @Param request body CreateCategoryRequest true "分类信息"
// @Success 200 {object} repository.AgentKnowledgeCategory
// @Router /api/v1/agent-knowledge/categories [post]
func (h *AgentKnowledgeHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &repository.AgentKnowledgeCategory{
		AgentID:     req.AgentID,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
	}

	if err := h.knowledgeRepo.CreateCategory(c.Request.Context(), category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategories 获取 Agent 的所有分类
// @Summary 获取 Agent 知识库分类列表
// @Tags Agent Knowledge
// @Produce json
// @Param agent_id path int true "Agent ID"
// @Success 200 {array} repository.AgentKnowledgeCategory
// @Router /api/v1/agent-knowledge/agents/{agent_id}/categories [get]
func (h *AgentKnowledgeHandler) GetCategories(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent_id"})
		return
	}

	categories, err := h.knowledgeRepo.GetCategoriesByAgentID(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// CreateKnowledgeItem 创建知识条目
// @Summary 创建 Agent 知识条目
// @Tags Agent Knowledge
// @Accept json
// @Produce json
// @Param request body CreateKnowledgeItemRequest true "知识条目信息"
// @Success 200 {object} repository.AgentKnowledgeItem
// @Router /api/v1/agent-knowledge/items [post]
func (h *AgentKnowledgeHandler) CreateKnowledgeItem(c *gin.Context) {
	var req CreateKnowledgeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &repository.AgentKnowledgeItem{
		AgentID:      req.AgentID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Content:      req.Content,
		Tags:         req.Tags,
		QualityScore: req.QualityScore,
	}

	if err := h.knowledgeRepo.CreateKnowledgeItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// GetKnowledgeItems 获取知识条目列表
// @Summary 获取 Agent 知识条目列表
// @Tags Agent Knowledge
// @Produce json
// @Param agent_id path int true "Agent ID"
// @Param category_id query int false "分类ID"
// @Param limit query int false "限制数量"
// @Success 200 {array} repository.AgentKnowledgeItem
// @Router /api/v1/agent-knowledge/agents/{agent_id}/items [get]
func (h *AgentKnowledgeHandler) GetKnowledgeItems(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent_id"})
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	var categoryID *int
	if catStr := c.Query("category_id"); catStr != "" {
		if cid, err := strconv.Atoi(catStr); err == nil {
			categoryID = &cid
		}
	}

	items, err := h.knowledgeRepo.GetKnowledgeItemsByAgentID(c.Request.Context(), agentID, categoryID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// SearchKnowledgeItems 搜索知识条目
// @Summary 搜索 Agent 知识条目
// @Tags Agent Knowledge
// @Produce json
// @Param agent_id path int true "Agent ID"
// @Param keyword query string true "搜索关键词"
// @Param limit query int false "限制数量"
// @Success 200 {array} repository.AgentKnowledgeItem
// @Router /api/v1/agent-knowledge/agents/{agent_id}/search [get]
func (h *AgentKnowledgeHandler) SearchKnowledgeItems(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent_id"})
		return
	}

	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keyword is required"})
		return
	}

	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	items, err := h.knowledgeRepo.SearchKnowledgeItems(c.Request.Context(), agentID, keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// GetKnowledgeStats 获取知识库统计
// @Summary 获取 Agent 知识库统计
// @Tags Agent Knowledge
// @Produce json
// @Param agent_id path int true "Agent ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-knowledge/agents/{agent_id}/stats [get]
func (h *AgentKnowledgeHandler) GetKnowledgeStats(c *gin.Context) {
	agentID, err := strconv.Atoi(c.Param("agent_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid agent_id"})
		return
	}

	stats, err := h.knowledgeRepo.GetKnowledgeStats(c.Request.Context(), agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateKnowledgeItem 更新知识条目
// @Summary 更新 Agent 知识条目
// @Tags Agent Knowledge
// @Accept json
// @Produce json
// @Param id path int true "条目ID"
// @Param request body CreateKnowledgeItemRequest true "知识条目信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-knowledge/items/{id} [put]
func (h *AgentKnowledgeHandler) UpdateKnowledgeItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	var req CreateKnowledgeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &repository.AgentKnowledgeItem{
		ID:           itemID,
		Title:        req.Title,
		Content:      req.Content,
		Tags:         req.Tags,
		QualityScore: req.QualityScore,
	}

	if err := h.knowledgeRepo.UpdateKnowledgeItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "knowledge item updated successfully"})
}

// DeleteKnowledgeItem 删除知识条目
// @Summary 删除 Agent 知识条目
// @Tags Agent Knowledge
// @Produce json
// @Param id path int true "条目ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/agent-knowledge/items/{id} [delete]
func (h *AgentKnowledgeHandler) DeleteKnowledgeItem(c *gin.Context) {
	itemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item_id"})
		return
	}

	if err := h.knowledgeRepo.DeleteKnowledgeItem(c.Request.Context(), itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "knowledge item deleted successfully"})
}
