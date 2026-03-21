# GoREveal Podman Containers

GoREveal is developed container-first.

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

`make test` is the main regression entrypoint. It now mounts baseline repositories and includes the differential Go test package in the normal `go test ./...` pass before running plugin tests.
Fixture-driven Python integrations and the differential report path use a built workspace binary at `.tmp/goreveal` to avoid flaky nested `go run` calls inside the dev container.
