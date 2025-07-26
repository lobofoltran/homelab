package jobs

import (
	"fmt"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func CheckSystem(central config.CentralConfig) {
	logger.Log.Info("[check_system] Coletando informações do sistema...")

	v, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	d, _ := disk.Usage("/")

	var tempInfo string
	if sensors, err := host.SensorsTemperatures(); err == nil && len(sensors) > 0 {
		for _, sensor := range sensors {
			if sensor.Temperature > 0 {
				tempInfo = sensor.SensorKey + ": " + formatTemp(sensor.Temperature)
				break
			}
		}
	}

	data := map[string]interface{}{
		"time":         time.Now().Format("2006-01-02 15:04:05"),
		"ram_used":     formatBytes(v.Used),
		"ram_total":    formatBytes(v.Total),
		"ram_percent":  v.UsedPercent,
		"cpu_percent":  cpuPercent[0],
		"disk_used":    formatBytes(d.Used),
		"disk_total":   formatBytes(d.Total),
		"disk_percent": d.UsedPercent,
		"temp":         tempInfo,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("system", data)
		logger.Log.Info("[check_system] (dev) Resultado salvo localmente em /result")
	} else {
		// TODO: enviar via HTTP para central.URL
		logger.Log.Infof("[check_system] (prod) Enviaria para %s", central.URL)
	}
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatTemp(temp float64) string {
	return fmt.Sprintf("%.1f°C", temp)
}
