# NovelForge AI - Makefile
# 更新时间: 2026-02-08

.PHONY: help build run test clean docker-up docker-down migrate lint format security

# 默认目标
help:
	@echo "NovelForge AI - 可用命令:"
	@echo "  make build          - 编译Go后端"
	@echo "  make run            - 运行后端服务"
	@echo "  make test           - 运行测试"
	@echo "  make test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  make lint           - 代码检查"
	@echo "  make format         - 格式化代码"
	@echo "  make security       - 安全扫描"
	@echo "  make docker-up      - 启动所有服务"
	@echo "  make docker-down    - 停止所有服务"
	@echo "  make migrate        - 运行数据库迁移"
	@echo "  make neo4j-init     - 初始化Neo4j索引"
	@echo "  make health         - 健康检查"
	@echo "  make clean          - 清理构建文件"

# 编译
build:
	@echo "🔨 编译后端..."
	cd backend && go build -o bin/server ./cmd/server

# 运行
run:
	@echo "🚀 启动服务..."
	cd backend && go run ./cmd/server/main.go

# 测试
test:
	@echo "🧪 运行测试..."
	cd backend && go test -v ./...

test-coverage:
	@echo "📊 生成测试覆盖率报告..."
	cd backend && go test -v -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告: backend/coverage.html"

# 代码检查
lint:
	@echo "🔍 运行代码检查..."
	cd backend && golangci-lint run ./...

# 格式化
format:
	@echo "✨ 格式化代码..."
	cd backend && gofmt -s -w .
	cd backend && goimports -w .

# 安全扫描
security:
	@echo "🔒 运行安全扫描..."
	cd backend && gosec ./...

# Docker相关
docker-up:
	@echo "🐳 启动Docker服务..."
	docker-compose up -d
	@echo "⏳ 等待服务就绪..."
	sleep 10
	@make health

docker-down:
	@echo "⬇️ 停止Docker服务..."
	docker-compose down

docker-logs:
	docker-compose logs -f backend

docker-rebuild:
	@echo "🔄 重新构建并启动..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# 数据库迁移
migrate:
	@echo "📦 运行数据库迁移..."
	@for file in backend/migrations/*.sql; do \
		echo "执行: $$file"; \
		docker exec -i novel_postgres psql -U postgres -d novel_forge < $$file; \
	done
	@echo "✅ 迁移完成"

# Neo4j索引初始化
neo4j-init:
	@echo "🕸️ 初始化Neo4j索引..."
	docker exec -i novel_neo4j cypher-shell -u neo4j -p neo4j_password < scripts/init_neo4j_indexes.cypher
	@echo "✅ Neo4j索引创建完成"

# 健康检查
health:
	@echo "❤️ 检查服务健康状态..."
	@curl -s http://localhost:8080/api/v1/health | jq '.' || echo "❌ 后端服务未就绪"

# 清理
clean:
	@echo "🧹 清理构建文件..."
	rm -rf backend/bin
	rm -f backend/coverage.out backend/coverage.html
	rm -rf backend/tmp
	@echo "✅ 清理完成"

# 开发环境设置
dev-setup:
	@echo "🛠️ 设置开发环境..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ 已创建 .env 文件，请填写配置"; \
	fi
	cd backend && go mod download
	@echo "✅ 开发环境设置完成"

# 安装工具
install-tools:
	@echo "📦 安装开发工具..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ 工具安装完成"

# 数据库备份
backup-db:
	@echo "💾 备份数据库..."
	@mkdir -p backups
	docker exec novel_postgres pg_dump -U postgres novel_forge > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✅ 备份完成"

# 查看日志
logs:
	docker-compose logs -f --tail=100

# 进入容器
shell-backend:
	docker exec -it novel_backend sh

shell-postgres:
	docker exec -it novel_postgres psql -U postgres -d novel_forge

shell-neo4j:
	docker exec -it novel_neo4j cypher-shell -u neo4j -p neo4j_password

shell-redis:
	docker exec -it novel_redis redis-cli
