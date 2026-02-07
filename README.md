# NovelForge AI - 智能小说创作平台

> 基于 Golang + Eino + 多 Agent 协作的智能小说创作平台

## 🚧 当前开发阶段

**Week 1-2: 基础功能开发** (进行中)

### 已完成

- [x] 项目目录结构搭建
- [x] Docker Compose 配置 (PostgreSQL + Neo4j + App)
- [x] 数据库连接封装
- [x] JWT 认证中间件
- [x] 用户注册/登录 API
- [x] 基础配置管理

### 待完成

- [ ] 项目管理 CRUD API
- [ ] 章节管理 CRUD API
- [ ] 前端 Layui 框架集成
- [ ] Monaco Editor 编辑器集成

## 🚀 快速开始

### 1. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入实际配置
```

### 2. 启动服务

```bash
docker-compose up -d
```

### 3. 访问服务

- **前端**: http://localhost:8080
- **API**: http://localhost:8080/api/v1
- **Neo4j 控制台**: http://localhost:7474

## 📚 API 文档

### 认证接口

#### 注册
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

#### 登录
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

#### 获取用户信息
```bash
GET /api/v1/profile
Authorization: Bearer <token>
```

## 💻 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Golang 1.21 + Gin |
| AI 框架 | Eino (计划中) |
| 数据库 | PostgreSQL 16 + pgvector |
| 图数据库 | Neo4j 5 |
| 前端 | HTML/JS/CSS + Layui |
| 编辑器 | Monaco Editor (计划中) |
| 部署 | Docker Compose |

## 📁 项目结构

```
novel-study/
├── backend/              # Go 后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   └── repository/
│   ├── migrations/
│   ├── go.mod
│   └── Dockerfile
├── frontend/            # 前端静态文件
│   └── index.html
├── docker-compose.yml
├── .env.example
├── 开发计划.md
└── README.md
```

## 📅 开发计划

详细开发计划请查看 [开发计划.md](./%E5%BC%80%E5%8F%91%E8%AE%A1%E5%88%92.md)

- **Week 1-2**: 基础功能开发 ✅ (进行中)
- **Week 3-5**: AI Agent 系统开发
- **Week 6-7**: 知识库与图谱系统
- **Week 8-9**: 前端界面优化
- **Week 10-12**: 测试与部署

## ❓ 常见问题

### 如何重置数据库？

```bash
docker-compose down -v
docker-compose up -d
```

### 如何查看日志？

```bash
docker-compose logs -f app
```

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

## 📝 许可证

MIT License

---

❤️ Made with Golang & AI
