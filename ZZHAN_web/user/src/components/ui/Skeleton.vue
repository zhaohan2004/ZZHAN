<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'card' | 'text'
    lines?: number
  }>(),
  { type: 'text', lines: 3 },
)
</script>

<template>
  <div class="skeleton" :class="type === 'card' ? 'skeleton-card' : 'skeleton-text'">
    <template v-if="type === 'card'">
      <div class="sk-cover h-40 rounded-xl bg-[var(--card-2)]" />
      <div class="sk-body flex flex-col gap-3 p-5">
        <div class="sk-line h-4 w-3/4 rounded-md bg-[var(--card-2)]" />
        <div class="sk-line h-4 w-full rounded-md bg-[var(--card-2)]" />
        <div class="sk-line h-4 w-1/2 rounded-md bg-[var(--card-2)]" />
      </div>
    </template>
    <template v-else>
      <div class="flex flex-col gap-3">
        <div
          v-for="n in lines"
          :key="n"
          class="sk-line h-4 rounded-md bg-[var(--card-2)]"
          :style="{ width: (100 - n * 8) + '%' }"
        />
      </div>
    </template>
  </div>
</template>
