package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations 执行所有数据库迁移
func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	migrations := []struct {
		name string
		sql  string
	}{
		{"001_extensions", migration001Extensions},
		{"002_users", migration002Users},
		{"003_projects", migration003Projects},
		{"004_chapters", migration004Chapters},
		{"005_characters_worldview", migration005CharactersWorldview},
		{"006_agents", migration006Agents},
		{"007_knowledge", migration007Knowledge},
		{"008_workflows", migration008Workflows},
		{"009_storylines", migration009Storylines},
		{"010_ai_logs", migration010AILogs},
	}

	// 创建迁移记录表
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id      SERIAL PRIMARY KEY,
			name    VARCHAR(100) UNIQUE NOT NULL,
			applied TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("创建迁移记录表失败: %w", err)
	}

	for _, m := range migrations {
		// 检查是否已执行
		var count int
		err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE name=$1", m.name).Scan(&count)
		if err != nil {
			return fmt.Errorf("检查迁移状态失败 [%s]: %w", m.name, err)
		}
		if count > 0 {
			continue
		}

		// 执行迁移
		log.Printf("  执行迁移: %s", m.name)
		if _, err := pool.Exec(ctx, m.sql); err != nil {
			return fmt.Errorf("迁移失败 [%s]: %w", m.name, err)
		}

		// 记录迁移
		if _, err := pool.Exec(ctx, "INSERT INTO schema_migrations (name) VALUES ($1)", m.name); err != nil {
			return fmt.Errorf("记录迁移状态失败 [%s]: %w", m.name, err)
		}
	}

	return nil
}

const migration001Extensions = `
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

const migration002Users = `
CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    avatar          VARCHAR(500) DEFAULT '',
    settings        JSONB DEFAULT '{}',
    api_key_encrypted VARCHAR(500) DEFAULT '',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

const migration003Projects = `
CREATE TABLE IF NOT EXISTS projects (
    id              SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title           VARCHAR(200) NOT NULL,
    type            VARCHAR(20) NOT NULL DEFAULT 'novel_long',
    genre           VARCHAR(50) DEFAULT '',
    description     TEXT DEFAULT '',
    cover_image     VARCHAR(500) DEFAULT '',
    status          VARCHAR(20) DEFAULT 'draft',
    word_count      INT DEFAULT 0,
    settings        JSONB DEFAULT '{}',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_projects_user ON projects(user_id);

CREATE TABLE IF NOT EXISTS project_collaborators (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id         INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(20) NOT NULL DEFAULT 'editor',
    invited_at      TIMESTAMP DEFAULT NOW(),
    UNIQUE(project_id, user_id)
);

CREATE TABLE IF NOT EXISTS volumes (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title           VARCHAR(200) NOT NULL,
    summary         TEXT DEFAULT '',
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_volumes_project ON volumes(project_id);
`

const migration004Chapters = `
CREATE TABLE IF NOT EXISTS chapters (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    volume_id       INT REFERENCES volumes(id) ON DELETE SET NULL,
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    word_count      INT DEFAULT 0,
    sort_order      INT DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'draft',
    locked_by       INT REFERENCES users(id) ON DELETE SET NULL,
    locked_at       TIMESTAMP,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_chapters_project ON chapters(project_id);
CREATE INDEX IF NOT EXISTS idx_chapters_volume ON chapters(volume_id);

CREATE TABLE IF NOT EXISTS chapter_versions (
    id              SERIAL PRIMARY KEY,
    chapter_id      INT NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    version_num     INT NOT NULL,
    content         TEXT DEFAULT '',
    delta_content   TEXT DEFAULT '',
    delta_position  JSONB DEFAULT '{}',
    agent_outputs   JSONB DEFAULT '{}',
    embedding_ids   INT[] DEFAULT '{}',
    graph_changes   JSONB DEFAULT '[]',
    created_by      INT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_chapter_versions ON chapter_versions(chapter_id, version_num);

CREATE TABLE IF NOT EXISTS content_embeddings (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL,
    chapter_id      INT REFERENCES chapters(id) ON DELETE CASCADE,
    chunk_text      TEXT NOT NULL,
    chunk_index     INT DEFAULT 0,
    embedding       VECTOR(1536),
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_content_emb_project ON content_embeddings(project_id);
`

const migration005CharactersWorldview = `
CREATE TABLE IF NOT EXISTS characters (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    avatar          VARCHAR(500) DEFAULT '',
    role_type       VARCHAR(20) DEFAULT 'supporting',
    personality     TEXT DEFAULT '',
    appearance      TEXT DEFAULT '',
    background      TEXT DEFAULT '',
    abilities       TEXT DEFAULT '',
    motivation      TEXT DEFAULT '',
    speech_style    TEXT DEFAULT '',
    notes           TEXT DEFAULT '',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_characters_project ON characters(project_id);

CREATE TABLE IF NOT EXISTS world_settings (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    category        VARCHAR(50) DEFAULT '',
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_world_settings_project ON world_settings(project_id);

CREATE TABLE IF NOT EXISTS outlines (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    parent_id       INT REFERENCES outlines(id) ON DELETE CASCADE,
    level           INT DEFAULT 0,
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_outlines_project ON outlines(project_id);
`

const migration006Agents = `
CREATE TABLE IF NOT EXISTS agents (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id) ON DELETE CASCADE,
    agent_key       VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(100) NOT NULL,
    icon            VARCHAR(50) DEFAULT '🤖',
    description     TEXT DEFAULT '',
    type            VARCHAR(20) NOT NULL DEFAULT 'extension',
    layer           VARCHAR(20) NOT NULL DEFAULT 'auxiliary',
    system_prompt   TEXT NOT NULL DEFAULT '',
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
`

const migration007Knowledge = `
CREATE TABLE IF NOT EXISTS agent_knowledge_categories (
    id              SERIAL PRIMARY KEY,
    agent_id        INT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    parent_id       INT REFERENCES agent_knowledge_categories(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT DEFAULT '',
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_akc_agent ON agent_knowledge_categories(agent_id);

CREATE TABLE IF NOT EXISTS agent_knowledge_items (
    id              SERIAL PRIMARY KEY,
    agent_id        INT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    category_id     INT REFERENCES agent_knowledge_categories(id) ON DELETE SET NULL,
    title           VARCHAR(200) NOT NULL,
    content         TEXT NOT NULL,
    tags            TEXT[] DEFAULT '{}',
    source          VARCHAR(50) DEFAULT 'manual',
    quality_score   FLOAT DEFAULT 0.5,
    use_count       INT DEFAULT 0,
    is_active       BOOLEAN DEFAULT TRUE,
    created_by      INT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_aki_agent ON agent_knowledge_items(agent_id);
CREATE INDEX IF NOT EXISTS idx_aki_category ON agent_knowledge_items(category_id);

CREATE TABLE IF NOT EXISTS agent_knowledge_embeddings (
    id              SERIAL PRIMARY KEY,
    item_id         INT NOT NULL REFERENCES agent_knowledge_items(id) ON DELETE CASCADE,
    agent_id        INT NOT NULL,
    chunk_text      TEXT NOT NULL,
    chunk_index     INT DEFAULT 0,
    embedding       VECTOR(1536),
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ake_agent ON agent_knowledge_embeddings(agent_id);
`

const migration008Workflows = `
CREATE TABLE IF NOT EXISTS workflows (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    description     TEXT DEFAULT '',
    type            VARCHAR(20) NOT NULL DEFAULT 'custom',
    category        VARCHAR(50) DEFAULT '',
    icon            VARCHAR(50) DEFAULT '🔧',
    is_active       BOOLEAN DEFAULT TRUE,
    version         INT DEFAULT 1,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workflow_nodes (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
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

CREATE TABLE IF NOT EXISTS workflow_edges (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_node_id    INT NOT NULL REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    to_node_id      INT NOT NULL REFERENCES workflow_nodes(id) ON DELETE CASCADE,
    edge_type       VARCHAR(20) DEFAULT 'normal',
    condition_expr  JSONB DEFAULT '{}',
    label           VARCHAR(100) DEFAULT '',
    sort_order      INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS workflow_executions (
    id              SERIAL PRIMARY KEY,
    workflow_id     INT NOT NULL REFERENCES workflows(id),
    project_id      INT REFERENCES projects(id),
    user_id         INT NOT NULL REFERENCES users(id),
    status          VARCHAR(20) DEFAULT 'running',
    input_data      JSONB DEFAULT '{}',
    output_data     JSONB DEFAULT '{}',
    current_node_id INT,
    error_message   TEXT DEFAULT '',
    started_at      TIMESTAMP DEFAULT NOW(),
    completed_at    TIMESTAMP
);

CREATE TABLE IF NOT EXISTS node_executions (
    id              SERIAL PRIMARY KEY,
    execution_id    INT NOT NULL REFERENCES workflow_executions(id) ON DELETE CASCADE,
    node_id         INT REFERENCES workflow_nodes(id),
    agent_id        INT,
    status          VARCHAR(20) DEFAULT 'pending',
    input_data      JSONB DEFAULT '{}',
    output_data     JSONB DEFAULT '{}',
    tokens_used     INT DEFAULT 0,
    duration_ms     INT DEFAULT 0,
    retry_count     INT DEFAULT 0,
    error_message   TEXT DEFAULT '',
    started_at      TIMESTAMP,
    completed_at    TIMESTAMP
);
`

const migration009Storylines = `
CREATE TABLE IF NOT EXISTS storylines (
    id              SERIAL PRIMARY KEY,
    project_id      INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    line_type       VARCHAR(20) NOT NULL,
    title           VARCHAR(200) NOT NULL,
    content         TEXT DEFAULT '',
    chapter_start   INT DEFAULT 0,
    chapter_end     INT DEFAULT 0,
    status          VARCHAR(20) DEFAULT 'planned',
    sort_order      INT DEFAULT 0,
    parent_id       INT REFERENCES storylines(id) ON DELETE CASCADE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_storylines_project ON storylines(project_id);

CREATE TABLE IF NOT EXISTS storyline_convergences (
    id                  SERIAL PRIMARY KEY,
    project_id          INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name                VARCHAR(200) NOT NULL,
    skyline_meaning     TEXT DEFAULT '',
    groundline_meaning  TEXT DEFAULT '',
    plotline_meaning    TEXT DEFAULT '',
    chapter_id          INT REFERENCES chapters(id),
    created_at          TIMESTAMP DEFAULT NOW()
);
`

const migration010AILogs = `
CREATE TABLE IF NOT EXISTS ai_interaction_logs (
    id              SERIAL PRIMARY KEY,
    user_id         INT REFERENCES users(id),
    project_id      INT REFERENCES projects(id),
    agent_id        INT REFERENCES agents(id),
    execution_id    INT REFERENCES workflow_executions(id),
    action_type     VARCHAR(50) DEFAULT '',
    input_prompt    TEXT DEFAULT '',
    output_response TEXT DEFAULT '',
    tokens_input    INT DEFAULT 0,
    tokens_output   INT DEFAULT 0,
    model           VARCHAR(50) DEFAULT '',
    duration_ms     INT DEFAULT 0,
    created_at      TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_ai_logs_user ON ai_interaction_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_logs_project ON ai_interaction_logs(project_id);
`
