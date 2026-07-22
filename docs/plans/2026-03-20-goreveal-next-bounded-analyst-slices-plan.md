# GoREveal Next Bounded Analyst Slices Plan

> Status: historical PM/risk planning note; superseded for execution by the
> [RT1 product design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
> and [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md)
> Date: 2026-03-20
> Purpose: choose the next bounded analyst-facing slice in the canonical CLI/runtime navigation surface without drifting into broad parser or heuristic rewrites.

## Why This Note Exists

`Sprint 12` has been productive because recent slices followed one rule:
- expose already-known truth
- do not broaden parser scope
- do not overclaim runtime-semantic support

That rule is still correct.

The current question is narrower:
- what is the next bounded analyst-facing slice now that function/string/export/plugin metadata have been strengthened?

This note evaluates the immediate options with product, risk, and execution criteria.

## Historical Product Read

At that checkpoint, the product was already strong in these analyst-facing areas:
- function navigation metadata
- string address truth
- stripped-runtime fallback visibility
- thin export and adapter behavior for canonical function/type/string metadata

The main remaining analyst-facing gap is not missing raw fields. It is missing explicit visibility into:
- when a source-tree node is backed by real file evidence
- when a source-tree node is a bounded fallback
- how much source-backed confidence a source-tree node should be given by an operator

This matters because the stripped-fixture path is now real product behavior, not only a spike.

## Candidate Directions

### Option A: Source-Tree Evidence Flags

Add explicit source-evidence fields to `SourcePackage` and possibly `SourceTree`, for example:
- `has_file_evidence`
- `fallback_empty_files`

These would be derived only from already-known truth:
- actual `files` length
- current bounded fallback path
- existing `external_packages` separation

**PM value**
- High.
- Directly improves analyst trust in `source-tree`.
- Makes stripped and rich fixture outputs easier to interpret without reading implementation notes.

**Risk**
- Low.
- No new recovery logic.
- No new heuristic scope.
- No runtime claim expansion.

**Speed**
- High.
- Schema + projection + tests + docs only.

### Option B: Runtime Summary/Source Enum

Add a higher-level summary field in `analysis.runtime`, for example:
- `firstmoduledata_source: "symbol" | "go_module_fallback"`
- or a more generic runtime-source classification

**PM value**
- Medium.
- Helpful, but narrower than `source-tree` evidence flags because one useful fallback bit already exists.

**Risk**
- Low.
- Mostly projection and naming.

**Speed**
- Medium to high.

### Option C: Package/Type Origin Classification

Add explicit origin/scope classification across package/type/function surfaces, for example:
- `origin: module | stdlib | runtime | external`

**PM value**
- Potentially high.

**Risk**
- Medium to high.
- This begins to harden heuristics that the current strategy explicitly wants to keep frozen until runtime naming/scope truth is stronger.
- Easy to overfit to current fixtures and convert weak inference into fake authority.

**Speed**
- Medium.

## Historical Recommendation

The recommendation below records the March checkpoint and no longer selects
active work.

Update after the latest checkpoint:
- **Option A is complete.**
- `source-tree` now exposes explicit `has_file_evidence`.
- `inspect packages` now exposes explicit `has_source_evidence`.
- the next bounded analyst-facing move should now shift away from source-tree/package trust flags and into runtime trust summarization.

The then-current recommendation was **Option B: Runtime Summary/Source Enum**.

Why:
- it adds new analyst-facing truth without changing recovery semantics
- it improves the already-landed `inspect runtime` surface instead of widening parser scope
- it makes the current bounded runtime contract easier to read for operators and adapters
- it stayed inside the then-current `Sprint 12` claim boundaries

Do **not** choose Option C yet.

That work should wait until runtime naming/scope truth is stronger than the current fixture-local state. Right now it would create more product ambiguity than value.

## Completed Slice Definition

The completed bounded slice did this:

1. Added explicit file-evidence metadata to `SourcePackage`.
2. Marked whether a `source-tree` package node is backed by real file evidence or by the bounded empty-files fallback.
3. Preserved current `files` behavior exactly as-is.
4. Kept `external_packages` and `packages` split unchanged.
5. Extended package navigation with explicit `has_source_evidence` derived only from existing source-tree correlation.

## Proposed Next Slice Definition

The next bounded slice should do only this:

1. Add a compact runtime trust/evidence summary field to `analysis.runtime`.
2. Expose the same summary through `inspect runtime`.
3. Keep all current raw runtime fields intact.
4. Avoid inventing any new parser semantics or cross-version claims.
5. Optionally mirror the same summary into export payloads if it remains thin and schema-driven.

This should **not**:
- change runtime recovery logic
- change any existing raw runtime fields
- introduce generic support claims
- add package/type origin enums
- reopen broad heuristic work

## PM/Risk Decision

The next lane should stay inside `Sprint 12`, but as a **bounded analyst-trust slice**, not as a parser slice.

That means:
- continue using current runtime/navigation truth
- improve operator readability and trust
- avoid heuristic expansion

This keeps the product on the current high-signal path:
- stronger UX from existing truth
- no claim inflation
- no parser explosion

## Nearest 5 Tasks

1. Add a compact bounded trust/evidence summary field to `schema.RuntimeMetadata`.
2. Expose it through `inspect runtime` for rich and stripped fixtures.
3. Add red-green tests for rich fixture and stripped fixture runtime trust behavior.
4. Decide whether the same summary should be mirrored into `IDA` / `Ghidra` exports.
5. Reassess whether the following slice should stay in runtime trust surfaces or return to capability transfer.

## Exit Criteria

The completed source-tree/package trust slice is complete because:
- rich fixture `source-tree` explicitly shows file-backed package nodes
- stripped fixture `source-tree` explicitly shows fallback-backed package nodes
- `inspect packages` now distinguishes source-backed package metadata from bounded fallback-backed metadata
- no recovery or grouping semantics changed

The next runtime-summary slice is complete when:
- rich fixture `inspect runtime` shows a compact trust/evidence summary
- stripped fixture `inspect runtime` shows the bounded fallback path through the same compact summary
- no raw runtime fields are removed or reinterpreted
- docs describe this as bounded analyst-facing truth, not as broader runtime-semantic support

## Strategic Follow-Up

If this slice lands cleanly, the next reassessment should compare:
- another source-tree trust slice
- a runtime-summary slice
- a temporary return to capability transfer

It should **not** jump straight into package/type origin heuristics unless runtime truth has materially improved first.
