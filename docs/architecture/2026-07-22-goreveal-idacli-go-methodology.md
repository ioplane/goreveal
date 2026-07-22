# Methodology: Go Binary Reverse Engineering with goreveal + idacli + IDA

**Date:** 2026-07-22
**Status:** research input; superseded as execution guidance by the
[RT1 product design](../superpowers/specs/2026-07-22-goreveal-rt1-product-design.md)
and [RT1 Horizon A plan](../superpowers/plans/2026-07-22-goreveal-rt1-horizon-a.md)
**Authors:** infra4 RE team

## Problem

IDA Pro's auto-analysis fails on large stripped Go binaries because:

1. **Go ABI ≠ C ABI** — IDA expects System V AMD64 ABI, Go uses register-based
   ABI (regabi since Go 1.17) with different register assignments
2. **Function boundaries wrong** — IDA's heuristic prologue/epilogue detection
   doesn't recognize Go's goroutine stack check (`cmp rsp, [r14+10h]; jbe morestack`)
3. **Function names missing** — stripped binary has no symbols; Go function
   names live in pclntab, which IDA's golang.so plugin may not fully recover
4. **Hex-Rays decompilation fails** — wrong function boundaries + Go ABI =
   Hex-Rays can't build valid control flow graph

Field observation in Teleport 18.10.0 (410 MB), pending the forced-plugin and
artifact-identity baseline defined by RT1-S1:
- IDA identified 248K functions; the earlier 54% ratio used the unverified
  goreveal entry count as its denominator
- goreveal emitted about 458,600 pclntab-derived function entries; completeness,
  uniqueness, boundary validity, and the denominator have not yet been proven
- Hex-Rays decompiled only 2 of 9 key license functions
- 3 functions couldn't be created (wrong boundaries)
- 4 functions decompilation failed (Go ABI)

These are research observations, not accepted product metrics. The 9/9 result
below is a target hypothesis, not a completed experiment.

## Solution: goreveal → idacli → IDA pipeline

### Proposed target architecture

The diagram describes a candidate workflow. It is not the current capability
map and must not override RT1 promotion gates.

```
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 1: GOREVEAL (Go-native, pclntab-derived provider evidence)  │
│                                                                     │
│  goreveal analyze <binary>                                         │
│    ├── core/pclntab → candidate function entries/ends/names         │
│    ├── core/types → truthful current type surface (mainly DWARF);   │
│    │                broad typelinks recovery is not yet claimed     │
│    ├── core/strings → embedded strings                             │
│    ├── core/runtime → moduledata, pclntab layout                   │
│    └── core/packages → Go package metadata                        │
│                                                                     │
│  goreveal export ida <binary>                                      │
│    └── IDAExport JSON (functions[], types[], strings[], runtime)   │
│                                                                     │
└──────────────────────────┬──────────────────────────────────────────┘
                           │ JSON file
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 2: IDA DATABASE CREATION (auto-analysis, may be incomplete) │
│                                                                     │
│  idat -A -B -o<binary.i64> <binary>                                │
│    ├── IDA auto-analysis (1h34m for 410MB)                         │
│    ├── 248K functions identified (54% — incomplete for Go)         │
│    └── Hex-Rays plugin loaded                                      │
│                                                                     │
└──────────────────────────┬──────────────────────────────────────────┘
                           │ .i64 database
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 3: GOREVEAL IMPORT (NEW — bridge goreveal → IDA)            │
│                                                                     │
│  idacli import-goreveal --input goreveal-export.json               │
│    For each goreveal function:                                      │
│      1. del_func(old_boundary)   — remove IDA's wrong boundary     │
│      2. add_func(entry, end)     — create with goreveal's boundary │
│      3. set_name(entry, name)    — rename with recovered Go name   │
│      4. set_cmt(entry, package)  — comment with Go package path    │
│      5. create_insn(entry)       — ensure code at entry           │
│                                                                     │
│    For each goreveal type:                                          │
│      6. apply Go type info to IDA type library                     │
│                                                                     │
│    For each goreveal string:                                        │
│      7. create_strlit(addr, value) — define string in IDB          │
│                                                                     │
└──────────────────────────┬──────────────────────────────────────────┘
                           │ corrected .i64
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 4: BATCH DECOMPILE (Hex-Rays now works!)                    │
│                                                                     │
│  idacli decompile --targets @license-functions.txt                  │
│    For each target address:                                         │
│      1. get_func_start(addr) — now finds correct function          │
│      2. hexrays.decompile(addr) — test whether boundary repair helps│
│      3. Write pseudocode to JSONL output                            │
│                                                                     │
│  Target: measure whether all 9 functions become decompilable        │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Why this works

| Problem | IDA alone | goreveal + idacli |
|---|---|---|
| Function boundaries | Heuristic, incomplete in the field observation | pclntab-derived provider evidence, subject to identity/range/host-conflict verification |
| Function names | Missing in the field observation | candidate pclntab names; coverage pending baseline validation |
| Function count | 248K observed | about 458,600 emitted entries; not a proven 100% denominator |
| Hex-Rays decompile | 2/9 observed | improvement is an RT1-S3 acceptance hypothesis |
| Go stack check | Not recognized in the observation | future classifier candidate after measured promotion, not current export truth |

### Implementation plan

#### Task 1: goreveal export enhancement

Add to `schema/export_ida.go`:

```go
// IDAFunction — add Go-specific fields
type IDAFunction struct {
    Name          string     `json:"name"`
    RefinedName   string     `json:"refined_name,omitempty"`
    Package       string     `json:"package,omitempty"`
    ImportPath    string     `json:"import_path,omitempty"`
    SourceFile    string     `json:"source_file,omitempty"`
    SourceLine    int        `json:"source_line,omitempty"`
    Autogenerated bool       `json:"autogenerated,omitempty"`
    Entry         uint64     `json:"entry"`
    End           uint64     `json:"end"`
    ModuleLocal   bool       `json:"module_local,omitempty"`
    Provenance    Provenance `json:"provenance"`
    // NEW: Go-specific metadata for idacli import
    GoPrologue    string     `json:"go_prologue,omitempty"`    // "stack_check" | "standard" | "unknown"
    GoStackCheck  bool       `json:"go_stack_check,omitempty"` // has goroutine stack check
    IsThunk       bool       `json:"is_thunk,omitempty"`       // small function, just JMP
    IsClosure     bool       `json:"is_closure,omitempty"`     // has bound variables
}
```

This was an early proposal for prologue detection. RT1 does not accept it as a
current claim: any classifier must be architecture/version aware, preserve raw
bytes and provenance, use explicit unsupported states, and pass corpus and
benchmark promotion gates before export.

#### Task 2: idacli new task `import-goreveal`

New file: `src/tasks/task_import_goreveal.h` + `.cpp`

```cpp
// Reads goreveal IDAExport JSON, applies to IDB:
// 1. For each function: fix boundaries + rename
// 2. For each type: apply Go type info
// 3. For each string: create string literal in IDB
// 4. For Go stack check: add comment "Go stack check (goroutine prologue)"

class task_import_goreveal : public task_base {
    task_outcome execute(const task_config &cfg) override;
    std::string_view name() const override { return "import-goreveal"; }
};
```

Implementation:
```cpp
// For each function in goreveal export:
for (const auto& fn : export_data["functions"]) {
    ea_t entry = parse_hex(fn["entry"]);
    ea_t end = parse_hex(fn["end"]);
    std::string name = fn["name"];

    // Remove existing function if boundaries don't match
    func_t* existing = get_func(entry);
    if (existing && (existing->start_ea != entry || existing->end_ea != end)) {
        del_func(existing->start_ea);
    }

    // Create function with correct boundaries from pclntab
    if (get_func_start(entry) == BADADDR) {
        create_insn(entry);  // ensure code, not data
        add_func(entry, end);
    }

    // Rename with recovered Go name
    set_name(entry, name.c_str(), SN_NOCHECK | SN_FORCE);

    // Comment with package info
    if (fn.contains("package")) {
        set_cmt(entry, fn["package"].get<std::string>().c_str(), false);
    }

    // Mark Go stack check
    if (fn.value("go_stack_check", false)) {
        set_cmt(entry, "Go goroutine stack check prologue", true);
    }
}
```

#### Task 3: Pipeline orchestration

Single command:
```bash
# Full Go binary RE pipeline
goreveal analyze binary && \
goreveal export ida binary > goreveal-export.json && \
idat -A -B -o binary.i64 binary && \
idat -A -S'-OIDACli:task=import-goreveal,input=goreveal-export.json' binary.i64 && \
idat -A -S'-OIDACli:task=decompile,targets=@targets.txt,output=decompiled.jsonl' binary.i64
```

Or as idacli multi-task pipeline:
```bash
idat -A -S'-OIDACli:task=import-goreveal+decompile,input=goreveal-export.json,targets=@targets.txt' binary.i64
```

#### Task 4: goreveal diff capability

Add `goreveal diff` to compare two goreveal analyses:
```bash
goreveal diff sqlite db1.db db2.db
# Shows: functions added/removed/changed, types changed, strings changed
```

This enables version comparison (18.7.2 vs 18.10.0) without IDA.

### Validation criteria

1. After a verified preview/apply, reconcile IDA actions against the exact
   artifact identity and every accepted goreveal entity; do not use raw count
   equality as proof of correctness
2. All license-related functions (FromPEM, checkLicense, IsExpired, etc.)
   should have correct boundaries and be decompilable
3. Function names should match goreveal's pclntab output
4. Measure Hex-Rays success on the fixed target set against the forced-plugin
   baseline; `>90%` is a research target, not a present result

### Alternative approaches considered

| Approach | Pros | Cons | Verdict |
|---|---|---|---|
| **goreveal + idacli** (proposed) | Go-native, pclntab-derived evidence, integrates with existing tools | Requires new idacli task | ✅ Recommended |
| IDA Go plugin (golang.so) | Built into IDA | Incomplete, may not load in batch mode | Insufficient |
| AlphaGolang plugin | Better Go support | Third-party, may not work in batch | Supplementary |
| Ghidra + Go plugin | Free, good Go support | Not Hex-Rays, different workflow | Alternative |
| Manual IDA + GoReSym | Works today | Manual, not scalable | Current fallback |

### Historical hypothesis sequence

This table is superseded by RT1. None of its rows is an active sprint or a
landed capability.

| Sprint | Tasks | Deliverable |
|---|---|---|
| Sprint A | goreveal: add prologue detection to export | Enhanced IDAExport with Go metadata |
| Sprint B | idacli: implement import-goreveal task | New task that fixes function boundaries |
| Sprint C | idacli: pipeline orchestration | Single-command Go RE workflow |
| Sprint D | goreveal: diff capability | Version comparison without IDA |
| Sprint E | Validation on Teleport 18.10.0 | Target: measure 9/9 license functions; no result claimed |
