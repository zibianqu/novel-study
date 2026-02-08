package collaboration

import (
	"context"
	"fmt"
	"time"
)

// ReviewFeedback 审核反馈
type ReviewFeedback struct {
	Approved    bool
	Score       float64 // 0-100
	Issues      []string
	Suggestions []string
	Comments    string
}

// RevisionResult 修改结果
type RevisionResult struct {
	Content      string
	Changes      []string
	Improved     bool
	RevisionNote string
}

// ReviewLoopConfig 审核循环配置
type ReviewLoopConfig struct {
	MaxIterations   int     // 最大迭代次数
	MinScore        float64 // 最低分数要求
	Timeout         time.Duration
	AutoApprove     bool // 超过迭代次数后自动通过
}

// ReviewLoop 审核-修改循环
type ReviewLoop struct {
	scheduler    *Scheduler
	messageBus   *MessageBus
	config       *ReviewLoopConfig
}

// NewReviewLoop 创建审核循环
func NewReviewLoop(
	scheduler *Scheduler,
	messageBus *MessageBus,
	config *ReviewLoopConfig,
) *ReviewLoop {
	if config == nil {
		config = &ReviewLoopConfig{
			MaxIterations: 3,
			MinScore:      80.0,
			Timeout:       5 * time.Minute,
			AutoApprove:   true,
		}
	}

	return &ReviewLoop{
		scheduler:  scheduler,
		messageBus: messageBus,
		config:     config,
	}
}

// Execute 执行审核-修改循环
func (rl *ReviewLoop) Execute(
	ctx context.Context,
	generatorAgentID int,
	reviewerAgentID int,
	initialContent string,
	context map[string]interface{},
) (*ReviewLoopResult, error) {
	startTime := time.Now()

	// 创建超时上下文
	ctxWithTimeout, cancel := context.WithTimeout(ctx, rl.config.Timeout)
	defer cancel()

	iterations := make([]*ReviewIteration, 0)
	currentContent := initialContent

	for i := 0; i < rl.config.MaxIterations; i++ {
		// 检查上下文
		select {
		case <-ctxWithTimeout.Done():
			return nil, fmt.Errorf("review loop timeout")
		default:
		}

		iterStartTime := time.Now()

		// 1. 审核阶段
		feedback, err := rl.review(ctxWithTimeout, reviewerAgentID, currentContent, context)
		if err != nil {
			return nil, fmt.Errorf("review failed: %w", err)
		}

		// 2. 记录迭代
		iteration := &ReviewIteration{
			Iteration:  i + 1,
			Content:    currentContent,
			Feedback:   feedback,
			StartTime:  iterStartTime,
			EndTime:    time.Now(),
		}

		// 3. 检查是否通过
		if feedback.Approved || feedback.Score >= rl.config.MinScore {
			iteration.Approved = true
			iterations = append(iterations, iteration)

			return &ReviewLoopResult{
				Success:       true,
				FinalContent:  currentContent,
				Iterations:    iterations,
				TotalDuration: time.Since(startTime),
				FinalScore:    feedback.Score,
			}, nil
		}

		// 4. 修改阶段
		revision, err := rl.revise(ctxWithTimeout, generatorAgentID, currentContent, feedback, context)
		if err != nil {
			return nil, fmt.Errorf("revision failed: %w", err)
		}

		iteration.Revision = revision
		iterations = append(iterations, iteration)

		// 5. 更新内容
		if revision.Improved {
			currentContent = revision.Content
		} else {
			// 如果没有改进，停止循环
			break
		}

		// 6. 发送消息
		rl.sendIterationMessage(generatorAgentID, reviewerAgentID, i+1, feedback, revision)
	}

	// 超过迭代次数
	if rl.config.AutoApprove {
		return &ReviewLoopResult{
			Success:       true,
			FinalContent:  currentContent,
			Iterations:    iterations,
			TotalDuration: time.Since(startTime),
			FinalScore:    0,
			AutoApproved:  true,
		}, nil
	}

	return &ReviewLoopResult{
		Success:       false,
		FinalContent:  currentContent,
		Iterations:    iterations,
		TotalDuration: time.Since(startTime),
		FinalScore:    0,
	}, fmt.Errorf("max iterations reached without approval")
}

// review 执行审核
func (rl *ReviewLoop) review(
	ctx context.Context,
	reviewerAgentID int,
	content string,
	context map[string]interface{},
) (*ReviewFeedback, error) {
	// 构建审核任务
	task := &AgentTask{
		ID:      fmt.Sprintf("review_%d", time.Now().Unix()),
		AgentID: reviewerAgentID,
		Type:    "review",
		Input:   content,
		Context: context,
	}

	// 执行审核
	if err := rl.scheduler.executeTask(ctx, task); err != nil {
		return nil, err
	}

	// 解析审核结果 (简化处理，实际应解析 JSON)
	feedback := &ReviewFeedback{
		Approved: false,
		Score:    75.0, // 示例分数
		Issues: []string{
			"部分描写不够生动",
		},
		Suggestions: []string{
			"增加环境细节描写",
		},
		Comments: task.Result,
	}

	return feedback, nil
}

// revise 执行修改
func (rl *ReviewLoop) revise(
	ctx context.Context,
	generatorAgentID int,
	content string,
	feedback *ReviewFeedback,
	context map[string]interface{},
) (*RevisionResult, error) {
	// 构建修改任务
	reviseContext := make(map[string]interface{})
	for k, v := range context {
		reviseContext[k] = v
	}
	reviseContext["original_content"] = content
	reviseContext["feedback"] = feedback

	task := &AgentTask{
		ID:      fmt.Sprintf("revise_%d", time.Now().Unix()),
		AgentID: generatorAgentID,
		Type:    "revise",
		Input:   fmt.Sprintf("请根据以下反馈修改内容\n\n原内容: %s\n\n反馈: %v", content, feedback.Issues),
		Context: reviseContext,
	}

	// 执行修改
	if err := rl.scheduler.executeTask(ctx, task); err != nil {
		return nil, err
	}

	revision := &RevisionResult{
		Content:      task.Result,
		Improved:     true,
		Changes:      []string{"增加了环境细节"},
		RevisionNote: "根据反馈进行了修改",
	}

	return revision, nil
}

// sendIterationMessage 发送迭代消息
func (rl *ReviewLoop) sendIterationMessage(
	generatorID, reviewerID, iteration int,
	feedback *ReviewFeedback,
	revision *RevisionResult,
) {
	// 发送反馈消息
	feedbackMsg := NewMessageBuilder().
		From(reviewerID).
		To(generatorID).
		Type("feedback").
		Content(feedback.Comments).
		Metadata("iteration", iteration).
		Metadata("score", feedback.Score).
		Metadata("approved", feedback.Approved).
		Build()

	rl.messageBus.Publish(feedbackMsg)

	// 发送修改消息
	if revision != nil {
		revisionMsg := NewMessageBuilder().
			From(generatorID).
			To(reviewerID).
			Type("revision").
			Content(revision.RevisionNote).
			Metadata("iteration", iteration).
			Metadata("improved", revision.Improved).
			Build()

		rl.messageBus.Publish(revisionMsg)
	}
}

// ReviewIteration 审核迭代
type ReviewIteration struct {
	Iteration  int
	Content    string
	Feedback   *ReviewFeedback
	Revision   *RevisionResult
	Approved   bool
	StartTime  time.Time
	EndTime    time.Time
}

// ReviewLoopResult 审核循环结果
type ReviewLoopResult struct {
	Success       bool
	FinalContent  string
	Iterations    []*ReviewIteration
	TotalDuration time.Duration
	FinalScore    float64
	AutoApproved  bool
}

// GetIterationCount 获取迭代次数
func (r *ReviewLoopResult) GetIterationCount() int {
	return len(r.Iterations)
}

// GetApprovedIteration 获取通过的迭代
func (r *ReviewLoopResult) GetApprovedIteration() *ReviewIteration {
	for _, iter := range r.Iterations {
		if iter.Approved {
			return iter
		}
	}
	return nil
}
