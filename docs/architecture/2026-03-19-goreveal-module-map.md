# GoREveal Module Map

> Status: architecture support document
> Date: 2026-03-19

## Purpose

This document fixes the intended module ownership and dependency rules for GoREveal.

## Planned Top-Level Areas

- `core`
  - owns binary parsing, low-level recovery primitives, and raw metadata extraction
- `schema`
  - owns canonical analysis data structures and provenance/confidence semantics
- `engine`
  - owns orchestration of recovery, normalization, enrichment, and export preparation
- `deobfuscation`
  - owns refinement passes over recovered data without mutating raw truth
- `cmd/goreveal`
  - owns CLI commands and presentation logic
- `storage`
  - owns SQLite persistence and analysis diffing for stored runs
- `api`
  - owns machine-facing service interfaces
- `plugins/ida`
  - owns IDA import adapter only
- `plugins/ghidra`
  - owns Ghidra import adapter only
- `bench`
  - owns performance harnesses and hotspot measurement
- `corpus`
  - owns fixtures, snapshots, fixture metadata, and baseline-oriented evidence

## Dependency Rules

Allowed dependency direction:
- `core -> schema`
- `deobfuscation -> core + schema`
- `engine -> core + deobfuscation + schema`
- `cmd/goreveal -> engine + schema`
- `storage -> schema` or `storage -> engine + schema` only where orchestration is unavoidable
- `api -> engine + schema`
- `plugins -> schema` and exported payload definitions, never direct recovery internals
- `bench -> engine/core/schema` as needed for benchmarks only

Forbidden dependency direction:
- `core -> cli`
- `core -> plugins`
- `core -> storage`
- `schema -> engine`
- `schema -> plugins`
- `plugins -> core` for recovery logic

## Ownership Rules

- Any code that parses binary internals belongs in `core`.
- Any code that defines stable result shape belongs in `schema`.
- Any code that sequences stages belongs in `engine`.
- Any code that improves readability but is not raw recovery belongs in `deobfuscation`.
- Any code that exists only to render or transport results belongs outside `core`.


## Development Environment Boundary

The repository must include OCI-compatible Podman container definitions under `deployments/docker`. Tooling, tests, linters, and benchmarks should be executed through the dev container rather than directly on the host.
