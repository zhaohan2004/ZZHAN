/**
 * 主题 store — 读写 localStorage `ct-theme`，同步 html.dark（与前台一致）。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'

const THEME_KEY = 'ct-theme'

export const useThemeStore = defineStore('theme', () => {
  const dark = ref(false)

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
