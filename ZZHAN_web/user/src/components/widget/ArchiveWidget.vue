<script setup lang="ts">
/**
 * 月度归档 widget（首页右侧）— 结构对齐静态原型 .arch-widget。
 */
import { onMounted, ref } from 'vue'
import { ArrowRight, ChevronDown } from 'lucide-vue-next'
import { getArchives } from '@/api/site'
import type { ArchiveItem } from '@/types/models'

const items = ref<ArchiveItem[]>([])
const open = ref<Set<string>>(new Set())

const keyOf = (y: string, m: string) => y + '-' + m

onMounted(async () => {
  try {
    items.value = await getArchives()
    const sorted = [...items.value].sort((a, b) => (b.year + b.month).localeCompare(a.year + a.month))
    sorted.slice(0, 3).forEach((x) => open.value.add(keyOf(x.year, x.month)))
  } catch {
    /* 静默 */
  }
})

function toggle(y: string, m: string): void {
  const k = keyOf(y, m)
  if (open.value.has(k)) open.value.delete(k)
  else open.value.add(k)
  open.value = new Set(open.value)
}
</script>

<template>
  <div class="widget arch-widget reveal in">
    <h3 class="widget-title">归档</h3>
    <div v-for="item in items" :key="keyOf(item.year, item.month)" class="aw-month" :class="{ open: open.has(keyOf(item.year, item.month)) }">
      <button class="aw-head" type="button" @click="toggle(item.year, item.month)">
        {{ item.year }} 年 {{ Number(item.month) }} 月
        <span class="aw-cnt">{{ item.count }} 篇</span>
        <ChevronDown class="aw-chev" :size="15" />
      </button>
      <div v-if="open.has(keyOf(item.year, item.month))" class="aw-body">
        <router-link
          v-for="a in item.articles.slice(0, 3)"
          :key="a.id"
          :to="`/article/${a.id}`"
        >
          {{ a.title }}
        </router-link>
      </div>
    </div>
    <router-link to="/archive" class="section-more mt-3 inline-flex" style="color:var(--accent)">
      <ArrowRight :size="14" /> 全部归档
    </router-link>
  </div>
</template>
