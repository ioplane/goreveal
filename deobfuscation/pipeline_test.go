package deobfuscation

import (
	"context"
	"testing"

	"github.com/ioplane/goreveal/schema"
)

func TestPipelinePreservesRawTruth(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Functions: []schema.Function{
			{Name: "main.main"},
		},
		Packages: []schema.Package{
			{Name: "main"},
		},
		Types: []schema.Type{
			{Name: "main.fixtureCounter"},
		},
		Strings: []schema.StringCandidate{
			{Value: "goreveal fixture"},
		},
	}

	refined, err := NewPipeline(renameFunctionPass{}).Run(context.Background(), analysis)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if analysis.Functions[0].Name != "main.main" {
		t.Fatalf("raw function mutated: %#v", analysis.Functions[0])
	}
	if refined.Functions[0].Name != "main.cleaned" {
		t.Fatalf("refined function name = %q", refined.Functions[0].Name)
	}
	if refined.Packages[0].Name != "main" {
		t.Fatalf("refined package name = %q", refined.Packages[0].Name)
	}
	if refined.Types[0].Name != "main.fixtureCounter" {
		t.Fatalf("refined type name = %q", refined.Types[0].Name)
	}
	if refined.Strings[0].Value != "goreveal fixture" {
		t.Fatalf("refined string value = %q", refined.Strings[0].Value)
	}
	if len(refined.Passes) != 1 || refined.Passes[0] != "rename-function" {
		t.Fatalf("refined passes = %#v", refined.Passes)
	}
}

type renameFunctionPass struct{}

func (renameFunctionPass) Name() string { return "rename-function" }

func (renameFunctionPass) Apply(_ context.Context, _ schema.Analysis, refined *schema.RefinedAnalysis) error {
	refined.Functions[0].Name = "main.cleaned"
	return nil
}
