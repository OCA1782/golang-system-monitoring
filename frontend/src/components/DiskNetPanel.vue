<template>
  <div class="dn-grid">
    <!-- Disk Usage -->
    <div class="dn-card">
      <div class="dn-title">💿 Disk Usage</div>
      <div v-if="disks.length" class="disk-list">
        <div v-for="d in disks" :key="d.mountpoint" class="disk-item">
          <div class="disk-top">
            <span class="disk-mount">{{ d.mountpoint }}</span>
            <span class="disk-info">{{ fmtBytes(d.used) }} / {{ fmtBytes(d.total) }}</span>
            <span class="disk-pct" :style="{ color: pctColor(d.percent) }">{{ d.percent.toFixed(1) }}%</span>
          </div>
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: d.percent + '%', background: pctColor(d.percent) }" />
          </div>
          <div class="disk-io">
            ↑ {{ fmtSpeed(d.write_bytes_sec) }}/s &nbsp; ↓ {{ fmtSpeed(d.read_bytes_sec) }}/s
          </div>
        </div>
      </div>
      <div v-else class="empty">No disk data</div>
    </div>

    <!-- Network -->
    <div class="dn-card">
      <div class="dn-title">🌐 Network I/O</div>
      <div v-if="networks.length" class="net-list">
        <div v-for="n in activeNets" :key="n.interface" class="net-item">
          <div class="net-name">{{ n.interface }}</div>
          <div class="net-stats">
            <span class="net-send">▲ {{ fmtSpeed(n.bytes_sent_sec) }}/s</span>
            <span class="net-recv">▼ {{ fmtSpeed(n.bytes_recv_sec) }}/s</span>
          </div>
        </div>
      </div>
      <div v-else class="empty">No network data</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useMetricsStore } from '../stores/metrics'

const store = useMetricsStore()

const disks = computed(() => store.latest?.disks ?? [])
const networks = computed(() => store.latest?.networks ?? [])
const activeNets = computed(() => networks.value.filter(n => n.bytes_sent_sec > 0 || n.bytes_recv_sec > 0).slice(0, 6))

function fmtBytes(b) {
  if (!b) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let v = b
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(1)} ${units[i]}`
}

function fmtSpeed(b) {
  if (!b || b < 1024) return `${b ?? 0} B`
  if (b < 1024 * 1024) return `${(b / 1024).toFixed(1)} KB`
  return `${(b / 1024 / 1024).toFixed(1)} MB`
}

function pctColor(p) {
  if (p >= 90) return '#ff4455'
  if (p >= 70) return '#ffaa00'
  return '#00d4ff'
}
</script>

<style scoped>
.dn-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.dn-card { background: #12122a; border: 1px solid #2a2a4a; border-radius: 10px; padding: 14px; }
.dn-title { font-size: 13px; font-weight: 700; color: #ccc; letter-spacing: 1px; margin-bottom: 12px; text-transform: uppercase; }
.disk-item { margin-bottom: 12px; }
.disk-top { display: flex; align-items: center; gap: 6px; margin-bottom: 4px; font-size: 12px; }
.disk-mount { font-weight: 700; color: #eee; flex: 1; }
.disk-info { color: #666; font-size: 11px; }
.disk-pct { font-family: monospace; font-weight: 700; }
.progress-bar { height: 5px; background: #2a2a4a; border-radius: 3px; overflow: hidden; margin-bottom: 3px; }
.progress-fill { height: 100%; border-radius: 3px; transition: width 0.8s; }
.disk-io { font-size: 10px; color: #666; }
.net-item { display: flex; align-items: center; gap: 10px; padding: 7px 0; border-bottom: 1px solid #1a1a2e; font-size: 12px; }
.net-name { flex: 1; color: #aaa; font-family: monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.net-stats { display: flex; gap: 12px; }
.net-send { color: #ff9944; }
.net-recv { color: #44aaff; }
.empty { color: #555; font-size: 13px; text-align: center; padding: 20px; }
@media (max-width: 700px) { .dn-grid { grid-template-columns: 1fr; } }
</style>
