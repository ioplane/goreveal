# GoREveal Protected Binary Comparison Plan

> Status: active next lane
> Date: 2026-04-01
> Purpose: turn the commercial-software protection research into a bounded comparison and corpus plan, so `GoREveal` can be evaluated against realistic protected Go binaries without prematurely expanding `core` claims.

## Why This Exists

The current repo now has maintained footholds on:
- `ELF`
- `PE`
- `Mach-O`

That changes the next question.

The next question is no longer only:
- which format foothold is still missing?

It is also:
- how well does `GoREveal` behave on realistic protected Go binaries built for commercial distribution?

The new draft on protecting commercial Go software points to a concrete target class:
- stripped binaries via `-s -w`
- build-path-reduced binaries via `-trimpath`
- `PIE` builds
- garbled binaries
- OSS + Enterprise style binaries with feature gates and license checks

This target class matters, but it should enter the roadmap as:
- a corpus/comparison lane first
- a bounded workflow/product lane second
- parser or deobfuscation work only after measured gaps are clear

## Scope

This plan is for:
- evidence-backed comparison
- target-binary taxonomy
- fixture and corpus growth
- analyst-workflow evaluation

This plan is not for:
- immediate anti-obfuscation claims
- immediate native garble deobfuscation parity
- broad parser rewrites inside `core`
- license-system implementation inside `GoREveal` itself

## Target Profiles

The first protected-binary matrix should cover these build profiles:

1. `plain`
   - normal Go build

2. `stripped`
   - `-ldflags="-s -w -buildid="`

3. `trimpath`
   - `-trimpath`

4. `stripped + trimpath`
   - `-trimpath -ldflags="-s -w -buildid="`

5. `pie`
   - `-buildmode=pie`

6. `stripped + trimpath + pie`
   - realistic hardened non-garbled baseline

7. `garble`
   - name obfuscation only

8. `garble + literals + tiny`
   - realistic bounded garble stress case

9. `enterprise-gated`
   - binary with feature-flag checks and a simple license-validation path

The first matrix does not need all profiles across all formats.
It should start with:
- `linux/amd64`
- `windows/amd64`
- `darwin/amd64`

Then expand to:
- `arm64`

## First Candidate Targets

Prefer targets with open source and easy reproducible builds.

Good initial classes:
- one small in-repo purpose-built sample with simple feature gating
- one real OSS server-style binary
- one CLI-style binary

Target selection rules:
- source must be available
- build recipe must be reproducible in the dev container
- license or feature-check logic must be simple enough to reason about
- avoid malware-only samples in the first pass

## Comparison Questions

For each protected-binary target, measure:
- format detection still correct?
- build info still present or absent?
- function/package/peeling surfaces still useful?
- source/file visibility degraded how much?
- current `classification_evidence` still meaningful?
- `matched_functions` and transfer surfaces still useful across neighboring builds?
- can license checks or enterprise feature-gate functions still be isolated as likely user code?
- where do baseline tools outperform `GoREveal` materially?

## Output Shape

Expected outputs from this lane:
- one repo-native comparison note
- one compact matrix by profile and format
- one short list of recurring protected-binary pain points
- one short list of capabilities already strong enough
- one weighted recommendation for:
  - `engine/peeling`
  - `storage/diff`
  - `deobfuscation`
  - or one more bounded recovery slice

## Decision Rules

Promote work into the active implementation lane only if:
- the protected-binary comparison reveals a repeated analyst pain point
- that pain point is not already solvable through current schema, peeling, diff, or export surfaces
- the next step can stay bounded and evidence-backed

Prefer next steps in this order:
1. `engine/peeling` refinement
2. `storage/diff` / transfer workflow refinement
3. thin source-visibility refinement over current line-table-backed fallback
4. `deobfuscation` orchestration or bounded refinement
5. one more bounded recovery slice

Prefer not to:
- widen `core` first
- brand a capability as anti-obfuscation before it has corpus-backed evidence
- treat `garble` handling as equivalent to general protected-binary support

## Relationship To Existing Plans

This lane comes after:
- the post-`PE` baseline comparison rerun
- the first package-level transfer-workflow polish increment
- the first thin source-visibility response

Current phase:
- the first purpose-built enterprise-gated matrix pass is now started and recorded in `docs/plans/2026-04-01-goreveal-protected-binary-initial-results.md`
- that pass now proves baseline-compared stability across the first non-garbled profile matrix on `linux/amd64`, `windows/amd64`, and `darwin/amd64`
- the operator-path drift in `task protected-matrix` is now fixed too, so the canonical task output matches the direct script run and shows stable file visibility on the full first non-garbled matrix
- the protected matrix now also prefers a local source-built `/repos/garble` checkout when available, so the previous `v0.15.0` release-gap is no longer the active blocker
- the first bounded `garble` rows are now measured on `linux/amd64`, and they currently collapse to `functions = 0` / `packages = 0` / `files = 0` while still preserving `runtime_trust_summary = "go_module_fallback"`, exposing `elf_pclntab_header_magic_kind = "unknown"`, and now carrying the explicit blocker `elf_function_recovery_blocker = "custom_pclntab_magic"`
- the `arm64` widening checkpoint is now landed
- the bounded `linux/arm64` triage pass is now also landed:
  - the old split turned out not to be a garble-specific offset failure
  - `moduledata_text_*` stays absent on the current `linux/arm64` family
  - a separate ELF `.text` section range is enough to make the bounded `elf_function_foothold = "address_only"` surface portable across `linux/amd64` and `linux/arm64`
- one more compact analyst-facing slice is now landed too:
  - `elf_function_foothold_text_source` distinguishes `moduledata_text` and `elf_text_section`
  - `elf_function_foothold_start_addr` / `elf_function_foothold_end_addr` project the bounded foothold span directly
  - the same compact runtime surface is now asserted in thin `IDA` / `Ghidra` export-contract tests and plugin consumer tests, so adapter consumers inherit it without local recomputation
- the next sub-step is therefore no longer architecture-specific triage, but a product decision:
  - treat the current protected lane as the evidence baseline
  - return to workflow/value work by default
  - only add one stronger but still truthful protected-specific analyst surface if a specific measured analyst pain point justifies it

This lane should inform:
- `Go code peeling / user-code isolation`
- `Go version tracking / build correlation`
- future bounded deobfuscation priorities
- enterprise feature-gate and license-check triage workflows

Related documents:
- `docs/tmp/draft/protecting-commercial-software-on-go.md`
- `docs/plans/2026-04-01-garble-go126-support-research.md`
- `docs/plans/2026-03-31-goreveal-baseline-comparison-plan.md`
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md`
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md`
