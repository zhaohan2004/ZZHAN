<script setup lang="ts">
/** 操作日志 — 列表 + 筛选（操作类型、操作对象、时间段）。 */
import { computed, onMounted, reactive, ref } from 'vue'
import { RotateCcw, Search } from 'lucide-vue-next'
import { listOperationLogs } from '@/api/operation-logs'
import type { OperationLogAdmin } from '@/types/models'
import { toast } from '@/composables/useToast'

const logs = ref<OperationLogAdmin[]>([])
const total = ref(0)
const loading = ref(false)

const query = reactive<{
  page: number
  pageSize: number
  action: string
  target: string
  startDate: string
  endDate: string
}>({
  page: 1,
  pageSize: 10,
  action: '',
  target: '',
  startDate: '',
  endDate: '',
})

/** 操作类型选项 */
const ACTION_OPTIONS = [
  { value: '', label: '全部类型' },
  { value: '新建', label: '新建' },
  { value: '更新', label: '更新' },
  { value: '删除', label: '删除' },
  { value: '发布', label: '发布' },
  { value: '存为草稿', label: '存为草稿' },
  { value: '下架', label: '下架' },
  { value: '启用', label: '启用' },
  { value: '禁用', label: '禁用' },
  { value: '解封', label: '解封' },
  { value: '封禁', label: '封禁' },
]

/** 操作对象选项（对应中间件 resourceNames） */
const TARGET_OPTIONS = [
  { value: '', label: '全部对象' },
  { value: '文章', label: '文章' },
  { value: '分类', label: '分类' },
  { value: '标签', label: '标签' },
  { value: '评论', label: '评论' },
  { value: '系统设置', label: '系统设置' },
]

async function load() {
  loading.value = true
  try {
    const res = await listOperationLogs({
      page: query.page,
      pageSize: query.pageSize,
      action: query.action || undefined,
      target: query.target || undefined,
      startDate: query.startDate || undefined,
      endDate: query.endDate || undefined,
    })
    logs.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载操作日志失败')
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
  query.action = ''
  query.target = ''
  query.startDate = ''
  query.endDate = ''
  query.page = 1
  load()
}

onMounted(() => {
  load()
})
</script>

<template>
  <div>
    <div class="admin-tools">
      <div style="flex: 1"></div>
      <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
        <select v-model="query.action" class="select" style="width:130px">
          <option v-for="opt in ACTION_OPTIONS" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <select v-model="query.target" class="select" style="width:130px">
          <option v-for="opt in TARGET_OPTIONS" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <input v-model="query.startDate" class="input" type="date" style="width:150px" />
        <span class="muted" style="font-size:13px">至</span>
        <input v-model="query.endDate" class="input" type="date" style="width:150px" />
        <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
        <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
      </div>
      <span class="muted" style="font-size: 13px">共 <b>{{ total }}</b> 条记录</span>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>操作者</th>
            <th>操作类型</th>
            <th>操作对象</th>
            <th>操作时间</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="log in logs" :key="log.id">
            <td style="font-weight: 600; color: var(--text-2)">{{ log.admin_name || '-' }}</td>
            <td><span class="st st-ok">{{ log.action }}</span></td>
            <td style="color: var(--text-2)">{{ log.target }}</td>
            <td style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px">{{ log.created_at }}</td>
          </tr>
          <tr v-if="!logs.length">
            <td colspan="4" style="text-align: center; color: var(--text-3); padding: 40px">暂无操作日志</td>
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
