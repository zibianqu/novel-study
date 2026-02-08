package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/service"
)

type ChapterHandler struct {
	service *service.ChapterService
}

func NewChapterHandler(service *service.ChapterService) *ChapterHandler {
	return &ChapterHandler{service: service}
}

// CreateChapter 创建章节
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := h.service.CreateChapter(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chapter)
}

// GetChapter 获取章节详情
func (h *ChapterHandler) GetChapter(c *gin.Context) {
	userID := c.GetInt("user_id")
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节 ID"})
		return
	}

	chapter, err := h.service.GetChapter(chapterID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// GetProjectChapters 获取项目的所有章节
func (h *ChapterHandler) GetProjectChapters(c *gin.Context) {
	userID := c.GetInt("user_id")
	projectID, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的项目 ID"})
		return
	}

	chapters, err := h.service.GetProjectChapters(projectID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
}

// UpdateChapter 更新章节
func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	userID := c.GetInt("user_id")
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节 ID"})
		return
	}

	var req model.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := h.service.UpdateChapter(chapterID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// DeleteChapter 删除章节
func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	userID := c.GetInt("user_id")
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节 ID"})
		return
	}

	if err := h.service.DeleteChapter(chapterID, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "章节不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "章节已删除"})
}

// LockChapter 锁定章节
func (h *ChapterHandler) LockChapter(c *gin.Context) {
	userID := c.GetInt("user_id")
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节 ID"})
		return
	}

	if err := h.service.LockChapter(chapterID, userID); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "章节已被锁定"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "章节已锁定"})
}

// UnlockChapter 解锁章节
func (h *ChapterHandler) UnlockChapter(c *gin.Context) {
	userID := c.GetInt("user_id")
	chapterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的章节 ID"})
		return
	}

	if err := h.service.UnlockChapter(chapterID, userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权解锁此章节"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "章节已解锁"})
}
