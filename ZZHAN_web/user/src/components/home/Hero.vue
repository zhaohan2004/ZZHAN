<script setup lang="ts">
/** 首页站点介绍与作者手记。 */
import { computed } from 'vue'
import { BookOpen, Calendar, Github, Link, MapPin, Twitter, User } from 'lucide-vue-next'
import type { SiteInfo } from '@/types/models'
import HeroTerminal from './HeroTerminal.vue'

const props = defineProps<{ site: SiteInfo | null }>()

const name = computed(() => props.site?.name ?? '小猫的个人博客')

/** 位置 — 来自 site.location。 */
const locationText = computed(() => props.site?.location || '北京')

/** 写博客第 N 年 — 来自 site.since 自动计算。 */
const yearCount = computed(() => {
  const base = props.site?.since ?? new Date().getFullYear() - 7
  const n = Math.max(1, new Date().getFullYear() - base)
  return n
})

/** 社交图标映射 */
const SOCIAL_ICON_MAP: Record<string, typeof Github> = {
  github: Github,
  twitter: Twitter,
  x: Twitter,
}

function getSocialIcon(icon: string) {
  return SOCIAL_ICON_MAP[icon?.toLowerCase()] || Link
}
</script>

<template>
  <section class="hero">
    <div class="container hero-inner">
      <div>
        <div v-if="site?.avatar || site?.tagline" class="hero-identity">
          <img v-if="site.avatar" class="hero-avatar" :src="site.avatar" alt="" />
          <p v-if="site.tagline" class="hero-sub">{{ site.tagline }}</p>
        </div>
        <h1 class="hero-title">{{ name }}</h1>
        <p v-if="site?.bio" class="hero-bio">{{ site.bio }}</p>
        <div class="hero-actions">
          <router-link to="/articles" class="btn btn-primary"><BookOpen :size="17" /> 开始阅读</router-link>
          <router-link to="/about" class="btn btn-ghost"><User :size="17" /> 关于我</router-link>
        </div>
        <div class="hero-meta">
          <span class="hm"><MapPin :size="14" /> {{ locationText }}</span>
          <span class="hm"><Calendar :size="14" /> 写博客第 {{ yearCount }} 年</span>
          <template v-if="site?.socials?.length">
            <a
              v-for="s in site.socials"
              :key="s.name"
              class="hm social-link"
              :href="s.url"
              target="_blank"
              rel="noopener"
              :title="s.name"
            >
              <component :is="getSocialIcon(s.icon)" :size="14" />
              {{ s.name }}
            </a>
          </template>
        </div>
      </div>
      <HeroTerminal />
    </div>
  </section>
</template>

<style scoped>
.social-link {
  text-decoration: none;
  transition: color 0.2s ease;
}
.social-link:hover {
  color: var(--accent);
}
</style>
