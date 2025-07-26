package jobs

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
	"github.com/lobofoltran/homelab/apps/agentd/internal/utils"
)

type FileInfo struct {
	Path       string `json:"path"`
	Exists     bool   `json:"exists"`
	SizeBytes  int64  `json:"size_bytes"`
	ModifiedAt string `json:"modified_at,omitempty"`
	IsDir      bool   `json:"is_dir"`
}

func CheckFiles(central config.CentralConfig) {
	logger.Log.Info("[check_files] Verificando arquivos e pastas sensíveis...")

	var paths []string

	if runtime.GOOS == "windows" {
		paths = []string{
			`C:\Windows\System32\config\SAM`,
			`C:\Windows\System32\drivers\etc\hosts`,
			`C:\Windows\Temp`,
			`C:\Windows\Logs`,
		}
	} else {
		paths = []string{
			"/etc/passwd",
			"/etc/hosts",
			"/var/log",
			"/tmp",
		}
	}

	var results []FileInfo
	for _, path := range paths {
		info := checkPath(path)
		results = append(results, info)
	}

	result := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"os":        runtime.GOOS,
		"files":     results,
	}

	if !config.IsProduction() {
		utils.SaveResultJSON("files", result)
		logger.Log.Info("[check_files] (dev) Resultado salvo localmente em /result")
	} else {
		logger.Log.Infof("[check_files] (prod) Enviaria para %s", central.URL)
	}
}

func checkPath(path string) FileInfo {
	fi := FileInfo{
		Path:   path,
		Exists: false,
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fi
	}

	fi.Exists = true
	fi.IsDir = stat.IsDir()
	fi.SizeBytes = getTotalSize(path, fi.IsDir)

	if !fi.IsDir {
		fi.ModifiedAt = stat.ModTime().Format("2006-01-02 15:04:05")
	}
	return fi
}

func getTotalSize(path string, isDir bool) int64 {
	if !isDir {
		stat, err := os.Stat(path)
		if err == nil {
			return stat.Size()
		}
		return 0
	}

	var size int64 = 0
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}
