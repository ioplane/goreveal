# GoREveal Standalone Qualification and R1 Release Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Qualify the completed standalone provider, make its CPU and performance decisions observable, produce reproducible release assets, and actually publish `RT1-R1` before any IDA implementation starts.

**Architecture:** `RT1-S2C` measures the already-frozen `RT1-S2B` behavior without reopening recovery or export semantics. Diagnostics and evidence orchestration remain outside `core`; a measured owner package may receive a scalar optimization and, only after explicit gates pass, an architecture-specific implementation behind immutable runtime dispatch. `RT1-R1` is a publication milestone over byte-reproducible S2C artifacts, not another feature sprint.

**Tech Stack:** Go 1.26, `golang.org/x/sys/cpu`, Go benchmarks/fuzzing/pprof, Python 3 evidence and release tooling, GNU `time -v`, SPDX JSON SBOM, Podman, GitHub Actions, GitHub Releases.

---

## Authority, Preconditions, and Stop Rules

Use these sources in order:

1. `../specs/2026-07-22-goreveal-standalone-release-ida-bootstrap-design.md`;
2. `../specs/2026-07-22-goreveal-rt1-product-design.md`;
3. `2026-07-22-goreveal-rt1-horizon-a.md` Tasks 0-17;
4. `AGENTS.md` and the architecture contracts it names.

Required execution skills are `@goreveal-navigation`, `@goreveal-perf-simd`,
`@goreveal-release-ops`, `@goreveal-export-contracts`,
`@goreveal-corpus-validation`, `@goreveal-differential-testing`,
`@goreveal-doc-sync`, `@test-driven-development`, and
`@verification-before-completion`.

Do not start Task 1 until `docs/evidence/rt1-s2b-closure.md` exists and proves:

- S0-S2B gates passed sequentially in the pinned Podman environment;
- v1 exact bytes remain frozen;
- the selected v2 envelope and reference verifier are public and measured;
- ELF, PIE ELF, PE, and Mach-O location fixtures pass;
- wrong binary/base/artifact/component/location cases reject;
- the registered real 410 MB-class input was measured.

If this precondition is false, stop and finish the Horizon A plan. Do not adapt
this plan to the current pre-S0 code.

Hard stop rules:

- profiling may produce a negative SIMD decision; that is a successful S2C result;
- no generic `simd/` package is authorized;
- no parser, address mapping, schema, confidence, provenance, or recovery decision becomes SIMD-dependent;
- no `GOAMD64=v3` default artifact is authorized;
- no hardware threshold is weakened after seeing results;
- unresolved licensing/compliance keeps `RT1-R1` open and blocks every IDA implementation task;
- R1 closes only after the remote release assets are downloaded and independently verified.

## File Responsibility Map

### Always-created S2C files

| File | Single responsibility |
| --- | --- |
| `schema/cpu_diagnostics.go` | stable machine-readable CPU, build, kernel, and run-class contract |
| `schema/cpu_diagnostics_test.go` | exact JSON and empty-collection contract |
| `internal/cpudiag/report.go` | pure deterministic report builder over injected inputs |
| `internal/cpudiag/features_amd64.go` | dispatch-relevant AMD64 feature inventory |
| `internal/cpudiag/features_arm64.go` | ASIMD plus inventory-only SVE/SVE2 facts |
| `internal/cpudiag/features_other.go` | explicit scalar-only fallback |
| `internal/cpudiag/report_test.go` | feature, selection, and rejection table tests |
| `cmd/goreveal/internal/diagnostics_cpu.go` | JSON CLI projection only |
| `cmd/goreveal/internal/diagnostics_cpu_test.go` | command bytes and exit behavior |
| `bench/standalone_bench_test.go` | named engine/provider/verifier benchmarks |
| `bench/profile/main.go` | evidence-only pprof driver over production entrypoints |
| `bench/artifacts.lock.json` | exact registered small/medium/large inputs |
| `bench/README.md` | repetitions, cache, host, and raw-evidence protocol |
| `scripts/evidence/measure_standalone_performance.py` | execute and summarize registered runs |
| `scripts/evidence/test_measure_standalone_performance.py` | subprocess, locale, and metric negative tests |
| `scripts/evidence/validate_rt1_s2c.py` | validate run records and promotion thresholds |
| `scripts/evidence/test_validate_rt1_s2c.py` | invalid and negative-result fixtures |
| `docs/evidence/rt1-s2c/result.schema.json` | machine-checkable run and decision schema |
| `docs/evidence/rt1-s2c/README.md` | reproducible protocol and external-raw-artifact rules |
| `docs/evidence/rt1-s2c/decision.md` | measured per-hotspot keep/reject record |
| `docs/evidence/rt1-s2c/closure.md` | trace every S2C obligation to evidence |

### Release-candidate and publication files

| File | Single responsibility |
| --- | --- |
| `schema/version.go` | stable version/provenance JSON contract |
| `cmd/goreveal/internal/version.go` | `goreveal version --json` projection |
| `scripts/release/build.py` | deterministic archives, OCI layouts, checksums, and SBOM invocation |
| `scripts/release/test_build.py` | reproducibility and normalized-metadata tests |
| `scripts/release/verify.py` | offline whole-release verifier |
| `scripts/release/test_verify.py` | changed/missing/extra/wrong-subject negative tests |
| `scripts/release/verify_publication.py` | remote inventory and digest verification |
| `scripts/release/test_verify_publication.py` | remote manifest negative tests |
| `docs/release/standalone-operator.md` | one-page standalone workflows |
| `docs/release/support-matrix.schema.json` | separates execution hosts from analyzed targets |
| `docs/release/support-matrix.json` | evidence-derived supported matrix |
| `docs/release/known-limitations.md` | explicit unqualified and unsupported states |
| `docs/release/compatibility.md` | v1 freeze and selected v2 references |
| `docs/release/compliance.md` | owned legal, license, dependency, and publication gates |
| `docs/evidence/rt1-r1/release-manifest.schema.json` | local candidate provenance contract |
| `docs/evidence/rt1-r1/publication.schema.json` | actual remote-publication contract |
| `.github/workflows/release.yml` | publish the preverified bytes for an approved tag |

### Conditional owner-package files

Create these only when Task 4 names the corresponding hotspot:

- diff set operations: `storage/diff/setops.go`, `setops_test.go`,
  `setops_fuzz_test.go`, `setops_bench_test.go`;
- repeated byte scanning: `core/internal/bytescan/scan.go`,
  `scan_generic.go`, `scan_test.go`, `scan_fuzz_test.go`,
  `scan_bench_test.go`;
- string classification: add the equivalent seam locally under `core/strings/`.

Architecture-specific files (`*_amd64.go`, `*_amd64.s`, `*_arm64.go`,
`*_arm64.s`) require a target-specific addendum reviewed after scalar profiling.
They are not authorized by this base plan alone.

## RT1-S2C Detailed Tasks

### Task 1: Enforce the closed-S2B dependency and evidence schema

**Files:**
- Create: `docs/evidence/rt1-s2c/result.schema.json`
- Create: `scripts/evidence/validate_rt1_s2c.py`
- Create: `scripts/evidence/test_validate_rt1_s2c.py`
- Modify: `scripts/dev/podman_runner.py`
- Modify: `scripts/dev/test_podman_runner.py`
- Modify: `Makefile`
- Modify: `Taskfile.yml`

- [ ] **Step 1: Write the dependency and validator tests**

The test fixture must reject an absent S2B closure, a registered input without
SHA-256/size/class, a run without commit/command/CPU/raw-output reference, and a
positive optimization claim without thresholds. It must accept a fully explicit
negative decision:

```python
decision = {
    "candidate": "none",
    "outcome": "reject",
    "reason": "no_product_hotspot",
    "release_implementation": "scalar",
}
self.assertEqual(validate_decision(decision), [])
```

- [ ] **Step 2: Run RED inside the dev container**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.evidence.test_validate_rt1_s2c -v
```

Expected: FAIL because the validator and schema do not exist.

- [ ] **Step 3: Implement the strict validator and JSON Schema**

The validator reads exact paths passed by the caller; it never searches the host.
Use lower-case `sha256:<64 hex>` digests, positive byte sizes, enumerated
`small|medium|large`, `correctness-only|performance-qualification` run classes,
and `keep|reject` decisions. A keep result requires all six approved thresholds;
a reject result requires a stable reason and selected scalar implementation.

- [ ] **Step 4: Add runner tasks and read-only artifact mounting**

Add `perf-smoke`, `perf-profile`, and `perf-evidence` as sequential task plans.
External artifacts mount at `/artifacts:ro`; outputs remain under the bind-mounted
workspace. The runner must require explicit `--artifacts-root` for profile/evidence
and must not infer `/srv`, `$HOME`, or another broad host path.

- [ ] **Step 5: Verify GREEN and wrapper parity**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.evidence.test_validate_rt1_s2c \
  scripts.dev.test_podman_runner -v
make -n perf-smoke
task --summary perf-smoke
```

Expected: all tests pass; Make and Task delegate to the runner and contain no
independent benchmark command.

- [ ] **Step 6: Commit**

```bash
git add docs/evidence/rt1-s2c scripts/evidence scripts/dev Makefile Taskfile.yml
git commit -m "test(perf): enforce RT1-S2C evidence gates"
```

### Task 2: Add deterministic CPU and dispatch diagnostics

**Files:**
- Create: `schema/cpu_diagnostics.go`
- Create: `schema/cpu_diagnostics_test.go`
- Create: `internal/cpudiag/report.go`
- Create: `internal/cpudiag/report_test.go`
- Create: `internal/cpudiag/features_amd64.go`
- Create: `internal/cpudiag/features_arm64.go`
- Create: `internal/cpudiag/features_other.go`
- Create: `cmd/goreveal/internal/diagnostics_cpu.go`
- Create: `cmd/goreveal/internal/diagnostics_cpu_test.go`
- Modify: `cmd/goreveal/main.go`
- Modify: `cmd/goreveal/main_test.go`
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Write exact contract tests**

Use this minimum canonical shape; arrays are sorted and always encode as `[]`:

```go
type CPUDiagnostics struct {
    Contract       string               `json:"contract"`
    GOOS           string               `json:"goos"`
    GOARCH         string               `json:"goarch"`
    GoVersion      string               `json:"go_version"`
    GoExperiments  []string             `json:"go_experiments"`
    Features       []CPUFeature         `json:"features"`
    Implementations []KernelSelection   `json:"implementations"`
    RunClass       string               `json:"run_class"`
    Virtualized    *bool                `json:"virtualized"`
}
```

Table cases cover AMD64 AVX2-only, AMD64 AVX-512, ARM64 ASIMD with inventory-only
SVE/SVE2, and an unknown architecture with scalar selection. Test that unavailable
optimized paths include a reason.

- [ ] **Step 2: Run RED**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./schema ./internal/cpudiag ./cmd/goreveal/... \
  -run 'TestCPU(Diagnostics|Report)|TestDiagnosticsCPU' -v
```

Expected: FAIL because the packages and command do not exist.

- [ ] **Step 3: Implement an immutable injected report builder**

`internal/cpudiag.Build(inputs)` receives GOOS/GOARCH/version/experiments,
detected feature booleans, compiled implementations, selection results, run class,
and optional virtualization evidence. Production inputs use `x/sys/cpu`; tests do
not mutate package globals. Promote `golang.org/x/sys` to a direct dependency.

- [ ] **Step 4: Implement `goreveal diagnostics cpu --json`**

The command writes exactly one JSON object plus LF. Unknown flags and a requested
forced implementation unsupported by the current CPU return non-zero before any
kernel call. Diagnostics never claim a virtualized host is a performance host.

- [ ] **Step 5: Verify CLI and cross-build fallback**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./schema ./internal/cpudiag ./cmd/goreveal/... -v
python3 -m scripts.dev.podman_runner task build-dev-bin
python3 -m scripts.dev.podman_runner exec -- \
  /workspace/.tmp/goreveal diagnostics cpu --json
python3 -m scripts.dev.podman_runner exec -- bash -lc \
  'GOOS=linux GOARCH=amd64 go build ./cmd/goreveal && GOOS=linux GOARCH=arm64 go build ./cmd/goreveal'
```

Expected: tests and builds pass; JSON selects scalar until a qualified kernel exists.

- [ ] **Step 6: Commit**

```bash
git add schema internal/cpudiag cmd/goreveal go.mod go.sum
git commit -m "feat(diagnostics): report CPU dispatch capabilities"
```

### Task 3: Build the reproducible standalone performance harness

**Files:**
- Create: `bench/standalone_bench_test.go`
- Create: `bench/profile/main.go`
- Create: `bench/artifacts.lock.json`
- Create: `bench/README.md`
- Create: `scripts/evidence/measure_standalone_performance.py`
- Create: `scripts/evidence/test_measure_standalone_performance.py`
- Modify: `scripts/dev/podman_runner.py`
- Modify: `scripts/dev/test_podman_runner.py`

- [ ] **Step 1: Write failing parser and subprocess tests**

Cover `LC_ALL=C` GNU `time -v`, non-zero children, missing output, wrong artifact
digest/size, partial repetitions, p50/p95 calculation, and external profile records.
The script must not parse localized labels.

- [ ] **Step 2: Run RED**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.evidence.test_measure_standalone_performance -v
```

Expected: FAIL because the measurement module is absent.

- [ ] **Step 3: Add named production-path benchmarks**

Benchmark the canonical engine analysis, selected v2 provider, and reference verifier.
Use registered fixture files and the S2B large synthetic constructor; do not create a
parallel analysis implementation. Every benchmark validates output digest/count before
reporting timing.

- [ ] **Step 4: Add the evidence-only pprof driver**

`bench/profile` accepts one registered artifact, output directory, and operation.
It writes CPU/heap profiles and a JSON manifest containing tool/input identities and
exit state. It is not copied into builder or release images.

- [ ] **Step 5: Implement measurement and lock validation**

For each registered class, capture command, commit, binary/tool identity, CPU report,
wall time, peak RSS, output bytes/digest/count, cache policy, run class, and raw-output
hash/reference. Use 20 measured repetitions after declared warmups for p95 claims.
Large or sensitive raw profiles remain outside git.

- [ ] **Step 6: Run the bounded KVM smoke**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.evidence.test_measure_standalone_performance -v
make perf-smoke
```

Expected: one fixture executes the real CLI provider/verifier and produces a
`correctness-only` record. No performance threshold is evaluated on KVM.

- [ ] **Step 7: Commit**

```bash
git add bench scripts/evidence scripts/dev
git commit -m "test(perf): add standalone profiling harness"
```

### Task 4: Profile first and freeze the scalar/SIMD decision

**Files:**
- Create: `docs/evidence/rt1-s2c/decision.md`
- Create from runs: `docs/evidence/rt1-s2c/runs/<run-id>/host.json`
- Create from runs: `docs/evidence/rt1-s2c/runs/<run-id>/workflow.json`
- Create from runs: `docs/evidence/rt1-s2c/runs/<run-id>/bench.txt`
- Modify: `docs/evidence/rt1-s2c/README.md`

- [ ] **Step 1: Pre-register the host and measurement protocol**

Record the bare-metal x86-64 host, bare-metal ARM64/NEON host, AVX2-only host,
power/governor policy, warm/cold rule, 20 repetitions, artifact locks, and commands
before collecting performance data. KVM Xeon Gold 6348 is correctness-only.

- [ ] **Step 2: Capture the unmodified closed-S2B scalar baseline**

```bash
python3 -m scripts.dev.podman_runner task perf-profile \
  --artifacts-root /srv/goreveal-rt1 \
  --output-root .tmp/evidence/rt1-s2c/profile \
  --run-class correctness-only
```

Run the equivalent registered command on both bare-metal performance hosts.
Expected: profiles and small/medium/large workflow records validate against the schema.

- [ ] **Step 3: Rank product hotspots**

Compare cumulative CPU and wall contribution, allocation/RSS impact, and stdlib use.
Reject custom SIMD immediately when no bounded data-parallel owner operation accounts
for at least 10% of the registered large workflow.

- [ ] **Step 4: Compare the best scalar and standard-library alternatives**

For a selected owner operation, benchmark current behavior against the simplest
correct scalar algorithm and relevant `bytes`, `strings`, `crypto`, or `hash`
primitive. Never use a deliberately weak scalar control.

- [ ] **Step 5: Record one immutable decision branch**

Choose exactly one:

- `reject=no_product_hotspot`: ship the existing scalar/stdlib path;
- `keep=optimized_scalar`: implement Task 5, then re-profile;
- `candidate=architecture_specific`: Task 5 must land first, then write and review
  a target-specific addendum naming exact owner files, instruction set, dispatch,
  tests, and hardware. Do not create architecture files from this plan alone.

- [ ] **Step 6: Validate and commit the decision**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 scripts/evidence/validate_rt1_s2c.py \
  --schema docs/evidence/rt1-s2c/result.schema.json \
  --evidence docs/evidence/rt1-s2c
git add docs/evidence/rt1-s2c
git commit -m "docs(perf): record standalone hotspot decision"
```

Expected: negative decisions pass; unsupported positive claims fail.

### Task 5: Implement only the selected scalar owner optimization

**Files:**
- Create/Modify: only the exact owner-package files selected by Task 4
- Test: matching `*_test.go`, `*_fuzz_test.go`, and `*_bench_test.go`
- Modify: `docs/evidence/rt1-s2c/decision.md`

Skip this task when Task 4 selected the existing scalar/stdlib implementation.

- [ ] **Step 1: Write scalar-oracle equivalence tests**

Cover empty input, one byte, boundary sizes around vector/word widths, all alignments,
tails, malformed data, overlapping patterns, and deterministic errors. The optimized
function and reference function must return identical values and errors.

- [ ] **Step 2: Run RED in the owner package**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./<owner-package> -run 'Test<Operation>Equivalence' -v
```

Expected: FAIL because the selected seam or optimized scalar implementation is absent.

- [ ] **Step 3: Add fuzz and benchmark tests before implementation**

The fuzz target always compares against the reference oracle. The benchmark reports
throughput and allocations for reference, current, and optimized scalar paths.

- [ ] **Step 4: Implement the smallest scalar improvement**

Keep reference truth callable by tests. Do not change schema bytes, error classes,
ordering, confidence, provenance, locations, or export envelope.

- [ ] **Step 5: Verify equivalence and improvement**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./<owner-package> -v
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./<owner-package> -run '^$' \
  -bench 'Benchmark<Operation>' -benchmem -count=10 -benchtime=3s
make fuzz
make perf-smoke
```

Expected: exact equivalence; benchmark evidence recorded. Re-profile the large
workflow before requesting any architecture-specific addendum.

- [ ] **Step 6: Commit**

```bash
git add <owner-package> docs/evidence/rt1-s2c/decision.md
git commit -m "perf(<owner>): optimize measured scalar path"
```

### Task 6: Qualify hardware, dispatch, and the final performance decision

**Files:**
- Modify: `schema/cpu_diagnostics.go`
- Modify: `internal/cpudiag/report.go`
- Modify: selected owner-package files only if an approved addendum exists
- Modify: `docs/evidence/rt1-s2c/decision.md`
- Create from runs: additional `docs/evidence/rt1-s2c/runs/<run-id>/...`

- [ ] **Step 1: Test forced selection without package-global mutation**

An injected immutable dispatcher must support forced scalar and forced named modes.
Forcing an unsupported implementation returns an error before invocation; normal mode
always selects a safe compiled implementation.

- [ ] **Step 2: Prove scalar and unsupported-feature paths**

Run generic cross-builds, injected feature tables, KVM functional dispatch, and the
real AVX2-only host. ARM64 runs must prove the scalar/ASIMD-selected behavior actually
shipped for that artifact.

- [ ] **Step 3: Collect representative bare-metal records**

```bash
python3 -m scripts.dev.podman_runner task perf-evidence \
  --artifacts-root /srv/goreveal-rt1 \
  --output-root .tmp/evidence/rt1-s2c/baremetal-<host-id> \
  --run-class performance-qualification
```

Required: bare-metal x86-64 and bare-metal ARM64/NEON. AVX-512 needs a same-host
AVX2 comparison and must still be tested for fallback on an AVX2-only host.

- [ ] **Step 4: Apply the pre-registered gates once**

Promotion requires hotspot >=10%, microbenchmark >=1.5x best scalar, large workflow
p50 or p95 >=10% faster, deterministic bytes/errors identical, scalar/unsupported
paths green, and unrelated supported workflows no more than 5% slower. Otherwise
record a reject reason and ship scalar.

- [ ] **Step 5: Commit only the qualified result**

```bash
git add schema internal/cpudiag <owner-package-if-any> docs/evidence/rt1-s2c
git commit -m "perf: qualify standalone dispatch decision"
```

Expected: diagnostics name the selected implementation and why alternatives were
unavailable or rejected. Experimental `GOEXPERIMENT=simd` results, if any, remain a
separate non-release evidence lane.

### Task 7: Freeze release contracts, tooling, and publication workflow

**Files:**
- Create: `schema/version.go`
- Create: `schema/version_test.go`
- Create: `cmd/goreveal/internal/version.go`
- Create: `cmd/goreveal/internal/version_test.go`
- Modify: `cmd/goreveal/main.go`
- Modify: `internal/version/version.go`
- Create after explicit owner decision: `VERSION`
- Create after explicit owner decision: `LICENSE`
- Create: `CHANGELOG.md`
- Modify: `deployments/docker/Containerfile.dev`
- Modify: `deployments/docker/Containerfile.builder`
- Modify: `deployments/docker/Containerfile.release`
- Modify: `deployments/docker/toolchain.lock.json`
- Create: `scripts/release/build.py`
- Create: `scripts/release/test_build.py`
- Create: `scripts/release/verify.py`
- Create: `scripts/release/test_verify.py`
- Create: `scripts/release/verify_publication.py`
- Create: `scripts/release/test_verify_publication.py`
- Create: `docs/release/standalone-operator.md`
- Create: `docs/release/support-matrix.schema.json`
- Create: `docs/release/support-matrix.json`
- Create: `docs/release/compatibility-manifest.schema.json`
- Create: `docs/release/known-limitations.md`
- Create: `docs/release/compatibility.md`
- Create: `docs/release/compliance.md`
- Create: `docs/evidence/rt1-r1/release-manifest.schema.json`
- Create: `docs/evidence/rt1-r1/publication.schema.json`
- Create: `.github/workflows/release.yml`
- Modify: `scripts/dev/podman_runner.py`
- Modify: `scripts/dev/test_podman_runner.py`
- Modify: `Makefile`
- Modify: `Taskfile.yml`

- [ ] **Step 1: Obtain the version, license, and compliance decisions**

The release owner selects the exact first public version and records it in `VERSION`
before candidate construction. Legal supplies the actual `LICENSE` and closes the
owned compliance checklist. If either decision is absent, stop release contract work
with `RT1-R1: blocked`; do not create placeholder text or a provisional version.

- [ ] **Step 2: Write version and provenance RED tests**

`goreveal version --json` must report semantic version, source commit, source epoch,
Go version, GOOS/GOARCH, and dirty state. It must not embed wall-clock `now`.

- [ ] **Step 3: Write contract-first release builder/verifier RED tests**

Test normalized sorted paths, modes, uid/gid 0, `SOURCE_DATE_EPOCH`, gzip without
current timestamps, no checksum self-reference, exact SBOM subject, no missing/extra
assets, changed-byte rejection, and strict validation against the support,
compatibility, candidate, and publication schemas created in this task. Two isolated
builds must hash identically.

- [ ] **Step 4: Write publication staging/verifier RED tests**

Reject wrong source SHA/version, draft release, missing/extra assets, changed checksum,
wrong SBOM subject, absent support matrix/limitations, a staging artifact produced by
another workflow/source, and any attempt to publish bytes not present in the approved
candidate manifest.

- [ ] **Step 5: Run RED**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  /usr/local/go/bin/go test ./schema ./cmd/goreveal/... -run 'TestVersion' -v
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.release.test_build scripts.release.test_verify \
  scripts.release.test_verify_publication -v
```

Expected: FAIL because the command and release modules do not exist.

- [ ] **Step 6: Implement version output and pin every release input**

Use the S1 tool lock; add digest-pinned builder/release bases and one pinned SPDX
SBOM tool. Correct OCI source labels to the actual repository. Treat a Go module-path
migration as a separate compatibility decision; do not silently rewrite it here.

- [ ] **Step 7: Implement deterministic multi-arch packaging from frozen contracts**

Build broad-baseline `linux/amd64` and `linux/arm64` CGO-free CLI archives and OCI
archives. Populate SPDX JSON per binary/image plus schema-validated support,
compatibility, and release manifests, then `SHA256SUMS`. Host availability is
distinct from formats and architectures that GoREveal analyzes.

The exact staged asset inventory is:

```text
goreveal_<version>_linux_amd64.tar.gz
goreveal_<version>_linux_arm64.tar.gz
goreveal_<version>_linux_amd64.oci.tar
goreveal_<version>_linux_arm64.oci.tar
goreveal_<version>_linux_amd64.spdx.json
goreveal_<version>_linux_arm64.spdx.json
goreveal_<version>_contracts.tar.gz
support-matrix.json
compatibility-manifest.json
release-manifest.json
known-limitations.md
SHA256SUMS
```

The contracts archive contains the frozen v1 bytes, selected v2 positive fixtures
for ELF/PIE ELF/PE/Mach-O, all published v2 mismatch fixtures, detached digests, and
a manifest binding each file to the R1 source and contract. `SHA256SUMS` does not
hash itself.

- [ ] **Step 8: Implement immutable GitHub staging and publication**

Pin actions by full commit SHA. One `workflow_dispatch` mode builds twice from an
explicit clean source SHA, verifies byte equality, and uploads exactly one immutable
Actions artifact named `rt1-r1-candidate-<source-sha>` containing `dist/` and the
candidate manifest. A separate publish mode accepts the approved candidate workflow
run ID plus manifest digest, downloads that exact Actions artifact, re-verifies it,
and uploads those existing bytes to the annotated tag release. Publish mode never
rebuilds.

- [ ] **Step 9: Add wrapper targets and verify tooling locally**

```bash
make release-candidate SOURCE_DATE_EPOCH=<source-commit-epoch>
make release-repro SOURCE_DATE_EPOCH=<same-epoch>
make verify-release DIST=dist
```

Expected `dist/` contains Linux amd64/arm64 archives, OCI archives, SPDX files,
support/compatibility/release manifests, and checksums; both builds are byte-identical.

- [ ] **Step 10: Commit every release input before final qualification**

```bash
git add VERSION LICENSE CHANGELOG.md schema cmd/goreveal internal/version \
  deployments/docker scripts/release scripts/dev docs/release \
  docs/evidence/rt1-r1 .github/workflows/release.yml Makefile Taskfile.yml
git commit -m "build(release): freeze standalone release pipeline"
```

Expected: publication workflow, schemas, version, license, tooling, and release docs
all exist in history before the future tagged source commit. Local `dist/` is only a
test output and is not the future approved candidate.

### Task 8: Close S2C and freeze the final candidate source commit

**Files:**
- Modify: `CHANGELOG.md`
- Modify: `docs/release/standalone-operator.md`
- Modify: `docs/release/support-matrix.json`
- Modify: `docs/release/known-limitations.md`
- Modify: `docs/release/compatibility.md`
- Modify: `docs/release/compliance.md`
- Create: `docs/evidence/rt1-s2c/closure.md`
- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/plans/2026-03-19-goreveal-feature-map.md`

- [ ] **Step 1: Finalize evidence-derived operator and support docs**

The one-page workflow covers analyze, inspect, SQLite persist, diff, verify, and
v1/v2 export without IDA, Ghidra, idacli, IDAPython, or plugins. The matrix separately
states release host GOOS/GOARCH and analyzable format/arch/Go/stripped/protected claims.

- [ ] **Step 2: Run the complete S2C gate sequentially**

```bash
make lint
make lint-scripts
make test
make test-differential
make test-differential-report
make test-plugins
make test-snapshots
make fuzz
make bench
make perf-smoke
git diff --check
```

Expected: all exit 0, no skip/no-op gates, v1 bytes unchanged, v2 verification green,
and standalone smoke proves every claim without an RE-tool integration. Candidate
bytes are deliberately not approved yet because the final source commit does not
exist until Step 4.

- [ ] **Step 3: Validate closure**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 scripts/evidence/validate_rt1_s2c.py \
  --schema docs/evidence/rt1-s2c/result.schema.json \
  --evidence docs/evidence/rt1-s2c
```

- [ ] **Step 4: Commit the final candidate source**

```bash
git add CHANGELOG.md docs README.md AGENTS.md
git commit -m "docs(release): close standalone qualification evidence"
git status --short
```

Expected: clean worktree. This commit, after merge, is the only source SHA eligible
for candidate staging. Do not modify tracked source or release docs between staging
and tag creation. `RT1-S2C` is closed; `RT1-R1` remains open.

## RT1-R1 Publication Milestone

### Task 9: Publish and independently verify the exact candidate bytes

**Files:**
- Use unchanged: `.github/workflows/release.yml`
- Use unchanged: `scripts/release/verify_publication.py`
- Use unchanged: `docs/evidence/rt1-r1/publication.schema.json`
- Create from staged candidate, after publication: `docs/evidence/rt1-r1/candidate.json`
- Create after publication: `docs/evidence/rt1-r1/<tag>.json`
- Modify after remote verification: `README.md`
- Modify after remote verification: `AGENTS.md`
- Modify after remote verification: `.agents/skills/goreveal-navigation/SKILL.md`
- Modify after remote verification: `skills/goreveal-navigation/SKILL.md`

- [ ] **Step 1: Merge Task 8 and identify the exact clean source SHA**

Confirm the merged SHA contains `.github/workflows/release.yml`, `VERSION`, `LICENSE`,
all release schemas/docs, S2C closure, and the tested build/verifier code. Record its
commit epoch. No later commit is eligible without rerunning S2C.

- [ ] **Step 2: Stage the candidate at that SHA in GitHub Actions**

Create and push a candidate branch that points exactly at `<source-sha>`; the workflow
still checks out and verifies the explicit SHA rather than trusting a movable branch.

```bash
git branch release/rt1-r1-candidate <source-sha>
git push origin release/rt1-r1-candidate
gh workflow run release.yml --ref release/rt1-r1-candidate \
  -f mode=candidate -f source_sha=<source-sha>
gh run watch <candidate-run-id> --exit-status
gh run download <candidate-run-id> -n rt1-r1-candidate-<source-sha> \
  -D .tmp/candidate/<source-sha>
```

Expected: the workflow builds twice from the same clean commit epoch, proves identical
bytes, verifies the whole release, and stages one immutable Actions artifact.

- [ ] **Step 3: Independently verify the downloaded staged candidate**

```bash
python3 -m scripts.dev.podman_runner exec -- \
  python3 -m unittest scripts.release.test_verify_publication -v
python3 -m scripts.dev.podman_runner task verify-release \
  --dist-root .tmp/candidate/<source-sha>/dist
```

Expected: PASS; embedded source/version/epoch, manifest, checksums, SBOMs, matrices,
and asset inventory all match the merged source and Task 8 evidence.

- [ ] **Step 4: Create the approved annotated tag at the verified source**

Create the owner-approved annotated tag from `VERSION` at the already-verified
`<source-sha>`. Do not merge anything else first and do not use a lightweight tag.

```bash
release_version=$(tr -d '\n' < VERSION)
test -n "$release_version"
git tag -a "v${release_version}" <source-sha> \
  -m "GoREveal v${release_version}"
git push origin "v${release_version}"
```

Expected: `git cat-file -t "v${release_version}"` prints `tag`, and its peeled commit
is exactly `<source-sha>`.

- [ ] **Step 5: Publish the exact staged artifact without rebuilding**

```bash
gh workflow run release.yml --ref <tag> -f mode=publish \
  -f source_sha=<source-sha> \
  -f candidate_run_id=<candidate-run-id> \
  -f candidate_manifest_sha256=<approved-manifest-digest>
gh run watch <publication-run-id> --exit-status
```

Expected: workflow downloads the named candidate artifact, re-verifies exact bytes,
checks tag/source/compliance/evidence, and uploads those bytes. It executes no build.

- [ ] **Step 6: Download and verify remotely**

```bash
gh release download <tag> --dir .tmp/publication/<tag>
gh release view <tag> \
  --json url,tagName,isDraft,isPrerelease,targetCommitish,publishedAt,assets \
  > .tmp/publication/<tag>/release.json
python3 -m scripts.dev.podman_runner task verify-publication \
  --tag <tag> \
  --download-root .tmp/publication/<tag> \
  --candidate-manifest .tmp/candidate/<source-sha>/dist/release-manifest.json \
  --release-json .tmp/publication/<tag>/release.json \
  --record-out docs/evidence/rt1-r1/<tag>.json
```

Expected: remote asset names and SHA-256 values exactly match the candidate; release
is non-draft; source/tag/workflow identities match; SBOM and matrices verify. The
verifier writes `<tag>.json` atomically only after all checks pass; that record is an
output of Step 6, never a precondition for its own verification.

- [ ] **Step 7: Record post-publication evidence and close R1**

Create `candidate.json` and `<tag>.json` on a post-release evidence branch. These
records reference the immutable tagged source and do not pretend to be part of it.
Update status/navigation only after remote verification.

Update status/navigation to `S0 -> S1 -> S2A -> S2B -> S2C -> R1(closed) -> S3A`.
This authorizes execution of the separately reviewed IDACli S3A plan; it does not
merge or start that implementation automatically.

- [ ] **Step 8: Commit post-publication evidence**

```bash
git add docs/evidence/rt1-r1 README.md AGENTS.md .agents/skills/goreveal-navigation/SKILL.md \
  skills/goreveal-navigation/SKILL.md
git commit -m "release: verify published standalone R1 assets"
```

Merge this evidence commit before treating the IDACli hard gate as satisfied.

## Final Acceptance

The plan is complete only when:

- S2C evidence is valid even when the final SIMD decision is negative;
- every shipped optimized path has scalar equivalence and safe runtime fallback;
- representative bare-metal x86-64 and ARM64 evidence exists;
- AVX-512 is either rejected or has same-host AVX2 comparison plus real AVX2-only fallback;
- release assets reproduce byte-for-byte and verify offline;
- the standalone operator workflow requires no RE tool;
- legal/compliance is closed;
- the release is actually public and its downloaded bytes match the recorded manifest;
- only then is IDACli `RT1-S3A` eligible to start.
