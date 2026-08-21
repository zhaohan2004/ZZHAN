/**
 * 前台认证接口（OAuth）：微信 / GitHub / 刷新令牌 / 退出 / 完善资料。
 */
import { request } from './http'
import type { AuthUser, LoginResult, RefreshResult, UpdateProfileResult } from '@/types/models'

/** 微信扫码登录 POST /auth/wechat */
export function wechatLogin(code: string): Promise<LoginResult> {
  return request<LoginResult>({ method: 'POST', url: '/auth/wechat', data: { code } })
}

/** GitHub 授权登录 POST /auth/github */
export function githubLogin(code: string): Promise<LoginResult> {
  return request<LoginResult>({ method: 'POST', url: '/auth/github', data: { code } })
}

/** 刷新令牌 POST /auth/refresh */
export function refreshToken(refreshToken: string): Promise<RefreshResult> {
  return request<RefreshResult>({ method: 'POST', url: '/auth/refresh', data: { refreshToken } })
}

/** 退出登录 POST /auth/logout */
export function logout(): Promise<null> {
  return request<null>({ method: 'POST', url: '/auth/logout' })
}

/** 完善资料 PUT /auth/profile */
export function updateProfile(p: Pick<AuthUser, 'nickname' | 'avatar'>): Promise<UpdateProfileResult> {
  return request<UpdateProfileResult>({ method: 'PUT', url: '/auth/profile', data: p })
}
