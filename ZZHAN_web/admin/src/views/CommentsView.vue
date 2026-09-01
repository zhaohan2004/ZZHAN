<script setup lang="ts">
/** 评论管理 — 1:1 复刻静态 comments.html（.admin-tools/.data-table/.pagination）。去 Element Plus。 */
import { computed, onMounted, reactive, ref } from 'vue'
import { Ban, RotateCcw, Search, ShieldCheck, Trash2 } from 'lucide-vue-next'
import { deleteComment, listAdminArticles, listComments, updateCommentStatus } from '@/api/admin'
import type { AdminArticle, CommentAdmin, CommentStatus } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'

const comments = ref<CommentAdmin[]>([])
const articles = ref<AdminArticle[]>([])
const total = ref(0)
const loading = ref(false)

const query = reactive<{
  page: number
  pageSize: number
  status: CommentStatus | 'all'
  article_id: number
  startDate: string
  endDate: string
}>({
  page: 1,
  pageSize: 10,
  status: 'all',
  article_id: 0,
  startDate: '',
  endDate: '',
})

const STATUS_MAP: Record<CommentStatus, { cls: string; label: string }> = {
  normal: { cls: 'st-ok', label: '正常' },
  banned: { cls: 'st-block', label: '封禁' },
}

async function load() {
  loading.value = true
  try {
    const res = await listComments({
      page: query.page,
      pageSize: query.pageSize,
      status: query.status,
      article_id: query.article_id || undefined,
      startDate: query.startDate || undefined,
      endDate: query.endDate || undefined,
    })
    comments.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载评论失败')
  } finally {
    loading.value = false
  }
}

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / query.pageSize)))

function goPage(p: number) {
  if (p < 1 || p > totalPages.value) return
  query.page = p
  load()
}
function search() {
  query.page = 1
  load()
}
function reset() {
  query.article_id = 0
  query.startDate = ''
  query.endDate = ''
  query.status = 'all'
  query.page = 1
  load()
}

async function setStatus(c: CommentAdmin, status: CommentStatus, label: string) {
  const ok = await confirm(`确定将这条评论${label}吗？`, '审核评论')
  if (!ok) return
  try {
    await updateCommentStatus(c.id, status)
    toast.success(`已${label}`)
    load()
  } catch {
    toast.error('操作失败')
  }
}
async function removeComment(c: CommentAdmin) {
  const ok = await confirm('确定删除这条评论吗？此操作不可恢复。', '删除评论')
  if (!ok) return
  try {
    await deleteComment(c.id)
    toast.success('已删除')
    if (comments.value.length === 1 && query.page > 1) query.page--
    load()
  } catch {
    toast.error('删除失败')
  }
}

onMounted(async () => {
  try {
    const res = await listAdminArticles({ page: 1, pageSize: 100 })
    articles.value = res.list
  } catch {
    /* 静默 */
  }
  load()
})
</script>

<template>
  <div>
    <div class="admin-tools">
      <div style="flex: 1"></div>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
        <select v-model.number="query.article_id" class="select" style="width:210px">
          <option :value="0">全部文章</option>
          <option v-for="a in articles" :key="a.id" :value="a.id">{{ a.title }}</option>
        </select>
        <input v-model="query.startDate" class="input" type="date" style="width:150px" />
        <span class="muted" style="font-size:13px">至</span>
        <input v-model="query.endDate" class="input" type="date" style="width:150px" />
        <select v-model="query.status" class="select" style="width:110px">
          <option value="all">全部状态</option>
          <option value="normal">正常</option>
          <option value="banned">封禁</option>
        </select>
        <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
        <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
      </div>
      <span class="muted" style="font-size: 13px">共 <b>{{ total }}</b> 条评论</span>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>评论用户</th>
            <th>头像</th>
            <th>评论内容</th>
            <th>所属文章</th>
            <th>点赞</th>
            <th>发布时间</th>
            <th>状态</th>
            <th style="text-align: right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in comments" :key="c.id">
            <td style="font-weight: 600; color: var(--text-2)">{{ c.user_name }}</td>
            <td>
              <img
                v-if="c.user_avatar"
                :src="c.user_avatar"
                :alt="c.user_name"
                style="width: 32px; height: 32px; border-radius: 8px; object-fit: cover; border: 1px solid var(--border)"
              />
              <span
                v-else
                style="display: inline-flex; align-items: center; justify-content: center; width: 32px; height: 32px; border-radius: 8px; font-size: 12px; font-weight: 600; color: #fff; background: var(--accent)"
              >{{ c.user_name?.[0] }}</span>
            </td>
            <td style="max-width: 320px; color: var(--text-2)">{{ c.content }}</td>
            <td style="color: var(--text-2)">{{ c.article_title }}</td>
            <td style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px">{{ c.like_count }}</td>
            <td style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px">{{ c.created_at }}</td>
            <td><span class="st" :class="STATUS_MAP[c.status].cls">{{ STATUS_MAP[c.status].label }}</span></td>
            <td style="text-align: right; white-space: nowrap">
              <button
                v-if="c.status === 'normal'"
                class="tb-btn edit"
                title="封禁"
                @click="setStatus(c, 'banned', '封禁')"
              >
                <Ban :size="15" />
              </button>
              <button
                v-else
                class="tb-btn view"
                title="恢复正常"
                @click="setStatus(c, 'normal', '恢复正常')"
              >
                <ShieldCheck :size="15" />
              </button>
              <button class="tb-btn del" title="删除" @click="removeComment(c)"><Trash2 :size="15" /></button>
            </td>
          </tr>
          <tr v-if="!comments.length">
            <td colspan="8" style="text-align: center; color: var(--text-3); padding: 40px">暂无评论</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="total > 0">
      <button class="page-btn" :class="{ disabled: query.page <= 1 }" @click="goPage(query.page - 1)">上一页</button>
      <button
        v-for="p in totalPages"
        :key="p"
        class="page-btn"
        :class="{ active: p === query.page }"
        @click="goPage(p)"
      >
        {{ p }}
      </button>
      <button class="page-btn" :class="{ disabled: query.page >= totalPages }" @click="goPage(query.page + 1)">
        下一页
      </button>
    </div>
  </div>
</template>
