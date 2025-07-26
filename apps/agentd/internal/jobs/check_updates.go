package jobs

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
)

type UpdateRequest struct {
	Hostname string            `json:"hostname"`
	Apps     []config.AppEntry `json:"apps"`
}

type UpdateResponse struct {
	Project   string `json:"project"`
	Type      string `json:"type"`
	Path      string `json:"target"`
	Download  string `json:"download"`
	Version   string `json:"version"`
	ServiceID string `json:"service"`
}

func CheckUpdates(central config.CentralConfig) {
	hostname, _ := os.Hostname()
	apps := []config.AppEntry{}

	for _, app := range config.Current.UpdateApps {
		ver := fileHash(app.Path)
		apps = append(apps, config.AppEntry{
			Type:    app.Type,
			Project: app.Project,
			Path:    app.Path,
			Version: ver,
			Service: app.Service,
		})
	}

	reqBody := UpdateRequest{
		Hostname: hostname,
		Apps:     apps,
	}

	jsonBody, _ := json.Marshal(reqBody)
	resp, err := http.Post(central.URL+"/check-updates", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		logger.Log.Errorf("[check_updates] Falha na comunicação com a central: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		logger.Log.Info("[check_updates] Todos os aplicativos estão atualizados.")
		return
	}

	if resp.StatusCode == 200 {
		var updates []UpdateResponse
		_ = json.NewDecoder(resp.Body).Decode(&updates)
		for _, upd := range updates {
			executeUpdate(upd)
		}
	} else {
		logger.Log.Warnf("[check_updates] Código HTTP inesperado: %d", resp.StatusCode)
	}
}

func executeUpdate(upd UpdateResponse) {
	logger.Log.Infof("[check_updates] [%s] Iniciando atualização (%s)...", upd.Project, upd.Type)

	resp := downloadNewVersion(upd)
	if resp == nil {
		logger.Log.Warnf("[check_updates] [%s] Atualização abortada — download inválido.", upd.Project)
		return
	}
	defer resp.Body.Close()

	backupDir := createBackup(upd)

	if !saveNewFile(upd.Path+".new", resp.Body, upd.Project) {
		return
	}

	payload := map[string]string{
		"project":   upd.Project,
		"type":      upd.Type,
		"path":      upd.Path,
		"service":   upd.ServiceID,
		"backupDir": backupDir,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		logger.Log.Errorf("[check_updates] [%s] Erro ao gerar payload de update: %v", upd.Project, err)
		return
	}

	respExec, err := http.Post("http://127.0.0.1:9898/execute-update", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		logger.Log.Errorf("[check_updates] [%s] Erro ao comunicar com executor: %v", upd.Project, err)
		return
	}
	defer respExec.Body.Close()

	if respExec.StatusCode != http.StatusOK {
		logger.Log.Warnf("[check_updates] [%s] Executor respondeu com status: %d", upd.Project, respExec.StatusCode)
		return
	}

	logger.Log.Infof("[check_updates] [%s] Executor acionado com sucesso.", upd.Project)
}

func createBackup(upd UpdateResponse) string {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupDir := filepath.Join("backup", upd.Project, upd.Version+"_"+timestamp)
	_ = os.MkdirAll(backupDir, os.ModePerm)

	originalFile, err := os.Open(upd.Path)
	if err == nil {
		defer originalFile.Close()
		backupPath := filepath.Join(backupDir, filepath.Base(upd.Path))
		backupFile, err := os.Create(backupPath)
		if err == nil {
			defer backupFile.Close()
			io.Copy(backupFile, originalFile)
			logger.Log.Infof("[check_updates] [%s] Backup criado em: %s", upd.Project, backupPath)
		}
	}
	return backupDir
}

func downloadNewVersion(upd UpdateResponse) *http.Response {
	logger.Log.Infof("[check_updates] [%s] Baixando nova versão de: %s", upd.Project, upd.Download)
	resp, err := http.Get(upd.Download)
	if err != nil {
		logger.Log.Errorf("[check_updates] [%s] Erro ao baixar: %v", upd.Project, err)
		return nil
	}
	if resp.StatusCode != 200 {
		logger.Log.Errorf("[check_updates] [%s] Download inválido! HTTP %d ao acessar %s", upd.Project, resp.StatusCode, upd.Download)
		return nil
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "text/html" || contentType == "application/json" {
		logger.Log.Errorf("[check_updates] [%s] Tipo de conteúdo inesperado: %s", upd.Project, contentType)
		return nil
	}
	logger.Log.Infof("[check_updates] [%s] Nova versão validada com sucesso", upd.Project)
	return resp
}

func saveNewFile(path string, body io.Reader, project string) bool {
	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		logger.Log.Errorf("[check_updates] [%s] Erro ao criar novo arquivo: %v", project, err)
		return false
	}

	_, err = io.Copy(out, body)
	if err != nil {
		out.Close()
		logger.Log.Errorf("[check_updates] [%s] Erro ao salvar novo arquivo: %v", project, err)
		return false
	}

	// if err := out.Sync(); err != nil {
	// 	out.Close()
	// 	logger.Log.Errorf("[check_updates] [%s] Erro ao sincronizar novo arquivo: %v", project, err)
	// 	return false
	// }

	out.Close()
	logger.Log.Infof("[check_updates] [%s] Arquivo salvo com sucesso: %s", project, path)
	return true
}

func fileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return "erro"
	}
	defer f.Close()
	h := sha256.New()
	_, _ = io.Copy(h, f)
	return fmt.Sprintf("%x", h.Sum(nil))[:12] // versão baseada em hash resumido
}
