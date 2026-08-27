<script setup lang="ts">
/** 单条评论 — 头像 / 名字 / 时间 / 正文 / 点赞 / 回复（楼中楼）。 */
import { ref } from 'vue'
import { ChevronDown, ChevronUp, MessageSquare, ThumbsUp } from 'lucide-vue-next'
import type { CommentItem } from '@/types/models'
import { initialsAvatar } from '@/utils/avatar'
import { toggleCommentLike, getReplies } from '@/api/comments'
import { useAuthStore } from '@/stores/auth'
import CommentForm from './CommentForm.vue'

const props = defineProps<{ comment: CommentItem; slug: string }>()
const emit = defineEmits<{ (e: 'done'): void }>()

const auth = useAuthStore()
const liked = ref(props.comment.liked)
const like_count = ref(props.comment.like_count)
const replying = ref(false)
const replyingTo = ref<number | null>(null)
const busy = ref(false)

// 子评论相关
const replies = ref<CommentItem[]>(props.comment.replies || [])
const replyTotal = ref(props.comment.reply_total || 0)
const hasMore = ref(props.comment.has_more_reply || false)
const loadingMore = ref(false)
const expanded = ref(false)

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

async function loadMore(): Promise<void> {
  if (loadingMore.value) return
  loadingMore.value = true
  try {
    // 加载下一页（每页10条）
    const nextPage = Math.floor(replies.value.length / 10) + 1
    const res = await getReplies(props.comment.id, { page: nextPage, page_size: 10 })
    replies.value = [...replies.value, ...res.list]
    hasMore.value = replies.value.length < (props.comment.reply_total || 0)
    expanded.value = !hasMore.value
  } catch {
    /* 静默 */
  } finally {
    loadingMore.value = false
  }
}

function collapse(): void {
  // 收起：只保留前2条
  replies.value = (props.comment.replies || []).slice(0, 2)
  hasMore.value = (props.comment.reply_total || 0) > 2
  expanded.value = false
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

      <CommentForm v-if="replying" :slug="slug" :parent-id="comment.id" :reply-to="null" class="mt-3" @done="replying = false; emit('done')" />

      <div v-if="replies.length" class="mt-3 space-y-3">
        <div v-for="r in replies" :key="r.id" class="comment-item reply" style="border-bottom:none;padding:10px 0">
          <img class="cm-avatar" :src="avatarOf(r.user_name, r.avatar)" :alt="r.user_name" style="width:34px;height:34px;border-radius:8px" />
          <div style="flex:1;min-width:0">
            <div class="cm-head">
              <span class="cm-name">{{ r.user_name }}</span>
              <span class="cm-time">{{ r.time }}</span>
            </div>
            <div class="cm-text">{{ r.content }}</div>
            <div class="cm-actions">
              <button type="button" @click="replyingTo = replyingTo === r.id ? null : r.id">
                <MessageSquare :size="13" /> 回复
              </button>
            </div>
            <CommentForm v-if="replyingTo === r.id" :slug="slug" :parent-id="comment.id" :reply-to="r.user_name" class="mt-3" @done="replyingTo = null; emit('done')" />
          </div>
        </div>

        <!-- 展开/收起按钮 -->
        <div class="expand-actions">
          <button
            v-if="hasMore && !expanded"
            type="button"
            class="expand-btn"
            :disabled="loadingMore"
            @click="loadMore"
          >
            <ChevronDown :size="14" />
            {{ loadingMore ? '加载中...' : `展开更多回复（共${replyTotal}条）` }}
          </button>
          <button
            v-if="expanded"
            type="button"
            class="expand-btn"
            @click="collapse"
          >
            <ChevronUp :size="14" />
            收起回复
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.expand-actions {
  display: flex;
  gap: 12px;
  margin-top: 4px;
}
.expand-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  font-size: 13px;
  color: var(--accent);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
}
.expand-btn:hover {
  opacity: 0.8;
}
.expand-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
