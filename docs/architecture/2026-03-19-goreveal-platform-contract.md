# GoREveal Platform Contract

> Status: draft architecture contract
> Date: 2026-03-19
> Purpose: capture approved product, architecture, and planning decisions early so implementation can return to a stable source of truth.

## 0. Delivery Model

GoREveal is developed container-first. Build, lint, test, fuzz, benchmark, and code-generation workflows must run inside a Podman-managed OCI-compatible container environment. The host machine is treated as orchestration-only.

Reference pattern: `gobfd/deployments/docker` with separate development, builder, and release container definitions.


## 1. Product Scope and Boundaries

`GoREveal` is a clean-room platform for analysis of Go binaries that functionally covers the strongest capabilities of `gore`, `redress`, `GoReSym`, `GoResolver`, `gostringungarbler`, and `AlphaGolang`, without copying their code or inheriting their architectural constraints. These projects are treated as behavioral references, sources of edge cases, and differential-testing baselines.

Scope for the first product line:
- recovery core
- canonical analysis schema
- CLI
- export formats
- deobfuscation pipeline
- SQLite cache/store
- server API
- IDA/Ghidra adapters

Out of scope for the first product line:
- rich TUI/web UI
- distributed analysis
- live instrumentation
- binary rewriting
- Windows kernel/debugger integration
- mandatory CGO components

Platform contract:
- accuracy first
- schema before integrations
- core independent from CLI/plugins/storage
- SIMD is an optimization layer, not a correctness layer
- baseline tools are references, not runtime dependencies

Priority order:
1. accuracy
2. convenience
3. speed

## 2. Platform Architecture and Module Map

`GoREveal` should be built as a monorepo with one product contract and a set of isolated modules around it.

Core modules:
- `core`: binary parsing and recovery of runtime metadata, functions, packages, types, strings, build information, and file layout
- `schema`: canonical data model of analysis results
- `deobfuscation`: garble-aware and string/name/CFG refinement passes over recovered data
- `engine`: orchestration layer that runs recovery, enrichment, validation, and normalization
- `cli`: command-line surface for analysis, inspection, export, diffing, and service mode
- `storage`: SQLite-first persistence for cache and analysis artifacts
- `api`: server interface for automation and integrations
- `plugins/ida`: thin IDA adapter consuming canonical outputs
- `plugins/ghidra`: thin Ghidra adapter consuming canonical outputs
- `bench`: benchmarks, hotspot microbenchmarks, SIMD experiments, performance regression checks
- `testdata` or `corpus`: golden binaries, expected snapshots, and cross-tool comparison fixtures

Dependency direction must stay strict:
- `core -> schema`
- `deobfuscation -> core + schema`
- `engine -> core + deobfuscation + schema`
- `cli/storage/api/plugins -> engine + schema`
- no reverse dependency into `core`

This separation is required so that correctness work, plugin work, storage work, and SIMD work can evolve independently without contaminating the recovery core.

## 3. Data Flow, Analysis Pipeline, and SIMD Strategy

The analysis pipeline should remain explicitly staged.

Stages:
1. `ingest`: open binary, file-format detection, mmap or stream abstraction, basic layout capture
2. `runtime recovery`: recover build info, pclntab, moduledata, typelinks, interface/type metadata, string regions, symbol hints
3. `semantic recovery`: build functions, methods, packages, source roots, compiler/runtime assumptions, higher-level recovered structures
4. `deobfuscation passes`: garble-aware string recovery, name refinement, package recovery, optional CFG-guided enrichment
5. `normalization`: map all results into canonical `schema`
6. `projection/export`: CLI tables, JSON/proto export, source-tree projection, SQLite persistence, plugin payloads
7. `comparison/validation`: differential validation against baseline tools and golden expectations

Every stage should preserve provenance and confidence metadata so users can distinguish:
- directly recovered data
- heuristically inferred data
- deobfuscation-enriched data

SIMD strategy:
- SIMD is not introduced into correctness-critical code until baseline scalar behavior is stable and benchmarked
- SIMD is allowed only in measured hotspots such as:
  - pattern scanning over large mapped sections
  - byte classification and delimiter scanning
  - hashing and fingerprinting
  - bulk string-candidate scanning
  - fast compare/filter kernels
- every optimized path must have:
  - scalar fallback
  - identical deterministic output
  - explicit benchmark coverage
  - runtime feature detection

Optimization ladder:
1. pure Go reference implementation
2. optimized scalar implementation
3. architecture-specific SIMD implementation
4. optional Go 1.26 `simd/archsimd` experiment where justified

## 4. Product Surface and First-Class User Workflows

The product should be organized around user workflows rather than around isolated utilities.

First-class workflows:
- analyze a binary and obtain a full canonical analysis result
- inspect recovered functions, types, packages, strings, build information, and compiler/runtime assumptions
- export results into IDA/Ghidra without manual glue work
- project a source-like tree from recovered metadata
- deobfuscate garble-like binaries to improve readability
- compare `GoREveal` output with baseline tools on the same input

Primary product surface:
- CLI commands:
  - `goreveal analyze`
  - `goreveal inspect functions|types|packages|strings|build`
  - `goreveal deobfuscate`
  - `goreveal export ida|ghidra|json|proto|sqlite`
  - `goreveal source-tree`
  - `goreveal diff baseline`
  - `goreveal serve`
- machine-readable outputs:
  - canonical JSON
  - protobuf
  - SQLite analysis database
- integrations:
  - IDA plugin consuming exported schema
  - Ghidra adapter consuming exported schema or service output
- developer workflows:
  - benchmark runs
  - corpus regression runs
  - differential validation runs

Non-goals for the early platform:
- plugin-specific semantics in the core
- storage as a mandatory prerequisite for analysis
- premature UI work before schema stabilization
- mixing the inspection CLI and automation API into an undefined surface

## 5. Quality Strategy, Differential Testing, and Scrum Roadmap Shape

Quality for `GoREveal` must be built around three truth layers:
- correctness against known binaries
- behavior parity or improvement versus baseline tools
- performance regression visibility

Quality strategy:
- `golden corpus`: binaries across Go versions, formats (`ELF/PE/Mach-O`), stripping modes, and obfuscation scenarios
- `snapshot tests`: canonical schema output compared against expected golden snapshots
- `differential tests`: comparisons against `GoReSym`, `gore`, `redress`, `GoResolver`, and `gostringungarbler` where applicable
- `property/fuzz tests`: parser and recovery fuzzing for layout, metadata, string extraction, and type decoding
- `benchmark gates`: hotspot and corpus-scale benchmarks so SIMD and parser changes remain measurable
- `confidence accounting`: all recovery results should carry provenance/confidence so improvements remain auditable

Scrum roadmap shape:
- `Epic 1`: Core binary recovery foundation
- `Epic 2`: Canonical schema and engine
- `Epic 3`: CLI and export surface
- `Epic 4`: Deobfuscation pipeline
- `Epic 5`: SQLite persistence and analysis DB
- `Epic 6`: Differential validation framework
- `Epic 7`: IDA/Ghidra integrations
- `Epic 8`: SIMD and performance acceleration
- `Epic 9`: Server/API surface
- `Epic 10`: Release engineering and corpus operations

Sprint policy:
- use capability sprints, not architecture-only sprints
- every sprint ends with a demonstrable increment
- correctness and evidence outrank feature breadth

## 6. Epic Breakdown and Sprint Structure

Backlog should be organized into:
- `Epic`: major capability area
- `Sprint`: demonstrable increment inside one or more epics
- `Task`: finished engineering unit with tests or benchmark evidence

Recommended epic set:
- `Epic 1: Core Recovery`
- `Epic 2: Schema + Engine`
- `Epic 3: Inspection CLI`
- `Epic 4: Source Projection`
- `Epic 5: Deobfuscation`
- `Epic 6: Storage`
- `Epic 7: Differential Validation`
- `Epic 8: Integrations`
- `Epic 9: Performance`
- `Epic 10: Service Surface`
- `Epic 11: Release and Operations`

Sprint structure should follow vertical capability slices rather than horizontal architecture-only milestones.

## 7. Definition of Done, Backlog Rules, and Planning Assumptions

Definition of done for any engineering task:
- code is integrated into the correct module boundary
- tests, snapshots, or benchmarks exist for the changed behavior
- docs or contract notes are updated if behavior or schema changes
- differential expectations are updated when recovery semantics change
- scalar fallback remains intact for optimized paths
- the change can be demonstrated as part of a sprint increment

Definition of done for any sprint increment:
- a user-visible or automation-visible capability exists through CLI, schema, export, or test flow
- the capability is proven on at least one golden binary
- follow-on work can continue without architectural rollback

Backlog rules:
- `accuracy work` outranks `feature breadth`
- `schema-changing tasks` require explicit review attention
- `plugin tasks` come only after stable export contracts
- `SIMD tasks` require profiler or benchmark evidence first
- `deobfuscation tasks` must add refined layers without mutating raw recovered truth
- `baseline parity` means behavior comparison, not code copying

Planning assumptions:
- team size is small, so tasks must stay narrow and vertical
- every major feature must include recovery logic, schema mapping, tests, and user-facing exposure if relevant
- every performance initiative must include measurement, scalar baseline, optimized path, and regression benchmark
- every integration initiative must include stable export contract and fixture-based validation

## 8. Recommended First Scrum Roadmap

Recommended sprint sequence:
- `Sprint 1`: analysis skeleton
- `Sprint 2`: function and package recovery
- `Sprint 3`: types and strings
- `Sprint 4`: source projection v1
- `Sprint 5`: deobfuscation v1
- `Sprint 6`: SQLite analysis store
- `Sprint 7`: plugin-ready exports
- `Sprint 8`: performance and SIMD foundation
- `Sprint 9`: service/API layer
- `Sprint 10`: hardening and release baseline

These sprints are defined as capability increments, not date-bound ceremonies.

## 9. Agent Docs and Mandatory Project Skills

`AGENTS.md` is the primary operational contract for agents and maintainers. It should cover:
- project purpose and clean-room boundary
- baseline references
- priorities and architecture invariants
- test rules
- backlog rules
- definition of done
- standard commands and repository map once scaffolded

Agent overlays should remain short:
- `CLAUDE.md`: planning and large-scale refactor guidance
- `CODEX.md`: implementation-heavy guidance with emphasis on tests, schema boundaries, and performance discipline
- `GEMINI.md`: research and baseline-comparison guidance

Mandatory project skills:
- `goreveal-navigation`
- `goreveal-cleanroom`
- `goreveal-corpus-validation`
- `goreveal-differential-testing`
- `goreveal-deobfuscation`
- `goreveal-perf-simd`
- `goreveal-export-contracts`
- `goreveal-release-ops`

## 10. Required Documentation and Planning Artifacts Before Implementation

Required artifacts before implementation begins:
- this platform contract
- module map document
- schema principles document
- testing strategy document
- Scrum implementation plan
- `AGENTS.md`, `CLAUDE.md`, `CODEX.md`, `GEMINI.md`
- required project skills

Planning order:
1. finalize architecture docs
2. write Scrum implementation plan
3. write agent docs and project skills
4. acquire or fork missing baseline repositories
5. scaffold code

## 11. Go 1.26 Best Practices and Project Patterns

Go 1.26 best practices are normative for this project and should be documented separately in `docs/architecture/2026-03-19-goreveal-go126-best-practices.md`.

High-level policy:
- small focused packages
- constructor-based wiring
- context-first APIs where cancellation and I/O matter
- explicit error wrapping
- table tests, fuzz tests, and benchmarks
- iterators where streaming avoids unnecessary allocations
- provenance/confidence as first-class schema concerns
- measure before optimizing
- scalar before SIMD
- no mutation of raw truth inside deobfuscation layers
