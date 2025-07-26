package webrtc

import (
	"log"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3"
)

func SendFramesToTrack(rtspURL string, track *webrtc.TrackLocalStaticRTP) error {
	log.Printf("Conectando à câmera RTSP: %s", rtspURL)

	client := gortsplib.Client{}

	u, err := base.ParseURL(rtspURL)
	if err != nil {
		return err
	}

	if err := client.Start(u.Scheme, u.Host); err != nil {
		return err
	}
	defer client.Close()

	desc, _, err := client.Describe(u)
	if err != nil {
		return err
	}

	var forma *format.H265
	media := desc.FindFormat(&forma)
	if media == nil {
		return err
	}

	if _, err := client.Setup(desc.BaseURL, media, 0, 0); err != nil {
		return err
	}

	// Callback para cada pacote RTP recebido
	client.OnPacketRTP(media, forma, func(pkt *rtp.Packet) {
		if pkt == nil {
			return
		}

		err := track.WriteRTP(pkt)
		if err != nil {
			log.Printf("Erro ao enviar para track WebRTC: %v", err)
		}

		log.Printf("Frame enviado: %d bytes", len(pkt.Payload))

		time.Sleep(60 * time.Millisecond) // ~30fps
	})

	if _, err := client.Play(nil); err != nil {
		return err
	}

	log.Println("Streaming iniciado via WebRTC")
	return client.Wait()
}
