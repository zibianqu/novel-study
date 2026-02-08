package service

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/repository"
)

type GraphService struct {
	neo4jRepo   *repository.Neo4jRepository
	projectRepo *repository.ProjectRepository
}

func NewGraphService(
	neo4jRepo *repository.Neo4jRepository,
	projectRepo *repository.ProjectRepository,
) *GraphService {
	return &GraphService{
		neo4jRepo:   neo4jRepo,
		projectRepo: projectRepo,
	}
}

// GetProjectGraph 获取项目知识图谱
func (s *GraphService) GetProjectGraph(ctx context.Context, projectID, userID int) (map[string]interface{}, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, fmt.Errorf("无权访问")
	}

	// 获取图谱数据
	nodes, relations, err := s.neo4jRepo.GetProjectGraph(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"nodes":     nodes,
		"relations": relations,
	}, nil
}

// CreateNode 创建图谱节点
func (s *GraphService) CreateNode(ctx context.Context, projectID, userID int, node *repository.GraphNode) error {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return fmt.Errorf("无权访问")
	}

	return s.neo4jRepo.CreateNode(ctx, projectID, node)
}

// CreateRelation 创建图谱关系
func (s *GraphService) CreateRelation(ctx context.Context, projectID, userID int, rel *repository.GraphRelation) error {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return fmt.Errorf("无权访问")
	}

	return s.neo4jRepo.CreateRelation(ctx, projectID, rel)
}
