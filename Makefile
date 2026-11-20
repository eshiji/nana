BINARY := nana
MODULE := github.com/eshiji/nana
VERSION ?= dev
LDFLAGS := -s -w -X $(MODULE)/internal/buildinfo.version=$(VERSION)
COVERAGE_MIN := 80
# Tools run with `go run module@version`: pinned, nothing to install globally.
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.10.0
ACTIONLINT := github.com/rhysd/actionlint/cmd/actionlint@v1.7.12

.DEFAULT_GOAL := help
.PHONY: help build test cover lint fmt golden vuln docs docs-check lint-ci snapshot clean

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*## "} /^[a-z-]+:.*## / {printf "  %-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build bin/$(BINARY) (override VERSION=vX.Y.Z to stamp a version)
	go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/$(BINARY)

test: ## Run all tests with the race detector
	go test -race ./...

cover: ## Run tests and fail below COVERAGE_MIN percent coverage
	go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | awk -v min=$(COVERAGE_MIN) '/^total:/ {sub("%", "", $$3); printf "coverage %s%% (minimum %s%%)\n", $$3, min; exit ($$3 + 0 < min)}'

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format the code
	golangci-lint fmt ./...

golden: ## Regenerate the testscript golden files after an intended CLI change
	go test ./cmd/... -update

vuln: ## Report known vulnerabilities reachable from the code (govulncheck)
	go run $(GOVULNCHECK) ./...

docs: ## Regenerate the CLI reference in docs/cli
	go run ./tools/gendocs docs/cli

docs-check: ## Fail if docs/cli is out of date
	go run ./tools/gendocs -check docs/cli

lint-ci: ## Lint GitHub workflows (actionlint, zizmor)
	go run $(ACTIONLINT) .github/workflows/*.yml
	zizmor --no-online-audits .github/workflows

snapshot: ## Build all release artifacts locally (binaries, archives, SBOMs, image) without publishing
	goreleaser release --snapshot --clean --skip=publish

clean: ## Remove build and coverage outputs
	rm -rf bin dist coverage.out
