---
title: GoREveal Baseline Sources
status: active
date: 2026-03-19
owners:
  - ioplane/goreveal-maintainers
tags:
  - baselines
  - clean-room
  - differential
---

# GoREveal Baseline Sources

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

> **Purpose.** Define the external reference set used for behavior study, differential tests,
  fixtures, and integration expectations.

## Rules

These repositories and tools are references only.
They are not runtime dependencies of GoREveal and must not be copied into the implementation.

## Reference Repositories

Check these out anywhere on disk and point the differential harness at the parent
directory with `GOREVEAL_BASELINES_HOST_ROOT`. Directory names must match the
first column. Absent baselines cause the differential suite to skip, not fail.

| Directory | Upstream | What GoREveal studies in it |
| --- | --- | --- |
| `gore` | [goretk/gore](https://github.com/goretk/gore) | Package, function, and type recovery behavior; library-oriented API expectations |
| `redress` | [goretk/redress](https://github.com/goretk/redress) | Source projection behavior, CLI presentation of packages and types |
| `GoReSym` | [mandiant/GoReSym](https://github.com/mandiant/GoReSym) | Runtime-aware recovery, `pclntab` and metadata extraction comparison |
| `GoResolver` | [volexity/GoResolver](https://github.com/volexity/GoResolver) | CFG-guided symbol recovery and deobfuscation comparison |
| `gostringungarbler` | [mandiant/gostringungarbler](https://github.com/mandiant/gostringungarbler) | `garble` literal deobfuscation comparison and fixture ideas |
| `AlphaGolang` | [SentineLabs/AlphaGolang](https://github.com/SentineLabs/AlphaGolang) | IDA integration expectations, annotation and import behavior |

Several of these are AGPL-licensed. That is precisely why the clean-room rule
above is a licensing constraint rather than a stylistic preference: reading them
to understand a binary format is fine, and importing their code is not.

## Usage Policy

Turn baseline study into one or more of:

- differential tests
- corpus fixtures
- architecture notes
- export contract notes
- documented divergences

Do not turn baseline study into copied implementation.

Which external capabilities are already absorbed, which are planned for native
transfer, and which remain intentionally external is tracked in the maintainers'
working notes rather than in this document.

## Current Differential Coverage

The first active differential surface is intentionally small and stable:

- `GoReSym` for build info parity, module-local source-file overlap, and narrow user-function
  overlap
- `redress` for module-path parity through `gomod`, plus package presence and narrow source-file and
  user-function overlap through `source`
- `gore` for build info parity, package presence, source-file basename overlap, and normalized
  function-name overlap

This is a v1 comparison set, not the final baseline matrix. Future expansions should add:

- `GoResolver` for deobfuscation-oriented comparison
- `gostringungarbler` for garble string refinement comparison
