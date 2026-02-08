package collaboration

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AgentTask Agent 任务
type AgentTask struct {
	ID          string
	AgentID     int
	Type        string                 // "generate", "review", "revise", "analyze"
	Input       string                 // 输入内容
	Context     map[string]interface{} // 上下文
	Status      string                 // "pending", "running", "completed", "failed"
	Result      string                 // 输出结果
	Error       error
	StartTime   time.Time
	EndTime     time.Time
	DependsOn   []string // 依赖的任务 ID
}

// WorkflowResult 工作流结果
type WorkflowResult struct {
	Success      bool
	FinalContent string
	Tasks        []*AgentTask
	TotalTime    time.Duration
	Metadata     map[string]interface{}
}

// AgentExecutor Agent 执行器接口
type AgentExecutor interface {
	Execute(ctx context.Context, agentID int, input string, context map[string]interface{}) (string, error)
}

// Scheduler Agent 协作调度器
type Scheduler struct {
	executor  AgentExecutor
	mu        sync.RWMutex
	tasks     map[string]*AgentTask
	workflows map[string]*Workflow
}

// NewScheduler 创建调度器
func NewScheduler(executor AgentExecutor) *Scheduler {
	return &Scheduler{
		executor:  executor,
		tasks:     make(map[string]*AgentTask),
		workflows: make(map[string]*Workflow),
	}
}

// ExecuteWorkflow 执行工作流
func (s *Scheduler) ExecuteWorkflow(ctx context.Context, workflow *Workflow) (*WorkflowResult, error) {
	startTime := time.Now()

	// 注册工作流
	s.mu.Lock()
	s.workflows[workflow.ID] = workflow
	s.mu.Unlock()

	// 按顺序执行任务
	for _, task := range workflow.Tasks {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// 等待依赖任务完成
		if err := s.waitForDependencies(ctx, task); err != nil {
			return nil, err
		}

		// 执行任务
		if err := s.executeTask(ctx, task); err != nil {
			return &WorkflowResult{
				Success:   false,
				Tasks:     workflow.Tasks,
				TotalTime: time.Since(startTime),
				Metadata: map[string]interface{}{
					"error": err.Error(),
				},
			}, err
		}
	}

	// 构建结果
	finalContent := s.getFinalContent(workflow.Tasks)

	return &WorkflowResult{
		Success:      true,
		FinalContent: finalContent,
		Tasks:        workflow.Tasks,
		TotalTime:    time.Since(startTime),
		Metadata: map[string]interface{}{
			"workflow_id": workflow.ID,
			"task_count":  len(workflow.Tasks),
		},
	}, nil
}

// executeTask 执行单个任务
func (s *Scheduler) executeTask(ctx context.Context, task *AgentTask) error {
	task.Status = "running"
	task.StartTime = time.Now()

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()

	// 执行 Agent
	result, err := s.executor.Execute(ctx, task.AgentID, task.Input, task.Context)

	task.EndTime = time.Now()

	if err != nil {
		task.Status = "failed"
		task.Error = err
		return fmt.Errorf("task %s failed: %w", task.ID, err)
	}

	task.Status = "completed"
	task.Result = result

	return nil
}

// waitForDependencies 等待依赖任务完成
func (s *Scheduler) waitForDependencies(ctx context.Context, task *AgentTask) error {
	if len(task.DependsOn) == 0 {
		return nil
	}

	// 检查所有依赖
	for _, depID := range task.DependsOn {
		s.mu.RLock()
		depTask, ok := s.tasks[depID]
		s.mu.RUnlock()

		if !ok {
			return fmt.Errorf("dependency task %s not found", depID)
		}

		// 等待依赖任务完成
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			s.mu.RLock()
			status := depTask.Status
			s.mu.RUnlock()

			if status == "completed" {
				break
			}

			if status == "failed" {
				return fmt.Errorf("dependency task %s failed", depID)
			}

			time.Sleep(100 * time.Millisecond)
		}

		// 将依赖任务的结果添加到上下文
		if task.Context == nil {
			task.Context = make(map[string]interface{})
		}
		task.Context[fmt.Sprintf("dependency_%s", depID)] = depTask.Result
	}

	return nil
}

// getFinalContent 获取最终内容
func (s *Scheduler) getFinalContent(tasks []*AgentTask) string {
	if len(tasks) == 0 {
		return ""
	}

	// 返回最后一个成功任务的结果
	for i := len(tasks) - 1; i >= 0; i-- {
		if tasks[i].Status == "completed" {
			return tasks[i].Result
		}
	}

	return ""
}

// GetTaskStatus 获取任务状态
func (s *Scheduler) GetTaskStatus(taskID string) (*AgentTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	return task, nil
}

// GetWorkflowStatus 获取工作流状态
func (s *Scheduler) GetWorkflowStatus(workflowID string) (*Workflow, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workflow, ok := s.workflows[workflowID]
	if !ok {
		return nil, fmt.Errorf("workflow %s not found", workflowID)
	}

	return workflow, nil
}

// Workflow 工作流定义
type Workflow struct {
	ID          string
	Name        string
	Description string
	Tasks       []*AgentTask
	CreatedAt   time.Time
}

// NewWorkflow 创建工作流
func NewWorkflow(id, name, description string) *Workflow {
	return &Workflow{
		ID:          id,
		Name:        name,
		Description: description,
		Tasks:       make([]*AgentTask, 0),
		CreatedAt:   time.Now(),
	}
}

// AddTask 添加任务
func (w *Workflow) AddTask(task *AgentTask) {
	w.Tasks = append(w.Tasks, task)
}
