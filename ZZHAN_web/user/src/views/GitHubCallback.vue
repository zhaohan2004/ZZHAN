<script setup lang="ts">
/**
 * GitHub OAuth 回调页面 — 处理 GitHub 重定向回来的 code 参数。
 */
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

const router = useRouter()
const auth = useAuthStore()
const { toast } = useToast()

onMounted(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')

  if (!code) {
    toast('缺少授权码', 'error')
    router.push('/')
    return
  }

  try {
    const success = await auth.loginWith('github')
    if (success) {
      toast('登录成功', 'success')
    } else {
      toast('登录失败', 'error')
    }
  } catch (error) {
    console.error('GitHub 登录失败:', error)
    toast('登录失败，请重试', 'error')
  } finally {
    // 跳转到首页
    router.push('/')
  }
})
</script>

<template>
  <div class="callback-page">
    <div class="loading-container">
      <div class="spinner"></div>
      <p>正在完成 GitHub 登录...</p>
    </div>
  </div>
</template>

<style scoped>
.callback-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
}

.loading-container {
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
