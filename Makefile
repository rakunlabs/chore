BINARY    := chore
MAIN_FILE := cmd/$(BINARY)/main.go
PKG       := $(shell go list -m)
VERSION   := $(or $(IMAGE_TAG),$(shell git describe --tags --first-parent --match "v*" 2> /dev/null || echo v0.0.0))
LOCAL_BIN_DIR := $(PWD)/bin

## swaggo configuration
SWAG_VERSION := $(shell grep -E 'swaggo/swag' go.mod | awk '{print $$2}')

## golangci configuration
GOLANGCI_CONFIG_URL   := https://raw.githubusercontent.com/worldline-go/guide/main/lint/.golangci.yml
GOLANGCI_LINT_VERSION := v1.55.1

.DEFAULT_GOAL := help

.PHONY: run run-front whoami build build-front build-all copy-front docs golangci lint env test coverage html html-gen html-wsl help

run: export CONFIG_FILE ?= ./_example/config/config.yml
run: export ENV ?= development
run: ## Run the application
	go run $(MAIN_FILE)

run-front: ## Run the front
	(cd _web && pnpm run dev --host)

whoami: ## Run whoami container
	docker run --rm -it --name="whoami" -p 9090:80 traefik/whoami

build: ## Build the binary file
	goreleaser build --snapshot --rm-dist --single-target

build-front: ## Build the front
	(cd _web && pnpm build-front)

build-all: build-front build ## Build front and binary file

copy-front: ## Copy the front
	@echo "> Copying frontend outputs"
	@rm -rf ./internal/server/dist/* 2> /dev/null
	@cp -a _web/dist ./internal/server/.

bin/swag-$(SWAG_VERSION):
	@echo "> downloading swag@$(SWAG_VERSION)"
	@GOBIN=$(LOCAL_BIN_DIR) go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)
	@mv $(LOCAL_BIN_DIR)/swag $(LOCAL_BIN_DIR)/swag-$(SWAG_VERSION)

docs: bin/swag-$(SWAG_VERSION) ## Generate swagger documentation
	@$(LOCAL_BIN_DIR)/swag-$(SWAG_VERSION) init -g handlers.go --dir internal/server,internal/api,models

.golangci.yml:
	@$(MAKE) golangci

golangci: ## Download .golangci.yml file
	@curl --insecure -o .golangci.yml -L'#' $(GOLANGCI_CONFIG_URL)

bin/golangci-lint-$(GOLANGCI_LINT_VERSION):
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCAL_BIN_DIR) $(GOLANGCI_LINT_VERSION)
	@mv $(LOCAL_BIN_DIR)/golangci-lint $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION)

lint-all: .golangci.yml bin/golangci-lint-$(GOLANGCI_LINT_VERSION) ## Lint Go files
	@$(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) --version
	@GOPATH="$(shell dirname $(PWD))" $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) run ./...

lint: .golangci.yml bin/golangci-lint-$(GOLANGCI_LINT_VERSION) ## Lint Go files
	@$(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) --version
	@GOPATH="$(shell dirname $(PWD))" $(LOCAL_BIN_DIR)/golangci-lint-$(GOLANGCI_LINT_VERSION) run --new-from-rev remotes/origin/main ./...

env: ## Create environment
	docker compose --project-name=chore --file=_example/chore/docker-compose.yml up

env-extra: ## Create environment with extra services
	docker compose --profile=extra --project-name=chore --file=_example/chore/docker-compose.yml up

test: ## Run unit tests
	@go test -v -race -cover ./...

coverage: ## Run unit tests with coverage
	@go test -v -race -cover -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

html: ## Show html coverage result
	@go tool cover -html=./coverage.out

html-gen: ## Export html coverage result
	@go tool cover -html=./coverage.out -o ./coverage.html

html-wsl: html-gen ## Open html coverage result in wsl
	@explorer.exe `wslpath -w ./coverage.html` || true

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
