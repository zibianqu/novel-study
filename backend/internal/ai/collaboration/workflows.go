package collaboration

import (
	"fmt"
)

// WorkflowTemplate 工作流模板
type WorkflowTemplate struct {
	ID          string
	Name        string
	Description string
	Builder     func(input string, context map[string]interface{}) *Workflow
}

// WorkflowRegistry 工作流注册表
type WorkflowRegistry struct {
	templates map[string]*WorkflowTemplate
}

// NewWorkflowRegistry 创建工作流注册表
func NewWorkflowRegistry() *WorkflowRegistry {
	registry := &WorkflowRegistry{
		templates: make(map[string]*WorkflowTemplate),
	}

	// 注册默认模板
	registry.RegisterDefaultTemplates()

	return registry
}

// RegisterDefaultTemplates 注册默认模板
func (wr *WorkflowRegistry) RegisterDefaultTemplates() {
	// 1. 续写工作流
	wr.Register(&WorkflowTemplate{
		ID:          "continue_write",
		Name:        "续写工作流",
		Description: "由旁白叙述者生成，审核导演审核",
		Builder:     BuildContinueWriteWorkflow,
	})

	// 2. 对话工作流
	wr.Register(&WorkflowTemplate{
		ID:          "dialogue",
		Name:        "对话工作流",
		Description: "由角色扮演者生成，审核导演审核",
		Builder:     BuildDialogueWorkflow,
	})

	// 3. 全文生成工作流
	wr.Register(&WorkflowTemplate{
		ID:          "full_generation",
		Name:        "全文生成工作流",
		Description: "总导演分解任务，多 Agent 协作生成",
		Builder:     BuildFullGenerationWorkflow,
	})

	// 4. 三线规划工作流
	wr.Register(&WorkflowTemplate{
		ID:          "storyline_planning",
		Name:        "三线规划工作流",
		Description: "三个线掌控者协作规划三线",
		Builder:     BuildStorylinePlanningWorkflow,
	})
}

// Register 注册模板
func (wr *WorkflowRegistry) Register(template *WorkflowTemplate) {
	wr.templates[template.ID] = template
}

// Get 获取模板
func (wr *WorkflowRegistry) Get(id string) (*WorkflowTemplate, error) {
	template, ok := wr.templates[id]
	if !ok {
		return nil, fmt.Errorf("workflow template %s not found", id)
	}
	return template, nil
}

// List 列出所有模板
func (wr *WorkflowRegistry) List() []*WorkflowTemplate {
	templates := make([]*WorkflowTemplate, 0, len(wr.templates))
	for _, t := range wr.templates {
		templates = append(templates, t)
	}
	return templates
}

// BuildContinueWriteWorkflow 构建续写工作流
func BuildContinueWriteWorkflow(input string, context map[string]interface{}) *Workflow {
	workflow := NewWorkflow("wf_continue_"+generateID(), "续写工作流", "续写小说内容")

	// Task 1: 旁白叙述者生成
	task1 := &AgentTask{
		ID:      "task_generate",
		AgentID: 1, // Agent 1 - 旁白叙述者
		Type:    "generate",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 2: 审核导演审核
	task2 := &AgentTask{
		ID:        "task_review",
		AgentID:   3, // Agent 3 - 审核导演
		Type:      "review",
		Input:     "", // 将从 task1 获取
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_generate"},
	}

	workflow.AddTask(task1)
	workflow.AddTask(task2)

	return workflow
}

// BuildDialogueWorkflow 构建对话工作流
func BuildDialogueWorkflow(input string, context map[string]interface{}) *Workflow {
	workflow := NewWorkflow("wf_dialogue_"+generateID(), "对话工作流", "生成角色对话")

	// Task 1: 角色扮演者生成对话
	task1 := &AgentTask{
		ID:      "task_dialogue",
		AgentID: 2, // Agent 2 - 角色扮演者
		Type:    "generate",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 2: 审核导演审核
	task2 := &AgentTask{
		ID:        "task_review",
		AgentID:   3,
		Type:      "review",
		Input:     "",
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_dialogue"},
	}

	workflow.AddTask(task1)
	workflow.AddTask(task2)

	return workflow
}

// BuildFullGenerationWorkflow 构建全文生成工作流
func BuildFullGenerationWorkflow(input string, context map[string]interface{}) *Workflow {
	workflow := NewWorkflow("wf_full_"+generateID(), "全文生成工作流", "多 Agent 协作生成全文")

	// Task 1: 总导演分解任务
	task1 := &AgentTask{
		ID:      "task_plan",
		AgentID: 0, // Agent 0 - 总导演
		Type:    "analyze",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 2: 旁白叙述者生成描写
	task2 := &AgentTask{
		ID:        "task_narration",
		AgentID:   1,
		Type:      "generate",
		Input:     "",
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_plan"},
	}

	// Task 3: 角色扮演者生成对话
	task3 := &AgentTask{
		ID:        "task_dialogue",
		AgentID:   2,
		Type:      "generate",
		Input:     "",
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_plan"},
	}

	// Task 4: 审核导演整合审核
	task4 := &AgentTask{
		ID:        "task_review",
		AgentID:   3,
		Type:      "review",
		Input:     "",
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_narration", "task_dialogue"},
	}

	workflow.AddTask(task1)
	workflow.AddTask(task2)
	workflow.AddTask(task3)
	workflow.AddTask(task4)

	return workflow
}

// BuildStorylinePlanningWorkflow 构建三线规划工作流
func BuildStorylinePlanningWorkflow(input string, context map[string]interface{}) *Workflow {
	workflow := NewWorkflow("wf_storyline_"+generateID(), "三线规划工作流", "规划天地剧三线")

	// Task 1: 天线掌控者
	task1 := &AgentTask{
		ID:      "task_skyline",
		AgentID: 4, // Agent 4 - 天线掌控者
		Type:    "analyze",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 2: 地线掌控者
	task2 := &AgentTask{
		ID:      "task_groundline",
		AgentID: 5, // Agent 5 - 地线掌控者
		Type:    "analyze",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 3: 剧情线掌控者
	task3 := &AgentTask{
		ID:      "task_plotline",
		AgentID: 6, // Agent 6 - 剧情线掌控者
		Type:    "analyze",
		Input:   input,
		Context: context,
		Status:  "pending",
	}

	// Task 4: 总导演整合
	task4 := &AgentTask{
		ID:        "task_integrate",
		AgentID:   0,
		Type:      "analyze",
		Input:     "",
		Context:   context,
		Status:    "pending",
		DependsOn: []string{"task_skyline", "task_groundline", "task_plotline"},
	}

	workflow.AddTask(task1)
	workflow.AddTask(task2)
	workflow.AddTask(task3)
	workflow.AddTask(task4)

	return workflow
}

// 辅助函数

var idCounter uint64

func generateID() string {
	idCounter++
	return fmt.Sprintf("%d", idCounter)
}
