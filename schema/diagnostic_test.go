package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStageDiagnosticJSON(t *testing.T) {
	t.Parallel()

	diagnostic := StageDiagnostic{
		Stage:   AnalysisStageFunctions,
		Status:  StageStatusFailed,
		Code:    "fixture_failure",
		Message: "fixture failure",
	}

	encoded, err := json.Marshal(diagnostic)
	if err != nil {
		t.Fatalf("json.Marshal(StageDiagnostic) error = %v", err)
	}
	want := `{"stage":"functions","status":"failed","code":"fixture_failure","message":"fixture failure"}`
	if string(encoded) != want {
		t.Fatalf("json.Marshal(StageDiagnostic) = %s, want %s", encoded, want)
	}

	analysisBytes, err := json.Marshal(Analysis{})
	if err != nil {
		t.Fatalf("json.Marshal(Analysis) error = %v", err)
	}
	if !strings.Contains(string(analysisBytes), `"diagnostics":[]`) {
		t.Fatalf("json.Marshal(Analysis{}) = %s, want diagnostics as []", analysisBytes)
	}
}

func TestAnalysisDiagnosticsRespectEncoderEscapeHTML(t *testing.T) {
	t.Parallel()

	analysis := Analysis{Diagnostics: []StageDiagnostic{{
		Stage:   AnalysisStageFunctions,
		Status:  StageStatusFailed,
		Message: "<failed>&",
	}}}

	defaultJSON, err := json.Marshal(analysis)
	if err != nil {
		t.Fatalf("json.Marshal(Analysis) error = %v", err)
	}
	if !strings.Contains(string(defaultJSON), `"message":"\u003cfailed\u003e\u0026"`) {
		t.Fatalf("json.Marshal(Analysis) = %s, want default HTML escaping", defaultJSON)
	}

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(analysis); err != nil {
		t.Fatalf("Encoder.Encode(Analysis) error = %v", err)
	}
	if !strings.Contains(output.String(), `"message":"<failed>&"`) {
		t.Fatalf("Encoder.Encode(Analysis) = %s, want disabled HTML escaping", output.Bytes())
	}
}
