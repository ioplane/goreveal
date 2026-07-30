package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/ioplane/goreveal/core/buildinfo"
	"github.com/ioplane/goreveal/core/functions"
	"github.com/ioplane/goreveal/core/ingest"
	"github.com/ioplane/goreveal/core/packages"
	recoveryruntime "github.com/ioplane/goreveal/core/runtime"
	recoverystrings "github.com/ioplane/goreveal/core/strings"
	"github.com/ioplane/goreveal/core/types"
	"github.com/ioplane/goreveal/deobfuscation"
	"github.com/ioplane/goreveal/deobfuscation/garble"
	"github.com/ioplane/goreveal/deobfuscation/refine"
	"github.com/ioplane/goreveal/engine/peeling"
	"github.com/ioplane/goreveal/engine/projection"
	"github.com/ioplane/goreveal/schema"
)

// Analyzer orchestrates the minimal ingest-to-schema flow.
type Analyzer struct{}

// New creates a minimal Sprint 1 analyzer.
func New() Analyzer {
	return Analyzer{}
}

// AnalyzeFile ingests a binary and maps it into the canonical schema.
func (Analyzer) AnalyzeFile(ctx context.Context, path string) (schema.Analysis, error) {
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
		Types: []schema.Type{},
		Provenance: schema.Provenance{
			Source:     "core.ingest",
			Confidence: "high",
		},
	}

	if info, err := buildinfo.Read(path); err == nil {
		analysis.BuildInfo = &info
	}

	switch analysis.Input.Format {
	case "elf":
		analysis = analyzeELFSurfaces(path, analysis)
	case "pe":
		analysis = analyzePESurfaces(path, analysis)
	case "macho":
		analysis = analyzeMachOSurfaces(path, analysis)
	}

	analysis.Peeling = peeling.Build(analysis)

	if refined, err := deobfuscation.NewPipeline(refine.Pass{}, garble.Pass{}).Run(ctx, analysis); err == nil {
		if hasRefinedContent(refined) {
			analysis.Refined = &refined
		}
	}

	return analysis, nil
}

func analyzeELFSurfaces(path string, analysis schema.Analysis) schema.Analysis {
	if recovered, err := recoveryruntime.ReadMetadata(path); err == nil {
		analysis.Runtime = &recovered
	}
	if recovered, err := functions.Recover(path); err == nil {
		analysis.Functions = functions.EnrichBuildInfoMetadata(recovered, analysis.BuildInfo)
		analysis.Packages = packages.Recover(recovered)
		analysis.Packages = packages.EnrichBuildInfoMetadata(analysis.Packages, analysis.BuildInfo)
	}
	if recovered, err := types.Recover(path); err == nil {
		analysis.Types = recovered
		analysis.Types = types.EnrichBuildInfoMetadata(analysis.Types, analysis.BuildInfo)
	}
	if recovered, err := recoverystrings.Recover(path); err == nil {
		analysis.StringRegions = recovered.Regions
		analysis.Strings = recovered.Candidates
	}
	if files, err := projection.ReadSourceFiles(path); err == nil {
		if tree, buildErr := projection.BuildSourceTree(analysis, files); buildErr == nil {
			if fallbackTree, fallbackErr := projection.BuildFunctionSourceTree(analysis); fallbackErr == nil &&
				shouldPreferFunctionSourceTree(tree, fallbackTree) {
				tree = fallbackTree
			}
			analysis.SourceTree = &tree
			analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
			analysis.Types = types.EnrichUserMetadata(analysis.Types, analysis.SourceTree)
		}
	} else if tree, buildErr := projection.BuildFunctionSourceTree(analysis); buildErr == nil {
		analysis.SourceTree = &tree
		analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
	} else if tree, fallbackErr := projection.BuildFallbackSourceTree(analysis); fallbackErr == nil {
		analysis.SourceTree = &tree
	}

	annotateELFFunctionRecoveryBlocker(&analysis)
	annotateELFFunctionFoothold(&analysis)

	return analysis
}

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

func analyzePESurfaces(path string, analysis schema.Analysis) schema.Analysis {
	if recovered, err := recoveryruntime.ReadMetadata(path); err == nil {
		analysis.Runtime = &recovered
	}
	if recovered, err := functions.Recover(path); err == nil {
		analysis.Functions = functions.EnrichBuildInfoMetadata(recovered, analysis.BuildInfo)
		analysis.Packages = packages.Recover(recovered)
		analysis.Packages = packages.EnrichBuildInfoMetadata(analysis.Packages, analysis.BuildInfo)
		if tree, buildErr := projection.BuildFunctionSourceTree(analysis); buildErr == nil {
			analysis.SourceTree = &tree
			analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
		}
	}

	return analysis
}

func analyzeMachOSurfaces(path string, analysis schema.Analysis) schema.Analysis {
	if recovered, err := functions.Recover(path); err == nil {
		analysis.Functions = functions.EnrichBuildInfoMetadata(recovered, analysis.BuildInfo)
		analysis.Packages = packages.Recover(recovered)
		analysis.Packages = packages.EnrichBuildInfoMetadata(analysis.Packages, analysis.BuildInfo)
		if tree, buildErr := projection.BuildFunctionSourceTree(analysis); buildErr == nil {
			analysis.SourceTree = &tree
			analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
		}
	}

	return analysis
}

func hasRefinedContent(refined schema.RefinedAnalysis) bool {
	return len(refined.Functions) > 0 ||
		len(refined.Packages) > 0 ||
		len(refined.Types) > 0 ||
		len(refined.Strings) > 0
}
