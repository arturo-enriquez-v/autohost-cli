

.PHONY: build clean release dev-up vm-run vm-update vm-delete incus-run incus-update incus-delete incus-start incus-stop incus-up

BINARY_NAME = autohost
REPO        = mazapanuwu13/autohost-cli
PLATFORMS   = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
INCUS_INSTANCE = autohost-test

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  = -s -w -X autohost-cli/cmd/autohost-cli.Version=$(VERSION)

build:
	@echo "🔨 Building $(BINARY_NAME) $(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) main.go
	@echo "✅ Build complete: ./$(BINARY_NAME)"


release:
	@CURRENT=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	echo "📌 Versión actual: $$CURRENT"; \
	printf "🔖 Nueva versión (ej. v1.2.3): "; \
	read NEW_VERSION; \
	if [ -z "$$NEW_VERSION" ]; then echo "❌ La versión no puede estar vacía"; exit 1; fi; \
	echo "🏷️  Creando tag $$NEW_VERSION..."; \
	git tag -a "$$NEW_VERSION" -m "Release $$NEW_VERSION"; \
	echo "🚀 Compilando release $$NEW_VERSION para: $(PLATFORMS)"; \
	mkdir -p dist; \
	for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; GOARCH=$${platform#*/}; \
		out="dist/$(BINARY_NAME)-$${GOOS}-$${GOARCH}"; \
		echo "  → $$out"; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "-s -w -X autohost-cli/cmd/autohost-cli.Version=$$NEW_VERSION" -o "$$out" main.go; \
	done; \
	echo "🔐 Generating checksums..."; \
	cd dist && sha256sum $(BINARY_NAME)-* > checksums_$${NEW_VERSION}.txt; \
	echo "✅ Release artifacts in dist/"; \
	ls -lh dist/; \
	echo ""; \
	echo "💡 Para publicar el tag ejecuta: git push origin $$NEW_VERSION"



# ===== Incus ====== #

# run-incus: build
# 	@echo "🚀 Creating Incus instance ($(INCUS_INSTANCE))..."
# 	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh run

update-incus:
	@echo "🔄 Updating autohost binary in Incus instance ($(INCUS_INSTANCE))..."
	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh update

# delete-incus:
# 	@echo "🧹 Deleting Incus instance ($(INCUS_INSTANCE))..."
# 	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh delete

# start-incus:
# 	@echo "▶️ Starting Incus instance ($(INCUS_INSTANCE))..."
# 	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh start

# stop-incus:
# 	@echo "⏹ Stopping Incus instance ($(INCUS_INSTANCE))..."
# 	@VM_NAME="$(INCUS_INSTANCE)" bash scripts/autohost-incus.sh stop

# incus-up: update-incus
# 	@echo "🧪 Testing autohost up inside Incus instance ($(INCUS_INSTANCE))..."
# 	@GATEWAY=$$(incus exec $(INCUS_INSTANCE) -- sh -c "ip route show default" | awk '/default/{print $$3}' | head -1); \
# 	 echo "   Container gateway (host): $$GATEWAY"; \
# 	 echo "   Cloud URL (browser): http://localhost:3000"; \
# 	 echo "   API URL (container -> host): http://$$GATEWAY:8080"; \
# 	 incus exec $(INCUS_INSTANCE) -- env AUTOHOST_CLOUD_URL="http://localhost:3000" AUTOHOST_API_URL="http://$$GATEWAY:8080" /usr/local/bin/autohost up

# ===== Remote Node Config (LAN) ===== #
NODE_IP   ?= 192.168.101.33
NODE_USER ?= arturo
NODE_HOST ?= $(NODE_USER)@$(NODE_IP)
NODE_NAME ?= mazapan
HOST_IP   ?= 192.168.101.2
API_PORT  ?= 8080
GRPC_PORT ?= 9090

.PHONY: node-update node-setup node-setup-hypervisor node-logs node-status node-shell

node-update: build
	@echo "🔨 Compilando autohost-agent..."
	@(cd ../autohost-agent && go build -ldflags "-s -w" -o autohost-agent cmd/agent/main.go)
	@echo "🚀 Transfiriendo binarios a $(NODE_HOST)..."
	@scp $(BINARY_NAME) $(NODE_HOST):/tmp/autohost
	@scp ../autohost-agent/autohost-agent $(NODE_HOST):/tmp/autohost-agent
	@echo "📦 Instalando binarios y reiniciando agente en $(NODE_HOST)..."
	@ssh $(NODE_HOST) "sudo install -m 0755 /tmp/autohost /usr/local/bin/autohost && \
		sudo install -m 0755 /tmp/autohost-agent /usr/local/bin/autohost-agent && \
		rm -f /tmp/autohost /tmp/autohost-agent && \
		sudo systemctl restart autohost-agent"
	@echo "✅ Binarios actualizados y servicio reiniciado en $(NODE_IP)!"
	@ssh $(NODE_HOST) "sudo systemctl status autohost-agent --no-pager | head -n 12"

node-setup: node-update
	@echo "🔑 Generando token de enrolamiento en base de datos local..."
	@TOKEN=$$(python3 -c "import os, base64, hashlib, subprocess; \
		raw = os.urandom(32); \
		plain = 'autohost-enr_' + base64.urlsafe_b64encode(raw).decode().rstrip('='); \
		h = hashlib.sha256(plain.encode()).hexdigest(); \
		cid = subprocess.check_output(['docker', 'ps', '-q', '-f', 'name=postgres']).decode().strip().split()[0]; \
		subprocess.run(['docker', 'exec', '-i', cid, 'psql', '-U', 'autohost', '-d', 'autohost', '-c', \
			f\"INSERT INTO enroll_tokens (token, user_id, expires_at) VALUES ('{h}', 'a5cd786b-6cb5-4a79-997c-efe78b0cfa24', NOW() + INTERVAL '24 hours');\"], \
			check=True, stdout=subprocess.DEVNULL); \
		print(plain)"); \
	 echo "   Token generado: $$TOKEN"; \
	 echo "🔗 Enrolando nodo $(NODE_NAME) en http://$(HOST_IP):$(API_PORT)..."; \
	 ssh $(NODE_HOST) "sudo autohost enroll link --api http://$(HOST_IP):$(API_PORT) --token $$TOKEN --name $(NODE_NAME) && \
		sudo sed -i 's|grpc_address:.*|grpc_address: \"$(HOST_IP):$(GRPC_PORT)\"|g' /etc/autohost/config.yaml && \
		if ! sudo grep -q 'grpc_insecure' /etc/autohost/config.yaml; then echo 'grpc_insecure: true' | sudo tee -a /etc/autohost/config.yaml; fi && \
		sudo systemctl daemon-reload && \
		sudo systemctl enable --now autohost-agent && \
		sudo systemctl restart autohost-agent"; \
	 echo "🎉 Nodo '$(NODE_NAME)' ($(NODE_IP)) configurado y conectado a AutoHost!"; \
	 echo "🌐 Disponible en dashboard: http://localhost:3000"

node-setup-hypervisor:
	@echo "⚡ Instalando KVM/libvirt en nodo remoto $(NODE_HOST)..."
	@ssh $(NODE_HOST) "yes | sudo autohost setup hypervisor && sudo usermod -aG libvirt,kvm autohost && sudo sed -i 's|SupplementaryGroups=docker.*|SupplementaryGroups=docker libvirt kvm|g' /etc/systemd/system/autohost-agent.service 2>/dev/null || true && sudo systemctl daemon-reload && sudo systemctl restart autohost-agent"
	@echo "✅ Hypervisor configurado y servicio reiniciado en $(NODE_IP)!"

node-logs:
	@ssh $(NODE_HOST) "sudo journalctl -u autohost-agent -f"

node-status:
	@ssh $(NODE_HOST) "sudo systemctl status autohost-agent"

node-shell:
	@ssh -t $(NODE_HOST) "bash -l"

