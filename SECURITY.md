---
title: Security Policy
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - security
  - policy
---

# Security Policy

<!-- shieldcn: subtle document metadata strip -->
<a href="https://github.com/ioplane/goreveal/security/advisories">
  <img
    src="https://shieldcn.dev/badge/policy-security-slate.svg?variant=outline&size=xs&logo=github"
    alt="Security policy" height="20"></a>
<a href="https://github.com/ioplane/goreveal/blob/main/docs/RELEASE.md">
  <img
    src="https://shieldcn.dev/badge/artifacts-signed-slate.svg?variant=outline&size=xs&logo=sigstore"
    alt="Signed artifacts" height="20"></a>

## Supported versions

GoREveal follows semantic versioning. Security fixes land on the latest minor
release line; older lines receive fixes only for critical issues.

| Version | Supported |
| --- | --- |
| `main` | Yes — development branch, fixes land here first |
| Latest `v0.x` release | Yes |
| Older `v0.x` releases | No — upgrade to the latest release |

While the project is pre-1.0, only the most recent tagged release is supported.

## Reporting a vulnerability

**Do not open a public issue for security problems.**

Report privately through
[GitHub Security Advisories](https://github.com/ioplane/goreveal/security/advisories/new).
That channel is monitored by the maintainers and keeps the report confidential
until a fix is published.

Please include:

- affected version or commit SHA
- the input artifact class (ELF / PE / Mach-O, stripped, obfuscated) — a minimal
  reproducer binary is ideal, but describe it rather than attach it if the sample
  is sensitive or not yours to share
- the command line and observed versus expected behavior
- impact assessment as you see it

### Response targets

| Stage | Target |
| --- | --- |
| Acknowledgement | 3 business days |
| Triage and severity assignment | 10 business days |
| Fix or documented mitigation | 90 days from acknowledgement |

We coordinate disclosure with the reporter and credit reporters in the advisory
unless anonymity is requested.

## Threat model

GoREveal parses **untrusted, potentially hostile binaries**. That is its purpose,
and it shapes what counts as a vulnerability.

### In scope

- memory-safety or panic conditions reachable from a crafted input artifact
  (out-of-range slice, unbounded allocation, infinite loop, stack exhaustion)
- unbounded resource consumption from an attacker-controlled size or count field
- path traversal or arbitrary file write through export, SQLite, or plugin paths
- SQL injection in the SQLite storage layer
- code execution triggered by analyzing an artifact
- secret or host-path leakage into analysis output, exports, or SBOMs
- supply-chain integrity gaps in the release pipeline (unsigned artifacts,
  reproducibility breaks, dependency confusion)

### Out of scope

- **Incorrect recovery output.** A wrong function name, missed package, or
  low-confidence classification is a correctness bug, not a vulnerability.
  Open a normal issue.
- **Refusal to analyze.** Returning `unavailable` or an explicit error rather
  than inventing truth is intended behavior, per the project's claim-boundary
  rules.
- Findings in the reference tools GoREveal is compared against.
- Vulnerabilities in IDA, Ghidra, or other host RE tools. The plugins under
  `plugins/` are thin consumers of exported artifacts; report host-tool issues
  to their vendors.
- Results from running GoREveal against artifacts you are not authorized to
  analyze.

## Handling untrusted input safely

GoREveal is a defensive analysis tool, but it is not a sandbox. When analyzing
samples you do not trust:

- run it in a container or VM with no network egress
- mount the sample read-only
- do not run the resulting exports through a host RE tool on a production
  workstation

The `deployments/docker/Containerfile.release` image runs as a non-root user and
contains no shell tooling beyond CA certificates for exactly this reason.

## Supply-chain integrity

Every tagged release publishes:

- SHA-256 checksums for all artifacts
- a Software Bill of Materials (SPDX JSON, produced by Syft)
- a Sigstore keyless signature over the checksum file
- SLSA build provenance attestations

Verification instructions are in [docs/RELEASE.md](docs/RELEASE.md). Verify
before deploying. If verification fails, treat the artifact as compromised and
report it through the advisory channel above.

## Security tooling in CI

Every pull request and the `main` branch are scanned by:

| Tool | Coverage |
| --- | --- |
| CodeQL | Go and Python semantic analysis, `security-extended` queries |
| Semgrep Managed Scanning | SAST via the Semgrep AppSec Platform GitHub App; full and diff-aware PR scans |
| `gosec` (via golangci-lint) | Go security linting, audit mode |
| `govulncheck` | Go vulnerability database, call-graph aware |
| OSV-Scanner | dependency advisories across ecosystems |
| Gitleaks | committed-secret detection |
| SonarQube Cloud | quality gate, security hotspots |
| OSSF Scorecard | repository supply-chain posture |
| Dependency Review | license and advisory gate on dependency changes |
| zizmor | GitHub Actions workflow security audit |

Findings surface in the repository's code-scanning dashboard.

Code scanning requires either a public repository or GitHub Advanced Security.
While neither applies, every scan still runs and still gates the build — only the
dashboard upload is skipped, with an explicit notice in the workflow summary. No
check silently passes because a report could not be filed.
