# GoREveal REHelp and RE Lab Inventory Notes

> Status: clean-room research note
> Date: 2026-04-01
> Purpose: record what `rehelp` and the current remote RE workstation reveal about adjacent tooling, so roadmap and interop decisions are based on real operator context rather than generic assumptions.

## Scope

This note is about:
- adjacent RE tooling available in the current operator environment
- clean-room product and roadmap implications for `GoREveal`
- what to borrow as workflow and integration ideas

This is not a proposal to move third-party recovery logic into `core`.

## Inputs

Local reference repo reviewed:
- `/opt/projects/repositories/rehelp`

Remote workstation path reviewed through `Teleport`:
- `tsh ssh root@oel-lab-gui`

Representative remote commands:
- `rehelp version --json`
- `rehelp categories --json`
- `rehelp pipeline static-go --json`
- `rehelp pipeline dynamic --json`
- `rehelp list --category=re-disassemblers --json`
- `rehelp list --category=re-tools --json`
- `rehelp search diff --json`
- `rehelp search go --json`
- `rehelp sync --json`

## What REHelp Confirms

`rehelp` is not only a tool catalog.
It already models:
- reproducible tool inventory
- named workflows and pipelines
- local and remote operator usage
- CLI and MCP-oriented automation posture

The strongest reusable pattern is operational, not algorithmic:
- one inventory
- one normalized operator view
- multiple execution paths (`CLI`, remote shell, later MCP)

That is directly relevant to `GoREveal` because the current repo is already strongest as:
- a Go-native recovery and trust layer
- a workflow and transfer layer
- an integration surface for larger RE workstations

## Remote Host Signals

The current RE lab host is a broad workstation, not a narrow Go-only box.

Measured categories:
- `re-disassemblers`: `7`
- `re-tools`: `8`
- `ida-plugins`: `31`
- `rizin-plugins`: `4`
- `java-re`: `12`
- `debuggers`: `7`
- `firmware`: `8`

Measured named pipelines:
- `static-go`
  - `goresym`
  - `redress`
  - `alphagolang`
  - `goresolver`
  - `gostringungarbler`
  - `golangci-lint`
  - `dlv`
- `dynamic`
  - `gdb`
  - `dlv`
  - `frida`
  - `bpftrace`
  - `valgrind`
  - `uftrace`
  - `angr`
  - `unicorn`
  - `qiling`
  - `afl-fuzz`

Measured adjacent tools on the host:
- Go-focused references:
  - `goresym`
  - `redress`
  - `goresolver`
  - `gostringungarbler`
  - `alphagolang`
  - `dlv`
- host platforms and import targets:
  - `ida-pro`
  - `ghidra`
  - `jeb`
  - `rizin`
  - `radare2`
  - `retdec`
  - `pyghidra`
  - `headless-ida`
- diffing and transfer adjacencies:
  - `diaphora`
  - `binexport`
  - `binsync`
  - `ida-pro-mcp`
- dynamic and symbolic adjacencies:
  - `frida`
  - `angr`
  - `unicorn`
  - `qiling`
  - `gdb`
  - `bpftrace`
  - `uftrace`
  - `z3`
- triage/signature/string adjacencies:
  - `capa`
  - `floss`
  - `yara`
  - `findcrypt-yara`
- Rizin ecosystem:
  - `rz-retdec`
  - `sigdb`
  - `angrcutter`

## What GoREveal Can Borrow

### Borrow as workflow/product pattern

- explicit remote-workstation guidance, not only local repo commands
- host inventory as a planning artifact
- named operator workflows instead of only low-level commands
- CLI and MCP parity where the same truth can be consumed by people and agents

### Borrow as roadmap signal

- `ida-pro-mcp` makes host-platform MCP interop concrete, not theoretical
- `diaphora` and `binexport` make future function-level diffing and transfer workflows more real
- `frida`, `angr`, `qiling`, `unicorn`, `uftrace`, and `z3` confirm a realistic external orchestration environment for later protected/deobfuscation work
- `rizin` plus its plugins make a thin future adapter more credible technically, even if it still stays below `JEB` and `Binary Ninja` in product priority

### Borrow as execution model

`GoREveal` should increasingly be treated as:
- the Go-native truth layer
- the schema/export/transfer layer
- the orchestrator-facing specialist that plugs into a richer RE workstation

That is a stronger and lower-risk product position than trying to absorb generic workstation capabilities into `core`.

## What Not To Borrow

Do not borrow:
- disassembler or decompiler logic into `core`
- debugger, emulator, or symbolic-execution logic as native recovery dependencies
- host-platform workspace semantics into canonical `GoREveal` schema
- a “generic RE Swiss Army knife” product direction

The clean-room and architecture boundary still holds:
- `core` remains native Go recovery
- `engine` remains orchestration and enrichment
- external tools remain references, adapters, or orchestrated sidecars

## Roadmap Implications

This research changes prioritization more than capability claims.

### Immediate

- keep the workflow/value lane primary
- add a bounded CLI/operator review projection next
- make MCP/interop notes more concrete now that `ida-pro-mcp` is known to exist on the real host

### Near-term backlog adjustments

- move host-platform MCP interop from generic idea to explicit `IDA`/`Ghidra`/`IDA MCP` handoff note
- treat `Diaphora` / `BinExport` as real external-reference inputs for future version-tracking work
- keep `Rizin` in backlog, but note that technical feasibility is stronger than market priority

### Later, only if measured need remains

- use the host dynamic/symbolic toolchain as an external orchestration lane for protected/deobfuscation work
- do not promote that lane ahead of current workflow/value work without a measured pain point

## Weighted Conclusion

The RE lab inventory does not argue for a new parser lane.
It argues for a stronger interop/orchestration roadmap.

Current weighted next order after this research:
1. bounded CLI/operator review projection over the existing transfer-review surfaces
2. explicit host-platform handoff projection over that same review state
3. explicit host-platform MCP interop hardening in roadmap and docs
4. later function-level diffing and metadata-transfer work informed by `Diaphora` / `BinExport`
5. protected/deobfuscation orchestration only when a measured gap justifies it
