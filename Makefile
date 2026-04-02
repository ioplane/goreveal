.PHONY: fmt build-dev-bin test test-differential test-differential-report test-plugins test-snapshots snapshot-update lint fuzz bench verify task-build-image lint-python lint-yaml lint-shell lint-scripts format-python protected-matrix

PODMAN ?= podman
DEV_IMAGE ?= localhost/goreveal:dev
DEV_WORKDIR ?= /workspace
GO ?= /usr/local/go/bin/go
GOFMT ?= /usr/local/go/bin/gofmt
GOLANGCI_LINT ?= /go/bin/golangci-lint
PYTHON ?= python3
DEV_BIN ?= $(DEV_WORKDIR)/.tmp/goreveal
DEV_RUN = $(PODMAN) run --rm -v $(CURDIR):$(DEV_WORKDIR):Z -w $(DEV_WORKDIR) $(DEV_IMAGE)
DEV_RUN_WITH_BASELINES = $(PODMAN) run --rm -v $(CURDIR):$(DEV_WORKDIR):Z -v /opt/projects/repositories:/repos:Z -w $(DEV_WORKDIR) -e GOREVEAL_BASELINES_ROOT=/repos $(DEV_IMAGE)

fmt:
	$(DEV_RUN) bash -c '$(GOFMT) -w $$(find . -type f -name "*.go" -not -path "./.git/*")'

build-dev-bin:
	$(DEV_RUN) bash -lc 'mkdir -p .tmp && $(GO) build -o .tmp/goreveal ./cmd/goreveal'

test:
	$(MAKE) build-dev-bin
	$(DEV_RUN_WITH_BASELINES) $(GO) test ./...
	$(DEV_RUN) $(PYTHON) -m unittest scripts.baseline.test_normalize
	$(DEV_RUN) $(PYTHON) -m unittest scripts.baseline.test_report
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ida.test_goreveal_ida'
	$(DEV_RUN) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m unittest plugins.ghidra.test_goreveal_ghidra'

test-plugins:
	$(DEV_RUN) $(PYTHON) -m unittest plugins.ida.test_goreveal_ida
	$(DEV_RUN) $(PYTHON) -m unittest plugins.ghidra.test_goreveal_ghidra

test-differential:
	$(DEV_RUN_WITH_BASELINES) $(GO) test ./tests/differential

test-differential-report:
	$(MAKE) build-dev-bin
	$(DEV_RUN_WITH_BASELINES) bash -lc 'GOREVEAL_BIN=$(DEV_BIN) $(PYTHON) -m scripts.baseline.generate_fixture_report corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin'

test-snapshots:
	$(DEV_RUN) $(GO) test ./tests/snapshots

snapshot-update:
	$(DEV_RUN) $(GO) test ./tests/snapshots -run TestAnalyzeFixtureSnapshot -update

lint:
	$(DEV_RUN) $(GOLANGCI_LINT) run

fuzz:
	$(DEV_RUN) $(GO) test -fuzz=Fuzz -run=^$$ ./...

bench:
	$(DEV_RUN) $(GO) test -bench=. -benchmem ./...

verify: fmt test

task-build-image:
	python3 -m scripts.dev.podman_runner build-image

format-python:
	python3 -m scripts.dev.podman_runner task format-python

lint-python:
	python3 -m scripts.dev.podman_runner task lint-python

lint-yaml:
	python3 -m scripts.dev.podman_runner task lint-yaml

lint-shell:
	python3 -m scripts.dev.podman_runner task lint-shell

lint-scripts: lint-python lint-yaml lint-shell

protected-matrix:
	python3 -m scripts.dev.podman_runner task protected-matrix
