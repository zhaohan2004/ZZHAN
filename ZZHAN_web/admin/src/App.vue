<template>
  <RouterView />
  <ToastHost />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import ToastHost from '@/components/common/ToastHost.vue'
import { useSettingsStore } from '@/stores/settings'
import { useAuthStore } from '@/stores/auth'

const settingsStore = useSettingsStore()
const authStore = useAuthStore()

onMounted(() => {
  settingsStore.load()
  // 登录状态下自动加载管理员资料（解决页面刷新后 profile 为 null 的问题）
  if (authStore.loggedIn) {
    authStore.loadProfile()
  }
})
</script>
