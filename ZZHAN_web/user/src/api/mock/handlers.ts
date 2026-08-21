/**
 * Mock 适配器 — 将 `method + ' ' + url` 规范化为路由模板（路径段 `{slug}`/`{id}` → `*`），
 * 在路由表查 handler，按 params/data 过滤与分页。未命中抛 40400。
 */
import type { RequestConfig } from '../http'
import type { CommentDraft, LoginResult, SiteInfo } from '@/types/models'
import { initialsAvatar } from '@/utils/avatar'
import {
  aboutData,
  archives,
  articleSummaries,
  categories,
  commentsForArticle,
  getArticleDetail,
  siteData,
  statsData,
  tags,
} from './data'

/** 单个 mock 请求的上下文 */
export interface MockContext {
  /** 路由通配符捕获的路径段（按出现顺序） */
  path: string[]
  /** query 参数 */
  params: Record<string, unknown>
  /** 请求体 */
  data: unknown
}

export type MockHandler = (ctx: MockContext) => unknown

/** 携带业务错误码的 mock 错误 */
export class MockError extends Error {
  code: number
  constructor(code: number, message: string) {
    super(message)
    this.name = 'MockError'
    this.code = code
  }
}

/** 资源不存在 */
export function notFound(): never {
  throw new MockError(40400, '资源不存在')
}

interface MockRoute {
  pattern: RegExp
  handler: MockHandler
}

/** 由路由模板（如 `GET /articles/*`）构建匹配正则：`*` → `([^/]+)` */
function buildPattern(template: string): RegExp {
  const withWild = template.replace(/\*/g, '__WILD__')
  const escaped = withWild.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return new RegExp('^' + escaped.replace(/__WILD__/g, '([^/]+)') + '$')
}

/* ---------- handlers ---------- */

function listArticles(ctx: MockContext) {
  const p = ctx.params as Record<string, string | number | undefined>
  let list = articleSummaries
  const category = p.category
  const tag = p.tag
  const q = p.q
  const sort = p.sort
  if (category) list = list.filter(a => a.category.slug === category)
  if (tag) list = list.filter(a => a.tags.includes(String(tag)))
  if (q) {
    const kw = String(q)
    list = list.filter(a => a.title.includes(kw) || a.summary.includes(kw))
  }
  let ordered = list
  if (sort === 'hot') ordered = [...list].sort((a, b) => b.views - a.views)
  const total = ordered.length
  const page = Number(p.page || 1)
  const pageSize = Number(p.page_size ?? p.pageSize ?? 10)
  const start = (page - 1) * pageSize
  return { list: ordered.slice(start, start + pageSize), total }
}

function getArticleHandler(ctx: MockContext) {
  const article = getArticleDetail(ctx.path[0])
  if (!article) notFound()
  return article
}

function listComments(ctx: MockContext) {
  const article = getArticleDetail(ctx.path[0])
  if (!article) notFound()
  const all = commentsForArticle(article.id)
  const p = ctx.params as Record<string, string | number | undefined>
  const page = Number(p.page || 1)
  const pageSize = Number(p.page_size ?? p.pageSize ?? 10)
  const start = (page - 1) * pageSize
  return { list: all.slice(start, start + pageSize), total: all.length }
}

function postComment(ctx: MockContext) {
  const body = (ctx.data ?? {}) as CommentDraft
  if (typeof body.content !== 'string' || !body.content.trim()) {
    throw new MockError(40001, '评论内容不能为空')
  }
  return { id: 9000 + Math.floor(Math.random() * 999), status: 'normal', message: '评论已提交' }
}

function articleLike(ctx: MockContext) {
  const article = articleSummaries.find(a => a.slug === ctx.path[0])
  if (!article) notFound()
  return { liked: true, likes: article.likes }
}

function loginResult(provider: 'wechat' | 'github'): LoginResult {
  // 模拟 OAuth 直接返回真实昵称+头像（前端不再提供编辑入口）
  const nickname = provider === 'wechat' ? '微信用户' : 'GitHub User'
  const colors: [string, string] =
    provider === 'wechat' ? ['#07c160', '#00d866'] : ['#24292e', '#58a6ff']
  const avatar = initialsAvatar(nickname.charAt(0), colors[0], colors[1], 160)
  return {
    access_token: 'mock-access-token-' + provider,
    refresh_token: 'mock-refresh-token-' + provider,
    expires_in: 7200,
    user: { id: provider === 'wechat' ? 1 : 2, provider, nickname, avatar },
    need_profile: false,
  }
}

/* ---------- 站点信息：localStorage 桥（后台 `/admin/settings` 写入 `ct-site-settings`） ---------- */
const SETTINGS_KEY = 'ct-site-settings'

/**
 * GET /site — 若后台保存过设置（localStorage），映射字段覆盖默认 siteData；
 * 否则返回内置默认值。键名映射：后台 SETTINGS → 前台 SiteInfo。
 */
function siteHandler(): SiteInfo {
  let stored: Record<string, unknown> = {}
  try {
    const raw = localStorage.getItem(SETTINGS_KEY)
    if (raw) stored = JSON.parse(raw) as Record<string, unknown>
  } catch {
    /* 忽略（无 localStorage 环境） */
  }
  const s = stored as Partial<Record<string, unknown>>
  const storedSocials = Array.isArray(s.socials) ? s.socials : null
  return {
    ...siteData,
    name: String(s.blog_name || siteData.name),
    logo_text: String(s.logo_text || siteData.logo_text),
    tagline: String(s.tagline || s.blog_desc || siteData.tagline),
    bio: String(s.author_intro || siteData.bio),
    github: String(s.github || siteData.github),
    email: String(s.email || siteData.email),
    author: String(s.author_name || siteData.author),
    role: String(s.author_role || siteData.role),
    motto: String(s.motto || siteData.motto),
    location: String(s.location || siteData.location),
    since: typeof s.since === 'number' ? s.since : siteData.since,
    avatar: String(s.avatar || ''),
    hero_terminal: String(s.hero_terminal || siteData.hero_terminal || ''),
    socials: storedSocials ?? siteData.socials,
  }
}

/* ---------- 路由表 ---------- */
const routes: MockRoute[] = [
  { pattern: buildPattern('GET /site'), handler: siteHandler },
  { pattern: buildPattern('GET /articles'), handler: listArticles },
  { pattern: buildPattern('GET /articles/*/comments'), handler: listComments },
  { pattern: buildPattern('GET /articles/*'), handler: getArticleHandler },
  { pattern: buildPattern('GET /categories'), handler: () => categories },
  { pattern: buildPattern('GET /tags'), handler: () => tags },
  { pattern: buildPattern('GET /archives'), handler: () => archives() },
  { pattern: buildPattern('GET /about'), handler: () => aboutData },
  { pattern: buildPattern('GET /stats'), handler: () => statsData() },
  { pattern: buildPattern('POST /articles/*/comments'), handler: postComment },
  { pattern: buildPattern('POST /articles/*/like'), handler: articleLike },
  { pattern: buildPattern('POST /comments/*/like'), handler: () => ({ liked: true, like_count: 13 }) },
  { pattern: buildPattern('POST /auth/wechat'), handler: () => loginResult('wechat') },
  { pattern: buildPattern('POST /auth/github'), handler: () => loginResult('github') },
  { pattern: buildPattern('POST /auth/refresh'), handler: () => ({ access_token: 'mock-refreshed-token', expires_in: 7200 }) },
  { pattern: buildPattern('POST /auth/logout'), handler: () => null },
  { pattern: buildPattern('PUT /auth/profile'), handler: (ctx) => ({ user: (ctx.data ?? {}) }) },
]

/** mock 请求入口：规范化 method+url → 查路由表 → 返回 data（非 ApiResponse 包装） */
export async function mockRequest<T>(cfg: RequestConfig): Promise<T> {
  const method = (cfg.method || 'GET').toUpperCase()
  const url = cfg.url || ''
  const key = method + ' ' + url
  for (const route of routes) {
    const match = route.pattern.exec(key)
    if (match) {
      const ctx: MockContext = {
        path: match.slice(1),
        params: (cfg.params ?? {}) as Record<string, unknown>,
        data: cfg.data,
      }
      return route.handler(ctx) as T
    }
  }
  notFound()
}
