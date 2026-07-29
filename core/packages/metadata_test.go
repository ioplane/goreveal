package packages

import (
	"testing"

	"github.com/ioplane/goreveal/schema"
)

func TestEnrichSourceMetadata(t *testing.T) {
	t.Parallel()

	pkgs := []Package{
		{Name: "main", FunctionCount: 3},
		{Name: "runtime", FunctionCount: 10},
		{Name: "example.com/gorevealfixture/pkg/sub", FunctionCount: 2},
	}

	tree := &schema.SourceTree{
		Root:               "example.com/gorevealfixture",
		SourceEvidenceKind: schema.SourceEvidenceKindDWARFPaths,
		Packages: []schema.SourcePackage{
			{
				Name:            "main",
				ImportPath:      "example.com/gorevealfixture",
				FunctionCount:   3,
				HasFileEvidence: true,
				Files:           []string{"main.go"},
			},
			{
				Name:            "sub",
				ImportPath:      "example.com/gorevealfixture/pkg/sub",
				FunctionCount:   2,
				HasFileEvidence: true,
				Files:           []string{"pkg/sub/a.go", "pkg/sub/b.go"},
			},
		},
		ExternalPackages: []schema.SourcePackage{
			{
				Name:            "runtime",
				ImportPath:      "runtime",
				HasFileEvidence: true,
				Files:           []string{"/usr/local/go/src/runtime/proc.go"},
			},
		},
	}

	got := EnrichSourceMetadata(pkgs, tree)

	if got[0].ImportPath != "example.com/gorevealfixture" || got[0].SourceFileCount != 1 || !got[0].ModuleLocal || !got[0].HasSourceEvidence {
		t.Fatalf("main package = %#v", got[0])
	}
	if got[0].SourceEvidenceKind != schema.SourceEvidenceKindDWARFPaths {
		t.Fatalf("main source evidence kind = %#v", got[0])
	}
	if got[1].ImportPath != "runtime" || got[1].SourceFileCount != 1 || got[1].ModuleLocal || !got[1].HasSourceEvidence {
		t.Fatalf("runtime package = %#v", got[1])
	}
	if got[1].SourceEvidenceKind != schema.SourceEvidenceKindDWARFPaths {
		t.Fatalf("runtime source evidence kind = %#v", got[1])
	}
	if got[2].ImportPath != "example.com/gorevealfixture/pkg/sub" || got[2].SourceFileCount != 2 || !got[2].ModuleLocal || !got[2].HasSourceEvidence {
		t.Fatalf("sub package = %#v", got[2])
	}
	if got[2].SourceEvidenceKind != schema.SourceEvidenceKindDWARFPaths {
		t.Fatalf("sub source evidence kind = %#v", got[2])
	}
}

func TestEnrichBuildInfoMetadata(t *testing.T) {
	t.Parallel()

	pkgs := []Package{
		{Name: "main", FunctionCount: 3},
		{Name: "runtime", FunctionCount: 10},
	}

	got := EnrichBuildInfoMetadata(pkgs, &schema.BuildInfo{Path: "example.com/gorevealfixture"})

	if got[0].ImportPath != "example.com/gorevealfixture" {
		t.Fatalf("main import path = %q, want %q", got[0].ImportPath, "example.com/gorevealfixture")
	}
	if !got[0].ModuleLocal {
		t.Fatal("main package should be module-local after build info enrichment")
	}
	if got[0].HasSourceEvidence {
		t.Fatal("main package should not claim source evidence from build info enrichment alone")
	}
	if got[0].SourceEvidenceKind != schema.SourceEvidenceKindPackageFallback {
		t.Fatalf("main source evidence kind = %#v", got[0])
	}
	if got[1].ImportPath != "" {
		t.Fatalf("runtime import path = %q, want empty", got[1].ImportPath)
	}
	if got[1].ModuleLocal {
		t.Fatal("runtime package should not become module-local from build info enrichment")
	}
	if got[1].HasSourceEvidence {
		t.Fatal("runtime package should not gain source evidence without source-tree correlation")
	}
	if got[1].SourceEvidenceKind != "" {
		t.Fatalf("runtime source evidence kind = %#v", got[1])
	}
}
