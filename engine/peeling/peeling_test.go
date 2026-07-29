package peeling

import (
	"testing"

	"github.com/ioplane/goreveal/schema"
)

func TestBuildClassifiesFunctionsFromCanonicalTruth(t *testing.T) {
	t.Parallel()

	analysis := schema.Analysis{
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample"},
		Functions: []schema.Function{
			{
				Name:        "main.main",
				Package:     "main",
				ImportPath:  "example.com/sample",
				SourceFile:  "main.go",
				SourceLine:  10,
				Entry:       0x1000,
				End:         0x1100,
				ModuleLocal: true,
			},
			{
				Name:       "runtime.newobject",
				Package:    "runtime",
				ImportPath: "runtime",
				Entry:      0x2000,
				End:        0x2100,
			},
			{
				Name:       "fmt.Println",
				Package:    "fmt",
				ImportPath: "fmt",
				Entry:      0x3000,
				End:        0x3100,
			},
			{
				Name:       "github.com/acme/lib.Helper",
				Package:    "github.com/acme/lib",
				ImportPath: "github.com/acme/lib",
				Entry:      0x4000,
				End:        0x4100,
			},
		},
	}

	got := Build(analysis)
	if got == nil {
		t.Fatal("Build() = nil")
	}
	if got.Provenance.Source != "engine.peeling" || got.Provenance.Confidence != "medium" {
		t.Fatalf("Build() provenance = %#v", got.Provenance)
	}

	tests := []struct {
		name     string
		want     schema.PeelingClass
		evidence schema.PeelingEvidence
	}{
		{name: "main.main", want: schema.PeelingClassUser, evidence: schema.PeelingEvidenceModuleLocal},
		{name: "runtime.newobject", want: schema.PeelingClassRuntime, evidence: schema.PeelingEvidenceRuntimeImportPath},
		{name: "fmt.Println", want: schema.PeelingClassStdlib, evidence: schema.PeelingEvidenceStdlibImportPath},
		{name: "github.com/acme/lib.Helper", want: schema.PeelingClassThirdParty, evidence: schema.PeelingEvidenceThirdPartyImportPath},
	}

	for _, tt := range tests {
		gotFn, ok := findPeelingFunction(got.Functions, tt.name)
		if !ok {
			t.Fatalf("Build() missing %q in %#v", tt.name, got.Functions)
		}
		if gotFn.Classification != tt.want {
			t.Fatalf("Build() classification for %q = %q, want %q", tt.name, gotFn.Classification, tt.want)
		}
		if gotFn.ClassificationEvidence != tt.evidence {
			t.Fatalf("Build() classification evidence for %q = %q, want %q", tt.name, gotFn.ClassificationEvidence, tt.evidence)
		}
	}

	userPkg, ok := findPeelingPackage(got.Packages, "main", "example.com/sample")
	if !ok {
		t.Fatalf("Build() missing user package summary in %#v", got.Packages)
	}
	if userPkg.PrimaryClassification != schema.PeelingClassUser || userPkg.UserFunctionCount != 1 || userPkg.FunctionCount != 1 || !userPkg.ModuleLocal {
		t.Fatalf("Build() user package summary = %#v", userPkg)
	}

	runtimePkg, ok := findPeelingPackage(got.Packages, "runtime", "runtime")
	if !ok {
		t.Fatalf("Build() missing runtime package summary in %#v", got.Packages)
	}
	if runtimePkg.PrimaryClassification != schema.PeelingClassRuntime || runtimePkg.RuntimeFunctionCount != 1 {
		t.Fatalf("Build() runtime package summary = %#v", runtimePkg)
	}

	stdlibPkg, ok := findPeelingPackage(got.Packages, "fmt", "fmt")
	if !ok {
		t.Fatalf("Build() missing stdlib package summary in %#v", got.Packages)
	}
	if stdlibPkg.PrimaryClassification != schema.PeelingClassStdlib || stdlibPkg.StdlibFunctionCount != 1 {
		t.Fatalf("Build() stdlib package summary = %#v", stdlibPkg)
	}

	thirdPartyPkg, ok := findPeelingPackage(got.Packages, "github.com/acme/lib", "github.com/acme/lib")
	if !ok {
		t.Fatalf("Build() missing third-party package summary in %#v", got.Packages)
	}
	if thirdPartyPkg.PrimaryClassification != schema.PeelingClassThirdParty || thirdPartyPkg.ThirdPartyFunctionCount != 1 {
		t.Fatalf("Build() third-party package summary = %#v", thirdPartyPkg)
	}
}

func TestBuildUsesBuildInfoAsBoundedModuleFallback(t *testing.T) {
	t.Parallel()

	got := Build(schema.Analysis{
		BuildInfo: &schema.BuildInfo{Path: "example.com/sample"},
		Functions: []schema.Function{
			{
				Name:       "example.com/sample/internal/pkg.Helper",
				Package:    "example.com/sample/internal/pkg",
				ImportPath: "example.com/sample/internal/pkg",
				Entry:      0x1000,
				End:        0x1100,
			},
		},
	})
	if got == nil || len(got.Functions) != 1 {
		t.Fatalf("Build() = %#v", got)
	}
	if got.Functions[0].Classification != schema.PeelingClassUser {
		t.Fatalf("Build() module fallback classification = %q, want %q", got.Functions[0].Classification, schema.PeelingClassUser)
	}
	if got.Functions[0].ClassificationEvidence != schema.PeelingEvidenceBuildInfoPath {
		t.Fatalf("Build() module fallback evidence = %q, want %q", got.Functions[0].ClassificationEvidence, schema.PeelingEvidenceBuildInfoPath)
	}
	gotPkg, ok := findPeelingPackage(got.Packages, "example.com/sample/internal/pkg", "example.com/sample/internal/pkg")
	if !ok {
		t.Fatalf("Build() missing module-fallback package summary in %#v", got.Packages)
	}
	if gotPkg.PrimaryClassification != schema.PeelingClassUser || gotPkg.UserFunctionCount != 1 {
		t.Fatalf("Build() module-fallback package summary = %#v", gotPkg)
	}
}

func TestUserOnlyViewFiltersToUserClassifiedEntries(t *testing.T) {
	t.Parallel()

	got := UserOnlyView(&schema.PeelingAnalysis{
		Functions: []schema.PeelingFunction{
			{Name: "main.main", Classification: schema.PeelingClassUser, Entry: 0x1000, End: 0x1100},
			{Name: "runtime.newobject", Classification: schema.PeelingClassRuntime, Entry: 0x2000, End: 0x2100},
		},
		Packages: []schema.PeelingPackage{
			{Name: "main", ImportPath: "example.com/sample", FunctionCount: 1, UserFunctionCount: 1, PrimaryClassification: schema.PeelingClassUser},
			{Name: "runtime", ImportPath: "runtime", FunctionCount: 1, RuntimeFunctionCount: 1, PrimaryClassification: schema.PeelingClassRuntime},
		},
		Provenance: schema.Provenance{Source: "engine.peeling", Confidence: "medium"},
	})
	if got == nil {
		t.Fatal("UserOnlyView() = nil")
	}
	if len(got.Functions) != 1 || got.Functions[0].Name != "main.main" {
		t.Fatalf("UserOnlyView() functions = %#v", got.Functions)
	}
	if len(got.Packages) != 1 || got.Packages[0].Name != "main" {
		t.Fatalf("UserOnlyView() packages = %#v", got.Packages)
	}
	if got.Provenance.Source != "engine.peeling.user_only" || got.Provenance.Confidence != "medium" {
		t.Fatalf("UserOnlyView() provenance = %#v", got.Provenance)
	}
}

func TestBuildUsesFingerprintAssistedStdlibAndRuntimeRefinement(t *testing.T) {
	t.Parallel()

	got := Build(schema.Analysis{
		Functions: []schema.Function{
			{
				Name:  "asyncPreempt",
				Entry: 0x1000,
				End:   0x1100,
			},
			{
				Name:       "printArg",
				SourceFile: "/usr/local/go/src/fmt/print.go",
				Entry:      0x2000,
				End:        0x2100,
			},
			{
				Name:       "mstart0",
				SourceFile: "runtime/proc.go",
				Entry:      0x3000,
				End:        0x3100,
			},
		},
	})
	if got == nil || len(got.Functions) != 3 {
		t.Fatalf("Build() = %#v", got)
	}

	runtimeByName, ok := findPeelingFunction(got.Functions, "asyncPreempt")
	if !ok {
		t.Fatalf("Build() missing runtime fingerprint entry in %#v", got.Functions)
	}
	if runtimeByName.Classification != schema.PeelingClassRuntime || runtimeByName.ClassificationEvidence != schema.PeelingEvidenceRuntimeNameFingerprint {
		t.Fatalf("Build() runtime name fingerprint = %#v", runtimeByName)
	}

	stdlibBySource, ok := findPeelingFunction(got.Functions, "printArg")
	if !ok {
		t.Fatalf("Build() missing stdlib source fingerprint entry in %#v", got.Functions)
	}
	if stdlibBySource.Classification != schema.PeelingClassStdlib || stdlibBySource.ClassificationEvidence != schema.PeelingEvidenceStdlibSourceFingerprint {
		t.Fatalf("Build() stdlib source fingerprint = %#v", stdlibBySource)
	}

	runtimeBySource, ok := findPeelingFunction(got.Functions, "mstart0")
	if !ok {
		t.Fatalf("Build() missing runtime source fingerprint entry in %#v", got.Functions)
	}
	if runtimeBySource.Classification != schema.PeelingClassRuntime || runtimeBySource.ClassificationEvidence != schema.PeelingEvidenceRuntimeSourceFingerprint {
		t.Fatalf("Build() runtime source fingerprint = %#v", runtimeBySource)
	}
}

func findPeelingFunction(funcs []schema.PeelingFunction, want string) (schema.PeelingFunction, bool) {
	for _, fn := range funcs {
		if fn.Name == want {
			return fn, true
		}
	}

	return schema.PeelingFunction{}, false
}

func findPeelingPackage(pkgs []schema.PeelingPackage, wantName, wantImportPath string) (schema.PeelingPackage, bool) {
	for _, pkg := range pkgs {
		if pkg.Name == wantName && pkg.ImportPath == wantImportPath {
			return pkg, true
		}
	}

	return schema.PeelingPackage{}, false
}
