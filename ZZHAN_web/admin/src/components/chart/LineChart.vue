<script setup lang="ts">
/** 折线图 — 平滑曲线 + 渐变面积 + 点标记。 */
import { computed } from 'vue'
import type { EChartsOption } from 'echarts'
import BaseChart from './BaseChart.vue'

const props = defineProps<{ series: { data: number[]; color: string; name?: string }[]; labels: string[] }>()

const option = computed<EChartsOption>(() => ({
  grid: { left: 44, right: 14, top: 26, bottom: 30 },
  tooltip: { trigger: 'axis', backgroundColor: 'rgba(13,18,30,.92)', borderColor: 'rgba(255,255,255,.1)', textStyle: { color: '#e7ebf3', fontSize: 12 } },
  legend: { show: props.series.length > 1, top: 0, textStyle: { color: 'var(--text-2)', fontSize: 12 }, itemWidth: 14, itemHeight: 8 },
  xAxis: { type: 'category', boundaryGap: false, data: props.labels, axisLine: { lineStyle: { color: 'var(--border)' } }, axisLabel: { color: 'var(--text-3)', fontSize: 10.5 } },
  yAxis: { type: 'value', splitLine: { lineStyle: { color: 'var(--border)' } }, axisLabel: { color: 'var(--text-3)', fontSize: 10 } },
  series: props.series.map((s) => ({
    name: s.name,
    type: 'line',
    smooth: true,
    symbolSize: 4,
    data: s.data,
    lineStyle: { width: 2.5, color: s.color },
    itemStyle: { color: s.color },
    areaStyle: {
      color: {
        type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
        colorStops: [
          { offset: 0, color: s.color + '55' },
          { offset: 1, color: s.color + '00' },
        ],
      },
    },
  })),
}))
</script>

<template>
  <BaseChart :option="option" />
</template>
