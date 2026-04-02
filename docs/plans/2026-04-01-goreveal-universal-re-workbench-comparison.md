# GoREveal Universal RE Workbench Comparison

> Status: functional comparison note
> Date: 2026-04-01
> Purpose: compare `GoREveal` with general-purpose RE workbenches by module and workflow role, so roadmap decisions do not confuse Go-native recovery depth with generic disassembly-suite breadth.

## Scope

This note is not an apples-to-apples output benchmark.

It compares:
- `GoREveal` as a Go-native recovery, trust, and transfer product
- universal RE workbenches present in the measured operator environment

It does not claim:
- feature parity with `IDA`, `Ghidra`, `JEB`, `Rizin`, or `RetDec`
- that generic workbenches should become native `core` dependencies

## Inputs

Fresh GoREveal empirical inputs:
- `.tmp/fresh-eval/goreveal-matrix.json`
- `.tmp/fresh-eval/go-baselines.json`
- `.tmp/intermediate-eval/*.json`
- current bounded timing sample on `rclone-linux-amd64`

Measured workstation / RE-lab inputs:
- `tsh ssh root@oel-lab-gui -- 'rehelp pipeline static-go --json'`
- `tsh ssh root@oel-lab-gui -- 'command -v rizin frida angr uftrace z3'`
- `docs/plans/2026-04-01-goreveal-rehelp-and-re-lab-inventory-notes.md`

## Product Role Split

The correct product reading is:
- `GoREveal` is the Go-native truth, trust, and transfer layer
- universal RE workbenches are analyst workspaces, disassembly/decompilation shells, and broader orchestration environments

That means the right comparison is role-based, not “who is the biggest Swiss Army knife”.

## Functional Matrix

| Capability / Module | GoREveal | IDA Pro | Ghidra | JEB | Rizin | RetDec | Current reading |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Go-native build info / package / function recovery | `high` | `low` standalone | `low` standalone | `medium` with plugins | `low` standalone | `low` | `GoREveal` is the specialist; generic workbenches need imports, scripts, or plugins |
| Go-native trust / provenance contract | `high` | `low` | `low` | `low` | `low` | `low` | canonical schema and bounded trust fields are a real differentiator |
| Go-specific transfer workflow | `high-medium` | `medium` with plugins | `medium` with scripts | `medium` | `low-medium` | `low` | `diff review sqlite` and `diff handoff sqlite` now make this a real product surface |
| Stored-run diffing / review queues | `high` | `medium` via external tooling | `medium` via external tooling | `medium` | `low` | `low` | `GoREveal` already has native review-oriented diff state; workbenches usually need sidecars such as `Diaphora`, `BinExport`, or `BinSync` |
| Workspace mutation / annotations | `low` by design | `high` | `high` | `high` | `medium` | `low` | this remains the domain of host platforms, not `GoREveal core` |
| General disassembly / decompilation breadth | `low` by design | `high` | `high` | `high` | `medium` | `medium` | `GoREveal` should not compete here |
| Headless automation / scripting | `high` | `medium-high` | `medium-high` | `medium` | `high` | `medium` | `GoREveal` is already strong in CLI/schema/export automation; `Rizin` is especially strong as a headless generic shell |
| MCP / agent-facing interop | `medium-high` | `medium` with `ida-pro-mcp` | `medium` future/server-specific | `low` | `low` | `low` | `GoREveal` now has a real handoff surface; the next step is explicit MCP/workstation hardening, not parser work |
| Dynamic / symbolic / protected-binary sidecars | `low` native, `medium` orchestrated | `medium` | `medium` | `medium` | `medium` | `low` | the workstation wins here through `frida`, `angr`, `qiling`, `unicorn`, `uftrace`, `z3`; `GoREveal` should orchestrate, not absorb |

## Module-By-Module Reading

| GoREveal module | Competitive position vs universal workbenches |
| --- | --- |
| `schema` | stronger than universal tools because it gives a canonical machine-readable Go-native contract instead of workspace-local state |
| `core` | narrower than universal workbenches by design, but already stronger for bounded Go-native recovery on current `ELF`/`PE`/`Mach-O` slices |
| `engine` | strategically strong because it can add Go-specific enrichment such as peeling without contaminating raw recovery |
| `storage` | now unusually strong for a specialist RE product because review queues, handoff projections, and stored-run diffing are native features |
| `cmd/goreveal` | stronger for repeatable automation than GUI-first suites, weaker for rich analyst interaction |
| `plugins` | intentionally thin; they are bridges into universal workbenches, not competitors to them |

## Weighted Conclusion

Fresh weighted reading:
- `GoREveal` is already ahead of universal workbenches on Go-native truth, trust, and transfer workflow shape
- universal workbenches remain far ahead on analyst workspace depth, mutation UX, and generic RE breadth
- the latest intermediate rerun and timing sample do not weaken that reading: they reinforce that the next bottleneck is not local CLI efficiency, because the current `rclone-linux-amd64` timing sample is tightly clustered across all four Go-native tools, but clearer source-confidence semantics and stronger workstation handoff
- that does **not** argue for a broader parser or decompiler lane
- it argues for better handoff, interop, and operator workflow hardening on top of already-landed Go-native truth

## Recommended Next Move

The next best product move after this comparison is:
1. keep using `GoREveal` as the Go-native recovery and review system of record
2. harden explicit host-platform MCP and workstation handoff on top of `diff handoff sqlite`
3. only then decide whether another Go-native semantic slice materially outranks more interop polish
