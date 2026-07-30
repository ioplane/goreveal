# GoREveal Agent Contract

<img
  src="https://shieldcn.dev/badge/contract-agents-slate.svg?variant=outline&size=xs"
  alt="contract: agents" height="20">

This file is the operational contract for automated contributors and maintainers
working in GoREveal. Read it before making changes.

Human-facing documentation lives in [README.md](README.md),
[CONTRIBUTING.md](CONTRIBUTING.md), and [docs/README.md](docs/README.md). This
file does not restate them; it states the constraints an agent must not violate.

## Project purpose

GoREveal is a clean-room reverse-engineering platform for Go binaries. It is
informed by `gore`, `redress`, `GoReSym`, `GoResolver`, `gostringungarbler`, and
`AlphaGolang` — and must not copy code from any of them.

Priority order, in this order, always:

1. accuracy
2. convenience
3. speed

## Required reading before substantial work

- [docs/architecture/2026-03-19-goreveal-platform-contract.md](docs/architecture/2026-03-19-goreveal-platform-contract.md)
- [docs/architecture/2026-03-19-goreveal-module-map.md](docs/architecture/2026-03-19-goreveal-module-map.md)
- [docs/architecture/2026-03-19-goreveal-schema-principles.md](docs/architecture/2026-03-19-goreveal-schema-principles.md)
- [docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md](docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md)
- [docs/architecture/2026-03-19-goreveal-go126-best-practices.md](docs/architecture/2026-03-19-goreveal-go126-best-practices.md)
- [docs/architecture/2026-03-19-goreveal-testing-strategy.md](docs/architecture/2026-03-19-goreveal-testing-strategy.md)

## Hard rule 1: clean-room boundary

Allowed:

- studying reference repositories for behavior, edge cases, data formats, and
  test ideas
- building differential tests against baseline tools
- documenting similarities and divergences

Forbidden:

- copying implementation code from baseline projects
- translating AGPL-licensed code into superficially different Go
- treating baseline output as infallible truth without validation

When baseline behavior is useful, convert it into a documented finding, a corpus
fixture, a differential test, or a consciously designed implementation. Never
into copied code.

## Hard rule 2: never invent recovery truth

When evidence is absent the correct output is `unavailable`, an empty collection,
or an explicit error. A plausible-looking guess is a defect.

Concretely:

- degrade to `[]`, never to JSON `null`, when there is no truthful surface
- `analyze` emits `types: []` rather than omitting the field
- `inspect runtime` returns `unavailable` rather than inventing runtime facts
- nodes without file evidence are marked `has_file_evidence: false` rather than
  presented as real paths
- `provenance` and `confidence` must accurately describe how each value was
  obtained, including which fallback path produced it

## Architecture invariants

Non-negotiable:

- `schema` is the canonical contract
- `core` stays independent of CLI, storage, API, and plugin concerns
- `deobfuscation` refines and must never overwrite raw recovered truth
- provenance and confidence remain first-class result fields
- SIMD is an optimization layer, never a correctness layer
- plugins consume exports; they do not implement recovery logic

[`.golangci.yml`](.golangci.yml) enforces part of this mechanically: `depguard`
denies `unsafe` and `os/exec` outside the baseline harness. Do not weaken that
policy to make a change pass — adapt the change.

## Delivery rules

Prefer vertical slices over architecture-only work: recovery logic, its schema
mapping, CLI or export exposure when user-visible, and evidence.

Every meaningful change leaves behind at least one of:

- corpus fixture coverage
- golden snapshot coverage
- differential comparison coverage
- fuzz coverage for parsing and recovery boundaries
- benchmark evidence
- updated docs or contract notes

A change is done only when it respects module boundaries, carries the right
evidence type, updates docs if behavior or contract changed, keeps a scalar
fallback for any optimized path, and can be demonstrated through the CLI, schema
output, tests, or benchmarks.

If a change affects recovery semantics, update the golden outputs, the
differential expectations, and the provenance/confidence behavior together.
**Refreshing a golden snapshot to make a test pass, without explaining the diff,
is prohibited.** That is how a correctness regression ships silently.

## Performance policy

Strict order:

1. pure Go reference implementation
2. optimized scalar implementation
3. architecture-specific SIMD
4. optional experimentation

No SIMD work without hotspot evidence, correctness-equivalence tests, a scalar
fallback, and documented feature gating.

## Verification

All development commands run inside the Podman dev container. Do not rely on host
Go or host linters.

```bash
podman build -f deployments/docker/Containerfile.dev -t localhost/goreveal:dev .
task build-image

task fmt
task lint
task test
task lint-scripts
task test-snapshots
task test-differential
task test-differential-report
task test-plugins
```

Equivalent `make` targets exist. Keep verification **sequential** — overlapping
runs against the same bind-mounted workspace produce flaky results.

Non-Go gates are part of the supported surface and run through `uv`:

```bash
uv sync --group dev
uv run ruff check .
uv run ruff format --check .
uv run ty check
uv run yamllint --strict .
```

`scripts/dev/podman_runner.py` is the canonical Podman automation entrypoint
behind `make` and `task`. Podman socket discovery may come from
`PODMAN_BASE_URL`, `CONTAINER_HOST`, or `DOCKER_HOST` — do not assume only
`XDG_RUNTIME_DIR`-based rootless layouts.

Differential tests need baseline checkouts. Point at them with
`GOREVEAL_BASELINES_HOST_ROOT`; absent baselines cause skips, not failures.

## Repository layout

```text
core/            recovery primitives (format, ingest, buildinfo, pclntab,
                 runtime, functions, packages, types, strings)
schema/          canonical analysis contract and export encoders
engine/          pipeline orchestration, peeling, projection
deobfuscation/   garble string and name refinement
storage/         SQLite persistence and build-to-build diffing
cmd/goreveal/    operator CLI
internal/        version identity and test helpers
plugins/         thin IDA and Ghidra adapters
corpus/          fixtures and the protected-binary sample
tests/           snapshot and differential suites
scripts/         Python automation
docs/            architecture and release documentation
deployments/     Containerfiles
```

## Commits and pull requests

Conventional Commits, enforced on pull request titles by CI. Allowed types:
`feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`, `chore`,
`revert`, `security`. Scope is a package path.

```text
feat(core/pclntab): recover funcnametab offsets on stripped ELF
fix(engine/peeling): keep classification evidence on empty package sets
```

Never introduce host paths, internal hostnames, secrets, or organization-internal
identifiers. This repository is public.

## Project skills

Repo-local skills live in [`.agents/skills/`](.agents/skills/). Use the one that
matches the task:

| Skill | Use when |
| --- | --- |
| `goreveal-navigation` | Orienting before any change |
| `goreveal-cleanroom` | Studying a reference tool |
| `goreveal-corpus-validation` | Touching fixtures or golden snapshots |
| `goreveal-differential-testing` | Comparing against baseline tools |
| `goreveal-deobfuscation` | Working on refinement layers |
| `goreveal-export-contracts` | Changing an export shape |
| `goreveal-perf-simd` | Any performance work |
| `goreveal-release-ops` | Release and operational claims |
| `goreveal-doc-sync` | Keeping docs aligned after a strategic change |

Repo-local subagents live in [`.codex/agents/`](.codex/agents/); project-scoped
Codex configuration is in `.codex/config.toml`.

## Overlay files

[CLAUDE.md](CLAUDE.md), [CODEX.md](CODEX.md), and [GEMINI.md](GEMINI.md) are
per-assistant overlays. They refine emphasis; they never override this file.

## Working notes

`docs/.local/` is git-ignored and holds maintainer planning material. Read it for
context when it exists, but never cite it in code, commits, or public
documentation — it is not part of the project, and contributors will not have it.
