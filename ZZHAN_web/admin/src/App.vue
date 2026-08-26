<template>
  <RouterView />
  <ToastHost />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import ToastHost from '@/components/common/ToastHost.vue'
import { useSettingsStore } from '@/stores/settings'

const settingsStore = useSettingsStore()
onMounted(() => settingsStore.load())
</script>
