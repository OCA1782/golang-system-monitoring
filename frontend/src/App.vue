<template>
  <div class="app" :class="theme">
    <!-- Title Bar -->
    <div class="title-bar">
      <div class="tb-left">
        <span class="app-icon">📊</span>
        <span class="app-title">SysMonitor</span>
        <span class="app-sub">System Resource Monitor</span>
      </div>
      <div class="tb-center">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >{{ tab.label }}</button>
      </div>
      <div class="tb-right">
        <button class="icon-btn" @click="exportCSV" title="Export CSV">⬇️</button>
        <button class="icon-btn" @click="theme = theme === 'dark' ? 'light' : 'dark'" title="Toggle Theme">
          {{ theme === 'dark' ? '☀️' : '🌙' }}
        </button>
      </div>
    </div>

    <!-- Main content -->
    <div class="main-content">

      <!-- OVERVIEW TAB -->
      <div v-show="activeTab === 'overview'" class="tab-content">
        <!-- Metric cards row -->
        <div class="cards-row">
          <MetricCard
            icon="🖥️" label="CPU Usage"
            :value="`${store.latest?.cpu?.per_core?.length ?? 0} cores @ ${(store.latest?.cpu?.freq_mhz ?? 0).toFixed(0)} MHz`"
            :pct="store.latest?.cpu?.percent ?? 0"
            :threshold="store.alertConfig.cpu_threshold"
          />
          <MetricCard
            icon="💾" label="RAM Usage"
            :value="fmtBytes(store.latest?.memory?.used) + ' / ' + fmtBytes(store.latest?.memory?.total)"
            :pct="store.latest?.memory?.percent ?? 0"
            :threshold="store.alertConfig.mem_threshold"
            :sub="store.latest?.memory?.swap_total > 0 ? `Swap: ${fmtBytes(store.latest?.memory?.swap_used)} / ${fmtBytes(store.latest?.memory?.swap_total)}` : ''"
          />
          <MetricCard v-if="store.latest?.gpu?.length"
            icon="🎮" label="GPU Usage"
            :value="store.latest?.gpu[0]?.name ?? 'GPU'"
            :pct="store.latest?.gpu[0]?.utilization ?? 0"
            :threshold="store.alertConfig.gpu_threshold"
            :sub="`VRAM: ${fmtBytes(store.latest?.gpu[0]?.memory_used)} / ${fmtBytes(store.latest?.gpu[0]?.memory_total)} — Temp: ${store.latest?.gpu[0]?.temperature}°C`"
          />
          <MetricCard v-else
            icon="🎮" label="GPU"
            value="No NVIDIA GPU detected"
            :pct="0"
          />
        </div>

        <!-- Charts -->
        <div class="charts-grid">
          <RealtimeChart
            title="CPU" icon="🖥️"
            :history="store.cpuHistory"
            :labels="store.timeLabels"
            color="#00d4ff"
            :threshold="store.alertConfig.cpu_threshold"
          />
          <RealtimeChart
            title="Memory" icon="💾"
            :history="store.memHistory"
            :labels="store.timeLabels"
            color="#aa44ff"
            :threshold="store.alertConfig.mem_threshold"
          />
          <RealtimeChart
            title="GPU Util" icon="🎮"
            :history="store.gpuHistory"
            :labels="store.timeLabels"
            color="#ff6644"
            :threshold="store.alertConfig.gpu_threshold"
          />
          <RealtimeChart
            title="Net Upload" icon="🌐"
            :history="store.netSentHistory"
            :labels="store.timeLabels"
            color="#ffaa00"
            unit=""
            :formatFn="fmtSpeed"
          />
        </div>

        <!-- CPU per-core mini bars -->
        <div v-if="store.latest?.cpu?.per_core?.length" class="core-section">
          <div class="section-label">CPU Cores</div>
          <div class="core-bars">
            <div v-for="(v, i) in store.latest.cpu.per_core" :key="i" class="core-bar">
              <div class="core-label">C{{ i }}</div>
              <div class="core-track">
                <div class="core-fill" :style="{ height: v + '%', background: coreColor(v) }" />
              </div>
              <div class="core-val">{{ v.toFixed(0) }}%</div>
            </div>
          </div>
        </div>

        <DiskNetPanel />
      </div>

      <!-- PROCESSES TAB -->
      <div v-show="activeTab === 'processes'" class="tab-content tab-full">
        <ProcessTable />
      </div>

      <!-- ALERTS TAB -->
      <div v-show="activeTab === 'alerts'" class="tab-content tab-two-col">
        <AlertPanel />
        <SettingsPanel />
      </div>

    </div>

    <StatusBar />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useMetricsStore } from './stores/metrics'
import { useWebSocket } from './composables/useWebSocket'
import MetricCard from './components/MetricCard.vue'
import RealtimeChart from './components/RealtimeChart.vue'
import ProcessTable from './components/ProcessTable.vue'
import AlertPanel from './components/AlertPanel.vue'
import SettingsPanel from './components/SettingsPanel.vue'
import DiskNetPanel from './components/DiskNetPanel.vue'
import StatusBar from './components/StatusBar.vue'

const store = useMetricsStore()
useWebSocket()

const theme = ref('dark')
const activeTab = ref('overview')

const tabs = [
  { id: 'overview', label: '📊 Overview' },
  { id: 'processes', label: '⚙️ Processes' },
  { id: 'alerts', label: '🔔 Alerts' },
]

function fmtBytes(b) {
  if (!b) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0; let v = b
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return `${v.toFixed(1)} ${units[i]}`
}

function fmtSpeed(b) {
  if (!b || b < 1024) return `${b ?? 0} B`
  if (b < 1024 * 1024) return `${(b / 1024).toFixed(1)} KB`
  return `${(b / 1024 / 1024).toFixed(2)} MB`
}

function coreColor(v) {
  if (v >= 80) return '#ff4455'
  if (v >= 50) return '#ffaa00'
  return '#00d4ff'
}

function exportCSV() {
  if (!store.latest) return
  const m = store.latest
  const rows = [
    ['timestamp', 'cpu_percent', 'mem_percent', 'mem_used_gb', 'uptime_s'],
    [new Date(m.timestamp).toISOString(), m.cpu.percent.toFixed(2), m.memory.percent.toFixed(2), (m.memory.used/1e9).toFixed(2), m.uptime],
  ]
  const csv = rows.map(r => r.join(',')).join('\n')
  const a = document.createElement('a')
  a.href = URL.createObjectURL(new Blob([csv], { type: 'text/csv' }))
  a.download = `sysmonitor_${Date.now()}.csv`
  a.click()
}
</script>

<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body { background: #080818; color: #eee; font-family: 'Segoe UI', system-ui, sans-serif; overflow: hidden; }
#app { height: 100vh; display: flex; flex-direction: column; }

.app { display: flex; flex-direction: column; height: 100vh; background: #080818; color: #eee; }
.app.light { background: #f0f2f8; color: #111; }
.app.light .title-bar { background: #e0e4f0; border-bottom-color: #c0c8e0; }
.app.light .main-content { background: #f0f2f8; }

/* Title Bar */
.title-bar {
  display: flex; align-items: center; gap: 12px;
  background: #0d0d20; border-bottom: 1px solid #2a2a4a;
  padding: 0 16px; height: 50px; flex-shrink: 0;
  user-select: none;
}
.tb-left { display: flex; align-items: center; gap: 8px; min-width: 200px; }
.app-icon { font-size: 22px; }
.app-title { font-size: 16px; font-weight: 800; color: #00d4ff; letter-spacing: 2px; }
.app-sub { font-size: 11px; color: #555; margin-left: 4px; }
.tb-center { display: flex; gap: 2px; flex: 1; justify-content: center; }
.tb-right { display: flex; gap: 8px; min-width: 80px; justify-content: flex-end; }

.tab-btn {
  background: transparent; border: none; color: #666;
  padding: 6px 16px; border-radius: 6px; cursor: pointer;
  font-size: 13px; font-weight: 600; transition: all 0.2s;
}
.tab-btn:hover { background: #1a1a3a; color: #aaa; }
.tab-btn.active { background: #1a1a3a; color: #00d4ff; border-bottom: 2px solid #00d4ff; }

.icon-btn { background: transparent; border: none; font-size: 18px; cursor: pointer; padding: 4px; border-radius: 6px; }
.icon-btn:hover { background: #1a1a3a; }

/* Main content */
.main-content { flex: 1; overflow-y: auto; padding: 12px 16px; display: flex; flex-direction: column; gap: 12px; }

/* Scrollbar */
::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: #080818; }
::-webkit-scrollbar-thumb { background: #2a2a4a; border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: #4a4a6a; }

/* Tab layouts */
.tab-content { display: flex; flex-direction: column; gap: 12px; }
.tab-full { flex: 1; height: 100%; }
.tab-two-col { display: grid !important; grid-template-columns: 1fr 320px; gap: 12px; align-items: start; }

/* Cards */
.cards-row { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 12px; }

/* Charts grid */
.charts-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(280px, 1fr)); gap: 12px; }

/* CPU Core bars */
.core-section { background: #12122a; border: 1px solid #2a2a4a; border-radius: 10px; padding: 14px; }
.section-label { font-size: 12px; font-weight: 700; color: #888; text-transform: uppercase; letter-spacing: 1px; margin-bottom: 10px; }
.core-bars { display: flex; gap: 8px; flex-wrap: wrap; }
.core-bar { display: flex; flex-direction: column; align-items: center; gap: 3px; }
.core-label { font-size: 10px; color: #666; }
.core-track { width: 20px; height: 60px; background: #2a2a4a; border-radius: 3px; overflow: hidden; display: flex; align-items: flex-end; }
.core-fill { width: 100%; border-radius: 3px; transition: height 0.8s ease; }
.core-val { font-size: 10px; font-family: monospace; color: #888; }

@media (max-width: 768px) {
  .tab-two-col { grid-template-columns: 1fr; }
  .tb-center { display: none; }
}
</style>
