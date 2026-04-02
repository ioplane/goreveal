---
name: goreveal-perf-simd
description: Apply GoREveal’s measure-first, scalar-first performance policy before SIMD or low-level optimization work.
---

# GoREveal Performance and SIMD

## Use When

- optimizing hot paths
- adding hashes/fingerprints over large datasets
- considering architecture-specific acceleration

## Required Order

1. measure the hotspot
2. confirm a scalar reference path
3. improve scalar path if worthwhile
4. add SIMD only if still justified

## Hard Rules

- no SIMD without scalar fallback
- no SIMD without equivalence tests
- no SIMD without benchmarks
- no architecture-specific path without gating
- no optimization that changes schema semantics

## Common Acceptable Targets

- pattern scanning
- byte classification
- hashing and fingerprinting
- bulk string-candidate scans
- compare/filter kernels
