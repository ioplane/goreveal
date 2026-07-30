---
title: GoREveal Schema Principles
status: active
date: 2026-03-19
owners:
  - ioplane/goreveal-maintainers
tags:
  - schema
  - contract
---

# GoREveal Schema Principles

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

## Purpose

This document defines the principles for the canonical GoREveal analysis schema.

## Core Principles

- The schema is the canonical product contract.
- Raw recovered truth must be representable without refinement.
- Refined or deobfuscated views must remain separate from raw truth.
- Provenance and confidence are mandatory where recovery is not absolute.
- The schema should be equally usable for CLI, JSON, protobuf, SQLite, and plugin exports.

## Required Concepts

Every major recovered entity should be able to represent:

- raw recovered value
- refined value if any
- provenance
- confidence
- source offsets or addresses where relevant

## Provenance Model

Provenance should distinguish at least:

- direct runtime metadata recovery
- structural parsing inference
- heuristic recovery
- deobfuscation/refinement pass
- imported baseline comparison or externally derived annotation

## Confidence Model

Confidence should be explicit enough to separate:

- strongly recovered facts
- probable inferences
- weak heuristics

## Evolution Rules

- Backward-incompatible schema changes require documentation.
- Silent field repurposing is forbidden.
- Export layers must reflect schema semantics, not invent their own.
- Plugins and APIs consume the schema; they do not redefine it.
