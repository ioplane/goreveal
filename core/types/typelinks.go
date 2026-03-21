package types

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"os"
)

func openDWARF(path string) (*dwarf.Data, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err != nil {
		return nil, fmt.Errorf("open ELF: %w", err)
	}

	data, err := ef.DWARF()
	if err != nil {
		return nil, fmt.Errorf("load DWARF: %w", err)
	}

	return data, nil
}
