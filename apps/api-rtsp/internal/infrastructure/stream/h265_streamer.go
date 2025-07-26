package stream

import (
	"bytes"
	"image/jpeg"
	"log"
	"time"

	"github.com/bluenviron/gortsplib/v4"
	"github.com/bluenviron/gortsplib/v4/pkg/base"
	"github.com/bluenviron/gortsplib/v4/pkg/format"
	"github.com/bluenviron/gortsplib/v4/pkg/format/rtph265"
	"github.com/bluenviron/mediacommon/v2/pkg/codecs/h265"
	"github.com/pion/rtp"
)

func StartH265Stream(rtspURL string, hub *Hub, name string) {
	log.Printf("[%s] conectando ao stream H265: %s", name, rtspURL)

	c := gortsplib.Client{}

	u, err := base.ParseURL(rtspURL)
	if err != nil {
		log.Printf("[%s] URL inválida: %v", name, err)
		return
	}

	if err := c.Start(u.Scheme, u.Host); err != nil {
		log.Printf("[%s] erro ao conectar: %v", name, err)
		return
	}
	defer c.Close()

	desc, _, err := c.Describe(u)
	if err != nil {
		log.Printf("[%s] erro ao descrever stream: %v", name, err)
		return
	}

	var forma *format.H265
	media := desc.FindFormat(&forma)
	if media == nil {
		log.Printf("[%s] stream não contém H265", name)
		return
	}

	rtpDec, err := forma.CreateDecoder()
	if err != nil {
		log.Printf("[%s] erro no decoder RTP: %v", name, err)
		return
	}

	h265Dec := &h265Decoder{}
	if err := h265Dec.initialize(); err != nil {
		log.Printf("[%s] erro ao inicializar decoder H265: %v", name, err)
		return
	}
	defer h265Dec.close()

	// Envia VPS, SPS e PPS se disponíveis
	if forma.VPS != nil {
		h265Dec.decode([][]byte{forma.VPS})
	}
	if forma.SPS != nil {
		h265Dec.decode([][]byte{forma.SPS})
	}
	if forma.PPS != nil {
		h265Dec.decode([][]byte{forma.PPS})
	}

	if _, err := c.Setup(desc.BaseURL, media, 0, 0); err != nil {
		log.Printf("[%s] erro no setup: %v", name, err)
		return
	}

	firstRandomAccess := false
	frameCount := 0
	lastSent := time.Now()

	c.OnPacketRTP(media, forma, func(pkt *rtp.Packet) {
		au, err := rtpDec.Decode(pkt)
		if err != nil && err != rtph265.ErrNonStartingPacketAndNoPrevious && err != rtph265.ErrMorePacketsNeeded {
			log.Printf("[%s] erro ao decodificar pacote RTP: %v", name, err)
			return
		}
		if au == nil {
			return
		}

		if !firstRandomAccess && !h265.IsRandomAccess(au) {
			log.Printf("[%s] aguardando I-frame", name)
			return
		}
		firstRandomAccess = true

		img, err := h265Dec.decode(au)
		if err != nil {
			log.Printf("[%s] erro ao decodificar H265: %v", name, err)
			return
		}
		if img == nil {
			log.Printf("[%s] quadro inválido", name)
			return
		}

		// Throttle
		if time.Since(lastSent) < 10*time.Millisecond {
			return
		}
		lastSent = time.Now()

		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60}); err != nil {
			log.Printf("[%s] erro ao codificar JPEG: %v", name, err)
			return
		}

		hub.Broadcast(buf.Bytes())
		frameCount++

		if frameCount%15 == 0 {
			log.Printf("[%s] %d frames enviados", name, frameCount)
		}
	})

	if _, err := c.Play(nil); err != nil {
		log.Printf("[%s] erro ao iniciar streaming: %v", name, err)
		return
	}

	log.Printf("[%s] stream H265 iniciado", name)
	_ = c.Wait()
	log.Printf("[%s] stream finalizado", name)
}
