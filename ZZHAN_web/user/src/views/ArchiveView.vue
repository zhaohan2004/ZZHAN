<script setup lang="ts">
/** 归档页 — 按年 → 按月分组，时间轴列表（条目带分类徽章 + 浏览数）。对齐静态原型。 */
import { computed, onMounted, ref } from 'vue'
import { Calendar, ChevronRight, Eye } from 'lucide-vue-next'
import { getArchives } from '@/api/site'
import type { ArchiveItem } from '@/types/models'
import { fmtNum } from '@/utils/format'
import { useReveal } from '@/composables/useReveal'

const root = ref<HTMLElement | null>(null)
useReveal(root)

const items = ref<ArchiveItem[]>([])

const years = computed(() => {
  const map = new Map<string, ArchiveItem[]>()
  items.value.forEach((it) => {
    const arr = map.get(it.year) ?? []
    arr.push(it)
    map.set(it.year, arr)
  })
  return [...map.entries()].sort((a, b) => b[0].localeCompare(a[0]))
})

function catColorOf(name: string): string {
  return getCategoryColor(name)
}

function getCategoryColor(name: string): string {
  const PALETTE: Record<string, string> = {
    Go: '#4a8eff', MySQL: '#38bdf8', Redis: '#fb7185', Linux: '#a3e635', Docker: '#60a5fa',
    Git: '#fb923c', Gin: '#34d399', WebSocket: '#93c5fd', SSE: '#fbbf24', '数据结构与算法': '#f472b6', '计算机基础': '#94a3b8',
  }
  return PALETTE[name] || '#4a8eff'
}

onMounted(async () => {
  try {
    items.value = await getArchives()
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div ref="root">
    <section style="padding:118px 0 48px">
      <div class="container" style="max-width:800px">
        <div class="anim-fade" style="text-align:center;margin-bottom:32px">
          <div style="display:flex;align-items:center;justify-content:center;gap:10px;font-size:13px;color:var(--text-3);margin-bottom:10px">
            <router-link to="/" style="color:var(--text-2)">首页</router-link>
            <ChevronRight :size="13" />
            <span class="grad-text" style="font-weight:600">归档</span>
          </div>
          <h1 style="font-size:28px;font-weight:800;letter-spacing:-.4px">文章归档</h1>
          <p class="muted" style="margin-top:8px;font-size:14px">时间会流逝，但文字会留下。</p>
        </div>

        <div v-for="[year, months] in years" :key="year" class="archive-year reveal">
          <div class="ay-badge"><Calendar :size="18" style="color:var(--accent)" />{{ year }} 年</div>
          <div v-for="m in [...months].sort((a, b) => b.month.localeCompare(a.month))" :key="m.month" class="archive-month">
            <div class="am-head"><span class="am-dot" />{{ Number(m.month) }} 月 · {{ m.count }} 篇</div>
            <router-link v-for="a in m.articles" :key="a.id" class="arc-item" :to="`/article/${a.slug}`">
              <span class="arc-date">{{ a.date.slice(5) }}</span>
              <div style="flex:1;min-width:0">
                <div style="font-size:14.5px;font-weight:600;line-height:1.5">{{ a.title }}</div>
                <span class="badge" style="margin-top:4px;font-size:11px" :style="{ background: catColorOf(a.category) + '12', color: catColorOf(a.category) }">{{ a.category }}</span>
              </div>
              <span style="display:inline-flex;align-items:center;gap:4px;font-size:12.5px;color:var(--text-3);flex-shrink:0;margin-top:2px"><Eye :size="13" />{{ fmtNum(a.views) }}</span>
            </router-link>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
