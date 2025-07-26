package jobs

import (
	"runtime"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
)

func CheckHostInfo(central config.CentralConfig) {
	logger.Log.Info("[check_hostinfo] Coletando informações do host...")

	info, err := host.Info()
	if err != nil {
		logger.Log.Errorf("[check_hostinfo] Erro ao obter host info: %v", err)
		return
	}

	cpuInfo, _ := cpu.Info()
	cpuCount, _ := cpu.Counts(false)
	logicalCount, _ := cpu.Counts(true)
	loc := time.Now().Location()

	result := map[string]interface{}{
		"timestamp":     time.Now().Format("2006-01-02 15:04:05"),
		"hostname":      info.Hostname,
		"host_id":       info.HostID,
		"os":            info.OS,
		"platform":      info.Platform,
		"platform_ver":  info.PlatformVersion,
		"kernel_ver":    info.KernelVersion,
		"arch":          runtime.GOARCH,
		"uptime_secs":   info.Uptime,
		"uptime_human":  formatUptime(info.Uptime),
		"cpu_model":     cpuInfo[0].ModelName,
		"cpu_cores":     cpuCount,
		"logical_cores": logicalCount,
		"timezone":      loc.String(),
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("hostinfo", result)
		logger.Log.Info("[check_hostinfo] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_hostinfo] (prod) Enviaria para %s", central.URL)
	}
}

func formatUptime(uptime uint64) string {
	d := time.Duration(uptime) * time.Second
	return d.Truncate(time.Second).String()
}
