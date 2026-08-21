/**
 * 认证 store — GitHub OAuth 登录、完善资料、退出、刷新令牌。
 * token 持久化到 localStorage；user / need_profile 为内存态。
 * `ensureAuth()` 未登录时置 `loginModalOpen = true`（LoginModal 组件消费该标志）并返回 false。
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { githubLogin, logout as apiLogout, refresh_token, updateProfile } from '@/api/auth'
import type { AuthUser } from '@/types/models'

const ACCESS_TOKEN_KEY = 'ct-access-token'
const REFRESH_TOKEN_KEY = 'ct-refresh-token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY))
  const user = ref<AuthUser | null>(null)
  const need_profile = ref(false)
  /** 登录弹窗开关（LoginModal 组件读取渲染） */
  const loginModalOpen = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  /** GitHub OAuth 登录：填写 token / user / need_profile，返回是否成功 */
  async function loginWith(provider: 'github'): Promise<boolean> {
    if (provider === 'github') {
      const res = await githubLogin()
      setTokens(res.access_token, res.refresh_token)
      user.value = res.user
      need_profile.value = res.need_profile
      loginModalOpen.value = false
      return true
    }
    return false
  }

  /** 保存双 token 到 state + localStorage */
  function setTokens(access: string, refresh: string): void {
    token.value = access
    refreshToken.value = refresh
    localStorage.setItem(ACCESS_TOKEN_KEY, access)
    localStorage.setItem(REFRESH_TOKEN_KEY, refresh)
  }

  /** 用 refresh_token 换新 access_token */
  async function refreshAccessToken(): Promise<boolean> {
    if (!refreshToken.value) return false
    try {
      const res = await refresh_token(refreshToken.value)
      token.value = res.access_token
      localStorage.setItem(ACCESS_TOKEN_KEY, res.access_token)
      return true
    } catch {
      // refresh 也失败了，需要重新登录
      clearAuth()
      return false
    }
  }

  /** 清空认证状态 */
  function clearAuth(): void {
    token.value = null
    refreshToken.value = null
    user.value = null
    need_profile.value = false
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
  }

  /** 完善资料：PUT /auth/profile，更新本地 user 并清除 need_profile */
  async function completeProfile(p: Pick<AuthUser, 'nickname' | 'avatar'>): Promise<boolean> {
    await updateProfile(p)
    user.value = { ...(user.value ?? ({} as AuthUser)), ...p }
    need_profile.value = false
    return true
  }

  /** 退出登录：调用 API 并清空本地态 */
  async function logout(): Promise<void> {
    try {
      await apiLogout()
    } catch {
      /* 静默失败，仍清空本地 */
    }
    clearAuth()
  }

  /** 登录门卫：已登录返回 true；未登录弹 LoginModal 并返回 false */
  async function ensureAuth(): Promise<boolean> {
    if (token.value) return true
    loginModalOpen.value = true
    return false
  }

  return {
    token, refreshToken, user, need_profile, loginModalOpen, isLoggedIn,
    loginWith, setTokens, refreshAccessToken, completeProfile, logout, ensureAuth,
  }
})
