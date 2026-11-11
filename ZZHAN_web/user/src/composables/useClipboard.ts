import { readonly, ref } from 'vue'

export type CopyKind = 'code' | 'link'

const manualCopy = ref<{ text: string; kind: CopyKind } | null>(null)

/** 代码与链接共用复制入口，失败时保留内容供用户手动复制。 */
export function useClipboard() {
  async function copyText(text: string, kind: CopyKind): Promise<boolean> {
    try {
      if (typeof navigator.clipboard?.writeText !== 'function') {
        throw new Error('Clipboard unavailable')
      }
      await navigator.clipboard.writeText(text)
      return true
    } catch {
      manualCopy.value = { text, kind }
      return false
    }
  }

  function closeManualCopy(): void {
    manualCopy.value = null
  }

  return { manualCopy: readonly(manualCopy), copyText, closeManualCopy }
}
