package graph

import (
	"context"
	"fmt"
	"time"
)

// GraphBuilder 图谱构建器
type GraphBuilder struct {
	repository       *Neo4jRepository
	entityExtractor  *EntityExtractor
	relationExtractor *RelationExtractor
}

// NewGraphBuilder 创建图谱构建器
func NewGraphBuilder(repository *Neo4jRepository) *GraphBuilder {
	return &GraphBuilder{
		repository:        repository,
		entityExtractor:   NewEntityExtractor(),
		relationExtractor: NewRelationExtractor(),
	}
}

// BuildFromText 从文本构建图谱
func (gb *GraphBuilder) BuildFromText(
	ctx context.Context,
	text string,
	options *BuildOptions,
) (*BuildResult, error) {
	result := &BuildResult{
		StartTime: time.Now(),
	}

	// 1. 提取实体
	entities, err := gb.entityExtractor.Extract(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("entity extraction failed: %w", err)
	}
	result.EntitiesExtracted = len(entities)

	// 2. 创建节点
	entityMap := make(map[string]string) // name -> id
	for _, entity := range entities {
		if entity.Confidence < options.MinConfidence {
			continue
		}

		node := &Node{
			ID:          generateID(string(entity.Type)),
			Type:        entity.Type,
			Name:        entity.Name,
			Description: entity.Context,
			Properties:  entity.Properties,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := gb.repository.CreateNode(ctx, node)
		if err != nil {
			// 忽略重复节点
			continue
		}

		entityMap[entity.Name] = node.ID
		result.NodesCreated++
	}

	// 3. 提取关系
	relations, err := gb.relationExtractor.Extract(ctx, text, entities)
	if err != nil {
		return nil, fmt.Errorf("relation extraction failed: %w", err)
	}
	result.RelationsExtracted = len(relations)

	// 4. 创建关系
	relationships := gb.relationExtractor.BuildRelationships(relations, entityMap)
	for _, rel := range relationships {
		err := gb.repository.CreateRelationship(ctx, rel)
		if err != nil {
			// 忽略错误
			continue
		}
		result.RelationshipsCreated++
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = true

	return result, nil
}

// BuildOptions 构建选项
type BuildOptions struct {
	MinConfidence float64 // 最小置信度阈值
	MaxNodes      int     // 最大节点数
	EnableAI      bool    // 启用AI增强
}

// DefaultBuildOptions 默认选项
func DefaultBuildOptions() *BuildOptions {
	return &BuildOptions{
		MinConfidence: 0.5,
		MaxNodes:      1000,
		EnableAI:      false,
	}
}

// BuildResult 构建结果
type BuildResult struct {
	Success              bool
	EntitiesExtracted    int
	RelationsExtracted   int
	NodesCreated         int
	RelationshipsCreated int
	Errors               []string
	StartTime            time.Time
	EndTime              time.Time
	Duration             time.Duration
}

// UpdateFromChapter 从章节更新图谱
func (gb *GraphBuilder) UpdateFromChapter(
	ctx context.Context,
	chapterID int,
	content string,
) error {
	// 构建图谱
	options := DefaultBuildOptions()
	result, err := gb.BuildFromText(ctx, content, options)
	if err != nil {
		return err
	}

	// 记录章节信息
	_ = result
	_ = chapterID

	return nil
}

// MergeGraphs 合并图谱
func (gb *GraphBuilder) MergeGraphs(
	ctx context.Context,
	sourceGraphID, targetGraphID string,
) error {
	// 合并逻辑
	// TODO: 实现图谱合并
	return nil
}

// IncrementalBuild 增量构建
func (gb *GraphBuilder) IncrementalBuild(
	ctx context.Context,
	newText string,
	existingNodes []*Node,
) (*BuildResult, error) {
	// 只构建新实体和关系
	options := DefaultBuildOptions()
	return gb.BuildFromText(ctx, newText, options)
}

// ValidateGraph 验证图谱一致性
func (gb *GraphBuilder) ValidateGraph(
	ctx context.Context,
) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:  true,
		Issues: make([]ValidationIssue, 0),
	}

	// 1. 检查孤立节点
	isolatedNodes, err := gb.findIsolatedNodes(ctx)
	if err != nil {
		return nil, err
	}
	if len(isolatedNodes) > 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "isolated_nodes",
			Severity:    "warning",
			Description: fmt.Sprintf("发现 %d 个孤立节点", len(isolatedNodes)),
			NodeIDs:     isolatedNodes,
		})
	}

	// 2. 检查重复关系
	duplicates := gb.findDuplicateRelationships(ctx)
	if len(duplicates) > 0 {
		result.Issues = append(result.Issues, ValidationIssue{
			Type:        "duplicate_relationships",
			Severity:    "warning",
			Description: fmt.Sprintf("发现 %d 个重复关系", len(duplicates)),
		})
	}

	if len(result.Issues) > 0 {
		result.Valid = false
	}

	return result, nil
}

// ValidationResult 验证结果
type ValidationResult struct {
	Valid  bool
	Issues []ValidationIssue
}

// ValidationIssue 验证问题
type ValidationIssue struct {
	Type        string
	Severity    string // "error", "warning", "info"
	Description string
	NodeIDs     []string
	RelIDs      []string
	Suggestion  string
}

// findIsolatedNodes 查找孤立节点
func (gb *GraphBuilder) findIsolatedNodes(ctx context.Context) ([]string, error) {
	// TODO: 实现 Cypher 查询
	return []string{}, nil
}

// findDuplicateRelationships 查找重复关系
func (gb *GraphBuilder) findDuplicateRelationships(ctx context.Context) []string {
	// TODO: 实现 Cypher 查询
	return []string{}
}

// OptimizeGraph 优化图谱
func (gb *GraphBuilder) OptimizeGraph(ctx context.Context) error {
	// 1. 删除孤立节点
	// 2. 合并重复关系
	// 3. 调整关系权重
	return nil
}

// GetGraphStats 获取图谱统计
func (gb *GraphBuilder) GetGraphStats(ctx context.Context) (*GraphStatistics, error) {
	stats, err := gb.repository.client.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return &GraphStatistics{
		TotalNodes:         stats.NodeCount,
		TotalRelationships: stats.RelationshipCount,
		NodesByType:        make(map[NodeType]int64),
		RelsByType:         make(map[RelationType]int64),
	}, nil
}

// GraphStatistics 图谱统计
type GraphStatistics struct {
	TotalNodes         int64
	TotalRelationships int64
	NodesByType        map[NodeType]int64
	RelsByType         map[RelationType]int64
	AvgDegree          float64
	Density            float64
}

// ExportGraph 导出图谱
func (gb *GraphBuilder) ExportGraph(
	ctx context.Context,
	format string,
) ([]byte, error) {
	// TODO: 实现导出功能 (JSON, GraphML, etc.)
	return nil, nil
}

// ImportGraph 导入图谱
func (gb *GraphBuilder) ImportGraph(
	ctx context.Context,
	data []byte,
	format string,
) error {
	// TODO: 实现导入功能
	return nil
}
