---
name: goreveal-export-contracts
description: Keep exports schema-driven and thin for CLI, SQLite, API, and RE-tool consumers.
---

# GoREveal Export Contracts

## Use When

- changing JSON output
- changing `IDA` / `Ghidra` payloads
- changing SQLite export shape
- deciding whether a new schema field should be projected to exports

## Rules

- exports derive from canonical schema
- plugins consume exports and do not recompute recovery logic
- backward-incompatible changes require explicit docs
- provenance/confidence must not disappear silently
- thin exports may project navigation metadata, but they must not invent truth

## Decision Questions

When schema changes:
1. Should the export inherit the field?
2. Is the field already canonical truth?
3. Would projecting it keep the adapter thinner?
4. If omitted, is that omission intentional and documented?
