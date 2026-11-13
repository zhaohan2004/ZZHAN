<script setup lang="ts">
/** 作者手记，沿用后台的「类型|文本」配置，逐行显示。 */
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { Terminal } from 'lucide-vue-next'
import { useSiteStore } from '@/stores/site'

interface TermLine {
  cls: string
  text: string
}

/** 兜底内容（后台未配置时） */
const DEFAULT_TERMINAL = [
  'tk|$ whoami',
  'fn|zhaohan',
  'tk|',
  'tk|$ cat role.txt',
  'fn|Gopher · 后端工程师',
  'tk|',
  'tk|$ ls tech/',
  'fn|Go  MySQL  Redis  Linux  Docker',
  'tk|',
  'tk|$ echo $MOTTO',
  'cm|写代码是跟计算机对话，写博客是跟自己对话。',
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

const lines = computed(() => parseLines(site.site?.hero_terminal || DEFAULT_TERMINAL))

const visibleCount = ref(0)
let timer: ReturnType<typeof setInterval> | null = null
let motionQuery: MediaQueryList | null = null

function stopReveal(): void {
  if (timer !== null) clearInterval(timer)
  timer = null
}

function startReveal(): void {
  stopReveal()
  visibleCount.value = motionQuery?.matches ? lines.value.length : 0
  if (visibleCount.value >= lines.value.length) return

  timer = setInterval(() => {
    visibleCount.value += 1
    if (visibleCount.value >= lines.value.length) stopReveal()
  }, 430)
}

// 后台配置异步返回或更新时，从新内容的第一行开始，避免沿用旧进度。
watch(lines, startReveal)

onMounted(() => {
  motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
  motionQuery.addEventListener('change', startReveal)
  startReveal()
})

onBeforeUnmount(() => {
  stopReveal()
  motionQuery?.removeEventListener('change', startReveal)
})
</script>

<template>
  <aside class="hero-terminal" aria-labelledby="note-label">
    <h2 id="note-label" class="note-label"><Terminal :size="16" aria-hidden="true" /> 作者手记</h2>
    <div class="note-lines">
      <div
        v-for="(line, index) in lines"
        :key="index"
        :class="[line.cls, { 'note-line-pending': index >= visibleCount }]"
      >{{ line.text }}</div>
    </div>
  </aside>
</template>

<style scoped>
/* 保留尚未显示的行占位，避免逐行出现时卡片和下方文章跳动。 */
.note-line-pending {
  visibility: hidden;
}
</style>
