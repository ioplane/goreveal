---
name: goreveal-cleanroom
description: Keep baseline-tool study clean-room and convert it into findings, fixtures, and tests instead of copied implementation.
---

# GoREveal Clean-Room

## Use When

- reading `gore`, `redress`, `GoReSym`, `GoResolver`, `gostringungarbler`, or `AlphaGolang`
- comparing outputs against baseline tools
- designing a native GoREveal feature inspired by an external project

## Allowed

- inspect behavior
- capture edge cases
- build fixtures
- write divergence notes
- write differential tests
- extract product lessons

## Forbidden

- copy code
- translate code structure line-for-line
- import AGPL behavior as unexamined truth
- justify a feature with “baseline does it” and no GoREveal design

## Required Output

Produce at least one of:
- a documented finding
- a fixture note
- a differential expectation
- a consciously designed native implementation task

## High-Value Targets

- `GoReSym`-style runtime truth
- `redress`-style source projection
- `gore`-style package/type usability
- bounded deobfuscation ideas that fit raw/refined separation
