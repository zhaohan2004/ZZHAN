<script setup lang="ts">
/** 首页 — Hero + 最新文章（横向全宽卡片）+ 最近动态 + 归档 widget。结构对齐静态原型。 */
import { computed, onMounted, ref } from 'vue'
import { Activity, ArrowRight, Sparkles } from 'lucide-vue-next'
import Hero from '@/components/home/Hero.vue'
import DynamicsTimeline from '@/components/home/DynamicsTimeline.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'
import ArchiveWidget from '@/components/widget/ArchiveWidget.vue'
import { useSiteStore } from '@/stores/site'
import { getArticles } from '@/api/articles'
import type { ArticleSummary } from '@/types/models'
import { useReveal } from '@/composables/useReveal'

const root = ref<HTMLElement | null>(null)
useReveal(root)

const site = useSiteStore()
const stats = computed(() => site.stats)

const latest = ref<ArticleSummary[]>([])

onMounted(async () => {
  site.fetchStats()
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

    <!-- 最近动态 + 归档 widget -->
    <div class="section" style="padding-top:10px">
      <div class="container">
        <div class="grid gap-8 lg:grid-cols-3" style="align-items:start">
          <div class="lg:col-span-2">
            <div class="section-head reveal" style="margin-bottom:18px">
              <h2 class="section-title"><span class="st-ico"><Activity :size="17" /></span>最近动态</h2>
              <router-link to="/archive" class="section-more">文章归档 <ArrowRight :size="15" /></router-link>
            </div>
            <DynamicsTimeline :dynamics="stats?.dynamics ?? []" />
          </div>
          <div class="lg:col-span-1">
            <ArchiveWidget />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
