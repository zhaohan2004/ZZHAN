<script setup lang="ts">
/** 最近动态时间线 — write / star / talk 图标 + 文本 + 时间。 */
import { MessageSquare, PenTool, Star } from 'lucide-vue-next'
import type { Dynamic } from '@/types/models'

defineProps<{ dynamics: Dynamic[] }>()

const ICONS: Record<string, any> = { write: PenTool, star: Star, talk: MessageSquare }
</script>

<template>
  <div class="timeline">
    <div v-for="(d, i) in dynamics" :key="i" class="tl-item done reveal in">
      <div class="tl-head">
        <span class="tl-date">{{ d.time }}</span>
        <component :is="ICONS[d.type] || MessageSquare" :size="13" style="color:var(--accent)" />
      </div>
      <span class="tl-title">
        <template v-if="d.link">
          {{ d.text.slice(0, d.text.indexOf('《')) }}<router-link :to="d.link">{{ d.text.slice(d.text.indexOf('《'), d.text.indexOf('》') + 1) }}</router-link>{{ d.text.slice(d.text.indexOf('》') + 1) }}
        </template>
        <template v-else>{{ d.text }}</template>
      </span>
    </div>
  </div>
</template>
