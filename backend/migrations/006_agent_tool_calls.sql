-- Agent 工具调用日志表
CREATE TABLE IF NOT EXISTS agent_tool_calls (
    id SERIAL PRIMARY KEY,
    agent_id INT REFERENCES agents(id) ON DELETE CASCADE,
    tool_name VARCHAR(100) NOT NULL,
    input_params JSONB DEFAULT '{}',
    output_result JSONB DEFAULT '{}',
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    duration_ms INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_agent_tool_calls_agent_id ON agent_tool_calls(agent_id);
CREATE INDEX idx_agent_tool_calls_tool_name ON agent_tool_calls(tool_name);
CREATE INDEX idx_agent_tool_calls_created_at ON agent_tool_calls(created_at);

-- Agent 知识库分类种子数据
-- Agent 1: 旁白叙述者的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(1, '环境描写', '自然/人文/室内/战场/奇幻场景描写技巧', 1),
(1, '动作描写', '武打/日常/微表情/群体动作描写', 2),
(1, '镜头语言', '远景/特写/蒙太奇/慢镜头等电影化叙事', 3),
(1, '文风范例', '古典/现代/悬疑/幽默/诗意等不同风格', 4),
(1, '五感描写', '视觉/听觉/嗅觉/触觉/味觉的细腻描写', 5),
(1, '心理描写', '内心独白/意识流/情绪递进', 6),
(1, '场景过渡', '时间跳转/空间转换/视角切换', 7),
(1, '氛围营造', '紧张/浪漫/悲伤/恐怖/史诗感的渲染', 8)
ON CONFLICT DO NOTHING;

-- Agent 2: 角色扮演者的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(2, '对话写作技巧', '节奏/潜台词/冲突对话/温情对话', 1),
(2, '语言风格库', '古风/现代/方言/职业术语', 2),
(2, '角色类型知识', '英雄型/反派型/智者型/少年型', 3),
(2, '情感表达', '愤怒/悲伤/喜悦/恐惧/复杂情感', 4),
(2, '社会阶层语言', '帝王/文人/商人/平民/军人语言特征', 5),
(2, '关系互动模式', '师徒/情侣/仇敌/兄弟/君臣互动', 6)
ON CONFLICT DO NOTHING;

-- Agent 3: 审核导演的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(3, '一致性检查规则', '角色/时间线/场景/前文矛盾检查标准', 1),
(3, '叙事质量标准', '衔接/节奏/冗余度评估标准', 2),
(3, '情节推进标准', '大纲符合度/伏笔/铺垫检查', 3),
(3, '角色表现标准', '对话区分度/动机合理性评估', 4)
ON CONFLICT DO NOTHING;

-- Agent 4: 天线掌控者的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(4, '世界大势设计', '时代背景/重大事件/天道命运规划', 1),
(4, '势力格局', '兴衰曲线/联盟对抗/资源流动', 2),
(4, '天线时间轴', '世界大事件的时间线管理', 3)
ON CONFLICT DO NOTHING;

-- Agent 5: 地线掌控者的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(5, '主角成长弧', '性格/能力/关系/信念/抉择演变', 1),
(5, '主角处境', '困境/资源/已知未知/情感状态', 2),
(5, '配角路线', '重要配角的成长和关系变化', 3)
ON CONFLICT DO NOTHING;

-- Agent 6: 剧情线掌控者的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(6, '节奏控制', '危机-行动-成长循环的节奏设计', 1),
(6, '伏笔管理', '伏笔埋设和回收技巧', 2),
(6, '章节规划', '每章情节点和冲突设计', 3)
ON CONFLICT DO NOTHING;

-- Agent 0: 总导演的知识库分类
INSERT INTO agent_knowledge_categories (agent_id, name, description, sort_order) VALUES
(0, '任务分解策略', '如何将用户指令分解为Agent任务', 1),
(0, '冲突仲裁案例', '三线冲突的历史仲裁案例', 2),
(0, '调度优化', 'Agent调度顺序和并行策略', 3)
ON CONFLICT DO NOTHING;

COMMENT ON TABLE agent_tool_calls IS 'Agent工具调用日志，记录每次工具执行的输入输出';
COMMENT ON COLUMN agent_tool_calls.duration_ms IS '工具执行耗时（毫秒）';
COMMENT ON COLUMN agent_tool_calls.success IS '工具执行是否成功';
