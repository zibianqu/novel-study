package service

import (
	"errors"
	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(userID int, req *model.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		UserID:      userID,
		Title:       req.Title,
		Type:        req.Type,
		Genre:       req.Genre,
		Description: req.Description,
		CoverImage:  req.CoverImage,
		Status:      "draft",
		WordCount:   0,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetProject(id, userID int) (*model.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 验证权限
	if project.UserID != userID {
		return nil, errors.New("无权访问此项目")
	}

	return project, nil
}

func (s *ProjectService) GetUserProjects(userID int) ([]*model.Project, error) {
	return s.repo.GetByUserID(userID)
}

func (s *ProjectService) UpdateProject(id, userID int, req *model.UpdateProjectRequest) (*model.Project, error) {
	// 首先获取项目并验证权限
	project, err := s.GetProject(id, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != "" {
		project.Title = req.Title
	}
	if req.Genre != "" {
		project.Genre = req.Genre
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.CoverImage != "" {
		project.CoverImage = req.CoverImage
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(id, userID int) error {
	return s.repo.Delete(id, userID)
}
