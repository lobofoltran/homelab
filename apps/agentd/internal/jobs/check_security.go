package jobs

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

func CheckSecurity(central config.CentralConfig) {
	logger.Log.Info("[check_security] Iniciando verificação de segurança...")

	var firewallStatus string
	var antivirusStatus string
	var updateStatus string
	var rootOrAdmin string
	var sensitiveFiles []string

	if runtime.GOOS == "windows" {
		firewallStatus = runAndParse("netsh", "advfirewall", "show", "allprofiles")
		antivirusStatus = runAndParse("powershell", "Get-CimInstance -Namespace root/SecurityCenter2 -ClassName AntivirusProduct")
		rootOrAdmin = detectAdminWindows()
		sensitiveFiles = []string{"C:\\Windows\\System32\\config\\SAM"}
	} else {
		firewallStatus = tryFirstValid(
			runAndParse("ufw", "status"),
			runAndParse("iptables", "-L"))
		antivirusStatus = detectClamAV()
		updateStatus = runAndParse("apt", "list", "--upgradable")
		if updateStatus == "" {
			updateStatus = runAndParse("dnf", "check-update")
		}
		rootOrAdmin = detectRootLogin()
		sensitiveFiles = []string{"/etc/shadow", "/etc/sudoers"}
	}

	permStatus := map[string]string{}
	for _, f := range sensitiveFiles {
		permStatus[f] = checkFilePermissions(f)
	}

	result := map[string]interface{}{
		"timestamp":        time.Now().Format("2006-01-02 15:04:05"),
		"os":               runtime.GOOS,
		"firewall_status":  firewallStatus,
		"antivirus":        antivirusStatus,
		"updates":          updateStatus,
		"root_or_admin":    rootOrAdmin,
		"file_permissions": permStatus,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("security", result)
		logger.Log.Info("[check_security] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_security] (prod) Enviaria para %s", central.URL)
	}
}

func runAndParse(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}

func tryFirstValid(cmds ...string) string {
	for _, out := range cmds {
		if out != "" {
			return out
		}
	}
	return "N/A"
}

func checkFilePermissions(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "inacessível"
	}
	mode := info.Mode()
	return mode.String()
}

func detectClamAV() string {
	if _, err := exec.LookPath("clamscan"); err == nil {
		return "ClamAV instalado"
	}
	return "Nenhum antivírus detectado"
}

func detectRootLogin() string {
	out, err := exec.Command("who").Output()
	if err != nil {
		return "erro ao verificar sessões"
	}
	if strings.Contains(string(out), "root") {
		return "root logado"
	}
	return "sem root logado"
}

func detectAdminWindows() string {
	out, err := exec.Command("net", "user").Output()
	if err != nil {
		return "erro ao verificar usuários"
	}
	if strings.Contains(string(out), "Administrator") {
		return "admin presente"
	}
	return "sem admin logado"
}
