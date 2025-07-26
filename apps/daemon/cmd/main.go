package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/lobofoltran/homelab/apps/daemon/daemon/server"
	"github.com/lobofoltran/homelab/libs/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger.Info("Iniciando homelabd...")
	logger.SetDebug(true)

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Erro ao executar daemon: %v", err)
	}
}
