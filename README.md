![autohost-cli](./autohost-cli.png)

# AutoHost CLI

CLI to automate self-hosting on Linux: install applications, manage the agent, and connect nodes to the [autohost-cloud-api](../autohost-cloud-api).

## Stack

- **Go** 1.23.0
- **Cobra** (CLI commands)
- **SQLite** (local state via modernc/sqlite)
- **Docker** (app backend)
- **Caddy / CoreDNS / Tailscale / Cloudflare Tunnel** (service exposure)

## Prerequisites

- Linux (amd64 or arm64)
- Go 1.23+ (only to build from source)
- Docker installed and running
- `sudo` or root privileges for system operations

---

## Installation

### From GitHub Releases (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

The script auto-detects your OS and architecture, downloads the latest binary from GitHub Releases, verifies the checksum, and installs it.

If the binary lands in `~/.local/bin`, add it to your PATH:

```bash
export PATH="$PATH:$HOME/.local/bin"
source ~/.bashrc
```

### Specific version

```bash
VERSION=v0.4.0 curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

### Verify installation

```bash
autohost --version
```

---

## 🛠 Basic Usage

### Example workflow

```bash
export PATH="$PATH:$HOME/.local/bin"   # also add this to ~/.bashrc

autohost setup docker      # install Docker and create the shared network
autohost agent install     # install the agent as a systemd service
autohost up                # authorize this node in the browser and connect it to the cloud
```

### Build from source

```bash
git clone https://github.com/mazapanuwu13/autohost-cli.git
cd autohost-cli
go build -o autohost main.go
```

---

## Available commands

### `autohost agent` — Agent management

```bash
autohost agent install   # Download and install autohost-agent as a systemd service
autohost agent update    # Update the agent to the latest version and restart it
```

After installing, edit `/etc/autohost/config.yaml` and start the service:

```bash
sudo nano /etc/autohost/config.yaml
sudo systemctl enable autohost-agent
sudo systemctl start autohost-agent
```

### `autohost enroll` — Node enrollment

```bash
autohost enroll link --token <ENROLLMENT_TOKEN> [--api <API_URL>] [--name <name>]
```

Links the current node to the cloud API. The command:
1. Auto-detects hostname, local IP, OS and architecture
2. Calls `POST /v1/enrollments/enroll` with the enrollment token (uses the configured default API or the one given with `--api`)
3. Saves the resulting `agent_token` to `/etc/autohost/config.yaml`

### `autohost config` — CLI configuration

View and change the CLI configuration persisted in `~/.autohost/config.yaml`.

```bash
# Show the active configuration
autohost config view

# Change the default URLs
autohost config set api_url http://localhost:8080
autohost config set cloud_url http://localhost:3000
autohost config set vpn_server https://hs.autohst.dev

# Get a single value
autohost config get api_url
```

#### Configuration precedence:
1. **CLI flags** (`--api`, `--cloud`, `--login-server`)
2. **Environment variables** (`AUTOHOST_API_URL`, `AUTOHOST_CLOUD_URL`, `AUTOHOST_VPN_SERVER`)
3. **Config file** (`~/.autohost/config.yaml`)
4. **Defaults** (`https://api.autohst.dev`, `https://cloud.autohst.dev`, `https://hs.autohst.dev`)

### Customizing providers (`providers.yaml`)

The recipes and URLs used to install dependencies (`docker`, `tailscale`, `hypervisor`, `agent`) are embedded in the binary by default (`providers.yaml`). If you need to customize scripts, repositories or packages for a specific or air-gapped environment, you can create an override file at `~/.autohost/providers.yaml`.

Any key defined in `~/.autohost/providers.yaml` transparently overrides the default values, without rebuilding the CLI.

---

## Project structure

```
.
├── cmd/autohost-cli/
│   ├── root.go           # CLI entry point, dependency injection and command registration
│   ├── agent/            # autohost agent (install, update)
│   ├── app/              # autohost app (ls, start, stop, status, remove)
│   ├── cc/               # autohost cc (create and list remote scripts)
│   ├── cfg/              # autohost config (view, set, get)
│   ├── enroll/           # autohost enroll (link node with the API)
│   ├── setup/            # autohost setup (Docker, Tailscale, KVM/QEMU)
│   ├── up/               # autohost up (quick interactive connection to the cloud)
│   └── vpn/              # autohost vpn (connect, disconnect, logout)
├── internal/
│   ├── core/             # Bounded contexts (domain, ports and services)
│   │   ├── apps/         # App lifecycle and custom commands (cc)
│   │   ├── node/         # Nodes, systemd agent, setup and cloud connection (up)
│   │   └── network/      # VPN connect, disconnect and status reporting (Headscale/Tailscale)
│   ├── adapters/         # Infrastructure adapters (with injectable Runner)
│   │   ├── agentconfig/  # Parsing and updating the agent's config.yaml
│   │   ├── docker/       # Docker daemon and Compose adapter
│   │   ├── enrollapi/    # HTTP client for the enrollment API
│   │   ├── hypervisor/   # KVM virtualization checks and installation
│   │   ├── infra/        # Platform, OS and network detection, interactive steps
│   │   ├── installed/    # Filesystem persistence of installed apps (JSON store)
│   │   └── tailscale/    # Tailscale/Headscale management adapter
│   ├── platform/         # Base infrastructure and low-level utilities
│   │   ├── config/       # CLI configuration (~/.autohost/config.yaml, env vars and defaults)
│   │   ├── env/          # Environment variable and placeholder substitution
│   │   ├── exec/         # System command abstraction and execution (Runner)
│   │   ├── fs/           # Filesystem helpers and template copying
│   │   ├── paths/        # System path management (~/.autohost)
│   │   └── providers/    # Centralized install recipes and provider URLs (Docker, Tailscale, KVM, Agent)
│   ├── security/         # Secure generation of cryptographic keys and secrets
│   └── ui/               # Terminal presentation (interactive prompts and validation)
├── scripts/
│   ├── install.sh        # CLI installer
│   ├── autohost-multipass.sh
│   └── autohost-incus.sh
├── Makefile
└── main.go
```

---

## Make commands (test environments)

```bash
# Multipass
make vm-run        # Create a test VM (autohost-test) with Multipass
make vm-update     # Update the binary in the VM
make vm-delete     # Delete the VM

# Incus
make incus-run     # Create a test instance with Incus
make incus-update  # Update the binary in the instance
make incus-delete  # Delete the instance
```

---

## Code conventions

- Go 1.23.0, `gofmt` / `goimports` required
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Idempotent commands: create if missing, update if changed, never duplicate
- No secrets in logs or error messages
- Detect root with `os.Geteuid() == 0`; fall back to `sudo` when available


## License

This project is free and open-source software, distributed under the **[GNU Affero General Public License v3.0 (AGPLv3)](LICENSE)**.

You are free to use it, modify it and host your own (self-hosted) instance at no cost. Under the terms of the AGPLv3, if you modify the software and offer it as a service over a network, you must make the corresponding source code available under the same license.

For alternative commercial licensing or questions about the managed SaaS version, contact the AutoHost team.
