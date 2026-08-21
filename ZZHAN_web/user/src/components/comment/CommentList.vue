<script setup lang="ts">
/** 评论列表 — 顶层评论（parent_id 为 null）+ 楼中楼 replies。 */
import type { CommentItem as CommentModel } from '@/types/models'
import CommentItem from './CommentItem.vue'

defineProps<{ comments: CommentModel[]; slug: string }>()
</script>

<template>
  <div v-if="comments.length">
    <CommentItem v-for="c in comments" :key="c.id" :comment="c" :slug="slug" />
  </div>
  <p v-else class="muted" style="padding:18px 0">还没有评论，快来抢沙发～</p>
</template>
