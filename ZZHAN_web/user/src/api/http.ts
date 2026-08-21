/**
 * HTTP 客户端 — 按 VITE_API_MODE 分流：
 * - mock：走 mock 适配器
 * - api：axios + baseURL + 注入 Authorization + 解包 { code, message, data }
 */
import axios, { type AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'
import { mockRequest } from './mock/handlers'

export interface RequestConfig extends AxiosRequestConfig {}

export async function request<T>(cfg: RequestConfig): Promise<T> {
  const mode = import.meta.env.VITE_API_MODE
  if (mode === 'mock') return mockRequest<T>(cfg)
  const { data } = await axios.request<ApiResponse<T>>({
    ...cfg,
    baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('ct-access-token') || ''}`,
      ...(cfg.headers || {}),
    },
  })
  if (data.code !== 0) throw new Error(data.message)
  return data.data
}
