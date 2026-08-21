<script setup lang="ts">
/** 标签页 — 标签云（统一字号）+ 标签卡片列表。 */
import { onMounted, ref } from 'vue'
import { ChevronRight, Hash } from 'lucide-vue-next'
import { getTags } from '@/api/site'
import type { Tag } from '@/types/models'

const tags = ref<Tag[]>([])

onMounted(async () => {
  try {
    tags.value = await getTags()
  } catch {
    /* 静默 */
  }
})
</script>

<template>
  <div>
    <section style="padding:118px 0 60px">
      <div class="container">
        <div class="anim-fade" style="text-align:center;margin-bottom:36px">
          <div style="display:flex;align-items:center;justify-content:center;gap:10px;font-size:13px;color:var(--text-3);margin-bottom:12px">
            <router-link to="/" style="color:var(--text-2)">首页</router-link>
            <ChevronRight :size="13" />
            <span class="grad-text" style="font-weight:600">标签</span>
          </div>
          <h1 style="font-size:32px;font-weight:800;letter-spacing:-.5px">标签云</h1>
          <p class="muted" style="margin-top:10px;font-size:14.5px">点击任意标签，查看该主题下的所有文章。</p>
        </div>

        <div class="widget" style="margin-bottom:34px">
          <div class="tag-cloud" style="justify-content:center;padding:26px 10px">
            <router-link
              v-for="t in tags"
              :key="t.id"
              :to="`/articles?tag=${encodeURIComponent(t.name)}`"
              style="font-size:15px"
            >
              {{ t.name }}
            </router-link>
          </div>
        </div>

        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <router-link
            v-for="t in tags"
            :key="t.id"
            :to="`/articles?tag=${encodeURIComponent(t.name)}`"
            class="glass-card card-hover"
            style="padding:18px 20px;display:flex;align-items:center;justify-content:space-between;gap:10px"
          >
            <span style="display:flex;align-items:center;gap:10px;font-size:14.5px;font-weight:600"><Hash :size="16" style="color:var(--accent)" />{{ t.name }}</span>
            <span class="badge" style="background:var(--grad-soft);color:var(--accent);border:1px solid var(--border)">{{ t.count ?? 0 }} 篇</span>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>
