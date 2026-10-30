![autohost-cli](./autohost-cli.png)

# AutoHost CLI

CLI para automatizar el self-hosting en Linux: instala aplicaciones, gestiona el agente, y conecta nodos con la [autohost-cloud-api](../autohost-cloud-api).

## Stack

- **Go** 1.23.0
- **Cobra** (comandos CLI)
- **SQLite** (estado local via modernc/sqlite)
- **Docker** (backend de apps)
- **Caddy / CoreDNS / Tailscale / Cloudflare Tunnel** (exposición de servicios)

## Requisitos previos

- Linux (amd64 o arm64)
- Go 1.23+ (solo para compilar desde fuente)
- Docker instalado y en ejecución
- `sudo` o permisos de root para operaciones de sistema

---

## Instalación

### Desde GitHub Releases (recomendado)

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
export PATH="$PATH:$HOME/.local/bin"   # añadir también a ~/.bashrc
```

### Compilar desde fuente

```bash
git clone https://github.com/mazapanuwu13/autohost-cli.git
cd autohost-cli
go build -o autohost main.go
```

---

## Comandos disponibles

### `autohost agent` — Gestión del agente

```bash
autohost agent install   # Descarga e instala el autohost-agent como servicio systemd
```

Tras la instalación, edita `/etc/autohost/config.yaml` y arranca el servicio:

```bash
sudo nano /etc/autohost/config.yaml
sudo systemctl enable autohost-agent
sudo systemctl start autohost-agent
```

### `autohost enroll` — Enrollment de nodos

```bash
autohost enroll link --token <ENROLLMENT_TOKEN> [--api <URL_API>] [--name <nombre>]
```

Enlaza el nodo actual con la cloud API. El comando:
1. Auto-detecta hostname, IP local, OS y arquitectura
2. Llama a `POST /v1/enrollments/enroll` con el token de enrollment (usa la API configurada por defecto o la especificada en `--api`)
3. Guarda el `agent_token` resultante en `/etc/autohost/config.yaml`

### `autohost config` — Configuración del CLI

Permite consultar y cambiar la configuración del CLI persistida en `~/.autohost/config.yaml`.

```bash
# Ver configuración activa
autohost config view

# Cambiar URLs por defecto
autohost config set api_url http://localhost:8080
autohost config set cloud_url http://localhost:3000
autohost config set vpn_server https://hs.autohst.dev

# Obtener un valor
autohost config get api_url
```

#### Precedencia de configuración:
1. **Flags CLI** (`--api`, `--cloud`, `--login-server`)
2. **Variables de entorno** (`AUTOHOST_API_URL`, `AUTOHOST_CLOUD_URL`, `AUTOHOST_VPN_SERVER`)
3. **Archivo de configuración** (`~/.autohost/config.yaml`)
4. **Valores por defecto** (`https://api.autohst.dev`, `https://cloud.autohst.dev`, `https://hs.autohst.dev`)

### Personalización de Proveedores (`providers.yaml`)

Las recetas y URLs para instalar dependencias (`docker`, `tailscale`, `hypervisor`, `agent`) vienen embebidas en el binario por defecto (`providers.yaml`). Si necesitas personalizar scripts, repositorios o paquetes para un entorno particular o air-gapped, puedes crear un archivo de override en `~/.autohost/providers.yaml`.

Cualquier clave definida en `~/.autohost/providers.yaml` sobreescribe de forma transparente los valores predeterminados sin requerir recompilar el CLI.

---

## Estructura del proyecto

```
.
├── cmd/autohost-cli/
│   ├── root.go           # Punto de entrada CLI, inyección de dependencias y registro de comandos
│   ├── agent/            # autohost agent (install, update)
│   ├── app/              # autohost app (ls, start, stop, status, remove)
│   ├── cc/               # autohost cc (creación y listado de scripts remotos)
│   ├── cfg/              # autohost config (view, set, get)
│   ├── enroll/           # autohost enroll (link de nodo con la API)
│   ├── setup/            # autohost setup (Docker, Tailscale, KVM/QEMU)
│   ├── up/               # autohost up (conexión interactiva rápida con Cloud)
│   └── vpn/              # autohost vpn (connect, disconnect, logout)
├── internal/
│   ├── core/             # Bounded Contexts (dominio, puertos y servicios)
│   │   ├── apps/         # Gestión de ciclo de vida de apps y comandos custom (cc)
│   │   ├── node/         # Gestión de nodos, agente systemd, setup y conexión cloud (up)
│   │   └── network/      # Conexión, desconexión y reporte de estado VPN (Headscale/Tailscale)
│   ├── adapters/         # Adaptadores de infraestructura (con Runner inyectable)
│   │   ├── agentconfig/  # Parseo y actualización de config.yaml del agente
│   │   ├── docker/       # Adaptador para Docker daemon y Compose
│   │   ├── enrollapi/    # Cliente HTTP para el API de enrollment
│   │   ├── hypervisor/   # Verificación e instalación de virtualización KVM
│   │   ├── infra/        # Detección de plataforma, SO, red y pasos interactivos
│   │   ├── installed/    # Persistencia en filesystem de apps instaladas (JSON store)
│   │   └── tailscale/    # Adaptador para gestión de Tailscale/Headscale
│   ├── platform/         # Infraestructura base y utilidades de bajo nivel
│   │   ├── config/       # Configuración del CLI (~/.autohost/config.yaml, env vars y defaults)
│   │   ├── env/          # Reemplazo de variables de entorno y placeholders
│   │   ├── exec/         # Abstracción y ejecución de comandos de sistema (Runner)
│   │   ├── fs/           # Helpers de filesystem y copia de templates
│   │   ├── paths/        # Gestión de rutas del sistema (~/.autohost)
│   │   └── providers/    # Recetas centralizadas de instalación y URLs de providers (Docker, Tailscale, KVM, Agent)
│   ├── security/         # Generación segura de claves criptográficas y secretos
│   └── ui/               # Presentación en terminal (prompts interactivos y validaciones)
├── scripts/
│   ├── install.sh        # Instalador del CLI
│   ├── autohost-multipass.sh
│   └── autohost-incus.sh
├── Makefile
└── main.go
```

---

## Comandos Make (entornos de prueba)

```bash
# Multipass
make vm-run        # Crea VM de prueba (autohost-test) con Multipass
make vm-update     # Actualiza el binario en la VM
make vm-delete     # Elimina la VM

# Incus
make incus-run     # Crea instancia de prueba con Incus
make incus-update  # Actualiza el binario en la instancia
make incus-delete  # Elimina la instancia
```

---

## Convenciones de código

- Go 1.23.0, `gofmt` / `goimports` obligatorio
- Errores envueltos con contexto: `fmt.Errorf("contexto: %w", err)`
- Comandos idempotentes: crear si no existe, actualizar si cambió, no duplicar
- Sin secretos en logs ni mensajes de error
- Detectar root con `os.Geteuid() == 0`; usar `sudo` como fallback si está disponible


## Licencia / License

Este proyecto es software libre y de código abierto (Open Source), distribuido bajo la licencia **[GNU Affero General Public License v3.0 (AGPLv3)](LICENSE)**.

Eres libre de usarlo, modificarlo y hospedar tu propia instancia (self-hosted) de forma gratuita. De acuerdo con los términos de AGPLv3, si modificas el software y ofreces servicios a través de una red, debes poner a disposición el código fuente correspondiente bajo la misma licencia.

Para licencias comerciales alternativas o consultas sobre la versión SaaS gestionada, contacta al equipo de AutoHost.
