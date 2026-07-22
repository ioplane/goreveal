package internalcmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dantte-lp/goreveal/schema"
	storesqlite "github.com/dantte-lp/goreveal/storage/sqlite"
)

func TestAnalysisPolicyCoversEveryCurrentAnalysisCommand(t *testing.T) {
	t.Parallel()

	commands := []analysisCommand{
		analysisCommandAnalyze,
		analysisCommandInspectFunctions,
		analysisCommandInspectPackages,
		analysisCommandInspectTypes,
		analysisCommandInspectStrings,
		analysisCommandInspectRuntime,
		analysisCommandInspectPeeling,
		analysisCommandPeel,
		analysisCommandSourceTree,
		analysisCommandDeobfuscate,
		analysisCommandExportSQLite,
		analysisCommandExportIDAV1,
		analysisCommandExportGhidraV1,
	}

	if len(commandPolicies) != len(commands) {
		t.Fatalf("commandPolicies has %d entries, want %d: %#v", len(commandPolicies), len(commands), commandPolicies)
	}
	for _, command := range commands {
		if _, ok := commandPolicies[command]; !ok {
			t.Fatalf("commandPolicies missing %q", command)
		}
	}
}

func TestEnforceAnalysisPolicy(t *testing.T) {
	t.Parallel()

	available := func(stages ...schema.AnalysisStage) schema.Analysis {
		diagnostics := make([]schema.StageDiagnostic, 0, len(stages))
		for _, stage := range stages {
			diagnostics = append(diagnostics, schema.StageDiagnostic{Stage: stage, Status: schema.StageStatusAvailable})
		}
		return schema.Analysis{Diagnostics: diagnostics}
	}

	tests := []struct {
		name     string
		command  analysisCommand
		analysis schema.Analysis
		wantCode string
	}{
		{
			name:    "analyze permits failed partial result",
			command: analysisCommandAnalyze,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusFailed, Code: "stage_failed"},
			}},
		},
		{
			name:     "inspect functions rejects missing stage",
			command:  analysisCommandInspectFunctions,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{}},
			wantCode: "stage_not_attempted",
		},
		{
			name:     "inspect functions accepts available",
			command:  analysisCommandInspectFunctions,
			analysis: available(schema.AnalysisStageFunctions),
		},
		{
			name:    "inspect functions rejects failed",
			command: analysisCommandInspectFunctions,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusFailed, Code: "fixture_failed"},
			}},
			wantCode: "fixture_failed",
		},
		{
			name:    "inspect types accepts unavailable as empty",
			command: analysisCommandInspectTypes,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageTypes, Status: schema.StageStatusUnavailable, Code: "dwarf_not_found"},
			}},
		},
		{
			name:    "inspect types rejects unsupported",
			command: analysisCommandInspectTypes,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageTypes, Status: schema.StageStatusUnsupported, Code: "dwarf_unsupported_container"},
			}},
			wantCode: "dwarf_unsupported_container",
		},
		{
			name:    "inspect runtime reports unavailable",
			command: analysisCommandInspectRuntime,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageRuntime, Status: schema.StageStatusUnavailable, Code: "runtime_metadata_not_found"},
			}},
			wantCode: "runtime_metadata_not_found",
		},
		{
			name:     "inspect packages requires derived stage",
			command:  analysisCommandInspectPackages,
			analysis: available(schema.AnalysisStageFunctions),
			wantCode: "stage_not_attempted",
		},
		{
			name:    "inspect peeling accepts both available stages",
			command: analysisCommandInspectPeeling,
			analysis: available(
				schema.AnalysisStageFunctions,
				schema.AnalysisStagePeeling,
			),
		},
		{
			name:    "deobfuscate requires all raw and refinement stages",
			command: analysisCommandDeobfuscate,
			analysis: available(
				schema.AnalysisStageFunctions,
				schema.AnalysisStagePackages,
				schema.AnalysisStageTypes,
				schema.AnalysisStageStrings,
			),
			wantCode: "stage_not_attempted",
		},
		{
			name:    "sqlite accepts optional unavailable stage",
			command: analysisCommandExportSQLite,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusAvailable},
				{Stage: schema.AnalysisStageTypes, Status: schema.StageStatusUnavailable, Code: "dwarf_not_found"},
			}},
		},
		{
			name:    "sqlite rejects any failed attempted stage",
			command: analysisCommandExportSQLite,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusAvailable},
				{Stage: schema.AnalysisStageRuntime, Status: schema.StageStatusFailed, Code: "stage_failed"},
			}},
			wantCode: "stage_failed",
		},
		{
			name:    "IDA v1 accepts optional unsupported stage",
			command: analysisCommandExportIDAV1,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusAvailable},
				{Stage: schema.AnalysisStageTypes, Status: schema.StageStatusUnsupported, Code: "dwarf_unsupported_container"},
			}},
		},
		{
			name:    "IDA v1 rejects failed projected stage",
			command: analysisCommandExportIDAV1,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusAvailable},
				{Stage: schema.AnalysisStageStrings, Status: schema.StageStatusFailed, Code: "stage_failed"},
			}},
			wantCode: "stage_failed",
		},
		{
			name:    "Ghidra v1 rejects unavailable functions",
			command: analysisCommandExportGhidraV1,
			analysis: schema.Analysis{Diagnostics: []schema.StageDiagnostic{
				{Stage: schema.AnalysisStageFunctions, Status: schema.StageStatusUnavailable, Code: "pclntab_not_found"},
			}},
			wantCode: "pclntab_not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := enforceAnalysisPolicy(tt.command, tt.analysis)
			if tt.wantCode == "" {
				if err != nil {
					t.Fatalf("enforceAnalysisPolicy() error = %v", err)
				}
				return
			}

			var policyError *analysisPolicyError
			if !errors.As(err, &policyError) {
				t.Fatalf("enforceAnalysisPolicy() error = %v, want *analysisPolicyError", err)
			}
			if policyError.Code != tt.wantCode {
				t.Fatalf("analysisPolicyError.Code = %q, want %q", policyError.Code, tt.wantCode)
			}
		})
	}
}

func TestEnforceAnalysisPolicyRejectsUnknownCommand(t *testing.T) {
	t.Parallel()

	err := enforceAnalysisPolicy(analysisCommand("future-command"), schema.Analysis{Diagnostics: []schema.StageDiagnostic{}})
	var policyError *analysisPolicyError
	if !errors.As(err, &policyError) {
		t.Fatalf("enforceAnalysisPolicy() error = %v, want *analysisPolicyError", err)
	}
	if policyError.Code != "unknown_analysis_policy" {
		t.Fatalf("analysisPolicyError.Code = %q, want unknown_analysis_policy", policyError.Code)
	}
}

func TestRunAnalyze(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.bin")
	if err := os.WriteFile(path, []byte{'M', 'Z', 0x90, 0x00}, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var out strings.Builder
	if err := RunAnalyze(context.Background(), &out, path); err != nil {
		t.Fatalf("RunAnalyze() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"path": `, `"format": "pe"`, `"source": "core.ingest"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunAnalyze() output missing %q in %q", want, got)
		}
	}
}

func TestRunAnalyzeIncludesBuildInfoAndFunctions(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunAnalyze(context.Background(), &out, path); err != nil {
		t.Fatalf("RunAnalyze() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"build_info": {`,
		`"path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "symbol_backed"`,
		`"elf_pclntab_header_magic": "f1ffffff"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"elf_pclntab_header_quantum": 1`,
		`"elf_pclntab_header_pointer_size": 8`,
		`"elf_pclntab_function_count_hint": `,
		`"elf_pclntab_file_count_hint": `,
		`"elf_pclntab_funcnametab_offset_hint": `,
		`"elf_pclntab_functab_offset_hint": `,
		`"elf_functab_last_pc_offset_hint": `,
		`"elf_functab_pc_offsets_monotonic": true`,
		`"elf_functab_first_pc_addr_hint": `,
		`"elf_functab_last_pc_addr_hint": `,
		`"elf_functab_pc_addr_hints_within_text": true`,
		`"elf_functab_pc_addr_sample": [`,
		`"elf_functab_pc_addr_sample_all_within_text": true`,
		`"firstmoduledata_addr": `,
		`"typelink_count": `,
		`"itablink_count": `,
		`"moduledata_pcheader_addr": `,
		`"moduledata_pcheader_matches_gopclntab": true`,
		`"moduledata_funcnametab_slice_word_index": 1`,
		`"moduledata_funcnametab_addr": `,
		`"moduledata_funcnametab_len": `,
		`"moduledata_funcnametab_cap": `,
		`"moduledata_funcnametab_within_gopclntab": true`,
		`"moduledata_cutab_slice_word_index": 4`,
		`"moduledata_cutab_addr": `,
		`"moduledata_cutab_len": `,
		`"moduledata_cutab_cap": `,
		`"moduledata_cutab_within_gopclntab": true`,
		`"moduledata_filetab_slice_word_index": 7`,
		`"moduledata_filetab_addr": `,
		`"moduledata_filetab_len": `,
		`"moduledata_filetab_cap": `,
		`"moduledata_filetab_within_gopclntab": true`,
		`"moduledata_pctab_slice_word_index": 10`,
		`"moduledata_pctab_addr": `,
		`"moduledata_pctab_len": `,
		`"moduledata_pctab_cap": `,
		`"moduledata_pctab_within_gopclntab": true`,
		`"moduledata_pclntable_slice_word_index": 13`,
		`"moduledata_pclntable_addr": `,
		`"moduledata_pclntable_len": `,
		`"moduledata_pclntable_cap": `,
		`"moduledata_pclntable_within_gopclntab": true`,
		`"typelink_sample": [`,
		`"typelink_min_offset": `,
		`"typelink_max_offset": `,
		`"typelink_non_negative_count": `,
		`"firstmoduledata_in_go_module": true`,
		`"go_module_word_sample": [`,
		`"moduledata_typelink_slice_word_index": `,
		`"moduledata_typelink_len": `,
		`"moduledata_typelink_cap": `,
		`"moduledata_itablink_slice_word_index": `,
		`"moduledata_itablink_len": `,
		`"moduledata_itablink_cap": `,
		`"moduledata_memory_ranges_word_index": `,
		`"moduledata_noptrdata_addr": `,
		`"moduledata_data_addr": `,
		`"moduledata_bss_addr": `,
		`"moduledata_noptrbss_addr": `,
		`"moduledata_rodata_word_index": `,
		`"moduledata_rodata_addr": `,
		`"moduledata_rodata_end": `,
		`"moduledata_text_word_index": `,
		`"moduledata_text_addr": `,
		`"moduledata_text_end_inclusive": `,
		`"typelink_resolved_base_addr": `,
		`"typelink_resolved_sample": [`,
		`"typelink_resolved_within_rodata_count": `,
		`"typelink_all_resolved_within_rodata": true`,
		`"moduledata_types_addr": `,
		`"moduledata_etypes_addr": `,
		`"moduledata_types_range_word_index": `,
		`"typelink_resolved_within_types_count": `,
		`"typelink_all_resolved_within_types": true`,
		`"functions": [`,
		`"name": "main.main"`,
		`"packages": [`,
		`"name": "main"`,
		`"types": [`,
		`"name": "main.fixtureCounter"`,
		`"strings": [`,
		`"value": "goreveal fixture"`,
		`"source_tree": {`,
		`"root": "example.com/gorevealfixture"`,
		`"peeling": {`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
		`"packages": [`,
		`"primary_classification": "user"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunAnalyze() output missing %q in %q", want, got)
		}
	}
}

func TestRunAnalyzeStrippedFixturePreservesBoundedRuntimeAndPackageContract(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunAnalyze(context.Background(), &out, path); err != nil {
		t.Fatalf("RunAnalyze() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"build_info": {`,
		`"path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "go_module_fallback"`,
		`"elf_pclntab_header_magic": "f1ffffff"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"elf_pclntab_function_count_hint": `,
		`"elf_functab_last_pc_offset_hint": `,
		`"elf_functab_pc_offsets_monotonic": true`,
		`"elf_functab_first_pc_addr_hint": `,
		`"elf_functab_last_pc_addr_hint": `,
		`"elf_functab_pc_addr_hints_within_text": true`,
		`"elf_functab_pc_addr_sample": [`,
		`"elf_functab_pc_addr_sample_all_within_text": true`,
		`"firstmoduledata_addr": `,
		`"firstmoduledata_from_go_module_fallback": true`,
		`"go_module_addr": `,
		`"firstmoduledata_in_go_module": true`,
		`"moduledata_pcheader_matches_gopclntab": true`,
		`"packages": [`,
		`"name": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"module_local": true`,
		`"source_tree": {`,
		`"root": "example.com/gorevealfixture"`,
		`"external_packages": [`,
		`"import_path": "runtime"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunAnalyze() stripped output missing %q in %q", want, got)
		}
	}
	if !strings.Contains(got, `"types": []`) {
		t.Fatalf("RunAnalyze() stripped output missing empty types array in %q", got)
	}
}

func TestRunInspectRuntimeStrippedFixtureShowsFallbackSource(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectRuntime(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectRuntime() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"trust_summary": "go_module_fallback"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"elf_pclntab_function_count_hint": `,
		`"elf_functab_last_pc_offset_hint": `,
		`"elf_functab_first_pc_addr_hint": `,
		`"elf_functab_last_pc_addr_hint": `,
		`"elf_functab_pc_addr_hints_within_text": true`,
		`"elf_functab_pc_addr_sample": [`,
		`"elf_functab_pc_addr_sample_all_within_text": true`,
		`"firstmoduledata_addr": `,
		`"firstmoduledata_from_go_module_fallback": true`,
		`"go_module_addr": `,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectRuntime() stripped output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectTypes(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectTypes(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectTypes() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"name": "main.fixtureCounter"`, `"kind": "`, `"package": "main"`, `"import_path": "example.com/gorevealfixture"`, `"source_file_count": 1`, `"module_local": true`, `"user_meaningful": true`} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectTypes() output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectPackagesStrippedFixture(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectPackages(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectPackages() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"name": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"module_local": true`,
		`"has_source_evidence": true`,
		`"source_evidence_kind": "line_table_files"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectPackages() stripped output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectTypesStrippedFixtureReturnsEmptyArray(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectTypes(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectTypes() error = %v", err)
	}

	got := strings.TrimSpace(out.String())
	if got != "[]" {
		t.Fatalf("RunInspectTypes() stripped output = %q, want %q", got, "[]")
	}
}

func TestRunInspectStrings(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectStrings(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectStrings() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"value": "goreveal fixture"`, `"region": "`, `"addr": `} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectStrings() output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectRuntime(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectRuntime(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectRuntime() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"trust_summary": "symbol_backed"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"firstmoduledata_addr": `,
		`"typelink_count": `,
		`"moduledata_pclntable_within_gopclntab": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectRuntime() output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectPeeling(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectPeeling(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectPeeling() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"functions": [`,
		`"name": "main.main"`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
		`"packages": [`,
		`"name": "main"`,
		`"primary_classification": "user"`,
		`"user_function_count": `,
		`"name": "runtime.newobject"`,
		`"classification": "runtime"`,
		`"classification_evidence": "runtime_import_path"`,
		`"name": "runtime"`,
		`"primary_classification": "runtime"`,
		`"runtime_function_count": `,
		`"name": "fmt.Fprintln"`,
		`"classification": "stdlib"`,
		`"classification_evidence": "stdlib_import_path"`,
		`"name": "fmt"`,
		`"primary_classification": "stdlib"`,
		`"stdlib_function_count": `,
		`"source": "engine.peeling"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectPeeling() output missing %q in %q", want, got)
		}
	}
}

func TestRunPeelReturnsOnlyUserOwnedProjection(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunPeel(context.Background(), &out, path); err != nil {
		t.Fatalf("RunPeel() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"functions": [`,
		`"name": "main.main"`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
		`"packages": [`,
		`"name": "main"`,
		`"primary_classification": "user"`,
		`"source": "engine.peeling.user_only"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunPeel() output missing %q in %q", want, got)
		}
	}
	for _, unwanted := range []string{
		`"name": "runtime.newobject"`,
		`"name": "fmt.Fprintln"`,
		`"name": "runtime"`,
		`"name": "fmt"`,
	} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("RunPeel() output unexpectedly contains %q in %q", unwanted, got)
		}
	}
}

func TestRunInspectRuntimeMalformedPEFailsExplicitly(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.bin")
	if err := os.WriteFile(path, []byte{'M', 'Z', 0x90, 0x00}, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var out strings.Builder
	err := RunInspectRuntime(context.Background(), &out, path)
	if err == nil {
		t.Fatal("RunInspectRuntime() error = nil, want explicit failed stage")
	}
	if !strings.Contains(err.Error(), `status="failed"`) || !strings.Contains(err.Error(), `code="stage_failed"`) {
		t.Fatalf("RunInspectRuntime() error = %q, want explicit failed stage", err)
	}
}

func TestRunInspectRuntimePEFixtureReturnsSectionHeuristic(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	var out strings.Builder
	if err := RunInspectRuntime(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectRuntime() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"trust_summary": "section_heuristic"`,
		`"pe_text_section_addr": `,
		`"pe_rdata_section_addr": `,
		`"pe_pclntab_magic_section": ".rdata"`,
		`"pe_pclntab_magic_addr": `,
		`"pe_pclntab_magic_count": `,
		`"pe_pclntab_header_section": ".rdata"`,
		`"pe_pclntab_header_magic": "f1ffffff"`,
		`"pe_pclntab_header_quantum": 1`,
		`"pe_pclntab_header_pointer_size": 8`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectRuntime() PE output missing %q in %q", want, got)
		}
	}
}

func TestRunAnalyzePEFixtureIncludesBoundedRuntimeSurface(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	var out strings.Builder
	if err := RunAnalyze(context.Background(), &out, path); err != nil {
		t.Fatalf("RunAnalyze() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"format": "pe"`,
		`"build_info": {`,
		`"path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "section_heuristic"`,
		`"pe_text_section_addr": `,
		`"pe_pclntab_header_magic": "f1ffffff"`,
		`"functions": [`,
		`"name": "main.main"`,
		`"packages": [`,
		`"module_local": true`,
		`"peeling": {`,
		`"classification": "user"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunAnalyze() PE output missing %q in %q", want, got)
		}
	}
}

func TestRunAnalyzeMachOFixtureIncludesBoundedFunctionFoothold(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-macho-buildinfo-darwin-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunAnalyze(context.Background(), &out, path); err != nil {
		t.Fatalf("RunAnalyze() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"format": "macho"`,
		`"build_info": {`,
		`"path": "example.com/gorevealfixture"`,
		`"functions": [`,
		`"name": "main.main"`,
		`"packages": [`,
		`"name": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"module_local": true`,
		`"peeling": {`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunAnalyze() Mach-O output missing %q in %q", want, got)
		}
	}
	if strings.Contains(got, `"runtime": {`) {
		t.Fatalf("RunAnalyze() Mach-O output unexpectedly includes runtime in %q", got)
	}
}

func TestRunSourceTree(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunSourceTree(context.Background(), &out, path); err != nil {
		t.Fatalf("RunSourceTree() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"root": "example.com/gorevealfixture"`, `"files": [`, `"main.go"`, `"function_count": 3`, `"has_file_evidence": true`, `"external_packages": [`, `"import_path": "runtime"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() output missing %q in %q", want, got)
		}
	}
	for _, want := range []string{`"source_evidence_kind": "dwarf_paths"`, `"source_evidence_summary": {`, `"tree_kind": "dwarf_paths"`, `"dwarf_path_package_count": `, `"dwarf_path_file_count": `} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() output missing %q in %q", want, got)
		}
	}
}

func TestRunSourceTreeStrippedFixtureReturnsBoundedFallback(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunSourceTree(context.Background(), &out, path); err != nil {
		t.Fatalf("RunSourceTree() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"root": "example.com/gorevealfixture"`,
		`"source_evidence_kind": "line_table_files"`,
		`"source_evidence_summary": {`,
		`"tree_kind": "line_table_files"`,
		`"line_table_package_count": `,
		`"line_table_file_count": `,
		`"pathless_file_evidence": true`,
		`"files": [`,
		`"main.go"`,
		`"packages": [`,
		`"name": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"function_count": 3`,
		`"has_file_evidence": true`,
		`"external_packages": [`,
		`"import_path": "runtime"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() stripped output missing %q in %q", want, got)
		}
	}
}

func TestRunSourceTreePEFixtureReturnsLineTableFallback(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	var out strings.Builder
	if err := RunSourceTree(context.Background(), &out, path); err != nil {
		t.Fatalf("RunSourceTree() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"root": "example.com/gorevealfixture"`,
		`"source_evidence_kind": "line_table_files"`,
		`"source_evidence_summary": {`,
		`"tree_kind": "line_table_files"`,
		`"line_table_package_count": `,
		`"line_table_file_count": `,
		`"pathless_file_evidence": true`,
		`"files": [`,
		`"main.go"`,
		`"packages": [`,
		`"name": "main"`,
		`"has_file_evidence": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() PE output missing %q in %q", want, got)
		}
	}
}

func TestRunSourceTreeMachOFixtureReturnsLineTableFallback(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-macho-buildinfo-darwin-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunSourceTree(context.Background(), &out, path); err != nil {
		t.Fatalf("RunSourceTree() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"root": "example.com/gorevealfixture"`,
		`"source_evidence_kind": "line_table_files"`,
		`"source_evidence_summary": {`,
		`"tree_kind": "line_table_files"`,
		`"line_table_package_count": `,
		`"line_table_file_count": `,
		`"pathless_file_evidence": true`,
		`"files": [`,
		`"main.go"`,
		`"packages": [`,
		`"name": "main"`,
		`"has_file_evidence": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() Mach-O output missing %q in %q", want, got)
		}
	}
}

func TestRunDeobfuscate(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunDeobfuscate(context.Background(), &out, path); err != nil {
		t.Fatalf("RunDeobfuscate() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"passes": [`, `"synthetic-function-names"`, `"string-segments"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDeobfuscate() output missing %q in %q", want, got)
		}
	}
}

func TestRunExportSQLite(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	dbPath := filepath.Join(t.TempDir(), "analysis.db")

	if err := RunExportSQLite(context.Background(), dbPath, path); err != nil {
		t.Fatalf("RunExportSQLite() error = %v", err)
	}

	store, err := storesqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	got, err := store.LoadAnalysis(context.Background(), 1)
	if err != nil {
		t.Fatalf("LoadAnalysis() error = %v", err)
	}
	if got.Input.Path != path {
		t.Fatalf("LoadAnalysis() path = %q, want %q", got.Input.Path, path)
	}
	if got.BuildInfo == nil || got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("LoadAnalysis() build info = %#v", got.BuildInfo)
	}
}

func TestRunExportIDA(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunExportIDA(context.Background(), &out, path); err != nil {
		t.Fatalf("RunExportIDA() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"contract": "goreveal.export.ida/v1"`,
		`"runtime": {`,
		`"trust_summary": "symbol_backed"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"peeling": {`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
		`"primary_classification": "user"`,
		`"functions": [`,
		`"refined_name": "main.main"`,
		`"types": [`,
		`"strings": [`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunExportIDA() output missing %q in %q", want, got)
		}
	}
}

func TestRunExportGhidra(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunExportGhidra(context.Background(), &out, path); err != nil {
		t.Fatalf("RunExportGhidra() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"contract": "goreveal.export.ghidra/v1"`,
		`"program": {`,
		`"module_path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "symbol_backed"`,
		`"elf_pclntab_header_magic_kind": "known"`,
		`"peeling": {`,
		`"classification": "user"`,
		`"classification_evidence": "module_local"`,
		`"primary_classification": "user"`,
		`"symbols": [`,
		`"refined_name": "main.main"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunExportGhidra() output missing %q in %q", want, got)
		}
	}
}

func TestRunExportIDAPEFixtureStaysThin(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	var out strings.Builder
	if err := RunExportIDA(context.Background(), &out, path); err != nil {
		t.Fatalf("RunExportIDA() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"contract": "goreveal.export.ida/v1"`,
		`"format": "pe"`,
		`"build_info": {`,
		`"path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "section_heuristic"`,
		`"pe_pclntab_magic_section": ".rdata"`,
		`"pe_pclntab_header_magic": "f1ffffff"`,
		`"functions": [`,
		`"name": "main.main"`,
		`"peeling": {`,
		`"classification": "user"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunExportIDA() PE output missing %q in %q", want, got)
		}
	}
}

func TestRunExportGhidraPEFixtureStaysThin(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	var out strings.Builder
	if err := RunExportGhidra(context.Background(), &out, path); err != nil {
		t.Fatalf("RunExportGhidra() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"contract": "goreveal.export.ghidra/v1"`,
		`"program": {`,
		`"format": "pe"`,
		`"module_path": "example.com/gorevealfixture"`,
		`"runtime": {`,
		`"trust_summary": "section_heuristic"`,
		`"pe_pclntab_magic_section": ".rdata"`,
		`"pe_pclntab_header_magic": "f1ffffff"`,
		`"symbols": [`,
		`"refined_name": "main.main"`,
		`"peeling": {`,
		`"classification": "user"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunExportGhidra() PE output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectFunctions(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectFunctions(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectFunctions() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"name": "main.main"`, `"entry": `} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectFunctions() output missing %q in %q", want, got)
		}
	}
	for _, want := range []string{
		`"package": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"module_local": true`,
		`"source_file": "main.go"`,
		`"source_line": 34`,
		`"autogenerated": true`,
		`"package": "runtime"`,
		`"import_path": "runtime"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectFunctions() output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectFunctionsStrippedFixturePreservesBoundedMetadata(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectFunctions(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectFunctions() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"name": "main.main"`,
		`"package": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"module_local": true`,
		`"source_file": "main.go"`,
		`"source_line": 34`,
		`"package": "runtime"`,
		`"import_path": "runtime"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectFunctions() stripped output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectPackages(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	var out strings.Builder
	if err := RunInspectPackages(context.Background(), &out, path); err != nil {
		t.Fatalf("RunInspectPackages() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{`"name": "main"`, `"function_count": `, `"import_path": "example.com/gorevealfixture"`, `"source_file_count": `, `"module_local": true`, `"has_source_evidence": true`, `"source_evidence_kind": "dwarf_paths"`} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectPackages() output missing %q in %q", want, got)
		}
	}
}

func TestRunDiffSQLite(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	dbPath := filepath.Join(t.TempDir(), "analysis.db")

	if err := RunExportSQLite(context.Background(), dbPath, path); err != nil {
		t.Fatalf("RunExportSQLite() first error = %v", err)
	}
	if err := RunExportSQLite(context.Background(), dbPath, path); err != nil {
		t.Fatalf("RunExportSQLite() second error = %v", err)
	}

	var out strings.Builder
	if err := RunDiffSQLite(context.Background(), &out, dbPath, 1, 2); err != nil {
		t.Fatalf("RunDiffSQLite() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"left_id": 1`,
		`"right_id": 2`,
		`"build_info_changed": false`,
		`"left_counts": {`,
		`"right_counts": {`,
		`"functions": `,
		`"matched_functions": [`,
		`"transfer_candidates": [`,
		`"accepted_transfers": [`,
		`"transfer_packages": [`,
		`"left_name": "main.main"`,
		`"score": 100`,
		`"reason": "exact_name"`,
		`"match_reason": "exact_name"`,
		`"disposition": "ready"`,
		`"accepted_by": "exact_name"`,
		`"projected_package": "main"`,
		`"candidate_count": `,
		`"highest_candidate_reason": "exact_name"`,
		`"projected_classification": "user"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffSQLite() output missing %q in %q", want, got)
		}
	}
}

func TestRunDiffReviewSQLite(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "analysis.db")
	store, err := storesqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	left := schema.Analysis{
		Input: schema.Input{Path: "left.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.main", Package: "main", ImportPath: "example.com/sample", SourceFile: "main.go", SourceLine: 10, ModuleLocal: true},
			{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 20, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.main",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "main.go",
					SourceLine:             10,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceModuleLocal,
				},
				{
					Name:                   "main.service",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             20,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}
	right := schema.Analysis{
		Input: schema.Input{Path: "right.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.main", Package: "main", ImportPath: "example.com/sample", SourceFile: "main.go", SourceLine: 10, ModuleLocal: true},
			{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 30, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.main",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "main.go",
					SourceLine:             10,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceModuleLocal,
				},
				{
					Name:                   "main.serviceV2",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             30,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}

	if _, err := store.SaveAnalysis(context.Background(), left); err != nil {
		t.Fatalf("SaveAnalysis(left) error = %v", err)
	}
	if _, err := store.SaveAnalysis(context.Background(), right); err != nil {
		t.Fatalf("SaveAnalysis(right) error = %v", err)
	}

	var out strings.Builder
	if err := RunDiffReviewSQLite(context.Background(), &out, dbPath, 1, 2); err != nil {
		t.Fatalf("RunDiffReviewSQLite() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"left_id": 1`,
		`"right_id": 2`,
		`"left_input": {`,
		`"path": "left.bin"`,
		`"right_input": {`,
		`"path": "right.bin"`,
		`"transfer_review": {`,
		`"review_count": 1`,
		`"review_package_count": 1`,
		`"transfer_review_packages": [`,
		`"name": "main"`,
		`"highest_match_reason": "source_file"`,
		`"transfer_review_focus": {`,
		`"action": "review_package"`,
		`"package": "main"`,
		`"import_path": "example.com/sample"`,
		`"item_count": 1`,
		`"handoff": {`,
		`"handoff_contract": "goreveal.review_handoff/v1"`,
		`"artifact_role": "review_handoff"`,
		`"mode": "host_platform_review"`,
		`"recommended_path": "export_then_import"`,
		`"recommended_targets": [`,
		`"ida"`,
		`"ghidra"`,
		`"artifacts": [`,
		`"id": "review_handoff"`,
		`"contract": "goreveal.review_handoff/v1"`,
		`"format": "json"`,
		`"id": "ida_export"`,
		`"contract": "goreveal.export.ida/v1"`,
		`"format": "ida"`,
		`"id": "ghidra_export"`,
		`"contract": "goreveal.export.ghidra/v1"`,
		`"format": "ghidra"`,
		`"target_profiles": [`,
		`"target": "ida"`,
		`"recommended_mcp_server": "ida-pro-mcp"`,
		`"export_format": "ida"`,
		`"export_contract": "goreveal.export.ida/v1"`,
		`"artifact_role": "go_metadata_export"`,
		`"binding_mode": "mcp_server"`,
		`"host_entrypoint": "ida-pro-mcp.import_export_payload"`,
		`"import_mode": "mcp_or_workspace_import"`,
		`"preferred_transport": "mcp"`,
		`"workspace_phase": "import_then_annotate"`,
		`"workspace_action": "apply_go_specific_annotations"`,
		`"expected_host_result": "annotated_workspace_review_ready"`,
		`"completion_signal": "names_comments_and_runtime_context_applied"`,
		`"host_actions": [`,
		`"import_export_payload"`,
		`"apply_names_and_comments"`,
		`"review_runtime_and_package_context"`,
		`"required_artifacts": [`,
		`"review_handoff"`,
		`"ida_export"`,
		`"target": "ghidra"`,
		`"export_format": "ghidra"`,
		`"export_contract": "goreveal.export.ghidra/v1"`,
		`"artifact_role": "go_metadata_export"`,
		`"binding_mode": "workspace_loader"`,
		`"host_entrypoint": "ghidra.workspace_import"`,
		`"import_mode": "workspace_import"`,
		`"preferred_transport": "file_or_workspace_import"`,
		`"workspace_phase": "import_then_annotate"`,
		`"expected_host_result": "annotated_workspace_review_ready"`,
		`"completion_signal": "names_comments_and_runtime_context_applied"`,
		`"required_artifacts": [`,
		`"review_handoff"`,
		`"ghidra_export"`,
		`"left_name": "main.service"`,
		`"right_name": "main.serviceV2"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffReviewSQLite() output missing %q in %q", want, got)
		}
	}
}

func TestRunDiffHandoffSQLite(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "analysis.db")
	store, err := storesqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	left := schema.Analysis{
		Input: schema.Input{Path: "left.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 20, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.service",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             20,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}
	right := schema.Analysis{
		Input: schema.Input{Path: "right.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 30, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.serviceV2",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             30,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}

	if _, err := store.SaveAnalysis(context.Background(), left); err != nil {
		t.Fatalf("SaveAnalysis(left) error = %v", err)
	}
	if _, err := store.SaveAnalysis(context.Background(), right); err != nil {
		t.Fatalf("SaveAnalysis(right) error = %v", err)
	}

	var out strings.Builder
	if err := RunDiffHandoffSQLite(context.Background(), &out, dbPath, 1, 2); err != nil {
		t.Fatalf("RunDiffHandoffSQLite() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"left_id": 1`,
		`"right_id": 2`,
		`"left_input": {`,
		`"path": "left.bin"`,
		`"right_input": {`,
		`"path": "right.bin"`,
		`"transfer_review_focus": {`,
		`"action": "review_package"`,
		`"handoff": {`,
		`"handoff_contract": "goreveal.review_handoff/v1"`,
		`"artifact_role": "review_handoff"`,
		`"mode": "host_platform_review"`,
		`"recommended_path": "export_then_import"`,
		`"recommended_targets": [`,
		`"ida"`,
		`"ghidra"`,
		`"recommended_mcp_servers": [`,
		`"ida-pro-mcp"`,
		`"recommended_exports": [`,
		`"artifacts": [`,
		`"id": "review_handoff"`,
		`"contract": "goreveal.review_handoff/v1"`,
		`"format": "json"`,
		`"id": "ida_export"`,
		`"contract": "goreveal.export.ida/v1"`,
		`"format": "ida"`,
		`"id": "ghidra_export"`,
		`"contract": "goreveal.export.ghidra/v1"`,
		`"format": "ghidra"`,
		`"review_command": "goreveal diff review sqlite `,
		`"export_commands": [`,
		`"goreveal export ida right.bin"`,
		`"goreveal export ghidra right.bin"`,
		`"target_profiles": [`,
		`"target": "ida"`,
		`"recommended_mcp_server": "ida-pro-mcp"`,
		`"export_format": "ida"`,
		`"export_contract": "goreveal.export.ida/v1"`,
		`"artifact_role": "go_metadata_export"`,
		`"binding_mode": "mcp_server"`,
		`"host_entrypoint": "ida-pro-mcp.import_export_payload"`,
		`"import_mode": "mcp_or_workspace_import"`,
		`"preferred_transport": "mcp"`,
		`"workspace_phase": "import_then_annotate"`,
		`"workspace_action": "apply_go_specific_annotations"`,
		`"expected_host_result": "annotated_workspace_review_ready"`,
		`"completion_signal": "names_comments_and_runtime_context_applied"`,
		`"host_actions": [`,
		`"import_export_payload"`,
		`"apply_names_and_comments"`,
		`"review_runtime_and_package_context"`,
		`"required_artifacts": [`,
		`"review_handoff"`,
		`"ida_export"`,
		`"target": "ghidra"`,
		`"export_format": "ghidra"`,
		`"export_contract": "goreveal.export.ghidra/v1"`,
		`"artifact_role": "go_metadata_export"`,
		`"binding_mode": "workspace_loader"`,
		`"host_entrypoint": "ghidra.workspace_import"`,
		`"import_mode": "workspace_import"`,
		`"preferred_transport": "file_or_workspace_import"`,
		`"workspace_phase": "import_then_annotate"`,
		`"expected_host_result": "annotated_workspace_review_ready"`,
		`"completion_signal": "names_comments_and_runtime_context_applied"`,
		`"required_artifacts": [`,
		`"review_handoff"`,
		`"ghidra_export"`,
		`"operator_steps": [`,
		`"handoff runtime/package review for main from left.bin to host platform MCP or workspace import"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffHandoffSQLite() output missing %q in %q", want, got)
		}
	}
}

func TestRunDiffNextSQLite(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "analysis.db")
	store, err := storesqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	left := schema.Analysis{
		Input: schema.Input{Path: "left.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 20, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.service",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             20,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}
	right := schema.Analysis{
		Input: schema.Input{Path: "right.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 30, ModuleLocal: true},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{
					Name:                   "main.serviceV2",
					Package:                "main",
					ImportPath:             "example.com/sample",
					SourceFile:             "service.go",
					SourceLine:             30,
					ModuleLocal:            true,
					Classification:         schema.PeelingClassUser,
					ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath,
				},
			},
		},
	}

	if _, err := store.SaveAnalysis(context.Background(), left); err != nil {
		t.Fatalf("SaveAnalysis(left) error = %v", err)
	}
	if _, err := store.SaveAnalysis(context.Background(), right); err != nil {
		t.Fatalf("SaveAnalysis(right) error = %v", err)
	}

	var out strings.Builder
	if err := RunDiffNextSQLite(context.Background(), &out, dbPath, 1, 2); err != nil {
		t.Fatalf("RunDiffNextSQLite() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"left_id": 1`,
		`"right_id": 2`,
		`"transfer_review_plan": [`,
		`"action": "review_package"`,
		`"package": "main"`,
		`"import_path": "example.com/sample"`,
		`"highest_match_reason": "source_file"`,
		`"item_count": 1`,
		`"items": [`,
		`"left_name": "main.service"`,
		`"right_name": "main.serviceV2"`,
		`"transfer_review_focus": {`,
		`"review_progress": {`,
		`"current_step": 1`,
		`"total_steps": 1`,
		`"current_package": "main"`,
		`"current_import_path": "example.com/sample"`,
		`"current_item_count": 1`,
		`"recommended_actions": [`,
		`"review_checklist": [`,
		`"review all 1 pending transfer items for main"`,
		`"confirm the strongest pending match reason for main remains source_file"`,
		`"emit or reuse the handoff bundle for main before host-platform review"`,
		`"review_snapshot": {`,
		`"current_package": "main"`,
		`"current_import_path": "example.com/sample"`,
		`"current_item_count": 1`,
		`"current_highest_match_score": 90`,
		`"current_highest_match_reason": "source_file"`,
		`"recommended_action_count": 4`,
		`"kind": "review_bundle"`,
		`"command": "goreveal diff review sqlite `,
		`"description": "review the current package bundle against the focused transfer items"`,
		`"kind": "handoff_bundle"`,
		`"command": "goreveal diff handoff sqlite `,
		`"description": "emit the workstation handoff artifact for the current package bundle"`,
		`"kind": "export_target"`,
		`"target": "ida"`,
		`"command": "goreveal export ida right.bin"`,
		`"target": "ghidra"`,
		`"command": "goreveal export ghidra right.bin"`,
		`"review_command": "goreveal diff review sqlite `,
		`"handoff_command": "goreveal diff handoff sqlite `,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffNextSQLite() output missing %q in %q", want, got)
		}
	}
}

func TestRunDiffNextSQLiteIncludesUpcomingPackageProgress(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "analysis.db")
	store, err := storesqlite.Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	left := schema.Analysis{
		Input:     schema.Input{Path: "left.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample"},
		Functions: []schema.Function{
			{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 40},
			{Name: "main.alpha", Package: "main", ImportPath: "example.com/sample", SourceFile: "alpha.go", SourceLine: 10},
			{Name: "example.com/sample/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/internal/app", SourceFile: "handler.go", SourceLine: 7},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.service", Package: "main", ImportPath: "example.com/sample", SourceFile: "service.go", SourceLine: 40, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "main.alpha", Package: "main", ImportPath: "example.com/sample", SourceFile: "alpha.go", SourceLine: 10, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/internal/app.Handler", Package: "internal/app", ImportPath: "example.com/sample/internal/app", SourceFile: "handler.go", SourceLine: 7, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}
	right := schema.Analysis{
		Input:     schema.Input{Path: "right.bin", Format: "elf"},
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample/v2"},
		Functions: []schema.Function{
			{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "service.go", SourceLine: 48},
			{Name: "main.alphaV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "alpha.go", SourceLine: 15},
			{Name: "example.com/sample/v2/internal/app.HandlerV2", Package: "internal/app", ImportPath: "example.com/sample/v2/internal/app", SourceFile: "handler.go", SourceLine: 21},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{
				{Name: "main.serviceV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "service.go", SourceLine: 48, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "main.alphaV2", Package: "main", ImportPath: "example.com/sample/v2", SourceFile: "alpha.go", SourceLine: 15, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceBuildInfoPath},
				{Name: "example.com/sample/v2/internal/app.HandlerV2", Package: "internal/app", ImportPath: "example.com/sample/v2/internal/app", SourceFile: "handler.go", SourceLine: 21, Classification: schema.PeelingClassUser, ClassificationEvidence: schema.PeelingEvidenceModuleLocal},
			},
		},
	}

	if _, err := store.SaveAnalysis(context.Background(), left); err != nil {
		t.Fatalf("SaveAnalysis(left) error = %v", err)
	}
	if _, err := store.SaveAnalysis(context.Background(), right); err != nil {
		t.Fatalf("SaveAnalysis(right) error = %v", err)
	}

	var out strings.Builder
	if err := RunDiffNextSQLite(context.Background(), &out, dbPath, 1, 2); err != nil {
		t.Fatalf("RunDiffNextSQLite() error = %v", err)
	}

	got := out.String()
	for _, want := range []string{
		`"up_next": {`,
		`"package": "internal/app"`,
		`"import_path": "example.com/sample/internal/app"`,
		`"review_count": 1`,
		`"item_count": 1`,
		`"sample_left_name": "example.com/sample/internal/app.Handler"`,
		`"sample_right_name": "example.com/sample/v2/internal/app.HandlerV2"`,
		`"upcoming_packages": [`,
		`"package": "internal/app"`,
		`"import_path": "example.com/sample/internal/app"`,
		`"review_count": 1`,
		`"item_count": 1`,
		`"highest_match_score": 90`,
		`"highest_match_reason": "source_file"`,
		`"sample_left_name": "example.com/sample/internal/app.Handler"`,
		`"sample_right_name": "example.com/sample/v2/internal/app.HandlerV2"`,
		`"review_progress": {`,
		`"current_step": 1`,
		`"total_steps": 2`,
		`"current_package": "main"`,
		`"current_item_count": 2`,
		`"remaining_package_count": 1`,
		`"remaining_review_item_count": 1`,
		`"next_package": "internal/app"`,
		`"next_import_path": "example.com/sample/internal/app"`,
		`"next_item_count": 1`,
		`"review_checklist": [`,
		`"review all 2 pending transfer items for main"`,
		`"confirm the strongest pending match reason for main remains source_file"`,
		`"emit or reuse the handoff bundle for main before host-platform review"`,
		`"after main, continue with internal/app"`,
		`"review_snapshot": {`,
		`"current_package": "main"`,
		`"current_import_path": "example.com/sample"`,
		`"current_item_count": 2`,
		`"current_highest_match_score": 90`,
		`"current_highest_match_reason": "source_file"`,
		`"next_package": "internal/app"`,
		`"next_import_path": "example.com/sample/internal/app"`,
		`"remaining_review_item_count": 1`,
		`"recommended_action_count": 4`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffNextSQLite() output missing %q in %q", want, got)
		}
	}
}
