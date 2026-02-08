# NovelForge AI 部署指南

## 📋 目录

- [系统要求](#系统要求)
- [开发环境部署](#开发环境部署)
- [生产环境部署](#生产环境部署)
- [Docker 部署](#docker-部署)
- [云平台部署](#云平台部署)
- [性能优化](#性能优化)
- [故障排查](#故障排查)

## 🖥️ 系统要求

### 最低配置
- CPU: 2 核
- 内存: 4GB
- 磁盘: 20GB
- OS: Linux / macOS / Windows

### 推荐配置
- CPU: 4 核+
- 内存: 8GB+
- 磁盘: 50GB+ (SSD)
- OS: Ubuntu 22.04 LTS

### 软件依赖
- Go 1.24+
- Docker 24.0+
- Docker Compose 2.20+
- PostgreSQL 16+ (或使用 Docker)
- Neo4j 5.x (或使用 Docker)
- Redis 7+ (可选)

## 🛠️ 开发环境部署

### 1. 克隆项目

```bash
git clone https://github.com/zibianqu/novel-study.git
cd novel-study
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件，填入必要配置：

```env
# OpenAI API Key (必需)
OPENAI_API_KEY=sk-your-api-key-here

# JWT Secret (建议修改)
JWT_SECRET=your_secure_random_string_here

# 数据库密码
DB_PASSWORD=your_secure_password
NEO4J_PASSWORD=your_neo4j_password
```

### 3. 启动数据库

```bash
make docker-up
# 或
docker-compose up -d
```

等待约 30 秒，确保服务启动完成。

### 4. 运行数据库迁移

```bash
make migrate
# 或手动执行
./scripts/init_db.sh
```

### 5. 启动后端服务

```bash
make dev
# 或
cd backend && go run cmd/server/main.go
```

### 6. 访问应用

打开浏览器访问：http://localhost:8080

## 🚀 生产环境部署

### 方案 1: 二进制部署

#### 1. 构建生产二进制

```bash
make build
```

生成的二进制文件位于 `bin/novelforge`。

#### 2. 配置生产环境变量

```bash
export ENVIRONMENT=production
export SERVER_PORT=8080
export OPENAI_API_KEY=sk-xxx
export JWT_SECRET=xxx
# ... 其他配置
```

#### 3. 启动服务

```bash
./bin/novelforge
```

#### 4. 配置 Systemd (推荐)

创建 `/etc/systemd/system/novelforge.service`：

```ini
[Unit]
Description=NovelForge AI Service
After=network.target postgresql.service neo4j.service

[Service]
Type=simple
User=novelforge
WorkingDirectory=/opt/novelforge
EnvironmentFile=/opt/novelforge/.env
ExecStart=/opt/novelforge/bin/novelforge
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable novelforge
sudo systemctl start novelforge
sudo systemctl status novelforge
```

### 方案 2: Docker 部署

#### 1. 构建 Docker 镜像

```bash
make build-docker
# 或
docker build -t novelforge-ai:latest .
```

#### 2. 运行容器

```bash
docker run -d \
  --name novelforge-api \
  -p 8080:8080 \
  --env-file .env \
  --network novelforge_network \
  novelforge-ai:latest
```

### 方案 3: Docker Compose (完整栈)

使用提供的 `docker-compose.yml`：

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## ☁️ 云平台部署

### AWS 部署

#### 使用 ECS (推荐)

1. 将镜像推送到 ECR
2. 创建 ECS Task Definition
3. 配置 ALB + Target Group
4. 部署 ECS Service

#### 使用 EC2

1. 创建 EC2 实例 (Ubuntu 22.04)
2. 安装 Docker
3. 配置安全组 (开放 8080 端口)
4. 部署应用

### Azure 部署

使用 Azure Container Instances 或 Azure Kubernetes Service (AKS)。

### Google Cloud 部署

使用 Cloud Run 或 GKE。

### 阿里云部署

使用容器服务 ACK 或 ECS。

## 🔧 Nginx 反向代理配置

创建 `/etc/nginx/sites-available/novelforge`：

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    # SSL 证书
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    # 静态文件
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # SSE 支持
    location /api/v1/ai/chat/stream {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Connection '';
        proxy_buffering off;
        proxy_cache off;
        chunked_transfer_encoding on;
    }

    # 文件上传限制
    client_max_body_size 50M;

    # Gzip 压缩
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/novelforge /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## ⚡ 性能优化

### 1. 数据库优化

**PostgreSQL**：

```sql
-- 创建必要的索引
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_chapters_project_id ON chapters(project_id);
CREATE INDEX idx_knowledge_project_id ON knowledge_base(project_id);

-- pgvector HNSW 索引
CREATE INDEX idx_vectors_embedding ON knowledge_vectors 
USING hnsw (embedding vector_cosine_ops);
```

**Neo4j**：

```cypher
// 创建索引
CREATE INDEX FOR (n:Character) ON (n.project_id);
CREATE INDEX FOR (n:Location) ON (n.project_id);
```

### 2. Redis 缓存

启用 Redis 缓存用户会话和热数据：

```go
// 缓存用户会话
rdb.Set(ctx, "session:"+sessionID, userData, 24*time.Hour)

// 缓存项目数据
rdb.Set(ctx, "project:"+projectID, projectJSON, 1*time.Hour)
```

### 3. 并发优化

使用 Goroutine 池处理 AI 请求：

```go
var wg sync.WaitGroup
semaphore := make(chan struct{}, 10) // 最多 10 个并发

for _, task := range tasks {
    wg.Add(1)
    go func(t Task) {
        defer wg.Done()
        semaphore <- struct{}{}
        defer func() { <-semaphore }()
        processTask(t)
    }(task)
}
wg.Wait()
```

## 🐛 故障排查

### 数据库连接失败

```bash
# 检查 PostgreSQL
docker logs novelforge_postgres
psql -h localhost -U novelforge -d novelforge_db -c "SELECT 1;"

# 检查 Neo4j
docker logs novelforge_neo4j
cypher-shell -u neo4j -p your_password "RETURN 1;"
```

### 服务无法启动

```bash
# 查看日志
sudo journalctl -u novelforge -f

# 检查端口占用
sudo lsof -i :8080
```

### OpenAI API 错误

- 检查 API Key 是否有效
- 检查余额是否充足
- 检查网络连接

### 性能问题

```bash
# 监控资源使用
top
htop
docker stats

# 查看慢查询
# PostgreSQL
SELECT * FROM pg_stat_activity WHERE state = 'active';
```

## 📊 监控和日志

### 日志级别

- `development`: DEBUG 级别
- `production`: INFO 级别

### 日志位置

- 应用日志: `stdout`
- Nginx 日志: `/var/log/nginx/`
- PostgreSQL 日志: Docker 容器内
- Neo4j 日志: Docker 容器内

### 监控指标

建议监控：
- API 响应时间
- 数据库连接数
- OpenAI API 调用次数和成本
- 内存使用率
- CPU 使用率
- 磁盘 I/O

## 🔒 安全建议

1. **使用 HTTPS**：通过 Let's Encrypt 获取免费证书
2. **定期更新依赖**：`go get -u all`
3. **限制 API 访问频率**：使用 rate limiting
4. **数据库备份**：每天自动备份
5. **敏感信息加密**：不要在代码中硬编码密钥
6. **最小权限原则**：数据库用户使用最小必要权限

## 📞 获取帮助

- GitHub Issues: https://github.com/zibianqu/novel-study/issues
- 文档: https://github.com/zibianqu/novel-study/docs

---

**部署成功后，不要忘记：**
- ⭐ Star 项目
- 📝 报告 Bug
- 💡 提出改进建议
