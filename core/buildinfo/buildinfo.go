package buildinfo

import (
	stdbuildinfo "debug/buildinfo"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"os"
	"strings"

	binaryformat "github.com/dantte-lp/goreveal/core/format"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
	"github.com/dantte-lp/goreveal/schema"
)

// Info is the canonical build information recovered from a Go binary.
type Info = schema.BuildInfo

// Read extracts build information directly from the binary.
func Read(path string) (Info, error) {
	info, err := stdbuildinfo.ReadFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "not a Go executable") && hasValidContainer(path) {
			return Info{}, recoveryerr.NewUnavailable(
				recoveryerr.CodeBuildInfoNotFound,
				"Go build info is absent",
				err,
			)
		}
		return Info{}, fmt.Errorf("read build info: %w", err)
	}

	return Info{
		GoVersion: info.GoVersion,
		Path:      info.Path,
		Provenance: schema.Provenance{
			Source:     "core.buildinfo",
			Confidence: "high",
		},
	}, nil
}

func hasValidContainer(path string) bool {
	kind, err := binaryformat.DetectFile(path)
	if err != nil {
		return false
	}

	fh, err := os.Open(path)
	if err != nil {
		return false
	}
	defer fh.Close()

	switch kind {
	case binaryformat.ELF:
		_, err = elf.NewFile(fh)
	case binaryformat.PE:
		_, err = pe.NewFile(fh)
	case binaryformat.MachO:
		_, err = macho.NewFile(fh)
	case binaryformat.Unknown:
		return false
	}

	return err == nil
}
