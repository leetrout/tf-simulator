# Makefile for the nimbus Terraform provider.
#
# All go commands require the project-local toolchain. Run from this directory;
# the recipes source ../../.toolchain/env.sh themselves.

SHELL := /bin/bash
BINARY := terraform-provider-nimbus
BIN_DIR := $(CURDIR)/bin
BIN := $(BIN_DIR)/$(BINARY)
ENV := source ../../.toolchain/env.sh

.PHONY: build vet test tidy clean install-dir

## build: compile the provider plugin to ./bin/terraform-provider-nimbus
build:
	mkdir -p $(BIN_DIR)
	$(ENV) && go build -o $(BIN) .
	@echo "built $(BIN)"

## vet: run go vet
vet:
	$(ENV) && go vet ./...

## test: run go test
test:
	$(ENV) && go test ./...

## tidy: tidy go.mod/go.sum
tidy:
	$(ENV) && go mod tidy

## install-dir: print the directory to use in ~/.terraformrc dev_overrides
install-dir: build
	@echo "Add this to ~/.terraformrc:"
	@echo ''
	@echo 'provider_installation {'
	@echo '  dev_overrides {'
	@echo '    "leetrout/nimbus" = "$(BIN_DIR)"'
	@echo '  }'
	@echo '  direct {}'
	@echo '}'

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR)
