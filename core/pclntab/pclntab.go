package pclntab

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	stdpe "debug/pe"
	"errors"
	"fmt"
	"os"

	"github.com/ioplane/goreveal/schema"
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

	return Data{}, errors.New("open binary: unsupported pclntab container")
}

// sectionNames is the per-format naming of the three sections pclntab recovery
// needs. ELF and Mach-O differ only in these names, not in the recovery itself.
type sectionNames struct {
	pcln   string
	text   string
	symtab string
}

var (
	elfSectionNames   = sectionNames{pcln: ".gopclntab", text: ".text", symtab: ".gosymtab"}
	machoSectionNames = sectionNames{pcln: "__gopclntab", text: "__text", symtab: "__gosymtab"}
)

// sectionLookup resolves a section name to its bytes and virtual address.
// It returns found=false when the section is simply absent, which the caller
// distinguishes from a read failure.
type sectionLookup func(name string) (data []byte, addr uint64, found bool, err error)

func elfSectionLookup(ef *elf.File) sectionLookup {
	return func(name string) ([]byte, uint64, bool, error) {
		section := ef.Section(name)
		if section == nil {
			return nil, 0, false, nil
		}
		data, err := section.Data()
		if err != nil {
			return nil, 0, true, fmt.Errorf("read %s: %w", section.Name, err)
		}
		return data, section.Addr, true, nil
	}
}

func machoSectionLookup(mf *macho.File) sectionLookup {
	return func(name string) ([]byte, uint64, bool, error) {
		section := mf.Section(name)
		if section == nil {
			return nil, 0, false, nil
		}
		data, err := section.Data()
		if err != nil {
			return nil, 0, true, fmt.Errorf("read %s: %w", section.Name, err)
		}
		return data, section.Addr, true, nil
	}
}

// readSectioned recovers pclntab data from any format whose sections can be
// resolved by name. The symbol table stays optional: a binary without
// .gosymtab / __gosymtab is normal, not an error.
func readSectioned(lookup sectionLookup, names sectionNames) (Data, error) {
	pcln, _, found, err := lookup(names.pcln)
	if err != nil {
		return Data{}, err
	}
	if !found {
		return Data{}, fmt.Errorf("section %q not found", names.pcln)
	}

	_, textAddr, found, err := lookup(names.text)
	if err != nil {
		return Data{}, err
	}
	if !found {
		return Data{}, fmt.Errorf("section %q not found", names.text)
	}

	symtab, _, _, err := lookup(names.symtab)
	if err != nil {
		return Data{}, err
	}

	return Data{
		PCLN:      pcln,
		Symtab:    symtab,
		TextStart: textAddr,
		Provenance: schema.Provenance{
			Source:     "core.pclntab",
			Confidence: "high",
		},
	}, nil
}

func readELF(ef *elf.File) (Data, error) {
	return readSectioned(elfSectionLookup(ef), elfSectionNames)
}

func readMachO(mf *macho.File) (Data, error) {
	return readSectioned(machoSectionLookup(mf), machoSectionNames)
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
		return Data{}, errors.New("PE pclntab header candidate not found")
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
