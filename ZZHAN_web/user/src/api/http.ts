/**
 * HTTP 客户端 — 按 VITE_API_MODE 分流：
 * - mock：走 mock 适配器
 * - api：axios + baseURL + 注入 Authorization + 解包 { code, message, data }
 * - 自动刷新 token：遇到 401 时尝试用 refresh_token 换新 access_token 再重试
 */
import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types/api'
import { mockRequest } from './mock/handlers'

export interface RequestConfig extends AxiosRequestConfig {}

const ACCESS_TOKEN_KEY = 'ct-access-token'
const REFRESH_TOKEN_KEY = 'ct-refresh-token'

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
})

/** 注入 Authorization 头 */
function authHeader(): Record<string, string> {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY)
  return token ? { Authorization: `Bearer ${token}` } : {}
}

/** 用 refresh_token 换新 access_token */
async function tryRefreshToken(): Promise<boolean> {
  const rt = localStorage.getItem(REFRESH_TOKEN_KEY)
  if (!rt) return false
  try {
    const { data } = await client.post<{ code: number; data?: { access_token: string } }>('/auth/refresh', { refresh_token: rt })
    if (data.code === 0 && data.data?.access_token) {
      localStorage.setItem(ACCESS_TOKEN_KEY, data.data.access_token)
      return true
    }
    return false
  } catch {
    return false
  }
}

/** 清空认证状态 */
function clearAuth(): void {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

export async function request<T>(cfg: RequestConfig): Promise<T> {
  const mode = import.meta.env.VITE_API_MODE
  if (mode === 'mock') return mockRequest<T>(cfg)

  const doRequest = async (): Promise<AxiosResponse<ApiResponse<T>>> => {
    return client.request<ApiResponse<T>>({
      ...cfg,
      headers: { ...authHeader(), ...(cfg.headers || {}) },
    })
  }

  let resp = await doRequest()

  // 如果返回 401，尝试刷新 token 后重试
  if (resp.data.code === 40100) {
    const refreshed = await tryRefreshToken()
    if (refreshed) {
      resp = await doRequest()
    } else {
      clearAuth()
      throw new Error('登录已过期，请重新登录')
    }
  }

  if (resp.data.code !== 0) throw new Error(resp.data.message)
  return resp.data.data
}
