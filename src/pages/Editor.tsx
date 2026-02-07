import { useState } from 'react';
import { volumes, sampleChapterContent } from '../data/mockData';
import { ChevronDown, ChevronRight, FileText, Save, Wand2, Sparkles, PenTool, MessageSquare, Search, Eye, BarChart3 } from 'lucide-react';

export default function Editor() {
  const [content, setContent] = useState(sampleChapterContent);
  const [activeChapter, setActiveChapter] = useState(8);
  const [expandedVolumes, setExpandedVolumes] = useState<Set<number>>(new Set([1, 2]));
  const [showAiPanel, setShowAiPanel] = useState(true);
  const [aiInput, setAiInput] = useState('');
  const [aiOutput, setAiOutput] = useState('');
  const [isAiGenerating, setIsAiGenerating] = useState(false);

  const wordCount = content.replace(/\s/g, '').length;

  const toggleVolume = (id: number) => {
    setExpandedVolumes(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  };

  const handleAiAction = async (action: string) => {
    setIsAiGenerating(true);
    setAiOutput('');
    await new Promise(r => setTimeout(r, 500));
    
    const responses: Record<string, string> = {
      continue: '　　密林深处，一道微弱的光芒从林远掌心的印记中透出。\n\n　　那是远古传承的力量在回应他的召唤。林远闭上眼睛，感受着体内灵力的流动——和普通修士不同，他的灵力并非循着经脉运行，而是如同潮水般在全身每一个细胞中涌动。\n\n　　"这就是远古修炼法的精髓吗？"他喃喃自语，手中的铁剑突然发出一声清鸣。\n\n　　剑身上那些细小的豁口，在灵力的灌注下竟然开始缓缓愈合。不，不是愈合——是整把剑都在发生质变。铁灰色的剑身渐渐泛起一层淡淡的青光。',
      polish: '　　夜幕低垂，银色月华倾泻而下，将天剑宗演武场渲染成一幅水墨丹青。青石地面在月光的抚摸下泛着柔和而清冷的光泽，宛如一面沉寂千年的古镜。\n\n　　林远伫立于演武场一隅，五指紧扣着一柄斑驳的铁剑。剑身遍布岁月的齿痕——那是杂役堂角落里落满灰尘的练功用剑。堂堂天剑宗弟子，却连一把堪用的兵器都不曾拥有。这便是外门末等弟子的宿命。',
      dialogue: '　　"林远。"王昊的声音从身后响起，带着居高临下的傲慢。\n\n　　林远没有转身。\n\n　　"我在跟你说话！"王昊加重了语气，两步走到他面前，"一个废脉的人，整天泡在演武场，你是想感动天道让你重新修炼吗？"\n\n　　他身后的两个跟班发出配合的笑声。\n\n　　林远终于抬起目光，平静得像一潭死水："王师兄，大比报名截止是明日。"\n\n　　王昊笑容微凝："你什么意思？"\n\n　　"字面意思。"林远收好铁剑，侧身走过他，"演武场不是你家的，我想来就来。"',
    };

    let result = responses[action] || '生成中...请选择具体的AI功能。';
    setAiOutput('');
    for (let i = 0; i < result.length; i++) {
      await new Promise(r => setTimeout(r, 15));
      setAiOutput(result.slice(0, i + 1));
    }
    setIsAiGenerating(false);
  };

  return (
    <div className="flex h-full">
      {/* Left: Chapter Tree */}
      <div className="w-56 border-r border-gray-200 bg-white flex flex-col flex-shrink-0">
        <div className="p-3 border-b border-gray-100">
          <h3 className="text-sm font-semibold text-gray-700">📚 章节目录</h3>
        </div>
        <div className="flex-1 overflow-y-auto p-2">
          {volumes.map(vol => (
            <div key={vol.id} className="mb-1">
              <button onClick={() => toggleVolume(vol.id)} className="w-full flex items-center gap-1.5 px-2 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50 rounded">
                {expandedVolumes.has(vol.id) ? <ChevronDown className="w-3.5 h-3.5" /> : <ChevronRight className="w-3.5 h-3.5" />}
                {vol.title}
              </button>
              {expandedVolumes.has(vol.id) && vol.chapters.map(ch => (
                <button
                  key={ch.id}
                  onClick={() => setActiveChapter(ch.id)}
                  className={`w-full flex items-center gap-1.5 px-2 py-1.5 ml-4 text-xs rounded transition-colors ${
                    activeChapter === ch.id ? 'bg-violet-50 text-violet-700 font-medium' : 'text-gray-600 hover:bg-gray-50'
                  }`}
                >
                  <FileText className="w-3 h-3 flex-shrink-0" />
                  <span className="truncate">{ch.title}</span>
                  {ch.status === 'final' && <span className="ml-auto text-emerald-500 text-[10px]">✓</span>}
                </button>
              ))}
            </div>
          ))}
        </div>
        {/* Characters quick list */}
        <div className="border-t border-gray-100 p-3">
          <h4 className="text-xs font-semibold text-gray-500 mb-2">👥 本章角色</h4>
          <div className="space-y-1">
            {['林远 (主角)', '王昊 (配角)', '长老 (配角)'].map(c => (
              <div key={c} className="text-xs text-gray-600 px-2 py-1 rounded hover:bg-gray-50 cursor-pointer">{c}</div>
            ))}
          </div>
        </div>
      </div>

      {/* Center: Editor */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Tab bar */}
        <div className="flex items-center gap-1 px-3 py-2 bg-gray-50 border-b border-gray-200">
          <div className="flex items-center gap-1.5 px-3 py-1.5 bg-white border border-gray-200 rounded-lg text-sm">
            <FileText className="w-3.5 h-3.5 text-violet-500" />
            <span className="font-medium text-gray-700">第八章 宗门大比</span>
            <span className="text-xs text-gray-400">•</span>
          </div>
          <div className="flex-1" />
          <button className="flex items-center gap-1.5 px-3 py-1.5 text-xs text-gray-500 hover:text-violet-600 hover:bg-violet-50 rounded-lg transition-colors">
            <Save className="w-3.5 h-3.5" />
            保存
          </button>
          <button className="flex items-center gap-1.5 px-3 py-1.5 text-xs text-gray-500 hover:text-violet-600 hover:bg-violet-50 rounded-lg transition-colors">
            <Eye className="w-3.5 h-3.5" />
            专注模式
          </button>
          <button 
            onClick={() => setShowAiPanel(!showAiPanel)}
            className={`flex items-center gap-1.5 px-3 py-1.5 text-xs rounded-lg transition-colors ${showAiPanel ? 'bg-violet-100 text-violet-700' : 'text-gray-500 hover:bg-gray-100'}`}
          >
            <Wand2 className="w-3.5 h-3.5" />
            AI助手
          </button>
        </div>

        {/* Editor area */}
        <div className="flex-1 overflow-y-auto bg-[#FDF6E3]">
          <div className="max-w-3xl mx-auto py-8 px-12">
            <textarea
              value={content}
              onChange={e => setContent(e.target.value)}
              className="w-full min-h-[600px] bg-transparent border-none outline-none resize-none text-gray-800 leading-[2] text-base font-serif"
              style={{ fontFamily: '"Noto Serif SC", "Source Han Serif CN", Georgia, serif' }}
            />
          </div>
        </div>

        {/* Status bar */}
        <div className="flex items-center gap-4 px-4 py-1.5 bg-gray-100 border-t border-gray-200 text-xs text-gray-500">
          <span className="flex items-center gap-1"><BarChart3 className="w-3 h-3" />本章: {wordCount} 字</span>
          <span>全书: 156,800 字</span>
          <span>第8章 / 共9章</span>
          <span className="flex items-center gap-1 text-emerald-600">
            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
            已自动保存
          </span>
        </div>
      </div>

      {/* Right: AI Panel */}
      {showAiPanel && (
        <div className="w-80 border-l border-gray-200 bg-white flex flex-col flex-shrink-0">
          <div className="p-3 border-b border-gray-100">
            <h3 className="text-sm font-semibold text-gray-700 flex items-center gap-2">
              <Wand2 className="w-4 h-4 text-violet-500" />
              AI 创作助手
            </h3>
          </div>

          {/* AI Actions */}
          <div className="p-3 border-b border-gray-100">
            <div className="grid grid-cols-3 gap-2">
              {[
                { label: '续写', icon: Sparkles, action: 'continue' },
                { label: '润色', icon: PenTool, action: 'polish' },
                { label: '对话', icon: MessageSquare, action: 'dialogue' },
                { label: '改写', icon: Wand2, action: 'rewrite' },
                { label: '建议', icon: Search, action: 'suggest' },
                { label: '检查', icon: Eye, action: 'check' },
              ].map(btn => {
                const Icon = btn.icon;
                return (
                  <button
                    key={btn.label}
                    onClick={() => handleAiAction(btn.action)}
                    disabled={isAiGenerating}
                    className="flex flex-col items-center gap-1 p-2 rounded-lg text-xs font-medium text-gray-600 hover:bg-violet-50 hover:text-violet-600 border border-gray-100 hover:border-violet-200 transition-colors disabled:opacity-50"
                  >
                    <Icon className="w-4 h-4" />
                    {btn.label}
                  </button>
                );
              })}
            </div>
          </div>

          {/* AI Instruction */}
          <div className="p-3 border-b border-gray-100">
            <textarea
              value={aiInput}
              onChange={e => setAiInput(e.target.value)}
              placeholder="输入创作指令... 如：让主角在客栈遇到神秘老者"
              className="w-full h-20 px-3 py-2 text-xs bg-gray-50 border border-gray-200 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-violet-300"
            />
          </div>

          {/* AI Output */}
          <div className="flex-1 overflow-y-auto p-3">
            {aiOutput ? (
              <>
                <div className="text-xs text-gray-400 mb-2 flex items-center gap-1">
                  {isAiGenerating && <span className="inline-block w-1.5 h-1.5 rounded-full bg-violet-500 animate-pulse" />}
                  {isAiGenerating ? 'AI 生成中...' : 'AI 输出：'}
                </div>
                <div className="text-sm text-gray-700 leading-relaxed whitespace-pre-wrap font-serif" style={{ fontFamily: '"Noto Serif SC", Georgia, serif' }}>
                  {aiOutput}
                </div>
                {!isAiGenerating && (
                  <div className="flex gap-2 mt-4">
                    <button onClick={() => setContent(prev => prev + '\n\n' + aiOutput)} className="flex-1 py-2 text-xs font-medium bg-violet-600 hover:bg-violet-700 text-white rounded-lg transition-colors">✅ 采纳</button>
                    <button onClick={() => handleAiAction('continue')} className="flex-1 py-2 text-xs font-medium bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">🔄 重新生成</button>
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-12 text-gray-400">
                <Wand2 className="w-8 h-8 mx-auto mb-3 opacity-30" />
                <p className="text-xs">选择AI功能或输入指令开始创作</p>
                <p className="text-xs mt-1">支持右键菜单快捷操作</p>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
