import type { Agent, Project, Volume, Workflow, Storyline } from '../types';

export const agents: Agent[] = [
  { id: 0, key: 'director', name: '总导演', icon: '🎬', type: 'core', layer: 'decision', description: '调度所有Agent / 用户对话入口 / 全局决策 / 推演协调', model: 'gpt-4o', temperature: 0.5, maxTokens: 4096, systemPrompt: '你是 NovelForge AI 的总导演...', isActive: true, knowledgeCount: 45 },
  { id: 1, key: 'narrator', name: '旁白叙述者', icon: '🎙️', type: 'core', layer: 'execution', description: '环境/动作/心理描写 / 叙事 / 氛围营造', model: 'gpt-4o', temperature: 0.8, maxTokens: 4096, systemPrompt: '你是一个专业的小说旁白叙述者...', isActive: true, knowledgeCount: 128 },
  { id: 2, key: 'character', name: '角色扮演者', icon: '🎭', type: 'core', layer: 'execution', description: '角色对话 / 角色行为 / 多角色互动 / 性格维持', model: 'gpt-4o', temperature: 0.9, maxTokens: 4096, systemPrompt: '你是一个角色扮演专家...', isActive: true, knowledgeCount: 96 },
  { id: 3, key: 'reviewer', name: '审核导演', icon: '👁️', type: 'core', layer: 'quality', description: '质量审核 / 一致性检查 / 修改指导', model: 'gpt-4o', temperature: 0.3, maxTokens: 2048, systemPrompt: '你是审核导演，负责审核创作内容...', isActive: true, knowledgeCount: 67 },
  { id: 4, key: 'skyline', name: '天线掌控者', icon: '🌍', type: 'core', layer: 'strategy', description: '世界命运 / 格局 / 大事件 / 势力消长', model: 'gpt-4o', temperature: 0.6, maxTokens: 4096, systemPrompt: '你是天线掌控者，负责世界命运线...', isActive: true, knowledgeCount: 53 },
  { id: 5, key: 'groundline', name: '地线掌控者', icon: '🛤️', type: 'core', layer: 'strategy', description: '主角路径 / 成长 / 关系 / 内心蜕变', model: 'gpt-4o', temperature: 0.6, maxTokens: 4096, systemPrompt: '你是地线掌控者，负责主角成长路径...', isActive: true, knowledgeCount: 48 },
  { id: 6, key: 'plotline', name: '剧情线掌控者', icon: '⚔️', type: 'core', layer: 'strategy', description: '危机→行动→升级节奏 / 伏笔管理 / 章节规划', model: 'gpt-4o', temperature: 0.7, maxTokens: 4096, systemPrompt: '你是剧情线掌控者，负责情节推进...', isActive: true, knowledgeCount: 72 },
  { id: 7, key: 'poem', name: '诗词Agent', icon: '📜', type: 'extension', layer: 'auxiliary', description: '专门负责小说中诗词歌赋的创作', model: 'gpt-4o', temperature: 0.9, maxTokens: 2048, systemPrompt: '你是一位精通中国古典诗词的创作大师...', isActive: true, knowledgeCount: 200 },
  { id: 8, key: 'cultivation', name: '修炼Agent', icon: '🔮', type: 'extension', layer: 'auxiliary', description: '管理修仙/武功体系、战力评估', model: 'gpt-3.5-turbo', temperature: 0.5, maxTokens: 2048, systemPrompt: '你是修炼体系专家...', isActive: true, knowledgeCount: 85 },
];

export const projects: Project[] = [
  { id: 1, title: '九天仙途', type: 'novel_long', genre: '仙侠', description: '废柴少年的修仙逆袭之路，天道崩塌，万族林立...', status: 'writing', wordCount: 156800, chapterCount: 42, updatedAt: '10分钟前' },
  { id: 2, title: '都市之巅', type: 'novel_long', genre: '都市', description: '重生归来的商业天才，这一世要改写所有人的命运', status: 'writing', wordCount: 89200, chapterCount: 28, updatedAt: '2小时前' },
  { id: 3, title: '星际迷途', type: 'novel_short', genre: '科幻', description: '一艘失联的星际飞船，船上的AI开始觉醒...', status: 'draft', wordCount: 12400, chapterCount: 5, updatedAt: '昨天' },
  { id: 4, title: '品牌故事文案集', type: 'copywriting', genre: '商业', description: '各类品牌故事和营销文案创作', status: 'writing', wordCount: 8600, chapterCount: 12, updatedAt: '3天前' },
];

export const volumes: Volume[] = [
  {
    id: 1, title: '卷一：初入修途',
    chapters: [
      { id: 1, title: '第一章 废柴少年', wordCount: 3200, status: 'final' },
      { id: 2, title: '第二章 坠崖奇遇', wordCount: 3800, status: 'final' },
      { id: 3, title: '第三章 远古传承', wordCount: 3100, status: 'final' },
      { id: 4, title: '第四章 初次修炼', wordCount: 2900, status: 'final' },
      { id: 5, title: '第五章 追杀之夜', wordCount: 3500, status: 'final' },
    ]
  },
  {
    id: 2, title: '卷二：宗门风云',
    chapters: [
      { id: 6, title: '第六章 天剑宗', wordCount: 3400, status: 'final' },
      { id: 7, title: '第七章 外门弟子', wordCount: 3600, status: 'final' },
      { id: 8, title: '第八章 宗门大比', wordCount: 4200, status: 'draft' },
      { id: 9, title: '第九章 一鸣惊人', wordCount: 0, status: 'draft' },
    ]
  }
];

export const workflows: Workflow[] = [
  { id: 1, name: '小说项目初始化', description: '创建新的长篇小说项目，自动构建天线/地线/剧情线', type: 'system', category: '长篇创作', icon: '🚀', nodes: [], edges: [], isActive: true },
  { id: 2, name: '章节创作（标准流程）', description: '完整的章节创作流程，包含三线协调、旁白对话生成、审核', type: 'system', category: '长篇创作', icon: '✍️', nodes: [], edges: [], isActive: true },
  { id: 3, name: '多章推演', description: '让三线Agent推演未来多章的走向', type: 'system', category: '长篇创作', icon: '🔮', nodes: [], edges: [], isActive: true },
  { id: 4, name: '三线调整', description: '修改天线/地线/剧情线规划，评估影响范围', type: 'system', category: '管理', icon: '🔧', nodes: [], edges: [], isActive: true },
  { id: 5, name: '角色创建', description: '多Agent协作创建完整的角色设定', type: 'system', category: '管理', icon: '👤', nodes: [], edges: [], isActive: true },
  { id: 6, name: '短篇小说创作', description: '简化流程的短篇小说一次性创作', type: 'system', category: '短篇创作', icon: '📝', nodes: [], edges: [], isActive: true },
  { id: 7, name: '文案生成', description: '快速文案生成流程', type: 'system', category: '文案', icon: '📋', nodes: [], edges: [], isActive: true },
  { id: 8, name: '一致性全书检查', description: '检查全书的一致性问题', type: 'system', category: '管理', icon: '🔍', nodes: [], edges: [], isActive: true },
];

export const storylines: Storyline[] = [
  {
    id: 1, type: 'skyline', title: '天线 - 世界命运', content: '修仙界灵气异变，百年大劫将至',
    status: 'active',
    items: [
      { id: 1, title: '灵气异变', content: '天地灵气开始紊乱，部分区域灵气枯竭', chapterRange: '1-10', status: 'completed' },
      { id: 2, title: '三大宗门争夺秘境', content: '上古秘境即将开启，各方势力暗中角力', chapterRange: '11-20', status: 'active' },
      { id: 3, title: '魔族封印松动', content: '远古封印出现裂痕，魔族蠢蠢欲动', chapterRange: '21-30', status: 'planned' },
      { id: 4, title: '正道联盟危机', content: '圣地出现叛徒，正道联盟面临分裂', chapterRange: '31-42', status: 'planned' },
    ]
  },
  {
    id: 2, type: 'groundline', title: '地线 - 主角路径', content: '废脉少年逆天改命的成长之路',
    status: 'active',
    items: [
      { id: 5, title: '废柴觉醒', content: '获得远古传承，从废柴变为天才', chapterRange: '1-5', status: 'completed' },
      { id: 6, title: '宗门历练', content: '加入天剑宗，从外门弟子做起', chapterRange: '6-15', status: 'active' },
      { id: 7, title: '身世之谜', content: '逐步揭开身世谜团，发现惊人真相', chapterRange: '16-25', status: 'planned' },
      { id: 8, title: '扛起大任', content: '面对大劫，承担拯救苍生的使命', chapterRange: '26-42', status: 'planned' },
    ]
  },
  {
    id: 3, type: 'plotline', title: '剧情线 - 升级节奏', content: '危机→行动→晋升的循环推进',
    status: 'active',
    items: [
      { id: 9, title: '初始篇', content: '被欺负→坠崖→获传承→首次觉醒', chapterRange: '1-5', status: 'completed' },
      { id: 10, title: '宗门篇', content: '入门→修炼→大比→崭露头角', chapterRange: '6-15', status: 'active' },
      { id: 11, title: '秘境篇', content: '探险→危机→突破→揭秘', chapterRange: '16-25', status: 'planned' },
      { id: 12, title: '大战篇', content: '集结→大战→牺牲→最终突破', chapterRange: '26-42', status: 'planned' },
    ]
  }
];

export const sampleChapterContent = `　　夜色如墨，月光洒在天剑宗的演武场上，将青石地面映得泛着清冷的光泽。

　　林远站在演武场的角落，手中紧握着一柄普通的铁剑。剑身上满是细小的豁口，这是他从杂役堂借来的练功用剑——连一把像样的兵器都没有，这就是一个外门弟子的处境。

　　"又来了。"他低声自语，目光落在演武场中央。

　　那里，几名内门弟子正在切磋剑法。剑光如匹练，破空之声不绝于耳。每一招每一式都蕴含着深厚的灵力，将周围的空气搅动得如同沸水。

　　"林远！"一道不耐烦的声音从身后传来。

　　他转过身，看到了王昊——内门弟子中最看不起他的人之一。王昊身后还跟着两个同伴，三人脸上都挂着讥讽的笑容。

　　"一个废脉的家伙，天天跑来演武场看什么？"王昊抱着胸，下巴微扬，"看了就能学会吗？"

　　林远没有说话，只是默默地将铁剑收好，准备离开。

　　"站住。"王昊伸手拦住了他，凑近了一步，压低声音道："听说宗门大比快开始了。你不会真打算参加吧？"

　　林远抬起头，平静地看着他："大比规则说，所有弟子皆可报名。"

　　王昊愣了一瞬，随即大笑起来。他身后的两人也跟着笑了，笑声在夜风中显得格外刺耳。

　　"好，好好好！"王昊连说了三个好字，眼中闪过一丝冷意，"那我等着在大比上，亲手教你什么叫差距。"

　　三人扬长而去，留下林远独自站在月光中。

　　他低头看了看自己的右手——掌心中央，一道古朴的纹路若隐若现，那是他在坠崖时获得的远古传承留下的印记。

　　"差距……"林远轻声重复了这个词，嘴角微微上扬，"他们不知道，真正的差距，很快就会让所有人看到。"

　　他转身走向后山的密林，那里有他秘密修炼的地方。月光在树梢间洒下斑驳的光影，像是为他铺就了一条通往未来的路。`;
