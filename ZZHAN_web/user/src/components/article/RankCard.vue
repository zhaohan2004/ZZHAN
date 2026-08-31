<script setup lang="ts">
/** 热门排行卡片 — 排名数字（前 3 高亮）+ 标题 + 阅读数 + 评论数。 */
import { Eye, MessageCircle } from 'lucide-vue-next'
import type { ArticleSummary } from '@/types/models'
import { fmtDate, fmtNum } from '@/utils/format'

defineProps<{ article: ArticleSummary; rank: number }>()
</script>

<template>
  <div class="rank-card reveal">
    <div class="rank-no">{{ rank }}</div>
    <div class="rank-info">
      <router-link class="rank-title" :to="`/article/${article.slug}`">{{ article.title }}</router-link>
      <div class="rank-meta"><Eye :size="12" /> {{ fmtNum(article.views) }} 阅读 · <MessageCircle :size="12" /> {{ fmtNum(article.comment_count) }} 评论 · {{ fmtDate(article.published_at) }}</div>
    </div>
  </div>
</template>
