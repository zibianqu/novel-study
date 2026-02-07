package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/novelforge/backend/internal/config"
	"github.com/novelforge/backend/internal/handler"
	"github.com/novelforge/backend/internal/middleware"
	"github.com/novelforge/backend/internal/repository"
	"github.com/novelforge/backend/internal/service"
)

// Setup 初始化并返回 Gin 路由引擎
func Setup(cfg *config.Config, db *pgxpool.Pool, neo4jDriver neo4j.DriverWithContext) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware(cfg.Server.CORS.AllowedOrigins))

	// 静态文件服务（前端）
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	// ========== 初始化各层 ==========

	// Repository 层
	userRepo := repository.NewUserRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	knowledgeRepo := repository.NewKnowledgeRepository(db)
	workflowRepo := repository.NewWorkflowRepository(db)
	storylineRepo := repository.NewStorylineRepository(db)

	// Service 层
	authService := service.NewAuthService(userRepo, cfg.JWT)
	projectService := service.NewProjectService(projectRepo)
	chapterService := service.NewChapterService(chapterRepo)
	agentService := service.NewAgentService(agentRepo)
	knowledgeService := service.NewKnowledgeService(knowledgeRepo)
	workflowService := service.NewWorkflowService(workflowRepo)
	storylineService := service.NewStorylineService(storylineRepo)
	_ = neo4jDriver // 后续阶段使用

	// Handler 层
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	chapterHandler := handler.NewChapterHandler(chapterService)
	agentHandler := handler.NewAgentHandler(agentService)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeService)
	workflowHandler := handler.NewWorkflowHandler(workflowService)
	storylineHandler := handler.NewStorylineHandler(storylineService)

	// ========== API 路由 ==========
	api := r.Group("/api/v1")
	{
		// 认证（无需Token）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 以下接口需要认证
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// 项目
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.List)
				projects.POST("", projectHandler.Create)
				projects.GET("/:id", projectHandler.Get)
				projects.PUT("/:id", projectHandler.Update)
				projects.DELETE("/:id", projectHandler.Delete)
				projects.POST("/:id/export", projectHandler.Export)
			}

			// 章节
			chapters := protected.Group("/chapters")
			{
				chapters.GET("/project/:projectId", chapterHandler.ListByProject)
				chapters.POST("", chapterHandler.Create)
				chapters.GET("/:id", chapterHandler.Get)
				chapters.PUT("/:id", chapterHandler.Update)
				chapters.GET("/:id/versions", chapterHandler.ListVersions)
				chapters.POST("/:id/rollback", chapterHandler.Rollback)
				chapters.POST("/:id/lock", chapterHandler.Lock)
				chapters.POST("/:id/unlock", chapterHandler.Unlock)
			}

			// Agent 管理
			agents := protected.Group("/agents")
			{
				agents.GET("", agentHandler.List)
				agents.POST("", agentHandler.Create)
				agents.GET("/:id", agentHandler.Get)
				agents.PUT("/:id", agentHandler.Update)
				agents.DELETE("/:id", agentHandler.Delete)
				agents.POST("/:id/test", agentHandler.Test)
			}

			// 知识库
			knowledge := protected.Group("/knowledge")
			{
				knowledge.GET("/categories", knowledgeHandler.ListCategories)
				knowledge.POST("/categories", knowledgeHandler.CreateCategory)
				knowledge.PUT("/categories/:id", knowledgeHandler.UpdateCategory)
				knowledge.DELETE("/categories/:id", knowledgeHandler.DeleteCategory)
				knowledge.GET("/items", knowledgeHandler.ListItems)
				knowledge.POST("/items", knowledgeHandler.CreateItem)
				knowledge.GET("/items/:id", knowledgeHandler.GetItem)
				knowledge.PUT("/items/:id", knowledgeHandler.UpdateItem)
				knowledge.DELETE("/items/:id", knowledgeHandler.DeleteItem)
				knowledge.POST("/import", knowledgeHandler.Import)
				knowledge.POST("/search", knowledgeHandler.Search)
			}

			// 工作流
			workflows := protected.Group("/workflows")
			{
				workflows.GET("", workflowHandler.List)
				workflows.POST("", workflowHandler.Create)
				workflows.GET("/:id", workflowHandler.Get)
				workflows.PUT("/:id", workflowHandler.Update)
				workflows.DELETE("/:id", workflowHandler.Delete)
				workflows.POST("/:id/execute", workflowHandler.Execute)
				workflows.GET("/executions/:id", workflowHandler.GetExecution)
			}

			// 三线管理
			storylines := protected.Group("/storylines")
			{
				storylines.GET("/project/:projectId", storylineHandler.ListByProject)
				storylines.POST("", storylineHandler.Create)
				storylines.PUT("/:id", storylineHandler.Update)
				storylines.DELETE("/:id", storylineHandler.Delete)
				storylines.POST("/adjust", storylineHandler.Adjust)
			}

			// AI 创作（后续阶段实现）
			ai := protected.Group("/ai")
			{
				ai.POST("/chat", handler.PlaceholderHandler("AI对话"))
				ai.POST("/forecast", handler.PlaceholderHandler("多章推演"))
				ai.POST("/continue", handler.PlaceholderHandler("续写"))
				ai.POST("/polish", handler.PlaceholderHandler("润色"))
				ai.POST("/rewrite", handler.PlaceholderHandler("改写"))
				ai.POST("/dialogue", handler.PlaceholderHandler("对话生成"))
				ai.POST("/consistency-check", handler.PlaceholderHandler("一致性检查"))
				ai.POST("/character/generate", handler.PlaceholderHandler("角色生成"))
				ai.POST("/outline/generate", handler.PlaceholderHandler("大纲生成"))
			}
		}
	}

	return r
}
