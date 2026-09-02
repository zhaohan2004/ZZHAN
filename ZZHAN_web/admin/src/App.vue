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
  // 加载设置（未登录时用公开接口获取基本信息）
  settingsStore.load()
  // 登录状态下自动加载管理员资料
  if (authStore.loggedIn) {
    authStore.loadProfile()
  }
})
</script>
