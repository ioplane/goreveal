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

func TestReadMachOFixture(t *testing.T) {
	t.Parallel()

	data, err := Read("../../corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin")
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

func TestReadPEFixture(t *testing.T) {
	t.Parallel()

	data, err := Read("../../corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe")
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
