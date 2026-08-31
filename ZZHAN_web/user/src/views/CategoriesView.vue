<script setup lang="ts">
/** 分类页 — 分类卡片网格（lucide 图标 / 颜色 / 描述 / 文章数 / 占比）。对齐静态原型。 */
import { onMounted, ref } from 'vue'
import {
  Activity,
  ChevronRight,
  Code2,
  Container,
  Cpu,
  Database,
  GitBranch,
  Network,
  Radio,
  Route,
  Terminal,
  Zap,
} from 'lucide-vue-next'
import { getCategories } from '@/api/site'
import type { Category } from '@/types/models'
import { useReveal } from '@/composables/useReveal'

const ICON_MAP: Record<string, any> = {
  'code-2': Code2,
  database: Database,
  zap: Zap,
  terminal: Terminal,
  container: Container,
  'git-branch': GitBranch,
  route: Route,
  radio: Radio,
  activity: Activity,
  network: Network,
  cpu: Cpu,
}

const root = ref<HTMLElement | null>(null)
useReveal(root)

const cats = ref<Category[]>([])
const total = ref(0)

onMounted(async () => {
  try {
    cats.value = (await getCategories()).sort((a, b) => (b.count ?? 0) - (a.count ?? 0))
    total.value = cats.value.reduce((n, c) => n + (c.count ?? 0), 0)
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div ref="root">
    <section style="padding:118px 0 60px">
      <div class="container">
        <div class="anim-fade" style="text-align:center;margin-bottom:40px">
          <div style="display:flex;align-items:center;justify-content:center;gap:10px;font-size:13px;color:var(--text-3);margin-bottom:12px">
            <router-link to="/" style="color:var(--text-2)">首页</router-link>
            <ChevronRight :size="13" />
            <span class="grad-text" style="font-weight:600">分类</span>
          </div>
          <h1 style="font-size:32px;font-weight:800;letter-spacing:-.5px">文章分类</h1>
          <p class="muted" style="margin-top:10px;font-size:14.5px">按技术领域整理知识，找到你感兴趣的方向。</p>
        </div>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <router-link
            v-for="c in cats"
            :key="c.slug"
            :id="'cat-' + c.slug"
            :to="`/articles?category=${c.slug}`"
            class="cat-card reveal"
            :style="{ '--cat-c': c.color }"
          >
            <div class="cat-ico" :style="{ background: 'linear-gradient(135deg,' + c.color + 'cc,' + c.color + '66)' }">
              <component :is="ICON_MAP[c.icon || ''] || Code2" :size="24" />
            </div>
            <h3>{{ c.name }}</h3>
            <p>{{ c.desc }}</p>
            <div class="cat-foot">
              <span><b :style="{ color: c.color }">{{ c.count ?? 0 }}</b> 篇文章</span>
              <span>{{ total ? Math.round(((c.count ?? 0) / total) * 100) : 0 }}% 占比</span>
            </div>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>
