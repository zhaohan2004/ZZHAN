<script setup lang="ts">
/** Hero 统计格 — 原创文章 / 总浏览量 / 分类数 / 标签数。 */
import { onMounted, ref } from 'vue'
import { getCategories, getTags } from '@/api/site'
import { fmtNum } from '@/utils/format'

const props = defineProps<{ articles: number; views: number }>()

const catCount = ref(0)
const tagCount = ref(0)

onMounted(async () => {
  try {
    const [c, t] = await Promise.all([getCategories(), getTags()])
    catCount.value = c.length
    tagCount.value = t.length
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div class="hero-stats">
    <div class="hero-stat"><b>{{ articles }}</b><span>原创文章</span></div>
    <div class="hero-stat"><b>{{ fmtNum(views) }}</b><span>总浏览量</span></div>
    <div class="hero-stat"><b>{{ catCount }}</b><span>技术分类</span></div>
    <div class="hero-stat"><b>{{ tagCount }}</b><span>标签</span></div>
  </div>
</template>
