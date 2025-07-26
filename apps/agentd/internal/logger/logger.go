package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

var currentLogDate string
var logFile *os.File

func InitLogger() {
	Log = logrus.New()

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	logDir := filepath.Join(exeDir, "logs")
	_ = os.MkdirAll(logDir, os.ModePerm)

	currentLogDate = time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, currentLogDate+".log")

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Erro ao criar arquivo de log: " + err.Error())
	}
	logFile = file

	isService := os.Getenv("SESSIONNAME") == "" // serviço não tem sessão visível

	if isService {
		// Somente arquivo quando rodando como serviço
		Log.SetOutput(logFile)
	} else {
		// Loga no terminal + arquivo quando rodando em modo CLI
		Log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     !isService,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	Log.SetLevel(logrus.InfoLevel)
}
