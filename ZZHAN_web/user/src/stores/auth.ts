/**
 * 认证 store — GitHub OAuth 登录、退出、刷新令牌。
 * token 持久化到 localStorage；user 为内存态。
 * `ensureAuth()` 未登录时置 `loginModalOpen = true`（LoginModal 组件消费该标志）并返回 false。
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { githubLogin, logout as apiLogout, refresh_token, getCurrentUser } from '@/api/auth'
import type { AuthUser } from '@/types/models'

const ACCESS_TOKEN_KEY = 'ct-access-token'
const REFRESH_TOKEN_KEY = 'ct-refresh-token'

/** GitHub OAuth 配置 */
const GITHUB_CLIENT_ID = import.meta.env.VITE_GITHUB_CLIENT_ID || ''
const GITHUB_REDIRECT_URI = import.meta.env.VITE_GITHUB_REDIRECT_URI || ''

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(ACCESS_TOKEN_KEY))
  const refreshToken = ref<string | null>(localStorage.getItem(REFRESH_TOKEN_KEY))
  const user = ref<AuthUser | null>(null)
  /** 登录弹窗开关（LoginModal 组件读取渲染） */
  const loginModalOpen = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  /** 初始化：如果有 token，获取用户信息 */
  async function init(): Promise<void> {
    if (token.value) {
      console.log('Token 存在，尝试获取用户信息...')
      try {
        const userData = await getCurrentUser()
        console.log('获取用户信息成功:', userData)
        user.value = userData
      } catch (error) {
        console.error('获取用户信息失败:', error)
        // 不清除登录状态，让用户继续使用
        // clearAuth()
      }
    } else {
      console.log('Token 不存在，跳过获取用户信息')
    }
  }

  /** GitHub OAuth 登录：处理回调或发起授权 */
  async function loginWith(provider: 'github'): Promise<boolean> {
    if (provider === 'github') {
      // 从 URL 获取 code 参数
      const urlParams = new URLSearchParams(window.location.search)
      const code = urlParams.get('code')

      if (!code) {
        // 没有 code，重定向到 GitHub 授权页面
        const githubAuthUrl = `https://github.com/login/oauth/authorize?client_id=${GITHUB_CLIENT_ID}&redirect_uri=${encodeURIComponent(GITHUB_REDIRECT_URI)}`
        window.location.href = githubAuthUrl
        return false
      }

      // 有 code，调用后端接口
      try {
        const res = await githubLogin(code)
        setTokens(res.access_token, res.refresh_token)
        user.value = res.user
        loginModalOpen.value = false

        // 清除 URL 中的 code 参数
        const url = new URL(window.location.href)
        url.searchParams.delete('code')
        window.history.replaceState({}, '', url.toString())

        return true
      } catch (error) {
        console.error('GitHub 登录失败:', error)
        return false
      }
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
    localStorage.removeItem(ACCESS_TOKEN_KEY)
    localStorage.removeItem(REFRESH_TOKEN_KEY)
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

  // 初始化时获取用户信息
  init()

  return {
    token, refreshToken, user, loginModalOpen, isLoggedIn,
    loginWith, setTokens, refreshAccessToken, logout, ensureAuth,
  }
})
