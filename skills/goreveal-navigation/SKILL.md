---
name: goreveal-navigation
description: Navigate GoREveal module ownership, the active RT1 train, and preserved compatibility boundaries before making changes.
metadata:
  short-description: Navigate GoREveal safely
---

# GoREveal Navigation

## Use When

- starting any non-trivial task in this repository
- deciding whether work belongs in `core`, `schema`, `engine`, `deobfuscation`, `storage`, `plugins`, `corpus`, or docs
- checking whether a proposed change matches the active sprint and roadmap

## Read First

- `AGENTS.md`
- `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- `docs/architecture/2026-03-19-goreveal-module-map.md`
- `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- `docs/superpowers/specs/2026-07-22-goreveal-rt1-product-design.md`
- `docs/superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md`

The March/April Scrum and strategic plans are historical evidence, not active
execution authority.

## Workflow

1. Identify the owning module.
2. State whether the work changes recovery truth, schema surface, enrichment logic, export behavior, or docs only.
3. Check whether the task is inside the active RT1 sequence
   `S0 -> S1 -> S2A -> S2B -> S3` and satisfies its promotion gate.
4. Check whether the task requires corpus, snapshot, differential, or benchmark evidence.
5. Refuse to put plugin, storage, or analyst UX logic into `core`.

## Guardrails

- `schema` is the canonical contract.
- `core` recovers truth from binaries; it does not own analyst-facing classification.
- `engine` owns orchestration and future code-peeling style enrichment.
- `deobfuscation` refines readability without overwriting raw truth.
- bounded Sprint 12 runtime fields remain compatibility/evidence slices, not
  an active queue or permission for a broad runtime parser rewrite.

## Deliverable

Before implementation, be able to say:
- which module owns the change
- which docs govern it
- which evidence type must change with it
