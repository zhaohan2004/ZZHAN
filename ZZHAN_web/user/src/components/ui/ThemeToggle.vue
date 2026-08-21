<script setup lang="ts">
import { Moon, Sun } from 'lucide-vue-next'
import { useThemeStore } from '@/stores/theme'

const theme = useThemeStore()
</script>

<template>
  <button
    class="icon-btn"
    type="button"
    :title="theme.dark ? '切换亮色' : '切换暗色'"
    data-theme-icon
    aria-label="切换主题"
    @click="theme.toggle()"
  >
    <Sun v-if="theme.dark" :size="18" />
    <Moon v-else :size="18" />
  </button>
</template>
