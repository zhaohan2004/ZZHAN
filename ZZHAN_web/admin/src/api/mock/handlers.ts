/**
 * 后台 Mock 适配器 — method + url 路由表 → handler；内存数组 CRUD。
 * 未命中抛 40400；登录校验账号 admin/123456 + 验证码非空。
 */
import type { RequestConfig } from '../http'
import { MockError, notFound } from './mock-util'
import {
  ADMIN_PROFILE,
  DASHBOARD_ARTICLES,
  DASHBOARD_CHARTS,
  DASHBOARD_COMMENTS,
  DASHBOARD_OPERATIONS,
  DASHBOARD_STATS,
  MOCK_ARTICLES,
  MOCK_CATEGORIES,
  MOCK_COMMENTS,
  MOCK_TAGS,
  SETTINGS,
  type MockArticle,
} from './data'

export type { MockContext } from './mock-util'

let articles: MockArticle[] = [...MOCK_ARTICLES]
let categories = [...MOCK_CATEGORIES]
let tags = [...MOCK_TAGS]
let comments = [...MOCK_COMMENTS]
let profile = { ...ADMIN_PROFILE }
/** 登录凭据（账号/密码，默认 admin/123456；可在「个人信息」中修改，mock 内存态刷新重置） */
let credentials = { username: 'admin', password: '123456' }

/* ---------- 设置：localStorage 桥（与前台共享 `ct-site-settings`） ---------- */
const SETTINGS_KEY = 'ct-site-settings'

/** 读取设置：优先 localStorage（后台改过 + 刷新保留），否则默认 SETTINGS。 */
function loadSettings(): Record<string, unknown> {
  try {
    const raw = localStorage.getItem(SETTINGS_KEY)
    if (raw) return { ...SETTINGS, ...(JSON.parse(raw) as object) }
  } catch {
    /* 忽略（SSR/隐私模式等无 localStorage 环境） */
  }
  return { ...SETTINGS }
}

let settings = loadSettings()

/** 写入 localStorage，让前台 `/site` 能读到后台保存的设置。 */
function persistSettings(): void {
  try {
    localStorage.setItem(SETTINGS_KEY, JSON.stringify(settings))
  } catch {
    /* 忽略 */
  }
}

function resetData(): void {
  articles = [...MOCK_ARTICLES]
  categories = [...MOCK_CATEGORIES]
  tags = [...MOCK_TAGS]
  comments = [...MOCK_COMMENTS]
  profile = { ...ADMIN_PROFILE }
  credentials = { username: 'admin', password: '123456' }
  settings = { ...SETTINGS }
  try {
    localStorage.removeItem(SETTINGS_KEY)
  } catch {
    /* 忽略 */
  }
}

export function mockRequest<T>(cfg: RequestConfig): Promise<T> {
  const method = (cfg.method || 'GET').toUpperCase()
  const url = (cfg.url || '').replace(/^\/+/, '/')
  const key = method + ' ' + url
  const params = (cfg.params ?? {}) as Record<string, string | number | undefined>
  const data = (cfg.data ?? {}) as Record<string, unknown>

  return Promise.resolve().then(() => {
    /* ---------- 认证 ---------- */
    if (key === 'POST /admin/auth/login') {
      const u = String(data.username || '')
      const p = String(data.password || '')
      const c = String(data.captcha || '')
      if (!u || !p) throw new MockError(40001, '请输入用户名和密码')
      if (!c) throw new MockError(40001, '请输入验证码')
      if (u !== credentials.username || p !== credentials.password) throw new MockError(40100, '用户名或密码错误')
      return { access_token: 'mock-admin-token', expires_in: 7200, profile: { ...profile, username: credentials.username } }
    }
    if (key === 'POST /admin/auth/logout') return null
    if (key === 'GET /admin/profile') return { ...profile, username: credentials.username }
    if (key === 'PUT /admin/profile') {
      profile = { ...profile, ...(data as object) }
      // 账号 / 密码仅在校验通过后更新凭据（模拟真实系统修改登录信息）
      if (typeof data.username === 'string' && String(data.username).trim()) {
        credentials.username = String(data.username).trim()
      }
      if (typeof data.password === 'string' && String(data.password).trim()) {
        credentials.password = String(data.password)
      }
      return { ...profile, username: credentials.username }
    }

    /* ---------- 仪表盘 ---------- */
    if (key === 'GET /admin/dashboard/stats') return DASHBOARD_STATS
    if (key === 'GET /admin/dashboard/charts') return DASHBOARD_CHARTS
    if (key === 'GET /admin/dashboard/articles') return DASHBOARD_ARTICLES
    if (key === 'GET /admin/dashboard/comments') return DASHBOARD_COMMENTS
    if (key === 'GET /admin/dashboard/operations') return DASHBOARD_OPERATIONS

    /* ---------- 文章 ---------- */
    if (key === 'GET /admin/articles') {
      let list = articles
      const q = String(params.q || '').toLowerCase()
      if (q) list = list.filter((a) => a.title.toLowerCase().includes(q))
      if (params.category && params.category !== 'all') list = list.filter((a) => a.category === params.category)
      if (params.status && params.status !== 'all') list = list.filter((a) => a.status === params.status)
      const page = Number(params.page || 1)
      const pageSize = Number(params.page_size ?? params.pageSize ?? 8)
      const start = (page - 1) * pageSize
      return { list: list.slice(start, start + pageSize), total: list.length }
    }
    const articleMatch = /^GET \/admin\/articles\/(\d+)$/.exec(key)
    if (articleMatch) {
      const a = articles.find((x) => x.id === Number(articleMatch[1]))
      if (!a) notFound()
      return a
    }
    if (key === 'POST /admin/articles') {
      const a: MockArticle = {
        id: Math.max(0, ...articles.map((x) => x.id)) + 1,
        slug: String(data.slug || 'article-' + Date.now()),
        title: String(data.title || '未命名文章'),
        summary: String(data.summary || ''),
        cover_image: String(data.cover_image || ''),
        category: String(data.category || 'Go'),
        tags: (data.tags as string[]) || [],
        status: (data.status as MockArticle['status']) || 'draft',
        views: 0,
        date: String(data.published_at || new Date().toISOString().slice(0, 10)),
        updated: new Date().toISOString().slice(0, 10),
        content: String(data.content || ''),
      }
      articles.unshift(a)
      return a
    }
    const putArticle = /^PUT \/admin\/articles\/(\d+)$/.exec(key)
    if (putArticle) {
      const a = articles.find((x) => x.id === Number(putArticle[1]))
      if (!a) notFound()
      Object.assign(a, {
        title: data.title !== undefined ? String(data.title) : a.title,
        summary: data.summary !== undefined ? String(data.summary) : a.summary,
        cover_image: data.cover_image !== undefined ? String(data.cover_image) : a.cover_image,
        category: data.category !== undefined ? String(data.category) : a.category,
        tags: data.tags !== undefined ? (data.tags as string[]) : a.tags,
        status: data.status !== undefined ? (data.status as MockArticle['status']) : a.status,
        content: data.content !== undefined ? String(data.content) : a.content,
        updated: new Date().toISOString().slice(0, 10),
      })
      return a
    }
    const statusMatch = /^PUT \/admin\/articles\/(\d+)\/status$/.exec(key)
    if (statusMatch) {
      const a = articles.find((x) => x.id === Number(statusMatch[1]))
      if (!a) notFound()
      a.status = data.status as MockArticle['status']
      a.updated = new Date().toISOString().slice(0, 10)
      return a
    }
    const delArticle = /^DELETE \/admin\/articles\/(\d+)$/.exec(key)
    if (delArticle) {
      articles = articles.filter((x) => x.id !== Number(delArticle[1]))
      return null
    }

    /* ---------- 分类 ---------- */
    if (key === 'GET /admin/categories') {
      let list = categories
      const q = String(params.q || '').toLowerCase()
      if (q) list = list.filter((c) => c.name.toLowerCase().includes(q))
      const minC = Number(params.min_count ?? params.minCount)
      const maxC = Number(params.max_count ?? params.maxCount)
      if (!Number.isNaN(minC)) list = list.filter((c) => c.count >= minC)
      if (!Number.isNaN(maxC)) list = list.filter((c) => c.count <= maxC)
      const page = Number(params.page || 1)
      const pageSize = Number(params.page_size ?? params.pageSize ?? 8)
      const start = (page - 1) * pageSize
      return { list: list.slice(start, start + pageSize), total: list.length }
    }
    if (key === 'POST /admin/categories') {
      const c = {
        id: Math.max(0, ...categories.map((x) => x.id)) + 1,
        name: String(data.name || ''),
        slug: String(data.slug || ''),
        icon: 'folder',
        desc: String(data.desc || ''),
        color: String(data.color || '#3b82f6'),
        count: 0,
        created_at: new Date().toISOString().slice(0, 10),
        updated_at: new Date().toISOString().slice(0, 10),
      }
      categories.push(c)
      return c
    }
    const putCat = /^PUT \/admin\/categories\/(\d+)$/.exec(key)
    if (putCat) {
      const c = categories.find((x) => x.id === Number(putCat[1]))
      if (!c) notFound()
      Object.assign(c, data)
      c.updated_at = new Date().toISOString().slice(0, 10)
      return c
    }
    const delCat = /^DELETE \/admin\/categories\/(\d+)$/.exec(key)
    if (delCat) {
      categories = categories.filter((x) => x.id !== Number(delCat[1]))
      return null
    }

    /* ---------- 标签 ---------- */
    if (key === 'GET /admin/tags') {
      let list = tags
      const q = String(params.q || '').toLowerCase()
      if (q) list = list.filter((t) => t.name.toLowerCase().includes(q))
      const minC = Number(params.min_count ?? params.minCount)
      const maxC = Number(params.max_count ?? params.maxCount)
      if (!Number.isNaN(minC)) list = list.filter((t) => t.count >= minC)
      if (!Number.isNaN(maxC)) list = list.filter((t) => t.count <= maxC)
      const page = Number(params.page || 1)
      const pageSize = Number(params.page_size ?? params.pageSize ?? 8)
      const start = (page - 1) * pageSize
      return { list: list.slice(start, start + pageSize), total: list.length }
    }
    if (key === 'POST /admin/tags') {
      const today = new Date().toISOString().slice(0, 10)
      const t = { id: Math.max(0, ...tags.map((x) => x.id)) + 1, name: String(data.name || ''), count: 0, created_at: today, updated_at: today }
      tags.push(t)
      return t
    }
    const putTag = /^PUT \/admin\/tags\/(\d+)$/.exec(key)
    if (putTag) {
      const t = tags.find((x) => x.id === Number(putTag[1]))
      if (!t) notFound()
      t.name = String(data.name || t.name)
      t.updated_at = new Date().toISOString().slice(0, 10)
      return t
    }
    const delTag = /^DELETE \/admin\/tags\/(\d+)$/.exec(key)
    if (delTag) {
      tags = tags.filter((x) => x.id !== Number(delTag[1]))
      return null
    }

    /* ---------- 评论 ---------- */
    if (key === 'GET /admin/comments') {
      let list = comments
      if (params.status && params.status !== 'all') list = list.filter((c) => c.status === params.status)
      const article_id = Number(params.article_id ?? params.article_id)
      if (!Number.isNaN(article_id) && article_id > 0) list = list.filter((c) => c.article_id === article_id)
      const startDate = String((params.start_date ?? params.startDate) || '')
      const endDate = String((params.end_date ?? params.endDate) || '')
      if (startDate) list = list.filter((c) => c.time >= startDate)
      if (endDate) list = list.filter((c) => c.time <= endDate + ' 23:59:59')
      const page = Number(params.page || 1)
      const pageSize = Number(params.page_size ?? params.pageSize ?? 10)
      const start = (page - 1) * pageSize
      return { list: list.slice(start, start + pageSize), total: list.length }
    }
    const putComment = /^PUT \/admin\/comments\/(\d+)\/status$/.exec(key)
    if (putComment) {
      const c = comments.find((x) => x.id === Number(putComment[1]))
      if (!c) notFound()
      c.status = data.status as (typeof c)['status']
      return c
    }
    const delComment = /^DELETE \/admin\/comments\/(\d+)$/.exec(key)
    if (delComment) {
      comments = comments.filter((x) => x.id !== Number(delComment[1]))
      return null
    }

    /* ---------- 设置 ---------- */
    if (key === 'GET /admin/settings') return settings
    if (key === 'PUT /admin/settings') {
      settings = { ...settings, ...(data as object) }
      persistSettings()
      return settings
    }

    notFound()
  }) as Promise<T>
}

export { resetData }
