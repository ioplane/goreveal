# GoREveal Protected Binary Initial Results

> Status: initial measured checkpoint
> Date: 2026-04-01
> Purpose: record the first protected-binary matrix pass over a purpose-built enterprise-gated sample, so the next lane decision can use measured profile behavior instead of assumptions.

## Scope

This first pass is intentionally narrow:
- one purpose-built open-source sample:
  - `corpus/protected/enterprise-sample/src`
- one first build-profile matrix:
  - `plain`
  - `stripped`
  - `trimpath`
  - `stripped-trimpath`
  - `pie` where supported
- three platforms:
  - `linux/amd64`
  - `windows/amd64`
  - `darwin/amd64`

This pass is not yet the full protected-binary comparison promised in the plan:
- `arm64` widening is now landed for `linux`, `windows`, and `darwin`
- the first `garble` and `garble + literals + tiny` rows now exist on both `linux/amd64` and `linux/arm64`
- current `garble` coverage uses a local source-built upstream `garble` checkout, not the older `v0.15.0` release binary

## Sample Shape

The current enterprise-gated sample includes bounded user-facing targets:
- `main.readLicenseToken`
- `main.auditFeatureGate`
- `main.runEnterpriseReport`

It also includes:
- a simple feature-gate path
- a simple license-check path
- module-local user code spread across:
  - `main`
  - `internal/features`
  - `internal/licensegate`

## Measured GoREveal Matrix

| Target | Profile | Format | Runtime Trust | Functions | Packages | Files | Pathless File Evidence | User Functions | User Packages | Relevant User Functions |
| --- | --- | --- | --- | ---: | ---: | ---: | --- | ---: | ---: | --- |
| `linux-amd64` | `plain` | `elf` | `symbol_backed` | `2118` | `58` | `3` | `true` | `9` | `3` | all three present |
| `linux-amd64` | `stripped` | `elf` | `go_module_fallback` | `2118` | `58` | `3` | `true` | `9` | `3` | all three present |
| `linux-amd64` | `trimpath` | `elf` | `symbol_backed` | `2118` | `58` | `3` | `false` | `9` | `3` | all three present |
| `linux-amd64` | `stripped-trimpath` | `elf` | `go_module_fallback` | `2118` | `58` | `3` | `true` | `9` | `3` | all three present |
| `linux-amd64` | `pie` | `elf` | `symbol_backed` | `2118` | `58` | `3` | `true` | `9` | `3` | all three present |
| `linux-amd64` | `garble` | `elf` | `go_module_fallback` | `0` | `0` | `0` | `false` | `0` | `0` | none |
| `linux-amd64` | `garble-literals-tiny` | `elf` | `go_module_fallback` | `0` | `0` | `0` | `false` | `0` | `0` | none |
| `windows-amd64` | `plain` | `pe` | `section_heuristic` | `2120` | `59` | `3` | `true` | `9` | `3` | all three present |
| `windows-amd64` | `stripped` | `pe` | `section_heuristic` | `2119` | `59` | `3` | `true` | `9` | `3` | all three present |
| `windows-amd64` | `trimpath` | `pe` | `section_heuristic` | `2120` | `59` | `3` | `true` | `9` | `3` | all three present |
| `windows-amd64` | `stripped-trimpath` | `pe` | `section_heuristic` | `2119` | `59` | `3` | `true` | `9` | `3` | all three present |
| `darwin-amd64` | `plain` | `macho` | `absent` | `2239` | `58` | `3` | `true` | `9` | `3` | all three present |
| `darwin-amd64` | `stripped` | `macho` | `absent` | `2238` | `58` | `3` | `true` | `9` | `3` | all three present |
| `darwin-amd64` | `trimpath` | `macho` | `absent` | `2239` | `58` | `3` | `true` | `9` | `3` | all three present |
| `darwin-amd64` | `stripped-trimpath` | `macho` | `absent` | `2238` | `58` | `3` | `true` | `9` | `3` | all three present |
| `darwin-amd64` | `pie` | `macho` | `absent` | `2239` | `58` | `3` | `true` | `9` | `3` | all three present |

Skipped:
- `windows-amd64` `pie`
  - `profile not enabled for platform`

## Reading

What this already shows:
- current `GoREveal` surfaces survive all first non-garbled protected profiles on the purpose-built sample
- the three analyst-relevant user functions remain visible across every measured profile
- the current user-code isolation layer keeps all three user packages visible across every measured profile
- `-s -w` changes the `ELF` runtime posture from `symbol_backed` to `go_module_fallback`, but does not collapse function/package/user-code visibility on this sample
- `PE` stays bounded at `section_heuristic`, but still preserves the relevant user-facing functions and packages on the sample
- `Mach-O` still has no runtime trust surface, but the current function/package/peeling foothold remains stable under `plain`, `stripped`, `trimpath`, and `pie`
- after fixing the operator-path drift in `task protected-matrix`, `linux/amd64` `plain` and `pie` now also expose the same three bounded module-local files as the other non-garbled profiles

What this does not yet prove:
- behavior on a larger OSS protected target beyond the purpose-built sample

Current measured reading for the new garbled rows:
- the protected matrix now prefers a local source-built upstream `garble` checkout when `/repos/garble` is available
- both `garble` and `garble + literals + tiny` now build and run on `linux/amd64`
- current `GoREveal` output on those rows is:
  - `runtime_trust_summary = "go_module_fallback"`
  - `elf_pclntab_header_magic_kind = "unknown"`
  - `elf_pclntab_header_magic = "08c37341"` for `garble`
  - `elf_pclntab_header_magic = "e79c734f"` for `garble + literals + tiny`
  - `elf_pclntab_function_count_hint = 2125` for `garble`
  - `elf_pclntab_function_count_hint = 3404` for `garble + literals + tiny`
  - `elf_pclntab_file_count_hint = 2265` for `garble`
  - `elf_pclntab_file_count_hint = 225` for `garble + literals + tiny`
  - `elf_functab_pc_offsets_monotonic = true`
  - `elf_functab_last_pc_offset_hint = 733313` for `garble`
  - `elf_functab_last_pc_offset_hint = 1541121` for `garble + literals + tiny`
  - `elf_functab_first_pc_addr_hint = 0x401000` for `garble`
  - `elf_functab_last_pc_addr_hint = 0x4b4011` for `garble`
  - `elf_functab_first_pc_addr_hint = 0x401000` for `garble + literals + tiny`
  - `elf_functab_last_pc_addr_hint = 0x579221` for `garble + literals + tiny`
  - `elf_functab_pc_addr_hints_within_text = true`
  - `elf_functab_pc_addr_sample_count = 8`
  - `elf_functab_pc_addr_sample_first = 0x401000`
  - `elf_functab_pc_addr_sample_all_within_text = true`
  - `elf_function_foothold = "address_only"`
  - `elf_function_foothold_count_hint = 2125` for `garble`
  - `elf_function_foothold_count_hint = 3404` for `garble + literals + tiny`
  - `elf_function_recovery_blocker = "custom_pclntab_magic"`
  - `moduledata_pcheader_matches_gopclntab = true`
  - `moduledata_funcnametab_within_gopclntab = true`
  - `moduledata_pclntable_within_gopclntab = true`
  - `build_info_path = ""`
  - `functions = 0`
  - `packages = 0`
  - `files = 0`
  - `user_functions = 0`
  - `user_packages = 0`
- current baseline posture on those same rows is also weak:
  - `GoReSym` fails
  - `gore` fails
  - `redress` returns `0` functions and `0` packages
- this means the external toolchain blocker is gone; the remaining gap is now a real measured `garble`-class recovery gap
- this also means the gap is no longer “no foothold at all”: the current garbled rows already preserve bounded header-level function/file-count hints, monotonic `functab` PC-offset hints, first and last absolute `PC` address hints within `.text`, a first sampled absolute `PC` foothold, a first explicit analyst-facing `address_only` function foothold, and the existing `moduledata -> pcheader/pclntable` bridges even while named-function recovery collapses
- the first workflow-shaped pain point is now explicit too:
  - the local analyst workflow depends on package/function anchors for `peel`, `diff review sqlite`, `diff handoff sqlite`, and `diff next sqlite`
  - the current `garble` rows preserve a truthful `address_only` foothold, but they do not yet produce a reviewable function/package bundle
  - in practice, this means the current protected gap is not “no evidence”, but “no review-ready anchor set for the existing workflow/value surfaces”
- the first bounded workflow rerun on neighboring garbled builds is now measured too through `.tmp/fresh-eval/protected-workflow-summary.json`:
  - `linux/amd64` `garble` vs `garble + literals + tiny`: `matched_functions = 0`, `transfer_candidates = 0`, `accepted_transfers = 0`, `transfer_packages = 0`
  - the same pair also yields `transfer_review_count = 0`, `transfer_review_packages = 0`, `transfer_review_focus = false`, `handoff_present = false`, `review_plan_count = 0`, `up_next_present = false`, `recommended_actions = 0`, `review_checklist_count = 0`, `review_snapshot_present = false`, and `review_progress_present = false`
  - `linux/arm64` `garble` vs `garble + literals + tiny` yields the same all-zero workflow surface
  - this closes the first `Sprint 16` DEV question as a measured negative result: the current truthful `address_only` foothold does not yet seed any review/transfer/handoff-ready anchor across neighboring garbled builds on either measured Linux architecture

Current measured reading after `arm64` widening:
- all first non-garbled `arm64` rows are now stable on the purpose-built sample:
  - `linux/arm64`: `plain`, `stripped`, `trimpath`, `stripped-trimpath`, and `pie` preserve all three relevant user functions, `2077` functions, `56` packages, and `3` file nodes
  - `windows/arm64`: `plain`, `stripped`, `trimpath`, and `stripped-trimpath` preserve all three relevant user functions, `2072-2073` functions, `57` packages, and `3` file nodes
  - `darwin/arm64`: `plain`, `stripped`, `trimpath`, `stripped-trimpath`, and `pie` preserve all three relevant user functions, `2200-2201` functions, `56` packages, and `3` file nodes
- the new `linux/arm64` garbled rows are now measured too:
  - `runtime_trust_summary = "go_module_fallback"`
  - `elf_pclntab_header_magic_kind = "unknown"`
  - `elf_pclntab_function_count_hint = 2083` for `garble`
  - `elf_pclntab_function_count_hint = 3348` for `garble + literals + tiny`
  - `elf_functab_pc_offsets_monotonic = true`
  - `elf_text_section_addr = 0x11000`
  - `elf_text_section_end_inclusive = 0xb55e3` for `garble`
  - `elf_text_section_end_inclusive = 0x15b2e3` for `garble + literals + tiny`
  - `moduledata_text_addr` remains absent on the current `linux/arm64` family
  - `elf_functab_pc_addr_sample_count = 8`
  - `elf_functab_pc_addr_sample_first = 0x11000`
  - `elf_functab_pc_addr_sample_all_within_text = true`
  - `elf_function_foothold = "address_only"`
  - `elf_function_foothold_count_hint = 2083` for `garble`
  - `elf_function_foothold_count_hint = 3348` for `garble + literals + tiny`
  - `elf_function_foothold_text_source = "elf_text_section"`
  - `elf_function_foothold_start_addr = 0x11000`
  - `elf_function_foothold_end_addr = 0xb55d1` for `garble`
  - `elf_function_foothold_end_addr = 0x15b2d1` for `garble + literals + tiny`
  - `elf_function_recovery_blocker = "custom_pclntab_magic"`
  - `moduledata_pcheader_matches_gopclntab = true`
  - `moduledata_funcnametab_within_gopclntab = true`
  - `moduledata_pclntable_within_gopclntab = true`
  - `functions = 0`
  - `packages = 0`
  - `files = 0`
- this changes the gap reading again:
  - the first `address_only` foothold is now real on both `linux/amd64` and `linux/arm64`
  - the `linux/arm64` gap turned out not to be garble-specific offset collapse, but missing `moduledata_text_*` evidence on the current family
  - the current bounded fix is section-backed rather than `moduledata`-backed, so the product still is not claiming recovered names or decrypted function entries
  - the new compact analyst-facing reading is now explicit too:
    - `linux/amd64` `garble` rows currently project `elf_function_foothold_text_source = "moduledata_text"`
    - `linux/arm64` `garble` rows currently project `elf_function_foothold_text_source = "elf_text_section"`
    - both families now carry a bounded foothold span through `elf_function_foothold_start_addr` / `elf_function_foothold_end_addr`

## Comparative Reading

What the baseline comparison already shows on this same sample:
- `GoREveal` keeps the widest function/package envelope across every measured non-garbled profile
- `GoReSym` still exposes the richest file list across every measured non-garbled profile
- `redress` and `gore` stay much narrower overall, but they still keep all three analyst-relevant user functions visible on this bounded sample
- the protected sample does not currently show an existential recovery failure under `-s -w`, `-trimpath`, or `pie`; the remaining differences are about breadth and source visibility, not basic user-function survival

## Tooling Status

New repo-native tooling now exists:
- `scripts/protected/profile_matrix.py`
- `scripts/protected/test_profile_matrix.py`
- `make protected-matrix`
- `task protected-matrix`

Current operator reality:
- direct dev-container execution of the profile matrix script is working
- the canonical `task protected-matrix` path is now also working through the dev container
- the canonical task now builds a fresh workspace-local `goreveal` binary in the same container step, so the emitted matrix no longer drifts behind current analysis behavior
- the matrix script now prefers a source-built local `/repos/garble` checkout when available, and falls back to the installed binary only when that repo is absent
- the dev image still includes `mvdan.cc/garble` as the fallback path, and the matrix script now has bounded `garble` / `garble + literals + tiny` profile support on `linux/amd64`
- the dev image now also includes `jq`, `yq`, `procps`, and `unzip`, which removes friction for matrix/debug/report inspection inside the canonical container path
- script linting is green
- the first measured pass also forced one real robustness fix in `engine/projection/source_tree.go`, removing a panic on absolute `.../src/main.go` paths discovered by the protected sample
- the old `garble v0.15.0` release gap is now documented separately in `docs/plans/2026-04-01-garble-go126-support-research.md`; it is no longer the active blocker for this lane
- the current compact protected runtime surface is now also locked into the thin `IDA` / `Ghidra` export contracts through `schema/export_test.go` and plugin consumer tests, so export consumers inherit the same bounded `address_only` + text-source + span reading without plugin-side recomputation

## Weighted Next Move

The next correct move is:
1. treat the newly measured `garble` rows as the active evidence baseline, not the earlier toolchain workaround narrative
2. treat the new explicit `elf_function_foothold = "address_only"` as the first truthful protected-binary recovery foothold, now portable across `linux/amd64` and `linux/arm64`, and still bounded to count/address-only semantics
3. treat the new `arm64` widening and its section-backed `.text` range response as completed checkpoints, not future todos
4. treat the first named protected pain point as workflow-shaped:
   - current `garble` rows do not yield a review-ready function/package anchor set for `peel`, `diff review sqlite`, `diff handoff sqlite`, or `diff next sqlite`
5. treat the neighboring-build workflow question as answered for the current sample family:
   - the first bounded rerun now shows no stable transfer/review foothold over the current workflow surfaces on either `linux/amd64` or `linux/arm64`
6. only after that choose between:
   - one stronger but still truthful analyst-facing protected surface
   - transfer workflow refinement
   - deobfuscation orchestration
   - or another bounded recovery slice

This first pass already reduces uncertainty:
- the immediate risk is not that current `GoREveal` surfaces collapse under basic protected build flags
- the immediate unknown is now narrower and more actionable: which next bounded protected response is most product-correct once the neighboring-build workflow question is already a measured “no”
