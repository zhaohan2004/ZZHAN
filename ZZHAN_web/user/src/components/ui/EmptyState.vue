<script setup lang="ts">
import { computed } from 'vue'
import { FileX, Inbox, PackageOpen, SearchX, type LucideIcon } from 'lucide-vue-next'

const props = withDefaults(
  defineProps<{
    icon?: 'search-x' | 'inbox' | 'file-x' | 'package-open' | string
    text?: string
    sub?: string
  }>(),
  { icon: 'search-x', text: '暂无数据', sub: '' },
)

const ICONS: Record<string, LucideIcon> = {
  'search-x': SearchX,
  inbox: Inbox,
  'file-x': FileX,
  'package-open': PackageOpen,
}

const Icon = computed(() => ICONS[props.icon] ?? SearchX)
</script>

<template>
  <div class="empty-state flex flex-col items-center justify-center gap-3 px-6 py-14 text-center">
    <div class="empty-icon flex h-16 w-16 items-center justify-center rounded-2xl bg-[var(--grad-soft)] text-[var(--accent)]">
      <component :is="Icon" :size="34" :stroke-width="1.6" />
    </div>
    <p class="empty-text text-[15px] font-semibold text-text-2">{{ text }}</p>
    <p v-if="sub" class="empty-sub text-[13px] text-text-3">{{ sub }}</p>
  </div>
</template>
