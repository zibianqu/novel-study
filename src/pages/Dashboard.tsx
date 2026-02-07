import { projects } from '../data/mockData';
import { Plus, BookOpen, FileText, PenTool, TrendingUp, Clock, BarChart3 } from 'lucide-react';

const typeConfig = {
  novel_long: { label: '长篇小说', color: 'bg-violet-100 text-violet-700', icon: BookOpen },
  novel_short: { label: '短篇小说', color: 'bg-blue-100 text-blue-700', icon: FileText },
  copywriting: { label: '文案', color: 'bg-amber-100 text-amber-700', icon: PenTool },
};

const statusConfig = {
  draft: { label: '草稿', color: 'bg-gray-100 text-gray-600' },
  writing: { label: '创作中', color: 'bg-emerald-100 text-emerald-700' },
  completed: { label: '已完结', color: 'bg-blue-100 text-blue-700' },
};

export default function Dashboard({ onNavigate }: { onNavigate: (page: string) => void }) {
  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* Stats */}
      <div className="grid grid-cols-4 gap-4 mb-8">
        {[
          { label: '总项目数', value: '4', icon: BookOpen, color: 'from-violet-500 to-indigo-600' },
          { label: '总字数', value: '267,000', icon: BarChart3, color: 'from-emerald-500 to-teal-600' },
          { label: '本周字数', value: '12,400', icon: TrendingUp, color: 'from-amber-500 to-orange-600' },
          { label: 'AI调用次数', value: '1,247', icon: PenTool, color: 'from-pink-500 to-rose-600' },
        ].map(stat => {
          const Icon = stat.icon;
          return (
            <div key={stat.label} className="bg-white rounded-xl p-5 border border-gray-100 shadow-sm">
              <div className="flex items-center justify-between mb-3">
                <div className={`w-10 h-10 rounded-xl bg-gradient-to-br ${stat.color} flex items-center justify-center`}>
                  <Icon className="w-5 h-5 text-white" />
                </div>
              </div>
              <p className="text-2xl font-bold text-gray-900">{stat.value}</p>
              <p className="text-sm text-gray-500 mt-1">{stat.label}</p>
            </div>
          );
        })}
      </div>

      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-xl font-bold text-gray-900">我的项目</h1>
          <p className="text-sm text-gray-500 mt-1">管理你的所有小说创作项目</p>
        </div>
        <button className="flex items-center gap-2 px-4 py-2.5 bg-violet-600 hover:bg-violet-700 text-white rounded-xl text-sm font-medium shadow-lg shadow-violet-600/20 transition-colors">
          <Plus className="w-4 h-4" />
          新建项目
        </button>
      </div>

      {/* Projects Grid */}
      <div className="grid grid-cols-2 gap-4">
        {projects.map(project => {
          const typeInfo = typeConfig[project.type];
          const statusInfo = statusConfig[project.status];
          const TypeIcon = typeInfo.icon;
          return (
            <div 
              key={project.id} 
              onClick={() => onNavigate('director')}
              className="bg-white rounded-xl border border-gray-100 p-5 hover:shadow-lg hover:border-violet-200 transition-all cursor-pointer group"
            >
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-violet-100 to-indigo-100 flex items-center justify-center">
                    <TypeIcon className="w-6 h-6 text-violet-600" />
                  </div>
                  <div>
                    <h3 className="font-bold text-gray-900 group-hover:text-violet-600 transition-colors">{project.title}</h3>
                    <div className="flex items-center gap-2 mt-1">
                      <span className={`text-xs px-2 py-0.5 rounded-full ${typeInfo.color}`}>{typeInfo.label}</span>
                      <span className={`text-xs px-2 py-0.5 rounded-full ${statusInfo.color}`}>{statusInfo.label}</span>
                      {project.genre && <span className="text-xs text-gray-400">{project.genre}</span>}
                    </div>
                  </div>
                </div>
              </div>
              <p className="text-sm text-gray-500 mb-4 line-clamp-2">{project.description}</p>
              <div className="flex items-center gap-4 text-xs text-gray-400">
                <span className="flex items-center gap-1"><BarChart3 className="w-3.5 h-3.5" />{project.wordCount.toLocaleString()} 字</span>
                <span className="flex items-center gap-1"><FileText className="w-3.5 h-3.5" />{project.chapterCount} 章</span>
                <span className="flex items-center gap-1 ml-auto"><Clock className="w-3.5 h-3.5" />{project.updatedAt}</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
