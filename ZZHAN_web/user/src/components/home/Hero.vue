<script setup lang="ts">
/**
 * Hero 区 — 头像 + 站点名（末段渐变）+ 标语 + bio + 按钮 + meta；
 * 右侧统计格（hero-stats）+ 终端。结构对齐静态原型 html/user/index.html。
 * 所有可编辑项均从 site 数据读取（来源：后台系统设置）。
 */
import { computed } from 'vue'
import { BookOpen, Calendar, Github, MapPin, User } from 'lucide-vue-next'
import type { SiteInfo } from '@/types/models'
import { initialsAvatar } from '@/utils/avatar'
import HeroTerminal from './HeroTerminal.vue'

const props = defineProps<{ site: SiteInfo | null }>()

const name = computed(() => props.site?.name ?? '小猫的个人博客')
const nameHead = computed(() => name.value.slice(0, Math.max(0, name.value.length - 4)))
const nameTail = computed(() => name.value.slice(-4))

/** 头像 — 优先取后台设置的自定义头像，否则按站点短名生成 initials。 */
const avatarSrc = computed(
  () => props.site?.avatar || initialsAvatar(props.site?.logo_text || props.site?.short_name || props.site?.name || 'CT', '#6b7280', '#9ca3af', 224),
)

/** 位置 — 来自 site.location。 */
const locationText = computed(() => props.site?.location || '北京')

/** 写博客第 N 年 — 来自 site.since 自动计算。 */
const yearCount = computed(() => {
  const base = props.site?.since ?? new Date().getFullYear() - 7
  const n = Math.max(1, new Date().getFullYear() - base)
  return n
})
</script>

<template>
  <section class="hero">
    <div class="hero-orb a" />
    <div class="hero-orb b" />
    <div class="container hero-inner">
      <div>
        <div class="hero-avatar-wrap">
          <img class="hero-avatar" :src="avatarSrc" alt="头像" />
          <div class="hero-avatar-ring" />
        </div>
        <div class="hero-sub"><span class="grad-text">{{ site?.tagline }}</span></div>
        <h1 class="hero-title">{{ nameHead }}<span class="grad-text">{{ nameTail }}</span></h1>
        <p class="hero-bio">{{ site?.bio }}</p>
        <div class="hero-actions">
          <router-link to="/articles" class="btn btn-primary"><BookOpen :size="17" /> 开始阅读</router-link>
          <router-link to="/about" class="btn btn-ghost"><User :size="17" /> 关于我</router-link>
        </div>
        <div class="hero-meta">
          <span class="hm"><MapPin :size="14" /> {{ locationText }}</span>
          <span class="hm">
            <Github :size="14" />
            {{ (site?.github || 'github.com/yourname').replace('https://github.com/', '') }}
          </span>
          <span class="hm"><Calendar :size="14" /> 写博客第 {{ yearCount }} 年</span>
        </div>
      </div>
      <div>
        <HeroTerminal />
      </div>
    </div>
  </section>
</template>
