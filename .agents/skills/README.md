# GoREveal Agent Skills

<img
  src="https://shieldcn.dev/badge/agents-skills-slate.svg?variant=outline&size=xs"
  alt="agents: skills" height="20">

Repo-local skills for automated contributors. Each skill is a narrow operational
workflow, not general policy — the root contract is [AGENTS.md](../../AGENTS.md).

Layout:

```text
.agents/skills/<skill-name>/SKILL.md
```

This is the single location. An earlier `skills/` compatibility mirror has been
removed; do not reintroduce it.

## Active skills

| Skill | Use when |
| --- | --- |
| `goreveal-navigation` | Orienting in the repository before any change |
| `goreveal-cleanroom` | Studying a reference tool without importing its code |
| `goreveal-corpus-validation` | Touching fixtures, golden snapshots, or recovery evidence |
| `goreveal-differential-testing` | Turning a baseline comparison into normalized evidence |
| `goreveal-deobfuscation` | Working on refinement layers over raw recovered truth |
| `goreveal-export-contracts` | Changing an export shape consumed by CLI, SQLite, or a host tool |
| `goreveal-perf-simd` | Any performance or low-level optimization work |
| `goreveal-release-ops` | Keeping release, CI, and operational claims aligned with reality |
| `goreveal-doc-sync` | Re-aligning docs, README, and agent configuration after a strategic change |

## Design principles

- Keep skills narrow and operational.
- Keep the `core` / `schema` / `engine` boundaries explicit in every workflow.
- Prefer references, checklists, and exact deliverables over generic policy text.
- Prefer validation steps that map to real repository commands, not invented ones.
- Treat [AGENTS.md](../../AGENTS.md) as the root contract: skills refine
  workflows, never project purpose or the two hard rules.
- Prefer Podman-first verification through `task …` or `make …`, backed by
  `scripts/dev/podman_runner.py`.
- When documenting automation work, name the actual strict lint surface: `ruff`,
  `ty`, `yamllint`, `shellcheck`, `golangci-lint`.

## Related configuration

| Path | Purpose |
| --- | --- |
| `.codex/agents/` | Repo-local subagent definitions |
| `.codex/config.toml` | Project-scoped Codex configuration |
| [`AGENTS.md`](../../AGENTS.md) | Root agent contract |
| [`CLAUDE.md`](../../CLAUDE.md), [`CODEX.md`](../../CODEX.md), [`GEMINI.md`](../../GEMINI.md) | Per-assistant overlays |
