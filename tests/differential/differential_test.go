package differential

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dantte-lp/goreveal/engine"
	"github.com/dantte-lp/goreveal/schema"
)

type goresymNormalized struct {
	BuildInfo struct {
		GoVersion string `json:"go_version"`
		Path      string `json:"path"`
	} `json:"build_info"`
	Files     []string `json:"files"`
	Functions []string `json:"functions"`
}

type redressNormalized struct {
	BuildInfo struct {
		Path string `json:"path"`
	} `json:"build_info"`
	Packages    []string `json:"packages"`
	SourceFiles []string `json:"source_files"`
	Functions   []string `json:"functions"`
}

type goreNormalized struct {
	BuildInfo struct {
		GoVersion string `json:"go_version"`
		Path      string `json:"path"`
	} `json:"build_info"`
	Packages    []string `json:"packages"`
	SourceFiles []string `json:"source_files"`
	Functions   []string `json:"functions"`
	Types       []string `json:"types"`
}

func TestFixtureDifferentialV1(t *testing.T) {
	t.Parallel()

	baselinesRoot := os.Getenv("GOREVEAL_BASELINES_ROOT")
	if baselinesRoot == "" {
		t.Skip("GOREVEAL_BASELINES_ROOT is not set")
	}

	fixture := filepath.Join("..", "..", "corpus", "fixtures", "go-elf-buildinfo-linux-amd64", "fixture.bin")
	analysis, err := engine.New().AnalyzeFile(context.Background(), fixture)
	if err != nil {
		t.Fatalf("AnalyzeFile() error = %v", err)
	}

	goresymOut := runScript(t, filepath.Join("..", "..", "scripts", "baseline", "run_goresym.sh"), fixture)
	var goresym goresymNormalized
	if err := json.Unmarshal([]byte(goresymOut), &goresym); err != nil {
		t.Fatalf("unmarshal GoReSym output: %v\n%s", err, goresymOut)
	}

	redressOut := runScript(t, filepath.Join("..", "..", "scripts", "baseline", "run_redress.sh"), fixture)
	var redress redressNormalized
	if err := json.Unmarshal([]byte(redressOut), &redress); err != nil {
		t.Fatalf("unmarshal redress output: %v\n%s", err, redressOut)
	}

	goreOut := runScript(t, filepath.Join("..", "..", "scripts", "baseline", "run_gore.sh"), fixture)
	var gore goreNormalized
	if err := json.Unmarshal([]byte(goreOut), &gore); err != nil {
		t.Fatalf("unmarshal gore output: %v\n%s", err, goreOut)
	}

	if analysis.BuildInfo == nil {
		t.Fatal("analysis build info is nil")
	}
	if analysis.BuildInfo.Path != goresym.BuildInfo.Path {
		t.Fatalf("build info path mismatch: goreveal=%q goresym=%q", analysis.BuildInfo.Path, goresym.BuildInfo.Path)
	}
	if analysis.BuildInfo.GoVersion != goresym.BuildInfo.GoVersion {
		t.Fatalf("go version mismatch: goreveal=%q goresym=%q", analysis.BuildInfo.GoVersion, goresym.BuildInfo.GoVersion)
	}
	if analysis.BuildInfo.Path != gore.BuildInfo.Path {
		t.Fatalf("build info path mismatch: goreveal=%q gore=%q", analysis.BuildInfo.Path, gore.BuildInfo.Path)
	}
	if analysis.BuildInfo.Path != redress.BuildInfo.Path {
		t.Fatalf("build info path mismatch: goreveal=%q redress=%q", analysis.BuildInfo.Path, redress.BuildInfo.Path)
	}
	if analysis.BuildInfo.GoVersion != gore.BuildInfo.GoVersion {
		t.Fatalf("go version mismatch: goreveal=%q gore=%q", analysis.BuildInfo.GoVersion, gore.BuildInfo.GoVersion)
	}

	if !slices.Contains(redress.Packages, "main") {
		t.Fatalf("redress packages missing main: %#v", redress.Packages)
	}
	if !slices.Contains(gore.Packages, "main") {
		t.Fatalf("gore packages missing main: %#v", gore.Packages)
	}
	if !containsPackage(analysis.Packages, "main") {
		t.Fatalf("goreveal packages missing main: %#v", analysis.Packages)
	}

	if analysis.SourceTree == nil {
		t.Fatal("analysis source tree is nil")
	}
	wantFile := analysis.SourceTree.Root + "/" + analysis.SourceTree.Files[0]
	if !slices.Contains(goresym.Files, wantFile) {
		t.Fatalf("GoReSym files missing %q", wantFile)
	}
	if !slices.Contains(goresym.Functions, "main.main") {
		t.Fatalf("GoReSym functions missing main.main: %#v", goresym.Functions)
	}
	if !slices.Contains(goresym.Functions, "main.helperAdd") {
		t.Fatalf("GoReSym functions missing main.helperAdd: %#v", goresym.Functions)
	}
	if !slices.Contains(goresym.Functions, "main.helperBanner") {
		t.Fatalf("GoReSym functions missing main.helperBanner: %#v", goresym.Functions)
	}
	if !slices.Contains(redress.SourceFiles, analysis.SourceTree.Files[0]) {
		t.Fatalf("redress source files missing %q: %#v", analysis.SourceTree.Files[0], redress.SourceFiles)
	}
	if !slices.Contains(redress.Functions, "main.main") {
		t.Fatalf("redress functions missing main.main: %#v", redress.Functions)
	}
	if !slices.Contains(redress.Functions, "main.helperBanner") {
		t.Fatalf("redress functions missing main.helperBanner: %#v", redress.Functions)
	}
	if !slices.Contains(gore.SourceFiles, analysis.SourceTree.Files[0]) {
		t.Fatalf("gore source files missing %q", analysis.SourceTree.Files[0])
	}
	if !slices.Contains(gore.Functions, "main.main") {
		t.Fatalf("gore functions missing main.main: %#v", gore.Functions)
	}
	if !slices.Contains(gore.Functions, "main.helperBanner") {
		t.Fatalf("gore functions missing main.helperBanner: %#v", gore.Functions)
	}
}

func runScript(t *testing.T, scriptPath, fixture string) string {
	t.Helper()

	cmd := exec.CommandContext(context.Background(), scriptPath, fixture)
	cmd.Env = append(os.Environ(), "LC_ALL=C")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s failed: %v\n%s", filepath.Base(scriptPath), err, out)
	}
	return strings.TrimSpace(string(out))
}

func containsPackage(pkgs []schema.Package, want string) bool {
	return slices.ContainsFunc(pkgs, func(pkg schema.Package) bool {
		return pkg.Name == want
	})
}
