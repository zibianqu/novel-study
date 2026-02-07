package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novelforge/backend/internal/middleware"
	"github.com/novelforge/backend/internal/model"
	"github.com/novelforge/backend/internal/service"
)

// ==================== 通用辅助 ====================

func ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, model.APIResponse{Code: 0, Message: "success", Data: data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, model.APIResponse{Code: 0, Message: "created", Data: data})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(code, model.APIResponse{Code: code, Message: msg})
}

func paramID(c *gin.Context, name string) (int, bool) {
	id, err := strconv.Atoi(c.Param(name))
	if err != nil {
		fail(c, 400, "无效的ID参数")
		return 0, false
	}
	return id, true
}

// PlaceholderHandler 生成占位Handler（后续阶段实现）
func PlaceholderHandler(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok(c, gin.H{"message": name + " 功能开发中...", "status": "todo"})
	}
}

// ==================== Auth Handler ====================

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		fail(c, 400, err.Error())
		return
	}

	created(c, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		fail(c, 401, err.Error())
		return
	}

	ok(c, resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := middleware.GetUserID(c)
	token, err := h.service.RefreshToken(c.Request.Context(), userID)
	if err != nil {
		fail(c, 500, "刷新Token失败")
		return
	}
	ok(c, gin.H{"token": token})
}

// ==================== Project Handler ====================

type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(s *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var pq model.PaginationQuery
	if err := c.ShouldBindQuery(&pq); err == nil {
		pq.Defaults()
	}

	projects, total, err := h.service.List(c.Request.Context(), userID, pq.Offset(), pq.PageSize)
	if err != nil {
		fail(c, 500, "获取项目列表失败")
		return
	}

	ok(c, gin.H{"items": projects, "total": total, "page": pq.Page, "page_size": pq.PageSize})
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	project, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		fail(c, 500, "创建项目失败: "+err.Error())
		return
	}

	created(c, project)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	project, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "项目不存在")
		return
	}

	ok(c, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	project, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "项目不存在")
		return
	}

	if err := c.ShouldBindJSON(project); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	if err := h.service.Update(c.Request.Context(), project); err != nil {
		fail(c, 500, "更新失败")
		return
	}

	ok(c, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		fail(c, 500, "删除失败")
		return
	}

	ok(c, nil)
}

func (h *ProjectHandler) Export(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	// TODO: 导出完整TXT
	ok(c, gin.H{"project_id": id, "message": "导出功能开发中"})
}

// ==================== Chapter Handler ====================

type ChapterHandler struct {
	service *service.ChapterService
}

func NewChapterHandler(s *service.ChapterService) *ChapterHandler {
	return &ChapterHandler{service: s}
}

func (h *ChapterHandler) ListByProject(c *gin.Context) {
	projectID, valid := paramID(c, "projectId")
	if !valid {
		return
	}

	volumes, chapters, err := h.service.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		fail(c, 500, "获取章节列表失败")
		return
	}

	ok(c, gin.H{"volumes": volumes, "chapters": chapters})
}

func (h *ChapterHandler) Create(c *gin.Context) {
	var req model.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	chapter, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		fail(c, 500, "创建章节失败")
		return
	}

	created(c, chapter)
}

func (h *ChapterHandler) Get(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	chapter, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "章节不存在")
		return
	}

	ok(c, chapter)
}

func (h *ChapterHandler) Update(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req model.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	chapter, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		fail(c, 500, "更新章节失败")
		return
	}

	ok(c, chapter)
}

func (h *ChapterHandler) ListVersions(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	versions, err := h.service.ListVersions(c.Request.Context(), id)
	if err != nil {
		fail(c, 500, "获取版本历史失败")
		return
	}

	ok(c, versions)
}

func (h *ChapterHandler) Rollback(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	var req struct {
		VersionID int `json:"version_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}

	chapter, err := h.service.Rollback(c.Request.Context(), id, req.VersionID)
	if err != nil {
		fail(c, 500, "回滚失败")
		return
	}

	ok(c, chapter)
}

func (h *ChapterHandler) Lock(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.service.Lock(c.Request.Context(), id, userID); err != nil {
		fail(c, 500, "锁定章节失败")
		return
	}

	ok(c, nil)
}

func (h *ChapterHandler) Unlock(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}

	userID := middleware.GetUserID(c)
	if err := h.service.Unlock(c.Request.Context(), id, userID); err != nil {
		fail(c, 500, "解锁章节失败")
		return
	}

	ok(c, nil)
}

// ==================== Agent Handler ====================

type AgentHandler struct {
	service *service.AgentService
}

func NewAgentHandler(s *service.AgentService) *AgentHandler {
	return &AgentHandler{service: s}
}

func (h *AgentHandler) List(c *gin.Context) {
	agents, err := h.service.List(c.Request.Context())
	if err != nil {
		fail(c, 500, "获取Agent列表失败")
		return
	}
	ok(c, agents)
}

func (h *AgentHandler) Get(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	agent, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "Agent不存在")
		return
	}
	ok(c, agent)
}

func (h *AgentHandler) Create(c *gin.Context) {
	var req model.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	agent, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		fail(c, 500, "创建Agent失败: "+err.Error())
		return
	}
	created(c, agent)
}

func (h *AgentHandler) Update(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	agent, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "Agent不存在")
		return
	}

	if err := c.ShouldBindJSON(agent); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	agent.ID = id

	if err := h.service.Update(c.Request.Context(), agent); err != nil {
		fail(c, 500, "更新Agent失败")
		return
	}
	ok(c, agent)
}

func (h *AgentHandler) Delete(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		fail(c, 400, err.Error())
		return
	}
	ok(c, nil)
}

func (h *AgentHandler) Test(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	// TODO: 测试Agent执行
	ok(c, gin.H{"agent_id": id, "message": "Agent测试功能开发中"})
}

// ==================== Knowledge Handler ====================

type KnowledgeHandler struct {
	service *service.KnowledgeService
}

func NewKnowledgeHandler(s *service.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{service: s}
}

func (h *KnowledgeHandler) ListCategories(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Query("agent_id"))
	if agentID == 0 {
		fail(c, 400, "缺少agent_id参数")
		return
	}
	cats, err := h.service.ListCategories(c.Request.Context(), agentID)
	if err != nil {
		fail(c, 500, "获取分类失败")
		return
	}
	ok(c, cats)
}

func (h *KnowledgeHandler) CreateCategory(c *gin.Context) {
	var cat model.KnowledgeCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	if err := h.service.CreateCategory(c.Request.Context(), &cat); err != nil {
		fail(c, 500, "创建分类失败")
		return
	}
	created(c, cat)
}

func (h *KnowledgeHandler) UpdateCategory(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	var cat model.KnowledgeCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	cat.ID = id
	if err := h.service.UpdateCategory(c.Request.Context(), &cat); err != nil {
		fail(c, 500, "更新分类失败")
		return
	}
	ok(c, cat)
}

func (h *KnowledgeHandler) DeleteCategory(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		fail(c, 500, "删除分类失败")
		return
	}
	ok(c, nil)
}

func (h *KnowledgeHandler) ListItems(c *gin.Context) {
	agentID, _ := strconv.Atoi(c.Query("agent_id"))
	if agentID == 0 {
		fail(c, 400, "缺少agent_id参数")
		return
	}
	var catID *int
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.Atoi(cid)
		catID = &v
	}
	var pq model.PaginationQuery
	_ = c.ShouldBindQuery(&pq)
	pq.Defaults()

	items, total, err := h.service.ListItems(c.Request.Context(), agentID, catID, pq.Offset(), pq.PageSize)
	if err != nil {
		fail(c, 500, "获取知识条目失败")
		return
	}
	ok(c, gin.H{"items": items, "total": total})
}

func (h *KnowledgeHandler) CreateItem(c *gin.Context) {
	var req model.CreateKnowledgeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "请求参数错误: "+err.Error())
		return
	}
	userID := middleware.GetUserID(c)
	item := &model.KnowledgeItem{
		AgentID:    req.AgentID,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		Tags:       req.Tags,
		Source:     "manual",
		CreatedBy:  &userID,
	}
	if err := h.service.CreateItem(c.Request.Context(), item); err != nil {
		fail(c, 500, "创建知识条目失败")
		return
	}
	created(c, item)
}

func (h *KnowledgeHandler) GetItem(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	item, err := h.service.GetItem(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "知识条目不存在")
		return
	}
	ok(c, item)
}

func (h *KnowledgeHandler) UpdateItem(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	item, err := h.service.GetItem(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "知识条目不存在")
		return
	}
	if err := c.ShouldBindJSON(item); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	item.ID = id
	if err := h.service.UpdateItem(c.Request.Context(), item); err != nil {
		fail(c, 500, "更新失败")
		return
	}
	ok(c, item)
}

func (h *KnowledgeHandler) DeleteItem(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	if err := h.service.DeleteItem(c.Request.Context(), id); err != nil {
		fail(c, 500, "删除失败")
		return
	}
	ok(c, nil)
}

func (h *KnowledgeHandler) Import(c *gin.Context) {
	// TODO: 批量导入知识
	ok(c, gin.H{"message": "批量导入功能开发中"})
}

func (h *KnowledgeHandler) Search(c *gin.Context) {
	// TODO: 向量检索
	ok(c, gin.H{"message": "向量检索功能开发中"})
}

// ==================== Workflow Handler ====================

type WorkflowHandler struct {
	service *service.WorkflowService
}

func NewWorkflowHandler(s *service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{service: s}
}

func (h *WorkflowHandler) List(c *gin.Context) {
	workflows, err := h.service.List(c.Request.Context())
	if err != nil {
		fail(c, 500, "获取工作流列表失败")
		return
	}
	ok(c, workflows)
}

func (h *WorkflowHandler) Get(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	wf, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "工作流不存在")
		return
	}
	ok(c, wf)
}

func (h *WorkflowHandler) Create(c *gin.Context) {
	var wf model.Workflow
	if err := c.ShouldBindJSON(&wf); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	userID := middleware.GetUserID(c)
	wf.UserID = &userID
	wf.Type = "custom"
	if err := h.service.Create(c.Request.Context(), &wf); err != nil {
		fail(c, 500, "创建工作流失败")
		return
	}
	created(c, wf)
}

func (h *WorkflowHandler) Update(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	wf, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "工作流不存在")
		return
	}
	if err := c.ShouldBindJSON(wf); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	wf.ID = id
	if err := h.service.Update(c.Request.Context(), wf); err != nil {
		fail(c, 500, "更新工作流失败")
		return
	}
	ok(c, wf)
}

func (h *WorkflowHandler) Delete(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		fail(c, 500, "删除失败（系统工作流不可删除）")
		return
	}
	ok(c, nil)
}

func (h *WorkflowHandler) Execute(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	// TODO: 工作流执行引擎（SSE流式输出）
	ok(c, gin.H{"workflow_id": id, "message": "工作流执行引擎开发中"})
}

func (h *WorkflowHandler) GetExecution(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	exec, err := h.service.GetExecution(c.Request.Context(), id)
	if err != nil {
		fail(c, 404, "执行记录不存在")
		return
	}
	ok(c, exec)
}

// ==================== Storyline Handler ====================

type StorylineHandler struct {
	service *service.StorylineService
}

func NewStorylineHandler(s *service.StorylineService) *StorylineHandler {
	return &StorylineHandler{service: s}
}

func (h *StorylineHandler) ListByProject(c *gin.Context) {
	projectID, valid := paramID(c, "projectId")
	if !valid {
		return
	}
	storylines, err := h.service.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		fail(c, 500, "获取三线失败")
		return
	}
	ok(c, storylines)
}

func (h *StorylineHandler) Create(c *gin.Context) {
	var s model.Storyline
	if err := c.ShouldBindJSON(&s); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	if err := h.service.Create(c.Request.Context(), &s); err != nil {
		fail(c, 500, "创建失败")
		return
	}
	created(c, s)
}

func (h *StorylineHandler) Update(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	var s model.Storyline
	if err := c.ShouldBindJSON(&s); err != nil {
		fail(c, 400, "请求参数错误")
		return
	}
	s.ID = id
	if err := h.service.Update(c.Request.Context(), &s); err != nil {
		fail(c, 500, "更新失败")
		return
	}
	ok(c, s)
}

func (h *StorylineHandler) Delete(c *gin.Context) {
	id, valid := paramID(c, "id")
	if !valid {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		fail(c, 500, "删除失败")
		return
	}
	ok(c, nil)
}

func (h *StorylineHandler) Adjust(c *gin.Context) {
	// TODO: 三线调整（联动影响分析）
	ok(c, gin.H{"message": "三线调整功能开发中"})
}
