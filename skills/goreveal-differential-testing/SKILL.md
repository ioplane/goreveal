---
name: goreveal-differential-testing
description: Compare GoREveal behavior against baseline tools and convert those comparisons into explicit expected outcomes.
metadata:
  short-description: Differential testing workflow
---

# GoREveal Differential Testing

## Purpose

Use this skill for behavior comparisons against baseline tools.

## Comparison Targets

- function recovery
- package recovery
- type recovery
- string recovery
- source projection
- deobfuscation behavior

## Rules

- Normalize tool outputs before comparing them.
- Record overlaps, divergences, and uncertainty explicitly.
- Better-than-baseline results are allowed, but they must be documented.
- Baseline mismatches should become tests or findings, not ad hoc opinions.
- Treat richer `GoREveal` schema surfaces like `external_packages`, package `module_local`, or type `import_path`/`source_file_count`/`module_local`/`user_meaningful` as potential product improvements unless a baseline exposes the same truthful overlap.
- Keep the machine-readable report path current when overlap or divergence policy changes.
