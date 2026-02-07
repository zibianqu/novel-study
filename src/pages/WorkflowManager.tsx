import { useState } from 'react';
import { workflows } from '../data/mockData';
import { Plus, Play, Settings, ChevronRight, Lock } from 'lucide-react';

const categoryColors: Record<string, string> = {
  '长篇创作': 'bg-violet-100 text-violet-700',
  '短篇创作': 'bg-blue-100 text-blue-700',
  '文案': 'bg-amber-100 text-amber-700',
  '管理': 'bg-emerald-100 text-emerald-700',
};

export default function WorkflowManager() {
  const [selectedWf, setSelectedWf] = useState(workflows[1]);

  // Simulated workflow nodes for "章节创作"
  const wfNodes = [
    { id: 'input', name: '用户输入', icon: '📥', type: 'input' },
    { id: 'director', name: '总导演分析', icon: '🎬', type: 'agent' },
    { id: 'plotline', name: '剧情线安排', icon: '⚔️', type: 'agent' },
    { id: 'skyline', name: '天线信息', icon: '🌍', type: 'agent' },
    { id: 'groundline', name: '地线信息', icon: '🛤️', type: 'agent' },
    { id: 'rag', name: 'RAG检索', icon: '📚', type: 'rag' },
    { id: 'neo4j', name: '图谱查询', icon: '🕸️', type: 'neo4j' },
    { id: 'narrator', name: '旁白叙述', icon: '🎙️', type: 'agent' },
    { id: 'character', name: '角色对话', icon: '🎭', type: 'agent' },
    { id: 'reviewer', name: '审核导演', icon: '👁️', type: 'agent' },
    { id: 'condition', name: '审核通过?', icon: '🔀', type: 'condition' },
    { id: 'confirm', name: '用户确认', icon: '👤', type: 'confirm' },
    { id: 'storage', name: '入库', icon: '💾', type: 'storage' },
  ];

  const nodeTypeColors: Record<string, string> = {
    input: 'bg-gray-100 border-gray-300 text-gray-700',
    agent: 'bg-violet-50 border-violet-300 text-violet-700',
    rag: 'bg-emerald-50 border-emerald-300 text-emerald-700',
    neo4j: 'bg-cyan-50 border-cyan-300 text-cyan-700',
    condition: 'bg-amber-50 border-amber-300 text-amber-700',
    confirm: 'bg-blue-50 border-blue-300 text-blue-700',
    storage: 'bg-rose-50 border-rose-300 text-rose-700',
  };

  return (
    <div className="flex h-full">
      {/* Workflow List */}
      <div className="w-72 border-r border-gray-200 bg-white flex flex-col flex-shrink-0">
        <div className="p-4 border-b border-gray-100">
          <div className="flex items-center justify-between mb-3">
            <h2 className="font-bold text-gray-900">工作流</h2>
            <button className="p-1.5 rounded-lg bg-violet-600 hover:bg-violet-700 text-white transition-colors">
              <Plus className="w-4 h-4" />
            </button>
          </div>
          <p className="text-xs text-gray-500">预置 {workflows.filter(w => w.type === 'system').length} 套 + 自定义</p>
        </div>
        <div className="flex-1 overflow-y-auto p-2 space-y-1">
          {workflows.map(wf => (
            <button
              key={wf.id}
              onClick={() => setSelectedWf(wf)}
              className={`w-full flex items-center gap-3 p-3 rounded-lg text-left transition-all ${
                selectedWf.id === wf.id ? 'bg-violet-50 border border-violet-200' : 'hover:bg-gray-50 border border-transparent'
              }`}
            >
              <span className="text-xl">{wf.icon}</span>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-1.5">
                  <h3 className="text-sm font-medium text-gray-900 truncate">{wf.name}</h3>
                  {wf.type === 'system' && <Lock className="w-3 h-3 text-gray-400 flex-shrink-0" />}
                </div>
                <span className={`text-[10px] px-1.5 py-0.5 rounded-full ${categoryColors[wf.category] || 'bg-gray-100 text-gray-600'}`}>{wf.category}</span>
              </div>
              <ChevronRight className="w-4 h-4 text-gray-300 flex-shrink-0" />
            </button>
          ))}
        </div>
      </div>

      {/* Workflow Detail */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Header */}
        <div className="p-4 border-b border-gray-200 bg-white flex items-center gap-4">
          <span className="text-2xl">{selectedWf.icon}</span>
          <div className="flex-1">
            <h2 className="font-bold text-gray-900">{selectedWf.name}</h2>
            <p className="text-xs text-gray-500">{selectedWf.description}</p>
          </div>
          <div className="flex gap-2">
            <button className="flex items-center gap-1.5 px-3 py-2 text-xs font-medium bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
              <Settings className="w-3.5 h-3.5" />
              编辑
            </button>
            <button className="flex items-center gap-1.5 px-4 py-2 text-xs font-medium bg-violet-600 hover:bg-violet-700 text-white rounded-lg shadow-lg shadow-violet-600/20 transition-colors">
              <Play className="w-3.5 h-3.5" />
              执行
            </button>
          </div>
        </div>

        {/* Workflow Canvas */}
        <div className="flex-1 overflow-auto bg-gray-50 p-8">
          <div className="max-w-4xl mx-auto">
            {/* Flow Visualization */}
            <div className="flex flex-col items-center gap-2">
              {wfNodes.map((node, idx) => (
                <div key={node.id} className="flex flex-col items-center">
                  <div className={`flex items-center gap-2 px-4 py-2.5 rounded-xl border-2 min-w-[180px] ${nodeTypeColors[node.type]} shadow-sm bg-white`}>
                    <span className="text-lg">{node.icon}</span>
                    <div className="flex-1">
                      <div className="text-sm font-medium">{node.name}</div>
                      <div className="text-[10px] opacity-70">{node.type === 'agent' ? 'Agent节点' : node.type === 'rag' ? 'RAG检索' : node.type === 'condition' ? '条件判断' : node.type}</div>
                    </div>
                  </div>
                  {idx < wfNodes.length - 1 && (
                    <div className="flex flex-col items-center my-1">
                      <div className="w-px h-4 bg-gray-300" />
                      {node.id === 'condition' ? (
                        <div className="flex gap-8 items-start">
                          <div className="text-center">
                            <div className="text-[10px] text-emerald-500 font-medium mb-1">✅ 通过</div>
                            <div className="w-px h-3 bg-emerald-300 mx-auto" />
                          </div>
                          <div className="text-center">
                            <div className="text-[10px] text-rose-500 font-medium mb-1">❌ 不通过 → 回到旁白</div>
                          </div>
                        </div>
                      ) : (
                        <svg width="8" height="8" viewBox="0 0 8 8"><polygon points="4,8 0,0 8,0" fill="#CBD5E1"/></svg>
                      )}
                    </div>
                  )}
                </div>
              ))}
            </div>

            {/* Node Types Legend */}
            <div className="mt-8 p-4 bg-white rounded-xl border border-gray-200">
              <h3 className="text-xs font-semibold text-gray-500 mb-3">节点类型说明</h3>
              <div className="grid grid-cols-4 gap-2">
                {[
                  { icon: '📥', label: '输入节点', desc: '接收用户数据' },
                  { icon: '🤖', label: 'Agent节点', desc: '调用AI Agent' },
                  { icon: '📚', label: 'RAG检索', desc: '向量检索知识' },
                  { icon: '🕸️', label: '图谱查询', desc: 'Neo4j查询' },
                  { icon: '🔀', label: '条件判断', desc: '分支路由' },
                  { icon: '👤', label: '用户确认', desc: '等待用户操作' },
                  { icon: '💾', label: '数据存储', desc: '写入数据库' },
                  { icon: '🔄', label: '循环节点', desc: '重试/迭代' },
                ].map(nt => (
                  <div key={nt.label} className="flex items-center gap-2 p-2 rounded-lg bg-gray-50">
                    <span className="text-sm">{nt.icon}</span>
                    <div>
                      <div className="text-xs font-medium text-gray-700">{nt.label}</div>
                      <div className="text-[10px] text-gray-400">{nt.desc}</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
