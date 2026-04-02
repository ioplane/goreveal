# GoREveal Agent Skills

This directory is the Codex-native home for repo-local skills.

Canonical location:
- `.agents/skills/<skill-name>/SKILL.md`

Compatibility location:
- `skills/<skill-name>/SKILL.md`

For now, `skills/` remains as the legacy mirror used by the current repository workflow, while `.agents/skills/` is the preferred layout for Codex CLI portability.

## Active Skills

- `goreveal-navigation`
- `goreveal-cleanroom`
- `goreveal-corpus-validation`
- `goreveal-differential-testing`
- `goreveal-deobfuscation`
- `goreveal-export-contracts`
- `goreveal-perf-simd`
- `goreveal-release-ops`
- `goreveal-doc-sync`
- `goreveal-sprint12-runtime`

## Design Principles

- Keep skills narrow and operational.
- Keep `core` / `schema` / `engine` boundaries explicit.
- Prefer references, checklists, and exact deliverables over generic policy text.
- Prefer validation steps that map to actual repo commands.
- Treat `AGENTS.md` as the root contract; skills refine workflows, not project purpose.
- Prefer Podman-first verification through `task ...` or `make ...`, backed by `scripts/dev/podman_runner.py`.
- When documenting script or automation work, include the supported strict lint surface:
  - `ruff`
  - `ty`
  - `yamllint`
  - `shellcheck`
