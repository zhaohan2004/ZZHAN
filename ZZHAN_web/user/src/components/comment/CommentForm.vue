<script setup lang="ts">
/** 评论输入表单 — 提交后提示「审核后展示」，成功后清空并通知父级刷新。 */
import { ref } from 'vue'
import { Send } from 'lucide-vue-next'
import { postComment } from '@/api/comments'
import { useToast } from '@/composables/useToast'

const props = withDefaults(defineProps<{ slug: string; parentId?: number | null }>(), { parentId: null })
const emit = defineEmits<{ (e: 'done'): void }>()

const { toast } = useToast()
const content = ref('')
const submitting = ref(false)

async function submit(): Promise<void> {
  const text = content.value.trim()
  if (!text) {
    toast('请先输入评论内容', 'error')
    return
  }
  submitting.value = true
  try {
    await postComment(props.slug, { content: text, parent_id: props.parentId })
    content.value = ''
    toast('评论已发布', 'success')
    emit('done')
  } catch {
    toast('评论提交失败，请重试', 'error')
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="comment-input-wrap" style="margin-bottom:0">
    <div style="flex:1">
      <textarea
        v-model="content"
        class="textarea"
        rows="3"
        :placeholder="parentId ? '写下你的回复…' : '写下你的看法…'"
      />
      <div style="display:flex;justify-content:flex-end;margin-top:10px">
        <button class="btn btn-primary btn-sm" type="button" :disabled="submitting" @click="submit">
          <Send :size="14" /> {{ submitting ? '提交中…' : '发表评论' }}
        </button>
      </div>
    </div>
  </div>
</template>
