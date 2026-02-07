import { useState, useRef, useEffect } from 'react';
import type { ChatMessage } from '../types';
import { Send, Sparkles, GitBranch, RefreshCw, PenTool, Loader2, Bot } from 'lucide-react';

const initialMessages: ChatMessage[] = [
  {
    id: '1', role: 'director', content: '你好！我是你的小说创作总导演 🎬\n\n我会协调旗下所有创作Agent为你打造完整的小说：\n\n• 🌍 天线掌控者 - 把控世界格局\n• 🛤️ 地线掌控者 - 规划主角路径\n• ⚔️ 剧情线掌控者 - 设计精彩情节\n• 🎙️ 旁白叙述者 - 撰写优美文笔\n• 🎭 角色扮演者 - 演绎鲜活对话\n• 👁️ 审核导演 - 保障内容质量\n\n当前项目「九天仙途」已完成8章，请告诉我接下来要做什么？',
    timestamp: '10:00', agent: '总导演'
  }
];

export default function DirectorChat() {
  const [messages, setMessages] = useState<ChatMessage[]>(initialMessages);
  const [input, setInput] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const simulateResponse = async (userMsg: string) => {
    setIsGenerating(true);
    
    // Thinking
    const thinkingId = Date.now().toString();
    setMessages(prev => [...prev, { id: thinkingId, role: 'director', content: '正在分析你的指令...', timestamp: new Date().toLocaleTimeString().slice(0, 5), agent: '总导演', status: 'thinking' }]);
    await new Promise(r => setTimeout(r, 1000));

    // Dispatching
    setMessages(prev => prev.map(m => m.id === thinkingId ? { ...m, content: '📡 正在调度Agent团队...', status: 'dispatching' } : m));
    await new Promise(r => setTimeout(r, 800));

    let response = '';
    if (userMsg.includes('第九章') || userMsg.includes('写')) {
      response = `好的，我来安排第九章「一鸣惊人」的创作！\n\n📋 **调度进度：**\n✅ 🌍 天线掌控者：秘境开启的前兆已设定\n✅ 🛤️ 地线掌控者：主角当前状态 — 实力隐藏，内心坚定\n✅ ⚔️ 剧情线掌控者：本章核心 — 宗门大比高潮\n\n📊 **本章规划：**\n- **天线**：灵气异变影响大比，天道考验降临\n- **地线**：林远在大比中展露传承力量，震惊全场\n- **剧情线**：危机→王昊挑衅→林远隐忍→最终爆发→晋升\n\n⚠️ **伏笔提醒：**\n- 第3章埋下的「远古传承印记」将在本章首次展现\n- 需要为第10章的「秘境开启」做铺垫\n\n准备开始正文创作吗？还是你想先调整一下本章规划？`;
    } else if (userMsg.includes('推演')) {
      response = `好的，我来组织全体Agent进行推演！\n\n🔮 **推演报告（第9-13章）：**\n\n🌍 **天线走向：**\n• 第9章: 灵气异变影响宗门大比\n• 第10章: 上古秘境开启信号出现\n• 第11章: 三大宗门争夺秘境入场资格\n• 第12章: 秘境内发现远古遗迹\n• 第13章: 遗迹中封印出现裂缝\n\n🛤️ **地线走向：**\n• 第9章: 大比中一鸣惊人，获内门资格\n• 第10章: 被选为秘境探索队成员\n• 第11章: 秘境外遭遇暗杀，识破阴谋\n• 第12章: 进入秘境，传承印记共鸣\n• 第13章: 发现远古传承的真正来源\n\n⚔️ **剧情线：**\n• 第9章: [爆发-燃] 大比高潮 🔥\n• 第10章: [过渡-期待] 新篇章开启\n• 第11章: [紧张-危机] 暗流涌动\n• 第12章: [探险-神秘] 未知世界\n• 第13章: [震撼-转折] 真相浮现\n\n你觉得这个走向如何？需要调整吗？`;
    } else {
      response = `收到！让我来分析一下你的需求...\n\n我可以帮你：\n1. 📝 **继续创作** - 写下一章内容\n2. 🔮 **推演未来** - 推演接下来几章的走向\n3. 🔧 **调整三线** - 修改天线/地线/剧情线\n4. 👤 **创建角色** - 设计新的角色\n5. 🔍 **一致性检查** - 检查已写内容的一致性\n\n请告诉我你想做什么？你也可以直接说出具体的创作指令，比如"写第九章"或"推演后5章"。`;
    }

    setMessages(prev => prev.map(m => m.id === thinkingId ? { ...m, content: response, status: 'complete' } : m));
    setIsGenerating(false);
  };

  const handleSend = () => {
    if (!input.trim() || isGenerating) return;
    const userMessage: ChatMessage = {
      id: Date.now().toString(),
      role: 'user',
      content: input,
      timestamp: new Date().toLocaleTimeString().slice(0, 5)
    };
    setMessages(prev => [...prev, userMessage]);
    const msg = input;
    setInput('');
    simulateResponse(msg);
  };

  const quickActions = [
    { label: '继续创作', icon: PenTool, action: '写第九章' },
    { label: '推演后5章', icon: Sparkles, action: '推演接下来5章的走向' },
    { label: '查看三线', icon: GitBranch, action: '展示当前三线状态' },
    { label: '一致性检查', icon: RefreshCw, action: '检查全书一致性' },
  ];

  return (
    <div className="flex flex-col h-full bg-gray-50">
      {/* Chat Area */}
      <div className="flex-1 overflow-y-auto px-4 py-6">
        <div className="max-w-3xl mx-auto space-y-4">
          {messages.map(msg => (
            <div key={msg.id} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
              <div className={`max-w-[85%] ${msg.role === 'user' ? '' : ''}`}>
                {msg.role === 'director' && (
                  <div className="flex items-center gap-2 mb-1.5">
                    <div className="w-7 h-7 rounded-lg bg-gradient-to-br from-violet-500 to-indigo-600 flex items-center justify-center">
                      {msg.status === 'thinking' || msg.status === 'dispatching' ? (
                        <Loader2 className="w-4 h-4 text-white animate-spin" />
                      ) : (
                        <Bot className="w-4 h-4 text-white" />
                      )}
                    </div>
                    <span className="text-xs font-medium text-gray-500">🎬 {msg.agent}</span>
                    <span className="text-xs text-gray-400">{msg.timestamp}</span>
                  </div>
                )}
                <div className={`rounded-2xl px-4 py-3 ${
                  msg.role === 'user' 
                    ? 'bg-violet-600 text-white rounded-br-md' 
                    : msg.status === 'thinking' || msg.status === 'dispatching'
                      ? 'bg-white border border-gray-200 text-gray-500 italic'
                      : 'bg-white border border-gray-100 shadow-sm text-gray-800 rounded-bl-md'
                }`}>
                  <div className="text-sm leading-relaxed whitespace-pre-wrap">
                    {msg.content.split('\n').map((line, i) => {
                      if (line.startsWith('**') && line.endsWith('**')) {
                        return <p key={i} className="font-bold mt-2">{line.replace(/\*\*/g, '')}</p>;
                      }
                      const boldParts = line.split(/(\*\*.*?\*\*)/);
                      return (
                        <p key={i} className={line === '' ? 'h-2' : ''}>
                          {boldParts.map((part, j) => 
                            part.startsWith('**') && part.endsWith('**') 
                              ? <strong key={j} className="font-semibold">{part.slice(2, -2)}</strong>
                              : <span key={j}>{part}</span>
                          )}
                        </p>
                      );
                    })}
                  </div>
                </div>
                {msg.role === 'user' && (
                  <div className="text-right mt-1">
                    <span className="text-xs text-gray-400">{msg.timestamp}</span>
                  </div>
                )}
              </div>
            </div>
          ))}
          <div ref={messagesEndRef} />
        </div>
      </div>

      {/* Quick Actions */}
      <div className="px-4 py-2 border-t border-gray-100 bg-white">
        <div className="max-w-3xl mx-auto flex gap-2">
          {quickActions.map(action => {
            const Icon = action.icon;
            return (
              <button
                key={action.label}
                onClick={() => { setInput(action.action); }}
                disabled={isGenerating}
                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-gray-600 bg-gray-50 hover:bg-violet-50 hover:text-violet-600 rounded-lg border border-gray-200 hover:border-violet-200 transition-colors disabled:opacity-50"
              >
                <Icon className="w-3.5 h-3.5" />
                {action.label}
              </button>
            );
          })}
        </div>
      </div>

      {/* Input */}
      <div className="px-4 py-3 bg-white border-t border-gray-200">
        <div className="max-w-3xl mx-auto flex gap-3">
          <div className="flex-1 relative">
            <input
              value={input}
              onChange={e => setInput(e.target.value)}
              onKeyDown={e => e.key === 'Enter' && !e.shiftKey && handleSend()}
              placeholder="与总导演对话... 例如：写第九章、推演后5章、创建一个反派角色"
              disabled={isGenerating}
              className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-violet-300 focus:border-violet-400 disabled:opacity-50 pr-12"
            />
            <button 
              onClick={handleSend} 
              disabled={!input.trim() || isGenerating}
              className="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 flex items-center justify-center rounded-lg bg-violet-600 hover:bg-violet-700 text-white disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors"
            >
              {isGenerating ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
