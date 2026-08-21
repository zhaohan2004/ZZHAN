/**
 * 仪表盘接口 — 5 个数据端点。
 */
import { request } from './http'
import type {
  DashboardArticles,
  DashboardCharts,
  DashboardStat,
  Operation,
  RecentComment,
} from '@/types/models'

export function getDashboardStats(): Promise<DashboardStat[]> {
  return request<DashboardStat[]>({ method: 'GET', url: '/admin/dashboard/stats' })
}
export function getDashboardCharts(): Promise<DashboardCharts> {
  return request<DashboardCharts>({ method: 'GET', url: '/admin/dashboard/charts' })
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
