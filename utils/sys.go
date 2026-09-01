package utils

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var BootTime = time.Now()

// BotSysStats returns Uptime, CPU, RAM, and Disk usage
func BotSysStats() (string, string, string, string) {
	uptime := GetReadableTime(int(time.Since(BootTime).Seconds()))

	cpuPercent, _ := cpu.Percent(0, false)
	cpuStr := "0%"
	if len(cpuPercent) > 0 {
		cpuStr = fmt.Sprintf("%.1f%%", cpuPercent[0])
	}

	vm, _ := mem.VirtualMemory()
	ramStr := fmt.Sprintf("%.1f%%", vm.UsedPercent)

	diskStat, _ := disk.Usage("/")
	diskStr := fmt.Sprintf("%.1f%%", diskStat.UsedPercent)

	return uptime, cpuStr, ramStr, diskStr
}
