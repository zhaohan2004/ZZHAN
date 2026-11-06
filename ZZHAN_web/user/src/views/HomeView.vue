<script setup lang="ts">
/** 首页 — Hero + 最新文章（横向全宽卡片） */
import { onMounted, ref } from 'vue'
import { ArrowRight, BookOpen } from 'lucide-vue-next'
import Hero from '@/components/home/Hero.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import Skeleton from '@/components/ui/Skeleton.vue'
import { useSiteStore } from '@/stores/site'
import { getArticles } from '@/api/articles'
import type { ArticleSummary } from '@/types/models'

const site = useSiteStore()

const latest = ref<ArticleSummary[]>([])
const loading = ref(true)
const loadFailed = ref(false)

async function loadLatest(): Promise<void> {
  loading.value = true
  loadFailed.value = false
  try {
    const l = await getArticles({ page: 1, size: 4, sort: 'latest' })
    latest.value = l.list
  } catch {
    loadFailed.value = true
  } finally {
    loading.value = false
  }
}

onMounted(loadLatest)
</script>

<template>
  <div>
    <Hero :site="site.site" />

    <!-- 最新文章（横向全宽卡片，图左/图右交替） -->
    <div class="section">
      <div class="container">
        <div class="section-head">
          <h2 class="section-title"><span class="st-ico"><BookOpen :size="17" aria-hidden="true" /></span>最新文章</h2>
          <router-link to="/articles" class="section-more">查看全部 <ArrowRight :size="15" /></router-link>
        </div>
        <div v-if="loading" class="glass-card p-5" role="status">
          <span class="sr-only">正在加载最新文章</span>
          <Skeleton aria-hidden="true" type="card" />
        </div>
        <div v-else-if="loadFailed" class="glass-card pb-6 text-center">
          <EmptyState icon="file-x" text="最新文章加载失败" />
          <button class="btn btn-ghost" type="button" @click="loadLatest">重试</button>
        </div>
        <div v-else-if="latest.length" id="latestGrid" class="grid gap-6 md:grid-cols-1">
          <ArticleCard v-for="(a, i) in latest" :key="a.id" :article="a" class="in" horizontal :reverse="i % 2 === 1" />
        </div>
        <EmptyState v-else icon="inbox" text="暂无已发布文章" />
      </div>
    </div>

  </div>
</template>
