---
name: goreveal-export-contracts
description: Keep GoREveal export formats schema-driven and stable for CLI, SQLite, API, IDA, and Ghidra consumers.
metadata:
  short-description: Export contract discipline
---

# GoREveal Export Contracts

## Purpose

Use this skill whenever user-visible or tool-facing outputs change.

## Rules

- all exports derive from canonical schema
- plugins consume exported contracts; they do not recompute recovery logic
- backward-incompatible export changes require explicit documentation
- provenance/confidence fields must not disappear silently
- if schema surfaces gain navigation metadata like `external_packages`, package `module_local`, or type `import_path`/`source_file_count`/`module_local`/`user_meaningful`, decide explicitly whether exports should inherit them or intentionally omit them

## Export Surfaces

- JSON
- protobuf
- SQLite
- IDA payloads
- Ghidra payloads
- service/API responses
