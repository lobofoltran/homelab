package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lobofoltran/homelab/apps/api-rtsp/internal/infrastructure/config"
	"github.com/lobofoltran/homelab/apps/api-rtsp/internal/infrastructure/stream"
)

var hubs map[string]*stream.Hub

func Start(cfg *config.AppConfig) {
	hubs = make(map[string]*stream.Hub)

	for i, cam := range cfg.Cameras {
		id := strconv.Itoa(i)
		hub := stream.NewHub()
		hubs[id] = hub

		go func(cam config.CameraConfig, id string, hub *stream.Hub) {
			log.Printf("[cam %s] Iniciando hub e stream", cam.Name)
			go hub.Run()

			go stream.StartFFmpegStream(cam.URL, hub, cam.Name)
		}(cam, id, hub)
	}

	http.HandleFunc("/cameras/", handleCameraStream)

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Servidor iniciado em http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleCameraStream(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/cameras/"):]

	hub, ok := hubs[id]
	if !ok {
		http.Error(w, "Câmera não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "close")

	frameChan := make(chan []byte, 10)
	hub.Register(frameChan)
	defer hub.Unregister(frameChan)

	log.Printf("[MJPEG] Cliente conectado para câmera %s", id)

	for {
		select {
		case <-r.Context().Done():
			log.Printf("[MJPEG] Cliente desconectado da câmera %s", id)
			return
		case frame := <-frameChan:
			fmt.Fprintf(w, "--frame\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n", len(frame))
			w.Write(frame)
			fmt.Fprint(w, "\r\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}
}
