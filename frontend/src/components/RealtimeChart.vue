<template>
  <div class="chart-card">
    <div class="chart-header">
      <span class="chart-icon">{{ icon }}</span>
      <span class="chart-title">{{ title }}</span>
      <span class="chart-value" :class="valueClass">{{ displayValue }}</span>
    </div>
    <apexchart type="area" height="120" :options="chartOptions" :series="series" />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: String,
  icon: String,
  history: Array,
  labels: Array,
  color: { type: String, default: '#00d4ff' },
  unit: { type: String, default: '%' },
  threshold: { type: Number, default: 80 },
  formatFn: { type: Function, default: null },
})

const displayValue = computed(() => {
  const last = props.history[props.history.length - 1] ?? 0
  return props.formatFn ? props.formatFn(last) : `${last.toFixed(1)}${props.unit}`
})

const valueClass = computed(() => {
  const last = props.history[props.history.length - 1] ?? 0
  if (props.unit !== '%') return ''
  if (last >= props.threshold) return 'danger'
  if (last >= props.threshold * 0.75) return 'warn'
  return 'ok'
})

const series = computed(() => [{ name: props.title, data: [...props.history] }])

const chartOptions = computed(() => ({
  chart: {
    background: 'transparent',
    toolbar: { show: false },
    animations: { enabled: true, easing: 'linear', dynamicAnimation: { speed: 800 } },
    sparkline: { enabled: false },
  },
  stroke: { curve: 'smooth', width: 2 },
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.4,
      opacityTo: 0.05,
      stops: [0, 100],
    },
  },
  colors: [props.color],
  xaxis: {
    categories: props.labels,
    labels: { show: false },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    min: 0,
    max: props.unit === '%' ? 100 : undefined,
    labels: {
      style: { colors: '#888', fontSize: '10px' },
      formatter: (v) => props.formatFn ? props.formatFn(v) : `${v.toFixed(0)}${props.unit}`,
    },
  },
  grid: { borderColor: '#2a2a3e', strokeDashArray: 3 },
  tooltip: {
    theme: 'dark',
    y: { formatter: (v) => props.formatFn ? props.formatFn(v) : `${v.toFixed(1)}${props.unit}` },
  },
  dataLabels: { enabled: false },
}))
</script>

<style scoped>
.chart-card {
  background: #12122a;
  border: 1px solid #2a2a4a;
  border-radius: 10px;
  padding: 12px 16px 4px;
  transition: border-color 0.3s;
}
.chart-card:hover { border-color: #4a4a8a; }
.chart-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.chart-icon { font-size: 18px; }
.chart-title { flex: 1; font-size: 13px; font-weight: 600; color: #ccc; text-transform: uppercase; letter-spacing: 1px; }
.chart-value { font-size: 20px; font-weight: 700; font-family: monospace; }
.chart-value.ok   { color: #00d4ff; }
.chart-value.warn { color: #ffaa00; }
.chart-value.danger { color: #ff4455; animation: pulse 1s infinite; }
@keyframes pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.6; } }
</style>
