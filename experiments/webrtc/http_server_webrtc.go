package webrtc

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"my-home-backend/internal/infrastructure/config"
// 	"my-home-backend/internal/infrastructure/stream"
// 	"my-home-backend/internal/webrtcHelper"

// 	"github.com/pion/webrtc/v3"
// )

// func Start(cfg *config.AppConfig) {
// 	http.HandleFunc("/ws/camera/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		if r.Method != http.MethodPost {
// 			log.Printf("Método não suportado")
// 			http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		parts := strings.Split(r.URL.Path, "/")
// 		if len(parts) < 4 {
// 			http.Error(w, "invalid camera path", http.StatusBadRequest)
// 			return
// 		}

// 		id, err := strconv.Atoi(parts[3])
// 		if err != nil || id < 1 || id > len(cfg.Cameras) {
// 			http.Error(w, "invalid camera id", http.StatusBadRequest)
// 			return
// 		}

// 		var msg struct {
// 			SDPOffer string `json:"sdp"`
// 		}

// 		// Lê o corpo da requisição
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Printf("Erro ao ler o corpo da requisição: %v", err)
// 			http.Error(w, "Erro ao ler o corpo", http.StatusBadRequest)
// 			return
// 		}

// 		err = json.Unmarshal(body, &msg)
// 		if err != nil {
// 			log.Printf("Erro ao decodificar JSON: %v", err)
// 			http.Error(w, "SDP inválido", http.StatusBadRequest)
// 			return
// 		}

// 		log.Printf("creating peer connection")
// 		// Criar conexão WebRTC
// 		peerConnection, videoTrack, err := webrtcHelper.CreatePeerConnection()
// 		if err != nil {
// 			log.Printf("failed to create peer connection")
// 			http.Error(w, "failed to create peer connection", http.StatusInternalServerError)
// 			return
// 		}
// 		defer peerConnection.Close()

// 		// Setar remote description (offer recebido)
// 		offer := webrtc.SessionDescription{
// 			Type: webrtc.SDPTypeOffer,
// 			SDP:  msg.SDPOffer,
// 		}
// 		if err := peerConnection.SetRemoteDescription(offer); err != nil {
// 			log.Printf("failed to set remote description")
// 			http.Error(w, "failed to set remote description", http.StatusInternalServerError)
// 			return
// 		}

// 		// Criar SDP answer e enviar de volta
// 		answer, err := peerConnection.CreateAnswer(nil)
// 		if err != nil {
// 			log.Printf("failed to create answer")
// 			http.Error(w, "failed to create answer", http.StatusInternalServerError)
// 			return
// 		}
// 		_ = peerConnection.SetLocalDescription(answer)

// 		w.Header().Set("Content-Type", "application/json")
// 		_ = json.NewEncoder(w).Encode(map[string]string{
// 			"sdp": answer.SDP,
// 		})

// 		// Iniciar streaming (adaptado para enviar para videoTrack)
// 		go stream.SendFramesToTrack(cfg.Cameras[id-1].URL, videoTrack)
// 	})

// 	addr := fmt.Sprintf("0.0.0.0:%d", cfg.ServerPort)
// 	log.Printf("Servidor WebRTC ouvindo em %s", addr)
// 	log.Fatal(http.ListenAndServe(addr, nil))
// }
