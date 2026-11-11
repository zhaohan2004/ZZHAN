import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useClipboard } from './useClipboard'

describe('useClipboard', () => {
  beforeEach(() => {
    useClipboard().closeManualCopy()
  })

  afterEach(() => {
    useClipboard().closeManualCopy()
    vi.unstubAllGlobals()
  })

  it.each(['code', 'link'] as const)('writes %s text before reporting success', async (kind) => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    vi.stubGlobal('navigator', { clipboard: { writeText } })
    const clipboard = useClipboard()
    const text = kind === 'link' ? 'http://example.test/article/go' : 'if (ok) {\n  run()\n}\n'

    await expect(clipboard.copyText(text, kind)).resolves.toBe(true)

    expect(writeText).toHaveBeenCalledExactlyOnceWith(text)
    expect(clipboard.manualCopy.value).toBeNull()
  })

  it.each([undefined, {}])('offers manual copying when the Clipboard API is unavailable: %s', async (api) => {
    vi.stubGlobal('navigator', { clipboard: api })
    const clipboard = useClipboard()
    const text = 'http://example.test/article/go?tag_id=2'

    await expect(clipboard.copyText(text, 'link')).resolves.toBe(false)

    expect(clipboard.manualCopy.value).toEqual({ text, kind: 'link' })
  })

  it('preserves the full code for manual copying when permission is denied', async () => {
    const writeText = vi.fn().mockRejectedValue(new DOMException('Permission denied', 'NotAllowedError'))
    vi.stubGlobal('navigator', { clipboard: { writeText } })
    const clipboard = useClipboard()
    const text = 'func main() {\n\tprintln("你好")\n\n}\n'

    await expect(clipboard.copyText(text, 'code')).resolves.toBe(false)

    expect(writeText).toHaveBeenCalledExactlyOnceWith(text)
    expect(clipboard.manualCopy.value).toEqual({ text, kind: 'code' })
  })

  it('offers manual copying when the browser throws synchronously', async () => {
    const writeText = vi.fn(() => { throw new Error('Clipboard blocked') })
    vi.stubGlobal('navigator', { clipboard: { writeText } })
    const clipboard = useClipboard()

    await expect(clipboard.copyText('http://example.test/article/go', 'link')).resolves.toBe(false)

    expect(clipboard.manualCopy.value).toEqual({ text: 'http://example.test/article/go', kind: 'link' })
  })

  it('shares fallback content with the dialog host and clears it on close', async () => {
    vi.stubGlobal('navigator', {})
    const trigger = useClipboard()
    const dialog = useClipboard()

    await trigger.copyText('const answer = 42\n', 'code')
    expect(dialog.manualCopy.value).toEqual({ text: 'const answer = 42\n', kind: 'code' })

    dialog.closeManualCopy()
    expect(trigger.manualCopy.value).toBeNull()
  })
})
