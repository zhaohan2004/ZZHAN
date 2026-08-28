/**
 * HTTP 客户端 — 按 VITE_API_MODE 分流 mock / axios。
 * 统一注入 Authorization，解包 { code, message, data }。
 * 支持 token 自动刷新：access_token 过期时用 refresh_token 换新 token。
 */
import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types/api'
import { mockRequest } from './mock/handlers'

export interface RequestConfig extends AxiosRequestConfig {}

const TOKEN_KEY = 'ct-admin-token'
const REFRESH_KEY = 'ct-admin-refresh-token'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
})

/* ---------- token 刷新状态 ---------- */

let refreshing = false
let pendingQueue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

/** 刷新完成后，用新 token 重放队列中所有等待的请求 */
function onRefreshed(newToken: string) {
  pendingQueue.forEach(({ resolve }) => resolve(newToken))
  pendingQueue = []
}

/** 刷新失败，拒绝队列中所有等待的请求 */
function onRefreshFailed(err: unknown) {
  pendingQueue.forEach(({ reject }) => reject(err))
  pendingQueue = []
}

/** 用 refresh_token 换取新的 access_token */
async function doRefreshToken(): Promise<string> {
  const refreshToken = localStorage.getItem(REFRESH_KEY)
  if (!refreshToken) throw new Error('无 refresh_token')

  const { data } = await instance.post<ApiResponse<{ access_token: string }>>('/admin/auth/refresh', {
    refresh_token: refreshToken,
  })

  if (data.code !== 0) throw new Error(data.message)

  const newToken = data.data.access_token
  localStorage.setItem(TOKEN_KEY, newToken)
  return newToken
}

/** 清除所有 token 并跳转登录页 */
function redirectToLogin() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_KEY)
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}

/* ---------- 响应拦截器 ---------- */

instance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalConfig = error.config as AxiosRequestConfig & { _retried?: boolean }
    const status = error.response?.status
    const code = error.response?.data?.code

    // 判断是否为 token 过期
    const isTokenExpired = status === 401 || code === 40100

    // 非 token 过期、已重试过、或正在刷新的请求本身失败 → 直接拒绝
    if (!isTokenExpired || originalConfig._retried || originalConfig.url?.includes('/admin/auth/refresh')) {
      if (isTokenExpired) redirectToLogin()
      // 用后端返回的 message 替换 axios 默认错误信息
      const serverMsg = error.response?.data?.message
      if (serverMsg) {
        return Promise.reject(new Error(serverMsg))
      }
      return Promise.reject(error)
    }

    // 标记已重试，防止无限循环
    originalConfig._retried = true

    if (refreshing) {
      // 已经有刷新请求在飞，挂起等待结果
      return new Promise<string>((resolve, reject) => {
        pendingQueue.push({ resolve, reject })
      }).then((newToken) => {
        // 拿到新 token 后重放原请求
        originalConfig.headers = {
          ...originalConfig.headers,
          Authorization: `Bearer ${newToken}`,
        }
        return instance.request(originalConfig)
      })
    }

    // 发起刷新
    refreshing = true
    try {
      const newToken = await doRefreshToken()
      onRefreshed(newToken)

      // 重放当前请求
      originalConfig.headers = {
        ...originalConfig.headers,
        Authorization: `Bearer ${newToken}`,
      }
      return instance.request(originalConfig)
    } catch (refreshErr) {
      onRefreshFailed(refreshErr)
      redirectToLogin()
      return Promise.reject(refreshErr)
    } finally {
      refreshing = false
    }
  },
)

/* ---------- 请求函数 ---------- */

export async function request<T>(cfg: RequestConfig): Promise<T> {
  const mode = import.meta.env.VITE_API_MODE
  if (mode === 'mock') return mockRequest<T>(cfg)
  const { data } = await instance.request<ApiResponse<T>>({
    ...cfg,
    headers: {
      Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY) || ''}`,
      ...(cfg.headers || {}),
    },
  })
  if (data.code !== 0) throw new Error(data.message)
  return data.data
}
