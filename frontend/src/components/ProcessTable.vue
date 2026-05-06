<template>
  <div class="proc-panel">
    <div class="proc-toolbar">
      <span class="panel-title">⚙️ Processes</span>
      <input v-model="search" class="search-input" placeholder="Filter processes..." />
      <select v-model="sortBy" class="sort-select">
        <option value="cpu_percent">Sort: CPU</option>
        <option value="memory_mb">Sort: RAM</option>
        <option value="name">Sort: Name</option>
        <option value="pid">Sort: PID</option>
      </select>
      <button class="btn-refresh" @click="fetchProcesses">↺ Refresh</button>
    </div>

    <div class="proc-table-wrap">
      <table class="proc-table">
        <thead>
          <tr>
            <th @click="sortBy='pid'">PID</th>
            <th @click="sortBy='name'">Name</th>
            <th @click="sortBy='cpu_percent'">CPU%</th>
            <th @click="sortBy='memory_mb'">RAM (MB)</th>
            <th @click="sortBy='mem_percent'">MEM%</th>
            <th>Status</th>
            <th>User</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in displayed" :key="p.pid" :class="rowClass(p)">
            <td class="mono">{{ p.pid }}</td>
            <td class="proc-name">{{ p.name }}</td>
            <td class="mono" :style="{ color: cpuColor(p.cpu_percent) }">{{ p.cpu_percent.toFixed(1) }}</td>
            <td class="mono">{{ p.memory_mb.toFixed(1) }}</td>
            <td class="mono">{{ p.mem_percent.toFixed(1) }}</td>
            <td><span class="status-badge" :class="p.status">{{ p.status }}</span></td>
            <td class="user-cell">{{ p.username }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="proc-footer">{{ filtered.length }} processes (showing {{ displayed.length }})</div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const processes = ref([])
const search = ref('')
const sortBy = ref('cpu_percent')
let timer = null

async function fetchProcesses() {
  try {
    const r = await fetch('/api/processes')
    processes.value = await r.json()
  } catch (_) {}
}

onMounted(() => {
  fetchProcesses()
  timer = setInterval(fetchProcesses, 3000)
})
onUnmounted(() => clearInterval(timer))

const filtered = computed(() => {
  const q = search.value.toLowerCase()
  return q
    ? processes.value.filter(p => p.name.toLowerCase().includes(q) || String(p.pid).includes(q))
    : processes.value
})

const sorted = computed(() =>
  [...filtered.value].sort((a, b) => {
    const va = a[sortBy.value], vb = b[sortBy.value]
    return typeof va === 'string' ? va.localeCompare(vb) : vb - va
  })
)

const displayed = computed(() => sorted.value.slice(0, 200))

function cpuColor(v) {
  if (v >= 50) return '#ff4455'
  if (v >= 20) return '#ffaa00'
  return '#00d4ff'
}

function rowClass(p) {
  if (p.cpu_percent >= 50) return 'row-hot'
  if (p.cpu_percent >= 20) return 'row-warm'
  return ''
}
</script>

<style scoped>
.proc-panel { background: #0d0d20; border: 1px solid #2a2a4a; border-radius: 10px; display: flex; flex-direction: column; overflow: hidden; }
.proc-toolbar { display: flex; align-items: center; gap: 8px; padding: 10px 14px; background: #12122a; border-bottom: 1px solid #2a2a4a; flex-wrap: wrap; }
.panel-title { font-size: 14px; font-weight: 700; color: #ccc; letter-spacing: 1px; flex-shrink: 0; }
.search-input { flex: 1; min-width: 120px; background: #0d0d20; border: 1px solid #3a3a5a; border-radius: 6px; color: #eee; padding: 5px 10px; font-size: 12px; outline: none; }
.search-input:focus { border-color: #00d4ff; }
.sort-select { background: #0d0d20; border: 1px solid #3a3a5a; border-radius: 6px; color: #eee; padding: 5px 8px; font-size: 12px; outline: none; }
.btn-refresh { background: #1a1a3a; border: 1px solid #3a3a5a; border-radius: 6px; color: #00d4ff; padding: 5px 12px; font-size: 12px; cursor: pointer; }
.btn-refresh:hover { background: #2a2a5a; }
.proc-table-wrap { overflow: auto; flex: 1; }
.proc-table { width: 100%; border-collapse: collapse; font-size: 12px; }
.proc-table th { background: #1a1a35; color: #888; padding: 8px 10px; text-align: left; border-bottom: 1px solid #2a2a4a; cursor: pointer; user-select: none; white-space: nowrap; }
.proc-table th:hover { color: #00d4ff; }
.proc-table td { padding: 6px 10px; border-bottom: 1px solid #1a1a2e; color: #ccc; white-space: nowrap; }
.proc-table tr:hover td { background: #1a1a35; }
.row-hot td { background: #2a0a0e; }
.row-warm td { background: #1a1500; }
.mono { font-family: monospace; }
.proc-name { max-width: 200px; overflow: hidden; text-overflow: ellipsis; }
.status-badge { font-size: 10px; padding: 2px 6px; border-radius: 10px; background: #1a1a3a; color: #888; }
.status-badge.running { background: #003a1a; color: #00cc66; }
.status-badge.sleep { background: #1a1a3a; color: #6688aa; }
.user-cell { color: #666; max-width: 100px; overflow: hidden; text-overflow: ellipsis; }
.proc-footer { padding: 6px 14px; font-size: 11px; color: #555; border-top: 1px solid #1a1a2e; background: #0d0d20; }
</style>
