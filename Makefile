GOPATH := $(shell go env GOPATH)
GO_BIN_DIR := $(GOPATH)/bin

GO_LINT := $(GO_BIN_DIR)/golangci-lint
GO_LINT_URI=github.com/golangci/golangci-lint/cmd/golangci-lint@latest

BENCH_COUNT ?= 10
BENCH_PKG := ./internal/benchmarks/suites/...
BENCH_TIME ?= 1s
BENCH_SUITES := $(shell go test -list . $(BENCH_PKG) | grep "Benchmark" | cut -d'_' -f1 | sed 's/Benchmark//' | sort -u)

ifeq ($(HEAVY),1)
    BENCH_COUNT := 20
    BENCH_TIME  := 3s
endif

.PHONY: build
build:
	@go build -v ./...

.PHONY: lint
lint:
	$(if $(GO_LINT), ,go install $(GO_LINT_URI))
	golangci-lint run -v

.PHONY: test
test:
	@go test -race ./...

.PHONY: test/cover
test/cover:
	@go test -race -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out

.PHONY: bench
bench:
	$(eval SUITES_TO_RUN := $(if $(SUITE),$(SUITE),$(BENCH_SUITES)))
	@for suite in $(SUITES_TO_RUN); do \
		printf "\ngo test -bench=%s -count=%s -benchtime=%s %s\n" \
			"$$suite" "$(BENCH_COUNT)" "$(BENCH_TIME)" "$(BENCH_PKG)"; \
		go test -bench=$$suite -count=$(BENCH_COUNT) -benchtime=$(BENCH_TIME) $(BENCH_PKG) \
		| tee /dev/stderr \
		| benchstat -row .fullname -col /Library /dev/stdin; \
	done
