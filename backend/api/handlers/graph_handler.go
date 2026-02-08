package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zibianqu/novel-study/internal/graph"
)

// GraphHandler 知识图谱 API 处理器
type GraphHandler struct {
	service *graph.GraphService
}

// NewGraphHandler 创建处理器
func NewGraphHandler(service *graph.GraphService) *GraphHandler {
	return &GraphHandler{
		service: service,
	}
}

// CreateGraph 创建知识图谱
// POST /api/graph/create
func (h *GraphHandler) CreateGraph(w http.ResponseWriter, r *http.Request) {
	var req graph.CreateGraphRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 设置默认值
	if req.MinConfidence == 0 {
		req.MinConfidence = 0.5
	}
	if req.MaxNodes == 0 {
		req.MaxNodes = 1000
	}

	resp, err := h.service.CreateKnowledgeGraph(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// QueryGraph 查询图谱
// POST /api/graph/query
func (h *GraphHandler) QueryGraph(w http.ResponseWriter, r *http.Request) {
	var query graph.GraphQuery
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.QueryGraph(r.Context(), &query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// FindPath 查找路径
// POST /api/graph/path
func (h *GraphHandler) FindPath(w http.ResponseWriter, r *http.Request) {
	var req graph.PathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 设置默认最大深度
	if req.MaxDepth == 0 {
		req.MaxDepth = 5
	}

	resp, err := h.service.FindPath(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetStatistics 获取统计信息
// GET /api/graph/stats
func (h *GraphHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStatistics(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// AnalyzeCharacter 分析人物
// GET /api/graph/character/:id/analyze
func (h *GraphHandler) AnalyzeCharacter(w http.ResponseWriter, r *http.Request) {
	characterID := r.URL.Query().Get("id")
	if characterID == "" {
		http.Error(w, "character_id is required", http.StatusBadRequest)
		return
	}

	analysis, err := h.service.AnalyzeCharacterRelations(r.Context(), characterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// DetectPlotHoles 检测剧情漏洞
// GET /api/graph/plot-holes
func (h *GraphHandler) DetectPlotHoles(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.DetectPlotHoles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// GenerateSuggestions 生成写作建议
// GET /api/graph/suggestions/:project_id
func (h *GraphHandler) GenerateSuggestions(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}

	suggestions, err := h.service.GenerateSuggestions(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

// SearchGraph 搜索图谱
// GET /api/graph/search
func (h *GraphHandler) SearchGraph(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}

	result, err := h.service.SearchGraph(r.Context(), keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ValidateConsistency 验证一致性
// GET /api/graph/validate
func (h *GraphHandler) ValidateConsistency(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ValidateConsistency(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// GetCharacterTimeline 获取人物时间线
// GET /api/graph/character/:id/timeline
func (h *GraphHandler) GetCharacterTimeline(w http.ResponseWriter, r *http.Request) {
	characterID := r.URL.Query().Get("id")
	if characterID == "" {
		http.Error(w, "character_id is required", http.StatusBadRequest)
		return
	}

	timeline, err := h.service.GetCharacterTimeline(r.Context(), characterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(timeline)
}

// HealthCheck 健康检查
// GET /api/graph/health
func (h *GraphHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := h.service.HealthCheck(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "knowledge_graph",
	})
}
