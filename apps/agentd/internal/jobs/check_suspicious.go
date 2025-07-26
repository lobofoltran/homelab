package jobs

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

func CheckSuspicious(central config.CentralConfig) {
	logger.Log.Info("[check_suspicious] Procurando sinais de anomalia...")

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"os":        runtime.GOOS,
	}

	if runtime.GOOS == "linux" {
		result["users_with_uid_0"] = findUID0Users()
		result["users_without_shell"] = findUsersWithoutShell()
		result["executables_in_tmp"] = findExecutablesInTmp()
		result["processes_from_tmp"] = findProcessesInTmp()
	} else {
		result["info"] = "Verificação limitada no Windows (não aplicável)"
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("suspicious", result)
		logger.Log.Info("[check_suspicious] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_suspicious] (prod) Enviaria para %s", central.URL)
	}
}

func findUID0Users() []string {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil
	}
	defer file.Close()

	var users []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}
		if parts[2] == "0" && parts[0] != "root" {
			users = append(users, parts[0])
		}
	}
	return users
}

func findUsersWithoutShell() []string {
	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil
	}
	defer file.Close()

	var users []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}
		shell := parts[6]
		if shell == "/usr/sbin/nologin" || shell == "/bin/false" || shell == "" {
			users = append(users, parts[0])
		}
	}
	return users
}

func findExecutablesInTmp() []string {
	cmd := exec.Command("find", "/tmp", "/var/tmp", "-type", "f", "-perm", "/111", "-exec", "ls", "-la", "{}", ";")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var paths []string
	for _, line := range lines {
		if strings.Contains(line, "/tmp") || strings.Contains(line, "/var/tmp") {
			fields := strings.Fields(line)
			if len(fields) > 8 {
				paths = append(paths, fields[len(fields)-1])
			}
		}
	}
	return paths
}

func findProcessesInTmp() []string {
	procs, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}
	var matches []string
	for _, entry := range procs {
		if !entry.IsDir() {
			continue
		}
		pid := entry.Name()
		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}
		exePath := filepath.Join("/proc", pid, "exe")
		target, err := os.Readlink(exePath)
		if err == nil && (strings.HasPrefix(target, "/tmp") || strings.HasPrefix(target, "/var/tmp")) {
			matches = append(matches, target)
		}
	}
	return matches
}
