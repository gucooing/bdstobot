package state

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"time"
)

func GetCpuPercent() string {
	percent, _ := cpu.Percent(time.Second, false)
	msg := fmt.Sprintf("%.2f", percent[0])
	return msg
}

func GetMemPercents() string {
	memInfo, _ := mem.VirtualMemory()
	msg := fmt.Sprintf("%.2f", memInfo.UsedPercent)
	return msg
}

func GetMemPercent() string {
	memInfo, _ := mem.VirtualMemory()
	total := memInfo.Total / 1024 / 1024 // 将总内存量转换为MB
	used := memInfo.Used / 1024 / 1024   // 将已使用内存量转换为MB
	fmt.Printf("Total Memory: %d MB\n", total)
	fmt.Printf("Used Memory: %d MB\n", used)
	msg := strconv.FormatUint(used, 10) + "MB /" + strconv.FormatUint(total, 10) + "MB"
	return msg
}

func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}
