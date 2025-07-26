package stream

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
)

func StartFFmpegStream(rtspURL string, hub *Hub, name string) {
	log.Printf("[%s] iniciando FFmpeg com GPU para %s", name, rtspURL)

	cmd := exec.Command("ffmpeg",
		"-hwaccel", "cuda",
		"-i", rtspURL,
		"-vf", "fps=15",
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-q:v", "5",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[%s] erro ao obter stdout: %v", name, err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Printf("[%s] erro ao iniciar ffmpeg: %v", name, err)
		return
	}

	log.Printf("[%s] ffmpeg iniciado com sucesso", name)

	// leitura de frames JPEG (delimitados por cabeçalho JPEG)
	reader := bufio.NewReader(stdout)
	frameCount := 0

	for {
		// procura por SOI (start of image) e EOI (end of image) markers
		frame, err := readJPEGFrame(reader)
		if err != nil {
			log.Printf("[%s] erro ao ler frame JPEG: %v", name, err)
			break
		}

		hub.Broadcast(frame)
		frameCount++

		if frameCount%15 == 0 {
			log.Printf("[%s] %d frames enviados", name, frameCount)
		}
	}

	cmd.Wait()
	log.Printf("[%s] ffmpeg finalizado", name)
}

func readJPEGFrame(r *bufio.Reader) ([]byte, error) {
	var buf bytes.Buffer
	inImage := false

	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		// JPEG start marker (0xFFD8)
		if b == 0xFF {
			next, _ := r.Peek(1)
			if len(next) == 1 && next[0] == 0xD8 {
				inImage = true
			}
		}

		if inImage {
			buf.WriteByte(b)
		}

		// JPEG end marker (0xFFD9)
		if b == 0xD9 && buf.Len() > 2 {
			break
		}
	}

	return buf.Bytes(), nil
}
