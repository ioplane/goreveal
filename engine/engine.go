package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/dantte-lp/goreveal/core/functions"
	"github.com/dantte-lp/goreveal/core/ingest"
	"github.com/dantte-lp/goreveal/core/packages"
	recoveryruntime "github.com/dantte-lp/goreveal/core/runtime"
	recoverystrings "github.com/dantte-lp/goreveal/core/strings"
	"github.com/dantte-lp/goreveal/core/types"
	"github.com/dantte-lp/goreveal/engine/peeling"
	"github.com/dantte-lp/goreveal/schema"
)

// Analyzer orchestrates the minimal ingest-to-schema flow.
type Analyzer struct {
	ops stageOps
}

// New creates a minimal Sprint 1 analyzer.
func New() Analyzer {
	return Analyzer{ops: productionStageOps()}
}

// AnalyzeFile ingests a binary and maps it into the canonical schema.
func (a Analyzer) AnalyzeFile(ctx context.Context, path string) (schema.Analysis, error) {
	if ctx == nil {
		return schema.Analysis{}, errors.New("analyze file: nil context")
	}

	file, err := ingest.Open(path)
	if err != nil {
		return schema.Analysis{}, fmt.Errorf("ingest %q: %w", path, err)
	}

	analysis := schema.Analysis{
		Input: schema.Input{
			Path:   file.Path,
			Size:   file.Size,
			Format: string(file.Format),
		},
		Types:       []schema.Type{},
		Diagnostics: make([]schema.StageDiagnostic, 0, 9),
		Provenance: schema.Provenance{
			Source:     "core.ingest",
			Confidence: "high",
		},
	}

	buildInfoAvailable := a.recoverBuildInfo(path, &analysis)
	a.recoverRuntime(path, &analysis)
	functionsAvailable := a.recoverFunctions(path, &analysis)
	a.recoverTypes(path, &analysis)
	a.recoverStrings(path, &analysis)
	a.recoverSourceTree(path, &analysis, buildInfoAvailable)
	annotateELFRecovery(&analysis)
	derivePeeling(&analysis, functionsAvailable)
	a.refine(ctx, &analysis)

	return analysis, nil
}

func (a Analyzer) recoverBuildInfo(path string, analysis *schema.Analysis) bool {
	info, available := executeStage(analysis, schema.AnalysisStageBuildInfo, func() (schema.BuildInfo, error) {
		return a.ops.buildInfo(path)
	})
	if available {
		analysis.BuildInfo = &info
	}
	return available
}

func (a Analyzer) recoverRuntime(path string, analysis *schema.Analysis) {
	metadata, available := executeStage(analysis, schema.AnalysisStageRuntime, func() (schema.RuntimeMetadata, error) {
		return a.ops.runtime(path)
	})
	if available {
		analysis.Runtime = &metadata
	}
}

func (a Analyzer) recoverFunctions(path string, analysis *schema.Analysis) bool {
	recovered, available := executeStage(analysis, schema.AnalysisStageFunctions, func() ([]schema.Function, error) {
		return a.ops.functions(path)
	})
	if !available {
		return false
	}

	analysis.Functions = functions.EnrichBuildInfoMetadata(recovered, analysis.BuildInfo)
	analysis.Packages = packages.EnrichBuildInfoMetadata(packages.Recover(recovered), analysis.BuildInfo)
	appendDerivedDiagnostic(analysis, schema.AnalysisStagePackages, schema.StageStatusAvailable, "", "")
	return true
}

func (a Analyzer) recoverTypes(path string, analysis *schema.Analysis) {
	recovered, available := executeStage(analysis, schema.AnalysisStageTypes, func() ([]schema.Type, error) {
		return a.ops.types(path)
	})
	if available {
		analysis.Types = types.EnrichBuildInfoMetadata(recovered, analysis.BuildInfo)
	}
}

func (a Analyzer) recoverStrings(path string, analysis *schema.Analysis) {
	recovered, available := executeStage(analysis, schema.AnalysisStageStrings, func() (recoverystrings.Result, error) {
		return a.ops.strings(path)
	})
	if available {
		analysis.StringRegions = recovered.Regions
		analysis.Strings = recovered.Candidates
	}
}

func (a Analyzer) recoverSourceTree(path string, analysis *schema.Analysis, buildInfoAvailable bool) {
	if !buildInfoAvailable {
		return
	}

	tree, sourceTreeAvailable := executeStage(analysis, schema.AnalysisStageSourceTree, func() (schema.SourceTree, error) {
		return a.ops.sourceTree(path, *analysis)
	})
	if !sourceTreeAvailable {
		return
	}

	analysis.SourceTree = &tree
	analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
	analysis.Types = types.EnrichUserMetadata(analysis.Types, analysis.SourceTree)
}

func annotateELFRecovery(analysis *schema.Analysis) {
	if analysis.Input.Format != "elf" {
		return
	}
	annotateELFFunctionRecoveryBlocker(analysis)
	annotateELFFunctionFoothold(analysis)
}

func derivePeeling(analysis *schema.Analysis, functionsAvailable bool) {
	if !functionsAvailable {
		return
	}

	analysis.Peeling = peeling.Build(*analysis)
	if analysis.Peeling == nil {
		appendDerivedDiagnostic(
			analysis,
			schema.AnalysisStagePeeling,
			schema.StageStatusUnavailable,
			"peeling_unavailable",
			"no function evidence is available for peeling",
		)
		return
	}
	appendDerivedDiagnostic(analysis, schema.AnalysisStagePeeling, schema.StageStatusAvailable, "", "")
}

func (a Analyzer) refine(ctx context.Context, analysis *schema.Analysis) {
	if !allStagesAvailable(*analysis, schema.AnalysisStageFunctions, schema.AnalysisStagePackages) {
		return
	}

	refined, available := executeStage(analysis, schema.AnalysisStageRefinement, func() (schema.RefinedAnalysis, error) {
		return a.ops.refine(ctx, *analysis)
	})
	if available && hasRefinedContent(refined) {
		analysis.Refined = &refined
	}
}

func hasRefinedContent(refined schema.RefinedAnalysis) bool {
	return len(refined.Functions) > 0 ||
		len(refined.Packages) > 0 ||
		len(refined.Types) > 0 ||
		len(refined.Strings) > 0
}

/*
	The functions below annotate bounded ELF evidence after all raw stage results
	have been recorded. They do not run recovery or create stage outcomes.
*/

func annotateELFFunctionRecoveryBlocker(analysis *schema.Analysis) {
	if analysis == nil || analysis.Runtime == nil {
		return
	}
	if analysis.Runtime.ELFPclntabHeaderMagicKind != "unknown" {
		return
	}
	if analysis.Runtime.ELFPclntabFunctionCountHint == 0 {
		return
	}
	if len(analysis.Functions) != 0 {
		return
	}

	analysis.Runtime.ELFFunctionRecoveryBlocker = "custom_pclntab_magic"
}

func annotateELFFunctionFoothold(analysis *schema.Analysis) {
	if analysis == nil || analysis.Runtime == nil {
		return
	}
	if analysis.Runtime.ELFPclntabHeaderMagicKind != "unknown" {
		return
	}
	if analysis.Runtime.ELFPclntabFunctionCountHint == 0 {
		return
	}
	if len(analysis.Functions) != 0 {
		return
	}
	if !analysis.Runtime.ELFFunctabPCOffsetsMonotonic {
		return
	}
	if !analysis.Runtime.ELFFunctabPCAddrHintsWithinText {
		return
	}
	if len(analysis.Runtime.ELFFunctabPCAddrSample) == 0 {
		return
	}
	if !analysis.Runtime.ELFFunctabPCAddrSampleAllWithinText {
		return
	}

	analysis.Runtime.ELFFunctionFoothold = "address_only"
	analysis.Runtime.ELFFunctionFootholdCountHint = analysis.Runtime.ELFPclntabFunctionCountHint
	analysis.Runtime.ELFFunctionFootholdTextSource = recoveryruntime.ELFTextSourceForProjection(analysis.Runtime)
	analysis.Runtime.ELFFunctionFootholdStartAddr = analysis.Runtime.ELFFunctabFirstPCAddrHint
	analysis.Runtime.ELFFunctionFootholdEndAddr = analysis.Runtime.ELFFunctabLastPCAddrHint
}

func shouldPreferFunctionSourceTree(dwarfTree, functionTree schema.SourceTree) bool {
	return len(dwarfTree.Files) == 0 && len(functionTree.Files) > 0
}
