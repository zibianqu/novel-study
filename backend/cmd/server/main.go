package main

import (
	"context"
	"fmt"
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
		log.Println("⚠️  警告: 未找到 .env 文件，使用系统环境变量")
	}

	// 加载配置
	cfg := config.Load()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❗ 配置验证失败: %v", err)
	}

	// 初始化数据库连接（使用优化版本）
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("❗ 数据库连接失败: %v", err)
	}
	defer db.Close()
	log.Println("✅ PostgreSQL 连接成功（连接池已优化）")

	// 输出连接池统计
	stats := repository.GetDBStats(db)
	log.Printf("📊 数据库连接池: MaxOpen=%d, MaxIdle=%d", stats.MaxOpenConnections, cfg.DBMaxIdleConnections)

	// 初始化 Neo4j 连接
	neo4jDriver, err := repository.NewNeo4jDriver(cfg)
	if err != nil {
		log.Fatalf("❗ Neo4j 连接失败: %v", err)
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
	
	// 使用 gin.New() 而不是 Default()，手动添加中间件
	router := gin.New()

	// ===== 全局中间件 =====
	router.Use(middleware.ErrorHandler())      // 错误处理
	router.Use(middleware.RequestLogger())     // 请求日志
	router.Use(gin.Recovery())                 // Panic恢复
	router.Use(middleware.CORS())              // CORS
	router.Use(middleware.TimeoutByPath())     // 超时控制
	router.Use(middleware.RateLimitByPath())   // 限流

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
			// 创建输入验证器
			validator := middleware.NewInputValidator()
			
			// 创建登录限流器
			loginLimiter := middleware.NewLoginLimiter()

			// 注册（带验证）
			auth.POST("/register", 
				validator.ValidateRegisterInput(),
				authHandler.Register,
			)

			// 登录（带验证和限流）
			auth.POST("/login",
				validator.ValidateLoginInput(),
				loginLimiter.CheckLimit(),
				authHandler.Login,
			)

			// 刷新Token
			auth.POST("/refresh", 
				middleware.JWTAuth(cfg.JWTSecret),
				authHandler.RefreshToken,
			)
		}

		// 需要认证的接口
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		protected.Use(middleware.APILogger()) // API详细日志
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
	
	printBanner(port, cfg)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❗ 服务器启动失败: %v", err)
	}
}

// printBanner 打印启动信息
func printBanner(port string, cfg *config.Config) {
	fmt.Println("")
	fmt.Println("✨ ========================================")
	fmt.Println("🚀 NovelForge AI 服务器启动成功")
	fmt.Println("✨ ========================================")
	fmt.Println("")
	fmt.Println("🎬  7 个核心 Agent 已就绪")
	fmt.Println("🧠  RAG 知识库系统已启用")
	fmt.Println("🕸️  Neo4j 知识图谱已连接")
	fmt.Println("✅  安全增强: 密码验证 + 登录限流")
	fmt.Println("✅  性能优化: 数据库索引 + 连接池")
	fmt.Println("✅  中间件: CORS / 超时 / 限流 / 日志")
	fmt.Println("")
	fmt.Printf("🔗 前端: http://localhost:%s\n", port)
	fmt.Printf("📚 API: http://localhost:%s/api/v1\n", port)
	fmt.Printf("❤️  Health: http://localhost:%s/api/v1/health\n", port)
	fmt.Println("")
	fmt.Printf("🌐 环境: %s\n", cfg.Environment)
	fmt.Println("✨ ========================================")
	fmt.Println("")
}
