package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/dantte-lp/goreveal/core/buildinfo"
	"github.com/dantte-lp/goreveal/core/functions"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
	recoveryruntime "github.com/dantte-lp/goreveal/core/runtime"
	recoverystrings "github.com/dantte-lp/goreveal/core/strings"
	"github.com/dantte-lp/goreveal/core/types"
	"github.com/dantte-lp/goreveal/deobfuscation"
	"github.com/dantte-lp/goreveal/deobfuscation/garble"
	"github.com/dantte-lp/goreveal/deobfuscation/refine"
	"github.com/dantte-lp/goreveal/engine/peeling"
	"github.com/dantte-lp/goreveal/engine/projection"
	"github.com/dantte-lp/goreveal/schema"
)

const stageFailureCode = "stage_failed"

const (
	stageCodePackagesNotFound      = "packages_not_found"
	stageCodePeelingUnavailable    = "peeling_unavailable"
	stageCodeRefinementUnavailable = "refinement_unavailable"
)

type stageOps struct {
	buildInfo  func(string) (schema.BuildInfo, error)
	runtime    func(string) (schema.RuntimeMetadata, error)
	functions  func(string) ([]schema.Function, error)
	types      func(string) ([]schema.Type, error)
	strings    func(string) (recoverystrings.Result, error)
	sourceTree func(string, schema.Analysis) (schema.SourceTree, error)
	peeling    func(schema.Analysis) *schema.PeelingAnalysis
	refine     func(context.Context, schema.Analysis) (schema.RefinedAnalysis, error)
}

func productionStageOps() stageOps {
	return stageOps{
		buildInfo:  buildinfo.Read,
		runtime:    recoveryruntime.ReadMetadata,
		functions:  functions.Recover,
		types:      types.Recover,
		strings:    recoverystrings.Recover,
		sourceTree: recoverSourceTree,
		peeling:    peeling.Build,
		refine: func(ctx context.Context, analysis schema.Analysis) (schema.RefinedAnalysis, error) {
			return deobfuscation.NewPipeline(refine.Pass{}, garble.Pass{}).Run(ctx, analysis)
		},
	}
}

func newAnalyzerForTest(ops stageOps) Analyzer {
	return Analyzer{ops: ops}
}

func executeStage[T any](
	analysis *schema.Analysis,
	stage schema.AnalysisStage,
	op func() (T, error),
	validate func(T) (code, message string),
) (T, bool) {
	value, err := op()
	if err != nil {
		analysis.Diagnostics = append(analysis.Diagnostics, diagnosticFromError(stage, err))
		var zero T
		return zero, false
	}
	if code, message := validate(value); code != "" {
		analysis.Diagnostics = append(analysis.Diagnostics, schema.StageDiagnostic{
			Stage:   stage,
			Status:  schema.StageStatusUnavailable,
			Code:    code,
			Message: message,
		})
		return value, false
	}

	analysis.Diagnostics = append(analysis.Diagnostics, schema.StageDiagnostic{
		Stage:  stage,
		Status: schema.StageStatusAvailable,
	})
	return value, true
}

func buildInfoEvidence(info schema.BuildInfo) (string, string) {
	if info.GoVersion != "" || info.Path != "" {
		return "", ""
	}
	return string(recoveryerr.CodeBuildInfoNotFound), "Go build info evidence is absent"
}

func runtimeEvidence(metadata schema.RuntimeMetadata) (string, string) {
	if metadata.TrustSummary != "" && metadata.TrustSummary != schema.RuntimeTrustSummaryAbsent {
		return "", ""
	}
	return string(recoveryerr.CodeRuntimeMetadataNotFound), "runtime metadata evidence is absent"
}

func functionEvidence(recovered []schema.Function) (string, string) {
	if len(recovered) != 0 {
		return "", ""
	}
	return string(recoveryerr.CodePclntabNotFound), "function evidence is absent"
}

func typeEvidence(recovered []schema.Type) (string, string) {
	if len(recovered) != 0 {
		return "", ""
	}
	return string(recoveryerr.CodeDWARFNotFound), "type evidence is absent"
}

func stringEvidence(recovered recoverystrings.Result) (string, string) {
	if len(recovered.Candidates) != 0 {
		return "", ""
	}
	if len(recovered.Regions) == 0 {
		return string(recoveryerr.CodeStringRegionsNotFound), "string evidence is absent"
	}
	return string(recoveryerr.CodeStringCandidatesNotFound), "string candidate evidence is absent"
}

func sourceTreeEvidence(tree schema.SourceTree) (string, string) {
	if tree.Root != "" || tree.SourceEvidenceKind != "" || len(tree.Files) != 0 ||
		len(tree.Packages) != 0 || len(tree.ExternalPackages) != 0 || tree.PathlessFileEvidence ||
		tree.SourceEvidenceSummary != (schema.SourceEvidenceSummary{}) {
		return "", ""
	}
	return string(recoveryerr.CodeSourceTreeNotFound), "source-tree evidence is absent"
}

func refinementEvidence(refined schema.RefinedAnalysis) (string, string) {
	if hasRefinedContent(refined) {
		return "", ""
	}
	return stageCodeRefinementUnavailable, "refinement evidence is absent"
}

func diagnosticFromError(stage schema.AnalysisStage, err error) schema.StageDiagnostic {
	diagnostic := schema.StageDiagnostic{
		Stage:   stage,
		Status:  schema.StageStatusFailed,
		Code:    stageFailureCode,
		Message: err.Error(),
	}

	var recoveryError *recoveryerr.Error
	if !errors.As(err, &recoveryError) {
		return diagnostic
	}

	diagnostic.Code = string(recoveryError.Code)
	diagnostic.Message = recoveryError.Message
	switch recoveryError.Kind {
	case recoveryerr.KindUnavailable:
		diagnostic.Status = schema.StageStatusUnavailable
	case recoveryerr.KindUnsupported:
		diagnostic.Status = schema.StageStatusUnsupported
	}

	return diagnostic
}

func appendDerivedDiagnostic(analysis *schema.Analysis, stage schema.AnalysisStage, status schema.StageStatus, code, message string) {
	analysis.Diagnostics = append(analysis.Diagnostics, schema.StageDiagnostic{
		Stage:   stage,
		Status:  status,
		Code:    code,
		Message: message,
	})
}

func allStagesAvailable(analysis schema.Analysis, stages ...schema.AnalysisStage) bool {
	for _, stage := range stages {
		available := false
		for _, diagnostic := range analysis.Diagnostics {
			if diagnostic.Stage == stage {
				available = diagnostic.Status == schema.StageStatusAvailable
				break
			}
		}
		if !available {
			return false
		}
	}

	return true
}

func recoverSourceTree(path string, analysis schema.Analysis) (schema.SourceTree, error) {
	if analysis.BuildInfo == nil || analysis.BuildInfo.Path == "" {
		return schema.SourceTree{}, recoveryerr.NewUnavailable(
			recoveryerr.CodeSourceTreeNotFound,
			"source-tree build path evidence is absent",
			nil,
		)
	}

	attemptErrors := make([]error, 0, 3)
	if analysis.Input.Format == "elf" {
		files, err := projection.ReadSourceFiles(path)
		if err == nil {
			tree, buildErr := projection.BuildSourceTree(analysis, files)
			if buildErr == nil {
				if functionTree, functionErr := projection.BuildFunctionSourceTree(analysis); functionErr == nil &&
					shouldPreferFunctionSourceTree(tree, functionTree) {
					tree = functionTree
				}
				return tree, nil
			}
			attemptErrors = append(attemptErrors, fmt.Errorf("build DWARF source tree: %w", buildErr))
		} else {
			attemptErrors = append(attemptErrors, fmt.Errorf("read DWARF source files: %w", err))
		}
	}

	functionTree, functionErr := projection.BuildFunctionSourceTree(analysis)
	if functionErr == nil {
		return functionTree, nil
	}
	attemptErrors = append(attemptErrors, fmt.Errorf("build function source tree: %w", functionErr))

	fallbackTree, fallbackErr := projection.BuildFallbackSourceTree(analysis)
	if fallbackErr == nil {
		return fallbackTree, nil
	}
	attemptErrors = append(attemptErrors, fmt.Errorf("build package fallback source tree: %w", fallbackErr))

	return schema.SourceTree{}, fmt.Errorf("recover source tree: %w", errors.Join(attemptErrors...))
}
