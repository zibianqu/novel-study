package graph

import (
	"context"
	"fmt"
	"time"
)

// GraphBuilder 图谱构建器
type GraphBuilder struct {
	repository       GraphRepository
	entityExtractor  *EntityExtractor
	relationExtractor *RelationExtractor
}

// NewGraphBuilder 创建图谱构建器
func NewGraphBuilder(
	repository GraphRepository,
	entityExtractor *EntityExtractor,
	relationExtractor *RelationExtractor,
) *GraphBuilder {
	return &GraphBuilder{
		repository:        repository,
		entityExtractor:   entityExtractor,
		relationExtractor: relationExtractor,
	}
}

// BuildFromText 从文本构建图谱
func (gb *GraphBuilder) BuildFromText(
	ctx context.Context,
	text string,
	chapter int,
) (*BuildResult, error) {
	// 1. 提取实体
	entities, err := gb.entityExtractor.Extract(ctx, text, chapter)
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}

	// 2. 创建节点
	entityMap := make(map[string]string) // name -> id
	for _, entity := range entities {
		node := gb.entityToNode(entity)
		err := gb.repository.CreateNode(ctx, node)
		if err != nil {
			// 忽略重复节点错误
			continue
		}
		entityMap[entity.Name] = entity.ID
	}

	// 3. 提取关系
	relations, err := gb.relationExtractor.Extract(ctx, text, entities)
	if err != nil {
		return nil, fmt.Errorf("relation extraction failed: %w", err)
	}

	// 4. 创建关系
	for _, rel := range relations {
		relationship := rel.ToRelationship(entityMap)
		if relationship != nil {
			err := gb.repository.CreateRelationship(ctx, relationship)
			if err != nil {
				// 忽略重复关系错误
				continue
			}
		}
	}

	return &BuildResult{
		NodesCreated:         len(entities),
		RelationshipsCreated: len(relations),
		Entities:            entities,
		Relations:           relations,
	}, nil
}

// BuildFromChapter 从章节构建图谱
func (gb *GraphBuilder) BuildFromChapter(
	ctx context.Context,
	chapterID int,
	chapterContent string,
) (*BuildResult, error) {
	return gb.BuildFromText(ctx, chapterContent, chapterID)
}

// UpdateGraph 更新图谱
func (gb *GraphBuilder) UpdateGraph(
	ctx context.Context,
	newText string,
	chapter int,
) (*BuildResult, error) {
	// 1. 提取新实体
	newEntities, err := gb.entityExtractor.Extract(ctx, newText, chapter)
	if err != nil {
		return nil, err
	}

	// 2. 检查实体是否已存在
	entityMap := make(map[string]string)
	for _, entity := range newEntities {
		existingNode, err := gb.repository.GetNode(ctx, entity.ID)
		if err != nil || existingNode == nil {
			// 创建新节点
			node := gb.entityToNode(entity)
			gb.repository.CreateNode(ctx, node)
			entityMap[entity.Name] = entity.ID
		} else {
			// 更新现有节点
			existingNode.Name = entity.Name
			gb.repository.UpdateNode(ctx, existingNode)
			entityMap[entity.Name] = entity.ID
		}
	}

	// 3. 提取并创建关系
	relations, _ := gb.relationExtractor.Extract(ctx, newText, newEntities)
	for _, rel := range relations {
		relationship := rel.ToRelationship(entityMap)
		if relationship != nil {
			gb.repository.CreateRelationship(ctx, relationship)
		}
	}

	return &BuildResult{
		NodesCreated:         len(newEntities),
		RelationshipsCreated: len(relations),
	}, nil
}

// MergeEntities 合并实体
func (gb *GraphBuilder) MergeEntities(
	ctx context.Context,
	entity1ID, entity2ID string,
) error {
	// 1. 获取两个节点
	node1, err := gb.repository.GetNode(ctx, entity1ID)
	if err != nil {
		return err
	}

	node2, err := gb.repository.GetNode(ctx, entity2ID)
	if err != nil {
		return err
	}

	// 2. TODO: 将 node2 的所有关系转移到 node1
	// 这里需要 Cypher 查询来实现

	// 3. 删除 node2
	err = gb.repository.DeleteNode(ctx, entity2ID)
	if err != nil {
		return err
	}

	_ = node1
	_ = node2

	return nil
}

// entityToNode 将实体转换为节点
func (gb *GraphBuilder) entityToNode(entity *Entity) *Node {
	node := &Node{
		ID:          entity.ID,
		Type:        entity.Type,
		Name:        entity.Name,
		Description: entity.Context,
		Properties:  make(map[string]interface{}),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 添加属性
	node.Properties["mentions"] = entity.Mentions
	node.Properties["confidence"] = entity.Confidence
	node.Properties["chapter"] = entity.Chapter

	for k, v := range entity.Attributes {
		node.Properties[k] = v
	}

	return node
}

// BuildResult 构建结果
type BuildResult struct {
	NodesCreated         int
	RelationshipsCreated int
	Entities            []*Entity
	Relations           []*ExtractedRelation
}

// ConsistencyChecker 一致性检查器
type ConsistencyChecker struct {
	repository GraphRepository
}

// NewConsistencyChecker 创建一致性检查器
func NewConsistencyChecker(repository GraphRepository) *ConsistencyChecker {
	return &ConsistencyChecker{
		repository: repository,
	}
}

// CheckConsistency 检查一致性
func (cc *ConsistencyChecker) CheckConsistency(
	ctx context.Context,
) ([]*ConsistencyIssue, error) {
	issues := make([]*ConsistencyIssue, 0)

	// 1. 检查孤立节点
	isolatedIssues := cc.checkIsolatedNodes(ctx)
	issues = append(issues, isolatedIssues...)

	// 2. 检查重复实体
	duplicateIssues := cc.checkDuplicateEntities(ctx)
	issues = append(issues, duplicateIssues...)

	// 3. 检查矛盾关系
	conflictIssues := cc.checkConflictingRelations(ctx)
	issues = append(issues, conflictIssues...)

	return issues, nil
}

// checkIsolatedNodes 检查孤立节点
func (cc *ConsistencyChecker) checkIsolatedNodes(
	ctx context.Context,
) []*ConsistencyIssue {
	issues := make([]*ConsistencyIssue, 0)

	// TODO: 查询没有任何关系的节点
	// MATCH (n) WHERE NOT (n)--() RETURN n

	return issues
}

// checkDuplicateEntities 检查重复实体
func (cc *ConsistencyChecker) checkDuplicateEntities(
	ctx context.Context,
) []*ConsistencyIssue {
	issues := make([]*ConsistencyIssue, 0)

	// TODO: 查询同名的节点
	// MATCH (n:Character) WITH n.name as name, collect(n) as nodes
	// WHERE size(nodes) > 1 RETURN name, nodes

	return issues
}

// checkConflictingRelations 检查矛盾关系
func (cc *ConsistencyChecker) checkConflictingRelations(
	ctx context.Context,
) []*ConsistencyIssue {
	issues := make([]*ConsistencyIssue, 0)

	// TODO: 检查矛盾的关系
	// 例如：A 是 B 的师父，但 B 也是 A 的师父

	return issues
}

// ConsistencyIssue 一致性问题
type ConsistencyIssue struct {
	Type        string   // "isolated", "duplicate", "conflict"
	Description string
	NodeIDs     []string
	Severity    string // "low", "medium", "high"
	Suggestion  string
}

// GraphService 图谱服务
type GraphService struct {
	builder  *GraphBuilder
	checker  *ConsistencyChecker
	repository GraphRepository
}

// NewGraphService 创建图谱服务
func NewGraphService(
	repository GraphRepository,
) *GraphService {
	return &GraphService{
		builder: NewGraphBuilder(
			repository,
			NewEntityExtractor(),
			NewRelationExtractor(),
		),
		checker:    NewConsistencyChecker(repository),
		repository: repository,
	}
}

// BuildGraph 构建图谱
func (gs *GraphService) BuildGraph(
	ctx context.Context,
	text string,
	chapter int,
) (*BuildResult, error) {
	return gs.builder.BuildFromText(ctx, text, chapter)
}

// CheckConsistency 检查一致性
func (gs *GraphService) CheckConsistency(
	ctx context.Context,
) ([]*ConsistencyIssue, error) {
	return gs.checker.CheckConsistency(ctx)
}

// QueryGraph 查询图谱
func (gs *GraphService) QueryGraph(
	ctx context.Context,
	query *GraphQuery,
) (*GraphResult, error) {
	nodes, err := gs.repository.FindNodes(ctx, query)
	if err != nil {
		return nil, err
	}

	return &GraphResult{
		Nodes:         nodes,
		Relationships: make([]*Relationship, 0),
	}, nil
}
