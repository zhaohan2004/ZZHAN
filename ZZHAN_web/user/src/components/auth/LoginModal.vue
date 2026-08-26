<script setup lang="ts">
/**
 * 登录弹窗 — GitHub OAuth 登录。
 */
import { computed, ref } from 'vue'
import { Loader2, X } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'
import { useToast } from '@/composables/useToast'

const auth = useAuthStore()
const site = useSiteStore()
const { toast } = useToast()

const loading = ref(false)

const open = computed(() => auth.loginModalOpen)

function close(): void {
  auth.loginModalOpen = false
}

async function doLogin(): Promise<void> {
  loading.value = true
  try {
    const success = await auth.loginWith('github')
    if (success) {
      toast('登录成功', 'success')
      close()
    }
  } catch {
    toast('登录失败，请重试', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="modal-overlay open login-modal" @click.self="close">
      <div class="modal">
        <div class="modal-head">
          <h3>登录 {{ site.site?.name ?? '小猫的个人博客' }}</h3>
          <button class="modal-close" type="button" aria-label="关闭" @click="close"><X :size="17" /></button>
        </div>
        <div class="modal-body">
          <div style="display:flex;flex-direction:column;align-items:center;padding:30px 0 20px">
            <svg height="56" width="56" viewBox="0 0 16 16" style="margin-bottom:20px;fill:var(--text-1)">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
            </svg>
            <p class="lm-hint" style="margin-bottom:16px">使用 GitHub 账号授权登录</p>
            <button class="btn btn-primary w-full" type="button" :disabled="loading" @click="doLogin">
              <Loader2 v-if="loading" :size="16" class="animate-spin" /> GitHub 授权登录
            </button>
          </div>
          <p class="lm-note">登录即代表同意本站服务条款</p>
        </div>
      </div>
    </div>
  </Teleport>
</template>
