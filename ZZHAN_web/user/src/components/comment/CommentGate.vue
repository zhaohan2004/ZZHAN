<script setup lang="ts">
/** 未登录评论门槛 — 结构对齐静态原型 .comment-gate（横排提示 + 登录按钮）。 */
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
</script>

<template>
  <div class="comment-gate">
    <div class="cg-text">
      登录后即可发表评论，<b>点击右侧按钮登录</b>
    </div>
    <button class="btn btn-primary btn-sm" type="button" @click="auth.ensureAuth()">立即登录</button>
  </div>
</template>
