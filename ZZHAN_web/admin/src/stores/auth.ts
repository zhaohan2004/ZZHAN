/**
 * 后台认证 store — 登录 / 资料 / 退出；token 存 localStorage `ct-admin-token`。
 */
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { adminLogin, adminLogout, getProfile, saveProfile } from '@/api/admin'
import type { AdminProfile } from '@/types/models'

const TOKEN_KEY = 'ct-admin-token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const profile = ref<AdminProfile | null>(null)

  const loggedIn = computed(() => !!token.value)

  async function login(username: string, password: string, captcha: string): Promise<boolean> {
    const res = await adminLogin(username, password, captcha)
    token.value = res.accessToken
    localStorage.setItem(TOKEN_KEY, res.accessToken)
    profile.value = res.profile
    return true
  }

  async function loadProfile(): Promise<void> {
    try {
      profile.value = await getProfile()
    } catch {
      /* 静默 */
    }
  }

  /** 更新资料（可含 password 用于改密码；mock 保存后刷新可验证新凭据） */
  async function updateProfile(p: AdminProfile & { password?: string }): Promise<void> {
    profile.value = await saveProfile(p)
  }

  async function logout(): Promise<void> {
    try {
      await adminLogout()
    } catch {
      /* 静默 */
    }
    token.value = null
    profile.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  return { token, profile, loggedIn, login, loadProfile, updateProfile, logout }
})
