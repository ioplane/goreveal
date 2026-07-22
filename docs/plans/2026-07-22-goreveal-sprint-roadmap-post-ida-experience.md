---
status: draft
date: 2026-07-22
owners:
  - maintainers
---

# Plan: goreveal improvements after IDA Pro RE experience

## Context

goreveal was used to analyse a 410 MB stripped Go 1.25 binary
(Teleport Enterprise 18.10.0). The analysis produced 458,600 functions
from pclntab, compared to IDA Pro 9.4's 248,583 functions (54 %) with
zero Go function names in batch mode.

The IDA Golang plugin was confirmed NOT to run in default batch mode
(`idat -A -B`). A forced build with `-Ogolang:force:force_regabi` is
underway to measure the actual baseline.

This document records the gaps discovered during the IDA workflow and
proposes sprints to close them.

## Observed gaps

### G1: No binary identity in export

`schema/export_ida.go` produces `IDAExport` with `input.path`,
`input.size`, `input.format` but no hash, build ID, or image base.
Without identity binding, an idacli consumer cannot verify that the
artifact describes the same binary as the loaded IDB.

**Current:**
```go
type Input struct {
    Path   string `json:"path"`
    Size   int64  `json:"size"`
    Format string `json:"format"`
}
```

**Needed:** `sha256`, `build_id`, `image_base`, `va_semantics`.

### G2: No prologue classification

goreveal recovers function entry and end from pclntab but does not
classify the function prologue. The IDA workflow showed that Go
functions with goroutine stack checks (`cmp rsp, [r14+10h]; jbe`)
are not recognised by IDA's heuristics. goreveal could flag these
prologue types in the export so downstream tools know which functions
need special handling.

**Current:** `Function` has no prologue field.

**Needed:** `go_prologue` field with values: `"stack_check"`,
`"standard"`, `"thunk"`, `"unknown"`.

### G3: No string length

`schema.StringCandidate` has `Value`, `Offset`, `Addr` but no length.
The idacli `create_strlit` call requires both address and length.
Without length, string import is unsafe.

**Current:**
```go
type StringCandidate struct {
    Value      string     `json:"value"`
    Region     string     `json:"region"`
    Addr       uint64     `json:"addr,omitempty"`
    Offset     uint64     `json:"offset"`
    Provenance Provenance `json:"provenance"`
}
```

**Needed:** `Length int` field.

### G4: No type layout

`schema.Type` has `Name`, `Package`, `Kind` but no struct fields,
sizes, or offsets. Without layout, applying Go types to IDA's type
library is impossible.

**Current:** Type has `Name`, `Kind`, `Package`, `ImportPath`,
`SourceFileCount`, `ModuleLocal`, `UserMeaningful`.

**Needed:** `Size uint64`, `Fields []Field` (for structs), with
`Field.Name`, `Field.Offset`, `Field.TypeName`, `Field.Size`.

### G5: No diff capability

goreveal can analyse two binaries but cannot compare the results.
Version comparison (e.g., 18.7.2 vs 18.10.0) requires manual diff of
JSON files. A built-in diff command would produce structured output:
functions added, removed, renamed, or with changed boundaries.

**Current:** No diff command.

**Needed:** `goreveal diff` command with SQLite-backed comparison.

### G6: No VA semantics

Function `Entry` and `End` are returned as absolute virtual addresses
from `debug/gosym`. However, the export does not state whether these
are absolute VAs, RVAs, or file offsets. For PIE binaries loaded at
different base addresses, this ambiguity can cause address mismatch
with IDA's IDB.

**Current:** No `va_semantics` field.

**Needed:** `runtime.va_semantics` with `"absolute"`, `"rva"`, or
`"file_offset"`.

### G7: Export v1 contract is not versioned for identity

`IDAExportContractV1` does not include identity fields. A v2
contract is needed with `input.sha256`, `input.build_id`,
`runtime.image_base`, and `artifact_sha256`.

### G8: No artifact self-digest

The JSON artifact has no SHA-256 of itself. idacli's `go-apply` task
needs to verify that the operator reviewed a specific preview before
applying. A self-digest enables digest-based apply authorization.

## Sprint plan

Methodology: vertical slices, each delivering a verifiable capability
that works on at least one real binary. Every slice includes tests or
evidence. Sprints are ordered by dependency: identity first, then
prologue, then the rest.

### Sprint A: Binary identity binding (G1, G6, G7, G8)

**Goal:** goreveal export includes verifiable binary identity.

| Task | Outcome | Evidence |
|---|---|---|
| A1 | Add `sha256` to `Input` — hash the binary file | Unit test: known binary → known SHA-256 |
| A2 | Add `build_id` to `BuildInfo` — extract from `debug/buildinfo` | Unit test: Go binary → non-empty build ID |
| A3 | Add `image_base` and `va_semantics` to `RuntimeMetadata` | Unit test: ELF binary → correct text section addr, `"absolute"` |
| A4 | Add `artifact_sha256` — self-digest of the JSON output | Unit test: stable digest for same input |
| A5 | Bump contract to `goreveal.export.ida/v2` | Schema test: v2 JSON has all identity fields |
| A6 | `goreveal export ida --v2` produces identity-bound JSON | Golden snapshot on fixture binary |

**Acceptance:** `goreveal export ida <binary>` output contains
`input.sha256`, `input.build_id`, `runtime.image_base`,
`runtime.va_semantics`, and `artifact_sha256`. All fields are
non-empty for a real Go binary. Existing v1 consumers continue to
work (v2 is additive).

### Sprint B: Function prologue classification (G2)

**Goal:** goreveal classifies Go function prologues.

| Task | Outcome | Evidence |
|---|---|---|
| B1 | Add `GoPrologue string` to `Function` schema | Schema test: field exists |
| B2 | Read first 8 bytes at each function entry | Unit test: bytes read correctly |
| B3 | Detect goroutine stack check: `49 3b 66 10` (cmp rsp, [r14+10h]) | Unit test: known Go binary → >80 % functions have `stack_check` |
| B4 | Detect thunk: function end - entry < 16 bytes and contains `JMP` or `RET` | Unit test: thunks classified |
| B5 | Classify as `"standard"` if prologue is `push rbp; mov rbp, rsp` or similar | Unit test |
| B6 | Default to `"unknown"` if no pattern matches | Unit test |
| B7 | Export prologue in `IDAFunction.GoPrologue` | Golden snapshot: prologue field populated |

**Acceptance:** `goreveal inspect functions <binary>` output includes
`go_prologue` for each function. At least 80 % of user functions in
a Go binary are classified as `"stack_check"`.

### Sprint C: String length recovery (G3)

**Goal:** goreveal exports string length for safe IDA import.

| Task | Outcome | Evidence |
|---|---|---|
| C1 | Add `Length int` to `StringCandidate` schema | Schema test |
| C2 | Compute length from string value (UTF-8 byte count) | Unit test: known string → correct length |
| C3 | Compute length from region boundary if value is truncated | Unit test: string at end of region |
| C4 | Export length in `IDAString` | Golden snapshot: length field populated |
| C5 | Verify: `create_strlit(addr, length)` parameters are safe | Integration test with fixture |

**Acceptance:** `goreveal inspect strings <binary>` output includes
`length` for each string. Length matches `len(value)` for ASCII
strings. For strings near region boundaries, length does not exceed
the region end.

### Sprint D: Type layout recovery (G4)

**Goal:** goreveal exports struct field layout for IDA type application.

| Task | Outcome | Evidence |
|---|---|---|
| D1 | Add `Size uint64` and `Fields []Field` to `Type` schema | Schema test |
| D2 | Define `Field` struct: `Name`, `Offset`, `TypeName`, `Size` | Schema test |
| D3 | Recover struct size from typelinks (if available in pclntab) | Unit test: known Go struct → correct size |
| D4 | Recover field offsets from type metadata | Unit test: known struct → correct field offsets |
| D5 | Export type layout in `IDAType` | Golden snapshot: struct types have fields |
| D6 | Handle non-struct types (interface, map, chan, func) — size only | Unit test |

**Acceptance:** `goreveal inspect types <binary>` output includes
`size` for all types. Struct types have `fields` with `name`,
`offset`, `type_name`, and `size`. Non-struct types have `size`
only. At least 10 user-defined struct types are fully recovered on
a real Go binary.

**Note:** This is the most complex sprint. Type layout recovery from
stripped Go binaries is an active research area. If pclntab does not
contain enough layout information, this sprint may be descoped or
reduced to size-only recovery.

### Sprint E: Version diff (G5)

**Goal:** goreveal can compare two analyses and report differences.

| Task | Outcome | Evidence |
|---|---|---|
| E1 | `goreveal analyze` persists to SQLite (`--sqlite output.db`) | Unit test: SQLite created with functions table |
| E2 | `goreveal diff sqlite <left.db> <right.db>` — compare function tables | Unit test: added/removed/renamed/changed-boundary functions |
| E3 | Diff output: JSON with `added[]`, `removed[]`, `renamed[]`, `boundary_changed[]` | Unit test: two known binaries → correct diff |
| E4 | Diff by address (same VA, different name = renamed) | Unit test |
| E5 | Diff by name (same name, different VA = moved) | Unit test |
| E6 | Diff report: summary counts + detailed lists | Golden snapshot: 18.7.2 vs 18.10.0 diff |

**Acceptance:** `goreveal diff sqlite <v18.7.2.db> <v18.10.0.db>`
produces a structured diff report showing functions added, removed,
renamed, and with changed boundaries. Diff is deterministic for the
same inputs.

## Sprint dependencies

```
Sprint A (identity) ──┐
                       ├──→ Sprint E (diff, needs SQLite + identity)
Sprint B (prologue) ──┤
                       │
Sprint C (strings) ────┤
                       │
Sprint D (types) ──────┘ (most complex, may slip)
```

Sprints A, B, C can run in parallel. Sprint D is the most complex
and may be descoped. Sprint E depends on A (identity for comparison).

## IDA Golang plugin baseline

A baseline measurement is running: `idat -A -B -Ogolang:force:force_regabi`
on the 410 MB binary. When complete, record:

| Metric | Without Golang plugin | With Golang plugin (forced) | goreveal |
|---|---|---|---|
| Function count | 248,583 | TBD | 458,600 |
| Go-named functions | 0 | TBD | 458,600 |
| Hex-Rays decompile success | ~20 % | TBD | N/A |

This baseline determines the actual delta that idacli's `go-apply`
task needs to recover. If the Golang plugin recovers 400K functions,
the goreveal export is still valuable for the remaining 58K, plus
for name recovery, prologue classification, and diff capability.

## Methodology

This plan follows the goreveal AGENTS.md delivery rules:

- Vertical slices: each sprint delivers a usable capability
- Evidence: every task has a test, snapshot, or benchmark
- Clean-room: no code copied from GoReSym, gore, or other tools
- Provenance: all new fields carry `Provenance` with source and confidence
- Schema-first: schema changes precede implementation
- Core independence: `core` stays independent from CLI and plugins
- Podman-first: all development in dev container

## Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Golang plugin recovers most functions, reducing goreveal value | Low | goreveal still provides names, prologues, diff, identity — not just function count |
| Type layout not available in pclntab (stripped binary) | High | Sprint D may be descoped to size-only; full layout needs DWARF or runtime introspection |
| SHA-256 of 410 MB binary takes too long | Low | Go crypto/sha256 benchmarks at >1 GB/s; 410 MB < 1 second |
| Prologue detection false positives | Medium | Validate against GoReSym output as differential test |
| Diff on 458K functions is slow | Low | SQLite indexed by entry address; diff is O(n log n) |
