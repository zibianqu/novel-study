import { useState } from 'react';
import { agents as agentData } from '../data/mockData';
import type { Agent } from '../types';
import { Plus, Settings, Database, Zap, ChevronRight, ToggleLeft, ToggleRight } from 'lucide-react';

const layerConfig = {
  decision: { label: '决策层', color: 'bg-amber-100 text-amber-700' },
  strategy: { label: '战略层', color: 'bg-blue-100 text-blue-700' },
  execution: { label: '执行层', color: 'bg-emerald-100 text-emerald-700' },
  quality: { label: '质量层', color: 'bg-rose-100 text-rose-700' },
  auxiliary: { label: '辅助层', color: 'bg-purple-100 text-purple-700' },
};

export default function AgentManager() {
  const [agents] = useState<Agent[]>(agentData);
  const [selected, setSelected] = useState<Agent | null>(null);
  const [showCreate, setShowCreate] = useState(false);

  const coreAgents = agents.filter(a => a.type === 'core');
  const extAgents = agents.filter(a => a.type === 'extension');

  const AgentCard = ({ agent }: { agent: Agent }) => {
    const layer = layerConfig[agent.layer];
    return (
      <div
        onClick={() => setSelected(agent)}
        className={`bg-white rounded-xl border p-4 cursor-pointer transition-all hover:shadow-md ${
          selected?.id === agent.id ? 'border-violet-400 shadow-md ring-2 ring-violet-100' : 'border-gray-100 hover:border-violet-200'
        }`}
      >
        <div className="flex items-start gap-3">
          <div className="text-3xl">{agent.icon}</div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-gray-900 text-sm">{agent.name}</h3>
              {agent.type === 'core' && <span className="text-[10px] px-1.5 py-0.5 bg-violet-100 text-violet-600 rounded-full font-medium">核心</span>}
            </div>
            <span className={`text-[10px] px-1.5 py-0.5 rounded-full ${layer.color} inline-block mt-1`}>{layer.label}</span>
            <p className="text-xs text-gray-500 mt-1.5 line-clamp-2">{agent.description}</p>
            <div className="flex items-center gap-3 mt-2 text-[11px] text-gray-400">
              <span className="flex items-center gap-1"><Database className="w-3 h-3" />{agent.knowledgeCount} 知识</span>
              <span className="flex items-center gap-1"><Zap className="w-3 h-3" />{agent.model}</span>
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="flex h-full">
      {/* Agent List */}
      <div className="flex-1 overflow-y-auto p-6">
        <div className="max-w-4xl mx-auto">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h1 className="text-xl font-bold text-gray-900">Agent 管理</h1>
              <p className="text-sm text-gray-500 mt-1">管理核心Agent和自定义扩展Agent</p>
            </div>
            <button 
              onClick={() => setShowCreate(!showCreate)}
              className="flex items-center gap-2 px-4 py-2.5 bg-violet-600 hover:bg-violet-700 text-white rounded-xl text-sm font-medium shadow-lg shadow-violet-600/20 transition-colors"
            >
              <Plus className="w-4 h-4" />
              创建 Agent
            </button>
          </div>

          {/* Architecture Diagram */}
          <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-5 mb-6 text-white">
            <h3 className="text-sm font-semibold mb-4 text-slate-300">Agent 协作架构</h3>
            <div className="flex flex-col items-center gap-2 text-xs">
              <div className="px-4 py-2 bg-amber-500/20 border border-amber-500/30 rounded-lg text-amber-300">🎬 总导演 (决策层)</div>
              <div className="w-px h-4 bg-slate-600" />
              <div className="flex gap-3">
                <div className="px-3 py-1.5 bg-blue-500/20 border border-blue-500/30 rounded-lg text-blue-300">🌍 天线</div>
                <div className="px-3 py-1.5 bg-blue-500/20 border border-blue-500/30 rounded-lg text-blue-300">🛤️ 地线</div>
                <div className="px-3 py-1.5 bg-blue-500/20 border border-blue-500/30 rounded-lg text-blue-300">⚔️ 剧情线</div>
              </div>
              <div className="text-slate-500 text-[10px]">战略层 — 三线联动</div>
              <div className="w-px h-4 bg-slate-600" />
              <div className="flex gap-3">
                <div className="px-3 py-1.5 bg-emerald-500/20 border border-emerald-500/30 rounded-lg text-emerald-300">🎙️ 旁白</div>
                <div className="px-3 py-1.5 bg-emerald-500/20 border border-emerald-500/30 rounded-lg text-emerald-300">🎭 角色</div>
              </div>
              <div className="text-slate-500 text-[10px]">执行层 — 内容生成</div>
              <div className="w-px h-4 bg-slate-600" />
              <div className="px-3 py-1.5 bg-rose-500/20 border border-rose-500/30 rounded-lg text-rose-300">👁️ 审核 (质量层)</div>
              <div className="w-px h-4 bg-slate-600" />
              <div className="flex gap-2">
                {extAgents.map(a => (
                  <div key={a.id} className="px-2 py-1 bg-purple-500/20 border border-purple-500/30 rounded text-purple-300 text-[10px]">{a.icon} {a.name}</div>
                ))}
                <div className="px-2 py-1 border border-dashed border-slate-600 rounded text-slate-500 text-[10px]">+ 更多</div>
              </div>
              <div className="text-slate-500 text-[10px]">辅助层 — 扩展能力</div>
            </div>
          </div>

          {/* Core Agents */}
          <h2 className="text-sm font-bold text-gray-700 mb-3 flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-violet-500" />
            核心 Agent ({coreAgents.length})
          </h2>
          <div className="grid grid-cols-2 gap-3 mb-8">
            {coreAgents.map(a => <AgentCard key={a.id} agent={a} />)}
          </div>

          {/* Extension Agents */}
          <h2 className="text-sm font-bold text-gray-700 mb-3 flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-purple-500" />
            扩展 Agent ({extAgents.length})
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {extAgents.map(a => <AgentCard key={a.id} agent={a} />)}
          </div>
        </div>
      </div>

      {/* Agent Detail Panel */}
      {selected && (
        <div className="w-80 border-l border-gray-200 bg-white overflow-y-auto flex-shrink-0">
          <div className="p-4 border-b border-gray-100">
            <div className="flex items-center gap-3">
              <span className="text-4xl">{selected.icon}</span>
              <div>
                <h2 className="font-bold text-gray-900">{selected.name}</h2>
                <span className={`text-[10px] px-1.5 py-0.5 rounded-full ${layerConfig[selected.layer].color}`}>{layerConfig[selected.layer].label}</span>
              </div>
            </div>
          </div>
          <div className="p-4 space-y-4">
            <div>
              <h4 className="text-xs font-semibold text-gray-500 mb-1">描述</h4>
              <p className="text-sm text-gray-700">{selected.description}</p>
            </div>
            <div>
              <h4 className="text-xs font-semibold text-gray-500 mb-2">模型配置</h4>
              <div className="space-y-2">
                <div className="flex justify-between text-sm"><span className="text-gray-500">模型</span><span className="text-gray-900 font-medium">{selected.model}</span></div>
                <div className="flex justify-between text-sm"><span className="text-gray-500">Temperature</span><span className="text-gray-900 font-medium">{selected.temperature}</span></div>
                <div className="flex justify-between text-sm"><span className="text-gray-500">Max Tokens</span><span className="text-gray-900 font-medium">{selected.maxTokens}</span></div>
              </div>
            </div>
            <div>
              <h4 className="text-xs font-semibold text-gray-500 mb-2">状态</h4>
              <div className="flex items-center gap-2 text-sm">
                {selected.isActive ? <ToggleRight className="w-6 h-6 text-emerald-500" /> : <ToggleLeft className="w-6 h-6 text-gray-400" />}
                <span className={selected.isActive ? 'text-emerald-600' : 'text-gray-500'}>{selected.isActive ? '已启用' : '已禁用'}</span>
              </div>
            </div>
            <div>
              <h4 className="text-xs font-semibold text-gray-500 mb-2">知识库</h4>
              <div className="text-sm text-gray-700 flex items-center gap-2">
                <Database className="w-4 h-4 text-violet-500" />
                {selected.knowledgeCount} 条知识
                <ChevronRight className="w-4 h-4 text-gray-400 ml-auto" />
              </div>
            </div>
            <div className="pt-2 space-y-2">
              <button className="w-full flex items-center justify-center gap-2 py-2.5 bg-violet-600 hover:bg-violet-700 text-white rounded-lg text-sm font-medium transition-colors">
                <Settings className="w-4 h-4" />
                编辑配置
              </button>
              <button className="w-full flex items-center justify-center gap-2 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg text-sm font-medium transition-colors">
                <Zap className="w-4 h-4" />
                测试 Agent
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
