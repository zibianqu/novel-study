package graph

import (
	"context"
	"regexp"
	"strings"
	"time"
)

// RelationExtractor 关系抽取器
type RelationExtractor struct {
	relationPatterns map[RelationType][]*RelationPattern
}

// RelationPattern 关系模式
type RelationPattern struct {
	Pattern    *regexp.Regexp
	Confidence float64
	Context    string
}

// ExtractedRelation 提取的关系
type ExtractedRelation struct {
	FromEntity string
	ToEntity   string
	Type       RelationType
	Confidence float64
	Context    string
}

// NewRelationExtractor 创建关系抽取器
func NewRelationExtractor() *RelationExtractor {
	return &RelationExtractor{
		relationPatterns: buildRelationPatterns(),
	}
}

// Extract 从文本中提取关系
func (r *RelationExtractor) Extract(
	ctx context.Context,
	text string,
	entities []*Entity,
) ([]*ExtractedRelation, error) {
	relations := make([]*ExtractedRelation, 0)

	// 1. 根据实体类型提取关系
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			entity1 := entities[i]
			entity2 := entities[j]

			// 提取关系
			extracted := r.extractRelationBetween(text, entity1, entity2)
			relations = append(relations, extracted...)
		}
	}

	// 2. 去重
	relations = r.deduplicateRelations(relations)

	return relations, nil
}

// extractRelationBetween 提取两个实体之间的关系
func (r *RelationExtractor) extractRelationBetween(
	text string,
	entity1, entity2 *Entity,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	// 根据实体类型组合决定可能的关系
	switch {
	case entity1.Type == NodeTypeCharacter && entity2.Type == NodeTypeCharacter:
		// 人物-人物关系
		relations = append(relations, r.extractCharacterRelations(text, entity1, entity2)...)

	case entity1.Type == NodeTypeCharacter && entity2.Type == NodeTypeLocation:
		// 人物-地点关系
		relations = append(relations, r.extractCharacterLocationRelations(text, entity1, entity2)...)

	case entity1.Type == NodeTypeCharacter && entity2.Type == NodeTypeEvent:
		// 人物-事件关系
		rel := &ExtractedRelation{
			FromEntity: entity1.Name,
			ToEntity:   entity2.Name,
			Type:       RelationParticipates,
			Confidence: 0.7,
		}
		relations = append(relations, rel)

	case entity1.Type == NodeTypeCharacter && entity2.Type == NodeTypeItem:
		// 人物-物品关系
		relations = append(relations, r.extractCharacterItemRelations(text, entity1, entity2)...)

	case entity1.Type == NodeTypeCharacter && entity2.Type == NodeTypeConcept:
		// 人物-概念关系
		rel := &ExtractedRelation{
			FromEntity: entity1.Name,
			ToEntity:   entity2.Name,
			Type:       RelationMasters,
			Confidence: 0.6,
		}
		relations = append(relations, rel)

	case entity1.Type == NodeTypeEvent && entity2.Type == NodeTypeLocation:
		// 事件-地点关系
		rel := &ExtractedRelation{
			FromEntity: entity1.Name,
			ToEntity:   entity2.Name,
			Type:       RelationHappensAt,
			Confidence: 0.7,
		}
		relations = append(relations, rel)
	}

	return relations
}

// extractCharacterRelations 提取人物关系
func (r *RelationExtractor) extractCharacterRelations(
	text string,
	entity1, entity2 *Entity,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	name1 := entity1.Name
	name2 := entity2.Name

	// 检查师徒关系
	if strings.Contains(text, name1+"的师父"+name2) ||
		strings.Contains(text, name2+"是"+name1+"的师父") {
		relations = append(relations, &ExtractedRelation{
			FromEntity: name2,
			ToEntity:   name1,
			Type:       RelationMasterOf,
			Confidence: 0.9,
		})
	}

	// 检查亲属关系
	if strings.Contains(text, name1+"的父亲"+name2) ||
		strings.Contains(text, name1+"的母亲"+name2) {
		relations = append(relations, &ExtractedRelation{
			FromEntity: name1,
			ToEntity:   name2,
			Type:       RelationFamilyOf,
			Confidence: 0.95,
		})
	}

	// 检查仇敵关系
	if strings.Contains(text, name1+"和"+name2+"是仇敌") ||
		strings.Contains(text, name1+"与"+name2+"成仇") {
		relations = append(relations, &ExtractedRelation{
			FromEntity: name1,
			ToEntity:   name2,
			Type:       RelationEnemyOf,
			Confidence: 0.85,
		})
	}

	// 检查盟友关系
	if strings.Contains(text, name1+"和"+name2+"结盟") ||
		strings.Contains(text, name1+"与"+name2+"结为盟友") {
		relations = append(relations, &ExtractedRelation{
			FromEntity: name1,
			ToEntity:   name2,
			Type:       RelationAllyOf,
			Confidence: 0.8,
		})
	}

	// 默认认识关系
	if len(relations) == 0 {
		// 如果两个人物在同一段落中出现
		if strings.Contains(text, name1) && strings.Contains(text, name2) {
			relations = append(relations, &ExtractedRelation{
				FromEntity: name1,
				ToEntity:   name2,
				Type:       RelationKnows,
				Confidence: 0.5,
			})
		}
	}

	return relations
}

// extractCharacterLocationRelations 提取人物-地点关系
func (r *RelationExtractor) extractCharacterLocationRelations(
	text string,
	entity1, entity2 *Entity,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	charName := entity1.Name
	locName := entity2.Name

	// 检查出生地
	if strings.Contains(text, charName+"出生于"+locName) ||
		strings.Contains(text, charName+"生于"+locName) {
		relations = append(relations, &ExtractedRelation{
			FromEntity: charName,
			ToEntity:   locName,
			Type:       RelationBornAt,
			Confidence: 0.9,
		})
	}

	// 检查居住地
	if strings.Contains(text, charName+"住在"+locName) ||
		strings.Contains(text, charName+"居住于"+locName) {
		relations = append(relations, &ExtractedRelation{
			FromEntity: charName,
			ToEntity:   locName,
			Type:       RelationLivesIn,
			Confidence: 0.85,
		})
	}

	// 默认位置关系
	if len(relations) == 0 {
		if strings.Contains(text, charName) && strings.Contains(text, locName) {
			relations = append(relations, &ExtractedRelation{
				FromEntity: charName,
				ToEntity:   locName,
				Type:       RelationLocatedAt,
				Confidence: 0.6,
			})
		}
	}

	return relations
}

// extractCharacterItemRelations 提取人物-物品关系
func (r *RelationExtractor) extractCharacterItemRelations(
	text string,
	entity1, entity2 *Entity,
) []*ExtractedRelation {
	relations := make([]*ExtractedRelation, 0)

	charName := entity1.Name
	itemName := entity2.Name

	// 检查拥有关系
	if strings.Contains(text, charName+"的"+itemName) ||
		strings.Contains(text, charName+"拥有"+itemName) {
		relations = append(relations, &ExtractedRelation{
			FromEntity: charName,
			ToEntity:   itemName,
			Type:       RelationOwns,
			Confidence: 0.8,
		})
	}

	// 检查使用关系
	if strings.Contains(text, charName+"使用"+itemName) ||
		strings.Contains(text, charName+"施展"+itemName) {
		relations = append(relations, &ExtractedRelation{
			FromEntity: charName,
			ToEntity:   itemName,
			Type:       RelationUses,
			Confidence: 0.75,
		})
	}

	return relations
}

// deduplicateRelations 去重关系
func (r *RelationExtractor) deduplicateRelations(
	relations []*ExtractedRelation,
) []*ExtractedRelation {
	seen := make(map[string]*ExtractedRelation)

	for _, rel := range relations {
		key := rel.FromEntity + "->" + string(rel.Type) + "->" + rel.ToEntity
		if existing, ok := seen[key]; ok {
			// 提高置信度
			if existing.Confidence < 0.95 {
				existing.Confidence += 0.05
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

// ToRelationship 转换为 Relationship
func (er *ExtractedRelation) ToRelationship(
	entityMap map[string]string,
) *Relationship {
	fromID := entityMap[er.FromEntity]
	toID := entityMap[er.ToEntity]

	if fromID == "" || toID == "" {
		return nil
	}

	return &Relationship{
		ID:         generateRelationshipID(fromID, toID, er.Type),
		Type:       er.Type,
		FromNodeID: fromID,
		ToNodeID:   toID,
		Weight:     er.Confidence,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// buildRelationPatterns 构建关系模式
func buildRelationPatterns() map[RelationType][]*RelationPattern {
	patterns := make(map[RelationType][]*RelationPattern)

	// 师徒关系
	patterns[RelationMasterOf] = []*RelationPattern{
		{
			Pattern:    regexp.MustCompile(`(\S+)的师父(\S+)`),
			Confidence: 0.9,
		},
	}

	// 亲属关系
	patterns[RelationFamilyOf] = []*RelationPattern{
		{
			Pattern:    regexp.MustCompile(`(\S+)的(?:父亲|母亲|兄弟|姐妹)(\S+)`),
			Confidence: 0.95,
		},
	}

	return patterns
}

func generateRelationshipID(fromID, toID string, relType RelationType) string {
	return fromID + "_" + string(relType) + "_" + toID
}
