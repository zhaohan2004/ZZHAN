<script setup lang="ts">
/** 分类管理 — 1:1 复刻静态 categories.html（.admin-tools/.data-table + .modal-overlay 新增/编辑弹窗）。去 Element Plus。 */
import { computed, onMounted, reactive, ref } from 'vue'
import { Info, Pencil, Plus, RotateCcw, Search, Trash2 } from 'lucide-vue-next'
import { createCategory, deleteCategory, listCategories, updateCategory } from '@/api/admin'
import type { CategoryAdmin } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'

const categories = ref<CategoryAdmin[]>([])
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
const form = reactive<{ name: string; desc: string }>({ name: '', desc: '' })

async function load() {
  loading.value = true
  try {
    const res = await listCategories({
      page: query.page,
      pageSize: query.pageSize,
      q: query.q.trim() || undefined,
      minCount: query.minCount === '' ? undefined : Number(query.minCount),
      maxCount: query.maxCount === '' ? undefined : Number(query.maxCount),
    })
    categories.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载分类失败')
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
  form.desc = ''
  showModal.value = true
}
function openEdit(c: CategoryAdmin) {
  editingId.value = c.id
  form.name = c.name
  form.desc = c.desc
  showModal.value = true
}
function closeModal() {
  showModal.value = false
}
async function save() {
  if (!form.name.trim()) {
    toast.warning('请输入分类名称')
    return
  }
  submitting.value = true
  try {
    if (editingId.value == null) {
      await createCategory({ name: form.name.trim(), slug: '', desc: form.desc.trim(), color: '#3b82f6' })
      toast.success('分类已创建')
    } else {
      await updateCategory(editingId.value, { name: form.name.trim(), desc: form.desc.trim() })
      toast.success('分类已更新')
    }
    showModal.value = false
    load()
  } catch {
    toast.error('保存失败')
  } finally {
    submitting.value = false
  }
}
async function removeCategory(c: CategoryAdmin) {
  const ok = await confirm(`确定删除分类「${c.name}」吗？该分类下的文章将变为未分类。`, '删除分类')
  if (!ok) return
  try {
    await deleteCategory(c.id)
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
      <button class="btn btn-primary" @click="openCreate"><Plus :size="16" /> 新增分类</button>
      <div style="flex: 1"></div>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
        <input v-model="query.q" class="input" type="text" placeholder="分类名称" style="width:150px" @keyup.enter="search" />
        <input v-model="query.minCount" class="input" type="number" min="0" placeholder="文章数 ≥" style="width:100px" />
        <input v-model="query.maxCount" class="input" type="number" min="0" placeholder="文章数 ≤" style="width:100px" />
        <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
        <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
      </div>
      <span class="muted" style="font-size: 13px">共 <b>{{ total }}</b> 个分类</span>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>分类名称</th>
            <th>分类描述</th>
            <th>文章数量</th>
            <th>创建时间</th>
            <th>更新时间</th>
            <th style="text-align: right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in categories" :key="c.id">
            <td>{{ c.id }}</td>
            <td style="font-weight: 600; color: var(--text-1)">
              <span style="display: inline-flex; align-items: center; gap: 8px">
                <span style="width: 10px; height: 10px; border-radius: 50%; flex: none" :style="{ background: c.color }"></span>
                {{ c.name }}
              </span>
            </td>
            <td style="color: var(--text-2); max-width: 280px">{{ c.desc || '—' }}</td>
            <td style="font-family: 'JetBrains Mono', monospace">{{ c.count }}</td>
            <td>{{ c.createdAt }}</td>
            <td>{{ c.updatedAt }}</td>
            <td style="text-align: right; white-space: nowrap">
              <button class="tb-btn edit" title="编辑" @click="openEdit(c)"><Pencil :size="15" /></button>
              <button class="tb-btn del" title="删除" @click="removeCategory(c)"><Trash2 :size="15" /></button>
            </td>
          </tr>
          <tr v-if="!categories.length">
            <td colspan="7" style="text-align: center; color: var(--text-3); padding: 40px">暂无分类</td>
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

    <!-- 新增 / 编辑分类 -->
    <div class="modal-overlay" :class="{ open: showModal }" @click.self="closeModal">
      <div class="modal">
        <div class="modal-head">
          <h3>{{ editingId == null ? '新增分类' : '编辑分类' }}</h3>
          <button class="modal-close" @click="closeModal"><i data-lucide="x" width="17" height="17"></i></button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">分类名称 <span class="req">*</span></label>
            <input v-model="form.name" class="input" type="text" placeholder="如：Go" />
          </div>
          <div class="form-group">
            <label class="form-label">分类描述</label>
            <textarea v-model="form.desc" class="textarea" rows="3" placeholder="一句话介绍该分类"></textarea>
          </div>
          <div class="form-hint"><Info :size="12" style="vertical-align: -2px" /> 图标与颜色在演示环境中由系统自动分配</div>
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
