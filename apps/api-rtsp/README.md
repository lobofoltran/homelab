# api-rtsp

**api-rtsp** is a service designed to consume RTSP streams from IP cameras, convert them to MJPEG or H.265, and serve them over HTTP for browser-friendly or downstream consumption. It is optimized for embedded setups, local monitoring, or as part of a broader video processing pipeline.

## Features

- Pulls video streams from RTSP cameras (configured via `config.json`)
- Supports **H.264/H.265 decoding** and **FFmpeg-based MJPEG conversion**
- Serves video streams as HTTP endpoints (MJPEG or raw)
- Built-in HTTP server for direct consumption or testing
- Modular architecture with support for new codecs or formats

## Project Structure

```bash
api-rtsp/
├── cmd/                         # Entrypoint (main.go)
├── internal/
│   ├── infrastructure/
│   │   ├── config/              # Configuration loader (config.go)
│   │   └── logger/              # Logging abstraction
│   ├── stream/
│   │   ├── ffmpeg_stream.go     # MJPEG conversion via FFmpeg
│   │   ├── h265_decoder.go      # Decode H.265 to raw
│   │   ├── h265_streamer.go     # Stream H.265 directly
│   │   └── hub.go               # Hub for managing multiple camera streams
│   └── interfaces/server/
│       └── http_server.go       # HTTP endpoints for stream serving
├── config.json                  # Configuration for cameras
├── Dockerfile                   # Containerization
├── go.mod / go.sum              # Go modules
└── README.md
```

## Use cases

- Embedded camera streaming via Raspberry Pi
- Serving MJPEG to frontend (React, browser, etc)
- Stream ingestion for analytics or event detection systems
- Hybrid setups with MJPEG for dashboard and H.265 for archival

## Roadmap

- [ ] Authentication on stream endpoints
- [ ] Stream transcoding fallback if codec unsupported
- [ ] Forward stream to MQTT or Kafka (experimental)
- [ ] Add Prometheus `/metrics` support

## Note

This project is part of the [lobofoltran/homelab](https://github.com/lobofoltran/homelab), an ecosystem of services written in Go for networking, observability, and infrastructure testing.