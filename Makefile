GOPATH := $(shell go env GOPATH)
GO_BIN_DIR := $(GOPATH)/bin

GO_LINT := $(GO_BIN_DIR)/golangci-lint
GO_LINT_URI=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: build
build:
	@go build -v ./...

.PHONY: lint
lint:
	$(if $(GO_LINT), ,go install $(GO_LINT_URI))
	golangci-lint run -v

.PHONY: test
test:
	@go test -race -v ./...

.PHONY: test/cover
test/cover:
	@go test -v -race -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out
