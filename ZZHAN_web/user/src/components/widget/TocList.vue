<script setup lang="ts">
/** 文章目录 TOC — 由 ArticleView 传入 buildTOC 结果与当前高亮 id。 */
export interface TocEntry {
  level: 2 | 3 | 4
  id: string
  text: string
}

defineProps<{ items: TocEntry[]; activeId: string }>()
const emit = defineEmits<{ (e: 'jump', id: string): void }>()
</script>

<template>
  <nav v-if="items.length" class="toc">
    <a
      v-for="item in items"
      :key="item.id"
      class="toc-link"
      :class="['lv' + item.level, { active: item.id === activeId }]"
      :href="'#' + item.id"
      @click.prevent="emit('jump', item.id)"
    >
      {{ item.text }}
    </a>
  </nav>
</template>
