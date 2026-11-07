<script setup lang="ts">
import { computed, ref } from 'vue'
import type { ArticleSummary } from '@/types/models'
import { coverArt } from '@/utils/cover'

const props = withDefaults(defineProps<{
  article: ArticleSummary
  alt?: string
  thumbnail?: boolean
}>(), { alt: '', thumbnail: false })

const failedUrl = ref('')
const fallback = computed(() => {
  const { title, category_name, id } = props.article
  return props.thumbnail ? coverArt(title, category_name, id, 96, 72) : coverArt(title, category_name, id)
})
const coverUrl = computed(() => {
  const url = props.article.cover_image
  return url && url !== failedUrl.value ? url : fallback.value
})

function showFallback(event: Event) {
  const image = event.currentTarget as HTMLImageElement
  if (image.getAttribute('src') === props.article.cover_image) {
    failedUrl.value = props.article.cover_image
  }
}
</script>

<template>
  <img :src="coverUrl" :alt="alt" loading="lazy" decoding="async" @error="showFallback" />
</template>
