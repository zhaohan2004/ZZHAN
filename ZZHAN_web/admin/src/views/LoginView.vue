<script setup lang="ts">
/** 后台登录页 — 服务端 base64 验证码 */
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Loader2, Lock, LogIn, Moon, RefreshCw, Sun, User } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { useSettingsStore } from '@/stores/settings'
import { toast } from '@/composables/useToast'
import { getCaptcha } from '@/api/admin'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const theme = useThemeStore()
const settings = useSettingsStore()

const username = ref('admin')
const password = ref('123456')
const captchaCode = ref('')
const remember = ref(true)
const loading = ref(false)

// 服务端验证码
const captchaId = ref('')
const captchaB64 = ref('')

async function loadCaptcha(): Promise<void> {
  try {
    const res = await getCaptcha()
    captchaId.value = res.captcha_id
    captchaB64.value = res.captcha_image
    captchaCode.value = ''
  } catch {
    toast.error('获取验证码失败')
  }
}

onMounted(loadCaptcha)

async function onLogin(): Promise<void> {
  if (!username.value.trim() || !password.value.trim()) {
    toast.error('请输入用户名和密码')
    return
  }
  if (!captchaCode.value.trim()) {
    toast.error('请输入验证码')
    return
  }
  loading.value = true
  try {
    await auth.login(username.value.trim(), password.value.trim(), captchaId.value, captchaCode.value.trim())
    toast.success('登录成功，欢迎回来')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e) {
    toast.error((e as Error).message || '登录失败')
    await loadCaptcha()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-body">
    <div class="hero-grid-bg" style="position: fixed" />
    <div class="hero-orb a" style="position: fixed" />
    <div class="hero-orb b" style="position: fixed" />

    <div style="position: fixed; top: 22px; right: 26px; z-index: 5; display: flex; gap: 8px; align-items: center">
      <button class="icon-btn" type="button" :title="theme.dark ? '切换浅色' : '切换深色'" @click="theme.toggle()">
        <component :is="theme.dark ? Sun : Moon" :size="18" />
      </button>
    </div>

    <div class="login-card anim-fade">
      <div class="login-logo">
        <span class="brand-logo">{{ settings.settings?.logo_text || 'CT' }}</span>
        <h1>{{ settings.settings?.blog_name || 'Blog' }} <span class="grad-text">管理后台</span></h1>
        <p>{{ settings.settings?.blog_desc || '后台管理系统' }}</p>
      </div>
      <form novalidate @submit.prevent="onLogin">
        <div class="form-group">
          <label class="form-label" for="loginUser">用户名 <span class="req">*</span></label>
          <div style="position: relative">
            <User :size="17" style="position: absolute; left: 13px; top: 50%; transform: translateY(-50%); color: var(--text-3)" />
            <input id="loginUser" v-model="username" class="input" type="text" placeholder="请输入用户名" style="padding-left: 40px" autocomplete="username">
          </div>
        </div>
        <div class="form-group">
          <label class="form-label" for="loginPass">密码 <span class="req">*</span></label>
          <div style="position: relative">
            <Lock :size="17" style="position: absolute; left: 13px; top: 50%; transform: translateY(-50%); color: var(--text-3)" />
            <input id="loginPass" v-model="password" class="input" type="password" placeholder="请输入密码" style="padding-left: 40px" autocomplete="current-password">
          </div>
        </div>
        <div class="form-group">
          <label class="form-label" for="loginCaptcha">图形验证码 <span class="req">*</span></label>
          <div class="captcha-row">
            <input id="loginCaptcha" v-model="captchaCode" class="input" type="text" maxlength="5" placeholder="验证码" style="text-transform: uppercase; font-family: 'JetBrains Mono', monospace" autocomplete="off">
            <img
              v-if="captchaB64"
              :src="captchaB64"
              alt="验证码"
              title="点击刷新验证码"
              style="cursor: pointer; height: 42px; border-radius: 12px; border: 1px solid var(--border-strong)"
              @click="loadCaptcha"
            >
            <div
              v-else
              style="width: 120px; height: 42px; border-radius: 12px; border: 1px solid var(--border-strong); display: flex; align-items: center; justify-content: center; color: var(--text-3); font-size: 12px"
            >
              加载中...
            </div>
          </div>
          <div class="form-hint"><RefreshCw :size="12" style="vertical-align: -2px" /> 点击图片可刷新验证码</div>
        </div>
        <div style="display: flex; align-items: center; justify-content: space-between; margin: 4px 0 20px">
          <label style="display: flex; align-items: center; gap: 8px; font-size: 13px; color: var(--text-2); cursor: pointer">
            <input v-model="remember" type="checkbox" style="accent-color: var(--accent); width: 15px; height: 15px"> 记住我
          </label>
          <a href="javascript:;" style="font-size: 13px; color: var(--accent)" @click="toast.info('演示环境，无找回密码功能')">忘记密码？</a>
        </div>
        <button class="btn btn-primary" type="submit" style="width: 100%; padding: 12px" :disabled="loading">
          <Loader2 v-if="loading" :size="17" class="spin" />
          <LogIn v-else :size="17" />
          {{ loading ? '登录中...' : '登 录' }}
        </button>
      </form>
      <div class="form-hint" style="text-align: center; margin-top: 16px; line-height: 1.9">
        推荐账号：admin / 123456
      </div>
    </div>
  </div>
</template>

<style scoped>
.spin { animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
