package monitor

import (
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemStats holds the real-time system statistics
type SystemStats struct {
	CpuUsage  float64 `json:"cpu_usage"`
	RamUsage  float64 `json:"ram_usage"`
	DiskUsage float64 `json:"disk_usage"`
	Timestamp int64   `json:"timestamp"`
}

// StatsChannel is where we send the gathered stats
var StatsChannel = make(chan SystemStats)

// StartMonitoring begins the stats collection loop
func StartMonitoring() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := gatherStats()
		// Non-blocking send to avoid hanging if no one is listening
		select {
		case StatsChannel <- stats:
		default:
		}
	}
}

func gatherStats() SystemStats {
	// CPU
	// 0 means calculate for all cores combined, false means don't separate per core
	percentages, err := cpu.Percent(0, false)
	cpuVal := 0.0
	if err == nil && len(percentages) > 0 {
		cpuVal = percentages[0]
	}

	// RAM
	v, err := mem.VirtualMemory()
	ramVal := 0.0
	if err == nil {
		ramVal = v.UsedPercent
	}

	// Disk (root)
	d, err := disk.Usage("/")
	diskVal := 0.0
	if err == nil {
		diskVal = d.UsedPercent
	}

	return SystemStats{
		CpuUsage:  cpuVal,
		RamUsage:  ramVal,
		DiskUsage: diskVal,
		Timestamp: time.Now().UnixMilli(),
	}
}
