<script setup lang="ts">
/** 系统设置 — 1:1 复刻静态 settings.html（.settings-grid 两列：基本信息 / 关于页&首页展示）。去 Element Plus。 */
import { onMounted, reactive, ref } from 'vue'
import { Info, Quote, RotateCcw, Save, Sparkles, Trash2, Upload } from 'lucide-vue-next'
import { getSettings, saveSettings } from '@/api/admin'
import type { SettingsKV } from '@/types/models'
import { SETTINGS } from '@/api/mock/data'
import { confirm, toast } from '@/composables/useToast'

const loading = ref(false)
const saving = ref(false)
const form = reactive<SettingsKV>({ ...(SETTINGS as SettingsKV) })

async function load() {
  loading.value = true
  try {
    const s = await getSettings()
    Object.assign(form, s)
  } catch {
    toast.error('加载设置失败')
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await saveSettings({ ...form })
    toast.success('设置已保存')
  } catch {
    toast.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function reset() {
  const ok = await confirm('确定恢复所有设置为默认值吗？', '恢复默认')
  if (!ok) return
  Object.assign(form, { ...(SETTINGS as SettingsKV) })
  toast.info('已恢复默认，记得点击保存')
}

/** 头像上传 — 读为 data URL 写入 form.avatar，同时即时预览。 */
function onAvatarChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    toast.error('头像文件不能超过 2MB')
    input.value = ''
    return
  }
  const reader = new FileReader()
  reader.onload = () => {
    form.avatar = String(reader.result || '')
  }
  reader.readAsDataURL(file)
}

function clearAvatar() {
  form.avatar = ''
  const el = document.getElementById('setAvatarInput') as HTMLInputElement | null
  if (el) el.value = ''
}

onMounted(load)
</script>

<template>
  <div>
    <div class="settings-grid">
      <!-- 基本信息 -->
      <div class="settings-sec reveal in">
        <h3>
          <span class="st-ico" style="width: 30px; height: 30px"><Info :size="15" /></span>基本信息
        </h3>
        <p class="sec-desc">博客的展示名称、简介与作者信息。</p>
        <div style="display: grid; gap: 4px; grid-template-columns: 1fr 1fr">
          <div class="form-group">
            <label class="form-label">博客名称</label>
            <input v-model="form.blogName" class="input" type="text" />
          </div>
          <div class="form-group">
            <label class="form-label">博客简介</label>
            <input v-model="form.blogDesc" class="input" type="text" />
          </div>
          <div class="form-group">
            <label class="form-label">Logo 文字</label>
            <div style="display: flex; gap: 10px; align-items: center">
              <input v-model="form.logoText" class="input" type="text" maxlength="3" style="width: 110px" />
              <span class="brand-logo" style="width: 38px; height: 38px; font-size: 15px">{{ form.logoText }}</span>
            </div>
          </div>
          <div class="form-group" style="grid-column: 1 / -1">
            <label class="form-label">作者头像</label>
            <div style="display: flex; gap: 12px; align-items: center; flex-wrap: wrap">
              <span class="brand-logo" style="width: 48px; height: 48px; font-size: 18px; flex: none; overflow: hidden; padding: 0">
                <img v-if="form.avatar" :src="form.avatar" alt="头像预览" style="width: 100%; height: 100%; object-fit: cover" />
                <span v-else>{{ form.authorName?.slice(0, 1) || 'U' }}</span>
              </span>
              <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap">
                <label class="btn btn-ghost btn-sm" style="cursor: pointer; margin: 0">
                  <Upload :size="14" /> 上传头像
                  <input id="setAvatarInput" type="file" accept="image/*" style="display: none" @change="onAvatarChange" />
                </label>
                <button v-if="form.avatar" type="button" class="btn btn-ghost btn-sm" @click="clearAvatar">
                  <Trash2 :size="14" /> 清除
                </button>
                <span class="form-hint" style="margin: 0">支持 JPG/PNG/WebP，最大 2MB</span>
              </div>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">作者昵称</label>
            <input v-model="form.authorName" class="input" type="text" />
          </div>
          <div class="form-group">
            <label class="form-label">作者头衔</label>
            <input v-model="form.authorRole" class="input" type="text" />
          </div>
          <div class="form-group" style="grid-column: 1 / -1">
            <label class="form-label">作者介绍</label>
            <textarea v-model="form.authorIntro" class="textarea" rows="2"></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">GitHub 地址</label>
            <input v-model="form.github" class="input" type="text" />
          </div>
          <div class="form-group">
            <label class="form-label">联系邮箱</label>
            <input v-model="form.email" class="input" type="email" />
          </div>
        </div>
      </div>

      <!-- 关于页与首页 hero 展示 -->
      <div class="settings-sec reveal in">
        <h3>
          <span class="st-ico" style="width: 30px; height: 30px"><Quote :size="15" /></span>关于页 & 首页展示
        </h3>
        <p class="sec-desc">这些内容同步展示在前台「关于我」页面与首页 Hero 区。</p>
        <div class="form-group">
          <label class="form-label">站点标语 <span class="muted" style="font-size: 11px; margin-left: 4px">（首页 hero 副标题）</span></label>
          <input v-model="form.tagline" class="input" type="text" placeholder="如：记录代码，分享技术，持续成长" />
        </div>
        <div class="form-group">
          <label class="form-label">座右铭 <span class="muted" style="font-size: 11px; margin-left: 4px">（关于页引用块）</span></label>
          <textarea v-model="form.motto" class="textarea" rows="2" placeholder="如：「写代码是跟计算机对话，写博客是跟自己对话。」"></textarea>
        </div>
        <div style="display: grid; gap: 4px; grid-template-columns: 1fr 1fr">
          <div class="form-group">
            <label class="form-label">所在位置</label>
            <input v-model="form.location" class="input" type="text" placeholder="如：北京" />
          </div>
          <div class="form-group">
            <label class="form-label">建站年份</label>
            <input v-model.number="form.since" class="input" type="number" min="2000" :max="new Date().getFullYear()" />
          </div>
        </div>
        <div class="form-hint" style="margin-top: 8px">
          <Sparkles :size="12" style="vertical-align: -1px; color: var(--accent)" />
          这两项已与「基本信息」section 互通，仅展示在此处便于集中修改。
        </div>
        <div class="form-group" style="margin-top: 18px; margin-bottom: 0">
          <label class="form-label">首页终端内容 <span class="muted" style="font-size: 11px; margin-left: 4px">（首页 hero 打字机）</span></label>
          <textarea v-model="form.heroTerminal" class="textarea" rows="8" style="font-family: 'JetBrains Mono', monospace; font-size: 12.5px"></textarea>
          <div class="form-hint">每行一条，格式：<code>类型|文本</code>（tk=蓝色命令 / cm=灰色注释 / fn=紫色输出；不带类型默认蓝色）。空行会被忽略。</div>
        </div>
      </div>
    </div>

    <div style="display: flex; gap: 10px; justify-content: flex-end; padding-bottom: 30px">
      <button class="btn btn-ghost" :disabled="saving" @click="reset"><RotateCcw :size="15" /> 恢复默认</button>
      <button class="btn btn-primary" :disabled="saving" @click="save"><Save :size="15" /> 保存设置</button>
    </div>
  </div>
</template>
