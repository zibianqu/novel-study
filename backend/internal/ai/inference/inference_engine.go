package inference

import (
	"context"
	"fmt"
	"time"
)

// ChapterContext 章节上下文
type ChapterContext struct {
	ChapterNumber int
	Title         string
	Outline       string
	KeyEvents     []string
	Characters    []string
	PlotPoints    []string
}

// InferenceResult 推演结果
type InferenceResult struct {
	ChapterNumber int
	Predictions   []*Prediction
	Conflicts     []*ConflictPrediction
	Suggestions   []string
	Confidence    float64
}

// Prediction 预测
type Prediction struct {
	Type        string  // "plot", "character", "conflict", "theme"
	Content     string
	Probability float64
	Reason      string
}

// ConflictPrediction 冲突预测
type ConflictPrediction struct {
	Type        string   // "logic", "character", "timeline"
	Description string
	AffectedChapters []int
	Severity    string   // "low", "medium", "high"
	Suggestion  string
}

// InferenceEngine 推演引擎
type InferenceEngine struct {
	analyzer *StorylineAnalyzer
}

// NewInferenceEngine 创建推演引擎
func NewInferenceEngine() *InferenceEngine {
	return &InferenceEngine{
		analyzer: NewStorylineAnalyzer(),
	}
}

// InferNextChapters 推演后续章节
func (ie *InferenceEngine) InferNextChapters(
	ctx context.Context,
	currentChapters []*ChapterContext,
	count int,
) ([]*InferenceResult, error) {
	if count <= 0 || count > 10 {
		return nil, fmt.Errorf("count must be between 1 and 10")
	}

	results := make([]*InferenceResult, 0, count)

	// 分析已有章节
	analysis := ie.analyzer.Analyze(currentChapters)

	// 推演每一章
	for i := 0; i < count; i++ {
		chapterNumber := len(currentChapters) + i + 1

		result := &InferenceResult{
			ChapterNumber: chapterNumber,
			Predictions:   make([]*Prediction, 0),
			Conflicts:     make([]*ConflictPrediction, 0),
			Suggestions:   make([]string, 0),
		}

		// 1. 剧情预测
		plotPredictions := ie.predictPlot(analysis, chapterNumber)
		result.Predictions = append(result.Predictions, plotPredictions...)

		// 2. 角色发展预测
		characterPredictions := ie.predictCharacterDevelopment(analysis, chapterNumber)
		result.Predictions = append(result.Predictions, characterPredictions...)

		// 3. 冲突预测
		conflicts := ie.predictConflicts(analysis, chapterNumber)
		result.Conflicts = append(result.Conflicts, conflicts...)

		// 4. 生成建议
		suggestions := ie.generateSuggestions(analysis, chapterNumber)
		result.Suggestions = append(result.Suggestions, suggestions...)

		// 5. 计算置信度
		result.Confidence = ie.calculateConfidence(result)

		results = append(results, result)
	}

	return results, nil
}

// predictPlot 预测剧情
func (ie *InferenceEngine) predictPlot(
	analysis *StorylineAnalysis,
	chapterNumber int,
) []*Prediction {
	predictions := make([]*Prediction, 0)

	// 根据当前剧情线预测
	if analysis.PlotlineProgression < 0.3 {
		// 故事前期
		predictions = append(predictions, &Prediction{
			Type:        "plot",
			Content:     "角色背景展开，世界观建立",
			Probability: 0.8,
			Reason:      "故事处于前期，需要铺垫",
		})
	} else if analysis.PlotlineProgression < 0.7 {
		// 故事中期
		predictions = append(predictions, &Prediction{
			Type:        "plot",
			Content:     "主要冲突升级，转折点出现",
			Probability: 0.85,
			Reason:      "故事进入中期，冲突应该加剧",
		})
	} else {
		// 故事后期
		predictions = append(predictions, &Prediction{
			Type:        "plot",
			Content:     "冲突解决，主线收尾",
			Probability: 0.9,
			Reason:      "故事接近结尾，需要收束",
		})
	}

	return predictions
}

// predictCharacterDevelopment 预测角色发展
func (ie *InferenceEngine) predictCharacterDevelopment(
	analysis *StorylineAnalysis,
	chapterNumber int,
) []*Prediction {
	predictions := make([]*Prediction, 0)

	// 主角发展预测
	for _, char := range analysis.MainCharacters {
		predictions = append(predictions, &Prediction{
			Type:        "character",
			Content:     fmt.Sprintf("%s 将面临重要选择或成长机会", char),
			Probability: 0.75,
			Reason:      "角色弧光发展规律",
		})
	}

	return predictions
}

// predictConflicts 预测潜在冲突
func (ie *InferenceEngine) predictConflicts(
	analysis *StorylineAnalysis,
	chapterNumber int,
) []*ConflictPrediction {
	conflicts := make([]*ConflictPrediction, 0)

	// 逻辑冲突检测
	if len(analysis.UnresolvedPlots) > 5 {
		conflicts = append(conflicts, &ConflictPrediction{
			Type:        "logic",
			Description: "过多未解决的情节线可能导致逻辑混乱",
			AffectedChapters: []int{chapterNumber, chapterNumber + 1},
			Severity:    "medium",
			Suggestion:  "建议在接下来的章节中解决部分情节线",
		})
	}

	// 角色冲突
	if analysis.CharacterCount > 10 {
		conflicts = append(conflicts, &ConflictPrediction{
			Type:        "character",
			Description: "角色过多可能分散读者注意力",
			AffectedChapters: []int{chapterNumber},
			Severity:    "low",
			Suggestion:  "聚焦核心角色，其他角色适当淡化",
		})
	}

	return conflicts
}

// generateSuggestions 生成建议
func (ie *InferenceEngine) generateSuggestions(
	analysis *StorylineAnalysis,
	chapterNumber int,
) []string {
	suggestions := make([]string, 0)

	// 基于剧情进度给建议
	if analysis.PlotlineProgression < 0.3 {
		suggestions = append(suggestions,
			"建议增加环境和世界观描写",
			"可以引入支线角色丰富故事",
		)
	} else if analysis.PlotlineProgression < 0.7 {
		suggestions = append(suggestions,
			"建议设置重要转折点",
			"可以深化主要矛盾",
		)
	} else {
		suggestions = append(suggestions,
			"建议开始收束情节线",
			"准备高潮和结局",
		)
	}

	// 基于节奏给建议
	if analysis.Pace == "too_fast" {
		suggestions = append(suggestions, "节奏偏快，建议增加细节描写")
	} else if analysis.Pace == "too_slow" {
		suggestions = append(suggestions, "节奏偏慢，建议加快剧情推进")
	}

	return suggestions
}

// calculateConfidence 计算置信度
func (ie *InferenceEngine) calculateConfidence(result *InferenceResult) float64 {
	// 基于预测数量和冲突数计算
	confidence := 0.7 // 基础置信度

	if len(result.Predictions) > 3 {
		confidence += 0.1
	}

	if len(result.Conflicts) > 2 {
		confidence -= 0.1 // 冲突多降低置信度
	}

	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// StorylineAnalyzer 故事线分析器
type StorylineAnalyzer struct{}

func NewStorylineAnalyzer() *StorylineAnalyzer {
	return &StorylineAnalyzer{}
}

// StorylineAnalysis 故事线分析结果
type StorylineAnalysis struct {
	PlotlineProgression float64  // 剧情线进度 0-1
	MainCharacters      []string
	CharacterCount      int
	UnresolvedPlots     []string
	Pace                string // "too_fast", "normal", "too_slow"
}

func (sa *StorylineAnalyzer) Analyze(chapters []*ChapterContext) *StorylineAnalysis {
	// 简化分析
	analysis := &StorylineAnalysis{
		PlotlineProgression: float64(len(chapters)) / 30.0, // 假设30章
		MainCharacters:      []string{"主角"},
		CharacterCount:      5,
		UnresolvedPlots:     []string{},
		Pace:                "normal",
	}

	return analysis
}

// InferenceReport 推演报告
type InferenceReport struct {
	GeneratedAt time.Time
	ProjectID   int
	Results     []*InferenceResult
	Summary     string
	Warnings    []string
}

// GenerateReport 生成推演报告
func (ie *InferenceEngine) GenerateReport(
	projectID int,
	results []*InferenceResult,
) *InferenceReport {
	report := &InferenceReport{
		GeneratedAt: time.Now(),
		ProjectID:   projectID,
		Results:     results,
		Warnings:    make([]string, 0),
	}

	// 生成摘要
	report.Summary = fmt.Sprintf(
		"推演了 %d 章内容，发现 %d 个潜在冲突",
		len(results),
		ie.countTotalConflicts(results),
	)

	// 收集警告
	for _, result := range results {
		for _, conflict := range result.Conflicts {
			if conflict.Severity == "high" {
				report.Warnings = append(report.Warnings,
					fmt.Sprintf("第%d章: %s", result.ChapterNumber, conflict.Description))
			}
		}
	}

	return report
}

func (ie *InferenceEngine) countTotalConflicts(results []*InferenceResult) int {
	count := 0
	for _, result := range results {
		count += len(result.Conflicts)
	}
	return count
}
