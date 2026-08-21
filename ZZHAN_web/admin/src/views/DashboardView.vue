<script setup lang="ts">
/** 仪表盘 — 1:1 复刻静态 index.html（.stat-grid/.quick-actions/.chart-card/.widget/.timeline）。
 *  组件外壳用 style.css 后台类；外层栅格沿用 Tailwind（与静态一致的混用风格）。去 Element Plus。 */
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  BarChart3,
  Eye,
  FileText,
  Folder,
  MessageCircle,
  PenLine,
  Tag,
  TrendingDown,
  TrendingUp,
} from 'lucide-vue-next'
import {
  getDashboardArticles,
  getDashboardCharts,
  getDashboardComments,
  getDashboardOperations,
  getDashboardStats,
} from '@/api/dashboard'
import type { DashboardArticles, DashboardCharts, DashboardStat, Operation, RecentComment } from '@/types/models'
import { fmtNum } from '@/utils/format'
import BarChart from '@/components/chart/BarChart.vue'
import LineChart from '@/components/chart/LineChart.vue'

const router = useRouter()

const stats = ref<DashboardStat[]>([])
const charts = ref<DashboardCharts | null>(null)
const articles = ref<DashboardArticles | null>(null)
const comments = ref<RecentComment[]>([])
const operations = ref<Operation[]>([])

const QUICK = [
  { to: '/editor', label: '写新文章', icon: PenLine },
  { to: '/articles', label: '文章管理', icon: FileText },
  { to: '/categories', label: '分类管理', icon: Folder },
  { to: '/tags', label: '标签管理', icon: Tag },
  { to: '/comments', label: '评论管理', icon: MessageCircle },
]

// mock 里 stat.icon 是 Element Plus 图标名，这里映射到 lucide
const ICONS: Record<string, unknown> = {
  Document: FileText,
  Folder: Folder,
  PriceTag: Tag,
  ChatDotRound: MessageCircle,
  View: Eye,
  TrendCharts: BarChart3,
}
function statIcon(name: string): unknown {
  return ICONS[name] || FileText
}

onMounted(async () => {
  try {
    const [s, c, a, cm, op] = await Promise.all([
      getDashboardStats(),
      getDashboardCharts(),
      getDashboardArticles(),
      getDashboardComments(),
      getDashboardOperations(),
    ])
    stats.value = s
    charts.value = c
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
    <!-- 快捷入口 -->
    <div class="quick-actions" style="margin-bottom: 22px">
      <router-link v-for="q in QUICK" :key="q.to" :to="q.to" class="qa-item">
        <component :is="q.icon" :size="20" />
        <div><b>{{ q.label }}</b><span>管理入口</span></div>
      </router-link>
    </div>

    <!-- 统计卡 -->
    <div class="stat-grid">
      <div v-for="s in stats" :key="s.key" class="stat-card" :style="{ '--sc': s.color }">
        <div class="stat-ico" :style="{ background: 'linear-gradient(135deg,' + s.color + 'dd,' + s.color + '77)' }">
          <component :is="statIcon(s.icon)" :size="24" />
        </div>
        <div class="stat-info"><b>{{ s.value }}</b><span>{{ s.label }}</span></div>
        <span class="stat-trend" :class="s.up ? 'up' : 'down'">
          <component :is="s.up ? TrendingUp : TrendingDown" :size="13" /> {{ s.trend }}
        </span>
      </div>
    </div>

    <!-- 图表区 -->
    <div style="display: grid; gap: 20px; grid-template-columns: repeat(auto-fit, minmax(360px, 1fr)); margin-bottom: 22px">
      <div class="chart-card">
        <div class="chart-title">最近 7 天访问量 <small>PV 数据</small></div>
        <div class="chart-box" style="height: 240px">
          <BarChart v-if="charts" :data="charts.weekVisits.map((d) => ({ label: d.d, value: d.pv }))" />
        </div>
      </div>
      <div class="chart-card">
        <div class="chart-title">
          <span>用户访问趋势 <small>近 7 天 PV / UV</small></span>
          <span style="display:inline-flex;gap:14px;align-items:center;font-size:12.5px;color:var(--text-2)">
            <span style="display:inline-flex;align-items:center;gap:6px"><i style="width:10px;height:10px;border-radius:3px;background:#818cf8;display:inline-block"></i>访问量 PV</span>
            <span style="display:inline-flex;align-items:center;gap:6px"><i style="width:10px;height:10px;border-radius:3px;background:#34d399;display:inline-block"></i>访客数 UV</span>
          </span>
        </div>
        <div class="chart-box" style="height: 240px">
          <LineChart
            v-if="charts"
            :labels="charts.weekVisits.map((d) => d.d)"
            :series="[
              { data: charts.weekVisits.map((d) => d.pv), color: '#60a5fa', name: '访问量 PV' },
              { data: charts.weekVisits.map((d) => d.uv), color: '#34d399', name: '访客数 UV' },
            ]"
          />
        </div>
      </div>
    </div>

    <!-- 底部列表 -->
    <div style="display: grid; gap: 20px; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))">
      <div class="widget">
        <h3 class="widget-title">最近发布文章</h3>
        <div
          v-for="a in articles?.recentPosts ?? []"
          :key="a.id"
          style="padding: 10px 0; border-bottom: 1px dashed var(--border)"
        >
          <div style="font-size: 13.5px; font-weight: 600; color: var(--text-2)">{{ a.title }}</div>
          <div style="font-size: 12px; color: var(--text-3); margin-top: 3px">{{ a.category }} · {{ a.date }} · {{ fmtNum(a.views) }} 阅读</div>
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
            <div style="font-size: 13px; font-weight: 600; color: var(--text-2)">{{ c.userName }} <span style="font-weight: 400; color: var(--text-3)">{{ c.time }}</span></div>
            <div style="font-size: 12.5px; color: var(--text-3); overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{ c.content }}</div>
          </div>
        </div>
      </div>
      <div class="widget">
        <h3 class="widget-title">最近操作记录</h3>
        <div class="timeline">
          <div v-for="op in operations" :key="op.time + op.target" style="display: flex; gap: 10px; font-size: 13px; padding: 4px 0">
            <span style="font-family: 'JetBrains Mono', monospace; font-size: 12px; color: var(--text-3)">{{ op.time }}</span>
            <span style="color: var(--text-2)"><b style="font-weight: 600">{{ op.user }}</b> {{ op.action }}：{{ op.target }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
