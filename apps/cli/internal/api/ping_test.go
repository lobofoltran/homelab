package api

import (
	"testing"
)

func TestPing_Integration(t *testing.T) {
	got, err := Ping()
	if err != nil {
		t.Fatalf("Erro ao chamar Ping(): %v", err)
	}
	if got != "pong" {
		t.Errorf("Esperado 'pong', mas retornou: %s", got)
	}
}
