---
name: goreveal-deobfuscation
description: Guide deobfuscation passes so they improve readability without corrupting raw recovered truth.
metadata:
  short-description: Deobfuscation rules
---

# GoREveal Deobfuscation

## Purpose

Use this skill when implementing or changing deobfuscation logic.

## Core Rules

- deobfuscation is a refinement layer
- raw recovered truth remains preserved
- refined outputs must carry provenance and confidence
- no deobfuscation pass may silently rewrite the canonical raw source of truth

## Typical Targets

- string ungarbling
- refined symbol names
- package refinement
- optional CFG-guided naming hints

## Required Evidence

- fixture demonstrating the obfuscation case
- test showing raw vs refined separation
- comparison note if behavior is inspired by a baseline tool
