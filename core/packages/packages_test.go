package packages

import (
	"testing"

	"github.com/ioplane/goreveal/core/functions"
)

func TestRecoverFixturePackages(t *testing.T) {
	t.Parallel()

	funcs, err := functions.Recover("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("functions.Recover() error = %v", err)
	}

	pkgs := Recover(funcs)
	if len(pkgs) == 0 {
		t.Fatal("Recover() returned no packages")
	}

	for _, want := range []string{"main", "runtime"} {
		if !containsPackage(pkgs, want) {
			t.Fatalf("Recover() missing package %q", want)
		}
	}

	runtimePkg, ok := findPackage(pkgs, "runtime")
	if !ok {
		t.Fatal("Recover() missing runtime package")
	}
	if runtimePkg.ImportPath != "runtime" {
		t.Fatalf("runtime import path = %q, want %q", runtimePkg.ImportPath, "runtime")
	}

	mainPkg, ok := findPackage(pkgs, "main")
	if !ok {
		t.Fatal("Recover() missing main package")
	}
	if mainPkg.ImportPath != "" {
		t.Fatalf("main import path = %q, want empty", mainPkg.ImportPath)
	}
}

func TestRecoverMachOFixturePackages(t *testing.T) {
	t.Parallel()

	funcs, err := functions.Recover("../../corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("functions.Recover() error = %v", err)
	}

	pkgs := Recover(funcs)
	if len(pkgs) == 0 {
		t.Fatal("Recover() returned no packages")
	}

	for _, want := range []string{"main", "runtime"} {
		if !containsPackage(pkgs, want) {
			t.Fatalf("Recover() missing package %q", want)
		}
	}
}

func TestRecoverPEFixturePackages(t *testing.T) {
	t.Parallel()

	funcs, err := functions.Recover("../../corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe")
	if err != nil {
		t.Fatalf("functions.Recover() error = %v", err)
	}

	pkgs := Recover(funcs)
	if len(pkgs) == 0 {
		t.Fatal("Recover() returned no packages")
	}

	for _, want := range []string{"main", "runtime"} {
		if !containsPackage(pkgs, want) {
			t.Fatalf("Recover() missing package %q", want)
		}
	}
}

func containsPackage(pkgs []Package, want string) bool {
	for _, pkg := range pkgs {
		if pkg.Name == want && pkg.FunctionCount > 0 {
			return true
		}
	}
	return false
}

func findPackage(pkgs []Package, want string) (Package, bool) {
	for _, pkg := range pkgs {
		if pkg.Name == want {
			return pkg, true
		}
	}
	return Package{}, false
}
