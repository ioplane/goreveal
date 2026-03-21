package deobfuscation

import (
	"context"
	"fmt"

	"github.com/dantte-lp/goreveal/schema"
)

// Pass refines the mutable deobfuscation layer without mutating raw recovered truth.
type Pass interface {
	Name() string
	Apply(ctx context.Context, analysis schema.Analysis, refined *schema.RefinedAnalysis) error
}

// Pipeline applies zero or more refinement passes to a raw analysis.
type Pipeline struct {
	passes []Pass
}

// NewPipeline constructs a refinement pipeline.
func NewPipeline(passes ...Pass) Pipeline {
	return Pipeline{passes: passes}
}

// Run builds the initial refined layer from raw truth, then applies passes in order.
func (p Pipeline) Run(ctx context.Context, analysis schema.Analysis) (schema.RefinedAnalysis, error) {
	refined := schema.RefinedAnalysis{
		Functions: make([]schema.RefinedFunction, 0, len(analysis.Functions)),
		Packages:  make([]schema.RefinedPackage, 0, len(analysis.Packages)),
		Types:     make([]schema.RefinedType, 0, len(analysis.Types)),
		Strings:   make([]schema.RefinedString, 0, len(analysis.Strings)),
	}

	for _, fn := range analysis.Functions {
		refined.Functions = append(refined.Functions, schema.RefinedFunction{Name: fn.Name})
	}
	for _, pkg := range analysis.Packages {
		refined.Packages = append(refined.Packages, schema.RefinedPackage{Name: pkg.Name})
	}
	for _, typ := range analysis.Types {
		refined.Types = append(refined.Types, schema.RefinedType{Name: typ.Name})
	}
	for _, str := range analysis.Strings {
		refined.Strings = append(refined.Strings, schema.RefinedString{Value: str.Value})
	}

	for _, pass := range p.passes {
		if err := pass.Apply(ctx, analysis, &refined); err != nil {
			return schema.RefinedAnalysis{}, fmt.Errorf("apply pass %q: %w", pass.Name(), err)
		}
		refined.Passes = append(refined.Passes, pass.Name())
	}

	return refined, nil
}
