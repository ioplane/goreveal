package types

import "testing"

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

func containsType(types []Type, want string) bool {
	for _, typ := range types {
		if typ.Name == want && typ.Kind != "" {
			return true
		}
	}
	return false
}
