<template>
  <div class="metric-card" :class="{ critical: pct >= threshold }">
    <div class="mc-top">
      <span class="mc-icon">{{ icon }}</span>
      <div class="mc-info">
        <div class="mc-label">{{ label }}</div>
        <div class="mc-value">{{ value }}</div>
      </div>
      <div class="mc-pct" :style="{ color: pctColor }">{{ pct.toFixed(1) }}%</div>
    </div>
    <div class="progress-bar">
      <div class="progress-fill" :style="{ width: pct + '%', background: pctColor }" />
    </div>
    <div v-if="sub" class="mc-sub">{{ sub }}</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  icon: String,
  label: String,
  value: String,
  pct: { type: Number, default: 0 },
  threshold: { type: Number, default: 80 },
  sub: String,
})

const pctColor = computed(() => {
  if (props.pct >= props.threshold) return '#ff4455'
  if (props.pct >= props.threshold * 0.75) return '#ffaa00'
  return '#00d4ff'
})
</script>

<style scoped>
.metric-card {
  background: #12122a;
  border: 1px solid #2a2a4a;
  border-radius: 10px;
  padding: 14px;
  transition: all 0.3s;
}
.metric-card.critical { border-color: #ff4455; box-shadow: 0 0 12px #ff445540; }
.mc-top { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.mc-icon { font-size: 24px; }
.mc-info { flex: 1; }
.mc-label { font-size: 11px; color: #888; text-transform: uppercase; letter-spacing: 1px; }
.mc-value { font-size: 15px; font-weight: 600; color: #eee; margin-top: 2px; }
.mc-pct { font-size: 22px; font-weight: 700; font-family: monospace; }
.progress-bar { height: 6px; background: #2a2a4a; border-radius: 3px; overflow: hidden; }
.progress-fill { height: 100%; border-radius: 3px; transition: width 0.8s ease, background 0.3s; }
.mc-sub { font-size: 11px; color: #666; margin-top: 6px; }
</style>
