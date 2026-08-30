<script setup lang="ts">
/** 关于我 — 头像/简介/座右铭/统计卡片。所有可编辑项从 site 数据读取（来源：后台系统设置）。 */
import { computed, onMounted, ref } from 'vue'
import { BookOpen, Calendar, Code2, Eye, Github, Mail, MapPin, MessageCircle, Quote } from 'lucide-vue-next'
import { useSiteStore } from '@/stores/site'
import { initialsAvatar } from '@/utils/avatar'
import { useReveal } from '@/composables/useReveal'

const root = ref<HTMLElement | null>(null)
useReveal(root)

const site = useSiteStore()

onMounted(() => {
  site.fetchSite()
  site.fetchStats()
})

/** 统计卡片 */
const yearCount = computed(() => {
  const base = site.site?.since ?? new Date().getFullYear() - 7
  return Math.max(1, new Date().getFullYear() - base)
})
const statsCards = computed(() => [
  { icon: BookOpen, label: '文章', value: site.stats?.articles ?? 0 },
  { icon: Eye, label: '浏览', value: site.stats?.views ?? 0 },
  { icon: MessageCircle, label: '评论', value: site.stats?.comments ?? 0 },
  { icon: Calendar, label: '写作年份', value: yearCount.value, suffix: '年' },
])

/** 头像 — 优先取后台设置的自定义头像，否则按作者名生成 initials。 */
const author_name = computed(() => site.site?.author || '小猫')
const author_role = computed(() => site.site?.role || 'Gopher · 后端工程师')
const authorLocation = computed(() => site.site?.location || '河南')
const authorMotto = computed(
  () => site.site?.motto || '「写代码是跟计算机对话，写博客是跟自己对话。」',
)
const avatarSrc = computed(() =>
  site.site?.avatar || initialsAvatar(author_name.value, '#6b7280', '#9ca3af', 260),
)
</script>

<template>
  <div ref="root">
    <section style="padding:118px 0 60px">
      <div class="container" style="max-width:1000px">
        <!-- 头部 -->
        <div class="glass-card about-head">
          <img class="about-avatar" :src="avatarSrc" :alt="author_name" />
          <div style="flex:1;min-width:0">
            <h1 style="font-size:28px;font-weight:800;letter-spacing:-.4px">
              <span class="grad-text">{{ author_name }}</span>
            </h1>
            <div style="display:flex;align-items:center;gap:8px;margin:8px 0 12px;flex-wrap:wrap">
              <span class="badge" style="background:var(--grad-soft);color:var(--accent);border:1px solid var(--border)">
                <Code2 :size="13" />{{ author_role }}
              </span>
              <span class="badge" style="background:rgba(148,163,184,.14);color:var(--text-2)">
                <MapPin :size="13" />{{ authorLocation }}
              </span>
            </div>
            <p class="muted" style="font-size:14.5px;line-height:1.9">{{ site.site?.bio }}</p>
            <div style="display:flex;gap:10px;margin-top:18px;flex-wrap:wrap">
              <a class="btn btn-primary btn-sm" :href="site.site?.github" target="_blank" rel="noopener">
                <Github :size="15" />{{ (site.site?.github || 'github.com').replace('https://github.com/', '@') }}
              </a>
              <a class="btn btn-ghost btn-sm" :href="'mailto:' + (site.site?.email || '')">
                <Mail :size="15" />邮箱联系
              </a>
            </div>
          </div>
        </div>

        <!-- 座右铭 -->
        <div class="about-motto">
          <Quote :size="18" style="color:var(--accent);vertical-align:-3px;margin-right:6px" />
          {{ authorMotto }}
        </div>

        <!-- 统计卡片 -->
        <div class="about-stats reveal">
          <div v-for="s in statsCards" :key="s.label" class="about-stat glass-card card-hover">
            <component :is="s.icon" :size="20" style="color:var(--accent);margin-bottom:10px" />
            <b>{{ typeof s.value === 'number' ? s.value.toLocaleString() : s.value }}{{ s.suffix ?? '' }}</b>
            <span>{{ s.label }}</span>
          </div>
        </div>

      </div>
    </section>
  </div>
</template>
