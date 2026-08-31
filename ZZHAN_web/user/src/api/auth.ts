/**
 * 前台认证接口（GitHub OAuth）：授权登录 / 刷新令牌 / 退出。
 */
import { request } from './http'
import type { AuthUser, LoginResult, RefreshResult } from '@/types/models'

/** GitHub 授权登录 POST /auth/github */
export function githubLogin(code: string, redirectUri: string): Promise<LoginResult> {
  return request<LoginResult>({ method: 'POST', url: '/auth/github', data: { code, redirect_uri: redirectUri } })
}

/** 刷新令牌 POST /auth/refresh */
export function refresh_token(refresh_token: string): Promise<RefreshResult> {
  return request<RefreshResult>({ method: 'POST', url: '/auth/refresh', data: { refresh_token } })
}

/** 退出登录 POST /auth/logout */
export function logout(): Promise<null> {
  return request<null>({ method: 'POST', url: '/auth/logout' })
}

/** 获取当前登录用户信息 GET /auth/me */
export function getCurrentUser(): Promise<AuthUser> {
  return request<AuthUser>({ method: 'GET', url: '/auth/me' })
}
