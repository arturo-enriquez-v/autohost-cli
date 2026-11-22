# autohost-cli

CLI `autohost` (ver `../CLAUDE.md`): prepara un servidor Linux (Docker, Tailscale, KVM), instala el agente,
lo enrola en la API, conecta la VPN y gestiona apps locales. Go 1.23, Cobra. Módulo `autohost-cli`. Salida,
ayuda y README en inglés.

## Comandos

```bash
make build        # ./autohost para esta máquina
make test         # go test ./...
go vet ./...
make build-all    # linux y darwin, amd64 y arm64
make release TAG=v1.2.3   # crea y sube el tag; CI publica (solo cuando se pida)
```

Para probar en un servidor real se usa el lab del agente: `make -C ../autohost-agent lab-update` copia
agente + CLI.

## Estructura

```
cmd/autohost-cli/<comando>/  un paquete por comando Cobra (agent, app, cc, cfg, enroll, setup, up, vpn);
                             root.go hace la inyección de dependencias y registra los comandos
internal/core/{apps,node,network}  lógica (servicios y puertos)
internal/adapters/           docker, tailscale, hypervisor, enrollapi (cliente HTTP de la API), agentconfig,
                             installed (estado local), infra (detección de SO y red)
internal/platform/           config (~/.autohost/config.yaml), exec (Runner), fs, paths, providers
                             (recetas de instalación y URLs)
internal/ui/prompt           prompts interactivos
```

## Reglas

- Los comandos de `cmd/` solo parsean flags y llaman a servicios de `internal/core`; la lógica no va en el
  comando.
- Todo comando de sistema pasa por el `Runner` de `platform/exec`, inyectable, para poder testear sin root ni
  red. Argumentos como lista, nunca un string para `sh -c` con datos del usuario.
- URLs y versiones de lo que se instala (Docker, Tailscale, agente) van en `platform/providers`, no
  repartidas por el código.
- Cuando la CLI llama a la API (`enrollapi`), el contrato es `../autohost-cloud-api/api/openapi.yaml`.
- Secretos (authkeys, tokens de enrolamiento) no se imprimen ni se escriben en logs; los ficheros con
  credenciales se crean con permisos 0600.
- Debe compilar para darwin además de linux: lo específico de Linux va detrás de build tags o detecta el SO.
