---
name: goreveal-release-ops
description: Keep release claims, corpus state, CI expectations, and operational docs aligned with actual supported capability.
---

# GoREveal Release and Operations

## Use When

- preparing releases
- updating CI or verification docs
- refreshing corpus state for release notes
- writing install or operations guidance
- changing automation entrypoints, dev container contents, or lint policy

## Required Checks

- architecture docs still match implementation
- `Taskfile.yml`, `Makefile`, and `scripts/dev/podman_runner.py` still describe the same operator surface
- corpus and snapshots are current
- differential tests are green or divergences are documented
- benchmark regressions are reviewed
- release claims do not outrun current bounded capability
- script-lint policy is explicit and reproducible:
  - `ruff`
  - `ty`
  - `yamllint`
  - `shellcheck`

## Scope

- release notes
- CI expectations
- install and operations docs
- corpus maintenance
- benchmark baseline notes
- Podman-first operator workflows
