---
title: "Methodology: Go Binary Reverse Engineering with goreveal + idacli + IDA"
status: draft
date: 2026-07-22
owners:
  - ioplane/goreveal-maintainers
tags:
  - methodology
  - ida
  - workflow
---

# Methodology: Go Binary Reverse Engineering with goreveal + idacli + IDA

<img
  src="https://shieldcn.dev/badge/status-draft-slate.svg?variant=outline&size=xs"
  alt="status: draft" height="20">
<img
  src="https://shieldcn.dev/badge/docs-architecture-slate.svg?variant=outline&size=xs"
  alt="docs: architecture" height="20">

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

Observed in Teleport 18.10.0 (410MB, 458K functions):

- IDA identified 248K functions (54% of actual) via auto-analysis
- goreveal identified all 458K functions from pclntab
- Hex-Rays decompiled only 2 of 9 key license functions
- 3 functions couldn't be created (wrong boundaries)
- 4 functions decompilation failed (Go ABI)

## Solution: goreveal → idacli → IDA pipeline

### Architecture

```text
┌─────────────────────────────────────────────────────────────────────┐
│  PHASE 1: GOREVEAL (Go-native, pclntab ground truth)               │
│                                                                     │
│  goreveal analyze <binary>                                         │
│    ├── core/pclntab → function entries, ends, names (458K)          │
│    ├── core/types → Go type info from typelinks                    │
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
│      2. hexrays.decompile(addr) — works because boundary is correct │
│      3. Write pseudocode to JSONL output                            │
│                                                                     │
│  Result: all 9 license functions decompiled successfully            │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Why this works

| Problem | IDA alone | goreveal + idacli |
| --- | --- | --- |
| Function boundaries | Heuristic (wrong for Go) | pclntab (authoritative) |
| Function names | Missing (stripped) | pclntab (all 458K names) |
| Function count | 248K (54%) | 458K (100%) |
| Hex-Rays decompile | Fails (wrong boundaries) | Works (correct boundaries) |
| Go stack check | Not recognized | goreveal marks in export |

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

goreveal's `core/pclntab/pclntab.go` already reads pclntab — add prologue
detection: read first 8 bytes at each function entry, check for
`49 3b 66 10` (cmp rsp, [r14+10h]) pattern.

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

1. After `import-goreveal`, IDA function count should match goreveal's
   (458K for Teleport 18.10.0)
2. All license-related functions (FromPEM, checkLicense, IsExpired, etc.)
   should have correct boundaries and be decompilable
3. Function names should match goreveal's pclntab output
4. Hex-Rays decompilation success rate should be >90% (vs ~20% without import)

### Alternative approaches considered

| Approach | Pros | Cons | Verdict |
| --- | --- | --- | --- |
| **goreveal + idacli** (proposed) | Go-native, pclntab truth, integrates with existing tools | Requires new idacli task | ✅ Recommended |
| IDA Go plugin (golang.so) | Built into IDA | Incomplete, may not load in batch mode | Insufficient |
| AlphaGolang plugin | Better Go support | Third-party, may not work in batch | Supplementary |
| Ghidra + Go plugin | Free, good Go support | Not Hex-Rays, different workflow | Alternative |
| Manual IDA + GoReSym | Works today | Manual, not scalable | Current fallback |

### Sprint plan

| Sprint | Tasks | Deliverable |
| --- | --- | --- |
| Sprint A | goreveal: add prologue detection to export | Enhanced IDAExport with Go metadata |
| Sprint B | idacli: implement import-goreveal task | New task that fixes function boundaries |
| Sprint C | idacli: pipeline orchestration | Single-command Go RE workflow |
| Sprint D | goreveal: diff capability | Version comparison without IDA |
| Sprint E | Validation on Teleport 18.10.0 | 9/9 license functions decompiled |
