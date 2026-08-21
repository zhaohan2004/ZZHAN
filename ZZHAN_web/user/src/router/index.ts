/**
 * 前台路由 — 7 个懒加载视图 + afterEach 设置 document.title。
 * 视图文件由后续任务补齐，此处用动态 import，构建时解析。
 */
import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useSiteStore } from '@/stores/site'

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('@/views/HomeView.vue'), meta: { title: '首页' } },
  { path: '/articles', name: 'articles', component: () => import('@/views/ArticlesView.vue'), meta: { title: '文章' } },
  { path: '/article/:slug', name: 'article', component: () => import('@/views/ArticleView.vue'), meta: { title: '文章详情' } },
  { path: '/categories', name: 'categories', component: () => import('@/views/CategoriesView.vue'), meta: { title: '分类' } },
  { path: '/tags', name: 'tags', component: () => import('@/views/TagsView.vue'), meta: { title: '标签' } },
  { path: '/archive', name: 'archive', component: () => import('@/views/ArchiveView.vue'), meta: { title: '归档' } },
  { path: '/about', name: 'about', component: () => import('@/views/AboutView.vue'), meta: { title: '关于我' } },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

router.afterEach((to) => {
  const siteStore = useSiteStore()
  const base = siteStore.site?.name ?? 'CodeThink'
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} - ${base}` : base
})

export default router
