---
name: goreveal-differential-testing
description: Turn baseline comparisons into normalized, explicit overlap and divergence evidence.
metadata:
  short-description: Differential testing workflow
---

# GoREveal Differential Testing

## Use When

- comparing behavior with baseline tools
- changing normalization logic
- changing overlap or divergence policy
- updating differential reports

## Workflow

1. Normalize baseline output first.
2. Compare only like-for-like fields.
3. Record overlap, divergence, and uncertainty explicitly.
4. Convert results into tests, fixtures, or findings.
5. Keep machine-readable report output current.

## Rules

- baseline mismatch is evidence, not drama
- richer GoREveal output can be a product improvement, not a bug
- if GoREveal exceeds the baseline, document why
- if a comparison is inherently lossy, say so
