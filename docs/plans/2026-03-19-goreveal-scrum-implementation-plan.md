# GoREveal Scrum Implementation Plan

> **Status (2026-07-22): historical delivery record.** Active product authority
> is the [RT1 design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md),
> and active file-level execution lives in the
> [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md).

> **Historical execution note:** the checkboxes below record the original
> delivery sequence. Do not start new work from them; use the active RT1 plan.

**Goal:** Build GoREveal as a clean-room Go binary reverse-engineering platform with accuracy-first recovery, schema-driven outputs, differential validation, and selective SIMD acceleration.

**Architecture:** The implementation is split into capability sprints. Early work establishes architecture docs, agent rules, baseline references, Podman-first repository scaffolding, and the first end-to-end recovery pipeline. Later sprints expand semantic recovery, deobfuscation, persistence, integrations, and performance acceleration while preserving a strict schema-first and clean-room boundary.

**Tech Stack:** Go 1.26, pure Go parsing where possible, SQLite, protobuf, `slog`, fuzzing, benchmarks, optional SIMD via architecture-specific implementations and Go 1.26 `simd/archsimd` experiments.

Quantified roadmap checkpoint:
- `docs/plans/2026-03-20-goreveal-progress-assessment.md`
- `docs/plans/2026-03-20-goreveal-next-bounded-analyst-slices-plan.md`
- `docs/plans/2026-03-20-goreveal-deferred-continuation.md`
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md`
- `docs/plans/2026-03-21-goreveal-runtime-modes-and-storage-ideas.md`
- `docs/plans/2026-03-31-goreveal-strategic-review.md`
- `docs/plans/2026-03-31-goreveal-baseline-comparison-plan.md`
- `docs/plans/2026-04-01-goreveal-protected-binary-comparison-plan.md`
- `docs/plans/2026-04-01-goreveal-next-execution-plan.md`
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`
- `docs/architecture/2026-03-31-goreveal-server-stack-decision.md`

Deferred-resume note:
- use `docs/plans/2026-03-20-goreveal-deferred-continuation.md` as the first handoff file before resuming implementation work

## Progress Snapshot

Completed:
- Chunk 0 / Tasks 0.1-0.4
- Sprint 1 / Tasks 1.1-1.2
- Sprint 2 / Tasks 2.1-2.2
- Sprint 3 / Tasks 3.1-3.2
- Sprint 4 / Task 4.1
- Sprint 5 / Tasks 5.1-5.2
- Sprint 6 / Tasks 6.1-6.2
- Sprint 8 / Task 8.2
- Sprint 11 / Task 11.1
- Sprint 11 / Task 11.2

In progress:
- Sprint 7 / Task 7.1 (maintenance lane)
- Sprint 12 / workflow-value lane after protected-binary stabilization
- Sprint 13 / workstation handoff contract hardening
- Sprint 14 / review workflow actionability operator-loop slice

Later ordered horizon:
- Sprint 15 / thin semantic and source confidence
- Sprint 16 / protected commercial Go workflows
- Sprint 17 / server control-plane foundations
- Sprint 18 / metadata and remote interop platform
- Sprint 19 / public release readiness and licensing
- Sprint 20 / evidence expansion and comparative automation
- Sprint 21 / build correlation and version tracking
- Sprint 22 / metadata knowledge network
- Sprint 23 / analyst workspace automation and replay
- Sprint 24 / comparative knowledge packs and decision support

## PM+DEV Sprint Task Model

Within this historical record, scope was tracked through paired PM and DEV tasks:

- `PM-*`
  Outcome-definition, stop-conditions, ranking, and support-policy work
- `DEV-*`
  Bounded implementation slices, verification, and doc-sync work

Task scoring model:
- `Value`
  `1-5`, where `5` means direct operator/product leverage
- `Risk`
  `1-5`, where `5` means likely scope drift or semantic instability
- `Evidence`
  `1-5`, where `5` means the task can be proven by existing tests/fixtures/comparison paths

Execution rule:
- prefer high-`Value`, low-`Risk`, high-`Evidence` tasks first
- if two tasks have similar value, choose the one that leaves behind the clearer stop-condition or operator-visible increment

Current PM+DEV bias:
- the then-active sprint work came from the backlog tables in `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`
- if a candidate task has `Value >= 4`, `Risk <= 2`, and `Evidence >= 4`, it is a default good next move
- if a task has high value but `Risk >= 4`, defer it unless it unblocks the active sprint

Post-Sprint12 sprint reset:
- use `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md` as the current near-term sprint baseline
- treat older `Sprint 13` deobfuscation notes as deferred follow-up, not the active next sprint
- treat `docs/plans/2026-03-22-goreveal-v2-ng-plan.md` as exploratory only, not the active sprint sequence
- treat `transfer_review_plan` plus `goreveal diff next sqlite ...` as the first bounded `Sprint 14` checkpoint, not as speculative backlog

Current implementation notes:
- All development and verification are container-first through Podman.
- repo-local operator entrypoints now exist through both `make ...` and `task ...`.
- `scripts/dev/podman_runner.py` is now the canonical Podman automation layer behind those entrypoints.
- repo-local Codex configuration is now standardized through `.agents/skills/`, `.codex/agents/`, and `.codex/config.toml`.
- script-facing verification now has an explicit strict baseline through `ruff`, `ty`, `yamllint`, and `shellcheck`.
- `task lint-scripts` and `make lint-scripts` are now green through the dev container.
- `analysis.runtime` now also exposes a compact `trust_summary`, so operators can distinguish symbol-backed and fallback/heuristic runtime posture without manually reading many raw fields.
- thin `IDA` and `Ghidra` exports now also mirror canonical `runtime.trust_summary`, keeping the export contract schema-driven and avoiding plugin-side runtime posture inference.
- a first bounded Windows `PE` checkpoint is now landed through a real Go-built `fixture.exe`, `debug/buildinfo` coverage, and `analyze` coverage without introducing any `PE` runtime recovery claim.
- that bounded `PE` checkpoint is now also locked by canonical snapshot coverage and thin export CLI coverage, so cross-format breadth is no longer resting on one narrow unit test path.
- that `PE` checkpoint now also includes the first bounded runtime section heuristic: `analysis.runtime` exposes `.text` / `.rdata` ranges, a raw `.rdata` `pclntab` magic candidate, and one header-looking `.rdata` `pclntab` candidate with raw `magic`, `quantum`, and `pointer_size` fields for the current fixture, while still avoiding any `moduledata` or generic `PE` parser claim.
- the first engine-owned code-peeling MVP slice is now landed: `analysis.peeling` exposes function-level `user | stdlib | runtime | third_party` classification derived only from existing canonical function/build-info truth, plus package-level summaries and per-class function counts, `inspect peeling` is available, `goreveal peel <binary>` now projects the bounded user-only view, and thin exports now mirror the canonical peeling surface without plugin-side inference.
- the next bounded code-peeling refinement is now also landed: function-level peeling output carries explicit `classification_evidence`, and `engine/peeling` now has a small bounded fingerprint-assisted refinement for `runtime` and `stdlib` classification when import-path truth is absent but known name/source fingerprints exist.
- `analyze`, `inspect functions`, and `inspect packages` are working.
- `inspect functions` now exposes bounded `package`, `import_path`, and `module_local` metadata derived from recovered function names plus `build_info` for `main`.
- `inspect functions` now also exposes bounded `source_file` truth from the existing `pclntab`/`gosym` line table, so the function surface is no longer only symbol names plus addresses.
- `inspect functions` now also exposes bounded `source_line` truth from the same `pclntab`/`gosym` line table.
- `inspect functions` now also exposes bounded `autogenerated` truth for synthetic functions whose current truthful source file is `<autogenerated>`.
- `inspect runtime` is now available as a direct CLI surface over the current bounded `analysis.runtime` contract.
- `inspect runtime` now also exposes `firstmoduledata_from_go_module_fallback`, making the current stripped-fixture fallback path explicit.
- `inspect types` and `inspect strings` are working.
- `inspect strings` now also exposes bounded absolute `addr` truth by combining already-known region base addresses with candidate offsets.
- `source-tree` package nodes now also expose explicit `has_file_evidence`, making rich vs fallback-backed package nodes distinguishable without reading implementation notes.
- fresh external reruns now confirm that the current bounded file-evidence path is not fixture-local: measured `ELF`, `PE`, and `Mach-O` targets all expose real file visibility in the current product surface.
- `inspect packages` now also exposes explicit `has_source_evidence`, making source-backed package metadata distinguishable from bounded fallback-backed package metadata without broadening package heuristics.
- `export ida` and `export ghidra` now also preserve canonical string `address`, so thin adapters do not need plugin-side recomputation for string locations.
- `export ida` and `export ghidra` now also preserve bounded function navigation metadata from canonical schema, so thin adapters can consume package/source locality without plugin-side inference.
- `export ida` and `export ghidra` now also preserve bounded type navigation metadata from canonical schema, so thin adapters can consume package/scope/source-evidence context without plugin-side inference.
- thin `IDA` and `Ghidra` adapters now also prefer canonical string `address` over `offset`, keeping runtime-agnostic plugin behavior aligned with the stronger export contract.
- on the stripped ELF fixture, `inspect packages` now preserves a useful `main` package surface through bounded `build_info.path` reuse, and `inspect types` now degrades to `[]` instead of JSON `null`.
- `source-tree` is working.
- Build info recovery is direct via `debug/buildinfo`.
- Initial runtime metadata recovery is now implemented for `ELF` as a read-only spike path.
- Function recovery is currently `ELF`-first via `.gopclntab`.
- Initial type recovery is implemented via `DWARF` as the minimal truthful path for the current fixture.
- Initial string recovery is implemented via `ELF` data-section scanning.
- Source-tree projection is implemented from `DWARF` file paths filtered by `build_info.path`.
- Source-tree projection now groups files by `import_path`, so one recovered package node can own multiple files instead of one file per package node.
- Source-tree package nodes now also carry `function_count`, giving the projection layer more useful package/file metadata without depending on baseline CLIs.
- Source-tree projection now also surfaces external source dependencies as explicit `external_packages` instead of silently dropping non-module `DWARF` file paths.
- Package metadata now exposes `import_path` directly from recovered non-`main` package names, is further enriched with `source_file_count` by correlating recovered packages with source-tree evidence, and now also carries explicit `has_source_evidence` so analysts can see whether package metadata is truly source-backed.
- Function metadata now also exposes `package`, direct non-`main` `import_path`, and bounded `module_local` truth for `main` through `build_info.path`, giving `inspect functions` a real navigation surface instead of only raw symbol names.
- Function metadata now also exposes bounded `source_file` truth from the existing line table, adding one more user-facing navigation signal without broadening parser scope.
- Function metadata now also exposes bounded `source_line` truth from the same line table, keeping the function surface navigable without broadening parser scope.
- Function metadata now also exposes bounded `autogenerated` truth, making synthetic helper/equality functions explicit without inventing new parser semantics.
- String candidates now also expose bounded absolute `addr` truth by combining existing region base addresses with candidate offsets, adding useful analyst navigation without widening the string-recovery parser.
- Type metadata is now enriched with `package` and `user_meaningful`, preserving the raw type list while separating module-local type signal from runtime-heavy background noise.
- Package metadata now also exposes `module_local`, making module-owned packages and external/runtime packages distinct in the canonical package surface.
- Type metadata now also exposes `import_path` and `module_local`, making user-facing type navigation consistent with the richer package surface.
- Type metadata now also exposes `source_file_count`, making the amount of source evidence behind a module-local type classification explicit.
- Runtime metadata now exposes `firstmoduledata_addr`, `.gopclntab`, `.typelink`, and `.go.module` addresses/sizes through `analysis.runtime`.
- Runtime metadata now also exposes raw `typelink_count` as the first bounded bridge from section presence into typelink evidence.
- Runtime metadata now also exposes a bounded `typelink_sample`, giving the first concrete typelink contents without semantic decoding.
- Runtime metadata now also exposes `typelink_min_offset` and `typelink_max_offset`, giving the first raw shape validation over typelink contents.
- Runtime metadata now also exposes `typelink_negative_count` and `typelink_non_negative_count`, giving the first whole-section sign-distribution check over typelink contents.
- Runtime metadata now also exposes a minimal `firstmoduledata`/`.go.module` cross-check and a bounded raw `.go.module` word sample.
- Runtime metadata now also exposes a minimal `moduledata` typelinks slice-header cross-check through `moduledata_typelink_slice_word_index`, `moduledata_typelink_len`, and `moduledata_typelink_cap`.
- Runtime metadata now also exposes `.itablink` section evidence and a bounded `moduledata` itablinks slice-header cross-check.
- Runtime metadata now also exposes a bounded `moduledata` memory-range block cross-check for `.noptrdata`, `.data`, `.bss`, and `.noptrbss`.
- Runtime metadata now also exposes a bounded `.rodata` range cross-check from `firstmoduledata`.
- Runtime metadata now also exposes a bounded `.text` range cross-check from `firstmoduledata`, with explicit inclusive-end semantics.
- Deobfuscation pipeline scaffolding is implemented and refined-layer separation is verified.
- First real deobfuscation passes are implemented for synthetic function-name refinement and string-segment extraction.
- `deobfuscate` CLI exposure is implemented.
- SQLite-backed analysis persistence is implemented with a pure-Go driver and `export sqlite` CLI support.
- Differential validation v1 is implemented for one real fixture using normalized `GoReSym` and `redress` runners.
- Differential validation now includes `GoReSym` overlap for build info, source files, and a narrow user-function set.
- Differential validation now includes a normalized `redress` runner for module-path, package, source-file, and user-function overlap.
- Differential validation now also includes a normalized `gore` runner for build info, package, source-file, and function overlap.
- `make test-differential` is green inside the Podman dev container.
- `make test` now includes the baseline-backed differential Go package in its main regression pass instead of silently skipping it.
- a shared Python normalization layer now covers the active baseline runners, with `GoReSym`, `redress`, and `gore` moved onto it and direct unit coverage added for normalization logic
- a machine-readable differential report path now exists through `scripts.baseline.generate_fixture_report` and `make test-differential-report`
- Python integration paths that need fixture-driven CLI output now use a built `GOREVEAL_BIN` in `.tmp/goreveal` instead of flaky nested `go run` calls
- The first native capability-transfer slice from `redress` is now in progress: source-tree projection has moved from per-file package nodes to grouped package/file aggregation.
- The second native capability-transfer slice from `redress`/`gore`-style package metadata is now landed: grouped source-package nodes are enriched with package-level `function_count`.
- The third native capability-transfer slice from `redress`-style source visibility is now landed: module-local source projection is complemented by explicit `external_packages` for runtime/stdlib or other non-module source paths.
- The first native capability-transfer slice from `gore`-style package metadata is now landed: recovered packages expose non-`main` `import_path` directly from the function surface and are further enriched with `source_file_count` from source-tree evidence.
- The second native capability-transfer slice from `gore`-style user-surface cleanup is now landed: raw types are enriched with `package` and `user_meaningful` instead of being destructively filtered.
- The third native capability-transfer slice from `gore`-style package clarity is now landed: recovered packages now carry `module_local`, making module-owned and external packages explicitly distinguishable.
- The fourth native capability-transfer slice from `gore`-style type clarity is now landed: raw types now carry `import_path` and `module_local`, so type navigation no longer has to infer scope from short package tokens alone.
- The fifth native capability-transfer slice from `gore`-style type evidence is now landed: raw types now carry `source_file_count`, exposing how much source-tree evidence backs the module-local type metadata.
- The first `Sprint 12` runtime slice is now landed as a read-only metadata probe: `analysis.runtime` carries low-risk ELF runtime layout evidence without changing existing recovery semantics.
- The second `Sprint 12` runtime slice is now landed: `analysis.runtime.typelink_count` exposes bounded typelink evidence without attempting full typelink decoding yet.
- The third `Sprint 12` runtime slice is now landed: `analysis.runtime.typelink_sample` exposes raw typelink offsets without attempting full typelink decoding yet.
- The fourth `Sprint 12` runtime slice is now landed: `analysis.runtime.typelink_min_offset` and `analysis.runtime.typelink_max_offset` expose bounded typelink shape evidence without semantic decoding.
- The fifth `Sprint 12` runtime slice is now landed: `analysis.runtime.typelink_negative_count` and `analysis.runtime.typelink_non_negative_count` expose a whole-section typelink sign-distribution check without semantic decoding.
- The sixth `Sprint 12` runtime slice is now landed: `analysis.runtime.firstmoduledata_in_go_module`, `analysis.runtime.firstmoduledata_go_module_offset`, and `analysis.runtime.go_module_word_sample` expose the first bounded `moduledata`/`.go.module` consistency bridge.
- The seventh `Sprint 12` runtime slice is now landed: `analysis.runtime.moduledata_typelink_slice_word_index`, `analysis.runtime.moduledata_typelink_len`, and `analysis.runtime.moduledata_typelink_cap` expose the first bounded `moduledata` typelinks slice-header parse.
- The eighth `Sprint 12` runtime slice is now landed: `analysis.runtime.itablink_addr`, `analysis.runtime.itablink_size`, `analysis.runtime.itablink_count`, `analysis.runtime.moduledata_itablink_slice_word_index`, `analysis.runtime.moduledata_itablink_len`, and `analysis.runtime.moduledata_itablink_cap` expose the first bounded `moduledata` itablinks slice-header parse.
- The ninth `Sprint 12` runtime slice is now landed: `analysis.runtime.moduledata_memory_ranges_word_index`, `analysis.runtime.moduledata_noptrdata_addr`, `analysis.runtime.moduledata_data_addr`, `analysis.runtime.moduledata_bss_addr`, and `analysis.runtime.moduledata_noptrbss_addr` (plus matching end fields) expose the first bounded `moduledata` memory-range block parse.
- The tenth `Sprint 12` runtime slice is now landed: `analysis.runtime.moduledata_rodata_word_index`, `analysis.runtime.moduledata_rodata_addr`, and `analysis.runtime.moduledata_rodata_end` expose the first bounded `.rodata` range parse from `firstmoduledata`.
- The eleventh `Sprint 12` runtime slice is now landed: `analysis.runtime.moduledata_text_word_index`, `analysis.runtime.moduledata_text_addr`, and `analysis.runtime.moduledata_text_end_inclusive` expose the first bounded `.text` range parse from `firstmoduledata`.
- The twelfth `Sprint 12` slice is now landed as the first tiny semantic decode: `analysis.runtime.typelink_resolved_base_addr`, `analysis.runtime.typelink_resolved_sample`, and `analysis.runtime.typelink_resolved_within_rodata_count` expose a bounded rodata-relative typelink resolution hypothesis for the canonical fixture.
- The thirteenth `Sprint 12` slice is now landed as a semantic confidence bit: `analysis.runtime.typelink_all_resolved_within_rodata` confirms that the current rodata-relative typelink hypothesis is fully self-consistent on the canonical fixture.
- The fourteenth `Sprint 12` slice is now landed as a bounded `types..etypes` bridge: `analysis.runtime.moduledata_types_range_word_index`, `analysis.runtime.moduledata_types_addr`, `analysis.runtime.moduledata_etypes_addr`, `analysis.runtime.typelink_resolved_within_types_count`, and `analysis.runtime.typelink_all_resolved_within_types` extend the current fixture-local semantic hypothesis without broadening it into a general typelink/type decoder.
- The fifteenth `Sprint 12` slice is now landed as a bounded `.gopclntab` semantic bridge: `analysis.runtime.moduledata_pcheader_addr`, `analysis.runtime.moduledata_pcheader_matches_gopclntab`, `analysis.runtime.moduledata_funcnametab_slice_word_index`, `analysis.runtime.moduledata_funcnametab_addr`, `analysis.runtime.moduledata_funcnametab_len`, `analysis.runtime.moduledata_funcnametab_cap`, and `analysis.runtime.moduledata_funcnametab_within_gopclntab` extend the current fixture-local runtime model without turning it into a generic pcln-table decoder.
- The sixteenth `Sprint 12` slice is now landed as the next bounded `.gopclntab` semantic bridge: `analysis.runtime.moduledata_cutab_slice_word_index`, `analysis.runtime.moduledata_cutab_addr`, `analysis.runtime.moduledata_cutab_len`, `analysis.runtime.moduledata_cutab_cap`, and `analysis.runtime.moduledata_cutab_within_gopclntab` extend the same fixture-local runtime model without turning it into a general `cutab` decoder.
- The seventeenth `Sprint 12` slice is now landed as the next bounded `.gopclntab` semantic bridge: `analysis.runtime.moduledata_filetab_slice_word_index`, `analysis.runtime.moduledata_filetab_addr`, `analysis.runtime.moduledata_filetab_len`, `analysis.runtime.moduledata_filetab_cap`, and `analysis.runtime.moduledata_filetab_within_gopclntab` extend the same fixture-local runtime model without turning it into a general `filetab` decoder.
- The eighteenth `Sprint 12` slice is now landed as the next bounded `.gopclntab` semantic bridge: `analysis.runtime.moduledata_pctab_slice_word_index`, `analysis.runtime.moduledata_pctab_addr`, `analysis.runtime.moduledata_pctab_len`, `analysis.runtime.moduledata_pctab_cap`, and `analysis.runtime.moduledata_pctab_within_gopclntab` extend the same fixture-local runtime model without turning it into a general `pctab` decoder.
- The nineteenth `Sprint 12` slice is now landed as the next bounded `.gopclntab` semantic bridge: `analysis.runtime.moduledata_pclntable_slice_word_index`, `analysis.runtime.moduledata_pclntable_addr`, `analysis.runtime.moduledata_pclntable_len`, `analysis.runtime.moduledata_pclntable_cap`, and `analysis.runtime.moduledata_pclntable_within_gopclntab` extend the same fixture-local runtime model without turning it into a general `pclntable` decoder.
- After these semantic slices, the next handoff still remains inside `Sprint 12`, but with a new checkpoint rule: the project should prefer a second fixture or a very small runtime-to-heuristic cross-check over another blind same-fixture `.gopclntab` bridge.
- That second-fixture move is now started: the stripped ELF fixture keeps `.go.module` while losing `runtime.firstmoduledata` from the rich symbol table, and `ReadMetadata()` now has a bounded `.go.module`-based fallback for this current ELF family.
- That second-fixture move now also makes the fallback source explicit through `analysis.runtime.firstmoduledata_from_go_module_fallback`, so stripped-fixture `firstmoduledata_addr` is no longer opaque to operators.
- That second-fixture move now also includes the first bounded runtime-to-heuristic UX cross-check: on the stripped fixture, `inspect packages` preserves a useful `main` package surface through `build_info.path` while still exposing external package `import_path` directly from function recovery, and `inspect types` degrades to `[]` instead of JSON `null`.
- That second-fixture move now also includes a bounded source-tree fallback: when DWARF file evidence is absent but build info and package truth still exist, `source-tree` returns a root, module-local package nodes, and external package nodes with empty file lists instead of failing outright.
- The first semantic typelink bridge does not yet change package/type heuristics: current package/type scope and usefulness metadata remain source-tree/DWARF-driven until runtime-semantic decoding yields stable naming or package truth.
- Package metadata now also has a bounded build-info fallback for `main` when source-tree evidence is absent, preserving a minimal module-local package surface on stripped fixtures.
- Type metadata now also has the same bounded build-info fallback for `main` when a truthful type surface exists without source-tree evidence, while non-`main` types keep direct `import_path` truth from their parsed type package.
- Canonical `analyze` output now exposes `types: []` when no truthful type surface exists, instead of silently omitting the field.
- Initial plugin-ready export contracts are implemented for `IDA` and `Ghidra`.
- `export ida <binary>` and `export ghidra <binary>` are available as stable v1 JSON payloads.
- Stored-run diffing is implemented for SQLite-backed analyses through `goreveal diff sqlite <db> <left-id> <right-id>`.
- stored-run diffing now also carries the first bounded version-tracking-adjacent function-matching surface: `matched_functions` records exact-name, source-location, source-file, and module-local normalized-name matches with `score`, `reason`, and optional peeling class context, without turning `storage/diff` into a decompiler-dependent matcher.
- stored-run diffing now also carries the first package-level transfer summary surface: `transfer_packages` aggregates candidate, ready, review, and accepted counts over the existing bounded transfer contract instead of introducing a broader matcher.
- A thin, container-testable `IDA` adapter is now started as an action-builder over the stable export contract.
- The `IDA` adapter now has fixture-driven validation against real `export ida` output from the canonical Go fixture.
- The `Ghidra` adapter now mirrors the same thin, fixture-driven contract pattern over `export ghidra`.
- `Sprint 12` has now accumulated enough bounded runtime density on the canonical fixture that the next strategic step should be the first very small semantic decode rather than another proof-only slice by default.
- that first semantic step is now started and should remain tightly fixture-scoped until broader evidence exists.
- the next handoff after this first semantic step should be another tiny semantic-runtime slice, not an early rewrite of package/type heuristics.
- `Sprint 7` remains important, but is now a maintenance lane for claim hygiene rather than the main execution lane.

## Criticality Snapshot

The project no longer has a “missing foundation” problem. It now has a “prove and stabilize what exists” problem. That changes sprint criticality.

## Priority Triage Snapshot

The current priority model is:
- `criticality`: how much project truth or downstream safety depends on the work
- `importance`: how much the work advances first-class platform value
- `speed`: how quickly the work can produce a demonstrable, low-regret increment

Current top candidates from this state:
- `Sprint 12 post-summary export decision and second-fixture checkpoint`
  - criticality: very high
  - importance: very high
  - speed: medium
- `Sprint 7 maintenance/evidence hygiene`
  - criticality: medium
  - importance: medium-high
  - speed: medium-high
- `Sprint 11 follow-on source/package refinement`
  - criticality: medium
  - importance: medium
  - speed: medium
- `Sprint 10 service/API`
  - criticality: medium
  - importance: high
  - speed: low
- `Sprint 9 SIMD/performance`
  - criticality: low
  - importance: medium
  - speed: low

`P0: highest current criticality`
- keep the new compact runtime trust/evidence summary stable across the bounded runtime and export surfaces
- choose the next bounded Windows `PE` slice beyond the current section/header heuristic only when it stays clearly evidence-backed
- keep the new `PE` runtime section heuristic clearly labeled as `section_heuristic`, not semantic runtime truth
- avoid reopening package/type heuristic work before the `PE` checkpoint strengthens evidence breadth

`P1: important after P0`
- keep `Sprint 7` healthy enough that new claims stay evidence-backed
- reopen `Sprint 11` only if new runtime truth changes package/type/source classification semantics
- widen code peeling only through bounded engine-owned slices over canonical truth
- continue service/API only after runtime-semantic truth is less heuristic
- keep repo automation, agent configuration, and script-lint policy synchronized as supporting infrastructure rather than a competing execution lane

`P2: deliberately deferred`
- aggressive SIMD/performance work
- rich plugin adapters
- release polish beyond current operator-grade docs
- function-level version tracking until code peeling MVP exists
- metadata-network work until code peeling and version-tracking foundations exist
- dual runtime mode work, `gorectl`, and server storage architecture remain valid future epics, but stay behind the current runtime/type accuracy lane
- MCP agent surfaces and object-store-backed artifact transfer are valid future platform epics, but should remain downstream of the core server/API decision rather than lead it

Recommended execution order from the current state:
1. keep the new runtime trust/evidence summary stable across runtime and export surfaces
2. keep `Sprint 7` as maintenance/evidence hygiene
3. treat the bounded Windows `PE` checkpoint as landed and deepen it only if one more very small evidence-backed slice is clearly justified
4. keep widening code peeling and transfer surfaces only through bounded engine/storage-owned slices over canonical truth
5. treat the protected-binary comparison lane as active, now with the old `garble v0.15.0` release-gap resolved for planning purposes when a local upstream checkout is available
6. treat the first bounded `elf_function_foothold = "address_only"` surface as landed on the current garbled `linux/amd64` rows
7. treat the `arm64` widening checkpoint as landed and the bounded Linux-architecture portability fix for `address_only` footholds as landed too
8. treat the compact protected runtime surface as stabilized across runtime, export-contract, and plugin-consumer boundaries
9. return by default to workflow/value work unless a newly measured protected-specific analyst pain point clearly outranks it
10. treat the first compact `transfer_review` queue and explicit `projected_package` transfer projection as the first concrete workflow/value checkpoint after protected-binary stabilization
11. treat the new package-first `transfer_review_packages` triage surface as the next bounded workflow/value checkpoint over the same existing transfer state, not as a new matcher lane
12. treat the new `transfer_review_focus` first-pass bundle as the first explicit recommended next step over that existing review state
13. treat the bounded CLI/operator projection for that focused review pass as landed through `goreveal diff review sqlite ...`
14. use the measured `rehelp` and RE-lab inventory as the current workstation/interop baseline, not as a reason to expand `core`
15. treat the compact machine-readable `handoff` block on `goreveal diff review sqlite ...` as the first workstation-facing review bridge, not the final interop shape
16. treat `goreveal diff handoff sqlite ...` as the first dedicated operator-facing handoff projection over that bridge
17. after landing that dedicated handoff projection, harden explicit host-platform MCP and workstation handoff planning around the now-measured `ida-pro-mcp`, `Diaphora`, `BinExport`, `rizin`, and dynamic/symbolic sidecars
18. keep service/API and performance behind accuracy work
19. use the fresh external comparison and universal-workbench comparison as the current product baseline for deciding whether another semantic slice truly outranks workflow/value and interop work
20. after the later server and remote-interop horizon, separate public-release/licensing hardening from evidence/comparison automation instead of blending them into one catch-all sprint

Current repo-ops checkpoint:
- Codex-native skills and subagents are now established as a portable repo contract.
- `Taskfile.yml` and `scripts/dev/podman_runner.py` now provide the main operator UX on top of the existing `Makefile`.
- strict Python/YAML/shell verification is now part of the expected dev-container workflow, not an ad hoc local-only check.

## Product/Risk Snapshot

Current product reading:
- `Sprint 7` now gives us a credible proof layer for the canonical fixture.
- `Sprint 11` is the first real move from “baseline comparison” to “native capability absorption”.
- the source-tree grouping change is small in code size but high-signal in product terms because it makes `GoREveal` projection structurally closer to a useful `redress`-style source view.
- the new `external_packages` projection is high-signal because it stops hiding runtime and stdlib source evidence from the user; this is useful product truth even before deeper source reconstruction exists.
- enriching packages with `import_path` and `source_file_count` is high-signal because it upgrades `inspect packages` from a presence/count surface into a more useful navigation surface without taking on typelink noise yet.
- enriching raw types with `package` and `user_meaningful` is high-signal because it gives operators a low-regret way to distinguish module-local type signal from runtime noise without losing raw evidence.
- enriching packages with `module_local` is high-signal because it turns the package list into a true navigation surface instead of forcing the operator to infer scope from names like `runtime` or `main`.
- enriching types with `import_path` and `module_local` is high-signal because it aligns the type surface with the package surface and makes nested-package types easier to reason about without filtering them away.
- enriching types with `source_file_count` is high-signal because it exposes how much source evidence backs the current metadata and makes the heuristic boundary visible instead of hidden.

Current risks:
- source projection is still ultimately `DWARF`-path driven, so it is deeper than before but not yet a full source-reconstruction engine
- package naming and package/file grouping still rely on relatively simple heuristics; this is good enough for the first transfer slice, but not yet broad parity
- package-level metadata is now richer, but it still depends on function-name-derived package recovery correlated with source-tree evidence rather than independent package metadata sources
- package-level metadata is now richer, but it still depends on function-name-derived package recovery correlated with source-tree evidence rather than independent package metadata sources
- external package classification is intentionally heuristic; standard-library paths are strong, but module-cache and non-standard source layouts are not yet normalized into a richer dependency model
- type-surface transfer is now less noisy and more navigable at the metadata level, but the underlying recovery source is still `DWARF` and not yet typelink/moduledata-driven
- some module-local decisions are still inferred from source-tree correlation, so stripped binaries with poorer file evidence may need a stronger future source of truth
- the next risk-reduction step should come from runtime metadata, not from adding more schema-only fields
- `analysis.runtime` is currently ELF-only and partly symbol-table-assisted, so it is a spike-quality truth source, not yet a cross-version compatibility claim
- `typelink_count` is currently derived from `.typelink` section size, so it is useful evidence but not yet a decoded semantic truth surface
- `typelink_sample` is raw offset evidence only; it improves observability, but it is not yet validated against decoded runtime type structures
- `typelink_min_offset` and `typelink_max_offset` improve plausibility checks, but they still do not validate the offsets against a decoded `moduledata.types` range
- `typelink_negative_count` and `typelink_non_negative_count` give better structural confidence, but they still do not validate typelinks against decoded runtime ranges
- the new `firstmoduledata`/`.go.module` bridge reduces ambiguity, but it is still only a bounded consistency check, not a parsed `moduledata` truth source
- the new `moduledata` typelinks slice-header parse is high-signal, but it is still a narrow cross-check against `.typelink` evidence, not yet a general `moduledata` field-decoding layer
- the new `itablinks` bridge is another high-signal consistency check, but it still does not imply broad semantic decoding of `itab` contents or cross-version `moduledata` support
- the new memory-range block bridge raises confidence in current `moduledata` positioning, but it is still bounded to one fixture/layout family and should not be overclaimed as generic memory-layout decoding
- the new `.rodata` bridge further increases bounded confidence in current `moduledata` range layout, but it is still a fixture-local cross-check rather than general code/data mapping support
- the new `.text` bridge is useful because it is explicit about inclusive-end semantics, but it still should not be stretched into a general executable-range or text-mapping claim beyond the current fixture/layout
- the project now risks over-investing in additional bounded bridges when it already has enough runtime density to justify one carefully chosen semantic decode
- the first semantic typelink bridge is promising, but it is still hypothesis-quality and must not be mistaken for general typelink/type decoding support
- the new all-within-rodata confidence bit strengthens the current fixture-local hypothesis, but it still must not be stretched into a cross-version or cross-format truth claim
- the new `types..etypes` bridge and `moduledata_types_range_word_index` add one more bounded semantic check, but they are still fixture-local and must not yet be mistaken for generic `moduledata.types` decoding support
- the new `pcHeader` / `funcnametab` bridge is high-signal, but it is still fixture-local and must not yet be mistaken for generic `pcHeader`, `funcnametab`, or pcln-table decoding support
- the new `cutab` bridge further raises confidence in the current pcln-layout hypothesis, but it is still fixture-local and must not yet be mistaken for generic `cutab` or pcln-table decoding support
- the new `filetab` bridge further raises confidence in the current pcln-layout hypothesis, but it is still fixture-local and must not yet be mistaken for generic `filetab` or pcln-table decoding support
- the new `pctab` bridge further raises confidence in the current pcln-layout hypothesis, but it is still fixture-local and must not yet be mistaken for generic `pctab` or pcln-table decoding support
- the new `pclntable` bridge further raises confidence in the current pcln-layout hypothesis, but it is still fixture-local and must not yet be mistaken for generic `pclntable` or pcln-table decoding support
- the current `.gopclntab` bridge chain is now dense enough that continued same-fixture expansion risks parser growth with diminishing product value
- the current semantic layer is still not strong enough to justify rewriting package/type heuristics or to claim typelinks-driven type recovery
- the first semantic typelink bridge is not yet strong enough to replace the current source-tree/DWARF-backed package/type heuristics, so changing those heuristics now would overfit the product to one fixture
- the imported `gobfd` `golangci-lint` policy is now fully integrated and `make lint` is green; future hygiene work should keep that baseline green instead of reopening config churn
- the protected-binary lane now has a first explicit garble-collapse explanation surface: garbled ELF rows preserve bounded runtime posture, expose `unknown` ELF `pclntab` magic, preserve nonzero header-level function/file-count hints, preserve monotonic `functab` PC-offset hints, preserve bounded absolute `PC` address hints within `.text`, preserve a first sampled absolute `PC` foothold, preserve the compact analyst-facing `elf_function_foothold = "address_only"` surface with its count hint on both `linux/amd64` and `linux/arm64`, preserve a compact projection of whether that foothold is backed by `moduledata_text` or `elf_text_section`, preserve the current `moduledata` bridges, and carry `elf_function_recovery_blocker = "custom_pclntab_magic"` when function recovery collapses

---

## Chunk 0: Documentation, Repo Foundations, and Clean-Room Controls

### Task 0.1: Finalize the initial architecture and planning package

**Files:**
- Maintain: `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- Maintain: `docs/architecture/2026-03-19-goreveal-go126-best-practices.md`
- Create: `docs/architecture/2026-03-19-goreveal-module-map.md`
- Create: `docs/architecture/2026-03-19-goreveal-schema-principles.md`
- Create: `docs/architecture/2026-03-19-goreveal-testing-strategy.md`
- Maintain: `docs/plans/2026-03-19-goreveal-scrum-implementation-plan.md`

- [x] Write the module map with exact repository boundaries and dependency rules.
- [x] Write schema principles covering raw truth, refined truth, provenance, and confidence.
- [x] Write testing strategy covering corpus, snapshots, differential tests, fuzzing, and benchmarks.
- [x] Review all architecture docs for consistency with the platform contract.

### Task 0.2: Create agent-facing docs and project skills

**Files:**
- Create: `AGENTS.md`
- Create: `CLAUDE.md`
- Create: `CODEX.md`
- Create: `GEMINI.md`
- Create: `skills/README.md`
- Create: `skills/goreveal-navigation/SKILL.md`
- Create: `skills/goreveal-cleanroom/SKILL.md`
- Create: `skills/goreveal-corpus-validation/SKILL.md`
- Create: `skills/goreveal-differential-testing/SKILL.md`
- Create: `skills/goreveal-deobfuscation/SKILL.md`
- Create: `skills/goreveal-perf-simd/SKILL.md`
- Create: `skills/goreveal-export-contracts/SKILL.md`
- Create: `skills/goreveal-release-ops/SKILL.md`

- [x] Write `AGENTS.md` as the single operational contract for clean-room work, architecture invariants, and definition of done.
- [x] Write short agent overlays that reference `AGENTS.md` instead of duplicating it.
- [x] Write the mandatory skills with concrete workflows and guardrails.
- [x] Ensure every skill points back to the architecture contract and best-practices doc.

### Task 0.3: Acquire baseline repositories and define the reference set

**Files:**
- Create: `docs/architecture/2026-03-19-goreveal-baseline-sources.md`
- Create: `corpus/baseline/README.md`
- Create: `scripts/baseline/README.md`

- [x] Verify the exact upstream repositories for `gore`, `redress`, `GoReSym`, `GoResolver`, `gostringungarbler`, and `AlphaGolang`.
- [x] Fork missing upstream repositories into the target GitHub namespace using `gh api`.
- [x] Clone missing baseline repositories into `/opt/projects/repositories`.
- [x] Document what each baseline is used for: parser truth, source projection, deobfuscation, plugin behavior, or differential comparison.

### Task 0.4: Initialize the repo as a Podman-first Go workspace monorepo

**Files:**
- Create: `go.work`
- Create: `go.work.sum`
- Create: `.gitignore`
- Create: `README.md`
- Create: `Makefile`
- Create: `.golangci.yml`
- Create: `tools.go`
- Create: `go.mod`
- Create: `internal/version/version.go`
- Create: `internal/testutil/testutil.go`
- Create: `deployments/docker/Containerfile.dev`
- Create: `deployments/docker/Containerfile.builder`
- Create: `deployments/docker/Containerfile.release`
- Create: `deployments/docker/README.md`

- [x] Initialize a git repository if one does not exist yet.
- [x] Create a top-level Go workspace and root module or tooling module.
- [x] Add standard repo automation for Podman-based format, lint, test, fuzz, and bench entrypoints.
- [x] Rebase `.golangci.yml` on the `gobfd` rule set and keep it in staged-adoption mode for `goreveal`.
- [x] Add version and shared test utilities.

---

## Chunk 1: Sprint 1 - Analysis Skeleton

### Task 1.1: Scaffold the core, schema, engine, and CLI modules

**Files:**
- Create: `core/go.mod`
- Create: `core/doc.go`
- Create: `core/ingest/file.go`
- Create: `core/format/format.go`
- Create: `schema/go.mod`
- Create: `schema/analysis.go`
- Create: `schema/provenance.go`
- Create: `engine/go.mod`
- Create: `engine/engine.go`
- Create: `cmd/goreveal/go.mod`
- Create: `cmd/goreveal/main.go`
- Create: `cmd/goreveal/internal/analyze.go`

**Test Files:**
- Create: `core/ingest/file_test.go`
- Create: `core/format/format_test.go`
- Create: `engine/engine_test.go`
- Create: `cmd/goreveal/internal/analyze_test.go`

- [x] Write the failing tests for opening a binary and detecting `ELF`, `PE`, and `Mach-O`.
- [x] Implement the minimal ingest and format-detection path.
- [x] Define the smallest viable `schema.Analysis` structure.
- [x] Wire `goreveal analyze <binary>` to emit canonical JSON.
- [x] Run `go test ./...` and store the first golden snapshot fixture.

### Task 1.2: Create the first corpus fixture and snapshot pipeline

**Files:**
- Create: `corpus/fixtures/README.md`
- Create: `corpus/fixtures/minimal-linux-amd64/fixture.json`
- Create: `corpus/fixtures/minimal-linux-amd64/expected.analysis.json`
- Create: `tests/snapshots/analyze_snapshot_test.go`

- [x] Add one known-good fixture binary and its metadata.
- [x] Write snapshot tests around `schema.Analysis` output.
- [x] Make snapshot update flow explicit in docs or Makefile targets.

---

## Chunk 2: Sprint 2 - Function and Package Recovery

### Task 2.1: Recover build info, pclntab, and functions

**Files:**
- Create: `core/buildinfo/buildinfo.go`
- Create: `core/buildinfo/buildinfo_test.go`
- Create: `core/pclntab/pclntab.go`
- Create: `core/pclntab/pclntab_test.go`
- Create: `core/functions/functions.go`
- Create: `core/functions/functions_test.go`
- Modify: `schema/analysis.go`
- Modify: `engine/engine.go`

- [x] Write failing tests for build info and function enumeration on corpus fixtures.
- [x] Implement build info recovery with explicit provenance.
- [x] Implement pclntab parsing and function extraction.
- [x] Map recovered functions into schema with raw offsets and confidence metadata.

### Task 2.2: Recover package information and expose inspection CLI

**Files:**
- Create: `core/packages/packages.go`
- Create: `core/packages/packages_test.go`
- Create: `cmd/goreveal/internal/inspect_functions.go`
- Create: `cmd/goreveal/internal/inspect_packages.go`
- Modify: `cmd/goreveal/main.go`

- [x] Write failing tests for package classification.
- [x] Implement package recovery and classification hooks.
- [x] Add `inspect functions` and `inspect packages` commands.
- [x] Verify CLI output shape against golden expectations.

---

## Chunk 3: Sprint 3 - Types and Strings

### Task 3.1: Add typelink and type recovery

**Files:**
- Create: `core/types/typelinks.go`
- Create: `core/types/types.go`
- Create: `core/types/types_test.go`
- Modify: `schema/analysis.go`
- Modify: `engine/engine.go`

- [x] Write failing tests for type recovery on supported fixtures.
- [x] Implement initial type model extraction for the current `ELF` fixture via `DWARF`; typelink traversal remains future work.
- [x] Map types into schema with provenance and confidence.

### Task 3.2: Add string region and string candidate recovery

**Files:**
- Create: `core/strings/regions.go`
- Create: `core/strings/candidates.go`
- Create: `core/strings/strings_test.go`
- Create: `cmd/goreveal/internal/inspect_types.go`
- Create: `cmd/goreveal/internal/inspect_strings.go`

- [x] Write failing tests for string candidate extraction.
- [x] Implement region scanning and candidate extraction.
- [x] Add `inspect types` and `inspect strings` commands.
- [x] Add corpus fixtures with string-heavy binaries.

---

## Chunk 4: Sprint 4 - Source Projection v1

### Task 4.1: Build package and file projection

**Files:**
- Create: `engine/projection/source_tree.go`
- Create: `engine/projection/source_tree_test.go`
- Create: `cmd/goreveal/internal/source_tree.go`
- Modify: `schema/analysis.go`

- [x] Write failing tests for source-tree projection on known fixtures.
- [x] Implement package-to-file projection using recovered package and type metadata.
- [x] Emit projection into schema and CLI.
- [x] Store expected source-tree snapshots.

---

## Chunk 5: Sprint 5 - Deobfuscation v1

### Task 5.1: Scaffold deobfuscation pipeline with refined layers

**Files:**
- Create: `deobfuscation/go.mod`
- Create: `deobfuscation/pipeline.go`
- Create: `deobfuscation/pipeline_test.go`
- Create: `schema/refined.go`
- Modify: `engine/engine.go`

- [x] Write tests proving raw truth and refined truth remain separate.
- [x] Add deobfuscation pass interface and pipeline orchestration.
- [x] Extend schema to carry refined names and provenance.

### Task 5.2: Implement string ungarbling and first name refinement pass

**Files:**
- Create: `deobfuscation/garble/strings.go`
- Create: `deobfuscation/garble/strings_test.go`
- Create: `deobfuscation/refine/names.go`
- Create: `deobfuscation/refine/names_test.go`
- Create: `cmd/goreveal/internal/deobfuscate.go`

- [x] Write failing fixtures for garble-like cases.
- [x] Implement string ungarbling v1.
- [x] Implement name refinement hooks.
- [x] Add CLI exposure for deobfuscation.

---

## Chunk 6: Sprint 6 - SQLite Analysis Store

### Task 6.1: Add SQLite-backed analysis persistence

**Files:**
- Create: `storage/go.mod`
- Create: `storage/sqlite/store.go`
- Create: `storage/sqlite/schema.sql`
- Create: `storage/sqlite/store_test.go`
- Create: `cmd/goreveal/internal/export_sqlite.go`

- [x] Write failing tests for saving and loading analyses.
- [x] Implement SQLite schema for canonical analysis persistence.
- [x] Add CLI export to SQLite.
- [x] Verify round-trip integrity against snapshot outputs.

### Task 6.2: Add comparison of stored analysis runs

**Files:**
- Create: `storage/diff/diff.go`
- Create: `storage/diff/diff_test.go`
- Create: `cmd/goreveal/internal/diff.go`

- [x] Write tests for comparing two stored analysis runs.
- [x] Implement diff logic at schema level.
- [x] Add CLI diff output for stored runs or baseline comparisons.

---

## Chunk 7: Sprint 7 - Differential Validation Framework

### Task 7.1: Create baseline runners and comparison harness

**Files:**
- Create: `scripts/baseline/run_gore.sh`
- Create: `scripts/baseline/run_redress.sh`
- Create: `scripts/baseline/run_goresym.sh`
- Create: `scripts/baseline/run_goresolver.sh`
- Create: `scripts/baseline/run_gostringungarbler.sh`
- Create: `tests/differential/differential_test.go`
- Create: `tests/differential/fixtures/README.md`
- Create: `docs/architecture/2026-03-19-goreveal-differential-v1-notes.md`

- [x] Write normalized baseline runner wrappers for the first overlap set: `GoReSym` and `redress`.
- [x] Write differential tests that compare `GoREveal` against baseline outputs where behavior should overlap for build info, source-tree files, and package presence.
- [x] Document allowed divergences and better-than-baseline cases for the current v1 overlap set.
- [x] Finish consolidating the growing shell baseline runners into a shared Python normalization layer once wrapper complexity exceeds thin orchestration.
- [x] Add a machine-readable differential report for the current fixture and overlap set.
- [ ] Expand wrappers and comparisons to the next overlap set: `GoResolver` and `gostringungarbler`.

---

## Chunk 8: Sprint 8 - Plugin-Ready Exports and Integrations

### Task 8.1: Stabilize export contracts for IDA and Ghidra

**Files:**
- Create: `schema/export_ida.go`
- Create: `schema/export_ghidra.go`
- Create: `schema/export_test.go`
- Create: `cmd/goreveal/internal/export_ida.go`
- Create: `cmd/goreveal/internal/export_ghidra.go`

- [x] Write tests for stable export payload shape.
- [x] Implement export serializers for IDA and Ghidra consumers.
- [x] Validate output on fixture-driven import simulations.

### Task 8.2: Add thin plugin adapters

**Files:**
- Create: `plugins/ida/goreveal_ida.py`
- Create: `plugins/ida/README.md`
- Create: `plugins/ghidra/GoRevealScript.java`
- Create: `plugins/ghidra/README.md`
- Deferred backlog: `plugins/jeb/` thin integration layer over canonical exports

- [x] Start the `IDA` adapter as a thin, schema-driven action builder over `export ida`.
- [x] Add container-testable validation for the thin `IDA` adapter behavior.
- [x] Add fixture-driven validation for the `IDA` adapter over real export payloads.
- [x] Keep the `IDA` adapter free of recovery logic.
- [x] Add the first `Ghidra` adapter with the same thin-contract rules.
- [ ] Deferred backlog: evaluate a thin `JEB` integration or plugin layer that consumes canonical exports without adding RE logic outside `GoREveal`.

---

## Chunk 9: Sprint 9 - Performance and SIMD Foundation

### Task 9.1: Build benchmark harness and identify hotspots

**Files:**
- Create: `bench/README.md`
- Create: `bench/analyze_bench_test.go`
- Create: `bench/hotspots/pattern_scan_bench_test.go`
- Create: `bench/hotspots/string_scan_bench_test.go`

- [ ] Write corpus-scale benchmarks for end-to-end analysis.
- [ ] Write microbenchmarks for likely hotspots.
- [ ] Capture baseline measurements before optimization.

### Task 9.2: Implement first optimized scalar path and first SIMD candidate

**Files:**
- Create: `core/scan/pattern_scalar.go`
- Create: `core/scan/pattern_scalar_test.go`
- Create: `core/scan/pattern_amd64.go`
- Create: `core/scan/pattern_amd64_test.go`
- Create: `core/scan/pattern_generic.go`
- Create: `core/scan/pattern_bench_test.go`

- [ ] Implement optimized scalar path first.
- [ ] Add deterministic equivalence tests.
- [ ] Add architecture-specific SIMD path only after benchmark evidence.
- [ ] Keep scalar fallback as the canonical behavior.

---

## Chunk 11: Sprint 11 - Native Capability Transfer From `redress` and `gore`

### Task 11.1: Deepen source reconstruction beyond the current narrow overlap

**Files:**
- Create: `engine/projection/source_tree_v2.go`
- Create: `engine/projection/source_tree_v2_test.go`
- Modify: `schema/analysis.go`
- Modify: `cmd/goreveal/internal/source_tree.go`

- [x] Add richer package-to-file grouping in canonical schema.
- [x] Enrich grouped source-package nodes with package-level function metadata.
- [x] Strengthen source-tree projection beyond the current module-local file set.
- [x] Validate the deeper projection against current `redress` fixture behavior without copying its presentation layer.

### Task 11.2: Improve package and user-type metadata quality

**Files:**
- Create: `core/packages/metadata.go`
- Create: `core/packages/metadata_test.go`
- Create: `core/types/user_types.go`
- Create: `core/types/user_types_test.go`

- [x] Improve package metadata quality beyond simple presence checks.
- [x] Separate user-meaningful type evidence from runtime-heavy background noise.
- [x] Use `gore` only as a behavior reference, not as a compatibility target.

---

## Chunk 12: Sprint 12 - Native Capability Transfer From `GoReSym`

### Task 12.1: Add richer runtime metadata recovery

**Files:**
- Create: `core/runtime/moduledata.go`
- Create: `core/runtime/moduledata_test.go`
- Modify: `engine/engine.go`
- Modify: `schema/analysis.go`

- [ ] Recover richer runtime metadata beyond current build info and narrow file/function overlap.
- [ ] Represent that metadata in schema with explicit provenance.
- [ ] Validate against current `GoReSym` evidence surfaces.

### Task 12.2: Move from DWARF-first type recovery to stronger native Go metadata

**Files:**
- Create: `core/types/typelinks.go`
- Create: `core/types/typelinks_test.go`
- Modify: `core/types/types.go`
- Modify: `core/types/types_test.go`

- [ ] Add `typelinks`-driven type recovery where possible.
- [ ] Keep `DWARF` as evidence/fallback, not the only type source.
- [ ] Expand multi-version recovery truth before claiming broad parity.

---

## Chunk 13: Sprint 13 - Selective Deobfuscation Transfer

> Legacy deferred note: this chunk is no longer the active next sprint. The active `Sprint 13` baseline is now workstation handoff contract hardening in `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`.

### Task 13.1: Bounded string refinement transfer from `gostringungarbler`

**Files:**
- Create: `deobfuscation/garble/ungarble.go`
- Create: `deobfuscation/garble/ungarble_test.go`
- Modify: `schema/refined.go`

- [ ] Add bounded garble-aware string refinement that preserves raw/refined separation.
- [ ] Validate only safe overlap surfaces before claiming parity.
- [ ] Do not start this task before the bounded Windows `PE` fixture checkpoint and first code-peeling MVP are complete.

### Task 13.2: Bounded name/symbol refinement transfer from `GoResolver`

**Files:**
- Create: `deobfuscation/refine/symbols.go`
- Create: `deobfuscation/refine/symbols_test.go`
- Create: `docs/architecture/2026-03-19-goreveal-deobfuscation-transfer-notes.md`

- [ ] Transfer only bounded, schema-friendly refinement ideas from `GoResolver`.
- [ ] Do not clone the full CFG-similarity engine into the early product line.
- [ ] Keep orchestration and comparison paths available while native refinement matures.
- [ ] Prefer external orchestration first if garble-constraint solving becomes active work later.

---

## Chunk 10: Sprint 10 - Service/API and Release Baseline

### Task 10.1: Add service mode and machine-facing API

**Files:**
- Create: `api/go.mod`
- Create: `api/proto/analyzer.proto`
- Create: `api/server/server.go`
- Create: `api/server/server_test.go`
- Create: `cmd/goreveal/internal/serve.go`

- [ ] Define protobuf or Connect-compatible analysis API.
- [ ] Implement service mode as a thin wrapper over engine/schema.
- [ ] Add API tests using fixtures and stored analyses.

### Task 10.2: Add release and operations baseline

**Files:**
- Create: `.goreleaser.yaml`
- Create: `.github/workflows/ci.yml`
- Create: `docs/INSTALL.md`
- Create: `docs/OPERATIONS.md`
- Create: `docs/RELEASE.md`

- [ ] Add CI for format, lint, test, fuzz smoke, and benchmark smoke.
- [ ] Add release automation for multi-platform binaries.
- [ ] Document install, corpus update, and release workflows.

---

## Standard Verification Commands

Run these through the Podman dev container as the repo matures:
- `podman exec goreveal-dev go test ./...`
- `podman exec goreveal-dev go test -run TestName ./path/to/package`
- `podman exec goreveal-dev go test -fuzz=Fuzz -run=^$ ./...`
- `podman exec goreveal-dev go test -bench=. -benchmem ./...`
- `podman exec goreveal-dev golangci-lint run`
- `podman exec goreveal-dev go test ./tests/differential/...`

## Planning Notes

- No code copying from baseline projects is allowed.
- Any recovery claim should be backed by corpus evidence.
- Any plugin or API surface must consume canonical schema outputs rather than reimplement recovery logic.
- Any SIMD work must prove identical semantics and benchmark value.
