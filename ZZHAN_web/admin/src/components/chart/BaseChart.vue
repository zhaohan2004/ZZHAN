<script setup lang="ts">
/**
 * ECharts 封装 — init / resize / dispose 生命周期。
 * 容器高度由父组件控制（如 style="height:280px"）。
 */
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'

const props = defineProps<{ option: EChartsOption }>()

const el = ref<HTMLDivElement | null>(null)
let chart: echarts.ECharts | null = null
let ro: ResizeObserver | null = null

onMounted(() => {
  if (!el.value) return
  chart = echarts.init(el.value)
  chart.setOption(props.option)
  if (typeof ResizeObserver !== 'undefined') {
    ro = new ResizeObserver(() => chart?.resize())
    ro.observe(el.value)
  }
})

watch(
  () => props.option,
  (o) => chart?.setOption(o, true),
  { deep: true },
)

onBeforeUnmount(() => {
  ro?.disconnect()
  chart?.dispose()
  chart = null
})
</script>

<template>
  <div ref="el" style="width: 100%; height: 100%" />
</template>
