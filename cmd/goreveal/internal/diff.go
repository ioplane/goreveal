package internalcmd

import (
	"context"
	"fmt"
	"io"

	storagediff "github.com/dantte-lp/goreveal/storage/diff"
	storesqlite "github.com/dantte-lp/goreveal/storage/sqlite"
)

type storedRunDiff struct {
	LeftID  int64               `json:"left_id"`
	RightID int64               `json:"right_id"`
	Summary storagediff.Summary `json:"summary"`
}

// RunDiffSQLite compares two stored analyses from the same SQLite database.
func RunDiffSQLite(ctx context.Context, stdout io.Writer, dbPath string, leftID, rightID int64) (err error) {
	store, err := storesqlite.Open(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("open sqlite store %q: %w", dbPath, err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close sqlite store %q: %w", dbPath, closeErr)
		}
	}()

	left, err := store.LoadAnalysis(ctx, leftID)
	if err != nil {
		return fmt.Errorf("load left analysis %d: %w", leftID, err)
	}
	right, err := store.LoadAnalysis(ctx, rightID)
	if err != nil {
		return fmt.Errorf("load right analysis %d: %w", rightID, err)
	}

	return writeJSON(stdout, storedRunDiff{
		LeftID:  leftID,
		RightID: rightID,
		Summary: storagediff.Compare(left, right),
	})
}
