package peeling

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ioplane/goreveal/schema"
)

// Build derives a bounded function-level code-peeling layer from canonical truth.
func Build(analysis schema.Analysis) *schema.PeelingAnalysis {
	if len(analysis.Functions) == 0 {
		return nil
	}

	out := &schema.PeelingAnalysis{
		Functions: make([]schema.PeelingFunction, 0, len(analysis.Functions)),
		Provenance: schema.Provenance{
			Source:     "engine.peeling",
			Confidence: "medium",
		},
	}

	for _, fn := range analysis.Functions {
		classification, evidence := classifyFunction(fn, analysis.BuildInfo)
		out.Functions = append(out.Functions, schema.PeelingFunction{
			Name:                   fn.Name,
			Package:                fn.Package,
			ImportPath:             fn.ImportPath,
			SourceFile:             fn.SourceFile,
			SourceLine:             fn.SourceLine,
			Entry:                  fn.Entry,
			End:                    fn.End,
			ModuleLocal:            fn.ModuleLocal,
			Classification:         classification,
			ClassificationEvidence: evidence,
		})
	}
	out.Packages = summarizePackages(out.Functions)

	return out
}

// UserOnlyView projects only user-classified peeling entries for analyst-facing workflows.
func UserOnlyView(peeling *schema.PeelingAnalysis) *schema.PeelingAnalysis {
	if peeling == nil {
		return nil
	}

	out := &schema.PeelingAnalysis{
		Functions:  make([]schema.PeelingFunction, 0, len(peeling.Functions)),
		Packages:   make([]schema.PeelingPackage, 0, len(peeling.Packages)),
		Provenance: schema.Provenance{Source: "engine.peeling.user_only", Confidence: peeling.Provenance.Confidence},
	}

	for _, fn := range peeling.Functions {
		if fn.Classification != schema.PeelingClassUser {
			continue
		}
		out.Functions = append(out.Functions, fn)
	}
	for _, pkg := range peeling.Packages {
		if pkg.PrimaryClassification != schema.PeelingClassUser {
			continue
		}
		out.Packages = append(out.Packages, pkg)
	}

	return out
}

func classifyFunction(fn schema.Function, info *schema.BuildInfo) (schema.PeelingClass, schema.PeelingEvidence) {
	if fn.ModuleLocal {
		return schema.PeelingClassUser, schema.PeelingEvidenceModuleLocal
	}
	if isModuleLocalImportPath(fn.ImportPath, info) {
		return schema.PeelingClassUser, schema.PeelingEvidenceBuildInfoPath
	}
	if fn.Package == "main" {
		return schema.PeelingClassUser, schema.PeelingEvidencePackageMain
	}
	if evidence, ok := runtimeEvidence(fn); ok {
		return schema.PeelingClassRuntime, evidence
	}
	if isThirdPartyImportPath(fn.ImportPath) {
		return schema.PeelingClassThirdParty, schema.PeelingEvidenceThirdPartyImportPath
	}
	if fn.ImportPath != "" {
		return schema.PeelingClassStdlib, schema.PeelingEvidenceStdlibImportPath
	}
	if evidence, ok := stdlibEvidence(fn); ok {
		return schema.PeelingClassStdlib, evidence
	}

	return schema.PeelingClassStdlib, schema.PeelingEvidenceDefaultStdlib
}

func isModuleLocalImportPath(importPath string, info *schema.BuildInfo) bool {
	if info == nil || info.Path == "" || importPath == "" {
		return false
	}

	return importPath == info.Path || strings.HasPrefix(importPath, info.Path+"/")
}

func runtimeEvidence(fn schema.Function) (schema.PeelingEvidence, bool) {
	if fn.ImportPath == "runtime" || strings.HasPrefix(fn.ImportPath, "runtime/") {
		return schema.PeelingEvidenceRuntimeImportPath, true
	}
	if fn.Package == "runtime" {
		return schema.PeelingEvidenceRuntimeImportPath, true
	}
	if strings.HasPrefix(fn.Name, "runtime.") || isRuntimeNameFingerprint(fn.Name) {
		return schema.PeelingEvidenceRuntimeNameFingerprint, true
	}
	if isRuntimeSourceFingerprint(fn.SourceFile) {
		return schema.PeelingEvidenceRuntimeSourceFingerprint, true
	}

	return "", false
}

func isThirdPartyImportPath(importPath string) bool {
	if importPath == "" {
		return false
	}

	firstSegment, _, _ := strings.Cut(importPath, "/")
	return strings.Contains(firstSegment, ".")
}

func isRuntimeNameFingerprint(name string) bool {
	switch name {
	case "asyncPreempt", "morestack", "morestack_noctxt", "duffcopy", "duffzero", "gcWriteBarrier", "typedmemmove", "typedslicecopy", "memequal128", "mapassign_faststr", "mapaccess1_faststr", "mapaccess2_faststr", "convT64", "convTstring":
		return true
	default:
		return false
	}
}

func isRuntimeSourceFingerprint(sourceFile string) bool {
	if sourceFile == "" {
		return false
	}

	sourceFile = filepath.ToSlash(sourceFile)
	return strings.HasPrefix(sourceFile, "runtime/") || strings.Contains(sourceFile, "/runtime/")
}

func stdlibEvidence(fn schema.Function) (schema.PeelingEvidence, bool) {
	if isStdlibNameFingerprint(fn.Name) {
		return schema.PeelingEvidenceStdlibNameFingerprint, true
	}
	if isStdlibSourceFingerprint(fn.SourceFile) {
		return schema.PeelingEvidenceStdlibSourceFingerprint, true
	}

	return "", false
}

func isStdlibNameFingerprint(name string) bool {
	for _, prefix := range []string{
		"fmt.",
		"bytes.",
		"strings.",
		"os.",
		"io.",
		"syscall.",
		"sync.",
		"crypto/",
		"net/http.",
		"net.",
		"path/filepath.",
		"internal/",
	} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}

func isStdlibSourceFingerprint(sourceFile string) bool {
	if sourceFile == "" {
		return false
	}

	sourceFile = filepath.ToSlash(sourceFile)
	for _, segment := range []string{
		"/src/fmt/",
		"/src/bytes/",
		"/src/strings/",
		"/src/os/",
		"/src/io/",
		"/src/syscall/",
		"/src/sync/",
		"/src/crypto/",
		"/src/net/",
		"/src/path/filepath/",
		"/src/internal/",
	} {
		if strings.Contains(sourceFile, segment) {
			return true
		}
	}
	for _, prefix := range []string{
		"fmt/",
		"bytes/",
		"strings/",
		"os/",
		"io/",
		"syscall/",
		"sync/",
		"crypto/",
		"net/",
		"path/filepath/",
		"internal/",
	} {
		if strings.HasPrefix(sourceFile, prefix) {
			return true
		}
	}

	return false
}

func summarizePackages(funcs []schema.PeelingFunction) []schema.PeelingPackage {
	if len(funcs) == 0 {
		return nil
	}

	type key struct {
		name       string
		importPath string
	}

	summaries := make(map[key]*schema.PeelingPackage)
	for _, fn := range funcs {
		if fn.Package == "" && fn.ImportPath == "" {
			continue
		}

		k := key{name: fn.Package, importPath: fn.ImportPath}
		summary := summaries[k]
		if summary == nil {
			summary = &schema.PeelingPackage{
				Name:        fn.Package,
				ImportPath:  fn.ImportPath,
				ModuleLocal: fn.ModuleLocal,
			}
			summaries[k] = summary
		}

		summary.FunctionCount++
		summary.ModuleLocal = summary.ModuleLocal || fn.ModuleLocal
		switch fn.Classification {
		case schema.PeelingClassUser:
			summary.UserFunctionCount++
		case schema.PeelingClassStdlib:
			summary.StdlibFunctionCount++
		case schema.PeelingClassRuntime:
			summary.RuntimeFunctionCount++
		case schema.PeelingClassThirdParty:
			summary.ThirdPartyFunctionCount++
		}
		summary.PrimaryClassification = primaryClassification(*summary)
	}

	out := make([]schema.PeelingPackage, 0, len(summaries))
	for _, summary := range summaries {
		out = append(out, *summary)
	}

	slices.SortFunc(out, func(a, b schema.PeelingPackage) int {
		if n := cmp.Compare(a.ImportPath, b.ImportPath); n != 0 {
			return n
		}
		return cmp.Compare(a.Name, b.Name)
	})

	return out
}

func primaryClassification(summary schema.PeelingPackage) schema.PeelingClass {
	bestClass := schema.PeelingClassUser
	bestCount := summary.UserFunctionCount

	type candidate struct {
		class schema.PeelingClass
		count int
	}

	for _, candidate := range []candidate{
		{class: schema.PeelingClassRuntime, count: summary.RuntimeFunctionCount},
		{class: schema.PeelingClassStdlib, count: summary.StdlibFunctionCount},
		{class: schema.PeelingClassThirdParty, count: summary.ThirdPartyFunctionCount},
	} {
		if candidate.count > bestCount {
			bestClass = candidate.class
			bestCount = candidate.count
		}
	}

	return bestClass
}
