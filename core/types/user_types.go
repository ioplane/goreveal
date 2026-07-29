package types

import (
	"path"
	"strings"

	"github.com/ioplane/goreveal/schema"
)

// EnrichUserMetadata annotates types with package and module-local usefulness hints.
func EnrichUserMetadata(types []Type, tree *schema.SourceTree) []Type {
	if len(types) == 0 {
		return types
	}

	modulePackages := make(map[string]schema.SourcePackage)
	externalPackages := make(map[string]schema.SourcePackage)
	if tree != nil {
		for _, pkg := range tree.Packages {
			addPackageKeys(modulePackages, pkg)
		}
		for _, pkg := range tree.ExternalPackages {
			addPackageKeys(externalPackages, pkg)
		}
	}

	enriched := make([]Type, 0, len(types))
	for _, typ := range types {
		pkg := typePackage(typ.Name)
		if pkg != "" {
			typ.Package = pkg
			if match, ok := modulePackages[pkg]; ok {
				typ.ImportPath = match.ImportPath
				typ.SourceFileCount = len(match.Files)
				typ.ModuleLocal = true
				typ.UserMeaningful = true
			} else if match, ok := externalPackages[pkg]; ok {
				typ.ImportPath = match.ImportPath
			}
		}
		enriched = append(enriched, typ)
	}

	return enriched
}

// EnrichBuildInfoMetadata adds low-risk type metadata when source-tree evidence
// is unavailable but build info still exposes the main module path.
func EnrichBuildInfoMetadata(types []Type, info *schema.BuildInfo) []Type {
	if len(types) == 0 || info == nil || info.Path == "" {
		return types
	}

	enriched := make([]Type, 0, len(types))
	for _, typ := range types {
		pkg := typePackage(typ.Name)
		if pkg != "" {
			typ.Package = pkg
			typ.ImportPath = pkg
		}
		if pkg == "main" {
			typ.ImportPath = info.Path
			typ.ModuleLocal = true
			typ.UserMeaningful = true
		}
		enriched = append(enriched, typ)
	}

	return enriched
}

func addPackageKeys(dst map[string]schema.SourcePackage, pkg schema.SourcePackage) {
	if pkg.Name != "" {
		dst[pkg.Name] = pkg
	}
	if pkg.ImportPath != "" {
		dst[pkg.ImportPath] = pkg
		dst[path.Base(pkg.ImportPath)] = pkg
	}
}

func typePackage(name string) string {
	cleaned := strings.TrimSpace(name)
	for {
		switch {
		case strings.HasPrefix(cleaned, "*"):
			cleaned = strings.TrimPrefix(cleaned, "*")
		case strings.HasPrefix(cleaned, "[]"):
			cleaned = strings.TrimPrefix(cleaned, "[]")
		case strings.HasPrefix(cleaned, "["):
			end := strings.Index(cleaned, "]")
			if end < 0 {
				return ""
			}
			cleaned = cleaned[end+1:]
		default:
			goto parsed
		}
	}

parsed:
	dot := strings.LastIndex(cleaned, ".")
	if dot <= 0 {
		return ""
	}

	return cleaned[:dot]
}
