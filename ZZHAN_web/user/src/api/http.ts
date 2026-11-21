/** 前台 HTTP 客户端。令牌刷新成功后仅重试一次请求。 */
import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ERROR_CODES, type ApiResponse } from '@/types/api'
import { mockRequest } from './mock/handlers'
import {
  AccountDisabledError, clearSession, isCurrentSession, readSession, updateAccessToken,
} from './session'

export interface RequestConfig extends AxiosRequestConfig {}

const client = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '/api/v1',
})

const refreshPath = '/auth/refresh'
type Session = ReturnType<typeof readSession>

async function send<T>(config: RequestConfig): Promise<AxiosResponse<ApiResponse<T>>> {
  try {
    return await client.request<ApiResponse<T>>(config)
  } catch (error) {
    if (axios.isAxiosError<ApiResponse<T>>(error) &&
      (error.response?.status === 401 || error.response?.status === 403)) {
      return error.response
    }
    throw error
  }
}

function checkAccountStatus(response: AxiosResponse<ApiResponse<unknown>>, session: Session): void {
  if (response.data?.code === ERROR_CODES.USER_DISABLED) {
    const error = new AccountDisabledError()
    clearSession(session, error)
    throw error
  }
}

function isUnauthorized(response: AxiosResponse<ApiResponse<unknown>>): boolean {
  return response.status === 401 || response.data?.code === 401 || response.data?.code === ERROR_CODES.UNAUTHORIZED
}

let pendingRefresh: { token: string; result: Promise<boolean> } | undefined

async function tryRefreshToken(session: Session): Promise<boolean> {
  const refreshToken = session.refreshToken
  if (!refreshToken || !isCurrentSession(session)) return false

  const currentAccessToken = readSession().accessToken
  if (currentAccessToken && currentAccessToken !== session.accessToken) return true
  if (pendingRefresh?.token === refreshToken) return pendingRefresh.result

  const result = (async () => {
    try {
      const response = await send<{ access_token: string }>({
        method: 'POST', url: refreshPath, data: { refresh_token: refreshToken },
      })
      checkAccountStatus(response, session)
      const accessToken = response.data?.data?.access_token
      return response.status < 400 && response.data?.code === ERROR_CODES.OK && !!accessToken &&
        updateAccessToken(session, accessToken)
    } catch (error) {
      if (error instanceof AccountDisabledError) throw error
      return false
    }
  })()
  pendingRefresh = { token: refreshToken, result }

  try {
    return await result
  } finally {
    if (pendingRefresh?.result === result) pendingRefresh = undefined
  }
}

export async function request<T>(config: RequestConfig): Promise<T> {
  if (import.meta.env.VITE_API_MODE === 'mock') return mockRequest<T>(config)

  const session = readSession()
  const doRequest = (accessToken: string | null) => send<T>({
    ...config,
    headers: { ...(accessToken ? { Authorization: 'Bearer ' + accessToken } : {}), ...config.headers },
  })

  let response = await doRequest(session.accessToken)
  checkAccountStatus(response, session)

  if (isUnauthorized(response)) {
    if (config.url !== refreshPath && await tryRefreshToken(session) && isCurrentSession(session)) {
      response = await doRequest(readSession().accessToken)
      checkAccountStatus(response, session)
    }

    if (isUnauthorized(response)) {
      clearSession(session)
      throw new Error('登录已过期，请重新登录')
    }
  }

  if (response.status >= 400 || response.data?.code !== ERROR_CODES.OK) {
    throw new Error(response.data?.message || '请求失败，请重试')
  }
  return response.data.data
}
