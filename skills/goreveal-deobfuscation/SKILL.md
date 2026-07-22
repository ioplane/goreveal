---
name: goreveal-deobfuscation
description: Keep deobfuscation as a refinement layer with preserved raw truth and explicit provenance.
metadata:
  short-description: Deobfuscation rules
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

## Sprint Rule

RT1-S0 keyed raw/refined correctness is the only active deobfuscation work
before promotion. String extents and later deobfuscation hypotheses remain
behind the RT1 gates; broad garble work is a protected later outcome, not an
implicit continuation of the historical Sprint 13 queue. If constraint solving
becomes relevant, validate it through external orchestration before accepting
native heavy dependencies.
