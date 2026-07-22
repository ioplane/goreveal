package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/dantte-lp/goreveal/core/recoveryerr"
	recoverystrings "github.com/dantte-lp/goreveal/core/strings"
	"github.com/dantte-lp/goreveal/schema"
)

func TestAnalyzerZeroValueMatchesNew(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	want, err := New().AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("New().AnalyzeFile() error = %v", err)
	}
	got, err := (Analyzer{}).AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("Analyzer{}.AnalyzeFile() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Analyzer{}.AnalyzeFile() differs from New():\nwant=%#v\ngot=%#v", want, got)
	}
}

func TestAnalyzerRejectsPartialStageOperations(t *testing.T) {
	t.Parallel()

	partial := Analyzer{ops: stageOps{
		buildInfo: func(string) (schema.BuildInfo, error) { return schema.BuildInfo{}, nil },
	}}
	_, err := partial.AnalyzeFile(context.Background(), writeIngestibleELF(t))
	if err == nil || !strings.Contains(err.Error(), "incomplete stage operations") || !strings.Contains(err.Error(), "runtime") {
		t.Fatalf("AnalyzeFile() error = %v, want explicit incomplete stage operations", err)
	}

	_, err = partial.AnalyzeFile(nil, writeIngestibleELF(t)) //nolint:staticcheck // Nil is intentional to verify error precedence.
	if err == nil || err.Error() != "analyze file: nil context" {
		t.Fatalf("AnalyzeFile(nil) error = %v, want nil context precedence", err)
	}
}

func TestAnalyzeFileRecordsStageFailure(t *testing.T) {
	t.Parallel()

	path := writeIngestibleELF(t)
	ops := successfulStageOps()
	ops.functions = func(string) ([]schema.Function, error) {
		return nil, errors.New("fixture failure")
	}

	got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), path)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	diagnostic, ok := diagnosticForStage(got.Diagnostics, schema.AnalysisStageFunctions)
	if !ok {
		t.Fatalf("AnalyzeFile() diagnostics = %#v, want functions diagnostic", got.Diagnostics)
	}
	if diagnostic.Status != schema.StageStatusFailed || diagnostic.Code != "stage_failed" {
		t.Fatalf("functions diagnostic = %#v, want failed/stage_failed", diagnostic)
	}
	if len(got.Functions) != 0 {
		t.Fatalf("AnalyzeFile() functions = %#v, want no failed-stage result", got.Functions)
	}
	if _, ok := diagnosticForStage(got.Diagnostics, schema.AnalysisStagePackages); ok {
		t.Fatalf("AnalyzeFile() diagnostics = %#v, want no derived packages claim", got.Diagnostics)
	}
	if _, ok := diagnosticForStage(got.Diagnostics, schema.AnalysisStagePeeling); ok {
		t.Fatalf("AnalyzeFile() diagnostics = %#v, want no derived peeling claim", got.Diagnostics)
	}
	if got.Packages != nil || got.Peeling != nil {
		t.Fatalf("AnalyzeFile() derived values = packages:%#v peeling:%#v", got.Packages, got.Peeling)
	}
}

func TestAnalyzeFileRefinesOnlyAvailableFamilies(t *testing.T) {
	t.Parallel()

	ops := successfulStageOps()
	ops.functions = func(string) ([]schema.Function, error) {
		return nil, recoveryerr.NewUnavailable(
			recoveryerr.CodePclntabNotFound,
			"function evidence is absent",
			nil,
		)
	}
	ops.refine = func(_ context.Context, analysis schema.Analysis) (schema.RefinedAnalysis, error) {
		if len(analysis.Functions) != 0 || len(analysis.Packages) != 0 {
			return schema.RefinedAnalysis{}, errors.New("unavailable function family leaked into refinement")
		}
		return schema.RefinedAnalysis{
			Functions: []schema.RefinedFunction{{Name: "invented.function"}},
			Packages:  []schema.RefinedPackage{{Name: "invented-package"}},
			Types:     []schema.RefinedType{{Name: analysis.Types[0].Name}},
			Strings:   []schema.RefinedString{{Value: analysis.Strings[0].Value}},
		}, nil
	}

	got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}
	if got.Refined == nil {
		t.Fatal("AnalyzeFile() refined = nil, want available type/string families")
	}
	if len(got.Refined.Functions) != 0 || len(got.Refined.Packages) != 0 ||
		len(got.Refined.Types) != 1 || len(got.Refined.Strings) != 1 {
		t.Fatalf("AnalyzeFile() refined = %#v, want only type/string families", got.Refined)
	}
	if diagnostic, ok := diagnosticForStage(got.Diagnostics, schema.AnalysisStageRefinement); !ok || diagnostic.Status != schema.StageStatusAvailable {
		t.Fatalf("AnalyzeFile() refinement diagnostic = %#v, exists=%v", diagnostic, ok)
	}
}

func TestAnalyzeFileRefinesAvailableStringsWithoutFunctionEvidence(t *testing.T) {
	t.Parallel()

	ops := successfulStageOps()
	ops.functions = func(string) ([]schema.Function, error) {
		return nil, recoveryerr.NewUnavailable(recoveryerr.CodePclntabNotFound, "function evidence is absent", nil)
	}
	ops.types = func(string) ([]schema.Type, error) {
		return nil, recoveryerr.NewUnavailable(recoveryerr.CodeDWARFNotFound, "type evidence is absent", nil)
	}
	ops.refine = func(_ context.Context, analysis schema.Analysis) (schema.RefinedAnalysis, error) {
		if len(analysis.Functions) != 0 || len(analysis.Packages) != 0 || len(analysis.Types) != 0 {
			return schema.RefinedAnalysis{}, errors.New("nonavailable raw family leaked into refinement")
		}
		return schema.RefinedAnalysis{
			Strings: []schema.RefinedString{{Value: analysis.Strings[0].Value}},
		}, nil
	}

	got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}
	if got.Refined == nil || len(got.Refined.Strings) != 1 || len(got.Refined.Functions) != 0 ||
		len(got.Refined.Packages) != 0 || len(got.Refined.Types) != 0 {
		t.Fatalf("AnalyzeFile() refined = %#v, want strings only", got.Refined)
	}
}

func TestAnalyzeFileStageStatusMatrix(t *testing.T) {
	t.Parallel()

	type stageCase struct {
		name   string
		stage  schema.AnalysisStage
		inject func(*stageOps, error)
	}

	stages := []stageCase{
		{
			name:  "build_info",
			stage: schema.AnalysisStageBuildInfo,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.buildInfo = func(string) (schema.BuildInfo, error) { return schema.BuildInfo{}, err }
				}
			},
		},
		{
			name:  "runtime",
			stage: schema.AnalysisStageRuntime,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.runtime = func(string) (schema.RuntimeMetadata, error) { return schema.RuntimeMetadata{}, err }
				}
			},
		},
		{
			name:  "functions",
			stage: schema.AnalysisStageFunctions,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.functions = func(string) ([]schema.Function, error) { return nil, err }
				}
			},
		},
		{
			name:  "types",
			stage: schema.AnalysisStageTypes,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.types = func(string) ([]schema.Type, error) { return nil, err }
				}
			},
		},
		{
			name:  "strings",
			stage: schema.AnalysisStageStrings,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.strings = func(string) (recoverystrings.Result, error) { return recoverystrings.Result{}, err }
				}
			},
		},
		{
			name:  "source_tree",
			stage: schema.AnalysisStageSourceTree,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.sourceTree = func(string, schema.Analysis) (schema.SourceTree, error) {
						return schema.SourceTree{}, err
					}
				}
			},
		},
		{
			name:  "refinement",
			stage: schema.AnalysisStageRefinement,
			inject: func(ops *stageOps, err error) {
				if err != nil {
					ops.refine = func(context.Context, schema.Analysis) (schema.RefinedAnalysis, error) {
						return schema.RefinedAnalysis{}, err
					}
				}
			},
		},
	}

	outcomes := []struct {
		name    string
		err     error
		status  schema.StageStatus
		code    string
		message string
	}{
		{name: "available", status: schema.StageStatusAvailable},
		{
			name:    "unavailable",
			err:     recoveryerr.NewUnavailable(recoveryerr.Code("fixture_unavailable"), "fixture unavailable", nil),
			status:  schema.StageStatusUnavailable,
			code:    "fixture_unavailable",
			message: "fixture unavailable",
		},
		{
			name:    "unsupported",
			err:     recoveryerr.NewUnsupported(recoveryerr.Code("fixture_unsupported"), "fixture unsupported", nil),
			status:  schema.StageStatusUnsupported,
			code:    "fixture_unsupported",
			message: "fixture unsupported",
		},
		{
			name:    "failed",
			err:     errors.New("fixture failure"),
			status:  schema.StageStatusFailed,
			code:    "stage_failed",
			message: "fixture failure",
		},
	}

	for _, stage := range stages {
		for _, outcome := range outcomes {
			t.Run(stage.name+"/"+outcome.name, func(t *testing.T) {
				t.Parallel()

				ops := successfulStageOps()
				stage.inject(&ops, outcome.err)

				got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
				if err != nil {
					t.Fatalf("AnalyzeFile() error = %v", err)
				}

				diagnostic, ok := diagnosticForStage(got.Diagnostics, stage.stage)
				if !ok {
					t.Fatalf("diagnostics = %#v, want %q", got.Diagnostics, stage.stage)
				}
				if diagnostic.Status != outcome.status || diagnostic.Code != outcome.code || diagnostic.Message != outcome.message {
					t.Fatalf(
						"diagnostic = %#v, want status=%q code=%q message=%q",
						diagnostic,
						outcome.status,
						outcome.code,
						outcome.message,
					)
				}
				if outcome.err != nil && stage.stage == schema.AnalysisStageSourceTree && got.SourceTree != nil {
					t.Fatalf("source-tree error published payload %#v", got.SourceTree)
				}
				if outcome.err != nil && stage.stage == schema.AnalysisStageRefinement && got.Refined != nil {
					t.Fatalf("refinement error published payload %#v", got.Refined)
				}
				assertOrderedUniqueDiagnostics(t, got.Diagnostics)
			})
		}
	}
}

func TestAnalyzeFileRecordsDerivedStageAvailability(t *testing.T) {
	t.Parallel()

	got, err := newAnalyzerForTest(successfulStageOps()).AnalyzeFile(context.Background(), writeIngestibleELF(t))
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	want := []schema.AnalysisStage{
		schema.AnalysisStageBuildInfo,
		schema.AnalysisStageRuntime,
		schema.AnalysisStageFunctions,
		schema.AnalysisStagePackages,
		schema.AnalysisStageTypes,
		schema.AnalysisStageStrings,
		schema.AnalysisStageSourceTree,
		schema.AnalysisStagePeeling,
		schema.AnalysisStageRefinement,
	}
	if len(got.Diagnostics) != len(want) {
		t.Fatalf("diagnostics = %#v, want stages %#v", got.Diagnostics, want)
	}
	for i, stage := range want {
		if got.Diagnostics[i].Stage != stage || got.Diagnostics[i].Status != schema.StageStatusAvailable {
			t.Fatalf("diagnostics[%d] = %#v, want %q/available", i, got.Diagnostics[i], stage)
		}
	}
}

func TestAnalyzeFileRejectsEmptyStageEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		stage   schema.AnalysisStage
		code    string
		message string
		inject  func(*stageOps)
		assert  func(*testing.T, schema.Analysis)
	}{
		{
			name:    "build info",
			stage:   schema.AnalysisStageBuildInfo,
			code:    string(recoveryerr.CodeBuildInfoNotFound),
			message: "Go build info evidence is absent",
			inject: func(ops *stageOps) {
				ops.buildInfo = func(string) (schema.BuildInfo, error) { return schema.BuildInfo{}, nil }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.BuildInfo != nil {
					t.Fatalf("BuildInfo = %#v, want nil", analysis.BuildInfo)
				}
			},
		},
		{
			name:    "runtime",
			stage:   schema.AnalysisStageRuntime,
			code:    string(recoveryerr.CodeRuntimeMetadataNotFound),
			message: "runtime metadata evidence is absent",
			inject: func(ops *stageOps) {
				ops.runtime = func(string) (schema.RuntimeMetadata, error) { return schema.RuntimeMetadata{}, nil }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.Runtime != nil {
					t.Fatalf("Runtime = %#v, want nil", analysis.Runtime)
				}
			},
		},
		{
			name:    "runtime explicit absent",
			stage:   schema.AnalysisStageRuntime,
			code:    string(recoveryerr.CodeRuntimeMetadataNotFound),
			message: "runtime metadata evidence is absent",
			inject: func(ops *stageOps) {
				ops.runtime = func(string) (schema.RuntimeMetadata, error) {
					return schema.RuntimeMetadata{TrustSummary: schema.RuntimeTrustSummaryAbsent}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.Runtime != nil {
					t.Fatalf("Runtime = %#v, want nil", analysis.Runtime)
				}
			},
		},
		{
			name:    "functions",
			stage:   schema.AnalysisStageFunctions,
			code:    string(recoveryerr.CodePclntabNotFound),
			message: "function evidence is absent",
			inject: func(ops *stageOps) {
				ops.functions = func(string) ([]schema.Function, error) { return nil, nil }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if len(analysis.Functions) != 0 || analysis.Packages != nil || analysis.Peeling != nil {
					t.Fatalf("empty functions derived payload: %#v", analysis)
				}
				for _, stage := range []schema.AnalysisStage{
					schema.AnalysisStagePackages,
					schema.AnalysisStagePeeling,
				} {
					if _, exists := diagnosticForStage(analysis.Diagnostics, stage); exists {
						t.Fatalf("diagnostics = %#v, want no derived %q claim", analysis.Diagnostics, stage)
					}
				}
				if analysis.Refined == nil || len(analysis.Refined.Functions) != 0 ||
					len(analysis.Refined.Packages) != 0 || len(analysis.Refined.Types) != 1 ||
					len(analysis.Refined.Strings) != 1 {
					t.Fatalf("Refined = %#v, want only available type/string families", analysis.Refined)
				}
			},
		},
		{
			name:    "packages",
			stage:   schema.AnalysisStagePackages,
			code:    "packages_not_found",
			message: "package evidence is absent",
			inject: func(ops *stageOps) {
				ops.functions = func(string) ([]schema.Function, error) {
					return []schema.Function{{Name: "type:.eq.fixture", Entry: 0x1000, End: 0x1100}}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.Packages != nil {
					t.Fatalf("Packages = %#v, want nil", analysis.Packages)
				}
				if analysis.Refined == nil || len(analysis.Refined.Packages) != 0 ||
					len(analysis.Refined.Functions) != 1 || len(analysis.Refined.Types) != 1 ||
					len(analysis.Refined.Strings) != 1 {
					t.Fatalf("Refined = %#v, want only available function/type/string families", analysis.Refined)
				}
			},
		},
		{
			name:    "types",
			stage:   schema.AnalysisStageTypes,
			code:    string(recoveryerr.CodeDWARFNotFound),
			message: "type evidence is absent",
			inject: func(ops *stageOps) {
				ops.types = func(string) ([]schema.Type, error) { return nil, nil }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if len(analysis.Types) != 0 {
					t.Fatalf("Types = %#v, want empty", analysis.Types)
				}
			},
		},
		{
			name:    "strings",
			stage:   schema.AnalysisStageStrings,
			code:    string(recoveryerr.CodeStringRegionsNotFound),
			message: "string evidence is absent",
			inject: func(ops *stageOps) {
				ops.strings = func(string) (recoverystrings.Result, error) { return recoverystrings.Result{}, nil }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.StringRegions != nil || analysis.Strings != nil {
					t.Fatalf("empty strings published payload: regions=%#v strings=%#v", analysis.StringRegions, analysis.Strings)
				}
			},
		},
		{
			name:    "string regions without candidates",
			stage:   schema.AnalysisStageStrings,
			code:    "string_candidates_not_found",
			message: "string candidate evidence is absent",
			inject: func(ops *stageOps) {
				ops.strings = func(string) (recoverystrings.Result, error) {
					return recoverystrings.Result{
						Regions: []schema.StringRegion{{Name: ".rodata", Addr: 0x2000, Size: 8}},
					}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if len(analysis.StringRegions) != 1 || analysis.Strings != nil {
					t.Fatalf(
						"region-only strings payload: regions=%#v strings=%#v",
						analysis.StringRegions,
						analysis.Strings,
					)
				}
			},
		},
		{
			name:    "source tree",
			stage:   schema.AnalysisStageSourceTree,
			code:    string(recoveryerr.CodeSourceTreeNotFound),
			message: "source-tree evidence is absent",
			inject: func(ops *stageOps) {
				ops.sourceTree = func(string, schema.Analysis) (schema.SourceTree, error) {
					return schema.SourceTree{}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.SourceTree != nil {
					t.Fatalf("SourceTree = %#v, want nil", analysis.SourceTree)
				}
			},
		},
		{
			name:    "source tree root only",
			stage:   schema.AnalysisStageSourceTree,
			code:    string(recoveryerr.CodeSourceTreeNotFound),
			message: "source-tree evidence is absent",
			inject: func(ops *stageOps) {
				ops.sourceTree = func(string, schema.Analysis) (schema.SourceTree, error) {
					return schema.SourceTree{Root: "example.com/fixture"}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.SourceTree != nil {
					t.Fatalf("SourceTree = %#v, want nil", analysis.SourceTree)
				}
			},
		},
		{
			name:    "source tree metadata only",
			stage:   schema.AnalysisStageSourceTree,
			code:    string(recoveryerr.CodeSourceTreeNotFound),
			message: "source-tree evidence is absent",
			inject: func(ops *stageOps) {
				ops.sourceTree = func(string, schema.Analysis) (schema.SourceTree, error) {
					return schema.SourceTree{
						Root:               "example.com/fixture",
						SourceEvidenceKind: schema.SourceEvidenceKindPackageFallback,
						SourceEvidenceSummary: schema.SourceEvidenceSummary{
							TreeKind: schema.SourceEvidenceKindPackageFallback,
						},
						PathlessFileEvidence: true,
					}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.SourceTree != nil {
					t.Fatalf("SourceTree = %#v, want nil", analysis.SourceTree)
				}
			},
		},
		{
			name:    "peeling",
			stage:   schema.AnalysisStagePeeling,
			code:    "peeling_unavailable",
			message: "peeling evidence is absent",
			inject: func(ops *stageOps) {
				ops.peeling = func(schema.Analysis) *schema.PeelingAnalysis { return &schema.PeelingAnalysis{} }
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.Peeling != nil {
					t.Fatalf("Peeling = %#v, want nil", analysis.Peeling)
				}
			},
		},
		{
			name:    "refinement",
			stage:   schema.AnalysisStageRefinement,
			code:    "refinement_unavailable",
			message: "refinement evidence is absent",
			inject: func(ops *stageOps) {
				ops.refine = func(context.Context, schema.Analysis) (schema.RefinedAnalysis, error) {
					return schema.RefinedAnalysis{}, nil
				}
			},
			assert: func(t *testing.T, analysis schema.Analysis) {
				t.Helper()
				if analysis.Refined != nil {
					t.Fatalf("Refined = %#v, want nil", analysis.Refined)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ops := successfulStageOps()
			tt.inject(&ops)
			got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
			if err != nil {
				t.Fatalf("AnalyzeFile() error = %v", err)
			}

			diagnostic, exists := diagnosticForStage(got.Diagnostics, tt.stage)
			if !exists {
				t.Fatalf("diagnostics = %#v, want %q", got.Diagnostics, tt.stage)
			}
			if diagnostic.Status != schema.StageStatusUnavailable || diagnostic.Code != tt.code || diagnostic.Message != tt.message {
				t.Fatalf("diagnostic = %#v, want unavailable/%s/%s", diagnostic, tt.code, tt.message)
			}
			tt.assert(t, got)
			assertOrderedUniqueDiagnostics(t, got.Diagnostics)
		})
	}
}

func TestAnalyzeFileAcceptsSourceTreeNodeEvidence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tree schema.SourceTree
	}{
		{
			name: "module root with fileless package node",
			tree: schema.SourceTree{
				Root:               "example.com/fixture",
				SourceEvidenceKind: schema.SourceEvidenceKindPackageFallback,
				Packages: []schema.SourcePackage{{
					Name:            "main",
					ImportPath:      "example.com/fixture",
					HasFileEvidence: false,
					Files:           []string{},
				}},
			},
		},
		{
			name: "file-backed tree",
			tree: schema.SourceTree{
				Root:  "example.com/fixture",
				Files: []string{"main.go"},
			},
		},
		{
			name: "external-only node tree",
			tree: schema.SourceTree{
				Root: "example.com/fixture",
				ExternalPackages: []schema.SourcePackage{{
					Name:       "fmt",
					ImportPath: "fmt",
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ops := successfulStageOps()
			ops.sourceTree = func(string, schema.Analysis) (schema.SourceTree, error) { return tt.tree, nil }
			got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
			if err != nil {
				t.Fatalf("AnalyzeFile() error = %v", err)
			}

			diagnostic, exists := diagnosticForStage(got.Diagnostics, schema.AnalysisStageSourceTree)
			if !exists || diagnostic.Status != schema.StageStatusAvailable {
				t.Fatalf("source-tree diagnostic = %#v, exists=%v, want available", diagnostic, exists)
			}
			if got.SourceTree == nil {
				t.Fatal("SourceTree = nil, want node-backed evidence")
			}
			assertOrderedUniqueDiagnostics(t, got.Diagnostics)
		})
	}
}

func TestRecoverSourceTreePreservesExternalDWARFCandidate(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Input:     schema.Input{Format: "elf"},
		BuildInfo: &schema.BuildInfo{Path: "example.com/fixture"},
	}
	dwarfTree := schema.SourceTree{
		Root: "example.com/fixture",
		ExternalPackages: []schema.SourcePackage{{
			Name:       "fmt",
			ImportPath: "fmt",
		}},
	}
	functionTree := schema.SourceTree{
		Root:  "example.com/fixture",
		Files: []string{"main.go"},
	}

	tests := []struct {
		name          string
		buildFunction func(schema.Analysis) (schema.SourceTree, error)
		want          schema.SourceTree
	}{
		{
			name: "stronger module-local function tree wins",
			buildFunction: func(schema.Analysis) (schema.SourceTree, error) {
				return functionTree, nil
			},
			want: functionTree,
		},
		{
			name: "external DWARF survives unavailable function tree",
			buildFunction: func(schema.Analysis) (schema.SourceTree, error) {
				return schema.SourceTree{}, recoveryerr.NewUnavailable(
					recoveryerr.CodeSourceTreeNotFound,
					"function source evidence is absent",
					nil,
				)
			},
			want: dwarfTree,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fallbackCalled := false
			got, err := recoverSourceTreeWithOps("fixture.bin", analysis, sourceTreeRecoveryOps{
				readSourceFiles: func(string) ([]string, error) { return []string{"/usr/local/go/src/fmt/print.go"}, nil },
				buildDWARFTree: func(schema.Analysis, []string) (schema.SourceTree, error) {
					return dwarfTree, nil
				},
				buildFunctionTree: tt.buildFunction,
				buildFallbackTree: func(schema.Analysis) (schema.SourceTree, error) {
					fallbackCalled = true
					return schema.SourceTree{Files: []string{"fallback.go"}}, nil
				},
			})
			if err != nil {
				t.Fatalf("recoverSourceTreeWithOps() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("recoverSourceTreeWithOps() = %#v, want %#v", got, tt.want)
			}
			if fallbackCalled {
				t.Fatal("recoverSourceTreeWithOps() called package fallback despite a truthful candidate")
			}
		})
	}
}

func TestAnalyzeFileSourceTreeFallbackTaxonomy(t *testing.T) {
	t.Parallel()

	unexpectedRead := errors.New("fixture DWARF reader corruption")
	absent := func(message string) error {
		return recoveryerr.NewUnavailable(recoveryerr.CodeSourceTreeNotFound, message, nil)
	}
	tests := []struct {
		name          string
		readErr       error
		functionTree  schema.SourceTree
		functionErr   error
		fallbackTree  schema.SourceTree
		fallbackErr   error
		wantStatus    schema.StageStatus
		wantCode      string
		wantCauseText string
	}{
		{
			name:         "all proven absent is unavailable",
			readErr:      absent("DWARF source evidence is absent"),
			functionErr:  absent("function source evidence is absent"),
			wantStatus:   schema.StageStatusUnavailable,
			wantCode:     string(recoveryerr.CodeSourceTreeNotFound),
			fallbackTree: schema.SourceTree{},
		},
		{
			name:          "unexpected reader failure remains failed",
			readErr:       unexpectedRead,
			functionErr:   absent("function source evidence is absent"),
			wantStatus:    schema.StageStatusFailed,
			wantCode:      stageFailureCode,
			wantCauseText: unexpectedRead.Error(),
		},
		{
			name:         "truthful function fallback masks richer reader failure",
			readErr:      unexpectedRead,
			functionTree: schema.SourceTree{Files: []string{"main.go"}},
			wantStatus:   schema.StageStatusAvailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recoveryOps := sourceTreeRecoveryOps{
				readSourceFiles: func(string) ([]string, error) { return nil, tt.readErr },
				buildDWARFTree: func(schema.Analysis, []string) (schema.SourceTree, error) {
					return schema.SourceTree{}, errors.New("unexpected DWARF builder call")
				},
				buildFunctionTree: func(schema.Analysis) (schema.SourceTree, error) {
					return tt.functionTree, tt.functionErr
				},
				buildFallbackTree: func(schema.Analysis) (schema.SourceTree, error) {
					return tt.fallbackTree, tt.fallbackErr
				},
			}
			if tt.wantCauseText != "" {
				_, directErr := recoverSourceTreeWithOps("fixture.bin", schema.Analysis{
					Input:     schema.Input{Format: "elf"},
					BuildInfo: &schema.BuildInfo{Path: "example.com/fixture"},
				}, recoveryOps)
				if !errors.Is(directErr, tt.readErr) {
					t.Fatalf("recoverSourceTreeWithOps() error = %v, want retained cause %v", directErr, tt.readErr)
				}
			}

			ops := successfulStageOps()
			ops.sourceTree = func(path string, analysis schema.Analysis) (schema.SourceTree, error) {
				return recoverSourceTreeWithOps(path, analysis, recoveryOps)
			}

			got, err := newAnalyzerForTest(ops).AnalyzeFile(context.Background(), writeIngestibleELF(t))
			if err != nil {
				t.Fatalf("AnalyzeFile() error = %v", err)
			}
			diagnostic, exists := diagnosticForStage(got.Diagnostics, schema.AnalysisStageSourceTree)
			if !exists || diagnostic.Status != tt.wantStatus || diagnostic.Code != tt.wantCode {
				t.Fatalf("source-tree diagnostic = %#v, exists=%v, want %q/%q", diagnostic, exists, tt.wantStatus, tt.wantCode)
			}
			if tt.wantCauseText != "" && !strings.Contains(diagnostic.Message, tt.wantCauseText) {
				t.Fatalf("source-tree diagnostic = %#v, want retained cause %q", diagnostic, tt.wantCauseText)
			}
		})
	}
}

func successfulStageOps() stageOps {
	return stageOps{
		buildInfo: func(string) (schema.BuildInfo, error) {
			return schema.BuildInfo{Path: "example.com/fixture"}, nil
		},
		runtime: func(string) (schema.RuntimeMetadata, error) {
			return schema.RuntimeMetadata{TrustSummary: schema.RuntimeTrustSummarySectionHeuristic}, nil
		},
		functions: func(string) ([]schema.Function, error) {
			return []schema.Function{{Name: "main.main", Package: "main", Entry: 0x1000, End: 0x1100}}, nil
		},
		types: func(string) ([]schema.Type, error) {
			return []schema.Type{{Name: "main.fixture", Kind: "struct"}}, nil
		},
		strings: func(string) (recoverystrings.Result, error) {
			return recoverystrings.Result{
				Regions:    []schema.StringRegion{{Name: ".rodata", Addr: 0x2000, Size: 8}},
				Candidates: []schema.StringCandidate{{Value: "fixture", Region: ".rodata", Addr: 0x2000}},
			}, nil
		},
		peeling: func(analysis schema.Analysis) *schema.PeelingAnalysis {
			if len(analysis.Functions) == 0 {
				return nil
			}
			return &schema.PeelingAnalysis{
				Functions: []schema.PeelingFunction{{Name: analysis.Functions[0].Name}},
			}
		},
		sourceTree: func(string, schema.Analysis) (schema.SourceTree, error) {
			return schema.SourceTree{
				Root:               "example.com/fixture",
				SourceEvidenceKind: schema.SourceEvidenceKindPackageFallback,
				Files:              []string{},
				Packages: []schema.SourcePackage{{
					Name:            "main",
					ImportPath:      "example.com/fixture",
					HasFileEvidence: false,
					Files:           []string{},
				}},
			}, nil
		},
		refine: func(_ context.Context, analysis schema.Analysis) (schema.RefinedAnalysis, error) {
			refined := schema.RefinedAnalysis{}
			for _, function := range analysis.Functions {
				refined.Functions = append(refined.Functions, schema.RefinedFunction{Name: function.Name})
			}
			for _, pkg := range analysis.Packages {
				refined.Packages = append(refined.Packages, schema.RefinedPackage{Name: pkg.Name})
			}
			for _, typ := range analysis.Types {
				refined.Types = append(refined.Types, schema.RefinedType{Name: typ.Name})
			}
			for _, candidate := range analysis.Strings {
				refined.Strings = append(refined.Strings, schema.RefinedString{Value: candidate.Value})
			}
			return refined, nil
		},
	}
}

func writeIngestibleELF(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "fixture.bin")
	data := []byte{0x7f, 'E', 'L', 'F', 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return path
}

func diagnosticForStage(diagnostics []schema.StageDiagnostic, stage schema.AnalysisStage) (schema.StageDiagnostic, bool) {
	for _, diagnostic := range diagnostics {
		if diagnostic.Stage == stage {
			return diagnostic, true
		}
	}
	return schema.StageDiagnostic{}, false
}

func assertOrderedUniqueDiagnostics(t *testing.T, diagnostics []schema.StageDiagnostic) {
	t.Helper()

	order := map[schema.AnalysisStage]int{
		schema.AnalysisStageBuildInfo:  0,
		schema.AnalysisStageRuntime:    1,
		schema.AnalysisStageFunctions:  2,
		schema.AnalysisStagePackages:   3,
		schema.AnalysisStageTypes:      4,
		schema.AnalysisStageStrings:    5,
		schema.AnalysisStageSourceTree: 6,
		schema.AnalysisStagePeeling:    7,
		schema.AnalysisStageRefinement: 8,
	}
	seen := make(map[schema.AnalysisStage]struct{}, len(diagnostics))
	previous := -1
	for _, diagnostic := range diagnostics {
		if _, ok := seen[diagnostic.Stage]; ok {
			t.Fatalf("duplicate diagnostic for stage %q in %#v", diagnostic.Stage, diagnostics)
		}
		seen[diagnostic.Stage] = struct{}{}
		current, ok := order[diagnostic.Stage]
		if !ok {
			t.Fatalf("unknown diagnostic stage %q", diagnostic.Stage)
		}
		if current <= previous {
			t.Fatalf("diagnostics out of order: %#v", diagnostics)
		}
		previous = current
	}
}

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
