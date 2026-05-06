<template>
  <div class="alert-panel">
    <div class="ap-header">
      <span class="panel-title">🔔 Alerts</span>
      <span class="alert-count" v-if="store.alerts.length">{{ store.alerts.length }}</span>
      <button class="btn-clear" @click="store.alerts = []">Clear</button>
    </div>
    <div class="alerts-list">
      <transition-group name="alert-fade">
        <div v-for="a in store.alerts.slice(0, 20)" :key="a.id" class="alert-item" :class="a.type">
          <span class="alert-icon">{{ icons[a.type] }}</span>
          <span class="alert-msg">{{ a.msg }}</span>
          <span class="alert-time">{{ a.time }}</span>
        </div>
      </transition-group>
      <div v-if="!store.alerts.length" class="no-alerts">No alerts — all systems normal</div>
    </div>
  </div>
</template>

<script setup>
import { useMetricsStore } from '../stores/metrics'
const store = useMetricsStore()
const icons = { cpu: '🔥', mem: '💾', disk: '💿', gpu: '🎮' }
</script>

<style scoped>
.alert-panel { background: #0d0d20; border: 1px solid #2a2a4a; border-radius: 10px; display: flex; flex-direction: column; overflow: hidden; }
.ap-header { display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #12122a; border-bottom: 1px solid #2a2a4a; }
.panel-title { flex: 1; font-size: 14px; font-weight: 700; color: #ccc; letter-spacing: 1px; }
.alert-count { background: #ff4455; color: #fff; font-size: 11px; font-weight: 700; padding: 2px 7px; border-radius: 10px; }
.btn-clear { background: #1a1a3a; border: 1px solid #3a3a5a; border-radius: 6px; color: #888; padding: 4px 10px; font-size: 11px; cursor: pointer; }
.btn-clear:hover { color: #ff4455; border-color: #ff4455; }
.alerts-list { overflow-y: auto; flex: 1; padding: 6px 0; }
.alert-item { display: flex; align-items: center; gap: 8px; padding: 7px 14px; border-bottom: 1px solid #1a1a2e; font-size: 12px; }
.alert-item.cpu { background: #1a0a0e; }
.alert-item.mem { background: #0a0a1a; }
.alert-item.disk { background: #0a1a0a; }
.alert-item.gpu { background: #1a0a1a; }
.alert-icon { font-size: 16px; }
.alert-msg { flex: 1; color: #ddd; }
.alert-time { color: #555; font-size: 11px; font-family: monospace; }
.no-alerts { padding: 20px; text-align: center; color: #555; font-size: 13px; }
.alert-fade-enter-active { animation: slideIn 0.3s ease; }
@keyframes slideIn { from { opacity: 0; transform: translateY(-10px); } to { opacity: 1; transform: none; } }
</style>
