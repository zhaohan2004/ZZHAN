/**
 * 设置 store — 加载 / 保存系统设置 KV。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSettings, saveSettings } from '@/api/admin'
import type { SettingsKV } from '@/types/models'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<SettingsKV | null>(null)

  async function load(): Promise<void> {
    try {
      settings.value = await getSettings()
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
