# agentd

**agentd** is a passive agent that runs as a daemon (background service), executing configurable local jobs, sending data to a central system (`tower`), and operating in distributed environments like servers, VMs, or embedded devices such as Raspberry Pi.

The project aims to solve data pipeline issues in environments with VPNs — that's why it is pull-based, allowing for extension and the creation of data pipelines.

## Features

- Executes **periodic jobs** based on a local JSON configuration
- Each job runs with its own interval
- Collects and stores local data
- Cross-platform distribuition as a service (Linux, MacOS and Windows)

## Project Structure

```bash
agentd/
├── cmd/                     # Service entrypoint (main.go)
├── internal/
│   ├── config/              # JSON config reader (config.go)
│   ├── jobs/                # Executable job modules
│   └── logger/              # Shared logging system
│   └── utils/               # Helpers: command execution, notifications, results
├── scripts/                 # Linux/Windows install scripts, build tools
├── config.example.json      # Configuration example
└── README.md
```

## Available Jobs

- `check_system.go` – CPU, memory, disk, uptime
- `check_processes.go` – Active processes
- `check_network.go` – Interfaces and connectivity
- `check_users.go` – Logged-in users and sessions
- `check_ports.go` – Open/listening ports
- `check_devices.go` – Devices and mounted disks
- `check_files.go` – Key file existence/validation
- `check_hostinfo.go` – Basic host information
- `check_services.go` – Service status (e.g., systemd)
- `check_suspicious.go` – Suspicious activity detection
- `check_updates.go` – Update checker (integration with tower)

## Roadmap

- [ ] Pull-based integration with **tower** project
- [ ] Self-updating agent with binary replacement
- [ ] Host CDN or use GitHub Releases for distribution
- [ ] Prometheus support with local `/metrics` endpoint
- [ ] Plugin support: external binary job execution
- [ ] Stealth mode for isolated deployments

## Note

This project is part of the [lobofoltran/homelab](https://github.com/lobofoltran/homelab), an ecosystem of tools in Go to study networking, distributed systems, and infrastructure automation in real environments.