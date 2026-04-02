# GoREveal External Binary Matrix Evaluation

> Status: empirical product/validation note
> Date: 2026-03-31
> Purpose: measure how current GoREveal behaves on one open-source multi-platform Go project outside the in-repo fixture corpus, so roadmap decisions can be based on observed cross-format behavior rather than only fixture-local success.

## Target Project

Evaluation target:
- `rclone` release `v1.73.3`
- source: `https://github.com/rclone/rclone`
- reason for choice:
  - open-source Go project
  - official release binaries for `linux`, `windows`, and `macOS`
  - both `amd64` and `arm64` are available for all three platforms

Downloaded local matrix:
- `.tmp/rclone-matrix/bin/linux-amd64/rclone`
- `.tmp/rclone-matrix/bin/linux-arm64/rclone`
- `.tmp/rclone-matrix/bin/darwin-amd64/rclone`
- `.tmp/rclone-matrix/bin/darwin-arm64/rclone`
- `.tmp/rclone-matrix/bin/windows-amd64/rclone.exe`
- `.tmp/rclone-matrix/bin/windows-arm64/rclone.exe`

## Method

Commands used:
- build local evaluator:
  - `python3 -m scripts.dev.podman_runner exec -- bash -lc 'mkdir -p .tmp && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal'`
- run matrix:
  - `./.tmp/goreveal analyze <binary>`
- one deeper user-code isolation check:
  - `./.tmp/goreveal peel .tmp/rclone-matrix/bin/linux-amd64/rclone`

This is an empirical smoke-check, not a benchmark and not a new broad support claim.

## Results Matrix

| Binary | Format | Build Info | Runtime Posture | Functions | Packages | Files | Peeling | Reading |
| --- | --- | --- | --- | ---: | ---: | ---: | --- |
| `linux-amd64` | `elf` | yes | `section_heuristic` | `83635` | `1774` | `431` | `83635` | strong stripped-ELF result with real file visibility |
| `linux-arm64` | `elf` | yes | `section_heuristic` | `83479` | `1770` | `431` | `83479` | strong stripped-ELF result with real file visibility |
| `darwin-amd64` | `macho` | yes | none | `93909` | `2955` | `432` | `93909` | first bounded external Mach-O function/package/peeling foothold with real file visibility |
| `darwin-arm64` | `macho` | yes | none | `93149` | `2945` | `431` | `93149` | first bounded external Mach-O function/package/peeling foothold with real file visibility |
| `windows-amd64` | `pe` | yes | `section_heuristic` | `82128` | `1737` | `423` | `82128` | first bounded PE function/package/peeling foothold on a real external sample with real file visibility |
| `windows-arm64` | `pe` | yes | `section_heuristic` | `81953` | `1733` | `423` | `81953` | first bounded PE function/package/peeling foothold on a real external sample with real file visibility |

Recovered shared metadata across the entire matrix:
- `build_info.go_version = "go1.25.8"`
- `build_info.path = "github.com/rclone/rclone"`

## Deeper Linux Check

`peel` on `linux-amd64` produced a meaningful user-only view:
- `12299` user-classified functions
- `285` user-classified packages
- first recovered package set includes:
  - `github.com/rclone/rclone`
  - `github.com/rclone/rclone/backend/alias`
  - `github.com/rclone/rclone/backend/archive`
  - `github.com/rclone/rclone/backend/azureblob`

Practical reading:
- the current `ELF` path is already useful on large real stripped Go binaries
- code peeling is not just fixture-local; it already yields actionable analyst reduction on real external samples

## Interpretation

### What already works well

- `ELF` stripped binaries on both `amd64` and `arm64`
- build info recovery across all tested formats and architectures
- first-pass user-code isolation on real `ELF` binaries
- first bounded `PE` function/package/peeling foothold plus `PE` runtime posture on both `amd64` and `arm64`
- real file visibility across the full fresh `rclone` matrix, not only functions/packages

### What is useful but still intentionally narrow

- `PE` support on both `amd64` and `arm64`
  - current value is honest and real
  - but it is still a first function/package/peeling foothold plus bounded runtime posture, not broader semantic recovery

### What is currently the clearest gap

- richer semantic/source confidence across formats
  - current practical result now includes real function/package/peeling footholds on `ELF`, `PE`, and `Mach-O`
  - current practical result also includes real file visibility on all tested binaries in the matrix
  - `runtime` remains absent on `Mach-O`
  - there is still no broader `Mach-O` semantic or runtime claim, and file-backed visibility remains stronger in some baseline tools

## Weighted Product Reading

This matrix changes the planning picture slightly.

Before this check, the strongest next move looked like:
- continue polishing the transfer workflow above `ELF`

After the first `PE` and `Mach-O` footholds, the weighted reading becomes:
- `ELF` is already strong enough for another workflow increment when the target user lives mostly in Linux/server Go binaries
- `PE` is no longer posture-only, but it still remains a bounded foothold rather than a broad semantic layer
- `Mach-O` is no longer build-info-only either, so the next best move shifts away from format foothold work and back toward a real comparison plus the strongest workflow, semantic/source-confidence, or interop gap that comparison exposes

## Recommendation

If the near-term product goal is:

1. best immediate value for current Linux-focused Go RE workflows:
   - continue with the first annotated or persisted transfer workflow over `accepted_transfers`
   - this remains the right immediate workflow step once the first bounded `Mach-O` foothold exists

2. best cross-platform product credibility:
   - keep the current bounded `Mach-O` foothold stable while avoiding widened runtime claims
  - use a real comparison pass against `gore`, `GoReSym`, and `redress` to decide whether the next increment should be transfer polish, workstation handoff hardening, runtime depth, or one more bounded cross-format slice

Current weighted recommendation:
- treat the first bounded `PE` and `Mach-O` function footholds as landed
- rerun the real comparison against the existing Go RE tools next
- use that comparison to choose between:
  - transfer-workflow polishing over the current `accepted_transfers`
  - thin source-visibility improvement
  - one more bounded semantic/runtime slice where the comparison still shows the largest gap

Reason:
- Linux `ELF` already proved useful on real external binaries
- `PE` already has an honest bounded foothold
- `Mach-O` is no longer absent either, so the highest-value next move becomes evidence-driven comparison rather than another format foothold by default
