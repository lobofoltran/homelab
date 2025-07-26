package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"github.com/lobofoltran/homelab/apps/agentd/internal/config"
	"github.com/lobofoltran/homelab/apps/agentd/internal/jobs"
	"github.com/lobofoltran/homelab/apps/agentd/internal/logger"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	logger.Log.Info("Iniciando agente...")

	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Log.Fatalf("Erro ao carregar configuração: %v", err)
		return
	}
	if err != nil {
		logger.Log.Fatalf("Erro ao carregar configuração: %v", err)
		return
	}

	logger.Log.Infof("Iniciando agente para AppID: %s", cfg.AppID)

	for _, job := range cfg.Jobs {
		job := job // evitar captura incorreta no loop
		go func() {
			logger.Log.Infof("Agendando job '%s' a cada %d minuto(s)", job.Job, job.Interval)
			ticker := time.NewTicker(time.Duration(job.Interval) * time.Minute)
			defer ticker.Stop()

			runJob(job, cfg.Central)

			for range ticker.C {
				runJob(job, cfg.Central)
			}
		}()
	}

	select {}
}

func (p *program) Stop(s service.Service) error {
	logger.Log.Info("agentd está parando...")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "agentdBeta2",
		DisplayName: "FT Agent Service",
		Description: "Agente de diagnóstico e atualização da Fiscaltech.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		err := service.Control(s, os.Args[1])
		if err != nil {
			log.Fatalf("Comando %s falhou: %v", os.Args[1], err)
		}
		return
	}

	logger.InitLogger()

	err = s.Run()
	if err != nil {
		logger.Log.Errorf("Erro ao executar serviço: %v", err)
	}
}

func runJob(job config.JobConfig, central config.CentralConfig) {
	switch job.Job {
	case "check_devices":
		jobs.CheckDevices(central)
	case "check_system":
		jobs.CheckSystem(central)
	case "check_hostinfo":
		jobs.CheckHostInfo(central)
	case "check_processes":
		jobs.CheckProcesses(central)
	case "check_network":
		jobs.CheckNetwork(central)
	case "check_users":
		jobs.CheckUsers(central)
	case "check_security":
		jobs.CheckSecurity(central)
	case "check_ports":
		jobs.CheckPorts(central)
	case "check_services":
		jobs.CheckServices(central)
	case "check_files":
		jobs.CheckFiles(central)
	case "check_suspicious":
		jobs.CheckSuspicious(central)
	case "check_internet":
		jobs.CheckInternet(central)
	case "check_databases":
		jobs.CheckDatabases(central)
	case "check_updates":
		jobs.CheckUpdates(central)
	default:
		logger.Log.Warnf("Job não reconhecido: %s", job.Job)
	}
}
