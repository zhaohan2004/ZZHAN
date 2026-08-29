<script setup lang="ts">
/** 仪表盘 — 统计卡 + 文章/评论/操作列表。
 *  组件外壳用 style.css 后台类；外层栅格沿用 Tailwind（与静态一致的混用风格）。去 Element Plus。 */
import { computed, onMounted, ref } from 'vue'
import {
  FileText,
  Folder,
  MessageCircle,
  Tag,
} from 'lucide-vue-next'
import {
  getDashboardArticles,
  getDashboardComments,
  getDashboardOperations,
  getDashboardStats,
} from '@/api/dashboard'
import type { DashboardArticles, DashboardStatData, Operation, RecentComment } from '@/types/models'
import { fmtNum } from '@/utils/format'

/** 统计卡展示配置 — 前端自管 */
const STAT_CARDS = [
  { key: 'articles', label: '文章总数', icon: FileText, color: '#5b9df6' },
  { key: 'categories', label: '分类数量', icon: Folder, color: '#38bdf8' },
  { key: 'tags', label: '标签数量', icon: Tag, color: '#a78bfa' },
  { key: 'comments', label: '评论数量', icon: MessageCircle, color: '#f59e0b' },
] as const

const rawData = ref<DashboardStatData | null>(null)
const articles = ref<DashboardArticles | null>(null)
const comments = ref<RecentComment[]>([])
const operations = ref<Operation[]>([])

/** 合并展示配置与后端数据 */
const stats = computed(() =>
  STAT_CARDS.map((c) => {
    const d = rawData.value?.[c.key]
    return { ...c, value: d?.value ?? 0 }
  }),
)

onMounted(async () => {
  try {
    const [s, a, cm, op] = await Promise.all([
      getDashboardStats(),
      getDashboardArticles(),
      getDashboardComments(),
      getDashboardOperations(),
    ])
    rawData.value = s
    articles.value = a
    comments.value = cm
    operations.value = op
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div>
    <!-- 统计卡 -->
    <div class="stat-grid">
      <div v-for="s in stats" :key="s.key" class="stat-card" :style="{ '--sc': s.color }">
        <div class="stat-ico" :style="{ background: 'linear-gradient(135deg,' + s.color + 'dd,' + s.color + '77)' }">
          <component :is="s.icon" :size="24" />
        </div>
        <div class="stat-info"><b>{{ s.value }}</b><span>{{ s.label }}</span></div>
      </div>
    </div>

    <!-- 底部列表 -->
    <div style="display: grid; gap: 20px; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))">
      <div class="widget">
        <h3 class="widget-title">最近发布文章</h3>
        <div
          v-for="a in articles?.recent_posts ?? []"
          :key="a.id"
          style="padding: 10px 0; border-bottom: 1px dashed var(--border)"
        >
          <div style="font-size: 14.5px; font-weight: 600; color: var(--text-2)">{{ a.title }}</div>
          <div style="font-size: 13px; color: var(--text-3); margin-top: 3px">{{ a.category }} · {{ a.date }} · {{ fmtNum(a.views) }} 阅读</div>
        </div>
      </div>
      <div class="widget">
        <h3 class="widget-title">最近评论</h3>
        <div
          v-for="c in comments"
          :key="c.id"
          style="display: flex; align-items: flex-start; gap: 10px; padding: 10px 0; border-bottom: 1px dashed var(--border)"
        >
          <span style="width: 32px; height: 32px; border-radius: 10px; flex: none" :style="{ background: c.avatar }" />
          <div style="min-width: 0">
            <div style="font-size: 14px; font-weight: 600; color: var(--text-2)">{{ c.user_name }} <span style="font-weight: 400; color: var(--text-3)">{{ c.time }}</span></div>
            <div style="font-size: 13.5px; color: var(--text-3); overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{ c.content }}</div>
          </div>
        </div>
      </div>
      <div class="widget">
        <h3 class="widget-title">最近操作记录</h3>
        <div class="timeline">
          <div v-for="op in operations" :key="op.time + op.target" style="display: flex; gap: 10px; font-size: 15px; padding: 4px 0">
            <span style="font-family: 'JetBrains Mono', monospace; font-size: 14px; color: var(--text-3)">{{ op.time }}</span>
            <span style="color: var(--text-2)"><b style="font-weight: 600">{{ op.user }}</b> {{ op.action }}：{{ op.target }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
