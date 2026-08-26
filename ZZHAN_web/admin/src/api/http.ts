/**
 * HTTP 客户端 — 按 VITE_API_MODE 分流 mock / axios。
 * 统一注入 Authorization，解包 { code, message, data }。
 */
import axios, { type AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'
import { mockRequest } from './mock/handlers'

export interface RequestConfig extends AxiosRequestConfig {}

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
})

// 响应拦截器：处理 token 过期（401）
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401 || error.response?.data?.code === 40100) {
      // token 过期，清除本地 token
      localStorage.removeItem('ct-admin-token')
      // 如果当前不在登录页，跳转到登录页
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)

export async function request<T>(cfg: RequestConfig): Promise<T> {
  const mode = import.meta.env.VITE_API_MODE
  if (mode === 'mock') return mockRequest<T>(cfg)
  const { data } = await instance.request<ApiResponse<T>>({
    ...cfg,
    headers: {
      Authorization: `Bearer ${localStorage.getItem('ct-admin-token') || ''}`,
      ...(cfg.headers || {}),
    },
  })
  if (data.code !== 0) throw new Error(data.message)
  return data.data
}
