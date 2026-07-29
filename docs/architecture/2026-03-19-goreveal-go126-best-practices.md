---
title: GoREveal Go 1.26 Best Practices
status: active
date: 2026-03-19
owners:
  - ioplane/goreveal-maintainers
tags:
  - go
  - engineering
  - performance
---

# GoREveal Go 1.26 Best Practices

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

> **Purpose.** Define implementation patterns, anti-patterns, and performance rules for GoREveal on
  Go 1.26.

## 0. Container-First Rule

All development workflows must run inside the project Podman dev container. Host-installed Go
tooling is not part of the supported development contract.

## 1. Core Engineering Patterns

Required patterns:

- small focused packages with one clear responsibility
- constructor-based wiring instead of dependency-injection frameworks
- `context.Context` as the first argument for long-running, cancelable, or I/O-heavy operations
- explicit error wrapping with `%w`
- table-driven tests for deterministic recovery behavior
- fuzz tests for parser and recovery boundaries
- benchmarks for hotspots and optimization-sensitive code
- stable exported structs only through the schema layer
- provenance and confidence as first-class fields in analysis results

Allowed modern Go patterns:

- generics only where they reduce duplication without obscuring behavior
- `log/slog` as the default logging API
- `errors.Is` and `errors.As` for operational error handling
- `internal/` to enforce package boundaries
- build tags for architecture-specific optimization layers
- pure Go baseline implementations before optimized specialization

## 2. Forbidden or Discouraged Patterns

Avoid or prohibit:

- god packages and mega-files
- hidden global mutable state
- plugin-specific logic inside `core`
- abstraction-first interfaces without a second real implementation
- mixing raw recovered truth and refined/deobfuscated truth in the same field
- optimization-driven API design
- premature pooling or unsafe zero-copy tricks without profile evidence
- SIMD-only implementations without scalar fallback
- ambiguous “best effort” results without provenance or confidence metadata

## 3. Parsing and Recovery Patterns

Use these rules:

- separate `ingest`, `recover`, `normalize`, `refine`, and `export`
- keep raw offsets/addresses and interpreted semantic values clearly distinguished
- make heuristics explicit and testable
- prefer deterministic passes over hidden side-effect pipelines
- ensure each recovery stage can be reproduced from input, options, and versioned schema rules

## 4. Performance Patterns

Optimization policy:

- measure first, optimize second
- pure Go reference implementation first
- optimized scalar implementation second
- SIMD implementation third
- use mmap or buffered scanning behind a stable abstraction
- prefer data-oriented layout in hotspots only when measurements justify it

Every performance change must include:

- before/after benchmarks
- correctness equivalence tests
- scalar fallback path
- documented CPU feature gating

## 5. Testing Patterns

Required test classes:

- corpus fixtures by binary family, Go version, format, and obfuscation profile
- snapshot tests at the schema boundary
- differential tests at the behavior boundary
- fuzz tests at parser boundaries
- benchmarks at hotspot boundaries

## 6. Repository Structure Patterns

Intended repository shape:

- `core` for recovery primitives
- `schema` for canonical data contracts
- `engine` for orchestration
- `deobfuscation` for refinement passes
- `cli`, `storage`, `api`, and `plugins` only as consumers of engine/schema
- `bench` and `corpus` as first-class engineering assets

## 7. SIMD Policy

SIMD is optional and subordinate to correctness.

Rules:

- no SIMD path without a scalar reference implementation
- no SIMD path without deterministic output equivalence tests
- no SIMD work before profiler or benchmark evidence identifies the hotspot
- any architecture-specific implementation must be isolated behind build tags or explicit runtime
  feature detection
- `simd/archsimd` experiments are allowed only after a stable scalar path exists and their behavior
  is fully benchmarked and validated
