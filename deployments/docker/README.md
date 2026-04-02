# GoREveal Podman Containers

GoREveal is developed container-first.

The dev image intentionally includes a small operator toolbox for current workflows:
- `jq` for structured JSON inspection of `analyze`, diff, and protected-matrix output
- `yq` for `Taskfile.yml`, CI, and YAML contract inspection
- `procps` for container/debug process inspection
- `unzip` for bounded artifact and baseline handling
- `shellcheck`, `ruff`, `ty`, `yamllint` for script verification

## Images

- `Containerfile.dev`: interactive development, lint, test, fuzz, benchmark, code generation
- `Containerfile.builder`: reproducible build stage for GoREveal binaries
- `Containerfile.release`: minimal OCI runtime image for packaged binaries

## Development Flow

Build dev image:
```bash
podman build -f deployments/docker/Containerfile.dev -t goreveal:dev .
```

Start long-running dev container:
```bash
podman run -d --name goreveal-dev -v "$PWD:/workspace:Z" -w /workspace goreveal:dev
```

Run commands inside container:
```bash
podman exec goreveal-dev /usr/local/go/bin/go test ./...
podman exec goreveal-dev golangci-lint run
podman exec goreveal-dev /usr/local/go/bin/go test -bench=. -benchmem ./...
```

The host machine is orchestration-only. Build, test, lint, fuzz, and benchmark steps should run inside Podman.

Repository shortcuts from the host:
```bash
make test
make test-differential
make test-differential-report
make test-snapshots
make snapshot-update
```

Taskfile shortcuts:
```bash
task build-image
task test
task lint
task verify
task lint-scripts
```

`make test` is the main regression entrypoint. It now mounts baseline repositories and includes the differential Go test package in the normal `go test ./...` pass before running plugin tests.
Fixture-driven Python integrations and the differential report path use a built workspace binary at `.tmp/goreveal` to avoid flaky nested `go run` calls inside the dev container.

## Python Automation

GoREveal now also includes a thin Python automation layer:
- `Taskfile.yml` provides the operator-facing task UX
- `scripts/dev/podman_runner.py` uses `podman-py` to build the dev image and run predefined tasks inside Podman

Install the Python dependency:
```bash
python3 -m pip install -e .
```

Script linting entrypoints:
```bash
make lint-scripts
task lint-scripts
```
