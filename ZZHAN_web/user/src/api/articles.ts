/**
 * 文章接口：列表 / 详情 / 点赞。
 * 参数与 docs/api.md 完全一致（page / page_size / category / tag / q / sort）。
 */
import { request } from './http'
import type { ArticleDetail, ArticleSummary, LikeResult } from '@/types/models'
import type { Paged } from '@/types/api'

export interface ArticleQuery {
  page?: number
  pageSize?: number
  category?: string
  tag?: string
  q?: string
  sort?: 'latest' | 'hot'
}

/** 文章列表 GET /articles */
export function getArticles(params: ArticleQuery = {}): Promise<Paged<ArticleSummary>> {
  const { pageSize, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  return request<Paged<ArticleSummary>>({ method: 'GET', url: '/articles', params: query })
}

/** 文章详情 GET /articles/{slug}（slug 或数字 id） */
export function getArticle(slug: string): Promise<ArticleDetail> {
  return request<ArticleDetail>({ method: 'GET', url: `/articles/${slug}` })
}

/** 文章点赞 POST /articles/{slug}/like（幂等切换） */
export function toggleLike(slug: string): Promise<LikeResult> {
  return request<LikeResult>({ method: 'POST', url: `/articles/${slug}/like` })
}

// 站点/分类/标签/关于我 — 供统一入口与 mock 测试从 './articles' 导入
export { getSite, getCategories, getTags, getAbout } from './site'
