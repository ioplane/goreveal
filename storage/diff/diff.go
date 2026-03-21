package diff

import (
	"slices"

	"github.com/dantte-lp/goreveal/schema"
)

// Counts summarizes cardinalities at the schema boundary.
type Counts struct {
	Functions int `json:"functions"`
	Packages  int `json:"packages"`
	Types     int `json:"types"`
	Strings   int `json:"strings"`
	Files     int `json:"files"`
}

// Summary is the stable v1 diff payload between two stored analyses.
type Summary struct {
	BuildInfoChanged     bool     `json:"build_info_changed"`
	LeftCounts           Counts   `json:"left_counts"`
	RightCounts          Counts   `json:"right_counts"`
	AddedFunctions       []string `json:"added_functions,omitempty"`
	RemovedFunctions     []string `json:"removed_functions,omitempty"`
	AddedPackages        []string `json:"added_packages,omitempty"`
	RemovedPackages      []string `json:"removed_packages,omitempty"`
	AddedTypes           []string `json:"added_types,omitempty"`
	RemovedTypes         []string `json:"removed_types,omitempty"`
	AddedStrings         []string `json:"added_strings,omitempty"`
	RemovedStrings       []string `json:"removed_strings,omitempty"`
	AddedSourceFiles     []string `json:"added_source_files,omitempty"`
	RemovedSourceFiles   []string `json:"removed_source_files,omitempty"`
	AddedRefinedPasses   []string `json:"added_refined_passes,omitempty"`
	RemovedRefinedPasses []string `json:"removed_refined_passes,omitempty"`
}

// Compare computes a stable summary diff between two canonical analyses.
func Compare(left, right schema.Analysis) Summary {
	out := Summary{
		BuildInfoChanged:     buildInfoChanged(left.BuildInfo, right.BuildInfo),
		LeftCounts:           counts(left),
		RightCounts:          counts(right),
		AddedFunctions:       difference(functionNames(right.Functions), functionNames(left.Functions)),
		RemovedFunctions:     difference(functionNames(left.Functions), functionNames(right.Functions)),
		AddedPackages:        difference(packageNames(right.Packages), packageNames(left.Packages)),
		RemovedPackages:      difference(packageNames(left.Packages), packageNames(right.Packages)),
		AddedTypes:           difference(typeNames(right.Types), typeNames(left.Types)),
		RemovedTypes:         difference(typeNames(left.Types), typeNames(right.Types)),
		AddedStrings:         difference(stringValues(right.Strings), stringValues(left.Strings)),
		RemovedStrings:       difference(stringValues(left.Strings), stringValues(right.Strings)),
		AddedSourceFiles:     difference(sourceFiles(right.SourceTree), sourceFiles(left.SourceTree)),
		RemovedSourceFiles:   difference(sourceFiles(left.SourceTree), sourceFiles(right.SourceTree)),
		AddedRefinedPasses:   difference(refinedPasses(right.Refined), refinedPasses(left.Refined)),
		RemovedRefinedPasses: difference(refinedPasses(left.Refined), refinedPasses(right.Refined)),
	}
	return out
}

func buildInfoChanged(left, right *schema.BuildInfo) bool {
	switch {
	case left == nil && right == nil:
		return false
	case left == nil || right == nil:
		return true
	default:
		return left.GoVersion != right.GoVersion || left.Path != right.Path
	}
}

func counts(analysis schema.Analysis) Counts {
	return Counts{
		Functions: len(analysis.Functions),
		Packages:  len(analysis.Packages),
		Types:     len(analysis.Types),
		Strings:   len(analysis.Strings),
		Files:     len(sourceFiles(analysis.SourceTree)),
	}
}

func functionNames(items []schema.Function) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func packageNames(items []schema.Package) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func typeNames(items []schema.Type) []string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	slices.Sort(names)
	return slices.Compact(names)
}

func stringValues(items []schema.StringCandidate) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	slices.Sort(values)
	return slices.Compact(values)
}

func sourceFiles(tree *schema.SourceTree) []string {
	if tree == nil {
		return nil
	}
	files := append([]string(nil), tree.Files...)
	slices.Sort(files)
	return slices.Compact(files)
}

func refinedPasses(refined *schema.RefinedAnalysis) []string {
	if refined == nil {
		return nil
	}
	passes := append([]string(nil), refined.Passes...)
	slices.Sort(passes)
	return slices.Compact(passes)
}

func difference(left, right []string) []string {
	var out []string
	for _, item := range left {
		if !slices.Contains(right, item) {
			out = append(out, item)
		}
	}
	return out
}
