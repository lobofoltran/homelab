package main

import (
	"log"

	"github.com/lobofoltran/homelab/apps/api-rtsp/internal/infrastructure/config"
	"github.com/lobofoltran/homelab/apps/api-rtsp/internal/infrastructure/logger"
	"github.com/lobofoltran/homelab/apps/api-rtsp/internal/interfaces/server"
)

func main() {
	logger.Init()

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("erro ao carregar config.json: %v", err)
	}

	server.Start(cfg)
}

// 	url := "rtsp://ajcc:ad4chp@192.168.1.101:554/onvif1"

// 	u, err := base.ParseURL(url)
// 	if err != nil {
// 		log.Fatal("URL inválida:", err)
// 	}

// 	c := gortsplib.Client{}
// 	if err := c.Start(u.Scheme, u.Host); err != nil {
// 		log.Fatal("Erro ao conectar:", err)
// 	}
// 	defer c.Close()

// 	desc, _, err := c.Describe(u)
// 	if err != nil {
// 		log.Fatal("Erro ao descrever o stream:", err)
// 	}

// 	log.Println("=== Medias e Formatos disponíveis ===")
// 	for i, media := range desc.Medias {
// 		log.Printf("Media %d: %s", i, media.Type)
// 		for j, format := range media.Formats {
// 			log.Printf("  Formato %d: %T", j, format)
// 		}
// 	}
// }
