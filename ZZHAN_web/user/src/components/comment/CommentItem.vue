<script setup lang="ts">
/** 单条评论 — 头像 / 名字 / 时间 / 正文 / 点赞 / 回复（楼中楼）。 */
import { ref } from 'vue'
import { MessageSquare, ThumbsUp } from 'lucide-vue-next'
import type { CommentItem } from '@/types/models'
import { initialsAvatar } from '@/utils/avatar'
import { toggleCommentLike } from '@/api/comments'
import { useAuthStore } from '@/stores/auth'
import CommentForm from './CommentForm.vue'

const props = defineProps<{ comment: CommentItem; slug: string }>()
const emit = defineEmits<{ (e: 'done'): void }>()

const auth = useAuthStore()
const liked = ref(props.comment.liked)
const like_count = ref(props.comment.like_count)
const replying = ref(false)
const busy = ref(false)

function avatarOf(name: string, avatar?: string): string {
  if (avatar) return avatar
  return initialsAvatar(name, '#6b7280', '#6b7280', 84)
}

async function onLike(): Promise<void> {
  if (!(await auth.ensureAuth())) return
  if (busy.value) return
  busy.value = true
  try {
    const res = await toggleCommentLike(props.comment.id)
    liked.value = res.liked
    like_count.value = res.like_count
  } catch {
    /* 静默 */
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="comment-item">
    <img class="cm-avatar" :src="avatarOf(comment.user_name, comment.avatar)" :alt="comment.user_name" />
    <div style="flex:1;min-width:0">
      <div class="cm-head">
        <span class="cm-name">{{ comment.user_name }}</span>
        <span class="cm-time">{{ comment.time }}</span>
      </div>
      <div class="cm-text">{{ comment.content }}</div>
      <div class="cm-actions">
        <button type="button" :style="liked ? 'color:var(--accent)' : ''" @click="onLike">
          <ThumbsUp :size="13" /> {{ like_count }}
        </button>
        <button type="button" @click="replying = !replying">
          <MessageSquare :size="13" /> 回复
        </button>
      </div>

      <CommentForm v-if="replying" :slug="slug" :parent-id="comment.id" class="mt-3" @done="replying = false; emit('done')" />

      <div v-if="comment.replies?.length" class="mt-3 space-y-3">
        <div v-for="r in comment.replies" :key="r.id" class="comment-item reply" style="border-bottom:none;padding:10px 0">
          <img class="cm-avatar" :src="avatarOf(r.user_name, r.avatar)" :alt="r.user_name" style="width:34px;height:34px;border-radius:10px" />
          <div style="flex:1;min-width:0">
            <div class="cm-head">
              <span class="cm-name">{{ r.user_name }}</span>
              <span class="cm-time">{{ r.time }}</span>
            </div>
            <div class="cm-text">{{ r.content }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
