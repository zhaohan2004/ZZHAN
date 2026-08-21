import MarkdownIt, { type RendererRule } from 'markdown-it'
import hljs from 'highlight.js'

/** Mirrors the static site's slugifier: keep CJK + ASCII word chars, collapse the rest to '-'. */
function slugify(s: string): string {
  return String(s)
    .trim()
    .toLowerCase()
    .replace(/[^\w一-龥]+/g, '-')
    .replace(/^-+|-+$/g, '') || 'h'
}

const md: InstanceType<typeof MarkdownIt> = new MarkdownIt({
  html: false,
  linkify: true,
  breaks: true,
  highlight(code, lang) {
    const langName = lang && hljs.getLanguage(lang) ? lang : ''
    const codeHtml = langName
      ? hljs.highlight(code, { language: langName }).value
      : md.utils.escapeHtml(code)
    return `<div class="code-block"><div class="code-head"><span class="lang">${langName || 'text'}</span><button class="code-copy" type="button">复制</button></div><pre><code class="language-${langName}">${codeHtml}</code></pre></div>`
  },
})

const defaultLink: RendererRule =
  md.renderer.rules.link_open ||
  ((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options))

md.renderer.rules.link_open = (tokens, idx, options, env, self) => {
  tokens[idx].attrSet('target', '_blank')
  tokens[idx].attrSet('rel', 'noopener')
  return defaultLink(tokens, idx, options, env, self)
}

// Shift every heading down one level (h1→h2 … h6→h7) so article prose never
// renders an <h1> (the page <h1> is the article title), and give every heading
// a stable id (`sec-{index}-{slug}`) so TOC anchors resolve.
md.core.ruler.after('inline', 'heading_shift_and_ids', (state) => {
  const tokens = state.tokens
  let headingIndex = 0
  for (let i = 0; i < tokens.length; i++) {
    const t = tokens[i]
    if (t.type !== 'heading_open') continue
    const inline = tokens[i + 1]
    let text = ''
    if (inline && inline.children) {
      for (const child of inline.children) {
        if (child.type === 'text' || child.type === 'code_inline') text += child.content
      }
    }
    if (!text && inline) text = inline.content
    const tag = 'h' + (t.tag === 'h1' ? 2 : Number(t.tag[1]) + 1)
    t.tag = tag
    t.attrSet('id', 'sec-' + headingIndex + '-' + slugify(text))
    if (tokens[i + 2] && tokens[i + 2].type === 'heading_close') {
      tokens[i + 2].tag = tag
    }
    headingIndex++
  }
})

export function renderMarkdown(src: string): string {
  return md.render(src)
}

export function buildTOC(container: HTMLElement): { level: 2 | 3; id: string; text: string }[] {
  const out: { level: 2 | 3; id: string; text: string }[] = []
  container.querySelectorAll('h2,h3').forEach((h) => {
    if (!h.id) h.id = 'sec-' + out.length + '-' + slugify(h.textContent!.trim())
    out.push({ level: h.tagName === 'H2' ? 2 : 3, id: h.id, text: h.textContent!.trim() })
  })
  return out
}

export function decorateCode(container: HTMLElement): void {
  container.querySelectorAll('pre code').forEach((el) => {
    const codeEl = el as HTMLElement
    const block = el.closest('.code-block') as HTMLElement | null
    if (!block || block.dataset.decorated) return
    block.dataset.decorated = '1'
    el.querySelectorAll('span.cl').forEach((n) => n.remove())
    const lines = el.innerHTML.split('\n')
    // highlight.js appends a trailing newline; drop trailing empty lines like the static site does
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
        Promise.resolve()
          .then(() => clip.writeText(codeEl.innerText))
          .then(done)
          .catch(done)
      } else {
        done()
      }
    })
  })
}
