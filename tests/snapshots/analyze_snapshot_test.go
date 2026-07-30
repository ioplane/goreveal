package snapshots

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ioplane/goreveal/engine"
	"github.com/ioplane/goreveal/internal/testutil"
	"github.com/ioplane/goreveal/schema"
)

var updateSnapshots = flag.Bool("update", false, "update analysis snapshots")

type fixtureMetadata struct {
	Fixture        string `json:"fixture"`
	NormalizedPath string `json:"normalized_path"`
}

func TestAnalyzeFixtureSnapshot(t *testing.T) {
	t.Parallel()

	fixtureNames, err := discoverSnapshotFixtures(filepath.Join("..", "..", "corpus", "fixtures"))
	testutil.MustNoErr(t, err)

	for _, fixtureName := range fixtureNames {
		t.Run(fixtureName, func(t *testing.T) {
			t.Parallel()

			fixtureDir := filepath.Join("..", "..", "corpus", "fixtures", fixtureName)

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

			gotJSON, gotMarshalErr := json.MarshalIndent(got, "", "  ")
			testutil.MustNoErr(t, gotMarshalErr)
			wantJSON, wantMarshalErr := json.MarshalIndent(want, "", "  ")
			testutil.MustNoErr(t, wantMarshalErr)

			if !bytes.Equal(gotJSON, wantJSON) {
				t.Fatalf("snapshot mismatch\nwant:\n%s\n\ngot:\n%s", wantJSON, gotJSON)
			}
		})
	}
}

func discoverSnapshotFixtures(fixturesRoot string) ([]string, error) {
	entries, err := os.ReadDir(fixturesRoot)
	if err != nil {
		return nil, err
	}

	fixtures := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		expectedPath := filepath.Join(fixturesRoot, entry.Name(), "expected.analysis.json")
		if _, statErr := os.Stat(expectedPath); statErr == nil {
			fixtures = append(fixtures, entry.Name())
		}
	}

	sort.Strings(fixtures)
	return fixtures, nil
}
