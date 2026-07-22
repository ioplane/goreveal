# GoREveal Strategic Review

> Status: historical strategic review; superseded for execution by the
> [RT1 product design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
> plus its
> [standalone-first refinement](../superpowers/specs/2026-07-22-goreveal-standalone-release-ida-bootstrap-design.md)
> and [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md)
> Date: 2026-03-31
> Purpose: consolidate the external strategic review into a repo-native planning document and translate it into concrete roadmap, architecture, and documentation updates.

## Summary

This review records the high-level product reading at that checkpoint:
- the project is materially ahead as a platform product
- the main remaining gap is still runtime-semantic accuracy depth
- `Sprint 12` remains the correct primary lane
- `Sprint 7` should remain maintenance-only
- `Sprint 11` should be treated as a completed checkpoint unless stronger runtime truth later forces package/type/source changes
- repo automation and agent workflow contracts are now in a healthy support state and should stay supporting, not become a competing roadmap lane

Current quantified checkpoint:
- platform as a product: `99%`
- accuracy engine: `74%`
- overall roadmap completion: `99%`

The current bounded bridge chain on the canonical ELF fixture is enough as a proof-of-concept.
The next best move is not more blind same-fixture bridge accumulation by default.
The first bounded `Mach-O` function foothold is now landed.
The first real comparison pass against `gore`, `GoReSym`, and `redress` is now also landed.
The first bounded `PE` function foothold is now landed too.
The current weighted next move is no longer the rerun itself, because that rerun is now completed for the current highest-signal slices. The first package-level transfer-workflow polish and the first thin source-visibility response are now also landed. The protected-binary matrix is now widened through `arm64`, and the old `garble v0.15.0` release-gap has been separated from current upstream reality. The first real `garble` rows are now measured through a local upstream checkout on both `linux/amd64` and `linux/arm64`. The latest bounded protected triage slice converted the old Linux architecture split into a portable section-backed `elf_function_foothold = "address_only"` projection, while keeping `moduledata_text_*` absent and leaving named-function recovery untouched. The newest compact analyst-facing increment now also projects whether that foothold is backed by `moduledata_text` or by the fallback ELF `.text` section, plus a bounded foothold span, and that same surface is now locked at export-contract and plugin-consumer boundaries. The workflow/value lane has now restarted in concrete form too: `storage/diff` exposes a compact `transfer_review` queue, a package-first `transfer_review_packages` triage surface, a bounded `transfer_review_focus` first-pass bundle over the already-landed matching and transfer primitives, a compact `transfer_review_plan` action queue, that same plan now carries package-ordered attached review items, an explicit `goreveal diff review sqlite ...` operator path over that same bounded review state, a compact machine-readable `handoff` block with left/right input context and recommended workstation targets, a structured per-target `target_profiles` contract for `ida` and `ghidra`, explicit export-contract IDs, preferred transport hints, artifact roles, workspace phases, host actions, a self-describing artifact bundle, binding-entrypoint hints, expected host-outcome hints, a dedicated `goreveal diff handoff sqlite ...` operator-facing projection for that bridge, and now a dedicated `goreveal diff next sqlite ...` next-step operator projection over the same review state with self-contained `recommended_actions` for review, handoff, and export follow-through, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress` over the remaining plan, a compact `up_next` snapshot for the next package after the current focus, and an `upcoming_packages` horizon with sample pair and strongest-match context. The current weighted next move is therefore cleaner still: keep building the workflow/value lane by default, treat `Sprint 14` review workflow actionability as materially underway, and only spend more on workstation-contract hardening if one last thin lock clearly pays for itself. The newest planning delta is that `Sprint 14` now also has a frozen stop-condition: once `diff next sqlite` supports the local loop `select -> review -> handoff -> export -> move next` without falling back to the raw queue, any further micro-surface must justify itself by removing one named remaining inference step rather than by general helpfulness. The selected post-`Sprint 14` target is now also explicit: a bounded source-evidence confidence contract over already-known `function`, `package`, and `source_tree` truth, because the freshest comparison still leaves `redress` as the strongest practical source/file-oriented baseline while `GoREveal` is no longer file-blind on the measured external targets. That lane is now materially executing through `source_evidence_kind`, a compact `source_evidence_summary`, and per-evidence-class file counts in that summary, and no new named remaining inference step is currently established. The current strategic decision is therefore to freeze `Sprint 15` by default and move the active PM queue into `Sprint 16`.
The new draft research on protecting commercial Go software should be folded into planning as a threat-model and target-binary expansion, not as a replacement for the current active lane.
The new `rehelp` plus remote RE-lab inventory research should also be folded into planning as an interop/orchestration signal, not as a reason to move third-party RE engines into `core`.
See also:
- `docs/plans/2026-03-31-goreveal-external-binary-matrix-evaluation.md`
- `docs/plans/2026-03-31-goreveal-baseline-comparison-plan.md`
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md`
- `docs/plans/2026-04-01-goreveal-next-execution-plan.md`
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`
- `docs/plans/2026-04-01-goreveal-universal-re-workbench-comparison.md`
- `docs/plans/2026-04-01-garble-go126-support-research.md`
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md`

## Immediate Actions

Supporting execution checkpoint already landed:
- repo-local Codex skills and subagents now live in a portable layout through `.agents/skills/` and `.codex/agents/`
- `Taskfile.yml` plus `scripts/dev/podman_runner.py` now make the Podman-first operator flow explicit
- the dev image now also bundles `jq`, `yq`, `procps`, and `unzip`, reducing friction in the canonical container workflow for structured output inspection and debug
- strict script verification now exists through `ruff`, `ty`, `yamllint`, and `shellcheck`
- this is enough repo-ops foundation for the current phase; the roadmap should stay focused on workflow/value and workstation handoff hardening unless a new measured accuracy gap clearly outranks it

### 1. Runtime trust/evidence summary

This bounded `Sprint 12` slice is now landed on `schema.RuntimeMetadata`.

Recommended values:
- `symbol_backed`
- `go_module_fallback`
- `section_heuristic`
- `absent`

This should be:
- one compact field, not a growing set of boolean flags
- exposed through `goreveal inspect runtime`
- validated on both the rich and stripped fixtures
- mirrored into `IDA` / `Ghidra` exports only if it stays thin and schema-driven

Immediate next question after landing:
- that export decision is now taken in favor of thin schema-driven mirroring for `IDA` and `Ghidra`
- keep future runtime export growth bounded so exports do not become a second recovery API

### 2. Public-release licensing blocker

The repository still lacks a finalized license decision.

Recommended release posture:
- `MIT` for the repository baseline

Reasoning:
- matches the clean-room boundary
- keeps thin commercial integrations viable
- keeps future server-mode commercialization open
- avoids unnecessary friction with closed-source host-platform workflows

This review treats explicit repository licensing as a public-release blocker.

Current delta:
- the repository still has no `LICENSE` file
- the README still does not carry an `MIT` badge because the in-repo decision artifact is still missing

## Strategic Actions

### 1. Add a third fixture: Windows PE

The current fixture matrix is too narrow.
Both active fixtures are `ELF`.

This checkpoint is now started:
- `corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe` exists
- the current bounded proof covers format detection, `debug/buildinfo`, and `analyze` without `PE` runtime claims
- snapshot coverage and thin export CLI coverage now also pin that checkpoint as a maintained product surface
- the current checkpoint now also includes a bounded `PE` runtime section heuristic over `.text` / `.rdata`, a raw `.rdata` `pclntab` magic candidate, and one header-looking `.rdata` `pclntab` candidate with raw `magic`, `quantum`, and `pointer_size`, still below any `moduledata` or generic `PE` parser claim
- the next `PE` step should remain evidence-first rather than jumping to a broad parser

Recommended next corpus addition:
- keep the widened protected matrix across `amd64` and `arm64` as the current evidence baseline, and treat the Linux-architecture portability checkpoint for `address_only` footholds as landed

Why:
- validates the cross-format design of `core`
- opens a major Go malware analysis surface
- gives `Sprint 12` a second real fixture instead of another same-shape ELF bridge

First PE slice should stay bounded:
- `PE` format detection
- build info via `debug/buildinfo`
- bounded `.gopclntab` / `.typelink` section evidence through PE sections
- no broad `moduledata` parser claim

### 2. Start Go code peeling after the PE checkpoint

The strongest near-term product differentiator remains:
- `Go code peeling / user-code isolation`

This is now started in bounded form:
- `analysis.peeling` exists as an engine-owned layer over canonical truth
- the first landed slice is function-level classification as `user | stdlib | runtime | third_party`
- that same bounded layer now also includes package-level summaries and per-class function counts
- the current implementation uses only existing `function` and `build_info` truth
- `inspect peeling`, `goreveal peel <binary>`, and thin export mirroring already expose that bounded surface
- the next bounded refinement is now also landed: function-level peeling output carries explicit `classification_evidence`, and the first tiny stdlib/runtime fingerprint-assisted refinement now exists inside `engine/peeling`

Recommended ownership:
- `engine` owns the peeling/classification layer
- `core` stays responsible only for recovery truth

Recommended MVP shape:
- classify functions as `user | stdlib | runtime | third_party`
- use existing function/package/type metadata first
- add a bounded fingerprint layer for stdlib/runtime recognition
- expose it through CLI and export payloads without mutating raw truth

Current delta:
- that first bounded fingerprint layer is now landed
- the first version-tracking-adjacent matching surface is now also landed in bounded form through `diff sqlite matched_functions`
- that bounded surface now covers `exact_name`, `source_location`, `source_file`, and `module_local_normalized_name` reasons
- the next decision is whether to widen fingerprints further or deepen the new bounded matching surface

Recommended next bounded order from the current state:
1. return to workflow/value work by default now that the protected lane is stabilized through runtime, export, and consumer boundaries
   - that first analyst-facing transfer review surface is now landed in compact CLI form through `goreveal diff review sqlite ...`
2. one more bounded protected-specific increment only if a new measured analyst pain point makes it clearly higher value than workflow/value work
3. function-level version tracking work only if the protected lane no longer dominates the measured gap

### 3. Freeze the server stack decision

The repo already leans toward the correct server-mode stack.
That decision should now be made explicit.

Recommended stack:
- `PostgreSQL 18`
- `pgx/v5 + sqlc`
- `goose`
- `ConnectRPC`
- `River`
- `S3-compatible object storage`
- local `SQLite` via `modernc.org/sqlite`
- `koanf v2`
- `cockroachdb/errors`

This review records that the server stack should be documented as an architecture decision rather than left as open brainstorming.

### 4. Expand deferred backlog selectively

Recommended deferred additions:
- thin `Rizin` adapter
- explicit MCP interop with host-platform MCP servers
- function-level version diffing after the protected-binary lane no longer dominates the measured gap
- a protected-commercial-binary corpus lane covering `-s -w`, `-trimpath`, `PIE`, `garble`, and enterprise feature-gate/license-check workflows

Updated reading from the current RE-lab inventory:
- host-platform MCP interop is now more concrete because `ida-pro-mcp` is already present in the measured operator environment
- function-level diffing is now better grounded as an external-reference/orchestration lane because `Diaphora` and `BinExport` are also present on the same host
- `Rizin` remains lower in product priority than `JEB` and `Binary Ninja`, but the technical path is more credible because the same host already carries `rizin` plus related plugins and headless tooling
- later protected/deobfuscation orchestration now has a realistic external tool environment through `frida`, `angr`, `qiling`, `unicorn`, `uftrace`, and `z3`

Recommended `Sprint 13` note:
- deobfuscation should not become the next active lane before the PE fixture and code peeling checkpoint
- if garble-constraint work resumes later, start with external orchestration before evaluating a native Z3 path

## Business and Release Notes

This review also records non-code decisions that affect roadmap shape:
- formal company setup should happen before public release
- sanctions/compliance constraints must be explicit before commercial rollout
- early customer targeting should favor EU security/compliance use cases where Go-native recovery and code peeling are directly valuable

These are planning constraints, not implementation tasks for `core`.

Current doc follow-up:
- commercialization and compliance notes are now tracked separately in `docs/plans/2026-03-31-goreveal-commercialization-and-compliance-notes.md`
- review-vs-repo deltas are tracked in `docs/plans/2026-03-31-goreveal-review-gap-checklist.md`

## What Not To Do

Do not:
- keep extending bounded `.gopclntab` bridges on the same fixture by default
- start SIMD work before profiling evidence exists
- begin a rich TUI or web UI before schema and core truth are stronger
- clone heavy external deobfuscation engines into `core`
- start server mode before the license decision and PE fixture checkpoint
- rewrite package/type heuristics before runtime-semantic truth is materially stronger
- treat metadata-network ideas as active implementation work before code peeling and version tracking foundations exist

## Historical Recommended Execution Order

This order is preserved as decision evidence only. It is not an executable
backlog; the active sequence is RT1-S0, S1, S2A, S2B, S2C, the standalone R1
release gate, then separate S3A headless-bootstrap and S3B plugin decisions.

1. treat `Sprint 14` as frozen by default for the current declared local scope unless one named remaining operator inference step can still be demonstrated
2. treat `Sprint 15` as frozen by default unless one new named source-confidence inference step is explicitly demonstrated
3. historically, move PM ranking into `Sprint 16` protected workflow/orchestration selection
4. treat any further `Sprint 13` workstation-contract work as optional and thin
5. repository license decision
6. one bounded protected-specific DEV increment only if the `Sprint 16` ranking names a first measured analyst pain point clearly enough to justify corpus/comparison work
7. server stack scaffold
8. public-release readiness and support posture
9. broader comparison/evidence automation
10. build correlation and version tracking
11. metadata knowledge network
12. analyst workspace automation and replay
13. comparative knowledge packs and decision support
9. build correlation and version tracking
9. metadata knowledge network

Fresh comparison addendum:
- the latest external rerun now shows real file visibility on measured `ELF`, `PE`, and `Mach-O` targets, so the strongest remaining product gap is no longer raw file absence
- the latest Go-native comparison plus the universal-workbench comparison make the product split clearer:
  - `GoREveal` should keep becoming the Go-native truth, trust, and transfer system of record
  - `IDA`, `Ghidra`, `JEB`, `Rizin`, and similar suites should remain the workspace and analyst-action layer
- this strengthens, rather than weakens, the case for explicit handoff/MCP hardening as the next default move
- the near-term sprint sequence should therefore be explicit too: workstation handoff contract first, workflow actionability second, thin semantic/source confidence third if still justified, and protected workflow/orchestration fourth
- the later sprint horizon should now be explicit as well: bounded server control-plane foundations first, then metadata and remote interop platform work, then public-release readiness and broader evidence/comparison automation as separate follow-on sprints rather than one blended backlog bucket, and only after that the longer-horizon differentiators of build correlation/version tracking and a metadata knowledge network

## Related Documents

- `docs/plans/2026-03-20-goreveal-deferred-continuation.md`
- `docs/plans/2026-03-19-goreveal-scrum-implementation-plan.md`
- `docs/plans/2026-03-20-goreveal-functional-assessment.md`
- `docs/plans/2026-03-19-goreveal-feature-map.md`
- `docs/plans/2026-03-20-goreveal-market-killer-features-brainstorm.md`
- `docs/plans/2026-03-21-goreveal-agent-mcp-and-artifact-transfer-ideas.md`
- `docs/plans/2026-03-31-goreveal-review-gap-checklist.md`
- `docs/plans/2026-03-31-goreveal-commercialization-and-compliance-notes.md`
- `docs/architecture/2026-03-31-goreveal-server-stack-decision.md`
