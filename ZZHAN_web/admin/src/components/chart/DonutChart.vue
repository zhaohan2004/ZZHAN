<script setup lang="ts">
/** 环形图 — 分类占比，中心显示总数。 */
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import BaseChart from './BaseChart.vue'

const props = defineProps<{ items: { name: string; value: number; color: string }[]; total: number }>()

const option = computed<EChartsOption>(() => ({
  tooltip: { trigger: 'item', backgroundColor: 'rgba(13,18,30,.92)', borderColor: 'rgba(255,255,255,.1)', textStyle: { color: '#e7ebf3', fontSize: 12 } },
  legend: { show: false },
  series: [
    {
      type: 'pie',
      radius: ['58%', '78%'],
      center: ['50%', '50%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 4, borderColor: 'var(--card)', borderWidth: 2 },
      label: { show: false },
      data: props.items.map((it) => ({ name: it.name, value: it.value, itemStyle: { color: it.color } })),
    },
  ],
  graphic: [
    {
      type: 'text',
      left: 'center',
      top: '40%',
      style: { text: String(props.total), textAlign: 'center', fill: 'var(--text)', fontSize: 24, fontWeight: 800, fontFamily: 'JetBrains Mono, monospace' },
    },
    {
      type: 'text',
      left: 'center',
      top: '56%',
      style: { text: '文章总数', textAlign: 'center', fill: 'var(--text-3)', fontSize: 11 },
    },
  ],
}))
</script>

<template>
  <div style="display: flex; align-items: center; gap: 12px">
    <div style="height: 200px; width: 200px; flex: none"><BaseChart :option="option" /></div>
    <ul style="min-width: 0; flex: 1; display: flex; flex-direction: column; gap: 6px; margin: 0; padding: 0; list-style: none">
      <li v-for="it in items" :key="it.name" style="display: flex; align-items: center; gap: 8px; font-size: 12.5px; color: var(--text-2)">
        <span style="height: 10px; width: 10px; flex: none; border-radius: 3px" :style="{ background: it.color }" />
        <span style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{ it.name }}</span>
        <span style="margin-left: auto; font-family: 'JetBrains Mono', monospace; color: var(--text-3)">{{ it.value }}</span>
      </li>
    </ul>
  </div>
</template>
