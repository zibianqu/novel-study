.PHONY: help dev build test clean docker-up docker-down migrate

# 默认目标
help:
	@echo "NovelForge AI - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make dev         - Start development server"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linter"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up   - Start all services (PostgreSQL, Neo4j, Redis)"
	@echo "  make docker-down - Stop all services"
	@echo "  make docker-logs - Show logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate     - Run database migrations"
	@echo "  make migrate-rollback - Rollback last migration"
	@echo ""
	@echo "Build:"
	@echo "  make build       - Build production binary"
	@echo "  make clean       - Clean build artifacts"

# 开发环境
dev:
	@echo "Starting development server..."
	cd backend && go run cmd/server/main.go

# 测试
test:
	@echo "Running tests..."
	cd backend && go test -v ./...

# 代码检查
lint:
	@echo "Running linter..."
	cd backend && golangci-lint run

# Docker 管理
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	sleep 10
	@echo "Services are ready!"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker volumes..."
	docker-compose down -v

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	./scripts/init_db.sh

migrate-rollback:
	@echo "Rollback not implemented yet"

# 构建
build:
	@echo "Building production binary..."
	cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../bin/novelforge cmd/server/main.go
	@echo "Binary created at: bin/novelforge"

build-docker:
	@echo "Building Docker image..."
	docker build -t novelforge-ai:latest .

# 清理
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf backend/tmp/

# 安装依赖
install:
	@echo "Installing dependencies..."
	cd backend && go mod download

# 格式化代码
fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...

# 生产环境部署
deploy:
	@echo "Deploying to production..."
	@echo "This is a placeholder. Implement your deployment strategy."

# 健康检查
health:
	@curl -f http://localhost:8080/api/v1/health || echo "Service is down"
