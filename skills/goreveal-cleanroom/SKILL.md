---
name: goreveal-cleanroom
description: Enforce clean-room work against baseline Go reverse-engineering projects while preserving useful behavioral references.
metadata:
  short-description: Clean-room workflow for GoREveal
---

# GoREveal Clean-Room Workflow

## Purpose

Use this skill whenever baseline repositories or tools are involved.

## Baseline Set

Current reference set:
- `gore`
- `redress`
- `GoReSym`
- `GoResolver`
- `gostringungarbler`
- `AlphaGolang`

## Allowed Use

- read behavior
- inspect outputs
- document edge cases
- build fixtures
- write differential tests
- compare capabilities and limitations

## Forbidden Use

- copy implementation code
- port code structure line-for-line
- translate AGPL code into “new” local code
- treat baseline output as unquestionable truth

## Required Output Style

When using baseline projects, produce one or more of:
- a documented finding
- a fixture note
- a differential expectation
- a consciously designed GoREveal implementation task

Current high-value clean-room targets:
- `redress`-style source reconstruction and source visibility
- `gore`-style package/type usability improvements
- `GoReSym`-style stronger runtime truth in a future sprint
