---
name: goreveal-navigation
description: Navigate GoREveal module ownership and claim boundaries before making changes.
---

# GoREveal Navigation

## Use When

- starting any non-trivial task in this repository
- deciding whether work belongs in `core`, `schema`, `engine`, `deobfuscation`, `storage`,
  `plugins`, `corpus`, or docs
- checking whether a proposed change respects the documented claim boundaries

## Read First

- `AGENTS.md`
- `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- `docs/architecture/2026-03-19-goreveal-module-map.md`
- `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- `docs/architecture/2026-03-19-goreveal-schema-principles.md`
- `CONTRIBUTING.md`

## Workflow

1. Identify the owning module.
2. State whether the work changes recovery truth, schema surface, enrichment logic, export behavior,
   or docs only.
3. Check whether the task stays inside the documented claim boundaries.
4. Check whether the task touches repo automation, agent files, skills, or verification entrypoints.
5. Check whether the task requires corpus, snapshot, differential, benchmark, or script-lint
   evidence.
6. Refuse to put plugin, storage, or analyst UX logic into `core`.

## Guardrails

- `schema` is the canonical contract.
- `core` recovers truth from binaries; it does not own analyst-facing classification.
- `engine` owns orchestration and future code-peeling style enrichment.
- `deobfuscation` refines readability without overwriting raw truth.
- bounded runtime fields are contract slices, not permission for a broad runtime
  parser rewrite.
- repo-local operator commands should stay aligned across `Taskfile.yml`, `Makefile`, and
  `scripts/dev/podman_runner.py`.

## Deliverable

Before implementation, be able to say:

- which module owns the change
- which docs govern it
- which evidence type must change with it
