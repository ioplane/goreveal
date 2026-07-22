---
status: research-input
date: 2026-07-22
owners:
  - maintainers
superseded_by:
  - ../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md
  - ../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md
---

# Proposal: goreveal improvements after IDA Pro RE field experience

> This field report is an RT1 research input, not an implementation queue.
> Its measured claims remain subject to the forced-plugin baseline defined in
> RT1-S1.

## Summary

goreveal was used to analyse a 410 MB stripped Go 1.25 binary as the
primary function recovery source for an IDA Pro 9.4 reverse-engineering
workflow. goreveal recovered 458,600 functions from pclntab; IDA's
auto-analysis recovered 248,583 (54 %) with zero Go function names in
default batch mode. Eight capability gaps were identified. This document
proposes five sprints to close them.

## What happened

A reverse-engineering project required decompiling specific functions
in a large stripped Go binary. The workflow was:

1. `idat -A -B -o binary.i64 binary` — IDA auto-analysis (1 h 34 m)
2. `goreveal inspect functions binary` — pclntab function recovery
3. Hex-Rays decompilation of target addresses via idalib

Results:

| Metric | IDA auto-analysis | goreveal (pclntab) |
|---|---|---|
| Functions recovered | 248,583 (54 %) | 458,600 (100 %) |
| Go function names | 0 (Golang plugin did not run in batch) | 458,600 |
| Hex-Rays decompilation success on 9 target functions | 2 of 9 (22 %) | N/A |
| Function boundary accuracy | Heuristic (frequently wrong for Go) | Authoritative (pclntab entry/end) |

Three decompilation failure modes were observed:

1. IDA split a large Go function (49 callees, 14 closures, ~20 KB
   expected pseudocode) into 140-byte fragments because the Go
   prologue was not recognised.
2. IDA classified a function entry address as data, preventing function
   creation.
3. IDA placed a function entry in the middle of an adjacent function
   because the preceding function's end was not detected.

goreveal's pclntab data correctly identified all function entries, ends,
and names, but the IDA database could not consume this data — there was
no bridge from goreveal's output to IDA's function table.

A separate proposal has been filed in the idacli repository
(`docs/planning/2026-07-22-go-function-recovery-task.md`) for a
`go-preview` / `go-apply` task that would consume goreveal's export to
correct IDA's function table. That proposal depends on goreveal adding
the capabilities described here.

## IDA Golang plugin baseline

IDA Pro 9.4 ships a Golang plugin that can extract functions from
pclntab, recover types from typelinks, and model the Go stack/register
ABI. It can be forced via `-Ogolang:force:force_regabi`.

In default batch mode (`idat -A -B`), the Golang plugin **did not
run** — zero Go function names were recovered. A forced build with
`-Ogolang:force:force_regabi` is underway to measure the actual
baseline. goreveal's value proposition depends on the delta between
the forced Golang plugin and pclntab truth:

- If the forced plugin recovers 400K of 458K functions, goreveal's
  primary value shifts to names, prologue classification, diff, and
  identity binding rather than raw function count.
- If the forced plugin recovers fewer (e.g., 300K), goreveal remains
  the primary function recovery tool.

Either way, goreveal's pclntab-based recovery is authoritative and
complements IDA's analysis.

## Identified gaps

### G1: No binary identity in export

`schema/export_ida.go` produces `IDAExport` with `input.path`,
`input.size`, `input.format` but no hash, build ID, or image base.
Without identity binding, a downstream consumer (e.g., idacli) cannot
verify that the artifact describes the same binary as the loaded IDB.

### G2: No function prologue classification

goreveal recovers function entry and end from pclntab but does not
classify the prologue. Go functions with goroutine stack checks
(`cmp rsp, [r14+10h]; jbe morestack`) are not recognised by IDA's
heuristics. goreveal could flag these so downstream tools know which
functions need special handling.

### G3: No string length

`schema.StringCandidate` has `Value`, `Offset`, `Addr` but no length.
IDA's `create_strlit(addr, length)` requires both address and length.
Without length, string import is unsafe.

### G4: No type layout

`schema.Type` has `Name`, `Package`, `Kind` but no struct fields,
sizes, or offsets. Without layout, applying Go types to IDA's type
library is impossible.

### G5: No diff capability

goreveal can analyse two binaries but cannot compare the results.
Version comparison requires manual JSON diff. A built-in diff command
would produce structured output: functions added, removed, renamed,
or with changed boundaries.

### G6: No VA semantics

Function `Entry` and `End` are absolute virtual addresses from
`debug/gosym`, but the export does not state this. For PIE binaries
loaded at different base addresses, this ambiguity can cause address
mismatch with IDA's IDB.

### G7: Export v1 contract is not versioned for identity

`IDAExportContractV1` does not include identity fields. A v2 contract
is needed with `input.sha256`, `input.build_id`,
`runtime.image_base`, and `artifact_sha256`.

### G8: No artifact self-digest

The JSON artifact has no SHA-256 of itself. A downstream consumer
needs to verify that an operator reviewed a specific preview before
applying changes. A self-digest enables digest-based authorization.

## Proposed sprints

Methodology: vertical slices per AGENTS.md. Each sprint delivers a
verifiable capability that works on at least one real binary. Every
task includes tests, golden snapshots, or differential evidence.
Sprints are ordered by dependency: identity first, then prologue,
then the rest.

### Sprint A: Binary identity binding (G1, G6, G7, G8)

Goal: goreveal export includes verifiable binary identity.

Tasks:

1. Add `sha256` to `Input` — hash the binary file via `crypto/sha256`.
2. Add `build_id` to `BuildInfo` — extract from `debug/buildinfo`
   (already partially implemented in `core/buildinfo`).
3. Add `image_base` and `va_semantics` to `RuntimeMetadata` —
   `image_base` from ELF `.text` section addr, `va_semantics` as
   `"absolute"` for ELF (current behaviour), `"rva"` for PE.
4. Add `artifact_sha256` — self-digest of the JSON output (SHA-256 of
   the serialized `IDAExport` struct).
5. Bump contract to `goreveal.export.ida/v2` — v2 is additive over v1;
   v1 consumers continue to work.
6. `goreveal export ida <binary>` produces identity-bound JSON with
   all new fields populated.

Acceptance: `goreveal export ida <binary>` output contains
`input.sha256`, `input.build_id`, `runtime.image_base`,
`runtime.va_semantics`, and `artifact_sha256`. All fields are
non-empty for a real Go binary. Existing v1 consumers are not broken.

### Sprint B: Function prologue classification (G2)

Goal: goreveal classifies Go function prologues.

Tasks:

1. Add `GoPrologue string` to `Function` schema with values:
   `"stack_check"`, `"standard"`, `"thunk"`, `"closure"`, `"unknown"`.
2. Read first 16 bytes at each function entry (file offset from VA via
   ELF section mapping).
3. Detect goroutine stack check: bytes `49 3b 66 10` at entry
   (`cmp rsp, [r14+10h]`) — classify as `"stack_check"`.
4. Detect thunk: `end - entry < 24` and contains `JMP` (0xE9/0xEB)
   or `RET` (0xC3) — classify as `"thunk"`.
5. Detect closure: function name contains `.func` suffix — classify
   as `"closure"`.
6. Detect standard C prologue: `55 48 89 e5` (`push rbp; mov rbp, rsp`)
   — classify as `"standard"`.
7. Default to `"unknown"` if no pattern matches.
8. Export prologue in `IDAFunction.GoPrologue`.

Acceptance: `goreveal inspect functions <binary>` output includes
`go_prologue` for each function. At least 80 % of non-runtime user
functions in a Go binary are classified as `"stack_check"`. Thunks
and closures are correctly identified. Differential test against
GoReSym output validates prologue counts.

### Sprint C: String length recovery (G3)

Goal: goreveal exports string length for safe downstream import.

Tasks:

1. Add `Length int` to `StringCandidate` schema.
2. Compute length from string value: `len([]byte(value))` for the
   UTF-8 byte count.
3. For strings near region boundaries, clamp length to not exceed the
   region end offset.
4. Export length in `IDAString.Length`.
5. Verify: `create_strlit(addr, length)` parameters are safe — length
   must be > 0 and not exceed region boundary.

Acceptance: `goreveal inspect strings <binary>` output includes
`length` for each string. Length matches `len(value)` for ASCII
strings. For strings at region boundaries, length does not exceed the
region end. Zero-length strings are filtered out.

### Sprint D: Type layout recovery (G4)

Goal: goreveal exports struct field layout for type application.

Tasks:

1. Add `Size uint64` and `Fields []Field` to `Type` schema.
2. Define `Field` struct: `Name string`, `Offset uint64`,
   `TypeName string`, `Size uint64`.
3. Recover struct size from typelinks if available.
4. Recover field offsets from type metadata (research needed —
   pclntab typelinks may not contain full layout for stripped
   binaries).
5. Export type layout in `IDAType`.
6. Handle non-struct types (interface, map, chan, func) — size only,
   no fields.

Acceptance: `goreveal inspect types <binary>` output includes `size`
for all types. Struct types have `fields` with `name`, `offset`,
`type_name`, and `size`. At least 10 user-defined struct types are
fully recovered on a real Go binary.

Note: this is the most complex sprint. Type layout recovery from
stripped Go binaries is an active research area. If pclntab does not
contain enough layout information, this sprint may be descoped to
size-only recovery or deferred until DWARF or runtime introspection
is available.

### Sprint E: Version diff (G5)

Goal: goreveal can compare two analyses and report differences.

Tasks:

1. `goreveal analyze` persists results to SQLite (`--sqlite output.db`)
   with a `functions` table (entry, end, name, package) and
   `types` table.
2. `goreveal diff sqlite <left.db> <right.db>` compares function
   tables by entry address.
3. Diff output: JSON with `added[]`, `removed[]`, `renamed[]` (same
   VA, different name), `boundary_changed[]` (same VA, different
   end), `moved[]` (same name, different VA).
4. Diff report: summary counts + detailed lists, sorted by VA.
5. Deterministic output for same inputs.

Acceptance: `goreveal diff sqlite <v1.db> <v2.db>` produces a
structured diff report. Diff is deterministic: same inputs always
produce the same output. At least 1,000 changed functions are
correctly classified on a real version comparison.

## Sprint dependencies

```
Sprint A (identity) ──┐
                       ├──→ Sprint E (diff, needs identity for comparison)
Sprint B (prologue) ──┤
                       │
Sprint C (strings) ────┤
                       │
Sprint D (types) ──────┘ (most complex, may be descoped)
```

Sprints A, B, C are independent and can run in parallel. Sprint D is
the most complex and may be descoped. Sprint E depends on A (identity
binding for comparison).

## Downstream consumer

The idacli repository has a proposal (`docs/planning/2026-07-22-go-
function-recovery-task.md`, revision 2) for a `go-preview` / `go-apply`
task that consumes goreveal's export to correct IDA's function table.
That proposal depends on:

- Sprint A: identity fields for binary verification
- Sprint B: prologue classification for stack-check marking
- Sprint C: string length for safe `create_strlit`
- Sprint D: type layout for type application (deferred in idacli v1)

The idacli proposal uses a provider-preview / provider-apply contract
where goreveal is the provider, idacli is the consumer. goreveal
remains a standalone Go binary; idacli does not embed Go runtime or
use IDAPython.

## Baseline measurement

The IDA Golang plugin baseline build with
`-Ogolang:force:force_regabi` is running on oel-lab-gui. When
complete, record:

| Metric | Without Golang plugin | With Golang plugin (forced) | goreveal |
|---|---|---|---|
| Function count | 248,583 | TBD | 458,600 |
| Go-named functions | 0 | TBD | 458,600 |
| Decompilation success rate | ~20 % | TBD | N/A |

This baseline determines the actual delta that idacli needs to
recover. goreveal's value is not solely function count — it also
provides names, prologue classification, string length, type layout,
diff capability, and identity binding that IDA's Golang plugin does
not offer.

## Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Golang plugin recovers most functions | Low | goreveal still provides names, prologues, diff, identity |
| Type layout not in pclntab (stripped) | High | Sprint D may be descoped to size-only |
| Prologue detection false positives | Medium | Differential test against GoReSym |
| Diff on 458K functions is slow | Low | SQLite indexed by entry; O(n log n) |
| SHA-256 of 410 MB binary slow | Low | Go crypto/sha256 >1 GB/s; <1 second |

## Methodology compliance

This proposal follows goreveal AGENTS.md rules:

- Vertical slices: each sprint delivers a usable capability
- Evidence: every task has a test, snapshot, or benchmark
- Clean-room: no code copied from GoReSym, gore, or other tools
- Provenance: all new fields carry `Provenance` with source and
  confidence
- Schema-first: schema changes precede implementation
- Core independence: `core` stays independent from CLI and plugins
- Podman-first: all development in dev container
