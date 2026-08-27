<script setup lang="ts">
/** 导航栏用户区 — 未登录显示登录按钮；已登录头像下拉（完善资料 / 退出）。 */
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { LogOut, UserRound } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { initialsAvatar } from '@/utils/avatar'

const auth = useAuthStore()
const { toast } = useToast()
const menuOpen = ref(false)

function avatarFor(): string {
  const u = auth.user
  if (!u) return initialsAvatar('?', '#3b82f6', '#38bdf8', 64)
  if (u.avatar) return u.avatar
  const color = '#24292f'
  return initialsAvatar(u.nickname || '我', color, color, 64)
}

function closeOnOutside(e: MouseEvent): void {
  const t = e.target as HTMLElement
  if (!t.closest('.user-drop')) menuOpen.value = false
}

async function onLogout(): Promise<void> {
  await auth.logout()
  menuOpen.value = false
  toast('已退出登录', 'info')
}

onMounted(() => document.addEventListener('click', closeOnOutside))
onBeforeUnmount(() => document.removeEventListener('click', closeOnOutside))
</script>

<template>
  <div class="user-drop">
    <!-- 未登录 -->
    <button v-if="!auth.isLoggedIn" class="nav-login ml-1" type="button" @click="auth.ensureAuth()">登录</button>

    <!-- 已登录 -->
    <template v-else>
      <button class="icon-btn ml-1" type="button" aria-label="用户菜单" @click.stop="menuOpen = !menuOpen">
        <img :src="avatarFor()" alt="头像" style="width:32px;height:32px;border-radius:9px" />
      </button>
      <div v-if="menuOpen" class="ud-menu">
        <div class="px-3 pb-2 pt-1">
          <div class="text-[14px] font-semibold text-text">{{ auth.user?.nickname || '用户' }}</div>
          <div class="text-[12px] text-text-3">GitHub 用户</div>
        </div>
        <button class="ud-item" type="button" @click="auth.loginModalOpen = true; menuOpen = false">
          <UserRound :size="15" /> 完善资料
        </button>
        <button class="ud-item" type="button" @click="onLogout">
          <LogOut :size="15" style="color:var(--danger)" /> 退出登录
        </button>
      </div>
    </template>
  </div>
</template>
