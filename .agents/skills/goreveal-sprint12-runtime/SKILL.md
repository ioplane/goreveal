---
name: goreveal-sprint12-runtime
description: Execute bounded Sprint 12 runtime work without drifting into generic parser claims or heuristic rewrites.
---

# GoREveal Sprint 12 Runtime

## Use When

- changing `schema.RuntimeMetadata`
- changing `core/runtime`
- exposing bounded runtime truth through CLI or exports
- selecting the next Sprint 12 slice

## Current Order

1. runtime trust/evidence summary
2. rich + stripped fixture tests
3. bounded Windows `PE` fixture checkpoint
4. only then future code-peeling work

## Hard Boundaries

- no broad `moduledata` parser
- no blind same-fixture bridge accumulation by default
- no package/type heuristic rewrite from fixture-local runtime hints
- no cross-version or cross-format claim inflation

## Acceptable Work

- compact trust summary fields
- thin CLI/runtime UX improvements over already-known truth
- bounded cross-checks on a new fixture
- docs that sharpen claim boundaries

## Required Evidence

- rich fixture coverage
- stripped fixture coverage where relevant
- explicit bounded wording in docs
