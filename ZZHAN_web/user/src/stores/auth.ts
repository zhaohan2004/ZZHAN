/**
 * 认证 store — 微信 / GitHub OAuth 登录、完善资料、退出。
 * token 持久化到 localStorage `ct-access-token`；user / needProfile 为内存态。
 * `ensureAuth()` 未登录时置 `loginModalOpen = true`（LoginModal 组件消费该标志）并返回 false。
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { githubLogin, logout as apiLogout, updateProfile, wechatLogin } from '@/api/auth'
import type { AuthUser } from '@/types/models'

export type LoginProvider = 'github' | 'wechat'

const TOKEN_KEY = 'ct-access-token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const user = ref<AuthUser | null>(null)
  const needProfile = ref(false)
  /** 登录弹窗开关（Task 17 的 LoginModal 读取渲染） */
  const loginModalOpen = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  /** OAuth 登录：填写 token / user / needProfile，返回是否成功 */
  async function loginWith(provider: LoginProvider, code: string): Promise<boolean> {
    const fn = provider === 'github' ? githubLogin : wechatLogin
    const res = await fn(code)
    token.value = res.accessToken
    localStorage.setItem(TOKEN_KEY, res.accessToken)
    user.value = res.user
    needProfile.value = res.needProfile
    loginModalOpen.value = false
    return true
  }

  /** 完善资料：PUT /auth/profile，更新本地 user 并清除 needProfile */
  async function completeProfile(p: Pick<AuthUser, 'nickname' | 'avatar'>): Promise<boolean> {
    await updateProfile(p)
    user.value = { ...(user.value ?? ({} as AuthUser)), ...p }
    needProfile.value = false
    return true
  }

  /** 退出登录：调用 API 并清空本地态 */
  async function logout(): Promise<void> {
    try {
      await apiLogout()
    } catch {
      /* 静默失败，仍清空本地 */
    }
    token.value = null
    user.value = null
    needProfile.value = false
    localStorage.removeItem(TOKEN_KEY)
  }

  /** 登录门卫：已登录返回 true；未登录弹 LoginModal 并返回 false */
  async function ensureAuth(): Promise<boolean> {
    if (token.value) return true
    loginModalOpen.value = true
    return false
  }

  return { token, user, needProfile, loginModalOpen, isLoggedIn, loginWith, completeProfile, logout, ensureAuth }
})
