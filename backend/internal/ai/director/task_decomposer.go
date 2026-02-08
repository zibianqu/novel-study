package director

import (
	"context"
	"fmt"
)

// SubTask 子任务
type SubTask struct {
	ID          string
	Type        string // "narration", "dialogue", "description", "action"
	Description string
	AgentID     int
	Priority    int
	EstimatedLength int
	Dependencies []string
	Context     map[string]interface{}
}

// DecompositionPlan 分解计划
type DecompositionPlan struct {
	TotalTasks int
	SubTasks   []*SubTask
	Strategy   string // "sequential", "parallel", "mixed"
	EstimatedTime int // 预估时间(秒)
}

// TaskDecomposer 任务分解器
type TaskDecomposer struct {
	intentAnalyzer *IntentAnalyzer
}

// NewTaskDecomposer 创建任务分解器
func NewTaskDecomposer(intentAnalyzer *IntentAnalyzer) *TaskDecomposer {
	return &TaskDecomposer{
		intentAnalyzer: intentAnalyzer,
	}
}

// Decompose 分解任务
func (td *TaskDecomposer) Decompose(
	ctx context.Context,
	intent *Intent,
	input string,
	context map[string]interface{},
) (*DecompositionPlan, error) {
	switch intent.Complexity {
	case "simple":
		return td.decomposeSimple(intent, input, context)
	case "medium":
		return td.decomposeMedium(intent, input, context)
	case "complex":
		return td.decomposeComplex(intent, input, context)
	default:
		return td.decomposeMedium(intent, input, context)
	}
}

// decomposeSimple 分解简单任务
func (td *TaskDecomposer) decomposeSimple(
	intent *Intent,
	input string,
	context map[string]interface{},
) (*DecompositionPlan, error) {
	// 简单任务：单个 Agent 处理
	agents := td.intentAnalyzer.GetRequiredAgents(intent)

	task := &SubTask{
		ID:          "task_1",
		Type:        intent.Type,
		Description: input,
		AgentID:     agents[0],
		Priority:    10,
		EstimatedLength: td.estimateLength(intent),
		Context:     context,
	}

	return &DecompositionPlan{
		TotalTasks:    1,
		SubTasks:      []*SubTask{task},
		Strategy:      "sequential",
		EstimatedTime: 30,
	}, nil
}

// decomposeMedium 分解中等任务
func (td *TaskDecomposer) decomposeMedium(
	intent *Intent,
	input string,
	context map[string]interface{},
) (*DecompositionPlan, error) {
	// 中等任务：生成 + 审核
	agents := td.intentAnalyzer.GetRequiredAgents(intent)

	subtasks := make([]*SubTask, 0)

	// Task 1: 生成
	task1 := &SubTask{
		ID:          "task_generate",
		Type:        "generate",
		Description: input,
		AgentID:     agents[0],
		Priority:    10,
		EstimatedLength: td.estimateLength(intent),
		Context:     context,
	}
	subtasks = append(subtasks, task1)

	// Task 2: 审核
	if len(agents) > 1 {
		task2 := &SubTask{
			ID:          "task_review",
			Type:        "review",
			Description: "审核生成的内容",
			AgentID:     agents[1],
			Priority:    9,
			Dependencies: []string{"task_generate"},
			Context:     context,
		}
		subtasks = append(subtasks, task2)
	}

	return &DecompositionPlan{
		TotalTasks:    len(subtasks),
		SubTasks:      subtasks,
		Strategy:      "sequential",
		EstimatedTime: 60,
	}, nil
}

// decomposeComplex 分解复杂任务
func (td *TaskDecomposer) decomposeComplex(
	intent *Intent,
	input string,
	context map[string]interface{},
) (*DecompositionPlan, error) {
	// 复杂任务：多 Agent 协作
	subtasks := make([]*SubTask, 0)

	switch intent.Type {
	case "plan":
		// 三线规划
		subtasks = td.decomposePlanningTask(input, context)

	default:
		// 全文生成
		subtasks = td.decomposeFullGenerationTask(input, context)
	}

	return &DecompositionPlan{
		TotalTasks:    len(subtasks),
		SubTasks:      subtasks,
		Strategy:      "mixed",
		EstimatedTime: 120,
	}, nil
}

// decomposePlanningTask 分解规划任务
func (td *TaskDecomposer) decomposePlanningTask(
	input string,
	context map[string]interface{},
) []*SubTask {
	subtasks := make([]*SubTask, 0)

	// 并行任务：三线分析
	subtasks = append(subtasks, &SubTask{
		ID:          "task_skyline",
		Type:        "analyze",
		Description: "分析天线（世界大势）",
		AgentID:     4,
		Priority:    10,
		Context:     context,
	})

	subtasks = append(subtasks, &SubTask{
		ID:          "task_groundline",
		Type:        "analyze",
		Description: "分析地线（主角成长）",
		AgentID:     5,
		Priority:    10,
		Context:     context,
	})

	subtasks = append(subtasks, &SubTask{
		ID:          "task_plotline",
		Type:        "analyze",
		Description: "分析剧情线（当前情节）",
		AgentID:     6,
		Priority:    10,
		Context:     context,
	})

	// 顺序任务：整合
	subtasks = append(subtasks, &SubTask{
		ID:          "task_integrate",
		Type:        "analyze",
		Description: "整合三线分析结果",
		AgentID:     0,
		Priority:    9,
		Dependencies: []string{"task_skyline", "task_groundline", "task_plotline"},
		Context:     context,
	})

	return subtasks
}

// decomposeFullGenerationTask 分解全文生成任务
func (td *TaskDecomposer) decomposeFullGenerationTask(
	input string,
	context map[string]interface{},
) []*SubTask {
	subtasks := make([]*SubTask, 0)

	// Task 1: 规划
	subtasks = append(subtasks, &SubTask{
		ID:          "task_plan",
		Type:        "analyze",
		Description: "分析并制定生成计划",
		AgentID:     0,
		Priority:    10,
		Context:     context,
	})

	// Task 2 & 3: 并行生成
	subtasks = append(subtasks, &SubTask{
		ID:          "task_narration",
		Type:        "generate",
		Description: "生成旁白叙述",
		AgentID:     1,
		Priority:    9,
		Dependencies: []string{"task_plan"},
		Context:     context,
	})

	subtasks = append(subtasks, &SubTask{
		ID:          "task_dialogue",
		Type:        "generate",
		Description: "生成角色对话",
		AgentID:     2,
		Priority:    9,
		Dependencies: []string{"task_plan"},
		Context:     context,
	})

	// Task 4: 审核
	subtasks = append(subtasks, &SubTask{
		ID:          "task_review",
		Type:        "review",
		Description: "整合并审核内容",
		AgentID:     3,
		Priority:    8,
		Dependencies: []string{"task_narration", "task_dialogue"},
		Context:     context,
	})

	return subtasks
}

// estimateLength 估计生成长度
func (td *TaskDecomposer) estimateLength(intent *Intent) int {
	// 从参数中获取
	if length, ok := intent.Parameters["length"].(int); ok {
		return length
	}

	// 默认值
	switch intent.Complexity {
	case "simple":
		return 300
	case "medium":
		return 500
	case "complex":
		return 1000
	default:
		return 500
	}
}

// OptimizePlan 优化计划
func (td *TaskDecomposer) OptimizePlan(plan *DecompositionPlan) *DecompositionPlan {
	// 1. 识别可并行任务
	plan.Strategy = td.determineStrategy(plan.SubTasks)

	// 2. 调整优先级
	plan.SubTasks = td.adjustPriorities(plan.SubTasks)

	// 3. 重新估计时间
	plan.EstimatedTime = td.estimateTime(plan)

	return plan
}

// determineStrategy 确定执行策略
func (td *TaskDecomposer) determineStrategy(tasks []*SubTask) string {
	hasParallel := false
	hasSequential := false

	for _, task := range tasks {
		if len(task.Dependencies) == 0 {
			hasParallel = true
		} else {
			hasSequential = true
		}
	}

	if hasParallel && hasSequential {
		return "mixed"
	} else if hasParallel {
		return "parallel"
	}
	return "sequential"
}

// adjustPriorities 调整优先级
func (td *TaskDecomposer) adjustPriorities(tasks []*SubTask) []*SubTask {
	// 根据依赖关系调整优先级
	for _, task := range tasks {
		if len(task.Dependencies) > 0 {
			// 有依赖的任务优先级降低
			task.Priority = task.Priority - 1
		}
	}
	return tasks
}

// estimateTime 估计总时间
func (td *TaskDecomposer) estimateTime(plan *DecompositionPlan) int {
	baseTime := 30 // 每个任务基础时间

	switch plan.Strategy {
	case "sequential":
		return baseTime * plan.TotalTasks
	case "parallel":
		return baseTime * 2 // 并行执行
	case "mixed":
		return baseTime * (plan.TotalTasks / 2 + 1)
	default:
		return baseTime * plan.TotalTasks
	}
}

// GetPlanSummary 获取计划摘要
func (td *TaskDecomposer) GetPlanSummary(plan *DecompositionPlan) string {
	return fmt.Sprintf(
		"任务分解计划: 总共 %d 个子任务，执行策略: %s，预估耗时: %d秒",
		plan.TotalTasks,
		plan.Strategy,
		plan.EstimatedTime,
	)
}
