package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zibianqu/novel-study/internal/config"
	_ "github.com/lib/pq"
)

// NewPostgresDB 创建优化的PostgreSQL连接
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	// 构建连接字符串
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// 打开数据库连接
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	// ===== 连接池配置 =====
	
	// 最大打开连接数
	maxOpenConns := cfg.DBMaxConnections
	if maxOpenConns == 0 {
		maxOpenConns = 25 // 默认值
	}
	db.SetMaxOpenConns(maxOpenConns)

	// 最大空闲连接数
	maxIdleConns := cfg.DBMaxIdleConnections
	if maxIdleConns == 0 {
		maxIdleConns = 5 // 默认值
	}
	db.SetMaxIdleConns(maxIdleConns)

	// 连接最大生命周期（5分钟）
	db.SetConnMaxLifetime(5 * time.Minute)

	// 连接最大空闲时间（3分钟）
	db.SetConnMaxIdleTime(3 * time.Minute)

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("数据库ping失败: %w", err)
	}

	// 验证pgvector扩展
	if err := verifyPgVector(db); err != nil {
		return nil, fmt.Errorf("pgvector验证失败: %w", err)
	}

	return db, nil
}

// verifyPgVector 验证pgvector扩展是否启用
func verifyPgVector(db *sql.DB) error {
	var exists bool
	query := `SELECT EXISTS(
		SELECT 1 FROM pg_extension WHERE extname = 'vector'
	)`

	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return fmt.Errorf("检查pgvector扩展失败: %w", err)
	}

	if !exists {
		return fmt.Errorf("pgvector扩展未启用,请执行 migrations/009_enable_pgvector.sql")
	}

	return nil
}

// GetDBStats 获取数据库连接池统计
func GetDBStats(db *sql.DB) sql.DBStats {
	return db.Stats()
}
