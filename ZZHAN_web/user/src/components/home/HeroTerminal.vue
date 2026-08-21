<script setup lang="ts">
/** 终端卡片 — 打字机逐行展示；内容来自后台「系统设置 → 首页终端内容」（多行 `类型|文本`）。 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useSiteStore } from '@/stores/site'

interface TermLine {
  cls: string
  text: string
}

/** 兜底内容（后台未配置时，与原静态文案一致） */
const DEFAULT_TERMINAL = [
  'tk|go run server.go',
  'cm|# 今天也在认真写代码',
  'fn|[小猫] ',
  'tk|server listening on :8080',
  'tk|curl http://localhost:8080/api/ping',
  'fn|{"pong":true, ',
  'fn|"uptime":"3d 14h"}',
  'tk|git push origin main ✓ done',
].join('\n')

const site = useSiteStore()

/** 解析 `类型|文本` 多行文本 → TermLine[]；纯文本默认 tk，空行忽略 */
function parseLines(raw: string): TermLine[] {
  return raw
    .split('\n')
    .map((l) => {
      const line = l.trim()
      if (!line) return null
      const idx = line.indexOf('|')
      if (idx > 0) {
        const cls = line.slice(0, idx).trim()
        const text = line.slice(idx + 1)
        if (cls === 'tk' || cls === 'cm' || cls === 'fn') return { cls, text }
      }
      return { cls: 'tk', text: line }
    })
    .filter((x): x is TermLine => x !== null)
}

const LINES = computed(() => parseLines(site.site?.heroTerminal || DEFAULT_TERMINAL))

const visible = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const typedLines = computed(() => LINES.value.slice(0, visible.value))

onMounted(() => {
  timer = setInterval(() => {
    if (visible.value < LINES.value.length) visible.value += 1
    else if (timer) clearInterval(timer)
  }, 430)
})
onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="hero-terminal">
    <div style="display:flex;gap:6px;margin-bottom:10px;align-items:center">
      <span class="dot" style="background:#f87171" />
      <span class="dot" style="background:#fbbf24" />
      <span class="dot" style="background:#34d399" />
      <span style="margin-left:8px" class="cm">~/blog · zsh</span>
    </div>
    <div v-for="(line, i) in typedLines" :key="i">
      <span :class="line.cls">{{ line.text }}</span>
    </div>
  </div>
</template>
