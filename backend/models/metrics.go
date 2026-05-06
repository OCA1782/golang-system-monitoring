package models

type CPUMetric struct {
	Percent  float64   `json:"percent"`
	PerCore  []float64 `json:"per_core"`
	Freq     float64   `json:"freq_mhz"`
}

type MemoryMetric struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	Percent     float64 `json:"percent"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapPercent float64 `json:"swap_percent"`
}

type DiskMetric struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Percent     float64 `json:"percent"`
	ReadBytes   uint64  `json:"read_bytes_sec"`
	WriteBytes  uint64  `json:"write_bytes_sec"`
}

type NetworkMetric struct {
	Interface  string `json:"interface"`
	BytesSent  uint64 `json:"bytes_sent_sec"`
	BytesRecv  uint64 `json:"bytes_recv_sec"`
	PacketsSent uint64 `json:"packets_sent_sec"`
	PacketsRecv uint64 `json:"packets_recv_sec"`
}

type GPUMetric struct {
	Name        string  `json:"name"`
	Utilization float64 `json:"utilization"`
	MemoryUsed  uint64  `json:"memory_used"`
	MemoryTotal uint64  `json:"memory_total"`
	Temperature float64 `json:"temperature"`
}

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryMB   float64 `json:"memory_mb"`
	MemPercent float32 `json:"mem_percent"`
	Status     string  `json:"status"`
	Username   string  `json:"username"`
}

type SystemMetrics struct {
	Timestamp int64           `json:"timestamp"`
	CPU       CPUMetric       `json:"cpu"`
	Memory    MemoryMetric    `json:"memory"`
	Disks     []DiskMetric    `json:"disks"`
	Networks  []NetworkMetric `json:"networks"`
	GPU       []GPUMetric     `json:"gpu"`
	Uptime    uint64          `json:"uptime"`
	HostInfo  HostInfo        `json:"host"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
	Platform string `json:"platform"`
	Arch     string `json:"arch"`
	KernelVersion string `json:"kernel_version"`
}

type AlertConfig struct {
	CPUThreshold    float64 `json:"cpu_threshold"`
	MemThreshold    float64 `json:"mem_threshold"`
	DiskThreshold   float64 `json:"disk_threshold"`
	GPUThreshold    float64 `json:"gpu_threshold"`
	SoundEnabled    bool    `json:"sound_enabled"`
	NotifyEnabled   bool    `json:"notify_enabled"`
}
