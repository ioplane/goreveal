package schema

import (
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

	analysisBytes, err := json.Marshal(Analysis{Diagnostics: []StageDiagnostic{}})
	if err != nil {
		t.Fatalf("json.Marshal(Analysis) error = %v", err)
	}
	if !strings.Contains(string(analysisBytes), `"diagnostics":[]`) {
		t.Fatalf("json.Marshal(Analysis) = %s, want diagnostics as []", analysisBytes)
	}
}
