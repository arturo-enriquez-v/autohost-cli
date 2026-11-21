SHELL := /bin/bash
.DEFAULT_GOAL := help

BINARY    := autohost
MAIN      := main.go
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS   := -s -w -X autohost-cli/cmd/autohost-cli.Version=$(VERSION)
AGENT_DIR := ../autohost-agent

.PHONY: help build build-all test clean release

help: ## Show this help
	@awk 'BEGIN{FS=":.*## "} /^##@/{printf "\n\033[1m%s\033[0m\n", substr($$0,5)} \
		/^[a-zA-Z0-9_%-]+:.*## /{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

##@ Build

build: ## Build ./autohost for this machine
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) $(MAIN)

build-all: ## Cross-compile linux and darwin (amd64, arm64) into dist/
	@mkdir -p dist
	@for p in $(PLATFORMS); do os=$${p%/*}; arch=$${p#*/}; out=dist/$(BINARY)-$$os-$$arch; \
		echo "  → $$out"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" -o $$out $(MAIN) || exit 1; \
	done

test: ## Run the tests
	go test ./...

clean: ## Remove build output
	rm -rf $(BINARY) dist

release: ## Tag a release and push the tag; CI builds and publishes it (TAG=v1.2.3)
	@tag="$(TAG)"; \
	[ -n "$$tag" ] || read -rp "New version (latest: $$(git describe --tags --abbrev=0 2>/dev/null || echo none)): " tag; \
	[[ "$$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$$ ]] || { echo "Invalid version '$$tag' (expected vX.Y.Z)"; exit 1; }; \
	[ -z "$$(git status --porcelain)" ] || { echo "Working tree is not clean"; exit 1; }; \
	git tag -a "$$tag" -m "Release $$tag" && git push origin "$$tag" && \
	echo "Pushed $$tag: GitHub Actions will build and publish the release"

##@ Lab (runs in ../autohost-agent, see `make -C ../autohost-agent help`)

lab-%: ## Any lab-* target, e.g. lab-update to push this CLI to the test servers
	@$(MAKE) --no-print-directory -C $(AGENT_DIR) $@

node-%: ## Any node-* target, e.g. node-update
	@$(MAKE) --no-print-directory -C $(AGENT_DIR) $@
