package schema

import (
	"bytes"
	"encoding/json"
	"errors"
)

// AnalysisStage identifies one recovery, derivation, or refinement stage.
type AnalysisStage string

const (
	AnalysisStageBuildInfo  AnalysisStage = "build_info"
	AnalysisStageRuntime    AnalysisStage = "runtime"
	AnalysisStageFunctions  AnalysisStage = "functions"
	AnalysisStagePackages   AnalysisStage = "packages"
	AnalysisStageTypes      AnalysisStage = "types"
	AnalysisStageStrings    AnalysisStage = "strings"
	AnalysisStageSourceTree AnalysisStage = "source_tree"
	AnalysisStagePeeling    AnalysisStage = "peeling"
	AnalysisStageRefinement AnalysisStage = "refinement"
)

// StageStatus records the outcome of one attempted analysis stage.
type StageStatus string

const (
	StageStatusAvailable   StageStatus = "available"
	StageStatusUnavailable StageStatus = "unavailable"
	StageStatusUnsupported StageStatus = "unsupported"
	StageStatusFailed      StageStatus = "failed"
)

// StageDiagnostic makes partial analysis outcomes explicit and machine-readable.
type StageDiagnostic struct {
	Stage   AnalysisStage `json:"stage"`
	Status  StageStatus   `json:"status"`
	Code    string        `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
}

// StageDiagnostics keeps the zero-value canonical analysis contract as an empty JSON array.
type StageDiagnostics []StageDiagnostic

// MarshalJSON emits nil diagnostics as [] while leaving HTML escaping to the outer encoder.
func (diagnostics StageDiagnostics) MarshalJSON() ([]byte, error) {
	if diagnostics == nil {
		return []byte("[]"), nil
	}

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode([]StageDiagnostic(diagnostics)); err != nil {
		return nil, err
	}

	encoded := output.Bytes()
	if len(encoded) == 0 || encoded[len(encoded)-1] != '\n' {
		return nil, errors.New("stage diagnostics encoder did not terminate JSON with a newline")
	}
	return bytes.Clone(encoded[:len(encoded)-1]), nil
}
