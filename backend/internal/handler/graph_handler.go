package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/novel-study/internal/repository"
	"github.com/zibianqu/novel-study/internal/service"
)

type GraphHandler struct {
	service *service.GraphService
}

func NewGraphHandler(service *service.GraphService) *GraphHandler {
	return &GraphHandler{service: service}
}

// GetProjectGraph 获取项目知识图谱
func (h *GraphHandler) GetProjectGraph(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	graphData, err := h.service.GetProjectGraph(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graphData)
}

// CreateNode 创建图谱节点
func (h *GraphHandler) CreateNode(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int                    `json:"project_id" binding:"required"`
		ID        string                 `json:"id" binding:"required"`
		Label     string                 `json:"label" binding:"required"`
		Type      string                 `json:"type" binding:"required"`
		Props     map[string]interface{} `json:"properties"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	node := &repository.GraphNode{
		ID:         req.ID,
		Label:      req.Label,
		Type:       req.Type,
		Properties: req.Props,
	}

	if err := h.service.CreateNode(c.Request.Context(), req.ProjectID, userID, node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功"})
}

// CreateRelation 创建图谱关系
func (h *GraphHandler) CreateRelation(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ProjectID int                    `json:"project_id" binding:"required"`
		Source    string                 `json:"source" binding:"required"`
		Target    string                 `json:"target" binding:"required"`
		Type      string                 `json:"type" binding:"required"`
		Props     map[string]interface{} `json:"properties"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rel := &repository.GraphRelation{
		Source:     req.Source,
		Target:     req.Target,
		Type:       req.Type,
		Properties: req.Props,
	}

	if err := h.service.CreateRelation(c.Request.Context(), req.ProjectID, userID, rel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "创建成功"})
}
