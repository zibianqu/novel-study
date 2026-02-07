export interface Agent {
  id: number;
  key: string;
  name: string;
  icon: string;
  type: 'core' | 'extension';
  layer: 'decision' | 'strategy' | 'execution' | 'quality' | 'auxiliary';
  description: string;
  model: string;
  temperature: number;
  maxTokens: number;
  systemPrompt: string;
  isActive: boolean;
  knowledgeCount: number;
}

export interface Project {
  id: number;
  title: string;
  type: 'novel_long' | 'novel_short' | 'copywriting';
  genre: string;
  description: string;
  status: 'draft' | 'writing' | 'completed';
  wordCount: number;
  chapterCount: number;
  updatedAt: string;
}

export interface Chapter {
  id: number;
  title: string;
  wordCount: number;
  status: 'draft' | 'final';
  volumeId?: number;
}

export interface Volume {
  id: number;
  title: string;
  chapters: Chapter[];
}

export interface ChatMessage {
  id: string;
  role: 'user' | 'director' | 'system';
  content: string;
  agent?: string;
  timestamp: string;
  status?: 'thinking' | 'dispatching' | 'streaming' | 'complete';
  agentOutputs?: { agent: string; content: string }[];
}

export interface WorkflowNode {
  id: string;
  type: string;
  name: string;
  icon: string;
  agentId?: number;
  x: number;
  y: number;
  config?: Record<string, unknown>;
}

export interface WorkflowEdge {
  id: string;
  from: string;
  to: string;
  label?: string;
  type: 'normal' | 'success' | 'failure';
}

export interface Workflow {
  id: number;
  name: string;
  description: string;
  type: 'system' | 'custom';
  category: string;
  icon: string;
  nodes: WorkflowNode[];
  edges: WorkflowEdge[];
  isActive: boolean;
}

export interface KnowledgeCategory {
  id: number;
  name: string;
  icon: string;
  itemCount: number;
  children?: KnowledgeCategory[];
}

export interface KnowledgeItem {
  id: number;
  title: string;
  content: string;
  tags: string[];
  useCount: number;
  createdAt: string;
}

export interface Storyline {
  id: number;
  type: 'skyline' | 'groundline' | 'plotline';
  title: string;
  content: string;
  status: 'planned' | 'active' | 'completed';
  items: StorylineItem[];
}

export interface StorylineItem {
  id: number;
  title: string;
  content: string;
  chapterRange: string;
  status: 'planned' | 'active' | 'completed';
}
