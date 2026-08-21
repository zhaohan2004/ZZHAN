/**
 * Toast 轻提示 — 在 body 的 `.toast-wrap` 追加 `.toast.toast-{type}`。
 * 2.6s 后加 `.out`，再 320ms 后移除（复刻 assets/js/ui.js showToast 行为）。
 */
import { h, render } from 'vue'
import { Check, Info, X, type LucideIcon } from 'lucide-vue-next'

export type ToastType = 'success' | 'error' | 'info'

const ICONS: Record<ToastType, LucideIcon> = {
  success: Check,
  error: X,
  info: Info,
}

export function useToast() {
  function toast(msg: string, type: ToastType = 'info'): void {
    let wrap = document.querySelector<HTMLElement>('.toast-wrap')
    if (!wrap) {
      wrap = document.createElement('div')
      wrap.className = 'toast-wrap'
      document.body.appendChild(wrap)
    }

    const t = document.createElement('div')
    t.className = 'toast ' + type

    const ico = document.createElement('span')
    ico.className = 't-ico'
    render(h(ICONS[type] ?? Info, { size: 13 }), ico)

    const text = document.createElement('span')
    text.textContent = msg

    t.appendChild(ico)
    t.appendChild(text)
    wrap.appendChild(t)

    setTimeout(() => {
      t.classList.add('out')
      setTimeout(() => t.remove(), 320)
    }, 2600)
  }

  return { toast }
}
