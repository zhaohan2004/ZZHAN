/**
 * 后台路由 — /login 公开；其余由 AdminLayout 包裹，未登录重定向到 /login。
 */
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { useAuthStore } from '@/stores/auth'

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue'), meta: { public: true, title: '登录' } },
  {
    path: '/',
    component: AdminLayout,
    children: [
      { path: '', name: 'dashboard', component: () => import('@/views/DashboardView.vue'), meta: { title: '仪表盘' } },
      { path: 'articles', name: 'a-articles', component: () => import('@/views/ArticlesView.vue'), meta: { title: '文章管理' } },
      { path: 'categories', name: 'a-categories', component: () => import('@/views/CategoriesView.vue'), meta: { title: '分类管理' } },
      { path: 'tags', name: 'a-tags', component: () => import('@/views/TagsView.vue'), meta: { title: '标签管理' } },
      { path: 'comments', name: 'a-comments', component: () => import('@/views/CommentsView.vue'), meta: { title: '评论管理' } },
      { path: 'users', name: 'a-users', component: () => import('@/views/UsersView.vue'), meta: { title: '用户管理' } },
      { path: 'operation-logs', name: 'a-operation-logs', component: () => import('@/views/OperationLogsView.vue'), meta: { title: '操作日志' } },
      { path: 'settings', name: 'a-settings', component: () => import('@/views/SettingsView.vue'), meta: { title: '系统设置' } },
    ],
  },
  { path: '/editor', name: 'editor', component: () => import('@/views/EditorView.vue'), meta: { title: '写文章' } },
  { path: '/editor/:id', name: 'editor-edit', component: () => import('@/views/EditorView.vue'), meta: { title: '编辑文章' } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (!to.meta.public && !auth.loggedIn) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
  if (to.path === '/login' && auth.loggedIn) {
    return { path: '/' }
  }
})

router.afterEach((to) => {
  document.title = `${String(to.meta.title || '管理后台')} - 管理后台`
})

export default router
