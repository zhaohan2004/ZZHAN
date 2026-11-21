/**
 * 认证 store — GitHub OAuth 登录、退出、刷新令牌。
 * token 持久化到 localStorage；user 为内存态。
 * `ensureAuth()` 未登录时置 `loginModalOpen = true`（LoginModal 组件消费该标志）并返回 false。
 */
import {defineStore} from 'pinia'
import {computed, onScopeDispose, ref} from 'vue'
import {getCurrentUser, getGitHubConfig, githubLogin, logout as apiLogout, refresh_token} from '@/api/auth'
import {AccountDisabledError, clearSession, isCurrentSession, onSessionChange, readSession, saveSession, updateAccessToken} from '@/api/session'
import {useToast} from '@/composables/useToast'
import type {AuthUser} from '@/types/models'

/** GitHub OAuth 配置 */
const GITHUB_REDIRECT_URI_PATH = import.meta.env.VITE_GITHUB_REDIRECT_URI || '/auth/github/callback'

export const useAuthStore = defineStore('auth', () => {
  const savedSession = readSession()
  const token = ref(savedSession.accessToken)
  const refreshToken = ref(savedSession.refreshToken)
  const user = ref<AuthUser | null>(null)
  const { toast } = useToast()
  /** 登录弹窗开关（LoginModal 组件读取渲染） */
  const loginModalOpen = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  onScopeDispose(onSessionChange((change) => {
    token.value = change.accessToken
    if (change.accessToken === null) {
      refreshToken.value = null
      user.value = null
      if (change.error) {
        loginModalOpen.value = false
        toast(change.error.message, 'error')
      }
    }
  }))

  /** 初始化：如果有 token，获取用户信息 */
  async function init(): Promise<void> {
    if (!savedSession.accessToken) return

    try {
      const userData = await getCurrentUser()
      if (isCurrentSession(savedSession)) {
        user.value = userData
      }
    } catch {
      clearSession(savedSession)
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
        const redirectUri = `${window.location.origin}${GITHUB_REDIRECT_URI_PATH}`
        const clientId = import.meta.env.VITE_GITHUB_CLIENT_ID?.trim() || (await getGitHubConfig()).client_id
        if (!clientId) throw new Error('GitHub 登录尚未配置')
        const authorizeUrl = new URL('https://github.com/login/oauth/authorize')
        authorizeUrl.searchParams.set('client_id', clientId)
        authorizeUrl.searchParams.set('redirect_uri', redirectUri)
        window.location.href = authorizeUrl.href
        return false
      }

      // 有 code，调用后端接口
      try {
        const redirectUri = `${window.location.origin}${GITHUB_REDIRECT_URI_PATH}`
        const res = await githubLogin(code, redirectUri)
        setTokens(res.access_token, res.refresh_token)
        user.value = res.user
        loginModalOpen.value = false

        // 清除 URL 中的 code 参数
        const url = new URL(window.location.href)
        url.searchParams.delete('code')
        window.history.replaceState({}, '', url.toString())

        return true
      } catch (error) {
        if (error instanceof AccountDisabledError) throw error
        return false
      }
    }
    return false
  }

  /** 保存双 token 到 state + localStorage */
  function setTokens(access: string, refresh: string): void {
    saveSession(access, refresh)
    token.value = access
    refreshToken.value = refresh
  }

  /** 用 refresh_token 换新 access_token */
  async function refreshAccessToken(): Promise<boolean> {
    const session = readSession()
    if (!session.refreshToken) return false
    try {
      const res = await refresh_token(session.refreshToken)
      return updateAccessToken(session, res.access_token)
    } catch (error) {
      if (error instanceof AccountDisabledError) throw error
      clearSession(session)
      return false
    }
  }

  /** 退出登录：调用 API 并清空本地态 */
  async function logout(): Promise<boolean> {
    const session = readSession()
    try {
      await apiLogout()
    } catch (error) {
      if (error instanceof AccountDisabledError) return false
    } finally {
      clearSession(session)
    }
    return true
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
