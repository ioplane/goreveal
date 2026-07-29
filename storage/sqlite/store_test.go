package sqlite

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ioplane/goreveal/schema"
)

func TestStoreSaveAndLoad(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "analysis.db")

	store, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	analysis := schema.Analysis{
		Input: schema.Input{
			Path:   "corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin",
			Format: "elf",
			Size:   123,
		},
		Provenance: schema.Provenance{
			Source:     "test",
			Confidence: "high",
		},
		Functions: []schema.Function{
			{Name: "main.main", Entry: 0x401000},
		},
	}

	id, err := store.SaveAnalysis(context.Background(), analysis)
	if err != nil {
		t.Fatalf("SaveAnalysis() error = %v", err)
	}
	if id == 0 {
		t.Fatal("SaveAnalysis() returned zero id")
	}

	got, err := store.LoadAnalysis(context.Background(), id)
	if err != nil {
		t.Fatalf("LoadAnalysis() error = %v", err)
	}

	if got.Input.Path != analysis.Input.Path {
		t.Fatalf("LoadAnalysis() path = %q, want %q", got.Input.Path, analysis.Input.Path)
	}
	if len(got.Functions) != 1 || got.Functions[0].Name != "main.main" {
		t.Fatalf("LoadAnalysis() functions = %#v", got.Functions)
	}
}

func TestStoreSnapshotRoundTripIntegrity(t *testing.T) {
	t.Parallel()

	dbPath := filepath.Join(t.TempDir(), "snapshot.db")

	store, err := Open(context.Background(), dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer func() {
		if closeErr := store.Close(); closeErr != nil {
			t.Fatalf("Close() error = %v", closeErr)
		}
	}()

	snapshotPath := filepath.Join("..", "..", "corpus", "fixtures", "minimal-linux-amd64", "expected.analysis.json")
	snapshotBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	var analysis schema.Analysis
	if unmarshalErr := json.Unmarshal(snapshotBytes, &analysis); unmarshalErr != nil {
		t.Fatalf("Unmarshal() error = %v", unmarshalErr)
	}

	id, err := store.SaveAnalysis(context.Background(), analysis)
	if err != nil {
		t.Fatalf("SaveAnalysis() error = %v", err)
	}

	got, err := store.LoadAnalysis(context.Background(), id)
	if err != nil {
		t.Fatalf("LoadAnalysis() error = %v", err)
	}

	if !reflect.DeepEqual(got, analysis) {
		t.Fatalf("round-trip mismatch\ngot: %#v\nwant: %#v", got, analysis)
	}
}
