package schema_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/schema"
)

func richExportAnalysis() schema.Analysis {
	return schema.Analysis{
		Input: schema.Input{
			Path:   "/tmp/rich-sample.bin",
			Size:   8192,
			Format: "elf",
		},
		BuildInfo: &schema.BuildInfo{
			GoVersion: "go1.26.1",
			Path:      "example.com/rich",
			Provenance: schema.Provenance{
				Source:     "core.buildinfo",
				Confidence: "high",
			},
		},
		Runtime: &schema.RuntimeMetadata{
			TrustSummary:                         schema.RuntimeTrustSummarySymbolBacked,
			ELFPclntabHeaderMagic:                "0xfffffff1",
			ELFPclntabHeaderMagicKind:            "go1.20",
			ELFPclntabHeaderQuantum:              1,
			ELFPclntabHeaderPointerSize:          8,
			ELFPclntabFunctionCountHint:          21,
			ELFPclntabFileCountHint:              3,
			ELFPclntabFuncnametabOffsetHint:      64,
			ELFPclntabCuOffsetHint:               128,
			ELFPclntabFiletabOffsetHint:          192,
			ELFPclntabPctabOffsetHint:            256,
			ELFPclntabFunctabOffsetHint:          320,
			ELFFunctabFirstPCOffsetHint:          16,
			ELFFunctabLastPCOffsetHint:           4096,
			ELFFunctabPCOffsetsMonotonic:         true,
			ELFTextSectionAddr:                   0x401000,
			ELFTextSectionEndInclusive:           0x402fff,
			ELFFunctabFirstPCAddrHint:            0x401010,
			ELFFunctabLastPCAddrHint:             0x402000,
			ELFFunctabPCAddrHintsWithinText:      true,
			ELFFunctabPCAddrSample:               []uint64{0x401010, 0x401100},
			ELFFunctabPCAddrSampleAllWithinText:  true,
			ELFFunctionFoothold:                  "address_only",
			ELFFunctionFootholdCountHint:         21,
			ELFFunctionFootholdTextSource:        "moduledata_text",
			ELFFunctionFootholdStartAddr:         0x401000,
			ELFFunctionFootholdEndAddr:           0x403000,
			ELFFunctionRecoveryBlocker:           "fixture_only",
			PETextSectionAddr:                    0x140001000,
			PETextSectionSize:                    4096,
			PERdataSectionAddr:                   0x140003000,
			PERdataSectionSize:                   2048,
			PEPclntabMagicSection:                ".rdata",
			PEPclntabMagicAddr:                   0x140003040,
			PEPclntabMagicCount:                  2,
			PEPclntabHeaderSection:               ".rdata",
			PEPclntabHeaderAddr:                  0x140003040,
			PEPclntabHeaderMagic:                 "0xfffffff1",
			PEPclntabHeaderQuantum:               1,
			PEPclntabHeaderPointerSize:           8,
			FirstModuleDataAddr:                  0x4f0000,
			FirstModuleDataFromGoModuleFallback:  true,
			GopclntabAddr:                        0x4a0000,
			GopclntabSize:                        16384,
			TypelinkAddr:                         0x4c0000,
			TypelinkSize:                         128,
			TypelinkCount:                        32,
			ItablinkAddr:                         0x4c1000,
			ItablinkSize:                         64,
			ItablinkCount:                        8,
			TypelinkSample:                       []int32{-16, 32},
			TypelinkResolvedBaseAddr:             0x4b0000,
			TypelinkResolvedSample:               []uint64{0x4afff0, 0x4b0020},
			TypelinkResolvedWithinRodataCount:    2,
			TypelinkAllResolvedWithinRodata:      true,
			TypelinkMinOffset:                    -16,
			TypelinkMaxOffset:                    32,
			TypelinkNegativeCount:                1,
			TypelinkNonNegativeCount:             1,
			GoModuleAddr:                         0x4e0000,
			GoModuleSize:                         4096,
			FirstModuleDataInGoModule:            true,
			FirstModuleDataGoModuleOffset:        128,
			GoModuleWordSize:                     8,
			GoModuleWordSample:                   []uint64{1, 2, 3},
			ModuledataPCHeaderAddr:               0x4a0000,
			ModuledataPCHeaderMatchesGopclntab:   true,
			ModuledataFuncnametabSliceWordIndex:  1,
			ModuledataFuncnametabAddr:            0x4a0100,
			ModuledataFuncnametabLen:             100,
			ModuledataFuncnametabCap:             120,
			ModuledataFuncnametabWithinGopclntab: true,
			ModuledataCutabSliceWordIndex:        4,
			ModuledataCutabAddr:                  0x4a0200,
			ModuledataCutabLen:                   200,
			ModuledataCutabCap:                   220,
			ModuledataCutabWithinGopclntab:       true,
			ModuledataFiletabSliceWordIndex:      7,
			ModuledataFiletabAddr:                0x4a0300,
			ModuledataFiletabLen:                 300,
			ModuledataFiletabCap:                 320,
			ModuledataFiletabWithinGopclntab:     true,
			ModuledataPctabSliceWordIndex:        10,
			ModuledataPctabAddr:                  0x4a0400,
			ModuledataPctabLen:                   400,
			ModuledataPctabCap:                   420,
			ModuledataPctabWithinGopclntab:       true,
			ModuledataPclntableSliceWordIndex:    13,
			ModuledataPclntableAddr:              0x4a0500,
			ModuledataPclntableLen:               500,
			ModuledataPclntableCap:               520,
			ModuledataPclntableWithinGopclntab:   true,
			ModuledataTypelinkSliceWordIndex:     16,
			ModuledataTypelinkLen:                32,
			ModuledataTypelinkCap:                32,
			ModuledataItablinkSliceWordIndex:     19,
			ModuledataItablinkLen:                8,
			ModuledataItablinkCap:                8,
			ModuledataMemoryRangesWordIndex:      22,
			ModuledataNoptrdataAddr:              0x4c2000,
			ModuledataNoptrdataEnd:               0x4c2100,
			ModuledataDataAddr:                   0x4c3000,
			ModuledataDataEnd:                    0x4c3100,
			ModuledataBssAddr:                    0x4c4000,
			ModuledataBssEnd:                     0x4c4100,
			ModuledataNoptrbssAddr:               0x4c5000,
			ModuledataNoptrbssEnd:                0x4c5100,
			ModuledataRodataWordIndex:            30,
			ModuledataRodataAddr:                 0x4b0000,
			ModuledataRodataEnd:                  0x4c0000,
			ModuledataTextWordIndex:              32,
			ModuledataTextAddr:                   0x401000,
			ModuledataTextEndInclusive:           0x402fff,
			ModuledataTypesRangeWordIndex:        34,
			ModuledataTypesAddr:                  0x4b1000,
			ModuledataETypesAddr:                 0x4b2000,
			TypelinkResolvedWithinTypesCount:     2,
			TypelinkAllResolvedWithinTypes:       true,
			Provenance: schema.Provenance{
				Source:     "core.runtime.fixture",
				Confidence: "medium",
			},
		},
		Functions: []schema.Function{{
			Name:          "main.main",
			Package:       "main",
			ImportPath:    "example.com/rich",
			SourceFile:    "cmd/rich/main.go",
			SourceLine:    42,
			Autogenerated: true,
			Entry:         0x401000,
			End:           0x401100,
			ModuleLocal:   true,
			Provenance: schema.Provenance{
				Source:     "core.pclntab",
				Confidence: "high",
			},
		}},
		Packages: []schema.Package{{
			Name:               "main",
			ImportPath:         "example.com/rich",
			SourceFileCount:    1,
			FunctionCount:      1,
			HasSourceEvidence:  true,
			SourceEvidenceKind: schema.SourceEvidenceKindDWARFPaths,
			ModuleLocal:        true,
			Provenance: schema.Provenance{
				Source:     "core.packages",
				Confidence: "high",
			},
		}},
		Types: []schema.Type{{
			Name:            "main.record",
			Package:         "main",
			ImportPath:      "example.com/rich",
			Kind:            "struct",
			SourceFileCount: 1,
			ModuleLocal:     true,
			UserMeaningful:  true,
			Provenance: schema.Provenance{
				Source:     "core.types",
				Confidence: "medium",
			},
		}},
		Strings: []schema.StringCandidate{{
			Value:  "hello, wire",
			Region: ".rodata",
			Addr:   0x4b0100,
			Offset: 256,
			Provenance: schema.Provenance{
				Source:     "core.strings",
				Confidence: "medium",
			},
		}},
		SourceTree: &schema.SourceTree{
			Root:               "example.com/rich",
			SourceEvidenceKind: schema.SourceEvidenceKindDWARFPaths,
			SourceEvidenceSummary: schema.SourceEvidenceSummary{
				TreeKind:                    schema.SourceEvidenceKindDWARFPaths,
				DWARFPathPackageCount:       1,
				DWARFPathFileCount:          1,
				LineTablePackageCount:       2,
				LineTableFileCount:          3,
				PackageFallbackPackageCount: 4,
				PackageFallbackFileCount:    5,
				MixedPackageEvidenceKinds:   true,
			},
			PathlessFileEvidence: true,
			Files:                []string{"cmd/rich/main.go"},
			Packages: []schema.SourcePackage{{
				Name:               "main",
				ImportPath:         "example.com/rich",
				FunctionCount:      1,
				HasFileEvidence:    true,
				SourceEvidenceKind: schema.SourceEvidenceKindDWARFPaths,
				Files:              []string{"cmd/rich/main.go"},
			}},
			ExternalPackages: []schema.SourcePackage{{
				Name:               "fmt",
				ImportPath:         "fmt",
				FunctionCount:      2,
				HasFileEvidence:    true,
				SourceEvidenceKind: schema.SourceEvidenceKindLineTableFiles,
				Files:              []string{"src/fmt/print.go"},
			}},
		},
		Peeling: &schema.PeelingAnalysis{
			Functions: []schema.PeelingFunction{{
				Name:                   "main.main",
				Package:                "main",
				ImportPath:             "example.com/rich",
				SourceFile:             "cmd/rich/main.go",
				SourceLine:             42,
				Entry:                  0x401000,
				End:                    0x401100,
				ModuleLocal:            true,
				Classification:         schema.PeelingClassUser,
				ClassificationEvidence: schema.PeelingEvidenceModuleLocal,
			}},
			Packages: []schema.PeelingPackage{{
				Name:                    "main",
				ImportPath:              "example.com/rich",
				ModuleLocal:             true,
				FunctionCount:           5,
				UserFunctionCount:       1,
				StdlibFunctionCount:     1,
				RuntimeFunctionCount:    1,
				ThirdPartyFunctionCount: 2,
				PrimaryClassification:   schema.PeelingClassUser,
			}},
			Provenance: schema.Provenance{
				Source:     "engine.peeling",
				Confidence: "medium",
			},
		},
		Refined: &schema.RefinedAnalysis{
			Functions: []schema.RefinedFunction{{Name: "main.refinedMain"}},
			Types:     []schema.RefinedType{{Name: "main.RefinedRecord"}},
			Strings:   []schema.RefinedString{{Value: "decoded wire"}},
			Passes:    []string{"names", "strings"},
		},
	}
}

func encodeExportV1(t *testing.T, value any) []byte {
	t.Helper()

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	return output.Bytes()
}

func TestIDAExportV1FrozenBytes(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile(filepath.Join("testdata", "export-v1", "ida.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := encodeExportV1(t, schema.NewIDAExport(richExportAnalysis()))
	if !bytes.Equal(got, want) {
		t.Fatalf("IDA v1 wire bytes changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestGhidraExportV1FrozenBytes(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile(filepath.Join("testdata", "export-v1", "ghidra.json"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	got := encodeExportV1(t, schema.NewGhidraExport(richExportAnalysis()))
	if !bytes.Equal(got, want) {
		t.Fatalf("Ghidra v1 wire bytes changed\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestFixtureIDAExportSimulation(t *testing.T) {
	t.Parallel()

	fixture := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	analysis, err := engine.New().AnalyzeFile(context.Background(), fixture)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	got := schema.NewIDAExport(analysis)

	if got.Contract != schema.IDAExportContractV1 {
		t.Fatalf("contract = %q", got.Contract)
	}
	if got.BuildInfo == nil || got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("build info = %#v", got.BuildInfo)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != schema.RuntimeTrustSummarySymbolBacked {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) == 0 {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if len(got.Peeling.Packages) == 0 {
		t.Fatalf("peeling packages = %#v", got.Peeling)
	}
	if len(got.Functions) == 0 || got.Functions[0].Entry == 0 {
		t.Fatalf("functions = %#v", got.Functions)
	}
	if len(got.Packages) == 0 || got.Packages[0].Name == "" {
		t.Fatalf("packages = %#v", got.Packages)
	}
}

func TestFixtureGhidraExportSimulation(t *testing.T) {
	t.Parallel()

	fixture := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	analysis, err := engine.New().AnalyzeFile(context.Background(), fixture)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	got := schema.NewGhidraExport(analysis)

	if got.Contract != schema.GhidraExportContractV1 {
		t.Fatalf("contract = %q", got.Contract)
	}
	if got.Program.ModulePath != "example.com/gorevealfixture" {
		t.Fatalf("module path = %q", got.Program.ModulePath)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != schema.RuntimeTrustSummarySymbolBacked {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) == 0 {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if len(got.Peeling.Packages) == 0 {
		t.Fatalf("peeling packages = %#v", got.Peeling)
	}
	if len(got.Symbols) == 0 || got.Symbols[0].Address == 0 {
		t.Fatalf("symbols = %#v", got.Symbols)
	}
	if got.SourceTree == nil || got.SourceTree.Root != "example.com/gorevealfixture" {
		t.Fatalf("source tree = %#v", got.SourceTree)
	}
}
