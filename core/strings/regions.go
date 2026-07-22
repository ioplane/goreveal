package recoverystrings

import (
	"debug/elf"
	"errors"
	"fmt"
	"os"

	binaryformat "github.com/dantte-lp/goreveal/core/format"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
	"github.com/dantte-lp/goreveal/schema"
)

// Region is a scan region used for string candidate extraction.
type Region = schema.StringRegion

func readRegions(path string) ([]Region, map[string][]byte, error) {
	kind, err := binaryformat.DetectFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("detect string-region container: %w", err)
	}
	switch kind {
	case binaryformat.PE, binaryformat.MachO:
		return nil, nil, recoveryerr.NewUnsupported(
			recoveryerr.CodeStringRegionsUnsupportedContainer,
			fmt.Sprintf("string region recovery does not support %s", kind),
			nil,
		)
	case binaryformat.Unknown:
		return nil, nil, errors.New("detect string-region container: unknown binary format")
	case binaryformat.ELF:
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err != nil {
		return nil, nil, fmt.Errorf("open ELF: %w", err)
	}

	regions := make([]Region, 0, len(ef.Sections))
	dataBySection := make(map[string][]byte)
	for _, section := range ef.Sections {
		if !scanSection(section) {
			continue
		}

		data, err := section.Data()
		if err != nil {
			return nil, nil, fmt.Errorf("read %s: %w", section.Name, err)
		}
		if len(data) == 0 {
			continue
		}

		regions = append(regions, Region{
			Name:   section.Name,
			Addr:   section.Addr,
			Size:   section.Size,
			Source: "elf.section",
		})
		dataBySection[section.Name] = data
	}
	if len(regions) == 0 {
		return nil, nil, recoveryerr.NewUnavailable(
			recoveryerr.CodeStringRegionsNotFound,
			"ELF string scan regions are absent",
			nil,
		)
	}

	return regions, dataBySection, nil
}

func scanSection(section *elf.Section) bool {
	switch section.Name {
	case ".rodata", ".noptrdata", ".data", ".go.buildinfo":
		return true
	default:
		return false
	}
}
