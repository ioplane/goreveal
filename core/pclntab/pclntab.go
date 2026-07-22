package pclntab

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	stdpe "debug/pe"
	"errors"
	"fmt"
	"os"

	binaryformat "github.com/dantte-lp/goreveal/core/format"
	"github.com/dantte-lp/goreveal/core/recoveryerr"
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
	kind, err := binaryformat.DetectFile(path)
	if err != nil {
		return Data{}, fmt.Errorf("detect pclntab container: %w", err)
	}
	if kind == binaryformat.Unknown {
		return Data{}, errors.New("detect pclntab container: unknown binary format")
	}

	fh, err := os.Open(path)
	if err != nil {
		return Data{}, fmt.Errorf("open file: %w", err)
	}
	defer fh.Close()

	switch kind {
	case binaryformat.ELF:
		ef, openErr := elf.NewFile(fh)
		if openErr != nil {
			return Data{}, fmt.Errorf("open ELF: %w", openErr)
		}
		return readELF(ef)
	case binaryformat.PE:
		pf, openErr := stdpe.NewFile(fh)
		if openErr != nil {
			return Data{}, fmt.Errorf("open PE: %w", openErr)
		}
		return readPE(pf)
	case binaryformat.MachO:
		mf, openErr := macho.NewFile(fh)
		if openErr != nil {
			return Data{}, fmt.Errorf("open Mach-O: %w", openErr)
		}
		return readMachO(mf)
	default:
		return Data{}, fmt.Errorf("detect pclntab container: unhandled format %q", kind)
	}
}

func readELF(ef *elf.File) (Data, error) {
	pclnSection := ef.Section(".gopclntab")
	if pclnSection == nil {
		return Data{}, recoveryerr.NewUnavailable(
			recoveryerr.CodePclntabNotFound,
			"ELF pclntab section is absent",
			nil,
		)
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
		return Data{}, recoveryerr.NewUnavailable(
			recoveryerr.CodePclntabNotFound,
			"Mach-O pclntab section is absent",
			nil,
		)
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
		return Data{}, recoveryerr.NewUnavailable(
			recoveryerr.CodePclntabNotFound,
			"PE pclntab header candidate is absent",
			nil,
		)
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
