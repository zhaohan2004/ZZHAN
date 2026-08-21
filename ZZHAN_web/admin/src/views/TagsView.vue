<script setup lang="ts">
/** 标签管理 — 1:1 复刻静态 tags.html（.admin-tools/.data-table + .modal-overlay 新增/编辑弹窗）。去 Element Plus。 */
import { computed, onMounted, reactive, ref } from 'vue'
import { Info, Pencil, Plus, RotateCcw, Search, Tag, Trash2 } from 'lucide-vue-next'
import { createTag, deleteTag, listTags, updateTag } from '@/api/admin'
import type { TagAdmin } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'

const tags = ref<TagAdmin[]>([])
const total = ref(0)
const loading = ref(false)

const query = reactive<{ page: number; pageSize: number; q: string; minCount: string; maxCount: string }>({
  page: 1,
  pageSize: 8,
  q: '',
  minCount: '',
  maxCount: '',
})

const showModal = ref(false)
const editingId = ref<number | null>(null)
const submitting = ref(false)
const form = reactive<{ name: string }>({ name: '' })

async function load() {
  loading.value = true
  try {
    const res = await listTags({
      page: query.page,
      pageSize: query.pageSize,
      q: query.q.trim() || undefined,
      minCount: query.minCount === '' ? undefined : Number(query.minCount),
      maxCount: query.maxCount === '' ? undefined : Number(query.maxCount),
    })
    tags.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载标签失败')
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
  query.q = ''
  query.minCount = ''
  query.maxCount = ''
  query.page = 1
  load()
}

function openCreate() {
  editingId.value = null
  form.name = ''
  showModal.value = true
}
function openEdit(t: TagAdmin) {
  editingId.value = t.id
  form.name = t.name
  showModal.value = true
}
function closeModal() {
  showModal.value = false
}
async function save() {
  if (!form.name.trim()) {
    toast.warning('请输入标签名称')
    return
  }
  submitting.value = true
  try {
    if (editingId.value == null) {
      await createTag(form.name.trim())
      toast.success('标签已创建')
    } else {
      await updateTag(editingId.value, form.name.trim())
      toast.success('标签已更新')
    }
    showModal.value = false
    load()
  } catch {
    toast.error('保存失败')
  } finally {
    submitting.value = false
  }
}
async function removeTag(t: TagAdmin) {
  const ok = await confirm(`确定删除标签「${t.name}」吗？使用该标签的文章将移除该标签。`, '删除标签')
  if (!ok) return
  try {
    await deleteTag(t.id)
    toast.success('已删除')
    load()
  } catch {
    toast.error('删除失败')
  }
}

onMounted(load)
</script>

<template>
  <div>
    <div class="admin-tools">
      <button class="btn btn-primary" @click="openCreate"><Plus :size="16" /> 新增标签</button>
      <div style="flex: 1"></div>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
        <input v-model="query.q" class="input" type="text" placeholder="标签名称" style="width:150px" @keyup.enter="search" />
        <input v-model="query.minCount" class="input" type="number" min="0" placeholder="次数 ≥" style="width:100px" />
        <input v-model="query.maxCount" class="input" type="number" min="0" placeholder="次数 ≤" style="width:100px" />
        <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
        <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
      </div>
      <span class="muted" style="font-size: 13px">共 <b>{{ total }}</b> 个标签</span>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>标签 ID</th>
            <th>标签名称</th>
            <th>使用次数</th>
            <th>创建时间</th>
            <th>更新时间</th>
            <th style="text-align: right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tags" :key="t.id">
            <td>{{ t.id }}</td>
            <td style="font-weight: 500; color: var(--text-2)">
              <span style="display: inline-flex; align-items: center; gap: 8px"><Tag :size="14" style="color: var(--accent)" /> {{ t.name }}</span>
            </td>
            <td style="font-family: 'JetBrains Mono', monospace">{{ t.count }}</td>
            <td>{{ t.createdAt }}</td>
            <td>{{ t.updatedAt }}</td>
            <td style="text-align: right; white-space: nowrap">
              <button class="tb-btn edit" title="编辑" @click="openEdit(t)"><Pencil :size="15" /></button>
              <button class="tb-btn del" title="删除" @click="removeTag(t)"><Trash2 :size="15" /></button>
            </td>
          </tr>
          <tr v-if="!tags.length">
            <td colspan="6" style="text-align: center; color: var(--text-3); padding: 40px">暂无标签</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="pagination" v-if="totalPages > 1">
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

    <!-- 新增 / 编辑标签 -->
    <div class="modal-overlay" :class="{ open: showModal }" @click.self="closeModal">
      <div class="modal">
        <div class="modal-head">
          <h3>{{ editingId == null ? '新增标签' : '编辑标签' }}</h3>
          <button class="modal-close" @click="closeModal"><i data-lucide="x" width="17" height="17"></i></button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">标签名称 <span class="req">*</span></label>
            <input v-model="form.name" class="input" type="text" placeholder="如：并发" maxlength="20" />
          </div>
          <div class="form-hint"><Info :size="12" style="vertical-align: -2px" /> 标签名称需唯一，支持中英文</div>
        </div>
        <div class="modal-foot">
          <button class="btn btn-ghost btn-sm" @click="closeModal">取消</button>
          <button class="btn btn-primary btn-sm" :disabled="submitting" @click="save">
            {{ submitting ? '保存中…' : '保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
