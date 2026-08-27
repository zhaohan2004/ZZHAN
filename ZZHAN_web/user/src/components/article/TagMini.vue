<script setup lang="ts">
/** 迷你标签 chip — 链接到 /articles?tag_id= */
defineProps<{ tagId: number; tagName: string }>()
</script>

<template>
  <router-link class="tag-mini" :to="`/articles?tag_id=${tagId}`">#{{ tagName }}</router-link>
</template>
