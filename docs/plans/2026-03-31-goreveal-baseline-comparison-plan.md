# GoREveal Baseline Comparison Plan

> Status: historical comparison plan; RT1-S1 now owns baseline evidence work
> Date: 2026-03-31
> Purpose: run a real, evidence-backed comparison between `GoREveal` and the current Go RE baselines after the first bounded `Mach-O` foothold, so the next roadmap step is chosen from measured gaps rather than intuition.

## Baselines

Primary comparison set:
- `GoReSym`
- `redress`
- `gore`

Comparison rule:
- compare operator-visible outputs and usable analyst truth
- do not treat any baseline as infallible ground truth
- normalize outputs into overlap, divergence, and unsupported buckets

## Targets

In-repo fixtures:
- `corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin`
- `corpus/fixtures/go-elf-stripped-linux-amd64/fixture.bin`
- `corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe`
- `corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin`

External binaries:
- `/opt/projects/repositories/ocserv-agent/bin/ocserv-agent-linux-amd64`
- `/opt/projects/repositories/ocserv-agent/bin/ocserv-agent-linux-arm64`
- `/opt/projects/repositories/GoReSym/testproject/testproject.exe`
- `/opt/projects/repositories/hashicorp-re/bin/Keygen.exe`
- `.tmp/rclone-matrix/bin/linux-amd64/rclone`
- `.tmp/rclone-matrix/bin/linux-arm64/rclone`
- `.tmp/rclone-matrix/bin/darwin-amd64/rclone`
- `.tmp/rclone-matrix/bin/darwin-arm64/rclone`
- `.tmp/rclone-matrix/bin/windows-amd64/rclone.exe`
- `.tmp/rclone-matrix/bin/windows-arm64/rclone.exe`

## Comparison Dimensions

For each target, compare:
- format detection
- build info
- function count and named-function overlap
- package/import-path recovery
- source file and source line truth when exposed
- runtime posture or runtime-semantic output when exposed
- user-code isolation or equivalent analyst-reduction output when exposed
- honest unsupported behavior

## Output Shape

Expected deliverables:
- one repo-native note with the normalized matrix
- one short finding list of where `GoREveal` is already stronger
- one short finding list of where `GoREveal` is still clearly behind
- one weighted decision for the next bounded roadmap slice

## Decision Rule

After the comparison:
- choose transfer-workflow polish if the main remaining gap is analyst throughput on already-strong `ELF` data
- choose one more bounded semantic or cross-format slice if the main remaining gap is still recovery truth on a supported format
- do not widen `core` claims unless the comparison plus fixtures justify it

## Current Outcome

The first comparison pass is now recorded in:
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md`

Current weighted result:
- `ELF` is already strong enough that more workflow polish is no longer the default next move
- the first bounded `Mach-O` foothold is already strong enough to stop being the highest-priority gap
- the first bounded `PE` function foothold is now landed
- the post-`PE` comparison rerun is now completed for the current highest-signal slices
- that rerun shifts the next weighted move away from parser expansion and toward:
  - transfer-workflow polish
  - thin source-visibility work
  - protected-binary comparison
- after that rerun, the next bounded comparison lane should expand into protected commercial-style Go binaries:
  - `-s -w`
  - `-trimpath`
  - `PIE`
  - `garble`
  - enterprise feature-gate and license-check workflows

Protected-binary follow-up plan:
- `docs/plans/2026-04-01-goreveal-protected-binary-comparison-plan.md`
