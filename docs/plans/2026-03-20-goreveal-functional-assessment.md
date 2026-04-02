# GoREveal Functional Assessment

> Status: working product/strategy note
> Date: 2026-03-20
> Purpose: reassess current GoREveal functionality against the feature map and baseline Go RE tools, then translate that into an updated roadmap and sprint strategy.

## Current Product Position

`GoREveal` is no longer “a scaffold competing with existing tools someday”.
It is already a working platform with:
- canonical schema-driven analysis output
- a real CLI surface
- persistence and stored-run diffing
- the first bounded version-tracking-adjacent function matching inside `diff sqlite`
- differential testing against multiple baselines
- thin IDA/Ghidra export/import paths
- staged `golangci-lint` policy imported from `gobfd`
- green `make lint` on top of the imported `gobfd` policy
- portable repo-local Codex skills and subagents
- Podman-first task orchestration through `Taskfile.yml` and `scripts/dev/podman_runner.py`
- a more complete dev-image operator toolbox through `jq`, `yq`, `procps`, and `unzip`
- strict script verification through `ruff`, `ty`, `yamllint`, and `shellcheck`
- an increasingly dense bounded runtime-truth layer in `Sprint 12`

The project is now strongest as a platform product and weakest as a runtime-semantic recovery engine.

That is the key strategic fact.

The new tactical product fact is this:
- the first engine-owned code-peeling MVP is now real, not hypothetical
- `analysis.peeling`, `inspect peeling`, and `goreveal peel <binary>` already prove that user-code isolation can advance without contaminating `core`
- the first bounded fingerprint-assisted `stdlib` / `runtime` refinement is now also landed inside `engine/peeling`, with explicit `classification_evidence` instead of hidden classifier drift
- the first bounded version-tracking-adjacent matching surface is now also landed in `storage/diff`, through exact-name, source-location, source-file, and module-local normalized-name function matches with score/reason output
- the first thin analyst-facing transfer preview is now also landed in `storage/diff`, through bounded `transfer_candidates` built only over existing user-classified matches with `ready | review` disposition and projected left-to-right truth
- the first deterministic accepted-transfer surface is now also landed in `storage/diff`, through bounded `accepted_transfers` derived only from `ready` candidates without introducing mutable review state yet
- the first package-level transfer queue view is now also landed in `storage/diff`, through bounded `transfer_packages` summaries aggregated only from existing candidates and accepted transfers
- the first thin source-visibility response is now also landed through line-table-backed `source_tree` fallback with explicit `pathless_file_evidence`, narrowing the practical gap with file-list-oriented baselines without broadening parser claims
- that shifts the best next move from “another tiny parser bridge by default” to “carefully improve the usefulness of the new engine-owned differentiation and transfer layer”
- the next strategic question is now workflow-shaped, not parser-shaped: how to turn bounded matches and peeling evidence into a first-class analyst-facing transfer surface

Current weighted answer:
- the next best measured move is no longer another architecture-specific triage pass
- reason: the protected matrix is now widened through `arm64`, the canonical operator path still matches the direct measured path, and the latest bounded triage slice removed the old Linux architecture split by adding a truthful section-backed ELF `.text` range for protected address hints; `address_only` is now portable across `linux/amd64` and `linux/arm64`, while named-function recovery still remains the real protected gap
- the newest compact analyst-facing increment now also makes that foothold easier to consume operationally by projecting whether the foothold is backed by `moduledata_text` or by the fallback ELF `.text` section, plus a bounded foothold span, without widening recovery claims
- that same compact protected runtime surface is now also locked at export-contract and plugin-consumer boundaries, so the current protected reading is no longer only a CLI/runtime detail
- the current decision point is therefore narrower:
  - either add one more very small protected-specific analyst surface
  - or deliberately return to workflow/value work instead of opening a broad parser or deobfuscation lane
- the weighted recommendation is now to return to workflow/value work by default
  - reason: the protected lane now looks more stabilized than blocked
  - only reopen it immediately if one new protected-specific analyst surface can be justified by a specific measured analyst pain point
- the strongest concrete next candidate inside that workflow/value lane is a first analyst-facing transfer review surface over already-landed `matched_functions`, `transfer_candidates`, `accepted_transfers`, and `transfer_packages`
- that first candidate is now landed in compact form through `storage/diff transfer_review`, plus explicit `projected_package` on the transfer surfaces
- that next analyst-facing step is now also landed in compact form through package-first `storage/diff transfer_review_packages`
- the first bounded operator-facing action bundle is now also landed through `storage/diff transfer_review_focus`
- that bounded CLI/operator projection is now landed too through `goreveal diff review sqlite <database> <left-id> <right-id>`
- that same dedicated CLI path now also carries a compact machine-readable `handoff` block with left/right input context and recommended workstation targets, so focused review no longer starts from raw queue state alone
- that same workstation-facing bridge is now also exposed as a dedicated operator surface through `goreveal diff handoff sqlite <database> <left-id> <right-id>`
- the first bounded next-step actionability slice is now also landed through `storage/diff transfer_review_plan` and `goreveal diff next sqlite <database> <left-id> <right-id>`
- that same next-step slice now also carries package-ordered attached review items plus self-contained `recommended_actions`, a compact `review_checklist`, a compact `review_snapshot`, explicit `review_progress`, a compact `up_next` package snapshot, and an `upcoming_packages` horizon with sample pair and strongest-match context, so operators no longer need to jump back to the larger review queue just to see what the next package bundle contains, reconstruct the obvious follow-on commands by hand, manually count what remains after the current package, guess what “done for this bundle” means, or parse the full queue to see what comes immediately after the focus bundle
- the next workflow/value question is therefore no longer “should we return there”, and no longer “how do we expose the focused review pass at all”, but how to harden the next workstation-facing handoff on top of that queue, triage layer, first-pass bundle, dedicated review CLI path, and dedicated handoff CLI path
- the current weighted recommendation remains to keep that work in bounded review-action and handoff surfaces, not matcher rules and not protected-specific semantic slices by default
- the new `rehelp` and remote RE-lab inventory research makes the broader environment clearer too:
  - `GoREveal` is not operating in a vacuum
  - the measured workstation already has `ida-pro`, `ghidra`, `jeb`, `rizin`, `diaphora`, `binexport`, `ida-pro-mcp`, and dynamic/symbolic sidecars such as `frida`, `angr`, `qiling`, `unicorn`, `uftrace`, and `z3`
  - that strengthens the case for treating `GoREveal` as the Go-native truth/transfer/orchestration layer inside a richer workstation, not as a generic RE suite
- the weighted next move therefore becomes slightly sharper:
  - keep the current workflow/value lane primary, but treat the focused review CLI path as landed rather than still pending
  - then harden explicit host-platform MCP and workstation handoff planning on top of the now-landed `diff handoff sqlite` bridge
  - only then widen protected-specific or parser work again by default
- the fresh external rerun sharpens that reading further:
  - `GoREveal` now shows real file visibility on the measured external `ELF`, `PE`, and `Mach-O` targets rather than only function/package breadth
  - `GoREveal` is ahead on the fresh Go-native comparison for function/package/peeling coverage across the rerun targets
  - `redress` remains the strongest compared source/file-oriented baseline, while universal workbenches still dominate workspace mutation and generic RE breadth
  - the remaining product gap is therefore more clearly workflow- and handoff-shaped than raw format-footprint-shaped
  - the current intermediate repo-local rerun confirms the same reading on the rich `ELF` fixture, the in-repo `PE` and `Mach-O` fixtures, and local `rclone-linux-amd64`
  - the latest bounded timing sample on that same `rclone-linux-amd64` binary is tightly clustered across all four tools (`redress 1.13s`, `GoReSym 1.18s`, `gore 1.22s`, `goreveal analyze 1.23s`), so efficiency is still not the next blocking product concern even though `GoREveal` is not claiming a timing lead on this runner
  - the next roadmap question is now mostly sequencing:
    - treat `Sprint 14` review workflow actionability as materially underway through `transfer_review_plan` and `diff next sqlite`, with `recommended_actions`, `review_progress`, `up_next`, and `upcoming_packages` now making that next-step bundle more self-contained for actual operator use
    - keep any further `Sprint 13` workstation-contract lock optional and very thin
    - keep `Sprint 15` and `Sprint 16` conditional rather than automatic

The new market-facing product fact is this:
- `GoREveal` should not try to replace `IDA`, `Ghidra`, `JEB`, or `Binary Ninja`
- it should become the best Go-native recovery, trust, and transfer layer that sits above or beside them

Quantified checkpoint:
- see `docs/plans/2026-03-20-goreveal-progress-assessment.md`
- current weighted reading is now `platform = 99%`, `accuracy engine = 74%`, `overall roadmap = 99%`
- the CLI surface is now wider too: current bounded runtime truth can be consumed directly through `inspect runtime`, not only through full `analyze`
- the CLI surface now also exposes the first user-code isolation view through `inspect peeling` and `peel`, and that surface now carries explicit `classification_evidence`
- package navigation now also exposes explicit `has_source_evidence`, separating file-backed package truth from bounded fallback metadata without broadening package heuristics
- `docs/plans/2026-03-31-goreveal-strategic-review.md` now records the next roadmap reshaping step after the current runtime checkpoint
- `docs/plans/2026-04-01-goreveal-initial-baseline-comparison-results.md` now records the fresh external Go-native rerun
- `docs/plans/2026-04-01-goreveal-universal-re-workbench-comparison.md` now records the role-based comparison with general-purpose RE workbenches
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md` now records the active post-`Sprint 12` sprint sequence

## Comparative Summary

### Against `gore`

`GoREveal` is already ahead in product shape:
- canonical schema
- first-class CLI instead of library-only posture
- SQLite persistence
- stored-run diffing
- better testing and differential discipline

`GoREveal` is roughly competitive on:
- build info
- package presence
- narrow function overlap
- user-facing function navigation metadata
- source-backed package/type metadata usability

`GoREveal` is still behind on:
- less heuristic package truth when source evidence is weak
- typelink-driven type recovery depth
- broader version-family confidence

### Against `redress`

`GoREveal` is already competitive or better on:
- machine-readable output
- grouped source-package projection
- explicit external-package visibility
- product-level reproducibility

`GoREveal` is still behind on:
- deeper source reconstruction breadth
- richer stripped-binary presentation workflows
- broader operator-facing source projection confidence

### Against `GoReSym`

`GoREveal` is now clearly on the board:
- `pclntab`-based functions on the canonical ELF fixture
- bounded `moduledata` evidence
- bounded typelink/itablink/range cross-checks
- file/function overlap on the current fixture

But it is still behind in the most important accuracy-first areas:
- true runtime-semantic decoding from `moduledata`
- typelinks-driven type recovery
- robustness across versions, formats, and malformed layouts
- explicit support claims across version families

This is now the main strategic gap.

### Against `GoResolver`

`GoREveal` is ahead in:
- clean schema boundary
- raw/refined separation
- platform discipline

`GoREveal` is behind in:
- CFG-guided symbol recovery
- obfuscation-oriented workflow depth
- multi-version reference orchestration

This remains a later-stage capability, not the next core priority.

### Against `gostringungarbler`

`GoREveal` is ahead in:
- product integration
- refined-layer design
- ability to make string refinement part of a broader analysis result

`GoREveal` is behind in:
- actual garble-literals deobfuscation depth
- targeted support for literal-obfuscated samples

This is still a bounded future transfer, not today’s main gap.

### Against `AlphaGolang`

`GoREveal` is already cleaner architecturally:
- thin adapters
- export-first plugin contract
- no plugin-side recovery logic

`GoREveal` is behind in:
- mature analyst workflow helpers inside IDA
- stepwise UX for malware reversing tasks

This is convenience work, not accuracy-critical work.

### Against The Major RE Platforms As A Category

`GoREveal` is not competing on:
- general-purpose decompilation breadth
- debugger maturity
- large GUI workflow surface

It can compete on:
- Go-native recovery depth
- Go-native trust/provenance clarity
- Go-native metadata reuse and transfer
- Go-native workflow acceleration across many similar binaries

## Feature Map Reassessment

### Where GoREveal Is Already Strong

- schema-first platform design
- reproducible container-first verification
- explicit operator automation through `make`, `task`, and a shared Podman runner
- portable agent/skill configuration instead of chat-history-only workflow state
- persistence and differential reporting
- export-driven integrations
- bounded runtime evidence collection
- first engine-owned code peeling over canonical truth
- first bounded version-tracking-adjacent matching over canonical function surfaces
- product-level documentation and sprint discipline

### Where GoREveal Is Competitive But Not Yet Dominant

- build metadata
- package recovery
- source-tree projection on the canonical ELF fixture
- user-facing package/type navigation metadata
- function recovery for the current fixture
- first-pass user-code isolation

### Where GoREveal Is Still Clearly Behind

- runtime-semantic `moduledata` decoding
- typelinks-driven type recovery
- cross-version runtime confidence
- PE/Mach-O depth
- deobfuscation depth
- mature RE-tool UX workflows

## Strategic Product Opportunity

The strongest product opportunities above the current roadmap are:

1. `Go code peeling / user-code isolation`
   This is the strongest near-term differentiator because it attacks the biggest practical pain in Go reversing: separating custom code from statically linked noise.

2. `Go version tracking / build correlation`
   This is the strongest workflow epic because it can transfer recovered truth and analyst work across related Go builds or malware-family variants.

3. `Go metadata knowledge network`
   This is the strongest long-term moat because it compounds trusted Go metadata over time in a way that typical point tools do not.

These opportunities should shape future roadmap epics, but they should not displace the current `Sprint 12` accuracy lane yet.

## Process Assessment

Overall process reading:
- repo operations are already mature for the current scope: containerized verification, task automation, lint discipline, plans, skills, contract docs, and the dev-image operator toolbox are all strong
- delivery style is also healthy: bounded vertical slices land with tests, snapshots, exports, and docs instead of architecture-only churn
- the weakest process area is not implementation discipline but evidence breadth: fixture families and cross-format checkpoints are still narrower than the strategy wants
- the main planning risk is stale narrative drift after rapid bounded slices, so progress and handoff docs must keep being updated together with capability changes

Weighted module/process reading:
- strongest areas: `schema`, `engine`, `cmd/goreveal`, and docs/ops
- solid but still incomplete areas: `core`, `storage`, `plugins`, and `corpus`, with `storage` now effectively acting like a product-strength workflow layer even though later mutable review/server state is still intentionally absent
- intentionally early areas: `deobfuscation`, service/API, and performance/SIMD

Planning consequence:
- the local product is now mature enough that the next roadmap question is mostly sequencing, not discovery
- that sequencing should stay explicit:
  - active now: `Sprint 16` PM-side protected workflow target ranking, with both `Sprint 14` and `Sprint 15` held at their frozen local stop-conditions unless measured gaps reopen them
  - conditional after that: one bounded `Sprint 16` DEV slice only if PM ranking names a first protected workflow gap clearly enough to justify corpus/comparison work
  - ordered later: `Sprint 17` server control-plane foundations, then `Sprint 18` metadata/remote interop platform
  - later after that: `Sprint 19` public-release/licensing hardening, then `Sprint 20` evidence/comparison automation
  - strategic horizon after that: `Sprint 21` build correlation/version tracking, then `Sprint 22` metadata knowledge network
  - long horizon after that: `Sprint 23` analyst workspace automation/replay, then `Sprint 24` comparative knowledge packs and decision support
- this keeps `PostgreSQL`, `gorectl`, remote MCP, and object-storage work in the roadmap without letting them displace the current working product

Current selected `Sprint 15` gap:
- one bounded source-evidence confidence contract over already-known `function`, `package`, and `source_tree` truth
- first landed slice:
  - `source_tree` and enriched package surfaces now project `source_evidence_kind = dwarf_paths | line_table_files | package_fallback`
  - `source_tree` now also projects a compact `source_evidence_summary`, so the high-level package evidence landscape is explicit without scanning every package node
  - that summary now also carries per-evidence-class file counts, so operators do not need to hand-sum file density across package nodes just to judge how rich each evidence class really is
  - this stays inside already-known truth and does not broaden parser claims
- remaining pain point:
  - operators can now see the direct evidence class, the compact tree-level summary, and the per-class file-density cue, so any next slice must remove one clearly named remaining inference step rather than merely restating the same truth in another shape
- weighted stop-condition:
  - if no such named inference step is left after the current three thin slices, freeze `Sprint 15` by default and move the PM lane to `Sprint 16` target ranking instead of inventing more explanation-only fields
- current weighted decision:
  - no such named remaining inference step is currently established in the planning/evidence set, so `Sprint 15` should now be treated as frozen-by-default for the current scope, not as an automatically continuing micro-surface lane
- immediate PM consequence:
  - the active planning queue should therefore move into `Sprint 16` target ranking instead of inventing another confidence-only field for `Sprint 15`
- why it outranks a protected return right now:
  - the protected lane is currently stabilized enough for its declared scope
  - `redress` remains the strongest practical source/file-oriented baseline
  - the gap can be addressed through confidence semantics over current truth rather than parser expansion
- next weighted move after the freeze:
  - `Sprint 16` PM ranking is now the active planning lane
  - the first protected workflow target should be `garble`-class workflows on the current enterprise-gated sample, because that class has the strongest measured ambiguity and commercial relevance among already-exercised protected cases
  - the next DEV move there should still be corpus/comparison-first and should describe the first protected analyst pain point in workflow terms before any new analyst-facing or deobfuscation surface is proposed

## Empirical External Binary Check

The current repo state is now strong enough for a bounded real-world check outside the fixture corpus.

External sample results:
- `/opt/projects/repositories/ocserv-agent/bin/ocserv-agent-linux-amd64`
  - stripped static Go `ELF`
  - `analyze` recovered `build_info.path = "github.com/dantte-lp/ocserv-agent/cmd/agent"`, strong `pclntab` function coverage, and bounded runtime evidence with `trust_summary = "section_heuristic"`
  - `peel` produced a useful user-only view with `main` package functions such as `main.main`, `main.setupLogger`, and `main.runGenCert`
- `/opt/projects/repositories/GoReSym/testproject/testproject.exe`
  - external Windows `PE`
  - `analyze` now recovers `build_info`, the bounded `.text` / `.rdata` plus `pclntab` header heuristic surface, and real `functions`, `packages`, and `peeling`
  - current result is no longer posture-only; it is now a real first `PE` function/package/peeling foothold, while still avoiding broad `PE` semantic claims
- `/opt/projects/repositories/hashicorp-re/bin/Keygen.exe`
  - external Windows `PE`
  - `analyze` again recovered `build_info` and the same bounded `PE` runtime posture surface
  - this still confirms low-confidence `PE` runtime posture beyond the in-repo fixture, but the product now also has a separate first `PE` function/package/peeling foothold on maintained `PE` samples

Practical reading:
- `GoREveal` is already meaningfully useful on real stripped Go `ELF` binaries for build info, function surfaces, and first-pass user-code isolation
- `GoREveal` is already useful on real `PE` binaries for build info, bounded runtime posture, and a first function/package/peeling foothold, but not yet for deep `PE` semantic recovery or broader source/file visibility on that format
- a first thin source-visibility increment is now also landed: when DWARF-backed source projection is unavailable, `source_tree` can still expose package/file evidence from recovered function line tables, clearly marked as `pathless_file_evidence`
- this is enough to say the tool already works outside the fixture corpus, but not enough to broaden support claims aggressively

Cross-platform matrix follow-up:
- `docs/plans/2026-03-31-goreveal-external-binary-matrix-evaluation.md`
- one open-source target, `rclone v1.73.3`, was checked across `linux/windows/macos` and `amd64/arm64`
- outcome:
  - `linux-amd64` and `linux-arm64`: strong stripped-`ELF` results with large function/package surfaces and meaningful `peel`
  - `windows-amd64` and `windows-arm64`: build-info plus bounded `PE` runtime posture and now large recovered `functions` / `packages` / `peeling` surfaces after the first `PE` foothold
  - `darwin-amd64` and `darwin-arm64`: after the new bounded `Mach-O` foothold, large function/package surfaces plus meaningful `peeling`, while `runtime` still remains absent

This changes the practical reading in one important way:
- the product is no longer best described as “fixture-proven only”
- it is now better described as:
  - real and useful on stripped Go `ELF`
  - early but already materially useful on `PE`
  - newly useful on `Mach-O` for function/package/peeling recovery, but still early on runtime/semantic depth

## Strategic Update

The previous roadmap emphasized:
1. `Sprint 7` proof expansion
2. `Sprint 11` source/package/type usability transfer
3. `Sprint 12` runtime depth

That was correct earlier.

It is no longer the best ordering.

### New Priority Order

1. `Sprint 12` becomes the primary lane.
   Reason: the biggest remaining product gap is now runtime-semantic accuracy, not documentation, not proof plumbing, and not source projection shape.

2. `Sprint 7` moves into maintenance mode.
   Reason: differential evidence is already good enough to support the current claims. It should still grow, but it is no longer the highest-leverage active lane.

3. `Sprint 11` is effectively a completed checkpoint.
   Reason: the current source/package/type usability surface is already materially better than the earlier baseline and no longer the main bottleneck.

4. `Sprint 13` is now re-baselined around workstation handoff contract work, not a deferred deobfuscation-first lane.
   Reason: deobfuscation is still strategically important, but accuracy-first means runtime/type truth comes first.

5. market-differentiation epics stay behind the current accuracy lane.
   Reason: user-code isolation, version tracking, and metadata-network work become much stronger once runtime/type truth is less heuristic.

6. later server/control-plane work is now ordered, but still clearly later.
   Reason: `PostgreSQL 18`, `pgx/v5 + sqlc`, `goose`, `ConnectRPC`, `River`, `gorectl`, and object storage are no longer undecided, but they still solve a later product phase, not the current local workflow phase.

## Sprint Translation

### Sprint 7

Status:
- active, but no longer the lead lane

Role now:
- maintain evidence quality
- add divergence notes only when new claims are introduced
- avoid letting wrapper/evidence work consume the main roadmap

### Sprint 11

Status:
- checkpoint complete

Role now:
- only reopen if runtime-driven truth later changes package/type/source classification rules

### Sprint 12

Status:
- primary execution lane

Sub-phases now:
- `12A`: bounded runtime evidence accumulation
  - effectively strong enough for the current fixture
- `12B`: first minimal semantic decode
  - should start from one tiny, testable semantic bridge
  - example directions:
    - bounded typelinks header semantics
    - bounded moduledata range semantics
    - one minimal decoded runtime structure field set with explicit fixture-local scope
- `12C`: runtime trust and second-fixture checkpoint
  - add a compact runtime trust/evidence summary
  - validate the next bounded step on a Windows `PE` fixture instead of extending the same ELF bridge chain again

### Sprint 13

Status:
- re-baselined
- now starts with workstation handoff contract hardening over already-landed review/handoff surfaces
- older deobfuscation-only `Sprint 13` reading is no longer the active roadmap baseline

## Recommended Next Move

The next best move is no longer “one more bounded bridge” or “one more raw matching rule” by default.

The weighted decision is now:
- keep the current bounded `PE` runtime checkpoint stable
- keep `Sprint 7` in maintenance mode
- keep the current `engine/peeling` and `storage/diff` surfaces stable
- rerun the same real comparison after the new `PE` foothold and use that rerun to decide whether workflow/source-visibility value now outranks another bounded semantic slice

Why this is the best risk-adjusted move:
- it converts the current bounded matches into usable analyst workflow value instead of leaving them as report-only output
- it improves the strongest practical differentiators, user-code isolation and build-to-build transfer
- it builds on already-known canonical truth instead of inventing new recovery claims
- it stays above `core`, so the architecture boundary remains clean
- it strengthens future version-tracking and host-tool workflows more than another tiny parser bridge would

That first bounded refinement is now landed.
That first bounded matching surface is now landed too.
That first transfer-oriented projection preview is now landed too.
That first deterministic accepted-transfer surface is now landed too.
That first bounded `Mach-O` function/package/peeling foothold is now landed too.
The next decision is no longer whether `Mach-O` or `PE` need a first foothold at all; it is whether the post-`PE` comparison rerun points next toward workflow polish, thin source-visibility improvement, or one more bounded semantic/cross-format slice.

The external-binary check sharpens that decision further:
- stripped `ELF` value is already good enough that more transfer polishing is no longer the only rational next move
- `PE` still needs more depth later, but the current real-world result no longer justifies reopening a broad `PE` parser lane by default
- `Mach-O` now has a real foothold too, so the clearest remaining practical gap shifts away from format presence and toward source/file visibility plus richer transfer workflow value

Immediate implication:
- do **not** rewrite package/type heuristics yet
- do **not** reopen broad parser work by default
- do **not** treat the current `PE` checkpoint as permission for generic `PE` runtime claims
- do use the new `peeling` surface plus the new `PE` and `Mach-O` footholds as inputs to the next evidence-driven comparison pass
- do treat stripped external `ELF` success as confirmation that Linux value is already real
- do treat source/file visibility and transfer usefulness as the strongest current product questions after all three format footholds are landed
- do use `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md` as the active next-sprint baseline

## Updated Weighted Decision

There are now two rational steps:

1. `real baseline comparison`
   Best when the goal is to choose the next increment based on measured gaps against `gore`, `GoReSym`, and `redress` rather than on intuition.

2. `annotated/persisted transfer workflow`
   Best when the comparison confirms that Linux/server analyst throughput is now a stronger gap than one more bounded semantic slice.

Current weighted recommendation after the new `PE` foothold:
- rerun the same comparison after that new `PE` surface
- if the rerun confirms current readings, choose annotated/persisted transfer work or a thin source-visibility improvement ahead of another parser expansion
- only choose another bounded semantic/cross-format slice first if the rerun shows a concrete measured gap that workflow polish cannot address

Reason:
- `ELF` already proved useful on real external binaries
- `PE` already has an honest bounded foothold
- the first `Mach-O` foothold is now landed too
- the next decision is therefore no longer a missing-format problem; it is a measured workflow/source-visibility problem, now on top of the first landed line-table-backed visibility fallback

## Immediate Task Queue

1. Hold broad package/type heuristics steady until runtime-semantic truth yields stable naming or scope evidence.
2. Keep the compact runtime trust/evidence summary stable across `inspect`, `analyze`, and thin exports.
3. Keep the stripped-fixture contract stable across `inspect`, `analyze`, and `source-tree`; treat it as a maintained product boundary, not an experimental side path.
4. Keep the bounded Windows `PE` checkpoint stable at its current claim level unless a clearly evidence-backed next slice appears.
5. Keep code-peeling and diff matching in `engine` / `storage`, not in `core` or `deobfuscation`.
6. Treat the new comparison findings as the current evidence baseline.
7. Prefer the next bounded increment to be the full post-`PE` comparison rerun, not another parser increment by default.
8. After that rerun, prefer transfer-workflow polish or thin source-visibility work unless measured gaps clearly point back to semantics.
9. Add more matching rules only when they are justified by concrete annotated-transfer workflow gaps rather than by curiosity-driven matcher expansion.

That next user-facing truth step is now landed for functions:
- `inspect functions` no longer exposes only raw names and addresses
- recovered functions now also carry bounded `package`, `import_path`, and `module_local` metadata
- recovered functions now also carry bounded `source_file` truth from the existing line table
- recovered functions now also carry bounded `source_line` truth from the same line table
- recovered functions now also carry bounded `autogenerated` truth for synthetic functions whose truthful source file is `<autogenerated>`
- recovered strings now also carry bounded absolute `addr` truth by combining already-known region base addresses with candidate offsets
- stripped `analysis.runtime` now also carries an explicit `firstmoduledata_from_go_module_fallback` bit, so operators can distinguish symbol-backed and fallback-backed `firstmoduledata_addr`
- `source-tree` package nodes now also carry explicit `has_file_evidence`, so rich and stripped outputs distinguish file-backed nodes from bounded empty-files fallback nodes
- `inspect packages` now carries explicit `has_source_evidence`, so package navigation distinguishes source-backed packages from bounded fallback-backed package metadata
- thin `IDA` / `Ghidra` exports now also preserve string `address` directly from canonical schema, avoiding plugin-side recomputation of already-known truth
- thin `IDA` / `Ghidra` exports now also preserve bounded function navigation metadata directly from canonical schema, again avoiding plugin-side recomputation of already-known truth
- thin `IDA` / `Ghidra` exports now also preserve bounded type navigation metadata directly from canonical schema, keeping package/scope/source-evidence context out of plugin-side heuristics
- thin `IDA` / `Ghidra` adapters now also prefer canonical string `address` over `offset`, so plugin behavior finally matches the stronger export contract
- thin `IDA` / `Ghidra` exports now also preserve canonical `runtime.trust_summary`, so runtime posture stays visible to host-tool workflows without new plugin-side classification logic
- this stays within the current strategy because it reuses existing function-name and `build_info` truth instead of expanding parser scope

## Post-Second-Semantic-Step Handoff

That second semantic step is now landed:
- the first typelink semantic bridge exists
- the current fixture also confirms that all resolved typelinks stay within `.rodata`
- bounded `types..etypes` and `pcHeader` / `funcnametab` bridges now strengthen the fixture-local runtime model without broadening it into generic runtime decoding
- a bounded `cutab` bridge now extends the same fixture-local `.gopclntab` layout hypothesis one slice further without turning it into a general parser
- a bounded `filetab` bridge now extends the same fixture-local `.gopclntab` layout hypothesis one slice further without turning it into a general parser
- a bounded `pctab` bridge now extends the same fixture-local `.gopclntab` layout hypothesis one slice further without turning it into a general parser
- a bounded `pclntable` bridge now extends the same fixture-local `.gopclntab` layout hypothesis one slice further without turning it into a general parser

This changes the handoff slightly:
- the project should continue in `Sprint 12`, not switch to `Sprint 10` yet
- the next `Sprint 12` work should no longer default to another same-fixture bridge
- the next candidate should preferably be either:
  - a second fixture validation move, or
  - a very small runtime-to-heuristic cross-check

Current recommendation:
1. keep `Sprint 7` in maintenance mode
2. keep `Sprint 11` heuristics unchanged except for bounded stripped-fixture UX fixes
3. treat the current `.gopclntab` bridge chain as a checkpoint, not as permission for blind continued expansion
4. choose the next `Sprint 12` move from:
   - a runtime trust summary
   - a second fixture
   - a very small heuristic-cross-check

That second-fixture move is now landed as a real checkpoint:
- a stripped ELF fixture exists
- a bounded Windows `PE` fixture also exists with build-info coverage, snapshot coverage, thin export coverage, and a low-confidence runtime section/header heuristic over `.text` / `.rdata` plus raw `pclntab` candidate fields
- the runtime layer now supports a bounded `.go.module`-based `firstmoduledata` fallback for this stripped ELF family
- `inspect packages` on the stripped fixture now preserves a useful `main` package surface by reusing `build_info.path` as bounded `import_path` truth while exposing external package `import_path` directly from function recovery
- `inspect types` on the stripped fixture now degrades to `[]` instead of JSON `null`, making the missing `DWARF`-backed type surface explicit without pretending runtime type decoding exists
- `source-tree` on the stripped fixture now returns a bounded fallback root plus module-local and external package nodes with empty file lists, instead of failing outright
- when a truthful type surface does exist without source-tree evidence, non-`main` types now keep direct `import_path` truth from parsed type packages and `main` still reuses bounded `build_info.path`
- canonical `analyze` now makes missing type truth explicit through `types: []` instead of silently omitting the field

## Risk Update

### Main Product Risk

The main risk is no longer underbuilding the platform.
The main risk is overfitting bounded runtime evidence and mistaking it for semantic support.

### Main Engineering Risk

The next semantic step could trigger parser explosion if it continues extending the same single-fixture `.gopclntab` layout without adding genuinely new information.

### Main PM Risk

If the roadmap keeps treating `Sprint 7` as the primary lane, the project will optimize for proof instead of closing the biggest accuracy gap.

## Bottom Line

`GoREveal` is already becoming a better product platform than the baseline tools in several areas.
It is not yet a better runtime-semantic recovery engine than `GoReSym`.

That gap should now dominate the strategy.

The new checkpoint conclusion is:
- `Sprint 12` is still the primary lane
- but the next move inside it should no longer be blind bridge accumulation on the same fixture or by default another format-foothold push
- the best next discriminator is now the post-`PE` comparison rerun, followed most likely by transfer-workflow polish or thin source-visibility work built on current `peeling`, `matched_functions`, `transfer_candidates`, and `accepted_transfers`, with the `PE` checkpoint and runtime trust summary held stable as bounded evidence anchors
