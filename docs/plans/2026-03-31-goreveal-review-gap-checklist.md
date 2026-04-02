# GoREveal Review Gap Checklist

> Status: review delta / action matrix
> Date: 2026-03-31
> Purpose: compare the external strategic review in `docs/tmp/draft/latest-comments.md` against the current repository state and distinguish what is already landed, what still needs code, what still needs documentation, and what still needs an explicit decision.

## Why This Note Exists

The external review was accurate when written, but the repository moved quickly afterward.
Without a delta note, the review now mixes:
- recommendations that are already landed
- recommendations that are still valid but not implemented
- recommendations that are valid but only weakly represented in repo docs

This file is the normalization layer.

## Already Landed And Now Stale In The Review

These review items are no longer pending work:
- compact `runtime trust/evidence summary` is already landed in canonical runtime metadata
- `inspect runtime` already exposes that bounded runtime posture directly
- rich and stripped fixture tests already cover the runtime trust surface
- thin `IDA` / `Ghidra` exports already mirror canonical `runtime.trust_summary`
- the third fixture checkpoint is already real through `go-pe-buildinfo-windows-amd64`
- the current `PE` checkpoint already includes build info, snapshot coverage, thin export coverage, and a bounded runtime section/header heuristic
- the first `code peeling` MVP is already landed through `analysis.peeling`, `inspect peeling`, `goreveal peel <binary>`, and thin export mirroring
- the server stack decision is already frozen in `docs/architecture/2026-03-31-goreveal-server-stack-decision.md`
- `modernize` and `sloglint` are already enabled in `.golangci.yml`

## Missing In Code Or Product

These remain real product gaps:
- repository `LICENSE` file is still absent
- README still has no `MIT` release badge because the license decision is not yet finalized in-repo
- code peeling still lacks the next bounded fingerprint-assisted `stdlib` / `runtime` refinement
- function-level version diffing with similarity scoring is still not implemented
- thin `Rizin` adapter/export flow is still not implemented
- server mode, `gorectl`, and object-store-backed transfer are still planning-only
- the later Sprint 13 deobfuscation path remains unimplemented

## Missing Or Underrepresented In Documentation

These are still doc gaps even where the strategic direction is already known:
- no dedicated repo-native note yet records the commercialization/compliance follow-up from the review
- license policy is still underrepresented because the repo has no final in-tree `LICENSE` decision artifact
- MCP interop with host-platform MCP servers is described only at a high level; the current docs should be clearer about the two-step `GoREveal MCP -> host-platform MCP` workflow
- the next code-peeling refinement is mentioned, but not yet collected as a tighter bounded design target
- version-diffing is present as a backlog idea, but not yet carried as a crisp implementation note in one current review delta document
- `Rizin` appears in backlog notes, but the repo benefits from a clearer statement that it stays below `JEB` and `Binary Ninja` because market priority is lower even if the technical adapter may be simpler

## Needs Decision

These items are blocked more by product/owner decision than by implementation mechanics:
- final repository license choice, currently recommended as `MIT`
- whether commercialization/compliance notes should stay as private planning notes or become part of a more formal public-facing release checklist

## Recommended Next Documentation Moves

1. Keep this file as the delta reference whenever the external review is cited.
2. Carry the commercialization/compliance part of the review into a dedicated repo-native note.
3. Make the host-platform MCP interop workflow explicit in the MCP planning doc.
4. Keep all roadmap docs treating fingerprint-assisted `engine/peeling` refinement as the next low-regret code-peeling step.
5. Do not add an `MIT` badge to the README until the repository actually contains the final `LICENSE` file.

## Related Documents

- `docs/tmp/draft/latest-comments.md`
- `docs/plans/2026-03-31-goreveal-strategic-review.md`
- `docs/plans/2026-03-20-goreveal-deferred-continuation.md`
- `docs/plans/2026-03-20-goreveal-functional-assessment.md`
- `docs/plans/2026-03-21-goreveal-agent-mcp-and-artifact-transfer-ideas.md`
