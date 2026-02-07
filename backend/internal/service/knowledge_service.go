package service

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/ai/rag"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/repository"
)

type KnowledgeService struct {
	repo        *repository.KnowledgeRepository
	projectRepo *repository.ProjectRepository
	retriever   *rag.Retriever
}

func NewKnowledgeService(
	repo *repository.KnowledgeRepository,
	projectRepo *repository.ProjectRepository,
	retriever *rag.Retriever,
) *KnowledgeService {
	return &KnowledgeService{
		repo:        repo,
		projectRepo: projectRepo,
		retriever:   retriever,
	}
}

func (s *KnowledgeService) CreateKnowledge(ctx context.Context, userID int, req *model.CreateKnowledgeRequest) (*model.KnowledgeBase, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问此项目")
	}

	kb := &model.KnowledgeBase{
		ProjectID: req.ProjectID,
		Title:     req.Title,
		Content:   req.Content,
		Type:      req.Type,
		Tags:      req.Tags,
	}

	if err := s.repo.Create(kb); err != nil {
		return nil, err
	}

	// 异步向量化（实际应用中应使用消息队列）
	go s.vectorizeKnowledge(context.Background(), kb)

	return kb, nil
}

func (s *KnowledgeService) GetKnowledge(id, userID int) (*model.KnowledgeBase, error) {
	kb, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 验证权限
	project, err := s.projectRepo.GetByID(kb.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问")
	}

	return kb, nil
}

func (s *KnowledgeService) GetProjectKnowledge(projectID, userID int) ([]*model.KnowledgeBase, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问")
	}

	return s.repo.GetByProjectID(projectID)
}

func (s *KnowledgeService) SearchKnowledge(ctx context.Context, projectID, userID int, query string, topK int) ([]*rag.Document, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问")
	}

	return s.retriever.Retrieve(ctx, projectID, query, topK)
}

func (s *KnowledgeService) DeleteKnowledge(id, userID int) error {
	kb, err := s.GetKnowledge(id, userID)
	if err != nil {
		return err
	}

	return s.repo.Delete(kb.ID)
}

// vectorizeKnowledge 向量化知识
func (s *KnowledgeService) vectorizeKnowledge(ctx context.Context, kb *model.KnowledgeBase) {
	// TODO: 实现向量化逻辑
	// 1. 生成嵌入
	// 2. 存储到向量数据库
	// 3. 标记为已向量化
}
