package handlers

import (
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"sysmonitor/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type MetricsCollector struct {
	mu            sync.RWMutex
	lastNetStats  map[string]net.IOCountersStat
	lastDiskStats map[string]disk.IOCountersStat
	lastCollect   time.Time
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		lastNetStats:  make(map[string]net.IOCountersStat),
		lastDiskStats: make(map[string]disk.IOCountersStat),
		lastCollect:   time.Now(),
	}
}

func (mc *MetricsCollector) Collect() (*models.SystemMetrics, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(mc.lastCollect).Seconds()
	if elapsed < 0.1 {
		elapsed = 1.0
	}

	metrics := &models.SystemMetrics{
		Timestamp: now.UnixMilli(),
	}

	// CPU
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		metrics.CPU.Percent = cpuPercent[0]
	}
	perCore, err := cpu.Percent(0, true)
	if err == nil {
		metrics.CPU.PerCore = perCore
	}
	freqs, err := cpu.Info()
	if err == nil && len(freqs) > 0 {
		metrics.CPU.Freq = freqs[0].Mhz
	}

	// Memory
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		metrics.Memory = models.MemoryMetric{
			Total:     vmStat.Total,
			Used:      vmStat.Used,
			Available: vmStat.Available,
			Percent:   vmStat.UsedPercent,
		}
	}
	swapStat, err := mem.SwapMemory()
	if err == nil {
		metrics.Memory.SwapTotal = swapStat.Total
		metrics.Memory.SwapUsed = swapStat.Used
		metrics.Memory.SwapPercent = swapStat.UsedPercent
	}

	// Disk
	parts, err := disk.Partitions(false)
	if err == nil {
		diskIOMap := make(map[string]disk.IOCountersStat)
		diskIO, _ := disk.IOCounters()
		for k, v := range diskIO {
			diskIOMap[k] = v
		}

		for _, p := range parts {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			dm := models.DiskMetric{
				Device:     p.Device,
				Mountpoint: p.Mountpoint,
				Total:      usage.Total,
				Used:       usage.Used,
				Free:       usage.Free,
				Percent:    usage.UsedPercent,
			}
			devName := deviceName(p.Device)
			if curr, ok := diskIOMap[devName]; ok {
				if prev, ok := mc.lastDiskStats[devName]; ok {
					dm.ReadBytes = uint64(float64(curr.ReadBytes-prev.ReadBytes) / elapsed)
					dm.WriteBytes = uint64(float64(curr.WriteBytes-prev.WriteBytes) / elapsed)
				}
				mc.lastDiskStats[devName] = curr
			}
			metrics.Disks = append(metrics.Disks, dm)
		}
	}

	// Network
	netStats, err := net.IOCounters(true)
	if err == nil {
		for _, ns := range netStats {
			nm := models.NetworkMetric{Interface: ns.Name}
			if prev, ok := mc.lastNetStats[ns.Name]; ok {
				nm.BytesSent = uint64(float64(ns.BytesSent-prev.BytesSent) / elapsed)
				nm.BytesRecv = uint64(float64(ns.BytesRecv-prev.BytesRecv) / elapsed)
				nm.PacketsSent = uint64(float64(ns.PacketsSent-prev.PacketsSent) / elapsed)
				nm.PacketsRecv = uint64(float64(ns.PacketsRecv-prev.PacketsRecv) / elapsed)
			}
			mc.lastNetStats[ns.Name] = ns
			metrics.Networks = append(metrics.Networks, nm)
		}
	}

	// GPU (nvidia-smi or fallback)
	metrics.GPU = collectGPU()

	// Host info
	hostInfo, err := host.Info()
	if err == nil {
		metrics.Uptime = hostInfo.Uptime
		metrics.HostInfo = models.HostInfo{
			Hostname:      hostInfo.Hostname,
			OS:            hostInfo.OS,
			Platform:      hostInfo.Platform,
			Arch:          runtime.GOARCH,
			KernelVersion: hostInfo.KernelVersion,
		}
	}

	mc.lastCollect = now
	return metrics, nil
}

func collectGPU() []models.GPUMetric {
	var gpus []models.GPUMetric

	out, err := exec.Command("nvidia-smi",
		"--query-gpu=name,utilization.gpu,memory.used,memory.total,temperature.gpu",
		"--format=csv,noheader,nounits").Output()
	if err != nil {
		return gpus
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ", ")
		if len(parts) < 5 {
			continue
		}
		util, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		memUsed, _ := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
		memTotal, _ := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 64)
		temp, _ := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
		gpus = append(gpus, models.GPUMetric{
			Name:        strings.TrimSpace(parts[0]),
			Utilization: util,
			MemoryUsed:  memUsed * 1024 * 1024,
			MemoryTotal: memTotal * 1024 * 1024,
			Temperature: temp,
		})
	}
	return gpus
}

func GetProcesses() ([]models.ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var result []models.ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpuPct, err := p.CPUPercent()
		if err != nil {
			continue
		}
		memInfo, err := p.MemoryInfo()
		if err != nil {
			continue
		}
		memPct, _ := p.MemoryPercent()
		status, _ := p.Status()
		username, _ := p.Username()

		statusStr := ""
		if len(status) > 0 {
			statusStr = status[0]
		}

		result = append(result, models.ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPct,
			MemoryMB:   float64(memInfo.RSS) / 1024 / 1024,
			MemPercent: memPct,
			Status:     statusStr,
			Username:   username,
		})
	}
	return result, nil
}

func deviceName(device string) string {
	if runtime.GOOS == "windows" {
		parts := strings.Split(device, "\\")
		return parts[len(parts)-1]
	}
	parts := strings.Split(device, "/")
	name := parts[len(parts)-1]
	// strip partition number for Linux
	name = strings.TrimRight(name, "0123456789")
	return name
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
