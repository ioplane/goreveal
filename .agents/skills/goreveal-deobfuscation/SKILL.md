---
name: goreveal-deobfuscation
description: Keep deobfuscation as a refinement layer with preserved raw truth and explicit provenance.
---

# GoREveal Deobfuscation

## Use When

- implementing garble or symbol-refinement logic
- adding refined-layer fields
- evaluating external deobfuscation orchestration

## Hard Rules

- raw truth must survive unchanged
- refined truth must stay separate
- provenance/confidence must remain explicit
- no pass may silently overwrite canonical raw recovery

## Evidence Required

- fixture showing the obfuscation case
- test proving raw/refined separation
- comparison note if inspired by an external tool

## Scope Rule

Keep refinement bounded: land the smallest pass that a fixture can prove, and do
not widen into general symbolic recovery on the strength of one working case.

If constraint solving becomes relevant later, prefer external orchestration over
pulling a heavy solver dependency into this module.
