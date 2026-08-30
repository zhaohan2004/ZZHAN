<script setup lang="ts">
/**
 * 后台顶栏 — 复刻静态 .admin-header（菜单按钮 / 标题 / 主题切换 / 头像下拉）。
 * 去 Element Plus：用自写 toast/confirm + lucide-vue-next 图标。
 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChevronDown, LogOut, Menu, Moon, Sun, Trash2, Upload, UserCog, X } from 'lucide-vue-next'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { useSettingsStore } from '@/stores/settings'
import { confirm, toast } from '@/composables/useToast'
import { initialsAvatar } from '@/utils/cover'
import { uploadImage } from '@/api/admin'
import type { AdminProfile } from '@/types/models'

const emit = defineEmits<{ toggle: [] }>()

const route = useRoute()
const router = useRouter()
const theme = useThemeStore()
const auth = useAuthStore()
const settings = useSettingsStore()
const menuOpen = ref(false)

const title = computed(() => String(route.meta.title || '管理后台'))

const avatar = computed(() => {
  const p = auth.profile
  if (p?.avatar) return p.avatar
  return initialsAvatar(p?.nickname || '轩', '#3b82f6', '#38bdf8', 64)
})

const profileOpen = ref(false)
const editUsername = ref('')
const editPassword = ref('')
const editPassword2 = ref('')
const avatarUrl = ref('')
const avatarFile = ref<File | null>(null)
function openProfile(): void {
  menuOpen.value = false
  editUsername.value = auth.profile?.username || ''
  editPassword.value = ''
  editPassword2.value = ''
  avatarUrl.value = auth.profile?.avatar || ''
  avatarFile.value = null
  profileOpen.value = true
}
/** 头像上传 — 记录文件对象，预览用临时 URL */
function onAvatarFile(e: Event): void {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    toast.error('头像文件不能超过 2MB')
    input.value = ''
    return
  }
  avatarFile.value = file
  avatarUrl.value = URL.createObjectURL(file)
}
function clearAvatar(): void {
  avatarUrl.value = ''
  avatarFile.value = null
}
async function saveProfile(): Promise<void> {
  if (!auth.profile) return
  const pwd = editPassword.value
  if (pwd && pwd.length < 6) {
    toast.error('新密码至少 6 位')
    return
  }
  if (pwd && pwd !== editPassword2.value) {
    toast.error('两次输入的密码不一致')
    return
  }
  try {
    let finalAvatar = avatarUrl.value
    // 如果选择了新头像文件，先上传获取 URL
    if (avatarFile.value) {
      const res = await uploadImage(avatarFile.value)
      finalAvatar = res.url
    }
    const payload: AdminProfile & { password?: string } = {
      ...auth.profile,
      username: editUsername.value,
      avatar: finalAvatar,
    }
    if (pwd) payload.password = pwd
    await auth.updateProfile(payload)
    profileOpen.value = false
    editPassword.value = ''
    editPassword2.value = ''
    avatarFile.value = null
    toast.success('资料已保存')
  } catch {
    toast.error('保存失败')
  }
}

async function onLogout(): Promise<void> {
  const ok = await confirm('确定要退出登录吗？', '退出确认')
  if (!ok) return
  menuOpen.value = false
  await auth.logout()
  toast.success('已退出登录')
  router.push('/login')
}

function onDocClick(): void {
  menuOpen.value = false
}
onMounted(() => document.addEventListener('click', onDocClick))
onBeforeUnmount(() => document.removeEventListener('click', onDocClick))
</script>

<template>
  <header class="admin-header">
    <button class="icon-btn" type="button" aria-label="菜单" @click="emit('toggle')"><Menu :size="19" /></button>
    <div>
      <div class="h-title"><span style="color: var(--accent)">{{ title }}</span></div>
      <div class="h-crumb">{{ settings.settings?.blog_name || 'Blog' }} 管理后台</div>
    </div>
    <div class="h-right">
      <button class="icon-btn" type="button" :title="theme.dark ? '切换浅色' : '切换深色'" @click="theme.toggle()">
        <component :is="theme.dark ? Sun : Moon" :size="18" />
      </button>
      <template v-if="auth.loggedIn">
        <div class="admin-user" :class="{ open: menuOpen }">
          <button class="icon-btn" type="button" title="账号" @click.stop="menuOpen = !menuOpen">
            <img :src="avatar" alt="头像" style="width: 32px; height: 32px; border-radius: 9px">
            <ChevronDown class="au-caret" :size="12" />
          </button>
          <div class="nav-dd" v-show="menuOpen" @click.stop>
            <div class="au-head">
              <img :src="avatar" alt="">
              <div><b>{{ auth.profile?.nickname || '管理员' }}</b><span>超级管理员</span></div>
            </div>
            <div class="dd-sep"></div>
            <button class="nav-dd-item" type="button" @click="openProfile"><UserCog :size="15" /> 个人信息</button>
            <div class="dd-sep"></div>
            <button class="nav-dd-item danger" type="button" @click="onLogout"><LogOut :size="15" /> 退出登录</button>
          </div>
        </div>
      </template>
      <template v-else>
        <button class="btn btn-primary btn-sm" type="button" @click="router.push('/login')">登录</button>
      </template>
    </div>
  </header>

  <div class="modal-overlay" :class="{ open: profileOpen }" @click.self="profileOpen = false">
    <div class="modal">
      <div class="modal-head">
        <h3>个人信息</h3>
        <button class="modal-close" type="button" @click="profileOpen = false"><X :size="17" /></button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label class="form-label">头像</label>
          <div style="display:flex;gap:12px;align-items:center;flex-wrap:wrap">
            <span class="brand-logo" style="width:52px;height:52px;font-size:20px;flex:none;overflow:hidden;padding:0;border-radius:12px">
              <img v-if="avatarUrl" :src="avatarUrl" alt="头像预览" style="width:100%;height:100%;object-fit:cover" />
              <span v-else>{{ (auth.profile?.nickname || 'U').slice(0, 1) }}</span>
            </span>
            <div style="display:flex;gap:8px;align-items:center;flex-wrap:wrap">
              <label class="btn btn-ghost btn-sm" style="cursor:pointer;margin:0">
                <Upload :size="14" /> 上传头像
                <input type="file" accept="image/*" style="display:none" @change="onAvatarFile" />
              </label>
              <button v-if="avatarUrl" type="button" class="btn btn-ghost btn-sm" @click="clearAvatar">
                <Trash2 :size="14" /> 清除
              </button>
              <span class="form-hint" style="margin:0">JPG/PNG/WebP，最大 2MB</span>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">账号</label>
          <input class="input" v-model="editUsername" type="text" placeholder="登录账号" />
        </div>
        <div class="form-group">
          <label class="form-label">新密码 <span class="form-hint" style="display:inline;margin-left:6px">留空则不修改</span></label>
          <input class="input" v-model="editPassword" type="password" placeholder="至少 6 位" autocomplete="new-password" />
        </div>
        <div class="form-group" style="margin-bottom:0">
          <label class="form-label">确认新密码</label>
          <input class="input" v-model="editPassword2" type="password" placeholder="再次输入新密码" autocomplete="new-password" />
        </div>
      </div>
      <div class="modal-foot">
        <button class="btn btn-ghost btn-sm" type="button" @click="profileOpen = false">取消</button>
        <button class="btn btn-primary btn-sm" type="button" @click="saveProfile">保存</button>
      </div>
    </div>
  </div>
</template>
