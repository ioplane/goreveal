package pclntab

import "testing"

func TestReadELFFixture(t *testing.T) {
	t.Parallel()

	data, err := Read("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if len(data.PCLN) == 0 {
		t.Fatal("Read() returned empty pclntab")
	}
	if data.TextStart == 0 {
		t.Fatal("Read() returned zero text start")
	}
	if data.Provenance.Source != "core.pclntab" {
		t.Fatalf("Read() provenance source = %q", data.Provenance.Source)
	}
	if data.Provenance.Confidence != "high" {
		t.Fatalf("Read() provenance confidence = %q", data.Provenance.Confidence)
	}
}
