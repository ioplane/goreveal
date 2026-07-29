// Package version exposes the build identity of the goreveal binary.
//
// The three variables are populated at link time with -X. When they are not (a
// plain `go build` or `go install`), Info falls back to the VCS stamps that the
// Go toolchain embeds in the build info, so `goreveal version` stays truthful
// instead of reporting "dev/unknown" for a binary that does have provenance.
package version

import (
	"runtime"
	"runtime/debug"
)

// Injected via -ldflags="-X github.com/ioplane/goreveal/internal/version.<Name>=...".
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

const unknown = "unknown"

// Info is the machine-readable build identity reported by `goreveal version`.
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
	Modified  bool   `json:"modified"`
}

// Get resolves the build identity, preferring link-time values and falling back
// to the toolchain's embedded VCS stamps.
func Get() Info {
	info := Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}

	build, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	if info.Version == "dev" && build.Main.Version != "" && build.Main.Version != "(devel)" {
		info.Version = build.Main.Version
	}

	for _, setting := range build.Settings {
		switch setting.Key {
		case "vcs.revision":
			if info.GitCommit == unknown {
				info.GitCommit = setting.Value
			}
		case "vcs.time":
			if info.BuildDate == unknown {
				info.BuildDate = setting.Value
			}
		case "vcs.modified":
			info.Modified = setting.Value == "true"
		}
	}

	return info
}

// String renders the identity as a single human-readable line.
func (i Info) String() string {
	out := "goreveal " + i.Version + " (" + i.GitCommit
	if i.Modified {
		out += "-dirty"
	}
	return out + ") built " + i.BuildDate + " with " + i.GoVersion + " for " + i.Platform
}
