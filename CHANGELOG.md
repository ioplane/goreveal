---
title: Changelog
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - release
---

# Changelog

<a href="https://keepachangelog.com/en/1.1.0/">
  <img
    src="https://shieldcn.dev/badge/changelog-Keep%20a%20Changelog-slate.svg?variant=outline&size=xs"
    alt="Keep a Changelog" height="20"></a>
<a href="https://semver.org/spec/v2.0.0.html">
  <img
    src="https://shieldcn.dev/badge/versioning-SemVer%202.0.0-slate.svg?variant=outline&size=xs"
    alt="Semantic Versioning" height="20"></a>

All notable changes to this project are documented here.

The format follows [Keep a Changelog 1.1.0](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

Per-release artifact notes, checksums, SBOMs, and signature verification steps
live with each [GitHub Release](https://github.com/ioplane/goreveal/releases) and
are described in [docs/RELEASE.md](docs/RELEASE.md).

## [Unreleased]

### Added

- Apache-2.0 licensing (`LICENSE`, `NOTICE`) with an explicit clean-room notice.
- Community health baseline: `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`,
  `SECURITY.md`, `SUPPORT.md`, `MAINTAINERS.md`, issue and pull-request
  templates, `CODEOWNERS`.
- `goreveal version [--json]` and `goreveal help` subcommands. Version identity
  falls back to the toolchain's embedded VCS stamps when link-time values are
  absent, so a `go install` build still reports truthful provenance.
- Continuous integration on GitHub Actions: Go lint and test matrix across
  Linux, macOS, and Windows; `uv`-managed Python lint and type checks (`ruff`,
  `ty`); YAML and shell linting; Markdown and link checking; workflow security
  auditing (`zizmor`, `actionlint`); container builds for all three
  Containerfiles.
- Security scanning: CodeQL (Go and Python, `security-extended`), Semgrep,
  Trivy vulnerability and secret scanning, `govulncheck`, OSV-Scanner, Gitleaks,
  SonarQube Cloud, dependency review, and OSSF Scorecard.
- Release automation via GoReleaser: reproducible cross-platform binaries,
  SPDX JSON SBOMs, SHA-256 checksums, Sigstore keyless signatures, SLSA build
  provenance attestations, Debian and RPM packages, and multi-architecture
  container images on GitHub Container Registry.
- `docs/RELEASE.md` documenting the release process and consumer-side
  verification of every published artifact.

### Changed

- Module path is now `github.com/ioplane/goreveal`, matching the canonical
  repository location.
- Python toolchain moved to `uv` with Python 3.14; `ruff`, `ty`, and `yamllint`
  are declared as a PEP 735 dependency group instead of being installed ad hoc.
- Container base images are digest-pinned and tool versions inside the dev
  container are explicit rather than `@latest`.
- `Containerfile.builder` now performs a real cross-compilable, cache-mounted,
  reproducible build (`-trimpath`, `-buildid=`) and self-verifies the output.
- Documentation restructured for public consumption: a task-oriented `README`,
  a documentation index, and uniform front matter across `docs/`.

### Removed

- Internal sprint, progress, assessment, and brainstorming documents. These were
  working artifacts for a single-maintainer development phase and are not part
  of the public project.
- The duplicated `skills/` mirror; `.agents/skills/` is the single source.

[Unreleased]: https://github.com/ioplane/goreveal/commits/main
