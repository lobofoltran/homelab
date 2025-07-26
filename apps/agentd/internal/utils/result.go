package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
)

func SaveResultJSON(name string, data any) {
	resultDir := "result"
	if _, err := os.Stat(resultDir); os.IsNotExist(err) {
		_ = os.MkdirAll(resultDir, os.ModePerm)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := name + "_" + timestamp + ".json"
	filePath := filepath.Join(resultDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		logger.Log.Errorf("Erro ao criar arquivo de resultado: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		logger.Log.Errorf("Erro ao salvar resultado JSON: %v", err)
	}
}
