package types

import (
	"errors"
	"testing"

	"github.com/dantte-lp/goreveal/core/recoveryerr"
)

func TestRecoverFixtureTypes(t *testing.T) {
	t.Parallel()

	recovered, err := Recover("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("Recover() error = %v", err)
	}

	if len(recovered) == 0 {
		t.Fatal("Recover() returned no types")
	}

	for _, want := range []string{"main.fixtureCounter", "main.fixtureGreeterImpl", "int", "[]string"} {
		if !containsType(recovered, want) {
			t.Fatalf("Recover() missing type %q", want)
		}
	}
}

func TestRecoverStrippedFixtureIsUnavailable(t *testing.T) {
	t.Parallel()

	_, err := Recover("../../corpus/fixtures/go-elf-stripped-linux-amd64/fixture.bin")
	if !errors.Is(err, recoveryerr.ErrUnavailable) {
		t.Fatalf("Recover(stripped ELF) error = %v, want unavailable", err)
	}
}

func TestRecoverMachOIsUnsupported(t *testing.T) {
	t.Parallel()

	_, err := Recover("../../corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin")
	if !errors.Is(err, recoveryerr.ErrUnsupported) {
		t.Fatalf("Recover(Mach-O) error = %v, want unsupported", err)
	}
}

func containsType(types []Type, want string) bool {
	for _, typ := range types {
		if typ.Name == want && typ.Kind != "" {
			return true
		}
	}
	return false
}
