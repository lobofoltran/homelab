package jobs

import (
	"bufio"
	"bytes"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

type OpenPort struct {
	Protocol string `json:"protocol"`
	LocalIP  string `json:"local_ip"`
	Port     string `json:"port"`
	PID      string `json:"pid"`
	Process  string `json:"process,omitempty"`
}

func CheckPorts(central config.CentralConfig) {
	logger.Log.Info("[check_ports] Mapeando portas abertas...")

	var ports []OpenPort

	if runtime.GOOS == "windows" {
		ports = getPortsWindows()
	} else {
		ports = getPortsLinux()
	}

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"ports":     ports,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("ports", result)
		logger.Log.Info("[check_ports] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_ports] (prod) Enviaria para %s", central.URL)
	}
}

func getPortsLinux() []OpenPort {
	cmd := exec.Command("ss", "-tulnp")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	scanner := bufio.NewScanner(&out)
	var ports []OpenPort
	skipFirst := true

	for scanner.Scan() {
		line := scanner.Text()
		if skipFirst {
			skipFirst = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		proto := fields[0]
		localAddr := fields[4]
		pidName := fields[len(fields)-1]

		ip, port := parseAddress(localAddr)
		pid, proc := parsePIDProcess(pidName)

		ports = append(ports, OpenPort{
			Protocol: proto,
			LocalIP:  ip,
			Port:     port,
			PID:      pid,
			Process:  proc,
		})
	}
	return ports
}

func getPortsWindows() []OpenPort {
	cmd := exec.Command("netstat", "-ano")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	scanner := bufio.NewScanner(&out)
	var ports []OpenPort

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "Proto") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		proto := strings.ToLower(fields[0])
		local := fields[1]
		pid := fields[len(fields)-1]

		ip, port := parseAddress(local)
		proc := getProcessNameByPID(pid)

		ports = append(ports, OpenPort{
			Protocol: proto,
			LocalIP:  ip,
			Port:     port,
			PID:      pid,
			Process:  proc,
		})
	}
	return ports
}

func parseAddress(addr string) (string, string) {
	if strings.Contains(addr, "[") { // IPv6
		addr = strings.Trim(addr, "[]")
	}
	split := strings.LastIndex(addr, ":")
	if split == -1 {
		return addr, ""
	}
	return addr[:split], addr[split+1:]
}

func parsePIDProcess(raw string) (string, string) {
	if strings.HasPrefix(raw, "users:(") {
		raw = strings.TrimPrefix(raw, "users:(")
		raw = strings.TrimSuffix(raw, ")")
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			if strings.Contains(part, "pid=") {
				pid := strings.Split(part, "=")[1]
				name := strings.Split(part, "=")[0]
				return pid, name
			}
		}
	}
	return "-", "-"
}

func getProcessNameByPID(pid string) string {
	p, err := strconv.Atoi(pid)
	if err != nil || p == 0 {
		return "-"
	}
	out, err := exec.Command("tasklist", "/FI", "PID eq "+pid).Output()
	if err != nil {
		return "-"
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 3 {
		fields := strings.Fields(lines[3])
		if len(fields) > 0 {
			return fields[0]
		}
	}
	return "-"
}
