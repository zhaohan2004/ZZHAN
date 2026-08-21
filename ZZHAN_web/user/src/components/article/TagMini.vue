<script setup lang="ts">
/** 迷你标签 chip — 链接到 /articles?tag= */
defineProps<{ tag: string }>()
</script>

<template>
  <router-link class="tag-mini" :to="`/articles?tag=${encodeURIComponent(tag)}`">#{{ tag }}</router-link>
</template>
