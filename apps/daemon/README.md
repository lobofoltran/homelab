# daemon

**daemon** is a local service that runs as a full-featured, socket-based API layer.  
It exposes a set of routes using REST or gRPC and serves as the runtime interface for interacting with system-level or custom application logic.

## Features

- Exposes APIs over Unix socket or TCP (REST and gRPC supported)
- Modular routing system for defining new endpoints
- Separation between transport layer (router) and business logic (internal/backend)
- Used as the local runtime controller for CLI tools and services like `agentd`
- Includes pluggable handlers like `ping` and `greeter`

## Project Structure

```bash
daemon/
├── cmd/
│   └── main.go               # Service entrypoint
├── daemon/
│   ├── run.go                # Entry to start server
│   ├── server/
│   │   ├── router/
│   │   │   ├── grpc/         # gRPC route definitions
│   │   │   └── rest/         # REST route definitions
│   │   └── server.go         # HTTP/gRPC server wrapper
│   └── router.go             # Router logic and init
├── internal/
│   └── backend/              # Business logic handlers (ping, greeter, etc.)
├── go.mod / go.sum
└── README.md
```

## Endpoints

- `GET /ping` – Returns `pong` (for testing liveness)
- `GET /greet` – Example route to return a greeting
- gRPC routes available via socket or configured port

## Example usage (local)

```bash
go run ./cmd
```

Will start the daemon and expose its routes.  
Example:

```bash
curl --unix-socket /tmp/homelab.sock http://unix/ping
# → pong
```

## Use cases

- Local system API layer for automating or exposing routines
- Integrates with `cli` tool for executing commands
- Serves as a gateway for services running on the machine
- Can evolve into a multi-protocol orchestrator for custom logic

## Roadmap

- [ ] Add job scheduler endpoints
- [ ] Integrate authentication middleware
- [ ] Expose system diagnostics and control commands
- [ ] Auto-reload config via signal
- [ ] Extend gRPC handlers and reflection

## Note

This daemon is part of the [lobofoltran/homelab](https://github.com/lobofoltran/homelab) project, a modular ecosystem for learning, experimentation, and production-ready service design in Go.