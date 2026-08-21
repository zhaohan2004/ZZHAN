/**
 * 主题 store — 读写 localStorage `ct-theme`（dark/light），同步 html.dark。
 * 默认白天（与 index.html 内联脚本一致）；`init()` 在应用启动时调用一次，`toggle()` 反转并持久化。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

const THEME_KEY = 'ct-theme'

export const useThemeStore = defineStore('theme', () => {
  const dark = ref(false)

  /** 读取持久化值（无记录时默认白天）并应用到 <html class="dark"> */
  function init(): void {
    let saved: string | null = null
    try {
      saved = localStorage.getItem(THEME_KEY)
    } catch {
      /* ignore */
    }
    dark.value = saved ? saved === 'dark' : false
    document.documentElement.classList.toggle('dark', dark.value)
  }

  /** 反转主题并持久化 */
  function toggle(): void {
    dark.value = !dark.value
    try {
      localStorage.setItem(THEME_KEY, dark.value ? 'dark' : 'light')
    } catch {
      /* ignore */
    }
    document.documentElement.classList.toggle('dark', dark.value)
  }

  return { dark, init, toggle }
})
