package internalcmd

import (
	"context"
	"fmt"
	"io"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/schema"
)

// RunExportIDA writes the stable IDA-oriented export payload.
func RunExportIDA(ctx context.Context, stdout io.Writer, path string) error {
	analysis, err := engine.New().AnalyzeFile(ctx, path)
	if err != nil {
		return fmt.Errorf("export ida %q: %w", path, err)
	}

	return writeJSON(stdout, schema.NewIDAExport(analysis))
}
