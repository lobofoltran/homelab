# cli

**cli** is a command-line interface tool for interacting with the services of the `homelab` ecosystem.  
It communicates directly with local daemons or agents, either through Unix sockets, HTTP APIs, or gRPC, depending on the command.

## Features

- Modular command system using [Cobra](https://github.com/spf13/cobra)
- Supports commands like `ping`, `version`, and more
- Can communicate with local daemons via socket (e.g. to test responsiveness)
- Useful for scripting, debugging, and triggering service-side actions

## Project Structure

```bash
cli/
├── cmd/               # Cobra command definitions
│   ├── ping.go        # `ping` command (e.g., checks daemon responsiveness)
│   ├── root.go        # Base/root command
│   └── version.go     # `version` command (prints version/build info)
├── internal/          # (Optional) Internal helper packages
├── main.go            # CLI entrypoint
├── Makefile           # Build and install shortcuts
├── go.mod / go.sum    # Go module dependencies
└── README.md
```

## Example Usage

```bash
homelab ping
> pong

homelab version
> v0.1.3 (build abc123)
```

## Use cases

- Check health/status of local daemons (e.g., `ping`)
- Trigger local service operations without needing external APIs
- Wrap as a binary tool and include in CI/CD or admin scripts
- Expandable to interact with `tower`, `agentd`, and more

## Roadmap

- [ ] Add command: `status` to query agentd’s current jobs
- [ ] Add command: `exec` to trigger remote job on tower
- [ ] Add shell autocompletion (bash/zsh/fish)
- [ ] Distribute compiled binaries via GitHub Releases

## Note

This CLI is part of the [lobofoltran/homelab](https://github.com/lobofoltran/homelab) ecosystem, designed for learning and experimenting with networking, observability, automation, and infrastructure written in Go.