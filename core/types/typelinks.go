package types

import (
	"debug/dwarf"
	"debug/elf"
	"errors"
	"fmt"
	"os"

	binaryformat "github.com/dantte-lp/goreveal/core/format"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
)

func openDWARF(path string) (*dwarf.Data, error) {
	kind, err := binaryformat.DetectFile(path)
	if err != nil {
		return nil, fmt.Errorf("detect DWARF container: %w", err)
	}
	switch kind {
	case binaryformat.PE, binaryformat.MachO:
		return nil, recoveryerr.NewUnsupported(
			recoveryerr.CodeDWARFUnsupportedContainer,
			fmt.Sprintf("DWARF type recovery does not support %s", kind),
			nil,
		)
	case binaryformat.Unknown:
		return nil, errors.New("detect DWARF container: unknown binary format")
	case binaryformat.ELF:
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err != nil {
		return nil, fmt.Errorf("open ELF: %w", err)
	}
	if section := ef.Section(".debug_info"); section == nil || section.Size == 0 {
		return nil, recoveryerr.NewUnavailable(
			recoveryerr.CodeDWARFNotFound,
			"ELF DWARF type evidence is absent",
			nil,
		)
	}

	data, err := ef.DWARF()
	if err != nil {
		return nil, fmt.Errorf("load DWARF: %w", err)
	}

	return data, nil
}
