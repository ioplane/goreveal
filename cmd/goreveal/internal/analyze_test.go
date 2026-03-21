package internalcmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	storesqlite "github.com/dantte-lp/goreveal/storage/sqlite"
)

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
		`"has_source_evidence": false`,
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
		`"firstmoduledata_addr": `,
		`"typelink_count": `,
		`"moduledata_pclntable_within_gopclntab": true`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunInspectRuntime() output missing %q in %q", want, got)
		}
	}
}

func TestRunInspectRuntimePEUnavailable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.bin")
	if err := os.WriteFile(path, []byte{'M', 'Z', 0x90, 0x00}, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var out strings.Builder
	err := RunInspectRuntime(context.Background(), &out, path)
	if err == nil {
		t.Fatal("RunInspectRuntime() error = nil, want unavailable error")
	}
	if !strings.Contains(err.Error(), "unavailable") {
		t.Fatalf("RunInspectRuntime() error = %q, want unavailable", err)
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
		`"files": []`,
		`"packages": [`,
		`"name": "main"`,
		`"import_path": "example.com/gorevealfixture"`,
		`"function_count": 3`,
		`"has_file_evidence": false`,
		`"external_packages": [`,
		`"import_path": "runtime"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunSourceTree() stripped output missing %q in %q", want, got)
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
		`"symbols": [`,
		`"refined_name": "main.main"`,
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunExportGhidra() output missing %q in %q", want, got)
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
	for _, want := range []string{`"name": "main"`, `"function_count": `, `"import_path": "example.com/gorevealfixture"`, `"source_file_count": `, `"module_local": true`, `"has_source_evidence": true`} {
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
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("RunDiffSQLite() output missing %q in %q", want, got)
		}
	}
}
