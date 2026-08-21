<script setup lang="ts">
/** 应用外壳 — 导航 + 路由视图 + 回顶 + 登录弹窗（无页脚，对齐静态原型）。 */
import { onMounted } from 'vue'
import SiteNav from '@/components/layout/SiteNav.vue'
import BackTop from '@/components/layout/BackTop.vue'
import LoginModal from '@/components/auth/LoginModal.vue'
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()

onMounted(() => {
  site.fetchSite()
})
</script>

<template>
  <SiteNav />
  <main class="min-h-screen">
    <router-view />
  </main>
  <BackTop />
  <LoginModal />
</template>
