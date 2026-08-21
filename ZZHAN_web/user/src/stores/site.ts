/**
 * 站点 store — 缓存站点信息 / 统计 / 关于我。拉取失败静默，避免页面崩溃。
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAbout, getSite, getStats } from '@/api/site'
import type { AboutData, SiteInfo, StatsData } from '@/types/models'

export const useSiteStore = defineStore('site', () => {
  const site = ref<SiteInfo | null>(null)
  const stats = ref<StatsData | null>(null)
  const about = ref<AboutData | null>(null)

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

  async function fetchAbout(): Promise<void> {
    try {
      about.value = await getAbout()
    } catch {
      /* 静默失败 */
    }
  }

  return { site, stats, about, fetchSite, fetchStats, fetchAbout }
})
