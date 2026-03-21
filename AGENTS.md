# GoREveal Agent Contract

This file is the primary operational contract for all agents and maintainers working in `GoREveal`.

## Project Purpose

`GoREveal` is a clean-room Go binary reverse-engineering platform. It is inspired by `gore`, `redress`, `GoReSym`, `GoResolver`, `gostringungarbler`, and `AlphaGolang`, but it must not copy code from them.

Primary priority order:
1. accuracy
2. convenience
3. speed

## Required Reference Docs

Read these before major work:
- `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- `docs/architecture/2026-03-19-goreveal-go126-best-practices.md`
- `docs/plans/2026-03-19-goreveal-scrum-implementation-plan.md`
- `docs/tmp/draft/go-bp.md`
- `docs/tmp/draft/simd-optimization.md`

## Clean-Room Boundary

Allowed:
- study reference repositories for behavior, edge cases, data formats, and test ideas
- build differential tests against baseline tools
- document similarities and divergences

Forbidden:
- copying implementation code from baseline projects
- translating AGPL code into slightly modified GoREveal code
- treating baseline output as infallible truth without validation

When baseline behavior is useful, turn it into:
- a documented finding
- a fixture
- a differential test
- a consciously designed GoREveal implementation

## Architecture Invariants

These rules are non-negotiable:
- `schema` is the canonical contract
- `core` must stay independent from CLI, storage, API, and plugin concerns
- `deobfuscation` must not overwrite raw recovered truth
- provenance and confidence must remain first-class result fields
- SIMD is an optimization layer, not a correctness layer
- plugins consume exports; they do not implement recovery logic

## Planning and Delivery Rules

Use capability increments, not architecture-only work.
Every meaningful change should leave behind at least one of:
- test coverage
- golden snapshot coverage
- differential comparison coverage
- benchmark evidence
- updated docs or contract notes

Prefer vertical slices:
- recovery logic
- schema mapping
- CLI/export exposure if user-visible
- tests/evidence

## Definition of Done

A task is done only if:
- it respects module boundaries
- it includes the right evidence type for the change
- docs are updated when behavior or contract changes
- scalar fallback still exists for optimized paths
- it can be demonstrated through CLI, schema output, tests, or benchmarks

A sprint increment is done only if:
- it exposes a usable or verifiable new capability
- it works on at least one golden binary
- it does not require architectural rollback for follow-up work

## Testing Policy

Required quality layers:
- corpus fixtures
- golden snapshots
- differential tests against baseline tools where relevant
- fuzz tests for parsing and recovery boundaries
- benchmarks for hotspots and any optimization work

If a change affects recovery semantics, update:
- golden outputs
- differential expectations
- provenance/confidence behavior if needed

## Performance Policy

Always follow this order:
1. pure Go reference implementation
2. optimized scalar implementation
3. architecture-specific SIMD implementation
4. optional `simd/archsimd` experimentation

No SIMD work is acceptable without:
- hotspot evidence
- correctness equivalence tests
- scalar fallback
- documented feature gating

## Repository Guidance

Current repository state is still early-stage. Until full code scaffolding exists, focus on:
- architecture docs
- plan docs
- baseline inventory
- agent docs and skills
- Podman-first development environment
- current native capability transfer in Sprint 11:
  - grouped source packages
  - `external_packages`
  - package `import_path`, `source_file_count`, `module_local`
  - type `package`, `import_path`, `source_file_count`, `module_local`, and `user_meaningful`
- current next priority after Sprint 11:
  - treat Sprint 7 as maintenance for evidence hygiene, not the main execution lane
  - advance Sprint 12 from bounded runtime evidence into the first very small semantic decode
  - use `docs/architecture/2026-03-19-goreveal-sprint12-runtime-spike-notes.md` as the initial runtime-spike reference
  - use `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md` when documenting or extending semantic-runtime claims
  - keep Sprint 12 bounded around field-specific `moduledata` cross-checks and tiny semantic steps such as the typelinks/itablinks slice headers, memory-range block, `.rodata` range, `.text` range, fixture-local typelink resolution, and the current `pcHeader` / `funcnametab` / `cutab` / `filetab` / `pctab` / `pclntable` bridges, not a broad parser
  - small user-facing Sprint 12 slices are also allowed when they only project already-known truth, such as bounded function/source metadata, string candidate absolute addresses, or thin export-layer projection of canonical function/type/string fields into RE-tool payloads
  - if thin adapters already receive stronger canonical location truth, they should consume it directly instead of recomputing or discarding it; fallback logic is fine, duplicate inference is not
  - after the current `.gopclntab` bridge chain checkpoint, do not add more blind same-fixture pcln bridges by default; prefer a second fixture or a very small runtime-to-heuristic cross-check
  - the stripped ELF checkpoint is now part of the active Sprint 12 contract:
    - `ReadMetadata()` may use bounded `.go.module` fallback as `firstmoduledata_addr` for the current stripped ELF family
    - `analysis.runtime` may expose `firstmoduledata_from_go_module_fallback` so operators can tell when that bounded stripped-fixture fallback path was used
    - `inspect functions` may expose bounded `package`, `import_path`, `module_local`, `source_file`, `source_line`, and `autogenerated` metadata derived from recovered function names plus `build_info.path` for `main`
    - `inspect runtime` may expose the bounded `analysis.runtime` contract directly and should return `unavailable` rather than inventing runtime truth when that contract is absent
    - `inspect packages` may preserve only the `main` package as module-local through `build_info.path` when source-tree evidence is absent, while external packages still expose direct `import_path` truth from function recovery; if source-tree correlation exists, packages may also expose explicit `has_source_evidence` from that already-known file-backed state
    - type metadata may preserve `main` package locality through `build_info.path` when a truthful type surface exists without source-tree evidence, while non-`main` types still expose direct `import_path` truth from parsed type packages
    - `inspect types` should degrade to `[]`, not JSON `null`, when no truthful type surface exists
    - canonical `analyze` should expose `types: []`, not omit the field, when no truthful type surface exists
    - `source-tree` may fall back to module root plus module-local and external package nodes with empty file lists when truthful file evidence is absent, but should mark those nodes with `has_file_evidence: false`
  - do not rewrite package/type heuristics from Sprint 11 until runtime-semantic work yields stable naming or scope truth beyond the current fixture-local bridge
  - use `docs/plans/2026-03-20-goreveal-functional-assessment.md` for the current product/strategy reassessment

Planned major areas:
- `core`
- `schema`
- `engine`
- `deobfuscation`
- `cli`
- `storage`
- `api`
- `plugins`
- `bench`
- `corpus`

## Expected Commands

All development commands must run inside the project Podman dev container.
Preferred flow once container scaffolding exists:
- `podman build -f deployments/docker/Containerfile.dev -t goreveal:dev .`
- `make lint`
- `make test`
- `make test-differential`
- `make test-differential-report`
- `make test-plugins`
- `make snapshot-update`
- differential and corpus-specific commands documented in project skills should also run through the dev container
- keep verification sequential; avoid overlapping build/test runs against the same bind-mounted workspace
- treat `.golangci.yml` as a staged policy imported from `gobfd`: preserve the shared rule philosophy, keep `make lint` green, and adapt future rule changes deliberately instead of silently weakening the policy

## Project Skills

Use the project skills in `skills/` whenever they match the task:
- `goreveal-navigation`
- `goreveal-cleanroom`
- `goreveal-corpus-validation`
- `goreveal-differential-testing`
- `goreveal-deobfuscation`
- `goreveal-perf-simd`
- `goreveal-export-contracts`
- `goreveal-release-ops`
