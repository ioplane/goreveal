package garble

import (
	"context"
	"testing"

	"github.com/dantte-lp/goreveal/deobfuscation"
	"github.com/dantte-lp/goreveal/schema"
)

func TestStringSegmentPassExtractsUsefulSubstring(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		Strings: []schema.StringCandidate{
			{Value: "goreveal fixture0123456789"},
		},
	}

	refined, err := deobfuscation.NewPipeline(Pass{}).Run(context.Background(), analysis)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !containsRefinedString(refined.Strings, "goreveal fixture") {
		t.Fatalf("refined strings = %#v", refined.Strings)
	}
}

func containsRefinedString(strings []schema.RefinedString, want string) bool {
	for _, str := range strings {
		if str.Value == want {
			return true
		}
	}
	return false
}
