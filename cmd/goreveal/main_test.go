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
	if !strings.Contains(err.Error(), "runtime") {
		t.Fatalf("runInspectCmd() error = %q, want runtime in usage", err)
	}
}

func TestErrUsageRootIncludesInspectRuntime(t *testing.T) {
	t.Parallel()

	if got := errUsageRoot().Error(); !strings.Contains(got, "inspect <functions|packages|types|strings|runtime>") {
		t.Fatalf("errUsageRoot() = %q, want runtime in inspect usage", got)
	}
}
