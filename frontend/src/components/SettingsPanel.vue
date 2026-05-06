<template>
  <div class="settings-panel">
    <div class="sp-header">
      <span class="panel-title">⚙️ Alert Thresholds</span>
      <button class="btn-save" @click="save" :class="{ saved }">{{ saved ? '✓ Saved' : 'Save' }}</button>
    </div>
    <div class="sp-body">
      <div class="sp-row" v-for="(item, key) in thresholds" :key="key">
        <label>{{ item.label }}</label>
        <div class="slider-wrap">
          <input type="range" min="10" max="100" step="1" v-model.number="cfg[key]" class="slider" :style="{ '--pct': cfg[key] + '%' }" />
          <span class="slider-val">{{ cfg[key] }}%</span>
        </div>
      </div>

      <div class="sp-divider" />

      <label class="toggle-row">
        <span>🔊 Sound Alerts</span>
        <div class="toggle" :class="{ on: cfg.sound_enabled }" @click="cfg.sound_enabled = !cfg.sound_enabled">
          <div class="toggle-knob" />
        </div>
      </label>
      <label class="toggle-row">
        <span>📣 Notifications</span>
        <div class="toggle" :class="{ on: cfg.notify_enabled }" @click="cfg.notify_enabled = !cfg.notify_enabled">
          <div class="toggle-knob" />
        </div>
      </label>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useMetricsStore } from '../stores/metrics'

const store = useMetricsStore()
const saved = ref(false)

const thresholds = {
  cpu_threshold: { label: '🖥️ CPU Threshold' },
  mem_threshold: { label: '💾 Memory Threshold' },
  disk_threshold: { label: '💿 Disk Threshold' },
  gpu_threshold: { label: '🎮 GPU Threshold' },
}

const cfg = reactive({ ...store.alertConfig })

onMounted(async () => {
  try {
    const r = await fetch('/api/alerts/config')
    const data = await r.json()
    Object.assign(cfg, data)
    Object.assign(store.alertConfig, data)
  } catch (_) {}
})

async function save() {
  try {
    await fetch('/api/alerts/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(cfg),
    })
    Object.assign(store.alertConfig, cfg)
    saved.value = true
    setTimeout(() => (saved.value = false), 2000)
  } catch (_) {}
}
</script>

<style scoped>
.settings-panel { background: #0d0d20; border: 1px solid #2a2a4a; border-radius: 10px; }
.sp-header { display: flex; align-items: center; padding: 10px 14px; background: #12122a; border-bottom: 1px solid #2a2a4a; }
.panel-title { flex: 1; font-size: 14px; font-weight: 700; color: #ccc; letter-spacing: 1px; }
.btn-save { background: #003a6a; border: 1px solid #0066aa; border-radius: 6px; color: #44aaff; padding: 5px 14px; font-size: 12px; cursor: pointer; transition: all 0.3s; }
.btn-save.saved { background: #003a1a; border-color: #006622; color: #00cc66; }
.btn-save:hover { background: #0055aa; }
.sp-body { padding: 14px; display: flex; flex-direction: column; gap: 14px; }
.sp-row label { display: block; font-size: 12px; color: #aaa; margin-bottom: 6px; }
.slider-wrap { display: flex; align-items: center; gap: 10px; }
.slider { flex: 1; -webkit-appearance: none; height: 6px; background: linear-gradient(to right, #00d4ff var(--pct), #2a2a4a var(--pct)); border-radius: 3px; outline: none; cursor: pointer; }
.slider::-webkit-slider-thumb { -webkit-appearance: none; width: 16px; height: 16px; background: #00d4ff; border-radius: 50%; cursor: pointer; }
.slider-val { font-family: monospace; font-size: 13px; color: #00d4ff; width: 36px; text-align: right; }
.sp-divider { height: 1px; background: #2a2a4a; margin: 4px 0; }
.toggle-row { display: flex; align-items: center; gap: 10px; font-size: 13px; color: #aaa; cursor: pointer; }
.toggle-row span { flex: 1; }
.toggle { width: 40px; height: 22px; background: #2a2a4a; border-radius: 11px; position: relative; cursor: pointer; transition: background 0.3s; flex-shrink: 0; }
.toggle.on { background: #0066aa; }
.toggle-knob { width: 16px; height: 16px; background: #fff; border-radius: 50%; position: absolute; top: 3px; left: 3px; transition: left 0.3s; }
.toggle.on .toggle-knob { left: 21px; }
</style>
