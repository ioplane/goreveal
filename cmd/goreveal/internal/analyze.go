package internalcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/engine/peeling"
	"github.com/dantte-lp/goreveal/schema"
)

// RunAnalyze executes the minimal analyze command and writes canonical JSON.
func RunAnalyze(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("analyze %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandAnalyze, analysis); err != nil {
		return fmt.Errorf("analyze %q: %w", path, err)
	}

	return writeJSON(stdout, analysis)
}

// RunInspectFunctions writes the recovered function set as canonical JSON.
func RunInspectFunctions(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect functions %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectFunctions, analysis); err != nil {
		return fmt.Errorf("inspect functions %q: %w", path, err)
	}

	return writeJSON(stdout, analysis.Functions)
}

// RunInspectPackages writes the recovered package set as canonical JSON.
func RunInspectPackages(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect packages %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectPackages, analysis); err != nil {
		return fmt.Errorf("inspect packages %q: %w", path, err)
	}

	return writeJSON(stdout, analysis.Packages)
}

// RunInspectTypes writes the recovered type set as canonical JSON.
func RunInspectTypes(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect types %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectTypes, analysis); err != nil {
		return fmt.Errorf("inspect types %q: %w", path, err)
	}
	if analysis.Types == nil {
		analysis.Types = []schema.Type{}
	}

	return writeJSON(stdout, analysis.Types)
}

// RunInspectStrings writes the recovered string candidate set as canonical JSON.
func RunInspectStrings(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect strings %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectStrings, analysis); err != nil {
		return fmt.Errorf("inspect strings %q: %w", path, err)
	}

	return writeJSON(stdout, analysis.Strings)
}

// RunInspectRuntime writes the recovered runtime metadata as canonical JSON.
func RunInspectRuntime(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect runtime %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectRuntime, analysis); err != nil {
		return fmt.Errorf("inspect runtime %q: %w", path, err)
	}
	if analysis.Runtime == nil {
		return fmt.Errorf("inspect runtime %q: unavailable", path)
	}

	return writeJSON(stdout, analysis.Runtime)
}

// RunInspectPeeling writes the bounded code-peeling layer as canonical JSON.
func RunInspectPeeling(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("inspect peeling %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandInspectPeeling, analysis); err != nil {
		return fmt.Errorf("inspect peeling %q: %w", path, err)
	}
	if analysis.Peeling == nil {
		return fmt.Errorf("inspect peeling %q: unavailable", path)
	}

	return writeJSON(stdout, analysis.Peeling)
}

// RunPeel writes a user-only projection derived from the bounded peeling layer.
func RunPeel(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("peel %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandPeel, analysis); err != nil {
		return fmt.Errorf("peel %q: %w", path, err)
	}
	projected := peeling.UserOnlyView(analysis.Peeling)
	if projected == nil {
		return fmt.Errorf("peel %q: unavailable", path)
	}

	return writeJSON(stdout, projected)
}

// RunSourceTree writes the projected source tree as canonical JSON.
func RunSourceTree(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("source tree %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandSourceTree, analysis); err != nil {
		return fmt.Errorf("source tree %q: %w", path, err)
	}
	if analysis.SourceTree == nil {
		return fmt.Errorf("source tree %q: unavailable", path)
	}

	return writeJSON(stdout, analysis.SourceTree)
}

func writeJSON(stdout io.Writer, value any) error {
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return fmt.Errorf("encode JSON: %w", err)
	}

	return nil
}
