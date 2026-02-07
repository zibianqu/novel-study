import { useState } from 'react';
import { agents as agentData } from '../data/mockData';
import { Plus, Search, Upload, Sparkles, Book, Tag, BarChart3, ChevronRight, ChevronDown } from 'lucide-react';

interface KnowledgeItem {
  id: number;
  title: string;
  preview: string;
  tags: string[];
  useCount: number;
}

const mockCategories = [
  { id: 1, name: '环境描写', icon: '🌄', count: 45, children: [
    { id: 11, name: '自然环境', count: 18 },
    { id: 12, name: '人文环境', count: 12 },
    { id: 13, name: '室内环境', count: 8 },
    { id: 14, name: '战场环境', count: 7 },
  ]},
  { id: 2, name: '动作描写', icon: '🏃', count: 38, children: [
    { id: 21, name: '武打战斗', count: 22 },
    { id: 22, name: '日常动作', count: 10 },
    { id: 23, name: '微表情', count: 6 },
  ]},
  { id: 3, name: '镜头语言', icon: '🎬', count: 25, children: [] },
  { id: 4, name: '文风范例', icon: '🎨', count: 20, children: [] },
  { id: 5, name: '五感描写', icon: '👃', count: 30, children: [] },
];

const mockItems: KnowledgeItem[] = [
  { id: 1, title: '雨中打斗场景描写技巧', preview: '雨丝可以作为视觉线索引导读者注意力，同时增加场景的动态感...', tags: ['武侠', '打斗', '雨景'], useCount: 23 },
  { id: 2, title: '雪景描写技巧与范例', preview: '大雪纷飞的场景描写要注意层次感，从远到近，从大到小...', tags: ['天气', '雪', '冬季'], useCount: 15 },
  { id: 3, title: '月夜氛围营造方法', preview: '月光作为古典意象，可以营造多种氛围：宁静、孤寂、浪漫...', tags: ['月光', '氛围', '夜景'], useCount: 31 },
  { id: 4, title: '战场描写的镜头切换技巧', preview: '大战场景建议采用"全景-中景-特写"的镜头切换方式...', tags: ['战争', '镜头', '全景'], useCount: 19 },
  { id: 5, title: '密室/封闭空间的压迫感营造', preview: '封闭空间描写要突出五感中的触觉和听觉，放大细微的声响...', tags: ['密室', '压迫', '五感'], useCount: 12 },
];

export default function KnowledgeBase() {
  const [selectedAgent, setSelectedAgent] = useState(agentData[1]); // Narrator
  const [selectedCategory, setSelectedCategory] = useState(1);
  const [expandedCats, setExpandedCats] = useState<Set<number>>(new Set([1, 2]));
  const [searchQuery, setSearchQuery] = useState('');

  const coreAgents = agentData.filter(a => a.type === 'core');
  const extAgents = agentData.filter(a => a.type === 'extension');

  return (
    <div className="flex h-full">
      {/* Left: Agent Selection */}
      <div className="w-56 border-r border-gray-200 bg-white flex flex-col flex-shrink-0">
        <div className="p-3 border-b border-gray-100">
          <h3 className="text-sm font-semibold text-gray-700">选择 Agent</h3>
        </div>
        <div className="flex-1 overflow-y-auto p-2">
          <div className="text-[10px] font-semibold text-gray-400 uppercase px-2 py-1">核心 Agent</div>
          {coreAgents.map(a => (
            <button
              key={a.id}
              onClick={() => setSelectedAgent(a)}
              className={`w-full flex items-center gap-2 px-2 py-2 rounded-lg text-sm transition-colors ${
                selectedAgent.id === a.id ? 'bg-violet-50 text-violet-700 font-medium' : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              <span>{a.icon}</span>
              <span className="truncate">{a.name}</span>
              <span className="ml-auto text-[10px] text-gray-400">{a.knowledgeCount}</span>
            </button>
          ))}
          <div className="text-[10px] font-semibold text-gray-400 uppercase px-2 py-1 mt-2">扩展 Agent</div>
          {extAgents.map(a => (
            <button
              key={a.id}
              onClick={() => setSelectedAgent(a)}
              className={`w-full flex items-center gap-2 px-2 py-2 rounded-lg text-sm transition-colors ${
                selectedAgent.id === a.id ? 'bg-violet-50 text-violet-700 font-medium' : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              <span>{a.icon}</span>
              <span className="truncate">{a.name}</span>
              <span className="ml-auto text-[10px] text-gray-400">{a.knowledgeCount}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Middle: Categories */}
      <div className="w-52 border-r border-gray-200 bg-gray-50 flex flex-col flex-shrink-0">
        <div className="p-3 border-b border-gray-100 flex items-center justify-between">
          <h3 className="text-sm font-semibold text-gray-700">知识分类</h3>
          <button className="p-1 rounded hover:bg-gray-200 text-gray-400"><Plus className="w-4 h-4" /></button>
        </div>
        <div className="flex-1 overflow-y-auto p-2">
          {mockCategories.map(cat => (
            <div key={cat.id}>
              <button
                onClick={() => {
                  setSelectedCategory(cat.id);
                  setExpandedCats(prev => {
                    const n = new Set(prev);
                    if (n.has(cat.id)) n.delete(cat.id); else n.add(cat.id);
                    return n;
                  });
                }}
                className={`w-full flex items-center gap-2 px-2 py-2 rounded-lg text-sm transition-colors ${
                  selectedCategory === cat.id ? 'bg-white text-violet-700 font-medium shadow-sm' : 'text-gray-600 hover:bg-white'
                }`}
              >
                {cat.children && cat.children.length > 0 ? (
                  expandedCats.has(cat.id) ? <ChevronDown className="w-3 h-3" /> : <ChevronRight className="w-3 h-3" />
                ) : <span className="w-3" />}
                <span>{cat.icon}</span>
                <span className="flex-1 text-left truncate">{cat.name}</span>
                <span className="text-[10px] text-gray-400">{cat.count}</span>
              </button>
              {expandedCats.has(cat.id) && cat.children?.map(sub => (
                <button
                  key={sub.id}
                  onClick={() => setSelectedCategory(sub.id)}
                  className={`w-full flex items-center gap-2 px-2 py-1.5 ml-6 rounded text-xs transition-colors ${
                    selectedCategory === sub.id ? 'text-violet-700 font-medium' : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  <span className="flex-1 text-left">{sub.name}</span>
                  <span className="text-[10px] text-gray-400">{sub.count}</span>
                </button>
              ))}
            </div>
          ))}
        </div>
        <div className="p-2 border-t border-gray-100">
          <button className="w-full py-2 text-xs text-gray-500 hover:text-violet-600 flex items-center justify-center gap-1">
            <Plus className="w-3 h-3" /> 新增分类
          </button>
        </div>
      </div>

      {/* Right: Knowledge Items */}
      <div className="flex-1 flex flex-col min-w-0">
        <div className="p-4 border-b border-gray-200 bg-white">
          <div className="flex items-center gap-3 mb-3">
            <div className="text-2xl">{selectedAgent.icon}</div>
            <div>
              <h2 className="font-bold text-gray-900">{selectedAgent.name} 知识库</h2>
              <p className="text-xs text-gray-500">共 {selectedAgent.knowledgeCount} 条知识</p>
            </div>
          </div>
          <div className="flex gap-2">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                value={searchQuery}
                onChange={e => setSearchQuery(e.target.value)}
                placeholder="搜索知识..."
                className="w-full pl-9 pr-3 py-2 text-sm bg-gray-50 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-violet-300"
              />
            </div>
            <button className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
              <Upload className="w-3.5 h-3.5" />
              导入
            </button>
            <button className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
              <Sparkles className="w-3.5 h-3.5" />
              AI生成
            </button>
            <button className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium bg-violet-600 hover:bg-violet-700 text-white rounded-lg transition-colors">
              <Plus className="w-3.5 h-3.5" />
              新增
            </button>
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-4">
          <div className="space-y-3">
            {mockItems.map(item => (
              <div key={item.id} className="bg-white rounded-xl border border-gray-100 p-4 hover:shadow-md hover:border-violet-200 transition-all cursor-pointer">
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <h3 className="text-sm font-semibold text-gray-900 flex items-center gap-2">
                      <Book className="w-4 h-4 text-violet-500" />
                      {item.title}
                    </h3>
                    <p className="text-xs text-gray-500 mt-1.5 line-clamp-2">{item.preview}</p>
                    <div className="flex items-center gap-2 mt-2">
                      {item.tags.map(tag => (
                        <span key={tag} className="flex items-center gap-0.5 text-[10px] px-1.5 py-0.5 bg-gray-100 text-gray-600 rounded-full">
                          <Tag className="w-2.5 h-2.5" />{tag}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div className="text-right flex-shrink-0 ml-4">
                    <div className="flex items-center gap-1 text-[10px] text-gray-400">
                      <BarChart3 className="w-3 h-3" />
                      使用 {item.useCount} 次
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
