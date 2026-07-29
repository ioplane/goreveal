package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	internalcmd "github.com/ioplane/goreveal/cmd/goreveal/internal"
	"github.com/ioplane/goreveal/internal/version"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	// Zero-argument and single-argument commands are handled before the arity
	// guard below, which exists for the binary-taking subcommands.
	if len(args) == 0 {
		return errUsageRoot()
	}

	switch args[0] {
	case "version", "--version", "-v":
		return runVersionCmd(os.Stdout, args)
	case "help", "--help", "-h":
		fmt.Fprintln(os.Stdout, usageRoot())
		return nil
	}

	if len(args) < 2 {
		return errUsageRoot()
	}

	switch args[0] {
	case "analyze":
		return runAnalyzeCmd(ctx, args)
	case "inspect":
		return runInspectCmd(ctx, args)
	case "peel":
		return runPeelCmd(ctx, args)
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

func runVersionCmd(out io.Writer, args []string) error {
	info := version.Get()

	if len(args) > 1 {
		if args[1] != "--json" {
			return errors.New("usage: goreveal version [--json]")
		}
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(info)
	}

	_, err := fmt.Fprintln(out, info)
	return err
}

func runAnalyzeCmd(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: goreveal analyze <binary>")
	}
	return internalcmd.RunAnalyze(ctx, os.Stdout, args[1])
}

func runInspectCmd(ctx context.Context, args []string) error {
	if len(args) != 3 {
		return errors.New("usage: goreveal inspect <functions|packages|types|strings|runtime|peeling> <binary>")
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
	case "peeling":
		return internalcmd.RunInspectPeeling(ctx, os.Stdout, args[2])
	default:
		return errors.New("usage: goreveal inspect <functions|packages|types|strings|runtime|peeling> <binary>")
	}
}

func runPeelCmd(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return errors.New("usage: goreveal peel <binary>")
	}

	return internalcmd.RunPeel(ctx, os.Stdout, args[1])
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

// diffRunner projects a stored-run comparison to out.
type diffRunner func(ctx context.Context, out io.Writer, database string, leftID, rightID int64) error

// diffTarget is the trailing operand triple shared by every diff subcommand:
// a database path and the two stored run IDs to compare.
type diffTarget struct {
	database string
	leftID   int64
	rightID  int64
}

// parseDiffTarget consumes exactly three operands. Taking a fixed-size array
// rather than a slice makes the bounds provable at the call site, so neither the
// reader nor a static analyzer has to correlate a length guard with the indices.
func parseDiffTarget(operands [3]string) (diffTarget, error) {
	leftID, err := strconv.ParseInt(operands[1], 10, 64)
	if err != nil {
		return diffTarget{}, fmt.Errorf("parse left id %q: %w", operands[1], err)
	}
	rightID, err := strconv.ParseInt(operands[2], 10, 64)
	if err != nil {
		return diffTarget{}, fmt.Errorf("parse right id %q: %w", operands[2], err)
	}
	return diffTarget{database: operands[0], leftID: leftID, rightID: rightID}, nil
}

// diffProjections maps each `diff <projection> sqlite …` subcommand onto its
// runner. The bare `diff sqlite …` form is handled separately because it has no
// projection word.
var diffProjections = map[string]diffRunner{
	"review":  internalcmd.RunDiffReviewSQLite,
	"handoff": internalcmd.RunDiffHandoffSQLite,
	"next":    internalcmd.RunDiffNextSQLite,
}

func runDiffCmd(ctx context.Context, args []string) error {
	// goreveal diff sqlite <database> <left-id> <right-id>
	//
	//nolint:gosec // G602 false positive: every index is guarded by the len check
	if len(args) == 5 && args[1] == "sqlite" {
		target, err := parseDiffTarget([3]string{args[2], args[3], args[4]})
		if err != nil {
			return err
		}
		return internalcmd.RunDiffSQLite(ctx, os.Stdout, target.database, target.leftID, target.rightID)
	}

	// goreveal diff <review|handoff|next> sqlite <database> <left-id> <right-id>
	//
	//nolint:gosec // G602 false positive: every index is guarded by the len check
	if len(args) == 6 && args[2] == "sqlite" {
		if run, ok := diffProjections[args[1]]; ok {
			target, err := parseDiffTarget([3]string{args[3], args[4], args[5]})
			if err != nil {
				return err
			}
			return run(ctx, os.Stdout, target.database, target.leftID, target.rightID)
		}
	}

	return errUsageDiff()
}

func errUsageDiff() error {
	return errors.New("usage: goreveal diff sqlite <database> <left-id> <right-id> | goreveal diff review sqlite <database> <left-id> <right-id> | goreveal diff handoff sqlite <database> <left-id> <right-id> | goreveal diff next sqlite <database> <left-id> <right-id>")
}

func usageRoot() string {
	return "usage: goreveal analyze <binary> | goreveal inspect <functions|packages|types|strings|runtime|peeling> <binary> | goreveal peel <binary> | goreveal source-tree <binary> | goreveal deobfuscate <binary> | goreveal export <sqlite|ida|ghidra> [args] | goreveal diff sqlite <database> <left-id> <right-id> | goreveal diff review sqlite <database> <left-id> <right-id> | goreveal diff handoff sqlite <database> <left-id> <right-id> | goreveal diff next sqlite <database> <left-id> <right-id> | goreveal version [--json] | goreveal help"
}

func errUsageRoot() error {
	return errors.New(usageRoot())
}
