package handlers

import (
	"encoding/json"
	"net/http"
	"novel-study/backend/internal/ai/inference"
	"strconv"
)

// InferenceHandler 推演 API 处理器
type InferenceHandler struct {
	service *inference.InferenceService
}

// NewInferenceHandler 创建处理器
func NewInferenceHandler() *InferenceHandler {
	return &InferenceHandler{
		service: inference.NewInferenceService(),
	}
}

// InferNextChapters 推演后续章节
// POST /api/inference/chapters
func (h *InferenceHandler) InferNextChapters(w http.ResponseWriter, r *http.Request) {
	var req inference.InferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 执行推演
	resp, err := h.service.Infer(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetInferenceReport 获取推演报告
// GET /api/inference/report/:project_id
func (h *InferenceHandler) GetInferenceReport(w http.ResponseWriter, r *http.Request) {
	// 从路径中获取 project_id
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}

	// 这里应该从数据库获取之前的推演结果
	// 暂时返回空报告
	report := &inference.InferenceReport{
		ProjectID: projectID,
		Summary:   "暂无推演报告",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// AnalyzeChapters 分析章节
// POST /api/inference/analyze
func (h *InferenceHandler) AnalyzeChapters(w http.ResponseWriter, r *http.Request) {
	var chapters []*inference.ChapterContext
	if err := json.NewDecoder(r.Body).Decode(&chapters); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	analysis, err := h.service.AnalyzeChapters(r.Context(), chapters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// DetectConflicts 检测冲突
// POST /api/inference/conflicts
func (h *InferenceHandler) DetectConflicts(w http.ResponseWriter, r *http.Request) {
	var chapters []*inference.ChapterContext
	if err := json.NewDecoder(r.Body).Decode(&chapters); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conflicts, err := h.service.DetectConflicts(r.Context(), chapters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"conflicts": conflicts,
		"count":     len(conflicts),
	})
}

// GenerateSuggestions 生成建议
// POST /api/inference/suggestions
func (h *InferenceHandler) GenerateSuggestions(w http.ResponseWriter, r *http.Request) {
	var chapters []*inference.ChapterContext
	if err := json.NewDecoder(r.Body).Decode(&chapters); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	suggestions, err := h.service.GenerateSuggestions(r.Context(), chapters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"suggestions": suggestions,
		"count":       len(suggestions),
	})
}

// PredictNextChapter 预测下一章
// POST /api/inference/predict
func (h *InferenceHandler) PredictNextChapter(w http.ResponseWriter, r *http.Request) {
	var chapters []*inference.ChapterContext
	if err := json.NewDecoder(r.Body).Decode(&chapters); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prediction, err := h.service.PredictNextChapter(r.Context(), chapters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prediction)
}
