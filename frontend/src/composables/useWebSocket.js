import { onMounted, onUnmounted } from 'vue'
import { useMetricsStore } from '../stores/metrics'

export function useWebSocket() {
  const store = useMetricsStore()
  let ws = null
  let retryTimer = null

  function connect() {
    const proto = location.protocol === 'https:' ? 'wss' : 'ws'
    ws = new WebSocket(`${proto}://${location.host}/ws`)

    ws.onopen = () => { store.connected = true }
    ws.onclose = () => {
      store.connected = false
      retryTimer = setTimeout(connect, 3000)
    }
    ws.onerror = () => ws.close()
    ws.onmessage = (e) => {
      try {
        const msg = JSON.parse(e.data)
        if (msg.type === 'metrics') store.update(msg.payload)
      } catch (_) {}
    }
  }

  onMounted(connect)
  onUnmounted(() => {
    clearTimeout(retryTimer)
    ws?.close()
  })
}
