package jobs

import (
	"bufio"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

func CheckNetwork(central config.CentralConfig) {
	logger.Log.Info("[check_network] Coletando informações de rede...")

	interfaces, _ := net.Interfaces()
	var activeIfaces []map[string]interface{}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		var ips []string
		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}
		activeIfaces = append(activeIfaces, map[string]interface{}{
			"name": iface.Name,
			"mac":  iface.HardwareAddr.String(),
			"ips":  ips,
		})
	}

	// Ping 8.8.8.8
	pingLatency := pingHost("8.8.8.8")

	// Gateway
	gateway := getGateway()

	// DNS
	dns := getDNS()

	result := map[string]interface{}{
		"timestamp":    time.Now().Format("2006-01-02 15:04:05"),
		"interfaces":   activeIfaces,
		"gateway":      gateway,
		"dns":          dns,
		"ping_8_8_8_8": pingLatency,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("network", result)
		logger.Log.Info("[check_network] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_network] (prod) Enviaria para %s", central.URL)
	}
}

func pingHost(host string) string {
	out, err := exec.Command("ping", "-c", "1", "-W", "1", host).Output()
	if err != nil {
		return "fail"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "time=") {
			return strings.TrimSpace(line)
		}
	}
	return "timeout"
}

func getGateway() string {
	out, err := exec.Command("ip", "route").Output()
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "default via") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				return fields[2]
			}
		}
	}
	return "not found"
}

func getDNS() []string {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return []string{}
	}
	defer file.Close()

	var dns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "nameserver") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				dns = append(dns, parts[1])
			}
		}
	}
	return dns
}
