package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zibianqu/novel-study/internal/repository"
)

type CacheService struct {
	redis *repository.RedisClient
}

func NewCacheService(redis *repository.RedisClient) *CacheService {
	return &CacheService{redis: redis}
}

// ====================================
// 项目缓存
// ====================================

// SetProject 缓存项目数据
func (s *CacheService) SetProject(ctx context.Context, projectID int, data interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("project:%d", projectID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, jsonData, ttl)
}

// GetProject 获取项目缓存
func (s *CacheService) GetProject(ctx context.Context, projectID int, result interface{}) error {
	key := fmt.Sprintf("project:%d", projectID)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// DeleteProject 删除项目缓存
func (s *CacheService) DeleteProject(ctx context.Context, projectID int) error {
	key := fmt.Sprintf("project:%d", projectID)
	return s.redis.Delete(ctx, key)
}

// ====================================
// 章节缓存
// ====================================

// SetChapter 缓存章节数据
func (s *CacheService) SetChapter(ctx context.Context, chapterID int, data interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("chapter:%d", chapterID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, jsonData, ttl)
}

// GetChapter 获取章节缓存
func (s *CacheService) GetChapter(ctx context.Context, chapterID int, result interface{}) error {
	key := fmt.Sprintf("chapter:%d", chapterID)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// DeleteChapter 删除章节缓存
func (s *CacheService) DeleteChapter(ctx context.Context, chapterID int) error {
	key := fmt.Sprintf("chapter:%d", chapterID)
	return s.redis.Delete(ctx, key)
}

// ====================================
// 用户会话缓存
// ====================================

// SetUserSession 缓存用户会话
func (s *CacheService) SetUserSession(ctx context.Context, userID int, sessionData interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("session:user:%d", userID)
	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, jsonData, ttl)
}

// GetUserSession 获取用户会话
func (s *CacheService) GetUserSession(ctx context.Context, userID int, result interface{}) error {
	key := fmt.Sprintf("session:user:%d", userID)
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// DeleteUserSession 删除用户会话
func (s *CacheService) DeleteUserSession(ctx context.Context, userID int) error {
	key := fmt.Sprintf("session:user:%d", userID)
	return s.redis.Delete(ctx, key)
}

// ====================================
// 通用缓存
// ====================================

// Set 设置任意缓存
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.redis.Set(ctx, key, jsonData, ttl)
}

// Get 获取任意缓存
func (s *CacheService) Get(ctx context.Context, key string, result interface{}) error {
	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), result)
}

// Delete 删除任意缓存
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
	return s.redis.Delete(ctx, keys...)
}

// ====================================
// 清理方法
// ====================================

// InvalidateProjectCache 失效项目相关的所有缓存
func (s *CacheService) InvalidateProjectCache(ctx context.Context, projectID int) error {
	keys := []string{
		fmt.Sprintf("project:%d", projectID),
	}
	
	// 查找并删除相关章节缓存
	chapterKeys, err := s.redis.Keys(ctx, fmt.Sprintf("chapter:*:project:%d", projectID))
	if err == nil {
		keys = append(keys, chapterKeys...)
	}
	
	if len(keys) > 0 {
		return s.redis.Delete(ctx, keys...)
	}
	return nil
}
