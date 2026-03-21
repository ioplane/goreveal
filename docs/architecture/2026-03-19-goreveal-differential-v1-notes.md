# GoREveal Differential Validation v1 Notes

> Status: active validation note
> Date: 2026-03-19

## Scope

Differential validation v1 intentionally covers only the smallest stable overlap between GoREveal and external baseline tools for the current real Go `ELF` fixture.

The current comparison set is:
- `GoReSym` for build info parity, source-file overlap, and user-function overlap
- `redress` for module-path parity, package-presence, source-file, and user-function overlap
- `gore` for build info parity, package presence, source-file basename overlap, and user-function overlap

This is a proof-of-harness milestone, not a final compatibility matrix.

## Assertions Currently Backed by Evidence

For the `go-elf-buildinfo-linux-amd64` fixture, GoREveal currently proves:
- build info `path` parity with `GoReSym`
- build info `go_version` parity with `GoReSym`
- function overlap for `main.main` and `main.helperAdd` against `GoReSym`
- function overlap for `main.helperBanner` against `GoReSym`
- build info `path` parity with `gore`
- build info `go_version` parity with `gore`
- source-tree overlap for the projected module-local file path
- build info `path` parity with `redress` through `gomod`
- package presence overlap for `main` against `redress`
- source-file overlap for `main.go` against `redress`
- function overlap for `main.main` and `main.helperBanner` against `redress`
- package presence overlap for `main` against `gore`
- source-file basename overlap for `main.go` against `gore`
- function overlap for `main.main` and `main.helperBanner` against `gore`

These checks run in the Podman dev container through:
- `make test-differential`

## Allowed Divergences in v1

The following divergences are currently expected and should not be treated as regressions by themselves:
- GoReSym reports a much broader runtime-aware file list than GoREveal currently projects into `source-tree`
- GoReSym overlap is currently asserted only for a narrow user-function set, not complete function parity
- `redress` source projection is currently used only as a narrow overlap surface for module-local file and user-function presence, not full source-reconstruction parity
- `redress` currently exposes only the module-local source surface for the canonical fixture; GoREveal's new `external_packages` projection is therefore treated as an intentional product improvement, not a parity requirement
- `redress` module-path parity currently comes from `gomod` output and does not imply broader build-metadata parity
- `gore` function names require normalization into package-qualified form before comparison
- `gore` type output for the current fixture is dominated by runtime-heavy surface area and is not yet a trustworthy user-type overlap surface, so it is intentionally not asserted
- GoREveal type recovery is currently `DWARF`-based for the active fixture, not full typelink/moduledata recovery
- GoREveal string recovery is currently `ELF` section-scan based, while other tools may use different heuristics or not expose the same surface directly
- GoREveal keeps raw and refined truth separate, so deobfuscation-oriented outputs should not be compared naively to raw recovery outputs from other tools
- GoREveal package classification intentionally suppresses synthetic pseudo-functions that are not meaningful package evidence

## Better-Than-Baseline Cases Already Present

The following properties are already stronger in GoREveal than in the current comparison setup:
- canonical schema output instead of tool-specific ad hoc text output
- explicit provenance-aware separation between raw and refined analysis layers
- SQLite persistence for analysis results and repeatable export flow
- container-first reproducibility for routine verification
- explicit `external_packages` visibility for non-module source evidence that the current `redress` fixture path does not surface

These are platform qualities, not yet broader recovery-superiority claims.

## Current Non-Claims

Differential validation v1 does not yet justify claims about:
- complete function parity with GoReSym
- source reconstruction parity with redress
- complete function parity with redress
- deobfuscation parity with GoResolver or gostringungarbler
- multi-format parity across `PE` or `Mach-O`
- broad compatibility across multiple Go version families

## Next Expansion Targets

The next differential-expansion pass should add one or more of:
- consolidation of shell baseline wrappers into a shared Python normalization layer before the harness grows much further
- `GoResolver` for deobfuscation-oriented comparison
- `gostringungarbler` for string refinement comparison

That consolidation is now started:
- `GoReSym`, `redress`, and `gore` normalization are now shared through `scripts.baseline.normalize`
- `gore` still uses a shell + temporary Go helper path for extraction, but no longer carries its own separate normalization logic

Each expansion should first define a narrow overlap surface, then add normalized wrappers, then add evidence-backed assertions.
