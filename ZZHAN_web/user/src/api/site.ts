/**
 * 前台公开接口：站点信息 / 分类 / 标签 / 归档 / 关于我 / 统计。
 * 路径与参数与 docs/api.md 完全一致。
 */
import { request } from './http'
import type { AboutData, ArchiveItem, Category, SiteInfo, StatsData, Tag } from '@/types/models'

/** 站点信息 GET /site */
export function getSite(): Promise<SiteInfo> {
  return request<SiteInfo>({ method: 'GET', url: '/site' })
}

/** 分类列表 GET /categories */
export function getCategories(): Promise<Category[]> {
  return request<Category[]>({ method: 'GET', url: '/categories' })
}

/** 标签列表 GET /tags */
export function getTags(): Promise<Tag[]> {
  return request<Tag[]>({ method: 'GET', url: '/tags' })
}

/** 按月归档 GET /archives */
export function getArchives(): Promise<ArchiveItem[]> {
  return request<ArchiveItem[]>({ method: 'GET', url: '/archives' })
}

/** 关于我 GET /about */
export function getAbout(): Promise<AboutData> {
  return request<AboutData>({ method: 'GET', url: '/about' })
}

/** 站点统计 GET /stats */
export function getStats(): Promise<StatsData> {
  return request<StatsData>({ method: 'GET', url: '/stats' })
}
