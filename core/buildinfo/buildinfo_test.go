package buildinfo

import (
	"strings"
	"testing"
)

func TestReadFixture(t *testing.T) {
	t.Parallel()

	info, err := Read("../../corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if !strings.HasPrefix(info.GoVersion, "go1.") {
		t.Fatalf("Read() GoVersion = %q, want prefix %q", info.GoVersion, "go1.")
	}
	if info.Path != "example.com/gorevealfixture" {
		t.Fatalf("Read() Path = %q, want %q", info.Path, "example.com/gorevealfixture")
	}
	if info.Provenance.Source != "core.buildinfo" {
		t.Fatalf("Read() provenance source = %q", info.Provenance.Source)
	}
	if info.Provenance.Confidence != "high" {
		t.Fatalf("Read() provenance confidence = %q", info.Provenance.Confidence)
	}
}

func TestReadPEFixture(t *testing.T) {
	t.Parallel()

	info, err := Read("../../corpus/fixtures/go-pe-buildinfo-windows-amd64/fixture.exe")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if !strings.HasPrefix(info.GoVersion, "go1.") {
		t.Fatalf("Read() GoVersion = %q, want prefix %q", info.GoVersion, "go1.")
	}
	if info.Path != "example.com/gorevealfixture" {
		t.Fatalf("Read() Path = %q, want %q", info.Path, "example.com/gorevealfixture")
	}
	if info.Provenance.Source != "core.buildinfo" {
		t.Fatalf("Read() provenance source = %q", info.Provenance.Source)
	}
	if info.Provenance.Confidence != "high" {
		t.Fatalf("Read() provenance confidence = %q", info.Provenance.Confidence)
	}
}

func TestReadMachOFixture(t *testing.T) {
	t.Parallel()

	info, err := Read("../../corpus/fixtures/go-macho-buildinfo-darwin-amd64/fixture.bin")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if !strings.HasPrefix(info.GoVersion, "go1.") {
		t.Fatalf("Read() GoVersion = %q, want prefix %q", info.GoVersion, "go1.")
	}
	if info.Path != "example.com/gorevealfixture" {
		t.Fatalf("Read() Path = %q, want %q", info.Path, "example.com/gorevealfixture")
	}
	if info.Provenance.Source != "core.buildinfo" {
		t.Fatalf("Read() provenance source = %q", info.Provenance.Source)
	}
	if info.Provenance.Confidence != "high" {
		t.Fatalf("Read() provenance confidence = %q", info.Provenance.Confidence)
	}
}
