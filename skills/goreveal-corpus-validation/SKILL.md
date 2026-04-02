---
name: goreveal-corpus-validation
description: Manage fixtures, snapshots, and recovery evidence updates without masking semantic changes.
metadata:
  short-description: Corpus and golden validation
---

# GoREveal Corpus Validation

## Use When

- adding a fixture
- changing snapshot expectations
- updating package/type/runtime semantics
- landing new Sprint 12 evidence surfaces

## Required Checks

- fixture metadata is present
- snapshot change is explained
- provenance/confidence meaning stays intact
- stripped vs rich behavior is still explicit
- package/type scope fields remain truthful
- runtime fields remain bounded and non-generic

## Evidence Layers

- corpus fixture
- snapshot output
- differential comparison where relevant
- fuzz or benchmark evidence if parser/perf behavior changed

## Current High-Signal Areas

- `external_packages`
- package `module_local` and `has_source_evidence`
- type `import_path`, `source_file_count`, `module_local`, `user_meaningful`
- bounded runtime trust and `moduledata_*` cross-check fields

## Rule

Do not refresh snapshots blindly.
Every meaningful snapshot change needs a reason.
