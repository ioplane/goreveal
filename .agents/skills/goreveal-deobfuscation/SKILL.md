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

## Sprint Rule

Do not expand `Sprint 13` ahead of the active `Sprint 12` lane.
If constraint solving becomes relevant later, prefer external orchestration before native heavy dependencies.
