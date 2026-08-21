<script setup lang="ts">
/** 柱状图 — 渐变柱 + 紧凑网格（对齐原手写 SVG 风格）。 */
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import BaseChart from './BaseChart.vue'

const props = withDefaults(defineProps<{ data: { label: string; value: number }[]; color?: [string, string] }>(), {
  color: () => ['#60a5fa', '#3b82f6'],
})

const option = computed<EChartsOption>(() => ({
  grid: { left: 44, right: 12, top: 18, bottom: 30 },
  tooltip: { trigger: 'axis', backgroundColor: 'rgba(13,18,30,.92)', borderColor: 'rgba(255,255,255,.1)', textStyle: { color: '#e7ebf3', fontSize: 12 } },
  xAxis: { type: 'category', data: props.data.map((d) => d.label), axisLine: { lineStyle: { color: 'var(--border)' } }, axisLabel: { color: 'var(--text-3)', fontSize: 10.5 } },
  yAxis: { type: 'value', splitLine: { lineStyle: { color: 'var(--border)' } }, axisLabel: { color: 'var(--text-3)', fontSize: 10 } },
  series: [
    {
      type: 'bar',
      data: props.data.map((d) => d.value),
      barMaxWidth: 34,
      itemStyle: {
        borderRadius: [5, 5, 0, 0],
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: props.color[0] },
            { offset: 1, color: props.color[1] },
          ],
        },
      },
    },
  ],
}))
</script>

<template>
  <BaseChart :option="option" />
</template>
