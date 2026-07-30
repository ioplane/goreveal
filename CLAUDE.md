# Claude Notes for GoREveal

<img
  src="https://shieldcn.dev/badge/overlay-claude-slate.svg?variant=outline&size=xs"
  alt="overlay: claude" height="20">

Read [`AGENTS.md`](AGENTS.md) first. This file is a planning overlay, not a second
source of truth, and it never overrides the two hard rules there.

## Focus areas

- design and architecture consolidation
- decomposing a capability into bounded, independently verifiable increments
- multi-file refactors that preserve module boundaries
- documenting clean-room findings as behavior notes, not code cargo-culting

## When working here

- clarify architectural impact before widening scope
- keep schema, provenance, and raw-versus-refined separation explicit
- do not let plugin, storage, or UI thinking bleed into `core`
- when a change touches recovery semantics, state the new claim and the evidence
  that justifies it before writing the code
- prefer one narrow slice with a fixture over a broad change with none

## What to reject

- a plan that produces no verifiable capability
- an abstraction introduced before a second real implementation exists
- output that fills an evidence gap with inference instead of `unavailable`
