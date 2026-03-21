# GoREveal IDA Marketplace Notes

> Status: research note
> Date: 2026-03-19

## Purpose

This note captures useful ideas from the current Hex-Rays plugin ecosystem for Go-related analysis, so `GoREveal` can borrow proven UX/import patterns without copying code or inheriting external architecture.

## Sources

- Hex-Rays plugin ecosystem / publication docs:
  - https://docs.hex-rays.com/user-guide/plugins/plugin-submission-guide
- Official Hex-Rays plugin repository index:
  - https://raw.githubusercontent.com/HexRaysSA/plugin-repository/v1/plugin-repository.json
- Official Hex-Rays tags index:
  - https://raw.githubusercontent.com/HexRaysSA/plugin-repository/v1/tags.json
- GoResolver plugin README:
  - https://raw.githubusercontent.com/volexity/GoResolver/main/Plugin/GoResolver/README.md

## Go-Related Marketplace Signals

The strongest directly relevant marketplace signals are:

- `GoResolver`
  - Present in the official Hex-Rays plugin repository.
  - Positioned as an IDA/Ghidra plugin that can either analyze the current file or import a previously generated report.
  - Useful ideas:
    - import-first workflow is valid and user-friendly
    - external toolchain management should stay outside the core plugin logic
    - a shared import layer for IDA and Ghidra is a good pattern

- `golang_loader_assist`
  - Appears in the official Hex-Rays tags index as a `favourite`.
  - Useful ideas:
    - Go-specific import workflows are explicitly valued in the Hex-Rays ecosystem
    - function discovery and renaming remain a core analyst need
    - historically important, but should be treated as methodology reference, not architectural template

## Useful Adjacent Marketplace Signals

- `funcfiletree`
  - Not Go-specific, but tagged with `go` in plugin metadata.
  - Useful idea:
    - source-file grouping is valuable enough to exist as a standalone plugin, which validates `GoREveal` package/source projection as a first-class import surface

- `HashDB`
  - Not Go-specific.
  - Useful idea:
    - enrichment plugins that consume already discovered data are valuable, but should stay downstream of core recovery

## What Is Most Useful For GoREveal

The most useful marketplace-derived patterns are:

1. Thin import adapters are the right architecture.
   - `GoReSym`, `GoResolver`, and other heavy analysis should stay outside IDA/Ghidra plugin code.
   - Plugins should ingest stable payloads, not reproduce recovery logic.

2. Importing previously generated results is a first-class workflow.
   - `GoREveal` should keep `export ida` and `export ghidra` stable and consumable offline.

3. Package/source grouping is a valuable UI feature, not “nice to have”.
   - The plugin layer should help analysts navigate user code quickly.

4. Function renaming remains the baseline expectation.
   - Any thin adapter should first do symbols/functions well before richer type UI work.

5. Shared adapter behavior across IDA and Ghidra is worth preserving.
   - Differences should be in apply-hooks, not in contract semantics.

## What We Should Not Copy

- no plugin-side recovery logic
- no plugin-side Go version management in the first implementation
- no monolithic “do everything in the plugin” workflow
- no dependency on marketplace-only packaging before the adapter semantics are stable

## Immediate Product Decisions

Based on these ecosystem signals, `GoREveal` should:

- keep `IDA` and `Ghidra` adapters thin and contract-driven
- prioritize function import, package grouping, and source-file grouping
- treat report import as the canonical plugin workflow
- defer richer plugin UX until export contracts and differential evidence are stronger
