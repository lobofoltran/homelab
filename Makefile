.PHONY: all cli daemon

all: cli daemon

cli:
	go build -o bin/homelab-cli ./cli

daemon:
	go build -o bin/homelabd ./homelab/cmd/
