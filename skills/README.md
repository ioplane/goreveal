# GoREveal Project Skills

These skills are local workflow guides for GoREveal.

Canonical Codex-native location:
- `.agents/skills/`

Compatibility mirror:
- `skills/`

Use them to keep work aligned with:
- `AGENTS.md`
- `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- `docs/architecture/2026-03-19-goreveal-go126-best-practices.md`
- `docs/superpowers/specs/2026-07-22-goreveal-rt1-product-design.md`
- `docs/superpowers/specs/2026-07-22-goreveal-standalone-release-ida-bootstrap-design.md`
- `docs/superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md`
- `docs/superpowers/plans/2026-07-22-goreveal-standalone-release.md`

The March/April Scrum and strategic plans remain historical evidence only.

Core skills:
- `goreveal-navigation`
- `goreveal-cleanroom`
- `goreveal-corpus-validation`
- `goreveal-differential-testing`
- `goreveal-deobfuscation`
- `goreveal-perf-simd`
- `goreveal-export-contracts`
- `goreveal-release-ops`

New Codex-native support skills:
- `goreveal-doc-sync`
- `goreveal-sprint12-runtime` — compatibility specialist for the preserved
  bounded runtime surface, not an active-roadmap selector

Related repo-local configuration:
- `.codex/config.toml`
- `.codex/agents/`

Operational verification baseline:
- prefer `task ...` or `make ...` entrypoints over ad hoc host commands
- Podman automation is centralized in `scripts/dev/podman_runner.py`
- script-quality checks use:
  - `ruff`
  - `ty`
  - `yamllint`
  - `shellcheck`
