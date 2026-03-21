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
- differential testing against multiple baselines
- thin IDA/Ghidra export/import paths
- staged `golangci-lint` policy imported from `gobfd`
- green `make lint` on top of the imported `gobfd` policy
- an increasingly dense bounded runtime-truth layer in `Sprint 12`

The project is now strongest as a platform product and weakest as a runtime-semantic recovery engine.

That is the key strategic fact.

The new market-facing product fact is this:
- `GoREveal` should not try to replace `IDA`, `Ghidra`, `JEB`, or `Binary Ninja`
- it should become the best Go-native recovery, trust, and transfer layer that sits above or beside them

Quantified checkpoint:
- see `docs/plans/2026-03-20-goreveal-progress-assessment.md`
- the CLI surface is now wider too: current bounded runtime truth can be consumed directly through `inspect runtime`, not only through full `analyze`
- package navigation now also exposes explicit `has_source_evidence`, separating file-backed package truth from bounded fallback metadata without broadening package heuristics

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
- persistence and differential reporting
- export-driven integrations
- bounded runtime evidence collection
- product-level documentation and sprint discipline

### Where GoREveal Is Competitive But Not Yet Dominant

- build metadata
- package recovery
- source-tree projection on the canonical ELF fixture
- user-facing package/type navigation metadata
- function recovery for the current fixture

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

4. `Sprint 13` remains deferred.
   Reason: deobfuscation is still strategically important, but accuracy-first means runtime/type truth comes first.

5. market-differentiation epics stay behind the current accuracy lane.
   Reason: user-code isolation, version tracking, and metadata-network work become much stronger once runtime/type truth is less heuristic.

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

### Sprint 13

Status:
- deferred until `12B` proves a real semantic-runtime foothold

## Recommended Next Move

The next best move is no longer “one more bounded bridge” by default.

The next best move is:
- begin the first very small semantic decode in `Sprint 12`
- keep it fixture-local and tightly cross-checked
- do not attempt a full `moduledata` parser

That is now the highest-value risk-adjusted step.

That first step is now landed through a bounded rodata-relative typelink resolution bridge.

Immediate implication:
- do **not** rewrite package/type heuristics yet
- the current package/type surface is still more trustworthy when driven by `DWARF + source-tree` correlation than by the new runtime-semantic bridge
- the next best move is a second tiny semantic-runtime step, not a premature heuristic rewrite

## Immediate Task Queue

1. Hold broad package/type heuristics steady until runtime-semantic truth yields stable naming or scope evidence.
2. Keep the stripped-fixture contract stable across `inspect`, `analyze`, and `source-tree`; treat it as a maintained product boundary, not an experimental side path.
3. Follow `docs/plans/2026-03-20-goreveal-next-bounded-analyst-slices-plan.md` and take the next bounded analyst-facing slice from source-tree evidence/trust surfaces before revisiting broader heuristics.
4. Keep `Sprint 7` in maintenance mode by documenting any new semantic claim boundaries.
5. Re-evaluate whether the next move should be another tiny `Sprint 12` slice or a temporary return to capability transfer.

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
   - a second fixture
   - a very small heuristic-cross-check
   - a temporary pivot to another lane

That second-fixture move is now landed as a real checkpoint:
- a stripped ELF fixture exists
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
- but the next move inside it should no longer be blind bridge accumulation on the same fixture
- the best next discriminator is a second fixture or a very small runtime-to-heuristic cross-check
