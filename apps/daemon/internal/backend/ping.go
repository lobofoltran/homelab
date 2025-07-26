package backend

import "github.com/lobofoltran/homelab/libs/logger"

func Ping() string {
	logger.Error("Chefe")
	return "pong"
}
