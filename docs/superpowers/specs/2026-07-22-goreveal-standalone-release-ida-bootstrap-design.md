# GoREveal Standalone Release, SIMD Qualification, and IDA Bootstrap Design

> Status: approved direction; written-spec review pending
> Date: 2026-07-22
> Decision owner: maintainers
> Refines: `2026-07-22-goreveal-rt1-product-design.md`
> Planning gate: no `RT1-S2C`, `RT1-S3A`, or `RT1-S3B` implementation plan is
> active until this written specification is reviewed and approved.

## Decision Summary

GoREveal is **standalone-first with optional integrations**.

The next release must be useful and contract-complete as a standalone Go
binary. It must not depend on IDA, `idacli`, IDAPython, idalib, or any plugin to
recover, preserve, inspect, store, compare, verify, or export its canonical
truth. IDA work begins only after that standalone release gate closes.

The approved sequence is:

1. complete the RT1 correctness, evidence, identity, location, and v2 export
   obligations in `RT1-S0` through `RT1-S2B`;
2. run a standalone performance qualification in `RT1-S2C`, including an
   explicit SIMD support decision based on real hotspots and modern CPU
   coverage;
3. close and publish the `RT1-R1` standalone release if its legal, correctness,
   compatibility, corpus, and performance gates pass;
4. build an external headless IDA database bootstrap in `RT1-S3A`, with
   `selective` as the default mode and `preseed` as an opt-in mode;
5. build a thin interactive IDA plugin in `RT1-S3B` over the same frozen export
   and action contracts.

`RT1-R1` is a release milestone, not another feature sprint. A later `RT1-S10`
may widen the supported matrix and public-release posture, but it no longer
blocks the first contract-complete standalone release.

## Product Invariant: Standalone Means Independent

A standalone GoREveal installation must provide, without a host RE tool:

- binary and Go build identity;
- raw recovery and explicit stage diagnostics;
- package, function, source, type, string, runtime, and refinement surfaces to
  the extent supported by evidence;
- provenance, confidence, ambiguity, and explicit unavailable states;
- deterministic CLI and schema output;
- SQLite persistence and build-to-build comparison;
- review and handoff artifacts;
- identity-bound IDA/Ghidra exports as optional outputs;
- corpus, snapshot, differential, fuzz, and benchmark evidence for claimed
  behavior.

No release capability may be documented as available if it exists only in an
IDA plugin. Plugins consume canonical exports; they do not implement recovery,
semantic inference, address translation, or identity policy.

The dependency direction remains:

```text
binary -> core -> schema -> engine -> CLI/storage/export
                                      |
                                      +-> optional host consumers
```

There is no reverse dependency from `core`, `schema`, or `engine` to IDA.

## RT1-R1 Standalone Release Boundary

### Required contract obligations

The release gate requires all of the following:

1. every attempted recovery stage reports `available`, `unavailable`,
   `unsupported`, or `failed` rather than silently discarding errors;
2. raw and refined entities use stable IDs rather than array positions;
3. only cardinality-safe `1:1` evidence may enter automatic diff or transfer
   tiers;
4. half-open address ranges and format-neutral location mapping are enforced;
5. binary identity, Go build provenance, analyzer identity, and exact provider
   artifact identity are distinct;
6. v1 wire bytes remain frozen while v2 is an explicit, independently
   negotiated contract;
7. the selected v2 envelope is deterministic, verifiable without implicit JSON
   reserialization, and measured on a large artifact;
8. wrong-binary, wrong-base, changed-byte, missing-component, and unmappable
   location cases fail closed;
9. every advertised verification command runs real work in the pinned Podman
   environment;
10. release claims link to corpus, snapshot, differential, fuzz, benchmark, or
    compatibility evidence as appropriate.

The existing `RT1-S0` through `RT1-S2B` definitions remain the source for the
detailed gates. `RT1-S2C` does not reopen their semantics; it qualifies the
finished standalone artifact and closes release evidence.

### Release outputs

If the gate passes, `RT1-R1` produces:

- a pinned standalone GoREveal build and container image;
- checksums and an SBOM;
- an explicit supported matrix for Go version, GOOS, GOARCH, binary format,
  stripped state, and protected-state claims;
- frozen v1 compatibility evidence and the selected v2 contract fixtures;
- small, medium, and large p50/p95 runtime, peak RSS, and output-size records;
- a CPU capability and selected-kernel report for every benchmark host;
- one-page standalone operator workflows for analyze, inspect, persist, diff,
  verify, and export;
- known limitations and unsupported states that match actual behavior.

Repository licensing and release-compliance posture remain a hard publication
gate. If engineering gates pass but the license decision is unresolved, the
release candidate may be reproducible internally but must not be described as a
published public release. `RT1-R1` remains blocked and IDA implementation does
not start until the release is actually published.

## SIMD Qualification

### Purpose

SIMD is a possible implementation of measured performance work, not a release
feature and never a correctness dependency. The operator-visible outcome is a
faster standalone analysis with identical schema and wire bytes.

The required order is:

1. profile the end-to-end standalone workflow on registered artifacts;
2. identify a hotspot that materially contributes to wall time or peak RSS;
3. preserve or add a pure-Go scalar reference implementation;
4. compare against optimized scalar and relevant standard-library primitives;
5. add an architecture-specific path only if it still has a measurable product
   benefit;
6. prove exact semantic equivalence and retain runtime dispatch to scalar.

The historical `docs/tmp/draft/simd-optimization.md` is research input only.
Its throughput and end-to-end latency estimates are not release claims.

### Candidate kernels

Only bounded data-parallel kernels are initial SIMD candidates:

- repeated byte-pattern or magic scanning;
- byte classification and bulk string-candidate filtering;
- bulk compare/filter kernels;
- hashing or fingerprinting where the hash contract is already defined.

Parsing control flow, address mapping, schema construction, confidence logic,
and recovery decisions remain scalar truth. SIMD may accelerate evidence
collection but cannot change what counts as evidence.

Before writing a custom kernel, the experiment must compare the existing code
with `bytes`, `strings`, `crypto`, `hash`, or other relevant standard-library
paths because those may already use architecture-specific assembly.

### Supported CPU policy

The release binary keeps the broadest supported scalar baseline. In particular,
the default `linux/amd64` artifact must not raise its global requirement to
`GOAMD64=v3` merely to obtain AVX2. Optional kernels use runtime dispatch.

The qualification matrix is:

| Architecture | Required release path | Optional qualified paths | Initial policy |
| --- | --- | --- | --- |
| `amd64` | scalar baseline | AVX2; AVX-512 only if it beats AVX2 end to end | detect with `golang.org/x/sys/cpu`; never execute an unsupported instruction |
| `arm64` | scalar baseline | ASIMD/NEON | ASIMD is part of ARM64, but the optimized path still needs equivalence and benchmark evidence |
| other supported architectures | scalar baseline | none initially | report scalar selection explicitly |

SVE and SVE2 may be inventoried on ARM64 hosts but are not an `RT1-R1`
implementation promise. AVX-512 is also not presumed faster: frequency,
thermal, virtual-machine, and workload effects require a same-host AVX2 versus
AVX-512 comparison.

Go 1.26's `simd/archsimd` is experimental, enabled with
`GOEXPERIMENT=simd`, and does not carry the Go 1 compatibility promise. It may
be evaluated in an experimental benchmark/build lane, but the default release
must not require it. Promotion into a normal release requires a later explicit
toolchain-compatibility decision. Stable alternatives may include existing
standard-library paths or small architecture-specific assembly behind the same
internal scalar contract.

### Runtime capability report

The standalone diagnostics surface must be machine-readable and include:

- GOOS, GOARCH, Go version, and relevant build experiment flags;
- detected CPU features used by dispatch;
- compiled kernel implementations;
- selected implementation for each optimized operation;
- why a faster-looking path was unavailable or rejected;
- whether the result is a correctness-only virtualized run or a performance
  qualification run.

The final CLI spelling is an implementation-plan decision. A conceptual shape
is `goreveal diagnostics cpu --json`; the schema is more important than this
provisional command name.

### Benchmark and equivalence gates

Every candidate must have:

- table-driven scalar-versus-optimized equivalence tests across boundary sizes,
  alignments, empty inputs, tails, malformed inputs, and overlapping patterns;
- fuzzing against the scalar oracle;
- a forced-scalar mode in tests and benchmarks;
- a forced-implementation mode that rejects unsupported hardware rather than
  risking an illegal instruction;
- microbenchmarks with throughput and allocations;
- end-to-end benchmarks on the same registered binary and host;
- cold and warm cache observations where I/O materially changes the result;
- recorded CPU model, virtualization state, governor/power context when
  available, Go version, command, commit, and raw benchmark output.

An optimized kernel is eligible for the default release only when all of these
conditions hold:

1. the original hotspot accounts for at least `10%` of the registered large
   workflow before optimization;
2. the optimized kernel is at least `1.5x` faster than the best scalar path in
   its microbenchmark;
3. the large end-to-end workflow improves by at least `10%` at p50 or p95 on a
   representative bare-metal host;
4. output and error semantics are byte-for-byte equivalent where the contract
   is deterministic;
5. scalar-only and unsupported-feature paths remain green;
6. unrelated supported workflows do not regress by more than `5%` without an
   explicitly accepted reason.

If a kernel misses the gate, keep the evidence and ship the scalar path. A
negative SIMD result still completes the qualification.

The current KVM development host exposes an Intel Xeon Gold 6348 model with
AVX2 and broad AVX-512 flags. It is suitable for dispatch and correctness
tests, but virtualization means it cannot by itself support a public
performance claim. `RT1-S2C` therefore needs at least one bare-metal x86-64
host and one ARM64/NEON host for representative performance evidence. An
AVX2-only x86-64 host is required before the project claims that fallback from
AVX-512 was tested on real hardware.

## IDA Database Bootstrap

### Feasibility boundary

GoREveal cannot safely write a valid `.i64` database by itself. IDA owns its
loader, segment model, analysis queues, database format, and save lifecycle.
The viable optimization is to let an installed IDA engine create a minimal
database with global autoanalysis disabled, import identity-verified GoREveal
facts, schedule only bounded analysis, and save the resulting database.

This avoids most full autoanalysis; it does not avoid IDA's initial file loader
or database creation. Documentation and benchmarks must use that exact claim.

The first-time workflow belongs in the existing native C++ IDACli/IDA SDK
consumer, invoked through the installed headless IDA executable. An ordinary
interactive UI opens too late to be the primary bootstrap path because IDA has
already opened a database and may already have started global analysis.

`idaclictl` may orchestrate process launch, artifacts, and exit status, but it
has no IDA SDK and cannot create or save the database by itself. The S3A
implementation must not introduce IDAPython. idalib remains a separately gated
future runtime alternative, not the default bootstrap architecture.

### Shared bootstrap flow

Both modes use the same ordered flow:

1. consume a frozen GoREveal v2 artifact and hash its exact bytes;
2. verify the input binary digest, size, format, architecture, and build identity;
3. launch the installed headless IDA executable with global automatic analysis
   disabled, load the native IDACli task, and let IDA's loader create the
   initial segments and database state;
4. read IDA's loaded segments and image base, then validate every location
   mapping before creating an action;
5. generate a deterministic preview containing the provider digest, IDA/input
   identity, normalized mapping, mode, selectors, and ordered actions;
6. require explicit approval of that exact preview digest before mutation;
7. create only non-conflicting functions and reviewed names;
8. schedule and wait for analysis only for the bounded ranges allowed by the
   selected mode;
9. validate outcomes, save the database through IDA, and write a machine-readable
   coverage/evidence manifest beside it.

No mode may use `del_func`, silently replace a function boundary, overwrite a
user name, ignore an unmappable address, or reconstruct an approved action plan
from a changed provider artifact.

### Mode 1: `selective` (default)

`selective` minimizes time to a useful analyst target.

The operator selects packages, functions, addresses, or a bounded named target
set. The consumer imports exact function boundaries and names for that set plus
only the minimum declared supporting range. IDA analysis queues are planned and
waited only for those ranges. The output manifest records precisely what was
selected, imported, analyzed, skipped, and left unknown.

The default must never silently expand into global autoanalysis. If a selector
cannot be resolved from canonical GoREveal truth, the preview reports it as
unresolved rather than guessing a broad range.

This mode is the recommended answer for very large binaries when the analyst
needs a small set of packages or functions quickly.

### Mode 2: `preseed` (opt-in)

`preseed` initializes the database with every exact, mappable function boundary
and safe default-name replacement supplied by GoREveal, but still does not run
global autoanalysis.

The import is deterministic, chunked, restart-aware, and evidence-recorded.
Conflicts are skipped and reported. After preseed, the analyst or plugin can
request bounded analysis on demand. The mode is useful when broad navigation
matters more than the shortest time to the first target, but it must be measured
because hundreds of thousands of IDA function actions may themselves be
expensive.

If preseed cannot beat or materially improve the conventional full-analysis
workflow on the registered large artifact, it remains experimental; selective
mode is not blocked by that result.

### Database reuse and cache identity

A bootstrap database may be reused only when its cache key matches:

- exact input binary SHA-256 and size;
- IDA product/version and processor/loader identity;
- GoREveal provider contract and exact artifact digest;
- image base and normalized segment mapping;
- bootstrap mode, selectors, and relevant analysis options;
- consumer/plugin contract version.

The coverage manifest is part of reuse validation. A selective database cannot
be represented as globally analyzed. A changed selector set produces an
incremental preview against the existing identity-bound database rather than an
implicit full rebuild.

### Failure and interruption behavior

Bootstrap fails closed on identity, contract, mapping, or capability mismatch.
Unsupported provider fields are explicit and do not become empty actions.
Individual boundary or naming conflicts may be skipped only when the preview
classified them that way and the final evidence records the same reason.

An interrupted run may retain a private checkpoint for restart, but it does not
publish the final output path or a `complete` manifest. Repeating the exact
approved plan against a completed database performs zero mutations.

### Large-binary acceptance experiment

The field case in
`docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md` is the first
registered large experiment: approximately 410 MB, a recorded conventional IDA
autoanalysis time of 1 hour 34 minutes, and 458,600 GoREveal-recovered
functions. The experiment must preserve the original binary and tool identity
records; counts alone are not evidence of correctness.

Compare three runs on the same host and IDA version:

1. conventional full automatic analysis;
2. `selective` bootstrap over a pre-registered analyst target set;
3. `preseed` bootstrap followed by the same target requests.

Record:

- loader time, provider verification time, import time, bounded analysis time,
  and time to first usable selected target;
- peak RSS, database size, total function actions, created/skipped/conflicting
  actions, and queued ranges;
- decompilation or navigation success for each pre-registered target;
- incorrect boundaries, overwritten names, identity failures, crashes, and
  rerun mutations;
- reopen time for the saved database.

`selective` is promoted only if it is at least `4x` faster than conventional
full analysis to the first usable registered target, all registered targets are
classified, unsafe mutations remain `0`, and the second apply performs `0`
mutations. `preseed` is promoted from experimental only if it provides a
measured operator benefit without weakening the same safety gates. Raw function
count is not a success metric.

## Thin Interactive IDA Integration

The interactive surface follows headless bootstrap and, by default, extends the
existing native IDACli plugin rather than creating another recovery plugin. It
uses the same provider verifier, preview contract, action classifier, and
coverage manifest. Its responsibilities are limited to:

- display identity and coverage status;
- request or preview an incremental selective import;
- navigate from GoREveal entities to mapped IDA addresses;
- show conflicts, unavailable evidence, and reasons for skipped actions;
- apply an explicitly approved bounded plan;
- request additional bounded IDA analysis;
- export identity-bound host observations back as external evidence.

The plugin does not parse Go runtime metadata, infer missing function
boundaries, calculate a second address mapping, or maintain a private semantic
schema. Headless and interactive paths must produce the same ordered action
classification for the same provider artifact, mapping, and selectors.

## Security and Compatibility

- The standalone release contains no IDA SDK or licensed IDA runtime.
- IDA consumers operate only with an operator-provided, installed, compatible
  IDA environment.
- Provider artifacts are untrusted input: parsers are bounded, exact-byte
  digests are verified, and malformed envelopes are negative fixtures.
- v1 consumers remain supported for the documented compatibility window but
  cannot request v2-only bootstrap actions.
- The v2 verifier rejects unknown required fields or contracts rather than
  silently downgrading.
- Saved databases and manifests are written to isolated output paths; the
  original analyst database is not mutated by default.

## Evidence and Definition of Done

### `RT1-S2C` is done when

- the standalone release obligations above are traced to passing evidence;
- real benchmark and fuzz targets replace every no-op gate;
- CPU capability reporting is deterministic and tested;
- each measured hotspot has a scalar baseline and a recorded keep/reject SIMD
  decision;
- any shipped optimized path passes equivalence, dispatch, hardware, and
  end-to-end gates;
- unsupported CPU paths fall back safely;
- release artifacts, checksums, SBOM, matrix, limitations, and operator docs
  are reproducible inside the project workflow;
- every unresolved publication gate is explicit; a blocked gate leaves
  `RT1-R1` open.

### `RT1-R1` closes when

- the qualified standalone artifact is actually published with its checksums,
  SBOM, supported matrix, compatibility evidence, and limitations;
- the release can perform every claimed standalone workflow without an IDA,
  Ghidra, `idacli`, idalib, IDAPython, or plugin dependency;
- release provenance identifies the source commit, toolchain, container, and
  evidence record;
- no unresolved legal, licensing, or compliance blocker remains.

### `RT1-S3A` is done when

- both bootstrap modes use the same frozen identity and action contracts;
- wrong binary, wrong base, changed artifact, unsupported contract, and
  unmappable location fixtures reject in every case;
- `selective` never schedules global analysis and passes its large-binary
  acceptance gate;
- `preseed` has a measured promote-or-remain-experimental decision;
- boundary replacement, user-name overwrite, and deletion counts are `0`;
- second apply performs `0` mutations;
- interrupted or partial output cannot masquerade as complete;
- saved-database reuse rejects every cache-key mismatch.

### `RT1-S3B` is done when

- plugin and headless classifications are identical on shared fixtures;
- the plugin contains no recovery or mapping logic duplicated from GoREveal;
- all mutations require an identity-verified preview and explicit approval;
- incremental selective import and coverage reporting work on a real IDB;
- host observations return only as identity-bound external evidence.

## Non-Goals

This design does not authorize:

- a native `.i64` writer in GoREveal;
- a generic disassembler or decompiler;
- unconditional global IDA autoanalysis after import;
- direct mutation of an analyst's only database copy;
- SIMD in parsers or semantic decisions merely because the CPU supports it;
- raising the default CPU baseline to AVX2 or AVX-512;
- making experimental `simd/archsimd` a release dependency;
- implementing SVE/SVE2 before a measured target and supported toolchain exist;
- moving IDA-specific concepts into `core`, `schema`, or engine recovery.

## Documentation and Planning Consequences

This specification refines the RT1 product design as follows:

- insert `RT1-S2C` and the `RT1-R1` standalone release before any IDA work;
- split the old `RT1-S3` into `RT1-S3A` headless bootstrap and `RT1-S3B` thin
  plugin;
- treat old `RT1-S3` tasks in the Horizon A implementation plan as suspended
  planning input, not executable work;
- retain `RT1-S0` through `RT1-S2B` ordering and detailed tasks;
- keep `RT1-S4` and later semantic/product lanes behind the new standalone and
  integration gates until the program plan is rewritten.

After written-spec approval, a new implementation plan must define file-level
tasks, owners, fixtures, commands, evidence storage, hardware access, and the
cross-repository `idacli` boundary. It must not combine standalone release work
and IDA mutation work in one sprint.

## Primary Sources

External API and toolchain facts were checked against official documentation on
2026-07-22:

- [Go 1.26 release notes](https://go.dev/doc/go1.26)
- [`simd/archsimd` package](https://pkg.go.dev/simd/archsimd)
- [`golang.org/x/sys/cpu`](https://pkg.go.dev/golang.org/x/sys/cpu)
- [IDA command-line switches](https://docs.hex-rays.com/9.0/user-guide/configuration/command-line-switches)
- [IDA Domain database API](https://ida-domain.docs.hex-rays.com/ref/database/)
- [IDAPython automatic-analysis API](https://python.docs.hex-rays.com/ida_auto/index.html)
- [IDAPython loader/database API](https://python.docs.hex-rays.com/ida_loader/index.html)
- [idalib documentation](https://docs.hex-rays.com/user-guide/idalib)

The IDA Domain, IDAPython, and idalib references validate loader,
autoanalysis, bounded-analysis, and save-lifecycle capabilities across official
IDA automation surfaces. They are not an implementation choice: the approved
S3A path follows the current native C++ IDACli/IDA SDK architecture and does
not add IDAPython.

Repository evidence and governing documents:

- `AGENTS.md`
- `docs/architecture/2026-03-19-goreveal-platform-contract.md`
- `docs/superpowers/specs/2026-07-22-goreveal-rt1-product-design.md`
- `docs/superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md`
- `docs/plans/2026-07-22-goreveal-proposal-post-ida-experience.md`
- `docs/tmp/draft/simd-optimization.md`
- `ioplane/idacli:docs/planning/2026-07-22-go-function-recovery-task.md`
