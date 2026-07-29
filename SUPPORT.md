---
title: Getting Support
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - support
  - community
---

# Getting Support

<a href="https://github.com/ioplane/goreveal/discussions">
  <img
    src="https://shieldcn.dev/badge/help-discussions-slate.svg?variant=outline&size=xs&logo=github"
    alt="Discussions" height="20"></a>
<a href="https://github.com/ioplane/goreveal/issues">
  <img
    src="https://shieldcn.dev/github/issues/ioplane/goreveal.svg?variant=outline&size=xs"
    alt="Open issues" height="20"></a>

GoREveal is maintained by volunteers. Pick the right channel and you will get a
faster, better answer.

## Where to go

| You have | Go to |
| --- | --- |
| A security vulnerability | [Private advisory](https://github.com/ioplane/goreveal/security/advisories/new) — never a public issue. See [SECURITY.md](SECURITY.md) |
| Wrong or missing recovery output | [Bug report](https://github.com/ioplane/goreveal/issues/new/choose) |
| A crash or panic on a binary | [Bug report](https://github.com/ioplane/goreveal/issues/new/choose) |
| A capability you want | [Feature request](https://github.com/ioplane/goreveal/issues/new/choose) |
| A usage question | [Discussions → Q&A](https://github.com/ioplane/goreveal/discussions) |
| An idea you want to sanity-check | [Discussions → Ideas](https://github.com/ioplane/goreveal/discussions) |
| A contribution question | [CONTRIBUTING.md](CONTRIBUTING.md), then Discussions |

## Before you ask

Most questions are answered by the documentation:

- [README](README.md) — install, CLI surface, quick start
- [docs/](docs/README.md) — architecture, schema principles, release process
- [docs/RELEASE.md](docs/RELEASE.md) — release verification, SBOM, signatures
- [CONTRIBUTING.md](CONTRIBUTING.md) — build, test, and verification workflow

## What to include in a report

Recovery bugs are only actionable with enough context to build a fixture:

```bash
goreveal version --json          # exact build identity
goreveal analyze <binary>        # canonical output, redacted as needed
```

Then describe:

- **input class** — ELF / PE / Mach-O, GOOS/GOARCH, stripped, `-trimpath`,
  obfuscated (and with what, if known)
- **what you expected** versus **what you got**
- whether a reference tool (`gore`, `redress`, `GoReSym`) produces something
  different, and what

Do not attach binaries you are not authorized to share. Describe the sample
instead, or build an equivalent reproducer — `corpus/fixtures/` shows the shape
we can act on.

## What GoREveal deliberately will not do

Understanding these avoids reports that get closed as intended behavior:

- **It will not guess.** When evidence is absent, output is `unavailable`, an
  empty list, or an explicit error. That is correct, not a bug.
- **It does not treat reference tools as ground truth.** A divergence from
  `GoReSym` or `redress` is a data point to investigate, not automatically a
  GoREveal defect.
- **It is not a sandbox.** It parses hostile input safely as a goal, but run
  unknown samples in an isolated container or VM regardless. See
  [SECURITY.md](SECURITY.md#handling-untrusted-input-safely).
- **It does not implement recovery in plugins.** `plugins/ida` and
  `plugins/ghidra` consume exported artifacts only. Missing data there almost
  always means missing data in the export contract.

## Response expectations

This is a community-maintained project with no commercial support tier. Issues
and discussions are handled on a best-effort basis. Security reports follow the
committed timelines in [SECURITY.md](SECURITY.md#response-targets).

A well-formed bug report with a reproducible fixture is the single most effective
way to get something fixed.
