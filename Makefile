.DEFAULT_GOAL := help

.PHONY: help fmt build-dev-bin build-image test test-differential \
        test-differential-report test-plugins test-snapshots snapshot-update \
        lint lint-python lint-yaml lint-shell lint-scripts format-python \
        fuzz bench verify protected-matrix release-check release-snapshot version

PODMAN ?= podman
DEV_IMAGE ?= localhost/goreveal:dev
DEV_WORKDIR ?= /workspace
GO ?= /usr/local/go/bin/go
GOFMT ?= /usr/local/go/bin/gofmt
GOLANGCI_LINT ?= /go/bin/golangci-lint
PYTHON ?= python3
UV ?= uv
DEV_BIN ?= $(DEV_WORKDIR)/.tmp/goreveal

# Host directory holding the reference-tool checkouts used by the differential
# suite. Override on the command line or in the environment:
#   make test-differential BASELINES_ROOT=/path/to/baselines
# See CONTRIBUTING.md for the expected directory names.
BASELINES_ROOT ?= $(if $(GOREVEAL_BASELINES_HOST_ROOT),$(GOREVEAL_BASELINES_HOST_ROOT),$(HOME)/goreveal-baselines)

DEV_RUN = $(PODMAN) run --rm -v $(CURDIR):$(DEV_WORKDIR):Z -w $(DEV_WORKDIR) $(DEV_IMAGE)
DEV_RUN_WITH_BASELINES = $(PODMAN) run --rm \
	-v $(CURDIR):$(DEV_WORKDIR):Z \
	-v $(BASELINES_ROOT):/repos:Z \
	-w $(DEV_WORKDIR) \
	-e GOREVEAL_BASELINES_ROOT=/repos \
	$(DEV_IMAGE)

help: ## Show available targets
	@printf 'GoREveal make targets\n\n'
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-26s\033[0m %s\n", $$1, $$2}'
	@printf '\nDevelopment is Podman-first; see CONTRIBUTING.md.\n'

# --- Build --------------------------------------------------------------------

build-image: ## Build the Podman dev image
	$(PYTHON) -m scripts.dev.podman_runner build-image

build-dev-bin: ## Build the workspace goreveal binary inside the dev container
	$(DEV_RUN) bash -lc 'mkdir -p .tmp && $(GO) build -o .tmp/goreveal ./cmd/goreveal'

version: ## Print the build identity of the workspace binary
	$(MAKE) build-dev-bin
	$(DEV_RUN) $(DEV_BIN) version

fmt: ## Format Go sources inside the dev container
	$(DEV_RUN) bash -c '$(GOFMT) -w $$(find . -type f -name "*.go" -not -path "./.git/*")'

# --- Test ---------------------------------------------------------------------

test: ## Run the main regression suite (Go plus Python unit tests)
	$(MAKE) build-dev-bin
	$(DEV_RUN_WITH_BASELINES) $(GO) test ./...
	$(DEV_RUN) $(PYTHON) -m unittest scripts.baseline.test_normalize
	$(DEV_RUN) $(PYTHON) -m unittest scripts.baseline.test_report
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ida.test_goreveal_ida'
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ghidra.test_goreveal_ghidra'

test-plugins: ## Run the IDA and Ghidra adapter tests
	$(MAKE) build-dev-bin
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ida.test_goreveal_ida'
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ghidra.test_goreveal_ghidra'

test-differential: ## Run differential tests against baseline tool checkouts
	$(DEV_RUN_WITH_BASELINES) $(GO) test ./tests/differential

test-differential-report: ## Generate the differential fixture report
	$(MAKE) build-dev-bin
	$(DEV_RUN_WITH_BASELINES) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m scripts.baseline.generate_fixture_report corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin'

test-snapshots: ## Run the golden snapshot suite
	$(DEV_RUN) $(GO) test ./tests/snapshots

snapshot-update: ## Refresh golden snapshots (explain every diff in the pull request)
	$(DEV_RUN) $(GO) test ./tests/snapshots -run TestAnalyzeFixtureSnapshot -update

fuzz: ## Run fuzz targets
	$(DEV_RUN) $(GO) test -fuzz=Fuzz -run=^$$ ./...

bench: ## Run benchmarks
	$(DEV_RUN) $(GO) test -bench=. -benchmem ./...

# --- Lint ---------------------------------------------------------------------

lint: ## Run golangci-lint
	$(DEV_RUN) $(GOLANGCI_LINT) run

format-python: ## Format Python sources with Ruff
	$(PYTHON) -m scripts.dev.podman_runner task format-python

lint-python: ## Run Ruff and ty
	$(PYTHON) -m scripts.dev.podman_runner task lint-python

lint-yaml: ## Run yamllint
	$(PYTHON) -m scripts.dev.podman_runner task lint-yaml

lint-shell: ## Run shellcheck
	$(PYTHON) -m scripts.dev.podman_runner task lint-shell

lint-scripts: lint-python lint-yaml lint-shell ## Run all non-Go linters

# --- Release ------------------------------------------------------------------

release-check: ## Validate the GoReleaser configuration
	$(DEV_RUN) bash -lc 'GOWORK=off goreleaser check'

release-snapshot: ## Build a full release locally without publishing or signing
	$(DEV_RUN) bash -lc 'GOWORK=off goreleaser release --snapshot --clean --skip=sign,publish'

# --- Aggregate ----------------------------------------------------------------

protected-matrix: ## Build the protected-binary profile matrix and emit a JSON report
	$(PYTHON) -m scripts.dev.podman_runner task protected-matrix

verify: fmt lint test lint-scripts test-snapshots ## Run every local gate CI runs
