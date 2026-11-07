<script setup lang="ts">
/**
 * 文章卡片 — 竖版 / 横向全宽（horizontal，图左/图右 reverse 交替）。
 * DOM 结构与静态原型 main.js articleCardHTML 完全一致。
 */
import { ArrowRight, Calendar, Eye, MessageCircle, ThumbsUp } from 'lucide-vue-next'
import type { ArticleSummary } from '@/types/models'
import { fmtDate, fmtNum } from '@/utils/format'
import TagMini from './TagMini.vue'
import ArticleCover from './ArticleCover.vue'

withDefaults(
  defineProps<{
    article: ArticleSummary
    horizontal?: boolean
    reverse?: boolean
    tagLimit?: number
  }>(),
  { horizontal: false, reverse: false, tagLimit: 3 },
)

</script>

<template>
  <article class="article-card reveal" :class="{ horizontal, reverse }">
    <router-link :to="`/article/${article.slug}`" class="ac-cover" :aria-label="article.title">
      <ArticleCover :article="article" :alt="article.title" />
    </router-link>

    <div class="ac-body">
      <div class="ac-title-wrap">
        <router-link :to="`/article/${article.slug}`" class="ac-title">{{ article.title }}</router-link>
        <span class="badge ac-cat">
          {{ article.category_name }}
        </span>
      </div>
      <p class="ac-summary">{{ article.summary }}</p>
      <div v-if="article.tags?.length" class="ac-tags">
        <TagMini v-for="t in article.tags.slice(0, tagLimit)" :key="t.id" :tag-id="t.id" :tag-name="t.name" />
      </div>
      <div class="ac-meta">
        <span><Calendar :size="13" /> {{ fmtDate(article.published_at) }}</span>
        <span><Eye :size="13" /> {{ fmtNum(article.views) }}</span>
        <span><ThumbsUp :size="13" /> {{ fmtNum(article.likes) }}</span>
        <span><MessageCircle :size="13" /> {{ fmtNum(article.comment_count) }}</span>
        <router-link
          :to="`/article/${article.slug}`"
          class="read-more"
          style="margin-left:auto;color:var(--accent);font-weight:600;display:inline-flex;align-items:center;gap:4px"
        >
          阅读 <ArrowRight :size="13" />
        </router-link>
      </div>
    </div>
  </article>
</template>
