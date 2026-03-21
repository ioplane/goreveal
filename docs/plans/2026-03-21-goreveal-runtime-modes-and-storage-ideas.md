# GoREveal Runtime Modes and Storage Ideas

> Status: product/architecture brainstorming note
> Date: 2026-03-21
> Purpose: capture ideas for dual runtime modes, server control plane, and storage strategy before implementation.

## Scope

This note is intentionally exploratory. It is not an implementation commitment.

The goal is to reason about three related questions:
- how `GoREveal` should work as a local autonomous tool
- how `GoREveal` should work as a multi-tenant server with `gorectl`
- whether `GoREveal` needs a custom artifact database, and which embedded database should back autonomous mode

## Product Framing

`GoREveal` should support two first-class operating modes:

1. `local autonomous mode`
   A single-user, cross-platform binary-first experience for Windows, Linux, and macOS on `x86_64` and `arm64`.

2. `server mode`
   A team and multi-tenant deployment model with durable storage, API tokens, job orchestration, and report lifecycle management.

These modes should share the same analysis core and canonical schema. They should diverge only in orchestration, persistence, auth, and collaboration capabilities.

The safest product principle is:
- one recovery engine
- one schema
- one export/report model
- multiple operating surfaces

## Recommended Product Shape

### Mode A: Local Autonomous

Recommended shape:
- a single `goreveal` binary
- cross-platform release artifacts for `windows/linux/darwin` and `amd64/arm64`
- no required external dependencies for the default experience
- local artifact/project storage
- optional local queue/executor for long-running analysis jobs

What this mode should optimize for:
- zero-friction startup
- portable offline use
- analyst workstation workflows
- easy import/export to `IDA`, `Ghidra`, later `JEB` and `Binary Ninja`
- reproducible local analysis results and reports

Recommended commands:
- `goreveal analyze ...`
- `goreveal inspect ...`
- `goreveal source-tree ...`
- `goreveal export ...`
- `goreveal report ...`
- `goreveal db ...` or `goreveal project ...` for local artifact management

### Mode B: Server

Recommended shape:
- `goreveal server`
- separate `gorectl` client
- API-first design with thin CLI over the same API
- multi-tenant project and artifact model
- async job execution for heavy analysis
- object storage for large artifacts and derived blobs
- relational database for metadata, access control, reports, and search

What this mode should optimize for:
- many users
- many artifacts
- repeatable analysis pipelines
- sharing and permissions
- durable report storage
- background and scheduled jobs
- external automation and CI integration

## `gorectl` Recommendation

A separate `gorectl` is the right direction.

Reasoning:
- it keeps `server` orchestration concerns out of the local autonomous UX
- it follows a proven daemon plus thin control-plane client model already visible in `gobfd`
- it makes auth, profiles, endpoints, and scripting cleaner than overloading the main local CLI

Recommended role split:
- `goreveal`
  - local autonomous execution
  - offline analysis
  - local reports and exports
  - optional embedded single-user project DB
- `goreveal server`
  - API server
  - job orchestration
  - storage coordination
  - tenancy and permissions
- `gorectl`
  - login/profile management
  - artifact upload/download
  - job submission and queue inspection
  - report retrieval
  - admin and operator workflows

## Transport Recommendation

Recommended default: `ConnectRPC/gRPC`.

Why:
- proven in nearby ecosystem patterns like `gobfd`
- efficient and strongly typed
- first-class Go support
- compatible with both CLI control-plane and future web/API consumers
- easier long-term schema governance than ad hoc HTTP JSON only

Recommended transport stance:
- canonical API contract in `proto`
- `ConnectRPC` for operational ergonomics
- `gRPC` compatibility preserved
- JSON/HTTP gateway only if later needed for browser/admin integrations

This is a better fit than inventing a custom binary protocol.

## Recommended Server Stack

### Relational DB

Recommended primary DB: `PostgreSQL 18`.

Why:
- strong transactional model
- excellent JSONB support
- mature indexing, partitioning, and concurrency model
- ideal for tenancy, artifacts, job metadata, reports, and search filters

Recommended categories of data in Postgres:
- users
- groups / workspaces / orgs
- API tokens and auth metadata
- projects
- artifacts and versions
- job definitions and job state
- report metadata
- sharing permissions
- audit trail
- cached analysis summaries and report indexes

Extensions worth evaluating later, not committing yet:
- `pg_trgm` for fuzzy search over names/paths/tags
- `btree_gin` / `btree_gist` depending on mixed index needs
- `pgcrypto` for token and secret handling helpers
- `uuid-ossp` only if native UUID generation policy needs it
- possibly `pgvector` only if semantic similarity use cases become real; not a default requirement

### Object Storage

Recommended artifact/blob store: `S3-compatible storage`, with `Garage` as a good self-hosted default.

Default transfer stance for server mode:
- `gorectl` should use `goreveal server` as the control plane
- large artifact bytes should go directly between `gorectl` and `S3/Garage`
- `goreveal server` should finalize metadata in `PostgreSQL` and coordinate object lifecycle
- server-proxied file transfer should exist only as a fallback or small-file convenience path


Why `Garage` is attractive:
- explicitly S3-compatible
- designed for self-hosted small-to-medium deployments
- geo-distributed and resilient by design
- a good fit for binary artifacts, derived blobs, exports, and report attachments

Recommended S3 bucket/data classes:
- raw uploaded binaries
- derived intermediate blobs
- large analysis payload snapshots
- generated exports for RE tools
- rendered reports (`md`, `json`, maybe `html` later)
- detached symbol/fixture packs if introduced later

### Queueing / Jobs

`River` is a strong candidate.

Why it fits:
- native Go ecosystem fit
- PostgreSQL-backed
- transactional enqueueing
- multiple queues and operational controls
- avoids adding another operational dependency like Redis unless proven necessary

Caveat:
- `River` is excellent if `GoREveal` server stays strongly Postgres-centric
- if future workloads become extremely high-throughput and decoupled from DB transactions, a separate queue may eventually become justified
- for current product direction, `River` is the simpler and lower-regret choice

Recommended job classes:
- artifact ingest
- analysis execution
- enrichment/deobfuscation
- export generation
- report generation
- periodic re-analysis after engine upgrades
- cleanup/retention tasks

## Multi-Tenant Server Product Model

Recommended hierarchy:
- tenant or workspace
- project
- artifact
- analysis run
- report/export

Recommended sharing model:
- private by default
- share to users
- share to groups/teams
- project-level and artifact-level permissions
- immutable run history where practical

Recommended API token model:
- personal tokens
- service tokens
- scoped tokens by workspace/project/action
- expirations and rotation support

Recommended report outputs:
- canonical DB-backed metadata records
- JSON for automation
- Markdown for analyst readability
- optional zipped evidence bundle containing report plus exported payloads

## Do We Need Our Own Artifact Database Like IDA?

Short answer: not yet.

Recommendation:
- do **not** start by designing a proprietary monolithic artifact database format like `IDA`.
- start with a layered storage model:
  - Postgres for relational metadata and indexes
  - S3-compatible object storage for heavy blobs
  - canonical schema snapshots as compressed blobs
  - optional local embedded DB for autonomous mode

Why not build a custom IDA-like DB first:
- it is a large design surface with low immediate product leverage
- it risks locking the product into an opaque storage format too early
- it creates migration and compatibility burden before schema semantics are mature enough
- Postgres plus S3 already gives strong durability, indexing, and separation of concerns

When a custom internal artifact pack might make sense later:
- if local mode needs very fast reopen performance for large workspaces
- if sync/replication between local and server becomes a key workflow
- if compressed snapshot shipping becomes a first-class feature

Even then, the safer direction is not “our own DB engine”, but a structured artifact container format on top of existing proven storage primitives.

## Autonomous Mode DB Options

The user specifically asked for embedded databases that can support `JSONB`-like workflows.

### Best Default Recommendation: SQLite

Recommended default for autonomous mode: `SQLite`.

Why:
- ubiquitous cross-platform support
- excellent operational simplicity
- mature indexing and transactional behavior for embedded use
- already partly aligned with current `GoREveal` local storage direction
- built-in JSON support by default in modern SQLite
- official `JSONB` storage support beginning with SQLite `3.45.0`

Important nuance:
- SQLite `JSONB` is not PostgreSQL JSONB
- it is an internal binary representation stored as a BLOB
- it is smaller and faster than text JSON in SQLite, but it does **not** offer PostgreSQL-style `O(1)` object lookup guarantees

What SQLite is good for in autonomous mode:
- project metadata
- artifact registry
- cached canonical summaries
- local reports and indexes
- small-to-medium report/query workloads
- embedding compressed JSON/JSONB snapshots

Recommended local architecture with SQLite:
- SQLite for metadata and queryable indexes
- filesystem or content-addressed local blob store for large binaries and exports
- optional compressed analysis snapshots stored either in SQLite or sidecar files depending on size thresholds

### Good Secondary Option: DuckDB

Recommended role: analytical companion, not primary operational DB.

Why:
- excellent for ad hoc analytics and report-style querying
- supports a `JSON` data type and JSON functions
- useful if we later want heavy analytical views over many runs

Why not the primary autonomous DB:
- weaker fit for OLTP-style app metadata, locking, and operational embedded app behavior than SQLite
- not a natural replacement for a project/workspace metadata DB
- no need to optimize for analytical SQL before local product ergonomics are finished

### Not Recommended As Primary Autonomous DB

- `Pebble`, `Badger`, `BoltDB`, `LMDB`
  - good key-value engines in some contexts
  - poor fit for relational metadata plus JSON-first querying
  - no native JSONB-first analyst/query experience
- inventing a custom binary artifact DB now
  - too early
  - too expensive in migration and maintenance cost

## Recommended Storage Architecture By Mode

### Local Mode

Recommended baseline:
- embedded `SQLite`
- local blob/CAS directory
- compressed canonical snapshots
- no required external services

Suggested storage split:
- `metadata.db`
  - projects
  - artifacts
  - runs
  - reports
  - tags
  - local sharing state if any
- `blobs/`
  - raw binaries
  - exports
  - large report bundles
- `cache/`
  - transient derived data

### Server Mode

Recommended baseline:
- `PostgreSQL 18`
- `Garage` or another S3-compatible object store
- `River`
- `goreveal server` API
- `gorectl`

Suggested storage split:
- Postgres
  - authoritative metadata and access model
- S3-compatible storage
  - heavy objects and immutable large payloads
- River
  - async analysis orchestration

## Architectural Boundaries

These ideas should preserve current `GoREveal` principles:
- analysis core remains independent from transport and storage choices
- server mode consumes the same schema as local mode
- exports and reports remain schema-driven
- multi-tenant logic stays out of core recovery packages
- no plugin-specific logic in core

## Suggested Future Epic Shape

If this direction is accepted later, it likely becomes:
- `Epic: Dual runtime modes`
- `Epic: Server API and gorectl`
- `Epic: Local embedded project store`
- `Epic: Multi-tenant artifact management`
- `Epic: Async analysis orchestration`
- `Epic: Report persistence and delivery`

## Recommended Decision Snapshot

If a short decision had to be made today:
- yes, support two operating modes
- yes, build a separate `gorectl`
- use `ConnectRPC/gRPC` as the control plane
- use `PostgreSQL 18` plus `Garage` plus likely `River` for server mode
- do **not** build a proprietary IDA-like database first
- use `SQLite` as the default embedded DB for autonomous mode
- keep `DuckDB` as an optional analytical adjunct, not the primary local store

## Source Notes

Useful reference points for these ideas:
- `docs/plans/2026-03-21-goreveal-agent-mcp-and-artifact-transfer-ideas.md` extends this note with MCP and transfer-layer recommendations
- `gobfd` shows a good daemon plus thin CLI pattern with a typed RPC API
- `scrapedoctl` shows a good thin-client ergonomics pattern for a separate control binary
- River official messaging emphasizes PostgreSQL-backed, transactional job enqueueing for Go applications
- Garage positions itself as S3-compatible self-hosted object storage for small-to-medium distributed deployments
- SQLite now has official JSONB support, but it is not equivalent to PostgreSQL JSONB
- DuckDB has a JSON type and strong analytical ergonomics, but is better suited as an analytical adjunct than as the primary embedded operational DB
