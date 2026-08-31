<script setup lang="ts">
/** 写文章 — 1:1 复刻静态 editor.html（.editor-page 全屏 + 分屏预览 + 发布弹窗 .tag-picker/.cover-upload/.pm-sec）。
 *  为与静态一致，本路由独立于 AdminLayout（无侧栏）。去 Element Plus / MarkdownEditor 组件 / Tailwind。 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import {
  ArrowLeft,
  Bold,
  Calendar,
  CheckSquare,
  Code,
  Code2,
  Heading,
  Image as ImageIcon,
  Info,
  Italic,
  Link,
  List,
  ListOrdered,
  Minus,
  Monitor,
  Quote,
  Save,
  Send,
  Strikethrough,
  Table,
  Terminal,
  Upload,
  X,
} from 'lucide-vue-next'
import { createAdminArticle, getAdminArticle, listCategories, listTags, updateAdminArticle } from '@/api/admin'
import type { AdminArticlePayload, ArticleStatus, CategoryAdmin, TagAdmin } from '@/types/models'
import { confirm, toast } from '@/composables/useToast'

const md = new MarkdownIt({ html: true, linkify: true, typographer: true })

const route = useRoute()
const router = useRouter()

const title = ref('')
const content = ref('')
const viewMode = ref<'split' | 'preview' | 'edit'>('split')
const textareaRef = ref<HTMLTextAreaElement | null>(null)

const categories = ref<CategoryAdmin[]>([])
const tags = ref<TagAdmin[]>([])
const editingId = ref<number | null>(null)

const previewHtml = computed(() => md.render(content.value))
const charCount = computed(() => content.value.length)
const workspaceClass = computed(() => (viewMode.value === 'split' ? '' : 'view-' + viewMode.value))

/* 发布弹窗 */
const showPublish = ref(false)
const pubCategory = ref('')
const selectedTags = ref<string[]>([])
const summary = ref('')
const pubDate = ref('')
const coverUrl = ref('')

function back() {
  router.push('/articles')
}

function toggleTag(name: string) {
  const i = selectedTags.value.indexOf(name)
  if (i >= 0) selectedTags.value.splice(i, 1)
  else selectedTags.value.push(name)
}

function onCoverPick(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files && input.files[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    coverUrl.value = String(reader.result)
  }
  reader.readAsDataURL(file)
}
function removeCover() {
  coverUrl.value = ''
}

function buildPayload(status: ArticleStatus): AdminArticlePayload {
  return {
    title: title.value.trim() || '未命名文章',
    summary: summary.value.trim(),
    cover_image: coverUrl.value,
    category: pubCategory.value,
    tags: [...selectedTags.value],
    content: content.value,
    status,
    published_at: pubDate.value || new Date().toISOString().slice(0, 10),
  }
}

async function saveDraft() {
  if (!title.value.trim()) {
    toast.warning('请先填写文章标题')
    return
  }
  if (!pubCategory.value) {
    toast.warning('请先选择分类')
    return
  }
  const payload = buildPayload('draft')
  try {
    if (editingId.value == null) await createAdminArticle(payload)
    else await updateAdminArticle(editingId.value, payload)
    toast.success('草稿已保存')
  } catch (e: any) {
    toast.error(e?.message || '保存失败')
  }
}

function openPublish() {
  if (!title.value.trim()) {
    toast.warning('请先填写文章标题')
    return
  }
  if (!pubCategory.value) pubCategory.value = ''
  if (!pubDate.value) pubDate.value = new Date().toISOString().slice(0, 10)
  showPublish.value = true
}

async function confirmPublish() {
  if (!pubCategory.value) {
    toast.warning('请选择分类')
    return
  }
  const ok = await confirm('确定发布这篇文章吗？', '发布文章')
  if (!ok) return
  const payload = buildPayload('published')
  try {
    if (editingId.value == null) await createAdminArticle(payload)
    else await updateAdminArticle(editingId.value, payload)
    toast.success('文章已发布')
    showPublish.value = false
    router.push('/articles')
  } catch (e: any) {
    toast.error(e?.message || '发布失败')
  }
}

async function loadForEdit() {
  const id = (route.params.id as string) || (route.query.id as string)
  if (!id) return
  try {
    const a = await getAdminArticle(id)
    editingId.value = a.id
    title.value = a.title
    content.value = a.content || ''
    pubCategory.value = a.category
    selectedTags.value = [...(a.tags || [])]
    summary.value = a.summary || ''
    pubDate.value = a.date
  } catch {
    toast.error('加载文章失败')
  }
}

function onKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
    e.preventDefault()
    saveDraft()
  }
}

// Markdown 格式工具栏功能
function insertText(before: string, after: string, defaultText: string = '') {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  const text = selected || defaultText
  const replacement = before + text + (after || '')
  content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
  const newPos = start + before.length + text.length
  setTimeout(() => {
    textarea.focus()
    textarea.setSelectionRange(newPos, newPos)
  }, 0)
}

function insertLine(text: string) {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const before = content.value.substring(0, start)
  const after = content.value.substring(start)
  const prefix = before.length > 0 && !before.endsWith('\n') ? '\n' : ''
  content.value = before + prefix + text + '\n' + after
  const newPos = start + prefix.length + text.length + 1
  setTimeout(() => {
    textarea.focus()
    textarea.setSelectionRange(newPos, newPos)
  }, 0)
}

function wrapSelection(wrapper: string, defaultText: string = '') {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  const text = selected || defaultText
  const beforeWrapper = content.value.substring(Math.max(0, start - wrapper.length), start)
  const afterWrapper = content.value.substring(end, end + wrapper.length)
  if (beforeWrapper === wrapper && afterWrapper === wrapper) {
    content.value = content.value.substring(0, start - wrapper.length) + text + content.value.substring(end + wrapper.length)
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start - wrapper.length, end - wrapper.length)
    }, 0)
  } else {
    const replacement = wrapper + text + wrapper
    content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + wrapper.length, end + wrapper.length)
    }, 0)
  }
}

function formatHeading() {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const lineStart = content.value.lastIndexOf('\n', start - 1) + 1
  const lineEnd = content.value.indexOf('\n', start)
  const line = content.value.substring(lineStart, lineEnd === -1 ? content.value.length : lineEnd)
  const match = line.match(/^(#{1,6})\s/)
  if (match) {
    const level = Math.min(match[1].length + 1, 6)
    const newLine = '#'.repeat(level) + ' ' + line.replace(/^#{1,6}\s/, '')
    content.value = content.value.substring(0, lineStart) + newLine + content.value.substring(lineEnd === -1 ? content.value.length : lineEnd)
  } else {
    const newLine = '## ' + line
    content.value = content.value.substring(0, lineStart) + newLine + content.value.substring(lineEnd === -1 ? content.value.length : lineEnd)
  }
  setTimeout(() => textarea.focus(), 0)
}

function formatLink() {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  if (selected) {
    const replacement = '[' + selected + '](url)'
    content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + selected.length + 3, start + selected.length + 6)
    }, 0)
  } else {
    insertText('[', '](url)', '链接文本')
  }
}

function formatImage() {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  if (selected) {
    const replacement = '![' + selected + '](image-url)'
    content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + selected.length + 4, start + selected.length + 13)
    }, 0)
  } else {
    insertText('![', '](image-url)', '图片描述')
  }
}

function formatCodeblock() {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selected = content.value.substring(start, end)
  if (selected) {
    const replacement = '```\n' + selected + '\n```'
    content.value = content.value.substring(0, start) + replacement + content.value.substring(end)
    setTimeout(() => {
      textarea.focus()
      textarea.setSelectionRange(start + 4, start + 4 + selected.length)
    }, 0)
  } else {
    insertLine('```\n代码\n```')
  }
}

function formatTable() {
  const table = '| 列1 | 列2 | 列3 |\n| --- | --- | --- |\n| 内容 | 内容 | 内容 |'
  insertLine(table)
}

function insertHeading(level: number) {
  const textarea = textareaRef.value
  if (!textarea) return
  const start = textarea.selectionStart
  const lineStart = content.value.lastIndexOf('\n', start - 1) + 1
  const lineEnd = content.value.indexOf('\n', start)
  const line = content.value.substring(lineStart, lineEnd === -1 ? content.value.length : lineEnd)
  const cleanLine = line.replace(/^#{1,6}\s/, '')
  const newLine = '#'.repeat(level) + ' ' + cleanLine
  content.value = content.value.substring(0, lineStart) + newLine + content.value.substring(lineEnd === -1 ? content.value.length : lineEnd)
  setTimeout(() => textarea.focus(), 0)
}

function handleToolbarAction(action: string) {
  const headingActions: Record<string, () => void> = {
    h1: () => insertHeading(1),
    h2: () => insertHeading(2),
    h3: () => insertHeading(3),
    h4: () => insertHeading(4),
    h5: () => insertHeading(5),
    h6: () => insertHeading(6),
  }
  if (headingActions[action]) {
    headingActions[action]()
    return
  }
  const actions: Record<string, () => void> = {
    bold: () => wrapSelection('**', '粗体文本'),
    italic: () => wrapSelection('*', '斜体文本'),
    strikethrough: () => wrapSelection('~~', '删除线文本'),
    link: formatLink,
    image: formatImage,
    code: () => wrapSelection('`', 'code'),
    codeblock: formatCodeblock,
    ul: () => insertLine('- 列表项'),
    ol: () => insertLine('1. 列表项'),
    quote: () => insertLine('> 引用文本'),
    hr: () => insertLine('---'),
    table: formatTable,
    tasklist: () => insertLine('- [ ] 任务项'),
  }
  if (actions[action]) {
    actions[action]()
  }
}

onMounted(async () => {
  window.addEventListener('keydown', onKeydown)
  try {
    ;[categories.value, tags.value] = await Promise.all([
      listCategories({ page: 1, pageSize: 100 }).then((r) => r.list),
      listTags({ page: 1, pageSize: 100 }).then((r) => r.list),
    ])
  } catch (e: any) {
    console.error('加载分类/标签失败:', e?.message)
  }
  await loadForEdit()
})
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="editor-page" style="overflow: hidden">
    <div class="editor-toolbar">
      <button class="icon-btn" title="返回" @click="back"><ArrowLeft :size="18" /></button>
      <label class="ed-title-label" for="mdInput">文章标题</label>
      <input id="edTitle" v-model="title" class="input" type="text" placeholder="输入文章标题…" autocomplete="off" />
      <div class="tb-actions">
        <button class="btn btn-ghost" @click="saveDraft"><Save :size="15" /> 保存草稿</button>
        <select v-model="viewMode" class="select" title="视图模式" style="width: 104px">
          <option value="split">分屏</option>
          <option value="preview">仅预览</option>
          <option value="edit">仅编辑</option>
        </select>
        <button class="btn btn-primary" @click="openPublish"><Send :size="15" /> 发布</button>
      </div>
    </div>

    <div class="editor-main">
      <div class="editor-workspace" :class="workspaceClass">
        <div class="editor-pane">
          <div class="editor-pane-head">
            <span class="eph-title"><Code2 :size="14" style="color: var(--accent)" /> Markdown 编辑</span>
            <span class="ep-hint">{{ charCount }} 字符</span>
          </div>
          <div class="md-toolbar">
            <div class="md-toolbar-group">
              <div class="md-toolbar-dropdown">
                <button class="md-toolbar-btn">
                  <Heading :size="16" />
                  <span class="md-toolbar-label">标题</span>
                </button>
                <div class="md-toolbar-dropdown-menu">
                  <button class="md-dropdown-item" @click="handleToolbarAction('h1')"><strong>H1</strong> 一级标题</button>
                  <button class="md-dropdown-item" @click="handleToolbarAction('h2')"><strong>H2</strong> 二级标题</button>
                  <button class="md-dropdown-item" @click="handleToolbarAction('h3')"><strong>H3</strong> 三级标题</button>
                  <button class="md-dropdown-item" @click="handleToolbarAction('h4')"><strong>H4</strong> 四级标题</button>
                  <button class="md-dropdown-item" @click="handleToolbarAction('h5')"><strong>H5</strong> 五级标题</button>
                  <button class="md-dropdown-item" @click="handleToolbarAction('h6')"><strong>H6</strong> 六级标题</button>
                </div>
              </div>
              <button class="md-toolbar-btn" @click="handleToolbarAction('bold')">
                <Bold :size="16" />
                <span class="md-toolbar-label">粗体</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('italic')">
                <Italic :size="16" />
                <span class="md-toolbar-label">斜体</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('strikethrough')">
                <Strikethrough :size="16" />
                <span class="md-toolbar-label">删除线</span>
              </button>
            </div>
            <div class="md-toolbar-divider"></div>
            <div class="md-toolbar-group">
              <button class="md-toolbar-btn" @click="handleToolbarAction('link')">
                <Link :size="16" />
                <span class="md-toolbar-label">链接</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('image')">
                <ImageIcon :size="16" />
                <span class="md-toolbar-label">图片</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('code')">
                <Code :size="16" />
                <span class="md-toolbar-label">代码</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('codeblock')">
                <Terminal :size="16" />
                <span class="md-toolbar-label">代码块</span>
              </button>
            </div>
            <div class="md-toolbar-divider"></div>
            <div class="md-toolbar-group">
              <button class="md-toolbar-btn" @click="handleToolbarAction('ul')">
                <List :size="16" />
                <span class="md-toolbar-label">无序列表</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('ol')">
                <ListOrdered :size="16" />
                <span class="md-toolbar-label">有序列表</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('quote')">
                <Quote :size="16" />
                <span class="md-toolbar-label">引用</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('hr')">
                <Minus :size="16" />
                <span class="md-toolbar-label">分割线</span>
              </button>
            </div>
            <div class="md-toolbar-divider"></div>
            <div class="md-toolbar-group">
              <button class="md-toolbar-btn" @click="handleToolbarAction('table')">
                <Table :size="16" />
                <span class="md-toolbar-label">表格</span>
              </button>
              <button class="md-toolbar-btn" @click="handleToolbarAction('tasklist')">
                <CheckSquare :size="16" />
                <span class="md-toolbar-label">任务列表</span>
              </button>
            </div>
          </div>
          <textarea id="mdInput" ref="textareaRef" v-model="content" spellcheck="false" placeholder="在这里输入 Markdown 内容…"></textarea>
        </div>
        <div class="editor-pane">
          <div class="editor-pane-head">
            <span class="eph-title"><Monitor :size="14" style="color: var(--accent)" /> 实时预览</span>
            <span class="ep-hint">Ctrl + S 保存草稿</span>
          </div>
          <div id="mdPreview" class="prose" v-html="previewHtml"></div>
        </div>
      </div>
    </div>

    <!-- 发布弹窗 -->
    <div class="modal-overlay" :class="{ open: showPublish }" @click.self="showPublish = false">
      <div class="modal login-modal">
        <div class="modal-head">
          <h3>发布文章</h3>
          <button class="modal-close" @click="showPublish = false"><X :size="18" /></button>
        </div>
        <div class="modal-body">
          <div class="pm-sec">
            <div class="pm-sec-title"><Info :size="14" /> 文章信息</div>
            <div class="form-group" style="text-align: left">
              <label class="form-label">文章标题</label>
              <div id="pubTitle" class="pub-title">{{ title || '未命名文章' }}</div>
            </div>
            <div class="form-group" style="text-align: left">
              <label class="form-label" for="edCat">文章分类 <span class="req">*</span></label>
              <select id="edCat" v-model="pubCategory" class="select">
                <option value="" disabled>请选择分类</option>
                <option v-for="c in categories" :key="c.id" :value="c.name">{{ c.name }}</option>
              </select>
            </div>
            <div class="form-group" style="text-align: left">
              <label class="form-label">标签</label>
              <div class="tag-picker">
                <div class="tp-selected" :class="{ 'has-chip': selectedTags.length }">
                  <template v-if="selectedTags.length">
                    <span v-for="t in selectedTags" :key="t" class="tp-chip active" @click="toggleTag(t)">
                      {{ t }} <span class="tp-x">×</span>
                    </span>
                  </template>
                  <span v-else class="tp-empty">暂未选择</span>
                </div>
                <div class="tp-source-title">点击下方标签添加</div>
                <div class="tp-source">
                  <span
                    v-for="t in tags"
                    :key="t.id"
                    class="tp-chip"
                    :class="{ active: selectedTags.includes(t.name) }"
                    @click="toggleTag(t.name)"
                  >
                    {{ t.name }}
                  </span>
                </div>
              </div>
            </div>
            <div class="form-group" style="text-align: left">
              <label class="form-label" for="edSummary">文章摘要</label>
              <textarea id="edSummary" v-model="summary" class="textarea" rows="3" placeholder="用于列表页与 SEO 的文章摘要…"></textarea>
            </div>
          </div>

          <div class="pm-sec">
            <div class="pm-sec-title"><Calendar :size="14" /> 发布设置</div>
            <div class="form-group" style="text-align: left">
              <label class="form-label" for="edDate">发布时间</label>
              <input id="edDate" v-model="pubDate" class="input" type="date" />
            </div>
            <div class="form-group" style="text-align: left">
              <label class="form-label">封面</label>
              <div class="cover-upload">
                <div v-if="coverUrl" class="cu-preview">
                  <img :src="coverUrl" alt="封面预览" />
                  <button class="cu-remove" type="button" title="移除封面" @click="removeCover"><X :size="14" /></button>
                </div>
                <div v-else class="cu-placeholder">
                  <ImageIcon :size="22" />
                  <span>未上传，将自动生成渐变封面</span>
                </div>
                <label class="btn btn-ghost" style="width: 100%; margin-top: 10px; justify-content: center">
                  <Upload :size="15" /> 上传封面图
                  <input type="file" accept="image/*" hidden @change="onCoverPick" />
                </label>
              </div>
            </div>
          </div>

          <button class="btn btn-primary" style="width: 100%; padding: 12px" @click="confirmPublish">
            <Send :size="16" /> 确认发布
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.editor-page :deep(.prose pre) {
  background: rgba(15, 23, 42, 0.55);
  border: 1px solid var(--border-strong);
  border-radius: 12px;
  padding: 14px 16px;
  overflow-x: auto;
  margin: 1.2em 0;
}
.editor-page :deep(.prose pre code) {
  background: transparent;
  padding: 0;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: #e2e8f0;
}
.editor-page :deep(.prose code.inline) {
  background: var(--card-2);
  padding: 2px 6px;
  border-radius: 6px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 0.88em;
}

/* Markdown 格式工具栏 */
.md-toolbar {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: var(--card-2);
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
  flex-shrink: 0;
}
.md-toolbar-group {
  display: flex;
  align-items: center;
  gap: 2px;
}
.md-toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border-strong);
  margin: 0 4px;
}
.md-toolbar-btn {
  display: inline-flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  min-width: 42px;
  height: 48px;
  padding: 6px 8px;
  border-radius: 8px;
  background: transparent;
  color: var(--text-2);
  transition: all 0.2s ease;
  position: relative;
}
.md-toolbar-label {
  font-size: 10px;
  line-height: 1;
  white-space: nowrap;
}
.md-toolbar-btn:hover {
  background: var(--grad-soft);
  color: var(--accent);
  transform: translateY(-1px);
}
.md-toolbar-btn:active {
  transform: translateY(0);
  background: rgba(59, 130, 246, 0.2);
}

/* 标题下拉菜单 */
.md-toolbar-dropdown {
  position: relative;
}
.md-toolbar-dropdown-menu {
  display: none;
  position: absolute;
  top: 100%;
  left: 0;
  background: var(--card);
  border: 1px solid var(--border-strong);
  border-radius: 10px;
  padding: 6px;
  box-shadow: var(--shadow-lg);
  z-index: 100;
  min-width: 160px;
}
.md-toolbar-dropdown:hover .md-toolbar-dropdown-menu {
  display: block;
}
.md-dropdown-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  border: none;
  background: transparent;
  color: var(--text);
  font-size: 13px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  text-align: left;
}
.md-dropdown-item:hover {
  background: var(--grad-soft);
  color: var(--accent);
}
.md-dropdown-item strong {
  color: var(--accent);
  font-weight: 700;
  min-width: 24px;
}

@media (max-width: 768px) {
  .md-toolbar {
    padding: 6px 8px;
    gap: 2px;
  }
  .md-toolbar-btn {
    min-width: 36px;
    height: 42px;
    padding: 4px 6px;
  }
  .md-toolbar-label {
    font-size: 9px;
  }
  .md-toolbar-divider {
    height: 16px;
    margin: 0 2px;
  }
}
</style>
