<script setup lang="ts">
/**
 * 后台侧边栏 — 复刻静态 .admin-sidebar（side-brand/side-nav/side-group/side-link/badge/side-foot）。
 * 去 Element Plus，图标改用 lucide-vue-next。
 */
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { FileText, Folder, Home, MessageCircle, ScrollText, Settings, Tag, Users } from 'lucide-vue-next'
import { listComments } from '@/api/admin'
import { useSettingsStore } from '@/stores/settings'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const route = useRoute()
const settings = useSettingsStore()
const bannedCount = ref(0)

onMounted(async () => {
  try {
    const res = await listComments({ status: 'banned', pageSize: 1 })
    bannedCount.value = res.total
  } catch {
    /* 静默 */
  }
})

interface NavItem {
  to: string
  label: string
  icon: unknown
  exact?: boolean
}

const groups: { label: string; items: NavItem[] }[] = [
  {
    label: '概览',
    items: [{ to: '/', label: '仪表盘', icon: Home, exact: true }],
  },
  {
    label: '内容管理',
    items: [
      { to: '/articles', label: '文章管理', icon: FileText },
      { to: '/categories', label: '分类管理', icon: Folder },
      { to: '/tags', label: '标签管理', icon: Tag },
      { to: '/comments', label: '评论管理', icon: MessageCircle },
      { to: '/users', label: '用户管理', icon: Users },
    ],
  },
  {
    label: '系统',
    items: [
      { to: '/operation-logs', label: '操作日志', icon: ScrollText },
      { to: '/settings', label: '系统设置', icon: Settings },
    ],
  },
]

function isActive(item: NavItem): boolean {
  if (item.exact) return route.path === item.to
  if (item.to === '/articles') return route.path === '/articles' || route.path.startsWith('/editor')
  return route.path.startsWith(item.to)
}
</script>

<template>
  <aside class="admin-sidebar" :class="{ open }">
    <div class="side-brand">
      <span class="brand-logo">
        <img v-if="settings.settings?.avatar" :src="settings.settings.avatar" alt="头像" />
        <span v-else>{{ (settings.settings?.blog_name || 'C').slice(0, 1) }}</span>
      </span>
      <span class="brand-text">{{ settings.settings?.blog_name || 'Blog' }}<span class="brand-dot">.</span><span style="font-size:11px;color:var(--text-3);font-weight:400;margin-left:8px">管理后台</span></span>
    </div>
    <nav class="side-nav">
      <template v-for="g in groups" :key="g.label">
        <div class="side-group">{{ g.label }}</div>
        <router-link
          v-for="item in g.items"
          :key="item.to"
          :to="item.to"
          class="side-link"
          :class="{ active: isActive(item) }"
          @click="emit('close')"
        >
          <component :is="item.icon" :size="18" />
          {{ item.label }}
          <span
            v-if="item.to === '/comments' && bannedCount"
            class="badge"
            style="margin-left:auto;background:rgba(239,68,68,.14);color:#ef4444"
          >{{ bannedCount }}</span>
        </router-link>
      </template>
    </nav>
    <div class="side-foot">
      <a href="/" target="_blank" rel="noopener"><Home :size="16" /> 返回前台</a>
    </div>
  </aside>
  <div class="side-backdrop" :class="{ open }" @click="emit('close')" />
</template>
