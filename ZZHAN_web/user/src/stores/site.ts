/**
 * 站点 store — 缓存站点信息 / 统计。拉取失败静默，避免页面崩溃。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getSite, getStats } from '@/api/site'
import type { SiteInfo, StatsData } from '@/types/models'

export const useSiteStore = defineStore('site', () => {
  const site = ref<SiteInfo | null>(null)
  const stats = ref<StatsData | null>(null)

  async function fetchSite(): Promise<void> {
    try {
      site.value = await getSite()
    } catch {
      /* 静默失败 */
    }
  }

  async function fetchStats(): Promise<void> {
    try {
      stats.value = await getStats()
    } catch {
      /* 静默失败 */
    }
  }

  return { site, stats, fetchSite, fetchStats }
})
