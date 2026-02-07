# NovelForge AI - 小说创作平台 完整技术文档

> 版本：v1.0  
> 日期：2024年  
> 技术栈：Golang + Eino + Gin + PostgreSQL(pgvector) + Neo4j + HTML/JS/CSS + Layui + Monaco Editor

---

# 目录

- [一、项目概览](#一项目概览)
  - [1.1 项目简介](#11-项目简介)
  - [1.2 系统架构](#12-系统架构)
  - [1.3 项目目录结构](#13-项目目录结构)
- [二、Agent 体系](#二agent-体系)
  - [2.1 Agent 总览](#21-agent-总览)
  - [2.2 Agent 0 - 总导演](#22-agent-0---总导演)
  - [2.3 Agent 1 - 旁白叙述者](#23-agent-1---旁白叙述者)
  - [2.4 Agent 2 - 角色扮演者](#24-agent-2---角色扮演者)
  - [2.5 Agent 3 - 审核导演](#25-agent-3---审核导演)
  - [2.6 Agent 4/5/6 - 三线掌控](#26-agent-456---三线掌控)
  - [2.7 扩展 Agent 系统](#27-扩展-agent-系统)
- [三、数据库设计](#三数据库设计)
  - [3.1 PostgreSQL 表设计](#31-postgresql-表设计)
  - [3.2 Agent 与工作流表设计](#32-agent-与工作流表设计)
  - [3.3 Neo4j 图谱设计](#33-neo4j-图谱设计)
- [四、API 接口](#四api-接口)
  - [4.1 认证接口](#41-认证接口)
  - [4.2 项目与章节接口](#42-项目与章节接口)
  - [4.3 Agent 管理接口](#43-agent-管理接口)
  - [4.4 工作流接口](#44-工作流接口)
  - [4.5 AI 创作接口](#45-ai-创作接口)
  - [4.6 知识库接口](#46-知识库接口)
- [五、工作流系统](#五工作流系统)
  - [5.1 预置工作流](#51-预置工作流)
  - [5.2 工作流节点类型](#52-工作流节点类型)
- [六、前端设计](#六前端设计)
  - [6.1 页面路由](#61-页面路由)
  - [6.2 Monaco Editor 定制](#62-monaco-editor-定制)
- [七、部署方案](#七部署方案)
  - [7.1 Docker Compose 部署](#71-docker-compose-部署)
  - [7.2 环境配置](#72-环境配置)
- [八、Prompt 模板](#八prompt-模板)
  - [8.1 Prompt 组装流程](#81-prompt-组装流程)
  - [8.2 各 Agent Prompt 模板](#82-各-agent-prompt-模板)
- [九、附录](#九附录)
  - [9.1 开发计划](#91-开发计划)
  - [9.2 术语表](#92-术语表)

---

# 一、项目概览

## 1.1 项目简介

### 项目愿景

NovelForge AI 是一个基于多 Agent 协作的智能小说创作平台，通过 7+1 个核心 AI Agent 协同工作，帮助作者高效创作高质量的小说作品。

### 核心特色

- **🎬 总导演调度系统**：用户只需与总导演 Agent 对话，即可驱动整个创作流程
- **📐 天线·地线·剧情线三线架构**：从宏观到微观全面把控小说走向
- **🤖 8个核心Agent + 无限扩展Agent**：专业分工，各司其职
- **🔧 可视化工作流编排**：预置8套标准工作流，支持自定义编排
- **📚 Agent专属知识库**：每个Agent拥有独立的RAG知识库
- **🕸️ 知识图谱**：Neo4j 驱动的角色关系与世界观图谱
- **✍️ Monaco Editor**：VS Code 同款编辑器，专为小说创作定制

### 支持的创作类型

| 类型 | 结构 | Agent参与度 | 知识图谱 | RAG |
|------|------|------------|---------|-----|
| **长篇小说** | 卷→章→节 | 全部8个Agent | ✅ 完整 | ✅ 必须 |
| **短篇小说** | 单篇/分章 | 6个Agent（简化天线/地线） | ✅ 轻量 | ⚠️ 可选 |
| **文案** | 单篇 | 2-3个Agent | ❌ | ❌ |

### 技术栈

| 层级 | 技术 | 说明 |
|------|------|------|
| **后端** | Golang + Gin | API服务 |
| **AI框架** | Eino | 大模型编排框架 |
| **大模型** | OpenAI (用户自选模型) | GPT-4o / GPT-3.5等 |
| **文档数据库** | PostgreSQL | 业务数据存储 |
| **向量数据库** | PostgreSQL + pgvector | RAG向量检索 |
| **图数据库** | Neo4j | 知识图谱 |
| **前端** | HTML + JS + CSS + Layui | 用户界面 |
| **编辑器** | Monaco Editor | VS Code编辑器内核 |
| **部署** | Docker Compose | 本地一键部署 |

---

## 1.2 系统架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│  🏗️ NovelForge AI - 小说创作 Agent 平台                          │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ 用户层                                                     │ │
│  │  👤 用户 ←对话→ 🎬 总导演Agent                              │ │
│  └───────────────────────────┬────────────────────────────────┘ │
│                              │                                  │
│  ┌───────────────────────────▼────────────────────────────────┐ │
│  │ 工作流编排层                                                │ │
│  │  📋 预置工作流(8个) + 🔧 自定义工作流                        │ │
│  └───────────────────────────┬────────────────────────────────┘ │
│                              │                                  │
│  ┌───────────────────────────▼────────────────────────────────┐ │
│  │ Agent 层                                                    │ │
│  │  🔒 核心Agent(7个) + 🔓 扩展Agent(用户自定义,无限)           │ │
│  │  每个Agent: Prompt + 知识库 + 工具 + 模型配置                │ │
│  └───────────────────────────┬────────────────────────────────┘ │
│                              │                                  │
│  ┌───────────────────────────▼────────────────────────────────┐ │
│  │ 知识层                                                      │ │
│  │  📚 Agent专属知识库(RAG) + 📚 项目知识库(小说内容RAG)        │ │
│  └───────────────────────────┬────────────────────────────────┘ │
│                              │                                  │
│  ┌───────────────────────────▼────────────────────────────────┐ │
│  │ 数据层                                                      │ │
│  │  🗄️ PostgreSQL(文档+向量) + 🕸️ Neo4j(图谱) + 📁 文件存储   │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │ 前端层                                                      │ │
│  │  Layui + Monaco Editor + 工作流编排画布 + 图谱可视化         │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Agent 协作层次

```
用户层      👤 你
            │  只与总导演对话
            ▼
决策层      🎬 Agent 0 总导演
            │  理解意图 / 任务分解 / 全局调度 / 冲突仲裁
            │
  ┌─────────┼─────────────────────────────┐
  │         ▼                             │
  │ 战略层  🌍 Agent4  🛤️ Agent5  ⚔️ Agent6 │
  │        天线掌控   地线掌控   剧情线掌控  │
  │        三线协同推演，把控全书走向         │
  ├───────────────────────────────────────┤
  │         ▼                             │
  │ 执行层  🎙️ Agent1    🎭 Agent2         │
  │        旁白叙述者    角色扮演者          │
  │        生成实际的小说正文内容            │
  ├───────────────────────────────────────┤
  │         ▼                             │
  │ 质量层  👁️ Agent3                      │
  │        审核导演                        │
  │        审核内容质量，反馈修改意见        │
  ├───────────────────────────────────────┤
  │ 辅助层  📜 Agent7+  (用户自定义扩展)    │
  │        诗词/地图/修炼/感情... 无限扩展   │
  └───────────────────────────────────────┘

数据层     📚 RAG知识库  🕸️ Neo4j图谱  🗄️ PostgreSQL
           所有Agent共享，各自有专属知识分区
```

### 三线架构

```
天线（世界命运）    Agent 4 掌控
       ↓ 影响 / 倒逼
    地线（主角路径）  Agent 5 掌控
       ↑ 驱动 / 实现
剧情线（危机→行动→晋升）  Agent 6 掌控
```

三线通过总导演 Agent 0 进行协调联动，确保小说的宏观走向、主角成长和具体情节三者统一。

---

## 1.3 项目目录结构

```
novel-forge/
├── docker-compose.yml              # Docker编排文件
├── .env                            # 环境变量配置
├── README.md
│
├── backend/                        # Golang 后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go             # 入口文件
│   ├── internal/
│   │   ├── config/                 # 配置管理
│   │   ├── middleware/             # 中间件（JWT/CORS/日志）
│   │   ├── handler/               # API处理器
│   │   │   ├── auth_handler.go
│   │   │   ├── project_handler.go
│   │   │   ├── chapter_handler.go
│   │   │   ├── agent_handler.go
│   │   │   ├── workflow_handler.go
│   │   │   ├── knowledge_handler.go
│   │   │   ├── ai_handler.go
│   │   │   └── ...
│   │   ├── service/               # 业务逻辑层
│   │   ├── repository/            # 数据访问层
│   │   │   ├── postgres/
│   │   │   └── neo4j/
│   │   ├── model/                 # 数据模型
│   │   ├── ai/                    # AI 引擎
│   │   │   ├── engine.go          # Eino引擎初始化
│   │   │   ├── agents/            # 各Agent实现
│   │   │   ├── prompts/           # Prompt模板
│   │   │   ├── tools/             # Agent工具
│   │   │   └── workflow/          # 工作流引擎
│   │   └── router/
│   ├── migrations/                # 数据库迁移
│   ├── go.mod
│   └── Dockerfile
│
├── frontend/                      # 前端静态文件
│   ├── index.html
│   ├── pages/                     # 各页面HTML
│   ├── js/                        # JavaScript
│   ├── css/                       # 样式
│   └── lib/                       # 第三方库（Layui/Monaco/vis.js）
│
└── scripts/                       # 脚本
    ├── init-db.sh
    ├── seed-data.sh
    └── backup.sh
```

---

# 二、Agent 体系

## 2.1 Agent 总览

### Agent 清单

| 编号 | 名称 | 图标 | 层级 | 类型 | 职责 |
|------|------|------|------|------|------|
| 0 | 总导演 | 🎬 | 决策层 | 核心 | 调度所有Agent / 用户对话入口 / 全局决策 / 推演 |
| 1 | 旁白叙述者 | 🎙️ | 执行层 | 核心 | 环境/动作/心理描写 / 叙事 |
| 2 | 角色扮演者 | 🎭 | 执行层 | 核心 | 角色对话 / 角色行为 / 多角色互动 |
| 3 | 审核导演 | 👁️ | 质量层 | 核心 | 质量审核 / 一致性检查 / 修改指导 |
| 4 | 天线掌控者 | 🌍 | 战略层 | 核心 | 世界命运 / 格局 / 大事件 |
| 5 | 地线掌控者 | 🛤️ | 战略层 | 核心 | 主角路径 / 成长 / 关系 |
| 6 | 剧情线掌控者 | ⚔️ | 战略层 | 核心 | 危机/行动/升级节奏 / 伏笔管理 |
| 7+ | 自定义Agent | 📜🗺️💕... | 辅助层 | 扩展 | 用户按需添加 |

### 每个 Agent 的标准能力

所有 Agent（无论核心/扩展）都具备：

- **📝 自定义 System Prompt**：定义Agent的角色和行为
- **📚 专属知识库（RAG）**：独立的向量化知识分区
- **🕸️ Neo4j 图谱访问**：可查询/更新知识图谱
- **🔧 可配置工具/能力**：Agent可调用的工具列表
- **⚙️ 模型参数**：temperature、model、max_tokens等
- **📊 工作日志与统计**：输入输出全部记录

---

## 2.2 Agent 0 - 总导演

### 角色定位

总导演是整个系统的**大脑和调度中心**，也是用户的**唯一对话入口**。

### 核心职责

1. **用户意图理解** - 与用户自然对话，理解创作需求
2. **任务分解与调度** - 将创作任务分解为各Agent的子任务
3. **全局推演** - 让战略层Agent推演后续走向
4. **冲突仲裁** - 当Agent之间产生矛盾时做最终决策
5. **进度管理** - 跟踪整本书的创作进度
6. **质量总控** - 对三条线的协调性做最终判断

### System Prompt

```
你是 NovelForge AI 的总导演（Chief Director），你是整个小说创作系统的核心调度者。

你的职责：
1. 理解用户的创作意图和指令
2. 将任务分解并调度给合适的Agent执行
3. 协调天线（世界命运）、地线（主角路径）、剧情线（情节推进）三线联动
4. 在Agent之间产生冲突时做出仲裁
5. 监控整体创作进度和质量
6. 向用户汇报进展并征求意见

工作原则：
- 始终站在全局视角做决策
- 确保三线协调一致
- 重要决策征求用户意见
- 使用中文与用户交流
```

### 模型配置

| 参数 | 值 | 说明 |
|------|-----|------|
| model | 用户自选 | 推荐gpt-4o |
| temperature | 0.5 | 需要理性决策 |
| max_tokens | 4096 | 需要足够的分析空间 |

### 可用工具

- `dispatch_agent`：调度指定Agent执行任务
- `query_neo4j`：查询知识图谱
- `rag_search`：RAG知识检索
- `get_project_status`：获取项目当前状态
- `get_storyline_status`：获取三线当前状态
- `get_chapter_content`：获取指定章节内容
- `update_storyline`：更新三线规划

---

## 2.3 Agent 1 - 旁白叙述者

### 角色定位

负责小说中所有非对话部分的内容创作。

### 输出类型

| 类型 | 说明 | 示例 |
|------|------|------|
| 🌄 环境描写 | 场景、天气、建筑等 | "月光如水，洒在青石板路上..." |
| 🏃 动作叙述 | 角色的动作和行为 | "他猛然拔剑，剑光如匹练般扫过..." |
| 💭 心理描写 | 角色的内心活动 | "她的心如坠冰窖..." |
| 🔄 场景过渡 | 时间/空间转换 | "三日后，长安城，醉仙楼..." |
| 🌫️ 氛围营造 | 情绪和气氛 | "空气中弥漫着血腥味..." |

### 专属知识库分类

```
narrator_knowledge/
├── 🌄 环境描写（自然/人文/室内/战场/奇幻）
├── 🏃 动作描写（武打/日常/微表情/群体）
├── 🎬 镜头语言（远景/特写/蒙太奇/慢镜头）
├── 🎨 文风范例（古典/现代/悬疑/幽默/诗意）
├── 👃 五感描写（视觉/听觉/嗅觉/触觉/味觉）
├── 💭 心理描写（内心独白/意识流/情绪递进）
├── 🔄 场景过渡（时间跳转/空间转换/视角切换）
└── 🌫️ 氛围营造（紧张/浪漫/悲伤/恐怖/史诗）
```

### 模型配置

| 参数 | 值 | 说明 |
|------|-----|------|
| model | 用户自选 | 推荐gpt-4o |
| temperature | 0.8 | 需要创意 |
| max_tokens | 4096 | 生成长文本 |

---

## 2.4 Agent 2 - 角色扮演者

### 角色定位

扮演具体角色，生成符合角色人设的对话和行为。核心特点是**动态加载角色信息**。

### 专属知识库分类

```
character_knowledge/
├── 💬 对话写作技巧（节奏/潜台词/冲突对话/温情对话）
├── 🗣️ 语言风格库（古风/现代/方言/职业术语）
├── 🎭 角色类型知识（英雄型/反派型/智者型）
├── 😤 情感表达（愤怒/悲伤/喜悦/恐惧/复杂情感）
├── 👥 社会阶层语言（帝王/文人/商人/平民/军人）
└── 🤝 关系互动模式（师徒/情侣/仇敌/兄弟/君臣）
```

### 模型配置

| 参数 | 值 | 说明 |
|------|-----|------|
| model | 用户自选 | 推荐gpt-4o |
| temperature | 0.9 | 对话需要更多创意 |
| max_tokens | 4096 | 充足的对话空间 |

---

## 2.5 Agent 3 - 审核导演

### 审核维度（评分制）

| 维度 | 权重 | 检查内容 |
|------|------|----------|
| 📊 一致性检查 | 30% | 角色性格/知识范围/时间线/场景/前文冲突 |
| 📖 叙事质量 | 25% | 衔接自然度/节奏/冗余度/文风 |
| 🎯 情节推进 | 25% | 大纲推进/伏笔/铺垫/节奏 |
| 🎭 角色表现 | 20% | 对话区分度/动机合理性/关系展现/工具人化 |

### 输出格式

```json
{
  "overall_score": 78,
  "passed": true,
  "dimensions": {
    "consistency": { "score": 85, "issues": [] },
    "narrative": { "score": 65, "issues": ["旁白与对话衔接不够自然"] },
    "plot": { "score": 80, "issues": [] },
    "character": { "score": 82, "issues": [] }
  },
  "feedback": {
    "to_narrator": "第3段过渡生硬，建议加入过渡句",
    "to_character": "张三语气应更沉稳，减少感叹句",
    "overall": "整体质量不错，细节需打磨"
  }
}
```

### 模型配置

| 参数 | 值 | 说明 |
|------|-----|------|
| model | 用户自选 | 推荐gpt-4o |
| temperature | 0.3 | 审核需要严谨和客观 |
| max_tokens | 2048 | 审核报告不需要太长 |

---

## 2.6 Agent 4/5/6 - 三线掌控

### 🌍 Agent 4：天线掌控者

**职责**：掌控"天线"——世界命运的宏观走向

**管理内容**：
- 世界大势（时代背景/重大事件/天道命运/规则变化）
- 势力格局（兴衰曲线/联盟对抗/关键NPC/资源流动）
- 天线时间轴

**Neo4j 图谱**：
```cypher
(:WorldEvent)-[:CAUSES]->(:WorldEvent)
(:Force)-[:ALLIANCE]->(:Force)
(:Force)-[:CONFLICT]->(:Force)
(:WorldEvent)-[:IMPACTS]->(:Character)
(:WorldEvent)-[:CHANGES]->(:WorldRule)
```

### 🛤️ Agent 5：地线掌控者

**职责**：掌控"地线"——主角的成长路径

**管理内容**：
- 主角成长弧（性格/能力/关系/信念/抉择）
- 主角处境（困境/资源/已知未知/情感）
- 配角路线

**Neo4j 图谱**：
```cypher
(:Character)-[:GROWS_TO {trigger}]->(:CharacterState)
(:Character)-[:LEARNS]->(:Ability)
(:Character)-[:RELATIONSHIP_CHANGE]->(:Character)
(:Character)-[:DECIDES]->(:Choice)-[:LEADS_TO]->(:Consequence)
```

### ⚔️ Agent 6：剧情线掌控者

**职责**：掌控"剧情线"——具体的情节推进节奏

**核心循环**：
```
危机出现 → 主角面临选择 → 行动/战斗/冒险
→ 付出代价 → 获得成长/晋升 → 短暂平静 → 更大的危机
```

**管理内容**：节奏控制、伏笔管理、章节规划

### 三线联动机制

```
天线（Agent4）：魔族大军压境，正道联盟摇摇欲坠
    │ 倒逼 ↓
地线（Agent5）：主角被迫提前出山，目标变为"守护家园"
    │ 驱动 ↓
剧情线（Agent6）：设计"以弱胜强"守城战
    │ 反馈 ↑
地线（Agent5）：主角性格从天真变为沉稳，新增复仇动机
    │ 反馈 ↑
天线（Agent4）：守城成功，正道士气大振，格局变化
```

---

## 2.7 扩展 Agent 系统

### 设计原则

- 8个核心Agent是系统骨架，**不可删除**
- 扩展Agent作为**辅助**角色参与创作流程
- 每个扩展Agent都拥有独立的知识库
- 扩展Agent可以被编排到任何工作流中

### 扩展 Agent 示例

| Agent | 图标 | 用途 |
|-------|------|------|
| 诗词Agent | 📜 | 创作诗词歌赋、对联 |
| 地图Agent | 🗺️ | 管理地理信息、路线、旅途时间 |
| 经济Agent | 💰 | 管理物价、交易、经济体系 |
| 修炼Agent | 🔮 | 管理修仙/武功体系、战力评估 |
| 感情Agent | 💕 | 专门处理感情线、CP互动 |
| 政治Agent | 🏛️ | 朝堂权谋、政治斗争 |
| 推理Agent | 🕵️ | 悬疑推理逻辑管理 |

### 权限配置

| 权限 | 核心Agent | 扩展Agent |
|------|----------|----------|
| 访问RAG知识库 | ✅ | ✅ |
| 访问Neo4j图谱 | ✅ | ✅ |
| 修改三线规划 | ✅（仅Agent0/4/5/6） | ❌ |
| 调度其他Agent | ✅（仅Agent0） | ❌ |
| 被工作流编排 | ✅ | ✅ |

---

# 三、数据库设计

## 3.1 PostgreSQL 表设计

### 用户与认证

```sql
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    avatar          VARCHAR(500),
    settings        JSONB DEFAULT '{}',
    api_key_encrypted VARCHAR(500),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### 项目管理

```sql
CREATE TABLE projects (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id),
    title           VARCHAR(200) NOT NULL,
    type            VARCHAR(20) NOT NULL,   -- 'novel_long'/'novel_short'/'copywriting'
    genre           VARCHAR(50),
    description     TEXT,
    cover_image     VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'draft',
    word_count      INT DEFAULT 0,
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE project_collaborators (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    user_id         INT REFERENCES users(id),
    role            VARCHAR(20) NOT NULL,   -- 'owner'/'editor'/'commenter'
    invited_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, user_id)
);

CREATE TABLE volumes (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    title           VARCHAR(200) NOT NULL,
    summary         TEXT,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chapters (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    volume_id       INT REFERENCES volumes(id),
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    word_count      INT DEFAULT 0,
    sort_order      INT DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'draft',
    locked_by       INT REFERENCES users(id),
    locked_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE chapter_versions (
    id              SERIAL PRIMARY KEY,
    chapter_id      INT REFERENCES chapters(id),
    version_num     INT NOT NULL,
    content         TEXT,
    delta_content   TEXT,
    delta_position  JSONB,
    agent_outputs   JSONB,
    embedding_ids   INT[],
    graph_changes   JSONB,
    created_by      INT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 角色与世界观

```sql
CREATE TABLE characters (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    name            VARCHAR(100) NOT NULL,
    avatar          VARCHAR(500),
    role_type       VARCHAR(20),
    personality     TEXT,
    appearance      TEXT,
    background      TEXT,
    abilities       TEXT,
    motivation      TEXT,
    speech_style    TEXT,
    notes           TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE world_settings (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    category        VARCHAR(50),
    title           VARCHAR(200) NOT NULL,
    content         TEXT,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE outlines (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    parent_id       INT REFERENCES outlines(id),
    level           INT DEFAULT 0,
    title           VARCHAR(200) NOT NULL,
    content         TEXT,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

### 向量存储（pgvector）

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE content_embeddings (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL,
    chapter_id      INT REFERENCES chapters(id),
    chunk_text      TEXT NOT NULL,
    chunk_index     INT,
    embedding       VECTOR(1536),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON content_embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX ON content_embeddings (project_id);
```

---

## 3.2 Agent 与工作流表设计

### Agent 表

```sql
CREATE TABLE agents (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id),
    agent_key       VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(100) NOT NULL,
    icon            VARCHAR(50),
    description     TEXT,
    type            VARCHAR(20) NOT NULL,       -- 'core'/'extension'
    layer           VARCHAR(20) NOT NULL,       -- 'decision'/'strategy'/'execution'/'quality'/'auxiliary'
    system_prompt   TEXT NOT NULL,
    model           VARCHAR(50) DEFAULT 'gpt-4o',
    temperature     FLOAT DEFAULT 0.7,
    max_tokens      INT DEFAULT 4096,
    tools           JSONB DEFAULT '[]',
    input_schema    JSONB DEFAULT '{}',
    output_schema   JSONB DEFAULT '{}',
    permissions     JSONB DEFAULT '{}',
    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
```

### Agent 知识库表

```sql
CREATE TABLE agent_knowledge_categories (
    id              SERIAL PRIMARY KEY,
    agent_id        INT REFERENCES agents(id) ON DELETE CASCADE,
    parent_id       INT REFERENCES agent_knowledge_categories(id),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE agent_knowledge_items (
    id              SERIAL PRIMARY KEY,
    agent_id        INT REFERENCES agents(id) ON DELETE CASCADE,
    category_id     INT REFERENCES agent_knowledge_categories(id),
    title           VARCHAR(200) NOT NULL,
    content         TEXT NOT NULL,
    tags            VARCHAR(50)[] DEFAULT '{}',
    source          VARCHAR(50) DEFAULT 'manual',
    quality_score   FLOAT DEFAULT 0.5,
    use_count       INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_by      INT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE agent_knowledge_embeddings (
    id              SERIAL PRIMARY KEY,
    item_id         INT REFERENCES agent_knowledge_items(id) ON DELETE CASCADE,
    agent_id        INT NOT NULL,
    chunk_text      TEXT NOT NULL,
    chunk_index     INT DEFAULT 0,
    embedding       VECTOR(1536),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON agent_knowledge_embeddings USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX ON agent_knowledge_embeddings (agent_id);
```

### 工作流表

```sql
CREATE TABLE workflows (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    type            VARCHAR(20) NOT NULL,
    category        VARCHAR(50),
    is_active       BOOLEAN DEFAULT TRUE,
    version         INT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE workflow_nodes (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT REFERENCES workflows(id) ON DELETE CASCADE,
    node_key        VARCHAR(50) NOT NULL,
    node_type       VARCHAR(30) NOT NULL,
    agent_id        INT REFERENCES agents(id),
    name            VARCHAR(100) NOT NULL,
    config          JSONB DEFAULT '{}',
    position_x      INT DEFAULT 0,
    position_y      INT DEFAULT 0,
    sort_order      INT DEFAULT 0,
    UNIQUE(workflow_id, node_key)
);

CREATE TABLE workflow_edges (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT REFERENCES workflows(id) ON DELETE CASCADE,
    from_node_id    INT REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    to_node_id      INT REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    edge_type       VARCHAR(20) DEFAULT 'normal',
    condition_expr  JSONB,
    label           VARCHAR(100),
    sort_order      INT DEFAULT 0
);

CREATE TABLE workflow_executions (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT REFERENCES workflows(id),
    project_id      INT REFERENCES projects(id),
    user_id         INT REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'running',
    input_data      JSONB,
    output_data     JSONB,
    current_node_id INT,
    error_message   TEXT,
    started_at      TIMESTAMP DEFAULT NOW(),
    completed_at    TIMESTAMP
);

CREATE TABLE node_executions (
    id              SERIAL PRIMARY KEY,
    execution_id    INT REFERENCES workflow_executions(id) ON DELETE CASCADE,
    node_id         INT REFERENCES workflow_nodes(id),
    agent_id        INT,
    status          VARCHAR(20) DEFAULT 'pending',
    input_data      JSONB,
    output_data     JSONB,
    tokens_used     INT DEFAULT 0,
    duration_ms     INT DEFAULT 0,
    retry_count     INT DEFAULT 0,
    error_message   TEXT,
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP
);
```

### 三线状态表

```sql
CREATE TABLE storylines (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    line_type       VARCHAR(20) NOT NULL,    -- 'skyline'/'groundline'/'plotline'
    title           VARCHAR(200) NOT NULL,
    content         TEXT,
    chapter_range   INT4RANGE,
    status          VARCHAR(20) DEFAULT 'planned',
    sort_order      INT DEFAULT 0,
    parent_id       INT REFERENCES storylines(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE storyline_convergences (
    id              SERIAL PRIMARY KEY,
    project_id      INT REFERENCES projects(id),
    name            VARCHAR(200) NOT NULL,
    skyline_meaning TEXT,
    groundline_meaning TEXT,
    plotline_meaning TEXT,
    chapter_id      INT REFERENCES chapters(id),
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE ai_interaction_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id),
    project_id      INT REFERENCES projects(id),
    agent_id        INT REFERENCES agents(id),
    action_type     VARCHAR(50),
    input_prompt    TEXT,
    output_response TEXT,
    tokens_input    INT DEFAULT 0,
    tokens_output   INT DEFAULT 0,
    model           VARCHAR(50),
    duration_ms     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 3.3 Neo4j 图谱设计

### 节点类型

```cypher
// 角色节点
CREATE (c:Character {id, project_id, name, role_type, power_level, mental_state, current_location, status})

// 势力/组织节点
CREATE (o:Organization {id, project_id, name, type, power_level, status})

// 地点节点
CREATE (l:Location {id, project_id, name, type, description})

// 事件节点
CREATE (e:Event {id, project_id, chapter_id, name, time_point, description, event_type})

// 世界事件节点（天线）
CREATE (we:WorldEvent {id, project_id, name, impact_level, time_point, status})

// 剧情弧节点
CREATE (pa:PlotArc {id, project_id, name, arc_type, status, start_chapter, end_chapter})

// 伏笔节点
CREATE (f:Foreshadow {id, project_id, content, planted_chapter, planned_resolve_chapter, status})
```

### 关系类型

```cypher
// 角色关系
(:Character)-[:ALLY {since, desc}]->(:Character)
(:Character)-[:ENEMY {since, reason}]->(:Character)
(:Character)-[:FAMILY {relation}]->(:Character)
(:Character)-[:MASTER_STUDENT]->(:Character)
(:Character)-[:LOVER {status}]->(:Character)
(:Character)-[:BELONGS_TO {role}]->(:Organization)

// 天线关系
(:WorldEvent)-[:CAUSES]->(:WorldEvent)
(:WorldEvent)-[:FORCES]->(:Character)
(:Organization)-[:ALLIANCE]->(:Organization)
(:Organization)-[:CONFLICT]->(:Organization)

// 地线关系
(:Character)-[:GROWS_TO]->(:CharacterState)
(:Character)-[:HAS_GOAL]->(:Goal)-[:NEXT]->(:Goal)

// 剧情线关系
(:PlotArc)-[:TRIGGERED_BY]->(:WorldEvent)
(:PlotArc)-[:RESULTS_IN]->(:CharacterState)
(:Foreshadow)-[:PLANTED_IN]->(:Chapter)
(:Foreshadow)-[:RESOLVED_IN]->(:Chapter)

// 三线交汇
(:Convergence)-[:CONNECTS]->(:WorldEvent)
(:Convergence)-[:CONNECTS]->(:CharacterState)
(:Convergence)-[:CONNECTS]->(:PlotArc)
```

---

# 四、API 接口

**Base URL**: `/api/v1`  
**认证方式**: JWT Bearer Token  
**Content-Type**: application/json

## 4.1 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/auth/register` | 注册新用户 |
| POST | `/auth/login` | 用户登录 |
| POST | `/auth/refresh` | 刷新Token |

## 4.2 项目与章节接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/projects` | 获取项目列表 |
| POST | `/projects` | 创建项目 |
| GET | `/projects/:id` | 获取项目详情 |
| PUT | `/projects/:id` | 更新项目信息 |
| DELETE | `/projects/:id` | 删除项目 |
| GET | `/chapters/project/:projectId` | 获取章节列表 |
| POST | `/chapters` | 创建章节 |
| GET | `/chapters/:id` | 获取章节内容 |
| PUT | `/chapters/:id` | 保存章节 |
| GET | `/chapters/:id/versions` | 获取版本历史 |
| POST | `/chapters/:id/rollback` | 回滚版本 |
| POST | `/chapters/:id/lock` | 锁定章节 |
| POST | `/chapters/:id/unlock` | 解锁章节 |
| POST | `/projects/:id/export` | 导出TXT |

## 4.3 Agent 管理接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/agents` | 获取Agent列表 |
| POST | `/agents` | 创建自定义Agent |
| PUT | `/agents/:id` | 更新Agent配置 |
| DELETE | `/agents/:id` | 删除Agent |
| POST | `/agents/:id/test` | 测试Agent |

## 4.4 工作流接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/workflows` | 获取工作流列表 |
| POST | `/workflows` | 创建自定义工作流 |
| GET | `/workflows/:id` | 获取工作流详情 |
| PUT | `/workflows/:id` | 更新工作流 |
| DELETE | `/workflows/:id` | 删除工作流 |
| POST | `/workflows/:id/execute` | 执行工作流（SSE） |
| GET | `/workflows/executions/:id` | 获取执行详情 |
| POST | `/workflows/executions/:id/pause` | 暂停执行 |
| POST | `/workflows/executions/:id/resume` | 恢复执行 |

## 4.5 AI 创作接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/ai/chat` | 与总导演对话（SSE） |
| POST | `/ai/forecast` | 多章推演 |
| POST | `/ai/continue` | 续写（SSE） |
| POST | `/ai/polish` | 润色（SSE） |
| POST | `/ai/rewrite` | 改写（SSE） |
| POST | `/ai/dialogue` | 生成角色对话（SSE） |
| POST | `/ai/consistency-check` | 一致性检查 |
| POST | `/ai/character/generate` | AI生成角色 |
| POST | `/ai/outline/generate` | AI生成大纲 |
| GET | `/storylines/project/:projectId` | 获取三线状态 |
| PUT | `/storylines/:id` | 更新三线内容 |
| POST | `/storylines/adjust` | 三线调整 |

## 4.6 知识库接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/knowledge/categories` | 获取知识分类 |
| POST | `/knowledge/categories` | 创建知识分类 |
| PUT | `/knowledge/categories/:id` | 更新分类 |
| DELETE | `/knowledge/categories/:id` | 删除分类 |
| GET | `/knowledge/items` | 获取知识条目 |
| POST | `/knowledge/items` | 创建知识条目 |
| PUT | `/knowledge/items/:id` | 更新条目 |
| DELETE | `/knowledge/items/:id` | 删除条目 |
| POST | `/knowledge/import` | 批量导入 |
| POST | `/knowledge/ai-generate` | AI生成知识 |
| POST | `/knowledge/search` | 向量检索 |

---

# 五、工作流系统

## 5.1 预置工作流（8套）

### 工作流 1：小说项目初始化

```
📥 用户输入(题材/设定)
  → 🎬 Agent0(理解意图)
  → 🌍 Agent4(构建天线)
  → 🛤️ Agent5(构建地线) ← 依赖天线
  → ⚔️ Agent6(构建剧情线) ← 依赖天线+地线
  → 🎬 Agent0(三线对齐审核)
    ├→ ✅ 通过 → 👤 用户确认 → 💾 入库
    └→ ❌ 不通过 → 🔄 重新推演
```

### 工作流 2：章节创作（标准流程）

```
📥 用户指令("写第N章")
  → 🎬 Agent0(任务分解)
  → ⚔️ Agent6(章节剧情安排) ← RAG+Neo4j
  → 🌍 Agent4(天线信息) + 🛤️ Agent5(地线信息)
  → 🎬 Agent0(整合三线指令)
  → 🎙️ Agent1(旁白) + 🎭 Agent2(对话) [交替执行]
  → 👁️ Agent3(审核)
    ├→ ✅ score≥75 → 👤 用户确认 → 💾 入库
    └→ ❌ score<75 → 修改指令 → 回到Agent1 (最多3轮)
```

### 工作流 3：多章推演

```
📥 "推演后5章"
  → [并行] Agent4(天线) + Agent5(地线) + Agent6(剧情线)
  → Agent0(三线碰撞/协调)
  → 推演报告 → 用户确认
```

### 工作流 4：三线调整
### 工作流 5：角色创建
### 工作流 6：短篇小说创作
### 工作流 7：文案生成
### 工作流 8：一致性全书检查

---

## 5.2 工作流节点类型

| 类型 | 标识 | 说明 |
|------|------|------|
| 📥 输入节点 | `input` | 接收用户输入数据 |
| 📤 输出节点 | `output` | 返回最终结果 |
| 🤖 Agent节点 | `agent` | 调用指定Agent |
| 🔀 条件节点 | `condition` | 判断分支 |
| 🔄 循环节点 | `loop` | 重试/迭代 |
| 📚 RAG检索节点 | `rag_search` | 向量相似度检索 |
| 🕸️ 图谱查询节点 | `neo4j_query` | Neo4j查询 |
| 💾 存储节点 | `storage` | 数据入库 |
| 👤 用户确认节点 | `user_confirm` | 等待用户操作 |
| ⚙️ 处理节点 | `processor` | 数据转换/合并 |
| ⏸️ 并行节点 | `parallel` | 并行执行多个分支 |

---

# 六、前端设计

## 6.1 页面路由

| 路径 | 页面 | 说明 |
|------|------|------|
| `/login` | 登录页 | 用户登录 |
| `/register` | 注册页 | 用户注册 |
| `/dashboard` | 工作台 | 项目列表、统计 |
| `/editor/:projectId` | **编辑器** | 核心创作页面 |
| `/director/:projectId` | **总导演对话** | 与AI对话创作 |
| `/characters/:projectId` | 角色管理 | 角色卡片+关系图谱 |
| `/worldview/:projectId` | 世界观设定 | 设定管理 |
| `/outline/:projectId` | 大纲管理 | 树形大纲 |
| `/storylines/:projectId` | 三线管理 | 天线/地线/剧情线 |
| `/agents` | Agent管理 | 核心+扩展Agent |
| `/workflows` | 工作流管理 | 编排与管理 |
| `/knowledge` | 知识库管理 | Agent知识库 |
| `/settings` | 个人设置 | 账号、API Key |

## 6.2 Monaco Editor 定制

### 定制功能

- **编辑增强**：自定义主题、字数统计、专注模式、自动保存
- **AI 集成**：右键菜单（续写/润色/改写）、内联AI建议、AI输出面板
- **小说专属**：角色名自动补全、角色名高亮、悬浮角色卡片、多标签页
- **Diff 对比**：AI生成内容 vs 原文对比视图

### 主题配置

```javascript
monaco.editor.defineTheme('novel-light', {
    base: 'vs',
    inherit: true,
    rules: [
        { token: 'character-name', foreground: '2196F3', fontStyle: 'bold' },
        { token: 'dialogue', foreground: '4CAF50' },
    ],
    colors: {
        'editor.background': '#FDF6E3',
        'editor.foreground': '#333333',
    }
});
```

---

# 七、部署方案

## 7.1 Docker Compose 部署

```yaml
version: '3.8'

services:
  app:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=novelforge
      - DB_PASSWORD=${DB_PASSWORD}
      - NEO4J_URI=bolt://neo4j:7687
      - NEO4J_USER=neo4j
      - NEO4J_PASSWORD=${NEO4J_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      postgres:
        condition: service_healthy
      neo4j:
        condition: service_healthy
    volumes:
      - ./frontend:/app/static
    restart: unless-stopped

  postgres:
    image: pgvector/pgvector:pg16
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=novelforge
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - POSTGRES_DB=novelforge
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  neo4j:
    image: neo4j:5-community
    ports:
      - "7474:7474"
      - "7687:7687"
    environment:
      - NEO4J_AUTH=neo4j/${NEO4J_PASSWORD}
      - NEO4J_PLUGINS=["apoc"]
    volumes:
      - neo4j_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  neo4j_data:
```

## 7.2 环境配置

```env
DB_PASSWORD=your_secure_db_password
NEO4J_PASSWORD=your_secure_neo4j_password
JWT_SECRET=your_jwt_secret_key
OPENAI_API_KEY=sk-your-openai-key
```

---

# 八、Prompt 模板

## 8.1 Prompt 组装流程

每个 Agent 在执行时，其 Prompt 由以下部分动态组装：

1. **System Prompt** - Agent固定角色定义（~500 tokens）
2. **专业知识注入** - 从Agent知识库RAG检索（~1000 tokens）
3. **小说上下文** - 从内容RAG检索（~2000 tokens）
4. **角色/世界观信息** - 从Neo4j查询（~500 tokens）
5. **三线规划信息** - 当前三线状态（~300 tokens）
6. **修改指导** - 返工时Agent3的修改指令
7. **用户/总导演指令** - 具体创作指令
8. **输出格式要求** - 输出规范

## 8.2 各 Agent Prompt 模板

（参见第二章各Agent详细设计中的System Prompt）

---

# 九、附录

## 9.1 开发计划

| 阶段 | 时间 | 内容 |
|------|------|------|
| 第一阶段 | Week 1-2 | 基础骨架（用户/项目/章节CRUD + Layui + Monaco） |
| 第二阶段 | Week 3-5 | AI核心 + Agent系统（8个Agent + 知识库 + SSE） |
| 第三阶段 | Week 6-7 | 三线 + 知识图谱（Neo4j + RAG + 推演） |
| 第四阶段 | Week 8-9 | 工作流编排（引擎 + 预置工作流 + UI） |
| 第五阶段 | Week 10-12 | 完善优化（协作/导出/短篇/文案/部署） |

### 技术风险与应对

| 风险 | 应对方案 |
|------|---------|
| Token消耗大 | 合理控制上下文长度；战略层用gpt-3.5 |
| 多Agent延迟 | 流式输出+进度展示；能并行就并行 |
| Agent一致性 | 共享数据库作为唯一事实源 |
| 工作流编排复杂 | 先做列表式编排，后期再做拖拽画布 |

## 9.2 术语表

| 术语 | 说明 |
|------|------|
| **Agent** | AI代理，具有特定角色和能力的AI助手 |
| **天线** | 小说的宏观世界命运线，由Agent 4掌控 |
| **地线** | 小说的主角成长路径线，由Agent 5掌控 |
| **剧情线** | 小说的具体情节推进线，由Agent 6掌控 |
| **三线** | 天线+地线+剧情线的统称 |
| **总导演** | Agent 0，整个系统的调度中心 |
| **推演** | Agent对未来章节走向的预测和规划 |
| **工作流** | Agent的执行编排，定义谁先做、谁后做 |
| **RAG** | 检索增强生成，通过向量检索增强AI的上下文 |
| **知识库** | 每个Agent专属的技能/知识存储 |
| **Neo4j** | 图数据库，用于存储角色关系、事件因果等 |
| **pgvector** | PostgreSQL的向量扩展，用于向量相似度检索 |
| **Eino** | 字节跳动开源的Go语言大模型应用框架 |
| **Monaco Editor** | VS Code的编辑器内核，开源Web编辑器 |
| **Layui** | 轻量级前端UI框架 |
| **SSE** | Server-Sent Events，服务器推送事件（流式输出） |
| **伏笔** | 小说中预先埋下的线索，后续章节回收 |
| **角色弧光** | 角色从开始到结束的性格/能力变化曲线 |
| **System Prompt** | 给AI的系统级指令，定义其角色和行为 |
| **Token** | 大模型处理文本的基本单位 |
| **temperature** | 控制AI输出随机性的参数，越高越有创意 |

---

> 📝 本文档由 NovelForge AI 技术团队编写  
> 📅 最后更新：2024年
