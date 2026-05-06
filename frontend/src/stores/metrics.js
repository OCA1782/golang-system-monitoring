import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

const HISTORY_SIZE = 60

function mkSeries() {
  return Array(HISTORY_SIZE).fill(0)
}

export const useMetricsStore = defineStore('metrics', () => {
  const connected = ref(false)
  const latest = ref(null)
  const alerts = ref([])

  const cpuHistory = ref(mkSeries())
  const memHistory = ref(mkSeries())
  const netSentHistory = ref(mkSeries())
  const netRecvHistory = ref(mkSeries())
  const gpuHistory = ref(mkSeries())

  const alertConfig = ref({
    cpu_threshold: 80,
    mem_threshold: 85,
    disk_threshold: 90,
    gpu_threshold: 85,
    sound_enabled: true,
    notify_enabled: true,
  })

  function push(arr, val) {
    arr.shift()
    arr.push(val)
  }

  function update(m) {
    latest.value = m
    push(cpuHistory.value, m.cpu.percent)
    push(memHistory.value, m.memory.percent)

    const totalSent = m.networks.reduce((s, n) => s + n.bytes_sent_sec, 0)
    const totalRecv = m.networks.reduce((s, n) => s + n.bytes_recv_sec, 0)
    push(netSentHistory.value, totalSent)
    push(netRecvHistory.value, totalRecv)

    if (m.gpu && m.gpu.length > 0) {
      push(gpuHistory.value, m.gpu[0].utilization)
    }

    checkAlerts(m)
  }

  function checkAlerts(m) {
    const cfg = alertConfig.value
    if (!cfg.notify_enabled) return

    const now = Date.now()
    const addAlert = (type, msg, pct) => {
      alerts.value.unshift({ id: now + Math.random(), type, msg, pct, time: new Date().toLocaleTimeString() })
      if (alerts.value.length > 50) alerts.value.pop()
      if (cfg.sound_enabled) playBeep(type)
    }

    if (m.cpu.percent >= cfg.cpu_threshold) addAlert('cpu', `CPU at ${m.cpu.percent.toFixed(1)}%`, m.cpu.percent)
    if (m.memory.percent >= cfg.mem_threshold) addAlert('mem', `RAM at ${m.memory.percent.toFixed(1)}%`, m.memory.percent)
    m.disks?.forEach(d => {
      if (d.percent >= cfg.disk_threshold) addAlert('disk', `Disk ${d.mountpoint} at ${d.percent.toFixed(1)}%`, d.percent)
    })
    m.gpu?.forEach(g => {
      if (g.utilization >= cfg.gpu_threshold) addAlert('gpu', `GPU at ${g.utilization.toFixed(1)}%`, g.utilization)
    })
  }

  function playBeep(type) {
    try {
      const ctx = new (window.AudioContext || window.webkitAudioContext)()
      const osc = ctx.createOscillator()
      const gain = ctx.createGain()
      osc.connect(gain)
      gain.connect(ctx.destination)
      osc.frequency.value = type === 'cpu' ? 880 : type === 'mem' ? 660 : 440
      gain.gain.setValueAtTime(0.3, ctx.currentTime)
      gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.4)
      osc.start(ctx.currentTime)
      osc.stop(ctx.currentTime + 0.4)
    } catch (_) {}
  }

  const timeLabels = computed(() => Array.from({ length: HISTORY_SIZE }, (_, i) => `-${HISTORY_SIZE - i}s`))

  return {
    connected, latest, alerts, alertConfig,
    cpuHistory, memHistory, netSentHistory, netRecvHistory, gpuHistory,
    timeLabels, update,
  }
})
