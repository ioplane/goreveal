# GoREveal Next Execution Plan

> Status: superseded on 2026-07-22 by the [RT1 design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md) and [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md)
> Date: 2026-04-01
> Purpose: turn the current roadmap reading into a concrete near-term execution order after the fresh external comparison, protected-binary stabilization, and workstation baseline research.

## Why This Plan Exists

The current repo no longer needs another default parser lane.

The strongest fresh signals are:
- `GoREveal` is already strong on real `ELF`, `PE`, and `Mach-O` binaries
- the protected lane now has a portable bounded `address_only` foothold on `linux/amd64` and `linux/arm64`
- workflow/value work is no longer hypothetical: `diff review sqlite` and `diff handoff sqlite` are landed
- the operator environment now clearly includes host platforms and adjacent tooling such as `ida-pro-mcp`, `rizin`, `diaphora`, `binexport`, `frida`, `angr`, and `z3`

That means the next best work is to harden the bridge between Go-native truth and the real analyst workstation.

## Recommended Execution Order

### Phase 1: Workstation Handoff Hardening

Primary goal:
- make the existing review/handoff surfaces easier to consume by operators and future MCP bridges

Primary ownership:
- `cmd/goreveal`
- `storage/diff`
- docs

Expected bounded outputs:
- clearer handoff JSON contract over the existing review state
- structured per-target workstation guidance rather than only flat recommendation lists
- explicit export-contract IDs and preferred transport hints per target
- explicit artifact roles, workspace phases, and host action lists per target
- explicit handoff-contract ID, self-describing artifact bundle, per-target binding semantics, and expected host-outcome semantics so future MCP/workstation adapters do not need to infer entrypoints, required payloads, or success conditions
- explicit host-target guidance for `IDA` / `Ghidra` / MCP-oriented flows
- one operator-facing export/review path that can be demonstrated end-to-end without changing recovery truth

Current reading:
- the basic shape is already landed through `diff handoff sqlite`, `target_profiles`, explicit contract/transport hints, artifact roles, workspace phases, and host actions
- the remaining work in this phase is now contract hardening and target-binding polish, not inventing the handoff model
- the handoff artifact is now already close to self-describing for future adapters: it carries its own contract ID, artifact bundle, required-artifact hints, binding semantics, and expected host-outcome semantics

Evidence gates:
- CLI tests
- schema/export contract checks if payload shape changes
- docs updates in `README`, roadmap notes, and MCP planning note

Do not do in this phase:
- new matcher rules by default
- new protected-specific semantic claims
- host-workspace mutation inside `GoREveal`

Immediate PM+DEV stack:
- `PM-A1` decide whether `Sprint 13` still needs one final lock or should be stopped
  Indicator: one explicit yes/no checkpoint tied to remaining adapter ambiguity
- `DEV-A1` if ambiguity remains, remove only that ambiguity
  Indicator: exactly one bindable contract gap closed without growing scope

Active queue:

| Task | Outcome | Value | Risk | Evidence |
| --- | --- | ---: | ---: | ---: |
| `PM-A1` | decide stop vs one-final-lock for `Sprint 13` | `5` | `1` | `5` |
| `DEV-A1` | close exactly one remaining contract ambiguity if it still exists | `3` | `2` | `4` |

### Phase 2: Host-Platform MCP / Interop Contract

Primary goal:
- make the workstation handoff explicit enough that a future MCP bridge is a transport problem, not a product-definition problem

Primary ownership:
- docs
- `cmd/goreveal` only if a tiny projection is needed
- plugins/export tests if a thin payload lock is justified

Expected bounded outputs:
- explicit `GoREveal -> host-platform MCP` handoff steps
- recommended export families per host target
- one stable operator artifact shape for workstation import/review

Evidence gates:
- planning docs
- README
- plugin/export tests only if a real downstream contract changes

Do not do in this phase:
- implement a broad MCP server
- move workspace semantics into `core`
- add server-mode scope by accident

Immediate PM+DEV stack:
- `PM-A2` keep MCP work phrased as transport and handoff, not product replacement
  Indicator: every new interop task maps back to an existing export or handoff artifact

Active queue:

| Task | Outcome | Value | Risk | Evidence |
| --- | --- | ---: | ---: | ---: |
| `PM-A2` | keep interop work transport-shaped and artifact-backed | `4` | `2` | `5` |

### Phase 3: Workflow/Value Polish

Primary goal:
- make the review workflow easier to act on before reopening semantic or protected-specific recovery work

Primary ownership:
- `storage/diff`
- `cmd/goreveal`

Expected bounded outputs:
- one more operator-facing review/action projection
- package-first or target-first prioritization improvements over existing review queues
- zero new recovery heuristics unless a measured blocker forces them

Current reading:
- this phase is now started through `transfer_review_plan` and `goreveal diff next sqlite ...`
- that same next-step projection now also carries package-ordered attached review items plus self-contained `recommended_actions`, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the active lane is no longer just prioritization but a real first-pass package bundle with operator-ready follow-through
- the next step inside this phase should build on that explicit next-step projection instead of reopening matcher or parser work

Evidence gates:
- `storage/diff` tests
- CLI tests
- fresh doc sync

Immediate PM+DEV stack:
- `PM-A4` hold the frozen operator stop-condition for `Sprint 14`
  Indicator: do not reopen the lane unless a measured local workflow gap actually violates the frozen rule
- `DEV-A3` keep adding only one projection at a time over existing state
  Indicator: every new `diff next sqlite` field reduces one concrete operator inference step

Frozen stop-condition:
- treat `Sprint 14` as complete enough once `diff next sqlite` supports the local loop `select -> review -> handoff -> export -> move next` without making the operator consult the larger raw queue or infer the obvious next package, remaining local progress, or completion cue by hand
- after that point, only continue `DEV-A3` if a new projection removes one named remaining inference step that the current fields still leave unresolved
- if no such step exists from measured local use, freeze the lane and advance the roadmap instead of inventing more queue surfaces

Active queue:

| Task | Outcome | Value | Risk | Evidence |
| --- | --- | ---: | ---: | ---: |
| `PM-A4` | hold the frozen `Sprint 14` stop-condition unless a measured local gap reopens it | `4` | `1` | `5` |
| `PM-A8` | landed: the first `Sprint 16` protected target matrix now ranks `garble`-class workflows first by analyst ambiguity, commercial relevance, and measured gap size | `5` | `2` | `5` |
| `PM-A9` | landed: `Sprint 16` now has an explicit stop-condition that keeps the lane workflow/orchestration-first and prevents broad deobfuscation drift | `5` | `1` | `5` |
| `DEV-A5` | landed: the first bounded rerun over neighboring first-ranked `garble`-class builds now shows no review-ready anchor set on either measured Linux architecture despite the existing `address_only` foothold | `4` | `2` | `5` |
| `PM-A11` | next: choose the first bounded post-rerun response now that the neighboring-build workflow question is a measured negative result | `5` | `2` | `5` |

## Deferred Until Re-Prioritized

Keep deferred by default:
- broad parser expansion
- broad deobfuscation lane
- another protected-specific increment without a measured analyst pain point
- service/API expansion
- UI-first work

## Ordered Later Horizon

After the current near-term sequence, the preferred later order is:

1. `Sprint 17` server control-plane foundations
2. `Sprint 18` metadata and remote interop platform
3. `Sprint 19` public release readiness and licensing
4. `Sprint 20` evidence expansion and comparative automation
5. `Sprint 21` build correlation and version tracking
6. `Sprint 22` metadata knowledge network

Why this order:
- `PostgreSQL 18`, `pgx/v5 + sqlc`, `goose`, `ConnectRPC`, `River`, and object storage are already chosen, but still belong to a later product phase
- `gorectl`, remote MCP, and artifact-session work only become low-regret once the local review/handoff flow is mature enough to deserve remote orchestration
- public-release hardening should happen only after the local product story and later-horizon architecture are coherent enough to describe honestly
- broader comparison automation should follow a clearer release/support posture, so evidence breadth keeps serving product truth instead of becoming a parallel vanity lane
- build correlation should follow those steps because it is the strongest next workflow moat once the current operator loop and support/evidence posture are stable
- metadata-network work should remain last because it compounds value from version tracking and remote metadata contracts rather than replacing them
- this keeps the product local-first and Go-native instead of prematurely shifting attention into server scope

## Reopen Conditions

Reopen protected-specific work only if:
- the current handoff/review workflow hits a real analyst limit on garbled samples
- a new bounded protected surface can be added without claiming named recovery

Reopen parser/runtime work only if:
- the workstation/handoff lane stops providing clear product value
- or a fresh comparison reveals a concrete cross-format regression

## Current Recommendation

If continuing immediately, do this:
1. treat the dedicated handoff surface as the canonical workstation artifact now
2. treat `Sprint 14` as frozen by default for the current scope because the stop-condition is now plausibly satisfied by the current `focus + recommended_actions + review_checklist + review_snapshot + review_progress + up_next + upcoming_packages` bundle
3. treat `Sprint 15` as frozen by default for the current scope: the chosen bounded source-evidence confidence contract over already-known `function`, `package`, and `source_tree` truth now has three thin landed slices through `source_evidence_kind`, a compact `source_evidence_summary`, and per-evidence-class file counts in that summary
4. treat the active PM queue inside `Sprint 16` as landed for the first checkpoint: `garble`-class workflows are now ranked first and the stop-condition is explicit
5. only reopen `Sprint 14` if a new measured local inference step appears
6. only add one final thin `Sprint 13` contract lock if a downstream binding ambiguity is still real
7. treat the first `Sprint 16` DEV rerun as landed too: neighboring first-ranked protected builds currently do not produce any review-ready transfer/review/handoff foothold on either measured Linux architecture, so the next move is no longer “measure whether”; it is “choose the smallest truthful response to that measured absence”
8. keep `Sprint 17` and `Sprint 18` deferred until those local product questions are answered by evidence, not by architecture preference
9. once that local sequence settles, prefer `Sprint 19` before `Sprint 20` so release posture hardens before comparison automation widens again
10. after that, prefer `Sprint 21` before `Sprint 22` so build-correlation value is proven before metadata-network work expands
11. only after those two should the roadmap widen into `Sprint 23` analyst workspace automation/replay and then `Sprint 24` comparative knowledge packs and decision support

PM note:
- `PM-A3` is now satisfied by an explicit stop-condition: if no new measured local inference step remains after the current `focus + recommended_actions + review_checklist + review_snapshot + review_progress + up_next + upcoming_packages` bundle, freeze local workflow polish and advance the roadmap instead of inventing extra queue surfaces
- `PM-A5` and `PM-A6` are now satisfied too: the selected `Sprint 15` target gap is stronger source-evidence confidence semantics over already-known file/package truth, because fresh comparison and the latest intermediate rerun still leave `redress` as the strongest source/file-oriented practical baseline while `GoREveal` is no longer file-blind
- the first `DEV-A4` slice is now landed through bounded `source_evidence_kind = dwarf_paths | line_table_files | package_fallback` projection on `source_tree` and enriched package surfaces, the next compact slice is landed through `source_evidence_summary` on `source_tree`, and the next refinement is landed too through per-evidence-class file counts in that summary
- `PM-A8` and `PM-A9` are now landed too: the first protected workflow matrix ranks `garble`-class workflows first, and `Sprint 16` now has an explicit stop-condition that keeps the lane workflow/orchestration-first
- `PM-A10` is now effectively landed too: the first protected workflow pain point is explicitly “no review-ready anchor set on current garbled rows for peel/review/handoff/next despite the existing address-only foothold”
- `DEV-A5` is now landed too through `.tmp/fresh-eval/protected-workflow-summary.json`: neighboring `garble` and `garble + literals + tiny` builds on both `linux/amd64` and `linux/arm64` currently yield zero matched functions, zero transfer/review packages, no review focus, no handoff, no review plan, and no next-step projection
- current weighted reading: no named remaining `Sprint 15` inference step is currently established, `Sprint 14` remains frozen, and `Sprint 16` should stay workflow/orchestration-first by choosing one bounded response to the now-measured absence of any review-ready anchor on neighboring garbled builds before any new analyst-facing or deobfuscation work is proposed
