package config

import (
	"encoding/json"
	"os"
)

type JobConfig struct {
	Job      string `json:"job"`
	Interval int    `json:"interval"` // minutes
}

type CentralConfig struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type DatabaseConfig struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServiceEntry struct {
	ServiceID string `json:"serviceId"`
}

type AppEntry struct {
	Type    string `json:"type"`
	Project string `json:"project"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Service string `json:"service"`
}

type SMTPConfig struct {
	Enabled  bool     `json:"enabled"`
	SMTPHost string   `json:"smtpHost"`
	SMTPPort string   `json:"smtpPort"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

type Config struct {
	AppID      string                  `json:"appId"`
	Production bool                    `json:"production"`
	Jobs       []JobConfig             `json:"jobs"`
	Central    CentralConfig           `json:"central"`
	Databases  []DatabaseConfig        `json:"databases"`
	SMTP       SMTPConfig              `json:"smtp"`
	Services   map[string]ServiceEntry `json:"services"`
	UpdateApps []AppEntry              `json:"updateApps"`
}

var Current Config

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Current)
	return Current, err
}

func IsProduction() bool {
	return Current.Production
}
