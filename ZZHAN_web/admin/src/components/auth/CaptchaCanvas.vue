<script setup lang="ts">
/** 图形验证码 — canvas 绘制，点击刷新。 */
import { onMounted, ref } from 'vue'
import { genCaptcha, type CaptchaResult } from '@/utils/captcha'

const emit = defineEmits<{ (e: 'change', text: string): void }>()

const canvas = ref<HTMLCanvasElement | null>(null)
let current: CaptchaResult = genCaptcha()

function draw(): void {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) {
    emit('change', current.text)
    return
  }
  current.draw(ctx, c.width, c.height)
  emit('change', current.text)
}

function refresh(): void {
  current = genCaptcha()
  draw()
}

onMounted(draw)

defineExpose({ refresh })
</script>

<template>
  <canvas
    ref="canvas"
    width="120"
    height="42"
    title="点击刷新验证码"
    style="cursor: pointer; border: 1px solid var(--border-strong); border-radius: 12px"
    @click="refresh"
  />
</template>
