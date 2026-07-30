package schema_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ioplane/goreveal/engine"
	"github.com/ioplane/goreveal/schema"
)

func TestFixtureIDAExportSimulation(t *testing.T) {
	t.Parallel()

	fixture := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	analysis, err := engine.New().AnalyzeFile(context.Background(), fixture)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	got := schema.NewIDAExport(analysis)

	if got.Contract != schema.IDAExportContractV1 {
		t.Fatalf("contract = %q", got.Contract)
	}
	if got.BuildInfo == nil || got.BuildInfo.Path != "example.com/gorevealfixture" {
		t.Fatalf("build info = %#v", got.BuildInfo)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != schema.RuntimeTrustSummarySymbolBacked {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) == 0 {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if len(got.Peeling.Packages) == 0 {
		t.Fatalf("peeling packages = %#v", got.Peeling)
	}
	if len(got.Functions) == 0 || got.Functions[0].Entry == 0 {
		t.Fatalf("functions = %#v", got.Functions)
	}
	if len(got.Packages) == 0 || got.Packages[0].Name == "" {
		t.Fatalf("packages = %#v", got.Packages)
	}
}

func TestFixtureGhidraExportSimulation(t *testing.T) {
	t.Parallel()

	fixture := filepath.Join("..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	analysis, err := engine.New().AnalyzeFile(context.Background(), fixture)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	got := schema.NewGhidraExport(analysis)

	if got.Contract != schema.GhidraExportContractV1 {
		t.Fatalf("contract = %q", got.Contract)
	}
	if got.Program.ModulePath != "example.com/gorevealfixture" {
		t.Fatalf("module path = %q", got.Program.ModulePath)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != schema.RuntimeTrustSummarySymbolBacked {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) == 0 {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if len(got.Peeling.Packages) == 0 {
		t.Fatalf("peeling packages = %#v", got.Peeling)
	}
	if len(got.Symbols) == 0 || got.Symbols[0].Address == 0 {
		t.Fatalf("symbols = %#v", got.Symbols)
	}
	if got.SourceTree == nil || got.SourceTree.Root != "example.com/gorevealfixture" {
		t.Fatalf("source tree = %#v", got.SourceTree)
	}
}
