<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { CodeXml, Link, MousePointer2, X } from 'lucide-vue-next'
import { useClipboard } from '@/composables/useClipboard'

const { manualCopy, closeManualCopy } = useClipboard()
const dialog = ref<HTMLDialogElement | null>(null)
const content = ref<HTMLTextAreaElement | null>(null)
const isCode = computed(() => manualCopy.value?.kind === 'code')
const kindLabel = computed(() => isCode.value ? '代码' : '链接')

function selectContent(): void {
  content.value?.focus({ preventScroll: true })
  content.value?.select()
}

watch(manualCopy, async (value) => {
  if (!value) {
    dialog.value?.close()
    return
  }
  await nextTick()
  if (!manualCopy.value || !dialog.value) return
  if (!dialog.value.open) dialog.value.showModal()
  selectContent()
}, { flush: 'post' })

function closeOnBackdrop(event: MouseEvent): void {
  if (event.target !== dialog.value || !dialog.value) return
  const rect = dialog.value.getBoundingClientRect()
  if (event.clientX < rect.left || event.clientX > rect.right || event.clientY < rect.top || event.clientY > rect.bottom) {
    closeManualCopy()
  }
}

onBeforeUnmount(() => dialog.value?.close())
</script>

<template>
  <Teleport to="body">
    <dialog
      ref="dialog"
      class="modal manual-copy-dialog"
      aria-labelledby="manual-copy-title"
      aria-describedby="manual-copy-description"
      @cancel.prevent="closeManualCopy"
      @close="closeManualCopy"
      @click="closeOnBackdrop"
    >
      <div class="modal-head">
        <div class="manual-copy-heading">
          <span class="manual-copy-icon" aria-hidden="true">
            <CodeXml v-if="isCode" :size="20" />
            <Link v-else :size="20" />
          </span>
          <h3 id="manual-copy-title">请手动复制{{ kindLabel }}</h3>
        </div>
        <button class="modal-close" type="button" aria-label="关闭手动复制窗口" @click="closeManualCopy">
          <X :size="18" />
        </button>
      </div>
      <div class="modal-body manual-copy-body">
        <p id="manual-copy-description" class="manual-copy-description">自动复制未成功，请手动复制下方{{ kindLabel }}。</p>
        <label class="form-label" for="manual-copy-content">{{ isCode ? '代码内容' : '文章链接' }}</label>
        <textarea
          id="manual-copy-content"
          ref="content"
          class="manual-copy-content resize-none"
          :class="{ 'is-code': isCode }"
          :value="manualCopy?.text ?? ''"
          :rows="isCode ? 10 : 3"
          :wrap="isCode ? 'off' : 'soft'"
          readonly
          spellcheck="false"
          aria-describedby="manual-copy-hint"
        />
        <p id="manual-copy-hint" class="manual-copy-hint">选中后按 Ctrl+C / ⌘C 复制，手机可长按内容复制。</p>
      </div>
      <div class="modal-foot">
        <button class="btn btn-ghost btn-sm" type="button" @click="closeManualCopy">关闭</button>
        <button class="btn btn-primary btn-sm" type="button" @click="selectContent">
          <MousePointer2 :size="15" aria-hidden="true" /> 全选{{ kindLabel }}
        </button>
      </div>
    </dialog>
  </Teleport>
</template>
