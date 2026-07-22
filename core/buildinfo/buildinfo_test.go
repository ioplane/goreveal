package buildinfo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dantte-lp/goreveal/core/recoveryerr"
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

func TestReadNonGoELFIsUnavailable(t *testing.T) {
	t.Parallel()

	_, err := Read("/bin/sh")
	if !errors.Is(err, recoveryerr.ErrUnavailable) {
		t.Fatalf("Read(/bin/sh) error = %v, want unavailable", err)
	}
}

func TestReadMalformedContainerRemainsFailure(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "malformed.bin")
	if err := os.WriteFile(path, []byte{0x7f, 'E', 'L', 'F', 0x02}, 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	_, err := Read(path)
	if err == nil {
		t.Fatal("Read() error = nil, want failure")
	}
	if errors.Is(err, recoveryerr.ErrUnavailable) || errors.Is(err, recoveryerr.ErrUnsupported) {
		t.Fatalf("Read() error = %v, want ordinary failure", err)
	}
}
