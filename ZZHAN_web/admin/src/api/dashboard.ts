/**
 * 仪表盘接口 — 4 个数据端点。
 */
import { request } from './http'
import type {
  DashboardArticles,
  DashboardStatData,
  Operation,
  RecentComment,
} from '@/types/models'

export function getDashboardStats(): Promise<DashboardStatData> {
  return request<DashboardStatData>({ method: 'GET', url: '/admin/dashboard/stats' })
}
export function getDashboardArticles(): Promise<DashboardArticles> {
  return request<DashboardArticles>({ method: 'GET', url: '/admin/dashboard/articles' })
}
export function getDashboardComments(): Promise<RecentComment[]> {
  return request<RecentComment[]>({ method: 'GET', url: '/admin/dashboard/comments' })
}
export function getDashboardOperations(): Promise<Operation[]> {
  return request<Operation[]>({ method: 'GET', url: '/admin/dashboard/operations' })
}
