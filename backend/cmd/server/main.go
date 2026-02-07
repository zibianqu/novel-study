package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	defer neo4jDriver.Close()
	log.Println("✅ Neo4j 连接成功")

	// 初始化 Repository
	projectRepo := repository.NewProjectRepository(db)
	chapterRepo := repository.NewChapterRepository(db)

	// 初始化 Service
	projectService := service.NewProjectService(projectRepo)
	chapterService := service.NewChapterService(chapterRepo, projectRepo)

	// 初始化 Handler
	authHandler := handler.NewAuthHandler(db, cfg)
	projectHandler := handler.NewProjectHandler(projectService)
	chapterHandler := handler.NewChapterHandler(chapterService)

	// 初始化 Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// CORS 中间件
	router.Use(middleware.CORS())

	// 静态文件服务
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")

	// API 路由组
	api := router.Group("/api/v1")
	{
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
			// 用户信息
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetInt("user_id")
				username := c.GetString("username")
				c.JSON(200, gin.H{
					"user_id":  userID,
					"username": username,
					"message":  "认证成功",
				})
			})

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
			protected.POST("/chapters/:id/lock", chapterHandler.LockChapter)
			protected.POST("/chapters/:id/unlock", chapterHandler.UnlockChapter)
		}
	}

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 服务器启动在端口 %s", port)
	log.Println("📚 API 文档: http://localhost:" + port + "/api/v1")
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
