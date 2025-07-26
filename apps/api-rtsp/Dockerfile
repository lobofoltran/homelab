FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y \
  pkg-config \
  ffmpeg \
  libavcodec-dev \
  libavutil-dev \
  libswscale-dev \
  build-essential \
  git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main ./cmd/main.go

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y \
  ffmpeg \
  libavcodec58 \
  libavutil56 \
  libswscale5 \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .

COPY config.json /app/config.json

EXPOSE 8080

CMD ["./main"]
