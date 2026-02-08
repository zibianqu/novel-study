package graph

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// RelationExtractor 关系提取器
type RelationExtractor struct {
	patterns map[RelationType]*RelationPattern
}

// RelationPattern 关系识别模式
type RelationPattern struct {
	Keywords []string
	Patterns []*regexp.Regexp
	Weight   float64
}

// ExtractedRelation 提取的关系
type ExtractedRelation struct {
	Type       RelationType
	FromEntity string
	ToEntity   string
	Context    string
	Confidence float64
	Properties map[string]interface{}
}

// NewRelationExtractor 创建关系提取器
func NewRelationExtractor() *RelationExtractor {
	return &RelationExtractor{
		patterns: make(map[RelationType]*RelationPattern),
	}
}

// Initialize 初始化关系模式
func (r *RelationExtractor) Initialize() {
	// 认识关系
	r.patterns[RelationKnows] = &RelationPattern{
		Keywords: []string{"认识", "见过", "听说"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)(认识|见过)(.+)`),
		},
		Weight: 0.5,
	}

	// 师徒关系
	r.patterns[RelationMasterOf] = &RelationPattern{
		Keywords: []string{"师傅", "弟子", "师父", "徒弟"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)是(.+)的师傅`),
			regexp.MustCompile(`(.+)拜(.+)为师`),
		},
		Weight: 0.9,
	}

	// 亲属关系
	r.patterns[RelationFamilyOf] = &RelationPattern{
		Keywords: []string{"父亲", "母亲", "兄弟", "姐妹"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)是(.+)的(父亲|母亲|兄弟|姐妹)`),
		},
		Weight: 1.0,
	}

	// 仇敵关系
	r.patterns[RelationEnemyOf] = &RelationPattern{
		Keywords: []string{"仇人", "敌人", "对手"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)和(.+)(成为仇敵|是仇人)`),
		},
		Weight: 0.8,
	}

	// 位置关系
	r.patterns[RelationLocatedAt] = &RelationPattern{
		Keywords: []string{"在", "位于", "处于"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)(在|位于)(.+)`),
		},
		Weight: 0.6,
	}

	// 拥有关系
	r.patterns[RelationOwns] = &RelationPattern{
		Keywords: []string{"拥有", "持有", "得到"},
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(`(.+)(拥有|持有|得到)(.+)`),
		},
		Weight: 0.7,
	}
}

// Extract 提取关系
func (r *RelationExtractor) Extract(
	ctx context.Context,
	text string,
	entities []*ExtractedEntity,
) ([]*ExtractedRelation, error) {
	if len(r.patterns) == 0 {
		r.Initialize()
	}

	relations := make([]*ExtractedRelation, 0)

	// 使用模式匹配
	for relType, pattern := range r.patterns {
		extracted := r.extractByPattern(text, entities, relType, pattern)
		relations = append(relations, extracted...)
	}

	// 使用共现分析
	coOccurrence := r.extractByCoOccurrence(text, entities)
	relations = append(relations, coOccurrence...)

	// 去重
	relations = r.deduplicateRelations(relations)

	return relations, nil
}

// extractByPattern 基于模式提取
func (r *RelationExtractor) extractByPattern(
	text string,
	entities []*ExtractedEntity,
	relType RelationType,
	pattern *RelationPattern,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	// 使用正则表达式匹配
	for _, re := range pattern.Patterns {
		matches := re.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			// 提取实体名
			fromName := strings.TrimSpace(match[1])
			toName := strings.TrimSpace(match[len(match)-1])

			if fromName == "" || toName == "" {
				continue
			}

			relation := &ExtractedRelation{
				Type:       relType,
				FromEntity: fromName,
				ToEntity:   toName,
				Context:    match[0],
				Confidence: pattern.Weight,
				Properties: make(map[string]interface{}),
			}

			relations = append(relations, relation)
		}
	}

	return relations
}

// extractByCoOccurrence 基于共现提取
func (r *RelationExtractor) extractByCoOccurrence(
	text string,
	entities []*ExtractedEntity,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	// 对每对实体检查共现
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			entity1 := entities[i]
			entity2 := entities[j]

			// 计算距离
			distance := abs(entity1.Position - entity2.Position)

			// 如果距离很近（在100个字符内）
			if distance < 100 {
				relation := &ExtractedRelation{
					Type:       RelationKnows, // 默认为认识关系
					FromEntity: entity1.Name,
					ToEntity:   entity2.Name,
					Confidence: 0.3, // 低置信度
					Properties: map[string]interface{}{
						"distance": distance,
					},
				}
				relations = append(relations, relation)
			}
		}
	}

	return relations
}

// deduplicateRelations 去重
func (r *RelationExtractor) deduplicateRelations(
	relations []*ExtractedRelation,
) []*ExtractedRelation {
	seen := make(map[string]*ExtractedRelation)

	for _, rel := range relations {
		key := fmt.Sprintf("%s:%s->%s", rel.Type, rel.FromEntity, rel.ToEntity)

		if existing, ok := seen[key]; ok {
			// 保留置信度更高的
			if rel.Confidence > existing.Confidence {
				seen[key] = rel
			}
		} else {
			seen[key] = rel
		}
	}

	result := make([]*ExtractedRelation, 0, len(seen))
	for _, rel := range seen {
		result = append(result, rel)
	}

	return result
}

// BuildRelationships 构建关系对象
func (r *RelationExtractor) BuildRelationships(
	extractedRelations []*ExtractedRelation,
	entityMap map[string]string, // name -> id 映射
) []*Relationship {
	relationships := make([]*Relationship, 0)

	for _, extracted := range extractedRelations {
		// 查找实体 ID
		fromID, okFrom := entityMap[extracted.FromEntity]
		toID, okTo := entityMap[extracted.ToEntity]

		if !okFrom || !okTo {
			continue
		}

		rel := &Relationship{
			ID:         generateID("rel"),
			Type:       extracted.Type,
			FromNodeID: fromID,
			ToNodeID:   toID,
			Weight:     extracted.Confidence,
			Properties: extracted.Properties,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		relationships = append(relationships, rel)
	}

	return relationships
}

// InferRelationType 推断关系类型
func (r *RelationExtractor) InferRelationType(
	fromType, toType NodeType,
	context string,
) RelationType {
	// 基于节点类型推断
	if fromType == NodeTypeCharacter && toType == NodeTypeCharacter {
		return RelationKnows
	}

	if fromType == NodeTypeCharacter && toType == NodeTypeLocation {
		return RelationLocatedAt
	}

	if fromType == NodeTypeCharacter && toType == NodeTypeItem {
		return RelationOwns
	}

	if fromType == NodeTypeCharacter && toType == NodeTypeConcept {
		return RelationMasters
	}

	if fromType == NodeTypeEvent && toType == NodeTypeLocation {
		return RelationHappensAt
	}

	return RelationKnows // 默认
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
