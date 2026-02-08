// NovelForge AI - Neo4j 图谱索引初始化
// 创建日期: 2026-02-08
// 使用方法: cypher-shell -u neo4j -p <password> -f init_neo4j_indexes.cypher

// =====================================================
// 1. 删除旧索引（如果存在）
// =====================================================

DROP INDEX character_project_idx IF EXISTS;
DROP INDEX character_name_idx IF EXISTS;
DROP INDEX organization_project_idx IF EXISTS;
DROP INDEX location_project_idx IF EXISTS;
DROP INDEX event_chapter_idx IF EXISTS;
DROP INDEX world_event_project_idx IF EXISTS;
DROP INDEX plot_arc_project_idx IF EXISTS;
DROP INDEX foreshadow_project_idx IF EXISTS;

// =====================================================
// 2. 创建节点索引
// =====================================================

// Character 节点索引
CREATE INDEX character_project_idx IF NOT EXISTS
FOR (c:Character)
ON (c.project_id);

CREATE INDEX character_name_idx IF NOT EXISTS
FOR (c:Character)
ON (c.name);

CREATE INDEX character_id_idx IF NOT EXISTS
FOR (c:Character)
ON (c.id);

// Organization 节点索引
CREATE INDEX organization_project_idx IF NOT EXISTS
FOR (o:Organization)
ON (o.project_id);

CREATE INDEX organization_name_idx IF NOT EXISTS
FOR (o:Organization)
ON (o.name);

// Location 节点索引
CREATE INDEX location_project_idx IF NOT EXISTS
FOR (l:Location)
ON (l.project_id);

CREATE INDEX location_name_idx IF NOT EXISTS
FOR (l:Location)
ON (l.name);

// Event 节点索引
CREATE INDEX event_chapter_idx IF NOT EXISTS
FOR (e:Event)
ON (e.chapter_id);

CREATE INDEX event_project_idx IF NOT EXISTS
FOR (e:Event)
ON (e.project_id);

CREATE INDEX event_time_idx IF NOT EXISTS
FOR (e:Event)
ON (e.time_point);

// WorldEvent 节点索引（天线）
CREATE INDEX world_event_project_idx IF NOT EXISTS
FOR (we:WorldEvent)
ON (we.project_id);

CREATE INDEX world_event_time_idx IF NOT EXISTS
FOR (we:WorldEvent)
ON (we.time_point);

CREATE INDEX world_event_impact_idx IF NOT EXISTS
FOR (we:WorldEvent)
ON (we.impact_level);

// PlotArc 节点索引（剧情线）
CREATE INDEX plot_arc_project_idx IF NOT EXISTS
FOR (pa:PlotArc)
ON (pa.project_id);

CREATE INDEX plot_arc_status_idx IF NOT EXISTS
FOR (pa:PlotArc)
ON (pa.status);

// Foreshadow 节点索引（伏笔）
CREATE INDEX foreshadow_project_idx IF NOT EXISTS
FOR (f:Foreshadow)
ON (f.project_id);

CREATE INDEX foreshadow_status_idx IF NOT EXISTS
FOR (f:Foreshadow)
ON (f.status);

// CharacterState 节点索引（地线）
CREATE INDEX character_state_chapter_idx IF NOT EXISTS
FOR (cs:CharacterState)
ON (cs.chapter_id);

// Goal 节点索引
CREATE INDEX goal_character_idx IF NOT EXISTS
FOR (g:Goal)
ON (g.character_id);

// Convergence 节点索引（三线交汇）
CREATE INDEX convergence_project_idx IF NOT EXISTS
FOR (conv:Convergence)
ON (conv.project_id);

CREATE INDEX convergence_chapter_idx IF NOT EXISTS
FOR (conv:Convergence)
ON (conv.chapter_id);

// =====================================================
// 3. 创建约束（唯一性）
// =====================================================

// Character 唯一约束（project_id + name）
CREATE CONSTRAINT character_unique IF NOT EXISTS
FOR (c:Character)
REQUIRE (c.project_id, c.name) IS UNIQUE;

// Organization 唯一约束
CREATE CONSTRAINT organization_unique IF NOT EXISTS
FOR (o:Organization)
REQUIRE (o.project_id, o.name) IS UNIQUE;

// Location 唯一约束
CREATE CONSTRAINT location_unique IF NOT EXISTS
FOR (l:Location)
REQUIRE (l.project_id, l.name) IS UNIQUE;

// =====================================================
// 4. 全文搜索索引
// =====================================================

// Character 名称全文搜索
CREATE FULLTEXT INDEX character_name_fulltext IF NOT EXISTS
FOR (c:Character)
ON EACH [c.name, c.personality, c.background];

// Location 描述全文搜索
CREATE FULLTEXT INDEX location_description_fulltext IF NOT EXISTS
FOR (l:Location)
ON EACH [l.name, l.description];

// Event 全文搜索
CREATE FULLTEXT INDEX event_fulltext IF NOT EXISTS
FOR (e:Event)
ON EACH [e.name, e.description];

// =====================================================
// 5. 验证索引创建
// =====================================================

CALL db.indexes();

// =====================================================
// 6. 示例查询（测试索引）
// =====================================================

// 查询项目中的所有角色
MATCH (c:Character {project_id: 1})
RETURN c.name, c.role_type
LIMIT 10;

// 查询角色关系
MATCH (c1:Character)-[r]->(c2:Character)
WHERE c1.project_id = 1
RETURN c1.name, type(r), c2.name
LIMIT 20;

// 全文搜索角色
CALL db.index.fulltext.queryNodes(
    'character_name_fulltext', 
    '主角'
) YIELD node, score
RETURN node.name, node.personality, score
LIMIT 5;

// =====================================================
// 完成
// =====================================================

RETURN 'Neo4j 索引初始化完成' AS status;
