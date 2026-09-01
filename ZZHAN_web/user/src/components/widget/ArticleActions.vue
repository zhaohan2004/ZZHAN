<script setup lang="ts">
/**
 * 文章底部操作 — 点赞（未登录弹登录门槛）+ 分享（复制链接）。对齐静态版（无收藏）。
 */
import { ref } from 'vue'
import { Share2, ThumbsUp } from 'lucide-vue-next'
import { toggleLike } from '@/api/articles'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

const props = defineProps<{ slug: string; likes: number; liked?: boolean }>()
const emit = defineEmits<{ (e: 'liked', likes: number): void }>()

const auth = useAuthStore()
const { toast } = useToast()
const liked = ref(props.liked ?? false)
const like_count = ref(props.likes)
const loading = ref(false)

async function onLike(): Promise<void> {
  if (!(await auth.ensureAuth())) return
  if (loading.value) return
  loading.value = true
  try {
    const res = await toggleLike(props.slug)
    liked.value = res.liked
    like_count.value = res.likes
    emit('liked', res.likes)
    toast(res.liked ? '点赞成功' : '已取消点赞', res.liked ? 'success' : 'info')
  } catch {
    toast('操作失败，请重试', 'error')
  } finally {
    loading.value = false
  }
}

async function onShare(): Promise<void> {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
    toast('链接已复制，去分享吧', 'success')
  } catch {
    toast(url, 'info')
  }
}
</script>

<template>
  <div class="post-actions">
    <button class="post-action-btn" :class="{ on: liked }" type="button" @click="onLike">
      <ThumbsUp :size="17" /> 点赞 <span>{{ like_count }}</span>
    </button>
    <button class="post-action-btn" type="button" @click="onShare">
      <Share2 :size="17" /> 分享
    </button>
  </div>
</template>
