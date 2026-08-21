import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'

const md: InstanceType<typeof MarkdownIt> = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: true,
  highlight(code, lang) {
    const langName = lang && hljs.getLanguage(lang) ? lang : ''
    const codeHtml = langName ? hljs.highlight(code, { language: langName }).value : md.utils.escapeHtml(code)
    return `<div class="code-block"><div class="code-head"><span class="lang">${langName || 'text'}</span><button class="code-copy" type="button">复制</button></div><pre><code class="language-${langName}">${codeHtml}</code></pre></div>`
  },
})

export function renderMarkdown(src: string): string {
  return md.render(src || '')
}

export function decorateCode(container: HTMLElement): void {
  container.querySelectorAll('pre code').forEach((el) => {
    const codeEl = el as HTMLElement
    const block = el.closest('.code-block') as HTMLElement | null
    if (!block || block.dataset.decorated) return
    block.dataset.decorated = '1'
    el.querySelectorAll('span.cl').forEach((n) => n.remove())
    const lines = el.innerHTML.split('\n')
    while (lines.length && /^\s*$/.test(lines[lines.length - 1])) lines.pop()
    el.innerHTML = lines.map((l) => `<span class="cl">${l || ' '}</span>`).join('\n')
    const btn = block.querySelector('.code-copy') as HTMLButtonElement | null
    btn?.addEventListener('click', () => {
      const done = () => {
        btn.textContent = '已复制'
        setTimeout(() => {
          btn.textContent = '复制'
        }, 1600)
      }
      const clip = navigator.clipboard
      if (clip && typeof clip.writeText === 'function') {
        Promise.resolve().then(() => clip.writeText(codeEl.innerText)).then(done).catch(done)
      } else {
        done()
      }
    })
  })
}
