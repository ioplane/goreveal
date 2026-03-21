# GoREveal Baseline Sources

> Status: working baseline inventory
> Date: 2026-03-19
> Purpose: define the external reference set used for behavior study, differential tests, fixtures, and integration expectations.

## Rules

These repositories and tools are references only.
They are not runtime dependencies of GoREveal and must not be copied into the implementation.

## Local Reference Repositories

- `gore`
  - Local path: `/opt/projects/repositories/gore`
  - Upstream family: GoRE toolkit library
  - Main use: package, function, and type recovery behavior; library-oriented API expectations

- `redress`
  - Local path: `/opt/projects/repositories/redress`
  - Upstream family: GoRE toolkit CLI
  - Main use: source projection behavior, CLI UX ideas, package/type presentation

- `GoReSym`
  - Local path: `/opt/projects/repositories/GoReSym`
  - Main use: runtime-aware recovery truth, pclntab and metadata extraction comparison

- `GoResolver`
  - Local path: `/opt/projects/repositories/GoResolver`
  - Main use: CFG-guided symbol recovery and deobfuscation comparison

- `gostringungarbler`
  - Local path: `/opt/projects/repositories/gostringungarbler`
  - Main use: garble literal deobfuscation comparison and fixture ideas

- `AlphaGolang`
  - Local path: `/opt/projects/repositories/AlphaGolang`
  - Main use: IDA integration expectations and annotation/import behavior

## Fork Policy

The following forks were created in the working namespace for long-term reference control:
- `dantte-lp/GoReSym`
- `dantte-lp/GoResolver`
- `dantte-lp/gostringungarbler`
- `dantte-lp/AlphaGolang`

## Usage Policy

Turn baseline study into one or more of:
- differential tests
- corpus fixtures
- architecture notes
- export contract notes
- documented divergences

Do not turn baseline study into copied implementation.

The current product strategy is tracked separately in:
- `docs/plans/2026-03-19-goreveal-capability-transfer-plan.md`

That document defines which external capabilities are already absorbed, which will be transferred natively next, and which remain intentionally external.

## Current Differential Coverage

The first active differential surface is intentionally small and stable:
- `GoReSym` for build info parity, module-local source-file overlap, and narrow user-function overlap
- `redress` for module-path parity through `gomod`, plus package presence and narrow source-file and user-function overlap through `source`
- `gore` for build info parity, package presence, source-file basename overlap, and normalized function-name overlap

This is a v1 comparison set, not the final baseline matrix. Future expansions should add:
- `GoResolver` for deobfuscation-oriented comparison
- `gostringungarbler` for garble string refinement comparison
