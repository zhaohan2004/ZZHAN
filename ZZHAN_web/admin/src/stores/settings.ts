/**
 * 设置 store — 加载 / 保存系统设置 KV。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSettings, saveSettings, getPublicSite } from '@/api/admin'
import type { SettingsKV } from '@/types/models'
import { useAuthStore } from './auth'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<SettingsKV | null>(null)

  async function load(): Promise<void> {
    try {
      const auth = useAuthStore()
      if (auth.loggedIn) {
        settings.value = await getSettings()
      } else {
        // 未登录时用公开接口获取基本信息
        const site = await getPublicSite()
        settings.value = site as SettingsKV
      }
    } catch {
      /* 静默 */
    }
  }

  async function save(): Promise<void> {
    if (!settings.value) return
    settings.value = await saveSettings(settings.value)
  }

  return { settings, load, save }
})
