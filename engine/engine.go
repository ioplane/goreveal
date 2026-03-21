package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/dantte-lp/goreveal/core/buildinfo"
	"github.com/dantte-lp/goreveal/core/functions"
	"github.com/dantte-lp/goreveal/core/ingest"
	"github.com/dantte-lp/goreveal/core/packages"
	recoveryruntime "github.com/dantte-lp/goreveal/core/runtime"
	recoverystrings "github.com/dantte-lp/goreveal/core/strings"
	"github.com/dantte-lp/goreveal/core/types"
	"github.com/dantte-lp/goreveal/deobfuscation"
	"github.com/dantte-lp/goreveal/deobfuscation/garble"
	"github.com/dantte-lp/goreveal/deobfuscation/refine"
	"github.com/dantte-lp/goreveal/engine/projection"
	"github.com/dantte-lp/goreveal/schema"
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

	if analysis.Input.Format == "elf" {
		analysis = analyzeELFSurfaces(path, analysis)
	}

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
			analysis.SourceTree = &tree
			analysis.Packages = packages.EnrichSourceMetadata(analysis.Packages, analysis.SourceTree)
			analysis.Types = types.EnrichUserMetadata(analysis.Types, analysis.SourceTree)
		}
	} else if tree, fallbackErr := projection.BuildFallbackSourceTree(analysis); fallbackErr == nil {
		analysis.SourceTree = &tree
	}

	return analysis
}

func hasRefinedContent(refined schema.RefinedAnalysis) bool {
	return len(refined.Functions) > 0 ||
		len(refined.Packages) > 0 ||
		len(refined.Types) > 0 ||
		len(refined.Strings) > 0
}
