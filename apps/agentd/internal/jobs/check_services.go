package jobs

import (
	"bufio"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

func CheckServices(central config.CentralConfig) {
	logger.Log.Info("[check_services] Coletando todos os serviços e seus status...")

	var services []ServiceStatus

	if runtime.GOOS == "windows" {
		services = getAllWindowsServices()
	} else {
		services = getAllLinuxServices()
	}

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"os":        runtime.GOOS,
		"services":  services,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("services", result)
		logger.Log.Info("[check_services] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_services] (prod) Enviaria para %s", central.URL)
	}
}

func getAllLinuxServices() []ServiceStatus {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--no-legend")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	var services []ServiceStatus
	scanner := bufio.NewScanner(&out)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			name := fields[0]
			status := fields[2] // 'active', 'inactive', etc.
			services = append(services, ServiceStatus{
				Name:   name,
				Status: status,
			})
		}
	}
	return services
}

func getAllWindowsServices() []ServiceStatus {
	cmd := exec.Command("sc", "query", "state=", "all")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	var services []ServiceStatus
	lines := strings.Split(out.String(), "\n")
	var currentName string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SERVICE_NAME:") {
			currentName = strings.TrimSpace(strings.TrimPrefix(line, "SERVICE_NAME:"))
		}
		if strings.HasPrefix(line, "STATE") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				status := strings.ToLower(parts[3]) // RUNNING, STOPPED, etc.
				services = append(services, ServiceStatus{
					Name:   currentName,
					Status: status,
				})
			}
		}
	}
	return services
}
