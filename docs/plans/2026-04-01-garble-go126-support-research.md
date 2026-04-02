# Garble Go 1.26 Support Research

> Status: clean-room external research note
> Date: 2026-04-01
> Purpose: record what current upstream `garble` actually supports for `Go 1.26.x`, so the protected-binary lane can distinguish release lag from a true upstream capability gap.

## Scope

This note is behavioral and planning-oriented only.
It does not copy `garble` implementation into `GoREveal`.

Studied repository:
- `https://github.com/burrowers/garble`
- local clone: `/opt/projects/repositories/garble`

Observed local upstream state:
- `HEAD`: `8503b1b`
- nearest tag: `v0.15.0-32-g8503b1b`

## High-Level Finding

The earlier protected-matrix blocker was a release-gap, not a current upstream support gap.

More precisely:
- released `mvdan.cc/garble v0.15.0` does not support the repo's `go1.26.1` toolchain
- current upstream `master` does support `Go 1.26.x`
- current upstream `master` still explicitly blocks `go1.27+`

So the correct framing is:
- `garble` supports `Go 1.26.x` on current upstream `master`
- the `v0.15.0` release lags that support work

## Concrete Evidence

### Version gate

Current upstream `main.go` sets:
- `minGoVersion = "go1.26.0"`
- `unsupportedGo = "go1.27"`

This is an explicit support window, not an inferred one.

### Linker patch lane

Current upstream includes:
- `internal/linker/patches/go1.26/0001-add-custom-magic-value.patch`
- `internal/linker/patches/go1.26/0002-add-unexported-function-name-removing.patch`
- `internal/linker/patches/go1.26/0003-add-entryOff-encryption.patch`

That means `Go 1.26` support is backed by a dedicated linker-patch set, not by only loosening version checks.

### Generated tables

Current upstream `scripts/gen_go_std_tables.go` uses:
- `go1.26.1`

And current upstream `go_std_tables.go` is regenerated from that version.

### Runtime patch drift

Current upstream `runtime_patch.go` was updated for the `Go 1.26` runtime shape.
The notable change is that the magic-value patching moved to the newer constant location instead of the older `moduledataverify1` pattern.

### Upstream test evidence

Focused containerized verification against current upstream `master` succeeded:
- `go test -run TestScript/goversion -count=1 ./...`
- `go build -o /tmp/garble . && /tmp/garble version`

That is enough to support:
- current upstream `master` works on `go1.26.1`

That is not enough to support:
- generic future-version claims
- `go1.27+`
- a complete `1.26.x` compatibility matrix

## Minimal Backport Set

If a release-branch backport were ever needed, the minimum work is not a one-line version bump.

The clean-room patch cluster is:
1. `529ee19 improve go_std_tables.go generation`
2. `41b45f9 windows/arm is gone for Go 1.26, use GOARCH=386 instead`
3. `7114403 drop Go 1.25, support Go 1.26`

Recommended to include as part of the same support lane:
4. `6f2a1f8 scripts: fix a sneaky bug in the compiler intrinsics generator`

Why this is not just version-string churn:
- `Go 1.26` needed new linker patches
- runtime magic patching changed
- linker `linkname` behavior became stricter
- std/runtime table generation needed broader platform coverage
- testscript expectations needed explicit updates

## What Changed In Practice For GoREveal

The weighted decision changes in two steps.

First:
- do not describe `garble` itself as lacking `Go 1.26.x` support
- describe the issue as `v0.15.0` release lag relative to upstream `master`

Second:
- prefer a pinned local or source-built upstream `garble` for the protected matrix when available
- keep an older-toolchain workaround as a fallback only, not as the default next move

## Planning Consequence

After switching the protected matrix to local upstream `garble`, the next real question is no longer toolchain compatibility.

It becomes:
- what `GoREveal` actually preserves or loses on measured `garble` binaries

That is a product gap investigation, not an environment gap investigation.
