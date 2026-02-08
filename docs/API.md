# NovelForge AI - API 文档

## 📋 目录

- [认证](#认证)
- [用户管理](#用户管理)
- [项目管理](#项目管理)
- [章节管理](#章节管理)
- [AI 功能](#ai-功能)
- [知识库](#知识库)
- [知识图谱](#知识图谱)
- [三线管理](#三线管理)
- [错误代码](#错误代码)

## 🔐 认证

所有需要认证的 API 都需要在请求头中包含 JWT Token：

```
Authorization: Bearer <token>
```

### 用户注册

**POST** `/api/v1/auth/register`

请求体：
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

响应：
```json
{
  "user_id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "created_at": "2026-02-08T00:00:00Z"
}
```

### 用户登录

**POST** `/api/v1/auth/login`

请求体：
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "username": "testuser",
  "expires_at": "2026-02-09T00:00:00Z"
}
```

## 📚 项目管理

### 获取项目列表

**GET** `/api/v1/projects`

响应：
```json
{
  "projects": [
    {
      "id": 1,
      "title": "我的第一部小说",
      "type": "novel_long",
      "genre": "现代都市",
      "description": "一个关于...",
      "status": "writing",
      "word_count": 30000,
      "created_at": "2026-02-01T00:00:00Z",
      "updated_at": "2026-02-08T00:00:00Z"
    }
  ]
}
```

### 创建项目

**POST** `/api/v1/projects`

请求体：
```json
{
  "title": "我的小说",
  "type": "novel_long",
  "genre": "玄幻修仙",
  "description": "项目简介"
}
```

### 更新项目

**PUT** `/api/v1/projects/:id`

### 删除项目

**DELETE** `/api/v1/projects/:id`

## ✍️ 章节管理

### 获取章节列表

**GET** `/api/v1/chapters/project/:projectId`

### 创建章节

**POST** `/api/v1/chapters`

请求体：
```json
{
  "project_id": 1,
  "title": "第一章 开端",
  "content": "章节内容...",
  "sort_order": 1
}
```

### 更新章节

**PUT** `/api/v1/chapters/:id`

### 删除章节

**DELETE** `/api/v1/chapters/:id`

## 🤖 AI 功能

### AI 对话

**POST** `/api/v1/ai/chat`

请求体：
```json
{
  "project_id": 1,
  "message": "帮我写一段主角登场的描写"
}
```

响应：
```json
{
  "content": "生成的内容...",
  "tokens_used": 500,
  "agent": "agent_1"
}
```

### 流式对话 (SSE)

**POST** `/api/v1/ai/chat/stream`

响应格式：Server-Sent Events

```
data: {"content": "生成", "done": false}
data: {"content": "的内容", "done": false}
data: {"done": true}
```

### 生成章节

**POST** `/api/v1/ai/generate/chapter`

请求体：
```json
{
  "project_id": 1,
  "chapter_title": "第一章",
  "outline": "章节大纲..."
}
```

### 质量检查

**POST** `/api/v1/ai/check/quality`

请求体：
```json
{
  "project_id": 1,
  "content": "要检查的内容..."
}
```

响应：
```json
{
  "score": 85,
  "issues": [
    {
      "type": "consistency",
      "message": "角色性格前后不一致",
      "severity": "medium"
    }
  ],
  "suggestions": [
    "建议增加环境描写"
  ]
}
```

## 🧠 知识库

### 获取知识列表

**GET** `/api/v1/knowledge/project/:projectId`

### 创建知识

**POST** `/api/v1/knowledge`

请求体：
```json
{
  "project_id": 1,
  "title": "主角设定",
  "type": "character",
  "content": "张三，25岁男性...",
  "tags": "主角,男性"
}
```

### RAG 检索

**POST** `/api/v1/knowledge/search`

请求体：
```json
{
  "project_id": 1,
  "query": "主角的性格特点",
  "top_k": 3
}
```

响应：
```json
{
  "results": [
    {
      "id": 1,
      "content": "主角张三性格坚毅...",
      "score": 0.95,
      "metadata": {}
    }
  ]
}
```

## 🕸️ 知识图谱

### 获取项目图谱

**GET** `/api/v1/graph/project/:projectId`

响应：
```json
{
  "nodes": [
    {
      "id": "char_001",
      "label": "张三",
      "type": "Character",
      "properties": {}
    }
  ],
  "relations": [
    {
      "source": "char_001",
      "target": "char_002",
      "type": "LOVES"
    }
  ]
}
```

### 创建节点

**POST** `/api/v1/graph/node`

请求体：
```json
{
  "project_id": 1,
  "id": "char_001",
  "label": "张三",
  "type": "Character",
  "properties": {
    "age": 25
  }
}
```

### 创建关系

**POST** `/api/v1/graph/relation`

请求体：
```json
{
  "project_id": 1,
  "source": "char_001",
  "target": "char_002",
  "type": "LOVES"
}
```

## 🎯 三线管理

### 获取三线数据

**GET** `/api/v1/storylines/project/:projectId`

响应：
```json
{
  "storylines": [
    {
      "id": 1,
      "project_id": 1,
      "type": "skyline",
      "title": "大陆战争爆发",
      "description": "三大势力对峙...",
      "sequence": 1,
      "status": "planning"
    }
  ]
}
```

### 创建故事线节点

**POST** `/api/v1/storylines`

请求体：
```json
{
  "project_id": 1,
  "type": "skyline",
  "title": "节点标题",
  "description": "节点描述",
  "sequence": 1
}
```

## ❌ 错误代码

| 代码 | 说明 |
|------|------|
| 400 | 请求参数错误 |
| 401 | 未认证或 Token 过期 |
| 403 | 无权访问 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |

错误响应格式：
```json
{
  "error": "错误描述",
  "code": "ERROR_CODE",
  "details": {}
}
```

## 📊 速率限制

- 认证接口: 5 次/分钟
- AI 接口: 10 次/分钟
- 其他接口: 60 次/分钟

## 🔗 基础 URL

- 开发环境: `http://localhost:8080/api/v1`
- 生产环境: `https://yourdomain.com/api/v1`
