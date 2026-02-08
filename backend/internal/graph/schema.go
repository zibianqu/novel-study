package graph

import "time"

// NodeType 节点类型
type NodeType string

const (
	NodeTypeCharacter NodeType = "Character" // 人物
	NodeTypeLocation  NodeType = "Location"  // 地点
	NodeTypeEvent     NodeType = "Event"     // 事件
	NodeTypeItem      NodeType = "Item"      // 物品
	NodeTypeConcept   NodeType = "Concept"   // 概念
)

// RelationType 关系类型
type RelationType string

const (
	// 人物关系
	RelationKnows      RelationType = "KNOWS"       // 认识
	RelationFamilyOf   RelationType = "FAMILY_OF"   // 亲属
	RelationMasterOf   RelationType = "MASTER_OF"   // 师徒
	RelationEnemyOf    RelationType = "ENEMY_OF"    // 仇敵
	RelationAllyOf     RelationType = "ALLY_OF"     // 盟友
	RelationLoves      RelationType = "LOVES"       // 爱慕

	// 位置关系
	RelationLocatedAt  RelationType = "LOCATED_AT"  // 位于
	RelationBornAt     RelationType = "BORN_AT"     // 出生于
	RelationLivesIn    RelationType = "LIVES_IN"    // 居住于

	// 事件关系
	RelationHappensAt  RelationType = "HAPPENS_AT"  // 发生于
	RelationParticipates RelationType = "PARTICIPATES" // 参与
	RelationCauses     RelationType = "CAUSES"      // 导致
	RelationLeadsTo    RelationType = "LEADS_TO"    // 引导至

	// 物品关系
	RelationOwns       RelationType = "OWNS"        // 拥有
	RelationUses       RelationType = "USES"        // 使用
	RelationCreates    RelationType = "CREATES"     // 创造

	// 概念关系
	RelationMasters    RelationType = "MASTERS"     // 掌握
	RelationBelongsTo  RelationType = "BELONGS_TO"  // 属于
)

// Node 节点基类
type Node struct {
	ID          string                 `json:"id"`
	Type        NodeType               `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Character 人物节点
type Character struct {
	Node
	Role        string   `json:"role"`         // 角色: protagonist, supporting, minor
	Age         int      `json:"age,omitempty"`
	Gender      string   `json:"gender,omitempty"`
	Personality []string `json:"personality,omitempty"`
	Appearance  string   `json:"appearance,omitempty"`
	Background  string   `json:"background,omitempty"`
}

// Location 地点节点
type Location struct {
	Node
	LocationType string `json:"location_type"` // city, building, room, etc.
	Climate      string `json:"climate,omitempty"`
	Culture      string `json:"culture,omitempty"`
}

// Event 事件节点
type Event struct {
	Node
	Timestamp   time.Time `json:"timestamp"`
	Chapter     int       `json:"chapter"`
	EventType   string    `json:"event_type"` // battle, meeting, discovery, etc.
	Importance  int       `json:"importance"` // 1-10
	Consequence string    `json:"consequence,omitempty"`
}

// Item 物品节点
type Item struct {
	Node
	ItemType   string `json:"item_type"` // weapon, treasure, artifact, etc.
	Power      int    `json:"power,omitempty"`
	Rarity     string `json:"rarity,omitempty"` // common, rare, legendary, etc.
	Origin     string `json:"origin,omitempty"`
}

// Concept 概念节点
type Concept struct {
	Node
	ConceptType string `json:"concept_type"` // skill, magic, organization, etc.
	Level       int    `json:"level,omitempty"`
	Category    string `json:"category,omitempty"`
}

// Relationship 关系
type Relationship struct {
	ID         string                 `json:"id"`
	Type       RelationType           `json:"type"`
	FromNodeID string                 `json:"from_node_id"`
	ToNodeID   string                 `json:"to_node_id"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Weight     float64                `json:"weight,omitempty"` // 关系强度 0-1
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// GraphQuery 图查询请求
type GraphQuery struct {
	NodeID     string                 `json:"node_id,omitempty"`
	NodeType   NodeType               `json:"node_type,omitempty"`
	RelType    RelationType           `json:"rel_type,omitempty"`
	Depth      int                    `json:"depth,omitempty"` // 查询深度
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
}

// GraphResult 图查询结果
type GraphResult struct {
	Nodes         []*Node         `json:"nodes"`
	Relationships []*Relationship `json:"relationships"`
	Paths         []Path          `json:"paths,omitempty"`
}

// Path 图路径
type Path struct {
	Nodes         []*Node         `json:"nodes"`
	Relationships []*Relationship `json:"relationships"`
	Length        int             `json:"length"`
}

// NodeBuilder 节点构建器接口
type NodeBuilder interface {
	Build() *Node
}

// CharacterBuilder 人物构建器
type CharacterBuilder struct {
	character *Character
}

func NewCharacterBuilder(id, name string) *CharacterBuilder {
	return &CharacterBuilder{
		character: &Character{
			Node: Node{
				ID:         id,
				Type:       NodeTypeCharacter,
				Name:       name,
				Properties: make(map[string]interface{}),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
}

func (b *CharacterBuilder) WithRole(role string) *CharacterBuilder {
	b.character.Role = role
	return b
}

func (b *CharacterBuilder) WithAge(age int) *CharacterBuilder {
	b.character.Age = age
	return b
}

func (b *CharacterBuilder) WithGender(gender string) *CharacterBuilder {
	b.character.Gender = gender
	return b
}

func (b *CharacterBuilder) WithDescription(desc string) *CharacterBuilder {
	b.character.Description = desc
	return b
}

func (b *CharacterBuilder) Build() *Character {
	return b.character
}

// LocationBuilder 地点构建器
type LocationBuilder struct {
	location *Location
}

func NewLocationBuilder(id, name string) *LocationBuilder {
	return &LocationBuilder{
		location: &Location{
			Node: Node{
				ID:         id,
				Type:       NodeTypeLocation,
				Name:       name,
				Properties: make(map[string]interface{}),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		},
	}
}

func (b *LocationBuilder) WithLocationType(locType string) *LocationBuilder {
	b.location.LocationType = locType
	return b
}

func (b *LocationBuilder) WithDescription(desc string) *LocationBuilder {
	b.location.Description = desc
	return b
}

func (b *LocationBuilder) Build() *Location {
	return b.location
}

// EventBuilder 事件构建器
type EventBuilder struct {
	event *Event
}

func NewEventBuilder(id, name string) *EventBuilder {
	return &EventBuilder{
		event: &Event{
			Node: Node{
				ID:         id,
				Type:       NodeTypeEvent,
				Name:       name,
				Properties: make(map[string]interface{}),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			Timestamp: time.Now(),
		},
	}
}

func (b *EventBuilder) WithChapter(chapter int) *EventBuilder {
	b.event.Chapter = chapter
	return b
}

func (b *EventBuilder) WithEventType(eventType string) *EventBuilder {
	b.event.EventType = eventType
	return b
}

func (b *EventBuilder) WithImportance(importance int) *EventBuilder {
	b.event.Importance = importance
	return b
}

func (b *EventBuilder) Build() *Event {
	return b.event
}
