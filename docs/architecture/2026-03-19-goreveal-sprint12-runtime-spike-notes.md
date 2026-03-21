# GoREveal Sprint 12 Runtime Spike Notes

> Status: initial spike note
> Date: 2026-03-19

## Purpose

Record the first bounded `Sprint 12` findings before any broader runtime or typelink refactor.

## Canonical Fixture Findings

For `corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin`, the current Podman-first probe confirmed:
- `runtime.firstmoduledata` is present in the ELF symbol table
- `.gopclntab` is present as a dedicated ELF section
- `.typelink` is present as a dedicated ELF section
- `.go.module` is present as a dedicated ELF section

Observed evidence from the fixture:
- `runtime.firstmoduledata` at `0x5781a0`
- `.gopclntab` at `0x4d5388`
- `.typelink` at `0x576a80`
- `.go.module` at `0x5781a0`

## What Was Implemented

The first `Sprint 12` slice is intentionally read-only:
- `core/runtime.ReadMetadata()` now extracts:
  - `firstmoduledata_addr`
  - `gopclntab_addr`
  - `gopclntab_size`
  - `typelink_addr`
  - `typelink_size`
  - `typelink_count`
  - `go_module_addr`
  - `go_module_size`
- the data is exposed in canonical schema as `analysis.runtime`
- the current provenance is:
  - `source = core.runtime.elf`
  - `confidence = medium`

This is not yet a `moduledata` parser and does not claim typelink traversal.

The second bounded slice is also now landed:
- raw `typelink_count` is exposed from `.typelink` size as a read-only evidence field
- for the canonical fixture this yields a non-zero typelink count and confirms that the section is structurally usable as the next runtime-truth bridge

The third bounded slice is now landed as well:
- a bounded `typelink_sample` is exposed as raw 32-bit offsets from `.typelink`
- this gives the project its first concrete typelink contents without crossing into type decoding yet

The fourth bounded slice is now landed too:
- raw `typelink_min_offset` and `typelink_max_offset` are exposed
- this gives the project its first shape validation over typelink contents without claiming semantic decoding

The fifth bounded slice is now landed as well:
- raw `typelink_negative_count` and `typelink_non_negative_count` are exposed
- this gives a whole-section sign-distribution check for typelink offsets without semantic decoding

The sixth bounded slice is now landed:
- `firstmoduledata_in_go_module` and `firstmoduledata_go_module_offset` are exposed
- `go_module_word_size` and a bounded `go_module_word_sample` are exposed
- this is the first real bridge between raw typelink evidence and a minimal `moduledata`/`.go.module` consistency check

The seventh bounded slice is now landed:
- a minimal `moduledata` typelinks slice-header cross-check is exposed
- `moduledata_typelink_slice_word_index`, `moduledata_typelink_len`, and `moduledata_typelink_cap` are exposed
- for the canonical fixture this confirms that `firstmoduledata` contains a typelinks slice header whose pointer and length/capacity agree with the existing `.typelink` section evidence

The eighth bounded slice is now landed:
- `.itablink` section evidence is exposed through `itablink_addr`, `itablink_size`, and `itablink_count`
- a minimal `moduledata` itablinks slice-header cross-check is exposed through `moduledata_itablink_slice_word_index`, `moduledata_itablink_len`, and `moduledata_itablink_cap`
- for the canonical fixture this confirms that `firstmoduledata` also carries an `itablinks` slice header whose pointer and length/capacity agree with the existing `.itablink` section evidence

The ninth bounded slice is now landed:
- a bounded `moduledata` memory-range block cross-check is exposed through:
  - `moduledata_memory_ranges_word_index`
  - `moduledata_noptrdata_addr` / `moduledata_noptrdata_end`
  - `moduledata_data_addr` / `moduledata_data_end`
  - `moduledata_bss_addr` / `moduledata_bss_end`
  - `moduledata_noptrbss_addr` / `moduledata_noptrbss_end`
- for the canonical fixture this confirms that `firstmoduledata` carries the expected contiguous range block for `.noptrdata`, `.data`, `.bss`, and `.noptrbss`, and that those ranges agree with the ELF section boundaries

The tenth bounded slice is now landed:
- a bounded `.rodata` range cross-check is exposed through:
  - `moduledata_rodata_word_index`
  - `moduledata_rodata_addr`
  - `moduledata_rodata_end`
- for the canonical fixture this confirms that `firstmoduledata` carries a `.rodata` range pair whose start and end agree with the ELF `.rodata` section boundaries

The eleventh bounded slice is now landed:
- a bounded `.text` range cross-check is exposed through:
  - `moduledata_text_word_index`
  - `moduledata_text_addr`
  - `moduledata_text_end_inclusive`
- for the canonical fixture this confirms that `firstmoduledata` carries a `.text` range pair whose start and inclusive end agree with the ELF `.text` section boundaries

The twelfth slice is the first very small semantic decode:
- `typelink_resolved_base_addr`, `typelink_resolved_sample`, and `typelink_resolved_within_rodata_count` are exposed
- on the canonical fixture, this models typelinks as bounded positive offsets resolved relative to the current `.rodata` base
- this is intentionally still fixture-local and does not yet claim that `.rodata` is the universal typelink base across versions or formats

The thirteenth slice strengthens that first semantic decode:
- `typelink_all_resolved_within_rodata` is exposed
- on the canonical fixture, this currently evaluates to true (`514/514`)
- this does not make the rodata-relative model universal, but it materially raises confidence that the current fixture-local semantic hypothesis is self-consistent

The fourteenth slice extends that semantic bridge without broad parser expansion:
- `moduledata_types_addr` and `moduledata_etypes_addr` are exposed
- `moduledata_types_range_word_index` is exposed
- `typelink_resolved_within_types_count` and `typelink_all_resolved_within_types` are exposed
- on the canonical fixture, this models `types..etypes` as the same bounded range currently observed through the `.rodata` pair in `firstmoduledata`
- this is still fixture-local and must not yet be stretched into a generic cross-version `moduledata.types` claim

The fifteenth slice adds the next bounded semantic bridge around `firstmoduledata` and `.gopclntab`:
- `moduledata_pcheader_addr` and `moduledata_pcheader_matches_gopclntab` are exposed
- `moduledata_funcnametab_slice_word_index`, `moduledata_funcnametab_addr`, `moduledata_funcnametab_len`, `moduledata_funcnametab_cap`, and `moduledata_funcnametab_within_gopclntab` are exposed
- on the canonical fixture, this confirms a minimal layout hypothesis:
  - the first `firstmoduledata` word points at the current `.gopclntab` base
  - the next three words form a bounded `funcnametab` slice that stays within `.gopclntab`
- this is still fixture-local and must not yet be stretched into a general `moduledata.pclntable` or `pcHeader` decoder

The sixteenth slice adds the next bounded semantic bridge inside the same `.gopclntab` layout:
- `moduledata_cutab_slice_word_index`, `moduledata_cutab_addr`, `moduledata_cutab_len`, `moduledata_cutab_cap`, and `moduledata_cutab_within_gopclntab` are exposed
- on the canonical fixture, this confirms that the next three `firstmoduledata` words form another bounded slice inside `.gopclntab`
- this is interpreted as the current fixture-local `cutab` bridge, not a general decoded table format

The seventeenth slice continues the same bounded `.gopclntab` layout:
- `moduledata_filetab_slice_word_index`, `moduledata_filetab_addr`, `moduledata_filetab_len`, `moduledata_filetab_cap`, and `moduledata_filetab_within_gopclntab` are exposed
- on the canonical fixture, this confirms that the next three `firstmoduledata` words also form a bounded slice inside `.gopclntab`
- this is interpreted as the current fixture-local `filetab` bridge, not a general decoded file table format

The eighteenth slice continues the same bounded `.gopclntab` layout:
- `moduledata_pctab_slice_word_index`, `moduledata_pctab_addr`, `moduledata_pctab_len`, `moduledata_pctab_cap`, and `moduledata_pctab_within_gopclntab` are exposed
- on the canonical fixture, this confirms that the next three `firstmoduledata` words also form a bounded slice inside `.gopclntab`
- this is interpreted as the current fixture-local `pctab` bridge, not a general decoded pc-data table format

The nineteenth slice continues the same bounded `.gopclntab` layout:
- `moduledata_pclntable_slice_word_index`, `moduledata_pclntable_addr`, `moduledata_pclntable_len`, `moduledata_pclntable_cap`, and `moduledata_pclntable_within_gopclntab` are exposed
- on the canonical fixture, this confirms that the next three `firstmoduledata` words also form a bounded slice inside `.gopclntab`
- this is interpreted as the current fixture-local `pclntable` bridge, not a general decoded pcln table format

The current checkpoint is no longer single-fixture only:
- `corpus/fixtures/go-elf-stripped-linux-amd64/fixture.bin` now exists as a stripped ELF variant
- on this stripped fixture, `.gopclntab`, `.typelink`, `.itablink`, and `.go.module` remain present while `runtime.firstmoduledata` is absent from the rich symbol table
- `ReadMetadata()` now falls back to `.go.module` address as `firstmoduledata_addr` for this current stripped ELF layout family
- `analysis.runtime` now also exposes `firstmoduledata_from_go_module_fallback`, so the bounded stripped-fixture source of truth is explicit instead of implicit
- the first bounded runtime-to-heuristic UX cross-check is now landed on top of that fallback:
  - `inspect functions` preserves bounded `package`, `import_path`, and `module_local` metadata from recovered function names plus `build_info.path` for `main`
  - `inspect functions` also preserves bounded `source_file` truth from the existing `pclntab` / `gosym` line table
  - `inspect functions` also preserves bounded `source_line` truth from the same line table
  - `inspect functions` also preserves bounded `autogenerated` truth by explicitly marking functions whose current truthful source file is `<autogenerated>`
  - `inspect packages` preserves a useful `main` package surface by reusing `build_info.path`
  - `inspect types` degrades to `[]` instead of JSON `null` when no truthful `DWARF` type surface exists
- a bounded source-tree UX fallback is now landed on top of the same stripped-fixture checkpoint:
  - `source-tree` returns the module root plus module-local and external package nodes with empty file lists when truthful file evidence is absent
- a further bounded UX slice is now landed on top of the same principle:
  - `inspect strings` preserves absolute candidate `addr` truth by combining already-known region base addresses with candidate offsets
- this reduces symbol-table dependence, but it is still a bounded fixture-family fallback and not a general discovery algorithm

## Why This Slice Was Chosen

This was the lowest-risk next step because it:
- reduces dependence on `DWARF + source-tree` correlation without replacing working recovery paths yet
- adds runtime truth in a form that is easy to test and reason about
- creates a stable bridge into deeper `Sprint 12` work

## Recommended Next Slice

The next bounded `Sprint 12` step should be one of:
1. decode a minimal `moduledata` header from `runtime.firstmoduledata`
2. validate `.typelink` bounds and count raw typelink entries without decoding every type

Recommendation:
- take `2` first if the goal is low-regret runtime evidence
- take `1` first if the goal is deeper long-term architecture alignment with `GoReSym`

Given the current strategy, `2` is the safer next move.

That safer move is now started:
- `.typelink` bounds are exposed
- raw `typelink_count` is exposed
- bounded raw `typelink_sample` is exposed
- raw `typelink_min_offset` and `typelink_max_offset` are exposed
- raw `typelink_negative_count` and `typelink_non_negative_count` are exposed
- `firstmoduledata_in_go_module` and `go_module_word_sample` are exposed
- `moduledata_typelink_slice_word_index`, `moduledata_typelink_len`, and `moduledata_typelink_cap` are exposed
- `itablink_addr`, `itablink_size`, and `itablink_count` are exposed
- `moduledata_itablink_slice_word_index`, `moduledata_itablink_len`, and `moduledata_itablink_cap` are exposed
- `moduledata_memory_ranges_word_index` is exposed
- bounded `moduledata` memory ranges for `.noptrdata`, `.data`, `.bss`, and `.noptrbss` are exposed
- bounded `moduledata` `.rodata` range fields are exposed
- bounded `moduledata` `.text` range fields are exposed with explicit inclusive-end semantics
- the first semantic typelink resolution fields are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `types..etypes` fields, `moduledata_types_range_word_index`, and typelink-in-types counters are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `pcHeader` / `funcnametab` bridge fields are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `cutab` bridge fields are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `filetab` bridge fields are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `pctab` bridge fields are exposed, but still under a bounded fixture-local interpretation
- bounded `moduledata` `pclntable` bridge fields are exposed, but still under a bounded fixture-local interpretation

Recommended next slice after this:
1. cross-check one more minimal `moduledata` slice or range against existing section evidence
2. keep widening raw evidence only if a specific ambiguity remains

Current recommendation:
- `1` is now the clearly better next move, but it should stay bounded and field-specific rather than becoming a full `moduledata` parser
- after the new `pcHeader` / `funcnametab` bridge, the next step should add one more small semantic structure or cross-check, not restart broad evidence accumulation
