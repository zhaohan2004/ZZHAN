<script setup lang="ts">
/**
 * 后台布局壳 — 复刻静态 .admin-layout（.admin-sidebar + .admin-main > .admin-header + .admin-content）。
 * 侧栏折叠状态在此管理，传给 SideNav / AdminHeader。
 */
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import SideNav from './SideNav.vue'
import AdminHeader from './AdminHeader.vue'

const sidebarOpen = ref(false)
function toggleSidebar(): void {
  sidebarOpen.value = !sidebarOpen.value
}
</script>

<template>
  <div class="admin-layout">
    <SideNav :open="sidebarOpen" @close="sidebarOpen = false" />
    <div class="admin-main">
      <AdminHeader @toggle="toggleSidebar" />
      <main class="admin-content">
        <RouterView />
      </main>
    </div>
  </div>
</template>
