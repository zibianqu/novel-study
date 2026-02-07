package model

import "time"

// ==================== 用户 ====================

type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Avatar          string    `json:"avatar"`
	Settings        any       `json:"settings"`
	APIKeyEncrypted string    `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// ==================== 项目 ====================

type Project struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	Genre        string    `json:"genre"`
	Description  string    `json:"description"`
	CoverImage   string    `json:"cover_image"`
	Status       string    `json:"status"`
	WordCount    int       `json:"word_count"`
	ChapterCount int       `json:"chapter_count,omitempty"`
	Settings     any       `json:"settings"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=novel_long novel_short copywriting"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
}

// ==================== 卷/章节 ====================

type Volume struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	SortOrder int       `json:"sort_order"`
	Chapters  []Chapter `json:"chapters,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Chapter struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	VolumeID  *int      `json:"volume_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content,omitempty"`
	WordCount int       `json:"word_count"`
	SortOrder int       `json:"sort_order"`
	Status    string    `json:"status"`
	LockedBy  *int      `json:"locked_by"`
	LockedAt  *time.Time `json:"locked_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateChapterRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	VolumeID  *int   `json:"volume_id"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content"`
}

type UpdateChapterRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Status  *string `json:"status"`
}

type ChapterVersion struct {
	ID            int       `json:"id"`
	ChapterID     int       `json:"chapter_id"`
	VersionNum    int       `json:"version_num"`
	Content       string    `json:"content"`
	DeltaContent  string    `json:"delta_content"`
	DeltaPosition any       `json:"delta_position"`
	AgentOutputs  any       `json:"agent_outputs"`
	EmbeddingIDs  []int     `json:"embedding_ids"`
	GraphChanges  any       `json:"graph_changes"`
	CreatedBy     *int      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

// ==================== 角色 ====================

type Character struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"project_id"`
	Name        string    `json:"name"`
	Avatar      string    `json:"avatar"`
	RoleType    string    `json:"role_type"`
	Personality string    `json:"personality"`
	Appearance  string    `json:"appearance"`
	Background  string    `json:"background"`
	Abilities   string    `json:"abilities"`
	Motivation  string    `json:"motivation"`
	SpeechStyle string    `json:"speech_style"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ==================== Agent ====================

type Agent struct {
	ID             int       `json:"id"`
	UserID         *int      `json:"user_id"`
	AgentKey       string    `json:"agent_key"`
	Name           string    `json:"name"`
	Icon           string    `json:"icon"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	Layer          string    `json:"layer"`
	SystemPrompt   string    `json:"system_prompt"`
	Model          string    `json:"model"`
	Temperature    float64   `json:"temperature"`
	MaxTokens      int       `json:"max_tokens"`
	Tools          any       `json:"tools"`
	InputSchema    any       `json:"input_schema"`
	OutputSchema   any       `json:"output_schema"`
	Permissions    any       `json:"permissions"`
	IsActive       bool      `json:"is_active"`
	SortOrder      int       `json:"sort_order"`
	KnowledgeCount int       `json:"knowledge_count,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateAgentRequest struct {
	AgentKey     string  `json:"agent_key" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Icon         string  `json:"icon"`
	Description  string  `json:"description"`
	Layer        string  `json:"layer" binding:"required,oneof=decision strategy execution quality auxiliary"`
	SystemPrompt string  `json:"system_prompt" binding:"required"`
	Model        string  `json:"model"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
}

// ==================== 知识库 ====================

type KnowledgeCategory struct {
	ID          int                  `json:"id"`
	AgentID     int                  `json:"agent_id"`
	ParentID    *int                 `json:"parent_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	SortOrder   int                  `json:"sort_order"`
	ItemCount   int                  `json:"item_count,omitempty"`
	Children    []KnowledgeCategory  `json:"children,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
}

type KnowledgeItem struct {
	ID           int       `json:"id"`
	AgentID      int       `json:"agent_id"`
	CategoryID   *int      `json:"category_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Tags         []string  `json:"tags"`
	Source       string    `json:"source"`
	QualityScore float64   `json:"quality_score"`
	UseCount     int       `json:"use_count"`
	IsActive     bool      `json:"is_active"`
	CreatedBy    *int      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateKnowledgeItemRequest struct {
	AgentID    int      `json:"agent_id" binding:"required"`
	CategoryID *int     `json:"category_id"`
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Tags       []string `json:"tags"`
}

// ==================== 工作流 ====================

type Workflow struct {
	ID          int            `json:"id"`
	UserID      *int           `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Category    string         `json:"category"`
	Icon        string         `json:"icon"`
	IsActive    bool           `json:"is_active"`
	Version     int            `json:"version"`
	Nodes       []WorkflowNode `json:"nodes,omitempty"`
	Edges       []WorkflowEdge `json:"edges,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type WorkflowNode struct {
	ID         int    `json:"id"`
	WorkflowID int    `json:"workflow_id"`
	NodeKey    string `json:"node_key"`
	NodeType   string `json:"node_type"`
	AgentID    *int   `json:"agent_id"`
	Name       string `json:"name"`
	Config     any    `json:"config"`
	PositionX  int    `json:"position_x"`
	PositionY  int    `json:"position_y"`
	SortOrder  int    `json:"sort_order"`
}

type WorkflowEdge struct {
	ID            int    `json:"id"`
	WorkflowID    int    `json:"workflow_id"`
	FromNodeID    int    `json:"from_node_id"`
	ToNodeID      int    `json:"to_node_id"`
	EdgeType      string `json:"edge_type"`
	ConditionExpr any    `json:"condition_expr"`
	Label         string `json:"label"`
	SortOrder     int    `json:"sort_order"`
}

type WorkflowExecution struct {
	ID            int        `json:"id"`
	WorkflowID    int        `json:"workflow_id"`
	ProjectID     *int       `json:"project_id"`
	UserID        int        `json:"user_id"`
	Status        string     `json:"status"`
	InputData     any        `json:"input_data"`
	OutputData    any        `json:"output_data"`
	CurrentNodeID *int       `json:"current_node_id"`
	ErrorMessage  string     `json:"error_message"`
	StartedAt     time.Time  `json:"started_at"`
	CompletedAt   *time.Time `json:"completed_at"`
}

// ==================== 三线 ====================

type Storyline struct {
	ID           int             `json:"id"`
	ProjectID    int             `json:"project_id"`
	LineType     string          `json:"line_type"`
	Title        string          `json:"title"`
	Content      string          `json:"content"`
	ChapterStart int             `json:"chapter_start"`
	ChapterEnd   int             `json:"chapter_end"`
	Status       string          `json:"status"`
	SortOrder    int             `json:"sort_order"`
	ParentID     *int            `json:"parent_id"`
	Children     []Storyline     `json:"children,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ==================== AI ====================

type ChatRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	SessionID string `json:"session_id"`
}

type AIWriteRequest struct {
	ProjectID int    `json:"project_id" binding:"required"`
	ChapterID int    `json:"chapter_id" binding:"required"`
	Action    string `json:"action" binding:"required"` // continue/polish/rewrite/dialogue
	Context   string `json:"context"`
	Instruction string `json:"instruction"`
}

// ==================== 通用 ====================

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PaginationQuery struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func (p *PaginationQuery) Defaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.PageSize
}
