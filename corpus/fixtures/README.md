# Fixtures

<img src="https://shieldcn.dev/badge/corpus-fixtures-slate.svg?variant=outline&size=xs" alt="corpus: fixtures" height="20">

This directory stores canonical corpus fixtures used by `GoREveal` snapshot,
golden, and differential tests.

Rules:
- Keep fixture metadata next to the binary under a dedicated fixture directory.
- Store expected canonical schema output as `expected.analysis.json`.
- Prefer the smallest binary that still exercises the intended behavior.
- Do not overwrite raw fixture provenance in expected outputs.
- Prefer adding fixture variants when the goal is validation breadth rather than new user-facing behavior.
- Snapshot tests auto-discover any fixture directory that contains `expected.analysis.json`.
- Current canonical fixture families now include `ELF`, bounded `PE`, and a first bounded `Mach-O` foothold.

Snapshot workflow:
- Run `make test-snapshots` from the repository root to verify current fixtures.
- Run `make snapshot-update` when the canonical schema intentionally changes and the
  updated snapshot has been reviewed.
