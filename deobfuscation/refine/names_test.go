package refine

import (
	"context"
	"testing"

	"github.com/ioplane/goreveal/deobfuscation"
	"github.com/ioplane/goreveal/schema"
)

func TestSyntheticNamePassRefinesTypeEqFunctions(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Functions: []schema.Function{
			{Name: "type:.eq.main.fixtureCounter"},
			{Name: "main.main"},
		},
	}

	refined, err := deobfuscation.NewPipeline(Pass{}).Run(context.Background(), analysis)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if refined.Functions[0].Name != "eq(main.fixtureCounter)" {
		t.Fatalf("refined function = %q", refined.Functions[0].Name)
	}
	if refined.Functions[1].Name != "main.main" {
		t.Fatalf("non-synthetic function changed = %q", refined.Functions[1].Name)
	}
}
