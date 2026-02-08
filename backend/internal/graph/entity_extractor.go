package graph

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// EntityExtractor 实体提取器
type EntityExtractor struct {
	patterns map[NodeType]*EntityPattern
}

// EntityPattern 实体识别模式
type EntityPattern struct {
	Keywords   []string          // 关键词
	Patterns   []*regexp.Regexp  // 正则表达式
	Extractor  func(string) []string // 自定义提取器
}

// ExtractedEntity 提取的实体
type ExtractedEntity struct {
	Type       NodeType
	Name       string
	Context    string
	Confidence float64
	Position   int // 在文本中的位置
	Properties map[string]interface{}
}

// NewEntityExtractor 创建实体提取器
func NewEntityExtractor() *EntityExtractor {
	return &EntityExtractor{
		patterns: make(map[NodeType]*EntityPattern),
	}
}

// Initialize 初始化提取模式
func (e *EntityExtractor) Initialize() {
	// 人物识别模式
	e.patterns[NodeTypeCharacter] = &EntityPattern{
		Keywords: []string{"人", "他", "她", "师傅", "弟子", "道友"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`[\p{Han}]{2,4}(说|道|想|笑|叹)`),
			regexp.MustCompile(`(少年|老者|青年|女子|男子)[\p{Han}]{0,2}`),
		},
	}

	// 地点识别模式
	e.patterns[NodeTypeLocation] = &EntityPattern{
		Keywords: []string{"山", "城", "宗", "殿", "阁", "谷", "峰", "岛"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`[\p{Han}]{2,6}(山|城|宗|殿|阁|谷|峰|岛)`),
			regexp.MustCompile(`(东|西|南|北|中)[\p{Han}]{1,3}(域|州|郡)`),
		},
	}

	// 事件识别模式
	e.patterns[NodeTypeEvent] = &EntityPattern{
		Keywords: []string{"战", "会", "典", "劫", "变"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`[\p{Han}]{2,6}(之战|大战|之会|大会)`),
		},
	}

	// 物品识别模式
	e.patterns[NodeTypeItem] = &EntityPattern{
		Keywords: []string{"剑", "刀", "丹", "符", "宝", "器"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`[\p{Han}]{2,6}(剑|刀|枪|斧|鼎|珠|玉)`),
		},
	}

	// 概念识别模式
	e.patterns[NodeTypeConcept] = &EntityPattern{
		Keywords: []string{"功", "法", "术", "道", "诀"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`[\p{Han}]{2,6}(功|法|术|诀)`),
		},
	}
}

// Extract 从文本提取实体
func (e *EntityExtractor) Extract(
	ctx context.Context,
	text string,
) ([]*ExtractedEntity, error) {
	if len(e.patterns) == 0 {
		e.Initialize()
	}

	entities := make([]*ExtractedEntity, 0)

	// 对每种实体类型进行提取
	for nodeType, pattern := range e.patterns {
		extracted := e.extractByType(text, nodeType, pattern)
		entities = append(entities, extracted...)
	}

	// 去重
	entities = e.deduplicateEntities(entities)

	// 排序（按置信度）
	entities = e.sortByConfidence(entities)

	return entities, nil
}

// extractByType 按类型提取
func (e *EntityExtractor) extractByType(
	text string,
	nodeType NodeType,
	pattern *EntityPattern,
) []*ExtractedEntity {
	entities := make([]*ExtractedEntity, 0)

	// 使用正则表达式提取
	for _, re := range pattern.Patterns {
		matches := re.FindAllStringIndex(text, -1)
		for _, match := range matches {
			name := text[match[0]:match[1]]
			
			// 过滤太短的名称
			if len([]rune(name)) < 2 {
				continue
			}

			entity := &ExtractedEntity{
				Type:       nodeType,
				Name:       name,
				Position:   match[0],
				Confidence: e.calculateConfidence(name, nodeType),
				Properties: make(map[string]interface{}),
			}

			// 提取上下文
			contextStart := max(0, match[0]-20)
			contextEnd := min(len(text), match[1]+20)
			entity.Context = text[contextStart:contextEnd]

			entities = append(entities, entity)
		}
	}

	return entities
}

// calculateConfidence 计算置信度
func (e *EntityExtractor) calculateConfidence(name string, nodeType NodeType) float64 {
	confidence := 0.5 // 基础置信度

	// 根据长度调整
	runes := []rune(name)
	if len(runes) >= 3 && len(runes) <= 5 {
		confidence += 0.2
	}

	// 根据类型特征调整
	pattern := e.patterns[nodeType]
	if pattern != nil {
		for _, keyword := range pattern.Keywords {
			if strings.Contains(name, keyword) {
				confidence += 0.1
				break
			}
		}
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// deduplicateEntities 去重
func (e *EntityExtractor) deduplicateEntities(entities []*ExtractedEntity) []*ExtractedEntity {
	seen := make(map[string]bool)
	result := make([]*ExtractedEntity, 0)

	for _, entity := range entities {
		key := fmt.Sprintf("%s:%s", entity.Type, entity.Name)
		if !seen[key] {
			seen[key] = true
			result = append(result, entity)
		}
	}

	return result
}

// sortByConfidence 按置信度排序
func (e *EntityExtractor) sortByConfidence(entities []*ExtractedEntity) []*ExtractedEntity {
	// 简单冒泡排序
	n := len(entities)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if entities[j].Confidence < entities[j+1].Confidence {
				entities[j], entities[j+1] = entities[j+1], entities[j]
			}
		}
	}
	return entities
}

// ExtractCharacters 提取人物
func (e *EntityExtractor) ExtractCharacters(
	ctx context.Context,
	text string,
) ([]*Character, error) {
	entities, err := e.Extract(ctx, text)
	if err != nil {
		return nil, err
	}

	characters := make([]*Character, 0)
	for _, entity := range entities {
		if entity.Type == NodeTypeCharacter {
			character := &Character{
				Node: Node{
					ID:   generateID("char"),
					Type: NodeTypeCharacter,
					Name: entity.Name,
					Description: entity.Context,
				},
				Role: e.inferCharacterRole(entity),
			}
			characters = append(characters, character)
		}
	}

	return characters, nil
}

// ExtractLocations 提取地点
func (e *EntityExtractor) ExtractLocations(
	ctx context.Context,
	text string,
) ([]*Location, error) {
	entities, err := e.Extract(ctx, text)
	if err != nil {
		return nil, err
	}

	locations := make([]*Location, 0)
	for _, entity := range entities {
		if entity.Type == NodeTypeLocation {
			location := &Location{
				Node: Node{
					ID:   generateID("loc"),
					Type: NodeTypeLocation,
					Name: entity.Name,
					Description: entity.Context,
				},
				LocationType: e.inferLocationType(entity.Name),
			}
			locations = append(locations, location)
		}
	}

	return locations, nil
}

// inferCharacterRole 推断角色类型
func (e *EntityExtractor) inferCharacterRole(entity *ExtractedEntity) string {
	// 简化实现
	if entity.Confidence > 0.8 {
		return "protagonist"
	} else if entity.Confidence > 0.6 {
		return "supporting"
	}
	return "minor"
}

// inferLocationType 推断地点类型
func (e *EntityExtractor) inferLocationType(name string) string {
	if strings.Contains(name, "宗") {
		return "sect"
	} else if strings.Contains(name, "城") {
		return "city"
	} else if strings.Contains(name, "山") {
		return "mountain"
	}
	return "unknown"
}

// 辅助函数

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

var entityIDCounter uint64

func generateID(prefix string) string {
	entityIDCounter++
	return fmt.Sprintf("%s_%d", prefix, entityIDCounter)
}
