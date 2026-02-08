package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zibianqu/novel-study/internal/model"
	"github.com/zibianqu/novel-study/internal/repository"
)

type CacheService struct {
	redis *repository.RedisClient
}

func NewCacheService(redis *repository.RedisClient) *CacheService {
	return &CacheService{redis: redis}
}

func (s *CacheService) GetProject(ctx context.Context, projectID int) (*model.Project, error) {
	key := fmt.Sprintf("project:%d", projectID)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var project model.Project
	if err := json.Unmarshal([]byte(data), &project); err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *CacheService) SetProject(ctx context.Context, project *model.Project) error {
	key := fmt.Sprintf("project:%d", project.ID)
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, 1*time.Hour)
}

func (s *CacheService) InvalidateProject(ctx context.Context, projectID int) error {
	key := fmt.Sprintf("project:%d", projectID)
	return s.redis.Del(ctx, key)
}

func (s *CacheService) GetChapter(ctx context.Context, chapterID int) (*model.Chapter, error) {
	key := fmt.Sprintf("chapter:%d", chapterID)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	var chapter model.Chapter
	if err := json.Unmarshal([]byte(data), &chapter); err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (s *CacheService) SetChapter(ctx context.Context, chapter *model.Chapter) error {
	key := fmt.Sprintf("chapter:%d", chapter.ID)
	data, err := json.Marshal(chapter)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, data, 30*time.Minute)
}

func (s *CacheService) InvalidateChapter(ctx context.Context, chapterID int) error {
	key := fmt.Sprintf("chapter:%d", chapterID)
	return s.redis.Del(ctx, key)
}

func (s *CacheService) InvalidateProjectChapters(ctx context.Context, projectID int) error {
	return nil
}

func (s *CacheService) GetUserSession(ctx context.Context, token string) (int, error) {
	key := fmt.Sprintf("session:%s", token)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return 0, err
	}
	var userID int
	if _, err := fmt.Sscanf(data, "%d", &userID); err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *CacheService) SetUserSession(ctx context.Context, token string, userID int, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return s.redis.Set(ctx, key, userID, expiration)
}

func (s *CacheService) DeleteUserSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return s.redis.Del(ctx, key)
}
