package internalcmd

import (
	"context"
	"fmt"

	"github.com/dantte-lp/goreveal/engine"
	storesqlite "github.com/dantte-lp/goreveal/storage/sqlite"
)

// RunExportSQLite analyzes a binary and persists the canonical result to SQLite.
func RunExportSQLite(ctx context.Context, dbPath, binaryPath string) (err error) {
	store, err := storesqlite.Open(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite store %q: %w", dbPath, err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close sqlite store %q: %w", dbPath, closeErr)
		}
	}()

	analysis, err := engine.New().AnalyzeFile(ctx, binaryPath)
	if err != nil {
		return fmt.Errorf("analyze %q: %w", binaryPath, err)
	}
	if err := enforceAnalysisPolicy(analysisCommandExportSQLite, analysis); err != nil {
		return fmt.Errorf("export sqlite %q: %w", binaryPath, err)
	}

	if _, err := store.SaveAnalysis(ctx, analysis); err != nil {
		return fmt.Errorf("save analysis to %q: %w", dbPath, err)
	}
	return nil
}
