package refine

import (
	"context"
	"strings"

	"github.com/ioplane/goreveal/schema"
)

// Pass normalizes selected synthetic compiler-generated names into a more readable form.
type Pass struct{}

func (Pass) Name() string { return "synthetic-function-names" }

func (Pass) Apply(_ context.Context, _ schema.Analysis, refined *schema.RefinedAnalysis) error {
	for i := range refined.Functions {
		refined.Functions[i].Name = refineFunctionName(refined.Functions[i].Name)
	}
	return nil
}

func refineFunctionName(name string) string {
	if trimmed, ok := strings.CutPrefix(name, "type:.eq."); ok {
		return "eq(" + trimmed + ")"
	}
	return name
}
