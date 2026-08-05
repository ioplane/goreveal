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

## [0.1.4](https://github.com/ioplane/goreveal/compare/v0.1.3...v0.1.4) (2026-08-05)


### Bug fixes

* **readme:** point the CI badge at the ci.yml workflow file ([e8fdda8](https://github.com/ioplane/goreveal/commit/e8fdda8bd915b9cb91023331576b98bb99770df3))

## [0.1.3](https://github.com/ioplane/goreveal/compare/v0.1.2...v0.1.3) (2026-07-30)


### Bug fixes

* **release:** publish only on the release event, not also on the tag push ([e467f9d](https://github.com/ioplane/goreveal/commit/e467f9ddf3fcd19cb8615a02b75089456c946f1d))

## [0.1.2](https://github.com/ioplane/goreveal/compare/v0.1.1...v0.1.2) (2026-07-30)


### Bug fixes

* **release:** pin cosign to v2 so sign-blob keeps the .sig/.pem format ([331aab1](https://github.com/ioplane/goreveal/commit/331aab14c79995d3b7807a08f56114132b6b6908))

## [0.1.1](https://github.com/ioplane/goreveal/compare/v0.1.0...v0.1.1) (2026-07-30)


### Features

* prepare repository for public enterprise release ([5fcddfc](https://github.com/ioplane/goreveal/commit/5fcddfc2ae900a7a400619c202b3ee945fb50053))


### Bug fixes

* **ci:** skip the conventional-title check for bot-authored pull requests ([372c970](https://github.com/ioplane/goreveal/commit/372c97095940682c0d7335537e831390d90642d8))
* **release:** ignore the Blacksmith buildkitd.toml so goreleaser sees a clean tree ([701ffb9](https://github.com/ioplane/goreveal/commit/701ffb9afb181c7ed815439c64199b0dc0d53240))
* **sonar:** run analysis instead of always skipping, refresh README badges ([3c7541e](https://github.com/ioplane/goreveal/commit/3c7541e6bc8bd80a2920e34bea1810d00c243e0b))


### Documentation

* align readme with gore-docs standard ([a2e05a2](https://github.com/ioplane/goreveal/commit/a2e05a2ed3fbb9047cf70b5f209c6bfc5c135003))
* **architecture:** note on Go pclntab-guided disassembler enrichment ([b82f091](https://github.com/ioplane/goreveal/commit/b82f091a4fd5d142d72f6a9cc84f843e8176c017))
* **readme:** keep OpenSSF Scorecard badge, drop the not-yet-computed Sonar gate ([c4f4954](https://github.com/ioplane/goreveal/commit/c4f4954db2e7ef127147319bfaccdddbf19669de))
* **readme:** unify OpenSSF badge style, enlarge badges, shrink header logo ([7fc88e6](https://github.com/ioplane/goreveal/commit/7fc88e64bdd7884b38097eb57f947da4babbf0ca))
* rename architecture docs to an ADR-numbered scheme, trim README ([e1e344f](https://github.com/ioplane/goreveal/commit/e1e344fba121371a4ef41e7ff136f60410da398a))


### Chores

* re-cut the 0.1.0 release on fixed main ([a410f3f](https://github.com/ioplane/goreveal/commit/a410f3f1a4c824a0b36d7db0451445d0d3bdd8bf))
* release 0.1.1 as the first published version ([d98d579](https://github.com/ioplane/goreveal/commit/d98d579182d314abad8b9d9314697374704af8f8))
* set the first release version to 0.1.0 ([37f9758](https://github.com/ioplane/goreveal/commit/37f9758d56b970f6763951454fb18748fce9bc82))

## [0.1.0](https://github.com/ioplane/goreveal/compare/v0.1.0...v0.1.0) (2026-07-30)


### Features

* prepare repository for public enterprise release ([5fcddfc](https://github.com/ioplane/goreveal/commit/5fcddfc2ae900a7a400619c202b3ee945fb50053))


### Bug fixes

* **ci:** skip the conventional-title check for bot-authored pull requests ([372c970](https://github.com/ioplane/goreveal/commit/372c97095940682c0d7335537e831390d90642d8))
* **release:** ignore the Blacksmith buildkitd.toml so goreleaser sees a clean tree ([701ffb9](https://github.com/ioplane/goreveal/commit/701ffb9afb181c7ed815439c64199b0dc0d53240))
* **sonar:** run analysis instead of always skipping, refresh README badges ([3c7541e](https://github.com/ioplane/goreveal/commit/3c7541e6bc8bd80a2920e34bea1810d00c243e0b))


### Documentation

* align readme with gore-docs standard ([a2e05a2](https://github.com/ioplane/goreveal/commit/a2e05a2ed3fbb9047cf70b5f209c6bfc5c135003))
* **architecture:** note on Go pclntab-guided disassembler enrichment ([b82f091](https://github.com/ioplane/goreveal/commit/b82f091a4fd5d142d72f6a9cc84f843e8176c017))
* **readme:** keep OpenSSF Scorecard badge, drop the not-yet-computed Sonar gate ([c4f4954](https://github.com/ioplane/goreveal/commit/c4f4954db2e7ef127147319bfaccdddbf19669de))
* **readme:** unify OpenSSF badge style, enlarge badges, shrink header logo ([7fc88e6](https://github.com/ioplane/goreveal/commit/7fc88e64bdd7884b38097eb57f947da4babbf0ca))
* rename architecture docs to an ADR-numbered scheme, trim README ([e1e344f](https://github.com/ioplane/goreveal/commit/e1e344fba121371a4ef41e7ff136f60410da398a))


### Chores

* re-cut the 0.1.0 release on fixed main ([a410f3f](https://github.com/ioplane/goreveal/commit/a410f3f1a4c824a0b36d7db0451445d0d3bdd8bf))
* set the first release version to 0.1.0 ([37f9758](https://github.com/ioplane/goreveal/commit/37f9758d56b970f6763951454fb18748fce9bc82))

## 0.1.0 (2026-07-30)


### Features

* prepare repository for public enterprise release ([5fcddfc](https://github.com/ioplane/goreveal/commit/5fcddfc2ae900a7a400619c202b3ee945fb50053))


### Bug fixes

* **ci:** skip the conventional-title check for bot-authored pull requests ([372c970](https://github.com/ioplane/goreveal/commit/372c97095940682c0d7335537e831390d90642d8))
* **sonar:** run analysis instead of always skipping, refresh README badges ([3c7541e](https://github.com/ioplane/goreveal/commit/3c7541e6bc8bd80a2920e34bea1810d00c243e0b))


### Documentation

* align readme with gore-docs standard ([a2e05a2](https://github.com/ioplane/goreveal/commit/a2e05a2ed3fbb9047cf70b5f209c6bfc5c135003))
* **architecture:** note on Go pclntab-guided disassembler enrichment ([b82f091](https://github.com/ioplane/goreveal/commit/b82f091a4fd5d142d72f6a9cc84f843e8176c017))
* **readme:** keep OpenSSF Scorecard badge, drop the not-yet-computed Sonar gate ([c4f4954](https://github.com/ioplane/goreveal/commit/c4f4954db2e7ef127147319bfaccdddbf19669de))
* **readme:** unify OpenSSF badge style, enlarge badges, shrink header logo ([7fc88e6](https://github.com/ioplane/goreveal/commit/7fc88e64bdd7884b38097eb57f947da4babbf0ca))
* rename architecture docs to an ADR-numbered scheme, trim README ([e1e344f](https://github.com/ioplane/goreveal/commit/e1e344fba121371a4ef41e7ff136f60410da398a))


### Chores

* set the first release version to 0.1.0 ([37f9758](https://github.com/ioplane/goreveal/commit/37f9758d56b970f6763951454fb18748fce9bc82))

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
- Security scanning: CodeQL (Go and Python, `security-extended`), Trivy
  vulnerability and secret scanning, `govulncheck`, OSV-Scanner, Gitleaks,
  SonarQube Cloud, dependency review, and OSSF Scorecard. Semgrep runs through
  Managed Scanning on the Semgrep AppSec Platform rather than a CI workflow.
- Release automation via GoReleaser: reproducible cross-platform binaries,
  SPDX JSON SBOMs, SHA-256 checksums, Sigstore keyless signatures, SLSA build
  provenance attestations, Debian and RPM packages, and multi-architecture
  container images on GitHub Container Registry.
- `docs/RELEASE.md` documenting the release process and consumer-side
  verification of every published artifact.
- `.github/actions/code-scanning-available`, a composite action that detects
  whether SARIF can be uploaded. Code scanning needs a public repository or
  GitHub Advanced Security; where neither holds, scans still run and still gate
  the build and only the dashboard upload is skipped, with an explicit notice.

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
- `Containerfile.release` addresses its binary through `TARGETPLATFORM`, matching
  the per-platform layout of GoReleaser's `dockers_v2` build context, so the
  image CI smoke-tests is the image releases publish.
- Secret scanning runs the Apache-2.0 `gitleaks` CLI directly rather than the
  wrapper action, which requires a paid key for organisation-owned repositories.
- Link checking splits into a strict offline pass that gates the build and an
  informational online pass that does not, so a third-party host resetting a
  connection cannot fail an unrelated change.

### Removed

- Internal sprint, progress, assessment, and brainstorming documents. These were
  working artifacts for a single-maintainer development phase and are not part
  of the public project.
- The duplicated `skills/` mirror; `.agents/skills/` is the single source.

[Unreleased]: https://github.com/ioplane/goreveal/commits/main
