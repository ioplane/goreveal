package recoverystrings

import "testing"

func TestRecoverFixtureStrings(t *testing.T) {
	t.Parallel()

	got, err := Recover("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("Recover() error = %v", err)
	}

	if len(got.Candidates) == 0 {
		t.Fatal("Recover() returned no string candidates")
	}
	if len(got.Regions) == 0 {
		t.Fatal("Recover() returned no regions")
	}

	for _, want := range []string{"goreveal fixture", "example.com/gorevealfixture"} {
		if !containsCandidate(got.Candidates, want) {
			t.Fatalf("Recover() missing string candidate %q", want)
		}
	}

	candidate, ok := findCandidate(got.Candidates, "goreveal fixture")
	if !ok {
		t.Fatal("Recover() missing goreveal fixture candidate for addr assertions")
	}
	region, ok := findRegion(got.Regions, candidate.Region)
	if !ok {
		t.Fatalf("Recover() missing region %q for addr assertions", candidate.Region)
	}
	if candidate.Addr != region.Addr+candidate.Offset {
		t.Fatalf("candidate addr = %#x, want %#x", candidate.Addr, region.Addr+candidate.Offset)
	}
}

func containsCandidate(candidates []Candidate, want string) bool {
	for _, candidate := range candidates {
		if candidate.Value == want && candidate.Region != "" {
			return true
		}
	}
	return false
}

func findCandidate(candidates []Candidate, want string) (Candidate, bool) {
	for _, candidate := range candidates {
		if candidate.Value == want {
			return candidate, true
		}
	}

	return Candidate{}, false
}

func findRegion(regions []Region, want string) (Region, bool) {
	for _, region := range regions {
		if region.Name == want {
			return region, true
		}
	}

	return Region{}, false
}
