# Gemini Notes for GoREveal

<img
  src="https://shieldcn.dev/badge/overlay-gemini-slate.svg?variant=outline&size=xs"
  alt="overlay: gemini" height="20">

Read [`AGENTS.md`](AGENTS.md) first. This file is a research and comparison
overlay and never overrides the two hard rules there.

## Focus areas

- baseline behavior mapping
- differential comparison summaries
- edge-case discovery across reference projects
- translating research into fixtures, findings, and tests

## When working here

- separate observed baseline behavior from GoREveal design decisions
- never turn baseline code structure into an implementation requirement
- turn research into corpus cases and expected-output notes, not into code
- a divergence from `gore`, `redress`, or `GoReSym` is a question to investigate,
  not automatically a GoREveal defect — those tools are evidence, not ground truth

## The clean-room line, concretely

Reading a reference implementation to learn that a header field is little-endian
is research. Reproducing its parsing routine, even reworded, is not. If you cannot
describe a finding as observable behavior without referring to how the other
project's code is organized, you have crossed the line.

See [`docs/architecture/0006-baseline-sources.md`](docs/architecture/0006-baseline-sources.md)
for the reference set and the licensing reason this boundary exists.
