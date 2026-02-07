package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/service"
)

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(service *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// CreateProject 创建项目
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.CreateProject(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建项目失败"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProject 获取项目详情
func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	project, err := h.service.GetProject(projectID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetProjects 获取用户的所有项目
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID := c.GetInt("user_id")

	projects, err := h.service.GetUserProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取项目列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// UpdateProject 更新项目
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	var req model.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.UpdateProject(projectID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// DeleteProject 删除项目
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	if err := h.service.DeleteProject(projectID, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "项目已删除"})
}
