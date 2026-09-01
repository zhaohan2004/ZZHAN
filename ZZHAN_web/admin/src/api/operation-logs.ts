/**
 * 操作日志接口 — 后台操作日志列表查询。
 */
import { request } from './http'
import type { Paged } from '@/types/api'
import type { OperationLogAdmin } from '@/types/models'

export interface OperationLogQuery {
  page?: number
  pageSize?: number
  action?: string
  target?: string
  startDate?: string
  endDate?: string
}

export function listOperationLogs(params: OperationLogQuery = {}): Promise<Paged<OperationLogAdmin>> {
  const { pageSize, startDate, endDate, ...rest } = params
  const query: Record<string, string | number | undefined> = { ...rest }
  if (pageSize !== undefined) query.page_size = pageSize
  if (startDate !== undefined) query.start_date = startDate
  if (endDate !== undefined) query.end_date = endDate
  return request<Paged<OperationLogAdmin>>({ method: 'GET', url: '/admin/operation-logs', params: query })
}
