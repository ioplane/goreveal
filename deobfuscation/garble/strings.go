package garble

import (
	"context"
	"regexp"
	"slices"
	"strings"

	"github.com/ioplane/goreveal/schema"
)

var usefulSegmentPattern = regexp.MustCompile(`[A-Za-z][A-Za-z0-9./_-]*(?: [A-Za-z][A-Za-z0-9./_-]*)*`)

// Pass extracts useful human-readable segments from larger raw string candidates.
type Pass struct{}

func (Pass) Name() string { return "string-segments" }

func (Pass) Apply(_ context.Context, analysis schema.Analysis, refined *schema.RefinedAnalysis) error {
	seen := make(map[string]struct{}, len(refined.Strings))
	for _, str := range refined.Strings {
		seen[str.Value] = struct{}{}
	}

	for _, raw := range analysis.Strings {
		for _, match := range usefulSegmentPattern.FindAllString(raw.Value, -1) {
			match = normalizeSegment(match)
			if len(match) < 4 {
				continue
			}
			if _, ok := seen[match]; ok {
				continue
			}
			seen[match] = struct{}{}
			refined.Strings = append(refined.Strings, schema.RefinedString{Value: match})
		}
	}

	slices.SortFunc(refined.Strings, func(a, b schema.RefinedString) int {
		if a.Value < b.Value {
			return -1
		}
		if a.Value > b.Value {
			return 1
		}
		return 0
	})

	return nil
}

func normalizeSegment(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimRight(value, "0123456789")
	return strings.TrimSpace(value)
}
