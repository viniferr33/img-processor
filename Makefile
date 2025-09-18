# Default goal
.DEFAULT_GOAL := help

ifneq (,$(wildcard ./.env))
	include .env
	export
endif

APP_NAME := img-processor
BINARY_NAME := bin/$(APP_NAME)

GO := go
GOFMT := gofmt
GOLINT := golint
GO_TEST_FLAGS := -v -race -cover

DC := docker compose
DC_FILE := docker-compose.yml

.PHONY: help
help: ## Display this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: env
env: ## Check environment variables (debug)
	@echo "Current environment:"
	@echo "  DB_HOST: $(DB_HOST)"
	@echo "  DB_PORT: $(DB_PORT)"
	@echo "  DB_USER: $(DB_USER)"
	@echo "  DB_NAME: $(DB_NAME)"
	@echo "  MINIO_ROOT_USER: $(MINIO_ROOT_USER)"

.PHONY: build
build: ## Build the application
	$(GO) build -o $(BINARY_NAME) ./cmd/server

.PHONY: run
run: ## Run the application
	$(GO) run ./cmd/main.go

.PHONY: test
test: ## Run tests
	$(GO) test $(GO_TEST_FLAGS) ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage report
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

.PHONY: lint
lint: ## Run linter
	$(GOFMT) -d .
	$(GOLINT) ./...

.PHONY: tidy
tidy: ## Tidy go.mod
	$(GO) mod tidy

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/ coverage.out

# Docker Compose targets
.PHONY: dc-up
dc-up: ## Start docker-compose services
	$(DC) -f $(DC_FILE) up -d

.PHONY: dc-down
dc-down: ## Stop docker-compose services
	$(DC) -f $(DC_FILE) down

.PHONY: dc-logs
dc-logs: ## View docker-compose logs
	$(DC) -f $(DC_FILE) logs -f

.PHONY: dc-ps
dc-ps: ## Show docker-compose status
	$(DC) -f $(DC_FILE) ps

.PHONY: dc-restart
dc-restart: ## Restart docker-compose services
	$(DC) -f $(DC_FILE) restart

.PHONY: dc-build
dc-build: ## Build docker-compose images
	$(DC) -f $(DC_FILE) build

# Development utilities
.PHONY: watch
watch: ## Watch for changes and restart server (requires air/gin)
	@if command -v air >/dev/null 2>&1; then \
		air; \
	elif command -v gin >/dev/null 2>&1; then \
		gin run ./cmd/server; \
	else \
		echo "Install air (go install github.com/cosmtrek/air@latest) or gin (go install github.com/codegangsta/gin@latest) for hot reload"; \
		exit 1; \
	fi

.PHONY: generate
generate: ## Run go generate
	$(GO) generate ./...

.PHONY: vendor
vendor: ## Vendor dependencies
	$(GO) mod vendor

# Combined targets
.PHONY: dev
dev: dc-up run ## Start development environment (db + app)

.PHONY: dev-full
dev-full: dc-up migrate-up run ## Start full dev environment with migrations

.PHONY: reset
reset: dc-down dc-up migrate-up ## Reset and rebuild everything

.PHONY: check
check: lint test ## Run all checks (lint + test)