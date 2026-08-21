/**
 * 全局 toast / confirm — 替代 Element Plus 的 ElMessage / ElMessageBox。
 * 后台去 EP 后统一走这里，UI 用 style.css 的 .modal / .btn 等原生类。
 */
import { reactive, readonly } from 'vue'

export type ToastType = 'success' | 'error' | 'info' | 'warning'
export interface ToastItem {
  id: number
  type: ToastType
  message: string
}

const state = reactive<{ toasts: ToastItem[] }>({ toasts: [] })
let seq = 0

function push(type: ToastType, message: string, duration = 2600): void {
  const id = ++seq
  state.toasts.push({ id, type, message })
  window.setTimeout(() => {
    const idx = state.toasts.findIndex((t) => t.id === id)
    if (idx >= 0) state.toasts.splice(idx, 1)
  }, duration)
}

export const toast = {
  success: (m: string) => push('success', m),
  error: (m: string) => push('error', m),
  info: (m: string) => push('info', m),
  warning: (m: string) => push('warning', m),
}

export function useToasts() {
  return readonly(state).toasts
}

interface ConfirmState {
  open: boolean
  title: string
  message: string
  resolve: (v: boolean) => void
}

const confirmState = reactive<ConfirmState>({
  open: false,
  title: '确认操作',
  message: '',
  resolve: () => {},
})

export function confirm(message: string, title = '确认操作'): Promise<boolean> {
  return new Promise<boolean>((resolve) => {
    confirmState.open = true
    confirmState.title = title
    confirmState.message = message
    confirmState.resolve = resolve
  })
}

export function useConfirm() {
  return confirmState
}
