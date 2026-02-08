package routes

import (
	"github.com/gorilla/mux"
	"github.com/zibianqu/novel-study/api/handlers"
)

// RegisterGraphRoutes 注册知识图谱路由
func RegisterGraphRoutes(router *mux.Router, handler *handlers.GraphHandler) {
	graph := router.PathPrefix("/api/graph").Subrouter()

	// 应用中间件W
	// graph.Use(middleware.Auth)
	// graph.Use(middleware.RateLimit)
	// graph.Use(middleware.Logging)

	// 图谱操作
	graph.HandleFunc("/create", handler.CreateGraph).Methods("POST")
	graph.HandleFunc("/query", handler.QueryGraph).Methods("POST")
	graph.HandleFunc("/path", handler.FindPath).Methods("POST")

	// 统计和分析
	graph.HandleFunc("/stats", handler.GetStatistics).Methods("GET")
	graph.HandleFunc("/validate", handler.ValidateConsistency).Methods("GET")

	// 人物分析
	graph.HandleFunc("/character/analyze", handler.AnalyzeCharacter).Methods("GET")
	graph.HandleFunc("/character/timeline", handler.GetCharacterTimeline).Methods("GET")

	// 写作辅助
	graph.HandleFunc("/plot-holes", handler.DetectPlotHoles).Methods("GET")
	graph.HandleFunc("/suggestions", handler.GenerateSuggestions).Methods("GET")

	// 搜索
	graph.HandleFunc("/search", handler.SearchGraph).Methods("GET")

	// 健康检查
	graph.HandleFunc("/health", handler.HealthCheck).Methods("GET")
}
