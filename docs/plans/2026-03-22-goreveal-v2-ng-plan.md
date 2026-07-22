# GoREveal v2 (NG) — Go-Native RE Platform

> Status: exploratory only, not active roadmap baseline
> Note: this document is intentionally parked. The active product baseline is
> the [RT1 product design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
> plus [RT1 Horizon A](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md).

> **Historical note for agentic workers:** do not execute the checkboxes in
> this exploratory document. Re-evaluate any idea through the RT1 gates first.

**Goal:** Transform GoREveal from a recovery+transfer layer into a self-sufficient Go binary reverse engineering platform that fully replaces IDA/Ghidra/JEB/Binary Ninja for Go binary analysis.

**Architecture:** Incremental build on top of existing v1 (3,206 LOC preserved intact). New packages add disassembly (Capstone Go bindings), CFG construction, multi-level IR lifting (LLIL→MLIL→HLIL), Go-aware pseudocode generation, cross-references, and analyst interfaces (CLI+TUI+Web). Each sprint produces a working, testable product increment.

**Tech Stack:** Go 1.26, gapstone (Capstone Go bindings) for disasm, pure Go for IR/lifting/codegen, bubbletea/lipgloss for TUI, templ+htmx for Web UI, modernc.org/sqlite for persistence, Podman-first dev.

**Targets:** amd64 + arm64 (covers 99% of production Go binaries).

**Existing backlog integration:** This is a historical exploratory mapping of hypothetical future sprints 13-20, not the active sprint baseline. Existing backlog epics (Code Peeling, Version Tracking, Metadata Network) are governed by the RT1 design and plan.

---

## Strategy Change: From Transfer Layer to Self-Sufficient Platform

### Previous Strategy (v1)
> "GoREveal should not try to beat IDA/Ghidra/JEB/Binary Ninja as a generic RE suite. It should become the best Go-native recovery, trust, and transfer layer that plugs into them."

### New Strategy (v2 NG)
> "For Go binaries specifically, GoREveal becomes the single tool an analyst needs. It replaces IDA/Ghidra/JEB/Binary Ninja by providing Go-aware disassembly, decompilation, and analysis that no general-purpose RE tool can match. The thin adapter plugins become export-only for analysts who still want IDA/Ghidra as a secondary view."

### Why This Is Now Feasible
1. Go compiler (gc) is the only production compiler — one target, not thousands
2. Go ABI is documented and stable (register-based since Go 1.17)
3. pclntab gives function boundaries for free — no heuristic function detection needed
4. Go calling conventions are uniform — no varargs, no inline asm, no exceptions
5. GoREveal v1 already recovers functions, types, packages, strings, runtime metadata
6. The missing piece is disassembly → lifting → pseudocode — a bounded, achievable scope for Go-only

---

## Module Map Update (v2)

```
core/                   ← v1 recovery primitives (PRESERVED)
  ├─ buildinfo/
  ├─ format/
  ├─ functions/
  ├─ ingest/
  ├─ packages/
  ├─ pclntab/
  ├─ runtime/
  ├─ strings/
  └─ types/

schema/                 ← v1 canonical model (EXTENDED with IR types)

engine/                 ← v1 pipeline (EXTENDED with lift stages)

disasm/                 ← NEW: instruction decoding
  ├─ disasm.go          ← Capstone wrapper: Decode(addr, bytes) → []Instruction
  ├─ instruction.go     ← Instruction model: Op, Operands, Size, Addr
  └─ arch/
     ├─ amd64.go        ← x86-64 Capstone config + Go ABI register names
     └─ arm64.go        ← ARM64 Capstone config + Go ABI register names

cfg/                    ← NEW: control flow graph
  ├─ graph.go           ← BasicBlock, Edge, Graph types
  ├─ builder.go         ← Build CFG from instruction stream
  ├─ dominator.go       ← Dominator tree (needed for SSA)
  └─ loop.go            ← Loop detection

lift/                   ← NEW: multi-level IR lifting
  ├─ llil/              ← Low Level IL (register operations)
  │  ├─ ops.go          ← LLIL operation types (Load, Store, Call, Branch, Ret, ...)
  │  ├─ lift_amd64.go   ← x86-64 instruction → LLIL mapping
  │  ├─ lift_arm64.go   ← ARM64 instruction → LLIL mapping
  │  └─ func.go         ← LLIL function representation
  ├─ mlil/              ← Medium Level IL (variables, types, SSA)
  │  ├─ ops.go          ← MLIL operation types
  │  ├─ ssa.go          ← SSA construction from LLIL
  │  ├─ propagate.go    ← Constant propagation, type propagation
  │  ├─ goabi.go        ← Go ABI: register→parameter mapping, return values
  │  └─ func.go         ← MLIL function representation
  └─ hlil/              ← High Level IL (expressions, control flow recovery)
     ├─ ops.go          ← HLIL operation types (IfElse, ForLoop, Switch, ...)
     ├─ fold.go         ← Expression folding, dead code elimination
     ├─ structurize.go  ← CFG → structured control flow (if/for/switch)
     ├─ goaware.go      ← Go-specific patterns: defer, goroutine, interface call, slice ops
     └─ func.go         ← HLIL function representation

xref/                   ← NEW: cross-reference engine
  ├─ xref.go            ← CodeRef, DataRef, XrefDB types
  └─ builder.go         ← Build xrefs from disasm + LLIL

codegen/                ← NEW: pseudocode generation
  ├─ printer.go         ← HLIL → Go-like text
  ├─ formatter.go       ← Indentation, line breaks, comments
  └─ annotate.go        ← Provenance annotations, address comments

tui/                    ← NEW: terminal UI
  ├─ app.go             ← bubbletea main model
  ├─ views/
  │  ├─ disasm.go       ← Disassembly view with highlighting
  │  ├─ pseudocode.go   ← Decompiled code view
  │  ├─ functions.go    ← Function list with search/filter
  │  ├─ xrefs.go        ← Cross-reference navigation
  │  └─ strings.go      ← String browser
  └─ keybinds.go        ← Navigation keybindings

web/                    ← NEW: web UI
  ├─ server.go          ← HTTP server (net/http)
  ├─ handlers/          ← API + page handlers
  ├─ templates/         ← templ components
  └─ static/            ← htmx + minimal CSS

deobfuscation/          ← v1 (PRESERVED, extended)
storage/                ← v1 (PRESERVED, extended)
plugins/                ← v1 (PRESERVED, demoted to export-only)
cmd/goreveal/           ← v1 (EXTENDED with new commands)
```

### Dependency Rules Update (v2)
```
core → schema                           (v1, preserved)
disasm → schema                         (new)
cfg → disasm + schema                   (new)
lift/llil → disasm + cfg + schema       (new)
lift/mlil → lift/llil + core + schema   (new: uses pclntab/types for Go awareness)
lift/hlil → lift/mlil + schema          (new)
xref → disasm + lift/llil + schema      (new)
codegen → lift/hlil + schema            (new)
engine → core + disasm + cfg + lift + xref + codegen + deobfuscation + schema
tui → engine + schema                   (new)
web → engine + schema                   (new)
```

---

## Sprint Plan

### Phase 1: Foundation (Sprints 13-14) — "See the code"
Disassembly + CFG + xrefs. Analyst can navigate Go binary in terminal.

### Phase 2: Understanding (Sprints 15-17) — "Understand the code"
LLIL + MLIL + HLIL lifting. Go-aware pseudocode generation.

### Phase 3: Experience (Sprints 18-19) — "Work with the code"
TUI + Web UI. Full analyst workflow without external tools.

### Phase 4: Intelligence (Sprint 20+) — "Know the code"
Code Peeling, Version Tracking, Metadata Network — backlog epics on top of IR.

---

## Sprint 13: Disassembly Engine (amd64)

**Goal:** `goreveal disasm <binary> [--addr 0x...] [--func name]` produces annotated x86-64 disassembly enriched with pclntab function names and Go ABI register annotations.

**Est. LOC:** ~1,500 new
**Duration:** 3-4 days
**Dependencies:** Capstone (gapstone Go bindings)

### Task 13.1: Capstone Integration + Instruction Model

**Files:**
- Create: `disasm/instruction.go`
- Create: `disasm/disasm.go`
- Create: `disasm/arch/amd64.go`
- Create: `disasm/disasm_test.go`
- Modify: `go.mod` (add gapstone dependency)
- Modify: `deployments/docker/Containerfile.dev` (add libcapstone-dev)

- [ ] Add `github.com/knightsc/gapstone` to go.mod
- [ ] Add libcapstone-dev to Containerfile.dev, rebuild dev image
- [ ] Define `Instruction` struct: `Addr uint64, Size int, Mnemonic string, OpStr string, Bytes []byte, Op uint, Operands []Operand`
- [ ] Define `Operand` struct: `Type OperandType, Reg string, Imm int64, Mem MemOperand`
- [ ] Write `disasm.Decode(arch string, addr uint64, code []byte, count int) ([]Instruction, error)`
- [ ] Write `arch/amd64.go`: Capstone CS_ARCH_X86/CS_MODE_64 config, Go ABI register name mapping (RAX→ret0, RBX→arg5, etc.)
- [ ] Test: decode `[]byte{0xB8, 0x01, 0x00, 0x00, 0x00, 0xC3}` → `[MOV EAX,1; RET]`
- [ ] Test: decode real function prologue from test fixture
- [ ] `make lint && make test`
- [ ] Commit: `feat(disasm): add Capstone-based x86-64 instruction decoder`

### Task 13.2: Function Disassembler

**Files:**
- Create: `disasm/function.go`
- Create: `disasm/function_test.go`
- Modify: `disasm/disasm.go`

- [ ] Write `DisassembleFunction(path string, entry uint64, end uint64) ([]Instruction, error)` — reads ELF, seeks to function offset, decodes until end or RET
- [ ] Write `DisassembleFunctionByName(path string, name string) ([]Instruction, error)` — uses pclntab to find entry/end, then disassembles
- [ ] Enrich instructions with pclntab: annotate CALL targets with function names from recovery
- [ ] Test: disassemble `main.main` from test fixture → verify instruction count matches
- [ ] Test: verify CALL targets have symbolic names
- [ ] `make lint && make test`
- [ ] Commit: `feat(disasm): function-level disassembly with pclntab name resolution`

### Task 13.3: Schema + Engine + CLI Integration

**Files:**
- Modify: `schema/analysis.go` — add `Disassembly` section
- Modify: `engine/engine.go` — add disasm stage (optional, on-demand)
- Modify: `cmd/goreveal/main.go` — add `disasm` subcommand
- Create: `cmd/goreveal/internal/disasm_cmd.go`

- [ ] Add `schema.DisasmFunction` struct: `Name, Entry, End, Instructions []DisasmInstruction, Provenance`
- [ ] Add `schema.DisasmInstruction`: `Addr, Size, Mnemonic, OpStr, Bytes, CallTarget, Comment`
- [ ] Wire `disasm` subcommand: `goreveal disasm <binary> --addr 0x... | --func name [--json]`
- [ ] Output: annotated disassembly text (addr | bytes | mnemonic operands ; comment)
- [ ] JSON mode: structured output matching schema
- [ ] Test: `goreveal disasm fixture --func main.main` produces non-empty output
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add disasm subcommand with Go-aware annotations`

### Task 13.4: ARM64 Support

**Files:**
- Create: `disasm/arch/arm64.go`
- Create: `disasm/arch/arm64_test.go`

- [ ] Write ARM64 Capstone config (CS_ARCH_ARM64/CS_MODE_ARM)
- [ ] Go ABI register mapping for ARM64 (R0-R15 → parameters, R26=g, etc.)
- [ ] Test: decode ARM64 MOV+RET sequence
- [ ] Auto-detect arch from ELF header in `DisassembleFunction`
- [ ] `make lint && make test`
- [ ] Commit: `feat(disasm): add ARM64 support`

---

## Sprint 14: Control Flow Graph + Cross-References

**Goal:** `goreveal cfg <binary> --func name` shows basic blocks, edges, dominators. `goreveal xrefs <binary> --func name` shows callers/callees.

**Est. LOC:** ~2,000 new
**Duration:** 3-4 days

### Task 14.1: CFG Builder

**Files:**
- Create: `cfg/graph.go` — BasicBlock, Edge, Graph types
- Create: `cfg/builder.go` — Build CFG from instruction stream
- Create: `cfg/builder_test.go`

- [ ] Define `BasicBlock`: `ID int, Addr uint64, End uint64, Instructions []disasm.Instruction, Successors []int, Predecessors []int`
- [ ] Define `Edge`: `From, To int, Type EdgeType` (EdgeFallthrough, EdgeBranch, EdgeCall, EdgeRet)
- [ ] Define `Graph`: `Blocks []BasicBlock, Edges []Edge, Entry int`
- [ ] Write `BuildCFG(instructions []disasm.Instruction) Graph` — split at branch targets, build edges
- [ ] Test: linear function (no branches) → 1 block
- [ ] Test: if-else pattern → 3+ blocks with correct edges
- [ ] Test: loop pattern → back-edge detected
- [ ] `make lint && make test`
- [ ] Commit: `feat(cfg): basic block and control flow graph construction`

### Task 14.2: Dominator Tree

**Files:**
- Create: `cfg/dominator.go`
- Create: `cfg/dominator_test.go`

- [ ] Implement Lengauer-Tarjan dominator tree algorithm
- [ ] `DominatorTree(graph Graph) map[int]int` — block → immediate dominator
- [ ] `DominanceFrontier(graph Graph, idom map[int]int) map[int][]int` — needed for SSA
- [ ] Test: diamond CFG → correct dominators
- [ ] Test: loop → correct dominator (loop header dominates body)
- [ ] `make lint && make test`
- [ ] Commit: `feat(cfg): Lengauer-Tarjan dominator tree and dominance frontiers`

### Task 14.3: Cross-Reference Engine

**Files:**
- Create: `xref/xref.go` — XrefDB, CodeRef, DataRef types
- Create: `xref/builder.go` — Build xrefs from disasm
- Create: `xref/builder_test.go`

- [ ] Define `CodeRef`: `From, To uint64, Type RefType, FromFunc, ToFunc string`
- [ ] Define `DataRef`: `From uint64, To uint64, Size int, FromFunc string`
- [ ] Define `XrefDB`: `CodeRefs []CodeRef, DataRefs []DataRef` + lookup methods
- [ ] Write `BuildXrefs(path string, functions []schema.Function) (*XrefDB, error)` — disassemble all functions, extract CALL/JMP targets + memory references
- [ ] `CallersOf(funcName string) []CodeRef`
- [ ] `CalleesOf(funcName string) []CodeRef`
- [ ] `DataRefsFrom(funcName string) []DataRef`
- [ ] Test: fixture `main.main` calls known functions → callee list correct
- [ ] Test: bidirectional — callee's callers include main.main
- [ ] `make lint && make test`
- [ ] Commit: `feat(xref): cross-reference engine with caller/callee tracking`

### Task 14.4: CLI Integration (cfg + xrefs)

**Files:**
- Create: `cmd/goreveal/internal/cfg_cmd.go`
- Create: `cmd/goreveal/internal/xrefs_cmd.go`
- Modify: `cmd/goreveal/main.go`

- [ ] `goreveal cfg <binary> --func name [--dot] [--json]` — show blocks + edges, optional DOT graph output
- [ ] `goreveal xrefs <binary> --func name [--callers] [--callees] [--json]` — show cross-references
- [ ] Text output: tree-style caller/callee display with addresses
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add cfg and xrefs subcommands`

---

## Sprint 15: Low-Level IL (LLIL)

**Goal:** Lift x86-64/ARM64 instructions into architecture-independent register operations. First IR level.

**Est. LOC:** ~3,500 new
**Duration:** 4-5 days
**This is the hardest sprint** — instruction semantics mapping.

### Task 15.1: LLIL Operation Model

**Files:**
- Create: `lift/llil/ops.go` — LLIL operation types
- Create: `lift/llil/func.go` — LLIL function + basic block representation
- Create: `lift/llil/ops_test.go`

- [ ] Define LLIL operations enum: `SetReg, GetReg, Load, Store, Add, Sub, Mul, Div, And, Or, Xor, Shl, Shr, Cmp, Test, Call, TailCall, Ret, Jump, BranchIf, Nop, Undef, Push, Pop, Lea, SignExtend, ZeroExtend, ...`
- [ ] Define `LLILExpr`: tree-based expression `{Op LLILOp, Left *LLILExpr, Right *LLILExpr, Reg string, Imm int64, Addr uint64, Size int}`
- [ ] Define `LLILInstr`: `{Addr uint64, Expr LLILExpr, Original *disasm.Instruction}`
- [ ] Define `LLILFunc`: `{Name string, Entry uint64, Blocks []LLILBlock}`, `LLILBlock`: `{ID int, Instrs []LLILInstr, Successors []int}`
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/llil): LLIL operation model and function representation`

### Task 15.2: x86-64 → LLIL Lifter

**Files:**
- Create: `lift/llil/lift_amd64.go`
- Create: `lift/llil/lift_amd64_test.go`

- [ ] Write `LiftAMD64(instrs []disasm.Instruction) []LLILInstr`
- [ ] Map common Go compiler output instructions (80% coverage):
  - Data movement: MOV, LEA, MOVZX, MOVSX, CMOV*
  - Arithmetic: ADD, SUB, IMUL, INC, DEC, NEG
  - Logic: AND, OR, XOR, NOT, SHL, SHR, SAR
  - Comparison: CMP, TEST
  - Control flow: JMP, Jcc, CALL, RET
  - Stack: PUSH, POP (rare in Go 1.17+ but still exists)
  - Go-specific: NOP padding, PCDATA/FUNCDATA markers
- [ ] Unrecognized instructions → `LLILOp_Undef` with original bytes preserved
- [ ] Test: `MOV EAX, 1; RET` → `SetReg(RAX, Imm(1)); Ret()`
- [ ] Test: `CMP RAX, 0; JE label; ...` → `BranchIf(Cmp(GetReg(RAX), Imm(0)), target)`
- [ ] Test: `CALL 0x12345` → `Call(Addr(0x12345))` with symbolic name if known
- [ ] Test: lift real fixture function → verify no panics, reasonable instruction count
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/llil): x86-64 to LLIL lifter with Go compiler pattern coverage`

### Task 15.3: ARM64 → LLIL Lifter

**Files:**
- Create: `lift/llil/lift_arm64.go`
- Create: `lift/llil/lift_arm64_test.go`

- [ ] Write `LiftARM64(instrs []disasm.Instruction) []LLILInstr`
- [ ] Map ARM64 instructions commonly emitted by Go compiler:
  - MOV, MOVZ, MOVK, LDR, STR, LDP, STP
  - ADD, SUB, MUL, CMP, TST
  - B, B.cond, BL, RET, CBZ, CBNZ, TBZ, TBNZ
- [ ] Test: `MOV X0, #1; RET` → `SetReg(X0, Imm(1)); Ret()`
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/llil): ARM64 to LLIL lifter`

### Task 15.4: LLIL CLI + Schema

**Files:**
- Modify: `schema/analysis.go` — add LLIL section
- Create: `cmd/goreveal/internal/llil_cmd.go`
- Modify: `cmd/goreveal/main.go`

- [ ] `goreveal llil <binary> --func name [--json]` — show LLIL for a function
- [ ] Text output: `0x1234: SetReg(RAX, Add(GetReg(RBX), Imm(8)))`
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add llil subcommand for low-level IL inspection`

---

## Sprint 16: Medium-Level IL (MLIL) + SSA

**Goal:** Lift LLIL to MLIL: registers → variables, Go ABI parameter mapping, SSA form, constant propagation, type propagation.

**Est. LOC:** ~3,500 new
**Duration:** 4-5 days

### Task 16.1: MLIL Operation Model

**Files:**
- Create: `lift/mlil/ops.go`
- Create: `lift/mlil/func.go`

- [ ] Define MLIL operations: `Assign, Load, Store, Call, Return, Branch, If, Phi, ...`
- [ ] Define `Variable`: `{Name string, Type string, SSAVersion int, Source VarSource}` where VarSource = Param/Local/Return/Temp
- [ ] Define `MLILExpr`: `{Op, Var *Variable, Left *MLILExpr, Right *MLILExpr, Imm int64, CallTarget string, Args []MLILExpr}`
- [ ] Define `MLILFunc`: `{Name, Params []Variable, Returns []Variable, Locals []Variable, Blocks []MLILBlock}`
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/mlil): MLIL operation model with variables and SSA support`

### Task 16.2: Go ABI Parameter Mapping

**Files:**
- Create: `lift/mlil/goabi.go`
- Create: `lift/mlil/goabi_test.go`

- [ ] Write Go register ABI mapping (Go 1.17+):
  - Integer params: RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11 (amd64)
  - Integer returns: RAX, RBX, RCX, RDI, RSI, R8, R9, R10, R11
  - Float params/returns: X0-X14
  - ARM64: R0-R15 for int, F0-F15 for float
- [ ] `MapParams(sig FuncSignature, llilFunc LLILFunc) []Variable` — map register reads in function prologue to named parameters
- [ ] `MapReturns(sig FuncSignature, llilFunc LLILFunc) []Variable` — map register writes before RET to named returns
- [ ] Use pclntab + recovered types for signature hints
- [ ] Test: function with 2 int params → var_arg0 (RAX), var_arg1 (RBX)
- [ ] Test: function returning error → var_ret0, var_ret1 (error)
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/mlil): Go register ABI parameter and return mapping`

### Task 16.3: SSA Construction

**Files:**
- Create: `lift/mlil/ssa.go`
- Create: `lift/mlil/ssa_test.go`

- [ ] Implement SSA construction from LLIL + dominator tree:
  1. Insert phi nodes at dominance frontiers for each variable
  2. Rename variables with SSA versions
  3. Build def-use chains
- [ ] `ConstructSSA(llilFunc LLILFunc, cfg cfg.Graph, idom map[int]int) MLILFunc`
- [ ] Test: single assignment → version 0
- [ ] Test: reassignment in branch → phi node at merge point
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/mlil): SSA construction with phi nodes and def-use chains`

### Task 16.4: Constant + Type Propagation

**Files:**
- Create: `lift/mlil/propagate.go`
- Create: `lift/mlil/propagate_test.go`

- [ ] Constant propagation: fold `var_0 = 5; var_1 = var_0 + 3` → `var_1 = 8`
- [ ] Type propagation: use Go recovered types to annotate variables (`*license.License`, `[]byte`, etc.)
- [ ] Dead store elimination: remove unused assignments
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/mlil): constant propagation, type propagation, dead store elimination`

### Task 16.5: MLIL CLI

- [ ] `goreveal mlil <binary> --func name [--ssa] [--json]`
- [ ] Text output: `var_arg0 = param0 (type: *x509.Certificate); var_1 = call x509.Verify(var_arg0, var_pool)`
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add mlil subcommand with SSA and type annotations`

---

## Sprint 17: High-Level IL (HLIL) + Pseudocode

**Goal:** `goreveal decompile <binary> --func name` produces Go-like pseudocode. **The headline feature.**

**Est. LOC:** ~3,000 new
**Duration:** 4-5 days

### Task 17.1: Control Flow Structurization

**Files:**
- Create: `lift/hlil/structurize.go`
- Create: `lift/hlil/structurize_test.go`

- [ ] Implement structural analysis: CFG → structured control flow
  - If-else: 2-way branch with merge → `IfElse{Cond, Then, Else}`
  - For loop: back-edge + header → `ForLoop{Init, Cond, Post, Body}`
  - Switch: multi-target branch → `Switch{Expr, Cases}`
  - Short-circuit: `&&` and `||` patterns
  - Infinite loop: unconditional back-edge → `ForLoop{nil, nil, nil, Body}`
- [ ] Handle irreducible control flow: fall back to goto
- [ ] Test: diamond pattern → IfElse
- [ ] Test: loop pattern → ForLoop
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/hlil): control flow structurization (if/for/switch recovery)`

### Task 17.2: Go-Aware Patterns

**Files:**
- Create: `lift/hlil/goaware.go`
- Create: `lift/hlil/goaware_test.go`

- [ ] Detect and annotate Go-specific patterns in HLIL:
  - `defer` calls: `runtime.deferproc` + `runtime.deferreturn` → `Defer{Call}`
  - `go` routines: `runtime.newproc` → `Go{Call}`
  - Interface calls: indirect call through itab → `InterfaceCall{Iface, Method, Args}`
  - Slice operations: `runtime.growslice` → `SliceAppend{...}`
  - Map operations: `runtime.mapassign`, `runtime.mapaccess` → `MapSet`, `MapGet`
  - String operations: `runtime.concatstrings` → `StringConcat`
  - Nil check: `runtime.panicnil` guard → annotate as nil check
  - Bounds check: `runtime.panicIndex` guard → annotate as bounds check
- [ ] Test: defer pattern recognized
- [ ] Test: interface call pattern recognized
- [ ] `make lint && make test`
- [ ] Commit: `feat(lift/hlil): Go-aware pattern detection (defer, go, interface, slice, map)`

### Task 17.3: Expression Folding + Pseudocode Generation

**Files:**
- Create: `lift/hlil/fold.go`
- Create: `lift/hlil/ops.go`
- Create: `codegen/printer.go`
- Create: `codegen/formatter.go`
- Create: `codegen/printer_test.go`

- [ ] Expression folding: collapse nested assignments into compound expressions
- [ ] Dead code elimination: remove panic paths, bounds checks (optional flag to show)
- [ ] `codegen.Print(hlilFunc HLILFunc) string` → Go-like pseudocode
- [ ] Output format:
  ```go
  func checkLicense(lf *licensefile.LicenseFile) error {
      cert := lf.KeyPair.Cert           // 0x1294203c
      err := cert.Verify(caPool)        // 0x12942080
      if err != nil {                   // 0x12942098
          return err
      }
      payload := getRawPayload(cert)    // 0x129420b0
      ...
  }
  ```
- [ ] Address comments on each line (configurable)
- [ ] Provenance annotations: `// recovered: pclntab + MLIL type propagation`
- [ ] Test: simple function → valid Go-like output
- [ ] Test: function with if-else → structured output
- [ ] `make lint && make test`
- [ ] Commit: `feat(codegen): Go-like pseudocode generator from HLIL`

### Task 17.4: `decompile` CLI Command

**Files:**
- Create: `cmd/goreveal/internal/decompile_cmd.go`
- Modify: `cmd/goreveal/main.go`

- [ ] `goreveal decompile <binary> --func name [--no-comments] [--show-bounds-checks] [--json]`
- [ ] Full pipeline: ingest → pclntab → disasm → CFG → LLIL → MLIL → HLIL → codegen
- [ ] `goreveal decompile <binary> --all` — decompile all user functions (with code peeling)
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add decompile subcommand — Go-like pseudocode from binary`

---

## Sprint 18: Terminal UI (TUI)

**Goal:** `goreveal tui <binary>` — interactive terminal-based analyst interface with disassembly, pseudocode, function list, xrefs, strings.

**Est. LOC:** ~2,500 new
**Duration:** 3-4 days
**Dependencies:** bubbletea, lipgloss, bubbles

### Task 18.1: TUI Framework + Function Browser

**Files:**
- Create: `tui/app.go`
- Create: `tui/views/functions.go`
- Create: `tui/keybinds.go`
- Modify: `go.mod` (add bubbletea, lipgloss)

- [ ] Main bubbletea app model with tab/split layout
- [ ] Function list view: searchable, filterable (user/stdlib/runtime)
- [ ] Go-to-function by name or address
- [ ] Keyboard navigation: j/k scroll, Enter select, / search, q quit
- [ ] `make lint && make test`
- [ ] Commit: `feat(tui): interactive function browser with search and filter`

### Task 18.2: Disassembly + Pseudocode Views

**Files:**
- Create: `tui/views/disasm.go`
- Create: `tui/views/pseudocode.go`

- [ ] Disassembly view: syntax-highlighted instructions, address column, bytes column
- [ ] Pseudocode view: Go-like code with syntax highlighting
- [ ] Toggle between disasm ↔ pseudocode with Tab
- [ ] Scroll sync: selecting a line in pseudocode highlights corresponding disasm range
- [ ] `make lint && make test`
- [ ] Commit: `feat(tui): disassembly and pseudocode views with synchronized navigation`

### Task 18.3: Xref Navigation + String Browser

**Files:**
- Create: `tui/views/xrefs.go`
- Create: `tui/views/strings.go`

- [ ] Xref panel: show callers/callees of selected function
- [ ] Click/Enter on xref → navigate to that function
- [ ] String browser: searchable string list with addresses
- [ ] Click/Enter on string → navigate to referencing function
- [ ] `make lint && make test`
- [ ] Commit: `feat(tui): xref navigation and string browser`

### Task 18.4: CLI Entry Point

**Files:**
- Modify: `cmd/goreveal/main.go`

- [ ] `goreveal tui <binary>` — launch TUI
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add tui subcommand`

---

## Sprint 19: Web UI

**Goal:** `goreveal web <binary> [--port 8080]` — browser-based analyst interface. Shareable, multi-user capable.

**Est. LOC:** ~3,000 new
**Duration:** 3-4 days
**Dependencies:** templ, htmx (no heavy JS framework)

### Task 19.1: HTTP Server + API

**Files:**
- Create: `web/server.go`
- Create: `web/handlers/api.go`
- Create: `web/handlers/pages.go`

- [ ] `GET /api/functions` — JSON function list with search/filter
- [ ] `GET /api/disasm/:addr` — disassembly for function at address
- [ ] `GET /api/decompile/:addr` — pseudocode for function
- [ ] `GET /api/xrefs/:addr` — callers/callees
- [ ] `GET /api/strings` — string list with search
- [ ] Page handlers: index, function detail, string browser
- [ ] `make lint && make test`
- [ ] Commit: `feat(web): HTTP server with REST API for binary analysis`

### Task 19.2: Frontend (templ + htmx)

**Files:**
- Create: `web/templates/*.templ`
- Create: `web/static/style.css`

- [ ] Layout: sidebar (functions) + main area (disasm/pseudocode) + bottom (strings/xrefs)
- [ ] htmx: click function → load disasm/pseudocode without page reload
- [ ] htmx: search box → filter functions/strings dynamically
- [ ] Tab switch: disasm ↔ pseudocode
- [ ] Syntax highlighting via CSS classes (no JS highlighting library)
- [ ] `make lint && make test`
- [ ] Commit: `feat(web): browser UI with function navigation and code views`

### Task 19.3: CLI Entry Point

- [ ] `goreveal web <binary> [--port 8080] [--open]` — start web server, optionally open browser
- [ ] `make lint && make test`
- [ ] Commit: `feat(cli): add web subcommand`

---

## Sprint 20: Code Peeling + Backlog Epic Integration

**Goal:** Integrate existing backlog epics on top of v2 IR. User-code isolation, version tracking foundation.

**Est. LOC:** ~2,000 new
**Duration:** 3-4 days

### Task 20.1: Code Peeling (User-Code Isolation)

- [ ] Classify functions into layers: `user`, `third-party`, `stdlib`, `runtime`
- [ ] `goreveal analyze --peel <binary>` — show only user code by default
- [ ] Filter in TUI/Web: toggle layers on/off
- [ ] Export: `export ida --peel` — only export user functions to IDA

### Task 20.2: Version Tracking Foundation

- [ ] `goreveal diff <binary-a> <binary-b>` — compare two Go builds
- [ ] Function matching: by name, by pclntab entry, by CFG similarity
- [ ] Output: added/removed/modified functions with diff summary
- [ ] Foundation for markup transfer in future sprints

### Task 20.3: Metadata Network Schema

- [ ] Define metadata exchange format (protobuf): function signatures, type info, package structure
- [ ] `goreveal export metadata <binary>` — export high-confidence metadata
- [ ] `goreveal import metadata <file>` — import and apply to current analysis
- [ ] Foundation for Lumina/WARP-style network in future sprints

---

## Sprint Estimates Summary

| Sprint | Name | Est. LOC | Duration | Cumulative |
|--------|------|----------|----------|------------|
| 13 | Disassembly Engine | 1,500 | 3-4 days | 1,500 |
| 14 | CFG + Cross-References | 2,000 | 3-4 days | 3,500 |
| 15 | Low-Level IL (LLIL) | 3,500 | 4-5 days | 7,000 |
| 16 | Medium-Level IL (MLIL+SSA) | 3,500 | 4-5 days | 10,500 |
| 17 | High-Level IL + Pseudocode | 3,000 | 4-5 days | 13,500 |
| 18 | Terminal UI (TUI) | 2,500 | 3-4 days | 16,000 |
| 19 | Web UI | 3,000 | 3-4 days | 19,000 |
| 20 | Code Peeling + Backlog Integration | 2,000 | 3-4 days | 21,000 |
| **Total** | | **~21,000** | **~28-36 days** | |

Plus existing v1 code (3,206 LOC) → total project: **~24,000 LOC**

---

## Quality Gates (QA)

Every sprint must pass:

1. **golangci-lint** — 0 issues (existing 55+ linter policy)
2. **Unit tests** — all new packages have `*_test.go`
3. **Fixture tests** — real Go binary analysis produces expected output
4. **Snapshot tests** — golden output comparison for regression
5. **Differential tests** — compare with GoReSym/redress/gore baselines
6. **No v1 regressions** — all existing `make test` stays green
7. **Container-first** — all dev/test in Podman

### New QA criteria for v2:
8. **Lift accuracy** — LLIL/MLIL/HLIL output for fixture functions compared against IDA/Ghidra pseudocode
9. **Decompile readability** — generated pseudocode reviewed for Go-likeness
10. **Performance** — full analysis of 100MB Go binary < 60 seconds (excluding first Capstone init)

---

## Risk Register

| Risk | Impact | Mitigation |
|------|--------|------------|
| Capstone CGo dependency | Build complexity, cross-compile | Pure Go disasm fallback (golang.org/x/arch) |
| x86-64 instruction coverage | Missing LLIL for rare instructions | Undef fallback, incremental coverage |
| SSA correctness | Wrong data flow | Extensive fixture tests vs IDA/Ghidra |
| Control flow structurization | Goto fallback too frequent | Iterative improvement, Go CFG is simpler than C++ |
| Stripped binary support | No pclntab in some builds | Fallback to address-based function detection |
| ARM64 coverage | Go compiler uses subset of ARM64 | Focus on gc-emitted subset only |
| TUI complexity | bubbletea state management | Keep views independent, simple model |

---

## Milestone Releases

| Version | Sprint | Headline |
|---------|--------|----------|
| v2.0.0-alpha.1 | 13 | Disassembly engine (amd64+arm64) |
| v2.0.0-alpha.2 | 14 | CFG + cross-references |
| v2.0.0-beta.1 | 15-16 | LLIL + MLIL lifting |
| v2.0.0-beta.2 | 17 | **Go-like pseudocode generation** |
| v2.0.0-rc.1 | 18-19 | TUI + Web UI |
| v2.0.0 | 20 | Code Peeling, production release |
