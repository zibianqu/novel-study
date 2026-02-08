package graph

import (
	"context"
	"fmt"
)

// GraphService 图谱服务
type GraphService struct {
	repository *Neo4jRepository
	builder    *GraphBuilder
}

// NewGraphService 创建图谱服务
func NewGraphService(client *Neo4jClient) *GraphService {
	repository := NewNeo4jRepository(client)
	builder := NewGraphBuilder(repository)

	return &GraphService{
		repository: repository,
		builder:    builder,
	}
}

// CreateKnowledgeGraph 创建知识图谱
func (s *GraphService) CreateKnowledgeGraph(
	ctx context.Context,
	req *CreateGraphRequest,
) (*CreateGraphResponse, error) {
	// 构建图谱
	options := &BuildOptions{
		MinConfidence: req.MinConfidence,
		MaxNodes:      req.MaxNodes,
	}

	result, err := s.builder.BuildFromText(ctx, req.Text, options)
	if err != nil {
		return nil, err
	}

	return &CreateGraphResponse{
		Success:              result.Success,
		NodesCreated:         result.NodesCreated,
		RelationshipsCreated: result.RelationshipsCreated,
		Duration:             result.Duration.String(),
	}, nil
}

// CreateGraphRequest 创建图谱请求
type CreateGraphRequest struct {
	ProjectID     int
	Text          string
	MinConfidence float64
	MaxNodes      int
}

// CreateGraphResponse 创建图谱响应
type CreateGraphResponse struct {
	Success              bool
	NodesCreated         int
	RelationshipsCreated int
	Duration             string
}

// QueryGraph 查询图谱
func (s *GraphService) QueryGraph(
	ctx context.Context,
	query *GraphQuery,
) (*GraphResult, error) {
	// 根据查询类型调用不同方法
	if query.NodeID != "" {
		// 查询子图
		if query.Depth > 0 {
			return s.repository.GetSubgraph(ctx, query.NodeID, query.Depth)
		}
		// 查询单个节点
		node, err := s.repository.GetNode(ctx, query.NodeID)
		if err != nil {
			return nil, err
		}
		return &GraphResult{
			Nodes: []*Node{node},
		}, nil
	}

	// 查询所有节点
	nodes, err := s.repository.FindNodes(ctx, query)
	if err != nil {
		return nil, err
	}

	return &GraphResult{
		Nodes: nodes,
	}, nil
}

// FindPath 查找路径
func (s *GraphService) FindPath(
	ctx context.Context,
	req *PathRequest,
) (*PathResponse, error) {
	var paths []Path
	var err error

	if req.Shortest {
		// 最短路径
		path, err := s.repository.FindShortestPath(ctx, req.FromNodeID, req.ToNodeID)
		if err != nil {
			return nil, err
		}
		paths = []Path{*path}
	} else {
		// 所有路径
		paths, err = s.repository.FindPath(ctx, req.FromNodeID, req.ToNodeID, req.MaxDepth)
		if err != nil {
			return nil, err
		}
	}

	return &PathResponse{
		Paths: paths,
		Count: len(paths),
	}, nil
}

// PathRequest 路径查询请求
type PathRequest struct {
	FromNodeID string
	ToNodeID   string
	MaxDepth   int
	Shortest   bool
}

// PathResponse 路径查询响应
type PathResponse struct {
	Paths []Path
	Count int
}

// GetStatistics 获取统计信息
func (s *GraphService) GetStatistics(ctx context.Context) (*GraphStatistics, error) {
	return s.builder.GetGraphStats(ctx)
}

// ValidateConsistency 验证一致性
func (s *GraphService) ValidateConsistency(ctx context.Context) (*ValidationResult, error) {
	return s.builder.ValidateGraph(ctx)
}

// AnalyzeCharacterRelations 分析人物关系
func (s *GraphService) AnalyzeCharacterRelations(
	ctx context.Context,
	characterID string,
) (*CharacterAnalysis, error) {
	// 获取邻居节点
	neighbors, err := s.repository.GetNeighbors(ctx, characterID)
	if err != nil {
		return nil, err
	}

	analysis := &CharacterAnalysis{
		CharacterID:   characterID,
		RelationCount: len(neighbors),
		Relations:     make(map[string]int),
	}

	// 统计关系类型
	for _, neighbor := range neighbors {
		analysis.Relations[string(neighbor.Type)]++
	}

	return analysis, nil
}

// CharacterAnalysis 人物分析
type CharacterAnalysis struct {
	CharacterID   string
	RelationCount int
	Relations     map[string]int
	Importance    float64
}

// DetectPlotHoles 检测剧情漏洞
func (s *GraphService) DetectPlotHoles(
	ctx context.Context,
) (*PlotHoleReport, error) {
	report := &PlotHoleReport{
		Holes: make([]*PlotHole, 0),
	}

	// 检测孤立节点
	validation, err := s.builder.ValidateGraph(ctx)
	if err != nil {
		return nil, err
	}

	for _, issue := range validation.Issues {
		if issue.Type == "isolated_nodes" {
			hole := &PlotHole{
				Type:        "isolated_character",
				Severity:    "medium",
				Description: issue.Description,
				Suggestion:  "建议添加与其他人物的关系",
			}
			report.Holes = append(report.Holes, hole)
		}
	}

	return report, nil
}

// PlotHoleReport 剧情漏洞报告
type PlotHoleReport struct {
	Holes      []*PlotHole
	TotalCount int
}

// PlotHole 剧情漏洞
type PlotHole struct {
	Type        string
	Severity    string
	Description string
	Suggestion  string
	Location    string
}

// GenerateWritingSuggestions 生成写作建议
func (s *GraphService) GenerateWritingSuggestions(
	ctx context.Context,
	projectID int,
) (*WritingSuggestions, error) {
	suggestions := &WritingSuggestions{
		Suggestions: make([]string, 0),
	}

	// 基于图谱分析生成建议
	stats, err := s.GetStatistics(ctx)
	if err != nil {
		return nil, err
	}

	// 如果人物关系过少
	if stats.TotalRelationships < stats.TotalNodes {
		suggestions.Suggestions = append(suggestions.Suggestions,
			"建议增加人物之间的互动和关系")
	}

	// 如果没有地点描写
	if stats.NodesByType[NodeTypeLocation] == 0 {
		suggestions.Suggestions = append(suggestions.Suggestions,
			"建议添加场景和地点描写")
	}

	return suggestions, nil
}

// WritingSuggestions 写作建议
type WritingSuggestions struct {
	Suggestions []string
	Priority    string
}

// GetCharacterTimeline 获取人物时间线
func (s *GraphService) GetCharacterTimeline(
	ctx context.Context,
	characterID string,
) (*Timeline, error) {
	// 查询人物相关事件
	// TODO: 实现 Cypher 查询

	return &Timeline{
		CharacterID: characterID,
		Events:      make([]*TimelineEvent, 0),
	}, nil
}

// Timeline 时间线
type Timeline struct {
	CharacterID string
	Events      []*TimelineEvent
}

// TimelineEvent 时间线事件
type TimelineEvent struct {
	EventID     string
	EventName   string
	Timestamp   string
	Chapter     int
	Description string
}

// SearchGraph 搜索图谱
func (s *GraphService) SearchGraph(
	ctx context.Context,
	keyword string,
) (*SearchResult, error) {
	// 搜索包含关键词的节点
	query := &GraphQuery{
		Limit: 50,
	}

	nodes, err := s.repository.FindNodes(ctx, query)
	if err != nil {
		return nil, err
	}

	// 过滤匹配的节点
	matched := make([]*Node, 0)
	for _, node := range nodes {
		if contains(node.Name, keyword) || contains(node.Description, keyword) {
			matched = append(matched, node)
		}
	}

	return &SearchResult{
		Nodes: matched,
		Total: len(matched),
	}, nil
}

// SearchResult 搜索结果
type SearchResult struct {
	Nodes []*Node
	Total int
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || (len(s) > len(substr) && 
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}

// HealthCheck 健康检查
func (s *GraphService) HealthCheck(ctx context.Context) error {
	return s.repository.client.HealthCheck(ctx)
}
