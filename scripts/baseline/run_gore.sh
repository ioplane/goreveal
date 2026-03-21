#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: run_gore.sh <binary>" >&2
  exit 2
fi

BASELINES_ROOT="${GOREVEAL_BASELINES_ROOT:-/repos}"
BINARY_PATH="$(realpath "$1")"
PROJECT_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
export PYTHONPATH="${PROJECT_ROOT}${PYTHONPATH:+:${PYTHONPATH}}"

TMP_GO="$(mktemp --suffix=.go)"
TMP_JSON="$(mktemp)"
TMP_ERR="$(mktemp)"
trap 'rm -f "$TMP_GO" "$TMP_JSON" "$TMP_ERR"' EXIT

cat >"${TMP_GO}" <<'GO'
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	gore "github.com/goretk/gore"
)

type normalized struct {
	BuildInfo struct {
		GoVersion string `json:"go_version"`
		Path      string `json:"path"`
	} `json:"build_info"`
	Packages    []string `json:"packages"`
	SourceFiles []string `json:"source_files"`
	Functions   []string `json:"functions"`
	Types       []string `json:"types"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: run_gore.sh <binary>")
		os.Exit(2)
	}

	f, err := gore.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "gore.Open: %v\n", err)
		os.Exit(1)
	}

	var out normalized
	if f.BuildInfo != nil {
		if f.BuildInfo.Compiler != nil {
			out.BuildInfo.GoVersion = f.BuildInfo.Compiler.Name
		}
		if f.BuildInfo.ModInfo != nil {
			out.BuildInfo.Path = f.BuildInfo.ModInfo.Path
		}
	}

	pkgs, err := f.GetPackages()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetPackages: %v\n", err)
		os.Exit(1)
	}

	sourceFiles := make(map[string]struct{})
	functions := make(map[string]struct{})
	for _, pkg := range pkgs {
		out.Packages = append(out.Packages, pkg.Name)
		for _, fn := range pkg.Functions {
			if fn == nil || fn.Name == "" {
				continue
			}
			name := fn.Name
			if pkg.Name != "" && !strings.Contains(name, ".") && !strings.Contains(name, "/") {
				name = pkg.Name + "." + name
			}
			functions[name] = struct{}{}
		}
		for _, sf := range f.GetSourceFiles(pkg) {
			if sf == nil || sf.Name == "" {
				continue
			}
			sourceFiles[sf.Name] = struct{}{}
		}
	}
	types, err := f.GetTypes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetTypes: %v\n", err)
		os.Exit(1)
	}
	typeNames := make(map[string]struct{})
	for _, typ := range types {
		if typ == nil || typ.Name == "" {
			continue
		}
		typeNames[typ.Name] = struct{}{}
	}

	sort.Strings(out.Packages)
	for name := range functions {
		out.Functions = append(out.Functions, name)
	}
	for name := range sourceFiles {
		out.SourceFiles = append(out.SourceFiles, name)
	}
	for name := range typeNames {
		out.Types = append(out.Types, name)
	}
	sort.Strings(out.Functions)
	sort.Strings(out.SourceFiles)
	sort.Strings(out.Types)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}
}
GO

if ! (
  cd "${BASELINES_ROOT}/gore"
  /usr/local/go/bin/go run "${TMP_GO}" "${BINARY_PATH}"
) >"${TMP_JSON}" 2>"${TMP_ERR}"; then
  cat "${TMP_ERR}" >&2
  exit 1
fi

python3 - "${TMP_JSON}" <<'PY'
import json
import sys

from scripts.baseline.normalize import dump_json, normalize_gore

with open(sys.argv[1], "r", encoding="utf-8") as fh:
    raw = json.load(fh)

print(dump_json(normalize_gore(raw)))
PY
