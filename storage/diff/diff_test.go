package diff

import (
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	left := schema.Analysis{
		Input: schema.Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		BuildInfo: &schema.BuildInfo{
			GoVersion: "go1.26.1",
			Path:      "example.com/sample",
		},
		Functions: []schema.Function{
			{Name: "main.main", Entry: 0x1000, End: 0x1100},
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 1},
		},
		Types: []schema.Type{
			{Name: "main.counter", Kind: "struct"},
		},
		Strings: []schema.StringCandidate{
			{Value: "hello", Region: ".rodata", Offset: 8},
		},
		SourceTree: &schema.SourceTree{
			Root:  "example.com/sample",
			Files: []string{"main.go"},
		},
		Refined: &schema.RefinedAnalysis{
			Passes: []string{"synthetic-function-names"},
		},
	}

	right := schema.Analysis{
		Input: schema.Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		BuildInfo: &schema.BuildInfo{
			GoVersion: "go1.26.2",
			Path:      "example.com/sample/v2",
		},
		Functions: []schema.Function{
			{Name: "main.main", Entry: 0x1000, End: 0x1100},
			{Name: "main.helper", Entry: 0x1200, End: 0x1220},
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 2},
			{Name: "internal/app", FunctionCount: 1},
		},
		Types: []schema.Type{
			{Name: "main.counter", Kind: "struct"},
			{Name: "main.extra", Kind: "string"},
		},
		Strings: []schema.StringCandidate{
			{Value: "world", Region: ".rodata", Offset: 24},
		},
		SourceTree: &schema.SourceTree{
			Root:  "example.com/sample/v2",
			Files: []string{"main.go", "helper.go"},
		},
		Refined: &schema.RefinedAnalysis{
			Passes: []string{"synthetic-function-names", "string-segments"},
		},
	}

	got := Compare(left, right)

	if !got.BuildInfoChanged {
		t.Fatal("BuildInfoChanged = false, want true")
	}
	if got.LeftCounts.Functions != 1 || got.RightCounts.Functions != 2 {
		t.Fatalf("function counts = %#v %#v", got.LeftCounts, got.RightCounts)
	}
	if len(got.AddedFunctions) != 1 || got.AddedFunctions[0] != "main.helper" {
		t.Fatalf("added functions = %#v", got.AddedFunctions)
	}
	if len(got.AddedPackages) != 1 || got.AddedPackages[0] != "internal/app" {
		t.Fatalf("added packages = %#v", got.AddedPackages)
	}
	if len(got.AddedTypes) != 1 || got.AddedTypes[0] != "main.extra" {
		t.Fatalf("added types = %#v", got.AddedTypes)
	}
	if len(got.RemovedStrings) != 1 || got.RemovedStrings[0] != "hello" {
		t.Fatalf("removed strings = %#v", got.RemovedStrings)
	}
	if len(got.AddedSourceFiles) != 1 || got.AddedSourceFiles[0] != "helper.go" {
		t.Fatalf("added source files = %#v", got.AddedSourceFiles)
	}
	if len(got.AddedRefinedPasses) != 1 || got.AddedRefinedPasses[0] != "string-segments" {
		t.Fatalf("added refined passes = %#v", got.AddedRefinedPasses)
	}
}
