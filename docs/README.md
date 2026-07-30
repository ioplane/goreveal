---
title: Documentation Index
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - documentation
  - index
---

# Documentation

<a href="https://github.com/ioplane/goreveal">
  <img
    src="https://shieldcn.dev/badge/docs-index-slate.svg?variant=outline&size=xs"
    alt="Documentation index" height="20"></a>

Start here to find the right document.

## By task

| I want to | Read |
| --- | --- |
| Install and run the tool | [README](../README.md) |
| Verify a downloaded release | [RELEASE.md — Verifying a release](RELEASE.md#verifying-a-release) |
| Understand what GoREveal promises and refuses to promise | [Platform contract](architecture/2026-03-19-goreveal-platform-contract.md) |
| Know which package owns what | [Module map](architecture/2026-03-19-goreveal-module-map.md) |
| Change the output shape | [Schema principles](architecture/2026-03-19-goreveal-schema-principles.md) |
| Know what a recovered field is allowed to claim | [Semantic claim boundaries](architecture/2026-03-20-goreveal-semantic-claim-boundaries.md) |
| Add a test or a fixture | [Testing strategy](architecture/2026-03-19-goreveal-testing-strategy.md) |
| Write Go that fits this codebase | [Go 1.26 practices](architecture/2026-03-19-goreveal-go126-best-practices.md) |
| Compare against a reference tool | [Baseline sources](architecture/2026-03-19-goreveal-baseline-sources.md) · [Differential notes](architecture/2026-03-19-goreveal-differential-v1-notes.md) |
| Cut a release | [RELEASE.md — Cutting a release](RELEASE.md#cutting-a-release) |
| Contribute a change | [CONTRIBUTING.md](../CONTRIBUTING.md) |
| Report a vulnerability | [SECURITY.md](../SECURITY.md) |

## Architecture

Durable design documents. Each one states its own status in front matter; a
document marked `draft` records a direction, not a commitment.

| Document | What it settles |
| --- | --- |
| [Platform contract](architecture/2026-03-19-goreveal-platform-contract.md) | Product scope, architecture decisions, and the boundaries the project holds itself to |
| [Module map](architecture/2026-03-19-goreveal-module-map.md) | Package ownership and the dependency direction between layers |
| [Schema principles](architecture/2026-03-19-goreveal-schema-principles.md) | How the canonical contract is shaped, versioned, and extended |
| [Semantic claim boundaries](architecture/2026-03-20-goreveal-semantic-claim-boundaries.md) | Exactly what a recovered field may assert, and what it may not |
| [Testing strategy](architecture/2026-03-19-goreveal-testing-strategy.md) | Fixtures, golden snapshots, differential comparison, fuzzing, benchmarks |
| [Go 1.26 practices](architecture/2026-03-19-goreveal-go126-best-practices.md) | Implementation patterns, anti-patterns, and the performance policy |
| [Baseline sources](architecture/2026-03-19-goreveal-baseline-sources.md) | The reference tool set and the clean-room rules governing its use |
| [Differential validation notes](architecture/2026-03-19-goreveal-differential-v1-notes.md) | What the current differential suite actually compares |
| [Server stack decision](architecture/2026-03-31-goreveal-server-stack-decision.md) | The frozen technology choice for the deferred server mode |

## Operations

| Document | Purpose |
| --- | --- |
| [RELEASE.md](RELEASE.md) | Release process, supply-chain guarantees, and consumer-side verification |
| [deployments/docker/README.md](../deployments/docker/README.md) | The three container images and how they relate |

## Component documentation

| Component | Document |
| --- | --- |
| IDA adapter | [plugins/ida/README.md](../plugins/ida/README.md) |
| Ghidra adapter | [plugins/ghidra/README.md](../plugins/ghidra/README.md) |
| Corpus fixtures | [corpus/fixtures/README.md](../corpus/fixtures/README.md) |
| Baseline corpus | [corpus/baseline/README.md](../corpus/baseline/README.md) |
| Baseline harness | [scripts/baseline/README.md](../scripts/baseline/README.md) |
| Agent skills | [.agents/skills/README.md](../.agents/skills/README.md) |

## Conventions

Every document under `docs/` carries YAML front matter:

```yaml
---
title: Human-readable title
status: active | draft | superseded
date: YYYY-MM-DD
owners:
  - github-team-or-handle
tags:
  - topic
---
```

Architecture documents are named `YYYY-MM-DD-goreveal-<topic>.md`. The date is
the document's origin, not its last edit — that lives in `git log`. A document is
never silently rewritten to reflect a new decision: it is either updated with its
`status` and `date` changed, or marked `superseded` with a pointer to what
replaced it.

`docs/.local/` is git-ignored. Maintainers keep working notes, planning drafts,
and sprint material there; none of it is part of the published project.
