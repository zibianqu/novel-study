-- Agent 专属知识库分类表
CREATE TABLE IF NOT EXISTS agent_knowledge_categories (
    id SERIAL PRIMARY KEY,
    agent_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    priority INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(agent_id, name)
);

CREATE INDEX idx_agent_knowledge_categories_agent_id ON agent_knowledge_categories(agent_id);

-- Agent 专属知识条目表（扩展现有 knowledge_base 表的功能）
CREATE TABLE IF NOT EXISTS agent_knowledge_items (
    id SERIAL PRIMARY KEY,
    agent_id INT NOT NULL,
    category_id INT REFERENCES agent_knowledge_categories(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    tags TEXT[],
    quality_score INT DEFAULT 0,
    usage_count INT DEFAULT 0,
    embedding vector(1536), -- OpenAI embedding dimension
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_agent_knowledge_items_agent_id ON agent_knowledge_items(agent_id);
CREATE INDEX idx_agent_knowledge_items_category_id ON agent_knowledge_items(category_id);
CREATE INDEX idx_agent_knowledge_items_embedding ON agent_knowledge_items USING ivfflat (embedding vector_cosine_ops);

-- Agent 1 (旁白叙述者) 8大分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, priority) VALUES
(1, '环境描写', '自然环境、人文环境、室内环境、战场环境、奇幻场景等环境描写技巧和范例', 10),
(1, '动作描写', '武打动作、日常动作、微表情、群体动作等动作描写技巧', 9),
(1, '镜头语言', '远景、近景、特写、蒙太奇等电影化叙事技巧', 8),
(1, '文风范例', '不同风格的文笔范例：古典、现代、玄幻、武侠、科幻等', 7),
(1, '五感描写', '视觉、听觉、嗅觉、触觉、味觉的细腻描写技巧', 9),
(1, '心理描写', '内心独白、心理活动、情绪变化的叙述技巧', 8),
(1, '场景过渡', '时间推移、空间转换、情境切换的衔接技巧', 7),
(1, '氛围营造', '紧张、悬疑、浪漫、悲伤等不同氛围的营造方法', 8);

-- Agent 2 (角色扮演者) 6大分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, priority) VALUES
(2, '对话写作技巧', '日常对话、冲突对话、说服对话、幽默对话等技巧', 10),
(2, '语言风格库', '不同时代、地域、阶层的语言风格和用词习惯', 9),
(2, '角色类型知识', '英雄、反派、导师、伙伴等不同角色原型的对话特点', 8),
(2, '情感表达', '愤怒、喜悦、悲伤、惊讶等情感的对话表现方式', 9),
(2, '社会阶层语言', '贵族、平民、江湖、朝廷等不同阶层的语言特征', 8),
(2, '关系互动模式', '师徒、恋人、仇敌、朋友等不同关系的对话模式', 9);

-- Agent 1 示例知识（环境描写类）
INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, quality_score) 
SELECT 
    1,
    id,
    '古代庭院环境描写',
    '青石板铺就的小径蜿蜒向前，两侧种满了翠竹，风过竹林，沙沙作响。庭院深处，一座八角凉亭静静矗立，飞檐翘角，雕梁画栋。亭下石桌上，还摆着未下完的棋局，黑白棋子落满棋盘，似乎主人刚刚离去不久。',
    ARRAY['古代', '庭院', '静谧', '写意'],
    95
FROM agent_knowledge_categories 
WHERE agent_id = 1 AND name = '环境描写'
LIMIT 1;

INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, quality_score) 
SELECT 
    1,
    id,
    '战场环境描写',
    '天空低垂，乌云翻滚如墨，空气中弥漫着血腥和硝烟的味道。大地被践踏得坑坑洼洼，到处散落着断裂的兵器和残破的旌旗。远处战鼓声震天，近处刀剑碰撞声此起彼伏，夹杂着战士们的呐喊和马匹的嘶鸣。',
    ARRAY['战场', '紧张', '血腥', '动态'],
    92
FROM agent_knowledge_categories 
WHERE agent_id = 1 AND name = '环境描写'
LIMIT 1;

-- Agent 1 示例知识（五感描写类）
INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, quality_score) 
SELECT 
    1,
    id,
    '清晨五感描写',
    '【视觉】晨曦透过薄雾，在山林间洒下斑驳的光影。【听觉】远处传来寺院的钟声，悠扬而空灵，伴随着鸟鸣声此起彼伏。【嗅觉】空气中弥漫着泥土的清香和松针的气息，还夹杂着淡淡的花香。【触觉】清晨的山风拂过面颊，带着微微的凉意，让人精神一振。【味觉】唇齿间还残留着清茶的苦涩和回甘。',
    ARRAY['五感', '清晨', '山林', '综合'],
    98
FROM agent_knowledge_categories 
WHERE agent_id = 1 AND name = '五感描写'
LIMIT 1;

-- Agent 2 示例知识（对话技巧类）
INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, quality_score) 
SELECT 
    2,
    id,
    '冲突对话技巧',
    '"你以为你是谁？"他冷笑一声，"凭什么对我指手画脚？"

"我是你师兄，"对方语气平静，却透着不容置疑的威严，"更是你救命恩人。这两个身份，够不够？"

"够！"他咬牙切齿，"够让我恶心的！"

技巧点评：
1. 使用反问句增强情绪冲突
2. 通过语气词和动作描写展现角色情绪
3. 对话简短有力，避免冗长
4. 善用停顿制造戏剧张力',
    ARRAY['冲突', '师兄弟', '对话技巧', '情绪表达'],
    96
FROM agent_knowledge_categories 
WHERE agent_id = 2 AND name = '对话写作技巧'
LIMIT 1;

INSERT INTO agent_knowledge_items (agent_id, category_id, title, content, tags, quality_score) 
SELECT 
    2,
    id,
    '不同阶层语言差异',
    '【贵族】"本侯今日设宴，诸位赏光，实乃荣幸之至。"

【江湖】"兄弟们今儿个给面子，咱老张感激不尽！来，干了这碗！"

【朝廷】"微臣愚钝，此事关乎国运，不敢妄断，还请陛下圣裁。"

【平民】"哎呀，大伙儿都来啦！快进屋坐，别客气！"

要点：
- 贵族：用词文雅，语气含蓄
- 江湖：直爽豪迈，多用俗语
- 朝廷：谦卑谨慎，讲究礼法
- 平民：朴实自然，生活气息浓',
    ARRAY['阶层', '语言风格', '社会差异'],
    94
FROM agent_knowledge_categories 
WHERE agent_id = 2 AND name = '社会阶层语言'
LIMIT 1;

-- 更新统计
COMMENT ON TABLE agent_knowledge_categories IS 'Agent专属知识库分类表，为每个Agent组织不同类型的专业知识';
COMMENT ON TABLE agent_knowledge_items IS 'Agent专属知识条目表，存储具体的知识内容和示例';
