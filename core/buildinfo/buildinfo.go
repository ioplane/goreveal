package buildinfo

import (
	stdbuildinfo "debug/buildinfo"
	"fmt"

	"github.com/ioplane/goreveal/schema"
)

// Info is the canonical build information recovered from a Go binary.
type Info = schema.BuildInfo

// Read extracts build information directly from the binary.
func Read(path string) (Info, error) {
	info, err := stdbuildinfo.ReadFile(path)
	if err != nil {
		return Info{}, fmt.Errorf("read build info: %w", err)
	}

	return Info{
		GoVersion: info.GoVersion,
		Path:      info.Path,
		Provenance: schema.Provenance{
			Source:     "core.buildinfo",
			Confidence: "high",
		},
	}, nil
}
