<script setup lang="ts">
/**
 * 侧边栏组件 — 结构对齐静态原型 main.js sidebarHTML：
 * 最新文章（缩略图 + 标题 + 日期）、分类列表、标签云。
 */
import { onMounted, ref } from 'vue'
import { FolderTree, History, Tags } from 'lucide-vue-next'
import { getArticles } from '@/api/articles'
import { getCategories, getTags } from '@/api/site'
import type { ArticleSummary, Category, Tag } from '@/types/models'
import { coverArt } from '@/utils/cover'

const props = withDefaults(defineProps<{ withCatTag?: boolean }>(), { withCatTag: true })

const recent = ref<ArticleSummary[]>([])
const cats = ref<Category[]>([])
const tags = ref<Tag[]>([])

onMounted(async () => {
  try {
    const [r, c, t] = await Promise.all([getArticles({ pageSize: 5, sort: 'latest' }), getCategories(), getTags()])
    recent.value = r.list
    cats.value = c
    tags.value = [...t].sort((a, b) => (b.count ?? 0) - (a.count ?? 0)).slice(0, 14)
  } catch {
    /* 静默 */
  }
})

function thumbOf(a: ArticleSummary): string {
  return a.coverImage || coverArt(a.title, a.category.name, a.id, 96, 72)
}
</script>

<template>
  <div class="space-y-6">
    <div class="widget reveal in">
      <h3 class="widget-title"><History :size="16" style="color:var(--accent)" />最新文章</h3>
      <ul class="widget-list">
        <li v-for="a in recent" :key="a.id" style="padding:7px 0">
          <router-link :to="`/article/${a.slug}`" style="width:100%">
            <span class="thumb-mini" :style="{ backgroundImage: 'url(' + thumbOf(a) + ')' }" />
            <span class="t">{{ a.title }}</span>
            <span class="d">{{ a.date.slice(5) }}</span>
          </router-link>
        </li>
      </ul>
    </div>

    <template v-if="withCatTag">
      <div class="widget reveal in">
        <h3 class="widget-title"><FolderTree :size="16" style="color:var(--accent)" />分类列表</h3>
        <ul class="widget-list">
          <li v-for="c in cats" :key="c.id">
            <router-link :to="`/articles?category=${c.slug}`"><i data-lucide="folder" :style="{ color: c.color }" />{{ c.name }}</router-link>
            <span class="cnt">{{ c.count ?? 0 }}</span>
          </li>
        </ul>
      </div>

      <div class="widget reveal in">
        <h3 class="widget-title"><Tags :size="16" style="color:var(--accent)" />标签云</h3>
        <div class="tag-cloud">
          <router-link v-for="t in tags" :key="t.id" :to="`/articles?tag=${encodeURIComponent(t.name)}`">{{ t.name }}</router-link>
        </div>
      </div>
    </template>
  </div>
</template>
