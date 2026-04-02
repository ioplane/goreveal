# GoREveal Initial Baseline Comparison Results

> Status: empirical comparison note
> Date: 2026-04-01
> Purpose: record the first real baseline comparison, the post-`PE` rerun, and the current intermediate rerun against `GoReSym`, `redress`, and `gore`, so the next roadmap move is chosen from measured gaps rather than from intuition.

## Method

Commands used:
- canonical normalized fixture report:
  - `python3 -m scripts.dev.podman_runner exec --with-baselines -- bash -lc 'cd /workspace && export GOREVEAL_BIN=/workspace/.tmp/goreveal && python3 -m scripts.baseline.generate_fixture_report corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin'`
- fresh combined intermediate rerun:
  - `python3 -m scripts.dev.podman_runner exec --with-baselines -- bash -lc 'cd /workspace && mkdir -p .tmp/fresh-eval && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal && python3 - <<\"PY\" ... PY'`
- real-world per-tool comparisons:
  - `goreveal analyze <binary>`
  - `scripts/baseline/run_goresym.sh <binary>`
  - `scripts/baseline/run_redress.sh <binary>`
  - `scripts/baseline/run_gore.sh <binary>`

Measured targets:
- canonical rich `ELF` fixture
- external stripped `ELF`: `ocserv-agent-linux-amd64`
- external `PE`: `GoReSym/testproject.exe`
- external `PE`: `hashicorp-re/bin/Keygen.exe`
- bounded `Mach-O` fixture: `go-macho-buildinfo-darwin-amd64`
- local extracted `rclone-linux-amd64`

## Rich ELF Fixture

Current normalized report result:
- `18 matched`
- `0 diverged`
- `3 skipped`

Reading:
- the current v1 overlap contract is fully green on the canonical rich `ELF` fixture
- no evidence currently suggests regressions versus the already-declared narrow overlap set

## Fresh External Target Comparison

| Target | Tool | Build Path | Functions | Packages | Files | Extra Reading |
| --- | --- | --- | ---: | ---: | ---: | --- |
| `ocserv-agent-linux-amd64` | `goreveal` | `github.com/dantte-lp/ocserv-agent/cmd/agent` | `14850` | `385` | `2` | `runtime=section_heuristic`, `peeling=14850` |
| `ocserv-agent-linux-amd64` | `GoReSym` | `github.com/dantte-lp/ocserv-agent/cmd/agent` | `6958` | `0` | `1065` | strong function and file-list breadth, but no package surface in the normalized path |
| `ocserv-agent-linux-amd64` | `redress` | `github.com/dantte-lp/ocserv-agent` | `316` | `10` | `13` | narrower source-oriented view |
| `ocserv-agent-linux-amd64` | `gore` | `github.com/dantte-lp/ocserv-agent/cmd/agent` | `34` | `8` | `13` | `types=6795` but not parity-usable |
| `testproject.exe` | `goreveal` | `command-line-arguments` | `1590` | `31` | `1` | `runtime=section_heuristic`, `peeling=1590`, first bounded `PE` function/package foothold |
| `testproject.exe` | `GoReSym` | `command-line-arguments` | `8` | `0` | `191` | much narrower function/package surface, but wider file list than the bounded GoREveal PE foothold |
| `testproject.exe` | `redress` | `command-line-arguments` | `4` | `1` | `1` | much narrower source-oriented view |
| `testproject.exe` | `gore` | `command-line-arguments` | `2` | `1` | `1` | `types=621` |
| `Keygen.exe` | `goreveal` | `hashi-gen` | `4264` | `75` | `5` | `runtime=section_heuristic`, `peeling=4264`, strong bounded `PE` result on a second external sample |
| `Keygen.exe` | `GoReSym` | `hashi-gen` | `685` | `0` | `374` | narrower function surface, but broader file list than the bounded GoREveal PE result |
| `Keygen.exe` | `redress` | `hashi-gen` | `10` | `5` | `5` | narrow source-oriented view |
| `Keygen.exe` | `gore` | `hashi-gen` | `4` | `3` | `5` | `types=1651` |
| `rclone-linux-amd64` | `goreveal` | `github.com/rclone/rclone` | `83635` | `1774` | `431` | `runtime=section_heuristic`, `peeling=83635` |
| `rclone-linux-amd64` | `GoReSym` | `github.com/rclone/rclone` | `66911` | `0` | `5216` | strong function and file-list breadth, but no package surface in the normalized path |
| `rclone-linux-amd64` | `redress` | `github.com/rclone/rclone` | `11637` | `261` | `432` | strong file visibility, much narrower package/function view |
| `rclone-linux-amd64` | `gore` | `github.com/rclone/rclone` | `1414` | `259` | `432` | `types=42532` |
| `rclone-darwin-amd64` | `goreveal` | `github.com/rclone/rclone` | `93909` | `2955` | `432` | `runtime` absent, `peeling=93909`, first broad external `Mach-O` foothold |
| `rclone-darwin-amd64` | `redress` | `github.com/rclone/rclone` | `12193` | `273` | `433` | strongest current compared file list on external `Mach-O` |

## Intermediate Validation Rerun

This is the freshest repo-local rerun used for the current checkpoint.

Method:
- rich `ELF` fixture via `scripts.baseline.generate_fixture_report`
- in-repo `PE` fixture via direct per-tool rerun
- in-repo `Mach-O` fixture via direct per-tool rerun
- local extracted `rclone-v1.73.3-linux-amd64` binary from `.tmp/rclone-matrix/raw`

Measured targets:
- `corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin`
- `corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe`
- `corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin`
- `.tmp/rclone-matrix/raw/rclone-v1.73.3-linux-amd64.zip`

Rich `ELF` rerun:
- `18 matched`
- `0 diverged`
- `3 skipped`

Fresh current counts:

| Target | Tool | Functions | Packages | Files | Extra Reading |
| --- | --- | ---: | ---: | ---: | --- |
| `PE` fixture | `goreveal` | `1953` | `39` | `1` | `runtime=section_heuristic` |
| `PE` fixture | `GoReSym` | `120` | `0` | `226` | wide file list, much narrower function/package surface |
| `PE` fixture | `redress` | `3` | `1` | `1` | very narrow source-oriented foothold |
| `PE` fixture | `gore` | `3` | `1` | `1` | `types=837` |
| `Mach-O` fixture | `goreveal` | `2059` | `37` | `1` | bounded `Mach-O` function/package foothold remains stable |
| `Mach-O` fixture | `GoReSym` | `116` | `0` | `260` | wide file list, much narrower function/package surface |
| `Mach-O` fixture | `redress` | `3` | `1` | `1` | very narrow source-oriented foothold |
| `Mach-O` fixture | `gore` | `3` | `1` | `1` | `types=843` |
| `rclone-linux-amd64` | `goreveal` | `83635` | `1774` | `431` | `runtime=section_heuristic` |
| `rclone-linux-amd64` | `GoReSym` | `66911` | `0` | `5216` | strongest file-list breadth in this rerun, but no package surface |
| `rclone-linux-amd64` | `redress` | `11637` | `261` | `432` | still useful as a source-oriented baseline, but much narrower than `GoREveal` here |
| `rclone-linux-amd64` | `gore` | `1414` | `259` | `432` | `types=42532` |

Bounded timing sample on the same `rclone-linux-amd64` binary:

| Tool | Time |
| --- | ---: |
| `redress` | `1.13s` |
| `GoReSym` | `1.18s` |
| `gore` | `1.22s` |
| `goreveal analyze` | `1.23s` |

Current reading from this rerun:
- the fresh rich `ELF` overlap contract remains green
- the bounded `PE` and `Mach-O` footholds remain stable in the current repo-local rerun
- on the measured `rclone-linux-amd64` sample, `GoREveal` remains strongest on function/package breadth while also keeping real file visibility
- `GoReSym` still shows the broadest file-list-oriented output on the same rerun, but not comparable package-level workflow shape
- the latest timing sample is now tightly clustered across all four tools on the current runner, so there is still no operational efficiency problem large enough to outrank workflow or source-confidence work

## Compact Matrix

| Slice | GoREveal | GoReSym | redress | gore | Current reading |
| --- | --- | --- | --- | --- | --- |
| rich `ELF` fixture | overlap green on the declared narrow contract | usable baseline | usable baseline | usable baseline | no measured regression on the canonical overlap set |
| stripped external `ELF` | strongest current function/package/peeling surface plus real file visibility | strong function overlap in some cases, but weak package/file normalization | strong file visibility, narrower package/function view | strong type/file counts, narrow function/package view | `GoREveal` is already productively strong on real Linux/server binaries and no longer file-blind on the measured external samples |
| external `PE` | strong bounded function/package/peeling surface plus real file visibility | weaker normalized function/package/file surface on measured samples | narrow source-oriented view | narrow library-oriented view | `PE` is no longer posture-only and already useful on more than one external sample |
| external `Mach-O` | strong bounded function/package/peeling foothold plus real file visibility | not part of the fresh normalized rerun | strong file visibility, much narrower package/function view | not part of the fresh normalized rerun | `Mach-O` foothold is no longer just a fixture-local proof; runtime and type depth remain the main gaps |

## Weighted Reading

The comparison changes the roadmap reading in a concrete way.

### Where GoREveal is already strongest

- stripped `ELF` function and package recovery on a real external binary
- first analyst-facing code peeling on real stripped `ELF`
- first bounded `Mach-O` function/package/peeling foothold

### Where baseline tools still have the clearest practical edge

- richer source-backed analyst visibility semantics rather than just raw file counts
- deeper type/runtime semantics where `GoREveal` still intentionally returns `types=[]`
- broader generic analyst convenience outside Go-native trust/transfer workflows

### What this means

The fresh comparison still does **not** point first toward more `ELF` workflow polish.

It now points first toward:
- explicit workstation/MCP handoff hardening over the now-landed review and handoff CLI surfaces
- then either one more workflow/value increment or one more thin source-visibility/semantic slice, depending on the measured operator pain point

Reason:
- `ELF` is already productively strong
- `PE` and `Mach-O` now both have real external footholds
- the clearest remaining practical gap is no longer raw file absence on these measured samples, but stronger semantic/source confidence and broader analyst workflow value

## Fresh Practical Reading

- `GoREveal` is now clearly ahead on function/package/peeling coverage across the fresh `ELF`, `PE`, and `Mach-O` targets that were rerun
- `GoREveal` is no longer file-blind on these measured external targets; real file visibility now exists on all fresh `GoREveal` rows in this note
- `redress` remains the most credible compared source/file-oriented baseline in the fresh external pass
- `gore` still exposes large type sets, but that current type surface is not yet trustworthy parity evidence for user-meaningful type recovery
- `GoReSym` remains important as a runtime-aware behavioral reference, but the fresh normalized external pass does not justify treating it as the dominant practical product baseline outside runtime-semantic depth

## Next Weighted Decision

Recommended next move after the fresh comparison:
- keep the current `ELF` / `PE` / `Mach-O` footholds stable
- treat the fresh external matrix as evidence that workflow/value and handoff work now outrank another default parser lane
- harden explicit workstation / MCP interop over `diff review sqlite` and `diff handoff sqlite`
- treat the current intermediate rerun as confirmation that this ordering still holds after the latest local product changes

## Post-PE Follow-Up

That first bounded `PE` function foothold is now landed.

Current measured checkpoint:
- in-repo Windows `PE` fixture: `1953` functions, `39` packages, `1953` peeling functions, `runtime` still present as bounded posture
- external `GoReSym/testproject.exe`: `1590` functions, `31` packages, `1590` peeling functions, `runtime` still present as bounded posture

Updated next move:
- rerun the same baseline comparison against `GoReSym`, `redress`, and `gore`
- then choose between transfer-workflow polish, thin source-visibility work, and the next bounded evidence-backed slice

## Post-PE Rerun

The rerun is now completed for the current highest-signal slices.

Rerun method:
- canonical rich `ELF` report rerun through `scripts.baseline.generate_fixture_report`
- focused per-tool rerun on:
  - external `PE`: `GoReSym/testproject.exe`
  - bounded `Mach-O` fixture: `go-macho-buildinfo-darwin-amd64`

Rerun checkpoint:
- rich `ELF` normalized report remains:
  - `18 matched`
  - `0 diverged`
  - `3 skipped`
- external `PE` rerun:
  - `GoREveal`: `1590` functions, `31` packages, `1` file, `1590` peeling, `runtime=section_heuristic`
  - `GoReSym`: `8` functions, `0` packages, `191` files
  - `redress`: `4` functions, `1` package
  - `gore`: `2` functions, `1` package
- bounded `Mach-O` rerun:
  - `GoREveal`: `2059` functions, `37` packages, `1` file, `2059` peeling
  - `GoReSym`: `116` functions, `0` packages, `260` files
  - `redress`: `3` functions, `1` package
  - `gore`: `3` functions, `1` package

Updated reading after the rerun:
- `PE` is no longer the clearest remaining format-gap lane for function/package footholds
- `GoREveal` is now ahead on function/package/peeling surface for the rerun `PE` and `Mach-O` checkpoints
- the fresh external rerun now also shows real file visibility on `GoREveal` rows, so the main remaining gap is no longer “files or no files”, but stronger semantic/source confidence and richer workstation integration
- the first thin source-visibility response is now landed too, through line-table-backed `source_tree` fallback with `pathless_file_evidence`
- the first transfer-workflow polish response is now landed too, through package-level `transfer_packages` summaries plus `diff review sqlite` and `diff handoff sqlite`
- because the function/package footholds are already landed across all three current formats, the next weighted move should now favor:
  - workstation / MCP handoff hardening
  - then either one more workflow/value increment or the next thin semantic/source-confidence improvement
  over another parser lane by default

## Adjacent RE Ecosystem Reading

The comparison against `gore`, `GoReSym`, and `redress` is still the main clean-room product comparison.
But the broader workstation context is now clearer too through `rehelp` and the current RE lab host.

Measured host signals:
- `static-go` pipeline already groups `goresym`, `redress`, `alphagolang`, `goresolver`, and `gostringungarbler`
- `ida-pro`, `ghidra`, `jeb`, `rizin`, and `retdec` are all present on the same remote workstation
- diffing adjacencies such as `diaphora` and `binexport` are available too
- `ida-pro-mcp` is present, making host-platform MCP interop concrete rather than hypothetical

What this changes:
- the baseline comparison should stay centered on Go-native utilities for recovery claims
- but the roadmap should now treat host-platform interop and future version-tracking handoff as more concrete
- `GoREveal` is best positioned as the Go-native truth and transfer layer inside a richer workstation, not as a replacement for that workstation

See also:
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md`
- `docs/plans/2026-04-01-goreveal-universal-re-workbench-comparison.md`

## Important Limits

- this note is a real comparison, not a parity claim
- `redress` still does not provide `go_version` in the normalized path
- current `gore` type counts are not treated as user-type parity evidence
- current `GoReSym` comparison still should not be misread as full-function parity proof
