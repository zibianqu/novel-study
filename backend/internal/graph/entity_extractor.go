package graph

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Entity 实体
type Entity struct {
	ID          string
	Name        string
	Type        NodeType
	Mentions    int     // 提及次数
	Confidence  float64 // 置信度
	Context     string  // 上下文
	Chapter     int     // 首次出现章节
	Attributes  map[string]string
}

// EntityExtractor 实体提取器
type EntityExtractor struct {
	characterPatterns []*regexp.Regexp
	locationPatterns  []*regexp.Regexp
	eventPatterns     []*regexp.Regexp
	itemPatterns      []*regexp.Regexp
	conceptPatterns   []*regexp.Regexp
}

// NewEntityExtractor 创建实体提取器
func NewEntityExtractor() *EntityExtractor {
	return &EntityExtractor{
		characterPatterns: compileCharacterPatterns(),
		locationPatterns:  compileLocationPatterns(),
		eventPatterns:     compileEventPatterns(),
		itemPatterns:      compileItemPatterns(),
		conceptPatterns:   compileConceptPatterns(),
	}
}

// Extract 从文本中提取实体
func (e *EntityExtractor) Extract(ctx context.Context, text string, chapter int) ([]*Entity, error) {
	entities := make([]*Entity, 0)

	// 1. 提取人物
	characters := e.extractCharacters(text, chapter)
	entities = append(entities, characters...)

	// 2. 提取地点
	locations := e.extractLocations(text, chapter)
	entities = append(entities, locations...)

	// 3. 提取事件
	events := e.extractEvents(text, chapter)
	entities = append(entities, events...)

	// 4. 提取物品
	items := e.extractItems(text, chapter)
	entities = append(entities, items...)

	// 5. 提取概念
	concepts := e.extractConcepts(text, chapter)
	entities = append(entities, concepts...)

	// 去重和合并
	entities = e.deduplicateEntities(entities)

	return entities, nil
}

// extractCharacters 提取人物
func (e *EntityExtractor) extractCharacters(text string, chapter int) []*Entity {
	entities := make([]*Entity, 0)

	// 关键词匹配
	characterKeywords := []string{
		"主角", "男主", "女主", "师父", "师兄", "师姐", "弟子",
		"长老", "宗主", "掌门", "前辈", "道友",
	}

	for _, keyword := range characterKeywords {
		if strings.Contains(text, keyword) {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeCharacter, keyword),
				Name:       keyword,
				Type:       NodeTypeCharacter,
				Mentions:   strings.Count(text, keyword),
				Confidence: 0.7,
				Chapter:    chapter,
				Attributes: make(map[string]string),
			}
			entities = append(entities, entity)
		}
	}

	// 人名模式匹配（简化版）
	// 匹配类似 "张三", "李四" 等两字或三字人名
	namePattern := regexp.MustCompile(`[\p{Han}]{2,3}(?:说道|说|道|想|笑|怒|叹)`)
	matches := namePattern.FindAllString(text, -1)

	for _, match := range matches {
		name := strings.TrimRight(match, "说道想笑怒叹")
		if len(name) >= 2 && len(name) <= 3 {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeCharacter, name),
				Name:       name,
				Type:       NodeTypeCharacter,
				Mentions:   strings.Count(text, name),
				Confidence: 0.8,
				Chapter:    chapter,
				Attributes: make(map[string]string),
			}
			entities = append(entities, entity)
		}
	}

	return entities
}

// extractLocations 提取地点
func (e *EntityExtractor) extractLocations(text string, chapter int) []*Entity {
	entities := make([]*Entity, 0)

	locationKeywords := []string{
		"宗门", "门派", "山峰", "大殿", "洞府", "密室",
		"城市", "村庄", "森林", "山脉", "江湖",
	}

	for _, keyword := range locationKeywords {
		if strings.Contains(text, keyword) {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeLocation, keyword),
				Name:       keyword,
				Type:       NodeTypeLocation,
				Mentions:   strings.Count(text, keyword),
				Confidence: 0.7,
				Chapter:    chapter,
				Attributes: make(map[string]string),
			}
			entities = append(entities, entity)
		}
	}

	// 地点模式匹配
	// 匹配 "XXX山", "XXX宗", "XXX城" 等
	locPattern := regexp.MustCompile(`[\p{Han}]{2,5}[山宗城峰谷洞府殿]`)
	matches := locPattern.FindAllString(text, -1)

	for _, match := range matches {
		entity := &Entity{
			ID:         generateEntityID(NodeTypeLocation, match),
			Name:       match,
			Type:       NodeTypeLocation,
			Mentions:   strings.Count(text, match),
			Confidence: 0.75,
			Chapter:    chapter,
			Attributes: make(map[string]string),
		}
		entities = append(entities, entity)
	}

	return entities
}

// extractEvents 提取事件
func (e *EntityExtractor) extractEvents(text string, chapter int) []*Entity {
	entities := make([]*Entity, 0)

	eventKeywords := []string{
		"战斗", "比武", "决斗", "突破", "修炼",
		"会议", "宴会", "拜师", "结盟",
	}

	for _, keyword := range eventKeywords {
		if strings.Contains(text, keyword) {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeEvent, keyword),
				Name:       keyword,
				Type:       NodeTypeEvent,
				Mentions:   strings.Count(text, keyword),
				Confidence: 0.6,
				Chapter:    chapter,
				Attributes: map[string]string{
					"event_type": keyword,
				},
			}
			entities = append(entities, entity)
		}
	}

	return entities
}

// extractItems 提取物品
func (e *EntityExtractor) extractItems(text string, chapter int) []*Entity {
	entities := make([]*Entity, 0)

	itemKeywords := []string{
		"宝剑", "法宝", "丹药", "灵石", "功法",
		"宝物", "神器", "令牌", "玉佩",
	}

	for _, keyword := range itemKeywords {
		if strings.Contains(text, keyword) {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeItem, keyword),
				Name:       keyword,
				Type:       NodeTypeItem,
				Mentions:   strings.Count(text, keyword),
				Confidence: 0.65,
				Chapter:    chapter,
				Attributes: make(map[string]string),
			}
			entities = append(entities, entity)
		}
	}

	return entities
}

// extractConcepts 提取概念
func (e *EntityExtractor) extractConcepts(text string, chapter int) []*Entity {
	entities := make([]*Entity, 0)

	conceptKeywords := []string{
		"剑法", "心法", "武学", "功夫", "阵法",
		"境界", "修为", "神通", "术法",
	}

	for _, keyword := range conceptKeywords {
		if strings.Contains(text, keyword) {
			entity := &Entity{
				ID:         generateEntityID(NodeTypeConcept, keyword),
				Name:       keyword,
				Type:       NodeTypeConcept,
				Mentions:   strings.Count(text, keyword),
				Confidence: 0.6,
				Chapter:    chapter,
				Attributes: make(map[string]string),
			}
			entities = append(entities, entity)
		}
	}

	return entities
}

// deduplicateEntities 去重实体
func (e *EntityExtractor) deduplicateEntities(entities []*Entity) []*Entity {
	seen := make(map[string]*Entity)

	for _, entity := range entities {
		key := string(entity.Type) + ":" + entity.Name
		if existing, ok := seen[key]; ok {
			// 合并提及次数
			existing.Mentions += entity.Mentions
			// 提高置信度
			if existing.Confidence < 0.9 {
				existing.Confidence += 0.1
			}
		} else {
			seen[key] = entity
		}
	}

	result := make([]*Entity, 0, len(seen))
	for _, entity := range seen {
		result = append(result, entity)
	}

	return result
}

// 编译正则表达式模式

func compileCharacterPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`[\p{Han}]{2,3}(?:说道|说|道|想|笑|怒|叹)`),
		regexp.MustCompile(`(?:师父|师兄|师姐|师弟|师妹)`),
	}
}

func compileLocationPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`[\p{Han}]{2,5}[山宗城峰谷洞府殿]`),
	}
}

func compileEventPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`(?:战斗|比武|决斗)`),
	}
}

func compileItemPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`[\p{Han}]{2,4}(?:剑|刀|枪|法宝|丹药)`),
	}
}

func compileConceptPatterns() []*regexp.Regexp {
	return []*regexp.Regexp{
		regexp.MustCompile(`[\p{Han}]{2,4}(?:心法|剑法|功法|神通)`),
	}
}

// generateEntityID 生成实体ID
func generateEntityID(nodeType NodeType, name string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(string(nodeType)), name)
}
