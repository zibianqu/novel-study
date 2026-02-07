package repository

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jRepository Neo4j图数据库操作封装
type Neo4jRepository struct {
	driver   neo4j.DriverWithContext
	database string
}

// NewNeo4jRepository 创建Neo4j仓库
func NewNeo4jRepository(driver neo4j.DriverWithContext, database string) *Neo4jRepository {
	return &Neo4jRepository{driver: driver, database: database}
}

func (r *Neo4jRepository) session(ctx context.Context) neo4j.SessionWithContext {
	return r.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: r.database,
		AccessMode:   neo4j.AccessModeWrite,
	})
}

// ==================== 初始化约束和索引 ====================

// InitSchema 初始化Neo4j Schema
func (r *Neo4jRepository) InitSchema(ctx context.Context) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	constraints := []string{
		"CREATE CONSTRAINT IF NOT EXISTS FOR (c:Character) REQUIRE c.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (o:Organization) REQUIRE o.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (l:Location) REQUIRE l.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (e:Event) REQUIRE e.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (we:WorldEvent) REQUIRE we.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (f:Foreshadow) REQUIRE f.id IS UNIQUE",
		"CREATE INDEX IF NOT EXISTS FOR (c:Character) ON (c.project_id)",
		"CREATE INDEX IF NOT EXISTS FOR (e:Event) ON (e.project_id)",
		"CREATE INDEX IF NOT EXISTS FOR (we:WorldEvent) ON (we.project_id)",
	}

	for _, c := range constraints {
		if _, err := session.Run(ctx, c, nil); err != nil {
			return fmt.Errorf("创建Neo4j约束失败: %w", err)
		}
	}
	return nil
}

// ==================== 角色图谱 ====================

// CharacterNode 角色节点
type CharacterNode struct {
	ID              int    `json:"id"`
	ProjectID       int    `json:"project_id"`
	Name            string `json:"name"`
	RoleType        string `json:"role_type"`
	PowerLevel      int    `json:"power_level"`
	MentalState     string `json:"mental_state"`
	CurrentLocation string `json:"current_location"`
	Status          string `json:"status"`
}

// UpsertCharacter 创建或更新角色节点
func (r *Neo4jRepository) UpsertCharacter(ctx context.Context, c CharacterNode) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (c:Character {id: $id})
		 SET c.project_id = $project_id, c.name = $name, c.role_type = $role_type,
			 c.power_level = $power_level, c.mental_state = $mental_state,
			 c.current_location = $current_location, c.status = $status`,
		map[string]interface{}{
			"id": c.ID, "project_id": c.ProjectID, "name": c.Name,
			"role_type": c.RoleType, "power_level": c.PowerLevel,
			"mental_state": c.MentalState, "current_location": c.CurrentLocation,
			"status": c.Status,
		})
	return err
}

// CharacterRelation 角色关系
type CharacterRelation struct {
	FromID   int    `json:"from_id"`
	ToID     int    `json:"to_id"`
	RelType  string `json:"rel_type"` // ALLY/ENEMY/FAMILY/MASTER_STUDENT/LOVER
	Since    string `json:"since"`
	Desc     string `json:"desc"`
}

// AddCharacterRelation 添加角色关系
func (r *Neo4jRepository) AddCharacterRelation(ctx context.Context, rel CharacterRelation) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	query := fmt.Sprintf(
		`MATCH (a:Character {id: $from_id}), (b:Character {id: $to_id})
		 MERGE (a)-[r:%s]->(b)
		 SET r.since = $since, r.description = $desc`, rel.RelType)

	_, err := session.Run(ctx, query, map[string]interface{}{
		"from_id": rel.FromID, "to_id": rel.ToID,
		"since": rel.Since, "desc": rel.Desc,
	})
	return err
}

// GetCharacterRelations 获取角色的所有关系
func (r *Neo4jRepository) GetCharacterRelations(ctx context.Context, projectID int) ([]map[string]interface{}, error) {
	session := r.session(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		`MATCH (a:Character {project_id: $pid})-[r]->(b:Character {project_id: $pid})
		 RETURN a.id as from_id, a.name as from_name, type(r) as rel_type, 
				r.description as description, b.id as to_id, b.name as to_name`,
		map[string]interface{}{"pid": projectID})
	if err != nil {
		return nil, err
	}

	var relations []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		rel := map[string]interface{}{
			"from_id":     record.Values[0],
			"from_name":   record.Values[1],
			"rel_type":    record.Values[2],
			"description": record.Values[3],
			"to_id":       record.Values[4],
			"to_name":     record.Values[5],
		}
		relations = append(relations, rel)
	}
	return relations, nil
}

// DeleteCharacter 删除角色节点及其关系
func (r *Neo4jRepository) DeleteCharacter(ctx context.Context, charID int) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MATCH (c:Character {id: $id}) DETACH DELETE c`,
		map[string]interface{}{"id": charID})
	return err
}

// ==================== 世界事件图谱（天线） ====================

// WorldEventNode 世界事件节点
type WorldEventNode struct {
	ID          int    `json:"id"`
	ProjectID   int    `json:"project_id"`
	Name        string `json:"name"`
	ImpactLevel int    `json:"impact_level"`
	TimePoint   string `json:"time_point"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// UpsertWorldEvent 创建或更新世界事件
func (r *Neo4jRepository) UpsertWorldEvent(ctx context.Context, e WorldEventNode) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (e:WorldEvent {id: $id})
		 SET e.project_id = $project_id, e.name = $name, 
			 e.impact_level = $impact_level, e.time_point = $time_point,
			 e.description = $description, e.status = $status`,
		map[string]interface{}{
			"id": e.ID, "project_id": e.ProjectID, "name": e.Name,
			"impact_level": e.ImpactLevel, "time_point": e.TimePoint,
			"description": e.Description, "status": e.Status,
		})
	return err
}

// LinkWorldEvents 链接世界事件因果关系
func (r *Neo4jRepository) LinkWorldEvents(ctx context.Context, causeID, effectID int) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MATCH (a:WorldEvent {id: $cause_id}), (b:WorldEvent {id: $effect_id})
		 MERGE (a)-[:CAUSES]->(b)`,
		map[string]interface{}{"cause_id": causeID, "effect_id": effectID})
	return err
}

// GetWorldEvents 获取项目的世界事件链
func (r *Neo4jRepository) GetWorldEvents(ctx context.Context, projectID int) ([]map[string]interface{}, error) {
	session := r.session(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		`MATCH (e:WorldEvent {project_id: $pid})
		 OPTIONAL MATCH (e)-[:CAUSES]->(next:WorldEvent)
		 RETURN e.id, e.name, e.impact_level, e.time_point, e.status, 
				collect(next.id) as causes`,
		map[string]interface{}{"pid": projectID})
	if err != nil {
		return nil, err
	}

	var events []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		events = append(events, map[string]interface{}{
			"id":           record.Values[0],
			"name":         record.Values[1],
			"impact_level": record.Values[2],
			"time_point":   record.Values[3],
			"status":       record.Values[4],
			"causes":       record.Values[5],
		})
	}
	return events, nil
}

// ==================== 伏笔图谱 ====================

// ForeshadowNode 伏笔节点
type ForeshadowNode struct {
	ID                   int    `json:"id"`
	ProjectID            int    `json:"project_id"`
	Content              string `json:"content"`
	PlantedChapter       int    `json:"planted_chapter"`
	PlannedResolveChapter int   `json:"planned_resolve_chapter"`
	Status               string `json:"status"` // planted/resolved/abandoned
}

// UpsertForeshadow 创建或更新伏笔
func (r *Neo4jRepository) UpsertForeshadow(ctx context.Context, f ForeshadowNode) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx,
		`MERGE (f:Foreshadow {id: $id})
		 SET f.project_id = $project_id, f.content = $content,
			 f.planted_chapter = $planted_chapter, 
			 f.planned_resolve_chapter = $planned_resolve_chapter,
			 f.status = $status`,
		map[string]interface{}{
			"id": f.ID, "project_id": f.ProjectID, "content": f.Content,
			"planted_chapter": f.PlantedChapter,
			"planned_resolve_chapter": f.PlannedResolveChapter,
			"status": f.Status,
		})
	return err
}

// GetUnresolvedForeshadows 获取未回收的伏笔
func (r *Neo4jRepository) GetUnresolvedForeshadows(ctx context.Context, projectID int) ([]map[string]interface{}, error) {
	session := r.session(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx,
		`MATCH (f:Foreshadow {project_id: $pid, status: 'planted'})
		 RETURN f.id, f.content, f.planted_chapter, f.planned_resolve_chapter
		 ORDER BY f.planted_chapter`,
		map[string]interface{}{"pid": projectID})
	if err != nil {
		return nil, err
	}

	var foreshadows []map[string]interface{}
	for result.Next(ctx) {
		record := result.Record()
		foreshadows = append(foreshadows, map[string]interface{}{
			"id":                      record.Values[0],
			"content":                 record.Values[1],
			"planted_chapter":         record.Values[2],
			"planned_resolve_chapter": record.Values[3],
		})
	}
	return foreshadows, nil
}

// ==================== 清理 ====================

// ClearProjectGraph 清除项目的所有图谱数据
func (r *Neo4jRepository) ClearProjectGraph(ctx context.Context, projectID int) error {
	session := r.session(ctx)
	defer session.Close(ctx)

	queries := []string{
		"MATCH (n:Character {project_id: $pid}) DETACH DELETE n",
		"MATCH (n:WorldEvent {project_id: $pid}) DETACH DELETE n",
		"MATCH (n:Event {project_id: $pid}) DETACH DELETE n",
		"MATCH (n:Location {project_id: $pid}) DETACH DELETE n",
		"MATCH (n:Organization {project_id: $pid}) DETACH DELETE n",
		"MATCH (n:Foreshadow {project_id: $pid}) DETACH DELETE n",
	}

	for _, q := range queries {
		if _, err := session.Run(ctx, q, map[string]interface{}{"pid": projectID}); err != nil {
			return fmt.Errorf("清理图谱失败: %w", err)
		}
	}
	return nil
}
