# GoREveal RT1 Horizon A Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore trustworthy GoREveal analysis and verification, publish an identity- and location-safe v2 provider contract, and validate the first conservative GoREveal-to-idacli function preview/apply workflow.

**Architecture:** GoREveal remains the Go-native evidence provider. Horizon A first makes stage status, raw/refined identity, diff cardinality, and verification truthful; it then adds binary/build identity and a format-neutral location contract without changing v1 implicitly. Stateful IDA mutation remains in idacli and is promoted through a separate reviewed plan after the provider contract is frozen.

**Tech Stack:** Go 1.26, `debug/buildinfo`, `debug/elf`, `debug/pe`, `debug/macho`, SHA-256, SQLite, JSON contracts, Python plugin fixtures, Podman, Taskfile, golangci-lint, Ruff, ty, yamllint, shellcheck.

---

## Plan Authority and Scope

Design authority:

- `docs/superpowers/specs/2026-07-22-goreveal-rt1-product-design.md`

This plan is implementation-ready for:

- `RT1-S0` truth restoration;
- `RT1-S1` reproducible evidence baseline;
- `RT1-S2A` identity, build provenance, location contract, and real PIE evidence;
- `RT1-S2B` measured v2 envelope, verifier, and consumer publication.

`RT1-S3` has an exact promotion and validation checklist, but idacli code work
must receive its own plan inside `/opt/projects/repositories/idacli` after the
v2 fixture is frozen. This prevents two repositories from being changed by one
atomic task or commit.

Later sprints are defined by outcomes, dependencies, and gates in the sprint
registry below. Create their file-level plans only when their Definition of
Ready is satisfied.

## Required Skills During Execution

- Use `@goreveal-navigation` at the start of every implementation task.
- Use `@goreveal-doc-sync` when a task changes contracts or active status.
- Use `@goreveal-cleanroom` for build ID, runtime, type, or baseline study.
- Use `@goreveal-corpus-validation` for fixture and snapshot changes.
- Use `@goreveal-differential-testing` for match or baseline changes.
- Use `@goreveal-deobfuscation` for raw/refined identity changes.
- Use `@goreveal-export-contracts` for v1/v2 and consumer work.
- Use `@goreveal-release-ops` for pinned tools, CI, license, and release claims.
- Use `@test-driven-development` before implementation changes.
- Use `@verification-before-completion` before every commit and handoff.

All development and verification commands run in the project Podman dev
container. Keep runs sequential.

## Program Sprint Registry

| Sprint | Status at plan publication | Outcome | Dependencies | Exit gate summary |
| --- | --- | --- | --- | --- |
| `RT1-S0` | ready | no silent loss or positional/collision-corrupted truth | approved RT1 design | typed stage status, keyed refinement, `1:1` auto-match only, correct half-open ranges |
| `RT1-S1` | ready after S0 | every green gate performs real reproducible work | S0 | pinned tools; lint/test/differential/plugins/snapshots/script lint/fuzz/bench green; IDA baseline recorded |
| `RT1-S2A` | planned | exact binary/build identity, format-neutral locations, and real PIE evidence | S0, S1 gate truth | identity/provenance parity, v1 freeze, four-format location round trips |
| `RT1-S2B` | blocked on S2A | measured explicit v2 envelope, reference verifier, and consumers | closed S2A | mismatch rejection, selected-envelope performance, four-format v2 fixtures |
| `RT1-S3` | promotion-gated | safe function-only GoREveal-to-IDA preview/apply | frozen S2B fixture; idacli plan | zero unsafe mutations; idempotency; measurable target improvement |
| `RT1-S4` | conditional | Go entity and source semantics | S0 stable IDs | exact entity decomposition; role/prologue precision >=99%; distinct full paths retained |
| `RT1-S5` | conditional | safe string extents and host-reference navigation | S2A locations, S3 host contract | exact extents; zero unsafe string actions; fixed xref/caller query parity |
| `RT1-S6` | research-gated | resilient metadata candidates and runtime type identity | S1 corpus gate, S2A locations | zero false selected candidates on negatives; exact supported type identity |
| `RT1-S7` | research-gated | layouts, methods, interfaces, preview-only type apply | S6 | exact two-version layout/edge manifests or documented reduced scope |
| `RT1-S8` | deferred | collision-safe build lineage and semantic diff | S0, S4, optional host fingerprints | 100% auto-accept precision; zero false collision accepts; deterministic output |
| `RT1-S9` | deferred | bounded protected/garbled anchor workflow | S8 lineage | twenty verified anchors with zero false accepts or a published negative result |
| `RT1-S10` | deferred | evidence-backed release baseline | S0-S3 minimum, supported-matrix decision | linked release claims, compatibility, pinned image/SBOM, performance envelope |

## File Responsibility Map

### RT1-S0

| File | Responsibility |
| --- | --- |
| `schema/export_v1_wire.go` | recursively frozen IDA/Ghidra v1 JSON DTOs, independent from evolving canonical schema |
| `schema/testdata/export-v1/*.json` | pre-RT1 byte-level compatibility fixtures |
| `schema/diagnostic.go` | typed analysis-stage status and diagnostic contract |
| `schema/entity_id.go` | deterministic local entity ID construction |
| `schema/analysis.go` | canonical IDs and diagnostics on raw analysis entities |
| `schema/refined.go` | keyed refinement identity and refinement kind/span |
| `engine/stages.go` | stage execution, error classification, and test injection |
| `engine/engine.go` | orchestration that records every attempted stage |
| `cmd/goreveal/internal/analyze.go` | command-specific availability enforcement |
| `deobfuscation/pipeline.go` | initialize keyed refined entities |
| `deobfuscation/refine/names.go` | preserve raw identity while changing display names |
| `deobfuscation/garble/strings.go` | emit raw-linked segment findings without sorting away identity |
| `schema/export_ida.go` | map primary refinements by raw ID, never by index |
| `schema/export_ghidra.go` | same keyed mapping rule for Ghidra |
| `storage/diff/diff.go` | cardinality-safe grouping, ambiguity output, ID-based consumed state |
| `core/runtime/moduledata.go` | correct half-open `.text` membership |
| corresponding `*_test.go` files | regression and contract evidence |

### RT1-S1

| File | Responsibility |
| --- | --- |
| `.golangci.yml` and currently reported Go files | restore staged lint policy without weakening it |
| `deployments/docker/Containerfile.dev` | pin Go and Python tooling versions |
| `scripts/dev/podman_runner.py` | real per-package fuzz and benchmark smoke steps |
| `scripts/dev/test_podman_runner.py` | verify commands are non-empty, bounded, and containerized |
| `Makefile`, `Taskfile.yml` | expose the canonical gates without duplicate logic |
| `core/runtime/moduledata_fuzz_test.go` | first bounded parser fuzz target |
| `core/runtime/moduledata_bench_test.go` | runtime metadata reference benchmark |
| `storage/diff/diff_bench_test.go` | large function-set diff benchmark |
| `tests/snapshots/analyze_snapshot_test.go` | discover all declared snapshot fixtures |
| `corpus/fixtures/go-elf-buildinfo-linux-amd64/expected.analysis.json` | rich ELF snapshot |
| `corpus/fixtures/go-elf-stripped-linux-amd64/expected.analysis.json` | stripped ELF snapshot |
| `.github/workflows/ci.yml` | Podman-first required CI sequence |
| `docs/evidence/ida-golang-baseline/README.md` | reproducible experiment contract |
| `docs/evidence/ida-golang-baseline/result.schema.json` | machine-checkable baseline record |

### RT1-S2A

| File | Responsibility |
| --- | --- |
| `docs/architecture/2026-07-22-goreveal-identity-and-location-contract.md` | stable identity/location ADR and envelope requirements handed to S2B |
| `schema/identity.go` | binary, analyzer, module, dependency, and build-setting contracts |
| `schema/location.go` | preferred base, VA/RVA/file-offset/section and `[start,end)` contract |
| `core/identity/identity.go` | streaming SHA-256, format/architecture dispatch |
| `core/identity/buildid.go` | clean-room Go build ID extraction with explicit unavailable state |
| `core/identity/identity_test.go` | known digest, architecture, and build-ID fixtures |
| `core/buildinfo/buildinfo.go` | preserve main module, deps, replacements, and settings |
| `core/location/location.go` | ELF/PE/Mach-O mapping implementation |
| `core/location/location_test.go` | round-trip and unmappable-address fixtures |
| `cmd/goreveal/internal/inspect_dependencies.go` | dependency inventory command |

### RT1-S2B

| File | Responsibility |
| --- | --- |
| `schema/export_ida_v2.go` | separate v2 payload and constructor |
| `schema/export_ghidra_v2.go` | format-neutral v2 host payload if promoted by the ADR |
| `schema/verify_ida_v2.go` | pure reference validation against exact artifact, binary, and loaded-base inputs |
| `cmd/goreveal/main.go` | explicit v1/v2 export selection |
| `cmd/goreveal/internal/export_ida.go` | versioned export dispatch |
| `cmd/goreveal/internal/verify_ida_export.go` | executable detached-digest and target-context reference verifier |
| `plugins/ida/goreveal_ida.py` | explicit contract negotiation or unsupported error |
| `plugins/ghidra/goreveal_ghidra.py` | explicit contract negotiation or unsupported error |
| plugin and schema fixture tests | v1 frozen and v2 explicit compatibility evidence |

## RT1-S0 Detailed Tasks

### Task 0: Freeze v1 wire bytes before canonical schema changes

**Files:**
- Create: `schema/export_v1_wire.go`
- Create: `schema/testdata/export-v1/ida.json`
- Create: `schema/testdata/export-v1/ghidra.json`
- Modify: `schema/export_ida.go`
- Modify: `schema/export_ghidra.go`
- Modify: `schema/export_test.go`
- Test: `schema/export_fixture_test.go`

- [ ] **Step 1: Capture pre-RT1 bytes from the current constructors**

Generate IDA and Ghidra v1 fixtures from the current rich test analysis before
adding diagnostics, entity IDs, build provenance, or locations. Review and
commit the exact JSON bytes; do not regenerate them from the future schema.

- [ ] **Step 2: Add byte-equality tests and observe the safe baseline**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run 'Test(IDA|Ghidra)ExportV1FrozenBytes' -v
```

Expected: PASS against current behavior. This is a characterization test, so
the red proof comes from Step 3 rather than inventing an incorrect fixture.

- [ ] **Step 3: Write the structural isolation test and verify RED**

Add `TestV1WireTypesDoNotEmbedCanonicalContracts`, recursively inspecting the
private v1 wire DTO field types and rejecting `Input`, `BuildInfo`,
`RuntimeMetadata`, `PeelingAnalysis`, `Package`, `SourceTree`, `Function`,
`Type`, `StringCandidate`, or aliases/slices/pointers to them.

Expected: compile FAIL because the private v1 wire DTOs do not exist yet.

- [ ] **Step 4: Isolate v1 recursively from mutable canonical structs**

Add private v1 wire DTOs and explicit projection/custom marshaling for every
nested v1 surface: input, build info, runtime, peeling, packages, source tree,
functions, types, strings, provenance, and refined summary. A v1 wire DTO must
not embed or alias a canonical struct that later RT1 tasks modify. Keep the
public constructors and contract IDs unchanged.

- [ ] **Step 5: Prove structure and bytes are frozen**

Run the recursive no-canonical-embedding guard and exact byte fixtures together.
Tasks 12 and 14 must add future-field leakage cases when dependency and location
fields actually exist.

- [ ] **Step 6: Run plugin compatibility and commit**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -v
make test-plugins
git add schema
git commit -m "test(export): freeze v1 wire contracts"
```

No later task may update these fixtures merely because canonical schema gained
fields. An intentional v1 change requires a separate compatibility decision.

### Task 1: Add explicit analysis-stage diagnostics

**Files:**
- Create: `schema/diagnostic.go`
- Create: `core/recoveryerr/error.go`
- Create: `core/recoveryerr/error_test.go`
- Create: `engine/stages.go`
- Modify: recovery readers under `core/buildinfo`, `core/runtime`, `core/functions`, `core/types`, and `core/strings`
- Modify: `schema/analysis.go:34-46`
- Modify: `engine/engine.go:23-115`
- Modify: `cmd/goreveal/internal/analyze.go:14-118`
- Create: `cmd/goreveal/internal/analysis_policy.go`
- Modify: `cmd/goreveal/internal/deobfuscate.go`
- Modify: `cmd/goreveal/internal/export_ida.go`
- Modify: `cmd/goreveal/internal/export_ghidra.go`
- Modify: `cmd/goreveal/internal/export_sqlite.go`
- Test: `engine/engine_test.go`
- Test: `cmd/goreveal/internal/analyze_test.go`
- Test: `schema/export_test.go`
- Update: existing `corpus/fixtures/*/expected.analysis.json` snapshots

- [ ] **Step 1: Write schema tests for typed status and JSON stability**

Add tests asserting the exact vocabulary:

```go
const (
    StageStatusAvailable   StageStatus = "available"
    StageStatusUnavailable StageStatus = "unavailable"
    StageStatusUnsupported StageStatus = "unsupported"
    StageStatusFailed      StageStatus = "failed"
)

type StageDiagnostic struct {
    Stage   AnalysisStage `json:"stage"`
    Status  StageStatus   `json:"status"`
    Code    string        `json:"code,omitempty"`
    Message string        `json:"message,omitempty"`
}
```

Assert that `Analysis.Diagnostics` serializes as `[]`, not `null`, when no
stages have run.

- [ ] **Step 2: Run the schema test and verify RED**

Run:

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run TestStageDiagnosticJSON -v
```

Expected: FAIL because `StageDiagnostic` and `Analysis.Diagnostics` do not
exist.

- [ ] **Step 3: Add the minimal schema contract**

Create the status/stage types in `schema/diagnostic.go`. Add
`Diagnostics []StageDiagnostic` to `schema.Analysis` and initialize it in
`engine.AnalyzeFile`.

- [ ] **Step 4: Add the core-owned recovery error taxonomy**

Create `core/recoveryerr` with typed `unavailable` and `unsupported` kinds,
stable machine codes, wrapped causes, and `errors.Is`/`errors.As` support.
Recovery packages use `unavailable` only when the artifact truthfully lacks
evidence and `unsupported` only for an identified format/version/architecture
outside the implemented matrix. Bounds, I/O, corrupt data, and unknown errors
remain failures. `engine` maps this taxonomy into schema diagnostics; core does
not import engine or CLI policy.

- [ ] **Step 5: Introduce injectable stage readers**

Refactor `Analyzer` to carry unexported stage operations:

```go
type stageOps struct {
    buildInfo func(string) (schema.BuildInfo, error)
    runtime   func(string) (schema.RuntimeMetadata, error)
    functions func(string) ([]schema.Function, error)
    types     func(string) ([]schema.Type, error)
    strings   func(string) (recoverystrings.Result, error)
    sourceTree func(string, schema.Analysis) (schema.SourceTree, error)
    refine    func(context.Context, schema.Analysis) (schema.RefinedAnalysis, error)
}
```

`New()` supplies production operations. The production source-tree operation
owns the current DWARF/function/fallback sequence so intermediate errors cannot
disappear between helpers. Tests use `newAnalyzerForTest` in the
`engine` package. Keep `core` independent from engine diagnostics.

- [ ] **Step 6: Write the complete status-matrix regression test**

Inject a functions reader returning `errors.New("fixture failure")`. Assert:

- `AnalyzeFile` returns a partial analysis and no top-level error;
- diagnostics contain `stage=functions`, `status=failed`;
- no packages or peeling claims are derived from the failed function stage;
- no empty result is presented as `available`.

For every attempted operation—build info, runtime, functions, types, strings,
source tree, and refinement—table-test success, typed unavailable, typed
unsupported, and unknown failure. Add source-tree and deobfuscation cases so
the two current nested `err == nil` chains cannot remain silent. Assert exactly
one ordered diagnostic per attempted stage and no derived claim from a
non-available prerequisite.

- [ ] **Step 7: Run the engine test and verify RED**

Run:

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./engine -run TestAnalyzeFileRecordsStageFailure -v
```

Expected: FAIL because errors are still silently discarded.

- [ ] **Step 8: Record every attempted stage**

Replace `if recovered, err := ...; err == nil` branches with stage helpers.
Use explicit `unsupported`/`unavailable` only for recognized sentinel errors;
unknown errors are `failed`. Do not turn a failure into `unavailable` merely to
keep output green.

- [ ] **Step 9: Freeze and implement command-specific status policy**

Add a table-driven policy covering every `inspect` subcommand and its required
stage dependencies. At minimum:

| Command | Required stage | Status policy |
| --- | --- | --- |
| `inspect functions` | functions | only `available` succeeds |
| `inspect types` | types | `available` returns data; `unavailable` returns `[]`; `unsupported`/`failed` error |
| `inspect runtime` | runtime | `available` returns data; `unavailable` returns explicit unavailable; `unsupported`/`failed` error |
| `inspect strings` | strings | only `available` succeeds |
| `inspect packages` | functions plus derived package stage | unavailable/unsupported/failed prerequisite cannot look successful |
| `inspect peeling` | functions plus peeling | unavailable/unsupported/failed prerequisite cannot look successful |
| `source-tree` | source tree | only `available` succeeds after its internal truthful fallback |
| `peel` | functions plus peeling | only available prerequisites succeed |
| `deobfuscate` | refinement plus every raw stage used by a finding | failed/unsupported prerequisite or failed refinement rejects |
| `export sqlite` | functions, valid entity IDs, and persistence | functions must be available; optional unavailable stages remain explicit in stored diagnostics; any failed stage rejects |
| `export ida` / `export ghidra` v1 | functions plus every projected stage | functions must be available; any failed projected stage rejects because v1 cannot carry stage diagnostics |

Extend the table when the current CLI exposes another top-level or inspect
surface; no command may fall through to an implicit default. Stored diff/review/
handoff/next commands do not rerun recovery, but must reject the invalid-analysis
error defined in Task 2. Full `analyze` is the only current user command allowed
to emit a partial result with failed stage diagnostics.

- [ ] **Step 10: Run focused and broad tests**

Run:

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./engine ./cmd/goreveal/internal ./schema -v
make snapshot-update
git diff -- corpus/fixtures
make test
```

Expected: PASS; every already-declared snapshot explicitly records stage
diagnostics. Task 9 still adds coverage for fixtures that had no snapshot at
plan publication. Review semantic changes before staging.

- [ ] **Step 11: Commit**

```bash
git add schema core/recoveryerr core/buildinfo core/runtime core/functions core/types core/strings engine cmd/goreveal/internal corpus/fixtures
git commit -m "fix(engine): expose analysis stage diagnostics"
```

### Task 2: Replace positional raw/refined association with stable IDs

**Files:**
- Create: `schema/entity_id.go`
- Modify: `schema/analysis.go:171-224`
- Modify: `schema/refined.go`
- Modify: `deobfuscation/pipeline.go`
- Modify: `deobfuscation/refine/names.go`
- Modify: `deobfuscation/garble/strings.go`
- Modify: `schema/export_ida.go:78-100`
- Modify: `schema/export_ghidra.go:80-102`
- Modify: `storage/sqlite/store.go`
- Modify: `storage/diff/diff.go`
- Modify: `engine/engine.go`
- Modify: `engine/engine_test.go`
- Modify: `cmd/goreveal/internal/diff.go`
- Modify: `cmd/goreveal/internal/analyze_test.go`
- Test: `deobfuscation/pipeline_test.go`
- Test: `deobfuscation/garble/strings_test.go`
- Test: `schema/export_test.go`
- Test: `schema/export_fixture_test.go`
- Test: `storage/sqlite/store_test.go`
- Test: `storage/diff/diff_test.go`
- Update: existing `corpus/fixtures/*/expected.analysis.json` snapshots

- [ ] **Step 1: Freeze the entity ID wire shape in tests**

Use domain-separated, length-prefixed SHA-256 inputs and this representation:

```text
goreveal:<entity-kind>:sha256:<64 lowercase hex>
```

Test determinism, part-boundary collision resistance (`["ab","c"]` differs
from `["a","bc"]`), and distinct string locations with equal values.

- [ ] **Step 2: Run the entity test and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run TestNewEntityID -v
```

Expected: FAIL because `EntityID` does not exist.

- [ ] **Step 3: Implement `EntityID` and raw ID assignment**

Add `ID EntityID` to `Function`, `Package`, `Type`, and `StringCandidate`.
Assign IDs from raw canonical fields before refinement:

- function: name, entry, end;
- package: name, import path;
- type: name, package/import path, kind;
- string: region, address, offset, exact raw value bytes.

Keep this a local artifact identity. Cross-build identity remains RT1-S8 work.

- [ ] **Step 4: Define and test the pre-RT1 zero-ID migration**

Expose `NormalizeEntityIDs(Analysis) (Analysis, error)` and call it after native
recovery, before refinement, before persistence, after
decoding stored `schema.Analysis`, and at the defensive entry to
`storage/diff.Compare` for direct callers. Validate every existing nonzero ID
against the current raw fields; a stale/mismatched ID is invalid rather than
silently trusted. Never insert the zero value into an ID-keyed map. Duplicate
IDs within an entity kind and still-unidentifiable entities fail analysis
validation. Persistence rejects them. Change the diff entrypoint to
`Compare(left, right schema.Analysis) (Summary, error)` and return a typed
`ErrInvalidAnalysis` before allocating match/candidate maps; CLI diff/review/
handoff/next return nonzero and emit no partial summary.

Legacy refined findings without both `ID` and `RawID` are never reconstructed
by array position or display value. Normalization keeps raw entities, removes
the unkeyed refined layer only from the normalized in-memory copy, and appends
`stage=refinement`, `status=unavailable`,
`code=legacy_unkeyed_refinement` to diagnostics. The persisted legacy JSON is
not rewritten on read. A later explicit re-analysis may regenerate keyed
refinements from raw evidence; diff/export cannot consume the quarantined
legacy findings.

Add a frozen pre-RT1 SQLite/JSON fixture with all IDs absent plus duplicate
display names. Assert load/diff backfills deterministic nonzero IDs, preserves
distinct function entries, and produces no zero-ID candidate or accepted
transfer. Add duplicate-derived-ID and mismatched-prepopulated-ID fixtures and
assert fail-closed behavior. Include reordered and segmented unkeyed refined
strings; assert they are not cross-bound, are unavailable in the normalized
copy with the stable diagnostic, and remain byte-unchanged in storage. Do not
rewrite old database rows merely by reading them.

- [ ] **Step 5: Write a reorder-and-segment regression test**

Construct two raw strings with different addresses. Run the garble pass so it
adds segments and reorders findings. Assert every refined string references the
correct raw ID and byte span.

- [ ] **Step 6: Run the regression and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./deobfuscation/... -run 'TestPipelinePreservesRawIDs|TestGarbleSegmentsReferenceRawString' -v
```

Expected: FAIL because refined strings have only `Value`.

- [ ] **Step 7: Add keyed refinement fields**

Use:

```go
type RefinementKind string

const (
    RefinementKindDisplay RefinementKind = "display"
    RefinementKindSegment RefinementKind = "segment"
)

type RefinedString struct {
    ID       EntityID       `json:"id"`
    RawID    EntityID       `json:"raw_id"`
    Kind     RefinementKind `json:"kind"`
    MatchedValue string      `json:"matched_value"`
    Value    string         `json:"value"`
    RawStart uint64         `json:"raw_start,omitempty"`
    RawEnd   uint64         `json:"raw_end,omitempty"`
    Provenance Provenance   `json:"provenance"`
}
```

Give functions/packages/types equivalent `ID`, `RawID`, and `Provenance`
fields and an evidence-derived ID that includes raw ID, finding kind, exact
span, matched bytes, and pass/provider identity. `MatchedValue` preserves the
exact raw substring while `Value` is only normalized display text. Provenance
includes source and confidence for each finding. Use
`FindAllStringIndex` in the garble pass. Adjust byte spans directly for trimmed
prefix/suffix bytes; never rediscover a normalized value with a second search.
Deduplicate only by `(RawID, Kind, RawStart, RawEnd, Provenance.Source,
MatchedValue, Value)`, never globally by display value, and sort by those keys
deterministically. Test equal display/matched values at distinct addresses and
normalization that trims leading/trailing bytes. A segment is a separate
finding, not the primary display value.

- [ ] **Step 8: Replace export index lookup with raw-ID lookup**

Build maps from `RawID` to the single `display` refinement. Ignore `segment`
findings for the singular `RefinedValue` field. If multiple display findings
exist for one raw ID, leave the singular field empty and surface a diagnostic;
do not choose by order.

- [ ] **Step 9: Run focused tests, migration tests, and snapshots**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./deobfuscation/... ./storage/diff ./storage/sqlite -v
make snapshot-update
git diff -- corpus/fixtures
make test
```

Expected: PASS with no address/value cross-binding; existing canonical
snapshots contain deterministic raw entity IDs and Task 0 v1 bytes do not move.

- [ ] **Step 10: Commit**

```bash
git add schema deobfuscation engine storage/diff storage/sqlite cmd/goreveal/internal corpus/fixtures
git commit -m "fix(schema): key refined evidence to raw entities"
```

### Task 3: Make all diff match tiers cardinality-safe

**Files:**
- Modify: `storage/diff/diff.go:216-310`
- Modify: `storage/diff/diff_test.go`
- Modify: `storage/sqlite/store_test.go` if persisted summary shape changes

- [ ] **Step 1: Add duplicate-name and duplicate-source regression fixtures**

Cover all cardinalities:

- `1:1` remains eligible;
- `1:N`, `N:1`, and `N:M` become ambiguity groups;
- duplicate exact names do not get overwritten;
- duplicate basename/package/line keys do not produce score-95 matches;
- no ambiguous member appears in `AcceptedTransfers`.

- [ ] **Step 2: Run the collision tests and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./storage/diff -run 'TestCompareRejects.*Collision|TestCompareEmitsAmbiguityGroup' -v
```

Expected: FAIL because current maps overwrite duplicate keys.

- [ ] **Step 3: Add explicit ambiguity output**

Define a bounded result shape in `storage/diff`:

```go
type FunctionAmbiguity struct {
    Reason    FunctionMatchReason `json:"reason"`
    Key       string              `json:"key"`
    LeftIDs   []schema.EntityID   `json:"left_ids"`
    RightIDs  []schema.EntityID   `json:"right_ids"`
}
```

Sort IDs and groups deterministically.

- [ ] **Step 4: Refactor exact-name and source-location matching into groups**

Use `map[key][]schema.Function` on both sides. Promote only groups where both
lengths are one. Track consumed functions by `EntityID`, not by name. Apply the
same rule to normalized-name and source-file tiers.

- [ ] **Step 5: Carry entity IDs through every downstream transfer type**

Add left/right entity IDs to `FunctionMatch`, `TransferCandidate`, review
actions, and `AcceptedTransfer`. Replace all name-keyed lookup, consumed-state,
candidate construction, package grouping, persistence reconstruction, and
accepted-transfer derivation with ID-keyed joins. Names remain display fields
only. Remove helpers such as `findFunctionByName` from correctness paths.

- [ ] **Step 6: Prove ambiguity and duplicate names cannot reach auto-accept**

Add a table test that passes ambiguity-containing matches through
`buildTransferCandidates` and `buildAcceptedTransfers`. Assert zero accepted
members from every ambiguity group. Add a separate case with two left and two
right functions sharing the same display name but having two unique source
`1:1` pairs; assert each candidate and accepted transfer retains the exact
entry/entity pair and can never cross-bind by name.

- [ ] **Step 7: Run focused and broad tests**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./storage/diff ./storage/sqlite -v
make test
```

Expected: PASS; existing unambiguous matches remain deterministic.

- [ ] **Step 8: Commit**

```bash
git add storage/diff/diff.go storage/diff/diff_test.go storage/sqlite/store_test.go
git commit -m "fix(diff): reject ambiguous automatic matches"
```

### Task 4: Correct half-open text-range validation

**Files:**
- Modify: `core/runtime/moduledata.go:448-490`
- Modify: `core/runtime/moduledata_test.go:356-398`

- [ ] **Step 1: Add exact-end regression cases**

Assert that:

- `start` is inside;
- `end-1` is inside;
- `end` is outside;
- addition overflow does not become an in-range address.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/runtime -run 'TestPopulateELFFunctabAddressHints.*ExclusiveEnd' -v
```

Expected: FAIL because current checks use `<= endExclusive` and
`> endExclusive`.

- [ ] **Step 3: Apply the minimal boundary fix**

Change membership to `addr < endExclusive`. Keep public inclusive-end fields
unchanged until the RT1-S2A location migration; centralize conversion before
comparison.

- [ ] **Step 4: Run tests and commit**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/runtime -v
git add core/runtime/moduledata.go core/runtime/moduledata_test.go
git commit -m "fix(runtime): reject exclusive text endpoint"
```

### Task 5: Close RT1-S0 with clean-room and evidence synchronization

**Files:**
- Verify: `docs/tmp/draft/go-bp.md` in full
- Modify: `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- Modify: `docs/plans/2026-03-20-goreveal-functional-assessment.md`
- Modify: `README.md`
- Modify: `AGENTS.md`

- [ ] **Step 1: Verify the historical research draft remains clean-room safe**

Review the entire draft, not only previously known line ranges. It must retain
the historical-draft warning and may discuss baseline behavior, but must not
recommend copying, translating, forking, or directly reusing baseline
implementation code. Any actionable recommendation must use upstream-Go
primary-source research, independently designed readers, fixture manifests,
and differential comparison.

- [ ] **Step 2: Record only landed S0 claims**

Update semantic boundaries and functional assessment after the code commits.
Do not mark S1/S2A/S2B planned fields as available.

- [ ] **Step 3: Run documentation consistency checks**

```bash
rg -n '(?i)copy|fork|translate|direct code|reuse|использ.*напрям' docs/tmp/draft/go-bp.md
rg -n 'Sprint (12|13|14|15|16)|Sprint [A-E]|RT1-S' README.md AGENTS.md docs/plans docs/superpowers
git diff --check
```

Expected: any remaining copy/fork language is clearly negative or historical;
README and AGENTS point to RT1 as active.

- [ ] **Step 4: Run S0 gate and commit**

```bash
make test
git add docs README.md AGENTS.md
git commit -m "docs: close RT1-S0 truth restoration"
```

Expected: tests pass; the closure record cites the four regression tests.

## RT1-S1 Detailed Tasks

### Task 6: Restore lint without weakening policy

**Files:**
- Modify: `core/runtime/moduledata_test.go`
- Modify: `tests/snapshots/analyze_snapshot_test.go`
- Modify: `core/pclntab/pclntab.go`
- Modify: `cmd/goreveal/internal/diff.go`
- Modify: `core/runtime/moduledata.go`
- Modify: `engine/projection/source_tree.go`
- Modify: `storage/diff/diff.go`
- Modify: `engine/engine.go`
- Modify: `cmd/goreveal/main.go`
- Modify: `storage/diff/diff_test.go`
- Modify: `schema/analysis.go`
- Preserve: `.golangci.yml`

Current pinned-image inventory is 51 findings: `copyloopvar=2`, `dupl=2`,
`funlen=4`, `gocritic=1`, `gosec=17`, `govet=1`, `modernize=4`,
`nonamedreturns=2`, `perfsprint=16`, `unconvert=2`.

- [ ] **Step 1: Save a fresh lint inventory in the task record**

Run `make lint`. Expected before fixes: exit `2` with the 51-finding inventory
above. Do not commit generated logs.

- [ ] **Step 2: Fix correctness-adjacent findings first**

Refactor `runDiffCmd` so length validation occurs before indexing and gosec can
prove it. Fix shadowing and range-copy findings. Run:

```bash
make lint
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./cmd/goreveal/... ./storage/diff ./core/runtime -v
```

- [ ] **Step 3: Refactor long/duplicate functions by responsibility**

Extract helpers instead of adding exclusions. Keep pclntab format readers thin,
split diff handoff composition, split source-tree grouping, and use the
cardinality helpers introduced in S0.

- [ ] **Step 4: Apply mechanical standard-library modernizations**

Replace only behavior-preserving constructs reported by the pinned linter.
Do not change lint configuration or add `nolint` without a specific documented
false positive.

- [ ] **Step 5: Run full lint and tests**

```bash
make lint
make test
```

Expected: both exit `0`; golangci-lint reports zero findings.

- [ ] **Step 6: Commit**

```bash
git add cmd core engine schema storage tests
git commit -m "chore: restore staged Go lint gate"
```

### Task 7: Pin every dev-image tool

**Files:**
- Modify: `deployments/docker/Containerfile.dev:28-42`
- Modify: `deployments/docker/README.md`
- Create: `deployments/docker/requirements-dev.lock`
- Create: `deployments/docker/requirements-host.lock`
- Create: `deployments/docker/toolchain.lock.json`
- Modify: `pyproject.toml`
- Test: `scripts/dev/test_podman_runner.py`

Pin the audited starting set:

- golangci-lint `v2.12.2`;
- benchstat `v0.0.0-20260709024250-82a0b07e230d`;
- govulncheck `v1.6.0`;
- gotestsum `v1.13.0`;
- garble `v0.16.0`;
- buf `v1.72.0`;
- protoc-gen-go `v1.36.11`;
- protoc-gen-connect-go `v1.20.0`;
- podman-py `5.8.0`;
- Ruff `0.15.22`;
- ty `0.0.62`;
- yq `4.1.2`;
- yamllint `1.38.0`.
- shellcheck `0.10.0`.
- GNU time at the exact Debian-snapshot version selected during the lock update.

- [ ] **Step 1: Add a failing image-policy test**

Read `Containerfile.dev` and fail on `@latest`, an unversioned Python package,
a base image without an OCI digest, a live rolling Debian mirror, or an apt
package without an exact version. Require a machine-readable lock record.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest scripts.dev.test_podman_runner.TestDevImagePolicy -v
```

Expected: FAIL on the current `@latest` and unpinned pip installs.

- [ ] **Step 3: Pin the whole build input and document the update procedure**

Pin `golang:1.26-trixie` by OCI digest, use a dated Debian snapshot, install apt
packages at exact versions, use a hash-checked Python requirements lock, and
use exact `module@version` declarations. Record architecture, source URL,
version, and digest where applicable in `toolchain.lock.json`. Document how to
refresh the lock deliberately and require full gates for upgrades.

The host-side Podman launcher is the single bootstrap exception to in-container
execution. Pin `podman==5.8.0` and all transitive wheels with hashes in
`requirements-host.lock`, pin the same exact dependency in `pyproject.toml`,
and make CI install that lock before invoking `scripts.dev.podman_runner`.
No Go build, test, lint, benchmark, or evidence script runs on the host.

- [ ] **Step 4: Rebuild and verify versions**

```bash
task build-image
python3 -m scripts.dev.podman_runner exec -- /go/bin/golangci-lint --version
python3 -m scripts.dev.podman_runner exec -- ruff --version
make lint
make test
```

Expected: declared versions and green gates.

Build the image twice from the same locks and compare the installed tool
manifest. OCI layer metadata may differ, but every declared executable version
and package-lock identity must match.

- [ ] **Step 5: Commit**

```bash
git add deployments/docker pyproject.toml scripts/dev/test_podman_runner.py
git commit -m "build: pin dev container tools"
```

### Task 8: Replace no-op fuzz and benchmark gates

**Files:**
- Create: `core/runtime/moduledata_fuzz_test.go`
- Create: `core/runtime/moduledata_bench_test.go`
- Create: `storage/diff/diff_bench_test.go`
- Modify: `scripts/dev/podman_runner.py:260-265`
- Modify: `scripts/dev/test_podman_runner.py`
- Modify: `Makefile:48-52`
- Modify: `Taskfile.yml:94-102`

- [ ] **Step 1: Test the runner contract first**

Assert fuzz steps target one package and one named fuzz function with bounded
`-fuzztime=2s`. Assert benchmark steps target packages containing named
benchmarks and use `-count=3` for the evidence run.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest scripts.dev.test_podman_runner -v
```

Expected: FAIL because the current fuzz command applies `-fuzz` to `./...` and
there are no real targets.

- [ ] **Step 3: Add the first parser fuzz target**

Seed valid pclntab headers plus truncated, oversized-count, and random inputs.
The property is no panic, bounded allocation, and either a validated result or
a controlled parse failure.

- [ ] **Step 4: Add two representative benchmarks**

Benchmark runtime metadata header parsing and `storage/diff.Compare` on a
deterministic large synthetic function set with collision groups.

- [ ] **Step 5: Wire bounded smoke and evidence commands**

`make fuzz` is a short smoke. Add `make fuzz-evidence` only if a longer corpus
run is required; do not make CI unbounded. `make bench` must print at least one
`Benchmark...` result.

- [ ] **Step 6: Verify**

```bash
make fuzz
make bench
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest scripts.dev.test_podman_runner -v
```

Expected: exit `0`; output includes at least one `Fuzz...` baseline/run and two
`Benchmark...` lines.

- [ ] **Step 7: Commit**

```bash
git add core/runtime storage/diff scripts/dev Makefile Taskfile.yml
git commit -m "test: add real fuzz and benchmark gates"
```

### Task 9: Add missing canonical snapshots and Podman-first CI

**Files:**
- Modify: `tests/snapshots/analyze_snapshot_test.go:79-98`
- Create: `corpus/fixtures/go-elf-buildinfo-linux-amd64/expected.analysis.json`
- Create: `corpus/fixtures/go-elf-stripped-linux-amd64/expected.analysis.json`
- Create: `tests/differential/baselines.lock.json`
- Create: `scripts/baseline/prepare_locked.py`
- Modify: `scripts/baseline/run_goresym.sh`
- Modify: `scripts/baseline/run_redress.sh`
- Modify: `scripts/baseline/run_gore.sh`
- Modify: `tests/differential/differential_test.go`
- Modify: `scripts/dev/podman_runner.py`
- Modify: `scripts/dev/test_podman_runner.py`
- Modify: `Makefile`
- Modify: `Taskfile.yml`
- Create: `.github/workflows/ci.yml`
- Modify: `README.md`

- [ ] **Step 1: Make snapshot discovery test declared fixtures**

Add a test that every fixture with `fixture.json` and `fixture.bin` either has
an expected snapshot or an explicit exclusion with a reason. The current
implementation silently discovers only fixtures that already have expected
files.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./tests/snapshots -run TestAllDeclaredFixturesHaveSnapshotDisposition -v
```

Expected: FAIL for rich and stripped ELF.

- [ ] **Step 3: Generate and review the two snapshots**

Run `make snapshot-update`, inspect every new semantic field, and verify the
stripped fixture still emits `types: []` and explicit fallback evidence.

- [ ] **Step 4: Lock and provision live differential baselines**

Create a machine-readable lock with URL, exact commit SHA, expected license,
build command, executable path, and executable SHA after the controlled build.
The audited starting source identities are:

- GoReSym `78c02cc73064da84dd528220a234e9bd9f133d81`;
- redress `fe38d961b5d8bf0a0ebf58421527da64422a7922`;
- gore `abfc7c568be817973509ef6a27386ba500a1edf4`.

`prepare-baselines` runs inside the pinned dev container during a networked
setup phase, clones into ignored `.tmp/baselines`, checks each HEAD against the
lock, builds/downloads dependencies, and writes a verified executable manifest.
The test phase mounts that directory read-only at `/baselines`, disables
container network, and runs only verified binaries/caches. Remove the implicit
`/opt/projects/repositories` dependency.

`make test` and `make test-differential` must fail with a preparation command if
the locked baseline bundle is absent or mismatched. Replace the current
`t.Skip` path with a hard failure when the differential gate is requested. A
developer may point preparation at an existing checkout only after its remote
URL and exact SHA pass the same lock validation.

- [ ] **Step 5: Add sequential Podman CI**

The workflow installs the hash-locked host Podman client, builds the pinned dev
image, prepares/validates the locked baselines, then runs in order:

1. `make lint`;
2. `make lint-scripts`;
3. `make test`;
4. `make test-differential`;
5. `make test-snapshots`;
6. `make fuzz`;
7. `make bench`.

Do not run workspace-writing jobs concurrently.

- [ ] **Step 6: Verify locally**

```bash
python3 -m scripts.dev.podman_runner task prepare-baselines
make lint
make lint-scripts
make test
make test-differential
make test-snapshots
make fuzz
make bench
git diff --check
```

Expected: all exit `0`, no differential test reports skip, the baseline manifest
contains the three exact SHAs, and all new snapshots are reviewed.

- [ ] **Step 7: Commit**

```bash
git add tests/snapshots tests/differential corpus/fixtures scripts/baseline scripts/dev Makefile Taskfile.yml .github/workflows/ci.yml README.md
git commit -m "ci: enforce reproducible evidence gates"
```

### Task 10: Record license decision and forced IDA baseline

**Files:**
- Create after maintainer decision: `LICENSE`
- Create: `docs/evidence/ida-golang-baseline/README.md`
- Create: `docs/evidence/ida-golang-baseline/result.schema.json`
- Create from experiment: `docs/evidence/ida-golang-baseline/2026-07-teleport-go125.json`
- Modify: `docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md`

- [ ] **Step 1: Obtain the explicit license decision**

Do not infer a license from dependencies or neighboring projects. Record the
decision owner and effective scope. If no decision is available, mark release
work blocked but continue non-release engineering.

- [ ] **Step 2: Define the baseline schema before running IDA**

Require binary SHA-256/size, GoREveal SHA/command, IDA version, Golang plugin
options, input/IDB base, elapsed time, function count, Go-named count, fixed
target list, per-target decompile outcome, and raw log artifact references.

- [ ] **Step 3: Validate the schema with negative fixtures**

Use a small Python unittest or `jsonschema`-free validator in
`scripts/evidence` only if the repo accepts that new helper. Missing tool or
binary identity must fail validation.

- [ ] **Step 4: Run the forced plugin experiment**

Run on the designated licensed workstation using the exact command recorded in
the evidence README. Preserve raw logs outside git if they contain licensed or
sensitive material; commit hashes and bounded results only.

- [ ] **Step 5: Replace proposal `TBD` with evidence references**

Do not rewrite the proposal as active. Link the measured result and state
whether RT1-S3 has a positive delta hypothesis.

- [ ] **Step 6: Verify and commit**

```bash
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest discover -s scripts -p 'test_*.py'
git diff --check
git add docs/evidence docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md
# Add LICENSE and scripts/evidence only when the preceding decision/validator
# steps actually created them.
git commit -m "docs(evidence): record release and IDA baseline decisions"
```

Omit paths that do not exist because the license decision is still blocked;
do not create a placeholder license.

## RT1-S2A Detailed Tasks — Identity and Location Timebox

### Task 11: Freeze the identity and location ADR with clean-room research

**Files:**
- Create: `docs/architecture/2026-07-22-goreveal-identity-and-location-contract.md`
- Create: `docs/architecture/findings/2026-07-22-go-build-id-format-finding.md`
- Create: `corpus/fixtures/identity-location/manifest.schema.json`
- Modify: `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`

- [ ] **Step 1: Study primary sources only**

Use upstream Go toolchain source for Go build ID behavior and standard-library
ELF/PE/Mach-O semantics. Baseline tools supply test ideas, not code.

- [ ] **Step 2: Record format-specific evidence and unsupported cases**

The finding must separate Go build ID from `debug/buildinfo`, define supported
ELF/PE/Mach-O encodings, and require `unavailable` when no supported encoding is
present.

- [ ] **Step 3: Freeze the location model**

The ADR defines preferred image base, VA, RVA, file offset, section/segment,
and half-open ranges. It states which coordinate is canonical in v2 and the
mapping required for rebased consumers.

- [ ] **Step 4: Freeze digest semantics**

Use `sha256:<lowercase-hex>`, exact bytes, no self-digest, and no implicit JSON
canonicalization.

- [ ] **Step 5: Independent architecture review and commit**

Dispatch one clean-room/export-contract reviewer. Resolve Important findings,
then:

```bash
git add docs/architecture corpus/fixtures/identity-location
git commit -m "docs(architecture): freeze identity and location v2"
```

### Task 12: Preserve full build provenance and expose dependencies

**Files:**
- Create: `schema/identity.go`
- Modify: `schema/analysis.go:34-46`
- Modify: `core/buildinfo/buildinfo.go`
- Modify: `core/buildinfo/buildinfo_test.go`
- Create: `cmd/goreveal/internal/inspect_dependencies.go`
- Modify: `cmd/goreveal/internal/analyze.go`
- Modify: `cmd/goreveal/main.go`
- Update: rich ELF, PE, and Mach-O fixture manifests/snapshots

- [ ] **Step 1: Write known-source provenance expectations**

Fixture expectations include main module version/sum, deps, replacements, and
selected settings: GOOS, GOARCH, CGO, build mode, VCS revision/time/modified.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/buildinfo -run TestReadPreservesModulesAndSettings -v
```

Expected: FAIL because current schema keeps only Go version and path.

- [ ] **Step 3: Add schema types and complete mapping**

Map `debug/buildinfo.BuildInfo.Main`, `Deps`, `Replace`, and `Settings` without
claiming a Go build ID. Sort settings and dependencies deterministically while
preserving replacement relations.

- [ ] **Step 4: Add `inspect dependencies`**

Return a stable object containing the main module, dependencies, replacements,
and relevant build settings. Return explicit unavailable status through the S0
diagnostic contract for non-Go binaries.

- [ ] **Step 5: Verify**

Populate dependencies/settings in the v1 fixture analysis before the
frozen-byte assertion; the new canonical fields must not leak into v1.

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/buildinfo ./cmd/goreveal/internal -v
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run 'Test(IDA|Ghidra)ExportV1FrozenBytes' -v
make test-snapshots
```

- [ ] **Step 6: Commit**

```bash
git add schema core/buildinfo cmd/goreveal corpus/fixtures
git commit -m "feat(buildinfo): expose module and dependency provenance"
```

### Task 13: Add exact binary identity

**Files:**
- Create: `core/identity/identity.go`
- Create: `core/identity/buildid.go`
- Create: `core/identity/identity_test.go`
- Modify: `schema/identity.go`
- Modify: `engine/engine.go`
- Modify: `engine/engine_test.go`
- Update: fixture manifests
- Update: existing `corpus/fixtures/*/expected.analysis.json` snapshots

- [ ] **Step 1: Write digest/architecture/build-ID and stage-policy tests**

Cover a known byte digest, ELF amd64, PE amd64, Mach-O amd64, absent build ID,
truncated notes, and a non-Go binary. Table-test the newly attempted
`entity_identity` stage as `available`, `unavailable`, `unsupported`, and
`failed`, following the S0 invariant of exactly one ordered diagnostic per
attempted stage. A malformed or contradictory identity is `failed`; it blocks
refinement, persistence, diff, and export rather than being published as an
empty available identity.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/identity -v
```

Expected: package does not exist.

- [ ] **Step 3: Implement streaming binary SHA-256 and format architecture**

Do not load a 410 MB binary solely to hash it. Emit lowercase
`sha256:<hex>`. Reuse parsed container facts where possible, but keep the
identity package independent from CLI/storage/plugins.

- [ ] **Step 4: Implement only ADR-approved build-ID readers**

Validate bounds before every read. Unsupported or absent encodings produce a
typed unavailable result, not an empty high-confidence value.

- [ ] **Step 5: Integrate as an analysis stage and verify**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/identity ./engine -v
make snapshot-update
git diff -- corpus/fixtures
make test
```

Expected: every canonical fixture snapshot records deterministic identity or
an explicit non-available identity diagnostic. Review all snapshot changes;
do not accept an empty high-confidence build ID.

- [ ] **Step 6: Commit**

```bash
git add core/identity schema/identity.go engine corpus/fixtures
git commit -m "feat(identity): bind analysis to exact binary bytes"
```

### Task 14: Add format-neutral location mapping

**Files:**
- Create: `schema/location.go`
- Create: `core/location/location.go`
- Create: `core/location/location_test.go`
- Modify: `core/pclntab/pclntab.go`
- Modify: `core/runtime/moduledata.go`
- Modify: `schema/analysis.go`
- Create: `corpus/fixtures/go-elf-pie-linux-amd64/src/go.mod`
- Create: `corpus/fixtures/go-elf-pie-linux-amd64/src/main.go`
- Create: `corpus/fixtures/go-elf-pie-linux-amd64/fixture.bin`
- Create: `corpus/fixtures/go-elf-pie-linux-amd64/fixture.json`
- Create: `corpus/fixtures/go-elf-pie-linux-amd64/expected.analysis.json`
- Update: existing `corpus/fixtures/*/expected.analysis.json` snapshots
- Create: `docs/evidence/rt1-s2a-closure.md`

- [ ] **Step 1: Write table tests before implementation**

Cover ELF ET_EXEC, PIE ELF, PE image base plus section RVA, Mach-O VM address,
section start/end, rebased VA, unmapped gaps, and arithmetic overflow.

Define the PIE fixture source and manifest at the same time. The manifest pins
Go/toolchain identity, source hashes, exact build flags, binary SHA-256, ELF
type, preferred base, sections/segments, and known VA/RVA/file-offset tuples.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/location -v
```

Expected: package does not exist.

- [ ] **Step 3: Implement the ADR model**

Use half-open ranges throughout. Keep file offset optional. Never infer image
base from `.text`. A mapping failure is explicit and carries section/segment
context when available.

Build the real PIE fixture in the pinned container:

```bash
python3 -m scripts.dev.podman_runner exec -- bash -lc \
  'cd corpus/fixtures/go-elf-pie-linux-amd64/src && \
   GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /usr/local/go/bin/go build \
   -buildvcs=false -trimpath -buildmode=pie -o ../fixture.bin .'
```

Verify the produced digest and ELF type against the manifest; a drifted binary
is a fixture failure, not an automatic manifest update.

- [ ] **Step 4: Replace duplicate address inference**

Adapt pclntab/runtime projections to consume the location mapper instead of
recomputing PE or ELF translations independently. Preserve raw runtime
evidence fields until a separate compatibility migration removes them.

Add an integration test that reads the real PIE binary, round-trips every
manifest tuple, applies a nonzero loaded-base delta, and rejects a wrong base,
unmapped gap, and changed binary. Generate/review its canonical analysis
snapshot explicitly. Task 16 creates the v2 export fixture after the v2 model
and selected envelope exist.

- [ ] **Step 5: Verify and commit**

Populate the canonical location/runtime additions in the v1 fixture analysis;
the pre-RT1 v1 bytes must remain unchanged.

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/location ./core/pclntab ./core/runtime -v
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run 'Test(IDA|Ghidra)ExportV1FrozenBytes' -v
make snapshot-update
git diff -- corpus/fixtures
make test
git add schema core/location core/pclntab core/runtime corpus/fixtures
git commit -m "feat(core): add explicit address location mapping"
```

- [ ] **Step 6: Close the S2A timebox before envelope work**

Run all established S1 gates sequentially plus the four-format identity/location
matrix and frozen v1 bytes. Record exact commands, fixture hashes, dependency
manifest parity, PIE build recipe, mismatch cases, and remaining unsupported
formats in `docs/evidence/rt1-s2a-closure.md`. Request independent review, then:

```bash
git add docs/evidence/rt1-s2a-closure.md
git commit -m "docs: close RT1-S2A identity and location"
```

If this evidence does not close, stop the timebox with `reduce`/`blocked`; do
not start Task 15.

## RT1-S2B Detailed Tasks — Envelope and Consumer Timebox

Start this timebox only after RT1-S2A closes with reviewed identity, build,
location, PIE, and unchanged-v1 evidence. Do not carry unfinished S2A work into
S2B under the same sprint label.

### Task 15: Build and measure v2 envelope candidates without changing v1

**Files:**
- Create: `schema/artifact_envelope.go`
- Create: `schema/artifact_envelope_test.go`
- Create: `schema/export_ida_v2.go`
- Create: `schema/export_ida_v2_test.go`
- Create: `schema/export_ida_v2_bench_test.go`
- Create: `schema/verify_ida_v2.go`
- Create: `schema/verify_ida_v2_test.go`
- Modify: `cmd/goreveal/internal/export_ida.go`
- Create: `cmd/goreveal/internal/verify_ida_export.go`
- Create: `cmd/goreveal/internal/verify_ida_export_test.go`
- Modify: `cmd/goreveal/main.go`
- Create: `internal/testutil/large_export.go`
- Create: `internal/evidence/envelopeprobe/main.go`
- Create: `scripts/evidence/measure_artifact_envelope.py`
- Create: `scripts/evidence/test_measure_artifact_envelope.py`

- [ ] **Step 1: Re-run the pre-RT1 v1 freeze**

Run the Task 0 IDA/Ghidra byte fixtures after Tasks 12-14 have populated new
canonical fields. They must remain byte-identical; do not update the goldens.

- [ ] **Step 2: Freeze candidate envelope integrity before encoding v2**

Test a format-neutral envelope abstraction with two candidates:

- single JSON: the external approval digest covers the exact file bytes;
- bundle: the external approval digest covers exact UTF-8 manifest bytes, and
  the manifest lists components in canonical order with safe relative ASCII
  name, media type, byte length, record count, and exact SHA-256.

The bundle manifest has no self-digest. Component records use UTF-8 NDJSON,
exactly one JSON record per LF-terminated line, deterministic field/record
ordering, and no implicit canonicalization. Verification rejects duplicate or
out-of-order names, absolute/path-traversal names, missing/extra files, length
or digest mismatch, malformed final newline, and changed manifest bytes. The
approved manifest digest plus verified component digests constitutes one bundle
identity.

- [ ] **Step 3: Write v2 model and verifier tests and verify RED**

Require binary identity, analyzer identity, address-space declaration,
function IDs/locations, build provenance, and diagnostics. The payload and
manifest have no `artifact_sha256` field. The pure verifier receives a parsed
v2 stream plus independently computed target context. Test:

- changed envelope bytes versus detached `sha256:<lowercase-hex>`;
- wrong binary digest/size/format/architecture;
- wrong or undeclared loaded image base;
- unmappable VA/RVA/file-offset relation;
- unknown contract and malformed digest;
- a correct rebased target context.

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./cmd/goreveal/internal -run 'Test(IDAExportV2|ArtifactEnvelope|VerifyIDAExportV2)' -v
```

Expected: v1 passes; v2/envelope/verifier tests fail because the types do not
exist.

- [ ] **Step 4: Implement the model, both candidate encoders, and verifier**

Keep v2 separate from recursively frozen v1. Sort every collection
deterministically. Implement both candidate encoders behind an explicitly
experimental envelope option used only by tests/evidence until Task 16 selects
one. The verifier reads records incrementally for the bundle candidate and
validates the envelope before exposing semantic records.

Add read-only experimental commands equivalent to:

```text
goreveal export ida --contract goreveal.export.ida/v2 \
  --experimental-envelope json|bundle --output <path> <binary>
goreveal verify ida-export --artifact <file-or-manifest> \
  --artifact-sha256 sha256:<64-lowercase-hex> \
  --binary fixture.bin --loaded-base 0x400000
```

The verifier hashes exact envelope bytes, derives binary identity with
`core/identity`, strict-decodes/streams v2, and performs no mutation.

- [ ] **Step 5: Add an executable provider-and-consumer measurement protocol**

`internal/testutil.LargeExportAnalysis(458600)` deterministically constructs
valid ordered v2 records, stable identity/locations, and a matching synthetic
target context without generating or committing a giant Go source tree. The
evidence-only `internal/evidence/envelopeprobe` binary invokes the same public
writer and pure verifier used by the CLI; it is not installed in release
images.

The containerized evidence script builds the actual `goreveal` CLI and the
probe once, then runs provider and reference-consumer subprocesses under the
pinned GNU `time -v`. It measures the probe on the deterministic 458,600-record
analysis, and the real CLI export/verifier on every corpus format plus the
identified 410 MB-class binary. For each envelope and input, record five
isolated runs:
command, binary/tool identities, exit status, wall time, output/component
bytes, record counts, maximum RSS, and verifier result. Measure a deterministic
458,600-function analysis plus the permitted real large binary. The large
binary itself need not be retained, but its identity, commands, and bounded
measurements are mandatory; Task 16 remains blocked if it cannot be measured.
The consumer must parse/validate every record, not merely hash the files. Also
measure the unchanged v1 provider path as the full-analysis time baseline.

Unit-test parsing of `time -v`, failed subprocesses, incomplete bundles, and
metric JSON. All script tests run through the Podman runner.

- [ ] **Step 6: Run candidate tests and commit the pre-decision machinery**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./cmd/goreveal/... -v
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest scripts.evidence.test_measure_artifact_envelope -v
git add schema cmd/goreveal internal/testutil internal/evidence scripts/evidence
git commit -m "feat(export): add measurable v2 envelope candidates"
```

Do not update IDA/Ghidra consumers or call v2 published in this task.

### Task 16: Select, freeze, and publish the v2 envelope

**Files:**
- Modify: v2 envelope/export/verifier files from Task 15
- Modify: `cmd/goreveal/internal/export_ida.go`
- Modify: `cmd/goreveal/internal/verify_ida_export.go`
- Create: `docs/evidence/rt1-s2-artifact-envelope.md`
- Modify: `plugins/ida/goreveal_ida.py`
- Modify: `plugins/ghidra/goreveal_ghidra.py`
- Modify: plugin tests and fixtures
- Update: ELF, PIE ELF, PE, and Mach-O v2 fixtures

- [ ] **Step 1: Run both candidate envelopes through the real protocol**

Run the Task 15 provider and reference consumer sequentially in the pinned
container. Attach the machine-readable measurements and commands to the
evidence record. If the identified real large binary cannot be measured, record
the external blocker and keep Task 16 blocked; the synthetic large case does
not waive the real-binary gate.

- [ ] **Step 2: Apply the decision rule once**

- select single JSON only if both provider and reference-consumer peak RSS are
  at most `2x` total artifact bytes on the large case and remain within the
  time envelope below;
- otherwise select manifest plus NDJSON;
- record go/reduce/kill, trade-offs, and the exact evidence SHA;
- do not change the transport later without a new contract version.

For a bundle, the externally approved digest is always the exact manifest
digest and every ordered component must pass its manifest length/digest/count
check. For single JSON, the approved digest is exact JSON bytes. `idacli` and
all other consumers use the same envelope rule.

Pre-registered time envelope on the pinned S1 reference machine:

- the 458,600-record probe provider and verifier each have median wall time at
  most `120s`, no individual run above `180s`, and no failed validation;
- the real large v2 provider median is at most `125%` of the measured unchanged
  v1 full-analysis/export median;
- the real large reference verifier median is at most `300s`;
- single JSON is selected only if it passes those limits and its median is no
  more than `120%` of the bundle median for both provider and verifier.

If neither candidate meets the hard envelope, take `reduce` and keep v2
unpublished while optimizing the measured bottleneck; do not relax thresholds
after seeing results.

- [ ] **Step 3: Remove experimental ambiguity and freeze CLI behavior**

Publish:

```text
goreveal export ida --contract goreveal.export.ida/v1 <binary>
goreveal export ida --contract goreveal.export.ida/v2 --output <path> <binary>
```

Keep no-flag default on v1 during the compatibility window. Make the selected
v2 envelope the only public v2 encoding; remove or test-hide the rejected
candidate. Freeze fixtures for exact envelope bytes and detached approval
digest.

- [ ] **Step 4: Make consumers explicit after selection**

Existing thin IDA/Ghidra consumers either parse the selected tested v2 subset
and validate its whole envelope or return an exact unsupported-contract error.
They must not treat v2 as v1 or validate only the manifest while ignoring
components.

- [ ] **Step 5: Verify the selected envelope across the real fixture matrix**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./cmd/goreveal/... -v
make test-plugins
make test-snapshots
```

Expected: pre-RT1 v1 bytes pass; exact v2 envelope fixtures pass for ELF, PIE
ELF, PE, and Mach-O; every artifact/binary/base/address mismatch rejects; the
rejected experimental encoding is not exposed as a public v2 alternative.

- [ ] **Step 6: Commit the decision and published contract**

```bash
git add schema cmd/goreveal plugins corpus/fixtures docs/evidence/rt1-s2-artifact-envelope.md
git commit -m "feat(export): publish measured identity-bound IDA v2"
```

### Task 17: Close RT1-S2B evidence and compatibility

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- Modify: `docs/plans/2026-03-19-goreveal-feature-map.md`
- Create: `docs/evidence/rt1-s2b-closure.md`

- [ ] **Step 1: Run all S2B mismatch and round-trip fixtures**

Wrong binary, modified provider byte, wrong base, unmappable location, and
unknown contract must all reject. ELF, PIE ELF, PE, and Mach-O must round-trip.

- [ ] **Step 2: Run complete gates sequentially**

```bash
make lint
make lint-scripts
make test
make test-differential
make test-snapshots
make fuzz
make bench
```

- [ ] **Step 3: Record measured closure and update only truthful claims**

Do not claim idacli preview/apply yet. Link the frozen provider fixture and
transport envelope.

- [ ] **Step 4: Commit**

```bash
git add README.md AGENTS.md docs
git commit -m "docs: close RT1-S2B export contract"
```

## RT1-S3 Promotion and Cross-Repository Plan Gate

### Task 18: Consolidate and promote the idacli implementation plan

**GoREveal inputs required before promotion:**

- frozen `goreveal.export.ida/v2` fixture;
- selected envelope fixture, exact provider approval digest, and, for a bundle,
  ordered component length/digest/count fixtures;
- address mapping fixtures including rebased and mismatch cases;
- function action fixture with missing, matching, named, unnamed, and boundary
  conflict cases;
- S1 forced Golang-plugin baseline.

**idacli documentation authority at RT1 publication:**

- existing draft to reconcile, not duplicate:
  `/opt/projects/repositories/idacli/docs/planning/2026-07-22-go-function-recovery-task.md`;
- approved implementation belongs under `docs/planning/active/` according to
  idacli's `docs/standards/documentation.md`;
- the promoted plan and any new ADR must be linked from `docs/index.md`;
- if a separate durable architecture decision is necessary, use idacli's
  `docs/decisions/adr-NNNN-kebab-case.md` convention, not GoREveal's
  `docs/superpowers/specs/` layout.

**Expected idacli implementation ownership:**

- `src/common/task_request.*`: top-level preview/apply request fields;
- `src/tasks/task_go_preview.*`: read-only identity/mapping/action comparison;
- `src/tasks/task_go_apply.*`: exact-preview application;
- `src/tasks/task_factory.cpp` and task enum: registered executable tasks;
- `src/cli/idb_isolation.*`: retained private IDB copy/checkpoint;
- focused C++ unit/fixture tests plus licensed integration evidence.

- [ ] **Step 1: Re-audit current idacli HEAD and instructions**

Do not assume the July draft matches the current task parser or runtime
discovery work. Read idacli `AGENTS.md`, its documentation skill, documentation
governance, naming rules, `docs/index.md`, current roadmap/backlog, and the
existing draft. At RT1 plan publication, the observed idacli HEAD was
`15a41d7`; always record the actual execution-time SHA.

- [ ] **Step 2: Reconcile the existing draft before promotion**

The draft currently contains stale self-digest and identity/base assumptions.
Revise or supersede it explicitly; do not create a second competing plan. Move
the approved implementation plan into `docs/planning/active/` with `git mv`,
valid lifecycle front matter, and atomic inbound-link/`docs/index.md` updates.
Create an ADR only if idacli maintainers decide the provider/preview/apply
boundary is a durable architecture decision.

- [ ] **Step 3: Freeze the idacli contract around exact preview bytes**

`go-preview` validates the complete selected provider envelope and writes an
immutable `idacli.go-preview/v1` artifact without a self-digest. `go-apply`
receives preview path plus expected
`sha256:<lowercase-hex>`, rehashes exact bytes, validates embedded provider and
IDB identities, and applies the embedded actions. Consume the S2B reference
verifier fixtures rather than reinterpreting GoREveal VA/RVA/base rules.

- [ ] **Step 4: Keep first apply conservative**

Allow create missing function, set default name, and additive comment only.
Skip boundary conflict and user name. Forbid `del_func` and automatic resize.

- [ ] **Step 5: Review and approve the separate idacli plan**

Use one plan-document reviewer and run idacli's naming, lifecycle, link,
Markdown, Mermaid, commit-range, and applicable product-contract gates. No
GoREveal commit contains idacli source changes.

### Task 19: Validate end-to-end function-only workflow

**Artifacts:**
- GoREveal v2 provider artifact;
- idacli preview artifact and external approval digest;
- isolated before/after IDB copies;
- fixed target list and before/after decompile results;
- `docs/evidence/rt1-s3-closure.md` in GoREveal;
- corresponding idacli evidence record.

- [ ] **Step 1: Run all identity rejection cases**

Wrong binary, changed provider artifact, changed preview artifact, wrong IDB,
wrong base, and stale preview must reject with zero mutation.

- [ ] **Step 2: Run the fixture preview classification**

Expected classes must match exactly and ordered action bytes must be
deterministic.

- [ ] **Step 3: Apply on an isolated fixture IDB twice**

First apply performs only reviewed safe operations. Second apply reports zero
mutations.

- [ ] **Step 4: Run the licensed real-binary experiment**

Compare forced Golang-plugin baseline to post-apply state on the predeclared
target set. Rollout requires at least one missing/conflicting target to become
usable and no previously usable target to regress.

- [ ] **Step 5: Record go/reduce/kill**

- `go`: positive delta and all safety gates pass;
- `reduce`: preview is useful but mutation is not yet justified;
- `kill`: no useful delta or an unsafe mutation occurs.

- [ ] **Step 6: Sync both repositories independently**

Commit evidence and docs separately in GoREveal and idacli. Do not create a
cross-repository atomic-history claim.

## Later Sprint Backlog — Promote Only After Definition of Ready

### RT1-S4: Go entity and source semantics

Tasks:

- parse raw symbol into receiver, method, closure lineage, generic
  instantiation, ABI wrapper, init, and autogenerated evidence;
- preserve full normalized source identity separately from display basename;
- expose package -> file -> function navigation;
- classify module, stdlib, dependency, runtime, and unknown roles;
- collect architecture/version-specific prologue bytes as raw evidence;
- promote automatic role/prologue labels only at >=99% precision.

### RT1-S5: Safe strings and host-evidence navigation

Tasks:

- recover exact raw extents independently from display strings;
- distinguish candidate, header/reference, decoded value, and refinement;
- add segment and overflow validation;
- define identity-bound host observation import;
- project string -> xref -> function -> caller and graph queries;
- add preview-only string actions before any apply.

### RT1-S6: Resilient candidates and runtime type identity

Tasks:

- create malformed/split/custom-magic/negative fixtures;
- enumerate pclntab/moduledata candidates with invariant scores and rejected
  reasons;
- select no candidate when evidence is ambiguous;
- decode supported runtime type address/name/kind/size/ptrdata;
- add `inspect type-at` as a diagnostic query;
- keep DWARF and runtime provenance separate.

### RT1-S7: Layouts, methods, and interfaces

Tasks:

- decode proven pointer/element/key/value relations;
- decode struct field names/types/offsets on supported layouts;
- decode interface methods, concrete method sets, and itab edges;
- compare against source-generated manifests across two Go versions;
- publish only proven fields;
- keep host type application preview-only.

### RT1-S8: Build lineage and semantic diff

Tasks:

- add cross-build entity identity distinct from local artifact IDs;
- normalize relocation/address evidence;
- import host fingerprints as hypotheses with provider/version/score;
- add package/module/type/string semantic deltas;
- preserve ambiguity and reason chains;
- measure auto-accept precision on neighboring and unrelated pairs;
- replay only explicitly accepted transfer plans.

### RT1-S9: Protected and garbled anchors

Tasks:

- pin neighboring clean/protected build pair and anchor manifest;
- orchestrate external providers without linking/copying solver code;
- store decoded strings and matches as keyed refined/provider evidence;
- export package/string/callsite anchor preview;
- verify twenty anchors with zero false accepts or publish/freeze a negative
  result.

### RT1-S10: Release baseline

Tasks:

- finalize license/compliance posture;
- generate pinned release image and SBOM;
- publish compatibility and supported-target policy;
- record small/medium/large performance envelope;
- provide one-page local and IDA-assisted operator workflows;
- generate release claims from evidence records only.

## Final Horizon A Verification

Run sequentially from a clean worktree:

```bash
git status --short
task build-image
make lint
make lint-scripts
make test
make test-differential
make test-differential-report
make test-plugins
make test-snapshots
make fuzz
make bench
git diff --check
```

Expected:

- every command exits `0`;
- fuzz and benchmark output proves real targets ran;
- v1 and v2 contract fixtures pass;
- mismatch fixtures reject;
- no uncommitted generated files remain;
- closure evidence links every claimed RT1-S0/S1/S2A/S2B gate.

Before merging Horizon A, request independent code review against the RT1
design and this plan. Fix every Critical and Important issue before proceeding
to RT1-S3.
