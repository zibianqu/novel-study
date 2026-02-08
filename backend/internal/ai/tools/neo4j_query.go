package tools

import (
	"context"
	"fmt"

	"github.com/zibianqu/novel-study/internal/repository"
)

// Neo4jQueryTool Neo4j查询工具
type Neo4jQueryTool struct {
	neo4jRepo *repository.Neo4jRepository
}

// NewNeo4jQueryTool 创建Neo4j查询工具
func NewNeo4jQueryTool(neo4jRepo *repository.Neo4jRepository) *Neo4jQueryTool {
	return &Neo4jQueryTool{
		neo4jRepo: neo4jRepo,
	}
}

func (t *Neo4jQueryTool) GetName() string {
	return "query_neo4j"
}

func (t *Neo4jQueryTool) GetDescription() string {
	return "查询知识图谱中的关系数据。参数: query_type(查询类型: character_relations|角色关系, world_events|世界事件, plot_arcs|剧情弧), entity_id(实体ID), project_id(项目ID)"
}

func (t *Neo4jQueryTool) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	queryType, ok := params["query_type"].(string)
	if !ok || queryType == "" {
		return nil, fmt.Errorf("missing required parameter: query_type")
	}

	projectID, ok := params["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing required parameter: project_id")
	}

	switch queryType {
	case "character_relations":
		return t.queryCharacterRelations(ctx, int(projectID), params)
	case "world_events":
		return t.queryWorldEvents(ctx, int(projectID), params)
	case "plot_arcs":
		return t.queryPlotArcs(ctx, int(projectID), params)
	case "character_state":
		return t.queryCharacterState(ctx, int(projectID), params)
	default:
		return nil, fmt.Errorf("unsupported query_type: %s", queryType)
	}
}

// queryCharacterRelations 查询角色关系
func (t *Neo4jQueryTool) queryCharacterRelations(ctx context.Context, projectID int, params map[string]interface{}) (interface{}, error) {
	characterID, ok := params["character_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing character_id for character_relations query")
	}

	relations, err := t.neo4jRepo.GetCharacterRelations(ctx, int(characterID), projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query character relations: %w", err)
	}

	return map[string]interface{}{
		"type":      "character_relations",
		"relations": relations,
	}, nil
}

// queryWorldEvents 查询世界事件
func (t *Neo4jQueryTool) queryWorldEvents(ctx context.Context, projectID int, params map[string]interface{}) (interface{}, error) {
	limit := 10
	if l, ok := params["limit"].(float64); ok {
		limit = int(l)
	}

	events, err := t.neo4jRepo.GetRecentWorldEvents(ctx, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query world events: %w", err)
	}

	return map[string]interface{}{
		"type":   "world_events",
		"events": events,
	}, nil
}

// queryPlotArcs 查询剧情弧
func (t *Neo4jQueryTool) queryPlotArcs(ctx context.Context, projectID int, params map[string]interface{}) (interface{}, error) {
	status := "active"
	if s, ok := params["status"].(string); ok {
		status = s
	}

	arcs, err := t.neo4jRepo.GetPlotArcs(ctx, projectID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query plot arcs: %w", err)
	}

	return map[string]interface{}{
		"type": "plot_arcs",
		"arcs": arcs,
	}, nil
}

// queryCharacterState 查询角色当前状态
func (t *Neo4jQueryTool) queryCharacterState(ctx context.Context, projectID int, params map[string]interface{}) (interface{}, error) {
	characterID, ok := params["character_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing character_id for character_state query")
	}

	state, err := t.neo4jRepo.GetCharacterCurrentState(ctx, int(characterID), projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query character state: %w", err)
	}

	return map[string]interface{}{
		"type":  "character_state",
		"state": state,
	}, nil
}
