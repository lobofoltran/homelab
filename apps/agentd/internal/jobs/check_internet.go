package jobs

import (
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

type PingResult struct {
	Host    string `json:"host"`
	Reach   bool   `json:"reachable"`
	Latency string `json:"latency,omitempty"`
}

func CheckInternet(central config.CentralConfig) {
	logger.Log.Info("[check_internet] Verificando conectividade...")

	targets := []string{
		"8.8.8.8",         // Google
		"1.1.1.1",         // Cloudflare
		"chat.openai.com", // OpenAI
	}

	// Adiciona host da central, se fornecido
	if central.URL != "" {
		target := extractDomain(central.URL)
		if target != "" {
			targets = append(targets, target)
		}
	}

	var results []PingResult
	for _, host := range targets {
		result := pingHostInternet(host)
		results = append(results, result)
	}

	report := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"os":        runtime.GOOS,
		"results":   results,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("internet", report)
		logger.Log.Info("[check_internet] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_internet] (prod) Enviaria para %s", central.URL)
	}
}

func pingHostInternet(host string) PingResult {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", host)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "1", host)
	}
	output, err := cmd.Output()
	if err != nil {
		return PingResult{Host: host, Reach: false}
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			lat := extractLatency(line)
			return PingResult{Host: host, Reach: true, Latency: lat}
		}
	}
	return PingResult{Host: host, Reach: false}
}

func extractLatency(line string) string {
	start := strings.Index(line, "time=")
	if start == -1 {
		return ""
	}
	lat := line[start+5:]
	end := strings.Index(lat, " ")
	if end != -1 {
		lat = lat[:end]
	}
	return lat
}

func extractDomain(url string) string {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	if strings.Contains(url, "/") {
		url = strings.Split(url, "/")[0]
	}
	return url
}
