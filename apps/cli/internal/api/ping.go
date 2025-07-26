package api

import (
	"context"
	"io"
	"net"
	"net/http"
)

const socketPath = "/tmp/homelab.sock"

func Ping() (string, error) {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
		},
	}
	client := &http.Client{Transport: transport}

	// Host é ignorado, mas precisa ser um endereço válido
	resp, err := client.Get("http://unix/ping")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
