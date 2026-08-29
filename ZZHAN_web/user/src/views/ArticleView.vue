<script setup lang="ts">
/**
 * 文章详情页 — 阅读进度条 / TOC scroll-spy / 点赞分享 / 评论。
 * 布局对齐静态原型：post-header-inner、无封面、无相关推荐、操作仅点赞+分享。
 */
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  Calendar,
  ChevronRight,
  Clock,
  Eye,
  MessageCircle,
} from 'lucide-vue-next'
import { getArticle } from '@/api/articles'
import { getComments } from '@/api/comments'
import { useSiteStore } from '@/stores/site'
import type { ArticleDetail, CommentItem } from '@/types/models'
import { renderMarkdown, buildTOC, decorateCode } from '@/utils/markdown'
import { fmtNum, readTime } from '@/utils/format'
import { useAuthStore } from '@/stores/auth'
import { useReveal } from '@/composables/useReveal'
import TocList, { type TocEntry } from '@/components/widget/TocList.vue'
import ArticleActions from '@/components/widget/ArticleActions.vue'
import CommentList from '@/components/comment/CommentList.vue'
import CommentForm from '@/components/comment/CommentForm.vue'
import CommentGate from '@/components/comment/CommentGate.vue'
import WidgetSidebar from '@/components/widget/WidgetSidebar.vue'

const route = useRoute()
const auth = useAuthStore()

const article = ref<ArticleDetail | null>(null)
const comments = ref<CommentItem[]>([])
const toc = ref<TocEntry[]>([])
const activeId = ref('')
const progress = ref(0)
const postBody = ref<HTMLElement | null>(null)

const root = ref<HTMLElement | null>(null)
useReveal(root)

const onScroll = () => {
  const h = document.documentElement.scrollHeight - window.innerHeight
  progress.value = h > 0 ? (window.scrollY / h) * 100 : 0
}

const catColor = computed(() => '#4a8eff') // 默认主题色

async function loadComments(): Promise<void> {
  if (!article.value) return
  try {
    const res = await getComments(article.value.slug)
    comments.value = res.list.filter((c) => c.parent_id == null)
  } catch {
    comments.value = []
  }
}

onMounted(async () => {
  window.addEventListener('scroll', onScroll, { passive: true })

  const slug = String(route.params.slug)
  try {
    article.value = await getArticle(slug)
    const siteStore = useSiteStore()
    document.title = article.value.title + ' - ' + (siteStore.site?.name ?? '小猫的个人博客')
    await nextTick()
    if (postBody.value) {
      postBody.value.innerHTML = renderMarkdown(article.value.content)
      decorateCode(postBody.value)
      toc.value = buildTOC(postBody.value)
      const io = new IntersectionObserver(
        (entries) => {
          entries.forEach((e) => {
            if (e.isIntersecting) activeId.value = e.target.id
          })
        },
        { rootMargin: '-80px 0px -70% 0px' },
      )
      postBody.value.querySelectorAll('h2,h3,h4').forEach((h) => io.observe(h))
    }
    await loadComments()
  } catch {
    /* 静默 */
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', onScroll)
})

function jumpTo(id: string): void {
  const el = document.getElementById(id)
  el?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}
</script>

<template>
  <div ref="root">
    <div class="progress-bar" :style="{ width: progress + '%' }" />

    <!-- 文章头部 -->
    <header class="post-header">
      <div class="container post-header-inner">
        <div class="flex flex-wrap items-center gap-2.5 text-[13px] text-text-3" style="margin-bottom:16px">
          <router-link to="/" class="muted">首页</router-link>
          <ChevronRight :size="13" />
          <router-link to="/articles" class="muted">文章</router-link>
          <ChevronRight :size="13" />
          <span v-if="article" class="badge" :style="{ background: catColor + '1f', color: catColor, border: '1px solid ' + catColor + '44' }">
            {{ article.category_name }}
          </span>
        </div>
        <h1 class="post-title">{{ article?.title }}</h1>
        <div class="post-summary">{{ article?.summary }}</div>
        <div v-if="article" class="post-meta">
          <span><Calendar :size="14" />发布于 {{ article.published_at }}</span>
          <span><Eye :size="14" />{{ fmtNum(article.views) }} 阅读</span>
          <span><Clock :size="14" />{{ readTime(article.content) }} 分钟</span>
          <span><MessageCircle :size="14" />{{ article.comment_count }} 评论</span>
        </div>
      </div>
    </header>

    <!-- 正文布局 -->
    <div v-if="article" class="container post-layout">
      <div class="post-body">
        <article ref="postBody" class="prose" />

        <ArticleActions :slug="article.slug" :likes="article.likes" />

        <div class="flex flex-wrap gap-2" style="margin-bottom:26px">
          <router-link
            v-for="t in article.tags"
            :key="t.id"
            class="chip"
            :to="`/articles?tag_id=${t.id}`"
          >#{{ t.name }}</router-link>
        </div>

        <!-- 评论区 -->
        <div class="comment-box">
          <div class="section-head">
            <h2 class="section-title" style="font-size:19px"><span class="st-ico"><MessageCircle :size="16" /></span>评论</h2>
            <span class="muted" style="font-size:13px">友善交流，理性发言</span>
          </div>
          <CommentForm v-if="auth.isLoggedIn" :slug="article.slug" @done="loadComments" />
          <CommentGate v-else />
          <CommentList :comments="comments" :slug="article.slug" @done="loadComments" />
        </div>
      </div>

      <!-- 右侧 TOC + 组件 -->
      <aside class="toc-wrap">
        <div class="widget" id="tocWrap">
          <h3 class="widget-title">目录</h3>
          <TocList :items="toc" :active-id="activeId" @jump="jumpTo" />
        </div>
        <WidgetSidebar />
      </aside>
    </div>
  </div>
</template>
