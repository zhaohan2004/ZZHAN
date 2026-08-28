<script setup lang="ts">
/** 文章管理 — 1:1 复刻静态 articles.html（.admin-tools/.data-table/.pagination）。去 Element Plus。 */
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowDownToLine, Download, Pencil, Plus, RotateCcw, Search, Send, Trash2 } from 'lucide-vue-next'
import {
  deleteAdminArticle,
  listAdminArticles,
  listCategories,
  listTags,
  setArticleStatus,
} from '@/api/admin'
import type { AdminArticle, ArticleStatus, CategoryAdmin, TagAdmin } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'
import { fmtNum } from '@/utils/format'

const router = useRouter()

const articles = ref<AdminArticle[]>([])
const categories = ref<CategoryAdmin[]>([])
const tags = ref<TagAdmin[]>([])
const total = ref(0)
const loading = ref(false)

interface ArticleQueryState {
  page: number
  pageSize: number
  keyword: string
  category: string
  tag: string
  status: ArticleStatus | 'all'
}
const query = reactive<ArticleQueryState>({ page: 1, pageSize: 8, keyword: '', category: 'all', tag: 'all', status: 'all' })

const STATUS_MAP: Record<ArticleStatus, { cls: string; label: string; next: ArticleStatus; nextLabel: string }> = {
  published: { cls: 'st-pub', label: '已发布', next: 'down', nextLabel: '下架' },
  draft: { cls: 'st-draft', label: '草稿', next: 'published', nextLabel: '发布' },
  down: { cls: 'st-down', label: '已下架', next: 'published', nextLabel: '发布' },
}

async function load() {
  loading.value = true
  try {
    const res = await listAdminArticles({
      page: query.page,
      pageSize: query.pageSize,
      keyword: query.keyword || undefined,
      category: query.category !== 'all' ? query.category : undefined,
      tag: query.tag !== 'all' ? query.tag : undefined,
      status: query.status !== 'all' ? query.status : undefined,
    })
    articles.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载文章失败')
  } finally {
    loading.value = false
  }
}

const totalPages = () => Math.max(1, Math.ceil(total.value / (query.pageSize || 8)))

function goPage(p: number) {
  if (p < 1 || p > totalPages()) return
  query.page = p
  load()
}

function search() {
  query.page = 1
  load()
}
function reset() {
  query.keyword = ''
  query.category = 'all'
  query.tag = 'all'
  query.status = 'all'
  query.page = 1
  load()
}

function newArticle() {
  router.push('/editor')
}
function editArticle(a: AdminArticle) {
  router.push('/editor?id=' + a.id)
}
async function toggleStatus(a: AdminArticle) {
  const map = STATUS_MAP[a.status]
  const ok = await confirm(`确定将《${a.title}》${map.nextLabel}吗？`, '切换状态')
  if (!ok) return
  try {
    await setArticleStatus(a.id, map.next)
    toast.success(`已${map.nextLabel}`)
    load()
  } catch {
    toast.error('操作失败')
  }
}
async function removeArticle(a: AdminArticle) {
  const ok = await confirm(`确定删除《${a.title}》吗？此操作不可恢复。`, '删除文章')
  if (!ok) return
  try {
    await deleteAdminArticle(a.id)
    toast.success('已删除')
    if (articles.value.length === 1 && query.page > 1) query.page--
    load()
  } catch {
    toast.error('删除失败')
  }
}
function exportArticles() {
  toast.info('演示环境暂不支持导出')
}

onMounted(async () => {
  try {
    const [catRes, tagRes] = await Promise.all([
      listCategories({ page: 1, pageSize: 100 }),
      listTags({ page: 1, pageSize: 100 }),
    ])
    categories.value = catRes.list
    tags.value = tagRes.list
  } catch {
    /* 静默 */
  }
  load()
})
</script>

<template>
  <div>
    <div class="admin-tools">
      <button class="btn btn-primary" @click="newArticle"><Plus :size="16" /> 新建文章</button>
      <button class="btn btn-ghost" @click="exportArticles"><Download :size="16" /> 导出</button>
      <div style="flex: 1"></div>
      <input
        v-model="query.keyword"
        class="input"
        type="search"
        placeholder="搜索文章标题…"
        style="width: 220px"
        @keyup.enter="search"
      />
      <select v-model="query.category" class="select">
        <option value="all">全部分类</option>
        <option v-for="c in categories" :key="c.id" :value="c.name">{{ c.name }}</option>
      </select>
      <select v-model="query.tag" class="select">
        <option value="all">全部标签</option>
        <option v-for="t in tags" :key="t.id" :value="t.name">{{ t.name }}</option>
      </select>
      <select v-model="query.status" class="select">
        <option value="all">全部状态</option>
        <option value="published">已发布</option>
        <option value="draft">草稿</option>
        <option value="down">已下架</option>
      </select>
      <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
      <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>文章标题</th>
            <th>分类</th>
            <th>标签</th>
            <th>状态</th>
            <th>阅读量</th>
            <th>发布时间</th>
            <th>更新时间</th>
            <th style="text-align: right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="a in articles" :key="a.id">
            <td>{{ a.id }}</td>
            <td style="font-weight: 500; color: var(--text-2)">{{ a.title }}</td>
            <td>{{ a.category }}</td>
            <td>
              <span
                v-for="t in a.tags"
                :key="t"
                style="
                  display: inline-flex;
                  align-items: center;
                  padding: 2px 8px;
                  margin: 2px 4px 2px 0;
                  border-radius: 6px;
                  font-size: 12px;
                  background: var(--glass);
                  color: var(--text-2);
                  border: 1px solid var(--border);
                "
                >{{ t }}</span
              >
              <span v-if="!a.tags.length" style="color: var(--text-3)">—</span>
            </td>
            <td>
              <span class="st" :class="STATUS_MAP[a.status].cls">{{ STATUS_MAP[a.status].label }}</span>
            </td>
            <td style="font-family: 'JetBrains Mono', monospace">{{ fmtNum(a.views) }}</td>
            <td>{{ a.date }}</td>
            <td>{{ a.updated }}</td>
            <td style="text-align: right; white-space: nowrap">
              <button class="tb-btn edit" title="编辑" @click="editArticle(a)"><Pencil :size="15" /></button>
              <button class="tb-btn view" :title="STATUS_MAP[a.status].nextLabel" @click="toggleStatus(a)">
                <Send v-if="STATUS_MAP[a.status].nextLabel === '发布'" :size="15" />
                <ArrowDownToLine v-else :size="15" />
              </button>
              <button class="tb-btn del" title="删除" @click="removeArticle(a)"><Trash2 :size="15" /></button>
            </td>
          </tr>
          <tr v-if="!articles.length">
            <td colspan="9" style="text-align: center; color: var(--text-3); padding: 40px">暂无文章</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="totalPages() > 1">
      <button class="page-btn" :class="{ disabled: query.page <= 1 }" @click="goPage(query.page - 1)">上一页</button>
      <button
        v-for="p in totalPages()"
        :key="p"
        class="page-btn"
        :class="{ active: p === query.page }"
        @click="goPage(p)"
      >
        {{ p }}
      </button>
      <button class="page-btn" :class="{ disabled: query.page >= totalPages() }" @click="goPage(query.page + 1)">
        下一页
      </button>
    </div>
  </div>
</template>
