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

  // 监听 storage 事件，当 token 被其他地方清除时自动更新状态
  function onStorage(e: StorageEvent) {
    if (e.key === TOKEN_KEY) {
      token.value = e.newValue
      if (!e.newValue) {
        profile.value = null
      }
    }
  }
  window.addEventListener('storage', onStorage)

  async function login(username: string, password: string, captchaId: string, captchaCode: string): Promise<boolean> {
    const res = await adminLogin(username, password, captchaId, captchaCode)
    token.value = res.access_token
    localStorage.setItem(TOKEN_KEY, res.access_token)
    // 登录后立即加载管理员资料
    await loadProfile()
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
