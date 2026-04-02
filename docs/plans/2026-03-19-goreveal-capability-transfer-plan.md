# GoREveal Capability Transfer Plan

> Status: working product-transfer plan
> Date: 2026-03-19
> Purpose: map each baseline Go RE utility to the specific capabilities that GoREveal already covers, will absorb natively, or will intentionally keep as external reference/orchestration.

## Core Principle

The long-term product direction is not to remain a shell around external tools.
The external tools are:
- behavior references
- differential baselines
- fixture sources
- temporary orchestration targets

The product goal is to absorb their practically useful capability into native `GoREveal` modules while preserving a clean-room boundary.

## Utility Matrix

### `gore`

Role today:
- library-oriented baseline for package, function, type, and source-file behavior

Already covered natively:
- build info recovery
- package presence
- source-file basename overlap
- narrow user-function overlap
- canonical schema and CLI instead of ad hoc library-shaped output
- package `import_path` and `source_file_count` enrichment from source-tree correlation
- package `module_local` scope metadata derived from source-tree correlation
- raw type metadata enriched with `package` and `user_meaningful` hints instead of destructive filtering
- raw type metadata enriched with `import_path` and `module_local` to match the richer package-navigation surface
- raw type metadata enriched with `source_file_count` so current source-backed scope heuristics are visible rather than implicit

Planned native transfers:
- broader package metadata recovery
- richer source-file to package correlation
- more usable user-type surface
- more complete coverage across version families
- less heuristic package metadata from independent metadata sources where possible
- stronger package-scope truth in binaries where source-file evidence is weak or absent
- stronger user-type evidence from typelinks/moduledata rather than `DWARF`-only heuristics
- less heuristic type-scope truth when source-file correlation is weak or absent

Intentionally not copied literally:
- `gore` library API shape as a compatibility target
- `gore` type output as-is, because current raw surface is too runtime-heavy for trustworthy parity claims

Sprint mapping:
- `Sprint 7`: keep `gore` as differential evidence for safe overlap surfaces
- `Future Sprint 11`: absorb broader package/type metadata natively

### `redress`

Role today:
- CLI-oriented baseline for source projection, package listing, and presentation patterns

Already covered natively:
- package recovery
- initial source-tree projection
- module path parity
- narrow source-file and function overlap
- grouped source-package nodes with multiple files per import path
- grouped source-package nodes enriched with package-level function counts
- explicit external-package projection for non-module source paths such as runtime/stdlib files

Planned native transfers:
- deeper source reconstruction
- stronger package/file tree projection
- more useful source-like grouping for reverse engineers
- richer projection metadata in schema/export layers
- better normalization for non-module and dependency source layouts beyond the current external-package heuristics

Intentionally not copied literally:
- `redress` CLI text presentation as a product target
- full “pretty CLI” behavior before schema/source reconstruction is deeper

Sprint mapping:
- `Sprint 7`: continue using `redress` as evidence for source projection claims
- `Future Sprint 11`: native `redress`-style source reconstruction depth

### `GoReSym`

Role today:
- runtime-aware truth baseline for `pclntab`, file mapping, and Go metadata recovery

Already covered natively:
- build info parity
- read-only ELF runtime metadata for `firstmoduledata`, `.gopclntab`, `.typelink`, and `.go.module`
- raw `typelink_count` evidence from the ELF `.typelink` section
- bounded raw `typelink_sample` evidence from the ELF `.typelink` section
- raw `typelink_min_offset` and `typelink_max_offset` evidence from the ELF `.typelink` section
- raw `typelink_negative_count` and `typelink_non_negative_count` evidence from the ELF `.typelink` section
- bounded `firstmoduledata`/`.go.module` cross-check plus raw `.go.module` word sample
- bounded `moduledata` typelinks slice-header cross-check against `.typelink` section evidence
- `.itablink` section evidence plus a bounded `moduledata` itablinks slice-header cross-check
- a bounded `moduledata` memory-range block cross-check against ELF `.noptrdata`, `.data`, `.bss`, and `.noptrbss` boundaries
- a bounded `.rodata` range cross-check against ELF section boundaries
- a bounded `.text` range cross-check against ELF section boundaries, using inclusive-end semantics
- the first bounded semantic typelink bridge via rodata-relative resolution on the canonical fixture
- a bounded semantic confidence bit showing that all current fixture typelinks resolve within `.rodata`
- module-local file overlap
- narrow user-function overlap
- direct `.gopclntab`-based function recovery for current `ELF` fixture

Planned native transfers:
- richer runtime metadata extraction
- stronger `pclntab` coverage across Go versions
- `moduledata`-driven recovery
- `typelinks`-driven type recovery
- better runtime-aware file and package evidence

Intentionally not copied literally:
- internal implementation details of `GoReSym`
- broad parity claims before multi-version evidence exists

Sprint mapping:
- `Sprint 7`: keep `GoReSym` as truth/evidence baseline
- `Future Sprint 12`: transfer richer runtime/type recovery natively

### `GoResolver`

Role today:
- deobfuscation and CFG-guided symbol-recovery reference

Already covered natively:
- nothing substantial yet beyond the product decision to keep refined output separate from raw truth

Planned native transfers:
- selective symbol/name refinement informed by stable semantics
- bounded deobfuscation passes that fit the schema/refined-layer model
- possibly orchestrated comparison or optional execution path

Intentionally not copied literally:
- full CFG-similarity engine in the early product line
- plugin-side or external-tool-side heavy analysis logic inside `GoREveal core`

Sprint mapping:
- `Future Sprint 13`: evaluate narrow overlap/orchestration first
- only later promote stable, bounded parts into native refinement passes

### `gostringungarbler`

Role today:
- string refinement and garble-related deobfuscation reference

Already covered natively:
- only very early refined string handling and string-segment extraction

Planned native transfers:
- bounded string refinement passes
- garble-aware heuristics where they can be expressed cleanly in the refined layer
- optional orchestration for comparison while the native passes are immature

Intentionally not copied literally:
- full external deobfuscator behavior before native confidence/provenance rules are preserved

Sprint mapping:
- `Future Sprint 13`: narrow overlap/orchestration and selective native transfer

### `AlphaGolang`

Role today:
- historical IDA integration and annotation/reference baseline

Already covered natively:
- thin IDA adapter
- thin Ghidra adapter
- stable export contracts

Planned native transfers:
- better import UX
- better symbol/package/type application flow in adapters
- possibly richer operator workflows around imports

Intentionally not copied literally:
- IDA-side Go analysis engine
- plugin-side recovery logic
- SDK-heavy native plugin work without concrete need

Sprint mapping:
- `Sprint 8`: already covers the first practical transfer wave
- later work should stay thin and contract-driven

### `Rizin`

Role today:
- deferred headless host-platform adapter opportunity for JSON payload import and automation workflows

Already covered natively:
- nothing directly; current host-platform transfer work is centered on `IDA` and `Ghidra`

Planned native transfers:
- thin `export rizin` payloads later if the same canonical export-first contract remains intact
- bounded headless workflows where `Rizin` is a consumer of canonical truth, not a recovery dependency

Intentionally not copied literally:
- `Rizin` analysis logic inside `core`
- any adapter-side recovery logic that duplicates canonical schema work

Sprint mapping:
- deferred backlog only
- keep below `JEB` and `Binary Ninja` in current host-platform priority, even if a future adapter may be technically simpler, because current market pull is weaker

## Adjacent Workstation Signals

Recent `rehelp` and RE-lab inventory work sharpens the surrounding transfer picture without changing the clean-room boundary.

Measured workstation adjacencies:
- host platforms:
  - `ida-pro`
  - `ghidra`
  - `jeb`
  - `rizin`
- diffing / metadata-transfer tools:
  - `diaphora`
  - `binexport`
  - `binsync`
- host-platform MCP signal:
  - `ida-pro-mcp`
- dynamic / symbolic sidecars:
  - `frida`
  - `angr`
  - `qiling`
  - `unicorn`
  - `uftrace`
  - `z3`

What this means for transfer planning:
- future function-level version tracking should be informed by `Diaphora` / `BinExport` as external-reference and workflow inputs, not by copying their implementations
- host-platform MCP interop is now a real near-term planning target, not only a generic architecture note
- a thin future `Rizin` adapter is technically more credible because the operator environment already has `rizin` plus related plugins, even if product priority still stays below `JEB` and `Binary Ninja`
- protected/deobfuscation work now has a realistic external orchestration environment if later measured needs justify it

See also:
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md`

## Prioritized Transfer Order

The most rational transfer order from the current codebase is now:

1. `GoReSym`-style runtime/type recovery depth
2. `gore`-style broader package/type metadata quality
3. `redress`-style deeper source reconstruction
4. bounded `gostringungarbler`-style string refinement
5. bounded `GoResolver`-style deobfuscation transfer

This order follows the project priority:
- accuracy first
- convenience second
- speed third

## Sprint Translation

### Near-Term

- `Sprint 7`
  - keep evidence healthy and divergences documented
  - treat it as a maintenance lane, not the main product lane

- `Sprint 8`
  - already established thin integration boundary

### Next Native Capability Sprints

- `Sprint 11`
  - native transfer from `redress` and `gore`
  - focus: source reconstruction depth, richer package/file metadata

- `Sprint 12`
  - native transfer from `GoReSym`
  - focus: runtime metadata, `moduledata`, `typelinks`, and the first minimal semantic decode from the accumulated bounded evidence

- `Sprint 13`
  - selective native transfer from `gostringungarbler` and `GoResolver`
  - focus: bounded deobfuscation, orchestration where native parity is not yet justified
  - do not start before the bounded `PE` fixture checkpoint and the first code-peeling MVP are complete

## What Counts As Success

For each utility, success is not:
- “we wrapped it”
- “we can call it”
- “we copied its output shape”

Success is:
- `GoREveal` covers the user-relevant capability natively
- the capability is represented in canonical schema
- the capability is validated against corpus and baselines
- the external tool becomes optional for the corresponding workflow
