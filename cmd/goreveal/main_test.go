package main

import (
	"context"
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
