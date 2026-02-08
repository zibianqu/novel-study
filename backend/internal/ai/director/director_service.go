package director

import (
	"context"
	"fmt"
)

// DirectorService 增强的总导演服务
type DirectorService struct {
	intentAnalyzer      *IntentAnalyzer
	taskDecomposer      *TaskDecomposer
	conflictArbitrator  *ConflictArbitrator
}

// NewDirectorService 创建总导演服务
func NewDirectorService() *DirectorService {
	intentAnalyzer := NewIntentAnalyzer()
	taskDecomposer := NewTaskDecomposer(intentAnalyzer)
	conflictArbitrator := NewConflictArbitrator()

	return &DirectorService{
		intentAnalyzer:     intentAnalyzer,
		taskDecomposer:     taskDecomposer,
		conflictArbitrator: conflictArbitrator,
	}
}

// ProcessRequest 处理用户请求
func (ds *DirectorService) ProcessRequest(
	ctx context.Context,
	userInput string,
	context map[string]interface{},
) (*DirectorResponse, error) {
	// 1. 分析意图
	intent, err := ds.intentAnalyzer.Analyze(ctx, userInput)
	if err != nil {
		return nil, fmt.Errorf("intent analysis failed: %w", err)
	}

	// 2. 分解任务
	plan, err := ds.taskDecomposer.Decompose(ctx, intent, userInput, context)
	if err != nil {
		return nil, fmt.Errorf("task decomposition failed: %w", err)
	}

	// 3. 优化计划
	plan = ds.taskDecomposer.OptimizePlan(plan)

	// 4. 构建响应
	return &DirectorResponse{
		Intent:       intent,
		Plan:         plan,
		WorkflowID:   ds.intentAnalyzer.GetWorkflowTemplate(intent),
		RequiredAgents: ds.intentAnalyzer.GetRequiredAgents(intent),
	}, nil
}

// ResolveConflicts 解决冲突
func (ds *DirectorService) ResolveConflicts(
	ctx context.Context,
	agentOutputs map[int]string,
) ([]*Resolution, error) {
	// 1. 检测冲突
	conflicts, err := ds.conflictArbitrator.DetectConflict(ctx, agentOutputs)
	if err != nil {
		return nil, fmt.Errorf("conflict detection failed: %w", err)
	}

	if len(conflicts) == 0 {
		return nil, nil
	}

	// 2. 仲裁每个冲突
	resolutions := make([]*Resolution, 0)
	for _, conflict := range conflicts {
		resolution, err := ds.conflictArbitrator.Arbitrate(ctx, conflict, agentOutputs)
		if err != nil {
			return nil, fmt.Errorf("conflict arbitration failed: %w", err)
		}
		resolutions = append(resolutions, resolution)
	}

	return resolutions, nil
}

// MakeDecision 做出决策
func (ds *DirectorService) MakeDecision(
	ctx context.Context,
	options []string,
	criteria map[string]interface{},
) (*Decision, error) {
	// 基于标准评估选项
	scores := make(map[int]float64)

	for i := range options {
		score := ds.evaluateOption(i, options[i], criteria)
		scores[i] = score
	}

	// 选择得分最高的
	bestOption := 0
	bestScore := 0.0

	for i, score := range scores {
		if score > bestScore {
			bestScore = score
			bestOption = i
		}
	}

	return &Decision{
		ChosenOption: bestOption,
		Score:        bestScore,
		Scores:       scores,
		Reason:       fmt.Sprintf("选项 %d 得分最高: %.2f", bestOption, bestScore),
	}, nil
}

// evaluateOption 评估选项
func (ds *DirectorService) evaluateOption(
	index int,
	option string,
	criteria map[string]interface{},
) float64 {
	score := 50.0 // 基础分

	// 根据长度加分
	length := len([]rune(option))
	if targetLength, ok := criteria["target_length"].(int); ok {
		lengthDiff := float64(abs(length - targetLength))
		score -= lengthDiff / 10.0
	}

	// 根据关键词加分
	if keywords, ok := criteria["keywords"].([]string); ok {
		for _, keyword := range keywords {
			if contains(option, keyword) {
				score += 5.0
			}
		}
	}

	return score
}

// CoordinateAgents 协调 Agent
func (ds *DirectorService) CoordinateAgents(
	ctx context.Context,
	plan *DecompositionPlan,
) (*CoordinationResult, error) {
	// 这个方法将与协作调度器集成
	return &CoordinationResult{
		Success:    true,
		TotalTasks: plan.TotalTasks,
		Strategy:   plan.Strategy,
	}, nil
}

// DirectorResponse 总导演响应
type DirectorResponse struct {
	Intent         *Intent
	Plan           *DecompositionPlan
	WorkflowID     string
	RequiredAgents []int
}

// Decision 决策结果
type Decision struct {
	ChosenOption int
	Score        float64
	Scores       map[int]float64
	Reason       string
}

// CoordinationResult 协调结果
type CoordinationResult struct {
	Success    bool
	TotalTasks int
	Strategy   string
}

// GetIntentAnalyzer 获取意图分析器
func (ds *DirectorService) GetIntentAnalyzer() *IntentAnalyzer {
	return ds.intentAnalyzer
}

// GetTaskDecomposer 获取任务分解器
func (ds *DirectorService) GetTaskDecomposer() *TaskDecomposer {
	return ds.taskDecomposer
}

// GetConflictArbitrator 获取冲突仲裁器
func (ds *DirectorService) GetConflictArbitrator() *ConflictArbitrator {
	return ds.conflictArbitrator
}

// 辅助函数

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) > len(substr) && 
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr))
}
