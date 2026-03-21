package schema

// RefinedFunction carries a deobfuscation-safe function view that does not mutate raw truth.
type RefinedFunction struct {
	Name string `json:"name"`
}

// RefinedPackage carries a refined package view.
type RefinedPackage struct {
	Name string `json:"name"`
}

// RefinedType carries a refined type view.
type RefinedType struct {
	Name string `json:"name"`
}

// RefinedString carries a refined string view.
type RefinedString struct {
	Value string `json:"value"`
}

// RefinedAnalysis is the mutable refinement layer derived from raw recovered truth.
type RefinedAnalysis struct {
	Functions []RefinedFunction `json:"functions,omitempty"`
	Packages  []RefinedPackage  `json:"packages,omitempty"`
	Types     []RefinedType     `json:"types,omitempty"`
	Strings   []RefinedString   `json:"strings,omitempty"`
	Passes    []string          `json:"passes,omitempty"`
}
