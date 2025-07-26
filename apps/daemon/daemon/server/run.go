package server

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/lobofoltran/homelab/libs/logger"
)

const socketPath = "/tmp/homelab.sock"

func Run(ctx context.Context) error {
	if err := os.RemoveAll(socketPath); err != nil {
		return err
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	defer listener.Close()

	if err := os.Chmod(socketPath, 0666); err != nil {
		return err
	}

	server := &http.Server{
		Handler: NewRouter(),
	}

	go func() {
		logger.Info("Daemon escutando no socket: %s", socketPath)
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Error("Erro ao servir: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Info("Encerrando daemon...")
	return server.Shutdown(context.Background())
}
