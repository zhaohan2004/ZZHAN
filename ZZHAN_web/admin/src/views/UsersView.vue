<script setup lang="ts">
/** 用户管理 */
import { computed, onMounted, reactive, ref } from 'vue'
import { RotateCcw, Search, ShieldCheck, ShieldOff, Trash2 } from 'lucide-vue-next'
import { deleteUser, listUsers, setUserStatus } from '@/api/admin'
import type { UserAdmin, UserStatus } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'

const users = ref<UserAdmin[]>([])
const total = ref(0)
const loading = ref(false)

const query = reactive<{
  page: number
  pageSize: number
  keyword: string
  status: UserStatus | 'all'
  startDate: string
  endDate: string
}>({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: 'all',
  startDate: '',
  endDate: '',
})

const STATUS_MAP: Record<number, { cls: string; label: string }> = {
  0: { cls: 'st-block', label: '禁用' },
  1: { cls: 'st-ok', label: '正常' },
}

const PROVIDER_MAP: Record<string, string> = {
  github: 'GitHub',
  wechat: '微信',
  google: 'Google',
}

async function load() {
  loading.value = true
  try {
    const res = await listUsers({
      page: query.page,
      pageSize: query.pageSize,
      keyword: query.keyword || undefined,
      status: query.status,
      startDate: query.startDate || undefined,
      endDate: query.endDate || undefined,
    })
    users.value = res.list
    total.value = res.total
  } catch {
    toast.error('加载用户失败')
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
  query.keyword = ''
  query.status = 'all'
  query.startDate = ''
  query.endDate = ''
  query.page = 1
  load()
}

async function handleSetStatus(u: UserAdmin, status: UserStatus, label: string) {
  const ok = await confirm(`确定将用户「${u.nickname}」${label}吗？`, '修改状态')
  if (!ok) return
  try {
    await setUserStatus(u.id, status)
    toast.success(`已${label}`)
    load()
  } catch {
    toast.error('操作失败')
  }
}

async function handleDelete(u: UserAdmin) {
  const ok = await confirm(`确定删除用户「${u.nickname}」吗？此操作不可恢复。`, '删除用户')
  if (!ok) return
  try {
    await deleteUser(u.id)
    toast.success('已删除')
    if (users.value.length === 1 && query.page > 1) query.page--
    load()
  } catch {
    toast.error('删除失败')
  }
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
        <input v-model="query.keyword" class="input" placeholder="搜索昵称" style="width:160px" @keyup.enter="search" />
        <input v-model="query.startDate" class="input" type="date" style="width:150px" />
        <span class="muted" style="font-size:13px">至</span>
        <input v-model="query.endDate" class="input" type="date" style="width:150px" />
        <select v-model="query.status" class="select" style="width:110px">
          <option value="all">全部状态</option>
          <option :value="1">正常</option>
          <option :value="0">禁用</option>
        </select>
        <button class="btn btn-ghost btn-sm" @click="search"><Search :size="14" /> 查询</button>
        <button class="btn btn-ghost btn-sm" @click="reset"><RotateCcw :size="14" /> 重置</button>
      </div>
      <span class="muted" style="font-size: 13px">共 <b>{{ total }}</b> 个用户</span>
    </div>

    <div class="data-table-wrap reveal in" v-if="!loading">
      <table class="data-table">
        <thead>
          <tr>
            <th>用户</th>
            <th>登录方式</th>
            <th>状态</th>
            <th>最后登录</th>
            <th>注册时间</th>
            <th style="text-align: right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>
              <div style="display: flex; align-items: center; gap: 10px">
                <img
                  v-if="u.avatar"
                  :src="u.avatar"
                  :alt="u.nickname"
                  style="width: 32px; height: 32px; border-radius: 10px; object-fit: cover"
                />
                <span
                  v-else
                  style="width: 32px; height: 32px; border-radius: 10px; background: var(--primary); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 600"
                >
                  {{ u.nickname?.charAt(0)?.toUpperCase() || '?' }}
                </span>
                <b style="font-weight: 600; color: var(--text-2)">{{ u.nickname }}</b>
              </div>
            </td>
            <td style="color: var(--text-2)">{{ PROVIDER_MAP[u.provider] || u.provider }}</td>
            <td><span class="st" :class="STATUS_MAP[u.status]?.cls">{{ STATUS_MAP[u.status]?.label }}</span></td>
            <td style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px">{{ u.last_login_at || '-' }}</td>
            <td style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px">{{ u.created_at }}</td>
            <td style="text-align: right; white-space: nowrap">
              <button
                v-if="u.status === 1"
                class="tb-btn edit"
                title="禁用"
                @click="handleSetStatus(u, 0, '禁用')"
              >
                <ShieldOff :size="15" />
              </button>
              <button
                v-else
                class="tb-btn view"
                title="启用"
                @click="handleSetStatus(u, 1, '启用')"
              >
                <ShieldCheck :size="15" />
              </button>
              <button class="tb-btn del" title="删除" @click="handleDelete(u)"><Trash2 :size="15" /></button>
            </td>
          </tr>
          <tr v-if="!users.length">
            <td colspan="6" style="text-align: center; color: var(--text-3); padding: 40px">暂无用户</td>
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
