ROOT_PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
COVERAGE_PATH := $(ROOT_PATH).coverage/

# Use the locally installed Go toolchain.
export GOTOOLCHAIN=local

# Ensure our local tools are preferred
export PATH := $(ROOT_PATH)tools/bin:$(PATH)

include $(ROOT_PATH)tools/tools.mk

BUILD_DIR := $(ROOT_PATH)build

CLI_REPO ?= https://github.com/e2engine/cli.git
CLI_REF ?= main

CLI_SRC := $(BUILD_DIR)/cli

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= $(shell uname -m)
ifeq ($(GOARCH), x86_64)
	override GOARCH = amd64
endif

CLI_BIN_NAME := e2engine-$(GOOS)-$(GOARCH)$(shell [ "$(GOOS)" = windows ] && echo .exe)
CLI_BIN   := $(ROOT_PATH)cli/bin/cli

.PHONY: build-cli
build-cli:
	@echo "Building e2engine-cli..."
	@rm -rf $(BUILD_DIR)
	@set -e; \
	trap 'rm -rf "$(BUILD_DIR)"' EXIT; \
	echo "Cloning CLI repository..."; \
    git clone $(CLI_REPO) $(CLI_SRC); \
    git -C $(CLI_SRC) checkout $(CLI_REF); \
	echo "Preparing e2engine-cli module..."; \
	cd $(CLI_SRC); \
	mkdir -p $(dir $(CLI_BIN)); \
	echo "Building $(CLI_BIN_NAME) binary..."; \
	$(MAKE) build-cli; \
	cp bin/$(CLI_BIN_NAME) $(CLI_BIN)
	@echo "$(CLI_BIN) built successfully"

SERVICES_BIN_PATH ?= $(ROOT_PATH)services/bin

.PHONY:  build-gateway-http-path-service
build-gateway-http-path-service:
	@echo "Building gateway HTTP path service..."
	@go build -o $(SERVICES_BIN_PATH)/gateway-http-path ./services/gateway-http-path
	@echo "Gateway HTTP path service built successfully at $(SERVICES_BIN_PATH)/gateway-http-path"

.PHONY:  build-gateway-http-body-service
build-gateway-http-body-service:
	@echo "Building gateway HTTP body service..."
	@go build -o $(SERVICES_BIN_PATH)/gateway-http-body ./services/gateway-http-body
	@echo "Gateway HTTP body service built successfully at $(SERVICES_BIN_PATH)/gateway-http-body"

.PHONY:  build-gateway-grpc-message-service
build-gateway-grpc-message-service:
	@echo "Building gateway gRPC message service..."
	@go build -o $(SERVICES_BIN_PATH)/gateway-grpc-message ./services/gateway-grpc-message
	@echo "Gateway gRPC message service built successfully at $(SERVICES_BIN_PATH)/gateway-grpc-message"

.PHONY:  build-gateway-grpc-method-service
build-gateway-grpc-method-service:
	@echo "Building gateway gRPC method service..."
	@go build -o $(SERVICES_BIN_PATH)/gateway-grpc-method ./services/gateway-grpc-method
	@echo "Gateway gRPC method service built successfully at $(SERVICES_BIN_PATH)/gateway-grpc-method"


.PHONY:  build-ok-http-service
build-ok-http-service:
	@echo "Building ok HTTP service..."
	@go build -o $(SERVICES_BIN_PATH)/ok-http ./services/ok-http
	@echo "Ok HTTP service built successfully at $(SERVICES_BIN_PATH)/ok-http"

.PHONY: build-ok-grpc-service
build-ok-grpc-service:
	@echo "Building ok gRPC service..."
	@go build -o $(SERVICES_BIN_PATH)/ok-grpc ./services/ok-grpc
	@echo "Ok gRPC service built successfully at $(SERVICES_BIN_PATH)/ok-grpc"

.PHONY:  protoc-ok-service
protoc-ok-service:
	@echo "Generating protobuf code for ok service..."
	@protoc \
      --go_out=. \
      --go_opt=module=github.com/e2engine/tests \
      --go-grpc_out=. \
      --go-grpc_opt=module=github.com/e2engine/tests \
      services/ok-grpc/ok.proto
	@echo "Protobuf code generated successfully for ok service."

.PHONY: build-image
build-image:
	@echo "Checking Docker availability..."
	@docker info >/dev/null 2>&1 || (echo "Docker is not available" && exit 1)
	@echo "Building e2engine Docker image..."
	@rm -rf $(BUILD_DIR)
	@set -e; \
	trap 'rm -rf "$(BUILD_DIR)"' EXIT; \
	echo "Cloning CLI repository..."; \
    git clone $(CLI_REPO) $(CLI_SRC); \
    git -C $(CLI_SRC) checkout $(CLI_REF); \
	echo "Preparing e2engine-cli module..."; \
	cd $(CLI_SRC); \
	echo "Building e2engine Docker image..."; \
	$(MAKE) build-image
	@echo "e2engine:local image built successfully"


DOCKER_SERVICES_BIN_PATH := $(ROOT_PATH)services/docker-bin

.PHONY: build-ok-http-image
build-ok-http-image:
	@echo "Checking Docker availability..."
	@docker info >/dev/null 2>&1 || (echo "Docker is not available" && exit 1)
	@echo "Building ok HTTP service for Linux..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build \
		-o $(DOCKER_SERVICES_BIN_PATH)/ok-http ./services/ok-http
	@file $(DOCKER_SERVICES_BIN_PATH)/ok-http
	@docker buildx build \
		--platform linux/$(GOARCH) \
		--load \
		--build-arg SERVICE=ok-http \
		-t e2engine-test-ok-http:local \
		-f services/Dockerfile \
		services
	@echo "ok-http image built successfully"


.PHONY: lint
lint: install-golangci-lint
	@echo "Running linter..."
	$(GOLANGCI_LINT) run

.PHONY: build
build: build-cli \
	build-gateway-http-path-service \
	build-gateway-http-body-service \
	build-gateway-grpc-message-service \
	build-gateway-grpc-method-service \
	build-ok-http-service \
	build-ok-grpc-service \
	build-image \
	build-ok-http-image
	@echo "All test binaries and images built successfully."

.PHONY: verify
verify: lint build
	@echo "All verifications passed successfully."

.PHONY: test-cli
test-cli: build-cli \
	build-gateway-http-path-service \
	build-gateway-http-body-service \
	build-gateway-grpc-message-service \
	build-gateway-grpc-method-service \
	build-ok-http-service \
	build-ok-grpc-service
	@echo "Running CLI tests..."
	@go test ./cli/... -count=1

.PHONY: test-docker
test-docker: build-image \
	build-ok-http-image
	@echo "Running tests with Docker image..."
	@go test ./docker/... -count=1
