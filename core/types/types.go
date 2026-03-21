package types

import (
	"debug/dwarf"
	"fmt"
	"slices"

	"github.com/dantte-lp/goreveal/schema"
)

// Type is the canonical recovered type representation.
type Type = schema.Type

// Recover extracts named types from DWARF metadata.
func Recover(path string) ([]Type, error) {
	data, err := openDWARF(path)
	if err != nil {
		return nil, fmt.Errorf("open DWARF: %w", err)
	}

	reader := data.Reader()
	seen := make(map[string]Type)

	for {
		entry, err := reader.Next()
		if err != nil {
			return nil, fmt.Errorf("read DWARF entry: %w", err)
		}
		if entry == nil {
			break
		}

		name, _ := entry.Val(dwarf.AttrName).(string)
		if name == "" {
			continue
		}

		kind, ok := kindFromTag(entry.Tag)
		if !ok {
			continue
		}

		candidate := Type{
			Name: name,
			Kind: kind,
			Provenance: schema.Provenance{
				Source:     "core.types.dwarf",
				Confidence: "high",
			},
		}

		existing, exists := seen[name]
		if !exists || kindRank(candidate.Kind) > kindRank(existing.Kind) {
			seen[name] = candidate
		}
	}

	types := make([]Type, 0, len(seen))
	for _, typ := range seen {
		types = append(types, typ)
	}

	slices.SortFunc(types, func(a, b Type) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	return types, nil
}

func kindFromTag(tag dwarf.Tag) (string, bool) {
	switch tag {
	case dwarf.TagArrayType:
		return "array", true
	case dwarf.TagBaseType:
		return "base", true
	case dwarf.TagInterfaceType:
		return "interface", true
	case dwarf.TagPointerType:
		return "pointer", true
	case dwarf.TagStructType:
		return "struct", true
	case dwarf.TagSubroutineType:
		return "func", true
	case dwarf.TagTypedef:
		return "typedef", true
	default:
		return "", false
	}
}

func kindRank(kind string) int {
	switch kind {
	case "typedef":
		return 7
	case "interface":
		return 6
	case "struct":
		return 5
	case "array":
		return 4
	case "pointer":
		return 3
	case "base":
		return 2
	case "func":
		return 1
	default:
		return 0
	}
}
