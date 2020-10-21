BIN_DIR := .tools/bin

GO := go
ifdef GO_BIN
	GO = $(GO_BIN)
endif

GOLANGCI_LINT_VERSION := 1.27.0
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint_$(GOLANGCI_LINT_VERSION)

## all: The default target. Build, test, lint
all: test fmt lint

## tidy: go mod tidy
tidy:
	$(GO) mod tidy -v

## fmt: format all go code
fmt:
	gofmt -s -w .

## build: build all files, including protoc if included
build:
	$(GO) build ./...

## test: Run all tests
test: build
	$(GO) test -cover -race -v ./...

## test-coverate: Run all tests and collect coverage
test-coverage:
	$(GO) test ./... -race -coverprofile=.testCoverage.txt && $(GO) tool cover -html=.testCoverage.txt

## setup: project setup, git hooks and validations
setup: setup-validations setup-git-hooks

setup-git-hooks:
	git config core.hooksPath .githooks

setup-validations:
	# A set of validations to help the first time use kick in
	scripts/preflight-checks.sh

## lint: lint all go code
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fast --enable-all -D wsl -D testpackage -D godot -D goerr113 -D gochecknoglobals -D prealloc -D gomnd

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)
	mv $(BIN_DIR)/golangci-lint $(GOLANGCI_LINT)

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
