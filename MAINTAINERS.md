---
title: Maintainers
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - governance
---

# Maintainers

<a href="https://github.com/ioplane/goreveal/graphs/contributors">
  <img
    src="https://shieldcn.dev/github/contributors/ioplane/goreveal.svg?variant=outline&size=xs"
    alt="Contributors" height="20"></a>

This file records who is accountable for what. Review routing is automated in
[`.github/CODEOWNERS`](.github/CODEOWNERS); this document explains the intent
behind it.

## Current maintainers

| Maintainer | GitHub | Areas |
| --- | --- | --- |
| Pavel Lavrukhin | [@dantte-lp](https://github.com/dantte-lp) | all areas; release management |

## Area ownership

| Area | Path | Primary concern |
| --- | --- | --- |
| Recovery primitives | `core/` | correctness on real binaries; no CLI/storage/plugin coupling |
| Canonical contract | `schema/` | backward compatibility of the analysis and export contracts |
| Pipeline | `engine/` | orchestration, peeling, projection; evidence preservation |
| Refinement | `deobfuscation/` | never overwrites raw recovered truth |
| Persistence and diffing | `storage/` | schema migrations, explainable build-to-build matches |
| Operator CLI | `cmd/goreveal/` | stable output shapes; machine-readable surfaces |
| Host-tool adapters | `plugins/` | thin consumers only; no recovery logic |
| Evidence | `corpus/`, `tests/` | fixtures and golden output integrity |
| Build and release | `.github/`, `.goreleaser.yaml`, `deployments/` | supply-chain integrity, reproducibility |
| Documentation | `docs/`, root `*.md` | claim accuracy; no overstated capability |

## Decision rules

These are the standing constraints a maintainer enforces on every review. They
exist because this is a reverse-engineering tool whose value is the
trustworthiness of its claims.

1. **Clean-room boundary.** No implementation code from `gore`, `redress`,
   `GoReSym`, `GoResolver`, `gostringungarbler`, or `AlphaGolang` enters this
   repository. Reference behavior becomes findings, fixtures, and differential
   tests — never copied or translated code.
2. **No invented truth.** Absent evidence produces `unavailable`, an empty list,
   or an explicit error. A plausible guess is a defect, not a feature.
3. **Accuracy outranks performance.** When the two conflict, accuracy wins.
4. **Evidence accompanies change.** See the evidence table in
   [CONTRIBUTING.md](CONTRIBUTING.md#every-change-leaves-evidence).
5. **Bounded increments.** Small, verifiable slices over speculative parser
   rewrites.

A change that violates rule 1 or 2 is rejected regardless of how much value it
appears to add.

## Release authority

Only maintainers tag releases. The tag triggers the automated pipeline in
[`.github/workflows/release.yml`](.github/workflows/release.yml); no artifact is
built or published by hand. The full process, including the supply-chain
guarantees and how consumers verify them, is documented in
[docs/RELEASE.md](docs/RELEASE.md).

## Becoming a maintainer

Maintainership follows sustained, high-quality contribution — typically several
merged changes that demonstrate fluency with the clean-room and claim-boundary
rules above, plus review participation. Existing maintainers extend the
invitation; there is no application process.

## Escalation

| Situation | Path |
| --- | --- |
| Security vulnerability | [Private advisory](https://github.com/ioplane/goreveal/security/advisories/new) |
| Code of Conduct concern | [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) reporting channel |
| Stalled review or unclear decision | Comment on the pull request, then open a [Discussion](https://github.com/ioplane/goreveal/discussions) |
