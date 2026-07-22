# GoREveal RT1 Product Design

> Status: approved direction
> Date: 2026-07-22
> Decision owner: maintainers
> Supersedes as active execution guidance:
> - `docs/plans/2026-03-20-goreveal-deferred-continuation.md`
> - `docs/plans/2026-03-31-goreveal-baseline-comparison-plan.md`
> - `docs/plans/2026-04-01-goreveal-protected-binary-comparison-plan.md`
> - `docs/plans/2026-04-01-goreveal-next-execution-plan.md`
> - `docs/plans/2026-04-01-goreveal-post-sprint12-sprint-plan.md`
> - `docs/plans/2026-07-22-goreveal-sprint-roadmap-post-ida-experience.md`
> - `docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md`

## Decision Summary

GoREveal will use an evidence-gated dual-track Scrumban model and will ship the
`RT1` execution train. GoREveal remains the Go-native semantic evidence
provider and build-to-build knowledge engine. IDA, Ghidra, and `idacli` remain
the stateful analyst workspace, decompiler, and controlled mutation layer.

The active product sequence is:

1. restore correctness and trustworthy gates;
2. define binary identity, build provenance, and address semantics;
3. deliver a safe function-only `GoREveal -> idacli -> IDA` preview/apply slice;
4. deepen Go-specific entity, string, runtime type, method, and interface truth;
5. strengthen build lineage, protected-binary workflow, and release evidence.

The project will not execute the July `A-E` proposal literally. That proposal
is retained as field-research input.

## Product Problem

GoREveal already has a strong platform shape:

- canonical schema and CLI output;
- SQLite persistence and stored-run diffing;
- review, handoff, and next-action projections;
- thin IDA and Ghidra consumers;
- bounded runtime evidence and function recovery;
- corpus, snapshot, and differential infrastructure.

The remaining product problem is not missing platform plumbing. It is that a
few correctness defects and unproven claims make downstream automation unsafe:

- analysis-stage errors can be silently discarded;
- refined entities can lose their raw identity;
- collision-prone diff keys can become high-confidence transfer candidates;
- address ranges and format-specific address semantics are not consistently
  explicit;
- the current export cannot bind a downstream action plan to an exact binary;
- runtime types, string extents, and host-tool observations do not yet have a
  sufficiently strong semantic contract.

The user-facing objective is therefore:

> Minimize the median time from a Go binary to an identity-verified,
> review-ready analyst workspace without overstating recovered truth.

## Product Position

### GoREveal owns

- binary and build identity;
- Go build, module, dependency, package, function, type, string, and runtime
  evidence;
- stable Go entity identities;
- raw facts, derived findings, confidence, provenance, ambiguity, and explicit
  unavailable states;
- stored analyses and build-to-build correlation;
- canonical export contracts;
- corpus fixtures, ground-truth manifests, snapshots, differential evidence,
  and performance envelopes.

### Host tools own

`idacli`, IDA, and Ghidra own:

- IDB/project lifecycle and isolation;
- decompilation and generic disassembly;
- xrefs, callgraphs, path analysis, SCC analysis, and host-generated function
  fingerprints;
- preview and mutation of workspace state;
- checkpoints, replay, deadlines, runtime discovery, and licensed-tool
  diagnostics.

### Integration rule

Host-tool observations may return to GoREveal only as identity-bound external
evidence. They do not overwrite GoREveal raw recovery. Plugins consume exports
and apply reviewed actions; they do not become a second Go recovery engine.

## Delivery Methodology

### Framework

Use evidence-gated dual-track Scrumban:

- discovery track: bounded research spikes, primary-source study, fixture
  design, baseline experiments, and claim-boundary decisions;
- delivery track: one operator-visible vertical capability at a time;
- maintenance lane: correctness, corpus, and gate repairs that must not be
  hidden inside product stories.

Default sprint timebox is ten working days. Outcome is fixed; scope is
variable. A discovery spike is normally limited to two working days and must
end in one of:

- `go`: evidence is sufficient to schedule a delivery slice;
- `reduce`: a smaller truthful capability can ship;
- `kill`: the claim is unsupported and the lane is frozen.

### WIP and sequencing

- maximum one semantic/product bet in delivery;
- maximum one maintenance lane;
- verification remains sequential in the bind-mounted workspace;
- do not run independent `A`, `B`, `C`, and `D` feature sprints in parallel;
- only the nearest two or three sprints receive file-level implementation
  plans;
- later sprints remain outcome-and-gate definitions until promoted.

### Prioritization

For normal work:

```text
priority =
  (user_value + risk_reduction + time_criticality)
  * evidence_confidence
  / effort
```

Correctness and unsafe-mutation defects are P0 and bypass the formula.

### Definition of Ready

A delivery task is ready only when it names:

- an analyst pain point;
- a bounded input fixture or real-binary experiment;
- the raw evidence source;
- the semantic claim being made;
- the module owner;
- the expected user-visible output;
- the required test, snapshot, differential comparison, or benchmark;
- a stop or kill condition.

### Definition of Done

In addition to the repository contract, an RT1 slice is done only when:

- unavailable and ambiguous states remain explicit;
- a downstream consumer cannot confuse binaries or address spaces;
- raw and refined truth remain keyed rather than positionally associated;
- every automatic transfer decision records its reasons;
- the relevant Podman-first gates execute real work and pass;
- active docs and release claims are synchronized with the evidence.

## Contract Decisions

### Stage diagnostics

`analyze` may return a partial analysis, but every attempted stage must produce
an explicit status. A stage result distinguishes:

- `available`;
- `unavailable` because the artifact lacks the evidence;
- `unsupported` for an unsupported format/version/architecture;
- `failed` for an unexpected recovery error.

Command-specific strictness decides whether a status is fatal. For example,
`inspect functions` cannot succeed without a truthful function surface, while
full `analyze` may emit a partial result with diagnostics.

### Stable entity identity

Functions, types, strings, and refined findings require stable identifiers.
Identifiers must not depend on array position. Refined evidence references the
raw entity identifier and may additionally carry a raw byte span or provider
observation identifier.

Diff matching is cardinality-safe:

- every match key resolves to a candidate group on both sides;
- only a unique `1:1` candidate may enter an automatic match tier;
- `1:N`, `N:1`, and `N:M` groups are emitted as explicit ambiguity groups;
- an ambiguity group records all members and the key/reason that produced it;
- ambiguous candidates are excluded from `accepted_transfers` and every other
  auto-apply or auto-accept projection;
- later match rules may disambiguate a group only by adding independent
  evidence and recording the complete reason chain.

### Artifact identity

The next export contract separates:

- binary identity: exact input byte digest, size, format, architecture;
- Go build provenance: Go version, real Go build ID when present, module,
  dependency, replacement, VCS, and build settings;
- analyzer identity: GoREveal version and schema contract;
- transport identity: digest of the exact serialized provider artifact.

The provider artifact does not contain a hash of itself. The consumer hashes
the exact bytes it receives. A preview digest binds:

- provider artifact digest;
- loaded workspace/binary identity;
- normalized address mapping;
- the deterministic ordered action plan.

All digest values use the wire representation:

```text
sha256:<64 lowercase hexadecimal characters>
```

Binary digests cover the exact input-file bytes. Provider artifact digests
cover the exact serialized artifact bytes. There is no JSON reserialization or
implicit canonicalization during verification.

The preview output is a separate immutable artifact with contract
`idacli.go-preview/v1`. It embeds the provider artifact digest, workspace
identity, normalized mapping, and ordered actions, but does not embed a digest
of itself. The operator approval value is the SHA-256 digest of the exact
preview artifact bytes. `go-apply` receives both the preview artifact and the
expected approval digest, hashes the exact bytes, validates the embedded
identities again, and applies that embedded plan. It must not reconstruct a
plan from the provider artifact and then compare an implementation-dependent
composite hash.

### Address semantics

The schema must not use `.text` as a synonym for image base. The address
contract distinguishes at least:

- preferred image base where the format defines one;
- virtual address;
- relative virtual address;
- file offset;
- segment or section identity;
- half-open ranges `[start, end)`.

Function exports use one declared canonical address space and carry enough
mapping metadata for a consumer to translate safely after rebasing.

### Export compatibility

`v1` and `v2` are separate explicit contracts. Adding a field does not imply
consumer compatibility. The rollout requires:

- an unchanged v1 fixture and constructor during the compatibility window;
- a dedicated v2 fixture and constructor;
- consumer negotiation or an explicit unsupported-contract error;
- a documented removal policy before v1 can be retired.

### Safe workspace mutation

The first IDA apply slice is function-only and conservative:

- default operation is read-only preview;
- wrong identity or unmappable address is rejected;
- missing functions may be created after review;
- default names may be replaced by reviewed Go names;
- existing user names remain untouched;
- boundary conflicts are reported and skipped;
- `del_func` and automatic boundary replacement are forbidden;
- apply uses an isolated IDB copy/checkpoint;
- applying the same reviewed plan twice is a no-op.

Strings and types remain separate later gates.

### Runtime types

Runtime type delivery is version-aware and layered:

1. candidate and raw layout evidence;
2. type address, name, kind, size, and pointer-data size;
3. pointer, element, key, and value relations;
4. struct fields and offsets;
5. interface methods, concrete method sets, and itab edges;
6. host type-application preview.

DWARF-derived and runtime-derived type facts retain distinct provenance. Field
layout is not exposed as high-confidence truth until it passes the supported
Go-version, architecture, format, and stripped-fixture matrix.

### Strings and host observations

The string model distinguishes:

- printable byte candidate;
- exact raw byte extent;
- Go string header/reference evidence;
- display refinement or decoded value;
- host-tool xref/caller observation.

A refined value never changes a raw address or length. Host observations are
stored as external evidence bound to binary and provider identity.

## RT1 Execution Train

### Horizon A: committed correctness and integration path

- `RT1-S0`: truth restoration;
- `RT1-S1`: reproducible evidence baseline;
- `RT1-S2`: identity, build provenance, and location contract v2;
- `RT1-S3`: safe Go-to-IDA function preview/apply MVP.

### Horizon B: conditional Go semantic depth

- `RT1-S4`: Go entity and source semantics;
- `RT1-S5`: safe strings and host-evidence navigation;
- `RT1-S6`: resilient runtime candidates and type identity;
- `RT1-S7`: layouts, methods, interfaces, and safe type preview.

### Horizon C: evidence-dependent product expansion

- `RT1-S8`: build lineage and semantic diff;
- `RT1-S9`: protected and garbled anchor workflow;
- `RT1-S10`: release baseline.

Detailed outcomes, tasks, dependencies, and gates live in the RT1 program plan.
Only Horizon A receives immediate file-level implementation plans.

## Sprint Exit Gates

| Sprint | Auditable exit gate |
| --- | --- |
| `RT1-S0` | injected failures for each attempted analysis stage appear as `failed`, `unsupported`, or `unavailable`; zero refined/raw misassociations in reorder/addition fixtures; zero automatic matches or accepted transfers from non-`1:1` keys; the exclusive `.text` endpoint is rejected |
| `RT1-S1` | pinned-container `lint`, `test`, differential, plugin, snapshot, script-lint, fuzz-smoke, and benchmark-smoke commands exit `0`; fuzz and benchmark discovery each find at least one real target; rich and stripped ELF canonical snapshots exist; the forced IDA Golang-plugin experiment records binary/tool identities, commands, and raw results |
| `RT1-S2` | wrong binary, changed artifact byte, wrong image base, and unmappable address fixtures are rejected in every case; v1 consumer fixtures remain green; v2 ELF, PIE ELF, PE, and Mach-O location round trips match fixture manifests exactly; dependency/build-setting output matches the known-source manifest |
| `RT1-S3` | all identity-mismatch cases are rejected; unsafe mutations and automatic boundary replacements remain `0`; second apply performs `0` mutations; fixture action classification is exact; real-binary before/after results are recorded, and rollout beyond the experiment requires at least one predeclared target to improve from missing/conflicting to usable without regressing any previously usable target |
| `RT1-S4` | known-source entity decomposition is exact for the supported fixture matrix; auto-labelled function-role/prologue claims meet at least `99%` precision, with unsupported cases left `unknown`; full source identity does not collapse distinct paths to one key |
| `RT1-S5` | exact raw address/length matches the fixture manifest for every exported string extent; zero zero-length, overflow, unmapped, or cross-segment actions enter preview; refined values cannot change the raw extent; string-to-function-to-caller queries reproduce the fixed host-observation fixture |
| `RT1-S6` | selected metadata candidates match fixture ground truth on the declared supported matrix; negative/malformed fixtures produce zero false selected candidates and zero panics; runtime type address/name/kind/size matches the fixture manifest exactly |
| `RT1-S7` | field offsets and method/interface edges match manifests exactly on at least two supported Go versions including one stripped fixture; otherwise the sprint takes the documented `reduce` exit and publishes only type identity and proven relations; type apply remains preview-only |
| `RT1-S8` | auto-accept precision is `100%` on known-source neighboring and unrelated negative build pairs; collision-derived false accepts are `0`; output is deterministic; lower-confidence or host-fingerprint-only matches remain review candidates |
| `RT1-S9` | Definition of Ready fixes the anchor set and false-positive budget before implementation; default minimum is twenty independently verified anchors across a neighboring clean/protected build pair with zero false accepted anchors; otherwise publish the negative result and freeze the lane |
| `RT1-S10` | every supported-target and release claim links to an evidence record; compatibility fixtures pass; release images and SBOM are pinned; small/medium p95 runtime and RSS do not regress more than `20%` from the S1 baseline unless the release record documents and accepts the tradeoff; the large-binary envelope is recorded |

Coverage below a precision gate is not repaired by guessing. It remains an
explicit unknown-rate metric.

## Roadmap Migration Crosswalk

| Previous lane or item | RT1 disposition |
| --- | --- |
| Sprint 12 bounded runtime checkpoint | preserve as evidence baseline; no blind extension; next semantic promotion occurs in `RT1-S6` |
| Sprint 13 workstation handoff contract | preserve landed handoff surfaces; close historical sprint; function-only provider/consumer binding continues in `RT1-S3` |
| Sprint 14 review workflow actionability | preserve and keep frozen; reuse review/next projections in `RT1-S8` rather than adding queue fields in Horizon A |
| Sprint 15 source-evidence confidence | preserve and keep frozen; stable full-path entity identity and source semantics continue in `RT1-S4` |
| Sprint 16 protected commercial workflow | migrate the evidence-first garble/protected hypothesis to `RT1-S9`; no protected delivery starts in Horizon A |
| Previous Sprint 17-18 server/control-plane themes | defer outside RT1 Horizon A/B; not active |
| Previous Sprint 19 release readiness | migrate to `RT1-S10` after local correctness and semantic gates |
| Previous Sprint 20 comparison automation | evidence maintenance remains continuous; build-family product work migrates to `RT1-S8` |
| Previous Sprint 21 build correlation | migrate to `RT1-S8` |
| Previous Sprint 22 metadata knowledge network | defer until after `RT1-S10` |
| Previous Sprint 23 workspace automation/replay | retain host-side replay as an `idacli` concern; broader automation deferred |
| Previous Sprint 24 comparative knowledge packs | defer until after `RT1-S10` |
| July Sprint A identity/VA/v2/digest | replace with corrected identity, exact-byte digest, compatibility, and location design in `RT1-S2` |
| July Sprint B prologue classification | split raw evidence from role/ABI/prologue semantics and move to `RT1-S4` |
| July Sprint C string length | replace display-value length with exact raw extent and host-reference work in `RT1-S5` |
| July Sprint D type layout | split into runtime type identity `RT1-S6` and gated layouts/methods/interfaces `RT1-S7` |
| July Sprint E new diff | cancel duplicate implementation; fix collision safety in `RT1-S0` and extend existing semantic diff in `RT1-S8` |

Already-landed CLI, storage, handoff, source-confidence, peeling, runtime, and
thin-export capabilities remain supported unless an RT1 migration task changes
their documented contract explicitly.

## Capability Transfer Policy

Transfer behavior and user outcomes, never implementation code:

- `GoReSym`: resilient metadata cases, stripped/malformed fixtures, runtime
  type behavior;
- `gore` and `redress`: build/package/type/source navigation behavior;
- `GoResolver`: provider-scored function-match hypotheses;
- `gostringungarbler`: decoded string evidence linked to raw locations;
- `AlphaGolang`: analyst expectations for function, type, and string
  application;
- `idacli`: preview/apply, IDB isolation, diagnostics, decompile, xref,
  callgraph, and replay workflow;
- BinDiff/Diaphora/Ghidra Version Tracking: confidence tiers, ambiguity,
  review, and annotation transfer behavior.

Every capability transfer must leave:

1. a clean-room finding;
2. a fixture and ground-truth manifest;
3. a differential expectation with accepted divergences;
4. an independently designed GoREveal task.

Baseline output is a comparator, not an oracle. Known-source fixture manifests
are the preferred oracle.

## Product Metrics

North-star metric:

- median elapsed operator time from input binary to an identity-verified,
  review-ready workspace.

Mandatory supporting metrics:

- function boundary and name precision, coverage, and unknown rate;
- type, method, and interface accuracy on the supported matrix;
- identity mismatch rejection rate, target `100%`;
- unsafe automatic mutations, target `0`;
- decompilation success on a fixed target corpus before and after apply;
- time and manual steps to the first business-logic function;
- diff auto-accept precision, review rate, and ambiguity rate;
- false transfer count, target `0` for the auto-accept tier;
- supported Go version, GOOS, GOARCH, format, stripped, and protected matrix;
- p50/p95 runtime, peak RSS, and artifact size on small, medium, and large
  binary classes;
- schema compatibility and consumer rejection behavior;
- gate validity and flake rate.

Raw function count is a diagnostic metric, not the north star.

## Explicit Deferrals

Do not start during Horizon A:

- generic decompiler or disassembler work;
- rich TUI/Web UI;
- server, multi-tenant, or distributed analysis platform;
- broad MCP control plane;
- debugger or live-instrumentation ownership;
- native generic crypto/API-hash/rule engine;
- broad GoResolver-like CFG engine;
- binary rewriting;
- broad native garble solver;
- SIMD without measured hotspot evidence;
- another same-fixture pclntab bridge without new semantic value.

## Documentation Authority

After this design is accepted and the plans are reviewed, authority is:

1. `AGENTS.md` for non-negotiable operational and architecture invariants;
2. this design for RT1 product and contract decisions;
3. the RT1 program plan for sprint ordering and gates;
4. the active Horizon implementation plan for exact tasks and commands;
5. sprint closure/evidence records for measured outcomes;
6. historical and research-input plans for context only.

`README.md` must point to the active design and program plan but must not copy a
long mutable sprint status snapshot. Historical plans keep their contents and
receive explicit status/pointer banners instead of silent rewrites.

## Open Decisions Assigned to Sprints

- `RT1-S1`: repository license and release-compliance posture;
- `RT1-S2`: single JSON v2 versus manifest plus streaming records, decided by
  the large-artifact memory/size benchmark;
- `RT1-S2`: exact supported Go build ID sources per format;
- `RT1-S3`: measurable delta after forced IDA Golang-plugin baseline;
- `RT1-S6`: smallest supported runtime type version matrix;
- `RT1-S9`: whether external provider orchestration is sufficient for garbled
  strings and function anchors.

These are bounded sprint decisions, not reasons to weaken the architecture
invariants above.
