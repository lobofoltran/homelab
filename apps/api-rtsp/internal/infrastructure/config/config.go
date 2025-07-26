package config

import (
	"encoding/json"
	"os"
)

type CameraConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AppConfig struct {
	ServerPort int            `json:"server_port"`
	Cameras    []CameraConfig `json:"cameras"`
}

func LoadConfig(path string) (*AppConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
