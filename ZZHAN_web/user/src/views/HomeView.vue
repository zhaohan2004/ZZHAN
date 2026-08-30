<script setup lang="ts">
/** 首页 — Hero + 最新文章（横向全宽卡片） */
import { onMounted, ref } from 'vue'
import { ArrowRight, Sparkles } from 'lucide-vue-next'
import Hero from '@/components/home/Hero.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'
import { useSiteStore } from '@/stores/site'
import { getArticles } from '@/api/articles'
import type { ArticleSummary } from '@/types/models'
import { useReveal } from '@/composables/useReveal'

const root = ref<HTMLElement | null>(null)
useReveal(root)

const site = useSiteStore()

const latest = ref<ArticleSummary[]>([])

onMounted(async () => {
  try {
    const l = await getArticles({ page: 1, size: 4, sort: 'latest' })
    latest.value = l.list
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div ref="root">
    <Hero :site="site.site" />

    <!-- 最新文章（横向全宽卡片，图左/图右交替） -->
    <div class="section">
      <div class="container">
        <div class="section-head reveal">
          <h2 class="section-title"><span class="st-ico"><Sparkles :size="17" /></span>最新文章</h2>
          <router-link to="/articles" class="section-more">查看全部 <ArrowRight :size="15" /></router-link>
        </div>
        <div id="latestGrid" class="grid gap-6 md:grid-cols-1">
          <ArticleCard v-for="(a, i) in latest" :key="a.id" :article="a" horizontal :reverse="i % 2 === 1" />
        </div>
      </div>
    </div>

  </div>
</template>
