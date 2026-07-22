package snapshots

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

	fixtureNames, err := discoverSnapshotFixtures(filepath.Join("..", "..", "corpus", "fixtures"))
	testutil.MustNoErr(t, err)

	for _, fixtureName := range fixtureNames {
		fixtureName := fixtureName
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
			normalizeAnalysisPaths(&got, fixturePath, metadata.NormalizedPath)

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

func TestNormalizeAnalysisPaths(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Input: schema.Input{Path: "../../corpus/fixture.bin"},
		Diagnostics: schema.StageDiagnostics{{
			Stage:   schema.AnalysisStageFunctions,
			Status:  schema.StageStatusFailed,
			Message: "read ../../corpus/fixture.bin: fixture failure",
		}},
	}

	normalizeAnalysisPaths(&analysis, "../../corpus/fixture.bin", "corpus/fixture.bin")
	if analysis.Input.Path != "corpus/fixture.bin" {
		t.Fatalf("Input.Path = %q, want normalized path", analysis.Input.Path)
	}
	if analysis.Diagnostics[0].Message != "read corpus/fixture.bin: fixture failure" {
		t.Fatalf("diagnostic message = %q, want normalized path", analysis.Diagnostics[0].Message)
	}
}

func normalizeAnalysisPaths(analysis *schema.Analysis, fixturePath, normalizedPath string) {
	analysis.Input.Path = normalizedPath
	for index := range analysis.Diagnostics {
		analysis.Diagnostics[index].Message = strings.ReplaceAll(
			analysis.Diagnostics[index].Message,
			fixturePath,
			normalizedPath,
		)
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
