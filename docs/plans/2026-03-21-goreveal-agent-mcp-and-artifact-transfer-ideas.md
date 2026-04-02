# GoREveal Agent MCP and Artifact Transfer Ideas

> Status: product/architecture brainstorming note
> Date: 2026-03-21
> Purpose: capture ideas for MCP surfaces for AI agents and for efficient artifact transfer between `gorectl` and `goreveal server`.

## Scope

This note extends the current server-mode brainstorming.

It focuses on two questions:
- how AI agents should interact with `goreveal` and `gorectl`
- how large artifacts and report bundles should move efficiently between client and server

This is not an implementation commitment.

## Current Context

The current product direction already assumes:
- `goreveal` as the local autonomous binary
- `goreveal server` as the multi-tenant server mode
- `gorectl` as a separate thin control-plane client
- `ConnectRPC/gRPC` as the leading candidate for the typed control plane
- `PostgreSQL 18` + `S3-compatible storage` + likely `River` as the leading server-side stack

The open design question is not whether remote operation should exist. It should.
The open design question is how to expose it cleanly to:
- human operators
- automation
- AI agents
- high-volume artifact transfer workflows

## Recommendation Summary

If a short answer had to be written today:
- keep `ConnectRPC/gRPC` as the canonical control plane
- do **not** make `MCP` the canonical server protocol
- add `MCP` as a thin agent-facing adapter layer over `gorectl` and local `goreveal`
- do **not** stream all large artifacts through the RPC API by default
- use a split control-plane/data-plane model
- for large uploads and downloads, prefer direct object-storage transfer with resumable sessions, manifests, and digest verification

## MCP for AI Agents

## Why MCP matters here

`GoREveal` is unusually well-positioned for MCP because its main surfaces are already converging on:
- canonical schema
- stable inspect/export/report operations
- explicit trust/provenance semantics
- project/artifact/job workflows

That is exactly the kind of structured tool surface that agents can use well.

## Recommended MCP Positioning

Recommended architecture:
- `ConnectRPC/gRPC` remains the real service API
- `gorectl` becomes the primary remote MCP bridge
- `goreveal` can also expose a local MCP surface in autonomous mode

This gives two MCP entry points:

1. `goreveal mcp`
   For local autonomous use.
   Example use cases:
   - analyze a local binary
   - inspect functions/packages/types/runtime/strings
   - generate exports and reports
   - manage local project state

2. `gorectl mcp`
   For remote server use.
   Example use cases:
   - authenticate with saved profile/token
   - upload artifact
   - create analysis job
   - watch queue state
   - fetch reports/exports
   - share artifacts/projects
   - administer tenants or groups if authorized

This is better than making the server speak MCP as its primary protocol.

## Why MCP should not be the canonical server protocol

Reasons:
- `MCP` is an integration surface for agent tooling, not a substitute for a stable service API
- server auth, tenancy, and long-running job control are better modeled in a typed RPC API
- non-agent clients will still need a stable public control plane
- direct server-native MCP risks mixing product API design with agent-framework concerns too early

Recommended principle:
- `MCP` sits at the edge
- `ConnectRPC/gRPC` sits at the center

## Host-Platform MCP Interop

If `IDA` or `Ghidra` MCP servers exist in the operator environment, `GoREveal` MCP should complement them, not compete with them.

Recommended workflow:
1. `goreveal mcp` or `gorectl mcp` produces canonical Go-specific analysis and export-ready truth.
2. The agent passes that result to a host-platform MCP server for annotation or import.

Recommended principle:
- `GoREveal` MCP is the Go-native knowledge source
- host-platform MCP remains the analyst workspace integration layer
- do not duplicate `IDA` / `Ghidra` workspace semantics inside `GoREveal` MCP

### Explicit IDA / Ghidra MCP Handoff

The intended workflow should be described more concretely:
1. agent calls `goreveal` MCP `analyze_binary` or a future `export_ida` / `export_ghidra`-style tool
2. agent receives canonical Go-specific schema or an export payload
3. agent passes that result to the host-platform MCP server
4. host-platform MCP applies annotations or imports the payload into the analyst workspace

Concrete tool names will depend on the host-platform MCP server.
The important contract is the handoff shape, not a hard-coded third-party tool name.

Current local product bridge:
- `goreveal diff handoff sqlite <database> <left-id> <right-id>` now exists as a thin operator-facing handoff artifact over the bounded review state
- that CLI path is still local and JSON-only; it does not mutate any host workspace and it does not pretend to be an MCP server
- this is the right current boundary: `GoREveal` prepares the handoff shape first, then future MCP/operator integration can carry it into `IDA`, `Ghidra`, or a workstation host
- the current near-term execution order is now captured separately in `docs/plans/2026-04-01-goreveal-next-execution-plan.md`, where workstation/MCP hardening is the next default move
- the current handoff artifact now also carries structured `target_profiles` for `ida` and `ghidra`, explicit export-contract IDs, preferred transport hints, artifact-role metadata, workspace phases, host action lists, explicit binding-entrypoint hints, required-artifact hints, and expected host-outcome hints, so the next MCP/operator step can bind to a clearer per-target contract instead of only flat recommendation lists

Example operator reading:
- `GoREveal MCP` answers Go-specific questions and prepares canonical payloads
- `IDA` / `Ghidra` MCP applies names, comments, structure hints, and other workspace-facing markup
- `GoREveal` should not try to mirror host-platform project/workspace semantics inside its own MCP surface

### Updated Real-Environment Reading

The host-platform MCP story is no longer only hypothetical.

Measured operator-environment signal:
- the current RE lab host exposes `ida-pro-mcp`
- the same host also exposes `ida-pro`, `ghidra`, `pyghidra`, `headless-ida`, `jeb`, and `rizin`
- `rehelp` already documents remote workstation usage through `Teleport`, which is a good operational pattern for agent-facing orchestration

That makes the intended split more concrete:
- `GoREveal MCP` should focus on Go-native truth, exports, and transfer workflows
- host-platform MCP should focus on workspace mutation and analyst interaction
- remote operator docs should explicitly acknowledge that this handoff may happen through a workstation host rather than on the same local machine

See also:
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md`

## Recommended MCP Tool Families

### Local MCP via `goreveal`

Recommended tools:
- `analyze_binary`
- `inspect_runtime`
- `inspect_functions`
- `inspect_packages`
- `inspect_types`
- `inspect_strings`
- `build_source_tree`
- `generate_report`
- `export_ida`
- `export_ghidra`
- later `export_jeb`
- later `export_binary_ninja`

### Remote MCP via `gorectl`

Recommended tools:
- `whoami`
- `list_projects`
- `create_project`
- `upload_artifact`
- `list_artifacts`
- `submit_analysis_job`
- `get_job_status`
- `list_reports`
- `download_report`
- `download_export`
- `share_artifact`
- `share_project`
- `list_group_memberships`
- `get_artifact_summary`

### MCP Response Style

Recommended response model:
- structured JSON only
- no hidden server-side heuristics that are not already in canonical schema
- large binaries/reports should return handles, URLs, IDs, or saved paths, not inline payloads

## Artifact Transfer Problem

The transfer problem is distinct from the control-plane problem.

The server API must handle:
- artifact creation and metadata registration
- auth and authorization
- session setup
- job creation
- finalization and verification

But the actual heavy bytes do not need to flow through the API server in the same way as metadata.

That distinction is the key design decision.

## Three Practical Transfer Approaches

### Approach A: Pure gRPC/Connect streaming

Model:
- `gorectl` streams file bytes to the server over the control-plane RPC channel
- downloads also stream back over the same API channel

Pros:
- simple mental model
- one transport
- easy to prototype

Cons:
- API servers become data movers instead of coordinators
- weaker resume/retry semantics for very large artifacts
- harder to parallelize efficiently
- harder to exploit object storage directly
- increases pressure on API ingress and memory/buffer tuning

Verdict:
- acceptable as a fallback path
- not the best default for large artifact workflows

### Approach B: Split control plane and data plane

Model:
- `gorectl` talks to the server over `ConnectRPC/gRPC`
- server creates an upload or download session
- actual bytes move directly to or from S3-compatible object storage
- server only finalizes metadata and ownership after verification

Pros:
- scalable
- clean separation of responsibilities
- natural fit with `Garage`
- resumable and parallelizable
- keeps API servers thin
- good for both human CLI and agents

Cons:
- more moving parts than pure RPC streaming
- requires upload-session logic and object-store integration

Verdict:
- recommended default

### Approach C: rsync-like sync protocol or directory mirroring

Model:
- `gorectl` tries to behave like `rsync` or a directory synchronizer
- transfer layer computes incremental differences at file-tree level

Pros:
- can be attractive for workstation sync mental model
- may save bytes in some niche cases

Cons:
- the product mostly deals with immutable artifacts and derived bundles, not mutable directory trees
- adds complexity far earlier than value
- not naturally aligned with object-store-first server storage
- weaker fit for multi-tenant API workflows

Verdict:
- not recommended as the default product transfer model
- useful only later for special sync workflows

## Recommended Transfer Architecture

Recommended default: `Approach B`.

## Default Upload/Download Path Decision

The intended default path is:
- `gorectl -> goreveal server` for control plane
- `gorectl -> S3/Garage` for large artifact bytes
- `goreveal server -> PostgreSQL` for metadata and ownership state
- `goreveal server -> S3/Garage` only for object-management coordination, sealing, and verification

This means the recommended production path is **not**:
- `gorectl -> goreveal server -> S3/Garage`

That server-proxied path may still exist as a fallback for:
- small files
- simple development mode
- constrained environments

But it should not be the default for large artifact workflows because it turns the API server into an avoidable bandwidth bottleneck.


### Upload flow

1. `gorectl` computes fast local metadata:
   - file size
   - content digest, ideally `BLAKE3` for speed plus `SHA-256` if needed for compatibility/policy
   - optional media/type hints
2. `gorectl` calls `BeginArtifactUpload`
3. server checks whether the artifact already exists by digest
4. if already present:
   - do not re-upload bytes
   - only create metadata ownership/reference records
5. if not present:
   - server creates upload session
   - returns object key, multipart/chunk parameters, and upload credentials or pre-signed URLs
6. `gorectl` uploads directly to object storage
7. `gorectl` calls `CompleteArtifactUpload`
8. server verifies digest, seals artifact record, optionally enqueues analysis job

### Download flow

1. `gorectl` requests artifact, export, or report
2. server authorizes access and decides delivery mode
3. for small payloads:
   - server may return direct RPC payload or small blob response
4. for large payloads:
   - server returns pre-signed download URL or download session
5. `gorectl` downloads directly from object storage
6. client verifies digest if provided

## Recommended Optimizations

### 1. Whole-file dedup first

This is the highest-value early optimization.

For many reverse-engineering workflows, binaries and generated bundles are immutable objects.
That means whole-file dedup by content digest is much more valuable, simpler, and safer than early block-delta sync.

Recommendation:
- artifact identity keyed by digest
- logical artifact records can point to the same stored blob
- same approach for large report/export bundles where appropriate

### 2. Multipart parallel upload/download

For large artifacts, use multipart or chunked parallel transfer.

Why:
- saturates fast links far better than single-stream upload
- supports retry of only failed parts
- fits object storage well

This is much closer to the useful lessons from `rclone` than from `rsync`.

### 3. Resumable sessions

Uploads should be resumable by session state, not by inventing a custom sync engine.

Recommendation:
- session ID
- chunk state or multipart part state
- expiration rules
- explicit finalize or abort semantics

### 4. Optional bundle mode for many small files

The `rsync` alternatives note is especially relevant when there are huge trees of tiny files.
That is usually not the main artifact shape for `GoREveal`, but it will matter for:
- report bundles
- evidence packs
- local export trees
- maybe project snapshot export/import

Recommendation:
- for multi-file exports or evidence bundles, prefer packaging first
- use `tar+zstd` or similar bundle format
- then transfer the bundle as one artifact

This is much lower complexity than trying to synchronize many small files individually over the product API.

### 5. Local cache and content-addressed blobs

`gorectl` should eventually maintain an optional local cache:
- digest-addressed blobs
- cached downloads
- already-uploaded artifact fingerprints

That gives:
- faster repeated workflows
- reduced network use
- better agent automation ergonomics

## What The rsync Alternatives Note Changes

The local note at `/opt/docs/rsync-alternative.md` is useful because it reinforces several product choices.

Important takeaways for `GoREveal`:
- `rsync` is bad as a default mental model for very many small files or high-speed links
- `tar+zstd` is extremely effective for first-copy bundle workflows
- parallel transfer matters more than clever delta logic in many real cases
- `rclone`-style concurrency and direct object-store data movement are closer to what a server product needs than `rsync`’s negotiation-heavy model
- block-delta sync is powerful, but it is not the first optimization to build for a system whose main unit is immutable artifact blobs

Practical product consequence:
- design `GoREveal` transfer around artifact objects, not filesystem mirroring
- use bundle formats for multi-file payloads
- use direct object-store multipart transfer for heavy data
- keep RPC for orchestration and metadata

## Should We Ever Add Delta Transfer?

Maybe later, but only after simpler wins exist.

Recommended order of sophistication:
1. whole-file dedup by digest
2. multipart parallel transfer
3. resumable sessions
4. bundle packaging for multi-file outputs
5. local cache/CAS
6. only then evaluate chunk-level dedup or delta transfer

Chunk-level dedup becomes more interesting if later product features include:
- repeated import/export of large project packs
- frequent near-duplicate binary families
- server-side artifact lineage and cross-build correlation

Even then, this should probably look like chunk-addressed storage or bundle dedup, not a generic rsync clone.

## MCP and Artifact Transfer Together

These two concerns should be connected, but not collapsed.

Good MCP behavior for large transfers:
- MCP tools should initiate and monitor transfers
- MCP tools should not inline large blobs into agent messages
- MCP tools should return session IDs, progress, artifact IDs, report IDs, and local saved paths

Examples:
- `upload_artifact` returns artifact ID, digest, upload status, and maybe job ID
- `download_report` returns saved path or download handle
- `create_export_bundle` returns bundle artifact ID, then agent can download it separately

This keeps agent workflows structured and efficient.

## Recommended Future Epics

If this direction is accepted later, likely epics are:
- `Epic: MCP surfaces for local and remote GoREveal workflows`
- `Epic: gorectl upload/download sessions`
- `Epic: object-store-backed artifact transfer`
- `Epic: resumable transfer and local cache`
- `Epic: project snapshot bundles`

## Recommended Decision Snapshot

If a short decision had to be made today:
- yes, add MCP, but as an edge adapter layer
- local MCP belongs in `goreveal`
- remote MCP belongs in `gorectl`
- no, do not make MCP the canonical service protocol
- no, do not default to pure gRPC streaming for large artifacts
- yes, use split control-plane/data-plane transfer
- yes, prefer digest-first dedup, multipart direct uploads, resumable sessions, and bundle packaging
- no, do not build rsync-like sync as the primary transfer model

## Source Notes

Useful references behind these ideas:
- `/opt/docs/rsync-alternative.md`
- `docs/plans/2026-03-21-goreveal-runtime-modes-and-storage-ideas.md`
- `gobfd` for daemon plus thin CLI RPC control pattern
- `scrapedoctl` for thin-client ergonomics
