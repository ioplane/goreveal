# GoREveal Market / Killer Features Brainstorm

> Status: historical PM/product brainstorming; superseded for execution by the
> [RT1 product design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
> and [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md)
> Date: 2026-03-20
> Purpose: compare `GoREveal` with major RE platforms and typical Go RE utilities, then identify which external-platform capabilities are worth absorbing into `GoREveal` as differentiated product features.

## Why This Note Exists

`GoREveal` is no longer a pure parser experiment.
It is already a product-shaped platform with:
- canonical schema
- CLI and stored analysis
- differential testing
- thin RE-tool integrations
- an increasingly dense bounded runtime layer

That changes the product question.

The next PM question is no longer:
- “can we recover more Go metadata?”

It is now:
- “which capabilities from the broader RE-tool market should be absorbed into `GoREveal`, and which should remain delegated to external platforms?”

This note is a brainstorming and prioritization pass, not a commitment to implement everything below.

## Updated Ecosystem Signal

Recent `rehelp` plus RE-lab inventory work makes the adjacent operator environment more concrete.

Measured real-environment signals now include:
- host platforms:
  - `ida-pro`
  - `ghidra`
  - `jeb`
  - `rizin`
- diff and metadata-transfer adjacencies:
  - `diaphora`
  - `binexport`
  - `binsync`
- host-platform MCP signal:
  - `ida-pro-mcp`
- dynamic and symbolic sidecars:
  - `frida`
  - `angr`
  - `qiling`
  - `unicorn`
  - `uftrace`
  - `z3`

That does not mean `GoREveal` should absorb those engines.
It does mean the product can plan more confidently around:
- workstation interop
- transfer workflows
- remote orchestration
- host-platform MCP handoff

## What The Big RE Platforms Already Do Well

### IDA Pro

Current differentiators from official Hex-Rays materials:
- broad disassembly/decompiler coverage across many processor families
- strong debugger support, including local and remote debugging
- rich plugin model through IDAPython and the C++ SDK
- Lumina for function recognition and metadata reuse
- Teams for collaborative diff/merge/sync workflows inside IDA
- headless processing via `idalib`

Practical product reading:
- IDA is not just an analyzer; it is a long-lived analyst workspace
- its biggest strengths are collaboration, metadata reuse, plugins, and analyst workflow durability

### Ghidra

Current differentiators from official Ghidra documentation:
- full SRE framework, not just a disassembler/decompiler
- headless and user-interactive modes
- extensibility with Java and Python
- project repository / shared server / versioned artifacts
- debugger traces that can be recorded, revisited, and correlated with static analysis
- version tracking sessions that transfer matches and markup between binaries

Practical product reading:
- Ghidra is very strong at projectization, version tracking, and static/dynamic linking
- its version tracking model is especially important for malware-family and multi-build workflows

### JEB

Current differentiators from PNF official materials:
- strong Android/Dalvik depth
- native and bytecode decompilers across several architectures
- robust deobfuscation / cleanup / refactoring workflows
- multi-artifact projects
- headless automation and plugin APIs
- strong “analysis pipeline” orientation rather than only desktop use

Practical product reading:
- JEB shines when the target class is complex and obfuscated, especially mobile
- the relevant lesson for `GoREveal` is not “build an Android RE platform”
- the relevant lesson is “productize deobfuscation and workflow automation as first-class features”

### Binary Ninja

Current differentiators from official Binary Ninja materials:
- strong automation API in C++, Python, and Rust
- BNIL as an architecture-independent IL stack
- on-prem collaboration and automation through Enterprise
- WARP for reusable function identification and metadata transfer
- good debugger integration
- headless analysis and plugin ecosystem

Practical product reading:
- Binary Ninja is strongest where clean automation, reusable semantics, and collaborative metadata meet
- WARP is especially relevant because it is not just “signatures”, it is metadata transfer for matched functions

## What Typical Go RE Utilities Usually Do Not Do

Typical Go-focused tools such as `gore`, `redress`, `GoReSym`, `GoResolver`, and `gostringungarbler` are usually excellent at one narrow layer:
- extracting Go runtime/build metadata
- reconstructing packages/files/types/strings
- handling one obfuscation or symbol-recovery niche
- serving as reference CLI tools

But they usually do **not** provide a complete analyst platform:
- no strong multi-binary project model
- no real collaboration layer
- no durable metadata reuse server like Lumina/WARP
- no serious version tracking between builds/families
- no explicit trust/provenance model across all recovered entities
- no first-class persistence / diff / corpus workflow
- thin or nonexistent RE-suite integration

This is exactly where `GoREveal` can differentiate.

## Important Product Constraint

`GoREveal` should **not** try to become:
- a general-purpose decompiler competitor to IDA/Ghidra/JEB/Binary Ninja
- a debugger replacement
- a generic GUI RE suite

Those markets are already occupied, capital-intensive, and outside the current product advantage.

The right question is:
- what specialized Go-native layer can sit above or beside those platforms and become difficult to replace?

## Brainstorm: Capabilities Worth Absorbing

### 1. Go Metadata Knowledge Network

Inspired by:
- Hex-Rays Lumina
- Binary Ninja WARP

Potential `GoREveal` version:
- private or team-hosted metadata service for Go binaries
- upload/download high-confidence package names, type names, function names, comments, source-tree hints, build-family fingerprints
- transfer recovered metadata across exact or near-exact matches
- separate public, private, and trusted metadata feeds

Why it matters:
- Go binaries contain huge amounts of repeated stdlib / runtime / third-party code
- analysts repeatedly spend time rediscovering the same package/function/type truth
- a Go-native metadata network could eliminate repeated manual recovery work

PM value:
- very high

Moat potential:
- very high

Risk:
- medium to high
- requires stable identity/matching and careful trust model

### 2. Go Version Tracking / Build Correlation

Inspired by:
- Ghidra Version Tracking
- IDA Teams diff/merge mindset

Potential `GoREveal` version:
- compare two Go binaries or two malware-family builds
- correlate packages, functions, types, and strings across versions
- transfer analyst markup / recovered truth / names / package structure from one build to another
- keep a session artifact showing accepted, rejected, and uncertain transfers

Why it matters:
- Go malware, agents, and enterprise binaries are often released in many near-neighbor builds
- analysts often work across families or revisions, not one sample at a time
- version-aware transfer is higher value than one-shot extraction

PM value:
- very high

Moat potential:
- very high

Risk:
- medium
- needs a strong matching model, but not necessarily a full decompiler

### 3. Go Code Peeling / User-Code Isolation

Inspired by:
- the practical analyst workflows in all major RE suites
- WARP/Lumina style metadata reuse

Potential `GoREveal` version:
- classify and collapse stdlib / runtime / third-party dependency code
- surface “likely user-owned code” as the primary analyst view
- provide peel layers:
  - user code
  - third-party modules
  - stdlib
  - runtime

Why it matters:
- statically linked Go binaries are huge
- the analyst usually wants the 3-10% of code that is actually custom
- this is a pain point that general RE suites do not solve in a Go-native way

PM value:
- extremely high

Moat potential:
- extremely high

Risk:
- medium
- needs better matching/classification confidence, but can be staged

### 4. Trust-Aware Recovery UX

Inspired by:
- not a direct single-platform feature
- this is closer to a product philosophy gap in the market

Potential `GoREveal` version:
- every package/type/function/string/source-tree node carries provenance and trust signals
- operators can switch between:
  - raw truth
  - refined truth
  - transferred truth
  - inferred truth
- CLI, exports, and later APIs all preserve those boundaries

Why it matters:
- current RE tooling often buries uncertainty behind UI polish
- for Go recovery, trust boundaries matter a lot

### 5. Protected Commercial Go Software Workflow

Inspired by:
- the real OSS + Enterprise deployment model in Go products
- current commercial-hardening practice around `garble`, `-s -w`, `-trimpath`, `PIE`, anti-debugging, and runtime license checks

Potential `GoREveal` version:
- build-profile awareness for protected Go binaries:
  - stripped
  - trimpathed
  - garbled
  - PIE
  - enterprise-gated
- focused analyst surfaces for:
  - user-code isolation in heavily linked binaries
  - locating license checks and feature gates
  - tracking enterprise-only code paths across builds
  - recognizing common hardening patterns without pretending anti-obfuscation parity

Why it matters:
- a large commercial Go target class is neither malware nor toy OSS fixtures
- enterprise binaries often combine static linking, partial stripping, build-path removal, and licensing logic
- this target class maps directly to high-value customer workflows better than generic “more parser depth” messaging

PM value:
- very high

Moat potential:
- high

Risk:
- medium
- easy to overclaim if the work starts as anti-obfuscation branding instead of corpus-backed bounded capability

Recommended rule:
- treat this as a corpus/comparison and workflow epic first
- only promote individual parser or deobfuscation slices into the then-active lane when protected-binary evidence shows a concrete analyst pain point
- this could become a defining product property if kept consistent

PM value:
- high

Moat potential:
- medium-high

Risk:
- low to medium

### 5. Go-Aware Static/Dynamic Bridge

Inspired by:
- Ghidra Debugger trace model
- IDA debugger workflows

Potential `GoREveal` version:
- import traces from `dlv`, `gdb`, `rr`, or external debuggers
- correlate runtime goroutines, call stacks, heap objects, interfaces, and channels with static recovery
- persist static + dynamic evidence together

Why it matters:
- Go has runtime-heavy semantics that are hard to reconstruct statically
- dynamic evidence could massively improve confidence for selected targets

PM value:
- high

Moat potential:
- high

Risk:
- high
- operationally much heavier than the current product stage

### 6. Obfuscation Workbench For Go

Inspired by:
- JEB’s refactoring/automation orientation
- GoResolver / gostringungarbler niche strengths

Potential `GoREveal` version:
- bounded deobfuscation session model
- preserve raw/refined/transferred truth
- let analysts apply and review refinement passes instead of getting one opaque output

Why it matters:
- obfuscated Go malware is a real market
- general RE platforms are powerful, but not Go-obfuscation-native

PM value:
- high

Moat potential:
- medium-high

Risk:
- medium-high

## What Could Become A Killer Feature

If the question is “what could actually be a market-defining differentiator?”, the strongest candidates are:

### Candidate A: User-Code Isolation For Go

Short version:
- “Show me the real user code and hide the runtime noise.”

Why this is strong:
- universal pain point in Go reversing
- immediately understandable to users
- useful in CLI, API, and RE-suite integrations
- does not require replacing IDA/Ghidra/JEB/BN
- complements all of them

This is the most obvious near-to-mid-term killer feature candidate.

### Candidate B: Go Version Tracking / Cross-Build Markup Transfer

Short version:
- “Transfer recovered truth across Go builds and malware-family variants.”

Why this is strong:
- highly valuable to malware RE, threat intel, and internal product security teams
- stronger than one-shot extraction
- under-served by Go-native tools
- can become a real team workflow moat

This is the strongest workflow killer feature candidate.

### Candidate C: Private Go Metadata Network

Short version:
- “Lumina/WARP, but Go-native and trust-aware.”

Why this is strong:
- compounds over time
- improves every later analysis
- creates team lock-in and measurable productivity gains

This is the strongest long-term moat candidate.

## PM Prioritization

### Best Near-Term Opportunity

**Go code peeling / user-code isolation**

Reason:
- highest user pain relief relative to implementation risk
- fits current `GoREveal` architecture
- can be built incrementally from current package/function/type/runtime truth
- makes `IDA`/`Ghidra`/`JEB`/`Binary Ninja` integrations more valuable without trying to outdo them at decompilation

### Best Mid-Term Opportunity

**Go version tracking / build correlation**

Reason:
- leverages current persistence, diffing, and schema strengths
- aligns with malware-family and enterprise build workflows
- differentiates from narrow Go metadata tools
- the current bounded `matched_functions` plus `classification_evidence` surfaces now make the first transfer-oriented workflow a realistic next increment rather than a distant concept

## Historical Weighted Decision

At the March repository checkpoint, the recommended next move was:
- do not open a broader parser lane
- do not add more raw matching rules by default
- do build the first accepted or annotated transfer workflow on top of current `engine/peeling` and `storage/diff` outputs

Why this wins:
- it advances both near-term differentiators at once: user-code isolation and build correlation
- it turns existing bounded signals into analyst workflow value instead of keeping them as report-only metadata
- it stays above `core`, so the clean-room recovery boundary remains intact

### Best Long-Term Moat

**Private Go metadata knowledge network**

Reason:
- strongest compounding value
- closest analog to Lumina/WARP
- difficult for point tools to match

## What Not To Build Soon

Do not prioritize:
- a full debugger
- a native decompiler competitor
- a heavy GUI
- generic collaboration that is not Go-aware
- a plugin-first workflow that duplicates recovery logic inside IDA/Ghidra/JEB/BN

Those are expensive and dilute the product.

## Historical Strategic Recommendation

The then-current recommended direction was:

1. Treat `Sprint 12` as effectively complete for the current declared recovery/value scope, not as the default next expansion lane.
2. Keep `Sprint 13` focused on **workstation handoff contract hardening** over already-landed review and handoff surfaces.
3. Move next to **review workflow actionability** rather than reopening parser work by inertia.
4. Keep **thin semantic/source confidence** as a conditional lane, only if it clearly improves analyst value after the handoff/review path settles.
5. Treat **protected commercial Go workflows** as an evidence-first orchestration lane, not as a default native deobfuscation sprint.
6. Keep `IDA` / `Ghidra` / `JEB` / `Binary Ninja` / `Rizin` as host environments and analyst surfaces, not enemies to replace.

Related review:
- `docs/plans/2026-03-31-goreveal-strategic-review.md`

Exploratory note:
- `docs/plans/2026-03-22-goreveal-v2-ng-plan.md` remains exploratory only and is not the active roadmap baseline
- `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md` is a superseded
  historical sequence; RT1 is the only active execution authority

## Historical Suggested Backlog Epics

These ideas require re-prioritization through RT1 and do not select current
work.

Historical candidate epics:

- `Epic: Workstation Handoff Contract`
- `Epic: Review Workflow Actionability`
- `Epic: Thin Semantic and Source Confidence`
- `Epic: Protected Commercial Go Workflows`
- `Epic: Go Version Tracking and Markup Transfer`
- `Epic: Go Metadata Knowledge Network`
- `Epic: Thin JEB / Binary Ninja / Rizin Integration`

Implementation note for code peeling:
- start with function classification from existing `import_path`, `module_local`, `build_info.path`, and `source_file` truth
- add bounded stdlib/runtime fingerprints later
- prefer a compact hash such as `xxhash` for first-pass function fingerprints
- the current landed first slice is exactly that bounded classification layer in `engine`, exposed through `analysis.peeling`, `inspect peeling`, `goreveal peel <binary>`, and thin exports, and it now also includes package-level summaries plus per-class function counts

Implementation note for version tracking:
- do not start before code peeling MVP
- first version should be function-level diffing with similarity scoring over canonical function surfaces rather than a decompiler-dependent workflow
- that first bounded foothold is now started through `diff sqlite matched_functions`, using exact-name, source-location, source-file, and module-local normalized-name matching with `score` and `reason`

## Sources

Primary references used for this brainstorming:

- Hex-Rays docs and product pages:
  - https://docs.hex-rays.com/user-guide
  - https://hex-rays.com/pricing
  - https://hex-rays.com/lumina
  - https://hex-rays.com/teams
  - https://docs.hex-rays.com/user-guide/teams
  - https://hex-rays.com/blog/ida-9.3-teams-moves-inside-ida
  - https://docs.hex-rays.com/9.1/release-notes/9_0
  - https://docs.hex-rays.com/9.1/user-guide/plugins

- Ghidra official/official-doc mirrors:
  - https://www.nsa.gov/ghidra
  - https://www.ghidradocs.com/11.0.1_PUBLIC/help/Debugger/help/topics/Debugger/Debugger.html
  - https://www.ghidradocs.com/10.3_PUBLIC/help/Base/help/topics/VersionControl/project_repository.htm
  - https://www.ghidradocs.com/11.0_PUBLIC/help/VersionTracking/help/topics/VersionTrackingPlugin/VT_Wizard.html
  - https://www.ghidradocs.com/11.0.2_PUBLIC/help/VersionTracking/help/topics/VersionTrackingPlugin/Version_Tracking_Intro.html
  - https://www.ghidradocs.com/11.1_PUBLIC/help/Base/help/topics/GhidraServer/GhidraServer.htm

- PNF Software / JEB:
  - https://pnfsoftware.com/jeb/
  - https://www.pnfsoftware.com/jeb/android
  - https://www.pnfsoftware.com/jeb/intel
  - https://www.pnfsoftware.com/jeb/arm
  - https://www.pnfsoftware.com/jeb/riscv
  - https://www.pnfsoftware.com/jeb/manual/dev/writing-client-scripts/

- Binary Ninja:
  - https://binary.ninja/
  - https://binary.ninja/features/
  - https://binary.ninja/enterprise/
  - https://docs.enterprise.binary.ninja/
  - https://docs.binary.ninja/dev/bnil-overview.html
  - https://docs.binary.ninja/guide/warp.html
  - https://binary.ninja/2025/07/24/5.1-helion.html
