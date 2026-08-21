/**
 * 滚动渐入动画 — 观察根元素内的 `.reveal`，进入视口时加 `.in`。
 * 使用方式：`const root = ref<HTMLElement|null>(null); useReveal(root)`。
 */
import { onBeforeUnmount, onMounted, type Ref } from 'vue'

export function useReveal(el: Ref<HTMLElement | null>): void {
  let io: IntersectionObserver | null = null

  function observe(): void {
    const root = el.value
    if (!root) return
    io = new IntersectionObserver(
      (entries) => {
        entries.forEach((e) => {
          if (e.isIntersecting) {
            e.target.classList.add('in')
            io?.unobserve(e.target)
          }
        })
      },
      { threshold: 0.08 },
    )
    root.querySelectorAll('.reveal:not(.in)').forEach((n) => io?.observe(n))
  }

  onMounted(() => {
    observe()
    // 视图内异步渲染完成后可能新增 .reveal 节点
    const timer = setTimeout(observe, 300)
    onBeforeUnmount(() => {
      clearTimeout(timer)
      io?.disconnect()
    })
  })
}
