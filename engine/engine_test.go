package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestAnalyzeFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.bin")
	if err := os.WriteFile(path, []byte{0x7f, 'E', 'L', 'F', 0x02}, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	analyzer := New()
	got, err := analyzer.AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	if got.Input.Path != path {
		t.Fatalf("AnalyzeFile() path = %q, want %q", got.Input.Path, path)
	}
	if got.Input.Format != "elf" {
		t.Fatalf("AnalyzeFile() format = %q, want %q", got.Input.Format, "elf")
	}
	if got.Provenance.Source != "core.ingest" {
		t.Fatalf("AnalyzeFile() provenance source = %q", got.Provenance.Source)
	}
	if got.Provenance.Confidence != "high" {
		t.Fatalf("AnalyzeFile() provenance confidence = %q", got.Provenance.Confidence)
	}
}

func TestAnalyzeGoFixtureIncludesBuildInfoAndFunctions(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")

	got, err := New().AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	if got.BuildInfo == nil {
		t.Fatal("AnalyzeFile() BuildInfo = nil")
	}
	if got.Runtime == nil {
		t.Fatal("AnalyzeFile() Runtime = nil")
	}
	if got.Runtime.FirstModuleDataAddr == 0 || got.Runtime.GopclntabAddr == 0 || got.Runtime.TypelinkAddr == 0 || got.Runtime.TypelinkCount == 0 || got.Runtime.ItablinkAddr == 0 || got.Runtime.ItablinkCount == 0 {
		t.Fatalf("AnalyzeFile() runtime metadata = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataPCHeaderAddr == 0 || !got.Runtime.ModuledataPCHeaderMatchesGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata pcheader bridge = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataFuncnametabSliceWordIndex != 1 || got.Runtime.ModuledataFuncnametabAddr == 0 || got.Runtime.ModuledataFuncnametabLen == 0 || got.Runtime.ModuledataFuncnametabCap == 0 || !got.Runtime.ModuledataFuncnametabWithinGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata funcnametab bridge = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataCutabSliceWordIndex != 4 || got.Runtime.ModuledataCutabAddr == 0 || got.Runtime.ModuledataCutabLen == 0 || got.Runtime.ModuledataCutabCap == 0 || !got.Runtime.ModuledataCutabWithinGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata cutab bridge = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataFiletabSliceWordIndex != 7 || got.Runtime.ModuledataFiletabAddr == 0 || got.Runtime.ModuledataFiletabLen == 0 || got.Runtime.ModuledataFiletabCap == 0 || !got.Runtime.ModuledataFiletabWithinGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata filetab bridge = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataPctabSliceWordIndex != 10 || got.Runtime.ModuledataPctabAddr == 0 || got.Runtime.ModuledataPctabLen == 0 || got.Runtime.ModuledataPctabCap == 0 || !got.Runtime.ModuledataPctabWithinGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata pctab bridge = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataPclntableSliceWordIndex != 13 || got.Runtime.ModuledataPclntableAddr == 0 || got.Runtime.ModuledataPclntableLen == 0 || got.Runtime.ModuledataPclntableCap == 0 || !got.Runtime.ModuledataPclntableWithinGopclntab {
		t.Fatalf("AnalyzeFile() runtime moduledata pclntable bridge = %#v", got.Runtime)
	}
	if len(got.Runtime.TypelinkSample) == 0 {
		t.Fatalf("AnalyzeFile() runtime typelink sample = %#v", got.Runtime)
	}
	if got.Runtime.TypelinkMinOffset == 0 || got.Runtime.TypelinkMaxOffset == 0 {
		t.Fatalf("AnalyzeFile() runtime typelink min/max = %#v", got.Runtime)
	}
	if got.Runtime.TypelinkNegativeCount+got.Runtime.TypelinkNonNegativeCount != got.Runtime.TypelinkCount {
		t.Fatalf("AnalyzeFile() runtime typelink sign counts = %#v", got.Runtime)
	}
	if !got.Runtime.FirstModuleDataInGoModule || got.Runtime.GoModuleWordSize == 0 || len(got.Runtime.GoModuleWordSample) == 0 {
		t.Fatalf("AnalyzeFile() runtime go.module cross-check = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataTypelinkSliceWordIndex == 0 || got.Runtime.ModuledataTypelinkLen != got.Runtime.TypelinkCount || got.Runtime.ModuledataTypelinkCap != got.Runtime.TypelinkCount {
		t.Fatalf("AnalyzeFile() runtime moduledata typelinks slice = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataItablinkSliceWordIndex == 0 || got.Runtime.ModuledataItablinkLen != got.Runtime.ItablinkCount || got.Runtime.ModuledataItablinkCap != got.Runtime.ItablinkCount {
		t.Fatalf("AnalyzeFile() runtime moduledata itablinks slice = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataMemoryRangesWordIndex == 0 ||
		got.Runtime.ModuledataNoptrdataAddr == 0 || got.Runtime.ModuledataNoptrdataEnd == 0 ||
		got.Runtime.ModuledataDataAddr == 0 || got.Runtime.ModuledataDataEnd == 0 ||
		got.Runtime.ModuledataBssAddr == 0 || got.Runtime.ModuledataBssEnd == 0 ||
		got.Runtime.ModuledataNoptrbssAddr == 0 || got.Runtime.ModuledataNoptrbssEnd == 0 {
		t.Fatalf("AnalyzeFile() runtime moduledata memory ranges = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataRodataWordIndex == 0 || got.Runtime.ModuledataRodataAddr == 0 || got.Runtime.ModuledataRodataEnd == 0 {
		t.Fatalf("AnalyzeFile() runtime moduledata rodata range = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataTextWordIndex == 0 || got.Runtime.ModuledataTextAddr == 0 || got.Runtime.ModuledataTextEndInclusive == 0 {
		t.Fatalf("AnalyzeFile() runtime moduledata text range = %#v", got.Runtime)
	}
	if got.Runtime.TypelinkResolvedBaseAddr == 0 || len(got.Runtime.TypelinkResolvedSample) == 0 || got.Runtime.TypelinkResolvedWithinRodataCount == 0 {
		t.Fatalf("AnalyzeFile() runtime typelink semantic bridge = %#v", got.Runtime)
	}
	if !got.Runtime.TypelinkAllResolvedWithinRodata {
		t.Fatalf("AnalyzeFile() runtime typelink semantic all-within-rodata = %#v", got.Runtime)
	}
	if got.Runtime.ModuledataTypesAddr == 0 || got.Runtime.ModuledataETypesAddr == 0 || got.Runtime.ModuledataTypesRangeWordIndex == 0 || got.Runtime.TypelinkResolvedWithinTypesCount == 0 || !got.Runtime.TypelinkAllResolvedWithinTypes {
		t.Fatalf("AnalyzeFile() runtime typelink types semantics = %#v", got.Runtime)
	}
	if got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("AnalyzeFile() BuildInfo.Path = %q", got.BuildInfo.Path)
	}
	if len(got.Functions) == 0 {
		t.Fatal("AnalyzeFile() returned no functions")
	}
	if len(got.Packages) == 0 {
		t.Fatal("AnalyzeFile() returned no packages")
	}
	mainPkg, ok := findPackage(got.Packages, "main")
	if !ok {
		t.Fatalf("AnalyzeFile() packages missing main: %#v", got.Packages)
	}
	if mainPkg.ImportPath != "example.com/gorevealfixture" {
		t.Fatalf("AnalyzeFile() main package import path = %q", mainPkg.ImportPath)
	}
	if mainPkg.SourceFileCount == 0 {
		t.Fatalf("AnalyzeFile() main package source file count = %d", mainPkg.SourceFileCount)
	}
	if !mainPkg.ModuleLocal {
		t.Fatalf("AnalyzeFile() main package module_local = %#v", mainPkg)
	}
	runtimePkg, ok := findPackage(got.Packages, "runtime")
	if !ok {
		t.Fatalf("AnalyzeFile() packages missing runtime: %#v", got.Packages)
	}
	if runtimePkg.ImportPath != "runtime" {
		t.Fatalf("AnalyzeFile() runtime package import path = %q", runtimePkg.ImportPath)
	}
	if runtimePkg.SourceFileCount == 0 {
		t.Fatalf("AnalyzeFile() runtime package source file count = %d", runtimePkg.SourceFileCount)
	}
	if runtimePkg.ModuleLocal {
		t.Fatalf("AnalyzeFile() runtime package module_local = %#v", runtimePkg)
	}
	if len(got.Types) == 0 {
		t.Fatal("AnalyzeFile() returned no types")
	}
	mainType, ok := findType(got.Types, "main.fixtureCounter")
	if !ok {
		t.Fatalf("AnalyzeFile() types missing main.fixtureCounter: %#v", got.Types)
	}
	if mainType.Package != "main" || mainType.ImportPath != "example.com/gorevealfixture" || mainType.SourceFileCount != 1 || !mainType.ModuleLocal || !mainType.UserMeaningful {
		t.Fatalf("AnalyzeFile() main type metadata = %#v", mainType)
	}
	runtimeType, ok := findType(got.Types, "**runtime.g")
	if !ok {
		t.Fatalf("AnalyzeFile() types missing **runtime.g: %#v", got.Types[:min(10, len(got.Types))])
	}
	if runtimeType.Package != "runtime" || runtimeType.ImportPath != "runtime" || runtimeType.SourceFileCount != 0 || runtimeType.ModuleLocal || runtimeType.UserMeaningful {
		t.Fatalf("AnalyzeFile() runtime type metadata = %#v", runtimeType)
	}
	if len(got.Strings) == 0 {
		t.Fatal("AnalyzeFile() returned no strings")
	}
	if got.SourceTree == nil {
		t.Fatal("AnalyzeFile() returned no source tree")
	}
	if got.Refined == nil {
		t.Fatal("AnalyzeFile() returned no refined layer")
	}
	if got.Refined.Functions[0].Name != got.Functions[0].Name {
		t.Fatalf("AnalyzeFile() refined function mismatch: raw=%q refined=%q", got.Functions[0].Name, got.Refined.Functions[0].Name)
	}
}

func findPackage(pkgs []schema.Package, want string) (schema.Package, bool) {
	for _, pkg := range pkgs {
		if pkg.Name == want {
			return pkg, true
		}
	}

	return schema.Package{}, false
}

func findType(types []schema.Type, want string) (schema.Type, bool) {
	for _, typ := range types {
		if typ.Name == want {
			return typ, true
		}
	}

	return schema.Type{}, false
}
