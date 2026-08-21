<script setup lang="ts">
/** 写文章 — 1:1 复刻静态 editor.html（.editor-page 全屏 + 分屏预览 + 发布弹窗 .tag-picker/.cover-upload/.pm-sec）。
 *  为与静态一致，本路由独立于 AdminLayout（无侧栏）。去 Element Plus / MarkdownEditor 组件 / Tailwind。 */
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MarkdownIt from 'markdown-it'
import {
  ArrowLeft,
  Calendar,
  Code2,
  Image as ImageIcon,
  Info,
  Monitor,
  Save,
  Send,
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
const targetStatus = ref<ArticleStatus>('published')
const viewMode = ref<'split' | 'preview' | 'edit'>('split')

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
    coverImage: coverUrl.value,
    category: pubCategory.value,
    tags: [...selectedTags.value],
    content: content.value,
    status,
    publishedAt: pubDate.value || new Date().toISOString().slice(0, 10),
  }
}

async function saveDraft() {
  const payload = buildPayload('draft')
  try {
    if (editingId.value == null) await createAdminArticle(payload)
    else await updateAdminArticle(editingId.value, payload)
    toast.success('草稿已保存')
  } catch {
    toast.error('保存失败')
  }
}

function openPublish() {
  if (!title.value.trim()) {
    toast.warning('请先填写文章标题')
    return
  }
  if (!pubCategory.value && categories.value.length) pubCategory.value = categories.value[0].name
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
  const payload = buildPayload(targetStatus.value)
  try {
    if (editingId.value == null) await createAdminArticle(payload)
    else await updateAdminArticle(editingId.value, payload)
    toast.success('文章已发布')
    showPublish.value = false
    router.push('/articles')
  } catch {
    toast.error('发布失败')
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
    targetStatus.value = a.status
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

onMounted(async () => {
  window.addEventListener('keydown', onKeydown)
  try {
    ;[categories.value, tags.value] = await Promise.all([
      listCategories({ page: 1, pageSize: 100 }).then((r) => r.list),
      listTags({ page: 1, pageSize: 100 }).then((r) => r.list),
    ])
  } catch {
    /* 静默 */
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
        <select v-model="targetStatus" class="select" style="width: 96px">
          <option value="draft">草稿</option>
          <option value="published">发布</option>
        </select>
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
          <textarea id="mdInput" v-model="content" spellcheck="false" placeholder="在这里输入 Markdown 内容…"></textarea>
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
</style>
