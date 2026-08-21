<script setup lang="ts">
/**
 * 登录弹窗 — 微信 / GitHub OAuth（mock）+ 完善资料 step。
 * 结构对齐静态原型（.wx-qr 二维码 / .avatar-grid 4 列）。
 */
import { computed, ref, watch } from 'vue'
import { Check, Github, Image as ImageIcon, Loader2, X } from 'lucide-vue-next'
import { useAuthStore, type LoginProvider } from '@/stores/auth'
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
const { toast } = useToast()

const step = ref<'login' | 'profile'>('login')
const tab = ref<LoginProvider>('wechat')
const loading = ref(false)
const nickname = ref('')
const picked = ref(0)
const avatarPreview = ref('')

const open = computed(() => auth.loginModalOpen)

watch(open, (v) => {
  if (v) {
    step.value = auth.needProfile ? 'profile' : 'login'
    if (step.value === 'profile' && auth.user) {
      nickname.value = auth.user.nickname || ''
      avatarPreview.value = auth.user.avatar || ''
    }
  }
})

function close(): void {
  auth.loginModalOpen = false
}

async function doLogin(provider: LoginProvider): Promise<void> {
  loading.value = true
  try {
    await auth.loginWith(provider, 'mock-code-' + provider)
    toast('登录成功', 'success')
    if (auth.needProfile) {
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
          <h3>{{ step === 'login' ? '登录 CodeThink' : '完善资料' }}</h3>
          <button class="modal-close" type="button" aria-label="关闭" @click="close"><X :size="17" /></button>
        </div>
        <div class="modal-body">
          <!-- 登录 step -->
          <template v-if="step === 'login'">
            <div class="lm-tabs">
              <button class="lm-tab" :class="{ active: tab === 'wechat' }" type="button" @click="tab = 'wechat'">
                <span class="lm-ico wx">微</span> 微信扫码
              </button>
              <button class="lm-tab" :class="{ active: tab === 'github' }" type="button" @click="tab = 'github'">
                <span class="lm-ico gh"><Github :size="14" /></span> GitHub
              </button>
            </div>

            <div class="lm-panel" :class="{ active: tab === 'wechat' }">
              <div class="wx-qr">
                <Github v-if="false" />
                <span class="wx-qr-ico" style="font-family:'JetBrains Mono',monospace;font-weight:800">CT</span>
              </div>
              <p class="lm-hint">使用微信扫一扫，模拟扫码登录</p>
              <button class="btn btn-primary w-full" type="button" :disabled="loading" @click="doLogin('wechat')">
                <Loader2 v-if="loading" :size="16" class="animate-spin" /> 模拟扫码成功
              </button>
            </div>

            <div class="lm-panel" :class="{ active: tab === 'github' }">
              <p class="lm-hint">使用 GitHub 账号授权登录</p>
              <button class="btn btn-primary w-full" type="button" :disabled="loading" @click="doLogin('github')">
                <Loader2 v-if="loading" :size="16" class="animate-spin" /> GitHub 授权登录
              </button>
            </div>
            <p class="lm-note">登录即代表同意本站服务条款（演示）</p>
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
