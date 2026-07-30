package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/ioplane/goreveal/schema"
)

//go:embed schema.sql
var schemaSQL string

// Store persists canonical analyses in a local SQLite database.
type Store struct {
	db *sql.DB
}

// Open opens or creates a SQLite-backed analysis store.
func Open(ctx context.Context, path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	store := &Store{db: db}
	if err := store.init(ctx); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, errors.Join(err, fmt.Errorf("close sqlite database after init failure: %w", closeErr))
		}
		return nil, err
	}

	return store, nil
}

func (s *Store) init(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}
	return nil
}

// Close closes the underlying database handle.
func (s *Store) Close() error {
	return s.db.Close()
}

// SaveAnalysis inserts one canonical analysis and returns its row id.
func (s *Store) SaveAnalysis(ctx context.Context, analysis schema.Analysis) (int64, error) {
	payload, err := json.Marshal(analysis)
	if err != nil {
		return 0, fmt.Errorf("marshal analysis: %w", err)
	}

	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO analyses (input_path, format, analysis_json) VALUES (?, ?, ?)`,
		analysis.Input.Path,
		analysis.Input.Format,
		string(payload),
	)
	if err != nil {
		return 0, fmt.Errorf("insert analysis: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("read inserted id: %w", err)
	}
	return id, nil
}

// LoadAnalysis loads one canonical analysis by row id.
func (s *Store) LoadAnalysis(ctx context.Context, id int64) (schema.Analysis, error) {
	var payload string
	if err := s.db.QueryRowContext(ctx, `SELECT analysis_json FROM analyses WHERE id = ?`, id).Scan(&payload); err != nil {
		return schema.Analysis{}, fmt.Errorf("query analysis %d: %w", id, err)
	}

	var analysis schema.Analysis
	if err := json.Unmarshal([]byte(payload), &analysis); err != nil {
		return schema.Analysis{}, fmt.Errorf("unmarshal analysis %d: %w", id, err)
	}
	return analysis, nil
}
