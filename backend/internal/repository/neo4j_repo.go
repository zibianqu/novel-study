package repository

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jRepository Neo4j 数据访问层
type Neo4jRepository struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jRepository(driver neo4j.DriverWithContext) *Neo4jRepository {
	return &Neo4jRepository{driver: driver}
}

// GraphNode 图节点
type GraphNode struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// GraphRelation 图关系
type GraphRelation struct {
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// CreateNode 创建节点
func (r *Neo4jRepository) CreateNode(ctx context.Context, projectID int, node *GraphNode) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		CREATE (n:%s {
			id: $id,
			project_id: $project_id,
			name: $name
		})
		RETURN n
	`, node.Type)

	_, err := session.Run(ctx, query, map[string]interface{}{
		"id":         node.ID,
		"project_id": projectID,
		"name":       node.Label,
	})

	return err
}

// CreateRelation 创建关系
func (r *Neo4jRepository) CreateRelation(ctx context.Context, projectID int, rel *GraphRelation) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MATCH (a {id: $source, project_id: $project_id})
		MATCH (b {id: $target, project_id: $project_id})
		CREATE (a)-[r:%s]->(b)
		RETURN r
	`, rel.Type)

	_, err := session.Run(ctx, query, map[string]interface{}{
		"source":     rel.Source,
		"target":     rel.Target,
		"project_id": projectID,
	})

	return err
}

// GetProjectGraph 获取项目图谱
func (r *Neo4jRepository) GetProjectGraph(ctx context.Context, projectID int) ([]*GraphNode, []*GraphRelation, error) {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// 获取节点
	nodesResult, err := session.Run(ctx,
		"MATCH (n) WHERE n.project_id = $project_id RETURN n",
		map[string]interface{}{"project_id": projectID},
	)
	if err != nil {
		return nil, nil, err
	}

	var nodes []*GraphNode
	for nodesResult.Next(ctx) {
		record := nodesResult.Record()
		nodeValue, ok := record.Get("n")
		if !ok {
			continue
		}

		if node, ok := nodeValue.(neo4j.Node); ok {
			props := node.Props
			nodes = append(nodes, &GraphNode{
				ID:         fmt.Sprintf("%v", props["id"]),
				Label:      fmt.Sprintf("%v", props["name"]),
				Type:       node.Labels[0],
				Properties: props,
			})
		}
	}

	// 获取关系
	relsResult, err := session.Run(ctx,
		"MATCH (a)-[r]->(b) WHERE a.project_id = $project_id RETURN a.id as source, b.id as target, type(r) as type",
		map[string]interface{}{"project_id": projectID},
	)
	if err != nil {
		return nodes, nil, err
	}

	var relations []*GraphRelation
	for relsResult.Next(ctx) {
		record := relsResult.Record()
		source, _ := record.Get("source")
		target, _ := record.Get("target")
		relType, _ := record.Get("type")

		relations = append(relations, &GraphRelation{
			Source: fmt.Sprintf("%v", source),
			Target: fmt.Sprintf("%v", target),
			Type:   fmt.Sprintf("%v", relType),
		})
	}

	return nodes, relations, nil
}

// DeleteNode 删除节点
func (r *Neo4jRepository) DeleteNode(ctx context.Context, projectID int, nodeID string) error {
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		"MATCH (n {id: $id, project_id: $project_id}) DETACH DELETE n",
		map[string]interface{}{
			"id":         nodeID,
			"project_id": projectID,
		},
	)

	return err
}
