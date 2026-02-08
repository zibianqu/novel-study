// NovelForge AI - Neo4j 索引初始化脚本
// 执行时间: 2026-02-08

// ====================================
// 节点索引
// ====================================

// Character 角色节点
CREATE INDEX character_project_id IF NOT EXISTS FOR (c:Character) ON (c.project_id);
CREATE INDEX character_name IF NOT EXISTS FOR (c:Character) ON (c.name);
CREATE INDEX character_role_type IF NOT EXISTS FOR (c:Character) ON (c.role_type);
CREATE CONSTRAINT character_unique IF NOT EXISTS FOR (c:Character) REQUIRE (c.id) IS UNIQUE;

// Organization 组织节点
CREATE INDEX organization_project_id IF NOT EXISTS FOR (o:Organization) ON (o.project_id);
CREATE INDEX organization_name IF NOT EXISTS FOR (o:Organization) ON (o.name);
CREATE CONSTRAINT organization_unique IF NOT EXISTS FOR (o:Organization) REQUIRE (o.id) IS UNIQUE;

// Location 地点节点
CREATE INDEX location_project_id IF NOT EXISTS FOR (l:Location) ON (l.project_id);
CREATE INDEX location_name IF NOT EXISTS FOR (l:Location) ON (l.name);
CREATE CONSTRAINT location_unique IF NOT EXISTS FOR (l:Location) REQUIRE (l.id) IS UNIQUE;

// Event 事件节点
CREATE INDEX event_project_id IF NOT EXISTS FOR (e:Event) ON (e.project_id);
CREATE INDEX event_chapter_id IF NOT EXISTS FOR (e:Event) ON (e.chapter_id);
CREATE INDEX event_time IF NOT EXISTS FOR (e:Event) ON (e.time_point);
CREATE CONSTRAINT event_unique IF NOT EXISTS FOR (e:Event) REQUIRE (e.id) IS UNIQUE;

// WorldEvent 世界事件（天线）
CREATE INDEX world_event_project_id IF NOT EXISTS FOR (we:WorldEvent) ON (we.project_id);
CREATE INDEX world_event_impact IF NOT EXISTS FOR (we:WorldEvent) ON (we.impact_level);
CREATE CONSTRAINT world_event_unique IF NOT EXISTS FOR (we:WorldEvent) REQUIRE (we.id) IS UNIQUE;

// PlotArc 剧情弧
CREATE INDEX plot_arc_project_id IF NOT EXISTS FOR (pa:PlotArc) ON (pa.project_id);
CREATE INDEX plot_arc_status IF NOT EXISTS FOR (pa:PlotArc) ON (pa.status);
CREATE CONSTRAINT plot_arc_unique IF NOT EXISTS FOR (pa:PlotArc) REQUIRE (pa.id) IS UNIQUE;

// Foreshadow 伏笔节点
CREATE INDEX foreshadow_project_id IF NOT EXISTS FOR (f:Foreshadow) ON (f.project_id);
CREATE INDEX foreshadow_status IF NOT EXISTS FOR (f:Foreshadow) ON (f.status);
CREATE CONSTRAINT foreshadow_unique IF NOT EXISTS FOR (f:Foreshadow) REQUIRE (f.id) IS UNIQUE;

// CharacterState 角色状态（地线）
CREATE INDEX character_state_character IF NOT EXISTS FOR (cs:CharacterState) ON (cs.character_id);
CREATE INDEX character_state_chapter IF NOT EXISTS FOR (cs:CharacterState) ON (cs.chapter_id);

// Goal 目标节点
CREATE INDEX goal_character IF NOT EXISTS FOR (g:Goal) ON (g.character_id);
CREATE INDEX goal_status IF NOT EXISTS FOR (g:Goal) ON (g.status);

// Convergence 三线交汇点
CREATE INDEX convergence_project_id IF NOT EXISTS FOR (c:Convergence) ON (c.project_id);
CREATE INDEX convergence_chapter IF NOT EXISTS FOR (c:Convergence) ON (c.chapter_id);
CREATE CONSTRAINT convergence_unique IF NOT EXISTS FOR (c:Convergence) REQUIRE (c.id) IS UNIQUE;

// ====================================
// 全文搜索索引
// ====================================

// 创建全文索引用于名称搜索
CREATE FULLTEXT INDEX character_name_fulltext IF NOT EXISTS 
FOR (c:Character) ON EACH [c.name, c.personality];

CREATE FULLTEXT INDEX location_name_fulltext IF NOT EXISTS 
FOR (l:Location) ON EACH [l.name, l.description];

CREATE FULLTEXT INDEX event_content_fulltext IF NOT EXISTS 
FOR (e:Event) ON EACH [e.name, e.description];

// ====================================
// 验证索引创建
// ====================================

// 列出所有索引
SHOW INDEXES;

// 列出所有约束
SHOW CONSTRAINTS;
