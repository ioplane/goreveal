---
title: Contributing to GoREveal
status: active
date: 2026-07-30
owners:
  - ioplane/goreveal-maintainers
tags:
  - contributing
  - process
---

# Contributing to GoREveal

<a href="https://github.com/ioplane/goreveal/blob/main/LICENSE">
  <img
    src="https://shieldcn.dev/badge/license-Apache--2.0-slate.svg?variant=outline&size=xs"
    alt="License" height="20"></a>
<a href="https://www.conventionalcommits.org/en/v1.0.0/">
  <img
    src="https://shieldcn.dev/badge/commits-conventional-slate.svg?variant=outline&size=xs&logo=git"
    alt="Conventional Commits" height="20"></a>
<a href="https://github.com/ioplane/goreveal/blob/main/CODE_OF_CONDUCT.md">
  <img
    src="https://shieldcn.dev/badge/CoC-Contributor%20Covenant-slate.svg?variant=outline&size=xs"
    alt="Code of Conduct" height="20"></a>

Thanks for considering a contribution. This document is the practical contract:
what we accept, how to build and verify, and the two rules that are genuinely
non-negotiable.

By participating you agree to the [Code of Conduct](CODE_OF_CONDUCT.md).
Contributions are licensed under [Apache-2.0](LICENSE).

## The two hard rules

### 1. Clean-room boundary

GoREveal studies the *observable behavior* of prior art — `gore`, `redress`,
`GoReSym`, `GoResolver`, `gostringungarbler`, `AlphaGolang` — and never copies
their implementation.

**Allowed:** reading those projects to understand a binary format, an edge case,
or a test idea; writing a differential test against their output; documenting
where GoREveal agrees or diverges.

**Not allowed:** copying code, translating code into Go with cosmetic changes,
or treating a reference tool's output as ground truth without validation.

If reference behavior is useful, land it as one of: a documented finding, a
corpus fixture, a differential test, or a deliberately designed implementation.
A pull request that reads like a port of an AGPL project will be closed.

### 2. Never invent recovery truth

GoREveal analyzes hostile, stripped, and obfuscated binaries. When evidence is
absent, the correct output is `unavailable`, an empty list, or an explicit error —
never a plausible guess.

- `schema` is the canonical contract; operator surfaces project it, they do not
  re-derive it
- `deobfuscation` refines; it must never overwrite raw recovered truth
- provenance and confidence are first-class result fields, not annotations
- degrade to `[]`, not JSON `null`, when there is no truthful surface

## Architecture invariants

| Rule | Meaning |
| --- | --- |
| `schema` is canonical | All outputs flow from it; no plugin-local contracts |
| `core` stays independent | No CLI, storage, API, or plugin imports in `core` |
| Plugins consume exports | `plugins/` never implements recovery logic |
| SIMD is optimization only | A scalar reference path must always exist |
| Bounded slices, not rewrites | Small, evidence-backed increments beat speculative parser work |

`.golangci.yml` enforces the module boundary mechanically via `depguard`
(`unsafe` and `os/exec` are denied outside the baseline harness).

## Development environment

Development is **Podman-first**. Do not rely on host Go or host linters — the
dev container pins every tool version.

```bash
git clone https://github.com/ioplane/goreveal.git
cd goreveal

# Build the dev image (Go 1.26 + linters + baseline tooling)
podman build -f deployments/docker/Containerfile.dev -t localhost/goreveal:dev .

# Or through the Python automation entrypoint
task build-image
```

### Python tooling (uv)

Python helpers under `scripts/` and `plugins/` are managed with
[uv](https://docs.astral.sh/uv/):

```bash
uv sync --group dev        # resolve and install from uv.lock
uv run ruff check .        # lint
uv run ruff format .       # format
uv run ty check            # type-check
uv run yamllint .          # YAML lint
```

`uv.lock` is committed and is the reproducibility boundary. If you change a
dependency, commit the updated lockfile.

### Documentation linting

Markdown is checked with a pinned `markdownlint-cli2`; there is no Node manifest
in the repository, so run it directly:

```bash
npx markdownlint-cli2@0.19.0
npx markdownlint-cli2@0.19.0 --fix
```

CI runs the same version through a SHA-pinned action.

### Verification

Run this before opening a pull request. CI runs the same gates.

```bash
task fmt            # gofmt + gofumpt inside the container
task lint           # golangci-lint
task test           # Go suite + Python unit tests
task lint-scripts   # ruff + ty + yamllint + shellcheck
task test-snapshots # golden snapshot suite
```

Equivalent `make` targets exist (`make fmt`, `make test`, `make lint`,
`make lint-scripts`). Keep verification **sequential** — overlapping runs against
the same bind-mounted workspace produce flaky results.

### Differential tests

`task test-differential` compares GoREveal against baseline tools and needs those
repositories on disk. Point the harness at them:

```bash
export GOREVEAL_BASELINES_HOST_ROOT=/path/to/your/baseline/repos
task test-differential
```

The directory should contain checkouts named `gore`, `redress`, `GoReSym`,
`GoResolver`, `gostringungarbler`, and `AlphaGolang`. These are reference-only;
they are never runtime dependencies. See
[docs/architecture/2026-03-19-goreveal-baseline-sources.md](docs/architecture/2026-03-19-goreveal-baseline-sources.md).

Differential tests are skipped, not failed, when baselines are absent — so CI
stays green without them.

## Making a change

### Prefer vertical slices

A good change carries recovery logic, its schema mapping, CLI or export exposure
if user-visible, and evidence. An architecture-only change with no capability
attached is usually not worth landing.

### Every change leaves evidence

Pick at least one, matched to what you changed:

| Change type | Required evidence |
| --- | --- |
| New or altered recovery | corpus fixture + golden snapshot |
| Changed recovery semantics | updated golden output **and** differential expectations |
| Parsing or boundary handling | fuzz target |
| Performance work | benchmark before/after (`benchstat`) |
| Contract or behavior change | updated docs |

If a change affects recovery semantics, also review whether
provenance/confidence behavior needs to change. Silently refreshing a golden
snapshot to make a test pass is how correctness regressions ship — explain the
diff in the pull request.

### Performance policy

Strict order, no shortcuts:

1. pure Go reference implementation
2. optimized scalar implementation
3. architecture-specific SIMD
4. optional experimentation

No SIMD without hotspot evidence, correctness-equivalence tests, a scalar
fallback, and documented feature gating. **Accuracy work outranks performance
work** — if you must choose, choose accuracy.

## Commits and pull requests

### Conventional Commits

Commit messages and **pull request titles** follow
[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/). CI
enforces the PR title, because it becomes the squash-merge commit and feeds the
generated changelog.

```text
<type>(<optional scope>): <description>
```

Allowed types: `feat`, `fix`, `perf`, `refactor`, `docs`, `test`, `build`, `ci`,
`chore`, `revert`, `security`.

```text
feat(core/pclntab): recover funcnametab offsets on stripped ELF
fix(engine/peeling): keep classification evidence on empty package sets
docs(architecture): record server stack decision
```

Breaking changes take a `!` before the colon and a `BREAKING CHANGE:` footer.

### Pull request checklist

- [ ] Title follows Conventional Commits
- [ ] Scoped to one logical change
- [ ] Evidence attached (see table above)
- [ ] `task lint`, `task test`, `task lint-scripts` pass locally
- [ ] Docs updated if behavior or contract changed
- [ ] No copied code from reference projects
- [ ] No host paths, secrets, or internal hostnames added

Every pull request runs the full CI matrix plus CodeQL, Trivy, `govulncheck`,
OSV-Scanner, Gitleaks, SonarQube Cloud, and dependency review, and is scanned by
Semgrep Managed Scanning on the Semgrep AppSec Platform.
Security findings surface in the code-scanning dashboard, not as raw log output.

### Review

Maintainers listed in [MAINTAINERS.md](MAINTAINERS.md) review changes.
`.github/CODEOWNERS` routes reviews by area. Expect questions about evidence and
claim boundaries — they are the point of the review, not friction.

## Reporting problems

| Kind | Where |
| --- | --- |
| Security vulnerability | [Private advisory](https://github.com/ioplane/goreveal/security/advisories/new) — see [SECURITY.md](SECURITY.md) |
| Incorrect recovery output | [Bug report issue](https://github.com/ioplane/goreveal/issues/new/choose) |
| Feature or capability request | [Feature request issue](https://github.com/ioplane/goreveal/issues/new/choose) |
| Question or usage help | [Discussions](https://github.com/ioplane/goreveal/discussions) — see [SUPPORT.md](SUPPORT.md) |

Wrong recovery output is a **bug**, not a vulnerability. Report it publicly so it
can get a fixture.

## Repository layout

```text
core/            recovery primitives (format, pclntab, runtime, types, strings)
schema/          canonical analysis contract and export encoders
engine/          pipeline orchestration, peeling, projection
deobfuscation/   garble string and name refinement (never overwrites raw truth)
storage/         SQLite persistence and build-to-build diffing
cmd/goreveal/    operator CLI
plugins/         thin IDA and Ghidra adapters (consume exports only)
corpus/          fixtures and the protected-binary sample
tests/           snapshot and differential suites
scripts/         Python automation (Podman runner, baseline harness)
docs/            architecture and release documentation
deployments/     Containerfiles (dev, builder, release)
```

`AGENTS.md` is the operational contract for automated contributors. If you are an
AI agent working in this repository, read it first.
