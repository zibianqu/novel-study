# Changelog

所有重要的项目更改将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且项目遵循 [Semantic Versioning](https://semver.org/lang/zh-CN/)。

## [1.0.0] - 2026-02-08

### 新增
- ✅ 完整的用户认证系统 (JWT)
- ✅ 7 个 AI Agent 协作框架
- ✅ 项目和章节管理 CRUD
- ✅ Monaco Editor 集成
- ✅ AI 对话助手（实时和流式）
- ✅ RAG 知识库系统 (pgvector)
- ✅ Neo4j 知识图谱可视化 (D3.js)
- ✅ 三线管理系统（天线/地线/剧情线）
- ✅ 自动保存机制 (30秒)
- ✅ 实时字数统计
- ✅ Docker Compose 部署支持
- ✅ 限流中间件
- ✅ 请求日志中间件
- ✅ 异常恢复中间件
- ✅ CI/CD 配置 (GitHub Actions)
- ✅ Nginx 配置示例
- ✅ 部署指南
- ✅ 完整 API 文档

### Agent 系统
- Agent 0: 总导演 (任务调度)
- Agent 1: 旁白叙述者 (环境/动作/心理描写)
- Agent 2: 角色扮演者 (对话和行为)
- Agent 3: 审核导演 (质量检查)
- Agent 4: 天线掌控者 (世界命运)
- Agent 5: 地线掌控者 (主角成长)
- Agent 6: 剧情线掌控者 (情节推进)

### 技术栈
- **后端**: Go 1.24, Gin, PostgreSQL 16, Neo4j 5.x, Redis 7
- **前端**: Layui 2.8.18, Monaco Editor 0.44.0, D3.js 7.8.5
- **AI**: OpenAI API (gpt-4o, text-embedding-3-small)
- **向量存储**: pgvector
- **图数据库**: Neo4j

### 文档
- README.md: 项目介绍和快速开始
- docs/DEPLOY.md: 详细部署指南
- docs/API.md: 完整 API 文档
- CHANGELOG.md: 版本更新日志

### 工具
- Makefile: 项目管理脚本
- scripts/init_db.sh: 数据库初始化
- Docker 支持: Dockerfile, docker-compose.yml
- CI/CD: GitHub Actions 工作流

## [0.1.0] - 2026-02-01

### 新增
- 初始项目架构
- 基础数据库设计
- 用户认证系统
- 项目 CRUD API

---

**全功能版本 1.0.0 已发布！** 🎉

查看 [GitHub Releases](https://github.com/zibianqu/novel-study/releases) 获取更新。
