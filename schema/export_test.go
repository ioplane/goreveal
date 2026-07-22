package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/token"
	"reflect"
	"strings"
	"testing"
)

type failingV1WireValue struct {
	err error
}

func (value failingV1WireValue) MarshalJSON() ([]byte, error) {
	return nil, value.err
}

func TestIDAExportV1RespectsEncoderEscapeHTML(t *testing.T) {
	t.Parallel()

	assertV1ExportRespectsEncoderEscapeHTML(t, NewIDAExport(htmlSensitiveAnalysis()))
}

func TestGhidraExportV1RespectsEncoderEscapeHTML(t *testing.T) {
	t.Parallel()

	assertV1ExportRespectsEncoderEscapeHTML(t, NewGhidraExport(htmlSensitiveAnalysis()))
}

func TestV1WireMarshalJSONReturnsOneValidValueWithoutNewline(t *testing.T) {
	t.Parallel()

	for name, marshaler := range map[string]json.Marshaler{
		"IDA":    NewIDAExport(htmlSensitiveAnalysis()),
		"Ghidra": NewGhidraExport(htmlSensitiveAnalysis()),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			encoded, err := marshaler.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if !json.Valid(encoded) {
				t.Fatalf("MarshalJSON() returned invalid JSON: %q", encoded)
			}
			if bytes.HasSuffix(encoded, []byte{'\n'}) {
				t.Fatalf("MarshalJSON() returned a trailing newline: %q", encoded)
			}
		})
	}
}

func TestMarshalV1WirePropagatesEncodingError(t *testing.T) {
	t.Parallel()

	want := errors.New("fixture encoding failure")
	_, got := marshalV1Wire(failingV1WireValue{err: want})
	if !errors.Is(got, want) {
		t.Fatalf("marshalV1Wire() error = %v, want %v", got, want)
	}
}

func assertV1ExportRespectsEncoderEscapeHTML(t *testing.T, export any) {
	t.Helper()

	for _, testCase := range []struct {
		name       string
		escapeHTML bool
		want       string
		reject     string
	}{
		{name: "default", escapeHTML: true, want: `\u003c\u003e\u0026`, reject: `<>&`},
		{name: "disabled", escapeHTML: false, want: `<>&`, reject: `\u003c\u003e\u0026`},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			encoder := json.NewEncoder(&output)
			encoder.SetEscapeHTML(testCase.escapeHTML)
			if err := encoder.Encode(export); err != nil {
				t.Fatalf("Encode() error = %v", err)
			}
			if !strings.Contains(output.String(), testCase.want) {
				t.Fatalf("Encode() output does not contain %q: %s", testCase.want, output.Bytes())
			}
			if strings.Contains(output.String(), testCase.reject) {
				t.Fatalf("Encode() output unexpectedly contains %q: %s", testCase.reject, output.Bytes())
			}
		})
	}
}

func htmlSensitiveAnalysis() Analysis {
	return Analysis{
		Input: Input{Path: "/tmp/<>&.bin", Format: "elf"},
		Functions: []Function{{
			Name: "main.<>&",
		}},
		Strings: []StringCandidate{{
			Value:  "<>&",
			Region: ".rodata",
		}},
	}
}

func TestV1WireTypesDoNotEmbedCanonicalContracts(t *testing.T) {
	t.Parallel()

	for name, root := range map[string]reflect.Type{
		"IDA":    reflect.TypeFor[idaExportV1Wire](),
		"Ghidra": reflect.TypeFor[ghidraExportV1Wire](),
	} {
		assertV1WireTypeIsolated(t, name, root)
	}
}

func TestV1WireGuardRejectsCanonicalNamedScalars(t *testing.T) {
	t.Parallel()

	type scalarContainer struct {
		Evidence PeelingEvidence
	}
	type scalarGuardFixture struct {
		Trust      RuntimeTrustSummary
		SourceKind *SourceEvidenceKind
		Classes    []PeelingClass
		Nested     scalarContainer
	}

	violations := v1WireTypeViolations(reflect.TypeFor[scalarGuardFixture]())
	want := map[reflect.Type]struct{}{
		reflect.TypeFor[RuntimeTrustSummary](): {},
		reflect.TypeFor[SourceEvidenceKind]():  {},
		reflect.TypeFor[PeelingClass]():        {},
		reflect.TypeFor[PeelingEvidence]():     {},
	}
	for _, violation := range violations {
		delete(want, violation)
	}
	if len(want) != 0 {
		t.Fatalf("v1 wire guard missed canonical named scalar types: %v", want)
	}
}

func assertV1WireTypeIsolated(t *testing.T, name string, root reflect.Type) {
	t.Helper()
	for _, violation := range v1WireTypeViolations(root) {
		t.Errorf("%s v1 wire recursively embeds mutable schema type %s", name, violation)
	}
}

func v1WireTypeViolations(root reflect.Type) []reflect.Type {
	schemaPackagePath := reflect.TypeFor[Analysis]().PkgPath()
	seen := make(map[reflect.Type]struct{})
	var violations []reflect.Type
	var visit func(reflect.Type)
	visit = func(current reflect.Type) {
		if _, found := seen[current]; found {
			return
		}
		seen[current] = struct{}{}
		if current.PkgPath() == schemaPackagePath && token.IsExported(current.Name()) {
			violations = append(violations, current)
			return
		}
		switch current.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Array:
			visit(current.Elem())
		case reflect.Map:
			visit(current.Key())
			visit(current.Elem())
		case reflect.Struct:
			for field := range current.Fields() {
				visit(field.Type)
			}
		default:
			return
		}
	}

	visit(root)
	return violations
}

func TestNewIDAExport(t *testing.T) {
	t.Parallel()

	analysis := Analysis{
		Input: Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		Runtime: &RuntimeMetadata{
			TrustSummary:                  RuntimeTrustSummarySymbolBacked,
			ELFFunctionFoothold:           "address_only",
			ELFFunctionFootholdCountHint:  2125,
			ELFFunctionFootholdTextSource: "moduledata_text",
			ELFFunctionFootholdStartAddr:  0x401000,
			ELFFunctionFootholdEndAddr:    0x4b4101,
			Provenance: Provenance{
				Source:     "core.runtime.elf",
				Confidence: "medium",
			},
		},
		BuildInfo: &BuildInfo{
			GoVersion: "go1.26.1",
			Path:      "example.com/sample",
			Provenance: Provenance{
				Source:     "core.buildinfo",
				Confidence: "high",
			},
		},
		Functions: []Function{
			{
				Name:          "main.main",
				Package:       "main",
				ImportPath:    "example.com/sample",
				SourceFile:    "main.go",
				SourceLine:    34,
				Autogenerated: false,
				Entry:         0x1000,
				End:           0x1100,
				ModuleLocal:   true,
				Provenance: Provenance{
					Source:     "core.pclntab",
					Confidence: "high",
				},
			},
		},
		Types: []Type{
			{
				Name:            "main.fixtureCounter",
				Package:         "main",
				ImportPath:      "example.com/sample",
				Kind:            "struct",
				SourceFileCount: 1,
				ModuleLocal:     true,
				UserMeaningful:  true,
				Provenance: Provenance{
					Source:     "core.types",
					Confidence: "medium",
				},
			},
		},
		Strings: []StringCandidate{
			{
				Value:  "hello",
				Region: ".rodata",
				Addr:   0x2010,
				Offset: 16,
				Provenance: Provenance{
					Source:     "core.strings",
					Confidence: "medium",
				},
			},
		},
		Peeling: &PeelingAnalysis{
			Functions: []PeelingFunction{{
				Name:                   "main.main",
				Package:                "main",
				ImportPath:             "example.com/sample",
				SourceFile:             "main.go",
				SourceLine:             34,
				Entry:                  0x1000,
				End:                    0x1100,
				ModuleLocal:            true,
				Classification:         PeelingClassUser,
				ClassificationEvidence: PeelingEvidenceModuleLocal,
			}},
			Packages: []PeelingPackage{{
				Name:                  "main",
				ImportPath:            "example.com/sample",
				ModuleLocal:           true,
				FunctionCount:         1,
				UserFunctionCount:     1,
				PrimaryClassification: PeelingClassUser,
			}},
			Provenance: Provenance{
				Source:     "engine.peeling",
				Confidence: "medium",
			},
		},
		Refined: &RefinedAnalysis{
			Functions: []RefinedFunction{{Name: "main.main"}},
			Types:     []RefinedType{{Name: "main.fixtureCounter"}},
			Strings:   []RefinedString{{Value: "hello"}},
		},
	}

	got := NewIDAExport(analysis)

	if got.Contract != IDAExportContractV1 {
		t.Fatalf("contract = %q, want %q", got.Contract, IDAExportContractV1)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != RuntimeTrustSummarySymbolBacked {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Runtime.ELFFunctionFoothold != "address_only" ||
		got.Runtime.ELFFunctionFootholdCountHint != 2125 ||
		got.Runtime.ELFFunctionFootholdTextSource != "moduledata_text" ||
		got.Runtime.ELFFunctionFootholdStartAddr != 0x401000 ||
		got.Runtime.ELFFunctionFootholdEndAddr != 0x4b4101 {
		t.Fatalf("runtime foothold = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) != 1 || got.Peeling.Functions[0].Classification != PeelingClassUser {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if got.Peeling.Functions[0].ClassificationEvidence != PeelingEvidenceModuleLocal {
		t.Fatalf("peeling evidence = %#v", got.Peeling.Functions[0])
	}
	if len(got.Peeling.Packages) != 1 || got.Peeling.Packages[0].PrimaryClassification != PeelingClassUser {
		t.Fatalf("peeling packages = %#v", got.Peeling.Packages)
	}
	if len(got.Functions) != 1 || got.Functions[0].RefinedName != "main.main" {
		t.Fatalf("functions = %#v", got.Functions)
	}
	if got.Functions[0].Provenance.Source != "core.pclntab" {
		t.Fatalf("function provenance = %#v", got.Functions[0].Provenance)
	}
	if got.Functions[0].Package != "main" || got.Functions[0].ImportPath != "example.com/sample" || got.Functions[0].SourceFile != "main.go" || got.Functions[0].SourceLine != 34 || !got.Functions[0].ModuleLocal {
		t.Fatalf("function metadata = %#v", got.Functions[0])
	}
	if len(got.Types) != 1 || got.Types[0].RefinedName != "main.fixtureCounter" {
		t.Fatalf("types = %#v", got.Types)
	}
	if got.Types[0].Package != "main" || got.Types[0].ImportPath != "example.com/sample" || got.Types[0].SourceFileCount != 1 || !got.Types[0].ModuleLocal || !got.Types[0].UserMeaningful {
		t.Fatalf("type metadata = %#v", got.Types[0])
	}
	if len(got.Strings) != 1 || got.Strings[0].RefinedValue != "hello" {
		t.Fatalf("strings = %#v", got.Strings)
	}
	if got.Strings[0].Address != 0x2010 {
		t.Fatalf("string address = %#v", got.Strings[0])
	}
}

func TestNewGhidraExport(t *testing.T) {
	t.Parallel()

	analysis := Analysis{
		Input: Input{
			Path:   "/tmp/sample.bin",
			Format: "elf",
		},
		Runtime: &RuntimeMetadata{
			TrustSummary:                  RuntimeTrustSummaryGoModuleFallback,
			ELFFunctionFoothold:           "address_only",
			ELFFunctionFootholdCountHint:  2083,
			ELFFunctionFootholdTextSource: "elf_text_section",
			ELFFunctionFootholdStartAddr:  0x11000,
			ELFFunctionFootholdEndAddr:    0xb55d1,
			Provenance: Provenance{
				Source:     "core.runtime.elf",
				Confidence: "medium",
			},
		},
		BuildInfo: &BuildInfo{
			GoVersion: "go1.26.1",
			Path:      "example.com/sample",
		},
		Functions: []Function{
			{
				Name:          "main.main",
				Package:       "main",
				ImportPath:    "example.com/sample",
				SourceFile:    "main.go",
				SourceLine:    34,
				Autogenerated: false,
				Entry:         0x1000,
				End:           0x1100,
				ModuleLocal:   true,
				Provenance: Provenance{
					Source:     "core.pclntab",
					Confidence: "high",
				},
			},
		},
		Peeling: &PeelingAnalysis{
			Functions: []PeelingFunction{{
				Name:                   "main.main",
				Package:                "main",
				ImportPath:             "example.com/sample",
				Entry:                  0x1000,
				End:                    0x1100,
				ModuleLocal:            true,
				Classification:         PeelingClassUser,
				ClassificationEvidence: PeelingEvidenceModuleLocal,
			}},
			Packages: []PeelingPackage{{
				Name:                  "main",
				ImportPath:            "example.com/sample",
				ModuleLocal:           true,
				FunctionCount:         1,
				UserFunctionCount:     1,
				PrimaryClassification: PeelingClassUser,
			}},
			Provenance: Provenance{
				Source:     "engine.peeling",
				Confidence: "medium",
			},
		},
		Refined: &RefinedAnalysis{
			Functions: []RefinedFunction{{Name: "main.main"}},
		},
	}

	got := NewGhidraExport(analysis)

	if got.Contract != GhidraExportContractV1 {
		t.Fatalf("contract = %q, want %q", got.Contract, GhidraExportContractV1)
	}
	if got.Runtime == nil || got.Runtime.TrustSummary != RuntimeTrustSummaryGoModuleFallback {
		t.Fatalf("runtime = %#v", got.Runtime)
	}
	if got.Runtime.ELFFunctionFoothold != "address_only" ||
		got.Runtime.ELFFunctionFootholdCountHint != 2083 ||
		got.Runtime.ELFFunctionFootholdTextSource != "elf_text_section" ||
		got.Runtime.ELFFunctionFootholdStartAddr != 0x11000 ||
		got.Runtime.ELFFunctionFootholdEndAddr != 0xb55d1 {
		t.Fatalf("runtime foothold = %#v", got.Runtime)
	}
	if got.Peeling == nil || len(got.Peeling.Functions) != 1 || got.Peeling.Functions[0].Classification != PeelingClassUser {
		t.Fatalf("peeling = %#v", got.Peeling)
	}
	if got.Peeling.Functions[0].ClassificationEvidence != PeelingEvidenceModuleLocal {
		t.Fatalf("peeling evidence = %#v", got.Peeling.Functions[0])
	}
	if len(got.Peeling.Packages) != 1 || got.Peeling.Packages[0].PrimaryClassification != PeelingClassUser {
		t.Fatalf("peeling packages = %#v", got.Peeling.Packages)
	}
	if got.Program.Path != "/tmp/sample.bin" {
		t.Fatalf("program path = %q", got.Program.Path)
	}
	if got.Program.ModulePath != "example.com/sample" {
		t.Fatalf("module path = %q", got.Program.ModulePath)
	}
	if len(got.Symbols) != 1 || got.Symbols[0].RefinedName != "main.main" {
		t.Fatalf("symbols = %#v", got.Symbols)
	}
	if got.Symbols[0].Address != 0x1000 {
		t.Fatalf("symbol address = %#v", got.Symbols[0])
	}
	if got.Symbols[0].Package != "main" || got.Symbols[0].ImportPath != "example.com/sample" || got.Symbols[0].SourceFile != "main.go" || got.Symbols[0].SourceLine != 34 || !got.Symbols[0].ModuleLocal {
		t.Fatalf("symbol metadata = %#v", got.Symbols[0])
	}
	analysis.Types = []Type{
		{
			Name:            "main.fixtureCounter",
			Package:         "main",
			ImportPath:      "example.com/sample",
			Kind:            "struct",
			SourceFileCount: 1,
			ModuleLocal:     true,
			UserMeaningful:  true,
			Provenance: Provenance{
				Source:     "core.types",
				Confidence: "medium",
			},
		},
	}
	got = NewGhidraExport(analysis)
	if len(got.Types) != 1 || got.Types[0].Package != "main" || got.Types[0].ImportPath != "example.com/sample" || got.Types[0].SourceFileCount != 1 || !got.Types[0].ModuleLocal || !got.Types[0].UserMeaningful {
		t.Fatalf("ghidra types = %#v", got.Types)
	}
	analysis.Strings = []StringCandidate{
		{
			Value:  "hello",
			Region: ".rodata",
			Addr:   0x2010,
			Offset: 16,
			Provenance: Provenance{
				Source:     "core.strings",
				Confidence: "medium",
			},
		},
	}
	got = NewGhidraExport(analysis)
	if len(got.Strings) != 1 || got.Strings[0].Address != 0x2010 {
		t.Fatalf("ghidra strings = %#v", got.Strings)
	}
}
