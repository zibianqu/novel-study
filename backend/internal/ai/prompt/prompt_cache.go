package prompt

import (
	"sync"
	"time"
)

// CacheEntry 缓存条目
type CacheEntry struct {
	Content   string
	Tokens    int
	Timestamp time.Time
	HitCount  int
}

// PromptCache Prompt 缓存
type PromptCache struct {
	cache      map[string]*CacheEntry
	mu         sync.RWMutex
	maxSize    int
	ttl        time.Duration
}

// NewPromptCache 创建缓存
func NewPromptCache(maxSize int, ttl time.Duration) *PromptCache {
	cache := &PromptCache{
		cache:   make(map[string]*CacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// 启动清理程序
	go cache.cleanupLoop()

	return cache
}

// Get 获取缓存
func (pc *PromptCache) Get(key string) (string, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	entry, ok := pc.cache[key]
	if !ok {
		return "", false
	}

	// 检查是否过期
	if time.Since(entry.Timestamp) > pc.ttl {
		return "", false
	}

	// 增加命中次数
	entry.HitCount++

	return entry.Content, true
}

// Set 设置缓存
func (pc *PromptCache) Set(key string, content string, tokens int) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	// 如果超过容量，删除最旧的条目
	if len(pc.cache) >= pc.maxSize {
		pc.evictOldest()
	}

	pc.cache[key] = &CacheEntry{
		Content:   content,
		Tokens:    tokens,
		Timestamp: time.Now(),
		HitCount:  0,
	}
}

// Delete 删除缓存
func (pc *PromptCache) Delete(key string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.cache, key)
}

// Clear 清空缓存
func (pc *PromptCache) Clear() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.cache = make(map[string]*CacheEntry)
}

// GetStats 获取统计信息
func (pc *PromptCache) GetStats() map[string]interface{} {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	totalHits := 0
	totalTokens := 0
	for _, entry := range pc.cache {
		totalHits += entry.HitCount
		totalTokens += entry.Tokens
	}

	return map[string]interface{}{
		"size":        len(pc.cache),
		"max_size":    pc.maxSize,
		"total_hits":  totalHits,
		"total_tokens": totalTokens,
		"ttl_seconds": pc.ttl.Seconds(),
	}
}

// evictOldest 删除最旧的条目
func (pc *PromptCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range pc.cache {
		if first || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
			first = false
		}
	}

	if oldestKey != "" {
		delete(pc.cache, oldestKey)
	}
}

// cleanupLoop 清理循环
func (pc *PromptCache) cleanupLoop() {
	ticker := time.NewTicker(pc.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		pc.cleanup()
	}
}

// cleanup 清理过期条目
func (pc *PromptCache) cleanup() {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for key, entry := range pc.cache {
		if now.Sub(entry.Timestamp) > pc.ttl {
			delete(pc.cache, key)
		}
	}
}

// ContextCache 上下文缓存管理器
type ContextCache struct {
	projectCache   *PromptCache
	characterCache *PromptCache
	knowledgeCache *PromptCache
}

// NewContextCache 创建上下文缓存管理器
func NewContextCache() *ContextCache {
	return &ContextCache{
		projectCache:   NewPromptCache(100, 30*time.Minute),
		characterCache: NewPromptCache(200, 1*time.Hour),
		knowledgeCache: NewPromptCache(500, 2*time.Hour),
	}
}

// GetProjectContext 获取项目上下文
func (cc *ContextCache) GetProjectContext(projectID int) (string, bool) {
	key := formatProjectKey(projectID)
	return cc.projectCache.Get(key)
}

// SetProjectContext 设置项目上下文
func (cc *ContextCache) SetProjectContext(projectID int, content string, tokens int) {
	key := formatProjectKey(projectID)
	cc.projectCache.Set(key, content, tokens)
}

// GetCharacterContext 获取角色上下文
func (cc *ContextCache) GetCharacterContext(projectID int, characterID int) (string, bool) {
	key := formatCharacterKey(projectID, characterID)
	return cc.characterCache.Get(key)
}

// SetCharacterContext 设置角色上下文
func (cc *ContextCache) SetCharacterContext(projectID int, characterID int, content string, tokens int) {
	key := formatCharacterKey(projectID, characterID)
	cc.characterCache.Set(key, content, tokens)
}

// GetKnowledgeContent 获取知识内容
func (cc *ContextCache) GetKnowledgeContent(agentID int, category string) (string, bool) {
	key := formatKnowledgeKey(agentID, category)
	return cc.knowledgeCache.Get(key)
}

// SetKnowledgeContent 设置知识内容
func (cc *ContextCache) SetKnowledgeContent(agentID int, category string, content string, tokens int) {
	key := formatKnowledgeKey(agentID, category)
	cc.knowledgeCache.Set(key, content, tokens)
}

// InvalidateProject 失效项目缓存
func (cc *ContextCache) InvalidateProject(projectID int) {
	key := formatProjectKey(projectID)
	cc.projectCache.Delete(key)
}

// GetAllStats 获取所有统计信息
func (cc *ContextCache) GetAllStats() map[string]interface{} {
	return map[string]interface{}{
		"project":   cc.projectCache.GetStats(),
		"character": cc.characterCache.GetStats(),
		"knowledge": cc.knowledgeCache.GetStats(),
	}
}

// 辅助函数

func formatProjectKey(projectID int) string {
	return fmt.Sprintf("project:%d", projectID)
}

func formatCharacterKey(projectID int, characterID int) string {
	return fmt.Sprintf("character:%d:%d", projectID, characterID)
}

func formatKnowledgeKey(agentID int, category string) string {
	return fmt.Sprintf("knowledge:%d:%s", agentID, category)
}

import "fmt"
