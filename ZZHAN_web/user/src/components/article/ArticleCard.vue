<script setup lang="ts">
/**
 * 文章卡片 — 竖版 / 横向全宽（horizontal，图左/图右 reverse 交替）。
 * DOM 结构与静态原型 main.js articleCardHTML 完全一致。
 */
import { computed } from 'vue'
import { ArrowRight, Calendar, Eye, MessageCircle, ThumbsUp } from 'lucide-vue-next'
import type { ArticleSummary } from '@/types/models'
import { coverArt } from '@/utils/cover'
import { fmtNum } from '@/utils/format'
import TagMini from './TagMini.vue'

const props = withDefaults(
  defineProps<{
    article: ArticleSummary
    horizontal?: boolean
    reverse?: boolean
    tagLimit?: number
  }>(),
  { horizontal: false, reverse: false, tagLimit: 3 },
)

const cover = computed(() =>
  props.article.cover_image || coverArt(props.article.title, props.article.category_name, props.article.id),
)
const catColor = computed(() => '#4a8eff') // 默认主题色
</script>

<template>
  <article class="article-card reveal" :class="{ horizontal, reverse }">
    <router-link :to="`/article/${article.slug}`" class="ac-cover" :aria-label="article.title">
      <img :src="cover" :alt="article.title" loading="lazy" />
    </router-link>

    <div class="ac-body">
      <div class="ac-title-wrap">
        <router-link :to="`/article/${article.slug}`" class="ac-title">{{ article.title }}</router-link>
        <span class="badge ac-cat" :style="{ background: catColor + '1f', color: catColor, border: '1px solid ' + catColor + '44' }">
          {{ article.category_name }}
        </span>
      </div>
      <p class="ac-summary">{{ article.summary }}</p>
      <div v-if="article.tags?.length" class="ac-tags">
        <TagMini v-for="t in article.tags.slice(0, tagLimit)" :key="t.id" :tag-id="t.id" :tag-name="t.name" />
      </div>
      <div class="ac-meta">
        <span><Calendar :size="13" /> {{ article.published_at }}</span>
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
