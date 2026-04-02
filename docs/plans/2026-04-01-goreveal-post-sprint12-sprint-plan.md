# GoREveal Post-Sprint12 Sprint Plan

> Status: active sprint plan
> Date: 2026-04-01
> Purpose: define the near-term sprint set after the current `Sprint 12` checkpoint using the fresh comparisons, strategic review, protected-binary findings, and workstation baseline.

## Why This Sprint Reset Exists

The repo no longer needs a speculative post-`Sprint 12` roadmap.

The current measured facts are:
- `GoREveal` is already strong on real `ELF`, `PE`, and `Mach-O` binaries
- the protected lane is stabilized enough that it is no longer the default next lane
- workflow/value work is real, not aspirational
- the analyst environment already includes meaningful host platforms and sidecars

That means the next sprint set should not be parser-first, deobfuscation-first, or GUI-first.

It should be:
- workstation-aware
- transfer/workflow-driven
- evidence-backed
- still bounded around the current product identity: Go-native truth, trust, and transfer

## Sprint Set Overview

### Sprint 13: Workstation Handoff Contract

Primary goal:
- finish the thin workstation-facing handoff contract over the already-landed review surfaces

Why this sprint is first:
- `diff review sqlite` and `diff handoff sqlite` are already landed
- `target_profiles` now exist and already carry target, transport, contract, artifact-role, workspace-phase, host-action, binding-entrypoint hints, required-artifact hints, and expected host-outcome hints
- the highest-value remaining gap is now making that contract easier to bind to real host platforms and future MCP bridges

Primary ownership:
- `cmd/goreveal`
- `storage/diff`
- docs

Success criteria:
- one stable handoff artifact shape for operator and future MCP use
- explicit host-target semantics for `IDA` and `Ghidra`
- no new recovery claims

Evidence:
- CLI tests
- `storage/diff` tests if payload composition changes
- doc sync across README, progress, strategic review, deferred continuation, and MCP note

PM+DEV task frame:
- `PM-13.1` define stop-condition for workstation handoff hardening
  Indicator: one explicit written rule for when `Sprint 13` is considered done instead of “maybe one more thin lock”
- `DEV-13.1` keep `diff handoff sqlite` contract stable and bindable
  Indicator: no ambiguity in target, transport, artifact bundle, or host-outcome fields for `IDA` and `Ghidra`
- `DEV-13.2` lock any remaining downstream contract edges only if they remove real adapter guesswork
  Indicator: each added field removes one concrete operator or adapter inference step

Current backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-13.1` | PM | freeze an explicit `Sprint 13` stop-condition | `5` | `1` | `5` |
| `DEV-13.1` | DEV | keep `diff handoff sqlite` stable and bindable | `4` | `2` | `5` |
| `DEV-13.2` | DEV | add one last contract lock only if it removes real adapter guesswork | `3` | `2` | `4` |

### Sprint 14: Review Workflow Actionability

Primary goal:
- turn the current review queue into a tighter operator workflow without inventing mutable server state or new matcher heuristics

Why it follows Sprint 13:
- the handoff artifact must stabilize before more workflow projections pile on top
- after that, the next best gain is making the review path easier to act on, not widening recovery semantics

Primary ownership:
- `storage/diff`
- `cmd/goreveal`

Likely bounded outputs:
- one more review/action projection over the current queue/focus state
- clearer package-first or target-first prioritization
- possibly a thin machine-readable action bundle for a future UI/server layer

Evidence:
- `storage/diff` tests
- CLI tests
- docs

Current reading:
- started in bounded form through `transfer_review_plan` and `goreveal diff next sqlite ...`
- the lane now has a dedicated next-step operator projection instead of only queue, package, and focus views
- that next-step projection now carries package-ordered attached review items plus self-contained `recommended_actions`, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so the lane has crossed from prioritization-only into a real first-pass package bundle with explicit follow-on operator steps, queue-position clarity, a clearer “done for this bundle” operator cue, and a smaller machine-friendly summary for future UI/MCP bindings

PM+DEV task frame:
- `PM-14.1` define the minimal local operator loop for a review bundle
  Indicator: one canonical loop described as `select -> review -> handoff -> export -> move next`
- `DEV-14.1` keep enriching `diff next sqlite` only through projection over existing state
  Indicator: no new matcher rules, no new mutable review state, no storage-semantic drift
- `DEV-14.2` make the current package bundle self-sufficient for first-pass review
  Indicator: operator no longer needs the larger raw queue to know what to do, what is next, and what counts as done
- `PM-14.3` hold the frozen stop-condition for leaving `Sprint 14`
  Indicator: do not reopen the lane unless a measured local workflow gap actually violates the frozen rule

Frozen stop-condition:
- treat `Sprint 14` as complete enough to leave when `goreveal diff next sqlite ...` already lets a local operator run the canonical loop `select -> review -> handoff -> export -> move next` without falling back to the raw queue or reconstructing missing next-step context by hand
- that means the active surface must already expose:
  - the current package bundle and its review items
  - what to do now
  - what counts as done for the current bundle
  - what package comes next
  - enough remaining-progress context to continue locally
- after that point, no new `Sprint 14` field should be added unless it removes one concrete remaining operator inference step that can be named before implementation
- if no such step can be named from measured local use, freeze `Sprint 14` and advance the roadmap instead of adding more queue-shaped polish

Current backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-14.3` | PM | keep the frozen `Sprint 14` stop-condition intact unless a measured local gap reopens it | `4` | `1` | `5` |
| `PM-14.4` | PM | declare whether the frozen stop-condition is already met by the current operator bundle | `5` | `1` | `5` |
| `DEV-14.2` | DEV | make the current bundle self-sufficient for first-pass review only if `PM-14.4` finds one unresolved inference step | `5` | `2` | `5` |
| `DEV-14.1` | DEV | keep adding only projection-only review surfaces and only after naming the exact removed inference step | `4` | `2` | `5` |
| `PM-14.1` | PM | lock the canonical local operator loop wording | `4` | `1` | `5` |

### Sprint 15: Thin Semantic and Source Confidence

Primary goal:
- add one more bounded analyst-facing confidence layer only if it clearly improves the product after the handoff/review path stabilizes

Why this is not earlier:
- fresh comparison says the practical gap is no longer “format foothold exists or not”
- it is stronger semantic/source confidence and operator value

Primary ownership:
- `engine`
- `schema`
- `cmd/goreveal` only if a user-facing projection is justified

Allowed scope:
- confidence/visibility increments over already-known truth
- thin source or semantic signals
- no broad parser rewrite

Disallowed scope:
- broad `moduledata` parser work
- wide type-heuristic rewrite
- decompiler-first work

Evidence:
- fixture/snapshot coverage
- differential or external-binary evidence if semantics move
- doc sync

PM+DEV task frame:
- `PM-15.1` nominate one semantic/source-confidence gap that beats workflow work on product value
  Indicator: one ranked target gap with a named analyst pain point and comparison evidence
- `DEV-15.1` ship only one bounded confidence increment at a time
  Indicator: new field or surface is derived from already-known truth and covered by snapshots or differential evidence

Selected target gap:
- `Sprint 15` now targets a bounded source-evidence confidence contract over already-known `function`, `package`, and `source_tree` truth
- named analyst pain point: on the fresh external comparison and the latest intermediate rerun, `GoREveal` is no longer file-blind, but operators still need a cleaner confidence contract to distinguish full DWARF-backed source evidence from line-table-only pathless file evidence or package-only fallback metadata
- why this wins now:
  - `redress` remains the strongest compared source/file-oriented baseline in practical analyst readability
  - the gap is now trust semantics over current file/package surfaces, not raw absence of files
  - this can be shipped as a bounded confidence/explanation layer without reopening parser work or broadening recovery claims
- bounded success shape:
  - one thin confidence vocabulary for source/file-backedness
  - one compact tree-level summary so operators can read the package evidence landscape without reconstructing it from every package node
  - one compact file-density cue inside that summary so operators can read the relative weight of each evidence class without hand-summing package file lists
  - derived only from already-known evidence classes
  - user-visible through canonical schema and operator surfaces only if that projection stays thin
- do not expand this target into:
  - broad source reconstruction
  - new package heuristics
  - typelink or runtime parser work

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-15.2` | PM | keep the selected source-evidence confidence target ranked above protected return unless new evidence changes the order | `4` | `1` | `5` |
| `DEV-15.1` | DEV | landed: bounded source-evidence vocabulary for `dwarf_paths`, `line_table_files`, and `package_fallback` over existing `source_tree` and package truth | `5` | `2` | `5` |
| `DEV-15.2` | DEV | landed: compact `source_evidence_summary` on `source_tree` so the package evidence landscape is explicit without scanning the full package list | `4` | `2` | `5` |
| `DEV-15.3` | DEV | landed: per-evidence-class file counts inside `source_evidence_summary`, so file-density shape is explicit without hand-summing package file lists | `4` | `2` | `5` |
| `DEV-15.4` | DEV | held: reopen only if one real remaining inference step is explicitly demonstrated | `4` | `2` | `5` |
| `PM-15.3` | PM | landed: freeze `Sprint 15` by default because no single named remaining inference step is currently established after the current three slices | `5` | `1` | `5` |

### Sprint 16: Protected Commercial Go Workflows

Primary goal:
- resume the protected lane only as an evidence-first workflow/orchestration epic

Why this is later:
- the protected lane has a truthful `address_only` footing already
- the next missing product value is now how to work with those binaries, not only how to explain why recovery is bounded

Primary ownership:
- corpus
- comparison scripts
- docs
- `deobfuscation` only if a measured external-orchestration gap justifies it

Allowed scope:
- protected corpus expansion
- comparison and orchestration
- bounded external-sidecar workflows

Disallowed scope:
- broad native deobfuscation engine by default
- moving third-party solver logic into `core`

Evidence:
- protected matrix
- comparison notes
- doc sync

PM+DEV task frame:
- `PM-16.1` define which protected workflows matter most commercially
  Indicator: one ordered target matrix covering `-s -w`, `trimpath`, `PIE`, `garble`, and enterprise-gated patterns
- `DEV-16.1` keep protected work corpus/comparison-first
  Indicator: every new implementation move corresponds to one measured protected gap, not a speculative deobfuscation idea

Current ordered target matrix:

| Rank | Target workflow | Why it ranks here | Current measured gap | Next allowed move |
| ---: | --- | --- | --- | --- |
| `1` | `garble` / `garble + literals + tiny` on the purpose-built enterprise-gated sample | highest analyst ambiguity and strongest commercial relevance among already-measured protected classes | truthful `address_only` foothold exists on `linux/amd64` and `linux/arm64`, but named functions/packages/files collapse to `0` under `custom_pclntab_magic`, so the current review/transfer/handoff loop has no review-ready anchor set | one bounded corpus/comparison rerun that asks whether neighboring garbled builds still yield any stable reviewable foothold over current workflow surfaces |
| `2` | enterprise-gated non-garbled profiles across `stripped`, `trimpath`, and `PIE` | closest current proxy for real commercial review workflows and license/feature-gate triage | current user-code and package visibility are already strong, but orchestration and transfer use on protected neighboring builds is not yet measured deeply | one bounded workflow/comparison pass around transfer/review usefulness, not parser expansion |
| `3` | larger OSS protected targets with reproducible builds | strongest route to external validity beyond the purpose-built sample | corpus breadth is still narrow; repeated analyst pain points on larger protected binaries are not yet recorded | expand corpus only after rank `1` or `2` proves which question matters first |
| `4` | non-garbled `arm64` widening beyond current sample families | already broadly stable and lower uncertainty than the rows above | no primary analyst pain point currently outranks higher-ranked targets here | hold unless a fresh comparison reopens architecture-specific drift |

Sprint 16 stop-condition:
- keep the lane PM-active until one first protected workflow target is explicitly ranked and one first measured gap is named
- do not open `DEV-16.*` implementation work until that first target and gap are explicit
- once the first protected rerun or analyst-facing surface is landed, freeze `Sprint 16` again unless a second distinct protected workflow gap is clearly stronger than `Sprint 17`

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-16.1` | PM | landed: rank protected workflow targets commercially and analytically, with `garble`-class workflows ranked first | `4` | `2` | `5` |
| `PM-16.2` | PM | landed: define the first protected workflow matrix by analyst frequency, commercial relevance, and current evidence gap size | `5` | `2` | `5` |
| `PM-16.3` | PM | landed: define the `Sprint 16` stop-condition so the lane stays workflow/orchestration-first and does not reopen broad deobfuscation by drift | `5` | `1` | `5` |
| `PM-16.4` | PM | landed: name the first protected workflow pain point explicitly as “no review-ready anchor set on current garbled rows for peel/review/handoff/next” | `5` | `1` | `5` |
| `DEV-16.1` | DEV | held: open one bounded protected corpus/comparison increment only after `PM-16.1` and `PM-16.2` name the first target gap clearly enough | `4` | `3` | `5` |
| `DEV-16.2` | DEV | landed: neighboring `garble`-class builds currently yield no stable review-ready transfer/review/handoff foothold over current workflow surfaces on either measured Linux architecture | `4` | `3` | `5` |
| `PM-16.5` | PM | next: choose the first bounded response now that the neighboring-build workflow question is a measured negative result | `5` | `2` | `5` |

### Sprint 17: Server Control Plane Foundations

Primary goal:
- turn the already-chosen server stack into a bounded product-foundation sprint without displacing the current local-first product

Why this is later:
- current product value is still local CLI, transfer workflow, and workstation handoff
- the server stack is chosen architecturally, but not yet justified as an active implementation lane
- this sprint should only start after the local transfer/workflow story is clearly mature enough to benefit from remote orchestration

Primary ownership:
- future `api`
- storage/control-plane docs
- future `gorectl` planning

Allowed scope:
- `PostgreSQL 18`, `pgx/v5 + sqlc`, `goose`, `ConnectRPC`, `River`, and object-storage scaffolding
- bounded metadata/control-plane models
- local-vs-server contract notes and first repository scaffolding

Disallowed scope:
- moving recovery into the server
- UI-first work
- broad multi-tenant/auth work by default

Evidence:
- architecture and planning docs
- bounded scaffolding only if sequencing later justifies it
- explicit local/server contract notes

PM+DEV task frame:
- `PM-17.1` define server scope as control-plane only
  Indicator: one explicit contract that recovery remains local/worker-bound and server work is metadata/orchestration-first
- `DEV-17.1` scaffold only the smallest viable control-plane foundation
  Indicator: stack choice, repo scaffolding, and metadata contracts exist without dragging active product work into server mode

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-17.1` | PM | lock control-plane-only scope | `4` | `2` | `5` |
| `DEV-17.1` | DEV | scaffold minimal control-plane foundation | `3` | `4` | `4` |

### Sprint 18: Metadata and Remote Interop Platform

Primary goal:
- turn local workflow/handoff gains into a future remote metadata and MCP platform without turning `GoREveal` into a generic RE suite

Why this is later:
- this sprint depends on the local review/handoff contract being mature
- it also depends on server/control-plane direction being explicit enough that MCP becomes a transport layer, not a moving product target

Primary ownership:
- future `api`
- `cmd/goreveal` / future `gorectl`
- docs and interop contracts

Allowed scope:
- remote MCP surfaces through `gorectl`
- metadata-network and artifact-transfer groundwork
- explicit host-platform and server-side handoff contracts

Disallowed scope:
- broad workspace mutation inside `GoREveal`
- general-purpose RE automation beyond Go-native truth and transfer
- broad server feature growth before operator value is proven

Evidence:
- interop planning docs
- contract docs
- bounded operator/API artifacts if sequencing later justifies them

PM+DEV task frame:
- `PM-18.1` define which remote interop stories are worth building first
  Indicator: one ranked list among `gorectl`, MCP transport, artifact sessions, and remote metadata APIs
- `DEV-18.1` keep remote interop transport-shaped rather than UI-shaped
  Indicator: new work reuses existing handoff/export contracts instead of inventing a second product surface

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-18.1` | PM | rank remote interop stories by operator value | `4` | `2` | `4` |
| `DEV-18.1` | DEV | reuse existing handoff/export contracts in remote interop | `3` | `4` | `4` |

### Sprint 19: Public Release Readiness and Licensing

Primary goal:
- remove the highest-friction public-release blockers without letting release/admin work displace product truth

Why this is later:
- the strategic review still treats repository licensing as a public-release blocker
- the product is now strong enough that release posture matters more than another speculative parser slice
- release/readiness work only makes sense once the local workflow and later server/interop direction are explicit enough to describe honestly

Primary ownership:
- docs
- repo root release artifacts
- light CLI/docs polish only if needed

Allowed scope:
- `LICENSE` decision and repository-baseline release posture
- README/release-note cleanup
- supported-scope wording for the current local-first product
- evidence-packaging guidance for comparisons, protected matrix, and operator workflows

Disallowed scope:
- feature growth for its own sake
- broad commercialization plumbing
- server implementation by stealth

PM+DEV task frame:
- `PM-19.1` close the repository release blockers
  Indicator: explicit license decision, supported-scope wording, and public-release checklist
- `DEV-19.1` implement only the smallest repo-facing release artifacts
  Indicator: no product drift, only release/trust-facing files and wording changes

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-19.1` | PM | close license and release-policy blockers | `5` | `2` | `5` |
| `DEV-19.1` | DEV | land minimal repo-facing release artifacts | `4` | `2` | `5` |

### Sprint 20: Evidence Expansion and Comparative Automation

Primary goal:
- turn the current strong comparison/evidence posture into a more repeatable product baseline instead of a one-off planning win

Why this is later:
- recent roadmap decisions are now driven by real external-binary and universal-workbench comparisons
- that comparison discipline is becoming one of the product's differentiators
- the next low-regret gain after release-readiness work is broader repeatable evidence, not speculative parser depth

Primary ownership:
- corpus
- scripts
- docs

Allowed scope:
- wider external-binary comparison coverage
- normalized comparison tables and reports
- differential/report automation improvements
- clearer module- and capability-level comparison tracking

Disallowed scope:
- copying external tool behavior into `core`
- turning baseline output into ground truth
- broad new recovery claims without fixture and differential backing

PM+DEV task frame:
- `PM-20.1` define the evidence coverage target
  Indicator: one explicit matrix for target binaries, baseline tools, and report outputs
- `DEV-20.1` automate comparative reporting without changing truth semantics
  Indicator: wider reports and scripts, but no new recovery claims without fixture and differential backing

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-20.1` | PM | define evidence-coverage matrix | `4` | `1` | `5` |
| `DEV-20.1` | DEV | automate wider comparative reporting | `4` | `2` | `5` |

### Sprint 21: Build Correlation and Version Tracking

Primary goal:
- turn the current bounded transfer and review surfaces into a stronger build-family workflow without jumping to broad semantic diffing

Why this is later:
- the current roadmap first needs local workflow maturity, later platform sequencing, and clearer release/evidence posture
- comparison and review work now show that the strongest workflow-shaped future differentiator is cross-build reuse of trusted Go-native truth

Primary ownership:
- `storage/diff`
- `engine`
- docs

Allowed scope:
- build-to-build transfer workflow deepening
- bounded version-tracking surfaces
- stronger analyst reuse of accepted transfer truth

Disallowed scope:
- decompiler-style semantic diffing
- broad matcher expansion without measured need
- generic binary diff ambitions beyond Go-native truth and transfer

PM+DEV task frame:
- `PM-21.1` define what “version tracking” means in product terms
  Indicator: one narrow scope around build-family reuse of trusted names, packages, and analyst actions
- `DEV-21.1` deepen transfer only where existing accepted/review surfaces already prove reuse value
  Indicator: one bounded build-correlation surface that extends current transfer workflow instead of replacing it

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-21.1` | PM | define narrow version-tracking scope in product terms | `5` | `2` | `4` |
| `DEV-21.1` | DEV | extend current transfer workflow into build correlation | `5` | `3` | `4` |

### Sprint 22: Metadata Knowledge Network

Primary goal:
- compound trusted Go metadata across analyses and families without turning `GoREveal` into a generic knowledge-graph platform

Why this is later:
- this sprint depends on version-tracking value being real first
- it also depends on remote metadata and later platform contracts being explicit enough to support durable reuse

Primary ownership:
- future `api`
- storage/metadata docs
- future `gorectl` / remote metadata planning

Allowed scope:
- metadata reuse models
- trusted cross-analysis linkage
- analyst-facing metadata-network planning and bounded artifacts

Disallowed scope:
- broad generic graph platform work
- non-Go-centric entity modeling by default
- replacing canonical per-analysis truth with aggregated guesses

PM+DEV task frame:
- `PM-22.1` define the first metadata-network unit of value
  Indicator: one explicit answer to “what reusable metadata entity gives analysts leverage first?”
- `DEV-22.1` keep the metadata network additive, not authoritative
  Indicator: aggregated knowledge never overwrites canonical per-analysis truth

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-22.1` | PM | define the first metadata-network unit of analyst value | `4` | `3` | `3` |
| `DEV-22.1` | DEV | keep metadata-network work additive to canonical truth | `4` | `4` | `3` |

### Sprint 23: Analyst Workspace Automation and Replay

Primary goal:
- turn the mature local CLI, review, and handoff surfaces into a replayable analyst-workspace automation layer without collapsing `GoREveal` into a generic RE-suite controller

Why this is later:
- it depends on server/control-plane foundations, remote interop, and build-correlation value already being real
- the current product still gains more from better truth and transfer semantics than from automating analyst sessions end-to-end

Primary ownership:
- future `api`
- future `gorectl`
- docs/ops

Allowed scope:
- replayable workspace/handoff manifests
- bounded automation around existing export and handoff contracts
- operator-safe session/review replay semantics

Disallowed scope:
- GUI-first orchestration platform work
- host-tool mutation logic inside `core`
- broad remote execution ambitions before control-plane maturity

PM+DEV task frame:
- `PM-23.1` define the first analyst replay unit of value
  Indicator: one explicit answer to “what repeatable analyst action bundle should be replayable first?”
- `DEV-23.1` keep workspace automation contract-first and replayable
  Indicator: one bounded replay artifact can drive an existing handoff/export flow without redefining recovery semantics

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-23.1` | PM | define the first replayable analyst-workspace unit of value | `4` | `3` | `3` |
| `DEV-23.1` | DEV | keep workspace automation contract-first and replayable | `4` | `4` | `3` |

### Sprint 24: Comparative Knowledge Packs and Decision Support

Primary goal:
- turn the accumulated comparison, corpus, and metadata posture into reusable decision-support artifacts without replacing canonical truth with opaque scoring

Why this is later:
- it depends on evidence automation, metadata reuse, and workspace automation being real first
- before that point, any “knowledge pack” would be mostly prose and stale too quickly

Primary ownership:
- docs/plans
- future metadata/storage surfaces
- future `api`

Allowed scope:
- reusable analyst-facing comparison packs
- bounded decision-support summaries over trusted evidence
- explicit ranking and recommendation artifacts derived from measured comparison state

Disallowed scope:
- hidden recommendation engines
- replacing per-analysis evidence with opaque aggregate scores
- non-Go-centric knowledge products by default

PM+DEV task frame:
- `PM-24.1` define the first decision-support artifact analysts would actually reuse
  Indicator: one named knowledge-pack artifact with a measurable analyst question it answers
- `DEV-24.1` keep knowledge-pack logic evidence-backed and inspectable
  Indicator: every recommendation can be traced back to measured comparison or metadata inputs

Seed backlog:

| Task | Type | Outcome | Value | Risk | Evidence |
| --- | --- | --- | ---: | ---: | ---: |
| `PM-24.1` | PM | define the first reusable comparative knowledge-pack artifact | `4` | `3` | `3` |
| `DEV-24.1` | DEV | keep decision-support logic inspectable and evidence-backed | `4` | `4` | `3` |

## Deferred Sprint Themes

Keep explicitly deferred:
- broad parser expansion
- broad deobfuscation sprint
- self-sufficient RE-suite / decompiler roadmap
- TUI/Web-first roadmap
- server/control-plane implementation beyond bounded foundation work

## Sprint Progress Reading

Current reading:
- `Sprint 12`: essentially complete for the current declared scope, but still the active umbrella lane until the handoff/value work fully settles
- `Sprint 13`: now re-baselined as workstation handoff contract work, with the first meaningful contract slices already landed and the sprint now clearly closer to completion than to definition
- `Sprint 14`: now materially underway through `transfer_review_plan` plus `diff next sqlite`, with package-ordered attached review items, self-contained `recommended_actions`, an explicit `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` snapshot, and an `upcoming_packages` horizon that already carries strongest-match and sample-pair context; the sprint now also has a frozen stop-condition, and the current weighted reading is that this threshold is likely already met for the declared local operator scope
- `Sprint 15`: active lane completed for the current bounded scope and now frozen by default, with three thin confidence slices landed and `DEV-15.4` held unless one new named inference step is demonstrated
- `Sprint 16`: PM ranking is now landed, with `garble`-class workflows ranked first and the first workflow-shaped pain point explicitly named; DEV remains evidence-driven and still held until that gap is tested through one bounded rerun
- `Sprint 17`: ordered but deferred behind local workflow maturity
- `Sprint 18`: ordered but clearly downstream of both local workflow and server/control-plane foundations
- `Sprint 19`: later-horizon release/readiness sprint
- `Sprint 20`: later-horizon evidence/comparison automation sprint
- `Sprint 21`: later-horizon build-correlation/version-tracking sprint
- `Sprint 22`: later-horizon metadata-knowledge-network sprint
- `Sprint 23`: later-horizon analyst-workspace-automation/replay sprint
- `Sprint 24`: later-horizon comparative-knowledge-pack sprint

## Current Recommendation

If continuing immediately:
1. treat the handoff artifact as sufficiently mature by default and keep any remaining `Sprint 13` work optional and thin
2. treat `Sprint 14` as frozen by default for the current scope unless `PM-14.4` can still name one concrete unresolved local inference step
3. keep `Sprint 15` frozen by default unless one clearly named remaining source-confidence inference step is demonstrated
4. treat `PM-16.1`, `PM-16.2`, and `PM-16.3` as landed, with `garble`-class workflows ranked first inside `Sprint 16`
5. keep `Sprint 16` evidence-first and external-orchestration-first, and do not open `DEV-16.*` work beyond the first bounded rerun until the named workflow-shaped protected gap is actually measured on neighboring garbled builds
6. keep `Sprint 17` and `Sprint 18` explicitly deferred until the local product and handoff lanes justify remote control-plane work
7. once the local product story is clearly mature, sequence `Sprint 19` before `Sprint 20` so release posture is explicit before comparison automation widens again
8. after that, sequence `Sprint 21` before `Sprint 22` so build correlation proves value before a broader metadata-network sprint
9. only after that, move into `Sprint 23` and then `Sprint 24`
