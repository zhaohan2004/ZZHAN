<script setup lang="ts">
/** 全局 toast 容器 + confirm 弹窗（替代 Element Plus）。 */
import { AlertTriangle, CheckCircle, Info, XCircle } from 'lucide-vue-next'
import { useConfirm, useToasts } from '@/composables/useToast'

const toasts = useToasts()
const confirmState = useConfirm()

const icons = { success: CheckCircle, error: XCircle, info: Info, warning: AlertTriangle }

function onOk(): void {
  confirmState.resolve(true)
  confirmState.open = false
}
function onCancel(): void {
  confirmState.resolve(false)
  confirmState.open = false
}
</script>

<template>
  <div class="toast-wrap">
    <transition-group name="tfade">
      <div v-for="t in toasts" :key="t.id" class="toast" :class="'toast-' + t.type">
        <component :is="icons[t.type]" :size="17" />
        <span>{{ t.message }}</span>
      </div>
    </transition-group>
  </div>

  <div v-if="confirmState.open" class="modal-overlay open" @click.self="onCancel">
    <div class="modal">
      <div class="modal-head"><h3>{{ confirmState.title }}</h3></div>
      <div class="modal-body">
        <p style="color: var(--text-2); font-size: 14px; line-height: 1.8">{{ confirmState.message }}</p>
      </div>
      <div class="modal-foot">
        <button class="btn btn-ghost btn-sm" type="button" @click="onCancel">取消</button>
        <button class="btn btn-danger btn-sm" type="button" @click="onOk">确认</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.toast-wrap {
  position: fixed;
  top: 18px;
  right: 18px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  pointer-events: none;
}
.toast {
  display: flex;
  align-items: center;
  gap: 9px;
  min-width: 220px;
  max-width: 360px;
  padding: 12px 16px;
  border-radius: 12px;
  font-size: 13.5px;
  font-weight: 500;
  color: var(--text);
  background: var(--glass-strong);
  border: 1px solid var(--border-strong);
  backdrop-filter: blur(18px);
  box-shadow: var(--shadow-lg);
  pointer-events: auto;
}
.toast-success { color: var(--success); }
.toast-error { color: var(--danger); }
.toast-info { color: var(--info); }
.toast-warning { color: var(--warning); }
.tfade-enter-active,
.tfade-leave-active { transition: all 0.3s ease; }
.tfade-enter-from,
.tfade-leave-to { opacity: 0; transform: translateX(20px); }
</style>
