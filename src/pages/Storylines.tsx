import { storylines } from '../data/mockData';
import { Globe, Route, Swords, ChevronRight, Edit3, AlertCircle } from 'lucide-react';

const lineConfig = {
  skyline: { label: '天线', sublabel: '世界命运', icon: Globe, color: 'from-amber-500 to-orange-600', bg: 'bg-amber-50', border: 'border-amber-200', text: 'text-amber-700' },
  groundline: { label: '地线', sublabel: '主角路径', icon: Route, color: 'from-blue-500 to-cyan-600', bg: 'bg-blue-50', border: 'border-blue-200', text: 'text-blue-700' },
  plotline: { label: '剧情线', sublabel: '升级节奏', icon: Swords, color: 'from-rose-500 to-pink-600', bg: 'bg-rose-50', border: 'border-rose-200', text: 'text-rose-700' },
};

const statusColors = {
  planned: 'bg-gray-200 text-gray-600',
  active: 'bg-emerald-100 text-emerald-700',
  completed: 'bg-blue-100 text-blue-700',
};

export default function Storylines() {
  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="mb-6">
        <h1 className="text-xl font-bold text-gray-900">三线管理</h1>
        <p className="text-sm text-gray-500 mt-1">管理天线（世界命运）、地线（主角路径）、剧情线（升级节奏）</p>
      </div>

      {/* Three Lines Architecture */}
      <div className="bg-gradient-to-br from-slate-800 to-slate-900 rounded-xl p-6 mb-6">
        <div className="flex items-center justify-center gap-4 text-white text-sm">
          <div className="text-center">
            <div className="px-4 py-2 bg-amber-500/20 border border-amber-500/30 rounded-lg mb-1">🌍 天线（世界命运）</div>
            <div className="text-[10px] text-slate-400">Agent 4 掌控</div>
          </div>
          <div className="flex flex-col items-center gap-1 text-[10px] text-slate-400">
            <span>↓ 影响/倒逼</span>
          </div>
          <div className="text-center">
            <div className="px-4 py-2 bg-blue-500/20 border border-blue-500/30 rounded-lg mb-1">🛤️ 地线（主角路径）</div>
            <div className="text-[10px] text-slate-400">Agent 5 掌控</div>
          </div>
          <div className="flex flex-col items-center gap-1 text-[10px] text-slate-400">
            <span>↑ 驱动/实现</span>
          </div>
          <div className="text-center">
            <div className="px-4 py-2 bg-rose-500/20 border border-rose-500/30 rounded-lg mb-1">⚔️ 剧情线（升级节奏）</div>
            <div className="text-[10px] text-slate-400">Agent 6 掌控</div>
          </div>
        </div>
      </div>

      {/* Storylines */}
      <div className="space-y-6">
        {storylines.map(line => {
          const config = lineConfig[line.type];
          const Icon = config.icon;
          return (
            <div key={line.id} className={`rounded-xl border ${config.border} ${config.bg} overflow-hidden`}>
              {/* Header */}
              <div className="p-4 flex items-center gap-3">
                <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${config.color} flex items-center justify-center`}>
                  <Icon className="w-5 h-5 text-white" />
                </div>
                <div className="flex-1">
                  <h2 className={`font-bold ${config.text}`}>{config.label} — {config.sublabel}</h2>
                  <p className="text-sm text-gray-600">{line.content}</p>
                </div>
                <button className="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-gray-600 bg-white hover:bg-gray-50 rounded-lg border border-gray-200 transition-colors">
                  <Edit3 className="w-3.5 h-3.5" />
                  编辑
                </button>
              </div>

              {/* Timeline */}
              <div className="px-4 pb-4">
                <div className="bg-white rounded-lg border border-gray-100 overflow-hidden">
                  {line.items.map((item, idx) => (
                    <div key={item.id} className={`flex items-center gap-4 p-3 ${idx > 0 ? 'border-t border-gray-50' : ''} hover:bg-gray-50 transition-colors`}>
                      {/* Timeline dot */}
                      <div className="flex flex-col items-center gap-1 flex-shrink-0">
                        <div className={`w-3 h-3 rounded-full ${item.status === 'completed' ? 'bg-emerald-500' : item.status === 'active' ? 'bg-violet-500 ring-4 ring-violet-100' : 'bg-gray-300'}`} />
                        {idx < line.items.length - 1 && <div className="w-px h-6 bg-gray-200" />}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <h4 className="text-sm font-semibold text-gray-900">{item.title}</h4>
                          <span className={`text-[10px] px-1.5 py-0.5 rounded-full ${statusColors[item.status]}`}>
                            {item.status === 'completed' ? '已完成' : item.status === 'active' ? '进行中' : '待开始'}
                          </span>
                        </div>
                        <p className="text-xs text-gray-500 mt-0.5">{item.content}</p>
                      </div>
                      <div className="text-xs text-gray-400 flex-shrink-0">Ch.{item.chapterRange}</div>
                      <ChevronRight className="w-4 h-4 text-gray-300 flex-shrink-0" />
                    </div>
                  ))}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Convergence Points */}
      <div className="mt-6 bg-white rounded-xl border border-gray-100 p-5">
        <h3 className="font-bold text-gray-900 mb-3 flex items-center gap-2">
          <AlertCircle className="w-4 h-4 text-violet-500" />
          三线交汇点
        </h3>
        <div className="space-y-3">
          {[
            { name: '宗门大比', sky: '各势力暗中博弈的缩影', ground: '主角首次展露实力', plot: '第一个大高潮', chapter: 8 },
            { name: '秘境探险', sky: '上古秘密即将揭晓', ground: '主角发现身世线索', plot: '探险升级弧', chapter: 16 },
          ].map(conv => (
            <div key={conv.name} className="flex items-center gap-4 p-3 bg-gray-50 rounded-lg">
              <div className="w-10 h-10 rounded-lg bg-violet-100 flex items-center justify-center text-violet-600 font-bold text-xs">Ch.{conv.chapter}</div>
              <div className="flex-1">
                <h4 className="text-sm font-semibold text-gray-900">{conv.name}</h4>
                <div className="flex gap-3 mt-1 text-[11px]">
                  <span className="text-amber-600">🌍 {conv.sky}</span>
                  <span className="text-blue-600">🛤️ {conv.ground}</span>
                  <span className="text-rose-600">⚔️ {conv.plot}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
