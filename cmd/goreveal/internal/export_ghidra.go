package internalcmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/schema"
)

// RunExportGhidra writes the stable Ghidra-oriented export payload.
func RunExportGhidra(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("export ghidra %q: %w", path, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandExportGhidraV1, analysis); err != nil {
		return fmt.Errorf("export ghidra %q: %w", path, err)
	}

	return writeJSON(stdout, schema.NewGhidraExport(analysis))
}
