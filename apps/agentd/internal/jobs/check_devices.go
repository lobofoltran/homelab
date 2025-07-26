package jobs

import (
	"bytes"
	"runtime"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

func CheckDevices(central config.CentralConfig) {
	logger.Log.Infof("[check_devices] Iniciando coleta de dispositivos (%s)...", runtime.GOOS)

	var output string
	var err error

	switch runtime.GOOS {
	case "linux":
		output, err = utils.RunCommand("lspci")
	case "windows":
		output, err = utils.RunCommand("wmic", "path", "Win32_PnPSignedDriver", "get", "DeviceName,DriverVersion")
	default:
		logger.Log.Warn("Sistema operacional não suportado para check_devices.")
		return
	}

	if err != nil {
		logger.Log.Errorf("[check_devices] Erro ao executar comando: %v", err)
		return
	}

	result := map[string]interface{}{
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		"os":         runtime.GOOS,
		"lines":      len(bytes.Split([]byte(output), []byte("\n"))),
		"raw_output": output,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("devices", result)
		logger.Log.Info("[check_devices] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_devices] (prod) Enviaria para %s", central.URL)
	}
}
