package pclntab

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	stdpe "debug/pe"
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

// Read extracts pclntab-related sections from a supported Go binary container.
func Read(path string) (Data, error) {
	fh, err := os.Open(path)
	if err != nil {
		return Data{}, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	ef, err := elf.NewFile(fh)
	if err == nil {
		return readELF(ef)
	}

	if _, seekErr := fh.Seek(0, 0); seekErr != nil {
		return Data{}, fmt.Errorf("rewind file: %w", seekErr)
	}

	pf, err := stdpe.NewFile(fh)
	if err == nil {
		return readPE(pf)
	}

	if _, seekErr := fh.Seek(0, 0); seekErr != nil {
		return Data{}, fmt.Errorf("rewind file: %w", seekErr)
	}

	mf, err := macho.NewFile(fh)
	if err == nil {
		return readMachO(mf)
	}

	return Data{}, fmt.Errorf("open binary: unsupported pclntab container")
}

func readELF(ef *elf.File) (Data, error) {
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

func readMachO(mf *macho.File) (Data, error) {
	pclnSection := mf.Section("__gopclntab")
	if pclnSection == nil {
		return Data{}, fmt.Errorf("section %q not found", "__gopclntab")
	}
	pcln, err := pclnSection.Data()
	if err != nil {
		return Data{}, fmt.Errorf("read %s: %w", pclnSection.Name, err)
	}

	textSection := mf.Section("__text")
	if textSection == nil {
		return Data{}, fmt.Errorf("section %q not found", "__text")
	}

	var symtab []byte
	if section := mf.Section("__gosymtab"); section != nil {
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

func readPE(pf *stdpe.File) (Data, error) {
	textSection := peSection(pf, ".text")
	if textSection == nil {
		return Data{}, fmt.Errorf("section %q not found", ".text")
	}

	var pcln []byte
	for _, sectionName := range []string{".rdata", ".text"} {
		section := peSection(pf, sectionName)
		if section == nil {
			continue
		}

		data, err := section.Data()
		if err != nil {
			return Data{}, fmt.Errorf("read %s: %w", section.Name, err)
		}

		offset, ok := findPEPclntabHeaderCandidate(data)
		if !ok {
			continue
		}

		pcln = data[offset:]
		break
	}

	if len(pcln) == 0 {
		return Data{}, fmt.Errorf("PE pclntab header candidate not found")
	}

	return Data{
		PCLN:      pcln,
		TextStart: peImageBase(pf) + uint64(textSection.VirtualAddress),
		Provenance: schema.Provenance{
			Source:     "core.pclntab",
			Confidence: "high",
		},
	}, nil
}

func peSection(pf *stdpe.File, name string) *stdpe.Section {
	for _, section := range pf.Sections {
		if section.Name == name {
			return section
		}
	}

	return nil
}

func peImageBase(pf *stdpe.File) uint64 {
	switch header := pf.OptionalHeader.(type) {
	case *stdpe.OptionalHeader32:
		return uint64(header.ImageBase)
	case *stdpe.OptionalHeader64:
		return header.ImageBase
	default:
		return 0
	}
}

func findPEPclntabHeaderCandidate(data []byte) (int, bool) {
	for _, magic := range [][]byte{
		{0xf0, 0xff, 0xff, 0xff},
		{0xf1, 0xff, 0xff, 0xff},
		{0xfa, 0xff, 0xff, 0xff},
		{0xfb, 0xff, 0xff, 0xff},
	} {
		start := 0
		for {
			index := bytes.Index(data[start:], magic)
			if index < 0 {
				break
			}

			absolute := start + index
			if isPEPclntabHeaderCandidate(data, absolute) {
				return absolute, true
			}

			start = absolute + 1
		}
	}

	return 0, false
}

func isPEPclntabHeaderCandidate(data []byte, offset int) bool {
	if len(data) < offset+8 {
		return false
	}
	if data[offset+4] != 0 || data[offset+5] != 0 {
		return false
	}

	switch data[offset+6] {
	case 1, 2, 4:
	default:
		return false
	}

	switch data[offset+7] {
	case 4, 8:
		return true
	default:
		return false
	}
}
