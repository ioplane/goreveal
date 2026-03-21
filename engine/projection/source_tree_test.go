package projection

import (
	"testing"

	"github.com/dantte-lp/goreveal/schema"
)

func TestBuildSourceTree(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/gorevealfixture",
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 3},
		},
	}

	got, err := BuildSourceTree(
		analysis,
		[]string{
			"runtime/proc.go",
			"example.com/gorevealfixture/main.go",
		},
	)
	if err != nil {
		t.Fatalf("BuildSourceTree() error = %v", err)
	}

	if got.Root != "example.com/gorevealfixture" {
		t.Fatalf("BuildSourceTree() root = %q", got.Root)
	}
	if len(got.Files) != 1 || got.Files[0] != "main.go" {
		t.Fatalf("BuildSourceTree() files = %#v", got.Files)
	}
	if len(got.Packages) != 1 {
		t.Fatalf("BuildSourceTree() packages = %#v", got.Packages)
	}
	if got.Packages[0].ImportPath != "example.com/gorevealfixture" {
		t.Fatalf("BuildSourceTree() import path = %q", got.Packages[0].ImportPath)
	}
	if got.Packages[0].Name != "main" {
		t.Fatalf("BuildSourceTree() package name = %q", got.Packages[0].Name)
	}
	if got.Packages[0].FunctionCount != 3 {
		t.Fatalf("BuildSourceTree() function count = %d", got.Packages[0].FunctionCount)
	}
	if !got.Packages[0].HasFileEvidence {
		t.Fatalf("BuildSourceTree() package evidence = %#v", got.Packages[0])
	}
}

func TestBuildSourceTreeAggregatesFilesByImportPath(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/gorevealfixture",
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 3},
			{Name: "example.com/gorevealfixture/pkg/sub", FunctionCount: 2},
		},
	}

	got, err := BuildSourceTree(
		analysis,
		[]string{
			"example.com/gorevealfixture/main.go",
			"example.com/gorevealfixture/pkg/sub/a.go",
			"example.com/gorevealfixture/pkg/sub/b.go",
		},
	)
	if err != nil {
		t.Fatalf("BuildSourceTree() error = %v", err)
	}

	if len(got.Packages) != 2 {
		t.Fatalf("BuildSourceTree() packages = %#v", got.Packages)
	}

	mainPkg := got.Packages[0]
	subPkg := got.Packages[1]

	if mainPkg.ImportPath != "example.com/gorevealfixture" || len(mainPkg.Files) != 1 || mainPkg.Files[0] != "main.go" {
		t.Fatalf("main package = %#v", mainPkg)
	}
	if subPkg.ImportPath != "example.com/gorevealfixture/pkg/sub" {
		t.Fatalf("sub package import path = %#v", subPkg)
	}
	if subPkg.Name != "sub" {
		t.Fatalf("sub package name = %#v", subPkg)
	}
	if subPkg.FunctionCount != 2 {
		t.Fatalf("sub package function count = %#v", subPkg)
	}
	if len(subPkg.Files) != 2 || subPkg.Files[0] != "pkg/sub/a.go" || subPkg.Files[1] != "pkg/sub/b.go" {
		t.Fatalf("sub package files = %#v", subPkg.Files)
	}
}

func TestBuildSourceTreeSeparatesExternalPackages(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/gorevealfixture",
		},
		Packages: []schema.Package{
			{Name: "main", FunctionCount: 3},
		},
	}

	got, err := BuildSourceTree(
		analysis,
		[]string{
			"/usr/local/go/src/runtime/proc.go",
			"/usr/local/go/src/runtime/mgc.go",
			"example.com/gorevealfixture/main.go",
		},
	)
	if err != nil {
		t.Fatalf("BuildSourceTree() error = %v", err)
	}

	if len(got.Files) != 1 || got.Files[0] != "main.go" {
		t.Fatalf("module files = %#v", got.Files)
	}
	if len(got.ExternalPackages) != 1 {
		t.Fatalf("external packages = %#v", got.ExternalPackages)
	}

	runtimePkg := got.ExternalPackages[0]
	if runtimePkg.ImportPath != "runtime" {
		t.Fatalf("runtime package import path = %#v", runtimePkg)
	}
	if runtimePkg.Name != "runtime" {
		t.Fatalf("runtime package name = %#v", runtimePkg)
	}
	if len(runtimePkg.Files) != 2 ||
		runtimePkg.Files[0] != "/usr/local/go/src/runtime/mgc.go" ||
		runtimePkg.Files[1] != "/usr/local/go/src/runtime/proc.go" {
		t.Fatalf("runtime package files = %#v", runtimePkg.Files)
	}
	if !runtimePkg.HasFileEvidence {
		t.Fatalf("runtime package evidence = %#v", runtimePkg)
	}
}

func TestBuildFallbackSourceTree(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		BuildInfo: &schema.BuildInfo{
			Path: "example.com/gorevealfixture",
		},
		Packages: []schema.Package{
			{Name: "main", ImportPath: "example.com/gorevealfixture", FunctionCount: 3, ModuleLocal: true},
			{Name: "runtime", FunctionCount: 10},
			{Name: "fmt", FunctionCount: 4},
		},
	}

	got, err := BuildFallbackSourceTree(analysis)
	if err != nil {
		t.Fatalf("BuildFallbackSourceTree() error = %v", err)
	}

	if got.Root != "example.com/gorevealfixture" {
		t.Fatalf("BuildFallbackSourceTree() root = %q", got.Root)
	}
	if len(got.Files) != 0 {
		t.Fatalf("BuildFallbackSourceTree() files = %#v", got.Files)
	}
	if len(got.Packages) != 1 {
		t.Fatalf("BuildFallbackSourceTree() packages = %#v", got.Packages)
	}
	if len(got.ExternalPackages) != 2 {
		t.Fatalf("BuildFallbackSourceTree() external packages = %#v", got.ExternalPackages)
	}

	mainPkg := got.Packages[0]
	if mainPkg.Name != "main" || mainPkg.ImportPath != "example.com/gorevealfixture" || mainPkg.FunctionCount != 3 {
		t.Fatalf("BuildFallbackSourceTree() main package = %#v", mainPkg)
	}
	if len(mainPkg.Files) != 0 {
		t.Fatalf("BuildFallbackSourceTree() main package files = %#v", mainPkg.Files)
	}
	if mainPkg.HasFileEvidence {
		t.Fatalf("BuildFallbackSourceTree() main package evidence = %#v", mainPkg)
	}

	runtimePkg := got.ExternalPackages[1]
	if runtimePkg.Name != "runtime" || runtimePkg.ImportPath != "runtime" || runtimePkg.FunctionCount != 10 {
		t.Fatalf("BuildFallbackSourceTree() runtime package = %#v", runtimePkg)
	}
	if len(runtimePkg.Files) != 0 {
		t.Fatalf("BuildFallbackSourceTree() runtime package files = %#v", runtimePkg.Files)
	}
	if runtimePkg.HasFileEvidence {
		t.Fatalf("BuildFallbackSourceTree() runtime package evidence = %#v", runtimePkg)
	}

	fmtPkg := got.ExternalPackages[0]
	if fmtPkg.Name != "fmt" || fmtPkg.ImportPath != "fmt" || fmtPkg.FunctionCount != 4 {
		t.Fatalf("BuildFallbackSourceTree() fmt package = %#v", fmtPkg)
	}
	if len(fmtPkg.Files) != 0 {
		t.Fatalf("BuildFallbackSourceTree() fmt package files = %#v", fmtPkg.Files)
	}
	if fmtPkg.HasFileEvidence {
		t.Fatalf("BuildFallbackSourceTree() fmt package evidence = %#v", fmtPkg)
	}
}
