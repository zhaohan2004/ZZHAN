/**
 * 后台接口 — 认证 / 文章 / 分类 / 标签 / 评论 / 设置 / 关于。
 * 路径与 docs/api.md 后台前缀 /admin 一致。
 */
import { request } from './http'
import type { Paged } from '@/types/api'
import type {
  AdminArticle,
  AdminArticlePayload,
  AdminProfile,
  ArticleStatus,
  CategoryAdmin,
  CategoryStatus,
  CommentAdmin,
  CommentStatus,
  OperationLogAdmin,
  SettingsKV,
  TagAdmin,
  TagStatus,
  UserAdmin,
  UserStatus,
} from '@/types/models'

export interface LoginResult {
  access_token: string
  refresh_token: string
  user: { id: number; provider: string; nickname: string; avatar: string }
}

export interface CaptchaResult {
  captcha_id: string
  captcha_image: string
}

/* ---------- 认证 ---------- */
export function getCaptcha(): Promise<CaptchaResult> {
  return request<CaptchaResult>({ method: 'GET', url: '/admin/auth/captcha' })
}

export function adminLogin(username: string, password: string, captchaId: string, captchaCode: string): Promise<LoginResult> {
  return request<LoginResult>({ method: 'POST', url: '/admin/auth/login', data: { username, password, captcha_id: captchaId, captcha: captchaCode } })
}
export function adminLogout(): Promise<null> {
  return request<null>({ method: 'POST', url: '/admin/auth/logout' })
}
export function getProfile(): Promise<AdminProfile> {
  return request<AdminProfile>({ method: 'GET', url: '/admin/profile' })
}
export function saveProfile(p: AdminProfile): Promise<AdminProfile> {
  return request<AdminProfile>({ method: 'PUT', url: '/admin/profile', data: p })
}

/* ---------- 文章 ---------- */
export interface AdminArticleQuery {
  page?: number
  pageSize?: number
  keyword?: string
  category?: string
  tag?: string
  status?: ArticleStatus | 'all'
}
export function listAdminArticles(params: AdminArticleQuery = {}): Promise<Paged<AdminArticle>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<AdminArticle>>({ method: 'GET', url: '/admin/articles', params: query })
}
export function getAdminArticle(id: number | string): Promise<AdminArticle> {
  return request<AdminArticle>({ method: 'GET', url: `/admin/articles/${id}` })
}
export function createAdminArticle(payload: AdminArticlePayload): Promise<AdminArticle> {
  return request<AdminArticle>({ method: 'POST', url: '/admin/articles', data: payload })
}
export function updateAdminArticle(id: number | string, payload: AdminArticlePayload): Promise<AdminArticle> {
  return request<AdminArticle>({ method: 'PUT', url: `/admin/articles/${id}`, data: payload })
}
export function deleteAdminArticle(id: number | string): Promise<null> {
  return request<null>({ method: 'DELETE', url: `/admin/articles/${id}` })
}
export function setArticleStatus(id: number | string, status: ArticleStatus): Promise<AdminArticle> {
  return request<AdminArticle>({ method: 'PUT', url: `/admin/articles/${id}/status`, data: { status } })
}

/* ---------- 分类 ---------- */
export interface CategoryQuery {
  page?: number
  pageSize?: number
  keyword?: string
  status?: CategoryStatus | 'all'
  minCount?: number
  maxCount?: number
}
export function listCategories(params: CategoryQuery = {}): Promise<Paged<CategoryAdmin>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<CategoryAdmin>>({ method: 'GET', url: '/admin/categories', params: query })
}
export function createCategory(p: { name: string; slug: string; desc: string; color: string }): Promise<CategoryAdmin> {
  return request<CategoryAdmin>({ method: 'POST', url: '/admin/categories', data: p })
}
export function updateCategory(id: number, p: Partial<{ name: string; slug: string; desc: string; color: string }>): Promise<CategoryAdmin> {
  return request<CategoryAdmin>({ method: 'PUT', url: `/admin/categories/${id}`, data: p })
}
export function deleteCategory(id: number): Promise<null> {
  return request<null>({ method: 'DELETE', url: `/admin/categories/${id}` })
}
export function setCategoryStatus(id: number, status: CategoryStatus): Promise<null> {
  return request<null>({ method: 'PUT', url: `/admin/categories/${id}/status`, data: { status } })
}

/* ---------- 标签 ---------- */
export interface TagQuery {
  page?: number
  pageSize?: number
  keyword?: string
  status?: TagStatus | 'all'
  minCount?: number
  maxCount?: number
}
export function listTags(params: TagQuery = {}): Promise<Paged<TagAdmin>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<TagAdmin>>({ method: 'GET', url: '/admin/tags', params: query })
}
export function createTag(name: string): Promise<TagAdmin> {
  return request<TagAdmin>({ method: 'POST', url: '/admin/tags', data: { name } })
}
export function updateTag(id: number, name: string): Promise<TagAdmin> {
  return request<TagAdmin>({ method: 'PUT', url: `/admin/tags/${id}`, data: { name } })
}
export function deleteTag(id: number): Promise<null> {
  return request<null>({ method: 'DELETE', url: `/admin/tags/${id}` })
}
export function setTagStatus(id: number, status: TagStatus): Promise<null> {
  return request<null>({ method: 'PUT', url: `/admin/tags/${id}/status`, data: { status } })
}

/* ---------- 评论 ---------- */
export function listComments(
  params: {
    page?: number
    pageSize?: number
    status?: CommentStatus | 'all'
    article_id?: number
    startDate?: string
    endDate?: string
  } = {},
): Promise<Paged<CommentAdmin>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<CommentAdmin>>({ method: 'GET', url: '/admin/comments', params: query })
}
export function updateCommentStatus(id: number, status: CommentStatus): Promise<CommentAdmin> {
  return request<CommentAdmin>({ method: 'PUT', url: `/admin/comments/${id}/status`, data: { status } })
}
export function deleteComment(id: number): Promise<null> {
  return request<null>({ method: 'DELETE', url: `/admin/comments/${id}` })
}

/* ---------- 用户 ---------- */
export interface UserQuery {
  page?: number
  pageSize?: number
  keyword?: string
  status?: UserStatus | 'all'
  startDate?: string
  endDate?: string
}
export function listUsers(params: UserQuery = {}): Promise<Paged<UserAdmin>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<UserAdmin>>({ method: 'GET', url: '/admin/users', params: query })
}
export function getUser(id: number): Promise<UserAdmin> {
  return request<UserAdmin>({ method: 'GET', url: `/admin/users/${id}` })
}
export function setUserStatus(id: number, status: UserStatus): Promise<UserAdmin> {
  return request<UserAdmin>({ method: 'PUT', url: `/admin/users/${id}/status`, data: { status } })
}
export function deleteUser(id: number): Promise<null> {
  return request<null>({ method: 'DELETE', url: `/admin/users/${id}` })
}

/* ---------- 操作日志 ---------- */
export function listOperationLogs(
  params: {
    page?: number
    pageSize?: number
    action?: string
    target?: string
    startDate?: string
    endDate?: string
  } = {},
): Promise<Paged<OperationLogAdmin>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<OperationLogAdmin>>({ method: 'GET', url: '/admin/operation-logs', params: query })
}

/* ---------- 上传 ---------- */
export function uploadImage(file: File): Promise<{ url: string }> {
  const form = new FormData()
  form.append('file', file)
  return request<{ url: string }>({ method: 'POST', url: '/upload/image', data: form })
}

/* ---------- 设置 ---------- */
export function getSettings(): Promise<SettingsKV> {
  return request<SettingsKV>({ method: 'GET', url: '/admin/settings' })
}
export function saveSettings(s: SettingsKV): Promise<SettingsKV> {
  return request<SettingsKV>({ method: 'PUT', url: '/admin/settings', data: s })
}
