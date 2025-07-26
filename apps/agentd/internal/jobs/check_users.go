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

type LoggedUser struct {
	User      string `json:"user"`
	Terminal  string `json:"terminal,omitempty"`
	Host      string `json:"host,omitempty"`
	LoginTime string `json:"login_time,omitempty"`
}

func CheckUsers(central config.CentralConfig) {
	logger.Log.Info("[check_users] Coletando sessões ativas...")

	var users []LoggedUser

	if runtime.GOOS == "windows" {
		users = parseWindowsQueryUser()
	} else {
		users = parseLinuxWho()
	}

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"users":     users,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("users", result)
		logger.Log.Info("[check_users] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_users] (prod) Enviaria para %s", central.URL)
	}
}

func parseLinuxWho() []LoggedUser {
	out, err := exec.Command("who").Output()
	if err != nil {
		logger.Log.Errorf("[check_users] Erro ao executar who: %v", err)
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var users []LoggedUser

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		entry := LoggedUser{
			User:      fields[0],
			Terminal:  fields[1],
			LoginTime: fields[2] + " " + fields[3],
		}
		if len(fields) >= 5 && strings.HasPrefix(fields[4], "(") {
			entry.Host = strings.Trim(fields[4], "()")
		}
		users = append(users, entry)
	}
	return users
}

func parseWindowsQueryUser() []LoggedUser {
	out, err := exec.Command("query", "user").Output()
	if err != nil {
		logger.Log.Errorf("[check_users] Erro ao executar query user: %v", err)
		return nil
	}
	lines := strings.Split(string(out), "\n")
	var users []LoggedUser

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "USERNAME") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 {
			users = append(users, LoggedUser{
				User:      fields[0],
				Terminal:  fields[1],
				LoginTime: strings.Join(fields[len(fields)-2:], " "), // aproximação do logon time
			})
		}
	}
	return users
}
