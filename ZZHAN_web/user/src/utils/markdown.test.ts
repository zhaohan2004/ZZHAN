import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { renderMarkdown, buildTOC, decorateCode } from './markdown'
import { useClipboard } from '@/composables/useClipboard'

describe('renderMarkdown', () => {
  it('renders headings, lists, tables, code', () => {
    const html = renderMarkdown('# 标题\n\n- a\n- b\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\n```ts\nconst x: number = 1\n```')
    expect(html).toContain('<h2 id="sec-0-标题">标题</h2>')
    expect(html).toContain('<ul>')
    expect(html).toContain('<table>')
    expect(html).toContain('class="code-block"')
  })

  it('escapes html and opens links in new tab', () => {
    const html = renderMarkdown('<script>alert(1)</script> [link](https://x.com)')
    expect(html).not.toContain('<script>')
    expect(html).toContain('target="_blank" rel="noopener"')
  })

  it('shifts heading levels and assigns stable ids', () => {
    const html = renderMarkdown('# 标题\n\n## 小节\n\n### 第三级')
    expect(html).toContain('<h2 id="sec-0-标题">标题</h2>')
    expect(html).toContain('<h3 id="sec-1-小节">小节</h3>')
    expect(html).toContain('<h4 id="sec-2-第三级">第三级</h4>')
  })

  it('escapes raw html inside inline code', () => {
    const html = renderMarkdown('`<div>`')
    expect(html).not.toContain('<div>')
  })
})

describe('buildTOC', () => {
  it('collects h2/h3 and assigns ids when missing', () => {
    const container = document.createElement('div')
    container.innerHTML = '<h2>第一章</h2><h3>第一节</h3><h2 id="custom">第二章</h2>'
    const toc = buildTOC(container)
    expect(toc).toEqual([
      { level: 2, id: 'sec-0-第一章', text: '第一章' },
      { level: 3, id: 'sec-1-第一节', text: '第一节' },
      { level: 2, id: 'custom', text: '第二章' },
    ])
  })

  it('uses renderMarkdown ids for toc anchors', () => {
    const html = renderMarkdown('# 标题\n\n## 小节')
    const container = document.createElement('div')
    container.innerHTML = html
    const toc = buildTOC(container)
    expect(toc).toEqual([
      { level: 2, id: 'sec-0-标题', text: '标题' },
      { level: 3, id: 'sec-1-小节', text: '小节' },
    ])
  })
})

describe('decorateCode', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    useClipboard().closeManualCopy()
  })

  afterEach(() => {
    useClipboard().closeManualCopy()
    vi.unstubAllGlobals()
    vi.useRealTimers()
  })

  it('wraps each line in span.cl and guards re-decoration', () => {
    const html = renderMarkdown('```ts\nconst a = 1\nconst b = 2\n```')
    const container = document.createElement('div')
    container.innerHTML = html
    decorateCode(container)
    const block = container.querySelector('.code-block') as HTMLElement
    expect(block.dataset.decorated).toBe('1')
    const code = container.querySelector('.code-block code') as HTMLElement
    expect(code.querySelectorAll('span.cl').length).toBe(2)
    const before = container.innerHTML
    decorateCode(container)
    expect(container.innerHTML).toBe(before)
  })

  it('copies the original code once and shows success only after the write resolves', async () => {
    let finishCopy!: () => void
    const pendingWrite = new Promise<void>((resolve) => { finishCopy = resolve })
    const writeText = vi.fn(() => pendingWrite)
    vi.stubGlobal('navigator', { clipboard: { writeText } })
    const source = 'function greet() {\n  const name = "<visitor>"\n\n  return name\n}\n\n'
    const html = renderMarkdown('```ts\n' + source + '```')
    const container = document.createElement('div')
    container.innerHTML = html
    decorateCode(container)
    decorateCode(container)
    const btn = container.querySelector('.code-copy') as HTMLButtonElement

    btn.click()
    expect(writeText).toHaveBeenCalledExactlyOnceWith(source)
    expect(btn.textContent).toBe('复制')

    finishCopy()
    await vi.advanceTimersByTimeAsync(0)
    expect(btn.textContent).toBe('已复制')
    expect(useClipboard().manualCopy.value).toBeNull()

    vi.advanceTimersByTime(1600)
    expect(btn.textContent).toBe('复制')
  })

  it('offers the original code for manual copying when clipboard writing is rejected', async () => {
    const writeText = vi.fn().mockRejectedValue(new DOMException('Permission denied', 'NotAllowedError'))
    vi.stubGlobal('navigator', { clipboard: { writeText } })
    const source = 'if (ready) {\n\trun()\n\n}\n'
    const container = document.createElement('div')
    container.innerHTML = renderMarkdown('```ts\n' + source + '```')
    decorateCode(container)
    const btn = container.querySelector('.code-copy') as HTMLButtonElement

    btn.click()
    await vi.advanceTimersByTimeAsync(0)

    expect(writeText).toHaveBeenCalledExactlyOnceWith(source)
    expect(btn.textContent).toBe('复制')
    expect(useClipboard().manualCopy.value).toEqual({ text: source, kind: 'code' })
  })

  it('offers manual copying without claiming success when the Clipboard API is missing', async () => {
    vi.stubGlobal('navigator', {})
    const container = document.createElement('div')
    container.innerHTML = renderMarkdown('```ts\nconst a = 1\n```')
    decorateCode(container)
    const btn = container.querySelector('.code-copy') as HTMLButtonElement

    btn.click()
    await vi.advanceTimersByTimeAsync(0)

    expect(btn.textContent).toBe('复制')
    expect(useClipboard().manualCopy.value).toEqual({ text: 'const a = 1\n', kind: 'code' })
  })
})
