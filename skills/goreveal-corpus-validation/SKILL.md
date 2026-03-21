---
name: goreveal-corpus-validation
description: Manage corpus fixtures, golden snapshots, and evidence expectations for GoREveal recovery behavior.
metadata:
  short-description: Corpus and golden validation
---

# GoREveal Corpus Validation

## Purpose

Use this skill when adding or changing corpus fixtures, golden snapshots, or recovery expectations.

## Required Quality Layers

- fixture metadata
- canonical snapshot output
- version or format coverage note
- provenance/confidence-sensitive review for changed fields
- review of package/type scope fields such as `external_packages`, package `module_local`, and type `import_path`/`source_file_count`/`module_local`/`user_meaningful` when recovery semantics shift
- review of bounded runtime fields such as `.typelink` / `.itablink` evidence, `firstmoduledata`/`.go.module` cross-checks, `moduledata_typelink_*` / `moduledata_itablink_*` slice-header fields, `moduledata_*range*` memory-block fields, `moduledata_rodata_*` / `moduledata_text_*` range fields, and the first fixture-local typelink semantic bridge when runtime semantics shift

## Rules

- Every important recovery claim should be demonstrable on at least one fixture.
- Snapshot changes must be explained, not blindly refreshed.
- Prefer adding focused fixtures over one giant “kitchen sink” corpus sample.
- Keep raw truth and refined truth distinguishable in expected outputs.
- Preserve truthful external/runtime evidence instead of silently dropping it from snapshots.
