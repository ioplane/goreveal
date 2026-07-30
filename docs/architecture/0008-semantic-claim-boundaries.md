---
title: GoREveal Semantic Claim Boundaries
status: active
date: 2026-03-20
owners:
  - ioplane/goreveal-maintainers
tags:
  - claims
  - runtime
---

# GoREveal Semantic Claim Boundaries

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

> **Purpose.** State exactly what the bounded semantic-runtime slices do and do
  not justify, so differential-evidence work stays in maintenance and the new
  fields are not over-claimed.

## Why This Note Exists

Runtime recovery has moved beyond bounded section evidence into the first tiny
semantic steps.

That is good. It also creates a new product risk:

- bounded fixture-local semantic truth can be mistaken for broad support

This note exists to prevent that.

## Current Semantic Claims We Can Make

On the canonical ELF fixture only, `GoREveal` now has:

- bounded rodata-relative typelink resolution
- a bounded count of resolved typelinks that fall within `.rodata`
- a confidence bit showing that all currently observed typelinks resolve within `.rodata`
- a bounded `types..etypes` bridge with an explicit `moduledata_types_range_word_index`
- a bounded count showing that currently resolved typelinks fall within that fixture-local
  `types..etypes` range
- a bounded `pcHeader` / `funcnametab` bridge showing that the first observed `firstmoduledata`
  words align with `.gopclntab` and a fixture-local `funcnametab` slice
- a bounded `cutab` bridge showing that the next observed `firstmoduledata` words also form a
  fixture-local slice within `.gopclntab`
- a bounded `filetab` bridge showing that the next observed `firstmoduledata` words also form a
  fixture-local slice within `.gopclntab`
- a bounded `pctab` bridge showing that the next observed `firstmoduledata` words also form a
  fixture-local slice within `.gopclntab`
- a bounded `pclntable` bridge showing that the next observed `firstmoduledata` words also form a
  fixture-local slice within `.gopclntab`

These are legitimate claims because they are:

- directly tested
- fixture-backed
- bounded in scope
- represented in canonical schema

On the stripped ELF fixture only, `GoREveal` now also has:

- bounded `.go.module`-based fallback for `firstmoduledata_addr` in the current ELF family
- bounded `main` package preservation through `build_info.path` when source-tree evidence is absent
- bounded empty-array type degradation (`[]` instead of JSON `null`) when no truthful `DWARF` type
  surface exists
- bounded raw ELF `pclntab` header evidence showing the observed `magic`, `quantum`, `pointer_size`,
  and whether the magic is one of the currently recognised standard values

On the current Windows `PE` fixture only, `GoREveal` now also has:

- bounded `.text` / `.rdata` range evidence
- a raw `.rdata` `pclntab` magic candidate count and first-hit address
- one header-looking `.rdata` `pclntab` candidate with raw `magic`, `quantum`, and `pointer_size`

## Claims We Cannot Make Yet

We cannot honestly claim:

- general typelink decoding support
- typelinks-driven type recovery
- general `pcHeader` decoding support
- general `funcnametab` decoding support
- general `cutab` decoding support
- general `filetab` decoding support
- general `pctab` decoding support
- general `pclntable` decoding support
- cross-version typelink-base stability
- cross-format typelink-base stability
- package/type heuristic replacement from runtime truth
- `GoReSym` parity on runtime semantics
- runtime-driven type recovery on stripped binaries
- broad package/type heuristic replacement on stripped binaries
- general `PE` `pcHeader` decoding support
- general `PE` `pclntab` decoding support
- `PE` `moduledata` recovery support
- general support for recovering functions from binaries whose ELF `pclntab` header uses a
  non-standard custom magic and encrypted `entryOff`

None of those claims follow from the current semantic slices.

In particular, the current `moduledata_types_range_word_index` claim means only:

- the canonical fixture exposes a stable word position for the current bounded `types..etypes` range
  hypothesis
- not that `types..etypes` has been generally decoded across Go versions or formats

The current `moduledata_pcheader_*` and `moduledata_funcnametab_*` claims mean only:

- the canonical fixture exposes a stable bounded bridge from `firstmoduledata` into the current
  `.gopclntab` layout
- not that `pcHeader`, `funcnametab`, or the broader pcln table have been generally decoded across
  Go versions or formats

The current `moduledata_cutab_*` claims mean only:

- the canonical fixture exposes the next bounded slice-shaped bridge inside the current `.gopclntab`
  layout
- not that `cutab` has been generally decoded across Go versions or formats

The current `moduledata_filetab_*` claims mean only:

- the canonical fixture exposes the next bounded slice-shaped bridge inside the current `.gopclntab`
  layout
- not that `filetab` has been generally decoded across Go versions or formats

The current `moduledata_pctab_*` claims mean only:

- the canonical fixture exposes the next bounded slice-shaped bridge inside the current `.gopclntab`
  layout
- not that `pctab` has been generally decoded across Go versions or formats

The current `moduledata_pclntable_*` claims mean only:

- the canonical fixture exposes the next bounded slice-shaped bridge inside the current `.gopclntab`
  layout
- not that `pclntable` has been generally decoded across Go versions or formats

## Product Boundary

Current semantic-runtime fields are:

- evidence-backed
- product-useful
- intentionally non-generic

They should be treated as:

- fixture-local runtime truth
- bounded stepping stones toward stronger semantic recovery

They should not be treated as:

- marketing claims
- compatibility claims
- proof that current runtime semantics are solved

## Differential-Evidence Maintenance Rule

When the differential-evidence docs or fixtures are updated:

- allow the new semantic-runtime fields to appear in status notes
- do not convert them into parity claims against `GoReSym`
- do not widen them into version-family support claims
- do not imply that package/type heuristics are now runtime-driven

## Practical Rule for Future Tasks

Any future semantic-runtime task must state one of:

- “fixture-local and bounded”
- “validated across more than one fixture”
- “validated across a version family”

If it cannot say one of those explicitly, it is not ready for broader claims.
