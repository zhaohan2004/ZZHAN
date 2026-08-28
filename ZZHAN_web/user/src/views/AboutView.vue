<script setup lang="ts">
/** 关于我 — 头像/简介/座右铭/统计卡片/技术栈/学习方向。所有可编辑项从 site 数据读取（来源：后台系统设置）。 */
import { computed, onMounted, ref } from 'vue'
import { BookOpen, Calendar, CheckCircle2, Code2, Eye, Github, Mail, MapPin, MessageCircle, Quote } from 'lucide-vue-next'
import { useSiteStore } from '@/stores/site'
import { getAbout } from '@/api/site'
import type { AboutData } from '@/types/models'
import { initialsAvatar } from '@/utils/avatar'
import { useReveal } from '@/composables/useReveal'

const root = ref<HTMLElement | null>(null)
useReveal(root)

const site = useSiteStore()
const about = ref<AboutData | null>(null)

const DIRECTIONS = ['Go 并发与调度', '分布式系统', '数据库内核', '云原生', '可观测性', 'AI 应用工程']

onMounted(async () => {
  site.fetchSite()
  site.fetchStats()
  site.fetchAbout()
  try {
    about.value = await getAbout()
  } catch {
    /* 静默 */
  }
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
  site.site?.avatar || initialsAvatar(site.site?.logo_text || author_name.value, '#6b7280', '#9ca3af', 260),
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

        <div class="grid gap-8 lg:grid-cols-2" style="align-items:start;margin-top:10px">
          <!-- 技术栈 -->
          <div class="widget reveal">
            <h3 class="widget-title">技术栈</h3>
            <div v-for="s in about?.skills ?? []" :key="s.name" class="skill-row">
              <span class="sk-name"><CheckCircle2 :size="15" style="color:var(--success)" />{{ s.name }}</span>
              <div class="skill-bar" style="flex:1"><i :style="{ width: s.level + '%' }" /></div>
              <span class="sk-val">{{ s.level }}%</span>
            </div>
          </div>

          <!-- 学习方向 -->
          <div class="widget reveal">
            <h3 class="widget-title">学习方向</h3>
            <div style="display:flex;flex-wrap:wrap;gap:10px;margin-bottom:18px">
              <span v-for="d in DIRECTIONS" :key="d" class="chip" style="cursor:default">{{ d }}</span>
            </div>
            <p class="muted" style="font-size:13.5px;line-height:1.9">当前在深入分布式缓存与消息队列，同时保持每周 3 道算法题的节奏。近期目标：啃完《Designing Data-Intensive Applications》。</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
