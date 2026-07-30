package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ioplane/goreveal/schema"
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
	if got.Runtime.TrustSummary != "symbol_backed" {
		t.Fatalf("AnalyzeFile() runtime trust summary = %q, want %q", got.Runtime.TrustSummary, "symbol_backed")
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
	if got.Peeling == nil {
		t.Fatal("AnalyzeFile() returned no peeling layer")
	}
	mainPeel, ok := findPeelingFunction(got.Peeling.Functions, "main.main")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing main.main: %#v", got.Peeling.Functions)
	}
	if mainPeel.Classification != schema.PeelingClassUser {
		t.Fatalf("AnalyzeFile() main.main peeling = %#v", mainPeel)
	}
	if mainPeel.ClassificationEvidence != schema.PeelingEvidenceModuleLocal {
		t.Fatalf("AnalyzeFile() main.main peeling evidence = %#v", mainPeel)
	}
	runtimePeel, ok := findPeelingFunction(got.Peeling.Functions, "runtime.newobject")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing runtime.newobject: %#v", got.Peeling.Functions[:min(10, len(got.Peeling.Functions))])
	}
	if runtimePeel.Classification != schema.PeelingClassRuntime {
		t.Fatalf("AnalyzeFile() runtime.newobject peeling = %#v", runtimePeel)
	}
	if runtimePeel.ClassificationEvidence != schema.PeelingEvidenceRuntimeImportPath {
		t.Fatalf("AnalyzeFile() runtime.newobject peeling evidence = %#v", runtimePeel)
	}
	fmtPeel, ok := findPeelingFunction(got.Peeling.Functions, "fmt.Fprintln")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing fmt.Fprintln: %#v", got.Peeling.Functions[:min(10, len(got.Peeling.Functions))])
	}
	if fmtPeel.Classification != schema.PeelingClassStdlib {
		t.Fatalf("AnalyzeFile() fmt.Fprintln peeling = %#v", fmtPeel)
	}
	if fmtPeel.ClassificationEvidence != schema.PeelingEvidenceStdlibImportPath {
		t.Fatalf("AnalyzeFile() fmt.Fprintln peeling evidence = %#v", fmtPeel)
	}
	mainPeelPkg, ok := findPeelingPackage(got.Peeling.Packages, "main", "example.com/gorevealfixture")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing main package summary: %#v", got.Peeling.Packages)
	}
	if mainPeelPkg.PrimaryClassification != schema.PeelingClassUser || mainPeelPkg.UserFunctionCount == 0 || !mainPeelPkg.ModuleLocal {
		t.Fatalf("AnalyzeFile() main package peeling summary = %#v", mainPeelPkg)
	}
	runtimePeelPkg, ok := findPeelingPackage(got.Peeling.Packages, "runtime", "runtime")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing runtime package summary: %#v", got.Peeling.Packages)
	}
	if runtimePeelPkg.PrimaryClassification != schema.PeelingClassRuntime || runtimePeelPkg.RuntimeFunctionCount == 0 {
		t.Fatalf("AnalyzeFile() runtime package peeling summary = %#v", runtimePeelPkg)
	}
	fmtPeelPkg, ok := findPeelingPackage(got.Peeling.Packages, "fmt", "fmt")
	if !ok {
		t.Fatalf("AnalyzeFile() peeling missing fmt package summary: %#v", got.Peeling.Packages)
	}
	if fmtPeelPkg.PrimaryClassification != schema.PeelingClassStdlib || fmtPeelPkg.StdlibFunctionCount == 0 {
		t.Fatalf("AnalyzeFile() fmt package peeling summary = %#v", fmtPeelPkg)
	}
	if got.Refined == nil {
		t.Fatal("AnalyzeFile() returned no refined layer")
	}
	if got.Refined.Functions[0].Name != got.Functions[0].Name {
		t.Fatalf("AnalyzeFile() refined function mismatch: raw=%q refined=%q", got.Functions[0].Name, got.Refined.Functions[0].Name)
	}
}

func TestAnalyzeFileStrippedFixtureRuntimeTrustSummary(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "corpus", "fixtures", "go-elf-stripped-linux-amd64", "fixture.bin")

	got, err := New().AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}
	if got.Runtime == nil {
		t.Fatal("AnalyzeFile() Runtime = nil")
	}
	if got.Runtime.TrustSummary != "go_module_fallback" {
		t.Fatalf("AnalyzeFile() stripped runtime trust summary = %q, want %q", got.Runtime.TrustSummary, "go_module_fallback")
	}
	if !got.Runtime.FirstModuleDataFromGoModuleFallback {
		t.Fatalf("AnalyzeFile() stripped fallback bit = %#v", got.Runtime)
	}
}

func TestAnalyzePEFixtureIncludesBoundedRuntimeSectionHeuristic(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "corpus", "fixtures", "go-pe-buildinfo-windows-amd64", "fixture.exe")

	got, err := New().AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}
	if got.Input.Format != "pe" {
		t.Fatalf("AnalyzeFile() format = %q, want %q", got.Input.Format, "pe")
	}
	if got.BuildInfo == nil {
		t.Fatal("AnalyzeFile() BuildInfo = nil")
	}
	if got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("AnalyzeFile() BuildInfo.Path = %q", got.BuildInfo.Path)
	}
	if got.Runtime == nil {
		t.Fatal("AnalyzeFile() Runtime = nil")
	}
	if got.Runtime.TrustSummary != "section_heuristic" {
		t.Fatalf("AnalyzeFile() runtime trust summary = %q, want %q", got.Runtime.TrustSummary, "section_heuristic")
	}
	if got.Runtime.PETextSectionAddr == 0 || got.Runtime.PERdataSectionAddr == 0 {
		t.Fatalf("AnalyzeFile() PE runtime sections = %#v", got.Runtime)
	}
	if got.Runtime.PEPclntabMagicSection != ".rdata" || got.Runtime.PEPclntabMagicAddr == 0 || got.Runtime.PEPclntabMagicCount == 0 {
		t.Fatalf("AnalyzeFile() PE runtime pclntab magic = %#v", got.Runtime)
	}
	if got.Runtime.PEPclntabHeaderSection != ".rdata" || got.Runtime.PEPclntabHeaderAddr == 0 {
		t.Fatalf("AnalyzeFile() PE runtime pclntab header = %#v", got.Runtime)
	}
	if got.Runtime.PEPclntabHeaderMagic != "f1ffffff" || got.Runtime.PEPclntabHeaderQuantum != 1 || got.Runtime.PEPclntabHeaderPointerSize != 8 {
		t.Fatalf("AnalyzeFile() PE runtime pclntab header fields = %#v", got.Runtime)
	}
	if len(got.Functions) == 0 {
		t.Fatal("AnalyzeFile() returned no PE functions")
	}
	if len(got.Packages) == 0 {
		t.Fatal("AnalyzeFile() returned no PE packages")
	}
	if got.SourceTree == nil {
		t.Fatal("AnalyzeFile() returned no PE source tree")
	}
	if !got.SourceTree.PathlessFileEvidence || len(got.SourceTree.Files) == 0 {
		t.Fatalf("AnalyzeFile() PE source tree = %#v", got.SourceTree)
	}
	if len(got.Strings) != 0 {
		t.Fatalf("AnalyzeFile() unexpected PE strings = %#v", got.Strings)
	}
	mainPkg, ok := findPackage(got.Packages, "main")
	if !ok {
		t.Fatalf("AnalyzeFile() PE packages missing main: %#v", got.Packages)
	}
	if mainPkg.ImportPath != "example.com/gorevealfixture" || !mainPkg.ModuleLocal {
		t.Fatalf("AnalyzeFile() PE main package = %#v", mainPkg)
	}
	if got.Peeling == nil {
		t.Fatal("AnalyzeFile() returned no PE peeling")
	}
	mainPeel, ok := findPeelingFunction(got.Peeling.Functions, "main.main")
	if !ok {
		t.Fatalf("AnalyzeFile() PE peeling missing main.main: %#v", got.Peeling.Functions[:min(10, len(got.Peeling.Functions))])
	}
	if mainPeel.Classification != schema.PeelingClassUser || mainPeel.ClassificationEvidence != schema.PeelingEvidenceModuleLocal {
		t.Fatalf("AnalyzeFile() PE main.main peeling = %#v", mainPeel)
	}
}

func TestShouldPreferFunctionSourceTree(t *testing.T) {
	t.Parallel()

	if !shouldPreferFunctionSourceTree(
		schema.SourceTree{},
		schema.SourceTree{Files: []string{"main.go"}, PathlessFileEvidence: true},
	) {
		t.Fatal("shouldPreferFunctionSourceTree() = false, want true when function tree adds module-local files")
	}

	if shouldPreferFunctionSourceTree(
		schema.SourceTree{Files: []string{"main.go"}},
		schema.SourceTree{Files: []string{"main.go"}, PathlessFileEvidence: true},
	) {
		t.Fatal("shouldPreferFunctionSourceTree() = true, want false when DWARF tree already has module-local files")
	}

	if shouldPreferFunctionSourceTree(
		schema.SourceTree{},
		schema.SourceTree{},
	) {
		t.Fatal("shouldPreferFunctionSourceTree() = true, want false when fallback adds no files")
	}
}

func TestAnnotateELFFunctionRecoveryBlocker(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:   "unknown",
			ELFPclntabFunctionCountHint: 123,
		},
	}

	annotateELFFunctionRecoveryBlocker(&analysis)
	if analysis.Runtime.ELFFunctionRecoveryBlocker != "custom_pclntab_magic" {
		t.Fatalf("annotateELFFunctionRecoveryBlocker() blocker = %q", analysis.Runtime.ELFFunctionRecoveryBlocker)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:   "unknown",
			ELFPclntabFunctionCountHint: 123,
		},
		Functions: []schema.Function{{Name: "main.main"}},
	}
	annotateELFFunctionRecoveryBlocker(&analysis)
	if analysis.Runtime.ELFFunctionRecoveryBlocker != "" {
		t.Fatalf("annotateELFFunctionRecoveryBlocker() unexpectedly set blocker = %q", analysis.Runtime.ELFFunctionRecoveryBlocker)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind: "unknown",
		},
	}
	annotateELFFunctionRecoveryBlocker(&analysis)
	if analysis.Runtime.ELFFunctionRecoveryBlocker != "" {
		t.Fatalf("annotateELFFunctionRecoveryBlocker() unexpectedly set blocker without function count hint = %q", analysis.Runtime.ELFFunctionRecoveryBlocker)
	}
}

func TestAnnotateELFFunctionFoothold(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:           "unknown",
			ELFPclntabFunctionCountHint:         123,
			ELFFunctabPCOffsetsMonotonic:        true,
			ELFFunctabPCAddrHintsWithinText:     true,
			ELFFunctabPCAddrSample:              []uint64{0x401000, 0x401120},
			ELFFunctabPCAddrSampleAllWithinText: true,
		},
	}

	annotateELFFunctionFoothold(&analysis)
	if analysis.Runtime.ELFFunctionFoothold != "address_only" {
		t.Fatalf("annotateELFFunctionFoothold() foothold = %q", analysis.Runtime.ELFFunctionFoothold)
	}
	if analysis.Runtime.ELFFunctionFootholdCountHint != 123 {
		t.Fatalf(
			"annotateELFFunctionFoothold() count hint = %d",
			analysis.Runtime.ELFFunctionFootholdCountHint,
		)
	}
	if analysis.Runtime.ELFFunctionFootholdTextSource != "" {
		t.Fatalf("annotateELFFunctionFoothold() unexpected text source without text range = %#v", analysis.Runtime)
	}
	if analysis.Runtime.ELFFunctionFootholdStartAddr != 0 || analysis.Runtime.ELFFunctionFootholdEndAddr != 0 {
		t.Fatalf("annotateELFFunctionFoothold() unexpected span without addr hints = %#v", analysis.Runtime)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:           "unknown",
			ELFPclntabFunctionCountHint:         123,
			ELFFunctabPCOffsetsMonotonic:        true,
			ELFFunctabPCAddrHintsWithinText:     true,
			ELFFunctabPCAddrSample:              []uint64{0x401000, 0x401120},
			ELFFunctabPCAddrSampleAllWithinText: true,
		},
		Functions: []schema.Function{{Name: "main.main"}},
	}
	annotateELFFunctionFoothold(&analysis)
	if analysis.Runtime.ELFFunctionFoothold != "" || analysis.Runtime.ELFFunctionFootholdCountHint != 0 {
		t.Fatalf("annotateELFFunctionFoothold() unexpectedly set foothold = %#v", analysis.Runtime)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:           "unknown",
			ELFPclntabFunctionCountHint:         123,
			ELFFunctabPCOffsetsMonotonic:        true,
			ELFFunctabPCAddrHintsWithinText:     true,
			ELFFunctabPCAddrSample:              []uint64{0x401000, 0x401120},
			ELFFunctabPCAddrSampleAllWithinText: false,
		},
	}
	annotateELFFunctionFoothold(&analysis)
	if analysis.Runtime.ELFFunctionFoothold != "" || analysis.Runtime.ELFFunctionFootholdCountHint != 0 {
		t.Fatalf("annotateELFFunctionFoothold() unexpectedly set foothold without within-text sample = %#v", analysis.Runtime)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:           "unknown",
			ELFPclntabFunctionCountHint:         123,
			ELFFunctabPCOffsetsMonotonic:        true,
			ELFTextSectionAddr:                  0x501000,
			ELFTextSectionEndInclusive:          0x5015ff,
			ELFFunctabPCAddrHintsWithinText:     true,
			ELFFunctabPCAddrSample:              []uint64{0x501000, 0x501120},
			ELFFunctabPCAddrSampleAllWithinText: true,
		},
	}
	annotateELFFunctionFoothold(&analysis)
	if analysis.Runtime.ELFFunctionFoothold != "address_only" || analysis.Runtime.ELFFunctionFootholdCountHint != 123 {
		t.Fatalf("annotateELFFunctionFoothold() ELF text fallback foothold = %#v", analysis.Runtime)
	}
	if analysis.Runtime.ELFFunctionFootholdTextSource != "elf_text_section" {
		t.Fatalf("annotateELFFunctionFoothold() ELF text fallback source = %q", analysis.Runtime.ELFFunctionFootholdTextSource)
	}
	if analysis.Runtime.ELFFunctionFootholdStartAddr != 0 || analysis.Runtime.ELFFunctionFootholdEndAddr != 0 {
		t.Fatalf("annotateELFFunctionFoothold() unexpected ELF text fallback span without addr hints = %#v", analysis.Runtime)
	}

	analysis = schema.Analysis{
		Runtime: &schema.RuntimeMetadata{
			ELFPclntabHeaderMagicKind:           "unknown",
			ELFPclntabFunctionCountHint:         123,
			ELFFunctabPCOffsetsMonotonic:        true,
			ModuledataTextAddr:                  0x401000,
			ModuledataTextEndInclusive:          0x4015ff,
			ELFFunctabFirstPCAddrHint:           0x401000,
			ELFFunctabLastPCAddrHint:            0x401300,
			ELFFunctabPCAddrHintsWithinText:     true,
			ELFFunctabPCAddrSample:              []uint64{0x401000, 0x401120},
			ELFFunctabPCAddrSampleAllWithinText: true,
		},
	}
	annotateELFFunctionFoothold(&analysis)
	if analysis.Runtime.ELFFunctionFoothold != "address_only" || analysis.Runtime.ELFFunctionFootholdCountHint != 123 {
		t.Fatalf("annotateELFFunctionFoothold() moduledata text foothold = %#v", analysis.Runtime)
	}
	if analysis.Runtime.ELFFunctionFootholdTextSource != "moduledata_text" {
		t.Fatalf("annotateELFFunctionFoothold() moduledata text source = %q", analysis.Runtime.ELFFunctionFootholdTextSource)
	}
	if analysis.Runtime.ELFFunctionFootholdStartAddr != 0x401000 || analysis.Runtime.ELFFunctionFootholdEndAddr != 0x401300 {
		t.Fatalf("annotateELFFunctionFoothold() moduledata text span = %#v", analysis.Runtime)
	}
}

func TestAnalyzeMachOFixtureIncludesBoundedFunctionFoothold(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "corpus", "fixtures", "go-macho-buildinfo-darwin-amd64", "fixture.bin")

	got, err := New().AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}
	if got.Input.Format != "macho" {
		t.Fatalf("AnalyzeFile() format = %q, want %q", got.Input.Format, "macho")
	}
	if got.BuildInfo == nil {
		t.Fatal("AnalyzeFile() BuildInfo = nil")
	}
	if got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("AnalyzeFile() BuildInfo.Path = %q", got.BuildInfo.Path)
	}
	if got.Runtime != nil {
		t.Fatalf("AnalyzeFile() unexpected Mach-O runtime = %#v", got.Runtime)
	}
	if len(got.Functions) == 0 {
		t.Fatal("AnalyzeFile() returned no Mach-O functions")
	}
	if len(got.Packages) == 0 {
		t.Fatal("AnalyzeFile() returned no Mach-O packages")
	}
	if got.SourceTree == nil {
		t.Fatal("AnalyzeFile() returned no Mach-O source tree")
	}
	if !got.SourceTree.PathlessFileEvidence || len(got.SourceTree.Files) == 0 {
		t.Fatalf("AnalyzeFile() Mach-O source tree = %#v", got.SourceTree)
	}
	mainPkg, ok := findPackage(got.Packages, "main")
	if !ok {
		t.Fatalf("AnalyzeFile() Mach-O packages missing main: %#v", got.Packages)
	}
	if mainPkg.ImportPath != "example.com/gorevealfixture" || !mainPkg.ModuleLocal {
		t.Fatalf("AnalyzeFile() Mach-O main package = %#v", mainPkg)
	}
	if got.Peeling == nil {
		t.Fatal("AnalyzeFile() returned no Mach-O peeling layer")
	}
	mainPeel, ok := findPeelingFunction(got.Peeling.Functions, "main.main")
	if !ok {
		t.Fatalf("AnalyzeFile() Mach-O peeling missing main.main: %#v", got.Peeling.Functions)
	}
	if mainPeel.Classification != schema.PeelingClassUser || mainPeel.ClassificationEvidence != schema.PeelingEvidenceModuleLocal {
		t.Fatalf("AnalyzeFile() Mach-O main.main peeling = %#v", mainPeel)
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

func findPeelingFunction(funcs []schema.PeelingFunction, want string) (schema.PeelingFunction, bool) {
	for _, fn := range funcs {
		if fn.Name == want {
			return fn, true
		}
	}

	return schema.PeelingFunction{}, false
}

func findPeelingPackage(pkgs []schema.PeelingPackage, wantName, wantImportPath string) (schema.PeelingPackage, bool) {
	for _, pkg := range pkgs {
		if pkg.Name == wantName && pkg.ImportPath == wantImportPath {
			return pkg, true
		}
	}

	return schema.PeelingPackage{}, false
}
