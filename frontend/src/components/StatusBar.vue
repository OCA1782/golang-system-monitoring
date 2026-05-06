<template>
  <div class="status-bar">
    <div class="sb-left">
      <span class="dot" :class="store.connected ? 'online' : 'offline'" />
      <span class="sb-label">{{ store.connected ? 'CONNECTED' : 'RECONNECTING...' }}</span>
    </div>
    <div class="sb-center">
      <span v-if="host">🖥️ {{ host.hostname }} &nbsp;|&nbsp; {{ host.os }} {{ host.platform }} &nbsp;|&nbsp; {{ host.arch }}</span>
    </div>
    <div class="sb-right">
      <span v-if="uptime">⏱️ Uptime: {{ fmtUptime(uptime) }}</span>
      &nbsp;|&nbsp;
      <span>🕐 {{ now }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMetricsStore } from '../stores/metrics'

const store = useMetricsStore()
const now = ref('')
let timer = null

const host = computed(() => store.latest?.host)
const uptime = computed(() => store.latest?.uptime)

function fmtUptime(s) {
  const d = Math.floor(s / 86400)
  const h = Math.floor((s % 86400) / 3600)
  const m = Math.floor((s % 3600) / 60)
  if (d > 0) return `${d}d ${h}h ${m}m`
  if (h > 0) return `${h}h ${m}m`
  return `${m}m`
}

onMounted(() => {
  const tick = () => { now.value = new Date().toLocaleTimeString() }
  tick()
  timer = setInterval(tick, 1000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.status-bar {
  display: flex;
  align-items: center;
  background: #080818;
  border-top: 1px solid #2a2a4a;
  padding: 4px 16px;
  font-size: 11px;
  color: #666;
  gap: 12px;
  flex-wrap: wrap;
}
.sb-left, .sb-right, .sb-center { display: flex; align-items: center; gap: 6px; }
.sb-center { flex: 1; justify-content: center; }
.dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.dot.online { background: #00cc66; box-shadow: 0 0 6px #00cc66; animation: glow 2s infinite; }
.dot.offline { background: #ff4455; }
.sb-label { font-weight: 700; color: #aaa; }
@keyframes glow { 0%,100% { box-shadow: 0 0 4px #00cc66; } 50% { box-shadow: 0 0 10px #00cc66; } }
</style>
