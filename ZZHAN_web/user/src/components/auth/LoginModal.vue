<script setup lang="ts">
/**
 * 登录弹窗 — GitHub OAuth + 完善资料 step。
 */
import { computed, ref, watch } from 'vue'
import { Check, Github, Image as ImageIcon, Loader2, X } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'
import { useToast } from '@/composables/useToast'
import { initialsAvatar } from '@/utils/avatar'

const AVATAR_PALETTE = [
  ['#3b82f6', '#38bdf8'],
  ['#8b5cf6', '#c084fc'],
  ['#ec4899', '#f472b6'],
  ['#ef4444', '#fb923c'],
  ['#10b981', '#34d399'],
  ['#f59e0b', '#fbbf24'],
  ['#06b6d4', '#22d3ee'],
  ['#64748b', '#94a3b8'],
]

const auth = useAuthStore()
const site = useSiteStore()
const { toast } = useToast()

const step = ref<'login' | 'profile'>('login')
const loading = ref(false)
const nickname = ref('')
const picked = ref(0)
const avatarPreview = ref('')

const open = computed(() => auth.loginModalOpen)

watch(open, (v) => {
  if (v) {
    step.value = auth.need_profile ? 'profile' : 'login'
    if (step.value === 'profile' && auth.user) {
      nickname.value = auth.user.nickname || ''
      avatarPreview.value = auth.user.avatar || ''
    }
  }
})

function close(): void {
  auth.loginModalOpen = false
}

async function doLogin(): Promise<void> {
  loading.value = true
  try {
    await auth.loginWith('github')
    toast('登录成功', 'success')
    if (auth.need_profile) {
      step.value = 'profile'
      nickname.value = auth.user?.nickname ?? ''
    }
  } catch {
    toast('登录失败，请重试', 'error')
  } finally {
    loading.value = false
  }
}

function pickPalette(i: number): void {
  picked.value = i
  const [c1, c2] = AVATAR_PALETTE[i]
  avatarPreview.value = initialsAvatar(nickname.value || '我', c1, c2, 160)
}

function onFile(e: Event): void {
  const f = (e.target as HTMLInputElement).files?.[0]
  if (!f) return
  const r = new FileReader()
  r.onload = () => {
    avatarPreview.value = String(r.result)
  }
  r.readAsDataURL(f)
}

async function complete(): Promise<void> {
  if (!nickname.value.trim()) {
    toast('请输入昵称', 'error')
    return
  }
  if (!avatarPreview.value) pickPalette(picked.value)
  try {
    await auth.completeProfile({ nickname: nickname.value.trim(), avatar: avatarPreview.value })
    toast('资料已完善', 'success')
    close()
  } catch {
    toast('保存失败，请重试', 'error')
  }
}
</script>

<template>
  <Teleport to="body">
    <div v-if="open" class="modal-overlay open login-modal" @click.self="close">
      <div class="modal">
        <div class="modal-head">
          <h3>{{ step === 'login' ? '登录 ' + (site.site?.name ?? '小猫的个人博客') : '完善资料' }}</h3>
          <button class="modal-close" type="button" aria-label="关闭" @click="close"><X :size="17" /></button>
        </div>
        <div class="modal-body">
          <!-- 登录 step -->
          <template v-if="step === 'login'">
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
          </template>

          <!-- 完善资料 step -->
          <template v-else>
            <div class="lm-profile">
              <h4>完善资料</h4>
              <p class="lm-psub">设置昵称与头像，方便他人认识你</p>
              <div class="form-group" style="text-align:left">
                <label class="form-label">昵称 <span style="color:var(--danger)">*</span></label>
                <input v-model="nickname" class="input" type="text" placeholder="请输入昵称" maxlength="20" />
              </div>
              <div class="avatar-grid" style="text-align:left">
                <button
                  v-for="(p, i) in AVATAR_PALETTE"
                  :key="i"
                  type="button"
                  class="av-opt"
                  :class="{ selected: picked === i }"
                  :style="{ background: 'linear-gradient(135deg,' + p[0] + ',' + p[1] + ')' }"
                  @click="pickPalette(i)"
                >
                  <Check v-if="picked === i" :size="16" style="color:#fff;margin:auto" />
                </button>
              </div>
              <label class="upload-avatar-btn">
                <ImageIcon :size="15" /> 上传自定义头像
                <input type="file" accept="image/*" class="hidden" @change="onFile" />
              </label>
              <div v-if="avatarPreview" class="mb-4 flex items-center justify-center gap-3">
                <img :src="avatarPreview" alt="头像预览" style="width:56px;height:56px;border-radius:14px;border:1px solid var(--border-strong)" />
                <span class="muted text-[12.5px]">头像预览</span>
              </div>
              <button class="btn btn-primary w-full" type="button" @click="complete">保存并进入</button>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>
