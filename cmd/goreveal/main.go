package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	internalcmd "github.com/dantte-lp/goreveal/cmd/goreveal/internal"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errUsageRoot()
	}

	switch args[0] {
	case "analyze":
		return runAnalyzeCmd(ctx, args)
	case "inspect":
		return runInspectCmd(ctx, args)
	case "source-tree":
		return runSourceTreeCmd(ctx, args)
	case "deobfuscate":
		return runDeobfuscateCmd(ctx, args)
	case "export":
		return runExportCmd(ctx, args)
	case "diff":
		return runDiffCmd(ctx, args)
	default:
		return errUsageRoot()
	}
}

func runAnalyzeCmd(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: goreveal analyze <binary>")
	}
	return internalcmd.RunAnalyze(ctx, os.Stdout, args[1])
}

func runInspectCmd(ctx context.Context, args []string) error {
	if len(args) != 3 {
		return errors.New("usage: goreveal inspect <functions|packages|types|strings|runtime> <binary>")
	}

	switch args[1] {
	case "functions":
		return internalcmd.RunInspectFunctions(ctx, os.Stdout, args[2])
	case "packages":
		return internalcmd.RunInspectPackages(ctx, os.Stdout, args[2])
	case "types":
		return internalcmd.RunInspectTypes(ctx, os.Stdout, args[2])
	case "strings":
		return internalcmd.RunInspectStrings(ctx, os.Stdout, args[2])
	case "runtime":
		return internalcmd.RunInspectRuntime(ctx, os.Stdout, args[2])
	default:
		return errors.New("usage: goreveal inspect <functions|packages|types|strings|runtime> <binary>")
	}
}

func runSourceTreeCmd(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: goreveal source-tree <binary>")
	}
	return internalcmd.RunSourceTree(ctx, os.Stdout, args[1])
}

func runDeobfuscateCmd(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: goreveal deobfuscate <binary>")
	}
	return internalcmd.RunDeobfuscate(ctx, os.Stdout, args[1])
}

func runExportCmd(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("usage: goreveal export <sqlite|ida|ghidra> [args]")
	}

	switch args[1] {
	case "sqlite":
		if len(args) != 4 {
			return errors.New("usage: goreveal export sqlite <database> <binary>")
		}
		return internalcmd.RunExportSQLite(ctx, args[2], args[3])
	case "ida":
		if len(args) != 3 {
			return errors.New("usage: goreveal export ida <binary>")
		}
		return internalcmd.RunExportIDA(ctx, os.Stdout, args[2])
	case "ghidra":
		if len(args) != 3 {
			return errors.New("usage: goreveal export ghidra <binary>")
		}
		return internalcmd.RunExportGhidra(ctx, os.Stdout, args[2])
	default:
		return errors.New("usage: goreveal export <sqlite|ida|ghidra> [args]")
	}
}

func runDiffCmd(ctx context.Context, args []string) error {
	if len(args) != 5 || args[1] != "sqlite" {
		return errors.New("usage: goreveal diff sqlite <database> <left-id> <right-id>")
	}

	leftID, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return fmt.Errorf("parse left id %q: %w", args[3], err)
	}
	rightID, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return fmt.Errorf("parse right id %q: %w", args[4], err)
	}
	return internalcmd.RunDiffSQLite(ctx, os.Stdout, args[2], leftID, rightID)
}

func errUsageRoot() error {
	return errors.New("usage: goreveal analyze <binary> | goreveal inspect <functions|packages|types|strings|runtime> <binary> | goreveal source-tree <binary> | goreveal deobfuscate <binary> | goreveal export <sqlite|ida|ghidra> [args] | goreveal diff sqlite <database> <left-id> <right-id>")
}
