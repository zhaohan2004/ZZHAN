<script setup lang="ts">
/** 文章列表页 — 搜索 + 分类/标签/时间筛选 + 分页 + 侧边栏。URL query 同步（q/category/tag）。 */
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Calendar, ChevronLeft, ChevronRight, Folder, Search, Tag as TagIcon } from 'lucide-vue-next'
import ArticleCard from '@/components/article/ArticleCard.vue'
import WidgetSidebar from '@/components/widget/WidgetSidebar.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import { getArticles } from '@/api/articles'
import { getCategories, getTags } from '@/api/site'
import type { ArticleSummary, Category, Tag } from '@/types/models'

const route = useRoute()
const router = useRouter()

/** 滚动渐入观察根 */
const root = ref<HTMLElement | null>(null)

/** 手动触发 reveal 动画 */
function triggerReveal(): void {
  nextTick(() => {
    setTimeout(() => {
      if (!root.value) return
      root.value.querySelectorAll('.reveal:not(.in)').forEach((el) => {
        el.classList.add('in')
      })
    }, 100)
  })
}

const state = ref({
  page: 1,
  size: 6,
  keyword: (route.query.keyword as string) || '',
  category_id: Number(route.query.category_id) || 0,
  tag_id: Number(route.query.tag_id) || 0,
  time: 'all',
})
const cats = ref<Category[]>([])
const tags = ref<Tag[]>([])
const list = ref<ArticleSummary[]>([])
const total = ref(0)
const loading = ref(true)

const pages = computed(() => Math.max(1, Math.ceil(total.value / state.value.size)))

async function load(): Promise<void> {
  loading.value = true
  try {
    const { list: l, total: t } = await getArticles({
      page: state.value.page,
      size: state.value.size,
      category_id: state.value.category_id || undefined,
      tag_id: state.value.tag_id || undefined,
      keyword: state.value.keyword || undefined,
    })
    let rows = l
    // 时间筛选在前端按年份过滤（对齐静态页行为）
    if (state.value.time !== 'all') rows = rows.filter((a) => a.published_at.startsWith(state.value.time))
    list.value = rows
    total.value = t
  } catch {
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
    triggerReveal()
  }
}

function applyQuery(next: Record<string, string | number | undefined>, resetPage = true): void {
  const q: Record<string, string | number | undefined> = {
    keyword: state.value.keyword || undefined,
    category_id: state.value.category_id || undefined,
    tag_id: state.value.tag_id || undefined,
    page: resetPage ? 1 : state.value.page,
    ...next,
  }
  router.push({ query: q })
}

function syncFromQuery(): void {
  state.value.keyword = (route.query.keyword as string) || ''
  state.value.category_id = Number(route.query.category_id) || 0
  state.value.tag_id = Number(route.query.tag_id) || 0
  state.value.page = Number(route.query.page) || 1
  load()
}

watch(() => route.query, syncFromQuery)

function setTime(v: string): void {
  state.value.time = v
  state.value.page = 1
  load()
}

function goPage(p: number): void {
  if (p < 1 || p > pages.value) return
  applyQuery({ page: p }, false)
}

onMounted(async () => {
  try {
    const [c, t] = await Promise.all([getCategories(), getTags()])
    cats.value = c
    tags.value = [...t].sort((a, b) => (b.count ?? 0) - (a.count ?? 0))
  } catch {
    /* 静默 */
  }
  load()
})
</script>

<template>
  <div ref="root">
    <!-- 页面头部 -->
    <section style="padding:118px 0 30px">
      <div class="container">
        <div class="anim-fade" style="display:flex;align-items:center;gap:10px;font-size:13px;color:var(--text-3);margin-bottom:12px">
          <router-link to="/" style="color:var(--text-2)">首页</router-link>
          <ChevronRight :size="13" />
          <span class="grad-text" style="font-weight:600">文章列表</span>
        </div>
        <h1 style="font-size:30px;font-weight:800;letter-spacing:-.4px">全部文章 <span class="grad-text" style="font-size:22px">共 {{ total }} 篇</span></h1>
        <p class="muted" style="margin-top:8px;font-size:14px">在这里沉淀 Go、MySQL、Redis 与后端工程的每一个细节。</p>
      </div>
    </section>

    <!-- 筛选区 -->
    <section style="padding:8px 0 26px">
      <div class="container">
        <div class="widget" style="padding:18px 20px">
          <div style="display:flex;gap:10px;flex-wrap:wrap;margin-bottom:14px">
            <input v-model="state.keyword" class="input" type="search" placeholder="搜索标题 / 摘要 / 标签…" style="max-width:320px" @keyup.enter="applyQuery({ keyword: state.keyword })" />
            <button class="btn btn-primary btn-sm" type="button" @click="applyQuery({ keyword: state.keyword })"><Search :size="15" />搜索</button>
          </div>
          <div style="font-size:12.5px;color:var(--text-3);margin-bottom:8px;display:flex;align-items:center;gap:6px"><Folder :size="13" />分类筛选</div>
          <div class="flex gap-2 flex-wrap" style="margin-bottom:14px">
            <button class="chip" :class="{ active: state.category_id === 0 }" type="button" @click="applyQuery({ category_id: undefined })">全部</button>
            <button
              v-for="c in cats"
              :key="c.id"
              class="chip"
              :class="{ active: state.category_id === c.id }"
              type="button"
              @click="applyQuery({ category_id: c.id })"
            >
              {{ c.name }}
            </button>
          </div>
          <div style="font-size:12.5px;color:var(--text-3);margin-bottom:8px;display:flex;align-items:center;gap:6px"><TagIcon :size="13" />标签筛选</div>
          <div class="flex gap-2 flex-wrap" style="margin-bottom:14px">
            <button class="chip" :class="{ active: state.tag_id === 0 }" type="button" @click="applyQuery({ tag_id: undefined })">全部</button>
            <button
              v-for="t in tags"
              :key="t.id"
              class="chip"
              :class="{ active: state.tag_id === t.id }"
              type="button"
              @click="applyQuery({ tag_id: t.id })"
            >
              {{ t.name }}
            </button>
          </div>
          <div style="font-size:12.5px;color:var(--text-3);margin-bottom:8px;display:flex;align-items:center;gap:6px"><Calendar :size="13" />时间筛选</div>
          <div class="flex gap-2 flex-wrap">
            <button class="chip" :class="{ active: state.time === 'all' }" type="button" @click="setTime('all')">全部时间</button>
            <button v-for="y in ['2026', '2025', '2024']" :key="y" class="chip" :class="{ active: state.time === y }" type="button" @click="setTime(y)">{{ y }} 年</button>
          </div>
        </div>
      </div>
    </section>

    <!-- 内容区 -->
    <section style="padding-bottom:60px">
      <div class="container">
        <div class="grid gap-8 lg:grid-cols-3" style="align-items:start">
          <div class="lg:col-span-2" style="min-width:0">
            <div v-if="loading" class="grid gap-6 md:grid-cols-2">
              <div v-for="n in 4" :key="n" class="glass-card overflow-hidden"><div class="h-40 bg-[var(--card-2)]" /><div class="space-y-3 p-5"><div class="h-4 w-3/4 rounded bg-[var(--card-2)]" /><div class="h-4 w-full rounded bg-[var(--card-2)]" /></div></div>
            </div>
            <template v-else>
              <div v-if="list.length" class="grid gap-6 md:grid-cols-1">
                <ArticleCard v-for="(a, i) in list" :key="a.id" :article="a" horizontal :reverse="i % 2 === 1" />
              </div>
              <EmptyState v-else icon="search-x" text="没有找到匹配的文章" sub="换个关键词或筛选条件试试" />
            </template>

            <div v-if="pages > 1" class="pagination">
              <button class="page-btn" :class="{ disabled: state.page <= 1 }" type="button" @click="goPage(state.page - 1)"><ChevronLeft :size="15" /></button>
              <button
                v-for="p in pages"
                :key="p"
                class="page-btn"
                :class="{ active: p === state.page }"
                type="button"
                @click="goPage(p)"
              >
                {{ p }}
              </button>
              <button class="page-btn" :class="{ disabled: state.page >= pages }" type="button" @click="goPage(state.page + 1)"><ChevronRight :size="15" /></button>
            </div>
          </div>
          <aside class="lg:col-span-1" style="position:sticky;top:88px;min-width:0">
            <WidgetSidebar />
          </aside>
        </div>
      </div>
    </section>
  </div>
</template>
