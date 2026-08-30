/**
 * 文章接口：列表 / 详情 / 点赞。
 * 参数与后端保持一致（page / size / category_id / tag_id / keyword）。
 */
import { request } from './http'
import type { ArticleDetail, ArticleSummary, LikeResult } from '@/types/models'
import type { Paged } from '@/types/api'

export interface ArticleQuery {
  page?: number
  size?: number
  category_id?: number
  tag_id?: number
  keyword?: string
  sort?: 'latest' | 'hot'
}

/** 文章列表 GET /articles */
export function getArticles(params: ArticleQuery = {}): Promise<Paged<ArticleSummary>> {
  return request<Paged<ArticleSummary>>({ method: 'GET', url: '/articles', params })
}

/** 文章详情 GET /articles/{slug}（slug 或数字 id） */
export function getArticle(slug: string): Promise<ArticleDetail> {
  return request<ArticleDetail>({ method: 'GET', url: `/articles/${slug}` })
}

/** 文章点赞 POST /like/article/{slug}（幂等切换） */
export function toggleLike(slug: string): Promise<LikeResult> {
  return request<LikeResult>({ method: 'POST', url: `/like/article/${slug}` })
}

// 站点/分类/标签 — 供统一入口与 mock 测试从 './articles' 导入
export { getSite, getCategories, getTags } from './site'
