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

type sourceTreeRecoveryOps struct {
	readSourceFiles   func(string) ([]string, error)
	buildDWARFTree    func(schema.Analysis, []string) (schema.SourceTree, error)
	buildFunctionTree func(schema.Analysis) (schema.SourceTree, error)
	buildFallbackTree func(schema.Analysis) (schema.SourceTree, error)
}

type sourceTreeCandidateKind uint8

const (
	sourceTreeCandidateEmpty sourceTreeCandidateKind = iota
	sourceTreeCandidateExternalOnly
	sourceTreeCandidateModuleFiles
)

func (ops stageOps) isZero() bool {
	return ops.buildInfo == nil &&
		ops.runtime == nil &&
		ops.functions == nil &&
		ops.types == nil &&
		ops.strings == nil &&
		ops.sourceTree == nil &&
		ops.peeling == nil &&
		ops.refine == nil
}

func (ops stageOps) missing() []string {
	missing := make([]string, 0, 8)
	if ops.buildInfo == nil {
		missing = append(missing, "build_info")
	}
	if ops.runtime == nil {
		missing = append(missing, "runtime")
	}
	if ops.functions == nil {
		missing = append(missing, "functions")
	}
	if ops.types == nil {
		missing = append(missing, "types")
	}
	if ops.strings == nil {
		missing = append(missing, "strings")
	}
	if ops.sourceTree == nil {
		missing = append(missing, "source_tree")
	}
	if ops.peeling == nil {
		missing = append(missing, "peeling")
	}
	if ops.refine == nil {
		missing = append(missing, "refinement")
	}
	return missing
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
	if sourceTreeHasEvidence(tree) {
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

func anyStageAvailable(analysis schema.Analysis, stages ...schema.AnalysisStage) bool {
	for _, stage := range stages {
		if stageAvailable(analysis, stage) {
			return true
		}
	}

	return false
}

func stageAvailable(analysis schema.Analysis, stage schema.AnalysisStage) bool {
	for _, diagnostic := range analysis.Diagnostics {
		if diagnostic.Stage == stage {
			return diagnostic.Status == schema.StageStatusAvailable
		}
	}

	return false
}

func recoverSourceTree(path string, analysis schema.Analysis) (schema.SourceTree, error) {
	return recoverSourceTreeWithOps(path, analysis, sourceTreeRecoveryOps{
		readSourceFiles:   projection.ReadSourceFiles,
		buildDWARFTree:    projection.BuildSourceTree,
		buildFunctionTree: projection.BuildFunctionSourceTree,
		buildFallbackTree: projection.BuildFallbackSourceTree,
	})
}

func recoverSourceTreeWithOps(
	path string,
	analysis schema.Analysis,
	ops sourceTreeRecoveryOps,
) (schema.SourceTree, error) {
	if analysis.BuildInfo == nil || analysis.BuildInfo.Path == "" {
		return schema.SourceTree{}, recoveryerr.NewUnavailable(
			recoveryerr.CodeSourceTreeNotFound,
			"source-tree build path evidence is absent",
			nil,
		)
	}

	attemptErrors := make([]error, 0, 3)
	unexpectedErrors := make([]error, 0, 3)
	var candidate *schema.SourceTree
	if analysis.Input.Format == "elf" {
		tree, err := recoverDWARFSourceTreeWithOps(path, analysis, ops)
		if err != nil {
			recordSourceTreeFailure(&attemptErrors, &unexpectedErrors, "recover DWARF source tree", err)
		} else {
			switch sourceTreeCandidate(tree) {
			case sourceTreeCandidateModuleFiles:
				return tree, nil
			case sourceTreeCandidateExternalOnly:
				candidate = &tree
			case sourceTreeCandidateEmpty:
				recordSourceTreeAbsence(&attemptErrors, "DWARF source tree contains no nodes")
			}
		}
	}

	functionTree, functionErr := ops.buildFunctionTree(analysis)
	if functionErr != nil {
		recordSourceTreeFailure(&attemptErrors, &unexpectedErrors, "build function source tree", functionErr)
	} else {
		switch sourceTreeCandidate(functionTree) {
		case sourceTreeCandidateModuleFiles:
			if len(unexpectedErrors) == 0 {
				return functionTree, nil
			}
			candidate = &functionTree
		case sourceTreeCandidateExternalOnly:
			if candidate == nil {
				candidate = &functionTree
			}
		case sourceTreeCandidateEmpty:
			recordSourceTreeAbsence(&attemptErrors, "function source tree contains no nodes")
		}
	}
	if candidate != nil {
		if len(unexpectedErrors) == 0 {
			return *candidate, nil
		}
		return schema.SourceTree{}, sourceTreeRecoveryError(attemptErrors, unexpectedErrors)
	}
	return recoverPackageSourceTreeWithOps(analysis, ops, attemptErrors, unexpectedErrors)
}

func recoverPackageSourceTreeWithOps(
	analysis schema.Analysis,
	ops sourceTreeRecoveryOps,
	attemptErrors, unexpectedErrors []error,
) (schema.SourceTree, error) {
	fallbackTree, fallbackErr := ops.buildFallbackTree(analysis)
	if fallbackErr == nil && sourceTreeHasEvidence(fallbackTree) {
		if len(unexpectedErrors) == 0 {
			return fallbackTree, nil
		}
		return schema.SourceTree{}, sourceTreeRecoveryError(attemptErrors, unexpectedErrors)
	}
	if fallbackErr != nil {
		recordSourceTreeFailure(&attemptErrors, &unexpectedErrors, "build package fallback source tree", fallbackErr)
	} else {
		recordSourceTreeAbsence(&attemptErrors, "package fallback source tree contains no nodes")
	}

	return schema.SourceTree{}, sourceTreeRecoveryError(attemptErrors, unexpectedErrors)
}

func sourceTreeRecoveryError(attemptErrors, unexpectedErrors []error) error {
	if len(unexpectedErrors) != 0 {
		return fmt.Errorf("recover source tree: %w", errors.Join(unexpectedErrors...))
	}

	return recoveryerr.NewUnavailable(
		recoveryerr.CodeSourceTreeNotFound,
		"source-tree evidence is absent",
		errors.Join(attemptErrors...),
	)
}

func recoverDWARFSourceTreeWithOps(
	path string,
	analysis schema.Analysis,
	ops sourceTreeRecoveryOps,
) (schema.SourceTree, error) {
	files, err := ops.readSourceFiles(path)
	if err != nil {
		return schema.SourceTree{}, fmt.Errorf("read DWARF source files: %w", err)
	}
	return ops.buildDWARFTree(analysis, files)
}

func sourceTreeCandidate(tree schema.SourceTree) sourceTreeCandidateKind {
	if len(tree.Files) != 0 {
		return sourceTreeCandidateModuleFiles
	}
	if sourceTreeHasEvidence(tree) {
		return sourceTreeCandidateExternalOnly
	}
	return sourceTreeCandidateEmpty
}

func recordSourceTreeFailure(attempts, unexpected *[]error, operation string, err error) {
	wrapped := fmt.Errorf("%s: %w", operation, err)
	*attempts = append(*attempts, wrapped)
	if !errors.Is(err, recoveryerr.ErrUnavailable) {
		*unexpected = append(*unexpected, wrapped)
	}
}

func recordSourceTreeAbsence(attempts *[]error, message string) {
	*attempts = append(*attempts, recoveryerr.NewUnavailable(
		recoveryerr.CodeSourceTreeNotFound,
		message,
		nil,
	))
}

func sourceTreeHasEvidence(tree schema.SourceTree) bool {
	return len(tree.Files) != 0 || len(tree.Packages) != 0 || len(tree.ExternalPackages) != 0
}
