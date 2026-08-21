import { describe, it, expect } from 'vitest'
import { renderMarkdown, buildTOC, decorateCode } from './markdown'

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

  it('copy button swaps text to 已复制 on click', async () => {
    const html = renderMarkdown('```ts\nconst a = 1\n```')
    const container = document.createElement('div')
    container.innerHTML = html
    decorateCode(container)
    const btn = container.querySelector('.code-copy') as HTMLButtonElement
    btn.click()
    await new Promise((r) => setTimeout(r, 0))
    expect(btn.textContent).toBe('已复制')
  })
})
