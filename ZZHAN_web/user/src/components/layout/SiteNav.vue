<script setup lang="ts">
/** 前台导航 — 品牌 + 6 链接 + 搜索/主题 + 登录用户区（头像下拉：退出）。无 GitHub 外链、无个人资料编辑。 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ChevronDown, LogIn, LogOut, Menu, Search, X } from 'lucide-vue-next'
import { useSiteStore } from '@/stores/site'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { initialsAvatar } from '@/utils/avatar'
import ThemeToggle from '@/components/ui/ThemeToggle.vue'

const route = useRoute()
const router = useRouter()
const site = useSiteStore()
const auth = useAuthStore()
const { toast } = useToast()

const scrolled = ref(false)
const menuOpen = ref(false)
const searchOpen = ref(false)
const searchInput = ref('')

const siteName = computed(() => site.site?.name ?? '小猫的个人博客')
const logoText = computed(() => site.site?.logo_text ?? 'CT')

const navLinks = [
  { to: '/', label: '首页', exact: true },
  { to: '/articles', label: '文章', exact: false },
  { to: '/categories', label: '分类', exact: false },
  { to: '/tags', label: '标签', exact: false },
  { to: '/archive', label: '归档', exact: false },
  { to: '/about', label: '关于我', exact: false },
]

function isActive(l: { to: string; exact: boolean }): boolean {
  return l.exact ? route.path === l.to : route.path.startsWith(l.to)
}

function onScroll(): void {
  scrolled.value = window.scrollY > 10
}

function submitSearch(): void {
  const q = searchInput.value.trim()
  if (!q) {
    toast('请输入搜索关键词', 'info')
    return
  }
  router.push({ path: '/articles', query: { q } })
  searchInput.value = ''
  searchOpen.value = false
  menuOpen.value = false
}

/* ---------- 用户区：头像下拉（退出登录） ---------- */
const userMenuOpen = ref(false)
const userAvatar = computed(() => {
  if (auth.user?.avatar) return auth.user.avatar
  return initialsAvatar(auth.user?.nickname || '用', '#3b82f6', '#38bdf8', 64)
})
const providerLabel = computed(() =>
  auth.user?.provider === 'github' ? 'GitHub 用户' : '访客',
)

function onDocClick(): void {
  userMenuOpen.value = false
}

/* 退出确认弹窗 */
const logoutOpen = ref(false)
function openLogout(): void {
  userMenuOpen.value = false
  logoutOpen.value = true
}
async function doLogout(): Promise<void> {
  await auth.logout()
  logoutOpen.value = false
  toast('已退出登录', 'success')
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  document.addEventListener('click', onDocClick)
  onScroll()
})
onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
  document.removeEventListener('click', onDocClick)
})
</script>

<template>
  <header class="nav" :class="{ scrolled }">
    <div class="container nav-inner">
      <router-link to="/" class="brand">
        <span class="brand-logo">{{ logoText }}</span>
        <span class="brand-text">
          {{ siteName }}<span class="brand-dot">.</span>
        </span>
      </router-link>

      <nav class="nav-links" :class="{ open: menuOpen }">
        <router-link
          v-for="l in navLinks"
          :key="l.to"
          :to="l.to"
          class="nav-link"
          :class="{ active: isActive(l) }"
          @click="menuOpen = false"
        >
          {{ l.label }}
        </router-link>
      </nav>

      <div class="nav-actions">
        <button class="icon-btn" type="button" title="搜索" aria-label="搜索" @click="searchOpen = !searchOpen">
          <Search :size="19" />
        </button>
        <ThemeToggle />
        <!-- 登录用户区：已登录显示头像下拉，未登录显示登录按钮 -->
        <div v-if="auth.isLoggedIn" class="admin-user" :class="{ open: userMenuOpen }">
          <button class="icon-btn user-avatar-btn" type="button" title="账号" aria-label="账号" @click.stop="userMenuOpen = !userMenuOpen">
            <img :src="userAvatar" alt="头像" />
            <ChevronDown class="au-caret" :size="12" />
          </button>
          <div class="nav-dd" v-show="userMenuOpen" @click.stop>
            <div class="au-head">
              <img :src="userAvatar" alt="" />
              <div><b>{{ auth.user?.nickname || '用户' }}</b><span>{{ providerLabel }}</span></div>
            </div>
            <div class="dd-sep"></div>
            <button class="nav-dd-item danger" type="button" @click="openLogout"><LogOut :size="15" /> 退出登录</button>
          </div>
        </div>
        <button v-else class="btn btn-primary" type="button" style="height:34px;padding:0 15px;font-size:13px" @click="auth.ensureAuth()">
          <LogIn :size="14" /> 登录
        </button>
        <button class="icon-btn menu-btn" type="button" title="菜单" aria-label="菜单" @click="menuOpen = !menuOpen">
          <X v-if="menuOpen" :size="20" />
          <Menu v-else :size="20" />
        </button>
      </div>
    </div>

    <div class="search-bar" :class="{ open: searchOpen }">
      <div class="container search-inner">
        <input
          v-model="searchInput"
          class="input"
          type="search"
          placeholder="搜索文章标题、标签、内容… 回车跳转"
          aria-label="搜索关键词"
          @keyup.enter="submitSearch"
        />
        <button class="btn btn-primary" type="button" @click="submitSearch">搜索</button>
      </div>
    </div>
  </header>

  <Teleport to="body">
    <!-- 退出确认弹窗 -->
    <div class="modal-overlay" :class="{ open: logoutOpen }" @click.self="logoutOpen = false">
      <div class="modal">
        <div class="modal-head">
          <h3>退出登录</h3>
          <button class="modal-close" type="button" @click="logoutOpen = false"><X :size="17" /></button>
        </div>
        <div class="modal-body">
          <p style="font-size:14px;color:var(--text-2);line-height:1.7">确定要退出当前账号吗？退出后评论、点赞等操作需要重新登录。</p>
        </div>
        <div class="modal-foot">
          <button class="btn btn-ghost btn-sm" type="button" @click="logoutOpen = false">取消</button>
          <button class="btn btn-danger btn-sm" type="button" @click="doLogout">退出登录</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
