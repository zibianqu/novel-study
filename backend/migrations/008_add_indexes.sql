-- NovelForge AI - 数据库索引优化
-- 创建时间: 2026-02-08
-- 目的: 提升查询性能30-50%

-- ====================================
-- 用户表索引
-- ====================================

-- 邮箱索引 (登录查询)
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 用户名索引 (登录查询)
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- 创建时间索引 (排序查询)
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);

-- ====================================
-- 项目表索引
-- ====================================

-- 用户ID索引 (查询用户的所有项目)
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

-- 项目状态索引 (按状态筛选)
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);

-- 用户+状态复合索引 (最常用的查询)
CREATE INDEX IF NOT EXISTS idx_projects_user_status ON projects(user_id, status);

-- 创建时间索引
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at DESC);

-- 更新时间索引
CREATE INDEX IF NOT EXISTS idx_projects_updated_at ON projects(updated_at DESC);

-- ====================================
-- 章节表索引
-- ====================================

-- 项目ID索引 (查询项目的所有章节)
CREATE INDEX IF NOT EXISTS idx_chapters_project_id ON chapters(project_id);

-- 卷ID索引
CREATE INDEX IF NOT EXISTS idx_chapters_volume_id ON chapters(volume_id);

-- 项目+排序复合索引 (按顺序查询章节)
CREATE INDEX IF NOT EXISTS idx_chapters_project_sort ON chapters(project_id, sort_order);

-- 章节状态索引
CREATE INDEX IF NOT EXISTS idx_chapters_status ON chapters(status);

-- 创建时间索引
CREATE INDEX IF NOT EXISTS idx_chapters_created_at ON chapters(created_at DESC);

-- ====================================
-- 知识库表索引
-- ====================================

-- 项目+分类复合索引
CREATE INDEX IF NOT EXISTS idx_knowledge_project_category ON knowledge_base(project_id, category);

-- 分类索引
CREATE INDEX IF NOT EXISTS idx_knowledge_category ON knowledge_base(category);

-- ====================================
-- API密钥表索引
-- ====================================

-- 用户ID索引
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);

-- 密钥哈希索引 (验证)
CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);

-- 状态索引
CREATE INDEX IF NOT EXISTS idx_api_keys_status ON api_keys(status);

-- 过期时间索引
CREATE INDEX IF NOT EXISTS idx_api_keys_expires_at ON api_keys(expires_at);

-- ====================================
-- 验证索引创建
-- ====================================

-- 分析表统计信息
ANALYZE users;
ANALYZE projects;
ANALYZE chapters;
ANALYZE knowledge_base;
ANALYZE api_keys;

COMMIT;
