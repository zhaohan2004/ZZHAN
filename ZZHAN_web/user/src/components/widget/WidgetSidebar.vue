<script setup lang="ts">
/**
 * 侧边栏组件 — 结构对齐静态原型 main.js sidebarHTML：
 * 最新文章（缩略图 + 标题 + 日期）、分类列表、标签云。
 */
import { computed, onMounted, ref } from 'vue'
import { FolderTree, History, Tags } from 'lucide-vue-next'
import { getArticles } from '@/api/articles'
import { getCategories, getTags } from '@/api/site'
import type { ArticleSummary, Category, Tag } from '@/types/models'
import ArticleCover from '@/components/article/ArticleCover.vue'
import { fmtDate } from '@/utils/format'

const props = withDefaults(defineProps<{
  withCatTag?: boolean
  categories?: Category[]
  tags?: Tag[]
}>(), { withCatTag: true })

const recent = ref<ArticleSummary[]>([])
const ownCategories = ref<Category[]>([])
const ownTags = ref<Tag[]>([])
const cats = computed(() => [...(props.categories ?? ownCategories.value)].sort((a, b) => (b.count ?? 0) - (a.count ?? 0)).slice(0, 14))
const tags = computed(() => [...(props.tags ?? ownTags.value)].sort((a, b) => (b.count ?? 0) - (a.count ?? 0)).slice(0, 14))

onMounted(() => {
  void getArticles({ page: 1, size: 5 }).then((r) => {
    recent.value = r.list
  }).catch(() => {})
  if (props.withCatTag && props.categories === undefined) {
    void getCategories().then((rows) => { ownCategories.value = rows }).catch(() => {})
  }
  if (props.withCatTag && props.tags === undefined) {
    void getTags().then((rows) => { ownTags.value = rows }).catch(() => {})
  }
})

</script>

<template>
  <div class="space-y-6">
    <div class="widget reveal in">
      <h3 class="widget-title"><History :size="16" style="color:var(--accent)" />最新文章</h3>
      <ul class="widget-list">
        <li v-for="a in recent" :key="a.id" style="padding:7px 0">
          <router-link :to="`/article/${a.slug}`" style="width:100%">
            <ArticleCover class="thumb-mini" :article="a" thumbnail />
            <span class="t">{{ a.title }}</span>
            <span class="d">{{ fmtDate(a.published_at) }}</span>
          </router-link>
        </li>
      </ul>
    </div>

    <template v-if="props.withCatTag">
      <div class="widget reveal in">
        <h3 class="widget-title"><FolderTree :size="16" style="color:var(--accent)" />分类列表</h3>
        <ul class="widget-list">
          <li v-for="c in cats" :key="c.id">
            <router-link :to="`/articles?category_id=${c.id}`"><i data-lucide="folder" :style="{ color: c.color }" />{{ c.name }}</router-link>
            <span class="cnt">{{ c.count ?? 0 }}</span>
          </li>
        </ul>
      </div>

      <div class="widget reveal in">
        <h3 class="widget-title"><Tags :size="16" style="color:var(--accent)" />标签云</h3>
        <div class="tag-cloud">
          <router-link v-for="t in tags" :key="t.id" :to="`/articles?tag_id=${t.id}`">{{ t.name }}</router-link>
        </div>
      </div>
    </template>
  </div>
</template>
