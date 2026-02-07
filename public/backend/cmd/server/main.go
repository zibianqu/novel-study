package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/novelforge/backend/internal/config"
	"github.com/novelforge/backend/internal/database"
	"github.com/novelforge/backend/internal/router"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化 PostgreSQL
	pgPool, err := database.InitPostgres(cfg.Database.Postgres)
	if err != nil {
		log.Fatalf("连接 PostgreSQL 失败: %v", err)
	}
	defer pgPool.Close()
	log.Println("✅ PostgreSQL 连接成功")

	// 3. 初始化 Neo4j
	neo4jDriver, err := database.InitNeo4j(cfg.Database.Neo4j)
	if err != nil {
		log.Fatalf("连接 Neo4j 失败: %v", err)
	}
	defer neo4jDriver.Close(context.Background())
	log.Println("✅ Neo4j 连接成功")

	// 4. 运行数据库迁移
	if err := database.RunMigrations(pgPool); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 5. 初始化路由
	r := router.Setup(cfg, pgPool, neo4jDriver)

	// 6. 启动 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second, // SSE 需要较长的写超时
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("🚀 NovelForge AI 服务启动在 http://localhost:%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 7. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭失败: %v", err)
	}
	log.Println("✅ 服务已安全关闭")
}
