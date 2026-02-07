package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StorylineHandler struct {
	db *sql.DB
}

func NewStorylineHandler(db *sql.DB) *StorylineHandler {
	return &StorylineHandler{db: db}
}

// GetProjectStorylines 获取项目三线
func (h *StorylineHandler) GetProjectStorylines(c *gin.Context) {
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目ID"})
		return
	}

	query := `
		SELECT id, project_id, type, title, description, sequence, status, created_at, updated_at
		FROM storylines
		WHERE project_id = $1
		ORDER BY type, sequence
	`

	rows, err := h.db.Query(query, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var storylines []map[string]interface{}
	for rows.Next() {
		var (
			id          int
			projectID   int
			lineType    string
			title       string
			description sql.NullString
			sequence    int
			status      string
			createdAt   string
			updatedAt   string
		)

		if err := rows.Scan(&id, &projectID, &lineType, &title, &description, &sequence, &status, &createdAt, &updatedAt); err != nil {
			continue
		}

		storylines = append(storylines, map[string]interface{}{
			"id":          id,
			"project_id":  projectID,
			"type":        lineType,
			"title":       title,
			"description": description.String,
			"sequence":    sequence,
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"storylines": storylines})
}

// CreateStoryline 创建故事线
func (h *StorylineHandler) CreateStoryline(c *gin.Context) {
	var req struct {
		ProjectID   int    `json:"project_id" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Sequence    int    `json:"sequence"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO storylines (project_id, type, title, description, sequence, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, 'planning', NOW(), NOW())
		RETURNING id
	`

	var id int
	err := h.db.QueryRow(query, req.ProjectID, req.Type, req.Title, req.Description, req.Sequence).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "创建成功"})
}
