/**
 * 前台数据模型 — 与 docs/api.md 契约字段一一对应。
 */

/** 分类（列表项，含文章数） */
export interface Category {
  id: number
  name: string
  slug: string
  icon?: string
  color: string
  desc?: string
  count?: number
}

/** 标签（列表项，含文章数） */
export interface Tag {
  id: number
  name: string
  count?: number
}

/** 文章内嵌分类对象：{ name, slug, color } */
export interface CategoryRef {
  name: string
  slug: string
  color: string
}

/** 文章摘要（列表项） */
export interface ArticleSummary {
  id: number
  slug: string
  title: string
  summary: string
  coverImage: string
  category: CategoryRef
  tags: string[]
  date: string
  updated: string
  views: number
  likes: number
  commentCount: number
  featured: boolean
  hot: boolean
}

/** 作者 */
export interface Author {
  nickname: string
  avatar: string
  role: string
}

/** 文章详情 = 摘要字段 + author + content */
export interface ArticleDetail extends ArticleSummary {
  author: Author
  content: string
}

/** 站点信息 */
export interface SiteInfo {
  name: string
  logoText: string
  tagline: string
  bio: string
  github: string
  email: string
  icp: string
  beian: string
  socials: { name: string; icon: string; url: string }[]
  /** 座右铭（关于页引用块） */
  motto?: string
  /** 所在地（关于页 + 首页 hero meta） */
  location?: string
  /** 建站年份（首页「写博客第 N 年」自动计算） */
  since?: number
  /** 头像 URL/数据 URI（关于页 + 首页 hero），为空时回退到 initials */
  avatar?: string
  /** 首页终端打字机内容（多行 `类型|文本`：tk/cm/fn，纯文本默认 tk） */
  heroTerminal?: string
  /** 站点短名（兼容字段；静态 SITE.shortName） */
  shortName?: string
  /** 作者昵称（兼容字段；静态 SITE.author） */
  author?: string
  /** 作者头衔（兼容字段；静态 SITE.role） */
  role?: string
}

/** 评论项（含楼中楼 replies） */
export interface CommentItem {
  id: number
  parentId: number | null
  userName: string
  avatar: string
  content: string
  time: string
  likeCount: number
  liked: boolean
  replies?: CommentItem[]
}

/** 发表评论请求体 */
export interface CommentDraft {
  content: string
  parentId?: number | null
  userName?: string
  email?: string
}

/** 发表评论响应 */
export interface CommentPostResult {
  id: number
  status: string
  message: string
}

/** 文章点赞响应 */
export interface LikeResult {
  liked: boolean
  likes: number
}

/** 评论点赞响应 */
export interface CommentLikeResult {
  liked: boolean
  likeCount: number
}

/** 关于我：技能 */
export interface Skill {
  name: string
  level: number
}

/** 关于我聚合数据（当前仅技能区） */
export interface AboutData {
  skills: Skill[]
}

/** 动态时间线条目 */
export interface Dynamic {
  type: string
  text: string
  time: string
  link?: string
}

/** 站点统计 */
export interface StatsData {
  articles: number
  views: number
  dynamics: Dynamic[]
}

/** 归档项 */
export interface ArchiveItem {
  year: string
  month: string
  count: number
  articles: { id: number; title: string; date: string; category: string; views: number }[]
}

/** 前台登录用户 */
export interface AuthUser {
  id: number
  provider: string
  nickname: string
  avatar: string
}

/** 登录结果（微信 / GitHub） */
export interface LoginResult {
  accessToken: string
  refreshToken: string
  expiresIn: number
  user: AuthUser
  needProfile: boolean
}

/** 刷新令牌结果 */
export interface RefreshResult {
  accessToken: string
  expiresIn: number
}

/** 完善资料结果 */
export interface UpdateProfileResult {
  user: Pick<AuthUser, 'nickname' | 'avatar'>
}
