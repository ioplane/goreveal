package internalcmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dantte-lp/goreveal/engine"
)

// RunDeobfuscate writes the refined layer as canonical JSON.
func RunDeobfuscate(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("deobfuscate %q: %w", path, err)
	}
	if analysis.Refined == nil {
		return fmt.Errorf("deobfuscate %q: no refined layer available", path)
	}

	return writeJSON(stdout, analysis.Refined)
}
