package packages

import (
	"slices"
	"strings"

	"github.com/ioplane/goreveal/schema"
)

// Package is the canonical recovered package representation.
type Package = schema.Package

// Recover classifies packages from recovered function names.
func Recover(funcs []schema.Function) []Package {
	counts := make(map[string]int)
	for _, fn := range funcs {
		name, ok := packageName(fn.Name)
		if !ok {
			continue
		}
		counts[name]++
	}

	pkgs := make([]Package, 0, len(counts))
	for name, count := range counts {
		importPath := ""
		if name != "main" {
			importPath = name
		}
		pkgs = append(pkgs, Package{
			Name:          name,
			ImportPath:    importPath,
			FunctionCount: count,
			Provenance: schema.Provenance{
				Source:     "core.packages.functions",
				Confidence: "medium",
			},
		})
	}

	slices.SortFunc(pkgs, func(a, b Package) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	return pkgs
}

func packageName(functionName string) (string, bool) {
	if strings.HasPrefix(functionName, "type:") {
		return "", false
	}

	lastSlash := strings.LastIndex(functionName, "/")
	searchStart := max(lastSlash+1, 0)

	dot := strings.Index(functionName[searchStart:], ".")
	if dot <= 0 {
		return "", false
	}

	name := functionName[:searchStart+dot]
	if name == "" {
		return "", false
	}

	return name, true
}
