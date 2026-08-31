<script setup lang="ts">
/** 标签页 — 标签列表（名称 + 文章数），按文章数降序排列。 */
import { onMounted, ref } from 'vue'
import { ChevronRight, Hash } from 'lucide-vue-next'
import { getTags } from '@/api/site'
import type { Tag } from '@/types/models'

const tags = ref<Tag[]>([])

onMounted(async () => {
  try {
    tags.value = (await getTags()).sort((a, b) => (b.count ?? 0) - (a.count ?? 0))
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div>
    <section style="padding:118px 0 60px">
      <div class="container" style="max-width:680px">
        <div class="anim-fade" style="margin-bottom:32px">
          <div style="display:flex;align-items:center;gap:10px;font-size:13px;color:var(--text-3);margin-bottom:10px">
            <router-link to="/" style="color:var(--text-2)">首页</router-link>
            <ChevronRight :size="13" />
            <span class="grad-text" style="font-weight:600">标签</span>
          </div>
          <h1 style="font-size:28px;font-weight:800;letter-spacing:-.4px">技术标签</h1>
          <p class="muted" style="margin-top:8px;font-size:14px">点击标签查看相关文章。</p>
        </div>

        <div class="tag-grid">
          <router-link
            v-for="t in tags"
            :key="t.id"
            :to="`/articles?tag_id=${t.id}`"
            class="tag-row"
          >
            <span class="tag-name"><Hash :size="15" />{{ t.name }}</span>
            <span class="tag-count">{{ t.count ?? 0 }} 篇</span>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.tag-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 0 40px;
}
.tag-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 0;
  border-bottom: 1px dashed var(--border);
  transition: color 0.2s;
}
.tag-row:last-child {
  border-bottom: none;
}
.tag-row:hover {
  color: var(--accent);
}
.tag-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
}
.tag-name svg {
  color: var(--text-3);
}
.tag-row:hover .tag-name svg {
  color: var(--accent);
}
.tag-count {
  font-size: 13px;
  color: var(--text-3);
  font-family: 'JetBrains Mono', monospace;
}
.tag-row:hover .tag-count {
  color: var(--accent);
}
</style>
