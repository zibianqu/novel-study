-- NovelForge AI - 性能优化索引
-- 创建日期: 2026-02-08

-- =====================================================
-- 1. Users 表索引
-- =====================================================

-- 邮箱索引（用于登录）
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 用户名索引（用于登录和搜索）
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 创建时间索引（用于排序）
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- =====================================================
-- 2. Projects 表索引
-- =====================================================

-- 用户ID索引（最常用的查询）
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

-- 类型索引
CREATE INDEX IF NOT EXISTS idx_projects_type ON projects(type);

-- 状态索引
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);

-- 组合索引（用户 + 状态）
CREATE INDEX IF NOT EXISTS idx_projects_user_status ON projects(user_id, status);

-- 创建时间索引
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);

-- 更新时间索引
CREATE INDEX IF NOT EXISTS idx_projects_updated_at ON projects(updated_at DESC);

-- =====================================================
-- 3. Chapters 表索引
-- =====================================================

-- 项目ID索引（最常用）
CREATE INDEX IF NOT EXISTS idx_chapters_project_id ON chapters(project_id);

-- 卷ID索引
CREATE INDEX IF NOT EXISTS idx_chapters_volume_id ON volumes(id) WHERE volume_id IS NOT NULL;

-- 排序索引
CREATE INDEX IF NOT EXISTS idx_chapters_sort_order ON chapters(project_id, sort_order);

-- 状态索引
CREATE INDEX IF NOT EXISTS idx_chapters_status ON chapters(status);

-- 锁定状态索引
CREATE INDEX IF NOT EXISTS idx_chapters_locked ON chapters(locked_by) WHERE locked_by IS NOT NULL;

-- 更新时间索引
CREATE INDEX IF NOT EXISTS idx_chapters_updated_at ON chapters(updated_at DESC);

-- =====================================================
-- 4. Knowledge Base 索引
-- =====================================================

-- 项目 + 类型组合索引
CREATE INDEX IF NOT EXISTS idx_knowledge_project_type ON knowledge_base(project_id, type);

-- 创建时间索引
CREATE INDEX IF NOT EXISTS idx_knowledge_created_at ON knowledge_base(created_at DESC);

-- =====================================================
-- 5. Storylines 索引
-- =====================================================

-- 项目 + 类型索引
CREATE INDEX IF NOT EXISTS idx_storylines_project_type ON storylines(project_id, line_type);

-- 状态索引
CREATE INDEX IF NOT EXISTS idx_storylines_status ON storylines(status);

-- =====================================================
-- 6. Agents 索引
-- =====================================================

-- Agent Key 索引（唯一性）
CREATE UNIQUE INDEX IF NOT EXISTS idx_agents_key ON agents(agent_key);

-- 类型索引
CREATE INDEX IF NOT EXISTS idx_agents_type ON agents(type);

-- 激活状态索引
CREATE INDEX IF NOT EXISTS idx_agents_active ON agents(is_active);

-- =====================================================
-- 7. Agent Knowledge 索引
-- =====================================================

-- Agent ID + Category 索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_agent_category ON agent_knowledge_items(agent_id, category_id);

-- 激活状态索引
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_active ON agent_knowledge_items(is_active);

-- 标签索引（GIN索引用于数组）
CREATE INDEX IF NOT EXISTS idx_agent_knowledge_tags ON agent_knowledge_items USING GIN(tags);

-- =====================================================
-- 8. Workflows 索引
-- =====================================================

-- 用户 + 类型索引
CREATE INDEX IF NOT EXISTS idx_workflows_user_type ON workflows(user_id, type);

-- 激活状态索引
CREATE INDEX IF NOT EXISTS idx_workflows_active ON workflows(is_active);

-- =====================================================
-- 9. Workflow Executions 索引
-- =====================================================

-- Workflow ID + 状态索引
CREATE INDEX IF NOT EXISTS idx_workflow_exec_workflow_status ON workflow_executions(workflow_id, status);

-- 项目ID索引
CREATE INDEX IF NOT EXISTS idx_workflow_exec_project ON workflow_executions(project_id);

-- 开始时间索引
CREATE INDEX IF NOT EXISTS idx_workflow_exec_started ON workflow_executions(started_at DESC);

-- =====================================================
-- 10. AI Interaction Logs 索引
-- =====================================================

-- 项目ID + 创建时间索引
CREATE INDEX IF NOT EXISTS idx_ai_logs_project_time ON ai_interaction_logs(project_id, created_at DESC);

-- Agent ID 索引
CREATE INDEX IF NOT EXISTS idx_ai_logs_agent ON ai_interaction_logs(agent_id);

-- 用户ID 索引
CREATE INDEX IF NOT EXISTS idx_ai_logs_user ON ai_interaction_logs(user_id);

-- =====================================================
-- 分析表统计信息
-- =====================================================

-- 为所有表更新统计信息
ANALYZE users;
ANALYZE projects;
ANALYZE chapters;
ANALYZE knowledge_base;
ANALYZE storylines;
ANALYZE agents;
ANALYZE agent_knowledge_items;
ANALYZE workflows;
ANALYZE workflow_executions;
ANALYZE ai_interaction_logs;

-- =====================================================
-- 索引创建完成
-- =====================================================

COMMENT ON SCHEMA public IS '性能优化索引已添加 - 2026-02-08';
