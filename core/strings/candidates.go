package recoverystrings

import (
	"fmt"
	"slices"

	"github.com/dantte-lp/goreveal/schema"
)

// Candidate is a canonical recovered string candidate.
type Candidate = schema.StringCandidate

// Result groups string scan regions and recovered candidates.
type Result struct {
	Regions    []Region
	Candidates []Candidate
}

// Recover scans ELF data sections for printable ASCII strings.
func Recover(path string) (Result, error) {
	regions, dataBySection, err := readRegions(path)
	if err != nil {
		return Result{}, fmt.Errorf("read string regions: %w", err)
	}

	candidates := make([]Candidate, 0, 128)
	seen := make(map[string]struct{})
	for _, region := range regions {
		data := dataBySection[region.Name]
		start := -1
		for i, b := range data {
			if isPrintableASCII(b) {
				if start == -1 {
					start = i
				}
				continue
			}
			candidates = flushCandidate(candidates, seen, region, data, start, i)
			start = -1
		}
		candidates = flushCandidate(candidates, seen, region, data, start, len(data))
	}

	slices.SortFunc(candidates, func(a, b Candidate) int {
		if a.Value < b.Value {
			return -1
		}
		if a.Value > b.Value {
			return 1
		}
		if a.Region < b.Region {
			return -1
		}
		if a.Region > b.Region {
			return 1
		}
		if a.Offset < b.Offset {
			return -1
		}
		if a.Offset > b.Offset {
			return 1
		}
		return 0
	})

	return Result{
		Regions:    regions,
		Candidates: candidates,
	}, nil
}

func flushCandidate(dst []Candidate, seen map[string]struct{}, region Region, data []byte, start, end int) []Candidate {
	if start == -1 || end-start < 4 {
		return dst
	}

	value := string(data[start:end])
	if !looksUseful(value) {
		return dst
	}
	key := region.Name + "\x00" + value
	if _, ok := seen[key]; ok {
		return dst
	}
	seen[key] = struct{}{}

	return append(dst, Candidate{
		Value:  value,
		Region: region.Name,
		Addr:   region.Addr + uint64(start),
		Offset: uint64(start),
		Provenance: schema.Provenance{
			Source:     "core.strings.elf",
			Confidence: "medium",
		},
	})
}

func isPrintableASCII(b byte) bool {
	return b >= 0x20 && b <= 0x7e
}

func looksUseful(value string) bool {
	hasLetter := false
	for _, b := range []byte(value) {
		if b >= 'A' && b <= 'Z' || b >= 'a' && b <= 'z' {
			hasLetter = true
			break
		}
	}
	return hasLetter
}
