---
name: goreveal-perf-simd
description: Apply GoREveal performance and SIMD rules so optimizations remain measurable, reversible, and semantically safe.
metadata:
  short-description: Performance and SIMD discipline
---

# GoREveal Performance and SIMD

## Purpose

Use this skill before any non-trivial performance or SIMD change.

## Required Order

1. measure hotspot
2. implement or confirm scalar reference path
3. improve scalar path if useful
4. add SIMD path only if still justified

## Hard Rules

- no SIMD without scalar fallback
- no SIMD without equivalence tests
- no SIMD without benchmark evidence
- no architecture-specific path without explicit gating
- no optimization that changes schema semantics or provenance meaning

## Acceptable Hotspots

- pattern scanning
- byte classification
- hashing and fingerprinting
- bulk string candidate scan
- large buffer compare/filter kernels
