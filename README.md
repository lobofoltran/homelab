# homelab

**homelab** is a modular and extensible project designed to explore and experiment with networking, infrastructure, observability, and agent-based communication using Go.  
The goal is to simulate a real-world, distributed environment with components that communicate over sockets, HTTP, gRPC, and more — both locally and remotely.

## Architecture Overview

This repository is organized as a **monorepo** using `go.work`, composed of multiple independent services, libraries, and experiments:

### Core Applications (apps/)

| Service       | Description |
|---------------|-------------|
| `agentd`      | Passive agent that runs as a daemon, collects metrics, and sends data to `tower` (pull-based) |
| `tower`       | Central command and API hub that receives data and dispatches tasks to `agentd` |
| `daemon`      | Local API layer (REST/gRPC) running on the host, acting as a socket server |
| `cli`         | Command-line tool for interacting with the homelab ecosystem |
| `api-rtsp`    | RTSP camera ingestion service that transcodes or exposes MJPEG/H.265 streams via HTTP |
| `frontend`    | Web dashboard (React) for live camera view or remote monitoring (WIP) |

### Experiments (experiments/)

| Module         | Description |
|----------------|-------------|
| `tcp-tests`    | Scripts and logic for testing TCP behavior between hosts or interfaces |
| `grpc`         | gRPC playground and route design |
| `mqtt`         | MQTT integration testing (agent communication) |
| `graphql`      | GraphQL API experiments and federation structure |

### Infrastructure (infrastructure/)

| Folder       | Description |
|--------------|-------------|
| `docker`     | Dockerfiles for individual apps |
| `compose`    | `docker-compose.yaml` files for local orchestration |
| `k8s`        | Kubernetes manifests and service definitions |
| `terraform`  | Infrastructure as code (cloud experiments, provisioning VMs, etc) |

### Shared Libraries (libs/)

| Library       | Description |
|---------------|-------------|
| `logger`      | Common logging interface for all services |
| `config`      | Shared configuration loader/parser for `config.json` |

## Goals

- Build a real-world inspired internal network and service mesh
- Explore patterns like service discovery, observability, and remote agents
- Learn Go deeply through modular architecture
- Experiment with protocols: HTTP, gRPC, MQTT, Unix sockets, Prometheus
- Self-host and monitor everything on a Raspberry Pi or home cloud

## Roadmap

- [ ] Integrate all agents with `tower` for secure pull-based control
- [ ] Add Prometheus metrics to all services
- [ ] Deploy entire stack to Raspberry Pi or cloud environment
- [ ] Add GitHub Actions workflows for CI/CD
- [ ] Release CLI and agent binaries for Windows/Linux

## Philosophy

The project is built to learn from scratch, document architecture patterns, and design reliable systems that could serve as foundations for infrastructure automation, edge computing, and home-grown observability platforms.

## License

MIT © Gustavo Lobo