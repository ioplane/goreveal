# GoREveal Deferred Continuation Handoff

> Status: historical continuation checkpoint; superseded on 2026-07-22 by the [RT1 design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
> Date: 2026-03-20
> Purpose: capture the current product/engineering state and the next active tasks so work can resume later without rebuilding context from chat history.

Strategic follow-up:
- `docs/plans/2026-03-31-goreveal-strategic-review.md` now extends this handoff with the next roadmap reshaping decisions
- `docs/plans/2026-03-31-goreveal-external-binary-matrix-evaluation.md` now records the first open-source cross-platform matrix evaluation outside the fixture corpus
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md` now records the first real post-`Mach-O` comparison against `GoReSym`, `redress`, and `gore`
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md` now records measured workstation and interop signals from `rehelp` and the remote RE lab host
- `docs/plans/2026-04-01-goreveal-universal-re-workbench-comparison.md` now records the role-based comparison between `GoREveal` and general-purpose RE workbenches
- `docs/plans/2026-04-01-goreveal-next-execution-plan.md` now records the concrete near-term execution order
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md` recorded the then-active sprint sequence after the `Sprint 12` checkpoint

## Current Checkpoint

The repo is in a clean verified state after the latest bounded `Sprint 12` protected-binary triage slice.

Latest landed slice:
- raw ELF protected-binary runtime evidence now includes `elf_pclntab_header_magic`, `elf_pclntab_header_magic_kind`, `elf_pclntab_header_quantum`, and `elf_pclntab_header_pointer_size`
- that same raw ELF header surface now also includes bounded `elf_pclntab_function_count_hint`, `elf_pclntab_file_count_hint`, and the current Go 1.20+ style table-offset hints from `pcHeader`
- that same bounded protected slice now also includes raw `elf_functab_first_pc_offset_hint`, `elf_functab_last_pc_offset_hint`, `elf_functab_pc_offsets_monotonic`, the first bounded absolute `elf_functab_*_pc_addr_hint` fields, a bounded sampled `elf_functab_pc_addr_sample`, and the first compact projected `elf_function_foothold = "address_only"` surface with `elf_function_foothold_count_hint`, keeping the new foothold address-only and count-only instead of pretending to recover names
- `engine` now also emits the bounded analyst-facing blocker `elf_function_recovery_blocker = "custom_pclntab_magic"` when garbled ELF rows keep runtime posture but collapse to zero functions
- this slice stays evidence-first and explainability-first: it does not claim general garble recovery or move deobfuscation logic into `core`

Recent repo-ops checkpoint already landed:
- repo-local Codex configuration now uses `.agents/skills/`, `.codex/agents/`, and `.codex/config.toml`
- `Taskfile.yml` now mirrors the main `Makefile` operator flows
- `scripts/dev/podman_runner.py` now owns Podman container execution for repo tasks
- the dev image now also includes `jq`, `yq`, `procps`, and `unzip` for canonical JSON/YAML/debug/artifact workflows inside Podman
- script-facing verification now has a strict baseline through `ruff`, `ty`, `yamllint`, and `shellcheck`

Recent bounded analyst-facing slices already landed:
- `analysis.runtime` now exposes compact `trust_summary` values for the current rich and stripped ELF families
- `export ida` and `export ghidra` now also mirror canonical `runtime.trust_summary`, keeping runtime posture visible to thin host-tool consumers without adapter-side inference
- a real Windows `PE` fixture now exists for the first bounded cross-format checkpoint, and current coverage keeps the claim narrow to format detection plus `debug/buildinfo` and `analyze` without runtime semantics
- that same bounded `PE` checkpoint now also has canonical snapshot coverage and thin export CLI coverage, so the repo treats it as a real maintained checkpoint instead of an ad hoc fixture experiment
- the same `PE` checkpoint now also has a bounded runtime section heuristic through `.text` / `.rdata` ranges, a raw `.rdata` `pclntab` magic candidate, and one header-looking `.rdata` `pclntab` candidate with raw `magic`, `quantum`, and `pointer_size`, with `trust_summary: "section_heuristic"` and still no `moduledata` claim
- the first engine-owned code-peeling slice is now landed through `analysis.peeling`, `inspect peeling`, `goreveal peel <binary>`, and thin export mirroring of bounded function classification as `user | stdlib | runtime | third_party`, and that same layer now also carries package-level summaries plus per-class function counts
- that same peeling layer now also carries explicit `classification_evidence`, and the first bounded fingerprint-assisted `runtime` / `stdlib` refinement is now landed for cases where import-path truth is absent but known name/source fingerprints exist
- the first bounded `Mach-O` function foothold is now landed: `Mach-O` analyses preserve `build_info`, recovered `functions`, recovered `packages`, and `peeling`, while still exposing no `runtime`, `types`, `strings`, or source-tree claims on that format
- the first bounded `PE` function foothold is now landed too: `PE` analyses preserve the existing bounded runtime posture and now also expose recovered `functions`, recovered `packages`, and `peeling`, while still exposing no `types`, `strings`, or source-tree claims on that format
- the first thin source-visibility increment is now landed too: when full DWARF-backed projection is unavailable, `source_tree` may fall back to package/file evidence from recovered function line tables and marks that bounded mode through `pathless_file_evidence: true`
- the protected-binary lane now also has a first explicit garble-collapse explanation surface: garbled ELF rows currently preserve `go_module_fallback` runtime posture, expose `unknown` ELF `pclntab` magic, preserve nonzero header-level function/file-count hints, preserve monotonic `functab` PC-offset hints, preserve bounded absolute `PC` address hints within `.text`, preserve a first sampled absolute `PC` foothold, preserve a first analyst-facing `elf_function_foothold = "address_only"` projection with a count hint, preserve the current `moduledata -> pcheader/pclntable` bridges, and carry `elf_function_recovery_blocker = "custom_pclntab_magic"` when the function/package envelope collapses
- `goreveal diff sqlite` now also carries the first bounded version-tracking-adjacent function-matching surface through `matched_functions`, with exact-name, source-location, source-file, and module-local normalized-name matches, `score`, `reason`, and optional peeling class context
- `goreveal diff sqlite` now also carries the first thin transfer preview through `transfer_candidates`, limited to user-classified matches with bounded projected metadata and simple `ready | review` disposition
- `goreveal diff sqlite` now also carries the first deterministic accepted-transfer surface through `accepted_transfers`, limited to `ready` candidates and still fully bounded by the existing matching/evidence contract
- `goreveal diff sqlite` now also carries package-level transfer workflow summaries through `transfer_packages`, aggregating current candidate, ready, review, and accepted counts without adding a new matcher lane
- `goreveal diff review sqlite <database> <left-id> <right-id>` now exists as the first dedicated operator-facing review projection over the compact `transfer_review`, `transfer_review_packages`, and `transfer_review_focus` surfaces
- that same dedicated operator path now also carries a compact machine-readable `handoff` block with left/right input context and recommended workstation targets, keeping the first host-platform handoff thin and review-shaped instead of turning it into a second diff API
- `goreveal diff handoff sqlite <database> <left-id> <right-id>` now also exists as the first dedicated operator-facing handoff projection over that same bounded review bridge
- `storage/diff` now also exposes `transfer_review_plan`, a compact action queue over the existing review-package triage state, and that plan now carries package-ordered attached review items instead of only package headers
- `goreveal diff next sqlite <database> <left-id> <right-id>` now exists as the first dedicated operator-facing next-step projection over that same bounded review bridge, with the current package bundle attached directly, self-contained `recommended_actions` for review, handoff, and export follow-through, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress` over the remaining package queue, a compact `up_next` snapshot for the next package after the current focus, and an `upcoming_packages` horizon with sample pair and strongest-match context
- that handoff projection now also carries structured per-target `target_profiles` for `ida` and `ghidra`, plus explicit export-contract IDs, preferred transport hints, artifact roles, workspace phases, and host actions, so workstation guidance is no longer only a flat list of recommendations
- the active local workflow lane is therefore now shaped more by sequencing than by missing primitives: the next practical decision is how long to keep polishing review-action surfaces before shifting to later-horizon release or platform work
- snapshot harness now auto-discovers fixture directories that carry `expected.analysis.json`, lowering the maintenance cost of the next fixture checkpoint
- `source-tree` package nodes expose `has_file_evidence`
- `inspect runtime` exists as a direct bounded surface
- stripped runtime exposes `firstmoduledata_from_go_module_fallback`
- `inspect functions` exposes `package`, `import_path`, `module_local`, `source_file`, `source_line`, `autogenerated`
- string candidates expose absolute `addr`
- `IDA` / `Ghidra` exports preserve string address and bounded function/type navigation metadata

## Verification Status

Latest successful verification for this checkpoint:
- `make fmt`
- `make test`
- `make lint`
- `task lint-scripts`
- `python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./storage/diff ./cmd/goreveal/internal ./cmd/goreveal`

Result:
- `make test` green
- `make lint` green with `0 issues`
- `task lint-scripts` green through the dev container
- focused transfer-workflow tests green through the dev container

Important environment note:
- `make lint` now calls `/go/bin/golangci-lint` explicitly from the dev container
- `task lint-scripts` and `make lint-scripts` now also run through the dev container instead of assuming host-installed linters
- Podman endpoint discovery for repo automation may come from `PODMAN_BASE_URL`, `CONTAINER_HOST`, or `DOCKER_HOST`, not only rootless `XDG_RUNTIME_DIR` layouts
- direct sandboxed `podman` calls may still need escalation because of Podman storage locking, but the repo itself is lint-clean

## Product/PM Read

Current PM reading:
- `Sprint 12` remains the primary execution lane
- the best low-regret work is now bounded analyst-facing truth from the new `engine`-owned `peeling` layer over already-known runtime/navigation state
- all three current format footholds now exist in maintained form: strong `ELF`, bounded but useful `Mach-O`, and bounded but useful `PE`
- the new commercial-software-protection research should now be treated as a target-binary/threat-model input, not as a reason to derail the current lane
- broad package/type heuristic expansion remains the wrong move
- `Sprint 7` remains maintenance-only, not the main lane
- `GoREveal` should not try to replace `IDA`, `Ghidra`, `JEB`, or `Binary Ninja` as a generic RE suite
- the strongest future product directions are now explicitly identified as:
  - Go code peeling / user-code isolation
  - Go version tracking / build correlation
  - Go metadata knowledge network
- the next stage-gate after the current runtime-trust slice is a bounded Windows `PE` fixture, not another blind same-fixture bridge
- public-release planning still requires an explicit repository license decision
- future corpus and comparison planning should explicitly include protected Go binaries built with `-s -w`, `-trimpath`, `PIE`, and `garble`, plus OSS+EE style binaries where enterprise gating and license checks are part of the analyst workflow
- that protected-binary lane is now started through a purpose-built enterprise sample and the first non-garbled profile matrix, and the canonical `task protected-matrix` path is now working
- `rehelp` and the current RE-lab host now also show a concrete adjacent workstation environment:
  - `ida-pro`, `ghidra`, `jeb`, `rizin`, `diaphora`, `binexport`, and `ida-pro-mcp` are real operator-environment inputs now, not generic backlog names
  - that strengthens the roadmap case for explicit host-platform MCP and workstation handoff notes
  - it also makes later function-level diffing and protected/deobfuscation orchestration more concrete without changing the clean-room boundary

Current percentage checkpoint:
- platform as a product: `99%`
- accuracy/recovery engine: `73%`
- overall roadmap completion: `99%`
- `Sprint 12`: `99%`

Fresh external comparison reading:
- the latest rerun now confirms real file visibility on measured external `ELF`, `PE`, and `Mach-O` targets, not only function/package breadth
- the latest Go-native comparison still shows `GoREveal` ahead on function/package/peeling coverage across the rerun targets
- the latest universal-workbench comparison makes the product split clearer: `GoREveal` should be the Go-native truth/transfer layer, while general-purpose suites remain the workspace and mutation layer
- that pushes the default next move further toward workflow/value and workstation handoff hardening rather than another parser-first increment

## Active Next 5 Tasks

These are the next **active** tasks, not passive guardrails.

1. Keep `Sprint 14` review workflow actionability active only while one named remaining operator inference step still exists.
   Goal: build on the current review queue, package triage, first-pass focus bundle, and the new `transfer_review_plan` / `diff next sqlite` path only where that removes one concrete local review-step guess instead of adding more matcher rules or queue polish by default.

2. Treat any remaining `Sprint 13` workstation-contract work as optional and thin.
   Goal: only add another handoff-contract lock if it clearly reduces downstream ambiguity for `IDA`, `Ghidra`, or future MCP binding.

3. Keep `Sprint 7` in maintenance mode as workflow/value slices grow.
   Goal: avoid letting richer review/handoff surfaces outrun evidence hygiene, regression discipline, or differential checks.

4. Keep the rich/stripped/runtime-export trust summary, peeling `classification_evidence`, and protected runtime surfaces stable as later workflow layers grow.
   Goal: preserve the compact truthful contracts that the newer operator-facing layers already depend on.

5. Treat the portable protected `address_only` foothold as the current evidence baseline, not as the default next implementation lane.
   Result: protected work is now evidence-backed enough that the next default move should be workstation and workflow value, unless a new specific analyst pain point clearly outranks it.
6. Treat the compact analyst-facing protected projection as completed too.
   Result: the current `address_only` foothold now also carries bounded `elf_function_foothold_text_source` and a compact foothold span, so operators can tell whether the foothold is `moduledata_text`-backed or `elf_text_section`-backed without expanding recovery semantics.
7. Treat the thin protected export-contract lock as completed too.
   Result: the same bounded protected runtime surface is now asserted in `schema/export_test.go`, so `IDA` and `Ghidra` payloads inherit it directly from canonical schema instead of relying on adapter-local behavior.
8. Treat the thin protected plugin-consumer lock as completed too.
   Result: `plugins/ida/test_goreveal_ida.py` and `plugins/ghidra/test_goreveal_ghidra.py` now assert that the new protected runtime surface survives into downstream consumer payloads without changing adapter scope.

Current weighted decision after that checkpoint:
- treat the protected lane as the current evidence baseline, not as the default next implementation lane
- return to workflow/value work by default
- treat the bounded CLI/operator projections for the focused review pass, handoff, and next-step planning as landed
- harden explicit host-platform MCP and workstation-interop planning on top of the now-landed `diff handoff sqlite` bridge only if one more contract lock still looks materially useful, while keeping `Sprint 14` actionability as the default active lane
- only reopen the protected lane immediately if one new analyst-facing protected surface is justified by a specific measured pain point

Latest workflow/value checkpoint:
- `storage/diff` now also exposes a compact `transfer_review` queue over the existing bounded transfer contract
- that same surface keeps scope narrow: it only projects pending human-review items and counts, without adding a new matcher lane or mutable review state
- transfer surfaces now also carry explicit `projected_package`, making the first analyst-facing review queue easier to consume than raw function-level matches alone
- `storage/diff` now also exposes `transfer_review_packages`, a package-first triage projection over those same pending review items with review counts, auto-accepted package context, and the strongest pending match per package
- `storage/diff` now also exposes `transfer_review_focus`, a bounded first-pass bundle for the recommended next operator step built over that existing review state
- the focused CLI/operator projection is now landed through `goreveal diff review sqlite <database> <left-id> <right-id>`
- that same focused CLI/operator projection now also carries a compact `handoff` block, and `goreveal diff handoff sqlite <database> <left-id> <right-id>` now exposes that bridge directly
- the first dedicated next-step actionability projection is now landed too through `transfer_review_plan` and `goreveal diff next sqlite <database> <left-id> <right-id>`, and that projection now carries package-ordered attached review items plus self-contained `recommended_actions`, an explicit `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the lane now has a frozen stop-condition: keep building on those review queues, first-pass bundles, and dedicated CLI paths only while one concrete local inference step still remains

Commercial-software follow-up to schedule immediately after that checkpoint:
- add one bounded protected-binary comparison lane focused on stripped/trimpathed/garbled Go binaries and OSS+EE deployment patterns
- keep it out of `core` until the comparison shows a specific missing surface that cannot be covered by current schema, peeling, diff, or future deobfuscation layers

Deferred integration backlog:
- add a thin `JEB` integration/plugin layer later, using the same export-first contract strategy as `IDA` and `Ghidra`
- add a thin `Binary Ninja` integration/plugin layer later, using the same export-first contract strategy
- add a thin `Rizin` integration/plugin layer later, using the same export-first contract strategy and keeping recovery logic outside the adapter
- keep it explicitly below current `Sprint 12` priority

Deferred strategic epics:
- `Epic: Go code peeling / user-code isolation`
- `Epic: Go version tracking / build correlation`
- `Epic: Go metadata knowledge network`
- `Epic: Protected commercial Go software workflows`
- `Epic: Dual runtime modes and gorectl`
- `Epic: Server storage and multi-tenant artifact platform`
- `Epic: MCP surfaces for local and remote agent workflows`
- `Epic: Object-store-backed artifact transfer and gorectl sessions`
- these are product-shaping opportunities, but they stay below the current runtime/type accuracy lane

## Guardrails For Resume

When resuming, do not:
- broaden package/type heuristics
- add a wide `moduledata` parser
- add another blind same-fixture `.gopclntab` bridge by default
- overclaim cross-version runtime support

When resuming, prefer:
- bounded user-facing truth
- schema-first changes
- rich + stripped fixture regression coverage
- protected-binary evidence additions through corpus/comparison first, not speculative anti-obfuscation claims in `core`
- docs updates in the same change
- product moves that make `GoREveal` a stronger Go-native layer on top of external RE suites, not a clone of them

## Files To Read First On Resume

- `docs/plans/2026-03-20-goreveal-progress-assessment.md`
- `docs/plans/2026-03-20-goreveal-functional-assessment.md`
- `docs/plans/2026-03-31-goreveal-strategic-review.md`
- `docs/plans/2026-03-20-goreveal-next-bounded-analyst-slices-plan.md`
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md`
- `docs/plans/2026-03-21-goreveal-runtime-modes-and-storage-ideas.md`
- `docs/plans/2026-03-21-goreveal-agent-mcp-and-artifact-transfer-ideas.md`
- `docs/architecture/2026-03-31-goreveal-server-stack-decision.md`
- default future server transfer path is now documented as `gorectl -> goreveal server` for control plus direct `gorectl <-> S3/Garage` for large artifact bytes, with `PostgreSQL` as metadata authority
- `docs/plans/2026-03-19-goreveal-scrum-implementation-plan.md`
- `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- `AGENTS.md`

## Resume Recommendation

If continuing later, start with:
- keeping the new `goreveal diff review sqlite ...` and `goreveal diff next sqlite ...` operator paths stable over the current `transfer_review`, `transfer_review_packages`, `transfer_review_focus`, and `transfer_review_plan` surfaces
- then harden explicit host-platform MCP and workstation handoff planning using the measured `rehelp` / RE-lab inventory as the real operator baseline and the now-landed `diff handoff sqlite` bridge as the review-path handoff
- then keep the fresh external comparison and universal-workbench comparison as the planning baseline so future roadmap decisions are judged against real workstation value, not only fixture-local recovery growth
- then follow `docs/plans/2026-04-01-goreveal-next-execution-plan.md` as the concrete near-term ordering artifact
- this handoff originally pointed to `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`; use RT1 instead
- then keep the new `classification_evidence`, `pathless_file_evidence`, and transfer-review surfaces stable while deciding whether workflow/value or another bounded semantic step has the highest leverage
- then any required rich/stripped/runtime-export contract stabilization
- then treat post-`Sprint 18` work as its own later horizon: public-release/licensing hardening first, broader evidence/comparison automation second

That is the next best bounded continuation from the current checkpoint.

After that next slice lands, the first PM re-evaluation question should be:
- keep the bounded `PE` fixture checkpoint stable
- then keep shaping the first future differentiator epic as `Go code peeling / user-code isolation` plus build-to-build transfer, not as parser expansion
