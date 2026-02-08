package service

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/repository"
)

type ChapterService struct {
	repo        *repository.ChapterRepository
	projectRepo *repository.ProjectRepository
}

func NewChapterService(repo *repository.ChapterRepository, projectRepo *repository.ProjectRepository) *ChapterService {
	return &ChapterService{repo: repo, projectRepo: projectRepo}
}

func (s *ChapterService) CreateChapter(userID int, req *model.CreateChapterRequest) (*model.Chapter, error) {
	// 验证项目权限
	project, err := s.projectRepo.GetByID(req.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权在此项目下创建章节")
	}

	// 计算字数
	wordCount := countWords(req.Content)

	chapter := &model.Chapter{
		ProjectID: req.ProjectID,
		VolumeID:  req.VolumeID,
		Title:     req.Title,
		Content:   req.Content,
		WordCount: wordCount,
		SortOrder: req.SortOrder,
		Status:    "draft",
	}

	if err := s.repo.Create(chapter); err != nil {
		return nil, err
	}

	// 更新项目总字数
	s.updateProjectWordCount(req.ProjectID)

	return chapter, nil
}

func (s *ChapterService) GetChapter(id, userID int) (*model.Chapter, error) {
	chapter, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 验证权限
	project, err := s.projectRepo.GetByID(chapter.ProjectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问此章节")
	}

	return chapter, nil
}

func (s *ChapterService) GetProjectChapters(projectID, userID int) ([]*model.Chapter, error) {
	// 验证权限
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, errors.New("无权访问此项目")
	}

	return s.repo.GetByProjectID(projectID)
}

func (s *ChapterService) UpdateChapter(id, userID int, req *model.UpdateChapterRequest) (*model.Chapter, error) {
	// 首先获取章节并验证权限
	chapter, err := s.GetChapter(id, userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != "" {
		chapter.Title = req.Title
	}
	if req.Content != "" {
		chapter.Content = req.Content
		chapter.WordCount = countWords(req.Content)
	}
	if req.Status != "" {
		chapter.Status = req.Status
	}
	if req.SortOrder > 0 {
		chapter.SortOrder = req.SortOrder
	}

	if err := s.repo.Update(chapter); err != nil {
		return nil, err
	}

	// 更新项目总字数
	s.updateProjectWordCount(chapter.ProjectID)

	return chapter, nil
}

func (s *ChapterService) DeleteChapter(id, userID int) error {
	// 验证权限
	chapter, err := s.GetChapter(id, userID)
	if err != nil {
		return err
	}

	projectID := chapter.ProjectID

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 更新项目总字数
	s.updateProjectWordCount(projectID)

	return nil
}

func (s *ChapterService) LockChapter(chapterID, userID int) error {
	return s.repo.Lock(chapterID, userID)
}

func (s *ChapterService) UnlockChapter(chapterID, userID int) error {
	return s.repo.Unlock(chapterID, userID)
}

// updateProjectWordCount 更新项目总字数
func (s *ChapterService) updateProjectWordCount(projectID int) {
	chapters, err := s.repo.GetByProjectID(projectID)
	if err != nil {
		return
	}

	totalWords := 0
	for _, chapter := range chapters {
		totalWords += chapter.WordCount
	}

	s.projectRepo.UpdateWordCount(projectID, totalWords)
}

// countWords 计算中文字数（简单版）
func countWords(text string) int {
	// 移除空白字符
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	// 计算 UTF-8 字符数（包括汉字、标点符号）
	// 这里使用简单算法，实际可以更精细
	return utf8.RuneCountInString(text)
}
