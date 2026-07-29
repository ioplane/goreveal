---
title: GoREveal Testing Strategy
status: active
date: 2026-03-19
owners:
  - ioplane/goreveal-maintainers
tags:
  - testing
  - evidence
---

# GoREveal Testing Strategy

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

## Purpose

This document defines the expected quality layers for GoREveal.

## Quality Layers

- unit tests for isolated recovery and transformation logic
- golden snapshot tests at the schema boundary
- differential tests against baseline tools where comparison is valid
- fuzz tests for parser and recovery boundaries
- benchmarks for hotspots and optimization-sensitive code

## Corpus Strategy

Corpus coverage should grow by:

- Go version family
- binary format: ELF, PE, Mach-O
- stripping profile
- obfuscation profile
- known edge-case runtime layouts

## Differential Testing Strategy

Use baseline tools to compare overlapping capabilities only.

Allowed outcomes:

- parity
- documented divergence
- GoREveal improvement
- baseline uncertainty

Never treat differential output as a substitute for raw GoREveal tests.

## Optimization Testing Strategy

Any performance change requires:

- correctness equivalence tests
- before/after benchmarks
- fallback-path validation
- explicit note if CPU feature gating is involved

## Release-Readiness Expectations

Before calling a capability stable, it should have:

- fixture coverage
- snapshot coverage
- differential coverage where relevant
- benchmark visibility if performance-sensitive

## Current Priority Interpretation

Given the current implementation progress, the highest testing priority is now:

1. differential validation expansion and divergence accounting
2. export-contract stability and fixture-backed payload shape checks for downstream consumers
3. persisted-analysis diffing across stored runs

This project now has enough capability breadth that additional features are lower priority than
proving and hardening existing claims.

## Current Coverage Snapshot

As of the current implementation baseline:

- fixture-backed coverage exists for minimal format detection and one real Go `ELF` binary
- snapshot coverage exists for the minimal `analyze` schema
- unit coverage exists for ingest, format detection, build info, pclntab extraction, function
  recovery, package recovery, initial type recovery, and initial string recovery
- unit coverage now also exists for initial ELF runtime metadata extraction
- unit coverage now also exists for bounded typelink-count evidence
- unit coverage now also exists for bounded raw typelink samples
- unit coverage now also exists for raw typelink min/max validation
- unit coverage now also exists for raw typelink sign-distribution validation
- unit coverage now also exists for the first `firstmoduledata`/`.go.module` consistency cross-check
- unit coverage now also exists for the first bounded `moduledata` typelinks slice-header parse and
  cross-check
- unit coverage now also exists for `.itablink` section evidence and the first bounded `moduledata`
  itablinks slice-header parse and cross-check
- unit coverage now also exists for the first bounded `moduledata` memory-range block parse and
  cross-check against ELF data/bss section boundaries
- unit coverage now also exists for the first bounded `.rodata` range parse and cross-check against
  ELF section boundaries
- unit coverage now also exists for the first bounded `.text` range parse and cross-check against
  ELF section boundaries, with explicit inclusive-end semantics
- unit coverage now also exists for the first tiny semantic typelink bridge via rodata-relative
  resolution and in-range counting
- unit coverage now also exists for the first typelink semantic confidence bit proving that the
  current fixture's resolved typelinks all remain within `.rodata`
- package metadata enrichment now has direct coverage for `import_path` and `source_file_count`
  correlation from source-tree evidence
- package metadata enrichment now has direct coverage for `module_local` scope classification so
  package navigation can distinguish module-owned packages from external/runtime packages
- type recovery currently uses `DWARF` as the minimal truthful metadata source; full
  typelink/moduledata traversal remains future work
- type metadata enrichment now has direct coverage for `package` and `user_meaningful` hints so
  user-facing tooling can distinguish module-local types from runtime-heavy background noise without
  filtering raw truth
- type metadata enrichment now also has direct coverage for `import_path` and `module_local`, so
  type navigation is consistent with the richer package surface
- type metadata enrichment now also has direct coverage for `source_file_count`, so the current
  source-backed heuristic is visible in tests and user-facing output
- string recovery currently uses `ELF` data-section scanning plus fixture-backed sentinels; richer
  Go-aware string modeling remains future work
- source-tree projection currently uses `DWARF` file paths filtered by `build_info.path`, but now
  also preserves non-module source evidence through explicit `external_packages`; richer
  package-to-file reconstruction and dependency normalization remain future work
- deobfuscation currently has scaffold-level coverage proving raw and refined layers remain separate
  before real passes are added
- first deobfuscation passes now have direct coverage for synthetic function-name refinement and
  string-segment extraction
- SQLite persistence now has direct coverage for save/load round-trips and CLI export flow
- stored-analysis diffing now has direct coverage for schema-level diff summaries and CLI output
  shape
- differential coverage now exists for one real Go `ELF` fixture against normalized `GoReSym`,
  `redress`, and `gore` outputs
- current `redress` overlap now covers module path, package presence, and a narrow source-file and
  user-function surface
- the currently proven user-function overlap set includes `main.main`, `main.helperAdd`, and
  `main.helperBanner` where the corresponding baseline surface exposes them
- a machine-readable differential report is available for the current fixture and overlap set
- export-contract coverage now exists for stable `IDA` and `Ghidra` v1 payload shapes
- thin plugin-adapter coverage now exists for the initial pure-Python `IDA` import layer
- thin plugin-adapter coverage now includes fixture-driven validation against real `export ida`
  output
- thin plugin-adapter coverage now includes fixture-driven validation against real `export ghidra`
  output
- all routine verification is expected to run through the Podman dev container entrypoints

## Containerized Execution

Test, fuzz, differential, and benchmark flows should be executed inside the Podman dev container so
results are reproducible and host toolchain drift is avoided.

For routine verification from the host:

- `make test` is the main regression path and includes baseline-backed differential Go tests in its
  `go test ./...` pass
- `make test-differential` remains available as the focused differential-only entrypoint
- `make test-differential-report` emits the current machine-readable differential summary for the
  canonical fixture
- keep Podman verification sequential; overlapping `go test` or build invocations against the same
  bind-mounted workspace can cause avoidable toolchain permission races
