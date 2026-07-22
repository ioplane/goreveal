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
- `RT1-S2` identity, build provenance, and location contract v2.

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
| `RT1-S2` | planned | exact binary/build identity and format-neutral locations in explicit v2 | S0, S1 gate truth | mismatch rejection, v1 compatibility, four-format v2 round trips, build provenance parity |
| `RT1-S3` | promotion-gated | safe function-only GoREveal-to-IDA preview/apply | frozen S2 fixture; idacli plan | zero unsafe mutations; idempotency; measurable target improvement |
| `RT1-S4` | conditional | Go entity and source semantics | S0 stable IDs | exact entity decomposition; role/prologue precision >=99%; distinct full paths retained |
| `RT1-S5` | conditional | safe string extents and host-reference navigation | S2 locations, S3 host contract | exact extents; zero unsafe string actions; fixed xref/caller query parity |
| `RT1-S6` | research-gated | resilient metadata candidates and runtime type identity | S1 corpus gate, S2 locations | zero false selected candidates on negatives; exact supported type identity |
| `RT1-S7` | research-gated | layouts, methods, interfaces, preview-only type apply | S6 | exact two-version layout/edge manifests or documented reduced scope |
| `RT1-S8` | deferred | collision-safe build lineage and semantic diff | S0, S4, optional host fingerprints | 100% auto-accept precision; zero false collision accepts; deterministic output |
| `RT1-S9` | deferred | bounded protected/garbled anchor workflow | S8 lineage | twenty verified anchors with zero false accepts or a published negative result |
| `RT1-S10` | deferred | evidence-backed release baseline | S0-S3 minimum, supported-matrix decision | linked release claims, compatibility, pinned image/SBOM, performance envelope |

## File Responsibility Map

### RT1-S0

| File | Responsibility |
| --- | --- |
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

### RT1-S2

| File | Responsibility |
| --- | --- |
| `docs/architecture/2026-07-22-goreveal-identity-and-location-contract.md` | stable v2 ADR |
| `schema/identity.go` | binary, analyzer, module, dependency, and build-setting contracts |
| `schema/location.go` | preferred base, VA/RVA/file-offset/section and `[start,end)` contract |
| `schema/export_ida_v2.go` | separate v2 payload and constructor |
| `schema/export_ghidra_v2.go` | format-neutral v2 host payload if promoted by the ADR |
| `core/identity/identity.go` | streaming SHA-256, format/architecture dispatch |
| `core/identity/buildid.go` | clean-room Go build ID extraction with explicit unavailable state |
| `core/identity/identity_test.go` | known digest, architecture, and build-ID fixtures |
| `core/buildinfo/buildinfo.go` | preserve main module, deps, replacements, and settings |
| `core/location/location.go` | ELF/PE/Mach-O mapping implementation |
| `core/location/location_test.go` | round-trip and unmappable-address fixtures |
| `cmd/goreveal/main.go` | explicit v1/v2 export selection and `inspect dependencies` |
| `cmd/goreveal/internal/export_ida.go` | versioned export dispatch |
| `cmd/goreveal/internal/inspect_dependencies.go` | dependency inventory command |
| `plugins/ida/goreveal_ida.py` | explicit contract negotiation or unsupported error |
| `plugins/ghidra/goreveal_ghidra.py` | explicit contract negotiation or unsupported error |
| plugin and schema fixture tests | v1 frozen and v2 explicit compatibility evidence |

## RT1-S0 Detailed Tasks

### Task 1: Add explicit analysis-stage diagnostics

**Files:**
- Create: `schema/diagnostic.go`
- Create: `engine/stages.go`
- Modify: `schema/analysis.go:34-46`
- Modify: `engine/engine.go:23-115`
- Modify: `cmd/goreveal/internal/analyze.go:14-118`
- Test: `engine/engine_test.go`
- Test: `cmd/goreveal/internal/analyze_test.go`
- Test: `schema/export_test.go`

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

- [ ] **Step 4: Introduce injectable stage readers**

Refactor `Analyzer` to carry an unexported reader set:

```go
type stageReaders struct {
    buildInfo  func(string) (schema.BuildInfo, error)
    runtime    func(string) (schema.RuntimeMetadata, error)
    functions  func(string) ([]schema.Function, error)
    types      func(string) ([]schema.Type, error)
    strings    func(string) (recoverystrings.Result, error)
    sourceFiles func(string) ([]string, error)
}
```

`New()` supplies production readers. Tests use `newAnalyzerForTest` in the
`engine` package. Keep `core` independent from engine diagnostics.

- [ ] **Step 5: Write the injected-failure regression test**

Inject a functions reader returning `errors.New("fixture failure")`. Assert:

- `AnalyzeFile` returns a partial analysis and no top-level error;
- diagnostics contain `stage=functions`, `status=failed`;
- no packages or peeling claims are derived from the failed function stage;
- no empty result is presented as `available`.

- [ ] **Step 6: Run the engine test and verify RED**

Run:

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./engine -run TestAnalyzeFileRecordsStageFailure -v
```

Expected: FAIL because errors are still silently discarded.

- [ ] **Step 7: Record every attempted stage**

Replace `if recovered, err := ...; err == nil` branches with stage helpers.
Use explicit `unsupported`/`unavailable` only for recognized sentinel errors;
unknown errors are `failed`. Do not turn a failure into `unavailable` merely to
keep output green.

- [ ] **Step 8: Add command-specific availability checks**

Add a helper in `cmd/goreveal/internal/analyze.go` that rejects a requested
surface whose stage is `failed` or `unsupported`. Preserve the existing
stripped-fixture rule that `inspect types` emits `[]` when no truthful type
surface exists. `inspect runtime` continues to return `unavailable` rather
than invented metadata.

- [ ] **Step 9: Run focused and broad tests**

Run:

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./engine ./cmd/goreveal/internal ./schema -v
make test
```

Expected: PASS; snapshot updates are deferred to RT1-S1.

- [ ] **Step 10: Commit**

```bash
git add schema/diagnostic.go schema/analysis.go engine/stages.go engine/engine.go engine/engine_test.go cmd/goreveal/internal/analyze.go cmd/goreveal/internal/analyze_test.go schema/export_test.go
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
- Test: `deobfuscation/pipeline_test.go`
- Test: `deobfuscation/garble/strings_test.go`
- Test: `schema/export_test.go`
- Test: `schema/export_fixture_test.go`

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

- [ ] **Step 4: Write a reorder-and-segment regression test**

Construct two raw strings with different addresses. Run the garble pass so it
adds segments and reorders findings. Assert every refined string references the
correct raw ID and byte span.

- [ ] **Step 5: Run the regression and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./deobfuscation/... -run 'TestPipelinePreservesRawIDs|TestGarbleSegmentsReferenceRawString' -v
```

Expected: FAIL because refined strings have only `Value`.

- [ ] **Step 6: Add keyed refinement fields**

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
    Value    string         `json:"value"`
    RawStart uint64         `json:"raw_start,omitempty"`
    RawEnd   uint64         `json:"raw_end,omitempty"`
}
```

Give functions/packages/types equivalent `ID` and `RawID` fields. Use
`FindAllStringIndex` in the garble pass. A segment is a separate finding, not
the primary display value.

- [ ] **Step 7: Replace export index lookup with raw-ID lookup**

Build maps from `RawID` to the single `display` refinement. Ignore `segment`
findings for the singular `RefinedValue` field. If multiple display findings
exist for one raw ID, leave the singular field empty and surface a diagnostic;
do not choose by order.

- [ ] **Step 8: Run focused tests and snapshots**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./deobfuscation/... -v
make test
```

Expected: PASS with no address/value cross-binding.

- [ ] **Step 9: Commit**

```bash
git add schema/entity_id.go schema/analysis.go schema/refined.go deobfuscation schema/export_ida.go schema/export_ghidra.go schema/*_test.go
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

- [ ] **Step 5: Prove ambiguity cannot reach auto-accept**

Add a table test that passes ambiguity-containing matches through
`buildTransferCandidates` and `buildAcceptedTransfers`. Assert zero accepted
members from every ambiguity group.

- [ ] **Step 6: Run focused and broad tests**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./storage/diff ./storage/sqlite -v
make test
```

Expected: PASS; existing unambiguous matches remain deterministic.

- [ ] **Step 7: Commit**

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
unchanged until the RT1-S2 location migration; centralize conversion before
comparison.

- [ ] **Step 4: Run tests and commit**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/runtime -v
git add core/runtime/moduledata.go core/runtime/moduledata_test.go
git commit -m "fix(runtime): reject exclusive text endpoint"
```

### Task 5: Close RT1-S0 with clean-room and evidence synchronization

**Files:**
- Modify: `docs/tmp/draft/go-bp.md:342-440`
- Modify: `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- Modify: `docs/plans/2026-03-20-goreveal-functional-assessment.md`
- Modify: `README.md`
- Modify: `AGENTS.md`

- [ ] **Step 1: Replace direct-code-reuse recommendations**

Rewrite the three GoReSym fork/copy recommendations as clean-room behavior
study, upstream-Go primary-source research, fixture manifests, and differential
comparison. Preserve the draft as a draft; do not promote its old advice.

- [ ] **Step 2: Record only landed S0 claims**

Update semantic boundaries and functional assessment after the code commits.
Do not mark S1/S2 planned fields as available.

- [ ] **Step 3: Run documentation consistency checks**

```bash
rg -n 'copy|fork|translate|direct code' docs/tmp/draft/go-bp.md
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

- [ ] **Step 1: Add a failing image-policy test**

Read `Containerfile.dev` and fail on `@latest` or an unversioned pip package.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- python3 -m unittest scripts.dev.test_podman_runner.TestDevImagePolicy -v
```

Expected: FAIL on the current `@latest` and unpinned pip installs.

- [ ] **Step 3: Pin versions and document the update procedure**

Use exact `package==version` and `module@version` declarations. Document how to
inspect installed module versions and require full gates for upgrades.

- [ ] **Step 4: Rebuild and verify versions**

```bash
task build-image
python3 -m scripts.dev.podman_runner exec -- /go/bin/golangci-lint --version
python3 -m scripts.dev.podman_runner exec -- ruff --version
make lint
make test
```

Expected: declared versions and green gates.

- [ ] **Step 5: Commit**

```bash
git add deployments/docker scripts/dev/test_podman_runner.py
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
python3 -m unittest scripts.dev.test_podman_runner -v
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
python3 -m unittest scripts.dev.test_podman_runner -v
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

- [ ] **Step 4: Add sequential Podman CI**

The workflow installs Podman, builds the pinned dev image, then runs in order:

1. `make lint`;
2. `make lint-scripts`;
3. `make test`;
4. `make test-differential`;
5. `make test-snapshots`;
6. `make fuzz`;
7. `make bench`.

Do not run workspace-writing jobs concurrently.

- [ ] **Step 5: Verify locally**

```bash
make lint
make lint-scripts
make test
make test-differential
make test-snapshots
make fuzz
make bench
git diff --check
```

Expected: all exit `0` and all new snapshots are reviewed.

- [ ] **Step 6: Commit**

```bash
git add tests/snapshots corpus/fixtures .github/workflows/ci.yml README.md
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
python3 -m unittest discover -s scripts -p 'test_*.py'
git diff --check
git add docs/evidence docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md
# Add LICENSE and scripts/evidence only when the preceding decision/validator
# steps actually created them.
git commit -m "docs(evidence): record release and IDA baseline decisions"
```

Omit paths that do not exist because the license decision is still blocked;
do not create a placeholder license.

## RT1-S2 Detailed Tasks

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

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/buildinfo ./cmd/goreveal/internal -v
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
- Update: fixture manifests

- [ ] **Step 1: Write digest/architecture/build-ID tests**

Cover a known byte digest, ELF amd64, PE amd64, Mach-O amd64, absent build ID,
truncated notes, and a non-Go binary.

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
make test
```

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

- [ ] **Step 1: Write table tests before implementation**

Cover ELF ET_EXEC, PIE ELF, PE image base plus section RVA, Mach-O VM address,
section start/end, rebased VA, unmapped gaps, and arithmetic overflow.

- [ ] **Step 2: Run and verify RED**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/location -v
```

Expected: package does not exist.

- [ ] **Step 3: Implement the ADR model**

Use half-open ranges throughout. Keep file offset optional. Never infer image
base from `.text`. A mapping failure is explicit and carries section/segment
context when available.

- [ ] **Step 4: Replace duplicate address inference**

Adapt pclntab/runtime projections to consume the location mapper instead of
recomputing PE or ELF translations independently. Preserve raw runtime
evidence fields until a separate compatibility migration removes them.

- [ ] **Step 5: Verify and commit**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./core/location ./core/pclntab ./core/runtime -v
make test
git add schema/location.go core/location core/pclntab core/runtime schema/analysis.go
git commit -m "feat(core): add explicit address location mapping"
```

### Task 15: Publish explicit v2 exports without changing v1

**Files:**
- Create: `schema/export_ida_v2.go`
- Create: `schema/export_ida_v2_test.go`
- Create if ADR promotes it: `schema/export_ghidra_v2.go`
- Modify: `cmd/goreveal/internal/export_ida.go`
- Modify: `cmd/goreveal/main.go`
- Modify: `plugins/ida/goreveal_ida.py`
- Modify: `plugins/ghidra/goreveal_ghidra.py`
- Modify: plugin tests and fixtures

- [ ] **Step 1: Freeze unchanged v1 bytes**

Add or retain a golden v1 fixture and assert the constructor emits contract
`goreveal.export.ida/v1` with the old field shape.

- [ ] **Step 2: Write v2 contract tests and verify RED**

Require binary identity, analyzer identity, address-space declaration,
function IDs/locations, build provenance, and diagnostics. The payload has no
`artifact_sha256` field.

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema -run 'TestIDAExportV1Frozen|TestIDAExportV2' -v
```

Expected: v1 passes; v2 test fails because no constructor exists.

- [ ] **Step 3: Implement a separate v2 type and constructor**

Do not add v2-only populated fields to the v1 wire struct. Sort all exported
collections deterministically.

- [ ] **Step 4: Add explicit CLI selection**

Support:

```text
goreveal export ida --contract goreveal.export.ida/v1 <binary>
goreveal export ida --contract goreveal.export.ida/v2 <binary>
```

Keep the no-flag default on v1 during the compatibility window.

- [ ] **Step 5: Make consumers explicit**

Existing thin consumers either parse a tested v2 subset or return an exact
unsupported-contract error. They must not accidentally treat v2 as v1.

- [ ] **Step 6: Verify contracts and plugins**

```bash
python3 -m scripts.dev.podman_runner exec -- /usr/local/go/bin/go test ./schema ./cmd/goreveal/... -v
make test-plugins
make test-snapshots
```

Expected: v1 and v2 fixtures pass; unknown contracts are rejected.

- [ ] **Step 7: Commit**

```bash
git add schema cmd/goreveal plugins corpus/fixtures
git commit -m "feat(export): publish identity-bound IDA v2 contract"
```

### Task 16: Decide JSON versus streaming bundle from evidence

**Files:**
- Create: `schema/export_ida_v2_bench_test.go`
- Create: `docs/evidence/rt1-s2-artifact-envelope.md`
- Modify only if threshold is exceeded: v2 export implementation and CLI

- [ ] **Step 1: Benchmark a large deterministic function export**

Record encoded bytes, elapsed time, allocations, and peak RSS for a synthetic
458,600-function payload plus the measured real-binary artifact when allowed.

- [ ] **Step 2: Apply the decision rule**

- keep single JSON if peak RSS is at most `2x` artifact bytes and both provider
  and reference consumer stay within the declared large-binary envelope;
- otherwise use a small manifest plus deterministic NDJSON record streams;
- do not invent chunking before the measurement.

- [ ] **Step 3: Freeze the selected envelope and rerun v2 tests**

The exact transport bytes remain what the consumer hashes.

- [ ] **Step 4: Commit**

```bash
git add schema docs/evidence/rt1-s2-artifact-envelope.md cmd/goreveal
git commit -m "perf(export): freeze RT1 artifact envelope"
```

### Task 17: Close RT1-S2 evidence and compatibility

**Files:**
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/architecture/2026-03-20-goreveal-semantic-claim-boundaries.md`
- Modify: `docs/plans/2026-03-19-goreveal-feature-map.md`
- Create: `docs/evidence/rt1-s2-closure.md`

- [ ] **Step 1: Run all S2 mismatch and round-trip fixtures**

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
git commit -m "docs: close RT1-S2 identity contract"
```

## RT1-S3 Promotion and Cross-Repository Plan Gate

### Task 18: Create and review the idacli implementation plan

**GoREveal inputs required before promotion:**

- frozen `goreveal.export.ida/v2` fixture;
- exact provider artifact digest fixture;
- address mapping fixtures including rebased and mismatch cases;
- function action fixture with missing, matching, named, unnamed, and boundary
  conflict cases;
- S1 forced Golang-plugin baseline.

**Expected idacli plan files:**

- Spec: `/opt/projects/repositories/idacli/docs/superpowers/specs/2026-07-22-go-provider-preview-apply-design.md`
- Plan: `/opt/projects/repositories/idacli/docs/superpowers/plans/2026-07-22-go-provider-preview-apply.md`

**Expected idacli implementation ownership:**

- `src/common/task_request.*`: top-level preview/apply request fields;
- `src/tasks/task_go_preview.*`: read-only identity/mapping/action comparison;
- `src/tasks/task_go_apply.*`: exact-preview application;
- `src/tasks/task_factory.cpp` and task enum: registered executable tasks;
- `src/cli/idb_isolation.*`: retained private IDB copy/checkpoint;
- focused C++ unit/fixture tests plus licensed integration evidence.

- [ ] **Step 1: Re-audit current idacli HEAD and instructions**

Do not assume the July draft matches the current task parser or runtime
discovery work.

- [ ] **Step 2: Write the idacli spec around exact preview bytes**

`go-preview` writes an immutable `idacli.go-preview/v1` artifact without a
self-digest. `go-apply` receives preview path plus expected
`sha256:<lowercase-hex>`, rehashes exact bytes, validates embedded provider and
IDB identities, and applies the embedded actions.

- [ ] **Step 3: Keep first apply conservative**

Allow create missing function, set default name, and additive comment only.
Skip boundary conflict and user name. Forbid `del_func` and automatic resize.

- [ ] **Step 4: Review and approve the separate idacli plan**

Use one plan-document reviewer and run idacli's own gates. No GoREveal commit
contains idacli source changes.

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
- closure evidence links every claimed RT1-S0/S1/S2 gate.

Before merging Horizon A, request independent code review against the RT1
design and this plan. Fix every Critical and Important issue before proceeding
to RT1-S3.
