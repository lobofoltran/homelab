package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Level string

const (
	INFO  Level = "INFO"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	DEBUG Level = "DEBUG"
)

var debugMode = false

func SetDebug(enabled bool) {
	debugMode = enabled
}

func logMessage(level Level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf(msg, args...)

	prefix := fmt.Sprintf("[%s] %s", level, timestamp)

	switch level {
	case ERROR:
		log.Printf("\033[31m%s\033[0m %s", prefix, formatted)
	case WARN:
		log.Printf("\033[33m%s\033[0m %s", prefix, formatted)
	case INFO:
		log.Printf("\033[34m%s\033[0m %s", prefix, formatted)
	case DEBUG:
		if debugMode {
			log.Printf("\033[36m%s\033[0m %s", prefix, formatted)
		}
	}
}

func Info(msg string, args ...any)  { logMessage(INFO, msg, args...) }
func Warn(msg string, args ...any)  { logMessage(WARN, msg, args...) }
func Error(msg string, args ...any) { logMessage(ERROR, msg, args...) }
func Debug(msg string, args ...any) { logMessage(DEBUG, msg, args...) }

func Fatal(msg string, args ...any) {
	logMessage(ERROR, msg, args...)
	os.Exit(1)
}
