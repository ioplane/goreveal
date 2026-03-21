package recoverystrings

import (
	"debug/elf"
	"fmt"
	"os"

	"github.com/dantte-lp/goreveal/schema"
)

// Region is a scan region used for string candidate extraction.
type Region = schema.StringRegion

func readRegions(path string) ([]Region, map[string][]byte, error) {
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
