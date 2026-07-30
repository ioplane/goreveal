# Codex Notes for GoREveal

<img
  src="https://shieldcn.dev/badge/overlay-codex-slate.svg?variant=outline&size=xs"
  alt="overlay: codex" height="20">

Read [`AGENTS.md`](AGENTS.md) first. This file is an implementation overlay and
never overrides the two hard rules there.

## Focus areas

- narrow vertical slices: recovery logic, schema mapping, operator surface, evidence
- test-backed changes
- schema-safe implementation
- benchmark-backed performance work

## When working here

- implement recovery logic together with its schema mapping and its evidence
- do not introduce an interface before a second real implementation exists
- do not optimize before a hotspot is measured
- keep scalar and optimized paths behaviorally identical
- absent evidence produces `unavailable`, `[]`, or an explicit error — never an
  inferred value

## Container rule

Do not rely on host Go or host linters. Run build, test, lint, fuzz, and bench
inside the Podman dev container, which pins every tool version:

```bash
task lint
task test
task lint-scripts
task test-snapshots
```

## Before opening a pull request

- `golangci-lint` reports zero issues; do not weaken `.golangci.yml` to get there
- a `//nolint` needs a specific linter and a real explanation — `nolintlint`
  enforces both, and a vague reason is a review failure
- golden snapshots have not drifted, or the diff is explained in the description
