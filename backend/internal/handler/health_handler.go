package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type HealthHandler struct {
	db          *sql.DB
	neo4jDriver neo4j.DriverWithContext
}

func NewHealthHandler(db *sql.DB, neo4jDriver neo4j.DriverWithContext) *HealthHandler {
	return &HealthHandler{
		db:          db,
		neo4jDriver: neo4jDriver,
	}
}

// HealthCheck 健康检查接口
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  make(map[string]string),
	}

	// 检查 PostgreSQL
	if err := h.db.PingContext(ctx); err != nil {
		health["services"].(map[string]string)["postgres"] = "unhealthy: " + err.Error()
		health["status"] = "unhealthy"
	} else {
		health["services"].(map[string]string)["postgres"] = "healthy"
	}

	// 检查 Neo4j
	if err := h.neo4jDriver.VerifyConnectivity(ctx); err != nil {
		health["services"].(map[string]string)["neo4j"] = "unhealthy: " + err.Error()
		health["status"] = "unhealthy"
	} else {
		health["services"].(map[string]string)["neo4j"] = "healthy"
	}

	statusCode := http.StatusOK
	if health["status"] == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// ReadinessCheck 就绪检查
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ready": true,
	})
}

// LivenessCheck 存活检查
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive": true,
	})
}
