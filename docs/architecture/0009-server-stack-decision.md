---
title: GoREveal Server Stack Decision
status: active
date: 2026-03-31
owners:
  - ioplane/goreveal-maintainers
tags:
  - architecture
  - decision
  - server
---

# GoREveal Server Stack Decision

<img
  src="https://shieldcn.dev/badge/status-active-slate.svg?variant=outline&size=xs"
  alt="status: active" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

> **Purpose.** Freeze the leading server-mode technology choices so server planning can proceed
  without re-litigating the base stack on every note.

## Decision Summary

Recommended server-mode baseline:

| Component | Decision | Reason |
| --- | --- | --- |
| Server database | `PostgreSQL 18` | strong JSONB, indexing, and future vector/search headroom |
| Go database access | `pgx/v5 + sqlc` | zero-CGo, explicit queries, type-safe boundaries |
| Migrations | `goose` | embeddable Go library and straightforward `embed.FS` integration |
| API | `ConnectRPC` | one handler family for Connect, gRPC, and gRPC-Web |
| Queue | `River` | PostgreSQL-backed transactional jobs fit the server model |
| Object storage | `S3-compatible` | clean split between metadata control plane and artifact/blob storage |
| Local embedded DB | `SQLite` via `modernc.org/sqlite` | zero-CGo local mode stays aligned with current storage layer |
| Configuration | `koanf v2` | multi-source config for local and server deployments |
| Error stack | `cockroachdb/errors` | strong stack/context handling for parser and orchestration debugging |

## Architectural Position

These choices apply to:

- `goreveal server`
- future `gorectl`
- multi-tenant artifact and job orchestration

These choices do not change the current product boundary:

- `core` remains independent from server/storage concerns
- local CLI analysis must continue working without a server dependency
- server mode is a product extension over canonical schema and engine outputs, not a parallel
  recovery stack

## Control Plane and Data Plane

Recommended transport shape:

- typed control plane over `ConnectRPC`
- artifact bytes over direct object-storage sessions by default

This keeps:

- API servers as coordinators, not large-byte data movers
- transfer logic resumable and scalable
- agent and CLI workflows aligned with the same artifact/job model

## Non-Goals

This decision does not yet commit to:

- final public API schemas
- a tenant or auth model
- a REST-only server surface
- UI work

It only freezes the baseline stack so future planning can stay consistent.

## Related Notes

- `docs/architecture/0001-platform-contract.md`
- `docs/architecture/0003-schema-principles.md`

Runtime-mode, storage, and agent-transport exploration that informed this
decision lives in the maintainers' working notes and is not published.
