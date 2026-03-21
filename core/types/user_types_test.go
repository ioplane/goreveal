package types

import (
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestEnrichUserMetadata(t *testing.T) {
	t.Parallel()

	tree := &schema.SourceTree{
		Root: "example.com/gorevealfixture",
		Packages: []schema.SourcePackage{
			{Name: "main", ImportPath: "example.com/gorevealfixture", Files: []string{"main.go"}},
			{Name: "sub", ImportPath: "example.com/gorevealfixture/pkg/sub", Files: []string{"pkg/sub/a.go", "pkg/sub/b.go"}},
		},
		ExternalPackages: []schema.SourcePackage{
			{Name: "runtime", ImportPath: "runtime", Files: []string{"/usr/local/go/src/runtime/proc.go"}},
		},
	}

	got := EnrichUserMetadata([]Type{
		{Name: "main.fixtureCounter", Kind: "struct"},
		{Name: "**runtime.g", Kind: "pointer"},
		{Name: "example.com/gorevealfixture/pkg/sub.Counter", Kind: "struct"},
	}, tree)

	if got[0].Package != "main" || got[0].ImportPath != "example.com/gorevealfixture" || got[0].SourceFileCount != 1 || !got[0].ModuleLocal || !got[0].UserMeaningful {
		t.Fatalf("main type = %#v", got[0])
	}
	if got[1].Package != "runtime" || got[1].ImportPath != "runtime" || got[1].SourceFileCount != 0 || got[1].ModuleLocal || got[1].UserMeaningful {
		t.Fatalf("runtime type = %#v", got[1])
	}
	if got[2].Package != "example.com/gorevealfixture/pkg/sub" || got[2].ImportPath != "example.com/gorevealfixture/pkg/sub" || got[2].SourceFileCount != 2 || !got[2].ModuleLocal || !got[2].UserMeaningful {
		t.Fatalf("sub type = %#v", got[2])
	}
}

func TestEnrichBuildInfoMetadata(t *testing.T) {
	t.Parallel()

	got := EnrichBuildInfoMetadata([]Type{
		{Name: "main.fixtureCounter", Kind: "struct"},
		{Name: "**runtime.g", Kind: "pointer"},
		{Name: "[]example.com/gorevealfixture/pkg/sub.Counter", Kind: "struct"},
	}, &schema.BuildInfo{Path: "example.com/gorevealfixture"})

	if got[0].Package != "main" || got[0].ImportPath != "example.com/gorevealfixture" || !got[0].ModuleLocal || !got[0].UserMeaningful {
		t.Fatalf("main type = %#v", got[0])
	}
	if got[1].Package != "runtime" || got[1].ImportPath != "runtime" || got[1].ModuleLocal || got[1].UserMeaningful {
		t.Fatalf("runtime type = %#v", got[1])
	}
	if got[2].Package != "example.com/gorevealfixture/pkg/sub" || got[2].ImportPath != "example.com/gorevealfixture/pkg/sub" || got[2].ModuleLocal || got[2].UserMeaningful {
		t.Fatalf("sub type = %#v", got[2])
	}
}

func TestTypePackage(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"main.fixtureCounter":                     "main",
		"**runtime.g":                             "runtime",
		"[]example.com/gorevealfixture/pkg/sub.T": "example.com/gorevealfixture/pkg/sub",
		"[4]sync.Pool":                            "sync",
		"string":                                  "",
	}

	for input, want := range cases {
		if got := typePackage(input); got != want {
			t.Fatalf("typePackage(%q) = %q, want %q", input, got, want)
		}
	}
}
