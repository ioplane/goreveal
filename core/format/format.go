package format

// Kind is the detected binary container format.
type Kind string

const (
	Unknown Kind = "unknown"
	ELF     Kind = "elf"
	PE      Kind = "pe"
	MachO   Kind = "macho"
)

// DetectHeader classifies a binary by its leading magic bytes.
func DetectHeader(header []byte) Kind {
	if len(header) >= 4 {
		if header[0] == 0x7f && header[1] == 'E' && header[2] == 'L' && header[3] == 'F' {
			return ELF
		}
		if header[0] == 'M' && header[1] == 'Z' {
			return PE
		}
	}

	if len(header) >= 4 {
		switch {
		case header[0] == 0xfe && header[1] == 0xed && header[2] == 0xfa && header[3] == 0xce:
			return MachO
		case header[0] == 0xce && header[1] == 0xfa && header[2] == 0xed && header[3] == 0xfe:
			return MachO
		case header[0] == 0xfe && header[1] == 0xed && header[2] == 0xfa && header[3] == 0xcf:
			return MachO
		case header[0] == 0xcf && header[1] == 0xfa && header[2] == 0xed && header[3] == 0xfe:
			return MachO
		case header[0] == 0xca && header[1] == 0xfe && header[2] == 0xba && header[3] == 0xbe:
			return MachO
		}
	}

	return Unknown
}
