---
name: goreveal-doc-sync
description: Keep architecture docs, README, and agent configuration synchronized after a decision changes a real boundary.
---

# GoREveal Doc Sync

## Use When

- a decision changes a documented boundary or contract
- an exploratory direction is promoted into a contract doc
- README and the architecture docs drift apart
- agent or skill configuration changes
- automation or verification entrypoints change

## Workflow

1. Update the source-of-truth decision doc, including its `status` and `date`
   front matter.
2. Update architecture docs only where the decision changes a real boundary.
3. Update `README.md` and `docs/README.md`.
4. Update `CHANGELOG.md` if the change is user-visible.
5. Update `AGENTS.md`, `.codex/agents/`, and the skill docs if the operator
   workflow changed.
6. Check that `status: active | draft | superseded` still reflects reality, and
   that a superseded doc names what replaced it.
7. Check that referenced verification commands still exist and match
   `Taskfile.yml`, `Makefile`, and `scripts/dev/podman_runner.py`.
8. Never cite anything under `docs/.local/`: it is untracked, and contributors
   will not have it.

## Borrowed Best Practices

- use explicit, checkable boundaries rather than vague prose
- keep exploratory documents marked `status: draft`
- keep operational docs separate from architecture decisions
- prefer checkable references and registries over vague prose
- keep agent and skill docs aligned with the actual repo automation surface
