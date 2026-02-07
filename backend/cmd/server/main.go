package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zibianqu/novel-study/internal/ai"
	"github.com/zibianqu/novel-study/internal/ai/rag"
	"github.com/zibianqu/novel-study/internal/config"
	"github.com/zibianqu/novel-study/internal/handler"
	"github.com/zibianqu/novel-study/internal/middleware"
	"github.com/zibianqu/novel-study/internal/repository"
	"github.com/zibianqu/novel-study/internal/service"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("警告: 未找到 .env 文件，使用系统环境变量")
	}

	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL 连接成功")

	// 初始化 Neo4j 连接
	neo4jDriver, err := repository.NewNeo4jDriver(cfg)
	if err != nil {
		log.Fatalf("Neo4j 连接失败: %v", err)
	}
	defer neo4jDriver.Close(context.Background())
	log.Println("✅ Neo4j 连接成功")

	// 初始化 AI 引擎
	aiEngine := ai.NewEngine(cfg)
	log.Printf("✅ AI 引擎初始化完成，已注册 %d 个 Agent", len(aiEngine.ListAgents()))

	// 初始化 RAG 系统
	embeddingService := rag.NewEmbeddingService(cfg.OpenAIAPIKey)
	vectorStore := rag.NewVectorStore(db)
	retriever := rag.NewRetriever(embeddingService, vectorStore)
	log.Println("✅ RAG 系统初始化完成")

	// 初始化 Repository
	projectRepo := repository.NewProjectRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	neo4jRepo := repository.NewNeo4jRepository(neo4jDriver)

	// 初始化 Service
	projectService := service.NewProjectService(projectRepo)
	chapterService := service.NewChapterService(chapterRepo, projectRepo)
	aiService := service.NewAIService(aiEngine, agentRepo, projectRepo)
	knowledgeService := service.NewKnowledgeService(knowledgeRepo, projectRepo, retriever)
	graphService := service.NewGraphService(neo4jRepo, projectRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(db, cfg)
	projectHandler := handler.NewProjectHandler(projectService)
	chapterHandler := handler.NewChapterHandler(chapterService)
	aiHandler := handler.NewAIHandler(aiService)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeService)
	graphHandler := handler.NewGraphHandler(graphService)
	storylineHandler := handler.NewStorylineHandler(db)
	healthHandler := handler.NewHealthHandler(db, neo4jDriver)

	// 初始化 Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS 中间件
	router.Use(middleware.CORS())

	// 静态文件服务
	router.Static("/css", "./frontend/css")
	router.Static("/js", "./frontend/js")
	router.StaticFile("/", "./frontend/index.html")
	router.StaticFile("/index.html", "./frontend/index.html")
	router.StaticFile("/dashboard.html", "./frontend/dashboard.html")
	router.StaticFile("/project.html", "./frontend/project.html")
	router.StaticFile("/editor.html", "./frontend/editor.html")
	router.StaticFile("/knowledge.html", "./frontend/knowledge.html")
	router.StaticFile("/graph.html", "./frontend/graph.html")
	router.StaticFile("/storyline.html", "./frontend/storyline.html")

	// API 路由组
	api := router.Group("/api/v1")
	{
		// 健康检查接口（公开）
		api.GET("/health", healthHandler.HealthCheck)
		api.GET("/ready", healthHandler.ReadinessCheck)
		api.GET("/alive", healthHandler.LivenessCheck)

		// 公开接口
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			// 项目管理
			protected.GET("/projects", projectHandler.GetProjects)
			protected.POST("/projects", projectHandler.CreateProject)
			protected.GET("/projects/:id", projectHandler.GetProject)
			protected.PUT("/projects/:id", projectHandler.UpdateProject)
			protected.DELETE("/projects/:id", projectHandler.DeleteProject)

			// 章节管理
			protected.GET("/chapters/project/:projectId", chapterHandler.GetProjectChapters)
			protected.POST("/chapters", chapterHandler.CreateChapter)
			protected.GET("/chapters/:id", chapterHandler.GetChapter)
			protected.PUT("/chapters/:id", chapterHandler.UpdateChapter)
			protected.DELETE("/chapters/:id", chapterHandler.DeleteChapter)

			// AI 功能
			protected.GET("/ai/agents", aiHandler.GetAgents)
			protected.POST("/ai/chat", aiHandler.Chat)
			protected.POST("/ai/chat/stream", middleware.SSE(), aiHandler.ChatStream)
			protected.POST("/ai/generate/chapter", aiHandler.GenerateChapter)
			protected.POST("/ai/check/quality", aiHandler.CheckQuality)

			// 知识库
			protected.GET("/knowledge/project/:projectId", knowledgeHandler.GetProjectKnowledge)
			protected.POST("/knowledge", knowledgeHandler.CreateKnowledge)
			protected.GET("/knowledge/:id", knowledgeHandler.GetKnowledge)
			protected.DELETE("/knowledge/:id", knowledgeHandler.DeleteKnowledge)
			protected.POST("/knowledge/search", knowledgeHandler.SearchKnowledge)

			// 知识图谱
			protected.GET("/graph/project/:projectId", graphHandler.GetProjectGraph)
			protected.POST("/graph/node", graphHandler.CreateNode)
			protected.POST("/graph/relation", graphHandler.CreateRelation)

			// 三线管理
			protected.GET("/storylines/project/:projectId", storylineHandler.GetProjectStorylines)
			protected.POST("/storylines", storylineHandler.CreateStoryline)
		}
	}

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("")
	log.Println("✨ ========================================")
	log.Printf("🚀 NovelForge AI 服务器启动成功")
	log.Printf("🎬 7 个核心 Agent 已就绪")
	log.Printf("🧠 RAG 知识库系统已启用")
	log.Printf("🕸️ Neo4j 知识图谱已连接")
	log.Printf("🔗 前端: http://localhost:%s", port)
	log.Printf("📚 API: http://localhost:%s/api/v1", port)
	log.Printf("❤️ Health: http://localhost:%s/api/v1/health", port)
	log.Println("✨ ========================================")
	log.Println("")

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
