<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { ArrowUp } from 'lucide-vue-next'

const show = ref(false)

function onScroll(): void {
  show.value = window.scrollY > 500
}

function toTop(): void {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})
onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <button class="back-top" :class="{ show }" type="button" title="回到顶部" aria-label="回到顶部" @click="toTop">
    <ArrowUp :size="19" />
  </button>
</template>
