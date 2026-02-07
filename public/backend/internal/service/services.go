package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/novelforge/backend/internal/config"
	"github.com/novelforge/backend/internal/middleware"
	"github.com/novelforge/backend/internal/model"
	"github.com/novelforge/backend/internal/repository"
)

// ==================== Auth Service ====================

type AuthService struct {
	userRepo *repository.UserRepository
	jwtCfg   config.JWTConfig
}

func NewAuthService(repo *repository.UserRepository, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{userRepo: repo, jwtCfg: jwtCfg}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.AuthResponse, error) {
	// 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 检查用户名是否已存在
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("用户名已被占用")
	}

	// 哈希密码
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 生成 Token
	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtCfg)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("邮箱或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("邮箱或密码错误")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtCfg)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID int) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	return middleware.GenerateToken(user.ID, user.Username, s.jwtCfg)
}

// ==================== Project Service ====================

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, userID int, req model.CreateProjectRequest) (*model.Project, error) {
	p := &model.Project{
		UserID:      userID,
		Title:       req.Title,
		Type:        req.Type,
		Genre:       req.Genre,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) Get(ctx context.Context, id int) (*model.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) List(ctx context.Context, userID, offset, limit int) ([]model.Project, int, error) {
	return s.repo.ListByUser(ctx, userID, offset, limit)
}

func (s *ProjectService) Update(ctx context.Context, p *model.Project) error {
	return s.repo.Update(ctx, p)
}

func (s *ProjectService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// ==================== Chapter Service ====================

type ChapterService struct {
	repo *repository.ChapterRepository
}

func NewChapterService(repo *repository.ChapterRepository) *ChapterService {
	return &ChapterService{repo: repo}
}

func (s *ChapterService) Create(ctx context.Context, req model.CreateChapterRequest) (*model.Chapter, error) {
	ch := &model.Chapter{
		ProjectID: req.ProjectID,
		VolumeID:  req.VolumeID,
		Title:     req.Title,
		Content:   req.Content,
	}
	if err := s.repo.Create(ctx, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *ChapterService) Get(ctx context.Context, id int) (*model.Chapter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ChapterService) ListByProject(ctx context.Context, projectID int) ([]model.Volume, []model.Chapter, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *ChapterService) Update(ctx context.Context, id int, req model.UpdateChapterRequest) (*model.Chapter, error) {
	ch, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		ch.Title = *req.Title
	}
	if req.Content != nil {
		ch.Content = *req.Content
	}
	if req.Status != nil {
		ch.Status = *req.Status
	}

	if err := s.repo.Update(ctx, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

func (s *ChapterService) Lock(ctx context.Context, chapterID, userID int) error {
	return s.repo.Lock(ctx, chapterID, userID)
}

func (s *ChapterService) Unlock(ctx context.Context, chapterID, userID int) error {
	return s.repo.Unlock(ctx, chapterID, userID)
}

func (s *ChapterService) ListVersions(ctx context.Context, chapterID int) ([]model.ChapterVersion, error) {
	return s.repo.ListVersions(ctx, chapterID)
}

func (s *ChapterService) Rollback(ctx context.Context, chapterID, versionID int) (*model.Chapter, error) {
	version, err := s.repo.GetVersion(ctx, versionID)
	if err != nil {
		return nil, err
	}

	ch, err := s.repo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	ch.Content = version.Content
	if err := s.repo.Update(ctx, ch); err != nil {
		return nil, err
	}
	return ch, nil
}

// ==================== Agent Service ====================

type AgentService struct {
	repo *repository.AgentRepository
}

func NewAgentService(repo *repository.AgentRepository) *AgentService {
	return &AgentService{repo: repo}
}

func (s *AgentService) List(ctx context.Context) ([]model.Agent, error) {
	return s.repo.List(ctx)
}

func (s *AgentService) Get(ctx context.Context, id int) (*model.Agent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AgentService) Create(ctx context.Context, userID int, req model.CreateAgentRequest) (*model.Agent, error) {
	a := &model.Agent{
		UserID:       &userID,
		AgentKey:     req.AgentKey,
		Name:         req.Name,
		Icon:         req.Icon,
		Description:  req.Description,
		Type:         "extension",
		Layer:        req.Layer,
		SystemPrompt: req.SystemPrompt,
		Model:        req.Model,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		IsActive:     true,
	}
	if a.Model == "" {
		a.Model = "gpt-4o"
	}
	if a.MaxTokens == 0 {
		a.MaxTokens = 4096
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AgentService) Update(ctx context.Context, a *model.Agent) error {
	return s.repo.Update(ctx, a)
}

func (s *AgentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// ==================== Knowledge Service ====================

type KnowledgeService struct {
	repo *repository.KnowledgeRepository
}

func NewKnowledgeService(repo *repository.KnowledgeRepository) *KnowledgeService {
	return &KnowledgeService{repo: repo}
}

func (s *KnowledgeService) ListCategories(ctx context.Context, agentID int) ([]model.KnowledgeCategory, error) {
	return s.repo.ListCategories(ctx, agentID)
}

func (s *KnowledgeService) CreateCategory(ctx context.Context, c *model.KnowledgeCategory) error {
	return s.repo.CreateCategory(ctx, c)
}

func (s *KnowledgeService) UpdateCategory(ctx context.Context, c *model.KnowledgeCategory) error {
	return s.repo.UpdateCategory(ctx, c)
}

func (s *KnowledgeService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *KnowledgeService) ListItems(ctx context.Context, agentID int, categoryID *int, offset, limit int) ([]model.KnowledgeItem, int, error) {
	return s.repo.ListItems(ctx, agentID, categoryID, offset, limit)
}

func (s *KnowledgeService) CreateItem(ctx context.Context, item *model.KnowledgeItem) error {
	return s.repo.CreateItem(ctx, item)
}

func (s *KnowledgeService) GetItem(ctx context.Context, id int) (*model.KnowledgeItem, error) {
	return s.repo.GetItem(ctx, id)
}

func (s *KnowledgeService) UpdateItem(ctx context.Context, item *model.KnowledgeItem) error {
	return s.repo.UpdateItem(ctx, item)
}

func (s *KnowledgeService) DeleteItem(ctx context.Context, id int) error {
	return s.repo.DeleteItem(ctx, id)
}

// ==================== Workflow Service ====================

type WorkflowService struct {
	repo *repository.WorkflowRepository
}

func NewWorkflowService(repo *repository.WorkflowRepository) *WorkflowService {
	return &WorkflowService{repo: repo}
}

func (s *WorkflowService) List(ctx context.Context) ([]model.Workflow, error) {
	return s.repo.List(ctx)
}

func (s *WorkflowService) Get(ctx context.Context, id int) (*model.Workflow, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *WorkflowService) Create(ctx context.Context, w *model.Workflow) error {
	return s.repo.Create(ctx, w)
}

func (s *WorkflowService) Update(ctx context.Context, w *model.Workflow) error {
	return s.repo.Update(ctx, w)
}

func (s *WorkflowService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *WorkflowService) GetExecution(ctx context.Context, id int) (*model.WorkflowExecution, error) {
	return s.repo.GetExecution(ctx, id)
}

// ==================== Storyline Service ====================

type StorylineService struct {
	repo *repository.StorylineRepository
}

func NewStorylineService(repo *repository.StorylineRepository) *StorylineService {
	return &StorylineService{repo: repo}
}

func (s *StorylineService) ListByProject(ctx context.Context, projectID int) ([]model.Storyline, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *StorylineService) Create(ctx context.Context, storyline *model.Storyline) error {
	return s.repo.Create(ctx, storyline)
}

func (s *StorylineService) Update(ctx context.Context, storyline *model.Storyline) error {
	return s.repo.Update(ctx, storyline)
}

func (s *StorylineService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
