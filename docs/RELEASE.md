---
title: Release Process and Artifact Verification
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - release
  - supply-chain
  - verification
---

# Release Process and Artifact Verification

<a href="https://slsa.dev/spec/v1.2/build-track-basics">
  <img
    src="https://shieldcn.dev/badge/SLSA-Build%20L2-slate.svg?variant=outline&size=xs"
    alt="SLSA Build Level 2" height="20"></a>
<a href="https://spdx.dev/">
  <img
    src="https://shieldcn.dev/badge/SBOM-SPDX%20JSON-slate.svg?variant=outline&size=xs"
    alt="SPDX JSON SBOM" height="20"></a>
<a href="https://www.sigstore.dev/">
  <img
    src="https://shieldcn.dev/badge/signing-Sigstore%20keyless-slate.svg?variant=outline&size=xs&logo=sigstore"
    alt="Sigstore keyless signing" height="20"></a>

This document has two audiences. **Consumers** should read
[Verifying a release](#verifying-a-release) — it is the part that matters if you
are about to run a GoREveal binary. **Maintainers** should read
[Cutting a release](#cutting-a-release) and
[One-time setup](#one-time-setup).

## What every release contains

A tagged release publishes the following, all produced by
[`.github/workflows/release.yml`](../.github/workflows/release.yml) and
[`.goreleaser.yaml`](../.goreleaser.yaml). Nothing is ever built or uploaded by
hand.

| Artifact | Purpose |
| --- | --- |
| `goreveal_<version>_<os>_<arch>.tar.gz` / `.zip` | Binary archive with `LICENSE`, `NOTICE`, and docs |
| `goreveal_<version>_<os>_<arch>.deb` / `.rpm` | Debian and RPM packages |
| `goreveal_<version>_source.tar.gz` | Source archive for the tag |
| `*.spdx.json` | SPDX JSON SBOM per archive, plus one for the source |
| `goreveal_<version>_checksums.txt` | SHA-256 over every artifact above |
| `goreveal_<version>_checksums.txt.sig` | Sigstore signature over the checksum file |
| `goreveal_<version>_checksums.txt.pem` | Signing certificate with the workflow identity |
| `ghcr.io/ioplane/goreveal:<version>` | Multi-architecture container image (`linux/amd64`, `linux/arm64`) |

Build provenance attestations are recorded against the repository rather than
uploaded as files; retrieve them with `gh attestation verify`.

### Target matrix

| OS | `amd64` | `arm64` |
| --- | --- | --- |
| Linux | Yes | Yes |
| macOS | Yes | Yes |
| Windows | Yes | No |

`windows/arm64` is deliberately excluded: there is no verified evidence lane for
it, and shipping an untested target would contradict the project's rule against
unsupported claims.

## Supply-chain guarantees

### SLSA Build Track

Releases satisfy **[SLSA v1.2](https://slsa.dev/spec/v1.2/build-track-basics)
Build Level 2**:

- **L1 — provenance exists.** The build runs from a scripted, reviewed workflow
  on a hosted platform, and provenance describing the builder, the process, and
  the top-level inputs is generated automatically.
- **L2 — provenance is signed by the platform.** GitHub's attestation service
  signs the provenance with an identity bound to
  `.github/workflows/release.yml` at the released tag. Consumers verify that
  signature with `gh attestation verify`.

The workflow also runs on ephemeral GitHub-hosted runners, and the signing
material comes from an OIDC token that the build steps cannot read — the two
properties Build L3 adds. We do not self-certify L3, because that claim depends
on the build platform's own attested guarantees rather than on this
repository's configuration. **Verify what the attestation actually says rather
than trusting a level badge**, here or anywhere else.

### What a signature does and does not prove

A successful verification proves that **this repository's release workflow, at
that tag, produced that byte sequence**. It does not prove the code is free of
defects, and it does not prove the release is fit for your purpose. It closes
the tampering-in-transit gap, nothing more.

## Verifying a release

Run all four steps. Step 1 alone proves nothing — an attacker who can replace
the artifact can replace the checksum file next to it.

```bash
VERSION=0.1.0
BASE="https://github.com/ioplane/goreveal/releases/download/v${VERSION}"
ARCHIVE="goreveal_${VERSION}_linux_amd64.tar.gz"

curl -fsSLO "${BASE}/${ARCHIVE}"
curl -fsSLO "${BASE}/goreveal_${VERSION}_checksums.txt"
curl -fsSLO "${BASE}/goreveal_${VERSION}_checksums.txt.sig"
curl -fsSLO "${BASE}/goreveal_${VERSION}_checksums.txt.pem"
```

### 1. Checksum integrity

```bash
sha256sum --check --ignore-missing "goreveal_${VERSION}_checksums.txt"
```

### 2. Signature authenticity

This is the step that establishes trust. The identity regexp pins the signature
to this repository's release workflow at a version tag; a signature from any
other workflow, repository, or ref will fail.

```bash
cosign verify-blob \
  --certificate "goreveal_${VERSION}_checksums.txt.pem" \
  --signature "goreveal_${VERSION}_checksums.txt.sig" \
  --certificate-identity-regexp \
    '^https://github\.com/ioplane/goreveal/\.github/workflows/release\.yml@refs/tags/v' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "goreveal_${VERSION}_checksums.txt"
```

Expected output ends with `Verified OK`. Anything else — including a warning you
do not understand — means stop and report it through
[SECURITY.md](../SECURITY.md).

### 3. Build provenance

```bash
gh attestation verify "${ARCHIVE}" --repo ioplane/goreveal
```

To inspect the provenance rather than only validate it:

```bash
gh attestation verify "${ARCHIVE}" --repo ioplane/goreveal --format json \
  | jq '.[].verificationResult.statement.predicate.buildDefinition'
```

### 4. Runtime identity

Confirm the binary you extracted reports the version you intended to install:

```bash
tar -xzf "${ARCHIVE}"
./goreveal version --json
```

```json
{
  "version": "0.1.0",
  "git_commit": "...",
  "build_date": "...",
  "go_version": "go1.26.5",
  "platform": "linux/amd64",
  "modified": false
}
```

`"modified": true` on a release artifact means the build tree was dirty and the
artifact should not be trusted.

### Container image

```bash
IMAGE="ghcr.io/ioplane/goreveal:${VERSION}"
docker pull "${IMAGE}"

gh attestation verify "oci://${IMAGE}" --repo ioplane/goreveal
docker run --rm "${IMAGE}" version --json
```

### Inspecting the SBOM

The SPDX JSON documents feed directly into standard scanners:

```bash
curl -fsSLO "${BASE}/${ARCHIVE}.spdx.json"

# Component and license inventory
jq -r '.packages[] | [.name, .versionInfo, (.licenseConcluded // "NOASSERTION")]
       | @tsv' "${ARCHIVE}.spdx.json" | column -t

# Scan the SBOM for known advisories
trivy sbom "${ARCHIVE}.spdx.json"
grype "sbom:${ARCHIVE}.spdx.json"
```

## Cutting a release

### Preconditions

- `main` is green: CI, CodeQL, Semgrep, Trivy, and supply-chain workflows all
  passing.
- `CHANGELOG.md` has an `Unreleased` section describing the release, and every
  entry is accurate. Overstating a capability here is a release defect.
- No golden-snapshot drift (`task test-snapshots` clean).
- Version number chosen per [SemVer 2.0.0](https://semver.org/spec/v2.0.0.html).
  While pre-1.0, breaking `schema` contract changes bump the minor version.

### Steps

```bash
# 1. Move Unreleased into a dated version section.
$EDITOR CHANGELOG.md
git commit -am "docs(changelog): prepare v0.1.0"

# 2. Dry run the full pipeline locally.
goreleaser check
goreleaser release --snapshot --clean --skip=sign,publish

# 3. Tag and push. The tag is what triggers the pipeline.
git tag -s v0.1.0 -m "goreveal v0.1.0"
git push origin main
git push origin v0.1.0
```

Prefer a **signed** tag (`-s`). It is not verified by the pipeline, but it
records maintainer intent independently of the platform.

### What the pipeline does

1. **Preflight** — verifies the module graph, runs the race-enabled test suite,
   asserts golden snapshots have not drifted, and validates the GoReleaser
   configuration. A tag pointing at a broken tree fails here, before anything is
   built.
2. **Build and publish** — GoReleaser cross-compiles all targets with
   `-trimpath` and `-buildid=`, builds archives and packages, generates SPDX
   SBOMs with Syft, computes checksums, signs the checksum file with Cosign
   keyless signing, and publishes the GitHub Release and GHCR images.
3. **Attest** — records SLSA build provenance and SBOM attestations.
4. **Self-verify** — re-checks the checksums, validates every SBOM's structure,
   verifies the Cosign signature from the consumer's side, and runs
   `goreveal version` from the published archive.

Step 4 exists because a release nobody can verify is worse than no release. If
it fails, the release is already published but must be treated as suspect: delete
the tag and the release, fix the cause, and re-tag with a new patch version.

### Manual dry run

```bash
gh workflow run release.yml --field dry_run=true
```

This builds and self-verifies everything without publishing or signing.

### Reproducing a build

`-trimpath` plus `-buildid=` plus a `mod_timestamp` derived from the commit make
the binaries reproducible for a given tag:

```bash
git checkout v0.1.0
goreleaser build --clean --single-target --snapshot
sha256sum dist/goreveal_*/goreveal
```

Compare against the published checksum for the same OS and architecture. A
mismatch that is not explained by a toolchain version difference is worth
reporting.

## One-time setup

These are repository-level prerequisites, not per-release steps.

### Required repository settings

| Setting | Value | Why |
| --- | --- | --- |
| Actions → Workflow permissions | Read repository contents (default) | Workflows request writes explicitly per job |
| Branch protection on `main` | Require the `CI success` check | Single aggregate gate; see `ci.yml` |
| Branch protection on `main` | Require signed commits (recommended) | Author accountability |
| Code security → Code scanning | Enabled | Receives CodeQL, Semgrep, Trivy, OSV, and Scorecard SARIF |
| Code security → Private vulnerability reporting | Enabled | The channel `SECURITY.md` points to |
| Code security → Secret scanning + push protection | Enabled | Complements Trivy and Gitleaks |

### Secrets

| Secret | Required | Effect when absent |
| --- | --- | --- |
| `GITHUB_TOKEN` | Provided automatically | — |
| `SEMGREP_APP_TOKEN` | Optional | Semgrep falls back to OSS rulesets instead of Pro rules |
| `SONAR_TOKEN` | Optional | SonarQube Cloud analysis is skipped with a notice |
| `CODECOV_TOKEN` | Optional | Coverage upload is skipped; CI still passes |
| `GITLEAKS_LICENSE` | Optional | Only needed for organization-level Gitleaks features |

Every optional integration degrades with an explicit notice rather than a silent
pass or a hard failure, so forks and unconfigured clones still get a meaningful
CI run.

### SonarQube Cloud

[`sonar-project.properties`](../sonar-project.properties) ships with `TODO`
placeholders, and the workflow skips analysis while they are present. To enable
it:

1. Import the repository at <https://sonarcloud.io> and choose the GitHub
   Actions analysis method.
2. Copy the organization key and project key from the project's
   **Information** panel.
3. Replace both `TODO` values in `sonar-project.properties`.
4. Add the generated token as the `SONAR_TOKEN` repository secret.

### Container registry

The GHCR package is created on the first release. Afterwards, set its visibility
to public and link it to this repository under the package settings, so
`docker pull` works without authentication.

## Post-release

- Confirm the release notes rendered the verification block correctly.
- Run the four consumer verification steps against the real published artifacts.
  Verifying your own release from the outside is the only way to know the
  instructions in this document are still true.
- Open a `docs(changelog)` pull request adding a fresh `Unreleased` section.

## Yanking a release

There is no unpublish. If a release must be withdrawn:

1. Mark the GitHub Release as a pre-release and prepend a prominent notice to
   its body explaining the defect and the recommended version.
2. If the cause is a security issue, publish a
   [security advisory](https://github.com/ioplane/goreveal/security/advisories)
   with the affected version range.
3. Release a fixed version promptly. Do not delete the tag of a release that has
   been publicly available — consumers may have pinned it, and a deleted tag
   turns a verifiable artifact into an unverifiable one.
