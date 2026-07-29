package main

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestRunInspectCmdUsageIncludesRuntime(t *testing.T) {
	t.Parallel()

	err := runInspectCmd(context.Background(), []string{"inspect", "bogus", "fixture.bin"})
	if err == nil {
		t.Fatal("runInspectCmd() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), "runtime") || !strings.Contains(err.Error(), "peeling") {
		t.Fatalf("runInspectCmd() error = %q, want runtime and peeling in usage", err)
	}
}

func TestErrUsageRootIncludesInspectRuntime(t *testing.T) {
	t.Parallel()

	if got := errUsageRoot().Error(); !strings.Contains(got, "inspect <functions|packages|types|strings|runtime|peeling>") || !strings.Contains(got, "peel <binary>") || !strings.Contains(got, "diff review sqlite") || !strings.Contains(got, "diff handoff sqlite") || !strings.Contains(got, "diff next sqlite") {
		t.Fatalf("errUsageRoot() = %q, want runtime/peeling inspect usage, peel command, and diff review/handoff/next commands", got)
	}
}

func TestRunVersionCmdPlainText(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	if err := runVersionCmd(&out, []string{"version"}); err != nil {
		t.Fatalf("runVersionCmd() error = %v, want nil", err)
	}

	got := out.String()
	if !strings.HasPrefix(got, "goreveal ") {
		t.Fatalf("runVersionCmd() = %q, want a line starting with %q", got, "goreveal ")
	}
	for _, want := range []string{"built", "linux", "darwin", "windows"} {
		if strings.Contains(got, want) {
			return
		}
	}
	t.Fatalf("runVersionCmd() = %q, want platform and build metadata", got)
}

func TestRunVersionCmdJSONIsMachineReadable(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	if err := runVersionCmd(&out, []string{"version", "--json"}); err != nil {
		t.Fatalf("runVersionCmd(--json) error = %v, want nil", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal(%q) error = %v, want valid JSON", out.String(), err)
	}
	for _, key := range []string{"version", "git_commit", "build_date", "go_version", "platform", "modified"} {
		if _, ok := decoded[key]; !ok {
			t.Errorf("version JSON missing key %q; got keys %v", key, decoded)
		}
	}
}

func TestRunVersionCmdRejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	err := runVersionCmd(&out, []string{"version", "--yaml"})
	if err == nil {
		t.Fatal("runVersionCmd(--yaml) error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), "--json") {
		t.Fatalf("runVersionCmd(--yaml) error = %q, want --json in usage", err)
	}
}

func TestRunDispatchesVersionAndHelpWithoutBinaryArgument(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"version"}, {"--version"}, {"-v"}, {"help"}, {"--help"}, {"-h"}} {
		if err := run(context.Background(), args); err != nil {
			t.Errorf("run(%v) error = %v, want nil", args, err)
		}
	}

	if err := run(context.Background(), nil); err == nil {
		t.Error("run(nil) error = nil, want usage error")
	}
}

func TestUsageRootIncludesVersionAndHelp(t *testing.T) {
	t.Parallel()

	got := usageRoot()
	if !strings.Contains(got, "version [--json]") || !strings.Contains(got, "goreveal help") {
		t.Fatalf("usageRoot() = %q, want version and help entries", got)
	}
}

func TestRunDiffCmdUsageIncludesReviewSQLite(t *testing.T) {
	t.Parallel()

	err := runDiffCmd(context.Background(), []string{"diff", "bogus"})
	if err == nil {
		t.Fatal("runDiffCmd() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), "diff review sqlite") || !strings.Contains(err.Error(), "diff handoff sqlite") || !strings.Contains(err.Error(), "diff next sqlite") || !strings.Contains(err.Error(), "diff sqlite") {
		t.Fatalf("runDiffCmd() error = %q, want diff sqlite, diff review sqlite, diff handoff sqlite, and diff next sqlite in usage", err)
	}
}
