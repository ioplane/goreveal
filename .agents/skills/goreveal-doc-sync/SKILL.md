---
name: goreveal-doc-sync
description: Keep plans, roadmap notes, architecture docs, README, and agent configuration synchronized after strategic changes.
---

# GoREveal Doc Sync

## Use When

- a strategic review changes roadmap ordering
- sprint priorities shift
- architecture decisions are promoted from brainstorming into contract docs
- README and plan docs drift apart
- agent or skill configuration changes
- automation or verification entrypoints change

## Workflow

1. Update the source-of-truth decision doc.
2. Update sprint and roadmap notes.
3. Update architecture docs only where the decision changes a real boundary.
4. Update `README.md`.
5. Update `AGENTS.md`, `.codex/agents/`, and skill READMEs if the operator workflow changed.
6. Check that active and exploratory documents are clearly distinguished.
7. Check that referenced verification commands still exist and match `Taskfile.yml`, `Makefile`, and `scripts/dev/podman_runner.py`.

## Borrowed Best Practices

- use explicit phase or lane boundaries
- keep exploratory plans marked as exploratory
- keep operational docs separate from architecture decisions
- prefer checkable references and registries over vague prose
- keep agent and skill docs aligned with the actual repo automation surface
