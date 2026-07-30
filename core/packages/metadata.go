package packages

import "github.com/ioplane/goreveal/schema"

// EnrichSourceMetadata correlates recovered packages with source-tree evidence.
func EnrichSourceMetadata(pkgs []Package, tree *schema.SourceTree) []Package {
	if len(pkgs) == 0 || tree == nil {
		return pkgs
	}

	moduleByImportPath := make(map[string]schema.SourcePackage, len(tree.Packages))
	moduleByName := make(map[string]schema.SourcePackage, len(tree.Packages))
	for _, pkg := range tree.Packages {
		moduleByImportPath[pkg.ImportPath] = pkg
		if pkg.Name != "" {
			moduleByName[pkg.Name] = pkg
		}
	}

	externalByImportPath := make(map[string]schema.SourcePackage, len(tree.ExternalPackages))
	externalByName := make(map[string]schema.SourcePackage, len(tree.ExternalPackages))
	for _, pkg := range tree.ExternalPackages {
		externalByImportPath[pkg.ImportPath] = pkg
		if pkg.Name != "" {
			externalByName[pkg.Name] = pkg
		}
	}

	enriched := make([]Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		match, scope, ok := matchSourcePackage(
			pkg,
			tree.Root,
			moduleByImportPath,
			moduleByName,
			externalByImportPath,
			externalByName,
		)
		if ok {
			pkg.ImportPath = match.ImportPath
			pkg.SourceFileCount = len(match.Files)
			pkg.HasSourceEvidence = match.HasFileEvidence
			pkg.SourceEvidenceKind = match.SourceEvidenceKind
			if pkg.SourceEvidenceKind == "" {
				pkg.SourceEvidenceKind = tree.SourceEvidenceKind
			}
			pkg.ModuleLocal = scope == sourceScopeModule
		}
		enriched = append(enriched, pkg)
	}

	return enriched
}

// EnrichBuildInfoMetadata adds low-risk package metadata when source-tree evidence
// is unavailable but build info still exposes the main module path.
func EnrichBuildInfoMetadata(pkgs []Package, info *schema.BuildInfo) []Package {
	if len(pkgs) == 0 || info == nil || info.Path == "" {
		return pkgs
	}

	enriched := make([]Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		if pkg.Name == "main" {
			pkg.ImportPath = info.Path
			pkg.ModuleLocal = true
			pkg.SourceEvidenceKind = schema.SourceEvidenceKindPackageFallback
		}
		enriched = append(enriched, pkg)
	}

	return enriched
}

type sourceScope int

const (
	sourceScopeUnknown sourceScope = iota
	sourceScopeModule
	sourceScopeExternal
)

func matchSourcePackage(
	pkg Package,
	root string,
	moduleByImportPath, moduleByName, externalByImportPath, externalByName map[string]schema.SourcePackage,
) (schema.SourcePackage, sourceScope, bool) {
	if pkg.Name == "main" {
		if match, ok := moduleByImportPath[root]; ok {
			return match, sourceScopeModule, true
		}
	}
	if match, ok := moduleByImportPath[pkg.Name]; ok {
		return match, sourceScopeModule, true
	}
	if match, ok := externalByImportPath[pkg.Name]; ok {
		return match, sourceScopeExternal, true
	}
	if match, ok := moduleByName[pkg.Name]; ok {
		return match, sourceScopeModule, true
	}
	if match, ok := externalByName[pkg.Name]; ok {
		return match, sourceScopeExternal, true
	}

	return schema.SourcePackage{}, sourceScopeUnknown, false
}
