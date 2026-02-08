package inference

import (
	"context"
	"fmt"
)

// InferenceService 推演服务
type InferenceService struct {
	engine *InferenceEngine
}

// NewInferenceService 创建推演服务
func NewInferenceService() *InferenceService {
	return &InferenceService{
		engine: NewInferenceEngine(),
	}
}

// InferRequest 推演请求
type InferRequest struct {
	ProjectID      int
	CurrentChapters []*ChapterContext
	InferCount     int
	Options        *InferOptions
}

// InferOptions 推演选项
type InferOptions struct {
	IncludePlot      bool
	IncludeCharacter bool
	IncludeConflict  bool
	DetailLevel      string // "simple", "detailed"
}

// InferResponse 推演响应
type InferResponse struct {
	Success bool
	Results []*InferenceResult
	Report  *InferenceReport
	Message string
}

// Infer 执行推演
func (is *InferenceService) Infer(
	ctx context.Context,
	req *InferRequest,
) (*InferResponse, error) {
	// 验证请求
	if err := is.validateRequest(req); err != nil {
		return nil, err
	}

	// 执行推演
	results, err := is.engine.InferNextChapters(
		ctx,
		req.CurrentChapters,
		req.InferCount,
	)
	if err != nil {
		return &InferResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// 生成报告
	report := is.engine.GenerateReport(req.ProjectID, results)

	return &InferResponse{
		Success: true,
		Results: results,
		Report:  report,
		Message: "推演成功",
	}, nil
}

// validateRequest 验证请求
func (is *InferenceService) validateRequest(req *InferRequest) error {
	if req.ProjectID <= 0 {
		return fmt.Errorf("invalid project_id")
	}

	if len(req.CurrentChapters) == 0 {
		return fmt.Errorf("no current chapters provided")
	}

	if req.InferCount <= 0 || req.InferCount > 10 {
		return fmt.Errorf("infer_count must be between 1 and 10")
	}

	return nil
}

// GetInferenceReport 获取推演报告
func (is *InferenceService) GetInferenceReport(
	ctx context.Context,
	projectID int,
	results []*InferenceResult,
) (*InferenceReport, error) {
	return is.engine.GenerateReport(projectID, results), nil
}

// AnalyzeChapters 分析章节
func (is *InferenceService) AnalyzeChapters(
	ctx context.Context,
	chapters []*ChapterContext,
) (*StorylineAnalysis, error) {
	analyzer := NewStorylineAnalyzer()
	return analyzer.Analyze(chapters), nil
}

// PredictNextChapter 预测下一章
func (is *InferenceService) PredictNextChapter(
	ctx context.Context,
	currentChapters []*ChapterContext,
) (*InferenceResult, error) {
	results, err := is.engine.InferNextChapters(ctx, currentChapters, 1)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no prediction generated")
	}

	return results[0], nil
}

// DetectConflicts 检测冲突
func (is *InferenceService) DetectConflicts(
	ctx context.Context,
	chapters []*ChapterContext,
) ([]*ConflictPrediction, error) {
	// 分析章节
	analyzer := NewStorylineAnalyzer()
	analysis := analyzer.Analyze(chapters)

	// 预测冲突
	conflicts := is.engine.predictConflicts(analysis, len(chapters)+1)

	return conflicts, nil
}

// GenerateSuggestions 生成建议
func (is *InferenceService) GenerateSuggestions(
	ctx context.Context,
	chapters []*ChapterContext,
) ([]string, error) {
	// 分析章节
	analyzer := NewStorylineAnalyzer()
	analysis := analyzer.Analyze(chapters)

	// 生成建议
	suggestions := is.engine.generateSuggestions(analysis, len(chapters)+1)

	return suggestions, nil
}
