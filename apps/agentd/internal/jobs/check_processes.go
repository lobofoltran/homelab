package jobs

import (
	"sort"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemoryMB   float32 `json:"memory_mb"`
	User       string  `json:"user"`
}

func CheckProcesses(central config.CentralConfig) {
	logger.Log.Info("[check_processes] Coletando informações dos processos...")

	procs, err := process.Processes()
	if err != nil {
		logger.Log.Errorf("[check_processes] Erro ao listar processos: %v", err)
		return
	}

	var top []ProcessInfo

	for _, p := range procs {
		cpuPercent, err1 := p.CPUPercent()
		memInfo, err2 := p.MemoryInfo()
		name, err3 := p.Name()
		user, err4 := p.Username()

		if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
			top = append(top, ProcessInfo{
				PID:        p.Pid,
				Name:       name,
				CPUPercent: cpuPercent,
				MemoryMB:   float32(memInfo.RSS) / 1024.0 / 1024.0,
				User:       user,
			})
		}
	}

	sort.Slice(top, func(i, j int) bool {
		return top[i].CPUPercent > top[j].CPUPercent
	})

	if len(top) > 10 {
		top = top[:10]
	}

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"top":       top,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("processes", result)
		logger.Log.Info("[check_processes] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_processes] (prod) Enviaria para %s", central.URL)
	}
}
