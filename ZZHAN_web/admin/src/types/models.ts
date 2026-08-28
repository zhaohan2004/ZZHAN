/**
 * 后台数据模型 — 与 docs/api.md 后台契约一致（单管理员，无 RBAC）。
 */

/** 分类 */
export type CategoryStatus = 'active' | 'inactive'

export interface CategoryAdmin {
  id: number
  name: string
  slug: string
  icon: string
  desc: string
  color: string
  status: CategoryStatus
  count: number
  created_at: string
  updated_at: string
}

/** 标签 */
export type TagStatus = 'active' | 'inactive'

export interface TagAdmin {
  id: number
  name: string
  status: TagStatus
  count: number
  created_at: string
  updated_at: string
}

/** 后台文章 */
export type ArticleStatus = 'published' | 'draft' | 'down'

export interface AdminArticle {
  id: number
  slug: string
  title: string
  summary: string
  cover_image: string
  category: string
  tags: string[]
  status: ArticleStatus
  views: number
  date: string
  updated: string
  content?: string
}

export interface AdminArticlePayload {
  title: string
  summary: string
  cover_image: string
  category: string
  tags: string[]
  content: string
  status: ArticleStatus
  published_at: string
}

/** 后台评论（两态：正常 / 封禁） */
export type CommentStatus = 'normal' | 'banned'

export interface CommentAdmin {
  id: number
  article_id: number
  article_title: string
  user_name: string
  avatar: string
  content: string
  ip: string
  time: string
  status: CommentStatus
}

/** 仪表盘统计卡 */
export interface DashboardStat {
  key: string
  label: string
  value: number | string
  trend: string
  up: boolean
  icon: string
  color: string
}

/** 仪表盘图表数据 */
export interface DashboardCharts {
  week_visits: { d: string; pv: number; uv: number }[]
  post_trend: { m: string; n: number }[]
  cat_dist: { name: string; value: number; color: string }[]
}

/** 仪表盘列表数据 */
export interface DashboardArticles {
  recent_posts: { id: number; title: string; category: string; date: string; views: number }[]
  hot_posts: { id: number; title: string; views: number }[]
}
export interface RecentComment {
  id: number
  user_name: string
  avatar: string
  content: string
  time: string
}
export interface Operation {
  time: string
  user: string
  action: string
  target: string
}

/** 管理员资料（password 不回传，仅前端本地输入） */
export interface AdminProfile {
  nickname: string
  avatar: string
  /** 登录账号（GET /admin/profile 返回，用于展示/修改） */
  username?: string
}

/** 系统设置 KV */
export interface SettingsKV {
  blog_name: string
  blog_desc: string
  logo_text: string
  author_name: string
  author_role: string
  author_intro: string
  github: string
  email: string
  /** 站点标语（首页 hero 副标题） */
  tagline: string
  /** 座右铭（关于页） */
  motto: string
  /** 所在地（关于页 + 首页 meta） */
  location: string
  /** 建站年份（首页「写博客第 N 年」） */
  since: number
  /** 头像 URL/数据 URI（关于页 + 首页 hero），为空时回退到 initials */
  avatar: string
  /** 首页终端打字机内容（多行，每行 `类型|文本`：tk 蓝 / cm 灰注释 / fn 紫，纯文本默认 tk） */
  hero_terminal: string
  /** 社交链接列表（前台关于页 / 首页展示） */
  socials: { name: string; icon: string; url: string }[]
  [key: string]: unknown
}
