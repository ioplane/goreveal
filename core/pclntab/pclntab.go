package pclntab

import (
	"debug/elf"
	"fmt"
	"os"

	"github.com/dantte-lp/goreveal/schema"
)

// Data captures the minimal pclntab-related inputs needed for function recovery.
type Data struct {
	PCLN       []byte
	Symtab     []byte
	TextStart  uint64
	Provenance schema.Provenance
}

// Read extracts pclntab-related sections from a Go ELF binary.
func Read(path string) (Data, error) {
	fh, err := os.Open(path)
	if err != nil {
		return Data{}, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err != nil {
		return Data{}, fmt.Errorf("open ELF: %w", err)
	}

	pclnSection := ef.Section(".gopclntab")
	if pclnSection == nil {
		return Data{}, fmt.Errorf("section %q not found", ".gopclntab")
	}
	pcln, err := pclnSection.Data()
	if err != nil {
		return Data{}, fmt.Errorf("read %s: %w", pclnSection.Name, err)
	}

	textSection := ef.Section(".text")
	if textSection == nil {
		return Data{}, fmt.Errorf("section %q not found", ".text")
	}

	var symtab []byte
	if section := ef.Section(".gosymtab"); section != nil {
		symtab, err = section.Data()
		if err != nil {
			return Data{}, fmt.Errorf("read %s: %w", section.Name, err)
		}
	}

	return Data{
		PCLN:      pcln,
		Symtab:    symtab,
		TextStart: textSection.Addr,
		Provenance: schema.Provenance{
			Source:     "core.pclntab",
			Confidence: "high",
		},
	}, nil
}
