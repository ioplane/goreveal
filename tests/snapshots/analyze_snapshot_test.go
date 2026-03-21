package snapshots

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/internal/testutil"
	"github.com/dantte-lp/goreveal/schema"
)

var updateSnapshots = flag.Bool("update", false, "update analysis snapshots")

type fixtureMetadata struct {
	Fixture        string `json:"fixture"`
	NormalizedPath string `json:"normalized_path"`
}

func TestAnalyzeFixtureSnapshot(t *testing.T) {
	t.Parallel()

	fixtureDir := filepath.Join("..", "..", "corpus", "fixtures", "minimal-linux-amd64")

	metadataBytes, err := os.ReadFile(filepath.Join(fixtureDir, "fixture.json"))
	testutil.MustNoErr(t, err)

	var metadata fixtureMetadata
	testutil.MustNoErr(t, json.Unmarshal(metadataBytes, &metadata))

	fixturePath := filepath.Join(fixtureDir, metadata.Fixture)
	if _, statErr := os.Stat(fixturePath); statErr != nil {
		t.Fatalf("fixture binary missing: %v", statErr)
	}

	got, err := engine.New().AnalyzeFile(context.Background(), fixturePath)
	testutil.MustNoErr(t, err)
	got.Input.Path = metadata.NormalizedPath

	expectedPath := filepath.Join(fixtureDir, "expected.analysis.json")
	if *updateSnapshots {
		gotJSON, marshalErr := json.MarshalIndent(got, "", "  ")
		testutil.MustNoErr(t, marshalErr)
		testutil.MustNoErr(t, os.WriteFile(expectedPath, append(gotJSON, '\n'), 0o644))
		return
	}

	expectedBytes, err := os.ReadFile(expectedPath)
	testutil.MustNoErr(t, err)

	var want schema.Analysis
	testutil.MustNoErr(t, json.Unmarshal(expectedBytes, &want))

	if !reflect.DeepEqual(got, want) {
		gotJSON, gotMarshalErr := json.MarshalIndent(got, "", "  ")
		testutil.MustNoErr(t, gotMarshalErr)
		wantJSON, wantMarshalErr := json.MarshalIndent(want, "", "  ")
		testutil.MustNoErr(t, wantMarshalErr)
		t.Fatalf("snapshot mismatch\nwant:\n%s\n\ngot:\n%s", wantJSON, gotJSON)
	}
}
